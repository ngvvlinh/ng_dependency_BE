package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandler("sql", GetProductSourceCategory)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", UpdateShopProductSourceCategory)
	bus.AddHandler("sql", RemoveShopProductSourceCategory)
}

func GetProductSourceCategory(ctx context.Context, query *catalogmodelx.GetProductSourceCategoryQuery) error {
	if query.CategoryID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	p := new(catalogmodel.ProductSourceCategory)

	s := x.Table("product_source_category").Where("id = ?", query.CategoryID)
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
	s := x.Table("product_source_category")
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.IDs != nil {
		s = s.In("id", query.IDs)
	}

	err := s.Find((*catalogmodel.ProductSourceCategories)(&query.Result.Categories))
	return err
}

func UpdateShopProductSourceCategory(ctx context.Context, cmd *catalogmodelx.UpdateShopProductSourceCategoryCommand) error {
	cat := &catalogmodel.ProductSourceCategory{
		ParentID: cmd.ParentID,
		Name:     cmd.Name,
	}
	if err := x.Table("product_source_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).ShouldUpdate(cat); err != nil {
		return err
	}

	query := &catalogmodelx.GetProductSourceCategoryQuery{
		CategoryID: cmd.ID,
		ShopID:     cmd.ShopID,
	}

	if err := GetProductSourceCategory(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func RemoveShopProductSourceCategory(ctx context.Context, cmd *catalogmodelx.RemoveShopProductSourceCategoryCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	return inTransaction(func(s Qx) error {
		if _, err := s.Table("product").Where("product_source_category_id = ?", cmd.ID).
			UpdateMap(M{"product_source_category_id": nil}); err != nil {
			return err
		}

		if err := s.Table("product_source_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).
			ShouldDelete(&catalogmodel.ProductSourceCategory{}); err != nil {
			return err
		}
		cmd.Result.Removed = 1
		return nil
	})
}
