package analyser

import (
	"fmt"
	"os"

	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/parser"
)

type analyser struct {
	classFile string
}

func New(classFile string) *analyser {
	return &analyser{
		classFile: classFile,
	}
}

func (a *analyser) GetClassFile() string {
	return a.classFile
}

func (a *analyser) Inspect() error {
	tokenCh := make(chan lexer.Token)

	l, err := lexer.New(a, tokenCh)
	if err != nil {
		return fmt.Errorf("could not create lexer: %w", err)
	}

	p := parser.New(a, tokenCh)

	go func() {
		if err := l.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: lexer: %v\n", err)
		}
	}()

	if err := p.Run(); err != nil {
		return fmt.Errorf("parser: %w", err)
	}

	lexerDone := true
	for t := range tokenCh {
		lexerDone = false
		fmt.Printf("Token: %v, Bytes: % X\n", t.Type, t.Bytes)
	}

	if !lexerDone {
		return fmt.Errorf("parser finished before receiving all tokens from lexer")
	}

	p.PrintData()

	return nil
}
