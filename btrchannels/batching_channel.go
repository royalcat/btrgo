package btrchannels

// BatchingChannel implements the Channel interface, with the change that instead of producing individual elements
// on Out(), it batches together the entire internal buffer each time. Trying to construct an unbuffered batching channel
// will panic, that configuration is not supported (and provides no benefit over an unbuffered NativeChannel).
type BatchingChannel[T any] struct {
	input  chan T
	output chan []T
	length chan int
	buffer []T
	size   BufferCap
}

func NewBatchingChannel[T any](size BufferCap) *BatchingChannel[T] {
	if size == None {
		panic("channels: BatchingChannel does not support unbuffered behaviour")
	}
	if size < 0 && size != Infinity {
		panic("channels: invalid negative size in NewBatchingChannel")
	}
	ch := &BatchingChannel[T]{
		input:  make(chan T),
		output: make(chan []T),
		length: make(chan int),
		size:   size,
	}
	go ch.batchingBuffer()
	return ch
}

func (ch *BatchingChannel[T]) In() chan<- T {
	return ch.input
}

// Out returns a <-chan interface{} in order that BatchingChannel conforms to the standard Channel interface provided
// by this package, however each output value is guaranteed to be of type []interface{} - a slice collecting the most
// recent batch of values sent on the In channel. The slice is guaranteed to not be empty or nil. In practice the net
// result is that you need an additional type assertion to access the underlying values.
func (ch *BatchingChannel[T]) Out() <-chan []T {
	return ch.output
}

func (ch *BatchingChannel[T]) Len() int {
	return <-ch.length
}

func (ch *BatchingChannel[T]) Cap() BufferCap {
	return ch.size
}

func (ch *BatchingChannel[T]) Close() {
	close(ch.input)
}

func (ch *BatchingChannel[T]) batchingBuffer() {
	var output chan []T
	var input, nextInput chan T
	nextInput = ch.input
	input = nextInput

	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				ch.buffer = append(ch.buffer, elem)
			} else {
				input = nil
				nextInput = nil
			}
		case output <- ch.buffer:
			ch.buffer = nil
		case ch.length <- len(ch.buffer):
		}

		if len(ch.buffer) == 0 {
			input = nextInput
			output = nil
		} else if ch.size != Infinity && len(ch.buffer) >= int(ch.size) {
			input = nil
			output = ch.output
		} else {
			input = nextInput
			output = ch.output
		}
	}

	close(ch.output)
	close(ch.length)
}
