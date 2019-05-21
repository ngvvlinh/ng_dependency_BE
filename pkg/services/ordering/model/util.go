package model

import (
	"fmt"
	"io"
)

func fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}
