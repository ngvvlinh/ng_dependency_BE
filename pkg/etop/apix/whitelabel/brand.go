package whitelabel

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/whitelabel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/capi/dot"
)

func (s *ImportService) Brands(ctx context.Context, r *BrandsEndpoint) error {
	if len(r.Brands) > MaximumItems {
		return cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, brand := range r.Brands {
		if brand.ExternalID == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "external_id should not be null")
		}
		shopBrand := &catalog.ShopBrand{
			PartnerID:   r.Context.AuthPartnerID,
			ShopID:      r.Context.Shop.ID,
			ExternalID:  brand.ExternalID,
			BrandName:   brand.BrandName,
			Description: brand.Description,
			CreatedAt:   brand.CreatedAt.ToTime(),
			UpdatedAt:   brand.UpdatedAt.ToTime(),
			DeletedAt:   brand.DeletedAt.ToTime(),
		}

		oldShopBrand, err := brandStoreFactory(ctx).ExternalID(brand.ExternalID).GetShopBrandDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			shopBrand.ID = id
			ids = append(ids, id)
			if _err := brandStoreFactory(ctx).CreateShopBrand(shopBrand); _err != nil {
				return err
			}
		case cm.NoError:
			shopBrand.ID = oldShopBrand.ID
			ids = append(ids, oldShopBrand.ID)
			if _err := brandStoreFactory(ctx).ExternalID(brand.ExternalID).UpdateShopBrand(shopBrand); _err != nil {
				return err
			}
		default:
			return err
		}

	}

	modelBrands, err := brandStoreFactory(ctx).IDs(ids...).ListShopBrandsDB()
	if err != nil {
		return err
	}

	var brandsResponse []*whitelabel.Brand
	for _, brand := range modelBrands {
		brandsResponse = append(brandsResponse, &whitelabel.Brand{
			ID:          brand.ID,
			PartnerID:   brand.PartnerID,
			ShopID:      brand.ShopID,
			ExternalID:  brand.ExternalID,
			BrandName:   brand.BrandName,
			Description: brand.Description,
			CreatedAt:   cmapi.PbTime(brand.CreatedAt),
			UpdatedAt:   cmapi.PbTime(brand.UpdatedAt),
			DeletedAt:   cmapi.PbTime(brand.DeletedAt),
		})
	}
	r.Result = &whitelabel.ImportBrandsResponse{Brands: brandsResponse}
	return nil
}
