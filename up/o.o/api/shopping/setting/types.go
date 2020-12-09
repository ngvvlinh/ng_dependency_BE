package setting

import (
	"time"

	"o.o/api/main/address"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

type ShopSetting struct {
	ShopID          dot.ID
	ReturnAddress   *address.Address
	ReturnAddressID dot.ID
	PaymentTypeID   shipping_payment_type.ShippingPaymentType
	TryOn           try_on.TryOnCode
	ShippingNote    string
	Weight          int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
