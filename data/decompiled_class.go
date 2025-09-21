package data

type DecompiledClass struct {
	Version      string
	ConstantPool []Data
	AccessFlags  AccessFlags
	ThisClass    *Data
	SuperClass   *Data
	Fields       []MemberInfo
	Methods      []MemberInfo
	Attributes   []AttributeHandle
	baseData
}

func (DecompiledClass) Tag() Tag {
	return DECOMPILED_CLASS
}

func (c *DecompiledClass) DecompiledClass() *DecompiledClass {
	return c
}

func (c DecompiledClass) String() string {
	str := "DecompiledClass {\n"
	str += "  Version: " + c.Version + "\n"
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
