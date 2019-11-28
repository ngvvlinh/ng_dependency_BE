package imcsv

import (
	"errors"
	"math"
	"strconv"
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

type ColumnDef struct {
	Name       string
	Display    string
	Norm       string
	Contains   string
	Line       bool
	Hidden     bool
	Optional   bool
	Translator func(name, value string) (interface{}, error)
}

type Schema []*ColumnDef

func (schema Schema) MustValidate() {
	for i, col := range schema {
		if col.Name == "" || col.Display == "" || col.Norm == "" {
			ll.Panic("invalid schema definition", l.Int("i", i))
		}
		if !strings.Contains(validate.NormalizeSearchSimple(col.Display), col.Norm) {
			ll.Panic("invalid schema definition (display must contain norm)", l.Int("i", i))
		}
	}
}

func (schema Schema) Indexer() func(string) int {
	return func(name string) int {
		for i, col := range schema {
			if col.Name == name {
				return i
			}
		}
		ll.Panic("No column", l.String("name", name))
		return -1
	}
}

// Indexer maps between column in schema and column in imported files
type Indexer struct {
	schema   Schema
	length   int
	mapIndex []int
}

func (idx Indexer) Length() int {
	return idx.length
}

func (idx Indexer) MapIndex(i int) int {
	if idx.mapIndex == nil {
		return i
	}
	if i < 0 {
		return -1
	}
	return idx.mapIndex[i]
}

func (idx Indexer) GetCell(row []string, i int) string {
	realIdx := idx.mapIndex[i]
	if realIdx < 0 {
		return ""
	}
	return row[realIdx]
}

type SchemaAndIndexer struct {
	Schema  Schema
	Indexer interface{}
}

func (schema Schema) MinColumns() int {
	minCol := 0
	for _, col := range schema {
		if !col.Optional {
			minCol++
		}
	}
	return minCol
}

func (schema Schema) ValidateSchema(headerRow *[]string) (idx Indexer, errs []error, _ error) {
	minCol := schema.MinColumns()
	if len(*headerRow) < minCol {
		return Indexer{}, nil, CellError(Indexer{}, 0, -1, "Số cột không đúng cấu trúc yêu cầu. Vui lòng tải lại file import hoặc liên hệ hotro@etop.vn.").
			WithMetap("expect_min", minCol).
			WithMetap("expect_max", len(schema)).
			WithMetap("actual", len(*headerRow))
	}

	hRow := *headerRow
	indexes := make([]int, len(schema))
	realIdx := 0
	for i, col := range schema {
		label := ""
		if realIdx < len(hRow) {
			label = hRow[realIdx]
		}
		if label == col.Display {
			indexes[i] = realIdx
			realIdx++
			continue
		}

		labelNorm := validate.NormalizeSearchSimple(label)
		if strings.Contains(labelNorm, col.Norm) &&
			(col.Contains == "" || strings.Contains(label, col.Contains)) {
			indexes[i] = realIdx
			realIdx++
			continue
		}
		if col.Optional {
			indexes[i] = -1
			continue
		}

		err := cellErrorWithCode(cm.InvalidArgument, nil, 0, realIdx, "Cột không đúng cấu trúc yêu cầu (mong đợi: %v).", col.Display)
		errs = append(errs, err)

		// increase realIdx even if the column is not found (and not optional)
		realIdx++
	}

	if errs == nil {
		// realIdx now is the number of columns in the actual file
		CleanRow(headerRow, realIdx)
		return Indexer{length: realIdx, mapIndex: indexes}, nil, nil
	}
	return Indexer{}, errs, nil
}

func ValidateAgainstSchemas(headerRow *[]string, schemas []Schema) (_ int, idx Indexer, errs []error, err error) {
	for i, schema := range schemas {
		idx, errs, err = schema.ValidateSchema(headerRow)
		if err == nil && len(errs) == 0 {
			return i, idx, nil, nil
		}
	}
	// Return the latest version error
	return len(schemas) - 1, Indexer{}, errs, err
}

func CellError(idx Indexer, row, col int, msg string, args ...interface{}) *xerrors.APIError {
	return CellErrorWithCode(idx, cm.InvalidArgument, nil, row, col, msg, args...)
}

func CellErrorWithCode(idx Indexer, code xerrors.Code, err error, row, col int, msg string, args ...interface{}) *xerrors.APIError {
	return cellErrorWithCode(code, err, row, idx.MapIndex(col), msg, args...)
}

func cellErrorWithCode(code xerrors.Code, err error, row, col int, msg string, args ...interface{}) *xerrors.APIError {
	_err := cm.NSErrorf(code, err, msg, args...).
		WithMeta("row_index", strconv.Itoa(row)).
		WithMeta("row", strconv.Itoa(row+1))
	if col >= 0 {
		_err = _err.
			WithMeta("col_index", strconv.Itoa(col)).
			WithMeta("col", ColName(col))
	}
	return _err
}

func ColName(col int) string {
	if col <= 'Z'-'A' {
		return string('A' + col)
	}
	col = col - ('Z' - 'A' + 1)
	return "A" + string('A'+col)
}

func CellName(row, col int) string {
	return ColName(col) + strconv.Itoa(row+1)
}

func ParseFloat(v string) (float64, error) {
	switch v {
	case "", "-":
		return 0, nil
	}

	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0, errors.New("Cần là số.")
	}
	if f < 0 {
		return 0, errors.New("Giá trị nhỏ hơn 0.")
	}
	return f, err
}

func ParseUint(v string) (int, error) {
	f, err := ParseFloat(v)
	if err != nil {
		return 0, err
	}
	if f-math.Floor(f) >= 0.001 {
		return 0, errors.New("Cần là số nguyên.")
	}
	return int(f), nil
}

func ParseBool(v string) (bool, error) {
	switch v {
	case "Có":
		return true, nil
	case "Không":
		return false, nil
	}

	switch validate.NormalizeSearchSimple(v) {
	case "co":
		return true, nil
	case "khong":
		return false, nil
	}
	return false, errors.New("Giá trị không hợp lệ, cần một trong các giá trị 'Có' hoặc 'Không'.")
}

func ParsePercent(v string) (float64, error) {
	if v == "" {
		return 0, errors.New("Giá trị phần trăm không hợp lệ.")
	}
	if v[len(v)-1] != '%' {
		return 0, errors.New("Giá trị phần trăm không hợp lệ.")
	}

	v = v[:len(v)-1]
	v = strings.TrimSpace(v)
	p, err := ParseFloat(v)
	return p / 100, err
}

func CleanRow(rows *[]string, minCol int) {
	rs := *rows
	nonEmptyCol := 0
	for i := range rs {
		rs[i] = strings.TrimSpace(rs[i])
		if rs[i] != "" {
			nonEmptyCol = i
		}
	}
	nonEmptyCol++
	if nonEmptyCol < minCol {
		nonEmptyCol = minCol
	}
	*rows = rs[:nonEmptyCol]
}

func IsRowEmpty(row []string) bool {
	for _, v := range row {
		if v != "" {
			return false
		}
	}
	return true
}

func GetFormValue(ss []string) string {
	if ss == nil {
		return ""
	}
	return ss[0]
}
