package priorityqueue

import "testing"

func TestNew(t *testing.T) {
	pq := New[int64](8)

	for i := 1; i <= 8; i++ {
		pq.Push(i)
	}
}
