package aggregate

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
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
		ProductType: args.ProductType,
	}
	if err := a.shopProduct(ctx).CreateShopProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (a *Aggregate) UpdateShopProductInfo(ctx context.Context, args *catalog.UpdateShopProductInfoArgs) (*catalog.ShopProduct, error) {
	productDB, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProductDB()
	if err != nil {
		return nil, err
	}
	updated := convert.UpdateShopProduct(productDB, args)
	if err = a.shopProduct(ctx).UpdateShopProduct(updated); err != nil {
		return nil, err
	}
	productReturn, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProduct()
	return productReturn, err
}

func (a *Aggregate) UpdateShopProductStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (int, error) {
	count, err := a.shopProduct(ctx).ShopID(args.ShopID).IDs(args.IDs...).UpdateStatusShopProducts(args.Status)
	return count, err
}

func (a *Aggregate) UpdateShopProductImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopProduct, error) {
	productDB, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ID).GetShopProduct()
	if err != nil {
		return nil, err
	}
	var newImageURLs []string
	newImageURLs = productDB.ImageURLs
	newImageURLs, err = Patch(newImageURLs, args.Updates)
	if err != nil {
		return nil, err
	}
	productDB.ImageURLs = newImageURLs
	if err = a.shopProduct(ctx).ShopID(productDB.ShopID).ID(productDB.ProductID).UpdateImageShopProduct(productDB); err != nil {
		return nil, err
	}
	return productDB, nil
}

func (a *Aggregate) DeleteShopProducts(ctx context.Context, args *shopping.IDsQueryShopArgs) (int, error) {
	deleted, err := a.shopProduct(ctx).ShopID(args.ShopID).IDs(args.IDs...).SoftDelete()
	return deleted, err
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
	if err = a.shopVariant(ctx).CreateShopVariant(variant); err != nil {
		return nil, err
	}
	return variant, nil
}

func (a *Aggregate) UpdateShopVariantInfo(ctx context.Context, args *catalog.UpdateShopVariantInfoArgs) (*catalog.ShopVariant, error) {
	variantDB, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	updated := convert.UpdateShopVariant(variantDB, args)

	if err = a.shopVariant(ctx).UpdateShopVariant(updated); err != nil {
		return nil, err
	}
	variant, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariant()
	return variant, err
}

func (a *Aggregate) DeleteShopVariants(ctx context.Context, args *shopping.IDsQueryShopArgs) (int, error) {
	deleted, err := a.shopVariant(ctx).ShopID(args.ShopID).IDs(args.IDs...).SoftDelete()
	return deleted, err
}

func (a *Aggregate) UpdateShopVariantStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (int, error) {
	update, err := a.shopVariant(ctx).ShopID(args.ShopID).IDs(args.IDs...).UpdateStatusShopVariant(args.Status)
	return update, err
}

func (a *Aggregate) UpdateShopVariantImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopVariant, error) {
	variantDB, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.ID).GetShopVariant()
	if err != nil {
		return nil, err
	}

	var newImageURLs []string
	newImageURLs = variantDB.ImageURLs
	newImageURLs, err = Patch(newImageURLs, args.Updates)
	if err != nil {
		return nil, err
	}
	variantDB.ImageURLs = newImageURLs
	if err = a.shopVariant(ctx).ShopID(variantDB.ShopID).ID(args.ID).UpdateImageShopVariant(variantDB); err != nil {
		return nil, err
	}
	return variantDB, nil
}

func (a *Aggregate) UpdateShopVariantAttributes(ctx context.Context, args *catalog.UpdateShopVariantAttributes) (*catalog.ShopVariant, error) {
	variantDB, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	if len(args.Attributes) <= 0 {
		return nil, cm.Error(cm.InvalidArgument, "Atributes is empty", nil)
	}
	var attributesUpdate model.ProductAttributes
	for _, value := range args.Attributes {
		attributesUpdate = append(attributesUpdate, &model.ProductAttribute{
			Name:  value.Value,
			Value: value.Name,
		})
	}
	variantDB.Attributes = attributesUpdate
	if err = a.shopVariant(ctx).UpdateShopVariant(variantDB); err != nil {
		return nil, err
	}
	return a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariant()
}

func Patch(source []string, update []*meta.UpdateSet) ([]string, error) {
	if OptionValue(update, meta.OpDeleteAll) != nil {
		return []string{}, nil
	}
	if replace := OptionValue(update, meta.OpReplaceAll); replace != nil {
		arr, _, err := replace.Update(source)
		if err != nil {
			return []string{}, err
		}
		return arr, nil
	}
	var arrReturn []string
	arrReturn = source
	if add := OptionValue(update, meta.OpAdd); add != nil {
		var err error
		arrReturn, _, err = add.Update(arrReturn)
		if err != nil {
			return []string{}, err
		}
	}
	if remove := OptionValue(update, meta.OpRemove); remove != nil {
		var err error
		arrReturn, _, err = remove.Update(arrReturn)
		if err != nil {
			return []string{}, err
		}
	}
	return arrReturn, nil
}

func OptionValue(update []*meta.UpdateSet, op meta.UpdateOp) *meta.UpdateSet {
	for i := 0; i < len(update); i++ {
		if update[i].Op == op {
			return update[i]
		}
	}
	return nil
}
