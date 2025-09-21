package data

import "fmt"

type Constant interface {
	Tag() Tag
	ConstantUtf8() *ConstantUtf8
	ConstantInteger() *ConstantInteger
	ConstantClass() *ConstantClass
	ConstantNameAndType() *ConstantNameAndType
	ConstantFieldref() *ConstantFieldref
	ConstantMethodref() *ConstantMethodref
	fmt.Stringer
}

type baseConstant struct{ Data }

func (b *baseConstant) ConstantUtf8() *ConstantUtf8       { panic(msg(b, "ConstantUtf8")) }
func (b *baseConstant) ConstantInteger() *ConstantInteger { panic(msg(b, "ConstantInteger")) }
func (b *baseConstant) ConstantClass() *ConstantClass     { panic(msg(b, "ConstantClass")) }
func (b *baseConstant) ConstantNameAndType() *ConstantNameAndType {
	panic(msg(b, "ConstantNameAndType"))
}
func (b *baseConstant) ConstantFieldref() *ConstantFieldref   { panic(msg(b, "ConstantFieldref")) }
func (b *baseConstant) ConstantMethodref() *ConstantMethodref { panic(msg(b, "ConstantMethodref")) }

type ConstantUtf8 struct {
	Value string
	baseConstant
}

func (c *ConstantUtf8) ConstantUtf8() *ConstantUtf8 {
	return c
}

func (*ConstantUtf8) Tag() Tag {
	return CP_UTF8
}

func (c ConstantUtf8) String() string {
	return fmt.Sprintf("%q", c.Value)
}

type ConstantInteger struct {
	Value int32
	baseConstant
}

func (c *ConstantInteger) ConstantInteger() *ConstantInteger {
	return c
}

func (*ConstantInteger) Tag() Tag {
	return CP_INTEGER
}

func (c ConstantInteger) String() string {
	return fmt.Sprint(c.Value)
}

type ConstantClass struct {
	Name *Data
	baseConstant
}

func (c *ConstantClass) ConstantClass() *ConstantClass {
	return c
}

func (*ConstantClass) Tag() Tag {
	return CP_CLASS
}

func (c ConstantClass) String() string {
	return fmt.Sprintf("<Class %s>", *c.Name)
}

type ConstantNameAndType struct {
	Name       *Data
	Descriptor *Data
	baseConstant
}

func (c *ConstantNameAndType) NameAndTypeInfo() *ConstantNameAndType {
	return c
}

func (*ConstantNameAndType) Tag() Tag {
	return CP_NAME_AND_TYPE
}

func (c ConstantNameAndType) String() string {
	return fmt.Sprintf("<NameAndType: %s, %s>", *c.Name, *c.Descriptor)
}

type ConstantFieldref struct {
	Class       *Data
	NameAndType *Data
	baseConstant
}

func (c *ConstantFieldref) FieldrefInfo() *ConstantFieldref {
	return c
}

func (*ConstantFieldref) Tag() Tag {
	return CP_FIELDREF
}

func (c ConstantFieldref) String() string {
	return fmt.Sprintf("<Fieldref: %s, %s>", *c.Class, *c.NameAndType)
}

type ConstantMethodref struct {
	Class       *Data
	NameAndType *Data
	baseConstant
}

func (c *ConstantMethodref) MethodrefInfo() *ConstantMethodref {
	return c
}

func (*ConstantMethodref) Tag() Tag {
	return CP_METHODREF
}

func (c ConstantMethodref) String() string {
	return fmt.Sprintf("<Methodref: %s, %s>", *c.Class, *c.NameAndType)
}
