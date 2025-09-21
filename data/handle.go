package data

import "fmt"

type Handle interface {
	Tag() Tag
	AttributeHandle() *AttributeHandle
	BytecodeHandle() *BytecodeHandle
	fmt.Stringer
}

type baseHandle struct{ Data }

func (b *baseHandle) AttributeHandle() *AttributeHandle { panic(msg(b, "AttributeHandle")) }
func (b *baseHandle) BytecodeHandle() *BytecodeHandle   { panic(msg(b, "BytecodeHandle")) }

type AttributeHandle struct {
	AttributeTag Tag
	Begin        int64
	baseAttribute
}

func (*AttributeHandle) Tag() Tag {
	return ATTRIBUTE_HANDLE
}

func (a *AttributeHandle) AttributeHandle() *AttributeHandle {
	return a
}

func (a AttributeHandle) String() string {
	return fmt.Sprintf("<%s @ %d>", a.AttributeTag, a.Begin)
}

type BytecodeHandle struct {
	Begin  int64
	Length uint32
	baseAttribute
}

func (*BytecodeHandle) Tag() Tag {
	return BYTECODE_HANDLE
}

func (c *BytecodeHandle) BytecodeHandle() *BytecodeHandle {
	return c
}

func (c BytecodeHandle) String() string {
	return fmt.Sprintf("@[%d.. %d]", c.Begin, c.Begin+int64(c.Length))
}
