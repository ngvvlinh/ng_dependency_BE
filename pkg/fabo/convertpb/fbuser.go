package convertpb

import (
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/fabo"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
)

func PbFbUserCombined(m *fbusering.FbExternalUserCombined) *fabo.FbUserCombined {
	if m == nil {
		return nil
	}
	return &fabo.FbUserCombined{
		ExternalID: m.FbExternalUser.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.FbExternalUser.ExternalInfo.Name,
			FirstName: m.FbExternalUser.ExternalInfo.FirstName,
			LastName:  m.FbExternalUser.ExternalInfo.LastName,
			ShortName: m.FbExternalUser.ExternalInfo.ShortName,
			ImageURL:  m.FbExternalUser.ExternalInfo.ImageURL,
		},
		Status:    m.FbExternalUser.Status,
		CreatedAt: m.FbExternalUser.CreatedAt,
		UpdatedAt: m.FbExternalUser.UpdatedAt,
	}
}

func PbFbUserCombineds(ms []*fbusering.FbExternalUserCombined) []*fabo.FbUserCombined {
	res := make([]*fabo.FbUserCombined, len(ms))
	for i, m := range ms {
		res[i] = PbFbUserCombined(m)
	}
	return res
}

func PbFbUserWithCustomer(m *fbusering.FbExternalUser, c *customering.ShopCustomer) *fabo.FbUserWithCustomer {
	if m == nil {
		return nil
	}
	var result = &fabo.FbUserWithCustomer{
		ExternalID: m.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.ExternalInfo.Name,
			FirstName: m.ExternalInfo.FirstName,
			LastName:  m.ExternalInfo.LastName,
			ShortName: m.ExternalInfo.ShortName,
			ImageURL:  m.ExternalInfo.ImageURL,
		},
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if c != nil {
		result.CustomerID = c.ID
		result.Customer = convertpb.PbShopCustomer(c)
	}
	return result
}

func PbFbUser(m *fbusering.FbExternalUserWithCustomer) *fabo.FbUserWithCustomer {
	if m == nil {
		return nil
	}
	var customerID dot.ID
	if m.ShopCustomer != nil {
		customerID = m.ShopCustomer.ID
	}
	return &fabo.FbUserWithCustomer{
		ExternalID: m.FbExternalUser.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      m.ExternalInfo.Name,
			FirstName: m.ExternalInfo.FirstName,
			LastName:  m.ExternalInfo.LastName,
			ShortName: m.ExternalInfo.ShortName,
			ImageURL:  m.ExternalInfo.ImageURL,
		},
		Status:     m.FbExternalUser.Status,
		CreatedAt:  m.FbExternalUser.CreatedAt,
		UpdatedAt:  m.FbExternalUser.UpdatedAt,
		CustomerID: customerID,
		Customer:   convertpb.PbShopCustomer(m.ShopCustomer),
	}
}

func PbFbUsers(ms []*fbusering.FbExternalUserWithCustomer) []*fabo.FbUserWithCustomer {
	res := make([]*fabo.FbUserWithCustomer, len(ms))
	for i, m := range ms {
		res[i] = PbFbUser(m)
	}
	return res
}
