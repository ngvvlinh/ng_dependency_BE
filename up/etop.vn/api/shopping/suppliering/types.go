package suppliering

import (
	"time"

	"etop.vn/api/meta"
	dot "etop.vn/capi/dot"
)

// +gen:event:topic=event/suppliering

type ShopSupplier struct {
	ID                dot.ID
	ShopID            dot.ID
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

	SupplierID dot.ID
	ShopID     dot.ID
}
