package data

import "fmt"

type AttributeHandle struct {
	AttributeTag Tag
	Begin        int64
	baseData
}

func (*AttributeHandle) Tag() Tag                            { return ATTRIBUTE_HANDLE }
func (a *AttributeHandle) AttributeHandle() *AttributeHandle { return a }
func (d *baseData) AttributeHandle() *AttributeHandle        { panic(msg(d, "AttributeHandle")) }

func (a AttributeHandle) String() string {
	return fmt.Sprintf("<%s @ %d>", a.AttributeTag, a.Begin)
}

type BytecodeHandle struct {
	Begin  int64
	Length uint32
	baseData
}

func (*BytecodeHandle) Tag() Tag                          { return BYTECODE_HANDLE }
func (b *BytecodeHandle) BytecodeHandle() *BytecodeHandle { return b }
func (d *baseData) BytecodeHandle() *BytecodeHandle       { panic(msg(d, "BytecodeHandle")) }

func (c BytecodeHandle) String() string {
	return fmt.Sprintf("@[%d.. %d]", c.Begin, c.Begin+int64(c.Length))
}
