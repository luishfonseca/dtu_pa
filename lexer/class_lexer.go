package lexer

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/state"
	"github.com/luishfonseca/dtu_pa/util"
)

// The magic item supplies the magic number identifying the class file format; it has the value 0xCAFEBABE.
func magic(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(4); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(MAGIC)

	return version
}

// The values of the minor_version and major_version items are the minor and major version numbers of this class file.
func version(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(2); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(MINOR_VERSION)

	if err := l.read(2); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(MAJOR_VERSION)

	return constantPoolCount
}

// The value of the constant_pool_count item is equal to the number of entries in the constant_pool table plus one.
func constantPoolCount(l *Lexer) state.Fn[*Lexer] {
	return count(CP_COUNT, func(l *Lexer) state.Fn[*Lexer] {
		l.sc.dec() // The constant_pool table is indexed from 1 to constant_pool_count-1
		return constantPool
	})
}

// Java Virtual Machine instructions do not rely on the run-time layout of classes, interfaces, class instances, or arrays. Instead, instructions refer to symbolic information in the constant_pool table.
func constantPool(l *Lexer) state.Fn[*Lexer] {
	return repeatUntil(constantPoolInfo, access)
}

// Each item in the constant_pool table must begin with a 1-byte tag indicating the kind of cp_info entry. The contents of the info array vary with the value of tag.
func constantPoolInfo(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(1); err != nil {
		return state.Fail[*Lexer](err)
	}

	tag := l.curr[0]
	l.emit(CP_INFO_TAG)

	switch tag {
	case 1: // CONSTANT_Utf8
		return constantUtf8Info
	case 3: // CONSTANT_Integer
		return constantIntegerInfo
	case 7: // CONSTANT_Class
		return constantPoolIndices(l, 1, constantPool)
	case 9, 10, 12: // CONSTANT_Fieldref, CONSTANT_Methodref, CONSTANT_NameAndType
		return constantPoolIndices(l, 2, constantPool)
	default:
		return state.Fail[*Lexer](fmt.Errorf("unknown cp_info_tag: %d. See https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140", int(tag)))
	}
}

// The CONSTANT_Utf8_info structure is used to represent constant string values
func constantUtf8Info(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(2); err != nil {
		return state.Fail[*Lexer](err)
	}

	var n uint16
	if err := util.Decode(l.curr, &n); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.curr = nil

	if err := l.read(int(n)); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(CP_UTF8)

	return constantPool
}

// The CONSTANT_Integer_info structure represents a 4-byte numeric (int) constant
func constantIntegerInfo(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(4); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(CP_INT)

	return constantPool
}

func access(l *Lexer) state.Fn[*Lexer] {
	return accessFlags(thisClass)
}

// The value of the this_class item must be a valid index into the constant_pool table.
func thisClass(l *Lexer) state.Fn[*Lexer] {
	return constantPoolIndices(l, 1, superClass)
}

// The value of the this_class item must be a valid index into the constant_pool table.
func superClass(l *Lexer) state.Fn[*Lexer] {
	if err := l.read(2); err != nil {
		return state.Fail[*Lexer](err)
	}

	l.emit(CP_NULLABLE_INDEX)

	return interfacesCount
}

// The value of the interfaces_count item gives the number of direct superinterfaces of this class or interface type.
func interfacesCount(l *Lexer) state.Fn[*Lexer] {
	return count(INTERFACES_COUNT, interfaces)
}

// Each value in the interfaces array must be a valid index into the constant_pool table.
func interfaces(l *Lexer) state.Fn[*Lexer] {
	return repeatUntil(func(l *Lexer) state.Fn[*Lexer] {
		return state.Fail[*Lexer](fmt.Errorf("interfaces not implemented"))
	}, fieldsCount)
}

// The value of the fields_count item gives the number of field_info structures in the fields table.
func fieldsCount(l *Lexer) state.Fn[*Lexer] {
	return count(FIELDS_COUNT, fields)
}

// Each value in the fields table must be a field_info structure giving a complete description of a field in this class or interface.
func fields(l *Lexer) state.Fn[*Lexer] {
	return repeatUntil(field_info, methodsCount)
}

func field_info(l *Lexer) state.Fn[*Lexer] {
	return accessFlags(constantPoolIndices(l, 2, fieldAttributesCount))
}

func fieldAttributesCount(l *Lexer) state.Fn[*Lexer] {
	return count(ATTRIBUTES_COUNT, fieldAttributes)
}

func fieldAttributes(l *Lexer) state.Fn[*Lexer] {
	return attributes(fields)
}

// The value of the methods_count item gives the number of method_info structures in the methods table.
func methodsCount(l *Lexer) state.Fn[*Lexer] {
	return count(METHODS_COUNT, methods)
}

// Each method, including each instance initialization method and the class or interface initialization method, is described by a method_info structure.
func methods(l *Lexer) state.Fn[*Lexer] {
	return repeatUntil(method_info, classAttributesCount)
}

func method_info(l *Lexer) state.Fn[*Lexer] {
	return accessFlags(constantPoolIndices(l, 2, methodAttributesCount))
}

func methodAttributesCount(l *Lexer) state.Fn[*Lexer] {
	return count(ATTRIBUTES_COUNT, methodAttributes)
}

func methodAttributes(l *Lexer) state.Fn[*Lexer] {
	return attributes(methods)
}

func classAttributesCount(l *Lexer) state.Fn[*Lexer] {
	return count(ATTRIBUTES_COUNT, classAttributes)
}

func classAttributes(l *Lexer) state.Fn[*Lexer] {
	return attributes(classEnd)
}
