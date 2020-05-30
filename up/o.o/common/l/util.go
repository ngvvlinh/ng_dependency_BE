package l

import (
	"fmt"
	"time"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

func Any(key string, val interface{}) zap.Field        { return zap.Any(key, val) }
func Bool(key string, val bool) zap.Field              { return zap.Bool(key, val) }
func Duration(key string, val time.Duration) zap.Field { return zap.Duration(key, val) }
func Float64(key string, val float64) zap.Field        { return zap.Float64(key, val) }
func Int(key string, val int) zap.Field                { return zap.Int(key, val) }
func Int32(key string, val int32) zap.Field            { return zap.Int(key, int(val)) }
func Int64(key string, val int64) zap.Field            { return zap.Int64(key, val) }
func Object(key string, val interface{}) zap.Field     { return zap.Stringer(key, Dump(val)) }
func Skip() zap.Field                                  { return zap.Skip() }
func Stack() zap.Field                                 { return zap.Stack("stack") }
func String(key string, val string) zap.Field          { return zap.String(key, val) }
func Stringer(key string, val fmt.Stringer) zap.Field  { return zap.Stringer(key, val) }
func Time(key string, val time.Time) zap.Field         { return zap.Time(key, val) }
func Uint(key string, val uint) zap.Field              { return zap.Uint(key, val) }
func Uint64(key string, val uint64) zap.Field          { return zap.Uint64(key, val) }
func Uintptr(key string, val uintptr) zap.Field        { return zap.Uintptr(key, val) }

// Error wraps error for zap.Error.
func Error(err error) zap.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}

func ID(key string, v interface{ Int64() int64 }) zap.Field {
	return zap.Int64(key, v.Int64())
}

type dd struct {
	v interface{}
}

func (d dd) String() string {
	return pp.Sprint(d.v)
}

// Dump renders object for debugging
func Dump(v interface{}) fmt.Stringer {
	return dd{v}
}
