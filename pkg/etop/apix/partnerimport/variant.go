package partnerimport

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/catalog/types"
	api "o.o/api/top/external/whitelabel"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Variants(ctx context.Context, r *api.ImportShopVariantsRequest) (*api.ImportShopVariantsResponse, error) {
	if len(r.Variants) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, variant := range r.Variants {
		var productID dot.ID

		if variant.ExternalId == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_id should not be null")
		}
		if variant.ExternalProductId == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_product_id should not be null")
		} else {
			product, err := s.shopProductStoreFactory(ctx).ExternalID(variant.ExternalProductId).GetShopProductDB()
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "external_product_id is invalid")
			}
			productID = product.ProductID
		}

		shopVariant := &catalog.ShopVariant{
			ExternalID:        variant.ExternalId,
			ExternalCode:      variant.ExternalCode,
			ExternalProductID: variant.ExternalProductId,
			PartnerID:         s.SS.Claim().AuthPartnerID,
			ShopID:            s.SS.Shop().ID,
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

		oldShopVariant, err := s.shopVariantStoreFactory(ctx).ExternalID(variant.ExternalId).GetShopVariantDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopVariant.VariantID = id
			if _err := s.shopVariantStoreFactory(ctx).CreateShopVariantImport(shopVariant); _err != nil {
				return nil, _err
			}
		case cm.NoError:
			shopVariant.VariantID = oldShopVariant.VariantID
			ids = append(ids, oldShopVariant.VariantID)
			if _err := s.shopVariantStoreFactory(ctx).UpdateShopVariantImport(shopVariant); _err != nil {
				return nil, _err
			}
		default:
			return nil, err
		}
	}

	modelVariants, err := s.shopVariantStoreFactory(ctx).IDs(ids...).ListShopVariantsDB()
	if err != nil {
		return nil, err
	}

	var variantsResponse []*api.ShopVariant
	for _, variant := range modelVariants {
		var attributes []*types.Attribute
		for _, attribute := range variant.Attributes {
			attributes = append(attributes, &types.Attribute{
				Name:  attribute.Name,
				Value: attribute.Value,
			})
		}
		variantsResponse = append(variantsResponse, &api.ShopVariant{
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
	result := &api.ImportShopVariantsResponse{Variants: variantsResponse}
	return result, nil
}
