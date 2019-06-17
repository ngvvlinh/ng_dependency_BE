package sqlstore

import (
	"context"
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
)

func init() {
	bus.AddHandlers("sql",
		CreateProductSource,
		CreateProductSourceCategory,
		CreateShopCollection,
		CreateVariant,
		GetAllShopExtendedsQuery,
		GetShop,
		GetShopCollection,
		GetShopCollections,
		GetShopExtended,
		GetShopProductSources,
		GetShops,
		GetShopWithPermission,
		RemoveShopCollection,
		UpdateProductsPSCategory,
		UpdateShopCollection,
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

func GetShopCollection(ctx context.Context, query *catalogmodelx.GetShopCollectionQuery) error {
	if query.CollectionID == 0 {
		return cm.Error(cm.InvalidArgument, "Thiếu CollectionID", nil)
	}

	s := x.Table("shop_collection").Where("id = ?", query.CollectionID)
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}

	collection := new(catalogmodel.ShopCollection)
	if has, err := s.Get(collection); err != nil {
		return err
	} else if !has {
		return cm.Error(cm.NotFound, "", nil)
	}

	query.Result = collection
	return nil
}

func GetShopCollections(ctx context.Context, query *catalogmodelx.GetShopCollectionsQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	s := x.Table("shop_collection").
		Where("shop_id = ?", query.ShopID)
	if query.CollectionIDs != nil {
		s = s.In("id", query.CollectionIDs)
	}

	if err := s.Find((*catalogmodel.ShopCollections)(&query.Result.Collections)); err != nil {
		return err
	}
	return nil
}

func CreateShopCollection(ctx context.Context, cmd *catalogmodelx.CreateShopCollectionCommand) error {
	collection := cmd.Collection
	if collection.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	collection.ID = cm.NewID()
	if _, err := x.Table("shop_collection").Insert(collection); err != nil {
		return err
	}
	cmd.Result = collection
	return nil
}

func UpdateShopCollection(ctx context.Context, cmd *catalogmodelx.UpdateShopCollectionCommand) error {
	collection := cmd.Collection
	if collection.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing CollectionID", nil)
	}

	s := x.Table("shop_collection").Where("id = ?", collection.ID)
	if collection.ShopID != 0 {
		s = s.Where("shop_id = ?", collection.ShopID)
	}
	if affected, err := s.Update(collection); err != nil {
		return err
	} else if affected == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	cmd.Result = collection
	return nil
}

func RemoveShopCollection(ctx context.Context, cmd *catalogmodelx.RemoveShopCollectionCommand) error {
	return inTransaction(func(s Qx) error {
		if cmd.CollectionID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing CollectionID", nil)
		}
		{
			s2 := s.Table("product_shop_collection").Where("collection_id = ?", cmd.CollectionID)
			if cmd.ShopID != 0 {
				s2 = s2.Where("shop_id = ?", cmd.ShopID)
			}
			if _, err := s2.Delete(&catalogmodel.ProductShopCollection{}); err != nil {
				return err
			}
		}
		{
			s2 := s.Table("shop_collection").Where("id = ?", cmd.CollectionID)
			if cmd.ShopID != 0 {
				s2 = s2.Where("shop_id = ?", cmd.ShopID)
			}

			if deleted, err := s2.Delete(&catalogmodel.ShopCollection{}); err != nil {
				return err
			} else if deleted == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
		}
		cmd.Result.Deleted = 1
		return nil
	})
}

