package btrchannels

import "testing"

func TestBlackHole(t *testing.T) {
	discard := NewBlackHole[int]()

	for i := 0; i < 1000; i++ {
		discard.In() <- i
	}

	discard.Close()

	if discard.Len() != 1000 {
		t.Error("blackhole expected 1000 was", discard.Len())
	}

	// no asserts here, this is just for the race detector's benefit
	ch := NewBlackHole[int]()
	go ch.Len()
	go ch.Cap()

	go func() {
		ch.In() <- 0
	}()
}
