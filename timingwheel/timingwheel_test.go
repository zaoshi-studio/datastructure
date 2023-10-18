package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	tw := New(time.Millisecond, 20)

	tw.start()
	defer tw.stop()
}
func TestTimingWheel_AfterFunc(t *testing.T) {
	tw := New(time.Millisecond, 20)
	tw.start()
	defer tw.stop()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		//10 * time.Millisecond,
		//50 * time.Millisecond,
		//100 * time.Millisecond,
		//500 * time.Millisecond,
		//1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			exitC := make(chan time.Time)

			start := time.Now().UTC()
			tw.afterFunc(d, func() {
				fmt.Println("execute")
				exitC <- time.Now().UTC()
			})

			got := (<-exitC).Truncate(time.Millisecond)
			min := start.Add(d).Truncate(time.Millisecond)

			err := 5 * time.Millisecond
			if got.Before(min) || got.After(min.Add(err)) {
				t.Errorf("Timer(%s) expiration: want [%s, %s], got %s", d, min, min.Add(err), got)
			}
		})
	}
}
