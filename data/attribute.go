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

func (a *AttributeCode) Tag() Tag                      { return ATTR_CODE }
func (a *AttributeCode) AttributeCode() *AttributeCode { return a }
func (d *baseData) AttributeCode() *AttributeCode      { panic(msg(d, "AttributeCode")) }

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

type AttributeSourceFile struct{ baseData }

func (a *AttributeSourceFile) Tag() Tag                                  { return ATTR_SOURCE_FILE }
func (a *AttributeSourceFile) AttributeSourceFile() *AttributeSourceFile { return a }
func (d *baseData) AttributeSourceFile() *AttributeSourceFile            { panic(msg(d, "AttributeSourceFile")) }

type AttributeRuntimeVisibleAnnotations struct{ baseData }

func (a *AttributeRuntimeVisibleAnnotations) Tag() Tag { return ATTR_RUNTIME_VISIBLE_ANNOTATIONS }
func (a *AttributeRuntimeVisibleAnnotations) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	return a
}
func (d *baseData) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	panic(msg(d, "AttributeRuntimeVisibleAnnotations"))
}

type AttributeInnerClasses struct{ baseData }

func (a *AttributeInnerClasses) Tag() Tag                                      { return ATTR_INNER_CLASSES }
func (a *AttributeInnerClasses) AttributeInnerClasses() *AttributeInnerClasses { return a }
func (d *baseData) AttributeInnerClasses() *AttributeInnerClasses {
	panic(msg(d, "AttributeInnerClasses"))
}

type AttributeLineNumberTable struct{ baseData }

func (a *AttributeLineNumberTable) Tag() Tag { return ATTR_LINE_NUMBER_TABLE }
func (a *AttributeLineNumberTable) AttributeLineNumberTable() *AttributeLineNumberTable {
	return a
}
func (d *baseData) AttributeLineNumberTable() *AttributeLineNumberTable {
	panic(msg(d, "AttributeLineNumberTable"))
}

type AttributeLocalVariableTable struct{ baseData }

func (a *AttributeLocalVariableTable) Tag() Tag { return ATTR_LOCAL_VARIABLE_TABLE }
func (a *AttributeLocalVariableTable) AttributeLocalVariableTable() *AttributeLocalVariableTable {
	return a
}
func (d *baseData) AttributeLocalVariableTable() *AttributeLocalVariableTable {
	panic(msg(d, "AttributeLocalVariableTable"))
}

type AttributeStackMapTable struct{ baseData }

func (a *AttributeStackMapTable) Tag() Tag { return ATTR_STACK_MAP_TABLE }
func (a *AttributeStackMapTable) AttributeStackMapTable() *AttributeStackMapTable {
	return a
}
func (d *baseData) AttributeStackMapTable() *AttributeStackMapTable {
	panic(msg(d, "AttributeStackMapTable"))
}
