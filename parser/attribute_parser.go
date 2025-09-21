package parser

import (
	"fmt"
	"io"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/state"
)

func attribute(attr data.AttributeHandle) state.Fn[*Parser] {
	return func(p *Parser) state.Fn[*Parser] {
		// seek to attribute position
		if _, err := p.input.Seek(attr.Begin, io.SeekStart); err != nil {
			return state.Fail[*Parser](err)
		}

		switch attr.AttributeTag {
		case data.ATTR_CODE:
			return attributeCode
		default:
			return state.Fail[*Parser](fmt.Errorf("unimplemented attribute: %s", attr.AttributeTag))
		}
	}
}

func attributeCode(p *Parser) state.Fn[*Parser] {
	return waitReq
}
