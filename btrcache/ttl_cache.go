package btrcache

import (
	"context"
	"sync"
	"time"
)

type TTLValue[V any] struct {
	m              sync.Mutex
	ttl            time.Duration
	updated        time.Time
	value          *V
	updateCallback func() V
}

func NewTTLValue[V any](ttl time.Duration, updateCallback func() V) *TTLValue[V] {
	return &TTLValue[V]{ttl: ttl, updateCallback: updateCallback}
}

func (self *TTLValue[V]) Get() V {
	self.m.Lock()
	defer self.m.Unlock()

	now := time.Now()
	if self.value == nil || now.After(self.updated.Add(self.ttl)) {
		newVal := self.updateCallback()
		self.value = &newVal
		self.updated = now
	}

	return *self.value
}

func (self *TTLValue[V]) Expire() {
	self.m.Lock()
	defer self.m.Unlock()

	self.value = nil
}

type TTLValueAsync[V any] struct {
	m              sync.Mutex
	ttl            time.Duration
	wg             sync.WaitGroup
	value          *V
	updateCallback func() V
}

func NewTTLValueAsync[V any](ctx context.Context, ttl time.Duration, updateCallback func() V) *TTLValueAsync[V] {
	val := &TTLValueAsync[V]{ttl: ttl, updateCallback: updateCallback}
	go func() {
		tick := time.NewTicker(ttl)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				val.wg.Add(1)

				// TODO what if the callback takes a long time?
				newVal := val.updateCallback()

				val.m.Lock()
				val.value = &newVal
				val.m.Unlock()

				val.wg.Done()
			}
		}
	}()

	return val
}

func (self *TTLValueAsync[V]) Get() V {
	if self.value == nil {
		self.wg.Wait()
	}

	self.m.Lock()
	defer self.m.Unlock()
	return *self.value
}

func (self *TTLValueAsync[V]) Expire() {
	self.m.Lock()
	defer self.m.Unlock()

	self.value = nil
}
