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
			return attributeCode(attr)
		default:
			return state.Fail[*Parser](fmt.Errorf("unimplemented attribute: %s", attr.AttributeTag))
		}
	}
}

func attributeCode(attr data.AttributeHandle) state.Fn[*Parser] {
	return func(p *Parser) state.Fn[*Parser] {
		var maxStack uint16
		if err := p.readDecode(&maxStack); err != nil {
			return state.Fail[*Parser](err)
		}

		var maxLocals uint16
		if err := p.readDecode(&maxLocals); err != nil {
			return state.Fail[*Parser](err)
		}

		var codeHandle data.BytecodeHandle
		if err := p.readDecode(&codeHandle.Length); err != nil {
			return state.Fail[*Parser](err)
		}

		if begin, err := p.input.Seek(0, io.SeekCurrent); err != nil {
			return state.Fail[*Parser](err)
		} else {
			codeHandle.Begin = begin
		}

		if _, err := p.input.Seek(int64(codeHandle.Length), io.SeekCurrent); err != nil {
			return state.Fail[*Parser](err)
		}

		var n uint16
		if err := p.readDecode(&n); err != nil {
			return state.Fail[*Parser](err)
		}

		exceptionTable := make([]data.ExceptionTableEntry, n)
		for i := range n {
			if err := p.readDecode(&exceptionTable[i].StartPC); err != nil {
				return state.Fail[*Parser](err)
			}

			if err := p.readDecode(&exceptionTable[i].EndPC); err != nil {
				return state.Fail[*Parser](err)
			}

			if err := p.readDecode(&exceptionTable[i].HandlerPC); err != nil {
				return state.Fail[*Parser](err)
			}

			var idx uint16
			if err := p.readDecode(&idx); err != nil {
				return state.Fail[*Parser](err)
			}

			if idx != 0 {
				fmt.Println("CatchType index:", idx)
				exceptionTable[i].CatchType = p.class.ConstantPool[idx-1].ConstantClass()
			}
		}

		if err := p.readDecode(&n); err != nil {
			return state.Fail[*Parser](err)
		}

		attrs := make([]data.AttributeHandle, n)
		for i := range n {
			if attr, err := parseAttribute(p); err != nil {
				return state.Fail[*Parser](err)
			} else {
				attrs[i] = *attr
			}
		}

		p.attributes[attr] = &data.AttributeCode{
			MaxStack:       maxStack,
			MaxLocals:      maxLocals,
			CodeHandle:     codeHandle,
			ExceptionTable: exceptionTable,
			Attributes:     attrs,
		}

		p.dataCh <- p.attributes[attr]

		return waitReq
	}
}
