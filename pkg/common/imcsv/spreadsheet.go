package imcsv

import "etop.vn/api/top/int/types/spreadsheet"

func ToSpreadsheetData(columns Schema, idx Indexer, rows [][]string, lastRow int) *spreadsheet.SpreadsheetData {
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

	pbRows := make([]*spreadsheet.Row, len(rows)-1)
	for r := 1; r < len(rows); r++ {
		pbRows[r-1] = &spreadsheet.Row{
			Cell: rows[r],
		}
	}
	resp := &spreadsheet.SpreadsheetData{
		Columns:    cols,
		RawColumns: rows[0],
		Rows:       pbRows,
	}
	return resp
}
