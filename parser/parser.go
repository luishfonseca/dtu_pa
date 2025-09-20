package parser

import (
	"fmt"
	"sort"

	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

type Parser struct {
	tokenCh <-chan lexer.Token
	cp      []cpInfo
	fields  []fieldInfo
	data    map[string]map[string]any
	err     error
}

type ConfigProvider interface {
}

func New(cfg ConfigProvider, tokenCh <-chan lexer.Token) *Parser {
	return &Parser{
		tokenCh: tokenCh,
	}
}

func (p *Parser) Fail(err error) {
	p.err = err
}

func (p *Parser) PrintData() {
	ks1 := make([]string, 0, len(p.data))
	for k := range p.data {
		ks1 = append(ks1, k)
	}

	sort.Strings(ks1)

	for _, k1 := range ks1 {
		fmt.Printf("=== %s ===\n", k1)

		ks2 := make([]string, 0, len(p.data[k1]))
		for k := range p.data[k1] {
			ks2 = append(ks2, k)
		}

		sort.Strings(ks2)

		for _, k2 := range ks2 {
			fmt.Printf("%s: %v\n", k2, p.data[k1][k2])
		}
	}
}

func (p *Parser) Run() error {
	state.Run(p, magic)

	if p.err != nil {
		return p.err
	}

	return nil
}

func (p *Parser) expect(t lexer.Type) ([]byte, error) {
	token, ok := <-p.tokenCh
	if !ok {
		return nil, fmt.Errorf("unexpected end of input, expected token type %d", t)
	}

	if token.Type != t {
		return nil, fmt.Errorf("unexpected token type %d, expected %d", token.Type, t)
	}

	return token.Bytes, nil
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
	mb, err := p.expect(lexer.MINOR_VERSION)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	Mb, err := p.expect(lexer.MAJOR_VERSION)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	var m uint16
	if err := util.Decode(mb, &m); err != nil {
		return state.Fail[*Parser](err)
	}

	var M uint16
	if err := util.Decode(Mb, &M); err != nil {
		return state.Fail[*Parser](err)
	}

	p.data = make(map[string]map[string]any)
	p.data["0. Class"] = make(map[string]any)
	p.data["0. Class"]["version"] = fmt.Sprintf("%d.%d", M, m)

	return constantPool
}

func constantPool(p *Parser) state.Fn[*Parser] {
	bn, err := p.expect(lexer.CP_COUNT)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	var n uint16
	if err := util.Decode(bn, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	p.cp = make([]cpInfo, n-1) // The constant_pool table is indexed from 1 to constant_pool_count-1
	for i := range n - 1 {
		b, err := p.expect(lexer.CP_INFO_TAG)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		tag := b[0]
		switch tag {
		case 1: // CONSTANT_Utf8
			b, err := p.expect(lexer.CP_UTF8)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			p.cp[i] = newConstantUtf8Info(b)
		case 3: // CONSTANT_Integer
			b, err := p.expect(lexer.CP_INT)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if info, err := newConstantIntegerInfo(b); err != nil {
				return state.Fail[*Parser](err)
			} else {
				p.cp[i] = info
			}
		case 7: // CONSTANT_Class
			b, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if info, err := newConstantClassInfo(b, p.cp); err != nil {
				return state.Fail[*Parser](err)
			} else {
				p.cp[i] = info
			}
		case 9: // CONSTANT_Fieldref
			bClass, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			bNameAndType, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if info, err := newConstantFieldrefInfo(bClass, bNameAndType, p.cp); err != nil {
				return state.Fail[*Parser](err)
			} else {
				p.cp[i] = info
			}
		case 10: // CONSTANT_Methodref
			bClass, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			bNameAndType, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if info, err := newConstantMethodrefInfo(bClass, bNameAndType, p.cp); err != nil {
				return state.Fail[*Parser](err)
			} else {
				p.cp[i] = info
			}
		case 12: // CONSTANT_NameAndType
			bName, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			bDescriptor, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if info, err := newConstantNameAndTypeInfo(bName, bDescriptor, p.cp); err != nil {
				return state.Fail[*Parser](err)
			} else {
				p.cp[i] = info
			}
		default:
			return state.Fail[*Parser](fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag)))
		}
	}

	p.data["0. Class"]["constant pool count"] = n - 1
	p.data["1. Constant Pool"] = make(map[string]any, len(p.cp))
	for i, v := range p.cp {
		p.data["1. Constant Pool"][fmt.Sprintf("%2d", i+1)] = v.String()
	}

	return access
}

func access(p *Parser) state.Fn[*Parser] {
	b, err := p.expect(lexer.ACCESS_FLAGS)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	var accessFlags accessFlags
	if err := util.Decode(b, &accessFlags); err != nil {
		return state.Fail[*Parser](err)
	}

	p.data["0. Class"]["access flags"] = accessFlags.String()

	return thisClass
}

func thisClass(p *Parser) state.Fn[*Parser] {
	b, err := p.expect(lexer.CP_INDEX)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	if info, err := resolveIndex(p.cp, b); err != nil {
		return state.Fail[*Parser](err)
	} else {
		p.data["0. Class"]["this class"] = (*info).String()
	}

	return superClass
}

func superClass(p *Parser) state.Fn[*Parser] {
	b, err := p.expect(lexer.CP_NULLABLE_INDEX)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	if len(b) == 2 && b[0] == 0 && b[1] == 0 {
		p.data["0. Class"]["super class"] = "None"
	} else if info, err := resolveIndex(p.cp, b); err != nil {
		return state.Fail[*Parser](err)
	} else {
		p.data["0. Class"]["super class"] = (*info).String()
	}

	return interfaces
}

func interfaces(p *Parser) state.Fn[*Parser] {
	bn, err := p.expect(lexer.INTERFACES_COUNT)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	var n uint16
	if err := util.Decode(bn, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		return state.Fail[*Parser](fmt.Errorf("interfaces not implemented"))
	}

	p.data["0. Class"]["interfaces count"] = n

	return fields
}

func fields(p *Parser) state.Fn[*Parser] {
	bn, err := p.expect(lexer.FIELDS_COUNT)
	if err != nil {
		return state.Fail[*Parser](err)
	}

	var n uint16
	if err := util.Decode(bn, &n); err != nil {
		return state.Fail[*Parser](err)
	}

	for range n {
		bAccessFlags, err := p.expect(lexer.ACCESS_FLAGS)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		bName, err := p.expect(lexer.CP_INDEX)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		bDescriptor, err := p.expect(lexer.CP_INDEX)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		info, err := newFieldInfo(bAccessFlags, bName, bDescriptor, p.cp)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		bn, err := p.expect(lexer.ATTRIBUTES_COUNT)
		if err != nil {
			return state.Fail[*Parser](err)
		}

		var n uint16
		if err := util.Decode(bn, &n); err != nil {
			return state.Fail[*Parser](err)
		}

		for range n {
			bAttributeName, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if err := info.appendAttribute(bAttributeName, p.cp); err != nil {
				return state.Fail[*Parser](err)
			}
		}

		p.fields = append(p.fields, *info)
	}

	p.data["0. Class"]["fields count"] = n
	p.data["2. Fields"] = make(map[string]any, len(p.fields))
	for i, v := range p.fields {
		p.data["2. Fields"][fmt.Sprintf("%2d", i+1)] = v.String()
	}

	return methods
}

func methods(p *Parser) state.Fn[*Parser] {
	return nil
}
