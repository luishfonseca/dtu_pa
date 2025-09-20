package data

import "fmt"

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
	Name        *CpInfo
	Descriptor  *CpInfo
	Attributes  []AttributeHandle
}

func (info MemberInfo) String() string {
	return fmt.Sprintf("<%s: %s %s %v> -> %v", info.MemberType, *info.Name, *info.Descriptor, info.AccessFlags, info.Attributes)
}
