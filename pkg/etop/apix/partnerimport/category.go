package partnerimport

import (
	"context"

	"o.o/api/main/catalog"
	api "o.o/api/top/external/whitelabel"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *ImportService) Categories(ctx context.Context, r *api.ImportCategoriesRequest) (*api.ImportCategoriesResponse, error) {
	if len(r.Categories) > MaximumItems {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cannot handle rather than 100 items at once")
	}

	var ids []dot.ID
	for _, category := range r.Categories {
		var parentID dot.ID
		if category.ExternalID == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_id should not be null")
		}
		if category.ExternalParentID != "" {
			parentCategory, err := s.categoryStoreFactory(ctx).ExternalID(category.ExternalParentID).GetShopCategoryDB()
			if err != nil {
				return nil, cm.Errorf(cm.InvalidArgument, err, "external_parent_id is invalid")
			}
			parentID = parentCategory.ID
		}
		shopCategory := &catalog.ShopCategory{
			PartnerID:        s.SS.Claim().AuthPartnerID,
			ShopID:           s.SS.Shop().ID,
			ExternalID:       category.ExternalID,
			ExternalParentID: category.ExternalParentID,
			ParentID:         parentID,
			Name:             category.Name,
			Status:           status3.P.Enum(),
			CreatedAt:        category.CreatedAt.ToTime(),
			UpdatedAt:        category.UpdatedAt.ToTime(),
			DeletedAt:        category.DeletedAt.ToTime(),
		}

		oldCategory, err := s.categoryStoreFactory(ctx).ExternalID(category.ExternalID).GetShopCategoryDB()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			id := cm.NewID()
			ids = append(ids, id)
			shopCategory.ID = id
			if _err := s.categoryStoreFactory(ctx).CreateShopCategory(shopCategory); _err != nil {
				return nil, _err
			}
		case cm.NoError:
			shopCategory.ID = oldCategory.ID
			ids = append(ids, oldCategory.ID)
			if _err := s.categoryStoreFactory(ctx).UpdateShopCategory(shopCategory); _err != nil {
				return nil, _err
			}
		default:
			return nil, err
		}
	}

	modelCategories, err := s.categoryStoreFactory(ctx).IDs(ids...).ListShopCategoriesDB()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*api.Category
	for _, category := range modelCategories {
		categoriesResponse = append(categoriesResponse, &api.Category{
			ID:               category.ID,
			ShopID:           category.ShopID,
			PartnerID:        category.PartnerID,
			ExternalID:       category.ExternalID,
			ExternalParentID: category.ExternalParentID,
			ParentID:         category.ParentID,
			Name:             category.Name,
			CreatedAt:        cmapi.PbTime(category.CreatedAt),
			UpdatedAt:        cmapi.PbTime(category.UpdatedAt),
			DeletedAt:        cmapi.PbTime(category.DeletedAt),
		})
	}
	result := &api.ImportCategoriesResponse{Categories: categoriesResponse}
	return result, nil
}
