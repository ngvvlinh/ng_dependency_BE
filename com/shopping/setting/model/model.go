package model

import (
	"time"

	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +convert:type=setting.ShopSetting
// +sqlgen
type ShopSetting struct {
	ShopID          dot.ID
	PaymentTypeID   shipping_payment_type.ShippingPaymentType
	ReturnAddressID dot.ID
	TryOn           try_on.TryOnCode
	ShippingNote    string
	Weight          int
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}
