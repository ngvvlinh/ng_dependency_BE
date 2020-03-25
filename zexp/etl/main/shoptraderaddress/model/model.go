package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	addressmodel "etop.vn/backend/com/main/address/model"
	"etop.vn/capi/dot"
)

// +sqlgen
type ShopTraderAddress struct {
	ID           dot.ID
	ShopID       dot.ID
	TraderID     dot.ID
	FullName     string
	Phone        string
	Email        string
	Company      string
	Address1     string
	Address2     string
	DistrictCode string
	WardCode     string
	Position     string
	IsDefault    bool
	Coordinates  *addressmodel.Coordinates
	CreatedAt    time.Time
	UpdatedAt    time.Time

	//Default status = 1
	Status status3.Status

	Rid dot.ID
}
