package pgrid

import (
	"strconv"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/pgevent"
)

var _ = selFoo(&Foo{})

// Foo is used for testing only
type Foo struct {
	ID     int64     `sel:" f.id"`
	ShopID string    `sel:" f.shop_id"`
	Time   time.Time `sel:"hf._time"`
}

type Meta struct {
	RID string `json:"rid"`
	Op  string `json:"_op"`
	Env string `json:"_env"`
}

// Implement IModel interface
func (m Meta) _meta() Meta { return m }

func ToMeta(e *pgevent.PgEvent) Meta {
	return Meta{
		RID: strconv.FormatInt(e.RID, 10),
		Op:  string(e.Op),
		Env: cm.Env(),
	}
}

func ignoreNotFound(ok bool, err error) error {
	return err
}

type IdentifyType string

type Identifier interface {
	identifyType() IdentifyType
}

func (t IdentifyType) identifyType() IdentifyType {
	return t
}

const (
	TypeUser     IdentifyType = "user"
	TypeShop     IdentifyType = "shop"
	TypeSupplier IdentifyType = "supplier"
)

type UserData struct {
	IdentifyType `json:"type"`

	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	ShortName string `json:"short_name,omitempty"`

	LatestShopID        string `json:"latest_shop_id,omitempty"`
	LatestShopName      string `json:"latest_shop_name,omitempty"`
	LatestShopCreatedAt int64  `json:"latest_shop_created_at,omitempty"`

	LatestSupplierID        string `json:"latest_supplier_id,omitempty"`
	LatestSupplierName      string `json:"latest_supplier_name,omitempty"`
	LatestSupplierCreatedAt int64  `json:"latest_supplier_created_at,omitempty"`
}

type IdentifyData struct {
	IdentifyType `json:"type"`

	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt int64  `json:"created_at"`
	Name      string `json:"name"`

	OwnerID        string `json:"owner_id"`
	OwnerEmail     string `json:"owner_email"`
	OwnerPhone     string `json:"owner_phone"`
	OwnerFullName  string `json:"owner_full_name"`
	OwnerShortName string `json:"owner_short_name"`
	OwnerCreatedAt int64  `json:"owner_created_at"`
}
