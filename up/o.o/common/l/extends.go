package l

import (
	"context"
	"runtime"

	"go.uber.org/zap/zapcore"

	"o.o/common/xerrors/logline"
)

// XError appends error to current context
func XError(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx, ok := ctx.(logline.LogAppender); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.AppendLog(logline.LogLine{
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
	if ctx, ok := ctx.(logline.LogAppender); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.AppendLog(logline.LogLine{
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
	if ctx, ok := ctx.(logline.LogAppender); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.AppendLog(logline.LogLine{
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
	if ctx, ok := ctx.(logline.LogAppender); ok {
		_, file, line, _ := runtime.Caller(1)
		ctx.AppendLog(logline.LogLine{
			File:    file,
			Line:    line,
			Fields:  fields,
			Message: msg,
		})
	} else {
		xl.Debug(msg, fields...)
	}
}
