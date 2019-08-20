package l

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type VerboseLogger struct {
	lvl    zapcore.Level
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func (l Logger) V(lvl int) VerboseLogger {
	v := l.verbose
	v.lvl = zapcore.Level(-lvl)
	return v
}

func (l VerboseLogger) Debug(msg string, fields ...zapcore.Field) {
	if ce := l.logger.Check(l.lvl, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l VerboseLogger) Debugf(template string, args ...interface{}) {
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
