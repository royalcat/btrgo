package btrchannels

import (
	"github.com/royalcat/btrgo/btrstruct"
)

// RingChannel implements the Channel interface in a way that never blocks the writer.
// Specifically, if a value is written to a RingChannel when its buffer is full then the oldest
// value in the buffer is discarded to make room (just like a standard ring-buffer).
// Note that Go's scheduler can cause discarded values when they could be avoided, simply by scheduling
// the writer before the reader, so caveat emptor.
// For the opposite behaviour (discarding the newest element, not the oldest) see OverflowingChannel.
type RingChannel[T any] struct {
	input, output chan T
	length        chan int
	buffer        *btrstruct.Queue[T]
	size          BufferCap
}

func NewRingChannel[T any](size BufferCap) *RingChannel[T] {
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewRingChannel")
	}
	ch := &RingChannel[T]{
		input:  make(chan T),
		output: make(chan T),
		buffer: btrstruct.NewQueue[T](),
		size:   size,
	}
	if size == None {
		go ch.overflowingDirect()
	} else {
		ch.length = make(chan int)
		go ch.ringBuffer()
	}
	return ch
}

func (ch *RingChannel[T]) In() chan<- T {
	return ch.input
}

func (ch *RingChannel[T]) Out() <-chan T {
	return ch.output
}

func (ch *RingChannel[T]) Len() int {
	if ch.size == None {
		return 0
	} else {
		return <-ch.length
	}
}

func (ch *RingChannel[T]) Cap() BufferCap {
	return ch.size
}

func (ch *RingChannel[T]) Close() {
	close(ch.input)
}

// for entirely unbuffered cases
func (ch *RingChannel[T]) overflowingDirect() {
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
func (ch *RingChannel[T]) ringBuffer() {
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
					ch.buffer.Add(elem)
					if ch.size != Infinity && ch.buffer.Length() > int(ch.size) {
						ch.buffer.Remove()
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
			//next = nil
		}
	}

	close(ch.output)
	close(ch.length)
}
