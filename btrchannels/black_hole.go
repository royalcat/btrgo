package btrchannels

// BlackHole implements the InChannel interface and provides an analogue for the "Discard" variable in
// the ioutil package - it never blocks, and simply discards every value it reads. The number of items
// discarded in this way is counted and returned from Len.
type BlackHole[T any] struct {
	input  chan T
	length chan int
	count  int
}

func NewBlackHole[T any]() *BlackHole[T] {
	ch := &BlackHole[T]{
		input:  make(chan T),
		length: make(chan int),
	}
	go ch.discard()
	return ch
}

func (ch *BlackHole[T]) In() chan<- T {
	return ch.input
}

func (ch *BlackHole[T]) Len() int {
	val, open := <-ch.length
	if open {
		return val
	} else {
		return ch.count
	}
}

func (ch *BlackHole[T]) Cap() BufferCap {
	return Infinity
}

func (ch *BlackHole[T]) Close() {
	close(ch.input)
}

func (ch *BlackHole[T]) discard() {
	for {
		select {
		case _, open := <-ch.input:
			if !open {
				close(ch.length)
				return
			}
			ch.count++
		case ch.length <- ch.count:
		}
	}
}
