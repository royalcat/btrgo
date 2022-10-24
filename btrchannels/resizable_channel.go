package btrchannels

import (
	"github.com/royalcat/btrgo/btrstruct"
)

// ResizableChannel implements the Channel interface with a resizable buffer between the input and the output.
// The channel initially has a buffer size of 1, but can be resized by calling Resize().
//
// Resizing to a buffer capacity of None is, unfortunately, not supported and will panic
// (see https://github.com/eapache/channels/issues/1).
// Resizing back and forth between a finite and infinite buffer is fully supported.
type ResizableChannel[T any] struct {
	input, output    chan T
	length           chan int
	capacity, resize chan BufferCap
	size             BufferCap
	buffer           *btrstruct.Queue[T]
}

func NewResizableChannel[T any]() *ResizableChannel[T] {
	ch := &ResizableChannel[T]{
		input:    make(chan T),
		output:   make(chan T),
		length:   make(chan int),
		capacity: make(chan BufferCap),
		resize:   make(chan BufferCap),
		size:     1,
		buffer:   btrstruct.NewQueue[T](),
	}
	go ch.magicBuffer()
	return ch
}

func (ch *ResizableChannel[T]) In() chan<- T {
	return ch.input
}

func (ch *ResizableChannel[T]) Out() <-chan T {
	return ch.output
}

func (ch *ResizableChannel[T]) Len() int {
	return <-ch.length
}

func (ch *ResizableChannel[T]) Cap() BufferCap {
	val, open := <-ch.capacity
	if open {
		return val
	} else {
		return ch.size
	}
}

func (ch *ResizableChannel[T]) Close() {
	close(ch.input)
}

func (ch *ResizableChannel[T]) Resize(newSize BufferCap) {
	if newSize == None {
		panic("channels: ResizableChannel does not support unbuffered behaviour")
	}
	if newSize < 0 && newSize != Infinity {
		panic("channels: invalid negative size trying to resize channel")
	}
	ch.resize <- newSize
}

func (ch *ResizableChannel[T]) magicBuffer() {
	var input, output, nextInput chan T
	var next T
	nextInput = ch.input
	input = nextInput

	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				ch.buffer.Add(elem)
			} else {
				input = nil
				nextInput = nil
			}
		case output <- next:
			ch.buffer.Remove()
		case ch.size = <-ch.resize:
		case ch.length <- ch.buffer.Length():
		case ch.capacity <- ch.size:
		}

		if ch.buffer.Length() == 0 {
			output = nil
			//next = nil
		} else {
			output = ch.output
			next = ch.buffer.Peek()
		}

		if ch.size != Infinity && ch.buffer.Length() >= int(ch.size) {
			input = nil
		} else {
			input = nextInput
		}
	}

	close(ch.output)
	close(ch.resize)
	close(ch.length)
	close(ch.capacity)
}
