package suppliering

import (
	"time"
)

type ShopSupplier struct {
	ID                int64
	ShopID            int64
	FullName          string
	Phone             string
	Code              string
	CodeNorm          int32
	Email             string
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
	Status            int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
