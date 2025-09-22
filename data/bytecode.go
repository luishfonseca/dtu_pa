package data

import (
	"fmt"
)

type OpCode byte

const (
	OP_NOP           OpCode = 0x00
	OP_ACONST_NULL   OpCode = 0x01
	OP_ICONST_0      OpCode = 0x03
	OP_ICONST_1      OpCode = 0x04
	OP_ICONST_2      OpCode = 0x05
	OP_ICONST_3      OpCode = 0x06
	OP_ICONST_4      OpCode = 0x07
	OP_ICONST_5      OpCode = 0x08
	OP_BIPUSH        OpCode = 0x10
	OP_SIPUSH        OpCode = 0x11
	OP_LDC           OpCode = 0x12
	OP_ILOAD         OpCode = 0x15
	OP_ALOAD         OpCode = 0x19
	OP_ILOAD_0       OpCode = 0x1a
	OP_ILOAD_1       OpCode = 0x1b
	OP_ILOAD_2       OpCode = 0x1c
	OP_ILOAD_3       OpCode = 0x1d
	OP_ALOAD_0       OpCode = 0x2a
	OP_ALOAD_1       OpCode = 0x2b
	OP_IALOAD        OpCode = 0x2e
	OP_CALOAD        OpCode = 0x34
	OP_ISTORE        OpCode = 0x36
	OP_ASTORE        OpCode = 0x3a
	OP_ISTORE_0      OpCode = 0x3b
	OP_ISTORE_1      OpCode = 0x3c
	OP_ISTORE_2      OpCode = 0x3d
	OP_ISTORE_3      OpCode = 0x3e
	OP_ASTORE_0      OpCode = 0x4b
	OP_ASTORE_1      OpCode = 0x4c
	OP_ASTORE_2      OpCode = 0x4d
	OP_IASTORE       OpCode = 0x4f
	OP_DUP           OpCode = 0x59
	OP_IADD          OpCode = 0x60
	OP_ISUB          OpCode = 0x64
	OP_IMUL          OpCode = 0x68
	OP_IDIV          OpCode = 0x6c
	OP_IREM          OpCode = 0x70
	OP_IINC          OpCode = 0x84
	OP_I2S           OpCode = 0x93
	OP_IFEQ          OpCode = 0x99
	OP_IFNE          OpCode = 0x9a
	OP_IFGE          OpCode = 0x9c
	OP_IFGT          OpCode = 0x9d
	OP_IF_ICMPEQ     OpCode = 0x9f
	OP_IF_ICMPNE     OpCode = 0xa0
	OP_IF_ICMPLT     OpCode = 0xa1
	OP_IF_ICMPGE     OpCode = 0xa2
	OP_IF_ICMPGT     OpCode = 0xa3
	OP_IF_ICMPLE     OpCode = 0xa4
	OP_GOTO          OpCode = 0xa7
	OP_IRETURN       OpCode = 0xac
	OP_ARETURN       OpCode = 0xb0
	OP_RETURN        OpCode = 0xb1
	OP_GETSTATIC     OpCode = 0xb2
	OP_PUTSTATIC     OpCode = 0xb3
	OP_INVOKEVIRTUAL OpCode = 0xb6
	OP_INVOKESPECIAL OpCode = 0xb7
	OP_INVOKESTATIC  OpCode = 0xb8
	OP_NEW           OpCode = 0xbb
	OP_NEWARRAY      OpCode = 0xbc
	OP_ARRAYLENGTH   OpCode = 0xbe
	OP_ATHROW        OpCode = 0xbf
)

