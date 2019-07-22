package sqlstore

import (
	"context"
	"strconv"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/services/catalog/convert"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
	catalogsqlstore "etop.vn/backend/pkg/services/catalog/sqlstore"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("sql",
		AddProductsToShopCollection,
		GetProductsExtended,
		GetShopVariant,
		RemoveProductsFromShopCollection,
		RemoveShopVariants,
		UpdateShopVariant,

		GetVariantsExtended,

		AddShopVariants,
		AddShopProducts,
		RemoveShopProducts,
		UpdateShopProduct,
		UpdateShopProductsTags,
	)
}

var (
	filterProductWhitelist = catalogsqlstore.FilterProductWhitelist
	filterVariantWhitelist = catalogsqlstore.FilterVariantWhitelist
	shopProductStore       catalogsqlstore.ShopProductStoreFactory
	shopVariantStore       catalogsqlstore.ShopVariantStoreFactory
)

func GetVariantByProductIDs(productIds []int64, filters []cm.Filter) ([]*catalogmodel.Variant, error) {
	s := x.Table("variant")
	s, _, err := Filters(s, filters, filterVariantWhitelist)
	if err != nil {
		return nil, err
	}
	var variants []*catalogmodel.Variant

	if err := s.Where("v.deleted_at is NULL").In("v.product_id", productIds).Find((*catalogmodel.Variants)(&variants)); err != nil {
		return nil, err
	}
	return variants, nil
}

func GetProductsExtended(ctx context.Context, query *catalogmodelx.DeprecatedGetProductsExtendedQuery) error {
	s := x.Table("product").Where("p.deleted_at is NULL")
	if query.ProductSourceType != "" {
		s = s.Where("ps.type = ?", query.ProductSourceType)
	}
	s = FilterStatus(s, "p.", query.StatusQuery)
	filtersProduct := make([]cm.Filter, 0, len(query.Filters))
	filtersVariant := make([]cm.Filter, 0, len(query.Filters))
	for _, filter := range query.Filters {
		if strings.Contains(filter.Name, "v.") {
			filter.Name = filter.Name[2:]
			filtersVariant = append(filtersVariant, filter)
		} else {
			filtersProduct = append(filtersProduct, filter)
		}
	}

	s, _, err := Filters(s, filtersProduct, filterProductWhitelist)
	if err != nil {
		return err
	}

	var products []*catalogmodel.Product
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "p.id", "created_at": "p.created_at", "updated_at": "p.updated_at", "name": "p.name"})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("p.id", query.IDs)
		}
		if err := s2.Find((*catalogmodel.Products)(&products)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&catalogmodel.Product{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}

	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}
	variants, err := GetVariantByProductIDs(productIDs, filtersVariant)
	if err != nil {
		return err
	}

	result := make([]*catalogmodel.ProductFtVariant, len(products))
	hashProductVariant := make(map[int64][]*catalogmodel.Variant)

	for _, v := range variants {
		hashProductVariant[v.ProductID] = append(hashProductVariant[v.ProductID], v)
	}

	for i, p := range products {
		result[i] = &catalogmodel.ProductFtVariant{
			Product:  p,
			Variants: hashProductVariant[p.ID],
		}
	}

	query.Result.Products = result
	return nil
}

