package btrchannels

import "testing"

func TestInfiniteChannel(t *testing.T) {
	var ch Channel[int]

	ch = NewInfiniteChannel[int]()
	testChannel(t, "infinite channel", ch)

	ch = NewInfiniteChannel[int]()
	testChannelPair(t, "infinite channel", ch, ch)

	ch = NewInfiniteChannel[int]()
	testChannelConcurrentAccessors(t, "infinite channel", ch)
}

func BenchmarkInfiniteChannelSerial(b *testing.B) {
	ch := NewInfiniteChannel[any]()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
	}
	for i := 0; i < b.N; i++ {
		<-ch.Out()
	}
}

func BenchmarkInfiniteChannelParallel(b *testing.B) {
	ch := NewInfiniteChannel[any]()
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch.Out()
		}
		ch.Close()
	}()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
	}
	<-ch.Out()
}

func BenchmarkInfiniteChannelTickTock(b *testing.B) {
	ch := NewInfiniteChannel[any]()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
		<-ch.Out()
	}
}
