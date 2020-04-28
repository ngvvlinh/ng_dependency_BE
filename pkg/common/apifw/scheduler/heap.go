package scheduler

import (
	"sync"
	"time"

	"o.o/common/l"
)

type taskItem struct {
	id    interface{}
	task  TaskRunner
	next  time.Time
	m     sync.Mutex
	index int

	recurrentAfter time.Duration
}

func (t *taskItem) Recurrent(after time.Duration) ItemHandler {
	if after <= 0 {
		ll.Panic("Invalid", l.Duration("after", after))
	}
	t.recurrentAfter = after
	return t
}

type taskHeap []*taskItem

func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	return h[i].next.Before(h[j].next)
}

func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *taskHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*taskItem)
	item.index = n
	*h = append(*h, item)
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[0 : n-1]
	return item
}
