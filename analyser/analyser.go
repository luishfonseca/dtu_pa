package analyser

import "fmt"

type Config struct {
	ClassFile string
}

type analyser struct {
	config Config
}

func New(classFile string) *analyser {
	return &analyser{
		config: Config{
			ClassFile: classFile,
		},
	}
}

func (a *analyser) Inspect() error {
	fmt.Printf("Inspecting %s\n", a.config.ClassFile)
	return nil
}
