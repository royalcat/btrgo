package btrchannels

import "github.com/royalcat/btrgo/btrstruct"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type InfiniteChannel[T any] struct {
	input, output chan T
	length        chan int
	buffer        *btrstruct.Queue[T]
}

func NewInfiniteChannel[T any]() *InfiniteChannel[T] {
	ch := &InfiniteChannel[T]{
		input:  make(chan T),
		output: make(chan T),
		length: make(chan int),
		buffer: btrstruct.NewQueue[T](),
	}
	go ch.infiniteBuffer()
	return ch
}

func (ch *InfiniteChannel[T]) In() chan<- T {
	return ch.input
}

func (ch *InfiniteChannel[T]) Out() <-chan T {
	return ch.output
}

func (ch *InfiniteChannel[T]) Len() int {
	return <-ch.length
}

func (ch *InfiniteChannel[T]) Cap() BufferCap {
	return Infinity
}

func (ch *InfiniteChannel[T]) Close() {
	close(ch.input)
}

func (ch *InfiniteChannel[T]) infiniteBuffer() {
	var input, output chan T
	var next T
	input = ch.input

	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				ch.buffer.Add(elem)
			} else {
				input = nil
			}
		case output <- next:
			ch.buffer.Remove()
		case ch.length <- ch.buffer.Length():
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
			next = ch.buffer.Peek()
		} else {
			output = nil
			//next = nil
		}
	}

	close(ch.output)
	close(ch.length)
}
