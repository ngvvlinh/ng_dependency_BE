package setting

import (
	"context"

	"o.o/api/main/address"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	UpdateShopSetting(context.Context, *UpdateShopSettingArgs) (*ShopSetting, error)
}

type QueryService interface {
	GetShopSetting(context.Context, *GetShopSettingArgs) (*ShopSetting, error)
}

//-- queries --//
type GetShopSettingArgs struct {
	ShopID dot.ID
}

//-- commands --//

// +convert:update=ShopSetting
type UpdateShopSettingArgs struct {
	ShopID        dot.ID
	ReturnAddress *address.Address
	PaymentTypeID shipping_payment_type.NullShippingPaymentType
	TryOn         try_on.NullTryOnCode
	ShippingNote  dot.NullString
	Weight        dot.NullInt
}
