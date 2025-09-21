package lexer

import (
	"io"

	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

func count(emit TokenType, next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return func(l *Lexer) state.Fn[*Lexer] {
		if err := l.read(2); err != nil {
			return state.Fail[*Lexer](err)
		}

		var n uint16
		if err := util.Decode(l.curr, &n); err != nil {
			return state.Fail[*Lexer](err)
		}

		l.sc.push(n)

		l.emit(emit)

		return next
	}
}

func repeatUntil(element state.Fn[*Lexer], next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return func(l *Lexer) state.Fn[*Lexer] {
		if l.sc.top() == 0 {
			l.sc.pop()
			return next
		}

		l.sc.dec()

		return element // must eventually return to this stateFn
	}
}

// Captures n indexes pointing the constant_pool table and continues to next.
func constantPoolIndices(l *Lexer, n int, next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return func(l *Lexer) state.Fn[*Lexer] {
		for range n {
			if err := l.read(2); err != nil {
				return state.Fail[*Lexer](err)
			}

			l.emit(CP_INDEX)
		}

		return next
	}
}

// The value of the access_flags item is a mask of flags used to denote access permission to and properties of this field.
func accessFlags(next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return func(l *Lexer) state.Fn[*Lexer] {
		if err := l.read(2); err != nil {
			return state.Fail[*Lexer](err)
		}

		l.emit(ACCESS_FLAGS)

		return next
	}
}

func attributes(next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return func(l *Lexer) state.Fn[*Lexer] {
		return repeatUntil(attributeHandle(l, attributes(next)), next)
	}
}

func attributeHandle(l *Lexer, next state.Fn[*Lexer]) state.Fn[*Lexer] {
	return constantPoolIndices(l, 1, func(l *Lexer) state.Fn[*Lexer] {
		if err := l.read(4); err != nil {
			return state.Fail[*Lexer](err)
		}

		var size uint32
		if err := util.Decode(l.curr, &size); err != nil {
			return state.Fail[*Lexer](err)
		}

		l.curr = nil

		begin, err := l.input.Seek(0, io.SeekCurrent) // Mark current position
		if err != nil {
			return state.Fail[*Lexer](err)
		}

		_, err = l.input.Seek(int64(size), io.SeekCurrent) // Skip attribute content
		if err != nil {
			return state.Fail[*Lexer](err)
		}

		bBegin, err := util.Encode(begin)
		if err != nil {
			return state.Fail[*Lexer](err)
		}
		l.curr = bBegin
		l.emit(ATTRIBUTE_BEGIN)

		return next
	})
}
