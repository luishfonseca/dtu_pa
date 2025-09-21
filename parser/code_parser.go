package parser

import (
	"io"

	"github.com/luishfonseca/dtu_pa/data"
	"github.com/luishfonseca/dtu_pa/state"
)

func bytecode(code data.BytecodeHandle) state.Fn[*Parser] {
	return func(p *Parser) state.Fn[*Parser] {
		if _, err := p.input.Seek(code.Begin, io.SeekStart); err != nil {
			return state.Fail[*Parser](err)
		}

		p.codes[code] = &data.Bytecode{}

		remaining := int(code.Length)

		for remaining > 0 {
			b, err := p.read(1)
			if err != nil {
				return state.Fail[*Parser](err)
			}

			op := data.Op{Code: data.OpCode(b[0])}
			nArgs, err := op.Code.NArgs()
			if err != nil {
				return state.Fail[*Parser](err)
			}

			if nArgs > 0 {
				if b, err := p.read(nArgs); err != nil {
					return state.Fail[*Parser](err)
				} else {
					op.Arg = b
				}
			}

			p.codes[code].Ops = append(p.codes[code].Ops, op)

			remaining -= 1 + nArgs
		}

		p.dataCh <- p.codes[code]
		return waitReq
	}
}
