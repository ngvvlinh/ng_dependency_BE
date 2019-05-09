package export

import "encoding/csv"

type CellWriter func(index int, header string) string

type TableWriter struct {
	headers    []string
	cells      []func() string
	middleware func(string) string

	writer *csv.Writer
	buf    []string
}

func NewTableWriter(csvWriter *csv.Writer, middleware func(string) string) *TableWriter {
	if middleware == nil {
		middleware = func(s string) string { return s }
	}
	return &TableWriter{writer: csvWriter, middleware: middleware}
}

func (tw *TableWriter) AddColumnExt(header string, cellWriter CellWriter) {
	index := len(tw.cells)
	tw.headers = append(tw.headers, header)
	tw.cells = append(tw.cells, func() string {
		return cellWriter(index, header)
	})
}

func (tw *TableWriter) AddColumn(header string, cellWriter func() string) {
	tw.headers = append(tw.headers, header)
	tw.cells = append(tw.cells, cellWriter)
}

func (tw *TableWriter) WriteHeader() error {
	return tw.writer.Write(tw.headers)
}

func (tw *TableWriter) WriteRow() error {
	if tw.buf == nil {
		tw.buf = make([]string, len(tw.headers))
	}
	buf, middleware := tw.buf, tw.middleware

	for i := range tw.headers {
		buf[i] = middleware(tw.cells[i]())
	}
	return tw.writer.Write(buf)
}
