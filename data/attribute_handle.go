package data

import (
	"fmt"
)

type AttributeHandle struct {
	AttributeTag Tag
	Begin        int64
	baseData
}

func (*AttributeHandle) Tag() Tag {
	return ATTRIBUTE_HANDLE
}

func (a *AttributeHandle) AttributeHandle() *AttributeHandle {
	return a
}

func (a AttributeHandle) String() string {
	return fmt.Sprintf("<%s @ %d>", a.AttributeTag, a.Begin)
}
