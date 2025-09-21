package parser

import (
	"fmt"
	"os"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

type Parser struct {
	lexer      *lexer.Lexer
	tokenCh    <-chan lexer.Token
	tokenReqCh chan<- data.Data
	dataCh     chan<- data.Data
	reqCh      <-chan data.Data
	attributes map[data.AttributeHandle]data.Data
	codes      map[data.BytecodeHandle]*data.Bytecode
	class      *data.DecompiledClass
	err        error
}

func New(file string, dataCh chan<- data.Data, reqCh <-chan data.Data) (*Parser, error) {
	tokenCh := make(chan lexer.Token)
	tokenReqCh := make(chan data.Data)

	lexer, err := lexer.New(file, tokenCh, tokenReqCh)
	if err != nil {
		return nil, err
	}

	return &Parser{
		lexer:      lexer,
		tokenCh:    tokenCh,
		tokenReqCh: tokenReqCh,
		dataCh:     dataCh,
		reqCh:      reqCh,
	}, nil
}

func (p *Parser) Fail(err error) {
	p.err = err
}

func (p *Parser) Run() error {
	defer close(p.dataCh)
	defer close(p.tokenReqCh)

	go func() {
		if err := p.lexer.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: lexer: %v\n", err)
		}
	}()

	p.class = &data.DecompiledClass{}
	state.Run(p, classStart)

	if p.err != nil {
		return p.err
	}

	return nil
}

func (p *Parser) expect(t lexer.TokenType) ([]byte, error) {
	token, ok := <-p.tokenCh
	if !ok {
		return nil, fmt.Errorf("unexpected end of input, expected token type %s", t)
	}

	if token.Type != t {
		return nil, fmt.Errorf("unexpected token type %s, expected %s", token.Type, t)
	}

	return token.Bytes, nil
}

func (p *Parser) expectDecode(t lexer.TokenType, v any) error {
	b, err := p.expect(t)
	if err != nil {
		return err
	}

	if err := util.Decode(b, v); err != nil {
		return err
	}

	return nil
}

func classStart(p *Parser) state.Fn[*Parser] {
	return magic
}

func classEnd(p *Parser) state.Fn[*Parser] {
	if _, err := p.expect(lexer.EOF); err != nil {
		return state.Fail[*Parser](err)
	}

	p.dataCh <- p.class

	return waitReq
}

func waitReq(p *Parser) state.Fn[*Parser] {
	req, ok := <-p.reqCh
	if !ok {
		return done
	}

	switch req.Tag() {
	case data.ATTRIBUTE_HANDLE:
		if data, ok := p.attributes[*req.AttributeHandle()]; ok {
			p.dataCh <- data
		} else {
			p.tokenReqCh <- req
			return attribute
		}
	case data.BYTECODE_HANDLE:
		if data, ok := p.codes[*req.BytecodeHandle()]; ok {
			p.dataCh <- data
		} else {
			p.tokenReqCh <- req
			return state.Fail[*Parser](fmt.Errorf("bytecode handle unimplemented"))
		}
	default:
		return state.Fail[*Parser](fmt.Errorf("unexpected request tag: %s", req.Tag()))
	}

	return waitReq
}

func attribute(p *Parser) state.Fn[*Parser] {
	return nil
}

func done(p *Parser) state.Fn[*Parser] {
	lexerDone := true
	for token, ok := <-p.tokenCh; ok; token, ok = <-p.tokenCh {
		lexerDone = false
		fmt.Fprintf(os.Stderr, "warning: unprocessed token: %s\n", token.Type)
	}

	if !lexerDone {
		return state.Fail[*Parser](fmt.Errorf("stopped before lexer was done"))
	}

	return nil
}
