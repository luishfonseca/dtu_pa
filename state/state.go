package state

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type fallible interface{ Fail(error) }

type Fn[T fallible] func(T) Fn[T]

func Run[T fallible](s T, start Fn[T]) {
	for state := start; state != nil; {
		state = state(s)
	}
}

func Fail[T fallible](err error) Fn[T] {
	_, file, line, _ := runtime.Caller(1)
	return func(s T) Fn[T] {
		s.Fail(fmt.Errorf("%w (at %s:%d)", err, filepath.Base(file), line))
		return nil
	}
}
