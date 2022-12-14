package btrchannels

import (
	"math/rand"
	"testing"
)

func TestResizableChannel(t *testing.T) {
	var ch *ResizableChannel[int]

	ch = NewResizableChannel[int]()
	testChannel(t, "default resizable channel", ch)

	ch = NewResizableChannel[int]()
	testChannelPair(t, "default resizable channel", ch, ch)

	ch = NewResizableChannel[int]()
	ch.Resize(Infinity)
	testChannel(t, "infinite resizable channel", ch)

	ch = NewResizableChannel[int]()
	ch.Resize(Infinity)
	testChannelPair(t, "infinite resizable channel", ch, ch)

	ch = NewResizableChannel[int]()
	ch.Resize(5)
	testChannel(t, "5-buffer resizable channel", ch)

	ch = NewResizableChannel[int]()
	ch.Resize(5)
	testChannelPair(t, "5-buffer resizable channel", ch, ch)

	ch = NewResizableChannel[int]()
	testChannelConcurrentAccessors(t, "resizable channel", ch)
}

func TestResizableChannelOnline(t *testing.T) {
	stopper := make(chan bool)
	ch := NewResizableChannel[int]()
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		<-stopper
		ch.Close()
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			ch.Resize(BufferCap(rand.Intn(50) + 1))
		}
		close(stopper)
	}()

	for i := 0; i < 1000; i++ {
		val := <-ch.Out()
		if i != val {
			t.Fatal("resizable channel expected", i, "but got", val)
		}
	}
}
