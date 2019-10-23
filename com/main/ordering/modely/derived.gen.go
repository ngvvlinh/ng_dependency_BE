// Code generated by goderive DO NOT EDIT.

package modely

import (
	"database/sql"
	"sync"

	model "etop.vn/backend/com/main/ordering/model"
	etop_vn_backend_com_main_shipping_model "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/cmsql"
	sq "etop.vn/backend/pkg/common/sq"
	core "etop.vn/backend/pkg/common/sq/core"
)

var __sqlModels []interface{ SQLVerifySchema(db *cmsql.Database) }
var __sqlonce sync.Once

func SQLVerifySchema(db *cmsql.Database) {
	__sqlonce.Do(func() {
		for _, m := range __sqlModels {
			m.SQLVerifySchema(db)
		}
	})
}

type SQLWriter = core.SQLWriter

// Type OrderExtended represents a join
func sqlgenOrderExtended(_ *OrderExtended, _ *model.Order, as sq.AS, t0 sq.JOIN_TYPE, _ *etop_vn_backend_com_main_shipping_model.Fulfillment, a0 sq.AS, c0 string) bool {
	__sqlOrderExtended_JoinTypes = []sq.JOIN_TYPE{t0}
	__sqlOrderExtended_As = as
	__sqlOrderExtended_JoinAs = []sq.AS{a0}
	__sqlOrderExtended_JoinConds = []string{c0}
	return true
}

type OrderExtendeds []*OrderExtended

var __sqlOrderExtended_JoinTypes []sq.JOIN_TYPE
var __sqlOrderExtended_As sq.AS
var __sqlOrderExtended_JoinAs []sq.AS
var __sqlOrderExtended_JoinConds []string

func (m *OrderExtended) SQLTableName() string  { return "order" }
func (m *OrderExtendeds) SQLTableName() string { return "order" }

func (m *OrderExtended) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *OrderExtendeds) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(OrderExtendeds, 0, 128)
	for rows.Next() {
		m := new(OrderExtended)
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

func (m *OrderExtended) SQLSelect(w SQLWriter) error {
	(*OrderExtended)(nil).__sqlSelect(w)
	w.WriteByte(' ')
	(*OrderExtended)(nil).__sqlJoin(w, __sqlOrderExtended_JoinTypes)
	return nil
}

func (m *OrderExtendeds) SQLSelect(w SQLWriter) error {
	return (*OrderExtended)(nil).SQLSelect(w)
}

func (m *OrderExtended) SQLJoin(w SQLWriter, types []sq.JOIN_TYPE) error {
	if len(types) == 0 {
		types = __sqlOrderExtended_JoinTypes
	}
	m.__sqlJoin(w, types)
	return nil
}

func (m *OrderExtendeds) SQLJoin(w SQLWriter, types []sq.JOIN_TYPE) error {
	return (*OrderExtended)(nil).SQLJoin(w, types)
}

func (m *OrderExtended) __sqlSelect(w SQLWriter) {
	w.WriteRawString("SELECT ")
	core.WriteCols(w, string(__sqlOrderExtended_As), (*model.Order)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, string(__sqlOrderExtended_JoinAs[0]), (*etop_vn_backend_com_main_shipping_model.Fulfillment)(nil).SQLListCols())
}

func (m *OrderExtended) __sqlJoin(w SQLWriter, types []sq.JOIN_TYPE) {
	if len(types) != 1 {
		panic("common/sql: expect 1 type to join")
	}
	w.WriteRawString("FROM ")
	w.WriteName("order")
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlOrderExtended_As))
	w.WriteByte(' ')
	w.WriteRawString(string(types[0]))
	w.WriteRawString(" JOIN ")
	w.WriteName((*etop_vn_backend_com_main_shipping_model.Fulfillment)(nil).SQLTableName())
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlOrderExtended_JoinAs[0]))
	w.WriteRawString(" ON ")
	w.WriteQueryString(__sqlOrderExtended_JoinConds[0])
}

func (m *OrderExtended) SQLScanArgs(opts core.Opts) []interface{} {
	args := make([]interface{}, 0, 64) // TODO: pre-calculate length
	m.Order = new(model.Order)
	args = append(args, m.Order.SQLScanArgs(opts)...)
	m.Fulfillment = new(etop_vn_backend_com_main_shipping_model.Fulfillment)
	args = append(args, m.Fulfillment.SQLScanArgs(opts)...)

	return args
}
