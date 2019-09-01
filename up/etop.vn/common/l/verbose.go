package l

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func V(verbosity int) zapcore.Level {
	if verbosity > 0 {
		verbosity = -verbosity
	}
	return zapcore.Level(verbosity)
}

type VerboseLogger struct {
	enabler zap.AtomicLevel
	lvl     zapcore.Level
	logger  *zap.Logger
	sugar   *zap.SugaredLogger
}

func (l Logger) V(lvl int) VerboseLogger {
	v := l.v // clone the struct
	v.lvl = zapcore.Level(-lvl)
	return v
}

func (l VerboseLogger) Enabled() bool {
	return l.enabler.Enabled(l.lvl)
}

func (l VerboseLogger) Debug(msg string, fields ...zapcore.Field) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l VerboseLogger) Debugf(template string, args ...interface{}) {
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

func (l VerboseLogger) Any(msg string, objs ...interface{}) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		fields := make([]zapcore.Field, len(objs))
		for i, obj := range objs {
			fields[i] = Any(strconv.Itoa(i), obj)
		}
		ce.Write(fields...)
	}
}

func (l VerboseLogger) Dump(msg string, objs ...interface{}) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		fields := make([]zapcore.Field, len(objs))
		for i, obj := range objs {
			fields[i] = Object(strconv.Itoa(i), obj)
		}
		ce.Write(fields...)
	}
}