func CreateProductSource(ctx context.Context, cmd *catalogmodelx.CreateProductSourceCommand) error {
	if cmd.Type == "" {
		return cm.Error(cm.InvalidArgument, "Type can not be nil", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	return inTransaction(func(s Qx) error {
		shop := new(model.Shop)
		if err := s.Table("shop").
			Where("id = ?", cmd.ShopID).
			ShouldGet(shop); err != nil {
			return err
		}

		// Reuse the existing product source
		if shop.ProductSourceID != 0 {
			var productSource catalogmodel.ProductSource
			if err := s.Table("product_source").
				Where("id = ?", shop.ProductSourceID).
				ShouldGet(&productSource); err != nil {
				return cm.Error(cm.Internal, "", err)
			}
			cmd.Result = &productSource
			return nil
		}

		source := &catalogmodel.ProductSource{
			Type:   cmd.Type,
			Name:   cmd.Name,
			Status: model.StatusActive,
		}
		source.ID = cm.NewID()
		if _, err := s.Table("product_source").Insert(source); err != nil {
			return err
		}
		cmd.Result = source

		if _, err := s.Table("shop").Where("id = ?", cmd.ShopID).UpdateMap(M{"product_source_id": source.ID}); err != nil {
			return err
		}
		return nil
	})
}

func CreateVariant(ctx context.Context, cmd *catalogmodelx.CreateVariantCommand) error {
	if cmd.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var productSource = new(catalogmodel.ProductSource)
	if err := x.Table("product_source").Where("id = ?", cmd.ProductSourceID).ShouldGet(productSource); err != nil {
		return cm.Error(cm.InvalidArgument, "ProductSource does not existed", nil)
	}
	if cmd.ProductID == 0 && cmd.ProductName == "" {
		return cm.Error(cm.InvalidArgument, "Missing ProductName", nil)
	}

	variant := &catalogmodel.Variant{
		ID:          cm.NewID(),
		ProductID:   cmd.ProductID,
		ShortDesc:   cmd.ShortDesc,
		Description: cmd.Description,
		DescHTML:    cmd.DescHTML,
		Status:      model.StatusActive,
		Code:        cmd.SKU,
		ListPrice:   cmd.ListPrice,
		ImageURLs:   cmd.ImageURLs,
		CostPrice:   cmd.CostPrice,
		Attributes:  cmd.Attributes,
	}
	if err := variant.BeforeInsert(); err != nil {
		return err
	}

	productID := cmd.ProductID
	if productID != 0 {
		// create variant on this productID
		_, err := x.Table("variant").Insert(variant)
		if _err := CheckErrorProductCode(err); _err != nil {
			return _err
		}

	} else {
		// create product + shop_product + variant
		errInsert := inTransaction(func(s Qx) error {
			product := &catalogmodel.Product{
				ID:              cm.NewID(),
				ProductSourceID: cmd.ProductSourceID,
				Name:            cmd.ProductName,
				ShortDesc:       cmd.ShortDesc,
				Description:     cmd.Description,
				Status:          model.StatusActive,
				Code:            cmd.Code,
				ImageURLs:       cmd.ImageURLs,
				DescHTML:        cmd.DescHTML,
			}
			if err := product.BeforeInsert(); err != nil {
				return err
			}
			_, err := s.Table("product").Insert(product)
			if _err := CheckErrorProductCode(err); _err != nil {
				return _err
			}

			shopProduct := ConvertProductToShopProduct(product)
			shopProduct.ShopID = cmd.ShopID
			shopProduct.Tags = cmd.Tags
			if err := s.Table("shop_product").ShouldInsert(shopProduct); err != nil {
				return err
			}

			variant.ProductID = product.ID
			variant.ProductSourceID = product.ProductSourceID
			productID = product.ID
			_, err2 := s.Table("variant").Insert(variant)
			if _err := CheckErrorProductCode(err2); _err != nil {
				return _err
			}

			return nil
		})
		if errInsert != nil {
			return errInsert
		}
	}
	{
		q := shopProductStore(ctx).ShopID(cmd.ShopID).ID(productID)
		product, err := q.GetShopProductWithVariants()
		if err != nil {
			return err
		}
		cmd.Result = product
	}
	return nil
}

func GetShopProductSources(ctx context.Context, query *catalogmodelx.GetShopProductSourcesCommand) error {
	if query.UserID == 0 && query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	var shops []*model.Shop
	s := x.Table("shop")
	if query.UserID != 0 {
		s = s.Where("owner_id = ?", query.UserID)
	}
	if query.ShopID != 0 {
		s = s.Where("id = ?", query.ShopID)
	}

	if err := s.Find((*model.Shops)(&shops)); err != nil {
		return err
	}
	if len(shops) == 0 {
		query.Result = nil
		return nil
	}

	productSourceIds := make([]int64, 0, len(shops))
	for _, shop := range shops {
		id := shop.ProductSourceID
		if id != 0 {
			productSourceIds = append(productSourceIds, id)
		}
	}
	if len(productSourceIds) == 0 {
		query.Result = nil
		return nil
	}

	var productSources []*catalogmodel.ProductSource
	if err := x.Table("product_source").In("id", productSourceIds).Find((*catalogmodel.ProductSources)(&productSources)); err != nil {
		return err
	}
	query.Result = productSources
	return nil
}

func CreateProductSourceCategory(ctx context.Context, cmd *catalogmodelx.CreateProductSourceCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
	}
	if cmd.ProductSourceType == "" {
		return cm.Error(cm.InvalidArgument, "Missing Product source type", nil)
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
		ID:                cm.NewID(),
		ProductSourceID:   cmd.ProductSourceID,
		ProductSourceType: cmd.ProductSourceType,
		ParentID:          cmd.ParentID,
		ShopID:            cmd.ShopID,
		Name:              name,
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
	if cmd.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
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
		Where("product_source_id = ?", cmd.ProductSourceID).
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
