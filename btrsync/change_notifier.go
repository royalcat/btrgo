package btrsync

import (
	"sync"

	"github.com/royalcat/btrgo"
)

type ChangeNotifier[T any] struct {
	mutex   sync.Mutex
	closed  bool
	last    T
	updates chan T

	subscribers []chan T
	handlers    []func(T)
}

func (p *ChangeNotifier[T]) Last() T {
	return p.last
}

func (p *ChangeNotifier[T]) Subscribe() <-chan T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return nil
	}

	c := make(chan T, 0)
	p.subscribers = append(p.subscribers, c)

	return c

}

func (p *ChangeNotifier[T]) AddHandler(h ...func(T)) {
	if p.closed {
		return
	}

	p.handlers = append(p.handlers, h...)
}

func (p *ChangeNotifier[T]) Update(updates ...T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return
	}

	for _, update := range updates {
		for _, sub := range p.subscribers {
			if len(sub) == 0 {
				sub <- update
			} else {
				btrgo.Drain(sub)
				sub <- update
			}
		}

		for _, handler := range p.handlers {
			handler(update)
		}
	}

	p.last = updates[len(updates)-1]

}

func (p *ChangeNotifier[T]) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return
	}

	for _, sub := range p.subscribers {
		if len(sub) == 1 { // give a chance to slower readers to read
			btrgo.Drain(sub)
			close(sub)
		}

		// close everything
		//drain the channel
		for range sub {
		}
		close(sub)
	}

	p.handlers = []func(T){}
}
