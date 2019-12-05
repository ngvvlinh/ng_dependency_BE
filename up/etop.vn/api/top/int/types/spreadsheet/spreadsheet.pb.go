package spreadsheet

import (
	"etop.vn/common/jsonx"
)

type SpreadsheetData struct {
	Columns    []string `json:"columns"`
	RawColumns []string `json:"raw_columns"`
	Rows       []*Row   `json:"rows"`
}

func (m *SpreadsheetData) Reset()         { *m = SpreadsheetData{} }
func (m *SpreadsheetData) String() string { return jsonx.MustMarshalToString(m) }

type Row struct {
	Cell []string `json:"cell"`
}

func (m *Row) Reset()         { *m = Row{} }
func (m *Row) String() string { return jsonx.MustMarshalToString(m) }
