package btrchannels

import "testing"

func testBatches(t *testing.T, chIn SimpleInChannel[int], chOut SimpleOutChannel[[]int]) {
	go func() {
		for i := 0; i < 1000; i++ {
			chIn.In() <- i
		}
		chIn.Close()
	}()

	i := 0
	for val := range chOut.Out() {
		for _, elem := range val {
			if i != elem {
				t.Fatal("batching channel expected", i, "but got", elem)
			}
			i++
		}
	}
}

func TestBatchingChannel(t *testing.T) {
	ch := NewBatchingChannel[int](Infinity)
	testBatches(t, ch, ch)

	ch = NewBatchingChannel[int](2)
	testBatches(t, ch, ch)

	// FIXME
	// ch = NewBatchingChannel[int](1)
	// testChannelConcurrentAccessors(t, "batching channel", ch)
}

func TestBatchingChannelCap(t *testing.T) {
	ch := NewBatchingChannel[int](Infinity)
	if ch.Cap() != Infinity {
		t.Error("incorrect capacity on infinite channel")
	}

	ch = NewBatchingChannel[int](5)
	if ch.Cap() != 5 {
		t.Error("incorrect capacity on infinite channel")
	}
}