func GetVariantsExtended(ctx context.Context, query *catalogmodelx.GetVariantsExtendedQuery) error {
	if query.SkipPaging {
		// IDs, Codes or EdCodes mut be provided
		if query.IDs == nil && query.EdCodes == nil && query.Codes == nil {
			return cm.Error(cm.InvalidArgument, "Neither id or code provided", nil)
		}
	}

	s := x.Table("variant").Where("v.deleted_at is NULL")
	s = FilterStatus(s, "v.", query.StatusQuery)

	s, _, err := Filters(s, query.Filters, filterProductWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s
		if query.IsPaging() {
			s2 = s.Clone()
		}

		s2, err := LimitSort(s2, query.Paging,
			Ms{"id": "", "created_at": "", "updated_at": "", "name": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("v.id", query.IDs)
		}
		if query.Codes != nil {
			s2 = s2.In("v.code", query.Codes)
		}
		if query.EdCodes != nil {
			s2 = s2.In("v.ed_code", query.EdCodes)
		}
		if err := s2.Find((*catalogmodel.VariantExtendeds)(&query.Result.Variants)); err != nil {
			return err
		}
	}
	if query.IsPaging() {
		total, err := s.Count(&catalogmodel.VariantExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetShopVariant(ctx context.Context, query *catalogmodelx.GetShopVariantQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	product := new(catalogmodel.ShopVariantExtended)
	if err := x.Where("sv.shop_id = ? AND sv.variant_id = ?", query.ShopID, query.VariantID).
		ShouldGet(product); err != nil {
		return err
	}
	query.Result = convert.ShopVariantExtended(product)
	return nil
}

func AddProductsToShopCollection(ctx context.Context, cmd *catalogmodelx.AddProductsToShopCollectionCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.CollectionID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing CollectionID", nil)
	}

	collection := new(catalogmodel.ShopCollection)
	if has, err := x.Table("shop_collection").
		Where("id = ? AND shop_id = ?", cmd.CollectionID, cmd.ShopID).
		Get(collection); err != nil {
		return nil
	} else if !has {
		return cm.Error(cm.NotFound, "Collection not found", nil)
	}

	errs := make([]error, len(cmd.ProductIDs))
	cmd.Result.Errors = errs
	updated := 0
	for i, id := range cmd.ProductIDs {
		productShopCollection := &catalogmodel.ProductShopCollection{
			ProductID:    id,
			CollectionID: cmd.CollectionID,
			ShopID:       cmd.ShopID,
		}
		if err := x.Table("product_shop_collection").ShouldInsert(productShopCollection); err != nil {
			errs[i] = err
			continue
		} else {
			updated++
		}
	}
	cmd.Result.Updated = updated
	return nil
}

func RemoveProductsFromShopCollection(ctx context.Context, cmd *catalogmodelx.RemoveProductsFromShopCollectionCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	if deleted, err := x.Table("product_shop_collection").Where("shop_id = ? AND collection_id = ?", cmd.ShopID, cmd.CollectionID).
		In("product_id", cmd.ProductIDs).Delete(&catalogmodel.ProductShopCollection{}); err != nil {
		return err
	} else if deleted == 0 {
		return cm.Error(cm.NotFound, "", nil)
	} else {
		cmd.Result.Updated = int(deleted)
	}
	return nil
}

func AddShopVariants(ctx context.Context, cmd *catalogmodelx.AddShopVariantsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ids", nil)
	}

	query := &catalogmodelx.GetVariantsExtendedQuery{
		IDs: cmd.IDs,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result.Variants) == 0 {
		return cm.Error(cm.NotFound, "No available product to add", nil)
	}
	products := query.Result.Variants

	var sql []byte
	counter, appendPlaceholder := sqlPlaceholder(0)

	sql = append(sql, `
INSERT INTO shop_variant("shop_id", "variant_id", "product_id", "retail_price", "created_at", "updated_at") VALUES(`...)
	for i, p := range products {
		if i > 0 {
			sql = append(sql, "),("...)
		}
		sql = strconv.AppendInt(sql, cmd.ShopID, 10)
		sql = append(sql, ',')
		sql = strconv.AppendInt(sql, p.ID, 10)
		sql = append(sql, ',')
		sql = strconv.AppendInt(sql, p.Product.ID, 10)
		sql = append(sql, ',')
		sql = strconv.AppendInt(sql, int64(p.ListPrice), 10)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
	}
	sql = append(sql, ") ON CONFLICT DO NOTHING"...)

	now := time.Now()
	args := make([]interface{}, *counter)
	for i := range args {
		args[i] = now
	}

	if res, err := x.Exec(string(sql), args...); err != nil {
		return err
	} else if updated, err := res.RowsAffected(); updated == 0 {
		return cm.Error(cm.AlreadyExists, "No product was added", err)
	}

	xproducts := make([]*catalogmodel.ShopVariantExtended, len(products))
	for i, p := range products {
		for _, id := range cmd.IDs {
			if p.ID == id {
				xproducts[i] = &catalogmodel.ShopVariantExtended{
					ShopVariant: &catalogmodel.ShopVariant{
						ShopID:      cmd.ShopID,
						VariantID:   id,
						ProductID:   p.Product.ID,
						RetailPrice: p.ListPrice,
					},
					Variant: p.Variant,
				}
				break
			}
		}
	}
	cmd.Result.Variants = convert.ShopVariantsExtended(xproducts)

	errors := make([]error, len(cmd.IDs))
	for i, id := range cmd.IDs {
		ok := false
		for _, p := range products {
			if p.ID == id {
				ok = true
				break
			}
		}
		if !ok {
			errors[i] = cm.Error(cm.NotFound, "", nil)
		}
	}
	cmd.Result.Errors = errors
	return nil
}

func updateOrInsertShopVariant(sv *catalogmodel.ShopVariant, productSourceID int64, x Qx) error {
	if sv.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing VariantID", nil)
	}

	if sv.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var shopProducts []*catalogmodel.ShopProductFtProductFtVariantFtShopVariant
	if err := x.Table("shop_product").
		Where("sp.shop_id = ? AND v.id = ?", sv.ShopID, sv.VariantID).
		Find((*catalogmodel.ShopProductFtProductFtVariantFtShopVariants)(&shopProducts)); err != nil {
		return err
	}

	if len(shopProducts) == 0 {
		// case: variant has product_source_id which type shop custom
		var variant = new(catalogmodel.Variant)
		if err := x.Table("variant").Where("id = ? AND product_source_id = ? AND status = 1 AND deleted_at is NULL", sv.VariantID, productSourceID).
			ShouldGet(variant); err != nil {
			return err
		}

		// add to table shop_product
		return inTransaction(func(x Qx) error {
			var product = new(catalogmodel.Product)
			if err := x.Table("product").Where("id = ? AND deleted_at is NULL", variant.ProductID).ShouldGet(product); err != nil {
				return err
			}

			shopProduct := ConvertProductToShopProduct(product)
			shopProduct.ShopID = sv.ShopID
			if err := x.Table("shop_product").ShouldInsert(shopProduct); err != nil {
				return err
			}

			variant, err := getVariant(sv.VariantID, x)
			if err != nil {
				return err
			}
			shopVariant := buildShopVariant(variant, sv)
			if err := x.Table("shop_variant").ShouldInsert(shopVariant); err != nil {
				return err
			}
			return nil
		})
	}

	if shopProducts[0].ShopVariant.VariantID != 0 {
		// case: variant in shop_variant
		if err := x.Table("shop_variant").Where("variant_id = ? AND shop_id = ?", sv.VariantID, sv.ShopID).
			ShouldUpdate(sv); err != nil {
			return err
		}
	} else {
		// case: variant not in table shop_variant, but product is in table shop_product
		variant, err := getVariant(sv.VariantID, x)
		if err != nil {
			return err
		}
		shopVariant := buildShopVariant(variant, sv)
		if err := x.Table("shop_variant").ShouldInsert(shopVariant); err != nil {
			return err
		}
	}
	return nil
}

func getVariant(id int64, x Qx) (*catalogmodel.Variant, error) {
	var variant = new(catalogmodel.Variant)
	if err := x.Table("variant").Where("id = ?", id).ShouldGet(variant); err != nil {
		return nil, err
	}
	return variant, nil
}

func buildShopVariant(v *catalogmodel.Variant, sv *catalogmodel.ShopVariant) *catalogmodel.ShopVariant {
	if v.ID != sv.VariantID {
		return nil
	}
	return &catalogmodel.ShopVariant{
		VariantID:   v.ID,
		Name:        v.GetName(),
		Description: cm.Coalesce(sv.Description, v.Description),
		DescHTML:    cm.Coalesce(sv.DescHTML, v.DescHTML),
		ShortDesc:   cm.Coalesce(sv.ShortDesc, v.ShortDesc),
		ImageURLs:   cm.CoalesceStrings(sv.ImageURLs, v.ImageURLs),
		Note:        sv.Note,
		RetailPrice: cm.CoalesceInt32(sv.RetailPrice, v.ListPrice),
		Status:      model.CoalesceStatus3(sv.Status, v.Status),
		ProductID:   cm.CoalesceInt64(sv.ProductID, v.ProductID),
		ShopID:      sv.ShopID,
	}
}

func UpdateShopVariant(ctx context.Context, cmd *catalogmodelx.UpdateShopVariantCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Variant == nil || cmd.Variant.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing product", nil)
	}

	sv := *cmd.Variant
	if err := sv.BeforeUpdate(); err != nil {
		return err
	}

	updateErr := inTransaction(func(x Qx) error {
		if err := updateOrInsertShopVariant(&sv, cmd.ProductSourceID, x); err != nil {
			return err
		}
		// update cost_price, inventory
		productSource := new(catalogmodel.ProductSource)
		if has, _ := x.Table("product_source").Where("id = ? AND type = ?", cmd.ProductSourceID, catalogmodel.ProductSourceCustom).
			Get(productSource); has {
			variant := &catalogmodel.Variant{
				ID: sv.VariantID,
			}
			if cmd.CostPrice != 0 {
				variant.CostPrice = cmd.CostPrice
			}
			if len(cmd.Attributes) > 0 {
				variant.Attributes = cmd.Attributes
			}
			if cmd.Code != "" {
				variant.Code = cmd.Code
			}
			_, err := x.Table("variant").Where("id = ? AND product_source_id = ?", sv.VariantID, cmd.ProductSourceID).Update(variant)
			if _err := CheckErrorProductCode(err); _err != nil {
				return _err
			}
		}
		return nil
	})
	if updateErr != nil {
		return updateErr
	}

	query := &catalogmodelx.GetShopVariantQuery{
		ShopID:    cmd.ShopID,
		VariantID: cmd.Variant.VariantID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.Error(cm.Internal, "", err)
	}

	cmd.Result = query.Result
	return nil
}

