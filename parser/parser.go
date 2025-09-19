package parser

import (
	"fmt"
	"sort"

	"github.com/luishfonseca/dtu_pa/lexer"
	"github.com/luishfonseca/dtu_pa/util"
)

type stateFn func(*Parser) stateFn

type Parser struct {
	tokenCh <-chan lexer.Token
	cp      []cpInfo
	data    []map[string]any
	err     error
}

type ConfigProvider interface {
}

func New(cfg ConfigProvider, tokenCh <-chan lexer.Token) *Parser {
	return &Parser{
		tokenCh: tokenCh,
	}
}

func (p *Parser) PrintData() {
	for i, m := range p.data {
		if i != 0 {
			fmt.Println("----")
		}

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf("%s: %v\n", k, m[k])
		}
	}
}

func (p *Parser) Run() error {
	for state := magic; state != nil; {
		state = state(p)
	}

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

func magic(p *Parser) stateFn {
	b, err := p.expect(lexer.MAGIC)
	if err != nil {
		p.err = err
		return nil
	}

	if string(b) != "\xCA\xFE\xBA\xBE" {
		p.err = fmt.Errorf("invalid magic number: %x", b)
		return nil
	}

	return version
}

func version(p *Parser) stateFn {
	mb, err := p.expect(lexer.MINOR_VERSION)
	if err != nil {
		p.err = err
		return nil
	}

	Mb, err := p.expect(lexer.MAJOR_VERSION)
	if err != nil {
		p.err = err
		return nil
	}

	var m uint16
	if err := util.Decode(mb, &m); err != nil {
		p.err = err
		return nil
	}

	var M uint16
	if err := util.Decode(Mb, &M); err != nil {
		p.err = err
		return nil
	}

	p.data = append(p.data, make(map[string]any))
	p.data[0]["version"] = fmt.Sprintf("%d.%d", M, m)

	return constantPool
}

func constantPool(p *Parser) stateFn {
	bn, err := p.expect(lexer.CP_COUNT)
	if err != nil {
		p.err = err
		return nil
	}

	var n uint16
	if err := util.Decode(bn, &n); err != nil {
		p.err = err
		return nil
	}

	p.cp = make([]cpInfo, n-1) // The constant_pool table is indexed from 1 to constant_pool_count-1
	for i := range n - 1 {
		b, err := p.expect(lexer.CP_INFO_TAG)
		if err != nil {
			p.err = err
			return nil
		}

		tag := b[0]
		switch tag {
		case 1: // CONSTANT_Utf8
			b, err := p.expect(lexer.CP_UTF8)
			if err != nil {
				p.err = err
				return nil
			}

			p.cp[i] = newConstantUtf8Info(b)
		case 3: // CONSTANT_Integer
			b, err := p.expect(lexer.CP_INT)
			if err != nil {
				p.err = err
				return nil
			}

			if info, err := newConstantIntegerInfo(b); err != nil {
				p.err = err
				return nil
			} else {
				p.cp[i] = info
			}
		case 7: // CONSTANT_Class
			b, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			if info, err := newConstantClassInfo(b, p.cp); err != nil {
				p.err = err
				return nil
			} else {
				p.cp[i] = info
			}
		case 9: // CONSTANT_Fieldref
			bClass, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			bNameAndType, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			if info, err := newConstantFieldrefInfo(bClass, bNameAndType, p.cp); err != nil {
				p.err = err
				return nil
			} else {
				p.cp[i] = info
			}
		case 10: // CONSTANT_Methodref
			bClass, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			bNameAndType, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			if info, err := newConstantMethodrefInfo(bClass, bNameAndType, p.cp); err != nil {
				p.err = err
				return nil
			} else {
				p.cp[i] = info
			}
		case 12: // CONSTANT_NameAndType
			bName, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			bDescriptor, err := p.expect(lexer.CP_INDEX)
			if err != nil {
				p.err = err
				return nil
			}

			if info, err := newConstantNameAndTypeInfo(bName, bDescriptor, p.cp); err != nil {
				p.err = err
				return nil
			} else {
				p.cp[i] = info
			}
		default:
			p.err = fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag))
			return nil
		}
	}

	p.data[0]["constant pool count"] = n - 1
	p.data = append(p.data, make(map[string]any))
	for i, v := range p.cp {
		p.data[1][fmt.Sprintf("%2d", i+1)] = v.String()
	}

	return access
}

func access(p *Parser) stateFn {
	b, err := p.expect(lexer.ACCESS_FLAGS)
	if err != nil {
		p.err = err
		return nil
	}

	var af uint16
	if err := util.Decode(b, &af); err != nil {
		p.err = err
		return nil
	}

	afs := accessFlagsFromMask(af)
	p.data[0]["access flags"] = make([]string, len(afs))
	for i, af := range afs {
		p.data[0]["access flags"].([]string)[i] = fmt.Sprintf("0x%04X", af.mask())
	}

	return nil
}
