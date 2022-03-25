package setting

import (
	"time"

	"o.o/api/main/address"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

type ShopSetting struct {
	ShopID                     dot.ID                                    `json:"shop_id"`
	ReturnAddress              *address.Address                          `json:"return_address"`
	ReturnAddressID            dot.ID                                    `json:"return_address_id"`
	PaymentTypeID              shipping_payment_type.ShippingPaymentType `json:"payment_type_id"`
	TryOn                      try_on.TryOnCode                          `json:"try_on"`
	ShippingNote               string                                    `json:"shipping_note"`
	Weight                     int                                       `json:"weight"`
	HideAllComments            dot.NullBool                              `json:"hide_all_comments"`
	CreatedAt                  time.Time                                 `json:"created_at"`
	UpdatedAt                  time.Time                                 `json:"updated_at"`
	AllowConnectDirectShipment bool                                      `json:"allow_connect_direct_shipment"`
}
