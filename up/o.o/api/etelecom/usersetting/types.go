package usersetting

import (
	"time"

	"o.o/api/top/types/etc/charge_type"
	"o.o/capi/dot"
)

type UserSetting struct {
	// userID
	ID                  dot.ID
	ExtensionChargeType charge_type.ChargeType
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