func RemoveShopVariants(ctx context.Context, cmd *catalogmodelx.RemoveShopVariantsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	return inTransaction(func(x Qx) error {
		if _, err := x.Table("shop_variant").
			Where("shop_id = ?", cmd.ShopID).In("variant_id", cmd.IDs).Delete(&catalogmodel.ShopVariant{}); err != nil {
			return err
		}

		now := time.Now()
		updated, err := x.Table("variant").
			Where("product_source_id = ? AND deleted_at is NULL", cmd.ProductSourceID).
			In("id", cmd.IDs).UpdateMap(M{"deleted_at": now})
		if err != nil {
			return err
		}

		cmd.Result.Removed = int(updated)
		return nil
	})
}

func AddShopProducts(ctx context.Context, cmd *catalogmodelx.AddShopProductsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ids", nil)
	}

	query := &catalogmodelx.DeprecatedGetProductsExtendedQuery{
		IDs: cmd.IDs,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result.Products) == 0 {
		return cm.Error(cm.NotFound, "No available product to add", nil)
	}
	products := query.Result.Products

	var sql []byte
	var args []interface{}
	now := time.Now()
	_, appendPlaceholder := sqlPlaceholder(0)

	sql = append(sql, `
INSERT INTO shop_product("shop_id", "product_id", "name", "description", "short_desc", "desc_html", "image_urls", "created_at", "updated_at") VALUES(`...)
	for i, p := range products {
		if i > 0 {
			sql = append(sql, "),("...)
		}
		sql = strconv.AppendInt(sql, cmd.ShopID, 10)
		sql = append(sql, ',')
		sql = strconv.AppendInt(sql, p.Product.ID, 10)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		sql = appendPlaceholder(sql)
		args = append(args, p.Product.Name, p.Product.Description, p.Product.ShortDesc, p.Product.DescHTML, x.Opts().Array(p.Product.ImageURLs), now, now)
	}
	sql = append(sql, ") ON CONFLICT DO NOTHING"...)

	if res, err := x.Exec(string(sql), args...); err != nil {
		return err
	} else if updated, err := res.RowsAffected(); updated == 0 {
		return cm.Error(cm.AlreadyExists, "No product was added", err)
	}

	xproducts := make([]*catalogmodel.ShopProduct, len(products))
	for i, p := range products {
		for _, id := range cmd.IDs {
			if p.Product.ID == id {
				xproducts[i] = &catalogmodel.ShopProduct{
					ShopID:      cmd.ShopID,
					ProductID:   p.Product.ID,
					Name:        p.Product.Name,
					Description: p.Product.Description,
					DescHTML:    p.DescHTML,
					ShortDesc:   p.ShortDesc,
					ImageURLs:   p.Product.ImageURLs,
					Status:      p.Product.Status,
				}
				break
			}
		}
	}
	cmd.Result.Products = xproducts

	errors := make([]error, len(cmd.IDs))
	for i, id := range cmd.IDs {
		ok := false
		for _, p := range products {
			if p.Product.ID == id {
				ok = true
				break
			}
		}
		if !ok {
			errors[i] = cm.Error(cm.NotFound, "", nil)
		}
	}
	cmd.Result.Errors = errors
	return nil
}

