package sqlstore

import (
	"context"
	"strings"
	"time"

	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
	catalogmodelx "o.o/backend/com/main/catalog/modelx"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
)

type ShopStoreInterface interface {

	CreateShopCategory(ctx context.Context, cmd *catalogmodelx.CreateShopCategoryCommand) error

	DeprecatedCreateVariant(ctx context.Context, cmd *catalogmodelx.DeprecatedCreateVariantCommand) error

	GetShop(ctx context.Context, query *identitymodelx.GetShopQuery) error

	GetShopExtended(ctx context.Context, query *identitymodelx.GetShopExtendedQuery) error

	GetShopWithPermission(ctx context.Context, query *identitymodelx.GetShopWithPermissionQuery) error

	GetShops(ctx context.Context, query *identitymodelx.GetShopsQuery) error

	UpdateProductsPSCategory(ctx context.Context, cmd *catalogmodelx.UpdateProductsShopCategoryCommand) error
}

type ShopStore struct {
	db *cmsql.Database
}

func NewShopStore(db com.MainDB) *ShopStore {
	s := &ShopStore{
		db: db,
	}
	return s
}

func (st *ShopStore) GetAllShopExtendedsQuery(ctx context.Context, query *identitymodelx.GetAllShopExtendedsQuery) error {
	s := st.db.Table("shop").Where("s.deleted_at is NULL")
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := sqlstore.LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"created_at": "s.created_at", "updated_at": "s.updated_at"})
		if err != nil {
			return err
		}
		var shops []*identitymodel.ShopExtended
		if err := s2.Find((*identitymodel.ShopExtendeds)(&shops)); err != nil {
			return err
		}
		query.Result.Shops = shops
	}
	return nil
}

func (st *ShopStore) GetShop(ctx context.Context, query *identitymodelx.GetShopQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	shop := new(identitymodel.Shop)
	if err := st.db.Where("id = ?", query.ShopID).
		Where("deleted_at is NULL").
		ShouldGet(shop); err != nil {
		return err
	}

	query.Result = shop
	return nil
}

func (st *ShopStore) GetShops(ctx context.Context, query *identitymodelx.GetShopsQuery) error {
	return st.db.Table("shop").
		In("id", query.ShopIDs).
		Find((*identitymodel.Shops)(&query.Result.Shops))
}

func (st *ShopStore) GetShopExtended(ctx context.Context, query *identitymodelx.GetShopExtendedQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "", nil)
	}

	var shop identitymodel.ShopExtended
	s := st.db.Where("s.id = ?", query.ShopID)
	if !query.IncludeDeleted {
		s = s.Where("s.deleted_at is NULL")
	}

	err := s.ShouldGet(&shop)
	query.Result = &shop
	return err
}

func (st *ShopStore) GetShopWithPermission(ctx context.Context, query *identitymodelx.GetShopWithPermissionQuery) error {
	if query.ShopID == 0 || query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	shop := new(identitymodel.Shop)
	if err := st.db.Where("id = ?", query.ShopID).
		ShouldGet(shop); err != nil {
		return err
	}
	query.Result.Shop = shop

	accUser := new(identitymodel.AccountUser)
	if err := st.db.
		Where("account_id = ? AND user_id = ?", query.ShopID, query.UserID).
		ShouldGet(accUser); err != nil {
		return err
	}
	query.Result.Permission = accUser.Permission
	return nil
}

func (st *ShopStore) DeprecatedCreateVariant(ctx context.Context, cmd *catalogmodelx.DeprecatedCreateVariantCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.ProductID == 0 && cmd.ProductName == "" {
		return cm.Error(cm.InvalidArgument, "Missing ProductName", nil)
	}

	productID := cmd.ProductID
	err := st.db.InTransaction(ctx, func(s cmsql.QueryInterface) error {
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
			Status:      status3.P,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			NameNorm:    validate.NormalizeSearch(cmd.Name),
		}
		attributes, attrNormKv := catalogmodel.NormalizeAttributes(cmd.Attributes)
		variant.Attributes = convert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(attributes)
		variant.AttrNormKv = attrNormKv
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
				ImageURLs:  cmd.ImageURLs,
				NameNorm:   validate.NormalizeSearch(cmd.ProductName),
				NameNormUa: validate.NormalizeUnaccent(cmd.ProductName),
			}
			variant.ProductID = product.ProductID
			productID = product.ProductID
			if err := st.db.ShouldInsert(product); err != nil {
				return err
			}
		}
		return st.db.ShouldInsert(variant)
	})
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "shop_product_shop_id_code_idx"):
			err = cm.Errorf(cm.FailedPrecondition, nil, "Mã sản phẩm %v đã tồn tại. Vui lòng chọn mã khác.", cmd.ProductCode)
		case strings.Contains(errMsg, "shop_variant_shop_id_code_idx"):
			err = cm.Errorf(cm.FailedPrecondition, nil, "Mã phiên bản sản phẩm %v đã tồn tại. Vui lòng chọn mã khác.", cmd.VariantCode)
		}
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

func (st *ShopStore) CreateShopCategory(ctx context.Context, cmd *catalogmodelx.CreateShopCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Name == "" {
		return cm.Error(cm.InvalidArgument, "Missing category name", nil)
	}
	name := strings.ToLower(cmd.Name)
	name = strings.Title(name)

	var productSourceCategory = new(catalogmodel.ShopCategory)
	s := st.db.Table("shop_category").Where("shop_id = ? AND name = ?", cmd.ShopID, name)
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

	if err := st.db.Table("shop_category").ShouldInsert(psCategory); err != nil {
		return err
	}
	cmd.Result = psCategory
	return nil
}

func (st *ShopStore) UpdateProductsPSCategory(ctx context.Context, cmd *catalogmodelx.UpdateProductsShopCategoryCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.CategoryID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShopCategoryID", nil)
	}

	category := new(catalogmodel.ShopCategory)
	if has, err := st.db.Table("shop_category").
		Where("id = ? AND shop_id = ?", cmd.CategoryID, cmd.ShopID).
		Get(category); err != nil {
		return nil
	} else if !has {
		return cm.Error(cm.NotFound, "ShopCategory not found", nil)
	}

	if updated, err := st.db.Table("shop_product").
		Where("shop_id = ?", cmd.ShopID).
		In("product_id", cmd.ProductIDs).
		UpdateMap(M{"category_id": cmd.CategoryID}); err != nil {
		return err
	} else if updated == 0 {
		return cm.Error(cm.NotFound, "No product updated", nil)
	} else {
		cmd.Result.Updated = updated
	}
	return nil
}
