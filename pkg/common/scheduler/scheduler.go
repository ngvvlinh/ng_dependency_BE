package scheduler

import (
	"container/heap"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/l"
)

var ll = l.New()

type TaskRunner func(id interface{}, p Planner) error

type Planner interface {
	Add(id interface{}, nextTime time.Time, task TaskRunner) ItemHandler
	AddAfter(id interface{}, after time.Duration, task TaskRunner) ItemHandler
}

type ItemHandler interface {
	Recurrent(after time.Duration) ItemHandler
}

type Scheduler struct {
	numWorkers int

	tasks map[interface{}]*taskItem
	heap  *taskHeap
	stop  chan struct{}
	ch    chan *taskItem
	m     sync.Mutex
	wg    sync.WaitGroup

	isRunning   bool
	actuallyRun bool
}

func New(numWorkers int) *Scheduler {
	if numWorkers <= 0 {
		ll.Panic("Invalid num workers")
	}
	return &Scheduler{
		numWorkers: numWorkers,

		tasks: make(map[interface{}]*taskItem),
		heap:  new(taskHeap),
		stop:  make(chan struct{}),
		ch:    make(chan *taskItem, numWorkers),
	}
}

func (s *Scheduler) AddAfter(id interface{}, after time.Duration, task TaskRunner) ItemHandler {
	return s.Add(id, time.Now().Add(after), task)
}

func (s *Scheduler) Add(id interface{}, nextTime time.Time, task TaskRunner) ItemHandler {
	if nextTime.IsZero() {
		ll.S.Panicf("Task %v: Zero time", id)
	}

	s.m.Lock()
	defer s.m.Unlock()

	var item *taskItem
	if item = s.tasks[id]; item != nil {
		item.next = nextTime
		heap.Fix(s.heap, item.index)

	} else {
		item = &taskItem{
			id:   id,
			task: task,
			next: nextTime,
		}
		heap.Push(s.heap, item)
		s.tasks[item.id] = item
	}

	s.tryStart()
	return item
}

func (s *Scheduler) Start() {
	s.m.Lock()
	defer s.m.Unlock()

	if s.isRunning {
		ll.Panic("scheduler: Already running")
	}

	// Start workers
	s.wg.Add(s.numWorkers)
	for i := 0; i < s.numWorkers; i++ {
		go s.runWorker(i)
	}

	s.isRunning = true
	s.tryStart()
}

func (s *Scheduler) Stop() {
	close(s.stop)
	s.wg.Wait()
	s.isRunning = false

	ll.Info("All workers stopped", l.Int("Remaining tasks", len(*s.heap)))
}

func (s *Scheduler) push(item *taskItem) {
	s.m.Lock()
	defer s.m.Unlock()

	heap.Push(s.heap, item)
	s.tasks[item.id] = item
}

func (s *Scheduler) pop() *taskItem {
	s.m.Lock()
	defer s.m.Unlock()

	item := heap.Pop(s.heap).(*taskItem)
	delete(s.tasks, item.id)
	return item
}

func (s *Scheduler) Peek(id interface{}) (time.Time, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	item := s.tasks[id]
	if item == nil {
		return time.Time{}, false
	}
	return item.next, true
}

func (s *Scheduler) Remove(id interface{}) bool {
	s.m.Lock()
	defer s.m.Unlock()

	item := s.tasks[id]
	if item == nil {
		return false
	}
	heap.Remove(s.heap, item.index)
	delete(s.tasks, item.id)
	return true
}

func (s *Scheduler) tryStart() {
	// Start producer, it only actually runs if there is at least one job.
	if s.isRunning && len(*s.heap) > 0 && !s.actuallyRun {
		s.actuallyRun = true
		go s.start()
	}
}

func (s *Scheduler) start() {
	ll.Info("scheduler: Actually start")

	for {
		s.m.Lock()
		if len(*s.heap) == 0 {
			s.actuallyRun = false
			s.m.Unlock()
			ll.Info("scheduler: No item remain, stop")
			return
		}

		now := time.Now()
		nextItem := (*s.heap)[0]
		delta := nextItem.next.Sub(now)
		if nextItem.recurrentAfter > 0 {
			nextItem.next = now.Add(nextItem.recurrentAfter)
			heap.Fix(s.heap, nextItem.index)
		} else {
			heap.Pop(s.heap)
			delete(s.tasks, nextItem.id)
		}
		s.m.Unlock()

		// Wait until the job ready
		time.Sleep(delta)
		s.ch <- nextItem
	}
}

func (s *Scheduler) runWorker(id int) {
	ll.S.Infof("Worker %v started", id)
	defer func() {
		ll.S.Infof("Worker %v stopped", id)
		s.wg.Done()
	}()

	for {
		select {
		case <-s.stop:
			return

		case item := <-s.ch:
			s.execTask(item)
		}
	}
}

func (s *Scheduler) execTask(item *taskItem) {
	t := time.Now()
	defer func() {
		err := recover()
		if err != nil {
			d := time.Since(t)
			ll.Error("Task panic", l.Any("id", item.id), l.Duration("d", d), l.Any("err", err), l.Stack())
		}
	}()

	err := item.task(item.id, s)
	d := time.Since(t)
	if err != nil {
		ll.Error("Task error", l.Any("id", item.id), l.Duration("d", d), l.Error(err))
	} else {
		ll.Info("Task done ", l.Any("id", item.id), l.Duration("d", d))
	}
}
