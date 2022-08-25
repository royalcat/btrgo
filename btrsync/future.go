package btrsync

import "sync"

type Future[v any] struct {
	value v
	err   error

	executingMutex sync.Mutex
}

func RunFuture[v any](f func() (v, error)) *Future[v] {
	future := &Future[v]{}
	future.executingMutex.Lock()
	go func() {
		future.value, future.err = f()
		future.executingMutex.Unlock()
	}()

	return future
}

func (f *Future[v]) Wait() (v, error) {
	f.executingMutex.Lock()
	defer f.executingMutex.Unlock()
	return f.value, f.err
}
