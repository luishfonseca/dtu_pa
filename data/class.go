package data

import "fmt"

type Class struct {
	Version      string
	ConstantPool []Data
	AccessFlags  AccessFlags
	ThisClass    ConstantClass
	SuperClass   *ConstantClass
	Fields       []MemberInfo
	Methods      []MemberInfo
	Attributes   map[Tag]*AttributeHandle
	baseData
}

func (c *Class) Tag() Tag         { return CLASS }
func (c *Class) Class() *Class    { return c }
func (d *baseData) Class() *Class { panic(msg(d, "Class")) }

func (c *Class) Method(name string, descriptor string) *MemberInfo {
	for _, method := range c.Methods {
		if method.Name.Value == name && method.Descriptor.Value == descriptor {
			return &method
		}
	}
	return nil
}

func (c Class) String() string {
	str := "Class {\n"
	str += fmt.Sprintln("  Version:", c.Version)
	str += "  ConstantPool: [\n"
	for i, constant := range c.ConstantPool {
		str += fmt.Sprintf("    %2d: %s\n", i+1, constant)
	}
	str += "  ]\n"
	str += fmt.Sprintln("  AccessFlags:", c.AccessFlags)
	str += fmt.Sprintln("  ThisClass:", c.ThisClass)
	if c.SuperClass != nil {
		str += fmt.Sprintln("  SuperClass:", *c.SuperClass)
	} else {
		str += "  SuperClass: None\n"
	}
	str += "  Fields: [\n"
	for _, field := range c.Fields {
		str += fmt.Sprintln("   ", field)
	}
	str += "  ]\n  Methods: [\n"
	for _, method := range c.Methods {
		str += fmt.Sprintln("   ", method)
	}
	str += "  ]\n  Attributes: [\n"
	for _, attr := range c.Attributes {
		str += fmt.Sprintln("   ", attr)
	}
	str += "  ]\n}"

	return str
}
