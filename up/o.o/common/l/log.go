package l

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ll, xl Logger

func init() {
	if err := zap.RegisterEncoder(ConsoleEncoderName,
		func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return NewConsoleEncoder(cfg), nil
		}); err != nil {
		panic(err)
	}

	ll = New()
	xl = New(WrapOption(zap.AddCallerSkip(1)))

	var envLog string
	for _, envKey := range envKeys {
		envLog = os.Getenv(envKey)
		if envLog != "" {
			break
		}
	}
	if envLog == "" {
		return
	}

	_envPatterns, errPattern, err := parseWildcardPatterns(envLog)
	if err != nil {
		ll.Fatal(fmt.Sprintf("unable to parse `%v`", EnvKey), String("invalid", errPattern), Error(err))
	}
	envPatterns = _envPatterns
	ll.Debug("enable debug log", String(EnvKey, envLog))
}

// Logger wraps zap.Logger
type Logger struct {
	opts    []zap.Option
	enabler AtomicLevel
	*zap.Logger

	S  *zap.SugaredLogger
	v  VerbosedLogger
	ch *Messenger
}

// New returns new zap.Logger
func New(opts ...Option) Logger {
	_, filename, _, _ := runtime.Caller(1)
	name := filepath.Dir(truncFilename(filename))

	var enabler AtomicLevel
	if e, ok := enablers[name]; ok {
		enabler = e
	} else {
		enabler = NewAtomicLevel(name)
		enablers[name] = enabler
	}

	setLogLevelFromPatterns(envPatterns, name, enabler)
	loggerConfig := zap.Config{
		Level:            enabler.AtomicLevel,
		Development:      false,
		Encoding:         ConsoleEncoderName,
		EncoderConfig:    DefaultConsoleEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	stacktraceLevel := zap.NewAtomicLevelAt(zapcore.PanicLevel)
	l := Logger{
		enabler: enabler,
		v: VerbosedLogger{
			enabler: enabler.AtomicLevel,
		},
	}
	for _, opt := range opts {
		opt(&l)
	}

	l.opts = append(l.opts, zap.AddStacktrace(stacktraceLevel))
	logger, err := loggerConfig.Build(l.opts...)
	if err != nil {
		panic(err)
	}

	verbose := logger.WithOptions(zap.AddCallerSkip(1))
	l.Logger = logger
	l.S = logger.Sugar()
	l.v.logger = verbose
	l.v.sugar = verbose.Sugar()
	return l
}

func (l Logger) Verbosed(verbosity int) bool {
	return l.enabler.Enabled(V(verbosity))
}

func (l Logger) Enabled(level zapcore.Level) bool {
	return l.enabler.Enabled(level)
}

func (l Logger) Watch(fn LevelWatcher) (unwatch func()) {
	return l.enabler.Watch(fn)
}

func (l Logger) Sync() {
	_ = l.Logger.Sync()
}

func (l Logger) Must(err error, msg string, fields ...zap.Field) {
	if err == nil {
		return
	}
	fs := make([]zap.Field, 0, len(fields)+1)
	fs = append(fs, Error(err))
	fs = append(fs, fields...)
	ll.Panic(msg, fields...)
}

func (l Logger) WithChannel(channel string) Logger {
	l.ch = getChannel(channel)
	return l
}

// SendMessage sends a message to Telegram
//
// TODO(vu): implement a watcher system instead
func (l *Logger) SendMessage(msg string) {
	if l.ch == nil {
		// retrieve the default channel
		l.ch = getChannel("")
	}
	(*l.ch).SendMessage(msg)
}
