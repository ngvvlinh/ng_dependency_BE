// +build !generator

// Code generated by generator sqlsel. DO NOT EDIT.

package sqlstore

import (
	"database/sql"

	core "o.o/backend/pkg/common/sql/sq/core"
)

type FfmByAreas []*FfmByArea

func (m *FfmByArea) SQLTableName() string  { return "" }
func (m *FfmByAreas) SQLTableName() string { return "" }

func (m *FfmByArea) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.Count),
		(*core.String)(&m.ProvinceCode),
		(*core.String)(&m.DistrictCode),
	}
}

func (m *FfmByArea) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *FfmByAreas) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(FfmByAreas, 0, 128)
	for rows.Next() {
		m := new(FfmByArea)
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

func (_ *FfmByArea) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT count(id), address_to_province_code, address_to_district_code`)
	return nil
}

func (_ *FfmByAreas) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT count(id), address_to_province_code, address_to_district_code`)
	return nil
}

type StaffOrders []*StaffOrder

func (m *StaffOrder) SQLTableName() string  { return "" }
func (m *StaffOrders) SQLTableName() string { return "" }

func (m *StaffOrder) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.String)(&m.UserName),
		&m.UserID,
		(*core.Int64)(&m.TotalCount),
		(*core.Int64)(&m.TotalAmount),
	}
}

func (m *StaffOrder) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *StaffOrders) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(StaffOrders, 0, 128)
	for rows.Next() {
		m := new(StaffOrder)
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

func (_ *StaffOrder) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT u.full_name, u.id, count(o.id) as total_amount, sum(o.total_amount) as order_count`)
	return nil
}

func (_ *StaffOrders) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT u.full_name, u.id, count(o.id) as total_amount, sum(o.total_amount) as order_count`)
	return nil
}

type TopSellItems []*TopSellItem

func (m *TopSellItem) SQLTableName() string  { return "" }
func (m *TopSellItems) SQLTableName() string { return "" }

func (m *TopSellItem) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.String)(&m.ProductCode),
		&m.ProductId,
		(*core.String)(&m.Name),
		(*core.Int64)(&m.Count),
		core.Array{&m.ImageUrls, opts},
	}
}

func (m *TopSellItem) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *TopSellItems) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(TopSellItems, 0, 128)
	for rows.Next() {
		m := new(TopSellItem)
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

func (_ *TopSellItem) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT sp.code, ol.product_id, sp.name, SUM(quantity) as sum, sp.image_urls`)
	return nil
}

func (_ *TopSellItems) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT sp.code, ol.product_id, sp.name, SUM(quantity) as sum, sp.image_urls`)
	return nil
}

type Totals []*Total

func (m *Total) SQLTableName() string  { return "" }
func (m *Totals) SQLTableName() string { return "" }

func (m *Total) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.TotalAmount),
		(*core.Int64)(&m.TotalOrder),
		(*core.Float64)(&m.AverageOrder),
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

func (_ *Total) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT SUM(total_amount), COUNT(id), AVG(total_amount)`)
	return nil
}

func (_ *Totals) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT SUM(total_amount), COUNT(id), AVG(total_amount)`)
	return nil
}

type TotalPerDates []*TotalPerDate

func (m *TotalPerDate) SQLTableName() string  { return "" }
func (m *TotalPerDates) SQLTableName() string { return "" }

func (m *TotalPerDate) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{}
}

func (m *TotalPerDate) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *TotalPerDates) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(TotalPerDates, 0, 128)
	for rows.Next() {
		m := new(TotalPerDate)
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

func (_ *TotalPerDate) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT `)
	return nil
}

func (_ *TotalPerDates) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT `)
	return nil
}
