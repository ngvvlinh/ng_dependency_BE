package sqlstore

import (
	"context"

	catalogmodel "etop.vn/backend/com/main/catalog/model"
	catalogmodelx "etop.vn/backend/com/main/catalog/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
)

func init() {
	bus.AddHandler("sql", GetShopCategory)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", UpdateShopShopCategory)
	bus.AddHandler("sql", RemoveShopShopCategory)
}

func GetShopCategory(ctx context.Context, query *catalogmodelx.GetShopCategoryQuery) error {
	if query.CategoryID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	p := new(catalogmodel.ShopCategory)

	s := x.Table("shop_category").Where("id = ?", query.CategoryID)
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if err := s.ShouldGet(p); err != nil {
		return err
	}

	query.Result = p
	return nil
}

func GetProductSourceCategories(ctx context.Context, query *catalogmodelx.GetProductSourceCategoriesQuery) error {
	s := x.Table("shop_category")
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.IDs != nil {
		s = s.In("id", query.IDs)
	}

	err := s.Find((*catalogmodel.ShopCategories)(&query.Result.Categories))
	return err
}

func UpdateShopShopCategory(ctx context.Context, cmd *catalogmodelx.UpdateShopCategoryCommand) error {
	cat := &catalogmodel.ShopCategory{
		ParentID: cmd.ParentID,
		Name:     cmd.Name,
	}
	if err := x.Table("shop_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).ShouldUpdate(cat); err != nil {
		return err
	}

	query := &catalogmodelx.GetShopCategoryQuery{
		CategoryID: cmd.ID,
		ShopID:     cmd.ShopID,
	}

	if err := GetShopCategory(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func RemoveShopShopCategory(ctx context.Context, cmd *catalogmodelx.RemoveShopCategoryCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	return inTransaction(func(s Qx) error {
		if _, err := s.Table("shop_product").Where("shop_category_id = ?", cmd.ID).
			UpdateMap(M{"shop_category_id": nil}); err != nil {
			return err
		}

		if err := s.Table("shop_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).
			ShouldDelete(&catalogmodel.ShopCategory{}); err != nil {
			return err
		}
		cmd.Result.Removed = 1
		return nil
	})
}
