package lexer

type stackedCounter struct {
	stack []uint16
}

func (sc *stackedCounter) empty() bool {
	return len(sc.stack) == 0
}

func (sc *stackedCounter) top() uint16 {
	return sc.stack[len(sc.stack)-1]
}

func (sc *stackedCounter) dec() {
	sc.stack[len(sc.stack)-1]--
}

func (sc *stackedCounter) push(n uint16) {
	sc.stack = append(sc.stack, n)
}

func (sc *stackedCounter) pop() {
	if sc.top() != 0 {
		panic("stackedCounter: pop called before counter reached zero")
	}

	sc.stack = sc.stack[:len(sc.stack)-1]
}
