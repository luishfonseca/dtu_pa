package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/luishfonseca/dtu_pa/lexer"
)

type stateFn func(*Parser) stateFn

type Parser struct {
	tokenCh <-chan lexer.Token
	data    []map[string]any
	err     error
}

type ConfigProvider interface {
}

func New(cfg ConfigProvider, tokenCh <-chan lexer.Token) *Parser {
	return &Parser{
		tokenCh: tokenCh,
	}
}

func (p *Parser) PrintData() {
	for i, m := range p.data {
		if i != 0 {
			fmt.Println("----")
		}

		for k, v := range m {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}

func (p *Parser) Run() error {
	for state := magic; state != nil; {
		state = state(p)
	}

	if p.err != nil {
		return p.err
	}

	return nil
}

func (p *Parser) expect(t lexer.Type) ([]byte, error) {
	token, ok := <-p.tokenCh
	if !ok {
		return nil, fmt.Errorf("unexpected end of input, expected token type %d", t)
	}

	if token.Type != t {
		return nil, fmt.Errorf("unexpected token type %d, expected %d", token.Type, t)
	}

	return token.Bytes, nil
}

func magic(p *Parser) stateFn {
	b, err := p.expect(lexer.MAGIC)
	if err != nil {
		p.err = err
		return nil
	}

	if string(b) != "\xCA\xFE\xBA\xBE" {
		p.err = fmt.Errorf("invalid magic number: %x", b)
		return nil
	}

	return version
}

func version(p *Parser) stateFn {
	mb, err := p.expect(lexer.MINOR_VERSION)
	if err != nil {
		p.err = err
		return nil
	}

	Mb, err := p.expect(lexer.MAJOR_VERSION)
	if err != nil {
		p.err = err
		return nil
	}

	// binary.Decode() big-endian
	var m, M uint16
	if nrd, err := binary.Decode(mb, binary.BigEndian, &m); err != nil {
		p.err = err
		return nil
	} else if nrd != 2 {
		p.err = fmt.Errorf("binary decode read %d bytes, expected 2", nrd)
		return nil
	}

	if nrd, err := binary.Decode(Mb, binary.BigEndian, &M); err != nil {
		p.err = err
		return nil
	} else if nrd != 2 {
		p.err = fmt.Errorf("binary decode read %d bytes, expected 2", nrd)
		return nil
	}

	p.data = append(p.data, make(map[string]any))
	p.data[0]["version"] = fmt.Sprintf("%d.%d", M, m)

	return constant_pool
}

func constant_pool(p *Parser) stateFn {
	return nil
}
