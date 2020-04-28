package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	addressmodel "o.o/backend/com/main/address/model"
	"o.o/capi/dot"
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
	Status status3.Status `sql_type:"int2"`

	Rid dot.ID
}
