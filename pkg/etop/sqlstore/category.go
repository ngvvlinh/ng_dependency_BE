package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
)

func init() {
	bus.AddHandler("sql", GetProductSourceCategory)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", CreateEtopCategory)
	bus.AddHandler("sql", GetEtopCategories)
	bus.AddHandler("sql", GetProductSourceCategories)
	bus.AddHandler("sql", GetProductSourceCategoriesExtended)
	bus.AddHandler("sql", UpdateShopProductSourceCategory)
	bus.AddHandler("sql", RemoveShopProductSourceCategory)
}

func GetProductSourceCategory(ctx context.Context, query *catalogmodelx.GetProductSourceCategoryQuery) error {
	if query.CategoryID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	p := new(catalogmodel.ProductSourceCategory)

	s := x.Table("product_source_category").Where("psc.id = ?", query.CategoryID)
	if query.SupplierID != 0 {
		s = s.Where("psc.supplier_id = ?", query.SupplierID)
	}
	if query.ShopID != 0 {
		s = s.Where("psc.shop_id = ?", query.ShopID)
	}
	if err := s.ShouldGet(p); err != nil {
		return err
	}

	query.Result = p
	return nil
}

func GetProductSourceCategoriesExtended(ctx context.Context, query *catalogmodelx.GetProductSourceCategoriesExtendedQuery) error {
	s := x.Table("product_source_category")
	if query.SupplierID != 0 {
		s = s.Where("psc.supplier_id = ?", query.SupplierID)
	}
	if query.ShopID != 0 {
		s = s.Where("psc.shop_id = ?", query.ShopID)
	}
	if query.IDs != nil {
		s = s.In("psc.id", query.IDs)
	}
	if query.ProductSourceType != "" {
		s = s.Where("psc.product_source_type = ?", query.ProductSourceType)
	}

	err := s.Find((*catalogmodel.ProductSourceCategories)(&query.Result.Categories))
	return err
}

func GetProductSourceCategories(ctx context.Context, query *catalogmodelx.GetProductSourceCategoriesQuery) error {
	s := x.Table("product_source_category")
	if query.SupplierID != 0 {
		s = s.Where("supplier_id = ?", query.SupplierID)
	}
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.IDs != nil {
		s = s.In("id", query.IDs)
	}
	if query.ProductSourceType != "" {
		s = s.Where("product_source_type = ?", query.ProductSourceType)
	}

	err := s.Find((*catalogmodel.ProductSourceCategories)(&query.Result.Categories))
	return err
}

func CreateEtopCategory(ctx context.Context, cmd *model.CreateEtopCategoryCommand) error {
	category := cmd.Category
	category.ID = cm.NewID()
	category.Status = model.StatusActive

	_, err := x.Table("etop_category").Insert(category)
	cmd.Result = category
	return err
}

func GetEtopCategories(ctx context.Context, query *model.GetEtopCategoriesQuery) error {
	s := x.Table("etop_category")
	if query.Status != nil {
		s = s.Where("status = ?", *query.Status)
	}

	err := s.Find((*model.EtopCategories)(&query.Result.Categories))
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
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
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
