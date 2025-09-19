package lexer

import (
	"io"
	"os"
)

type Type int

const (
	EOF Type = iota
	MAGIC
	MINOR_VERSION
	MAJOR_VERSION
	UNKNOWN
)

type Token struct {
	Type  Type
	Bytes []byte
}

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input   io.Reader
	tokenCh chan<- Token
	curr    []byte
	err     error
}

type ConfigProvider interface {
	GetClassFile() string
}

func New(cfg ConfigProvider, tokenCh chan<- Token) (*Lexer, error) {
	input, err := os.Open(cfg.GetClassFile())
	if err != nil {
		return nil, err
	}

	return &Lexer{
		input:   input,
		tokenCh: tokenCh,
		err:     nil,
		curr:    nil,
	}, nil
}

func (l *Lexer) Run() error {
	for state := magic; state != nil; {
		state = state(l)
	}

	close(l.tokenCh)

	if l.err != nil {
		return l.err
	}

	return nil
}

func (l *Lexer) read(n int) int {
	l.curr = make([]byte, n)

	if nrd, err := io.ReadFull(l.input, l.curr); err != nil {
		l.err = err
		return nrd
	}

	return n
}

func (l *Lexer) emit(t Type) {
	if l.curr == nil {
		panic("cannot emit a nil token")
	}

	l.tokenCh <- Token{
		Type:  t,
		Bytes: l.curr,
	}

	l.curr = nil
}

// The magic item supplies the magic number identifying the class file format; it has the value 0xCAFEBABE.
func magic(l *Lexer) stateFn {
	if l.read(4) != 4 {
		return nil
	}

	l.emit(MAGIC)

	return version
}

// The values of the minor_version and major_version items are the minor and major version numbers of this class file.
func version(l *Lexer) stateFn {
	if l.read(2) != 2 {
		return nil
	}

	l.emit(MINOR_VERSION)

	if l.read(2) != 2 {
		return nil
	}

	l.emit(MAJOR_VERSION)

	return nil
}
