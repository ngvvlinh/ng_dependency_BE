package sqlstore

import (
	"context"

	com "o.o/backend/com/main"
	catalogmodel "o.o/backend/com/main/catalog/model"
	catalogmodelx "o.o/backend/com/main/catalog/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
)

type CategoryStoreInterface interface {
	GetProductSourceCategories(ctx context.Context, query *catalogmodelx.GetProductSourceCategoriesQuery) error

	GetShopCategory(ctx context.Context, query *catalogmodelx.GetShopCategoryQuery) error

	RemoveShopShopCategory(ctx context.Context, cmd *catalogmodelx.RemoveShopCategoryCommand) error

	UpdateShopShopCategory(ctx context.Context, cmd *catalogmodelx.UpdateShopCategoryCommand) error
}

type CategoryStore struct {
	DB com.MainDB
	db *cmsql.Database `wire:"-"`
}

func BindCategoryStore(s *CategoryStore) (to CategoryStoreInterface) {
	s.db = s.DB
	return s
}

func (st *CategoryStore) NewCategoryStore(db com.MainDB) *CategoryStore {
	s := &CategoryStore{
		db: db,
	}
	return s
}

func (st *CategoryStore) GetShopCategory(ctx context.Context, query *catalogmodelx.GetShopCategoryQuery) error {
	if query.CategoryID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	p := new(catalogmodel.ShopCategory)

	s := st.db.Table("shop_category").Where("id = ?", query.CategoryID)
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if err := s.ShouldGet(p); err != nil {
		return err
	}

	query.Result = p
	return nil
}

func (st *CategoryStore) GetProductSourceCategories(ctx context.Context, query *catalogmodelx.GetProductSourceCategoriesQuery) error {
	s := st.db.Table("shop_category")
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}
	if query.IDs != nil {
		s = s.In("id", query.IDs)
	}

	err := s.Find((*catalogmodel.ShopCategories)(&query.Result.Categories))
	return err
}

func (st *CategoryStore) UpdateShopShopCategory(ctx context.Context, cmd *catalogmodelx.UpdateShopCategoryCommand) error {
	cat := &catalogmodel.ShopCategory{
		ParentID: cmd.ParentID,
		Name:     cmd.Name,
	}
	if err := st.db.Table("shop_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).ShouldUpdate(cat); err != nil {
		return err
	}

	query := &catalogmodelx.GetShopCategoryQuery{
		CategoryID: cmd.ID,
		ShopID:     cmd.ShopID,
	}

	if err := st.GetShopCategory(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func (st *CategoryStore) RemoveShopShopCategory(ctx context.Context, cmd *catalogmodelx.RemoveShopCategoryCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	return inTransaction(st.db, func(tx Qx) error {
		if _, err := tx.Table("shop_product").Where("category_id = ?", cmd.ID).
			UpdateMap(M{"category_id": nil}); err != nil {
			return err
		}

		if err := tx.Table("shop_category").Where("id = ? AND shop_id = ?", cmd.ID, cmd.ShopID).
			ShouldDelete(&catalogmodel.ShopCategory{}); err != nil {
			return err
		}
		cmd.Result.Removed = 1
		return nil
	})
}
