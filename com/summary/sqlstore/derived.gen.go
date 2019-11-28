// Code generated by goderive DO NOT EDIT.

package sqlstore

import (
	sql "database/sql"

	core "etop.vn/backend/pkg/common/sq/core"
)

type SQLWriter = core.SQLWriter

func selTotal(_ ...interface{}) bool { return true }

type Totals []*Total

func (m *Total) SQLTableName() string  { return "" }
func (m *Totals) SQLTableName() string { return "" }

func (m *Total) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		core.JSON{&m.TotalAmount},
		core.JSON{&m.TotalOrder},
		core.JSON{&m.AverageOrder},
	}
}

func (m *Total) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Totals) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Totals, 0, 128)
	for rows.Next() {
		m := new(Total)
		args := m.SQLScanArgs(opts)
		if err := rows.Scan(args...); err != nil {
			return err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

func (_ *Total) SQLSelect(w SQLWriter) error {
	w.WriteRawString(`SELECT SUM(total_amount), COUNT(id), AVG(total_amount)`)
	return nil
}

func (_ *Totals) SQLSelect(w SQLWriter) error {
	w.WriteRawString(`SELECT SUM(total_amount), COUNT(id), AVG(total_amount)`)
	return nil
}
