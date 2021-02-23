package model

import (
	"time"

	"o.o/api/top/types/etc/charge_type"
	"o.o/capi/dot"
)

// +sqlgen
type UserSetting struct {
	ID                  dot.ID
	ExtensionChargeType charge_type.ChargeType
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
}
