package parser

import (
	"fmt"

	"github.com/luishfonseca/dtu_pa/util"
)

type cpInfo fmt.Stringer

func resolveIndex(cp []cpInfo, bIndex []byte) (*cpInfo, error) {
	var index uint16
	if err := util.Decode(bIndex, &index); err != nil {
		return nil, fmt.Errorf("could not decode constant pool index: %w", err)
	}

	if int(index) == 0 || int(index) >= len(cp) {
		return nil, fmt.Errorf("constant pool index out of bounds: %d", index)
	}

	return &cp[index-1], nil
}

type constantUtf8Info struct {
	str string
}

func newConstantUtf8Info(b []byte) *constantUtf8Info {
	return &constantUtf8Info{str: string(b)}
}

func (c constantUtf8Info) String() string {
	return fmt.Sprintf("<Utf8: \"%s\">", c.str)
}

type constantIntegerInfo struct {
	i int32
}

func newConstantIntegerInfo(b []byte) (*constantIntegerInfo, error) {
	info := &constantIntegerInfo{}

	if err := util.Decode(b, &info.i); err != nil {
		return nil, fmt.Errorf("could not decode integer constant: %w", err)
	}

	return info, nil
}

func (c constantIntegerInfo) String() string {
	return fmt.Sprintf("<Int: %d>", c.i)
}

type constantClassInfo struct {
	name *cpInfo
}

func newConstantClassInfo(b []byte, cp []cpInfo) (*constantClassInfo, error) {
	if name, err := resolveIndex(cp, b); err != nil {
		return nil, fmt.Errorf("could not resolve class info name index: %w", err)
	} else {
		return &constantClassInfo{name: name}, nil
	}
}

func (c constantClassInfo) String() string {
	return fmt.Sprintf("<Class: %s>", *c.name)
}

type constantNameAndTypeInfo struct {
	name       *cpInfo
	descriptor *cpInfo
}

func newConstantNameAndTypeInfo(bName []byte, bDescriptor []byte, cp []cpInfo) (*constantNameAndTypeInfo, error) {
	var name, descriptor *cpInfo

	if n, err := resolveIndex(cp, bName); err != nil {
		return nil, fmt.Errorf("could not resolve nameAndType info name index: %w", err)
	} else {
		name = n
	}

	if d, err := resolveIndex(cp, bDescriptor); err != nil {
		return nil, fmt.Errorf("could not resolve nameAndType info descriptor index: %w", err)
	} else {
		descriptor = d
	}

	return &constantNameAndTypeInfo{
		name:       name,
		descriptor: descriptor,
	}, nil
}

func (c constantNameAndTypeInfo) String() string {
	return fmt.Sprintf("<NameAndType: %s, %s>", *c.name, *c.descriptor)
}

type constantFieldrefInfo struct {
	class       *cpInfo
	nameAndType *cpInfo
}

func newConstantFieldrefInfo(bClass []byte, bNameAndType []byte, cp []cpInfo) (*constantFieldrefInfo, error) {
	var class, nameAndType *cpInfo

	if c, err := resolveIndex(cp, bClass); err != nil {
		return nil, fmt.Errorf("could not resolve fieldref info class index: %w", err)
	} else {
		class = c
	}

	if n, err := resolveIndex(cp, bNameAndType); err != nil {
		return nil, fmt.Errorf("could not resolve fieldref info nameAndType index: %w", err)
	} else {
		nameAndType = n
	}

	return &constantFieldrefInfo{
		class:       class,
		nameAndType: nameAndType,
	}, nil
}

func (c constantFieldrefInfo) String() string {
	return fmt.Sprintf("<Fieldref: %s, %s>", *c.class, *c.nameAndType)
}

type constantMethodrefInfo struct {
	class       *cpInfo
	nameAndType *cpInfo
}

func newConstantMethodrefInfo(bClass []byte, bNameAndType []byte, cp []cpInfo) (*constantMethodrefInfo, error) {
	var class, nameAndType *cpInfo

	if c, err := resolveIndex(cp, bClass); err != nil {
		return nil, fmt.Errorf("could not resolve methodref info class index: %w", err)
	} else {
		class = c
	}

	if n, err := resolveIndex(cp, bNameAndType); err != nil {
		return nil, fmt.Errorf("could not resolve methodref info nameAndType index: %w", err)
	} else {
		nameAndType = n
	}

	return &constantMethodrefInfo{
		class:       class,
		nameAndType: nameAndType,
	}, nil
}

func (c constantMethodrefInfo) String() string {
	return fmt.Sprintf("<Methodref: %s, %s>", *c.class, *c.nameAndType)
}
