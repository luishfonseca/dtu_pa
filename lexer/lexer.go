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
	CP_NULLABLE_INDEX
	CP_INFO_TAG
	CP_UTF8
	CP_INT
	ACCESS_FLAGS
	INTERFACES_COUNT
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

	if !l.sc.empty() {
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

	return constantPoolCount
}

func count(emit Type, next stateFn) stateFn {
	return func(l *Lexer) stateFn {
		if err := l.read(2); err != nil {
			l.err = err
			return nil
		}

		var n uint16
		if err := util.Decode(l.curr, &n); err != nil {
			l.err = err
			return nil
		}

		l.sc.push(n)

		l.emit(emit)

		return next
	}
}

// The value of the constant_pool_count item is equal to the number of entries in the constant_pool table plus one.
func constantPoolCount(l *Lexer) stateFn {
	return count(CP_COUNT, func(l *Lexer) stateFn {
		l.sc.dec() // The constant_pool table is indexed from 1 to constant_pool_count-1
		return constantPool
	})
}

// Java Virtual Machine instructions do not rely on the run-time layout of classes, interfaces, class instances, or arrays. Instead, instructions refer to symbolic information in the constant_pool table.
func constantPool(l *Lexer) stateFn {
	if l.sc.top() == 0 {
		l.sc.pop()
		return accessFlags
	}

	l.sc.dec()

	if err := l.read(1); err != nil {
		l.err = err
		return nil
	}

	tag := l.curr[0]
	l.emit(CP_INFO_TAG)

	switch tag {
	case 1: // CONSTANT_Utf8
		return constantUtf8Info
	case 3: // CONSTANT_Integer
		return constantIntegerInfo
	case 7: // CONSTANT_Class
		return constantPoolIndices(1, constantPool)
	case 9, 10, 12: // CONSTANT_Fieldref, CONSTANT_Methodref, CONSTANT_NameAndType
		return constantPoolIndices(2, constantPool)
	default:
		l.err = fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag))
		return nil
	}
}

// Captures n indexes pointing the constant_pool table and continues to next.
func constantPoolIndices(n int, next stateFn) stateFn {
	return func(l *Lexer) stateFn {
		for range n {
			if err := l.read(2); err != nil {
				l.err = err
				return nil
			}

			l.emit(CP_INDEX)
		}

		return next
	}
}

// The CONSTANT_Utf8_info structure is used to represent constant string values
func constantUtf8Info(l *Lexer) stateFn {
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

	return constantPool
}

// The CONSTANT_Integer_info structure represents a 4-byte numeric (int) constant
func constantIntegerInfo(l *Lexer) stateFn {
	if err := l.read(4); err != nil {
		l.err = err
		return nil
	}

	l.emit(CP_INT)

	return constantPool
}

// The value of the access_flags item is a mask of flags used to denote access permission to and properties of this field.
func accessFlags(l *Lexer) stateFn {
	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	l.emit(ACCESS_FLAGS)

	return thisClass
}

// The value of the this_class item must be a valid index into the constant_pool table.
func thisClass(l *Lexer) stateFn {
	return constantPoolIndices(1, superClass)
}

// The value of the this_class item must be a valid index into the constant_pool table.
func superClass(l *Lexer) stateFn {
	if err := l.read(2); err != nil {
		l.err = err
		return nil
	}

	l.emit(CP_NULLABLE_INDEX)

	return interfacesCount
}

// The value of the interfaces_count item gives the number of direct superinterfaces of this class or interface type.
func interfacesCount(l *Lexer) stateFn {
	return count(INTERFACES_COUNT, interfaces)
}

// Each value in the interfaces array must be a valid index into the constant_pool table.
func interfaces(l *Lexer) stateFn {
	if l.sc.top() == 0 {
		l.sc.pop()
		return fieldsCount
	}

	l.err = fmt.Errorf("interfaces not implemented")
	return nil
}

// The value of the fields_count item gives the number of field_info structures in the fields table.
func fieldsCount(l *Lexer) stateFn {
	return nil
}
