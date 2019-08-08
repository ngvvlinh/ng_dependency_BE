package sqlstore

import (
	"context"
	"strings"
	"time"

	catalogmodel "etop.vn/backend/com/main/catalog/model"
	catalogmodelx "etop.vn/backend/com/main/catalog/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("sql",
		CreateShopCategory,
		DeprecatedCreateVariant,
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

func DeprecatedCreateVariant(ctx context.Context, cmd *catalogmodelx.DeprecatedCreateVariantCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.ProductID == 0 && cmd.ProductName == "" {
		return cm.Error(cm.InvalidArgument, "Missing ProductName", nil)
	}

	productID := cmd.ProductID
	err := x.InTransaction(ctx, func(s cmsql.QueryInterface) error {
		variant := &catalogmodel.ShopVariant{
			ShopID:      cmd.ShopID,
			VariantID:   cm.NewID(),
			ProductID:   cmd.ProductID,
			Code:        cmd.VariantCode,
			Name:        cmd.Name,
			Description: cmd.Description,
			DescHTML:    cmd.DescHTML,
			ShortDesc:   cmd.ShortDesc,
			ImageURLs:   cmd.ImageURLs,
			Note:        "",
			Tags:        nil,
			CostPrice:   cmd.CostPrice,
			ListPrice:   cmd.ListPrice,
			RetailPrice: cmd.RetailPrice,
			Status:      model.StatusActive,
			Attributes:  cmd.Attributes,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			NameNorm:    validate.NormalizeSearch(cmd.Name),
		}
		variant.Attributes, variant.AttrNormKv = catalogmodel.NormalizeAttributes(cmd.Attributes)

		if cmd.ProductID != 0 {
			_, err := shopProductStore(ctx).ShopID(cmd.ShopID).ID(cmd.ProductID).GetShopProductDB()
			if err != nil {
				return err
			}

		} else {
			product := &catalogmodel.ShopProduct{
				ShopID:     cmd.ShopID,
				ProductID:  cm.NewID(),
				Code:       cmd.ProductCode,
				Name:       cmd.ProductName,
				NameNorm:   validate.NormalizeSearch(cmd.ProductName),
				NameNormUa: validate.NormalizeUnaccent(cmd.ProductName),
			}
			variant.ProductID = product.ProductID
			productID = product.ProductID
		}

		return x.ShouldInsert(variant)
	})
	if err != nil {
		return err
	}

	q := shopProductStore(ctx).ShopID(cmd.ShopID).ID(productID)
	product, err := q.GetShopProductWithVariants()
	if err != nil {
		return err
	}
	cmd.Result = product
	return nil
}

func CreateShopCategory(ctx context.Context, cmd *catalogmodelx.CreateShopCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Name == "" {
		return cm.Error(cm.InvalidArgument, "Missing category name", nil)
	}
	name := strings.ToLower(cmd.Name)
	name = strings.Title(name)

	var productSourceCategory = new(catalogmodel.ShopCategory)
	s := x.Table("shop_category").Where("shop_id = ? AND name = ?", cmd.ShopID, name)
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

	psCategory := &catalogmodel.ShopCategory{
		ID:       cm.NewID(),
		ParentID: cmd.ParentID,
		ShopID:   cmd.ShopID,
		Name:     name,
	}

	if err := x.Table("shop_category").ShouldInsert(psCategory); err != nil {
		return err
	}
	cmd.Result = psCategory
	return nil
}

func UpdateProductsPSCategory(ctx context.Context, cmd *catalogmodelx.UpdateProductsShopCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.CategoryID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShopCategoryID", nil)
	}

	category := new(catalogmodel.ShopCategory)
	if has, err := x.Table("shop_category").
		Where("id = ? AND shop_id = ?", cmd.CategoryID, cmd.ShopID).
		Get(category); err != nil {
		return nil
	} else if !has {
		return cm.Error(cm.NotFound, "ShopCategory not found", nil)
	}

	if updated, err := x.Table("shop_product").
		Where("product_source_id = ?", cmd.ShopID).
		In("product_id", cmd.ProductIDs).
		UpdateMap(M{"shop_category_id": cmd.CategoryID}); err != nil {
		return err
	} else if updated == 0 {
		return cm.Error(cm.NotFound, "No product updated", nil)
	} else {
		cmd.Result.Updated = int(updated)
	}
	return nil
}
