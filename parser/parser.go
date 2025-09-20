package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

type Parser struct {
	tokenCh <-chan lexer.Token
	reqCh   chan<- data.AttributeHandle
	Class   data.Class
	err     error
}

type ConfigProvider interface {
}

func New(cfg ConfigProvider, tokenCh <-chan lexer.Token, reqCh chan<- data.AttributeHandle) *Parser {
	return &Parser{
		tokenCh: tokenCh,
		reqCh:   reqCh,
	}
}

func (p *Parser) Fail(err error) {
	p.err = err
}

func (p *Parser) Run() error {
	state.Run(p, magic)

	if p.err != nil {
		return p.err
	}

	return nil
}

func (p *Parser) expect(t lexer.TokenType) ([]byte, error) {
	token, ok := <-p.tokenCh
	if !ok {
		return nil, fmt.Errorf("unexpected end of input, expected token type %s", t)
	}

	if token.Type != t {
		return nil, fmt.Errorf("unexpected token type %s, expected %s", token.Type, t)
	}

	return token.Bytes, nil
}

func (p *Parser) expectDecode(t lexer.TokenType, v any) error {
	b, err := p.expect(t)
	if err != nil {
		return err
	}

	if err := util.Decode(b, v); err != nil {
		return err
	}

	return nil
}

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

	p.Class.Version = fmt.Sprintf("%d.%d", M, m)

	return constantPool
}

func constantPool(p *Parser) state.Fn[*Parser] {
	var n uint16
	if err := p.expectDecode(lexer.CP_COUNT, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	// The constant_pool table is indexed from 1 to constant_pool_count-1
	p.Class.ConstantPool = make([]data.CpInfo, n-1)

	for i := range n - 1 {
		var tag uint8
		if err := p.expectDecode(lexer.CP_INFO_TAG, &tag); err != nil {
			return state.Fail[*Parser](err)
		}

		switch tag {
		case 1: // CONSTANT_Utf8
			info := data.CpUtf8Info{}

			if b, err := p.expect(lexer.CP_UTF8); err != nil {
				return state.Fail[*Parser](err)
			} else {
				info.Value = string(b)
			}

			p.Class.ConstantPool[i] = info
		case 3: // CONSTANT_Integer
			info := data.CpIntegerInfo{}

			if err := p.expectDecode(lexer.CP_INT, &info.Value); err != nil {
				return state.Fail[*Parser](err)
			}

			p.Class.ConstantPool[i] = info
		case 7: // CONSTANT_Class
			info := data.CpClassInfo{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.Class.ConstantPool[cpIndex-1]
			p.Class.ConstantPool[i] = info
		case 9: // CONSTANT_Fieldref
			info := data.CpFieldrefInfo{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Class = &p.Class.ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.Class.ConstantPool[cpIndex-1]
			p.Class.ConstantPool[i] = info
		case 10: // CONSTANT_Methodref
			info := data.CpMethodrefInfo{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Class = &p.Class.ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.NameAndType = &p.Class.ConstantPool[cpIndex-1]
			p.Class.ConstantPool[i] = info
		case 12: // CONSTANT_NameAndType
			info := data.CpNameAndTypeInfo{}

			var cpIndex uint16
			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Name = &p.Class.ConstantPool[cpIndex-1]

			if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
				return state.Fail[*Parser](err)
			}

			info.Descriptor = &p.Class.ConstantPool[cpIndex-1]
			p.Class.ConstantPool[i] = info
		default:
			return state.Fail[*Parser](fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag)))
		}
	}

	return access
}

func access(p *Parser) state.Fn[*Parser] {
	if err := p.expectDecode(lexer.ACCESS_FLAGS, &p.Class.AccessFlags); err != nil {
		return state.Fail[*Parser](err)
	}

	return thisClass
}

func thisClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	p.Class.ThisClass = &p.Class.ConstantPool[cpIndex-1]

	return superClass
}

func superClass(p *Parser) state.Fn[*Parser] {
	var cpIndex uint16
	if err := p.expectDecode(lexer.CP_NULLABLE_INDEX, &cpIndex); err != nil {
		return state.Fail[*Parser](err)
	}

	if cpIndex != 0 {
		p.Class.SuperClass = &p.Class.ConstantPool[cpIndex-1]
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

	info.Name = &p.Class.ConstantPool[cpIndex-1]

	if err := p.expectDecode(lexer.CP_INDEX, &cpIndex); err != nil {
		return nil, err
	}

	info.Descriptor = &p.Class.ConstantPool[cpIndex-1]

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

	name := &p.Class.ConstantPool[cpIndex-1]

	var begin int64
	if err := p.expectDecode(lexer.ATTRIBUTE_BEGIN, &begin); err != nil {
		return nil, err
	}

	return data.NewAttributeHandle((*name).Utf8Info().Value, begin), nil
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
			p.Class.Fields = append(p.Class.Fields, *field)
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
			p.Class.Methods = append(p.Class.Methods, *method)
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
			p.Class.Attributes = append(p.Class.Attributes, *attr)
		}
	}

	return end
}

func end(p *Parser) state.Fn[*Parser] {
	_, err := p.expect(lexer.EOF)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	return nil
}
