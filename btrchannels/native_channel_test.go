package btrchannels

import "testing"

func TestNativeChannels(t *testing.T) {
	var ch Channel[int]

	ch = NewNativeChannel[int](None)
	testChannel(t, "bufferless native channel", ch)

	ch = NewNativeChannel[int](None)
	testChannelPair(t, "bufferless native channel", ch, ch)

	ch = NewNativeChannel[int](5)
	testChannel(t, "5-buffer native channel", ch)

	ch = NewNativeChannel[int](5)
	testChannelPair(t, "5-buffer native channel", ch, ch)

	ch = NewNativeChannel[int](None)
	testChannelConcurrentAccessors(t, "native channel", ch)
}

func TestNativeInOutChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	Pipe[int](NativeOutChannel[int](ch1), NativeInChannel[int](ch2))
	NativeInChannel[int](ch1).Close()
}

func TestDeadChannel(t *testing.T) {
	ch := NewDeadChannel[int]()

	if ch.Len() != 0 {
		t.Error("dead channel length not 0")
	}
	if ch.Cap() != 0 {
		t.Error("dead channel cap not 0")
	}

	select {
	case <-ch.Out():
		t.Error("read from a dead channel")
	default:
	}

	select {
	case ch.In() <- 0:
		t.Error("wrote to a dead channel")
	default:
	}

	ch.Close()

	ch = NewDeadChannel[int]()
	testChannelConcurrentAccessors(t, "dead channel", ch)
}
