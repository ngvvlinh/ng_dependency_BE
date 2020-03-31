package l

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Foreground colors.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

var (
	_levelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  Magenta,
		zapcore.InfoLevel:   Blue,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.DPanicLevel: Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
	}
	_unknownLevelColor = Red

	_levelToLowercaseColorString = make(map[zapcore.Level]string, len(_levelToColor))
	_levelToCapitalColorString   = make(map[zapcore.Level]string, len(_levelToColor))
)

func init() {
	for level, c := range _levelToColor {
		_levelToLowercaseColorString[level] = c.Add(level.String())
		_levelToCapitalColorString[level] = c.Add(level.CapitalString())
	}
	debugColor := _levelToColor[zapcore.DebugLevel]
	for level := zapcore.DebugLevel - 1; level >= -MaxVerbosity; level-- {
		_levelToLowercaseColorString[level] = debugColor.Add(stringLevel(level))
		_levelToCapitalColorString[level] = debugColor.Add(capitalLevel(level))
	}
}

func stringLevel(level zapcore.Level) string {
	if level >= zapcore.DebugLevel || level < -MaxVerbosity {
		return level.String()
	}
	switch level {
	case -2:
		return "debug-2"
	case -3:
		return "debug-3"
	case -4:
		return "debug-4"
	case -5:
		return "debug-5"
	case -6:
		return "debug-6"
	case -7:
		return "debug-7"
	case -8:
		return "debug-8"
	case -9:
		return "debug-9"
	default:
		return level.String()
	}
}

func capitalLevel(level zapcore.Level) string {
	if level >= zapcore.DebugLevel {
		return level.CapitalString()
	}
	switch level {
	case -2:
		return "DEBUG-2"
	case -3:
		return "DEBUG-3"
	case -4:
		return "DEBUG-4"
	case -5:
		return "DEBUG-5"
	case -6:
		return "DEBUG-6"
	case -7:
		return "DEBUG-7"
	case -8:
		return "DEBUG-8"
	case -9:
		return "DEBUG-9"
	default:
		return level.CapitalString()
	}
}

func unmarshalLevel(s string) (zapcore.Level, bool) {
	if strings.HasPrefix(s, "DEBUG-") || strings.HasPrefix(s, "debug-") {
		levelStr := s[len("DEBUG-"):]
		lvl, err := strconv.Atoi(levelStr)
		if err != nil {
			return 0, false
		}
		if lvl < 1 || lvl > MaxVerbosity {
			return 0, false
		}
		return zapcore.Level(-lvl), true
	}

	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(s))
	return lvl, err == nil
}

func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(capitalLevel(l))
}

func CapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := _levelToCapitalColorString[l]
	if !ok {
		s = _unknownLevelColor.Add(l.CapitalString())
	}
	enc.AppendString(s)
}

func trimPath(c zapcore.EntryCaller) string {
	index := strings.Index(c.File, pathPrefix)
	if index < 0 {
		return c.TrimmedPath()
	}
	return c.File[index+len(pathPrefix):]
}

// ShortColorCallerEncoder encodes caller information with sort path filename and enable
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	const gray, resetColor = "\x1b[90m", "\x1b[0m"
	callerStr := gray + "→ " + trimPath(caller) + ":" + strconv.Itoa(caller.Line) + resetColor
	enc.AppendString(callerStr)
}

func truncFilename(filename string) string {
	index := strings.Index(filename, pathPrefix)
	return filename[index+len(pathPrefix):]
}
