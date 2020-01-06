package convertpb

import (
	admin "etop.vn/api/top/int/admin"
	identitymodel "etop.vn/backend/com/main/identity/model"
)

func CreatePartnerRequestToModel(m *admin.CreatePartnerRequest) *identitymodel.Partner {
	p := m.Partner
	isTest := 0
	if p.IsTest {
		isTest = 1
	}
	return &identitymodel.Partner{
		ID:             0,
		OwnerID:        p.OwnerId,
		Status:         0,
		IsTest:         isTest,
		Name:           p.Name,
		PublicName:     p.PublicName,
		Phone:          p.Phone,
		Email:          p.Email,
		ImageURL:       p.ImageUrl,
		WebsiteURL:     p.WebsiteUrl,
		ContactPersons: ContactPersonsToModel(p.ContactPersons),
	}
}
