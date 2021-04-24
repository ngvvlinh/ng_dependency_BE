package webhook

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
	cm "o.o/backend/pkg/common"
)

type Func func(context.Context, interface{}) error

type JobKeeper struct {
	mChan map[string]chan bool
	mu    sync.Mutex
}

func NewJobKeeper() *JobKeeper {
	return &JobKeeper{
		mChan: make(map[string]chan bool),
	}
}

func (jk *JobKeeper) AddJob(id string, jobFunc Func, args interface{}) error {
	jk.mu.Lock()
	defer jk.mu.Unlock()

	if _, ok := jk.mChan[id]; ok {
		return cm.Errorf(cm.InvalidArgument, nil, "id exists")
	}
	ll.Info(fmt.Sprintf("Start job with id (%s)", id))

	ch := make(chan bool)
	jk.mChan[id] = ch

	go func(_ch chan bool, _func Func, _args interface{}) {
		ctx, ctxCancel := context.WithCancel(context.Background())

		if err := _func(ctx, _args); err != nil {
			return
		}

		for {
			select {
			case <-_ch:
				ctxCancel()
				return
			default:
				// no-op
			}
		}
	}(ch, jobFunc, args)

	return nil
}

func (jk *JobKeeper) StopJob(id string) error {
	jk.mu.Lock()
	defer jk.mu.Unlock()
	ch, ok := jk.mChan[id]
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "id doesn't exist")
	}
	ll.Info(fmt.Sprintf("Stop job with id (%s)", id))

	ch <- true

	delete(jk.mChan, id)

	return nil
}