func (o OpCode) String() string {
	switch o {
	case OP_NOP:
		return "nop"
	case OP_ACONST_NULL:
		return "aconst_null"
	case OP_ICONST_0:
		return "iconst_0"
	case OP_ICONST_1:
		return "iconst_1"
	case OP_ICONST_2:
		return "iconst_2"
	case OP_ICONST_3:
		return "iconst_3"
	case OP_ICONST_4:
		return "iconst_4"
	case OP_ICONST_5:
		return "iconst_5"
	case OP_BIPUSH:
		return "bipush"
	case OP_SIPUSH:
		return "sipush"
	case OP_LDC:
		return "ldc"
	case OP_ILOAD:
		return "iload"
	case OP_ALOAD:
		return "aload"
	case OP_ILOAD_0:
		return "iload_0"
	case OP_ILOAD_1:
		return "iload_1"
	case OP_ILOAD_2:
		return "iload_2"
	case OP_ILOAD_3:
		return "iload_3"
	case OP_ALOAD_0:
		return "aload_0"
	case OP_ALOAD_1:
		return "aload_1"
	case OP_IALOAD:
		return "iaload"
	case OP_CALOAD:
		return "caload"
	case OP_ISTORE:
		return "istore"
	case OP_ASTORE:
		return "astore"
	case OP_ISTORE_0:
		return "istore_0"
	case OP_ISTORE_1:
		return "istore_1"
	case OP_ISTORE_2:
		return "istore_2"
	case OP_ASTORE_0:
		return "astore_0"
	case OP_ASTORE_1:
		return "astore_1"
	case OP_ASTORE_2:
		return "astore_2"
	case OP_IASTORE:
		return "iastore"
	case OP_DUP:
		return "dup"
	case OP_IADD:
		return "iadd"
	case OP_ISUB:
		return "isub"
	case OP_IMUL:
		return "imul"
	case OP_IDIV:
		return "idiv"
	case OP_IREM:
		return "irem"
	case OP_IINC:
		return "iinc"
	case OP_I2S:
		return "i2s"
	case OP_IFEQ:
		return "ifeq"
	case OP_IFNE:
		return "ifne"
	case OP_IFGE:
		return "ifge"
	case OP_IFGT:
		return "ifgt"
	case OP_IF_ICMPEQ:
		return "if_icmpeq"
	case OP_IF_ICMPNE:
		return "if_icmpne"
	case OP_IF_ICMPLT:
		return "if_icmplt"
	case OP_IF_ICMPGE:
		return "if_icmpge"
	case OP_IF_ICMPGT:
		return "if_icmpgt"
	case OP_IF_ICMPLE:
		return "if_icmple"
	case OP_GOTO:
		return "goto"
	case OP_IRETURN:
		return "ireturn"
	case OP_GETSTATIC:
		return "getstatic"
	case OP_PUTSTATIC:
		return "putstatic"
	case OP_ARETURN:
		return "areturn"
	case OP_RETURN:
		return "return"
	case OP_INVOKEVIRTUAL:
		return "invokevirtual"
	case OP_INVOKESPECIAL:
		return "invokespecial"
	case OP_INVOKESTATIC:
		return "invokestatic"
	case OP_NEW:
		return "new"
	case OP_NEWARRAY:
		return "newarray"
	case OP_ARRAYLENGTH:
		return "arraylength"
	case OP_ATHROW:
		return "athrow"
	default:
		return fmt.Sprintf("UNKNOWN_OP_%02X", byte(o))
	}
}

func (op OpCode) NArgs() (int, error) {
	switch byte(op) {
	case 0x00, 0x01, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x12, 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b,
		0x2e, 0x34, 0x3b, 0x3c, 0x3d, 0x3e, 0x4b, 0x4c, 0x4d, 0x4f, 0x59, 0x60, 0x64, 0x68, 0x6c,
		0x70, 0x93, 0xac, 0xb0, 0xb1, 0xbe, 0xbf:
		return 0, nil
	case 0x10, 0x15, 0x19, 0x36, 0x3a, 0xbc:
		return 1, nil
	case 0x11, 0x84, 0x99, 0x9a, 0x9c, 0x9d, 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa7, 0xb2, 0xb3,
		0xb6, 0xb7, 0xb8, 0xbb:
		return 2, nil
	default:
		return -1, fmt.Errorf("unimplemented bytecode: 0x%02x", byte(op))
	}
}

type Op struct {
	Code OpCode
	Arg  []byte
}

func (o Op) String() string {
	if o.Arg != nil {
		return fmt.Sprintf("%s %v", o.Code, o.Arg)
	}
	return o.Code.String()
}

type Bytecode struct {
	Ops []Op
	baseData
}

func (b *Bytecode) Tag() Tag            { return BYTECODE }
func (b *Bytecode) Bytecode() *Bytecode { return b }
func (d *baseData) Bytecode() *Bytecode { panic(msg(d, "Bytecode")) }

func (b Bytecode) String() string {
	str := "Bytecode["
	for _, op := range b.Ops {
		str += fmt.Sprint("\n ", op)
	}
	str += "]"
	return str
}
