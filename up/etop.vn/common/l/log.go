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
	xl = New(zap.AddCallerSkip(1))

	envLog := os.Getenv(EnvKey)
	if envLog == "" {
		envLog = os.Getenv(deprecatedEnvKey)
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
	enabler AtomicLevel
	*zap.Logger

	S *zap.SugaredLogger
	v VerboseLogger
}

// New returns new zap.Logger
func New(opts ...zap.Option) Logger {
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

	opts = append(opts, zap.AddStacktrace(stacktraceLevel))
	logger, err := loggerConfig.Build(opts...)
	if err != nil {
		panic(err)
	}

	verbose := logger.WithOptions(zap.AddCallerSkip(1))
	l := Logger{
		enabler: enabler,
		Logger:  logger,
		S:       logger.Sugar(),
		v: VerboseLogger{
			enabler: enabler.AtomicLevel,
			logger:  verbose,
			sugar:   verbose.Sugar(),
		},
	}
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
