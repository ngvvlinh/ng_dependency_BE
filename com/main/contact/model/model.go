package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type Contact struct {
	ID           dot.ID `paging:"id"`
	ShopID       dot.ID
	FullName     string
	FullNameNorm string
	Phone        string
	PhoneNorm    string
	WLPartnerID  dot.ID

	CreatedAt time.Time `sq:"create" paging:"created_at"`
	UpdatedAt time.Time `sq:"update" paging:"updated_at"`
	DeletedAt time.Time
}
