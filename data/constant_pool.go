package data

import "fmt"

type CpTag int

const (
	cpUnknown CpTag = iota
	cpUtf8
	cpInteger
	cpClass
	cpNameAndType
	cpFieldref
	cpMethodref
)

type CpUtf8Info struct {
	Value string
	baseCpInfo
}

func (i CpUtf8Info) Utf8Info() CpUtf8Info {
	return i
}

func (i CpUtf8Info) Tag() CpTag {
	return cpUtf8
}

func (i CpUtf8Info) String() string {
	return fmt.Sprintf("\"%s\"", i.Value)
}

type CpIntegerInfo struct {
	Value int32
	baseCpInfo
}

func (i CpIntegerInfo) IntegerInfo() CpIntegerInfo {
	return i
}

func (i CpIntegerInfo) Tag() CpTag {
	return cpInteger
}

func (i CpIntegerInfo) String() string {
	return fmt.Sprint(i.Value)
}

type CpClassInfo struct {
	Name *CpInfo
	baseCpInfo
}

func (i CpClassInfo) ClassInfo() CpClassInfo {
	return i
}

func (i CpClassInfo) Tag() CpTag {
	return cpClass
}

func (i CpClassInfo) String() string {
	return fmt.Sprintf("<Class %s>", *i.Name)
}

type CpNameAndTypeInfo struct {
	Name       *CpInfo
	Descriptor *CpInfo
	baseCpInfo
}

func (i CpNameAndTypeInfo) NameAndTypeInfo() CpNameAndTypeInfo {
	return i
}

func (i CpNameAndTypeInfo) Tag() CpTag {
	return cpNameAndType
}

func (i CpNameAndTypeInfo) String() string {
	return fmt.Sprintf("<NameAndType: %s, %s>", *i.Name, *i.Descriptor)
}

type CpFieldrefInfo struct {
	Class       *CpInfo
	NameAndType *CpInfo
	baseCpInfo
}

func (i CpFieldrefInfo) FieldrefInfo() CpFieldrefInfo {
	return i
}

func (i CpFieldrefInfo) Tag() CpTag {
	return cpFieldref
}

func (i CpFieldrefInfo) String() string {
	return fmt.Sprintf("<Fieldref: %s, %s>", *i.Class, *i.NameAndType)
}

type CpMethodrefInfo struct {
	Class       *CpInfo
	NameAndType *CpInfo
	baseCpInfo
}

func (i CpMethodrefInfo) MethodrefInfo() CpMethodrefInfo {
	return i
}

func (i CpMethodrefInfo) Tag() CpTag {
	return cpMethodref
}

func (i CpMethodrefInfo) String() string {
	return fmt.Sprintf("<Methodref: %s, %s>", *i.Class, *i.NameAndType)
}

type CpInfo interface {
	Tag() CpTag
	Utf8Info() CpUtf8Info
	IntegerInfo() CpIntegerInfo
	ClassInfo() CpClassInfo
	NameAndTypeInfo() CpNameAndTypeInfo
	FieldrefInfo() CpFieldrefInfo
	MethodrefInfo() CpMethodrefInfo
	fmt.Stringer
}

type baseCpInfo struct{}

func (baseCpInfo) TagInfo() CpTag                     { return cpUnknown }
func (baseCpInfo) Utf8Info() CpUtf8Info               { panic("not a Utf8Info") }
func (baseCpInfo) IntegerInfo() CpIntegerInfo         { panic("not an IntegerInfo") }
func (baseCpInfo) ClassInfo() CpClassInfo             { panic("not a ClassInfo") }
func (baseCpInfo) NameAndTypeInfo() CpNameAndTypeInfo { panic("not a NameAndTypeInfo") }
func (baseCpInfo) FieldrefInfo() CpFieldrefInfo       { panic("not a FieldrefInfo") }
func (baseCpInfo) MethodrefInfo() CpMethodrefInfo     { panic("not a MethodrefInfo") }
func (baseCpInfo) StringInfo() string                 { return "<BaseCpInfo>" }
