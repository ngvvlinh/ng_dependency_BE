package whitelabel

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/catalog/types"
	"etop.vn/api/top/external/whitelabel"
	"etop.vn/api/top/types/etc/status3"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/capi/dot"
)

func (s *ImportService) Variants(ctx context.Context, r *VariantsEndpoint) error {
	if len(r.Variants) > MaximumItems {
		return cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, variant := range r.Variants {
		var productID dot.ID

		if variant.ExternalId == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "external_id should not be null")
		}
		if variant.ExternalProductId == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "external_product_id should not be null")
		} else {
			product, err := shopProductStoreFactory(ctx).ExternalID(variant.ExternalProductId).GetShopProductDB()
			if err != nil {
				return cm.Errorf(cm.InvalidArgument, err, "external_product_id is invalid")
			}
			productID = product.ProductID
		}

		shopVariant := &catalog.ShopVariant{
			ExternalID:        variant.ExternalId,
			ExternalCode:      variant.ExternalCode,
			ExternalProductID: variant.ExternalProductId,
			PartnerID:         r.Context.AuthPartnerID,
			ShopID:            r.Context.Shop.ID,
			ProductID:         productID,
			Name:              variant.Name,
			ShortDesc:         variant.ShortDesc,
			Description:       variant.Description,
			ImageURLs:         variant.ImageUrls,
			Status:            status3.Z,
			Attributes:        variant.Attributes,
			CostPrice:         variant.CostPrice,
			ListPrice:         variant.ListPrice,
			RetailPrice:       variant.RetailPrice,
			Note:              variant.Note,
			CreatedAt:         variant.CreatedAt.ToTime(),
			UpdatedAt:         variant.UpdatedAt.ToTime(),
			DeletedAt:         variant.DeletedAt.ToTime(),
		}

		oldShopVariant, err := shopVariantStoreFactory(ctx).ExternalID(variant.ExternalId).GetShopVariantDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopVariant.VariantID = id
			if _err := shopVariantStoreFactory(ctx).CreateShopVariantImport(shopVariant); _err != nil {
				return _err
			}
		case cm.NoError:
			shopVariant.VariantID = oldShopVariant.VariantID
			ids = append(ids, oldShopVariant.VariantID)
			if _err := shopVariantStoreFactory(ctx).UpdateShopVariantImport(shopVariant); _err != nil {
				return _err
			}
		default:
			return err
		}
	}

	modelVariants, err := shopVariantStoreFactory(ctx).IDs(ids...).ListShopVariantsDB()
	if err != nil {
		return err
	}

	var variantsResponse []*whitelabel.ShopVariant
	for _, variant := range modelVariants {
		var attributes []*types.Attribute
		for _, attribute := range variant.Attributes {
			attributes = append(attributes, &types.Attribute{
				Name:  attribute.Name,
				Value: attribute.Value,
			})
		}
		variantsResponse = append(variantsResponse, &whitelabel.ShopVariant{
			ExternalId:        variant.ExternalID,
			ExternalCode:      variant.ExternalCode,
			ExternalProductId: variant.ExternalProductID,
			Id:                variant.VariantID,
			ProductID:         variant.ProductID,
			ShopID:            variant.ShopID,
			PartnerID:         variant.PartnerID,
			Code:              variant.Code,
			Name:              variant.Name,
			Description:       variant.Description,
			ShortDesc:         variant.ShortDesc,
			ImageUrls:         variant.ImageURLs,
			ListPrice:         variant.ListPrice,
			RetailPrice:       variant.RetailPrice,
			Note:              variant.Note,
			Status:            variant.Status,
			CostPrice:         variant.CostPrice,
			Attributes:        attributes,
			CreatedAt:         cmapi.PbTime(variant.CreatedAt),
			UpdatedAt:         cmapi.PbTime(variant.UpdatedAt),
			DeletedAt:         cmapi.PbTime(variant.DeletedAt),
		})
	}
	r.Result = &whitelabel.ImportShopVariantsResponse{Variants: variantsResponse}
	return nil
}
