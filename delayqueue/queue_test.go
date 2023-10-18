package delayqueue

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

}

func TestDelayQueuePutAndTake(t *testing.T) {
	queue := New(8)

	wg := sync.WaitGroup{}
	now := time.Now()
	element := NewElement(5, now.Add(time.Second*5).UnixMilli())
	wg.Add(1)
	go func() {

		queue.Take(context.Background())
		fmt.Println(time.Now().UnixMilli())

		wg.Done()
	}()

	queue.put(element)
	fmt.Println(now.UnixMilli())
	fmt.Println(now.Add(time.Second * 5).UnixMilli())
	wg.Wait()
}
