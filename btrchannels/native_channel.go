package btrchannels

// NativeInChannel implements the InChannel interface by wrapping a native go write-only channel.
type NativeInChannel[T any] chan<- T

func (ch NativeInChannel[T]) In() chan<- T {
	return ch
}

func (ch NativeInChannel[T]) Len() int {
	return len(ch)
}

func (ch NativeInChannel[T]) Cap() BufferCap {
	return BufferCap(cap(ch))
}

func (ch NativeInChannel[T]) Close() {
	close(ch)
}

// NativeOutChannel implements the OutChannel interface by wrapping a native go read-only channel.
type NativeOutChannel[T any] <-chan T

func (ch NativeOutChannel[T]) Out() <-chan T {
	return ch
}

func (ch NativeOutChannel[T]) Len() int {
	return len(ch)
}

func (ch NativeOutChannel[T]) Cap() BufferCap {
	return BufferCap(cap(ch))
}

// NativeChannel implements the Channel interface by wrapping a native go channel.
type NativeChannel[T any] chan T

// NewNativeChannel makes a new NativeChannel with the given buffer size. Just a convenience wrapper
// to avoid having to cast the result of make().
func NewNativeChannel[T any](size BufferCap) NativeChannel[T] {
	return make(chan T, size)
}

func (ch NativeChannel[T]) In() chan<- T {
	return ch
}

func (ch NativeChannel[T]) Out() <-chan T {
	return ch
}

func (ch NativeChannel[T]) Len() int {
	return len(ch)
}

func (ch NativeChannel[T]) Cap() BufferCap {
	return BufferCap(cap(ch))
}

func (ch NativeChannel[T]) Close() {
	close(ch)
}

// DeadChannel is a placeholder implementation of the Channel interface with no buffer
// that is never ready for reading or writing. Closing a dead channel is a no-op.
// Behaves almost like NativeChannel(nil) except that closing a nil NativeChannel will panic.
type DeadChannel[T any] struct{}

func NewDeadChannel[T any]() DeadChannel[T] {
	return DeadChannel[T]{}
}

func (ch DeadChannel[T]) In() chan<- T {
	return nil
}

func (ch DeadChannel[T]) Out() <-chan T {
	return nil
}

func (ch DeadChannel[T]) Len() int {
	return 0
}

func (ch DeadChannel[T]) Cap() BufferCap {
	return BufferCap(0)
}

func (ch DeadChannel[T]) Close() {
}
