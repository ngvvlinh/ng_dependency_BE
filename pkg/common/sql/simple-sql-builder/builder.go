package builder

import (
	"fmt"
	"strconv"
	"strings"

	cm "o.o/backend/pkg/common"
)

type Raw string

type SimpleSQLBuilder struct {
	buf strings.Builder
	err error
}

func (b *SimpleSQLBuilder) Printf(query string, args ...interface{}) {
	if b.err != nil {
		return
	}

	argIndex := 0
	queryIndex := 0
	for i := 0; i < len(query); i++ {
		c := query[i]
		if c != '?' {
			continue
		}

		b.buf.WriteString(query[queryIndex:i])
		queryIndex = i + 1
		if len(args) <= argIndex {
			b.err = fmt.Errorf("not enough argument at index %v", argIndex)
			return
		}
		arg := args[argIndex]
		argIndex++
		switch arg := arg.(type) {
		case Raw:
			b.buf.WriteString(string(arg))
		case string:
			arg, b.err = singleQuote(arg)
			b.buf.WriteString(arg)
		case bool:
			s := strconv.FormatBool(arg)
			b.buf.WriteString(s)
		case int64:
			b.buf.WriteString(strconv.FormatInt(arg, 10))
		case int32:
			b.buf.WriteString(strconv.FormatInt(int64(arg), 10))
		case int16:
			b.buf.WriteString(strconv.FormatInt(int64(arg), 10))
		case int8:
			b.buf.WriteString(strconv.FormatInt(int64(arg), 10))
		case int:
			b.buf.WriteString(strconv.FormatInt(int64(arg), 10))
		default:
			b.err = fmt.Errorf("unexpected argument of type %T (%v)", arg, arg)
		}
	}
	b.buf.WriteString(query[queryIndex:])
	if len(args) != argIndex {
		b.err = fmt.Errorf("expected %v arguments but %v arguments were provided", argIndex, len(args))
	}
}

func (b *SimpleSQLBuilder) String() (string, error) {
	return b.buf.String(), b.err
}

// singleQuote does not support strings with quote
func singleQuote(value string) (string, error) {
	if strings.Contains(value, `"`) || strings.Contains(value, `'`) {
		return "", cm.Errorf(cm.InvalidArgument, nil, "value string contains single or double-quote")
	}
	return strings.ReplaceAll(strconv.Quote(value), `"`, `'`), nil
}
