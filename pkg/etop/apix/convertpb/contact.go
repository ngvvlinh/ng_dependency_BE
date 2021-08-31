package convertpb

import (
	"o.o/api/main/contact"
	typesx "o.o/api/top/external/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_Contact_To_apix_Contact(in *contact.Contact) *typesx.Contact {
	if in == nil {
		return nil
	}
	res := &typesx.Contact{
		ID:        in.ID,
		ShopID:    in.ShopID,
		FullName:  in.FullName,
		Phone:     in.Phone,
		CreatedAt: cmapi.PbTime(in.CreatedAt),
		UpdatedAt: cmapi.PbTime(in.UpdatedAt),
	}

	return res
}

func Convert_core_Contacts_To_apix_Contacts(items []*contact.Contact) []*typesx.Contact {
	result := make([]*typesx.Contact, len(items))
	for i, item := range items {
		result[i] = Convert_core_Contact_To_apix_Contact(item)
	}
	return result
}
