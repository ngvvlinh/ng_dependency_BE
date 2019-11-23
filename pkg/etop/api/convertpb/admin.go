package convertpb

import (
	"etop.vn/api/pb/etop/admin"
	"etop.vn/backend/pkg/etop/model"
)

func CreatePartnerRequestToModel(m *admin.CreatePartnerRequest) *model.Partner {
	p := m.Partner
	isTest := 0
	if p.IsTest {
		isTest = 1
	}
	return &model.Partner{
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
