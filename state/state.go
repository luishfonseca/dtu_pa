package state

type fallible interface{ Fail(error) }

type Fn[T fallible] func(T) Fn[T]

func Run[T fallible](s T, start Fn[T]) {
	for state := start; state != nil; {
		state = state(s)
	}
}