func ConvertProductToShopProduct(p *catalogmodel.Product) *catalogmodel.ShopProduct {
	return &catalogmodel.ShopProduct{
		ProductID:   p.ID,
		Name:        p.Name,
		Description: p.Description,
		DescHTML:    p.DescHTML,
		ShortDesc:   p.ShortDesc,
		ImageURLs:   p.ImageURLs,
		Status:      p.Status,
	}
}

func RemoveShopProducts(ctx context.Context, cmd *catalogmodelx.RemoveShopProductsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Misssing IDs", nil)
	}

	return inTransaction(func(x Qx) error {
		var productsCount uint64
		{
			s1 := x.Table("product").
				In("p.id", cmd.IDs).
				Where("p.product_source_id = ? AND sp.product_id is NULL AND p.deleted_at is NULL", cmd.ProductSourceID)
			var err error
			productsCount, err = s1.Count(&catalogmodel.ProductFtShopProduct{})
			if err != nil {
				return err
			}
		}

		deletedCount, err := x.Table("shop_product").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			Delete(&catalogmodel.ShopProduct{})
		if err != nil {
			return nil
		}

		if _, err2 := x.Table("shop_variant").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			Delete(&catalogmodel.ShopVariant{}); err2 != nil {
			return nil
		}

		now := time.Now()
		if _, err := x.Table("product").
			Where("product_source_id = ? AND deleted_at is NULL", cmd.ProductSourceID).
			In("id", cmd.IDs).
			UpdateMap(M{"deleted_at": now}); err != nil {
			return err
		}

		if _, err := x.Table("variant").
			Where("product_source_id = ? AND deleted_at is NULL", cmd.ProductSourceID).
			In("product_id", cmd.IDs).
			UpdateMap(M{"deleted_at": now}); err != nil {
			return err
		}
		cmd.Result.Removed = int(productsCount) + int(deletedCount)
		return nil
	})
}

