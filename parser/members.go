package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/util"
)

type memberType int

const (
	FIELD memberType = iota
	METHOD
)

func (m memberType) String() string {
	switch m {
	case FIELD:
		return "Field"
	case METHOD:
		return "Method"
	default:
		return "Unknown"
	}
}

type memberInfo struct {
	memberType  memberType
	accessFlags accessFlags
	name        *cpInfo
	descriptor  *cpInfo
	attributes  []cpInfo
}

func newMemberInfo(m memberType, bAccessFlags []byte, bName []byte, bDescriptor []byte, cp []cpInfo) (*memberInfo, error) {
	var accessFlags accessFlags
	if err := util.Decode(bAccessFlags, &accessFlags); err != nil {
		return nil, err
	}

	name, err := resolveIndex(cp, bName)
	if err != nil {
		return nil, err
	}

	descriptor, err := resolveIndex(cp, bDescriptor)
	if err != nil {
		return nil, err
	}

	return &memberInfo{
		memberType:  m,
		accessFlags: accessFlags,
		name:        name,
		descriptor:  descriptor,
	}, nil
}

func (f *memberInfo) appendAttribute(bName []byte, cp []cpInfo) error {
	name, err := resolveIndex(cp, bName)
	if err != nil {
		return err
	}

	f.attributes = append(f.attributes, *name)

	return nil
}

func (f memberInfo) String() string {
	return fmt.Sprintf("<%s: %s %s %v> -> %v", f.memberType, *f.name, *f.descriptor, f.accessFlags, f.attributes)
}
