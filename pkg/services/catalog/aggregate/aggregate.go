package aggregate

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
)

var _ catalog.Aggregate = &Aggregate{}

type Aggregate struct {
}

func (a *Aggregate) CreateShopProduct(ctx context.Context, args *catalog.CreateShopProductArgs) (*catalog.ShopProduct, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopProductInfo(ctx context.Context, args *catalog.UpdateShopProductInfoArgs) (*catalog.ShopProduct, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopProductStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (*catalog.ShopProduct, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopProductImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopProduct, error) {
	return nil, cm.ErrTODO
}

func (a *Aggregate) DeleteShopProducts(ctx context.Context, args *catalog.IDsShopArgs) (*meta.Empty, error) {
	panic("TODO")
}

func (a *Aggregate) CreateShopVariant(ctx context.Context, args *catalog.CreateShopVariantArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopVariantInfo(ctx context.Context, args *catalog.UpdateShopVariantInfoArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}

func (a *Aggregate) DeleteShopVariants(ctx context.Context, args *catalog.IDsShopArgs) (*meta.Empty, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopVariantStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopVariantImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}
