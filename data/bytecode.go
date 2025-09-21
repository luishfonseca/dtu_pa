package data

import "fmt"

type Bytecode struct {
	Ops []Data // TODO: create ops folder and add a new type erased interface Op
	baseData
}

func (b *Bytecode) Bytecode() *Bytecode {
	return b
}

func (*Bytecode) Tag() Tag {
	return BYTECODE
}

func (b Bytecode) String() string {
	str := "Bytecode["
	for _, op := range b.Ops {
		str += fmt.Sprint("\n ", op)
	}
	str += "]"
	return str
}
