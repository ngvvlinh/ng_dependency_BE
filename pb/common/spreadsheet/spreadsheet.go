package spreadsheet

import (
	"etop.vn/backend/pkg/common/imcsv"
)

func ToSpreadsheetData(columns imcsv.Schema, idx imcsv.Indexer, rows [][]string, lastRow int) *SpreadsheetData {
	rows = rows[:lastRow+1]
	length := idx.Length()
	if length == 0 {
		length = len(columns)
	}
	cols := make([]string, length)
	for i, col := range columns {
		realIdx := idx.MapIndex(i)
		if realIdx < 0 {
			continue
		}
		if !col.Hidden {
			cols[realIdx] = col.Display
		}
	}

	pbRows := make([]*Row, len(rows)-1)
	for r := 1; r < len(rows); r++ {
		pbRows[r-1] = &Row{
			Cell: rows[r],
		}
	}
	resp := &SpreadsheetData{
		Columns:    cols,
		RawColumns: rows[0],
		Rows:       pbRows,
	}
	return resp
}
