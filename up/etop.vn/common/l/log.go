package l

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const prefix = "etop.vn/backend/"

// ConsoleEncoderName ...
const ConsoleEncoderName = "custom_console"

var ll, xl Logger

// Short-hand functions for logging.
var (
	Bool     = zap.Bool
	Duration = zap.Duration
	Float64  = zap.Float64
	Int      = zap.Int
	Int64    = zap.Int64
	Skip     = zap.Skip
	String   = zap.String
	Stringer = zap.Stringer
	Time     = zap.Time
	Uint     = zap.Uint
	Uint64   = zap.Uint64
	Uintptr  = zap.Uintptr
)

func ID(key string, v interface{ Int64() int64 }) zap.Field {
	return zap.Int64(key, v.Int64())
}

// DefaultConsoleEncoderConfig ...
var DefaultConsoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   ShortColorCallerEncoder,
}

// DefaultTextEncoderConfig ...
var DefaultTextEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// Error wraps error for zap.Error.
func Error(err error) zapcore.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}

// Any ...
func Any(key string, val interface{}) zapcore.Field {
	return zap.Any(key, val)
}

// Stack ...
func Stack() zapcore.Field {
	return zap.Stack("stack")
}

// Int32 ...
func Int32(key string, val int32) zapcore.Field {
	return zap.Int(key, int(val))
}

// Object ...
func Object(key string, val interface{}) zapcore.Field {
	return zap.Stringer(key, Dump(val))
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

func trimPath(c zapcore.EntryCaller) string {
	index := strings.Index(c.File, prefix)
	if index < 0 {
		return c.TrimmedPath()
	}
	return c.File[index+len(prefix):]
}

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	const gray, resetColor = "\x1b[90m", "\x1b[0m"
	callerStr := gray + "→ " + trimPath(caller) + ":" + strconv.Itoa(caller.Line) + resetColor
	enc.AppendString(callerStr)
}

func truncFilename(filename string) string {
	index := strings.Index(filename, prefix)
	return filename[index+len(prefix):]
}
