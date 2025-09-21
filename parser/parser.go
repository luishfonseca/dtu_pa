package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/state"
)

type Parser struct {
	input      io.ReadSeekCloser
	dataCh     chan<- data.Data
	reqCh      <-chan data.Data
	attributes map[data.AttributeHandle]data.Data
	codes      map[data.BytecodeHandle]*data.Bytecode
	class      *data.DecompiledClass
	err        error
}

func New(file string, dataCh chan<- data.Data, reqCh <-chan data.Data) (*Parser, error) {
	input, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &Parser{
		input:      input,
		dataCh:     dataCh,
		reqCh:      reqCh,
		attributes: make(map[data.AttributeHandle]data.Data),
		codes:      make(map[data.BytecodeHandle]*data.Bytecode),
	}, nil
}

func (p *Parser) read(n int) ([]byte, error) {
	token := make([]byte, n)
	if _, err := io.ReadFull(p.input, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (p *Parser) readDecode(v any) error {
	n := binary.Size(v)

	b, err := p.read(n)
	if err != nil {
		return err
	}

	if nrd, err := binary.Decode(b, binary.BigEndian, v); err != nil {
		return err
	} else if nrd != n {
		return fmt.Errorf("binary decode read %d bytes, expected %d", nrd, n)
	}

	return nil
}

func (p *Parser) Fail(err error) {
	p.err = err
}

func (p *Parser) Run() error {
	defer close(p.dataCh)

	p.class = &data.DecompiledClass{}
	state.Run(p, classStart)

	if p.err != nil {
		return p.err
	}

	return nil
}

func classStart(p *Parser) state.Fn[*Parser] {
	return magic
}

func classEnd(p *Parser) state.Fn[*Parser] {
	if _, err := p.read(1); err != io.EOF {
		return state.Fail[*Parser](fmt.Errorf("expected EOF, got more data"))
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
		if attr, ok := p.attributes[*req.AttributeHandle()]; ok {
			p.dataCh <- attr
		} else {
			return attribute(*req.AttributeHandle())
		}
	case data.BYTECODE_HANDLE:
		if bc, ok := p.codes[*req.BytecodeHandle()]; ok {
			p.dataCh <- bc
		} else {
			return state.Fail[*Parser](fmt.Errorf("bytecode handle unimplemented"))
		}
	default:
		return state.Fail[*Parser](fmt.Errorf("unexpected request tag: %s", req.Tag()))
	}

	return waitReq
}

func done(p *Parser) state.Fn[*Parser] {
	return nil
}
