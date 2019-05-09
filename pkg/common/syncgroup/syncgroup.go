package syncgroup

import (
	"context"
	"errors"
	"sync"
)

type Errors []error

func (e Errors) N() int {
	return len(e)
}

func (e Errors) NErrors() int {
	c := 0
	for _, err := range e {
		if err != nil {
			c++
		}
	}
	return c
}

func (e Errors) IsAll() bool {
	if len(e) == 0 {
		return false
	}
	for _, err := range e {
		if err == nil {
			return false
		}
	}
	return true
}

func (e Errors) All() error {
	for _, err := range e {
		if err == nil {
			return nil
		}
	}
	return e.Any()
}

func (e Errors) Any() error {
	var b []byte
	for _, err := range e {
		if err != nil {
			if b == nil {
				b = make([]byte, 0, 1024)
			} else {
				b = append(b, "; "...)
			}
			b = append(b, err.Error()...)
		}
	}
	if b == nil {
		return nil
	}
	return errors.New(string(b))
}

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

func (g *Group) Wait() Errors {
	g.wg.Wait()
	return g.errs
}
