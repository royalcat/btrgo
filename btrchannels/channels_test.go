package btrchannels

import (
	"math/rand"
	"testing"
	"time"
)

func testChannel(t *testing.T, name string, ch Channel[int]) {
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-ch.Out()
		if i != val {
			t.Fatal(name, "expected", i, "but got", val)
		}
	}
}

func testChannelPair(t *testing.T, name string, in InChannel[int], out OutChannel[int]) {
	go func() {
		for i := 0; i < 1000; i++ {
			in.In() <- i
		}
		in.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-out.Out()
		if i != val {
			t.Fatal("pair", name, "expected", i, "but got", val)
		}
	}
}

func testChannelConcurrentAccessors(t *testing.T, name string, ch Channel[int]) {
	// no asserts here, this is just for the race detector's benefit
	go ch.Len()
	go ch.Cap()

	go func() {
		ch.In() <- 0
	}()

	go func() {
		<-ch.Out()
	}()
}

func TestPipe(t *testing.T) {
	a := NewNativeChannel[int](None)
	b := NewNativeChannel[int](None)

	Pipe[int](a, b)

	testChannelPair(t, "pipe", a, b)
}

func TestWeakPipe(t *testing.T) {
	a := NewNativeChannel[int](None)
	b := NewNativeChannel[int](None)

	WeakPipe[int](a, b)

	testChannelPair(t, "pipe", a, b)
}

func testMultiplex(t *testing.T, multi func(output SimpleInChannel[int], inputs ...SimpleOutChannel[int])) {
	a := NewNativeChannel[int](None)
	b := NewNativeChannel[int](None)

	multi(b, a)

	testChannelPair(t, "simple multiplex", a, b)

	a = NewNativeChannel[int](None)
	inputs := []Channel[int]{
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
	}

	multi(a, inputs[0], inputs[1], inputs[2], inputs[3])

	go func() {
		rand.Seed(time.Now().Unix())
		for i := 0; i < 1000; i++ {
			inputs[rand.Intn(len(inputs))].In() <- i
		}
		for i := range inputs {
			inputs[i].Close()
		}
	}()
	for i := 0; i < 1000; i++ {
		val := <-a.Out()
		if i != val {
			t.Fatal("multiplexing expected", i, "but got", val)
		}
	}
}

func TestMultiplex(t *testing.T) {
	testMultiplex(t, Multiplex[int])
}

func TestWeakMultiplex(t *testing.T) {
	testMultiplex(t, WeakMultiplex[int])
}

func testTee(t *testing.T, tee func(input SimpleOutChannel[int], outputs ...SimpleInChannel[int])) {
	a := NewNativeChannel[int](None)
	b := NewNativeChannel[int](None)

	tee(a, b)

	testChannelPair(t, "simple tee", a, b)

	a = NewNativeChannel[int](None)
	outputs := []Channel[int]{
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
	}

	tee(a, outputs[0], outputs[1], outputs[2], outputs[3])

	go func() {
		for i := 0; i < 1000; i++ {
			a.In() <- i
		}
		a.Close()
	}()
	for i := 0; i < 1000; i++ {
		for _, output := range outputs {
			val := <-output.Out()
			if i != val {
				t.Fatal("teeing expected", i, "but got", val)
			}
		}
	}
}

func TestTee(t *testing.T) {
	testTee(t, Tee[int])
}

func TestWeakTee(t *testing.T) {
	testTee(t, WeakTee[int])
}

func testDistribute(t *testing.T, dist func(input SimpleOutChannel[int], outputs ...SimpleInChannel[int])) {
	a := NewNativeChannel[int](None)
	b := NewNativeChannel[int](None)

	dist(a, b)

	testChannelPair(t, "simple distribute", a, b)

	a = NewNativeChannel[int](None)
	outputs := []Channel[int]{
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
		NewNativeChannel[int](None),
	}

	dist(a, outputs[0], outputs[1], outputs[2], outputs[3])

	go func() {
		for i := 0; i < 1000; i++ {
			a.In() <- i
		}
		a.Close()
	}()

	received := make([]bool, 1000)
	for range received {
		var val interface{}
		select {
		case val = <-outputs[0].Out():
		case val = <-outputs[1].Out():
		case val = <-outputs[2].Out():
		case val = <-outputs[3].Out():
		}
		if received[val.(int)] {
			t.Fatal("distribute got value twice", val.(int))
		}
		received[val.(int)] = true
	}
	for i := range received {
		if !received[i] {
			t.Fatal("distribute missed", i)
		}
	}
}

func TestDistribute(t *testing.T) {
	testDistribute(t, Distribute[int])
}

func TestWeakDistribute(t *testing.T) {
	testDistribute(t, WeakDistribute[int])
}

func TestWrap(t *testing.T) {
	rawChan := make(chan int, 5)
	ch := Wrap(rawChan)

	for i := 0; i < 5; i++ {
		rawChan <- i
	}
	close(rawChan)

	for i := 0; i < 5; i++ {
		x := <-ch.Out()
		if x != i {
			t.Error("Wrapped value", x, "was expecting", i)
		}
	}
	_, ok := <-ch.Out()
	if ok {
		t.Error("Wrapped channel didn't close")
	}
}

func TestUnwrap(t *testing.T) {
	rawChan := make(chan int)
	ch := NewNativeChannel[int](5)
	Unwrap[int](ch, rawChan)

	for i := 0; i < 5; i++ {
		ch.In() <- i
	}
	ch.Close()

	for i := 0; i < 5; i++ {
		x := <-rawChan
		if x != i {
			t.Error("Unwrapped value", x, "was expecting", i)
		}
	}
	_, ok := <-rawChan
	if ok {
		t.Error("Unwrapped channel didn't close")
	}
}

func ExampleChannel() {
	var ch Channel[any]

	ch = NewInfiniteChannel[any]()

	for i := 0; i < 10; i++ {
		ch.In() <- nil
	}

	for i := 0; i < 10; i++ {
		<-ch.Out()
	}
}
