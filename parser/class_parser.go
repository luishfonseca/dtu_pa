package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/state"
)

func magic(p *Parser) state.Fn[*Parser] {
	b, err := p.expect(lexer.MAGIC)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	if string(b) != "\xCA\xFE\xBA\xBE" {
		return state.Fail[*Parser](fmt.Errorf("invalid magic number: %x", b))
	}

	return version
}

func version(p *Parser) state.Fn[*Parser] {
	var m, M uint16

	if err := p.expectDecode(lexer.MINOR_VERSION, &m); err != nil {
		return state.Fail[*Parser](err)
	}

	if err := p.expectDecode(lexer.MAJOR_VERSION, &M); err != nil {
		return state.Fail[*Parser](err)
	}

	p.class.DecompiledClass().Version = fmt.Sprintf("%d.%d", M, m)

	return constantPool
}

func constantPool(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.CP_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	// The constant_pool table is indexed from 1 to constant_pool_count-1
	p.class.DecompiledClass().ConstantPool = make([]data.Data, n-1)

	for i := range n - 1 {
		var tag uint8
		if err := p.expectDecode(lexer.CP_INFO_TAG, &tag); err != nil {
			return state.Fail[*Parser](err)
		}

		switch tag {
		case 1: // CONSTANT_Utf8
			info := &data.ConstantUtf8{}

			if b, err := p.expect(lexer.CP_UTF8); err != nil {
				return state.Fail[*Parser](err)
			} else {
				info.Value = string(b)
			}

			p.class.DecompiledClass().ConstantPool[i] = info
		case 3: // CONSTANT_Integer
			info := &data.ConstantInteger{}

			if err := p.expectDecode(lexer.CP_INT, &info.Value); err != nil {
				return state.Fail[*Parser](err)
			}

			p.class.DecompiledClass().ConstantPool[i] = info
		case 7: // CONSTANT_Class
			info := &data.ConstantClass{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.class.DecompiledClass().ConstantPool[cpIndex-1]
			p.class.DecompiledClass().ConstantPool[i] = info
		case 9: // CONSTANT_Fieldref
			info := &data.ConstantFieldref{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Class = &p.class.DecompiledClass().ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.class.DecompiledClass().ConstantPool[cpIndex-1]
			p.class.DecompiledClass().ConstantPool[i] = info
		case 10: // CONSTANT_Methodref
			info := &data.ConstantMethodref{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Class = &p.class.DecompiledClass().ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.class.DecompiledClass().ConstantPool[cpIndex-1]
			p.class.DecompiledClass().ConstantPool[i] = info
		case 12: // CONSTANT_NameAndType
			info := &data.ConstantNameAndType{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.class.DecompiledClass().ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Descriptor = &p.class.DecompiledClass().ConstantPool[cpIndex-1]
			p.class.DecompiledClass().ConstantPool[i] = info
		default:
			return state.Fail[*Parser](fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag)))
		}
	}

	return access
}

func access(p *Parser) state.Fn[*Parser] {
	if err := p.expectDecode(lexer.ACCESS_FLAGS, &p.class.DecompiledClass().AccessFlags); err != nil {
		return state.Fail[*Parser](err)
	}

	return thisClass
}

func thisClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	p.class.DecompiledClass().ThisClass = *p.class.DecompiledClass().ConstantPool[cpIndex-1].ConstantClass()

	return superClass
}

func superClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_NULLABLE_INDEX, &cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	if cpIndex != 0 {
		p.class.DecompiledClass().SuperClass = p.class.DecompiledClass().ConstantPool[cpIndex-1].ConstantClass()
	}

	return interfaces
}

func interfaces(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.INTERFACES_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		return state.Fail[*Parser](fmt.Errorf("interfaces not implemented"))
	}

	return fields
}

func parseMember(p *Parser, m data.MemberType) (*data.MemberInfo, error) {
	info := &data.MemberInfo{
		MemberType: m,
	}

	if err := p.expectDecode(lexer.ACCESS_FLAGS, &info.AccessFlags); err != nil {
		return nil, err
	}

	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return nil, err
	}

	info.Name = *p.class.DecompiledClass().ConstantPool[cpIndex-1].ConstantUtf8()

	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return nil, err
	}

	info.Descriptor = *p.class.DecompiledClass().ConstantPool[cpIndex-1].ConstantUtf8()

	var n uint16
	if err := p.expectDecode(lexer.ATTRIBUTES_COUNT, &n); err != nil {
		return nil, err
	}

	for range n {
		if attr, err := parseAttribute(p); err != nil {
			return nil, err
		} else {
			info.Attributes = append(info.Attributes, *attr)
		}
	}

	return info, nil
}

func parseAttribute(p *Parser) (*data.AttributeHandle, error) {
	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return nil, err
	}

	name := p.class.DecompiledClass().ConstantPool[cpIndex-1].ConstantUtf8().Value

	var begin int64
	if err := p.expectDecode(lexer.ATTRIBUTE_BEGIN, &begin); err != nil {
		return nil, err
	}

	var tag data.Tag
	switch name {
	case "Code":
		tag = data.ATTR_CODE
	case "RuntimeVisibleAnnotations":
		tag = data.ATTR_RUNTIME_VISIBLE_ANNOTATIONS
	case "SourceFile":
		tag = data.ATTR_SOURCE_FILE
	case "InnerClasses":
		tag = data.ATTR_INNER_CLASSES
	default:
		return nil, fmt.Errorf("unknown attribute name: %s", name)
	}

	return &data.AttributeHandle{
		AttributeTag: tag,
		Begin:        begin,
	}, nil
}

func fields(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.FIELDS_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		if field, err := parseMember(p, data.FIELD); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.DecompiledClass().Fields = append(p.class.DecompiledClass().Fields, *field)
		}
	}

	return methods
}

func methods(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.METHODS_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		if method, err := parseMember(p, data.METHOD); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.DecompiledClass().Methods = append(p.class.DecompiledClass().Methods, *method)
		}
	}

	return attributes
}

func attributes(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.ATTRIBUTES_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		if attr, err := parseAttribute(p); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.DecompiledClass().Attributes = append(p.class.DecompiledClass().Attributes, *attr)
		}
	}

	return classEnd
}
