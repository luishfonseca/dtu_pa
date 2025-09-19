package lexer

type stackedCounter struct {
	stack []uint16
}

func (sc *stackedCounter) Top() uint16 {
	return sc.stack[len(sc.stack)-1]
}

func (sc *stackedCounter) Dec() {
	sc.stack[len(sc.stack)-1]--
}

func (sc *stackedCounter) Push(n uint16) {
	sc.stack = append(sc.stack, n)
}

func (sc *stackedCounter) Pop() {
	if sc.Top() != 0 {
		panic("stackedCounter: pop called before counter reached zero")
	}

	sc.stack = sc.stack[:len(sc.stack)-1]
}
