package priorityqueue

import (
	"container/heap"
	"github.com/zaoshi00/datastructure"
)

type PriorityQueue[integer datastructure.Integer] struct {
	elements []*Element[integer]
}

func New[integer datastructure.Integer](size int) *PriorityQueue[integer] {
	return &PriorityQueue[integer]{
		elements: make([]*Element[integer], 0, size),
	}
}

type Element[integer datastructure.Integer] struct {
	Value    any
	Priority integer
	Index    int
}

func NewElement[integer datastructure.Integer](value any, priority integer) *Element[integer] {
	return &Element[integer]{
		Value:    value,
		Priority: priority,
	}
}

func (pq PriorityQueue[integer]) Len() int {
	return len(pq.elements)
}

func (pq PriorityQueue[integer]) Less(i, j int) bool {
	return pq.elements[i].Priority < pq.elements[j].Priority
}

func (pq PriorityQueue[integer]) Swap(i, j int) {
	pq.elements[i], pq.elements[j] = pq.elements[j], pq.elements[i]
	pq.elements[i].Index = i
	pq.elements[j].Index = j
}

func (pq *PriorityQueue[integer]) Push(x any) {
	n := len(pq.elements)
	c := cap(pq.elements)
	if n+1 > c {
		npq := make([]*Element[integer], n, c*2)
		copy(npq, pq.elements)
		pq.elements = npq
	}
	pq.elements = pq.elements[0 : n+1]
	item := x.(*Element[integer])
	item.Index = n
	pq.elements[n] = item
}

func (pq *PriorityQueue[integer]) Pop() any {
	n := len(pq.elements)
	c := cap(pq.elements)
	if n < (c/2) || c > 25 {
		npq := make([]*Element[integer], n, c/2)
		copy(npq, pq.elements)
		pq.elements = npq
	}
	item := pq.elements[n-1]
	item.Index = -1
	pq.elements = pq.elements[:n-1]
	return item
}

func (pq *PriorityQueue[integer]) PeekAndShift(priority integer) (*Element[integer], integer) {
	if pq.Len() == 0 {
		return nil, integer(0)
	}
	elem := pq.elements[0]
	if elem.Priority > priority {
		return nil, elem.Priority - priority
	}
	heap.Remove(pq, 0)
	return elem, 0
}

func (pq *PriorityQueue[integer]) Offer(element *Element[integer]) {
	heap.Push(pq, element)
}

func (pq *PriorityQueue[integer]) Peek() *Element[integer] {
	if pq.Len() == 0 {
		return nil
	}
	elem := pq.elements[0]
	return elem
}

func (pq *PriorityQueue[integer]) Poll() *Element[integer] {
	return nil
}
