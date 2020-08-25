package partnerimport

import (
	"context"

	"o.o/api/main/catalog"
	api "o.o/api/top/external/whitelabel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Brands(ctx context.Context, r *api.ImportBrandsRequest) (*api.ImportBrandsResponse, error) {
	if len(r.Brands) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, brand := range r.Brands {
		if brand.ExternalID == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_id should not be null")
		}
		shopBrand := &catalog.ShopBrand{
			PartnerID:   s.SS.Claim().AuthPartnerID,
			ShopID:      s.SS.Shop().ID,
			ExternalID:  brand.ExternalID,
			BrandName:   brand.BrandName,
			Description: brand.Description,
			CreatedAt:   brand.CreatedAt.ToTime(),
			UpdatedAt:   brand.UpdatedAt.ToTime(),
			DeletedAt:   brand.DeletedAt.ToTime(),
		}

		oldShopBrand, err := s.brandStoreFactory(ctx).ExternalID(brand.ExternalID).GetShopBrandDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			shopBrand.ID = id
			ids = append(ids, id)
			if _err := s.brandStoreFactory(ctx).CreateShopBrand(shopBrand); _err != nil {
				return nil, err
			}
		case cm.NoError:
			shopBrand.ID = oldShopBrand.ID
			ids = append(ids, oldShopBrand.ID)
			if _err := s.brandStoreFactory(ctx).ExternalID(brand.ExternalID).UpdateShopBrand(shopBrand); _err != nil {
				return nil, err
			}
		default:
			return nil, err
		}

	}

	modelBrands, err := s.brandStoreFactory(ctx).IDs(ids...).ListShopBrandsDB()
	if err != nil {
		return nil, err
	}

	var brandsResponse []*api.Brand
	for _, brand := range modelBrands {
		brandsResponse = append(brandsResponse, &api.Brand{
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
	result := &api.ImportBrandsResponse{Brands: brandsResponse}
	return result, nil
}
