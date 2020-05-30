package l

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level
type Field = zapcore.Field

func V(verbosity int) Level {
	if verbosity > 0 {
		verbosity = -verbosity
	}
	return Level(verbosity)
}

type VerbosedLogger struct {
	enabler zap.AtomicLevel
	lvl     Level
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
}

func (l Logger) V(lvl int) VerbosedLogger {
	v := l.v // clone the struct
	v.lvl = Level(-lvl)
	return v
}

func (l VerbosedLogger) Enabled() bool {
	return l.enabler.Enabled(l.lvl)
}

func (l VerbosedLogger) Debug(msg string, fields ...zapcore.Field) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l VerbosedLogger) Debugf(template string, args ...interface{}) {
	if !l.logger.Core().Enabled(l.lvl) {
		return
	}

	msg := template
	if template == "" {
		msg = fmt.Sprint(args...)
	} else if len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		ce.Write()
	}
}

func (l VerbosedLogger) Any(msg string, objs ...interface{}) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		fields := make([]zapcore.Field, len(objs))
		for i, obj := range objs {
			fields[i] = Any(strconv.Itoa(i), obj)
		}
		ce.Write(fields...)
	}
}

func (l VerbosedLogger) Dump(msg string, objs ...interface{}) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		fields := make([]zapcore.Field, len(objs))
		for i, obj := range objs {
			fields[i] = Object(strconv.Itoa(i), obj)
		}
		ce.Write(fields...)
	}
}
