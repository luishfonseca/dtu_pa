package data

import (
	"fmt"
)

type CodeHandle struct {
	Begin int64
	End   int64
	baseData
}

func (CodeHandle) Tag() Tag {
	return CODE_HANDLE
}

func (c *CodeHandle) CodeHandle() *CodeHandle {
	return c
}

func (c CodeHandle) String() string {
	return fmt.Sprintf("<@ [%d.. %d] >", c.Begin, c.End)
}
