package idemp

import (
	"context"
	"errors"
	"sync"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/redis"
)

var ErrAnotherLock = errors.New("idemp: waiting for another lock")
var ErrAnotherInstance = errors.New("idemp: waiting for another instance")

const DefaultTTL = 1 * 60 // 1 minutes

type TaskFunc func() (interface{}, error)
type ExecFunc func(taskKey string, timeout time.Duration, fn TaskFunc) (v interface{}, cached bool, err error, idempErr error)

type LockItem struct {
	SubKey string
}

// RegisGroup provides a cross-instance duplicate suppression.
type RedisGroup struct {
	prefix string
	rd     redis.Store
	g      *Group
	ttl    int
	sync.Mutex
}

func NewRedisGroup(rd redis.Store, prefix string, ttl int) *RedisGroup {
	if ttl == 0 {
		ttl = DefaultTTL
	}
	return &RedisGroup{
		prefix: prefix + ":",
		rd:     rd,
		g:      NewGroup(),
		ttl:    ttl,
	}
}

func (rg *RedisGroup) Shutdown() {
	rg.Lock()
	defer rg.Unlock()

	rg.g.Shutdown()
}

// DoAndWrap is short-hand for calling AcquireLock() and then exec. It also automatically
// removes the key if an error occurs.
func (rg *RedisGroup) DoAndWrap(
	ctx context.Context, key string, timeout time.Duration, msg string, fn TaskFunc,
) (v interface{}, cached bool, err error) {
	execFn, err := rg.AcquireLock(key, "")
	if err != nil {
		return nil, false, WrapError(ctx, err, msg)
	}

	v, cached, err, idempErr := execFn(key, timeout, fn)
	if idempErr != nil {
		return nil, false, WrapError(ctx, idempErr, msg)
	}
	if err != nil {
		rg.forget(key)
	}
	return v, cached, err
}

// DoAndWrapWithSubkey is short-hand for calling AcquireLock() and then exec
func (rg *RedisGroup) DoAndWrapWithSubkey(ctx context.Context, key string, subkey string, timeout time.Duration, fn TaskFunc, msg string) (v interface{}, cached bool, err error) {
	execFn, err := rg.AcquireLock(key, subkey)
	if err != nil {
		return nil, false, WrapError(ctx, err, msg)
	}

	v, cached, err, idempErr := execFn(key, timeout, fn)
	if idempErr != nil {
		return nil, false, WrapError(ctx, idempErr, msg)
	}
	if err != nil {
		rg.forget(key)
	}
	return v, cached, err
}

func (rg *RedisGroup) Acquire(groupKey, subkey string) (err error) {
	rg.Lock()
	defer rg.Unlock()

	storedKey, err := rg.get(groupKey)
	switch err {
	case nil:
		if storedKey != subkey {
			return ErrAnotherLock
		}
		return ErrAnotherInstance

	case redis.ErrNil:
		rg.set(groupKey, subkey)
		return nil

	default:
		return err
	}
}

func (rg *RedisGroup) AcquireLock(groupKey, subkey string) (exec ExecFunc, err error) {
	rg.Lock()
	storedKey, err := rg.get(groupKey)
	switch err {
	case nil:
		rg.Unlock()

		// key exists in Redis but subkey does not match
		if storedKey != subkey {
			// temporary set key and acquire lock
			return rg.setKeyAndAquireLock(groupKey, subkey)

			// TODO: handle lock in another instance
			// return nil, ErrAnotherLock
		}

		// key exists in Redis, subkey matches - now wait
		return rg.wait, nil

	case redis.ErrNil:
		// key does not exist in Redis, set key and acquire lock
		return rg.setKeyAndAquireLock(groupKey, subkey)

	default:
		rg.Unlock()
		return nil, err
	}
}

func (rg *RedisGroup) setKeyAndAquireLock(groupKey, subkey string) (exec ExecFunc, err error) {
	// key does not exist in Redis, set key and acquire lock
	rg.set(groupKey, subkey)
	rg.Unlock()
	return func(taskKey string, timeout time.Duration, fn TaskFunc) (v interface{}, cached bool, err error, idempErr error) {
		v, err, _ = rg.g.DoAndCleanup(taskKey, timeout, fn, func() {
			rg.ReleaseKey(groupKey, subkey)
		})
		return v, false, err, nil
	}, nil
}

func (rg *RedisGroup) wait(taskKey string, timeout time.Duration, fn TaskFunc) (v interface{}, cached bool, err error, idempErr error) {
	rg.g.Lock()

	// key exists in Redis, subkey matches
	if c, ok := rg.g.m[taskKey]; ok {
		// ...and the call execute in this instance
		c.dups++
		rg.g.Unlock()
		c.wg.Wait()
		return c.val, true, c.err, nil
	}

	// ...but the call does not execute in this instance
	rg.g.Unlock()
	return nil, false, nil, ErrAnotherInstance
}

func (rg *RedisGroup) set(key string, subkey string) {
	rg.rd.SetStringWithTTL(rg.prefix+key, subkey, rg.ttl)
}

func (rg *RedisGroup) get(key string) (string, error) {
	return rg.rd.GetString(rg.prefix + key)
}

func (rg *RedisGroup) ReleaseKey(groupKey, subkey string) {
	skey, _ := rg.get(groupKey)
	if skey == subkey {
		rg.rd.Del(rg.prefix + groupKey)
	}
}

func (rg *RedisGroup) forget(key string) {
	rg.g.Forget(key)
}

func WrapError(ctx context.Context, err error, msg string) error {
	switch err {
	case ErrAnotherLock:
		err = cm.Errorf(cm.FailedPrecondition, err, "Một người khác đang %v. Vui lòng chờ một lúc trước khi thử lại. Nếu cần thêm thông tin vui lòng liên hệ %v", msg, wl.X(ctx).CSEmail)
	case ErrAnotherInstance:
		err = cm.Errorf(cm.FailedPrecondition, err, "Thao tác %v đang được thực hiện. Vui lòng chờ một lúc trước khi thử lại. Nếu cần thêm thông tin vui lòng liên hệ %v.", msg, wl.X(ctx).CSEmail)
	}
	return err
}
