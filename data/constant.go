package data

import "fmt"

type ConstantUtf8 struct {
	Value string
	baseData
}

func (c *ConstantUtf8) Tag() Tag                    { return CP_UTF8 }
func (c *ConstantUtf8) ConstantUtf8() *ConstantUtf8 { return c }
func (d *baseData) ConstantUtf8() *ConstantUtf8     { panic(msg(d, "ConstantUtf8")) }

func (c ConstantUtf8) String() string {
	return fmt.Sprintf("%q", c.Value)
}

type ConstantInteger struct {
	Value int32
	baseData
}

func (c *ConstantInteger) Tag() Tag                          { return CP_INTEGER }
func (c *ConstantInteger) ConstantInteger() *ConstantInteger { return c }
func (d *baseData) ConstantInteger() *ConstantInteger        { panic(msg(d, "ConstantInteger")) }

func (c ConstantInteger) String() string {
	return fmt.Sprint(c.Value)
}

type ConstantClass struct {
	Name *Data
	baseData
}

func (c *ConstantClass) Tag() Tag                      { return CP_CLASS }
func (c *ConstantClass) ConstantClass() *ConstantClass { return c }
func (d *baseData) ConstantClass() *ConstantClass      { panic(msg(d, "ConstantClass")) }

func (c ConstantClass) String() string {
	return fmt.Sprintf("<Class %s>", *c.Name)
}

type ConstantNameAndType struct {
	Name       *Data
	Descriptor *Data
	baseData
}

func (c *ConstantNameAndType) Tag() Tag                                  { return CP_NAME_AND_TYPE }
func (c *ConstantNameAndType) ConstantNameAndType() *ConstantNameAndType { return c }
func (d *baseData) ConstantNameAndType() *ConstantNameAndType            { panic(msg(d, "ConstantNameAndType")) }

func (c ConstantNameAndType) String() string {
	return fmt.Sprintf("<NameAndType: %s, %s>", *c.Name, *c.Descriptor)
}

type ConstantFieldref struct {
	Clazz       *Data
	NameAndType *Data
	baseData
}

func (c *ConstantFieldref) Tag() Tag                            { return CP_FIELDREF }
func (c *ConstantFieldref) ConstantFieldref() *ConstantFieldref { return c }
func (d *baseData) ConstantFieldref() *ConstantFieldref         { panic(msg(d, "ConstantFieldref")) }

func (c ConstantFieldref) String() string {
	return fmt.Sprintf("<Fieldref: %s, %s>", *c.Clazz, *c.NameAndType)
}

type ConstantMethodref struct {
	Clazz       *Data
	NameAndType *Data
	baseData
}

func (c *ConstantMethodref) Tag() Tag                              { return CP_METHODREF }
func (c *ConstantMethodref) ConstantMethodref() *ConstantMethodref { return c }
func (d *baseData) ConstantMethodref() *ConstantMethodref          { panic(msg(d, "ConstantMethodref")) }

func (c ConstantMethodref) String() string {
	return fmt.Sprintf("<Methodref: %s, %s>", *c.Clazz, *c.NameAndType)
}
