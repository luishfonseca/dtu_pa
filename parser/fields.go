package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/util"
)

type fieldInfo struct {
	accessFlags accessFlags
	name        *cpInfo
	descriptor  *cpInfo
	attributes  []*cpInfo
}

func newFieldInfo(bAccessFlags []byte, bName []byte, bDescriptor []byte, cp []cpInfo) (*fieldInfo, error) {
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

	return &fieldInfo{
		accessFlags: accessFlags,
		name:        name,
		descriptor:  descriptor,
	}, nil
}

func (f *fieldInfo) appendAttribute(bName []byte, cp []cpInfo) error {
	name, err := resolveIndex(cp, bName)
	if err != nil {
		return err
	}

	f.attributes = append(f.attributes, name)

	return nil
}

func (f fieldInfo) String() string {
	return fmt.Sprintf("<Field: %s %s %v>", *f.name, *f.descriptor, f.accessFlags)
}
