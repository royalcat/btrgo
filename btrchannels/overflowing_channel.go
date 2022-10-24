package btrchannels

import (
	"github.com/royalcat/btrgo/btrstruct"
)

// OverflowingChannel implements the Channel interface in a way that never blocks the writer.
// Specifically, if a value is written to an OverflowingChannel when its buffer is full
// (or, in an unbuffered case, when the recipient is not ready) then that value is simply discarded.
// Note that Go's scheduler can cause discarded values when they could be avoided, simply by scheduling
// the writer before the reader, so caveat emptor.
// For the opposite behaviour (discarding the oldest element, not the newest) see RingChannel.
type OverflowingChannel[T any] struct {
	input, output chan T
	length        chan int
	buffer        *btrstruct.Queue[T]
	size          BufferCap
}

func NewOverflowingChannel[T any](size BufferCap) *OverflowingChannel[T] {
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewOverflowingChannel")
	}
	ch := &OverflowingChannel[T]{
		input:  make(chan T),
		output: make(chan T),
		length: make(chan int),
		size:   size,
	}
	if size == None {
		go ch.overflowingDirect()
	} else {
		ch.buffer = btrstruct.NewQueue[T]()
		go ch.overflowingBuffer()
	}
	return ch
}

func (ch *OverflowingChannel[T]) In() chan<- T {
	return ch.input
}

func (ch *OverflowingChannel[T]) Out() <-chan T {
	return ch.output
}

func (ch *OverflowingChannel[T]) Len() int {
	if ch.size == None {
		return 0
	} else {
		return <-ch.length
	}
}

func (ch *OverflowingChannel[T]) Cap() BufferCap {
	return ch.size
}

func (ch *OverflowingChannel[T]) Close() {
	close(ch.input)
}

// for entirely unbuffered cases
func (ch *OverflowingChannel[T]) overflowingDirect() {
	for elem := range ch.input {
		// if we can't write it immediately, drop it and move on
		select {
		case ch.output <- elem:
		default:
		}
	}
	close(ch.output)
}

// for all buffered cases
func (ch *OverflowingChannel[T]) overflowingBuffer() {
	var input, output chan T
	var next T
	input = ch.input

	for input != nil || output != nil {
		select {
		// Prefer to write if possible, which is surprisingly effective in reducing
		// dropped elements due to overflow. The naive read/write select chooses randomly
		// when both channels are ready, which produces unnecessary drops 50% of the time.
		case output <- next:
			ch.buffer.Remove()
		default:
			select {
			case elem, open := <-input:
				if open {
					if ch.size == Infinity || ch.buffer.Length() < int(ch.size) {
						ch.buffer.Add(elem)
					}
				} else {
					input = nil
				}
			case output <- next:
				ch.buffer.Remove()
			case ch.length <- ch.buffer.Length():
			}
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
			next = ch.buffer.Peek()
		} else {
			output = nil
			//FIXME next = nil
		}
	}

	close(ch.output)
	close(ch.length)
}
