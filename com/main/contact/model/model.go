package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Contact struct {
	ID          dot.ID
	ShopID      dot.ID
	FullName    string
	Phone       string
	PhoneNorm   string
	WLPartnerID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}
