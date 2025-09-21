package data

import "fmt"

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
	baseData
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
	baseData
}

func (a *AttributeSourceFile) AttributeSourceFile() *AttributeSourceFile {
	return a
}

func (*AttributeSourceFile) Tag() Tag {
	return ATTR_SOURCE_FILE
}

type AttributeRuntimeVisibleAnnotations struct {
	baseData
}

func (a *AttributeRuntimeVisibleAnnotations) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	return a
}

func (*AttributeRuntimeVisibleAnnotations) Tag() Tag {
	return ATTR_RUNTIME_VISIBLE_ANNOTATIONS
}

type AttributeInnerClasses struct {
	baseData
}

func (a *AttributeInnerClasses) AttributeInnerClasses() *AttributeInnerClasses {
	return a
}

func (*AttributeInnerClasses) Tag() Tag {
	return ATTR_INNER_CLASSES
}

type AttributeLineNumberTable struct {
	baseData
}

func (a *AttributeLineNumberTable) AttributeLineNumberTable() *AttributeLineNumberTable {
	return a
}

func (*AttributeLineNumberTable) Tag() Tag {
	return ATTR_LINE_NUMBER_TABLE
}

type AttributeLocalVariableTable struct {
	baseData
}

func (a *AttributeLocalVariableTable) AttributeLocalVariableTable() *AttributeLocalVariableTable {
	return a
}

func (*AttributeLocalVariableTable) Tag() Tag {
	return ATTR_LOCAL_VARIABLE_TABLE
}

type AttributeStackMapTable struct {
	baseData
}

func (a *AttributeStackMapTable) AttributeStackMapTable() *AttributeStackMapTable {
	return a
}

func (*AttributeStackMapTable) Tag() Tag {
	return ATTR_STACK_MAP_TABLE
}
