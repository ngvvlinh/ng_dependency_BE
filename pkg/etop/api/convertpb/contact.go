package convertpb

import (
	"o.o/api/main/contact"
	shoptypes "o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_Contact_to_api_Contact(in *contact.Contact) *shoptypes.Contact {
	if in == nil {
		return nil
	}
	res := &shoptypes.Contact{
		ID:        in.ID,
		ShopID:    in.ShopID,
		FullName:  in.FullName,
		Phone:     in.Phone,
		CreatedAt: cmapi.PbTime(in.CreatedAt),
		UpdatedAt: cmapi.PbTime(in.UpdatedAt),
	}
	return res
}
