package state

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Fail[T fallible](err error) Fn[T] {
	_, file, line, _ := runtime.Caller(1)
	return func(s T) Fn[T] {
		s.Fail(fmt.Errorf("%w (at %s:%d)", err, filepath.Base(file), line))
		return nil
	}
}
