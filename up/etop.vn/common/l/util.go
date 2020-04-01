package l

import (
	"fmt"
	"time"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Any(key string, val interface{}) zapcore.Field        { return zap.Any(key, val) }
func Bool(key string, val bool) zapcore.Field              { return zap.Bool(key, val) }
func Duration(key string, val time.Duration) zapcore.Field { return zap.Duration(key, val) }
func Float64(key string, val float64) zapcore.Field        { return zap.Float64(key, val) }
func Int(key string, val int) zapcore.Field                { return zap.Int(key, val) }
func Int32(key string, val int32) zapcore.Field            { return zap.Int(key, int(val)) }
func Int64(key string, val int64) zapcore.Field            { return zap.Int64(key, val) }
func Object(key string, val interface{}) zapcore.Field     { return zap.Stringer(key, dump(val)) }
func Skip() zapcore.Field                                  { return zap.Skip() }
func Stack() zapcore.Field                                 { return zap.Stack("stack") }
func String(key string, val string) zapcore.Field          { return zap.String(key, val) }
func Stringer(key string, val fmt.Stringer) zapcore.Field  { return zap.Stringer(key, val) }
func Time(key string, val time.Time) zapcore.Field         { return zap.Time(key, val) }
func Uint(key string, val uint) zapcore.Field              { return zap.Uint(key, val) }
func Uint64(key string, val uint64) zapcore.Field          { return zap.Uint64(key, val) }
func Uintptr(key string, val uintptr) zapcore.Field        { return zap.Uintptr(key, val) }

// Error wraps error for zap.Error.
func Error(err error) zapcore.Field {
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
func dump(v interface{}) fmt.Stringer {
	return dd{v}
}
