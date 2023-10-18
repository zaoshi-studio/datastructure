package timingwheel

import (
	"context"
	"fmt"
	"github.com/zaoshi00/datastructure"
	"github.com/zaoshi00/datastructure/delayqueue"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type TimingWheel struct {
	tickMs      int64
	wheelSize   int64
	interval    int64
	currentTime int64

	buckets  []*bucket
	queue    *delayqueue.DelayQueue
	overflow unsafe.Pointer

	eg errgroup.Group
	wg sync.WaitGroup

	ctx    context.Context
	cancel func()
	ch     chan any
}

func New(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	startMs := datastructure.TimeToMs(time.Now().UTC())

	ctx, cancel := context.WithCancel(context.Background())

	return newTimingWheel(
		tickMs,
		wheelSize,
		startMs,
		delayqueue.New(int(wheelSize)),
		ctx,
		cancel,
		make(chan any, 1),
	)
}

func newTimingWheel(tickMs, wheelSize, startMs int64, queue *delayqueue.DelayQueue, ctx context.Context, cancel func(), ch chan any) *TimingWheel {
	buckets := make([]*bucket, wheelSize)
	for idx := range buckets {
		buckets[idx] = newBucket()
	}

	return &TimingWheel{
		tickMs:      tickMs,
		wheelSize:   wheelSize,
		currentTime: datastructure.Align(startMs, tickMs),
		interval:    tickMs * wheelSize,
		buckets:     buckets,
		queue:       queue,
		ctx:         ctx,
		cancel:      cancel,
		ch:          ch,
	}
}

func (tw *TimingWheel) start() {
	//tw.eg.Go(func() error {
	//
	//})

	tw.wg.Add(1)
	go func() {

		defer func() {
			tw.wg.Done()
			fmt.Println("take done")
		}()
		for {
			fmt.Println("take")
			//select {
			//case <-tw.ctx.Done():
			//	return
			//default:
			//	v := tw.queue.Take(tw.ctx)
			//	fmt.Println("take one")
			//	c := <-v
			//	tw.ch <- c.(*delayqueue.Element).Value
			//}

			v := tw.queue.Take(tw.ctx).Value
			if v == nil {
				return
			}
			tw.ch <- v
			fmt.Println("take one")
		}

	}()

	//tw.eg.Go(func() error {
	//})

	tw.wg.Add(1)
	go func() {
		defer func() {
			tw.wg.Done()
			fmt.Println("refresh done")
		}()
		for {
			select {
			case <-tw.ctx.Done():
				fmt.Println("refresh ctx done")
				return
			case v := <-tw.ch:
				if v == nil {
					return
				}
				b := v.(*bucket)
				tw.timeFix(b.getExpiration())
				b.refresh(tw.add)
			}
		}
	}()
}

func (tw *TimingWheel) stop() {
	tw.cancel()
	tw.ctx.Done()
	close(tw.ch)
	tw.wg.Wait()
	//err := tw.eg.Wait()
	//if err != nil {
	//	fmt.Println(err)
	//}

}

func (tw *TimingWheel) afterFunc(duration time.Duration, f func()) *timer {
	t := &timer{
		expiration: datastructure.TimeToMs(time.Now().Add(duration)),
		task:       f,
	}
	tw.add(t)
	return t
}

func (tw *TimingWheel) add(t *timer) {
	if !tw.addTimer(t) {
		go t.task()
	}
}

func (tw *TimingWheel) addTimer(t *timer) bool {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if t.expiration < currentTime+tw.tickMs {
		return false
	} else if t.expiration < currentTime+tw.interval {
		slot := t.expiration / tw.tickMs
		b := tw.buckets[slot%tw.wheelSize]
		b.addTimer(t)

		if b.setExpiration(slot * tw.tickMs) {
			tw.queue.Offer(delayqueue.NewElement(b, b.getExpiration()))
		}
		return true
	} else {
		overflow := atomic.LoadPointer(&tw.overflow)
		if overflow == nil {

			ctx, cancel := context.WithCancel(tw.ctx)
			atomic.CompareAndSwapPointer(
				&tw.overflow,
				nil,
				unsafe.Pointer(newTimingWheel(
					tw.tickMs,
					tw.wheelSize,
					currentTime,
					tw.queue,
					ctx,
					cancel,
					tw.ch,
				)),
			)
			overflow = atomic.LoadPointer(&tw.overflow)
		}

		return (*TimingWheel)(overflow).addTimer(t)
	}
}

func (tw *TimingWheel) timeFix(expiration int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= currentTime+tw.tickMs {
		currentTime = datastructure.Align(expiration, tw.tickMs)
		atomic.StoreInt64(&tw.currentTime, currentTime)

		overflow := atomic.LoadPointer(&tw.overflow)
		if overflow != nil {
			(*TimingWheel)(overflow).timeFix(expiration)
		}
	}
}
