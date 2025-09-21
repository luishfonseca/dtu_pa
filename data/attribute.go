package data

import "fmt"

type Attribute interface {
	Tag() Tag
	AttributeCode() *AttributeCode
	AttributeSourceFile() *AttributeSourceFile
	AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations
	AttributeInnerClasses() *AttributeInnerClasses
	AttributeLineNumberTable() *AttributeLineNumberTable
	AttributeLocalVariableTable() *AttributeLocalVariableTable
	AttributeStackMapTable() *AttributeStackMapTable
	fmt.Stringer
}

type baseAttribute struct{ Data }

func (b *baseAttribute) AttributeCode() *AttributeCode { panic(msg(b, "AttributeCode")) }
func (b *baseAttribute) AttributeSourceFile() *AttributeSourceFile {
	panic(msg(b, "AttributeSourceFile"))
}
func (b *baseAttribute) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	panic(msg(b, "AttributeRuntimeVisibleAnnotations"))
}
func (b *baseAttribute) AttributeInnerClasses() *AttributeInnerClasses {
	panic(msg(b, "AttributeInnerClasses"))
}
func (b *baseAttribute) AttributeLineNumberTable() *AttributeLineNumberTable {
	panic(msg(b, "AttributeLineNumberTable"))
}
func (b *baseAttribute) AttributeLocalVariableTable() *AttributeLocalVariableTable {
	panic(msg(b, "AttributeLocalVariableTable"))
}
func (b *baseAttribute) AttributeStackMapTable() *AttributeStackMapTable {
	panic(msg(b, "AttributeStackMapTable"))
}

type ExceptionTableEntry struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	CatchType *ConstantClass
}

func (e ExceptionTableEntry) String() string {
	c := "Any"
	if e.CatchType != nil {
		c = e.CatchType.String()
	}

	return fmt.Sprintf("<[%d, %d] -> %d \\ %s>", e.StartPC, e.EndPC, e.HandlerPC, c)
}

type AttributeCode struct {
	MaxStack       uint16
	MaxLocals      uint16
	CodeHandle     BytecodeHandle
	ExceptionTable []ExceptionTableEntry
	Attributes     []AttributeHandle
	baseAttribute
}

func (a *AttributeCode) AttributeCode() *AttributeCode {
	return a
}

func (*AttributeCode) Tag() Tag {
	return ATTR_CODE
}

func (a AttributeCode) String() string {
	str := "AttributeCode {"
	str += fmt.Sprint("\n  MaxStack: ", a.MaxStack)
	str += fmt.Sprint("\n  MaxLocals: ", a.MaxLocals)
	str += fmt.Sprint("\n  Code: ", a.CodeHandle)
	str += "\n  ExceptionTable: ["
	for _, entry := range a.ExceptionTable {
		str += fmt.Sprint("\n    ", entry)
	}
	str += "]\n  Attributes: ["
	for _, attr := range a.Attributes {
		str += fmt.Sprint("\n    ", attr)
	}
	str += "]"
	str += "\n}"
	return str
}

type AttributeSourceFile struct {
	baseAttribute
}

func (a *AttributeSourceFile) AttributeSourceFile() *AttributeSourceFile {
	return a
}

func (*AttributeSourceFile) Tag() Tag {
	return ATTR_SOURCE_FILE
}

type AttributeRuntimeVisibleAnnotations struct {
	baseAttribute
}

func (a *AttributeRuntimeVisibleAnnotations) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	return a
}

func (*AttributeRuntimeVisibleAnnotations) Tag() Tag {
	return ATTR_RUNTIME_VISIBLE_ANNOTATIONS
}

type AttributeInnerClasses struct {
	baseAttribute
}

func (a *AttributeInnerClasses) AttributeInnerClasses() *AttributeInnerClasses {
	return a
}

func (*AttributeInnerClasses) Tag() Tag {
	return ATTR_INNER_CLASSES
}

type AttributeLineNumberTable struct {
	baseAttribute
}

func (a *AttributeLineNumberTable) AttributeLineNumberTable() *AttributeLineNumberTable {
	return a
}

func (*AttributeLineNumberTable) Tag() Tag {
	return ATTR_LINE_NUMBER_TABLE
}

type AttributeLocalVariableTable struct {
	baseAttribute
}

func (a *AttributeLocalVariableTable) AttributeLocalVariableTable() *AttributeLocalVariableTable {
	return a
}

func (*AttributeLocalVariableTable) Tag() Tag {
	return ATTR_LOCAL_VARIABLE_TABLE
}

type AttributeStackMapTable struct {
	baseAttribute
}

func (a *AttributeStackMapTable) AttributeStackMapTable() *AttributeStackMapTable {
	return a
}

func (*AttributeStackMapTable) Tag() Tag {
	return ATTR_STACK_MAP_TABLE
}
