package whitelabel

import (
	"context"

	"etop.vn/api/top/external/whitelabel"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/capi/dot"
)

const shopID = 1057650360222204339

func (s *ImportService) Products(ctx context.Context, r *ProductsEndpoint) error {
	if len(r.Products) > MaximumItems {
		return cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, product := range r.ImportProductsRequest.Products {
		var brandID, categoryID dot.ID

		if product.ExternalID == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "id should not be null")
		}

		if product.ExternalBrandID != "" {
			brand, err := brandStoreFactory(ctx).ExternalID(product.ExternalBrandID).GetShopBrandDB()
			if err != nil {
				return cm.Errorf(cm.InvalidArgument, err, "brand_id is invalid")
			}
			brandID = brand.ID
		}

		if product.ExternalCategoryID != "" {
			category, err := categoryStoreFactory(ctx).ExternalID(product.ExternalCategoryID).GetShopCategoryDB()
			if err != nil {
				return cm.Errorf(cm.InvalidArgument, err, "category_id is invalid")
			}
			categoryID = category.ID
		}

		shopProduct := &model.ShopProduct{
			ShopID:             r.Context.Shop.ID,
			PartnerID:          r.Context.AuthPartnerID,
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

		oldShopProduct, err := shopProductStoreFactory(ctx).ExternalID(product.ExternalID).GetShopProductDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopProduct.ProductID = id
			if _err := shopProductStoreFactory(ctx).CreateShopProductImport(shopProduct); _err != nil {
				return _err
			}
		case cm.NoError:
			shopProduct.ProductID = oldShopProduct.ProductID
			ids = append(ids, oldShopProduct.ProductID)
			if _err := shopProductStoreFactory(ctx).UpdateShopProduct(shopProduct); _err != nil {
				return _err
			}
		default:
			return err
		}
	}

	modelProducts, err := shopProductStoreFactory(ctx).IDs(ids...).ListShopProductsDB()
	if err != nil {
		return err
	}

	var productsResponse []*whitelabel.Product
	for _, product := range modelProducts {
		productsResponse = append(productsResponse, &whitelabel.Product{
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
	r.Result = &whitelabel.ImportProductsResponse{
		Products: productsResponse,
	}
	return nil
}
