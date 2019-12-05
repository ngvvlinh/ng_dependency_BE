package shop

import (
	"context"

	"etop.vn/api/top/int/shop"

	"etop.vn/api/main/identity"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandler("api", accountService.UpdateExternalAccountAhamoveVerificationImages)
	bus.AddHandler("api", productSourceService.GetShopProductSources)
	bus.AddHandler("api", productSourceService.CreateProductSource)

	bus.AddHandler("api", collectionService.GetCollection)
	bus.AddHandler("api", collectionService.GetCollections)

}

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

	query := &model.GetUserByIDQuery{
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
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}
