package pgrid

import (
	"strconv"
	"time"

	"etop.vn/backend/com/handler/pgevent"
	cm "etop.vn/backend/pkg/common"
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

type IdentifyType string

type Identifier interface {
	identifyType() IdentifyType
}

func (t IdentifyType) identifyType() IdentifyType {
	return t
}

const (
	TypeUser IdentifyType = "user"
	TypeShop IdentifyType = "shop"
)
