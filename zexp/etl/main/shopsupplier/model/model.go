package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopSupplier struct {
	ID                dot.ID
	ShopID            dot.ID
	FullName          string
	Phone             string
	Email             string
	Code              string
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
	Status            int
	CreatedAt         time.Time
	UpdatedAt         time.Time

	Rid dot.ID
}
