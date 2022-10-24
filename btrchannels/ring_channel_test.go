package btrchannels

import "testing"

func TestRingChannel(t *testing.T) {
	var ch Channel[int]

	ch = NewRingChannel[int](Infinity) // yes this is rather silly, but it should work
	testChannel(t, "infinite ring-buffer channel", ch)

	ch = NewRingChannel[int](None)
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()
	prev := -1
	for i := range ch.Out() {
		if prev >= i {
			t.Fatal("ring channel prev", prev, "but got", i)
		}
	}

	ch = NewRingChannel[int](10)
	for i := 0; i < 1000; i++ {
		ch.In() <- i
	}
	ch.Close()
	for i := 990; i < 1000; i++ {
		val := <-ch.Out()
		if i != val {
			t.Fatal("ring channel expected", i, "but got", val)
		}
	}
	if val, open := <-ch.Out(); open == true {
		t.Fatal("ring channel expected closed but got", val)
	}

	ch = NewRingChannel[int](None)
	ch.In() <- 0
	ch.Close()
	if val, open := <-ch.Out(); open == true {
		t.Fatal("ring channel expected closed but got", val)
	}

	ch = NewRingChannel[int](2)
	testChannelConcurrentAccessors(t, "ring channel", ch)
}
