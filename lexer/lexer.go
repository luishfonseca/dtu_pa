package lexer

import (
	"fmt"
	"io"
	"os"

	"github.com/luishfonseca/dtu_pa/state"
)

type Lexer struct {
	input   io.ReadSeekCloser
	tokenCh chan<- Token
	sc      stackedCounter
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
	}, nil
}

func (l *Lexer) Fail(err error) {
	l.err = err
}

func (l *Lexer) Run() error {
	state.Run(l, classStart)

	close(l.tokenCh)

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

func waitReq(l *Lexer) state.Fn[*Lexer] {
	return nil
}
