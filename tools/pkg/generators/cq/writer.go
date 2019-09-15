package cq

import (
	"bytes"
	"go/types"
	"io"
	"regexp"
)

type MultiWriter struct {
	*Writer
	WriteArgs     bytes.Buffer
	WriteIface    bytes.Buffer
	WriteDispatch bytes.Buffer
}

func (ws *MultiWriter) GetImportWriter(w io.Writer) ImportWriter {
	return importWriterImpl{w, ws.Writer}
}

type ImportWriter interface {
	io.Writer
	Importer
}

type importWriterImpl struct {
	io.Writer
	Importer
}

type Importer interface {
	Import(name, path string)
	TypeString(types.Type) string
}

type Writer struct {
	PackageName string
	PackagePath string

	Importer
	Writer io.Writer
}

func NewWriter(pkgName string, pkgPath string, im Importer, w io.Writer) *Writer {
	return &Writer{
		PackageName: pkgName,
		PackagePath: pkgPath,
		Importer:    im,
		Writer:      w,
	}
}

// /v1 /v1a, /v1beta, /v1/foo
var reVx = regexp.MustCompile(`[a-z0-9]+/v[0-9]+[A-z]*(/[_0-9A-z]+)?$`)

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}
