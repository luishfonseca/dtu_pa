package data

import "fmt"

type Class struct {
	Version      string
	ConstantPool []CpInfo
	AccessFlags  AccessFlags
	ThisClass    *CpInfo
	SuperClass   *CpInfo
	Fields       []MemberInfo
	Methods      []MemberInfo
	Attributes   []AttributeHandle
}

func (c Class) String() string {
	str := "Class {\n"
	str += "  Version: " + c.Version + "\n"
	str += "  ConstantPool: [\n"
	for i, cp := range c.ConstantPool {
		str += "    #" + fmt.Sprintf("%d", i+1) + " = " + cp.String() + "\n"
	}
	str += "  ]\n"
	str += "  AccessFlags: " + c.AccessFlags.String() + "\n"
	str += "  ThisClass: " + (*c.ThisClass).String() + "\n"
	str += "  SuperClass: " + func() string {
		if c.SuperClass != nil {
			return (*c.SuperClass).String()
		}
		return "None"
	}() + "\n"
	str += "  Fields: [\n"
	for _, field := range c.Fields {
		str += "    " + field.String() + "\n"
	}
	str += "  ]\n"
	str += "  Methods: [\n"
	for _, method := range c.Methods {
		str += "    " + method.String() + "\n"
	}
	str += "  ]\n"
	str += "  Attributes: [\n"
	for _, attr := range c.Attributes {
		str += "    " + attr.String() + "\n"
	}
	str += "  ]\n"
	str += "}\n"

	return str
}
