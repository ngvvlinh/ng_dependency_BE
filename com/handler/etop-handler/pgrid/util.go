package pgrid

import (
	"strconv"
	"time"

	"o.o/backend/com/handler/pgevent"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi/dot"
)

// Foo is used for testing only
//
// +sqlsel
type Foo struct {
	ID     dot.ID    `sel:" f.id"`
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
		Op:  e.Op.String(),
		Env: cmenv.Env().String(),
	}
}
