package data

import (
	"fmt"
	"maps"
	"slices"
)

type MemberType int

const (
	FIELD MemberType = iota
	METHOD
)

func (m MemberType) String() string {
	switch m {
	case FIELD:
		return "Field"
	case METHOD:
		return "Method"
	default:
		return "Unknown"
	}
}

type MemberInfo struct {
	MemberType  MemberType
	AccessFlags AccessFlags
	Name        ConstantUtf8
	Descriptor  ConstantUtf8
	Attributes  map[Tag]*AttributeHandle
}

func (m MemberInfo) String() string {
	return fmt.Sprintf("<%s: %s %s %v> -> %v", m.MemberType, m.Name, m.Descriptor, m.AccessFlags, slices.Collect(maps.Values(m.Attributes)))
}
