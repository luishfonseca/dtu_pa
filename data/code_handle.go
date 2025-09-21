package data

import (
	"fmt"
)

type BytecodeHandle struct {
	Begin  int64
	Length uint32
	baseData
}

func (*BytecodeHandle) Tag() Tag {
	return BYTECODE_HANDLE
}

func (c *BytecodeHandle) BytecodeHandle() *BytecodeHandle {
	return c
}

func (c BytecodeHandle) String() string {
	return fmt.Sprintf("<@ [%d.. %d] >", c.Begin, c.Begin+int64(c.Length))
}
