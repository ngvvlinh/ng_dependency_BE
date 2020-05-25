package shop

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
)

// deprecated
func (s *ProductSourceService) CreateProductSource(ctx context.Context, q *CreateProductSourceEndpoint) error {
	q.Result = &shop.ProductSource{
		Id:     q.Context.Shop.ID,
		Status: 1,
	}
	return nil
}

// deprecated: 2018.07.24+14
func (s *ProductSourceService) GetShopProductSources(ctx context.Context, q *GetShopProductSourcesEndpoint) error {
	q.Result = &shop.ProductSourcesResponse{
		ProductSources: []*shop.ProductSource{
			{
				Id:     q.Context.Shop.ID,
				Status: 1,
			},
		},
	}
	return nil
}

// deprecated
func (s *AccountService) UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, r *UpdateExternalAccountAhamoveVerificationImagesEndpoint) error {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return err
	}

	query := &identitymodelx.GetUserByIDQuery{
		UserID: r.Context.Shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.UpdateExternalAccountAhamoveVerificationCommand{
		OwnerID:             user.ID,
		Phone:               phone,
		IDCardFrontImg:      r.IdCardFrontImg,
		IDCardBackImg:       r.IdCardBackImg,
		PortraitImg:         r.PortraitImg,
		WebsiteURL:          r.WebsiteUrl,
		FanpageURL:          r.FanpageUrl,
		CompanyImgs:         r.CompanyImgs,
		BusinessLicenseImgs: r.BusinessLicenseImgs,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}
