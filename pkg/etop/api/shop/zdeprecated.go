package shop

import (
	"context"

	"etop.vn/api/main/identity"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandler("api", UpdateExternalAccountAhamoveVerificationImages)
	bus.AddHandler("api", GetShopProductSources)
	bus.AddHandler("api", CreateProductSource)

	bus.AddHandler("api", CreateCollection)
	bus.AddHandler("api", DeleteCollection)
	bus.AddHandler("api", GetCollection)
	bus.AddHandler("api", GetCollections)
	bus.AddHandler("api", GetCollectionsByIDs)
}

// deprecated
func CreateProductSource(ctx context.Context, q *wrapshop.CreateProductSourceEndpoint) error {
	q.Result = &pbshop.ProductSource{
		Id:     q.Context.Shop.ID,
		Status: 1,
	}
	return nil
}

// deprecated: 2018.07.24+14
func GetShopProductSources(ctx context.Context, q *wrapshop.GetShopProductSourcesEndpoint) error {
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
func UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, r *wrapshop.UpdateExternalAccountAhamoveVerificationImagesEndpoint) error {
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
func GetCollection(ctx context.Context, q *wrapshop.GetCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func GetCollectionsByIDs(ctx context.Context, q *wrapshop.GetCollectionsByIDsEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func GetCollections(ctx context.Context, q *wrapshop.GetCollectionsEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func CreateCollection(ctx context.Context, q *wrapshop.CreateCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func DeleteCollection(ctx context.Context, q *wrapshop.DeleteCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func RemoveProductsCollection(ctx context.Context, q *wrapshop.RemoveProductsCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func UpdateCollection(ctx context.Context, q *wrapshop.UpdateCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func UpdateProductsCollection(ctx context.Context, q *wrapshop.UpdateProductsCollectionEndpoint) error {
	return cm.ErrREMOVED
}

// deprecated
func AddProducts(ctx context.Context, q *wrapshop.AddProductsEndpoint) error {
	return cm.ErrREMOVED
}
