package l

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelWatcher func(name string, level zapcore.Level)

type AtomicLevel struct {
	zap.AtomicLevel
	name  string
	watch map[int]LevelWatcher
	m     *sync.RWMutex
	id    *int // this is intentionally declared as pointer
}

func NewAtomicLevel(name string) AtomicLevel {
	id := 0
	return AtomicLevel{
		AtomicLevel: zap.NewAtomicLevel(),

		name:  name,
		watch: make(map[int]LevelWatcher),
		m:     new(sync.RWMutex),
		id:    &id,
	}
}

func (lvl AtomicLevel) SetLevel(level zapcore.Level) {
	lvl.m.RLock()
	defer lvl.m.RUnlock()

	lvl.AtomicLevel.SetLevel(level)
	for _, watch := range lvl.watch {
		watch(lvl.name, level)
	}
}

func (lvl AtomicLevel) Watch(fn LevelWatcher) (unwatch func()) {
	lvl.m.Lock()
	defer lvl.m.Unlock()
	if fn == nil {
		panic("nil function")
	}

	*lvl.id++                 // generate new id
	id := *lvl.id             // capture the value
	lvl.watch[id] = fn        // store the watch function
	fn(lvl.name, lvl.Level()) // call the watch function with current level
	return func() {
		lvl.unwatch(id)
	}
}

func (lvl AtomicLevel) unwatch(id int) {
	lvl.m.Lock()
	defer lvl.m.Unlock()
	delete(lvl.watch, id)
}
