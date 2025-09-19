package analyser

import (
	"fmt"
	"os"

	"github.com/luishfonseca/dtu_pa/lexer"
)

type analyser struct {
	classFile string
}

func (a *analyser) GetClassFile() string {
	return a.classFile
}

func New(classFile string) *analyser {
	return &analyser{
		classFile: classFile,
	}
}

func (a *analyser) Inspect() error {
	fmt.Printf("Inspecting %s\n", a.classFile)

	tokenCh := make(chan lexer.Token)

	l, err := lexer.New(a, tokenCh)
	if err != nil {
		return fmt.Errorf("could not create lexer: %w", err)
	}

	go func() {
		if err := l.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "lexer error: %v\n", err)
		}
	}()

	for token := range tokenCh {
		fmt.Printf("Token: %+v\n", token)
	}

	return nil
}
