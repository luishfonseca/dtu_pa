package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/state"
)

func magic(p *Parser) state.Fn[*Parser] {
	b, err := p.read(4)
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

	if err := p.readDecode(&m); err != nil {
		return state.Fail[*Parser](err)
	}

	if err := p.readDecode(&M); err != nil {
		return state.Fail[*Parser](err)
	}

	p.class.Version = fmt.Sprintf("%d.%d", M, m)

	return constantPool
}

func constantPool(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.readDecode(&n); err != nil {
		return state.Fail[*Parser](err)
	}

	// The constant_pool table is indexed from 1 to constant_pool_count-1
	p.class.ConstantPool = make([]data.Data, n-1)

	for i := range n - 1 {
		var tag uint8
		if err := p.readDecode(&tag); err != nil {
			return state.Fail[*Parser](err)
		}

		switch tag {
		case 1: // CONSTANT_Utf8
			info := &data.ConstantUtf8{}

			var n uint16
			if err := p.readDecode(&n); err != nil {
				return state.Fail[*Parser](err)
			}

			if b, err := p.read(int(n)); err != nil {
				return state.Fail[*Parser](err)
			} else {
				info.Value = string(b)
			}

			p.class.ConstantPool[i] = info
		case 3: // CONSTANT_Integer
			info := &data.ConstantInteger{}

			if err := p.readDecode(&info.Value); err != nil {
				return state.Fail[*Parser](err)
			}

			p.class.ConstantPool[i] = info
		case 7: // CONSTANT_Class
			info := &data.ConstantClass{}

			var cpIndex uint16
			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.class.ConstantPool[cpIndex-1]
			p.class.ConstantPool[i] = info
		case 9: // CONSTANT_Fieldref
			info := &data.ConstantFieldref{}

			var cpIndex uint16
			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Clazz = &p.class.ConstantPool[cpIndex-1]

			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.class.ConstantPool[cpIndex-1]
			p.class.ConstantPool[i] = info
		case 10: // CONSTANT_Methodref
			info := &data.ConstantMethodref{}

			var cpIndex uint16
			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Clazz = &p.class.ConstantPool[cpIndex-1]

			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.class.ConstantPool[cpIndex-1]
			p.class.ConstantPool[i] = info
		case 12: // CONSTANT_NameAndType
			info := &data.ConstantNameAndType{}

			var cpIndex uint16
			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.class.ConstantPool[cpIndex-1]

			if err := p.readDecode(&cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Descriptor = &p.class.ConstantPool[cpIndex-1]
			p.class.ConstantPool[i] = info
		default:
			return state.Fail[*Parser](fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag)))
		}
	}

	return access
}

func access(p *Parser) state.Fn[*Parser] {
	if err := p.readDecode(&p.class.AccessFlags); err != nil {
		return state.Fail[*Parser](err)
	}

	return thisClass
}

func thisClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.readDecode(&cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	p.class.ThisClass = *p.class.ConstantPool[cpIndex-1].ConstantClass()

	return superClass
}

func superClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.readDecode(&cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	if cpIndex != 0 {
		p.class.SuperClass = p.class.ConstantPool[cpIndex-1].ConstantClass()
	}

	return interfaces
}

func interfaces(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.readDecode(&n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		return state.Fail[*Parser](fmt.Errorf("interfaces not implemented"))
	}

	return fields
}

func fields(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.readDecode(&n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		if field, err := parseMember(p, data.FIELD); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.Fields = append(p.class.Fields, *field)
		}
	}

	return methods
}

func methods(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.readDecode(&n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		if method, err := parseMember(p, data.METHOD); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.Methods = append(p.class.Methods, *method)
		}
	}

	return attributes
}

func attributes(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.readDecode(&n); err != nil {
		return state.Fail[*Parser](err)
	}

	p.class.Attributes = make(map[data.Tag]*data.AttributeHandle)
	for range n {
		if attr, err := parseAttribute(p); err != nil {
			return state.Fail[*Parser](err)
		} else {
			p.class.Attributes[attr.AttributeTag] = attr
		}
	}

	return classEnd
}
