package data

import "fmt"

type Tag int

const (
	UNKNOWN Tag = iota
	DECOMPILED_CLASS
	ATTR_HANDLE
	CODE_HANDLE
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
)

func (t Tag) String() string {
	switch t {
	case UNKNOWN:
		return "Unknown"
	case DECOMPILED_CLASS:
		return "DecompiledClass"
	case ATTR_HANDLE:
		return "AttributeHandle"
	case CODE_HANDLE:
		return "CodeHandle"
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
	default:
		return fmt.Sprintf("Tag(%d)", int(t))
	}
}

type Data interface {
	Tag() Tag
	DecompiledClass() *DecompiledClass
	AttributeHandle() *AttributeHandle
	CodeHandle() *CodeHandle
	ConstantUtf8() *ConstantUtf8
	ConstantInteger() *ConstantInteger
	ConstantClass() *ConstantClass
	ConstantNameAndType() *ConstantNameAndType
	ConstantFieldref() *ConstantFieldref
	ConstantMethodref() *ConstantMethodref
	AttributeCode() *AttributeCode
	AttributeSourceFile() *AttributeSourceFile
	AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations
	AttributeInnerClasses() *AttributeInnerClasses
	fmt.Stringer
}

type baseData struct{}

func msg(b *baseData, expected string) string {
	return fmt.Sprintf("expected %s got %s", expected, b.Tag())
}

func (baseData) Tag() Tag                                     { return UNKNOWN }
func (b *baseData) DecompiledClass() *DecompiledClass         { panic(msg(b, "DecompiledClass")) }
func (b *baseData) AttributeHandle() *AttributeHandle         { panic(msg(b, "AttributeHandle")) }
func (b *baseData) CodeHandle() *CodeHandle                   { panic(msg(b, "CodeHandle")) }
func (b *baseData) ConstantUtf8() *ConstantUtf8               { panic(msg(b, "ConstantUtf8")) }
func (b *baseData) ConstantInteger() *ConstantInteger         { panic(msg(b, "ConstantInteger")) }
func (b *baseData) ConstantClass() *ConstantClass             { panic(msg(b, "ConstantClass")) }
func (b *baseData) ConstantNameAndType() *ConstantNameAndType { panic(msg(b, "ConstantNameAndType")) }
func (b *baseData) ConstantFieldref() *ConstantFieldref       { panic(msg(b, "ConstantFieldref")) }
func (b *baseData) ConstantMethodref() *ConstantMethodref     { panic(msg(b, "ConstantMethodref")) }
func (b *baseData) AttributeCode() *AttributeCode             { panic(msg(b, "AttributeCode")) }
func (b *baseData) AttributeSourceFile() *AttributeSourceFile { panic(msg(b, "AttributeSourceFile")) }
func (b *baseData) AttributeRuntimeVisibleAnnotations() *AttributeRuntimeVisibleAnnotations {
	panic(msg(b, "AttributeRuntimeVisibleAnnotations"))
}
func (b *baseData) AttributeInnerClasses() *AttributeInnerClasses {
	panic(msg(b, "AttributeInnerClasses"))
}
func (baseData) String() string { return "<unknown>" }
