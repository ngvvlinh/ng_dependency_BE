package suppliering

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	dot "etop.vn/capi/dot"
)

// +gen:event:topic=event/suppliering

type ShopSupplier struct {
	ID                dot.ID
	ShopID            dot.ID
	FullName          string
	Phone             string
	Code              string
	CodeNorm          int
	Email             string
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
	Status            status3.Status
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type VariantSupplierDeletedEvent struct {
	meta.EventMeta

	SupplierID dot.ID
	ShopID     dot.ID
}
