package data

import "fmt"

type DecompiledClass struct {
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

func (c *DecompiledClass) DecompiledClass() *DecompiledClass {
	return c
}

func (*DecompiledClass) Tag() Tag {
	return DECOMPILED_CLASS
}

func (c *DecompiledClass) Method(name string, descriptor string) *MemberInfo {
	for _, method := range c.Methods {
		if method.Name.Value == name && method.Descriptor.Value == descriptor {
			return &method
		}
	}
	return nil
}

func (c DecompiledClass) String() string {
	str := "DecompiledClass {\n"
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
	str += "  ]\n}\n"

	return str
}
