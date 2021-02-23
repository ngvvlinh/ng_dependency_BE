package etelecom

import (
	"o.o/api/etelecom/usersetting"
	etelecomapi "o.o/api/top/int/etelecom"
)

// +gen:convert: o.o/api/etelecom -> o.o/api/top/int/etelecom/types

func Convert_usersetting_UserSetting_api_UserSetting(in *usersetting.UserSetting) *etelecomapi.EtelecomUserSetting {
	if in == nil {
		return nil
	}
	return &etelecomapi.EtelecomUserSetting{
		ID:                  in.ID,
		ExtensionChargeType: in.ExtensionChargeType,
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
	}
}