func UpdateOrInsertShopProduct(sp *catalogmodel.ShopProduct, productSourceID int64) error {
	return updateOrInsertShopProduct(sp, productSourceID, x)
}

func updateOrInsertShopProduct(sp *catalogmodel.ShopProduct, productSourceID int64, x Qx) error {
	if sp.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductID", nil)
	}

	if sp.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var shopProduct = new(catalogmodel.ShopProduct)
	if has, err := x.Table("shop_product").
		Where("product_id = ? AND shop_id = ?", sp.ProductID, sp.ShopID).
		Get(shopProduct); err != nil {
		return err
	} else if has {
		if err := x.Table("shop_product").
			Where("product_id = ? AND shop_id = ?", sp.ProductID, sp.ShopID).
			ShouldUpdate(sp); err != nil {
			return err
		}
		return nil
	}
	var product = new(catalogmodel.Product)
	if err := x.Table("product").Where("id = ? AND product_source_id = ? AND deleted_at is NULL", sp.ProductID, productSourceID).
		ShouldGet(product); err != nil {
		return err
	}

	if err := x.Table("shop_product").ShouldInsert(sp); err != nil {
		return err
	}
	return nil
}

func UpdateShopProduct(ctx context.Context, cmd *catalogmodelx.UpdateShopProductCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Product == nil || cmd.Product.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing product", nil)
	}
	if err := cmd.Product.BeforeUpdate(); err != nil {
		return err
	}

	sp := *cmd.Product
	if err := sp.BeforeUpdate(); err != nil {
		return err
	}

	errUpdate := inTransaction(func(x Qx) error {
		if err := updateOrInsertShopProduct(&sp, cmd.ProductSourceID, x); err != nil {
			return err
		}

		productSource := new(catalogmodel.ProductSource)
		if has, _ := x.Table("product_source").Where("id = ? AND type = ?", cmd.ProductSourceID, catalogmodel.ProductSourceCustom).
			Get(productSource); has {
			if cmd.Code != "" {
				_, err2 := x.Table("product").Where("id = ? AND product_source_id = ?", cmd.Product.ProductID, cmd.ProductSourceID).
					UpdateMap(M{
						"ed_code": cmd.Code,
					})
				if _err := CheckErrorProductCode(err2); _err != nil {
					return _err
				}
			}
		}
		return nil
	})
	if errUpdate != nil {
		return errUpdate
	}

	{
		q := shopProductStore(ctx).ShopID(cmd.ShopID).ID(cmd.Product.ProductID)
		product, err := q.GetShopProductWithVariants()
		if err != nil {
			return err
		}
		cmd.Result = product
	}
	return nil
}

