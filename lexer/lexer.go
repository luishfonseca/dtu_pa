package lexer

import (
	"fmt"
	"io"
	"os"

	"github.com/luishfonseca/dtu_pa/util"
)

type Type int

const (
	EOF Type = iota
	MAGIC
	MINOR_VERSION
	MAJOR_VERSION
	CP_COUNT
	CP_INDEX
	CP_INFO_TAG
	CP_UTF8
	CP_INT
)

type Token struct {
	Type  Type
	Bytes []byte
}

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input   io.Reader
	tokenCh chan<- Token
	sc      stackedCounter
	curr    []byte
	err     error
}

type ConfigProvider interface {
	GetClassFile() string
}

func New(cfg ConfigProvider, tokenCh chan<- Token) (*Lexer, error) {
	input, err := os.Open(cfg.GetClassFile())
	if err != nil {
		return nil, err
	}

	return &Lexer{
		input:   input,
		tokenCh: tokenCh,
	}, nil
}

func (l *Lexer) Run() error {
	for state := magic; state != nil; {
		state = state(l)
	}

	close(l.tokenCh)

	if l.err != nil {
		return l.err
	}

	if !l.sc.Empty() {
		return fmt.Errorf("lexer: stacked counter not empty at end of input")
	}

	return nil
}

func (l *Lexer) read(n int) error {
	l.curr = make([]byte, n)

	if _, err := io.ReadFull(l.input, l.curr); err != nil {
		return err
	}

	return nil
}

func (l *Lexer) emit(t Type) {
	if l.curr == nil {
		panic("cannot emit a nil token")
	}

	l.tokenCh <- Token{
		Type:  t,
		Bytes: l.curr,
	}

	l.curr = nil
}

// The magic item supplies the magic number identifying the class file format; it has the value 0xCAFEBABE.
func magic(l *Lexer) stateFn {
	if err := l.read(4); err != nil {
		l.err = err
		return nil
	}

	l.emit(MAGIC)

	return version
}

// The values of the minor_version and major_version items are the minor and major version numbers of this class file.
func version(l *Lexer) stateFn {
	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	l.emit(MINOR_VERSION)

	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	l.emit(MAJOR_VERSION)

	return constant_pool_count
}

// The value of the constant_pool_count item is equal to the number of entries in the constant_pool table plus one.
func constant_pool_count(l *Lexer) stateFn {
	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	var n uint16
	if err := util.Decode(l.curr, &n); err != nil {
		l.err = err
		return nil
	}

	// The constant_pool table is indexed from 1 to constant_pool_count-1
	l.sc.Push(n - 1)

	l.emit(CP_COUNT)

	return constant_pool
}

func constant_pool(l *Lexer) stateFn {
	if l.sc.Top() == 0 {
		l.sc.Pop()
		return access_flags
	}

	l.sc.Dec()

	if err := l.read(1); err != nil {
		l.err = err
		return nil
	}

	tag := l.curr[0]
	l.emit(CP_INFO_TAG)

	switch tag {
	case 1: // CONSTANT_Utf8
		return constant_utf8_info
	case 3: // CONSTANT_Integer
		return constant_integer_info
	case 7: // CONSTANT_Class
		return constant_pool_indices(1)
	case 9, 10, 12: // CONSTANT_Fieldref, CONSTANT_Methodref, CONSTANT_NameAndType
		return constant_pool_indices(2)
	default:
		l.err = fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag))
		return nil
	}
}

func constant_pool_indices(n int) stateFn {
	return func(l *Lexer) stateFn {
		for range n {
			if err := l.read(2); err != nil {
				l.err = err
				return nil
			}

			l.emit(CP_INDEX)
		}

		return constant_pool
	}
}

func constant_utf8_info(l *Lexer) stateFn {
	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	var n uint16
	if err := util.Decode(l.curr, &n); err != nil {
		l.err = err
		return nil
	}

	l.curr = nil

	if err := l.read(int(n)); err != nil {
		l.err = err
		return nil
	}

	l.emit(CP_UTF8)

	return constant_pool
}

func constant_integer_info(l *Lexer) stateFn {
	if err := l.read(4); err != nil {
		l.err = err
		return nil
	}

	l.emit(CP_INT)

	return constant_pool
}

func access_flags(l *Lexer) stateFn {
	return nil
}
