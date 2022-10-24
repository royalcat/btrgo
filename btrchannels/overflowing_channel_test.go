package btrchannels

import "testing"

func TestOverflowingChannel(t *testing.T) {
	var ch Channel[int]

	ch = NewOverflowingChannel[int](Infinity) // yes this is rather silly, but it should work
	testChannel(t, "infinite overflowing channel", ch)

	ch = NewOverflowingChannel[int](None)
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()
	prev := -1
	for i := range ch.Out() {
		if prev >= i {
			t.Fatal("overflowing channel prev", prev, "but got", i)
		}
	}

	ch = NewOverflowingChannel[int](10)
	for i := 0; i < 1000; i++ {
		ch.In() <- i
	}
	ch.Close()
	for i := 0; i < 10; i++ {
		val := <-ch.Out()
		if i != val {
			t.Fatal("overflowing channel expected", i, "but got", val)
		}
	}
	if val, open := <-ch.Out(); open == true {
		t.Fatal("overflowing channel expected closed but got", val)
	}

	ch = NewOverflowingChannel[int](None)
	ch.In() <- 0
	ch.Close()
	if val, open := <-ch.Out(); open == true {
		t.Fatal("overflowing channel expected closed but got", val)
	}

	ch = NewOverflowingChannel[int](2)
	testChannelConcurrentAccessors(t, "overflowing channel", ch)
}
