package product

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/catalog/types"
	"o.o/api/main/inventory"
	"o.o/api/meta"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	catalogmodelx "o.o/backend/com/main/catalog/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type ProductService struct {
	session.Session

	CatalogQuery   catalog.QueryBus
	CatalogAggr    catalog.CommandBus
	InventoryQuery inventory.QueryBus
}

func (s *ProductService) Clone() api.ProductService { res := *s; return &res }

func (s *ProductService) UpdateVariant(ctx context.Context, q *api.UpdateVariantRequest) (*api.ShopVariant, error) {
	shopID := s.SS.Shop().ID
	var attributes *types.Attributes = nil
	if q.Attributes != nil {
		attributesRequest := types.ValidateAttributesEmpty(q.Attributes)
		attributes = (*types.Attributes)(&attributesRequest)
	}
	cmd := &catalog.UpdateShopVariantInfoCommand{
		ShopID:    shopID,
		VariantID: q.Id,
		Name:      q.Name,
		Code:      q.Code,
		Note:      q.Note,

		ShortDesc:    q.ShortDesc,
		Descripttion: q.Description,
		DescHTML:     q.DescHtml,

		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,
		Attributes:  attributes,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopVariant(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateVariantAttributes(ctx context.Context, q *api.UpdateVariantAttributesRequest) (*api.ShopVariant, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopVariantAttributesCommand{
		ShopID:     shopID,
		VariantID:  q.VariantId,
		Attributes: q.Attributes,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopVariant(cmd.Result)
	return result, nil
}

func (s *ProductService) RemoveVariants(ctx context.Context, q *api.RemoveVariantsRequest) (*pbcm.RemovedResponse, error) {
	cmd := &catalog.DeleteShopVariantsCommand{
		ShopID: s.SS.Shop().ID,
		IDs:    q.Ids,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return result, nil
}

func (s *ProductService) GetProduct(ctx context.Context, q *pbcm.IDRequest) (*api.ShopProduct, error) {
	shopID := s.SS.Shop().ID
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	productPb, err := GetProductQuantity(ctx, s.InventoryQuery, shopID, query.Result)
	if err != nil {
		return nil, err
	}
	applyProductInfoForVariants([]*shop.ShopProduct{productPb})
	result := productPb
	return result, nil
}

func (s *ProductService) GetProductsByIDs(ctx context.Context, q *pbcm.IDsRequest) (*api.ShopProductsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    q.Ids,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	products, err := GetProductsQuantity(ctx, s.InventoryQuery, shopID, query.Result.Products)
	if err != nil {
		return nil, err
	}
	result := &api.ShopProductsResponse{
		Products: products,
	}
	return result, nil
}

func (s *ProductService) GetProducts(ctx context.Context, q *api.GetVariantsRequest) (*api.ShopProductsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	shopID := s.SS.Shop().ID
	var fullTextSearch filter.FullTextSearch = ""
	if q.Filter != nil {
		fullTextSearch = q.Filter.Name
	}
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
		Name:    fullTextSearch,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	products, err := GetProductsQuantity(ctx, s.InventoryQuery, shopID, query.Result.Products)
	if err != nil {
		return nil, err
	}

	applyProductInfoForVariants(products)

	result := &api.ShopProductsResponse{
		Paging: cmapi.PbPaging(cm.Paging{
			Limit: query.Result.Paging.Limit,
			Sort:  query.Result.Paging.Sort,
		}),
		Products: products,
	}
	return result, nil
}

func applyProductInfoForVariants(products []*shop.ShopProduct) {
	for _, product := range products {
		productID := product.Id
		productName := product.Name
		for _, variant := range product.Variants {
			variant.ProductId = productID
			variant.Product = &shop.ShopShortProduct{
				Id:   productID,
				Name: productName,
			}
		}
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, q *api.CreateProductRequest) (*api.ShopProduct, error) {
	metaFields := []*catalog.MetaField{}

	for _, metaField := range q.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	cmd := &catalog.CreateShopProductCommand{
		ShopID:    s.SS.Shop().ID,
		Code:      q.Code,
		Name:      q.Name,
		Unit:      q.Unit,
		ImageURLs: q.ImageUrls,
		Note:      q.Note,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   q.ShortDesc,
			Description: q.Description,
			DescHTML:    q.DescHtml,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   q.CostPrice,
			ListPrice:   q.ListPrice,
			RetailPrice: q.RetailPrice,
		},
		BrandID:     q.BrandId,
		ProductType: q.ProductType.Apply(0),
		MetaFields:  metaFields,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductService) RemoveProducts(ctx context.Context, q *api.RemoveVariantsRequest) (*pbcm.RemovedResponse, error) {
	cmd := &catalog.DeleteShopProductsCommand{
		ShopID: s.SS.Shop().ID,
		IDs:    q.Ids,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return result, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, q *api.UpdateProductRequest) (*api.ShopProduct, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopProductInfoCommand{
		ShopID:    shopID,
		ProductID: q.Id,
		Code:      q.Code,
		Name:      q.Name,
		Unit:      q.Unit,
		Note:      q.Note,
		BrandID:   q.BrandId,

		ShortDesc:   q.ShortDesc,
		Description: q.Description,
		DescHTML:    q.DescHtml,
		CategoryID:  q.CategoryID,
		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,
		ProductType: q.ProductType,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateProductsStatus(ctx context.Context, q *api.UpdateProductStatusRequest) (*api.UpdateProductStatusResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopProductStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.UpdateProductStatusResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ProductService) UpdateVariantsStatus(ctx context.Context, q *api.UpdateProductStatusRequest) (*api.UpdateProductStatusResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopVariantStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.UpdateProductStatusResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ProductService) UpdateProductsTags(ctx context.Context, q *api.UpdateProductsTagsRequest) (*pbcm.UpdatedResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalogmodelx.UpdateShopProductsTagsCommand{
		ShopID:     shopID,
		ProductIDs: q.Ids,
		Update: &model.UpdateListRequest{
			Adds:       q.Adds,
			Deletes:    q.Deletes,
			ReplaceAll: q.ReplaceAll,
			DeleteAll:  q.DeleteAll,
		},
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}

func (s *ProductService) populateVariantInfos(ctx context.Context, shopID dot.ID, variants []*shop.ShopVariant) error {
	shopIDs := make([]dot.ID, 0, len(variants))
	query := &catalog.ListShopProductsByIDsQuery{
		IDs:    shopIDs,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	mapShopProduct := make(map[dot.ID]*catalog.ShopProduct)
	for _, product := range query.Result.Products {
		mapShopProduct[product.ProductID] = product
	}

	for _, variant := range variants {
		variant.Product = &shop.ShopShortProduct{
			Id:   variant.ProductId,
			Name: mapShopProduct[variant.ProductId].Name,
		}
	}
	return nil
}

func (s *ProductService) GetVariant(ctx context.Context, q *api.GetVariantRequest) (*api.ShopVariant, error) {
	query := &catalog.GetShopVariantQuery{
		Code:      q.Code,
		VariantID: q.ID,
		ShopID:    s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shopVariantPb := PbShopVariant(query.Result)
	if err := s.populateVariantInfos(ctx, s.SS.Shop().ID, []*shop.ShopVariant{shopVariantPb}); err != nil {
		return nil, err
	}
	result := shopVariantPb

	return result, nil
}

func (s *ProductService) GetVariantsByIDs(ctx context.Context, q *pbcm.IDsRequest) (*api.ShopVariantsResponse, error) {
	query := &catalog.ListShopVariantsWithProductByIDsQuery{
		IDs:    q.Ids,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &api.ShopVariantsResponse{Variants: PbShopVariantsWithProducts(query.Result.Variants)}

	return result, nil
}

func (s *ProductService) CreateVariant(ctx context.Context, q *api.CreateVariantRequest) (*api.ShopVariant, error) {
	cmd := &catalog.CreateShopVariantCommand{
		ShopID:     s.SS.Shop().ID,
		ProductID:  q.ProductId,
		Code:       q.Code,
		Name:       q.Name,
		ImageURLs:  q.ImageUrls,
		Note:       q.Note,
		Attributes: types.ValidateAttributesEmpty(q.Attributes),
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   q.ShortDesc,
			Description: q.Description,
			DescHTML:    q.DescHtml,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   q.CostPrice,
			ListPrice:   q.ListPrice,
			RetailPrice: q.RetailPrice,
		},
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopVariant(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateProductImages(ctx context.Context, q *api.UpdateVariantImagesRequest) (*api.ShopProduct, error) {
	shopID := s.SS.Shop().ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op: meta.OpDeleteAll,
		})
	}
	if q.ReplaceAll != nil {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op:      meta.OpReplaceAll,
			Changes: q.ReplaceAll,
		})
	}
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpAdd,
		Changes: q.Adds,
	})
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpRemove,
		Changes: q.Deletes,
	})

	cmd := catalog.UpdateShopProductImagesCommand{
		ShopID:  shopID,
		ID:      q.Id,
		Updates: metaUpdate,
	}

	if err := s.CatalogAggr.Dispatch(ctx, &cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateProductMetaFields(ctx context.Context, q *api.UpdateProductMetaFieldsRequest) (*api.ShopProduct, error) {
	metaFields := []*catalog.MetaField{}
	for _, metaField := range q.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	cmd := catalog.UpdateShopProductMetaFieldsCommand{
		ID:         q.Id,
		ShopID:     s.SS.Shop().ID,
		MetaFields: metaFields,
	}
	if err := s.CatalogAggr.Dispatch(ctx, &cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateVariantImages(ctx context.Context, q *api.UpdateVariantImagesRequest) (*api.ShopVariant, error) {
	shopID := s.SS.Shop().ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op: meta.OpDeleteAll,
		})
	}
	if q.ReplaceAll != nil {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op:      meta.OpReplaceAll,
			Changes: q.ReplaceAll,
		})
	}
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpAdd,
		Changes: q.Adds,
	})
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpRemove,
		Changes: q.Deletes,
	})

	cmd := catalog.UpdateShopVariantImagesCommand{
		ShopID:  shopID,
		ID:      q.Id,
		Updates: metaUpdate,
	}
	if err := s.CatalogAggr.Dispatch(ctx, &cmd); err != nil {
		return nil, err
	}
	result := PbShopVariant(cmd.Result)
	return result, nil
}

func (s *ProductService) UpdateProductCategory(ctx context.Context, q *api.UpdateProductCategoryRequest) (*api.ShopProduct, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopProductCategoryCommand{
		ProductID:  q.ProductId,
		CategoryID: q.CategoryId,
		ShopID:     shopID,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductService) AddProductCollection(ctx context.Context, r *api.AddShopProductCollectionRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        s.SS.Shop().ID,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ProductService) RemoveProductCollection(ctx context.Context, r *api.RemoveShopProductCollectionRequest) (*pbcm.RemovedResponse, error) {
	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        s.SS.Shop().ID,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{Removed: cmd.Result}
	return result, nil
}

func (s *ProductService) RemoveProductCategory(ctx context.Context, r *pbcm.IDRequest) (*api.ShopProduct, error) {
	cmd := &catalog.RemoveShopProductCategoryCommand{
		ShopID:    s.SS.Shop().ID,
		ProductID: r.Id,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func GetProductsQuantity(ctx context.Context, inventoryQuery inventory.QueryBus, shopID dot.ID, products []*catalog.ShopProductWithVariants) ([]*shop.ShopProduct, error) {
	var variantIDs []dot.ID
	for _, valueProduct := range products {
		for _, valueVariant := range valueProduct.Variants {
			variantIDs = append(variantIDs, valueVariant.VariantID)
		}
	}
	inventoryVariants, err := getVariantsQuantity(ctx, inventoryQuery, shopID, variantIDs)
	if err != nil {
		return nil, err
	}
	return PbProductsQuantity(products, inventoryVariants), nil
}

func GetProductQuantity(ctx context.Context, inventoryQuery inventory.QueryBus, shopID dot.ID, shopProduct *catalog.ShopProductWithVariants) (*shop.ShopProduct, error) {
	var variantIDs []dot.ID
	for _, variant := range shopProduct.Variants {
		variantIDs = append(variantIDs, variant.VariantID)
	}
	inventoryVariants, err := getVariantsQuantity(ctx, inventoryQuery, shopID, variantIDs)
	if err != nil {
		return nil, err
	}
	shopProductPb := PbProductQuantity(shopProduct, inventoryVariants)
	return shopProductPb, nil
}

func getVariantsQuantity(ctx context.Context, inventoryQuery inventory.QueryBus, shopID dot.ID, variantIDs []dot.ID) (map[dot.ID]*inventory.InventoryVariant, error) {

	var mapInventoryVariant = make(map[dot.ID]*inventory.InventoryVariant)
	if len(variantIDs) == 0 {
		return mapInventoryVariant, nil
	}
	q := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: variantIDs,
	}
	if err := inventoryQuery.Dispatch(ctx, q); err != nil {
		return nil, err
	}

	for _, value := range q.Result.InventoryVariants {
		mapInventoryVariant[value.VariantID] = value
	}
	return mapInventoryVariant, nil
}

func (s *ProductService) GetVariantsBySupplierID(ctx context.Context, q *api.GetVariantsBySupplierIDRequest) (*api.ShopVariantsResponse, error) {
	query := &catalog.GetVariantsBySupplierIDQuery{
		SupplierID: q.SupplierId,
		ShopID:     s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	shopVariantsPb := PbShopVariants(query.Result.Variants)
	if err := s.populateVariantInfos(ctx, s.SS.Shop().ID, shopVariantsPb); err != nil {
		return nil, err
	}
	result := &api.ShopVariantsResponse{Variants: shopVariantsPb}
	return result, nil
}
