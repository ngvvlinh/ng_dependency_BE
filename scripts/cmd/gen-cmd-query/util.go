package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func must(err error) {
	if err != nil {
		fatalf("%v\n", err)
	}
}

func debugf(format string, args ...interface{}) {
	if *flV {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
	}
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

func mustWrite(w io.Writer, p []byte) {
	_, err := w.Write(p)
	must(err)
}

func formatFileName(format string, fileName string) string {
	if strings.Contains(format, "{}") {
		return strings.Replace(format, "{}", fileName, 1)
	}
	if strings.HasPrefix(format, ".") {
		return fileName + format
	}
	return format
}
