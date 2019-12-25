package syncgroup

import (
	"context"
	"sync"

	"etop.vn/common/xerrors"
)

type Group struct {
	ctx  context.Context
	errs []error
	wg   sync.WaitGroup
	m    sync.Mutex
}

func New(n int) *Group {
	g := &Group{}
	if n > 0 {
		g.errs = make([]error, 0, n)
	}
	return g
}

func (g *Group) Go(fn func() error) {
	g.m.Lock()
	i := len(g.errs)
	g.errs = append(g.errs, nil)
	g.m.Unlock()

	g.wg.Add(1)
	go func(i int) {
		defer g.wg.Done()
		err := fn()
		if err != nil {
			g.m.Lock()
			g.errs[i] = err
			g.m.Unlock()
		}
	}(i)
}

func (g *Group) Wait() xerrors.Errors {
	g.wg.Wait()
	return g.errs
}
