package lexer

import (
	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

func count(emit Type, next state.Fn[*Lexer]) state.Fn[*Lexer] {
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
