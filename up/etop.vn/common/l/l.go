package l

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"etop.vn/common/bus"
	"etop.vn/common/xerrors/logline"
)

const prefix = "etop.vn/backend/"

// ConsoleEncoderName ...
const ConsoleEncoderName = "custom_console"

var ll, xl Logger

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
	S *zap.SugaredLogger
}

// XError appends error to current context
func XError(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx, ok := ctx.(*bus.NodeContext); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.Logs = append(ctx.Logs, logline.LogLine{
			//Level:  "error",
			File:   file,
			Line:   line,
			Fields: fields,
		})
	} else {
		xl.Error(msg, fields...)
	}
}

// XWarn appends debug to current context
func XWarn(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx, ok := ctx.(*bus.NodeContext); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.Logs = append(ctx.Logs, logline.LogLine{
			//Level:  "warn",
			File:   file,
			Line:   line,
			Fields: fields,
		})
	} else {
		xl.Warn(msg, fields...)
	}
}

// XInfo appends debug to current context
func XInfo(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx, ok := ctx.(*bus.NodeContext); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.Logs = append(ctx.Logs, logline.LogLine{
			//Level:   "info",
			File:    file,
			Line:    line,
			Fields:  fields,
			Message: msg,
		})
	} else {
		xl.Info(msg, fields...)
	}
}

// XDebug appends debug to current context
func XDebug(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx, ok := ctx.(*bus.NodeContext); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.Logs = append(ctx.Logs, logline.LogLine{
			//Level:   "debug",
			File:    file,
			Line:    line,
			Fields:  fields,
			Message: msg,
		})
	} else {
		xl.Debug(msg, fields...)
	}
}

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

// DefaultConsoleEncoderConfig ...
var DefaultConsoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
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
	EncodeLevel:    zapcore.CapitalLevelEncoder,
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

func trimPath(c zapcore.EntryCaller) string {
	index := strings.Index(c.File, prefix)
	if index < 0 {
		return c.TrimmedPath()
	}
	return c.File[index+len(prefix):]
}

// // Observer ...
// type Observer interface {
// 	LogObserver(t time.Time, level, msg string, fields []zapcore.Field)
// }

// var observer Observer

// // RegisterObserver ...
// func RegisterObserver(fn Observer) {
// 	if observer != nil {
// 		ll.Panic("Already register")
// 	}
// 	observer = fn
// }

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	const gray, resetColor = "\x1b[90m", "\x1b[0m"
	callerStr := gray + "→ " + trimPath(caller) + ":" + strconv.Itoa(caller.Line) + resetColor
	enc.AppendString(callerStr)
}

// New returns new zap.Logger
func New(opts ...zap.Option) Logger {
	_, filename, _, _ := runtime.Caller(1)
	name := filepath.Dir(truncFilename(filename))

	var enabler zap.AtomicLevel
	if e, ok := enablers[name]; ok {
		enabler = e
	} else {
		enabler = zap.NewAtomicLevel()
		enablers[name] = enabler
	}

	setLogLevelFromEnv(name, enabler)
	loggerConfig := zap.Config{
		Level:            enabler,
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
	return Logger{logger, logger.Sugar()}
}

// // Debug ...
// func (l Logger) Debug(msg string, fields ...zapcore.Field) {
// 	l.Logger.Debug(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "debug", msg, fields)
// 	}
// }

// // Info ...
// func (l Logger) Info(msg string, fields ...zapcore.Field) {
// 	l.Logger.Info(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "info", msg, fields)
// 	}
// }

// // Warn ...
// func (l Logger) Warn(msg string, fields ...zapcore.Field) {
// 	l.Logger.Warn(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "warn", msg, fields)
// 	}
// }

// // Error ...
// func (l Logger) Error(msg string, fields ...zapcore.Field) {
// 	l.Logger.Error(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "error", msg, fields)
// 	}
// }

// // DPanic ...
// func (l Logger) DPanic(msg string, fields ...zapcore.Field) {
// 	l.Logger.DPanic(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "dpanic", msg, fields)
// 	}
// }

// // Panic ...
// func (l Logger) Panic(msg string, fields ...zapcore.Field) {
// 	l.Logger.Panic(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "panic", msg, fields)
// 	}
// }

// // Fatal ...
// func (l Logger) Fatal(msg string, fields ...zapcore.Field) {
// 	l.Logger.Fatal(msg, fields...)
// 	if observer != nil {
// 		observer.LogObserver(time.Now(), "fatal", msg, fields)
// 	}
// }

func (l Logger) Sync() {
	_ = l.Logger.Sync()
}

// ServeHTTP supports logging level with an HTTP request.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	type payload struct {
		Name  string         `json:"name"`
		Level *zapcore.Level `json:"level"`
	}

	enc := json.NewEncoder(w)

	switch r.Method {
	case "GET":
		var payloads []payload
		for k, e := range enablers {
			lvl := e.Level()
			payloads = append(payloads, payload{
				Name:  k,
				Level: &lvl,
			})
		}
		enc.Encode(payloads)

	case "PUT":
		var req payload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: fmt.Sprintf("Request body must be valid JSON: %v", err),
			})
			return
		}

		enabler, ok := enablers[req.Name]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: errEnablerNotFound.Error(),
			})
			return
		}

		if req.Level == nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: errLevelNil.Error(),
			})
			return
		}

		enabler.SetLevel(*req.Level)
		enc.Encode(req)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		enc.Encode(errorResponse{
			Error: "Only GET and PUT are supported.",
		})
	}
}

var (
	errEnablerNotFound = errors.New("enabler not found")
	errLevelNil        = errors.New("must specify a logging level")

	enablers = make(map[string]zap.AtomicLevel)
)

func truncFilename(filename string) string {
	index := strings.Index(filename, prefix)
	return filename[index+len(prefix):]
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

var envPatterns []*regexp.Regexp

func init() {
	zap.RegisterEncoder(ConsoleEncoderName, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewConsoleEncoder(cfg), nil
	})

	envLog := os.Getenv("ETOP_LOG_DEBUG")
	if envLog == "" {
		return
	}

	ll = New()
	xl = New(zap.AddCallerSkip(1))

	var errPattern string
	envPatterns, errPattern = initPatterns(envLog)
	if errPattern != "" {
		ll.Fatal("Unable to parse ETOP_LOG_DEBUG. Please set it to a proper value.", String("invalid", errPattern))
	}

	ll.Info("Enable debug log", String("ETOP_LOG_DEBUG", envLog))
}

func initPatterns(envLog string) ([]*regexp.Regexp, string) {
	patterns := strings.Split(envLog, ",")
	result := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		r, err := parsePattern(p)
		if err != nil {
			return nil, p
		}

		result[i] = r
	}
	return result, ""
}

func parsePattern(p string) (*regexp.Regexp, error) {
	p = strings.Replace(strings.Trim(p, " "), "*", ".*", -1)
	return regexp.Compile(p)
}

func setLogLevelFromEnv(name string, enabler zap.AtomicLevel) {
	for _, r := range envPatterns {
		if r.MatchString(name) {
			enabler.SetLevel(zap.DebugLevel)
		}
	}
}
