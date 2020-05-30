package lifecycle

import (
	"context"
	"sync"

	"o.o/common/l"
)

type Shutdowner interface {
	context.Context
	Register(...func())
	Wait()
}

type shutdowner struct {
	context.Context
	once      sync.Once
	funcs     []func()
	ctxCancel func()
}

func WithCancel(ctx context.Context) (Shutdowner, func()) {
	newCtx, cancel := context.WithCancel(ctx)
	s := &shutdowner{Context: newCtx, ctxCancel: cancel}
	return s, s.ctxCancel
}

func (s *shutdowner) Register(funcs ...func()) {
	select {
	case <-s.Context.Done():
		for _, fn := range funcs {
			fn()
		}
	default:
		s.funcs = append(s.funcs, funcs...)
	}
}

func (s *shutdowner) shutdownAll() {
	s.once.Do(func() {
		s.ctxCancel()
		for _, fn := range s.funcs {
			fn()
		}
	})
}

func (s *shutdowner) Wait() {
	r := recover()
	if r != nil {
		s.shutdownAll()
		ll.Panic("panic while waiting", l.Any("err", r)) // rethrow
	}
	<-s.Context.Done()
	s.shutdownAll()
}
