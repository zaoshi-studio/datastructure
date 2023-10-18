package timingwheel

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type bucket struct {
	expiration atomic.Int64
	mu         sync.Mutex
	timers     *list.List
}

func newBucket() *bucket {
	b := &bucket{
		timers:     list.New(),
		expiration: atomic.Int64{},
	}

	b.expiration.Store(-1)

	return b
}

func (b *bucket) getExpiration() int64 {
	return b.expiration.Load()
}

func (b *bucket) setExpiration(expiration int64) bool {
	return b.expiration.Swap(expiration) != expiration
}

func (b *bucket) addTimer(t *timer) {
	b.mu.Lock()
	element := b.timers.PushBack(t)
	t.setBucket(b)
	t.element = element
	b.mu.Unlock()
}

func (b *bucket) removeTimer(t *timer) bool {
	if t.getBucket() != b {
		return false
	}
	b.timers.Remove(t.element)
	t.setBucket(nil)
	t.element = nil
	return true
}

func (b *bucket) Remove(t *timer) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.removeTimer(t)
}

func (b *bucket) refresh(f func(t *timer)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for e := b.timers.Front(); e != nil; {
		next := e.Next()
		t := e.Value.(*timer)
		b.removeTimer(t)
		f(t)
		e = next
	}

	b.setExpiration(-1)
}
