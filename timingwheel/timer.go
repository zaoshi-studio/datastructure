package timingwheel

import (
	"container/list"
	"sync/atomic"
)

type timer struct {
	expiration int64
	task       func()

	b atomic.Pointer[bucket]

	element *list.Element
}

func (t *timer) getBucket() *bucket {
	return t.b.Load()
}

func (t *timer) setBucket(b *bucket) {
	t.b.Store(b)
}

func (t *timer) stop() bool {
	var stopped bool

	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stopped = b.Remove(t)
	}
	return stopped
}
