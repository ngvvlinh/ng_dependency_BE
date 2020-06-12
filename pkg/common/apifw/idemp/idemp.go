// Package idemp is the same as singleflight, except that it caches the result.
package idemp

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// call is an in-flight or completed singleflight.Do call
type call struct {
	wg sync.WaitGroup

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val interface{}
	err error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int
	chans []chan<- Result

	cleanup func()
}

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group struct {
	sync.Mutex // protects m

	m map[string]*call
}

// Result holds the results of Do, so they can be passed
// on a channel.
type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}

// NewGroup ...
func NewGroup() *Group {
	return &Group{m: make(map[string]*call)}
}

func (g *Group) Do(key string, timeout time.Duration, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return g.DoAndCleanup(key, timeout, fn, nil)
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (g *Group) DoAndCleanup(key string, timeout time.Duration, fn func() (interface{}, error), cleanup func()) (v interface{}, err error, shared bool) {
	g.Lock()
	if c, ok := g.m[key]; ok {
		c.dups++
		g.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := &call{cleanup: cleanup}
	c.wg.Add(1)
	g.m[key] = c
	g.Unlock()

	g.doCall(c, key, timeout, fn)
	return c.val, c.err, c.dups > 0
}

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
func (g *Group) DoChan(key string, timeout time.Duration, fn func() (interface{}, error), cleanup func()) <-chan Result {
	ch := make(chan Result, 1)
	g.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}, cleanup: cleanup}
	c.wg.Add(1)
	g.m[key] = c
	g.Unlock()

	go g.doCall(c, key, timeout, fn)

	return ch
}

// doCall handles the single call for a key.
func (g *Group) doCall(c *call, key string, timeout time.Duration, fn func() (_ interface{}, _err error)) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("idemp: panic (recovered)", e)
			debug.PrintStack()
			c.err = fmt.Errorf("%v", e)
		}

		c.wg.Done()
		for _, ch := range c.chans {
			ch <- Result{c.val, c.err, c.dups > 0}
		}

		if timeout >= 0 {
			time.AfterFunc(timeout, func() { g.Forget(key) })
		} else {
			g.Forget(key)
		}
	}()
	c.val, c.err = fn()
}

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (g *Group) Forget(key string) {
	g.Lock()
	defer g.Unlock()

	if c, ok := g.m[key]; ok {
		delete(g.m, key)
		if c.cleanup != nil {
			c.cleanup()
		}
	}
}

func (g *Group) Shutdown() {
	g.Lock()
	defer g.Unlock()

	for _, c := range g.m {
		if c.cleanup != nil {
			c.cleanup()
			c.cleanup = nil
		}
	}
}
