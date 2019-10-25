package shop

import (
	"context"

	"etop.vn/api/main/identity"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandler("api", s.UpdateExternalAccountAhamoveVerificationImages)
	bus.AddHandler("api", s.GetShopProductSources)
	bus.AddHandler("api", s.CreateProductSource)

	bus.AddHandler("api", s.GetCollection)
	bus.AddHandler("api", s.GetCollections)

}

// deprecated
func (s *Service) CreateProductSource(ctx context.Context, q *wrapshop.CreateProductSourceEndpoint) error {
	q.Result = &pbshop.ProductSource{
		Id:     q.Context.Shop.ID,
		Status: 1,
	}
	return nil
}

// deprecated: 2018.07.24+14
func (s *Service) GetShopProductSources(ctx context.Context, q *wrapshop.GetShopProductSourcesEndpoint) error {
	q.Result = &pbshop.ProductSourcesResponse{
		ProductSources: []*pbshop.ProductSource{
			{
				Id:     q.Context.Shop.ID,
				Status: 1,
			},
		},
	}
	return nil
}

// deprecated
func (s *Service) UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, r *wrapshop.UpdateExternalAccountAhamoveVerificationImagesEndpoint) error {
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

// deprecated
func (s *Service) AddProducts(ctx context.Context, q *wrapshop.AddProductsEndpoint) error {
	return cm.ErrREMOVED
}
