package util

import (
	"o.o/api/summary"
	smry "o.o/backend/pkg/etop/logic/summary"
)

func BuildResponse(tables []*smry.Table) []*summary.SummaryTable {
	res := make([]*summary.SummaryTable, len(tables))
	for i, table := range tables {
		stable := &summary.SummaryTable{
			Label: table.Label,
			Tags:  table.Tags,
			Cols:  buildCols(table.Cols),
			Rows:  buildRows(table.Rows),
			Data:  buildData(table.Data),
		}
		res[i] = stable
	}
	return res
}

func buildCols(cols []smry.Predicator) []summary.SummaryColRow {
	res := make([]summary.SummaryColRow, len(cols))
	for i, col := range cols {
		res[i] = summary.SummaryColRow{
			Label:  col.GetLabel(),
			Spec:   col.GetSpec(),
			Unit:   "",
			Indent: 0,
		}
	}
	return res
}

func buildRows(rows []smry.Subject) []summary.SummaryColRow {
	res := make([]summary.SummaryColRow, len(rows))
	for i, row := range rows {
		res[i] = summary.SummaryColRow{
			Label:  row.GetLabel(),
			Spec:   row.GetSpec(),
			Unit:   row.Unit,
			Indent: row.Ident,
		}
	}
	return res
}

func buildData(data []smry.Cell) []summary.SummaryItem {
	res := make([]summary.SummaryItem, len(data))
	for i, item := range data {
		res[i] = summary.SummaryItem{
			Spec:  item.Subject.GetSpec(),
			Value: item.Value,
			Unit:  item.Subject.Unit,
		}
	}
	return res
}
