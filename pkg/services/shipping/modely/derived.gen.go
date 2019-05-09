// Code generated by goderive DO NOT EDIT.

package modely

import (
	"database/sql"

	sq "etop.vn/backend/pkg/common/sql"
	core "etop.vn/backend/pkg/common/sql/core"
	model "etop.vn/backend/pkg/etop/model"
	etop_vn_backend_pkg_services_moneytx_model "etop.vn/backend/pkg/services/moneytx/model"
)

type SQLWriter = core.SQLWriter

// Type FulfillmentExtended represents a join
func sqlgenFulfillmentExtended(_ *FulfillmentExtended, _ *model.Fulfillment, as sq.AS, t0 sq.JOIN_TYPE, _ *model.Shop, a0 sq.AS, c0 string, t1 sq.JOIN_TYPE, _ *model.Order, a1 sq.AS, c1 string, t2 sq.JOIN_TYPE, _ *etop_vn_backend_pkg_services_moneytx_model.MoneyTransactionShipping, a2 sq.AS, c2 string) bool {
	__sqlFulfillmentExtended_JoinTypes = []sq.JOIN_TYPE{t0, t1, t2}
	__sqlFulfillmentExtended_As = as
	__sqlFulfillmentExtended_JoinAs = []sq.AS{a0, a1, a2}
	__sqlFulfillmentExtended_JoinConds = []string{c0, c1, c2}
	return true
}

type FulfillmentExtendeds []*FulfillmentExtended

var __sqlFulfillmentExtended_JoinTypes []sq.JOIN_TYPE
var __sqlFulfillmentExtended_As sq.AS
var __sqlFulfillmentExtended_JoinAs []sq.AS
var __sqlFulfillmentExtended_JoinConds []string

func (m *FulfillmentExtended) SQLTableName() string  { return "fulfillment" }
func (m *FulfillmentExtendeds) SQLTableName() string { return "fulfillment" }

func (m *FulfillmentExtended) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *FulfillmentExtendeds) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(FulfillmentExtendeds, 0, 128)
	for rows.Next() {
		m := new(FulfillmentExtended)
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

func (m *FulfillmentExtended) SQLSelect(w SQLWriter) error {
	(*FulfillmentExtended)(nil).__sqlSelect(w)
	w.WriteByte(' ')
	(*FulfillmentExtended)(nil).__sqlJoin(w, __sqlFulfillmentExtended_JoinTypes)
	return nil
}

func (m *FulfillmentExtendeds) SQLSelect(w SQLWriter) error {
	return (*FulfillmentExtended)(nil).SQLSelect(w)
}

func (m *FulfillmentExtended) SQLJoin(w SQLWriter, types []sq.JOIN_TYPE) error {
	if len(types) == 0 {
		types = __sqlFulfillmentExtended_JoinTypes
	}
	m.__sqlJoin(w, types)
	return nil
}

func (m *FulfillmentExtendeds) SQLJoin(w SQLWriter, types []sq.JOIN_TYPE) error {
	return (*FulfillmentExtended)(nil).SQLJoin(w, types)
}

func (m *FulfillmentExtended) __sqlSelect(w SQLWriter) {
	w.WriteRawString("SELECT ")
	core.WriteCols(w, string(__sqlFulfillmentExtended_As), (*model.Fulfillment)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, string(__sqlFulfillmentExtended_JoinAs[0]), (*model.Shop)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, string(__sqlFulfillmentExtended_JoinAs[1]), (*model.Order)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, string(__sqlFulfillmentExtended_JoinAs[2]), (*etop_vn_backend_pkg_services_moneytx_model.MoneyTransactionShipping)(nil).SQLListCols())
}

func (m *FulfillmentExtended) __sqlJoin(w SQLWriter, types []sq.JOIN_TYPE) {
	if len(types) != 3 {
		panic("common/sql: expect 3 types to join")
	}
	w.WriteRawString("FROM ")
	w.WriteName("fulfillment")
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlFulfillmentExtended_As))
	w.WriteByte(' ')
	w.WriteRawString(string(types[0]))
	w.WriteRawString(" JOIN ")
	w.WriteName((*model.Shop)(nil).SQLTableName())
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlFulfillmentExtended_JoinAs[0]))
	w.WriteRawString(" ON ")
	w.WriteQueryString(__sqlFulfillmentExtended_JoinConds[0])
	w.WriteByte(' ')
	w.WriteRawString(string(types[1]))
	w.WriteRawString(" JOIN ")
	w.WriteName((*model.Order)(nil).SQLTableName())
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlFulfillmentExtended_JoinAs[1]))
	w.WriteRawString(" ON ")
	w.WriteQueryString(__sqlFulfillmentExtended_JoinConds[1])
	w.WriteByte(' ')
	w.WriteRawString(string(types[2]))
	w.WriteRawString(" JOIN ")
	w.WriteName((*etop_vn_backend_pkg_services_moneytx_model.MoneyTransactionShipping)(nil).SQLTableName())
	w.WriteRawString(" AS ")
	w.WriteRawString(string(__sqlFulfillmentExtended_JoinAs[2]))
	w.WriteRawString(" ON ")
	w.WriteQueryString(__sqlFulfillmentExtended_JoinConds[2])
}

func (m *FulfillmentExtended) SQLScanArgs(opts core.Opts) []interface{} {
	args := make([]interface{}, 0, 64) // TODO: pre-calculate length
	m.Fulfillment = new(model.Fulfillment)
	args = append(args, m.Fulfillment.SQLScanArgs(opts)...)
	m.Shop = new(model.Shop)
	args = append(args, m.Shop.SQLScanArgs(opts)...)
	m.Order = new(model.Order)
	args = append(args, m.Order.SQLScanArgs(opts)...)
	m.MoneyTransactionShipping = new(etop_vn_backend_pkg_services_moneytx_model.MoneyTransactionShipping)
	args = append(args, m.MoneyTransactionShipping.SQLScanArgs(opts)...)

	return args
}
