package sqlstore

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/common/sql/core"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandlers("sql",
		AddProductsToShopCollection,
		AddShopVariants,
		GetProduct,
		GetProductsExtended,
		GetVariantExternalsFromID,
		GetShopVariant,
		GetShopVariants,
		RemoveProductsEtopCategory,
		RemoveProductsFromShopCollection,
		RemoveShopVariants,
		ScanVariantExternals,
		UpdateProduct,
		UpdateVariantImages,
		UpdateVariantPrice,
		UpdateVariants,
		UpdateProductsEtopCategory,
		UpdateVariantsStatus,
		UpdateShopVariant,
		UpdateShopVariantsStatus,
		UpdateShopVariantsTags,

		GetVariant,
		GetVariantsExtended,
		GetVariantExternals,
		UpdateVariant,
		UpdateProductImages,
		UpdateProductsStatus,

		AddShopProducts,
		GetShopProduct,
		GetShopProducts,
		RemoveShopProducts,
		UpdateShopProduct,
		UpdateShopProductsStatus,
		UpdateShopProductsTags,
		GetAllShopVariants,

		GetProducts,
		GetVariants,
	)
}

var (
	filterProductWhitelist = FilterWhitelist{
		Arrays:   []string{},
		Contains: []string{"name"},
		Equals:   []string{"supplier_id", "etop_category_id", "name"},
		Status:   []string{"ed_status", "status", "etop_status"},
		Numbers:  []string{"wholesale_price", "list_price", "retail_price_min", "retail_price_max", "ed_wholesale_price", "ed_list_price", "supplier_retail_price_min", "ed_retail_price_max"},
		Dates:    []string{"created_at", "updated_at"},
		Unaccent: []string{"name"},
		PrefixOrRename: map[string]string{
			"name":       "p",
			"status":     "p",
			"created_at": "p",
			"updated_at": "p",

			"wholesale_price":           "v",
			"list_price":                "v",
			"retail_price_min":          "v",
			"retail_price_max":          "v",
			"ed_wholesale_price":        "v",
			"ed_list_price":             "v",
			"supplier_retail_price_min": "v",
			"ed_retail_price_max":       "v",
		},
	}

	filterVariantWhitelist = FilterWhitelist{
		Arrays:   []string{},
		Contains: []string{"name"},
		Equals:   []string{"supplier_id", "name"},
		Status:   []string{"ed_status", "status", "etop_status"},
		Numbers:  []string{"wholesale_price", "list_price", "retail_price_min", "retail_price_max", "ed_wholesale_price", "ed_list_price", "supplier_retail_price_min", "ed_retail_price_max"},
	}

	filterShopProductWhitelist = FilterWhitelist{
		Arrays:   []string{"tags"},
		Contains: []string{"supplier_name", "external_name", "name"},
		Equals:   []string{"external_code", "external_base_id", "external_id", "supplier_id", "collection_id"},
		Status:   []string{"external_status", "ed_status", "status", "etop_status"},
		Numbers:  []string{"retail_price"},
		Dates:    []string{"created_at", "updated_at"},
		Unaccent: []string{"product.name"},

		PrefixOrRename: map[string]string{
			"name":       "sp",
			"status":     "sp",
			"created_at": "sp",
			"updated_at": "sp",

			"product.name": "p.name_norm_ua",
		},
	}
)

func (ft ProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func GetVariantByProductIDs(productIds []int64, filters []cm.Filter) ([]*model.VariantExternalExtended, error) {
	s := x.Table("variant")
	s, _, err := Filters(s, filters, filterVariantWhitelist)
	if err != nil {
		return nil, err
	}
	var variants []*model.VariantExternalExtended

	if err := s.Where("v.deleted_at is NULL").In("v.product_id", productIds).Find((*model.VariantExternalExtendeds)(&variants)); err != nil {
		return nil, err
	}

	return variants, nil
}

func GetProduct(ctx context.Context, query *model.GetProductQuery) error {
	if query.ProductID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	s := x.Table("product").Where("p.deleted_at is NULL")
	if query.SupplierID != 0 {
		s = s.Where("p.supplier_id = ?", query.SupplierID)
	}

	p := new(model.ProductExtended)
	has, err := s.Where("p.id = ?", query.ProductID).Get(p)
	if err != nil {
		return err
	}
	if !has {
		return cm.Error(cm.NotFound, "", nil)
	}

	variants, _ := GetVariantByProductIDs([]int64{p.Product.ID}, []cm.Filter{})

	query.Result = &model.ProductFtVariant{
		ProductExtended: *p,
		Variants:        variants,
	}
	return nil
}

func GetProductsExtended(ctx context.Context, query *model.GetProductsExtendedQuery) error {
	s := x.Table("product").Where("p.deleted_at is NULL")
	if query.SupplierID != 0 {
		s = s.Where("p.supplier_id = ?", query.SupplierID)
	}
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

	var products []*model.ProductExtended
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
		if err := s2.Find((*model.ProductExtendeds)(&products)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.ProductExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}

	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.Product.ID
	}
	variants, err := GetVariantByProductIDs(productIDs, filtersVariant)
	if err != nil {
		return err
	}

	result := make([]*model.ProductFtVariant, len(products))
	hashProductVariant := make(map[int64][]*model.VariantExternalExtended)

	for _, v := range variants {
		hashProductVariant[v.Variant.ProductID] = append(hashProductVariant[v.Variant.ProductID], v)
	}

	for i, p := range products {
		result[i] = &model.ProductFtVariant{
			ProductExtended: *p,
			Variants:        hashProductVariant[p.Product.ID],
		}
	}

	query.Result.Products = result
	return nil
}

