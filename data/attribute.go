package data

import (
	"fmt"
	"os"
)

type AttributeType int

const (
	UNKNOWN AttributeType = iota
	CODE
	SOURCE_FILE
	RUNTIME_VISIBLE_ANNOTATIONS
	INNER_CLASSES
)

func (at AttributeType) String() string {
	switch at {
	case UNKNOWN:
		return "Unknown"
	case CODE:
		return "Code"
	case SOURCE_FILE:
		return "SourceFile"
	case RUNTIME_VISIBLE_ANNOTATIONS:
		return "RuntimeVisibleAnnotations"
	case INNER_CLASSES:
		return "InnerClasses"
	default:
		return fmt.Sprintf("AttributeType(%d)", int(at))
	}
}

func AttributeTypeFromName(name string) AttributeType {
	switch name {
	case "Code":
		return CODE
	case "SourceFile":
		return SOURCE_FILE
	case "RuntimeVisibleAnnotations":
		return RUNTIME_VISIBLE_ANNOTATIONS
	case "InnerClasses":
		return INNER_CLASSES
	default:
		fmt.Fprintf(os.Stderr, "warning: unknown attribute name: %q\n", name)
		return UNKNOWN
	}
}

type AttributeHandle struct {
	Type  AttributeType
	Begin int64
}

func NewAttributeHandle(t string, begin int64) *AttributeHandle {
	return &AttributeHandle{
		Type:  AttributeTypeFromName(t),
		Begin: begin,
	}
}

func (a AttributeHandle) String() string {
	return fmt.Sprintf("<%s @ %d>", a.Type, a.Begin)
}
