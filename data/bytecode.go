package data

type Bytecode struct {
	Ops []Data
	baseData
}

func (b *Bytecode) Bytecode() *Bytecode {
	return b
}

func (Bytecode) Tag() Tag {
	return BYTECODE
}

func (b Bytecode) String() string {
	str := "Bytecode["
	for _, op := range b.Ops {
		str += "\n  " + op.String()
	}
	str += "]"
	return str
}