func GetVariant(ctx context.Context, query *model.GetVariantQuery) error {
	if query.VariantID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	s := x.Table("variant")
	if query.SupplierID != 0 {
		s = s.Where("v.supplier_id = ?", query.SupplierID)
	}

	v := new(model.VariantExtended)
	has, err := s.Where("v.id = ? AND v.deleted_at is NULL", query.VariantID).Get(v)
	if err != nil {
		return err
	}
	if !has {
		return cm.Error(cm.NotFound, "", nil)
	}

	query.Result = v
	return nil
}

func GetVariantsExtended(ctx context.Context, query *model.GetVariantsExtendedQuery) error {
	if query.SkipPaging {
		// IDs, Codes or EdCodes mut be provided
		if query.IDs == nil && query.EdCodes == nil && query.Codes == nil {
			return cm.Error(cm.InvalidArgument, "Neither id or code provided", nil)
		}
	}

	s := x.Table("variant").Where("v.deleted_at is NULL")
	if query.SupplierID != 0 {
		s = s.Where("v.supplier_id = ?", query.SupplierID)
	}
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
		if err := s2.Find((*model.VariantExtendeds)(&query.Result.Variants)); err != nil {
			return err
		}
	}
	if query.IsPaging() {
		total, err := s.Count(&model.VariantExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetVariantExternals(ctx context.Context, query *model.GetVariantExternalsQuery) error {
	s := x.NewQuery()
	if query.SupplierID != 0 {
		s = s.Where("v.supplier_id = ?", query.SupplierID)
	}
	s = FilterStatus(s, "v.", query.StatusQuery)

	s, _, err := Filters(s, query.Filters, filterProductWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "", "created_at": "", "updated_at": "", "name": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("v.id", query.IDs)
		}
		if err := s2.Find((*model.VariantExternalExtendeds)(&query.Result.Variants)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.VariantExtended{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetVariantExternalsFromID(ctx context.Context, query *model.GetVariantExternalsFromIDQuery) error {
	var variants model.VariantExternalExtendeds
	if err := x.Where("v.id > ?", query.FromID).
		OrderBy("v.id").
		Limit(uint64(query.Limit)).
		Find(&variants); err != nil {
		return err
	}

	query.Result.MaxID = variants[len(variants)-1].Variant.ID
	query.Result.Variants = variants
	return nil
}

func ScanVariantExternals(ctx context.Context, query *model.ScanVariantExternalsQuery) error {
	n := 0
	id := query.FromID
	if query.Limit == 0 {
		query.Limit = 1<<63 - 1
	}
	if query.PageSize == 0 {
		query.PageSize = 1000
	}

	limit := cm.MinInt(query.PageSize, query.Limit)
	if limit == 0 {
		return cm.Error(cm.FailedPrecondition, "Nothing to query", nil)
	}

	// Query the first batch
	q := &model.GetVariantExternalsFromIDQuery{
		FromID: id,
		Limit:  limit,
	}
	if err := bus.Dispatch(ctx, q); err != nil {
		return err
	}

	ch := make(chan model.ScanVariantExternalsResult, 1)
	query.Result = ch
	ch <- q.Result

	// Query the next remaining variants
	go func() {
		for {
			n += len(q.Result.Variants)
			id = q.Result.MaxID
			limit = cm.MinInt(query.PageSize, query.Limit-n)
			if limit == 0 {
				return
			}

			q = &model.GetVariantExternalsFromIDQuery{
				FromID: id,
				Limit:  limit,
			}

			if err := bus.Dispatch(ctx, q); err != nil {
				close(ch)
				if cm.ErrorCode(err) != cm.NotFound {
					ll.Error("Error scanning products", l.Error(err))
				}
				return
			}

			ch <- q.Result
		}
	}()
	return nil
}

func UpdateProduct(ctx context.Context, cmd *model.UpdateProductCommand) error {
	if err := cmd.Product.BeforeUpdate(); err != nil {
		return err
	}

	ft := NewProductFilters("")
	if err := x.
		Where(
			ft.NotDeleted(),
			ft.ByID(cmd.Product.ID),
			ft.BySupplierID(cmd.SupplierID).Optional(),
		).
		ShouldUpdate(cmd.Product); err != nil {
		return err
	}

	query := &model.GetProductQuery{
		ProductID: cmd.Product.ID,
	}
	if pErr := GetProduct(ctx, query); pErr != nil {
		return pErr
	}

	cmd.Result = query.Result
	return nil
}

func UpdateVariant(ctx context.Context, cmd *model.UpdateVariantCommand) error {
	if cmd.Variant == nil || cmd.Variant.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Thiếu VariantID", nil)
	}
	if err := cmd.Variant.BeforeUpdate(); err != nil {
		return err
	}

	s := x.Where("id = ?", cmd.Variant.ID)
	if cmd.SupplierID != 0 {
		s = s.Where("supplier_id = ?", cmd.SupplierID)
	}

	// cmd.Product.BeforeUpdate()
	updated, err := s.Update(cmd.Variant)
	if _err := CheckErrorProductCode(err); _err != nil {
		return _err
	}
	if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	cmd.Result = new(model.VariantExtended)
	if has, err := x.Where("v.id = ?", cmd.Variant.ID).
		Get(cmd.Result); err != nil || !has {
		return cm.Error(cm.Internal, "", err)
	}
	return nil
}

func UpdateVariantImages(ctx context.Context, cmd *model.UpdateVariantImagesCommand) error {
	if cmd.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Thiếu VariantID", nil)
	}

	s := x.Table("variant").Where("id = ? AND deleted_at is NULL", cmd.VariantID)
	if cmd.SupplierID != 0 {
		s = s.Where("supplier_id = ?", cmd.SupplierID)
	}

	if err := s.
		ShouldUpdateMap(M{
			"image_urls": core.Array{V: cmd.ImageURLs},
		}); err != nil {
		return err
	}

	cmd.Result = new(model.VariantExtended)
	if has, err := x.Where("v.id = ?", cmd.VariantID).
		Get(cmd.Result); err != nil || !has {
		return cm.Error(cm.Internal, "", err)
	}
	return nil
}

func UpdateProductImages(ctx context.Context, cmd *model.UpdateProductImagesCommand) error {
	if cmd.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Thiếu ProductID", nil)
	}

	s := x.Table("product").Where("id = ? AND deleted_at is NULL", cmd.ProductID)
	if cmd.SupplierID != 0 {
		s = s.Where("supplier_id = ?", cmd.SupplierID)
	}

	if err := s.
		ShouldUpdateMap(M{
			"image_urls": core.Array{V: cmd.ImageURLs},
		}); err != nil {
		return err
	}

	query := &model.GetProductQuery{
		ProductID: cmd.ProductID,
	}
	if pErr := GetProduct(ctx, query); pErr != nil {
		return pErr
	}

	cmd.Result = query.Result
	return nil
}

func UpdateVariantPrice(ctx context.Context, cmd *model.UpdateVariantPriceCommand) error {
	if cmd.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing VariantID", nil)
	}

	// Always update price even if price is 0
	if updated, err := x.
		Table("variant").
		Where("id = ?", cmd.VariantID).
		UpdateAll().
		Update(cmd.PriceDef); err != nil {
		return err
	} else if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	return nil
}

func UpdateVariants(ctx context.Context, cmd *model.UpdateVariantsCommand) error {
	toUpdates := 0
	ids := make([]int64, len(cmd.Variants))
	errs := make([]error, len(cmd.Variants))
	cmd.Result.Errors = errs

	for _, v := range cmd.Variants {
		if err := v.BeforeUpdate(); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	for i, p := range cmd.Variants {
		ids[i] = p.ID
		if p.ID == 0 {
			errs[i] = cm.Error(cm.NotFound, "", nil)
			continue
		}
		toUpdates++

		wg.Add(1)
		go func(i int, p *model.Variant) {
			defer wg.Done()
			s := x.Where("id = ?", p.ID)
			if cmd.SupplierID != 0 {
				s = s.Where("supplier_id = ?", cmd.SupplierID)
			}

			// p.BeforeUpdate()
			count, err := s.Update(p)
			if err != nil {
				errs[i] = err
			} else if count <= 0 {
				errs[i] = cm.Error(cm.NotFound, "", nil)
			}
		}(i, p)
	}
	wg.Wait()

	if toUpdates == 0 {
		return cm.Error(cm.NotFound, "Nothing to update", nil)
	}

	countErrors := 0
	for _, err := range errs {
		if err != nil {
			countErrors++
		}
	}
	if countErrors == len(errs) {
		return cm.Error(cm.Unknown, "Can not update variants", errs[0])
	}

	if err := x.In("id", ids).
		Find((*model.VariantExtendeds)(&cmd.Result.Variants)); err != nil {
		return cm.Error(cm.Unknown, "Can not retrieve variants", err)
	}
	return nil
}

func UpdateVariantsStatus(ctx context.Context, cmd *model.UpdateVariantsStatusCommand) error {
	s := x.Table("variant").In("id", cmd.IDs).Where("v.deleted_at is NULL")
	if cmd.SupplierID != 0 {
		s = s.Where("supplier_id = ?", cmd.SupplierID)
	}
	s = FilterStatus(s, "", cmd.StatusQuery)

	m := make(map[string]interface{})
	if cmd.Update.SupplierStatus != nil {
		m["ed_status"] = cmd.Update.SupplierStatus
	}
	if cmd.Update.EtopStatus != nil {
		m["etop_status"] = cmd.Update.EtopStatus
	}

	updated, err := s.UpdateMap(m)
	if err != nil {
		return cm.Error(cm.Unknown, "Unable to update status", err)
	}
	cmd.Result.Updated = int(updated)
	return nil
}

func GetShopVariant(ctx context.Context, query *model.GetShopVariantQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	product := new(model.ShopVariantExtended)
	if err := x.Table("shop_variant").
		Where("sv.shop_id = ? AND sv.variant_id = ?", query.ShopID, query.VariantID).
		ShouldGet(product); err != nil {
		return err
	}
	query.Result = product
	return nil
}

func GetShopVariants(ctx context.Context, query *model.GetShopVariantsQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	s := x.Table("shop_variant").
		Where("sv.shop_id = ?", query.ShopID)
	if query.ShopVariantStatus != nil {
		s = s.Where("sv.status = ?", *query.ShopVariantStatus)
	}

	s, _, err := Filters(s, query.Filters, filterShopProductWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"product_id": "", "created_at": "", "updated_at": ""})
		if err != nil {
			return err
		}
		if query.VariantIDs != nil {
			s2 = s2.In("sv.variant_id", query.VariantIDs)
		}
		if err := s2.Find((*model.ShopVariantExtendeds)(&query.Result.Variants)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.ShopVariantExtendeds{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func AddProductsToShopCollection(ctx context.Context, cmd *model.AddProductsToShopCollectionCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.CollectionID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing CollectionID", nil)
	}

	collection := new(model.ShopCollection)
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
		productShopCollection := &model.ProductShopCollection{
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

func RemoveProductsFromShopCollection(ctx context.Context, cmd *model.RemoveProductsFromShopCollectionCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	if deleted, err := x.Table("product_shop_collection").Where("shop_id = ? AND collection_id = ?", cmd.ShopID, cmd.CollectionID).
		In("product_id", cmd.ProductIDs).Delete(&model.ProductShopCollection{}); err != nil {
		return err
	} else if deleted == 0 {
		return cm.Error(cm.NotFound, "", nil)
	} else {
		cmd.Result.Updated = int(deleted)
	}
	return nil
}

func AddShopVariants(ctx context.Context, cmd *model.AddShopVariantsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ids", nil)
	}

	query := &model.GetVariantsExtendedQuery{
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

	xproducts := make([]*model.ShopVariantExtended, len(products))
	for i, p := range products {
		for _, id := range cmd.IDs {
			if p.ID == id {
				xproducts[i] = &model.ShopVariantExtended{
					VariantExtended: *p,
					ShopVariant: &model.ShopVariant{
						ShopID:      cmd.ShopID,
						VariantID:   id,
						RetailPrice: p.ListPrice,
					},
				}
				break
			}
		}
	}
	cmd.Result.Variants = xproducts

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

func UpdateOrInsertShopVariant(sv *model.ShopVariant, productSourceID int64) error {
	return updateOrInsertShopVariant(sv, productSourceID, x)
}

func updateOrInsertShopVariant(sv *model.ShopVariant, productSourceID int64, x Qx) error {
	if sv.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing VariantID", nil)
	}

	if sv.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var shopProducts []*model.ShopProductFtProductFtVariantFtShopVariant
	if err := x.Table("shop_product").
		Where("sp.shop_id = ? AND v.id = ?", sv.ShopID, sv.VariantID).
		Find((*model.ShopProductFtProductFtVariantFtShopVariants)(&shopProducts)); err != nil {
		return err
	}

	if len(shopProducts) == 0 {
		// case: variant has product_source_id which type shop custom
		var variant = new(model.Variant)
		if err := x.Table("variant").Where("id = ? AND product_source_id = ? AND status = 1 AND deleted_at is NULL", sv.VariantID, productSourceID).
			ShouldGet(variant); err != nil {
			return err
		}

		// add to table shop_product
		return inTransaction(func(x Qx) error {
			var product = new(model.Product)
			if err := x.Table("product").Where("id = ? AND deleted_at is NULL", variant.ProductID).ShouldGet(product); err != nil {
				return err
			}

			shopProduct := ConvertProductToShopProduct(product)
			shopProduct.ShopID = sv.ShopID
			if err := x.Table("shop_product").ShouldInsert(shopProduct); err != nil {
				return err
			}

			err, variant := getVariant(sv.VariantID, x)
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
		err, variant := getVariant(sv.VariantID, x)
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

func getVariant(id int64, x Qx) (error, *model.Variant) {
	var variant = new(model.Variant)
	if err := x.Table("variant").Where("id = ?", id).ShouldGet(variant); err != nil {
		return err, nil
	}
	return nil, variant
}

func buildShopVariant(v *model.Variant, sv *model.ShopVariant) *model.ShopVariant {
	if v.ID != sv.VariantID {
		return nil
	}
	return &model.ShopVariant{
		VariantID:   v.ID,
		Name:        v.GetName(),
		Description: cm.Coalesce(sv.Description, v.Description, v.EdDescription),
		DescHTML:    cm.Coalesce(sv.DescHTML, v.DescHTML, v.EdDescHTML),
		ShortDesc:   cm.Coalesce(sv.ShortDesc, v.ShortDesc, v.EdShortDesc),
		ImageURLs:   cm.CoalesceStrings(sv.ImageURLs, v.ImageURLs),
		Note:        sv.Note,
		RetailPrice: cm.CoalesceInt(sv.RetailPrice, v.ListPrice),
		Status:      model.CoalesceStatus3(sv.Status, v.Status),
		ProductID:   cm.CoalesceInt64(sv.ProductID, v.ProductID),
		ShopID:      sv.ShopID,
	}
}

func UpdateShopVariant(ctx context.Context, cmd *model.UpdateShopVariantCommand) error {
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
		productSource := new(model.ProductSource)
		if has, _ := x.Table("product_source").Where("id = ? AND type = ?", cmd.ProductSourceID, model.ProductSourceCustom).
			Get(productSource); has {
			variant := &model.Variant{
				ID:              sv.VariantID,
				ProductSourceID: cmd.ProductSourceID,
			}
			if cmd.CostPrice != 0 {
				variant.CostPrice = cmd.CostPrice
			}
			if cmd.Inventory != 0 {
				variant.QuantityAvailable = cmd.Inventory
				variant.QuantityOnHand = cmd.Inventory
			}
			if len(cmd.Attributes) > 0 {
				variant.Attributes = cmd.Attributes
			}
			if cmd.EdCode != "" {
				variant.EdCode = cmd.EdCode
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

	query := &model.GetShopVariantQuery{
		ShopID:    cmd.ShopID,
		VariantID: cmd.Variant.VariantID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.Error(cm.Internal, "", err)
	}

	cmd.Result = query.Result
	return nil
}

func UpdateShopVariants(ctx context.Context, cmd *model.UpdateShopVariantsCommand) error {
	return cm.ErrTODO
}

func RemoveShopVariants(ctx context.Context, cmd *model.RemoveShopVariantsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	return inTransaction(func(x Qx) error {
		if _, err := x.Table("shop_variant").
			Where("shop_id = ?", cmd.ShopID).In("variant_id", cmd.IDs).Delete(&model.ShopVariant{}); err != nil {
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

func UpdateShopVariantsStatus(ctx context.Context, cmd *model.UpdateShopVariantsStatusCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	updated := 0
	for _, id := range cmd.VariantIDs {
		sp := &model.ShopVariant{
			ShopID:    cmd.ShopID,
			VariantID: id,
			Status:    *cmd.Update.Status,
		}

		if err := UpdateOrInsertShopVariant(sp, cmd.ProductSourceID); err == nil {
			updated++
		}
	}

	if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	cmd.Result.Updated = updated
	return nil
}

func UpdateShopVariantsTags(ctx context.Context, cmd *model.UpdateShopVariantsTagsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	req := cmd.Update
	if err := req.Verify(); err != nil {
		return err
	}

	for i, tag := range req.Adds {
		tag, ok := validate.NormalizeTag(tag)
		if !ok {
			return cm.Error(cm.InvalidArgument, "Invalid tag: "+tag, nil)
		}
		req.Adds[i] = tag
	}
	for i, tag := range req.ReplaceAll {
		tag, ok := validate.NormalizeTag(tag)
		if !ok {
			return cm.Error(cm.InvalidArgument, "Invalid tag: "+tag, nil)
		}
		req.ReplaceAll[i] = tag
	}

	var products []*model.ShopVariant
	if err := x.Where("shop_id = ?", cmd.ShopID).
		In("variant_id", cmd.VariantIDs).
		Find((*model.ShopVariants)(&products)); err != nil {
		return err
	}

	if len(products) == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	countUpdated := 0
	var savedError error
	for _, p := range products {
		tags := req.Patch(p.Tags)
		updated, err := x.
			Table("shop_variant").
			Where("shop_id = ? AND variant_id = ?", cmd.ShopID, p.VariantID).
			UpdateMap(M{
				"tags": x.Opts().Array(model.TagsJoin(tags)),
			})
		if err != nil {
			savedError = err
			continue
		}
		if updated > 0 {
			countUpdated++
		}
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

func UpdateProductsEtopCategory(ctx context.Context, cmd *model.UpdateProductsEtopCategoryCommand) error {
	if cmd.EtopCategoryID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing CategoryID", nil)
	}
	if len(cmd.ProductIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Nothing to update", nil)
	}

	updated, err := x.Table("product").
		In("id", cmd.ProductIDs).
		Where("deleted_at is NULL").
		UpdateMap(M{"category_id": cmd.EtopCategoryID})
	if err != nil {
		return err
	}
	if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	cmd.Result.Updated = int(updated)
	return nil
}

func RemoveProductsEtopCategory(ctx context.Context, cmd *model.RemoveProductsEtopCategoryCommand) error {
	if len(cmd.ProductIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Nothing to remove", nil)
	}

	updated, err := x.Table("product").
		In("id", cmd.ProductIDs).
		Where("deleted_at is NULL").
		UpdateMap(M{"category_id": 0})
	if err != nil {
		return err
	}
	if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	cmd.Result.Updated = int(updated)
	return nil
}

func UpdateProductsStatus(ctx context.Context, cmd *model.UpdateProductsStatusCommand) error {
	s := x.Table("product").In("id", cmd.IDs).Where("deleted_at is NULL")
	if cmd.SupplierID != 0 {
		s = s.Where("supplier_id = ?", cmd.SupplierID)
	}
	s = FilterStatus(s, "", cmd.StatusQuery)

	m := make(map[string]interface{})
	if cmd.Update.EtopStatus != nil {
		m["status"] = cmd.Update.EtopStatus
	}

	updated, err := s.UpdateMap(m)
	if err != nil {
		return cm.Error(cm.Unknown, "Unable to update status", err)
	}
	cmd.Result.Updated = int(updated)
	return nil
}

func AddShopProducts(ctx context.Context, cmd *model.AddShopProductsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ids", nil)
	}

	query := &model.GetProductsExtendedQuery{
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

	xproducts := make([]*model.ShopProduct, len(products))
	for i, p := range products {
		for _, id := range cmd.IDs {
			if p.Product.ID == id {
				xproducts[i] = &model.ShopProduct{
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

func GetShopVariantByProductIDs(productIds []int64) ([]*model.ShopVariantExtended, error) {
	s := x.Table("shop_variant")
	var variants []*model.ShopVariantExtended

	if err := s.In("sv.product_id", productIds).Find((*model.ShopVariantExtendeds)(&variants)); err != nil {
		return nil, err
	}

	return variants, nil
}

func GetShopProduct(ctx context.Context, query *model.GetShopProductQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	if query.ProductID == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	// get products from table shop_product
	var shopProducts []*model.ShopProductFtProductFtVariantFtShopVariant
	if err := x.Table("shop_product").
		Where("sp.product_id = ? AND sp.shop_id = ? AND v.status = 1 AND v.deleted_at is NULL", query.ProductID, query.ShopID).
		Find((*model.ShopProductFtProductFtVariantFtShopVariants)(&shopProducts)); err != nil {
		return err
	}

	if len(shopProducts) > 0 {
		res := convertShopProductFtProductFtVariantFtShopVariants(query.ShopID, shopProducts, nil)
		// map collection_ids
		productIds := make([]int64, len(res))
		for i, p := range res {
			productIds[i] = p.ID
		}
		hashCollections := getCollectionByProductIDs(productIds)
		for _, p := range res {
			p.CollectionIDs = hashCollections[p.ID]
		}
		query.Result = res[0]
		return nil
	}

	// get product from product source
	if query.ProductSourceID == 0 {
		return cm.Error(cm.NotFound, "Not found product", nil)
	}

	var products []*model.ProductFtVariantFtShopProduct
	if err := x.Table("product").
		Where("p.id = ? AND v.status = 1 AND p.product_source_id = ? AND sp.product_id is NULL AND p.deleted_at is NULL AND v.deleted_at is NULL",
			query.ProductID, query.ProductSourceID).
		Find((*model.ProductFtVariantFtShopProducts)(&products)); err != nil {
		return err
	}

	if len(products) == 0 {
		return cm.Error(cm.NotFound, "Not found product", nil)
	}
	res := convertProductFtVariantFtShopProduct(query.ShopID, products, nil)
	query.Result = res[0]
	return nil
}

func getCollectionByProductIDs(ids []int64) map[int64][]int64 {
	var res = make(map[int64][]int64)
	if len(ids) == 0 {
		return res
	}
	var productCollections []*model.ProductShopCollection
	if err := x.Table("product_shop_collection").In("product_id", ids).Find((*model.ProductShopCollections)(&productCollections)); err != nil {
		return res
	}
	for _, pCollection := range productCollections {
		pID := pCollection.ProductID
		res[pID] = append(res[pID], pCollection.CollectionID)
	}
	return res
}

func convertProductFtVariantFtShopProduct(shopID int64, products []*model.ProductFtVariantFtShopProduct, productSource *model.ProductSource) []*model.ShopProductFtVariant {
	result := make([]*model.ShopProductFtVariant, 0, len(products))
	hashProductVariant := make(map[int64][]*model.ShopVariantExtended)
	hashShopProduct := make(map[int64]*model.ShopProduct)
	hashProduct := make(map[int64]*model.Product)
	for _, p := range products {
		pID := p.Product.ID
		hashProduct[pID] = p.Product
		if hashShopProduct[pID] == nil {
			hashShopProduct[pID] = ConvertProductToShopProduct(p.Product)
		}
		hashProductVariant[pID] = append(hashProductVariant[pID], &model.ShopVariantExtended{
			ShopVariant: convertVariantToShopVariant(shopID, p.Variant),
			VariantExtended: model.VariantExtended{
				Variant:         p.Variant,
				Product:         p.Product,
				VariantExternal: p.VariantExternal,
			},
		})
	}

	for i, value := range hashShopProduct {
		pdSourceID := hashProduct[i].ProductSourceID
		value.ProductSourceID = pdSourceID
		if productSource != nil && pdSourceID == productSource.ID {
			value.ProductSourceName = productSource.Name
			value.ProductSourceType = productSource.Type
		}
		result = append(result, &model.ShopProductFtVariant{
			ShopProduct: value,
			Product:     hashProduct[i],
			Variants:    hashProductVariant[i],
		})
	}
	return result
}

func convertShopProductFtProductFtVariantFtShopVariants(shopID int64, products []*model.ShopProductFtProductFtVariantFtShopVariant, productSource *model.ProductSource) []*model.ShopProductFtVariant {
	result := make([]*model.ShopProductFtVariant, 0, len(products))
	hashProductVariant := make(map[int64][]*model.ShopVariantExtended)
	hashShopProduct := make(map[int64]*model.ShopProduct)
	hashProduct := make(map[int64]*model.Product)
	for _, p := range products {
		pID := p.ShopProduct.ProductID
		hashShopProduct[pID] = p.ShopProduct
		hashProduct[pID] = p.Product
		if p.ShopVariant.VariantID != 0 {
			hashProductVariant[pID] = append(hashProductVariant[pID], &model.ShopVariantExtended{
				ShopVariant: p.ShopVariant,
				VariantExtended: model.VariantExtended{
					Variant:         p.Variant,
					Product:         p.Product,
					VariantExternal: p.VariantExternal,
				},
			})
		} else {
			hashProductVariant[pID] = append(hashProductVariant[pID], &model.ShopVariantExtended{
				ShopVariant: convertVariantToShopVariant(shopID, p.Variant),
				VariantExtended: model.VariantExtended{
					Variant:         p.Variant,
					Product:         p.Product,
					VariantExternal: p.VariantExternal,
				},
			})
		}
	}

	for i, value := range hashShopProduct {
		pdSourceID := hashProduct[i].ProductSourceID
		value.ProductSourceID = pdSourceID
		if productSource != nil && pdSourceID == productSource.ID {
			value.ProductSourceName = productSource.Name
			value.ProductSourceType = productSource.Type
		}
		result = append(result, &model.ShopProductFtVariant{
			ShopProduct: value,
			Product:     hashProduct[i],
			Variants:    hashProductVariant[i],
		})
	}

	return result
}

func convertVariantToShopVariant(shopID int64, v *model.Variant) *model.ShopVariant {
	return &model.ShopVariant{
		ShopID:      shopID,
		VariantID:   v.ID,
		Name:        v.GetName(),
		Description: cm.Coalesce(v.Description, v.EdDescription),
		DescHTML:    cm.Coalesce(v.DescHTML, v.EdDescHTML),
		ShortDesc:   cm.Coalesce(v.ShortDesc, v.EdShortDesc),
		ImageURLs:   v.ImageURLs,
		Note:        "",
		RetailPrice: v.ListPrice,
		Status:      v.Status,
	}
}

func ConvertProductToShopProduct(p *model.Product) *model.ShopProduct {
	return &model.ShopProduct{
		ProductID:   p.ID,
		Name:        cm.Coalesce(p.Name, p.EdName),
		Description: cm.Coalesce(p.Description, p.EdDescription),
		DescHTML:    cm.Coalesce(p.DescHTML, p.EdDescHTML),
		ShortDesc:   cm.Coalesce(p.ShortDesc, p.EdShortDesc),
		ImageURLs:   p.ImageURLs,
		Status:      p.Status,
	}
}

func GetShopProducts(ctx context.Context, query *model.GetShopProductsQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var productSource = new(model.ProductSource)
	ok, err := x.Table("product_source").Where("id = ?", query.ProductSourceID).Get(productSource)
	if err != nil {
		return err
	}
	if !ok {
		return nil // not found
	}

	// get products from table shop_product
	var shopProducts []*model.ShopProductFtProductFtVariantFtShopVariant

	s := x.Table("shop_product").
		Where("sp.shop_id = ? AND v.status = 1 AND v.deleted_at is NULL", query.ShopID)
	if query.ShopProductStatus != nil {
		s = s.Where("sp.status = ?", *query.ShopProductStatus)
	}
	if query.ProductIDs != nil {
		s = s.In("sp.product_id", query.ProductIDs)
	}
	s, _, err = Filters(s, query.Filters, filterShopProductWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"product_id": "sp.product_id", "created_at": "sp.created_at", "updated_at": "sp.updated_at"})
		if err != nil {
			return err
		}
		if err := s2.Find((*model.ShopProductFtProductFtVariantFtShopVariants)(&shopProducts)); err != nil {
			return err
		}
	}
	res := convertShopProductFtProductFtVariantFtShopVariants(query.ShopID, shopProducts, productSource)
	// map collection_ids
	productIds := make([]int64, len(res))
	for i, p := range res {
		productIds[i] = p.ID
	}
	hashCollections := getCollectionByProductIDs(productIds)
	for _, p := range res {
		p.CollectionIDs = hashCollections[p.ID]
	}

	// get product from product source
	if query.ProductSourceID == 0 {
		query.Result.Products = res
		return nil
	}
	var products []*model.ProductFtVariantFtShopProduct
	s3 := x.Table("product").
		Where("p.product_source_id = ? AND v.status = 1 AND sp.product_id is NULL AND p.deleted_at is NULL AND v.deleted_at is NULL", query.ProductSourceID)
	s3 = s3.In("p.id", productIds)
	if err != nil {
		return err
	}
	{
		s4 := s3.Clone()
		// s4 = LimitSort(s4, query.Paging, "product_id", "created_at", "updated_at")
		if err := s4.Find((*model.ProductFtVariantFtShopProducts)(&products)); err != nil {
			return err
		}
	}
	res2 := convertProductFtVariantFtShopProduct(query.ShopID, products, productSource)
	query.Result.Products = append(res, res2...)
	query.Result.Total = len(query.Result.Products)
	return nil
}

// Get all from source kiotviet + source shop
func GetAllShopVariants(ctx context.Context, query *model.GetAllShopVariantsQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	// get products from table shop_product
	var shopProducts []*model.ShopProductFtProductFtVariantFtShopVariant

	s := x.Table("shop_product").
		Where("sp.shop_id = ? AND v.status = 1 AND v.deleted_at is NULL", query.ShopID)
	if query.VariantIDs != nil {
		s = s.In("v.id ", query.VariantIDs)
	}
	if err := s.Find((*model.ShopProductFtProductFtVariantFtShopVariants)(&shopProducts)); err != nil {
		return err
	}
	res := make([]*model.ShopVariantExtended, len(shopProducts))
	for i, p := range shopProducts {
		if p.ShopVariant.VariantID != 0 {
			res[i] = &model.ShopVariantExtended{
				ShopVariant: p.ShopVariant,
				ShopProduct: p.ShopProduct,
				VariantExtended: model.VariantExtended{
					Variant: p.Variant,
					Product: p.Product,
				},
			}
		} else {
			res[i] = &model.ShopVariantExtended{
				ShopVariant: convertVariantToShopVariant(query.ShopID, p.Variant),
				ShopProduct: p.ShopProduct,
				VariantExtended: model.VariantExtended{
					Variant: p.Variant,
					Product: p.Product,
				},
			}
		}
	}
	// get product from product source
	if query.ProductSourceID == 0 {
		query.Result.Variants = res
		return nil
	}

	var products []*model.ProductFtVariantFtShopProduct
	s3 := x.Table("product").
		Where("p.product_source_id = ? AND v.status = 1 AND sp.product_id is NULL AND p.deleted_at is NULL AND v.deleted_at is NULL", query.ProductSourceID)
	if query.VariantIDs != nil {
		s3 = s3.In("v.id", query.VariantIDs)
	}
	if err := s3.Find((*model.ProductFtVariantFtShopProducts)(&products)); err != nil {
		return err
	}
	res2 := make([]*model.ShopVariantExtended, len(products))
	for i, p := range products {
		res2[i] = &model.ShopVariantExtended{
			ShopVariant: convertVariantToShopVariant(query.ShopID, p.Variant),
			VariantExtended: model.VariantExtended{
				Variant: p.Variant,
				Product: p.Product,
			},
		}
	}
	query.Result.Variants = append(res, res2...)
	return nil
}

func RemoveShopProducts(ctx context.Context, cmd *model.RemoveShopProductsCommand) error {
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
			productsCount, err = s1.Count(&model.ProductFtShopProduct{})
			if err != nil {
				return err
			}
		}

		deletedCount, err := x.Table("shop_product").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			Delete(&model.ShopProduct{})
		if err != nil {
			return nil
		}

		if _, err2 := x.Table("shop_variant").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			Delete(&model.ShopVariant{}); err2 != nil {
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

func UpdateOrInsertShopProduct(sp *model.ShopProduct, productSourceID int64) error {
	return updateOrInsertShopProduct(sp, productSourceID, x)
}

func updateOrInsertShopProduct(sp *model.ShopProduct, productSourceID int64, x Qx) error {
	if sp.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductID", nil)
	}

	if sp.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var shopProduct = new(model.ShopProduct)
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
	var product = new(model.Product)
	if err := x.Table("product").Where("id = ? AND product_source_id = ? AND deleted_at is NULL", sp.ProductID, productSourceID).
		ShouldGet(product); err != nil {
		return err
	}

	if err := x.Table("shop_product").ShouldInsert(sp); err != nil {
		return err
	}
	return nil
}

func UpdateShopProduct(ctx context.Context, cmd *model.UpdateShopProductCommand) error {
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

		productSource := new(model.ProductSource)
		if has, _ := x.Table("product_source").Where("id = ? AND type = ?", cmd.ProductSourceID, model.ProductSourceCustom).
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

	query := &model.GetShopProductQuery{
		ShopID:    cmd.ShopID,
		ProductID: cmd.Product.ProductID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.Error(cm.Internal, "", err)
	}

	cmd.Result = query.Result
	return nil
}

func UpdateShopProductsStatus(ctx context.Context, cmd *model.UpdateShopProductsStatusCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	updated := 0
	for _, id := range cmd.ProductIDs {
		sp := &model.ShopProduct{
			ShopID:    cmd.ShopID,
			ProductID: id,
			Status:    *cmd.Update.Status,
		}

		if err := UpdateOrInsertShopProduct(sp, cmd.ProductSourceID); err == nil {
			updated++
		}
	}

	if updated == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}

	cmd.Result.Updated = updated
	return nil
}

func UpdateShopProductsTags(ctx context.Context, cmd *model.UpdateShopProductsTagsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	req := cmd.Update
	if err := req.Verify(); err != nil {
		return err
	}

	var products []*model.ShopProduct
	if err := x.Where("shop_id = ?", cmd.ShopID).
		In("product_id", cmd.ProductIDs).
		Find((*model.ShopProducts)(&products)); err != nil {
		return err
	}

	productMap := make(map[int64]*model.ShopProduct)
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
		sp := &model.ShopProduct{
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

func GetProducts(ctx context.Context, query *model.GetProductsQuery) error {
	if query.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
	}

	s := x.Table("product").
		Where("product_source_id = ?", query.ProductSourceID)

	count := 0
	if query.EdCodes != nil {
		s = s.In("ed_code", query.EdCodes)
		count++
	}
	if query.NameNormUas != nil {
		s = s.In("name_norm_ua", query.NameNormUas)
		count++
	}
	if count == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	if !query.IncludeDeleted {
		s = s.Where("deleted_at IS NULL")
	}
	if query.ExcludeEdCode {
		s = s.Where("ed_code IS NULL")
	}
	return s.Find((*model.Products)(&query.Result.Products))
}

func GetVariants(ctx context.Context, query *model.GetVariantsQuery) error {
	if query.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
	}

	s := x.Table("variant").
		Where("product_source_id = ?", query.ProductSourceID)

	if query.Inclusive {
		if query.EdCodes == nil && query.AttrNorms == nil {
			return cm.Error(cm.InvalidArgument, "Must provide both params when using with inclusive", nil)
		}
		s = s.Where(sq.Or{
			sq.In("ed_code", query.EdCodes),
			sq.Ins([]string{"product_id", "attr_norm_kv"}, query.AttrNorms...),
		})

	} else {
		count := 0
		if query.EdCodes != nil {
			s = s.In("ed_code", query.EdCodes)
			count++
		}
		if query.AttrNorms != nil {
			s = s.In("attr_norm_kv", query.AttrNorms)
			count++
		}
		if count == 0 {
			return cm.Error(cm.InvalidArgument, "Missing required params", nil)
		}
	}

	if !query.IncludeDeleted {
		s = s.Where("deleted_at IS NULL")
	}
	return s.Find((*model.Variants)(&query.Result.Variants))
}
