package partnerimport

import (
	"context"

	api "o.o/api/top/external/whitelabel"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/catalog/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Products(ctx context.Context, r *api.ImportProductsRequest) (*api.ImportProductsResponse, error) {
	if len(r.Products) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, product := range r.Products {
		var brandID, categoryID dot.ID

		if product.ExternalID == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "id should not be null")
		}

		if product.ExternalBrandID != "" {
			brand, err := s.brandStoreFactory(ctx).ExternalID(product.ExternalBrandID).GetShopBrandDB()
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "brand_id is invalid")
			}
			brandID = brand.ID
		}

		if product.ExternalCategoryID != "" {
			category, err := s.categoryStoreFactory(ctx).ExternalID(product.ExternalCategoryID).GetShopCategoryDB()
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "category_id is invalid")
			}
			categoryID = category.ID
		}

		shopProduct := &model.ShopProduct{
			ShopID:             s.SS.Shop().ID,
			PartnerID:          s.SS.Claim().AuthPartnerID,
			ExternalID:         product.ExternalID,
			ExternalCode:       product.ExternalCode,
			ExternalBrandID:    product.ExternalBrandID,
			ExternalCategoryID: product.ExternalCategoryID,
			Code:               product.ExternalCode,
			Name:               product.Name,
			Unit:               product.Unit,
			ImageURLs:          product.ImageUrls,
			Note:               product.Note,
			ShortDesc:          product.ShortDesc,
			Description:        product.Description,
			CostPrice:          product.CodePrice,
			ListPrice:          product.ListPrice,
			RetailPrice:        product.RetailPrice,
			BrandID:            brandID,
			CategoryID:         categoryID,
			Status:             status3.P,
			CreatedAt:          product.CreatedAt.ToTime(),
			UpdatedAt:          product.UpdatedAt.ToTime(),
			DeletedAt:          product.DeletedAt.ToTime(),
		}

		oldShopProduct, err := s.shopProductStoreFactory(ctx).ExternalID(product.ExternalID).GetShopProductDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopProduct.ProductID = id
			if _err := s.shopProductStoreFactory(ctx).CreateShopProductImport(shopProduct); _err != nil {
				return nil, _err
			}
		case cm.NoError:
			shopProduct.ProductID = oldShopProduct.ProductID
			ids = append(ids, oldShopProduct.ProductID)
			if _err := s.shopProductStoreFactory(ctx).UpdateShopProduct(shopProduct); _err != nil {
				return nil, _err
			}
		default:
			return nil, err
		}
	}

	modelProducts, err := s.shopProductStoreFactory(ctx).IDs(ids...).ListShopProductsDB()
	if err != nil {
		return nil, err
	}

	var productsResponse []*api.Product
	for _, product := range modelProducts {
		productsResponse = append(productsResponse, &api.Product{
			Id:                 product.ProductID,
			PartnerID:          product.PartnerID,
			ShopID:             product.ShopID,
			ExternalId:         dot.String(product.ExternalID),
			ExternalCode:       dot.String(product.ExternalCode),
			ExternalBrandID:    product.ExternalBrandID,
			ExternalCategoryID: product.ExternalCategoryID,
			Name:               dot.String(product.Name),
			Description:        dot.String(product.Description),
			ShortDesc:          dot.String(product.ShortDesc),
			ImageUrls:          product.ImageURLs,
			Note:               dot.String(product.Note),
			Status:             status3.WrapStatus(product.Status),
			ListPrice:          dot.Int(product.ListPrice),
			RetailPrice:        dot.Int(product.RetailPrice),
			CreatedAt:          cmapi.PbTime(product.CreatedAt),
			UpdatedAt:          cmapi.PbTime(product.UpdatedAt),
			BrandId:            dot.WrapID(product.BrandID),
			CategoryId:         dot.WrapID(product.CategoryID),
		})
	}
	result := &api.ImportProductsResponse{
		Products: productsResponse,
	}
	return result, nil
}
