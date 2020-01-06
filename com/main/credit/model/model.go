package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenCredit(&Credit{})

type Credit struct {
	ID        dot.ID
	Amount    int
	ShopID    dot.ID
	Type      string
	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	PaidAt    time.Time
}

var _ = sqlgenCreditExtended(
	&CreditExtended{}, &Credit{}, "c",
	sq.LEFT_JOIN, &identitymodel.Shop{}, "s", "s.id = c.shop_id",
)

type CreditExtended struct {
	*Credit
	Shop *identitymodel.Shop
}
