package model

import (
	"time"

	"o.o/api/top/types/etc/charge_type"
	"o.o/capi/dot"
)

// +sqlgen
type UserSetting struct {
	ID                  dot.ID `paging:"id"`
	ExtensionChargeType charge_type.ChargeType
	CreatedAt           time.Time `sq:"create" paging:"created_at"`
	UpdatedAt           time.Time `sq:"update"`
}
