package lexer

import "fmt"

type TokenType int

const (
	EOF TokenType = iota
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
	FIELDS_COUNT
	METHODS_COUNT
	ATTRIBUTES_COUNT
	ATTRIBUTE_BEGIN
)

func (tt TokenType) String() string {
	switch tt {
	case EOF:
		return "EOF"
	case MAGIC:
		return "MAGIC"
	case MINOR_VERSION:
		return "MINOR_VERSION"
	case MAJOR_VERSION:
		return "MAJOR_VERSION"
	case CP_COUNT:
		return "CP_COUNT"
	case CP_INDEX:
		return "CP_INDEX"
	case CP_NULLABLE_INDEX:
		return "CP_NULLABLE_INDEX"
	case CP_INFO_TAG:
		return "CP_INFO_TAG"
	case CP_UTF8:
		return "CP_UTF8"
	case CP_INT:
		return "CP_INT"
	case ACCESS_FLAGS:
		return "ACCESS_FLAGS"
	case INTERFACES_COUNT:
		return "INTERFACES_COUNT"
	case FIELDS_COUNT:
		return "FIELDS_COUNT"
	case METHODS_COUNT:
		return "METHODS_COUNT"
	case ATTRIBUTES_COUNT:
		return "ATTRIBUTES_COUNT"
	case ATTRIBUTE_BEGIN:
		return "ATTRIBUTE_BEGIN"
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}

type Token struct {
	Type  TokenType
	Bytes []byte
}
