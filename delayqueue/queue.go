package delayqueue

import (
	"context"
	"github.com/zaoshi00/datastructure/priorityqueue"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"time"
)

type Element struct {
	Value      any
	Callback   func(element *Element)
	Expiration int64
}

func NewElement(value any, expiration int64) *Element {
	return &Element{
		Value:      value,
		Expiration: expiration,
	}
}

type DelayQueue struct {
	pq *priorityqueue.PriorityQueue[int64]
	mu sync.Mutex

	running atomic.Int32
	wakeup  chan struct{}
	eg      errgroup.Group
}

func New(size int) *DelayQueue {
	return &DelayQueue{
		pq: priorityqueue.New[int64](size),
	}
}

func (dq *DelayQueue) notify() {

}

func (dq *DelayQueue) add(element *Element) bool {
	return dq.Offer(element)
}

func (dq *DelayQueue) Offer(element *Element) bool {
	pqElement := priorityqueue.NewElement[int64](element, element.Expiration)
	dq.mu.Lock()
	dq.pq.Offer(pqElement)
	if pqElement.Index == 0 {
		dq.notify()
	}
	dq.mu.Unlock()
	return true
}

func (dq *DelayQueue) put(element *Element) {
	dq.Offer(element)
}

func (dq *DelayQueue) poll() *Element {
	dq.mu.Lock()
	pqElement := dq.pq.Peek()
	if pqElement == nil || pqElement.Priority > 0 {
		return nil
	}
	return dq.pq.Poll().Value.(*Element)
}

func (dq *DelayQueue) Take(ctx context.Context) *Element {

	for {
		dq.mu.Lock()
		pqElement, t := dq.pq.PeekAndShift(time.Now().UnixMilli())
		if pqElement != nil {
			defer dq.mu.Unlock()
			return pqElement.Value.(*Element)
		}
		if t <= 0 {
			dq.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil
			}
			continue
		} else {
			select {
			case <-time.After(time.Duration(t) * time.Millisecond):
			}
			dq.mu.Unlock()
			continue
		}
	}
}

func (dq *DelayQueue) peek() *Element {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	return dq.pq.Peek().Value.(*Element)
}

func (dq *DelayQueue) size() int {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	return dq.pq.Len()
}
