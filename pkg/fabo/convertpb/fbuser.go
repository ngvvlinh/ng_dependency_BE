package convertpb

import (
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/fabo"
	convertxmin "o.o/backend/pkg/etop/apix/convertpb/_min"
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
		result.Customer = convertxmin.PbShopCustomer(c)
	}
	result.TagIDS = m.TagIDs
	return result
}

func PbExternalUserWithCustomer(m *fbusering.FbExternalUserWithCustomer) *fabo.FbUserWithCustomer {
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
		Customer:   convertxmin.PbShopCustomer(m.ShopCustomer),
		TagIDS:     m.TagIDs,
	}
}

func PbExternalUsersWithCustomer(ms []*fbusering.FbExternalUserWithCustomer) []*fabo.FbUserWithCustomer {
	res := make([]*fabo.FbUserWithCustomer, len(ms))
	for i, m := range ms {
		res[i] = PbExternalUserWithCustomer(m)
	}
	return res
}

func PbCustomerWithFbUser(customer *fbusering.ShopCustomerWithFbExternalUser) *fabo.CustomerWithFbUserAvatars {
	if customer == nil {
		return nil
	}
	return &fabo.CustomerWithFbUserAvatars{
		Id:           customer.ID,
		ShopId:       customer.ShopID,
		ExternalId:   dot.String(customer.ExternalID),
		ExternalCode: dot.String(customer.ExternalCode),
		FullName:     dot.String(customer.FullName),
		Code:         dot.String(customer.Code),
		Note:         dot.String(customer.Note),
		Phone:        dot.String(customer.Phone),
		Email:        dot.String(customer.Email),
		Gender:       customer.Gender.Wrap(),
		Type:         customer.Type.Wrap(),
		Birthday:     dot.String(customer.Birthday),
		CreatedAt:    dot.Time(customer.CreatedAt),
		UpdatedAt:    dot.Time(customer.UpdatedAt),
		Status:       customer.Status.Wrap(),
		Deleted:      customer.Deleted,
		FbUsers:      PbFbUsers(customer.FbUsers),
	}
}
func PbFbUsers(ms []*fbusering.FbExternalUser) []*fabo.FbUser {
	res := make([]*fabo.FbUser, len(ms))
	for i, m := range ms {
		res[i] = PbFbUser(m)
	}
	return res
}

func PbFbUser(fbUser *fbusering.FbExternalUser) *fabo.FbUser {
	if fbUser == nil {
		return nil
	}
	return &fabo.FbUser{
		ExternalID: fbUser.ExternalID,
		ExternalInfo: &fabo.ExternalFbUserInfo{
			Name:      fbUser.ExternalInfo.Name,
			FirstName: fbUser.ExternalInfo.FirstName,
			LastName:  fbUser.ExternalInfo.LastName,
			ShortName: fbUser.ExternalInfo.ShortName,
			ImageURL:  fbUser.ExternalInfo.ImageURL,
		},
		Status:    fbUser.Status,
		CreatedAt: fbUser.CreatedAt,
		UpdatedAt: fbUser.UpdatedAt,
	}
}

func PbCustomersWithFbUsers(customers []*fbusering.ShopCustomerWithFbExternalUser) []*fabo.CustomerWithFbUserAvatars {
	res := make([]*fabo.CustomerWithFbUserAvatars, len(customers))
	for i, m := range customers {
		res[i] = PbCustomerWithFbUser(m)
	}
	return res
}
