package sqlstore

import (
	"context"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("sql",
		CreateProductSourceCategory,
		CreateVariant,
		GetAllShopExtendedsQuery,
		GetShop,
		GetShopExtended,
		GetShops,
		GetShopWithPermission,
		UpdateProductsPSCategory,
	)
}

func GetAllShopExtendedsQuery(ctx context.Context, query *model.GetAllShopExtendedsQuery) error {
	s := x.Table("shop").Where("s.deleted_at is NULL")
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"created_at": "s.created_at", "updated_at": "s.updated_at"})
		if err != nil {
			return err
		}
		var shops []*model.ShopExtended
		if err := s2.Find((*model.ShopExtendeds)(&shops)); err != nil {
			return err
		}
		query.Result.Shops = shops
	}
	{
		total, err := s.Count(&model.ShopExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func (ft ShopFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func GetShop(ctx context.Context, query *model.GetShopQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	shop := new(model.Shop)
	if err := x.Where("id = ?", query.ShopID).
		Where("deleted_at is NULL").
		ShouldGet(shop); err != nil {
		return err
	}

	query.Result = shop
	return nil
}

func GetShops(ctx context.Context, query *model.GetShopsQuery) error {
	return x.Table("shop").
		In("id", query.ShopIDs).
		Find((*model.Shops)(&query.Result.Shops))
}

func GetShopExtended(ctx context.Context, query *model.GetShopExtendedQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	var shop model.ShopExtended
	s := x.Where("s.id = ?", query.ShopID)
	if !query.IncludeDeleted {
		s = s.Where("s.deleted_at is NULL")
	}

	err := s.ShouldGet(&shop)
	query.Result = &shop
	return err
}

func GetShopWithPermission(ctx context.Context, query *model.GetShopWithPermissionQuery) error {
	if query.ShopID == 0 || query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	shop := new(model.Shop)
	if err := x.Where("id = ?", query.ShopID).
		ShouldGet(shop); err != nil {
		return err
	}
	query.Result.Shop = shop

	accUser := new(model.AccountUser)
	if err := x.
		Where("account_id = ? AND user_id = ?", query.ShopID, query.UserID).
		ShouldGet(accUser); err != nil {
		return err
	}
	query.Result.Permission = accUser.Permission
	return nil
}

func CreateVariant(ctx context.Context, cmd *catalogmodelx.CreateVariantCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.ProductID == 0 && cmd.ProductName == "" {
		return cm.Error(cm.InvalidArgument, "Missing ProductName", nil)
	}

	variant := &catalogmodel.ShopVariant{
		ShopID:      cmd.ShopID,
		VariantID:   cm.NewID(),
		ProductID:   cmd.ProductID,
		Code:        cmd.SKU,
		Name:        cmd.Name,
		Description: cmd.Description,
		DescHTML:    cmd.DescHTML,
		ShortDesc:   cmd.ShortDesc,
		ImageURLs:   cmd.ImageURLs,
		Note:        "",
		Tags:        nil,
		ListPrice:   cmd.ListPrice,
		CostPrice:   cmd.CostPrice,
		RetailPrice: 0,
		Status:      model.StatusActive,
		Attributes:  cmd.Attributes,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		NameNorm:    "",
		AttrNormKv:  "",
	}
	if err := variant.BeforeInsert(); err != nil {
		return err
	}
	if err := x.ShouldInsert(variant); err != nil {
		return err
	}
	{
		q := shopProductStore(ctx).ShopID(cmd.ShopID).ID(cmd.ProductID)
		product, err := q.GetShopProductWithVariants()
		if err != nil {
			return err
		}
		cmd.Result = product
	}
	return nil
}

func CreateProductSourceCategory(ctx context.Context, cmd *catalogmodelx.CreateProductSourceCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Name == "" {
		return cm.Error(cm.InvalidArgument, "Missing category name", nil)
	}
	name := strings.ToLower(cmd.Name)
	name = strings.Title(name)

	var productSourceCategory = new(catalogmodel.ProductSourceCategory)
	s := x.Table("product_source_category").Where("shop_id = ? AND name = ?", cmd.ShopID, name)
	if cmd.ParentID != 0 {
		s = s.Where("parent_id = ?", cmd.ParentID)
	}
	has, err := s.Get(productSourceCategory)
	if err != nil {
		return err
	}
	if has {
		cmd.Result = productSourceCategory
		return nil
	}

	psCategory := &catalogmodel.ProductSourceCategory{
		ID:       cm.NewID(),
		ParentID: cmd.ParentID,
		ShopID:   cmd.ShopID,
		Name:     name,
	}

	if err := x.Table("product_source_category").ShouldInsert(psCategory); err != nil {
		return err
	}
	cmd.Result = psCategory
	return nil
}

func UpdateProductsPSCategory(ctx context.Context, cmd *catalogmodelx.UpdateProductsProductSourceCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.CategoryID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceCategoryID", nil)
	}

	category := new(catalogmodel.ProductSourceCategory)
	if has, err := x.Table("product_source_category").
		Where("id = ? AND shop_id = ?", cmd.CategoryID, cmd.ShopID).
		Get(category); err != nil {
		return nil
	} else if !has {
		return cm.Error(cm.NotFound, "ProductSourceCategory not found", nil)
	}

	if updated, err := x.Table("product").
		Where("product_source_id = ?", cmd.ShopID).
		In("id", cmd.ProductIDs).
		UpdateMap(M{"product_source_category_id": cmd.CategoryID}); err != nil {
		return err
	} else if updated == 0 {
		return cm.Error(cm.NotFound, "No product updated", nil)
	} else {
		cmd.Result.Updated = int(updated)
	}
	return nil
}
