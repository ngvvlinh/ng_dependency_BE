package util

import (
	"time"

	"o.o/api/summary"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/sq"
	smry "o.o/backend/pkg/etop/logic/summary"
	"o.o/capi/dot"
)

const dateLayout = "2006-01-02"

func BuildKey(
	shopID dot.ID, dateFrom, dateTo time.Time,
	keycode, currentVersion string) string {
	key := "summary/" + keycode + ":version=" + currentVersion +
		",shop=" + shopID.String() +
		",from=" + dateFrom.Format(dateLayout) +
		",to=" + dateTo.Format(dateLayout)
	return key
}

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

func BuildRowPerDate(dateFrom time.Time, dateTo time.Time) []smry.Predicator {
	var result []smry.Predicator
	var timeStart = dateFrom
	for timeStart.Before(dateTo) {
		result = append(result, smry.Predicate{
			Spec:  timeStart.Format(dateLayout),
			Label: timeStart.Format(dateLayout),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", timeStart, timeStart.Add(24*time.Hour)),
		})
		timeStart = timeStart.Add(24 * time.Hour)
	}
	return result
}

//
// merge tables functions
//

type mergeOption string

const (
	SameRows mergeOption = "rows"
	SameCols mergeOption = "cols"
	SameBoth mergeOption = "both"
)

func MergeTables(tables []*smry.Table, _mergeOption mergeOption) (newTable *smry.Table, err error) {
	if len(tables) <= 1 {
		return nil, cm.Errorf(cm.Internal, nil, "must have rather than 1 table")
	}
	newTable = tables[0]

	for i := 1; i < len(tables); i++ {
		currentTable := tables[i]
		switch _mergeOption {
		case SameBoth:
			newTable, err = mergeTwoTablesSameBoth(newTable, currentTable)
			return nil, err
		case SameCols:
			newTable, err = mergeTwoTablesSameCols(newTable, currentTable)
			if err != nil {
				return nil, err
			}
		case SameRows:
			fallthrough
		default:
			return nil, cm.Errorf(cm.Internal, nil, "unsupported merge option %s", _mergeOption)
		}
	}

	return newTable, nil
}

func mergeTwoTablesSameBoth(firstTable, secondTable *smry.Table) (*smry.Table, error) {
	mapColTagsFirstTable := make(map[string]bool)
	mapRowTagsFirstTable := make(map[string]bool)

	if len(firstTable.Cols) != len(secondTable.Cols) {
		return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because number of cols not same")
	}
	if len(firstTable.Rows) != len(secondTable.Rows) {
		return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because number of rows not same")
	}

	for _, col := range firstTable.Cols {
		mapColTagsFirstTable[col.GetLabel()] = true
	}
	for _, row := range firstTable.Rows {
		mapRowTagsFirstTable[row.GetLabel()] = true
	}

	for _, col := range secondTable.Cols {
		if _, ok := mapColTagsFirstTable[col.GetLabel()]; !ok {
			return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because col tag[%s] not found in first table", col.GetLabel())
		}
	}
	for _, row := range secondTable.Rows {
		if _, ok := mapRowTagsFirstTable[row.GetLabel()]; !ok {
			return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because row tag[%s] not found in first table", row.GetLabel())
		}
	}

	newTable := firstTable.Clone(firstTable.Label)
	var secondData []smry.Cell
	copy(secondTable.Data, secondData)
	newTable.Data = append(newTable.Data, secondData...)

	return newTable, nil
}

func mergeTwoTablesSameCols(firstTable, secondTable *smry.Table) (*smry.Table, error) {
	mapColTagsFirstTable := make(map[string]bool)
	numOfRowsFirstTable := len(firstTable.Rows)
	numOfRowsSecondTable := len(secondTable.Rows)

	if len(firstTable.Cols) != len(secondTable.Cols) {
		return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because number of cols not same")
	}

	for _, col := range firstTable.Cols {
		mapColTagsFirstTable[col.GetLabel()] = true
	}

	for _, col := range secondTable.Cols {
		if _, ok := mapColTagsFirstTable[col.GetLabel()]; !ok {
			return nil, cm.Errorf(cm.Internal, nil, "Can't merge two tables because col tag[%s] not found in first table", col.GetLabel())
		}
	}

	newTable := firstTable.Clone(firstTable.Label)

	firstData := make([]smry.Cell, len(firstTable.Data))
	secondData := make([]smry.Cell, len(secondTable.Data))
	secondRows := make([]smry.Subject, len(secondTable.Rows))

	copy(firstData, firstTable.Data)
	copy(secondData, secondTable.Data)
	copy(secondRows, secondTable.Rows)

	newTable.Rows = append(newTable.Rows, secondRows...)

	sizeNewTableRows := numOfRowsFirstTable + numOfRowsSecondTable
	sizeNewTableData := len(firstData) + len(secondData)
	newTableData := make([]smry.Cell, 0, sizeNewTableData)
	for i := 0; i < sizeNewTableData/sizeNewTableRows; i++ {
		for j := 0; j < numOfRowsFirstTable; j++ {
			newTableData = append(newTableData, firstData[i*numOfRowsFirstTable+j])
		}
		for j := 0; j < numOfRowsSecondTable; j++ {
			newTableData = append(newTableData, secondData[i*numOfRowsSecondTable+j])
		}
	}
	newTable.Data = newTableData

	return newTable, nil
}
