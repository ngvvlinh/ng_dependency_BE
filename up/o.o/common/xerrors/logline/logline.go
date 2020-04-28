package logline

import (
	"strconv"

	zc "go.uber.org/zap/zapcore"

	"o.o/common/jsonx"
)

type LogAppender interface {
	AppendLog(LogLine)
}

func ValueOf(f zc.Field) interface{} {
	switch {
	case f.Integer != 0:
		return f.Integer
	case f.String != "":
		return f.String
	}
	return f.Interface
}

// LogLine ...
type LogLine struct {
	File    string
	Line    int
	Message string
	Fields  []zc.Field
}

// MarshalJSON ...
func (l LogLine) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 512)
	return l.MarshalTo(b), nil
}

// MarshalTo ...
func (l LogLine) MarshalTo(b []byte) []byte {
	b = append(b, '{')

	b = append(b, `"@msg":`...)
	b = append(b, marshal(l.Message)...)
	b = append(b, ',')

	b = append(b, `"@file":`...)
	b = append(b, marshal(l.File+":"+strconv.Itoa(l.Line))...)

	for _, field := range l.Fields {
		b = append(b, ',')
		b = append(b, marshal(field.Key)...)
		b = append(b, ':')

		if field.Integer != 0 {
			b = append(b, strconv.Itoa(int(field.Integer))...)
		} else if field.String != "" {
			b = append(b, marshal(field.String)...)
		} else {
			b = append(b, marshal(field.Interface)...)
		}
	}
	b = append(b, '}')
	return b
}

func marshal(v interface{}) []byte {
	data, err := jsonx.Marshal(v)
	if err != nil {
		data, _ = jsonx.Marshal(err)
	}
	return data
}
