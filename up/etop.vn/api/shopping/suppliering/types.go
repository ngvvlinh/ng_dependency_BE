package suppliering

import (
	"time"

	"etop.vn/api/meta"
)

// +gen:event:topic=event/suppliering

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

type VariantSupplierDeletedEvent struct {
	meta.EventMeta

	SupplierID int64
	ShopID     int64
}
