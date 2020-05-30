package l

import "go.uber.org/zap"

type Option func(*Logger)

func WrapOption(zapOpts ...zap.Option) Option {
	return func(l *Logger) {
		l.opts = append(l.opts, zapOpts...)
	}
}
