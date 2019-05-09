package admin

import (
	etopP "etop.vn/backend/pb/etop"
	"etop.vn/backend/pkg/etop/model"
)

func (m *CreatePartnerRequest) ToModel() *model.Partner {
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
		ContactPersons: etopP.ContactPersonsToModel(p.ContactPersons),
	}
}
