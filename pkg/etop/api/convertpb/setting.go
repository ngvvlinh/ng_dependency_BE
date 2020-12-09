package convertpb

import (
	"o.o/api/shopping/setting"
	shoptypes "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_Setting_to_api_Setting(in *setting.ShopSetting) *shoptypes.ShopSetting {
	if in == nil {
		return nil
	}
	res := &shoptypes.ShopSetting{
		ShopID:          in.ShopID,
		PaymentTypeID:   in.PaymentTypeID,
		ReturnAddressID: in.ReturnAddressID,
		ReturnAddress:   Convert_core_Address_To_api_Address(in.ReturnAddress),
		TryOn:           in.TryOn,
		ShippingNote:    in.ShippingNote,
		Weight:          in.Weight,
		CreatedAt:       cmapi.PbTime(in.CreatedAt),
		UpdatedAt:       cmapi.PbTime(in.UpdatedAt),
	}
	return res
}
