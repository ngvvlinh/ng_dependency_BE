package aggregate

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ catalog.Aggregate = &Aggregate{}

type Aggregate struct {
	shopProduct sqlstore.ShopProductStoreFactory
	shopVariant sqlstore.ShopVariantStoreFactory
}

func New(db cmsql.Database) *Aggregate {
	return &Aggregate{
		shopProduct: sqlstore.NewShopProductStore(db),
		shopVariant: sqlstore.NewShopVariantStore(db),
	}
}

func (a *Aggregate) MessageBus() catalog.CommandBus {
	b := bus.New()
	return catalog.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShopProduct(ctx context.Context, args *catalog.CreateShopProductArgs) (*catalog.ShopProduct, error) {
	product := &catalog.ShopProduct{
		ProductID: cm.NewID(),
		ShopID:    args.ShopID,
		Code:      args.Code,
		Name:      args.Name,
		Unit:      args.Unit,
		ImageURLs: args.ImageURLs,
		Note:      args.Note,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   args.ShortDesc,
			Description: args.Description,
			DescHTML:    args.DescHTML,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   args.CostPrice,
			ListPrice:   args.ListPrice,
			RetailPrice: args.RetailPrice,
		},
	}
	if err := a.shopProduct(ctx).CreateShopProduct(product); err != nil {
		return nil, err
	}
	return product, nil
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

func (a *Aggregate) DeleteShopProducts(ctx context.Context, args *shopping.IDsQueryShopArgs) (*meta.Empty, error) {
	panic("TODO")
}

func (a *Aggregate) CreateShopVariant(ctx context.Context, args *catalog.CreateShopVariantArgs) (*catalog.ShopVariant, error) {
	_, err := a.shopProduct(ctx).
		ShopID(args.ShopID).
		ID(args.ProductID).
		GetShopProductDB()
	if err != nil {
		return nil, err
	}

	variant := &catalog.ShopVariant{
		ShopID:    args.ShopID,
		ProductID: args.ProductID,
		VariantID: cm.NewID(),
		Code:      args.Code,
		Name:      args.Name,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   args.ShortDesc,
			Description: args.Description,
			DescHTML:    args.DescHTML,
		},
		ImageURLs:  args.ImageURLs,
		Status:     0,
		Attributes: args.Attributes,
		PriceInfo: catalog.PriceInfo{
			CostPrice:   args.CostPrice,
			ListPrice:   args.ListPrice,
			RetailPrice: args.RetailPrice,
		},
		Note: args.Note,
	}
	if err := a.shopVariant(ctx).CreateShopVariant(variant); err != nil {
		return nil, err
	}
	return variant, nil
}

func (a *Aggregate) UpdateShopVariantInfo(ctx context.Context, args *catalog.UpdateShopVariantInfoArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}

func (a *Aggregate) DeleteShopVariants(ctx context.Context, args *shopping.IDsQueryShopArgs) (*meta.Empty, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopVariantStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}

func (a *Aggregate) UpdateShopVariantImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopVariant, error) {
	panic("TODO")
}
