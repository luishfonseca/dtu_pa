package data

import "fmt"

type Tag int

const (
	UNKNOWN Tag = iota
	DECOMPILED_CLASS
	BYTECODE
	ATTRIBUTE_HANDLE
	BYTECODE_HANDLE
	CP_UTF8
	CP_INTEGER
	CP_CLASS
	CP_NAME_AND_TYPE
	CP_FIELDREF
	CP_METHODREF
	ATTR_CODE
	ATTR_SOURCE_FILE
	ATTR_RUNTIME_VISIBLE_ANNOTATIONS
	ATTR_INNER_CLASSES
	ATTR_LINE_NUMBER_TABLE
	ATTR_LOCAL_VARIABLE_TABLE
	ATTR_STACK_MAP_TABLE
)

func (t Tag) String() string {
	switch t {
	case UNKNOWN:
		return "Unknown"
	case DECOMPILED_CLASS:
		return "DecompiledClass"
	case BYTECODE:
		return "Bytecode"
	case ATTRIBUTE_HANDLE:
		return "AttributeHandle"
	case BYTECODE_HANDLE:
		return "BytecodeHandle"
	case CP_UTF8:
		return "ConstantUtf8"
	case CP_INTEGER:
		return "ConstantInteger"
	case CP_CLASS:
		return "ConstantClass"
	case CP_NAME_AND_TYPE:
		return "ConstantNameAndType"
	case CP_FIELDREF:
		return "ConstantFieldref"
	case CP_METHODREF:
		return "ConstantMethodref"
	case ATTR_CODE:
		return "AttributeCode"
	case ATTR_SOURCE_FILE:
		return "AttributeSourceFile"
	case ATTR_RUNTIME_VISIBLE_ANNOTATIONS:
		return "AttributeRuntimeVisibleAnnotations"
	case ATTR_INNER_CLASSES:
		return "AttributeInnerClasses"
	case ATTR_LINE_NUMBER_TABLE:
		return "AttributeLineNumberTable"
	case ATTR_LOCAL_VARIABLE_TABLE:
		return "AttributeLocalVariableTable"
	case ATTR_STACK_MAP_TABLE:
		return "AttributeStackMapTable"
	default:
		return fmt.Sprintf("Tag(%d)", int(t))
	}
}

type Data interface {
	Tag() Tag
	DecompiledClass() *DecompiledClass
	Bytecode() *Bytecode
	Constant
	Attribute
	Handle
	fmt.Stringer
}

type baseData struct {
	baseAttribute
	baseConstant
	baseHandle
}

func msg(b Data, expected string) string {
	return fmt.Sprintf("expected %s got %s", expected, b.Tag())
}

func (baseData) Tag() Tag                             { return UNKNOWN }
func (b *baseData) DecompiledClass() *DecompiledClass { panic(msg(b, "DecompiledClass")) }
func (b *baseData) Bytecode() *Bytecode               { panic(msg(b, "Bytecode")) }
func (b baseData) String() string                     { return b.Tag().String() }
