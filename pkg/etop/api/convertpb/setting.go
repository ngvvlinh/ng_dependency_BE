package convertpb

import (
	"o.o/api/shopping/setting"
	shoptypes "o.o/api/top/int/shop"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_Setting_to_api_Setting(in *setting.ShopSetting) *shoptypes.ShopSetting {
	if in == nil {
		return nil
	}
	paymentTypeID := shipping_payment_type.NullShippingPaymentType{
		Enum:  in.PaymentTypeID,
		Valid: in.PaymentTypeID != shipping_payment_type.None,
	}
	tryOn := try_on.NullTryOnCode{
		Enum:  in.TryOn,
		Valid: in.TryOn != try_on.Unknown,
	}

	res := &shoptypes.ShopSetting{
		ShopID:          in.ShopID,
		PaymentTypeID:   paymentTypeID,
		ReturnAddressID: in.ReturnAddressID,
		ReturnAddress:   Convert_core_Address_To_api_Address(in.ReturnAddress),
		TryOn:           tryOn,
		ShippingNote:    in.ShippingNote,
		Weight:          in.Weight,
		HideAllComments: in.HideAllComments,
		CreatedAt:       cmapi.PbTime(in.CreatedAt),
		UpdatedAt:       cmapi.PbTime(in.UpdatedAt),
	}
	return res
}
