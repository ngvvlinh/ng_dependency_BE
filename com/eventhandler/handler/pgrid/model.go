package pgrid

import (
	"time"

	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/pkg/common/sql/cmsql"
)

// We must always use string for id, shop_id, order_id, etc.

type IModel interface {
	Query(*cmsql.Database, *pgevent.PgEvent) (bool, error)
	_meta() Meta
}

// +sqlsel
type UserEvent struct {
	Meta
	Time      time.Time `json:"_time"      sel:"h._time"`
	CreatedAt time.Time `json:"created_at" sel:"u.created_at"`

	ID    string `json:"id"    sel:"u.id"`
	Phone string `json:"phone" sel:"u.phone"`
	Email string `json:"email" sel:"u.email"`

	FullName  string `json:"full_name"  sel:"u.full_name"`
	ShortName string `json:"short_name" sel:"u.short_name"`

	EmailVerified     time.Time `json:"email_verified"    sel:"u.email_verified_at    IS NOT NULL"`
	PhoneVerifiedAt   time.Time `json:"phone_verified"    sel:"u.phone_verified_at    IS NOT NULL"`
	AgreedEmailInfoAt time.Time `json:"agreed_email_info" sel:"u.agreed_email_info_at IS NOT NULL"`
}

func (m *UserEvent) Query(db *cmsql.Database, event *pgevent.PgEvent) (bool, error) {
	sql := `
FROM "user" AS u
JOIN history."user" AS h
  ON h.id = u.id
`
	m.Meta = ToMeta(event)
	return db.SQL(sql).Where("h.rid = ?", event.RID).Get(m)
}

// +sqlsel
type FulfillmentEvent struct {
	Meta
	Time          time.Time `json:"_time"          sel:"h._time"`
	ID            string    `json:"id"             sel:"f.id"`
	ShopID        string    `json:"shop_id"        sel:"f.shop_id"`
	OrderID       string    `json:"order_id"       sel:"o.id"`
	OrderCode     string    `json:"order_code"     sel:"o.code"`
	ShippingCode  string    `json:"shipping_code"  sel:"f.shipping_code"`
	ShippingState string    `json:"shipping_state" sel:"h.shipping_state"`
}

func (m *FulfillmentEvent) Query(db *cmsql.Database, event *pgevent.PgEvent) (bool, error) {
	sql := `
FROM fulfillment AS f
JOIN history.fulfillment AS h
  ON h.id = f.id
JOIN "order" AS o
  ON f.order_id = o.id
`
	m.Meta = ToMeta(event)
	return db.SQL(sql).Where("h.rid = ?", event.RID).Get(m)
}

// +sqlsel
type ShopEvent struct {
	Meta
	Time      time.Time `json:"_time"      sel:"h._time"`
	CreatedAt time.Time `json:"created_at" sel:"s.created_at"`

	ID       string `json:"id"        sel:"s.id"`
	Name     string `json:"name"      sel:"s.name"`
	ImageURL string `json:"image_url" sel:"s.image_url"`

	OwnerID        string    `json:"owner_id"    sel:"u.id"`
	OwnerFullName  string    `json:"owner_name"  sel:"u.full_name"`
	OwnerShortName string    `json:"owner_name"  sel:"u.short_name"`
	OwnerEmail     string    `json:"owner_email" sel:"u.email"`
	OwnerPhone     string    `json:"owner_phone" sel:"u.phone"`
	OwnerCreatedAt time.Time `json:"owner_created_at" sel:"u.created_at"`
}

func (m *ShopEvent) Query(db *cmsql.Database, event *pgevent.PgEvent) (bool, error) {
	sql := `
FROM shop AS s
JOIN history.shop AS h
  ON h.id = s.id
JOIN "user" AS u
  ON s.owner_id = u.id
`
	m.Meta = ToMeta(event)
	return db.SQL(sql).Where("h.rid = ?", event.RID).Get(m)
}

// +sqlsel
type ShopProductEvent struct {
	Meta
	Time      time.Time `json:"_time"      sel:"h._time"`
	CreatedAt time.Time `json:"created_at" sel:"sp.created_at"`

	ShopID    string `json:"shop_id"    sel:"sp.shop_id"`
	ProductID string `json:"product_id" sel:"sp.product_id"`
}

func (m *ShopProductEvent) Query(db *cmsql.Database, event *pgevent.PgEvent) (bool, error) {
	sql := `
FROM shop_product AS sp
JOIN history.shop_product AS h
  ON h.shop_id    = sp.shop_id
 AND h.product_id = sp.product_id
`
	m.Meta = ToMeta(event)
	return db.SQL(sql).Where("h.rid = ?", event.RID).Get(m)
}

// +sqlsel
type OrderEvent struct {
	Meta
	Time      time.Time `json:"_time"      sel:"h._time"`
	CreatedAt time.Time `json:"created_at" sel:"o.created_at"`

	ID     string `json:"id"      sel:"o.id"`
	ShopID string `json:"shop_id" sel:"o.shop_id"`
}

func (m *OrderEvent) Query(db *cmsql.Database, event *pgevent.PgEvent) (bool, error) {
	sql := `
FROM "order" AS o
JOIN history."order" AS h
  ON h.id = o.id
`
	m.Meta = ToMeta(event)
	return db.SQL(sql).Where("h.rid = ?", event.RID).Get(m)
}
