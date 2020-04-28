package api

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"o.o/common/l"
)

var ll = l.New()

// /v1 /v1a, /v1beta, /v1/foo
var reVx = regexp.MustCompile(`[a-z0-9]+/v[0-9]+[A-z]*(/[_0-9A-z]+)?$`)

func must(err error) {
	if err != nil {
		fatalf("%v\n", err)
	}
}

func debugf(format string, args ...interface{}) {
	ll.V(1).Debugf(format, args...)
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var storedErrors []string

func errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	storedErrors = append(storedErrors, msg)
}

func mustNoError(format string, args ...interface{}) {
	count := len(storedErrors)
	if count == 0 {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	for _, msg := range storedErrors {
		msg = strings.TrimRight(msg, "\n")
		_, _ = fmt.Fprintf(os.Stderr, "    %v\n", msg)
	}
	switch {
	case count == 1:
		fatalf("stopped due to %v error\n", count)
	case count > 1:
		fatalf("stopped due to %v errors\n", count)
	}
}

func p(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	must(err)
}