func UpdateShopProductsTags(ctx context.Context, cmd *catalogmodelx.UpdateShopProductsTagsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	req := cmd.Update
	if err := req.Verify(); err != nil {
		return err
	}

	var products []*catalogmodel.ShopProduct
	if err := x.Where("shop_id = ?", cmd.ShopID).
		In("product_id", cmd.ProductIDs).
		Find((*catalogmodel.ShopProducts)(&products)); err != nil {
		return err
	}

	productMap := make(map[int64]*catalogmodel.ShopProduct)
	for _, p := range products {
		productMap[p.ProductID] = p
	}

	countUpdated := 0
	var savedError error
	for _, id := range cmd.ProductIDs {
		p := productMap[id]
		var pTag []string
		if p != nil {
			pTag = p.Tags
		}
		tags, tErr := model.PatchTag(pTag, *req)
		if tErr != nil {
			savedError = tErr
			continue
		}
		sp := &catalogmodel.ShopProduct{
			ShopID:    cmd.ShopID,
			ProductID: id,
			Tags:      tags,
		}

		if err := UpdateOrInsertShopProduct(sp, cmd.ProductSourceID); err != nil {
			savedError = err
			continue
		}
		countUpdated++
	}
	if countUpdated > 0 {
		cmd.Result.Updated = countUpdated
		return nil
	}
	if savedError != nil {
		return savedError
	}
	return cm.Error(cm.NotFound, "No product updated", nil)
}
