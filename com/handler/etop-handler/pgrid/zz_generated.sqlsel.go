// +build !generator

// Code generated by generator sqlsel. DO NOT EDIT.

package pgrid

import (
	"database/sql"

	core "o.o/backend/pkg/common/sql/sq/core"
)

type Foos []*Foo

func (m *Foo) SQLTableName() string  { return "" }
func (m *Foos) SQLTableName() string { return "" }

func (m *Foo) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.ShopID),
		(*core.Time)(&m.Time),
	}
}

func (m *Foo) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Foos) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Foos, 0, 128)
	for rows.Next() {
		m := new(Foo)
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

func (_ *Foo) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT f.id, f.shop_id, hf._time`)
	return nil
}

func (_ *Foos) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT f.id, f.shop_id, hf._time`)
	return nil
}

type FulfillmentEvents []*FulfillmentEvent

func (m *FulfillmentEvent) SQLTableName() string  { return "" }
func (m *FulfillmentEvents) SQLTableName() string { return "" }

func (m *FulfillmentEvent) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Time)(&m.Time),
		(*core.String)(&m.ID),
		(*core.String)(&m.ShopID),
		(*core.String)(&m.OrderID),
		(*core.String)(&m.OrderCode),
		(*core.String)(&m.ShippingCode),
		(*core.String)(&m.ShippingState),
	}
}

func (m *FulfillmentEvent) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *FulfillmentEvents) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(FulfillmentEvents, 0, 128)
	for rows.Next() {
		m := new(FulfillmentEvent)
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

func (_ *FulfillmentEvent) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, f.id, f.shop_id, o.id, o.code, f.shipping_code, h.shipping_state`)
	return nil
}

func (_ *FulfillmentEvents) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, f.id, f.shop_id, o.id, o.code, f.shipping_code, h.shipping_state`)
	return nil
}

type OrderEvents []*OrderEvent

func (m *OrderEvent) SQLTableName() string  { return "" }
func (m *OrderEvents) SQLTableName() string { return "" }

func (m *OrderEvent) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Time)(&m.Time),
		(*core.Time)(&m.CreatedAt),
		(*core.String)(&m.ID),
		(*core.String)(&m.ShopID),
	}
}

func (m *OrderEvent) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *OrderEvents) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(OrderEvents, 0, 128)
	for rows.Next() {
		m := new(OrderEvent)
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

func (_ *OrderEvent) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, o.created_at, o.id, o.shop_id`)
	return nil
}

func (_ *OrderEvents) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, o.created_at, o.id, o.shop_id`)
	return nil
}

type ShopEvents []*ShopEvent

func (m *ShopEvent) SQLTableName() string  { return "" }
func (m *ShopEvents) SQLTableName() string { return "" }

func (m *ShopEvent) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Time)(&m.Time),
		(*core.Time)(&m.CreatedAt),
		(*core.String)(&m.ID),
		(*core.String)(&m.Name),
		(*core.String)(&m.ImageURL),
		(*core.String)(&m.OwnerID),
		(*core.String)(&m.OwnerFullName),
		(*core.String)(&m.OwnerShortName),
		(*core.String)(&m.OwnerEmail),
		(*core.String)(&m.OwnerPhone),
		(*core.Time)(&m.OwnerCreatedAt),
	}
}

func (m *ShopEvent) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopEvents) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopEvents, 0, 128)
	for rows.Next() {
		m := new(ShopEvent)
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

func (_ *ShopEvent) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, s.created_at, s.id, s.name, s.image_url, u.id, u.full_name, u.short_name, u.email, u.phone, u.created_at`)
	return nil
}

func (_ *ShopEvents) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, s.created_at, s.id, s.name, s.image_url, u.id, u.full_name, u.short_name, u.email, u.phone, u.created_at`)
	return nil
}

type ShopProductEvents []*ShopProductEvent

func (m *ShopProductEvent) SQLTableName() string  { return "" }
func (m *ShopProductEvents) SQLTableName() string { return "" }

func (m *ShopProductEvent) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Time)(&m.Time),
		(*core.Time)(&m.CreatedAt),
		(*core.String)(&m.ShopID),
		(*core.String)(&m.ProductID),
	}
}

func (m *ShopProductEvent) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopProductEvents) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopProductEvents, 0, 128)
	for rows.Next() {
		m := new(ShopProductEvent)
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

func (_ *ShopProductEvent) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, sp.created_at, sp.shop_id, sp.product_id`)
	return nil
}

func (_ *ShopProductEvents) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, sp.created_at, sp.shop_id, sp.product_id`)
	return nil
}

type UserEvents []*UserEvent

func (m *UserEvent) SQLTableName() string  { return "" }
func (m *UserEvents) SQLTableName() string { return "" }

func (m *UserEvent) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Time)(&m.Time),
		(*core.Time)(&m.CreatedAt),
		(*core.String)(&m.ID),
		(*core.String)(&m.Phone),
		(*core.String)(&m.Email),
		(*core.String)(&m.FullName),
		(*core.String)(&m.ShortName),
		(*core.Time)(&m.EmailVerified),
		(*core.Time)(&m.PhoneVerifiedAt),
		(*core.Time)(&m.AgreedEmailInfoAt),
	}
}

func (m *UserEvent) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *UserEvents) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(UserEvents, 0, 128)
	for rows.Next() {
		m := new(UserEvent)
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

func (_ *UserEvent) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, u.created_at, u.id, u.phone, u.email, u.full_name, u.short_name, u.email_verified_at IS NOT NULL, u.phone_verified_at IS NOT NULL, u.agreed_email_info_at IS NOT NULL`)
	return nil
}

func (_ *UserEvents) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT h._time, u.created_at, u.id, u.phone, u.email, u.full_name, u.short_name, u.email_verified_at IS NOT NULL, u.phone_verified_at IS NOT NULL, u.agreed_email_info_at IS NOT NULL`)
	return nil
}