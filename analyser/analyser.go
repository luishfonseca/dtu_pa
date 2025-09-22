package analyser

import (
	"fmt"
	"os"

	"github.com/luishfonseca/dtu_pa/data"
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

func (a *analyser) Inspect() error {
	dataCh := make(chan data.Data)

	reqCh := make(chan data.Data)
	defer close(reqCh)

	p, err := parser.New(a.classFile, dataCh, reqCh)
	if err != nil {
		return fmt.Errorf("error creating parser: %w", err)
	}

	go func() {
		if err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: parser: %v\n", err)
		}
	}()

	d, ok := <-dataCh
	if !ok {
		return fmt.Errorf("error: no data received from parser")
	}

	class := d.Class()
	fmt.Println(class)

	for _, method := range class.Methods {
		reqCh <- method.Attributes[data.ATTR_CODE]

		d, ok = <-dataCh
		if !ok {
			return fmt.Errorf("error: no data received from parser")
		}

		attr := d.AttributeCode()
		fmt.Println(method.Name, method.Descriptor, "->", attr)

		reqCh <- &attr.CodeHandle

		d, ok = <-dataCh
		if !ok {
			return fmt.Errorf("error: no data received from parser")
		}

		fmt.Println(d.Bytecode())
		fmt.Println()
	}

	return nil
}
