package lexer

import (
	"fmt"
	"io"
	"os"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/state"
)

type Lexer struct {
	input   io.ReadSeekCloser
	tokenCh chan<- Token
	reqCh   <-chan data.Data
	sc      stackedCounter
	curr    []byte
	err     error
}

type ConfigProvider interface {
	GetClassFile() string
}

func New(file string, tokenCh chan<- Token, reqCh <-chan data.Data) (*Lexer, error) {
	input, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		input:   input,
		tokenCh: tokenCh,
		reqCh:   reqCh,
	}, nil
}

func (l *Lexer) Fail(err error) {
	l.err = err
}

func (l *Lexer) Run() error {
	defer close(l.tokenCh)

	state.Run(l, classStart)

	if l.err != nil {
		return l.err
	}

	if !l.sc.empty() {
		return fmt.Errorf("stacked counter not empty at end of input")
	}

	if err := l.input.Close(); err != nil {
		return fmt.Errorf("error closing input: %w", err)
	}

	return nil
}

func (l *Lexer) read(n int) error {
	l.curr = make([]byte, n)

	if _, err := io.ReadFull(l.input, l.curr); err != nil {
		return err
	}

	return nil
}

func (l *Lexer) emit(t TokenType) {
	if l.curr == nil {
		panic("cannot emit a nil token")
	}

	l.tokenCh <- Token{
		Type:  t,
		Bytes: l.curr,
	}

	l.curr = nil
}

func classStart(l *Lexer) state.Fn[*Lexer] {
	return magic
}

func classEnd(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(1); err != io.EOF {
		if err != nil {
			return state.Fail[*Lexer](err)
		}

		return state.Fail[*Lexer](fmt.Errorf("expected EOF"))
	}

	l.emit(EOF)

	return waitReq
}

func waitReq(p *Lexer) state.Fn[*Lexer] {
	req, ok := <-p.reqCh
	if !ok {
		return done
	}

	switch req.Tag() {
	case data.ATTRIBUTE_HANDLE:
		switch req.AttributeHandle().AttributeTag {
		case data.ATTR_CODE:
			return attributeCode
		default:
			return state.Fail[*Lexer](fmt.Errorf("unexpected attribute tag: %s", req.AttributeHandle().AttributeTag))
		}
	case data.BYTECODE_HANDLE:
		return state.Fail[*Lexer](fmt.Errorf("bytecode handle unimplemented"))
	default:
		return state.Fail[*Lexer](fmt.Errorf("unexpected request tag: %s", req.Tag()))
	}
}

func attributeCode(p *Lexer) state.Fn[*Lexer] {
	return waitReq
}

func done(p *Lexer) state.Fn[*Lexer] {
	return nil
}
