package imcsv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	topintshop "o.o/api/top/int/shop"
	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
	catalogmodelx "o.o/backend/com/main/catalog/modelx"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/validate"
	apishop "o.o/backend/pkg/etop/api/shop"
	"o.o/capi/dot"
)

var maxPaging = meta.Paging{Limit: 5000}

// - check if product code exists
//   + if not exist, product information must not be empty
// - load all categories and collections from database
// - if any category or collection does not exist, create it and fill the id
func (im *Import) loadAndCreateProducts(
	ctx context.Context,
	schema imcsv.Schema,
	idx indexes,
	mode Mode,
	codeMode CodeMode,
	shop *identitymodel.Shop,
	rowProducts []*RowProduct,
	debug Debug,
	user *identitymodelx.SignedInUser,
) (stocktakeId dot.ID, msgs []string, _errs []error, _cellErrs []error, _err error) {
	var categories *Categories
	// var collections map[string]*catalogmodel.ShopCollection
	var products map[string]*catalog.ShopProduct
	var variantByCode, variantByAttr map[string]*catalogmodel.ShopVariant
	chErr := make(chan error)
	go func() {
		var err error
		categories, err = loadCategories(ctx, shop.ID)
		if err != nil {
			err = cm.Error(cm.Internal, "", err).
				WithMeta("step", "category")
		}
		chErr <- err
	}()
	go func() {
		var err error
		// collections, err = loadCollections(ctx, shop.ID)
		// if err != nil {
		// 	err = cm.Error(cm.Internal, "", err).
		// 		WithMeta("step", "collection")
		// }
		chErr <- err
	}()
	go func() {
		var err error
		{
			productKeys := make([]string, len(rowProducts))
			for i, p := range rowProducts {
				productKeys[i] = p.GetProductKey()
			}
			products, err = im.loadProducts(ctx, codeMode, shop.ID, productKeys)
			if err != nil {
				err = cm.Error(cm.Internal, "", err).
					WithMeta("step", "product")
			}
			chErr <- err
		}
		{
			codes := make([]string, len(rowProducts))
			attrNorms := make([]interface{}, 0, len(rowProducts)*2)
			for i, p := range rowProducts {
				codes[i] = p.VariantCode

				product := products[p.GetProductKey()]
				if product != nil {
					attrNorms = append(attrNorms, product.ProductID, p.GetVariantAttrNorm())
				}
			}
			variantByCode, variantByAttr, err = im.loadVariants(ctx, codeMode, shop.ID, codes, attrNorms)
			if err != nil {
				err = cm.Error(cm.Internal, "", err).
					WithMeta("step", "variant")
			}
			chErr <- err
		}
	}()
	for i := 0; i < 4; i++ {
		if err := <-chErr; err != nil {
			_err = err
		}
	}
	if _err != nil {
		return
	}

	// Validate product name for creating new product
	{
		_productCodes := make(map[string]struct{})
		for _, rowProduct := range rowProducts {
			if p := products[rowProduct.ProductCode]; p == nil {
				if rowProduct.ProductName == "" {
					if _, ok := _productCodes[rowProduct.ProductCode]; ok {
						// do not duplicate the error
						continue
					}
					err := imcsv.CellError(idx.indexer, rowProduct.RowIndex, idx.productName, "Mã sản phẩm %v chưa tồn tại, cần cung cấp tên sản phẩm để tạo sản phẩm mới.", rowProduct.ProductCode)
					_cellErrs = append(_cellErrs, err)
					if len(_cellErrs) >= MaxCellErrors {
						return
					}
					_productCodes[rowProduct.ProductCode] = struct{}{}
				}
			}
		}
		if len(_cellErrs) > 0 {
			return
		}
	}

	// validate variant ed_code and variant attribute uniqueness
	{
		for _, rowProduct := range rowProducts {
			if rowProduct.VariantCode != "" && variantByCode[rowProduct.VariantCode] != nil {
				err := imcsv.CellError(idx.indexer, rowProduct.RowIndex, idx.variantCode, `Mã phiên bản sản phẩm "%v" đã tồn tại. Vui lòng sử dụng mã khác hoặc xoá phiên bản này.`, rowProduct.VariantCode)
				_cellErrs = append(_cellErrs, err)
				if len(_cellErrs) >= MaxCellErrors {
					return
				}
			}

			if v := variantByAttr[rowProduct.GetVariantAttrNorm()]; v != nil {
				err := imcsv.CellError(idx.indexer, rowProduct.RowIndex, idx.attributes, `Một phiên bản của sản phẩm "%v" với thuộc tính "%v" đã tồn tại. Vui lòng sử dụng thuộc tính khác hoặc xoá phiên bản này.`, rowProduct.GetProductCodeOrName(), v.Attributes.ShortLabel())
				_cellErrs = append(_cellErrs, err)
				if len(_cellErrs) >= MaxCellErrors {
					return
				}
			}
		}
		if len(_cellErrs) > 0 {
			return
		}
	}

	// Create new categories and collections
	for _, rowProduct := range rowProducts {
		{
			var category *catalogmodel.ShopCategory
			var err error
			cc := normalizeCategory(rowProduct.Category)
			category, msgs, err = ensureCategory(ctx, msgs, categories.Sort, shop, rowProduct.Category, cc)
			if err != nil {
				err = imcsv.CellErrorWithCode(idx.indexer, cm.Internal, err, rowProduct.RowIndex, -1, "Không thể tạo danh mục \"%v\": %v", rowProduct.Category, err)
				_errs = append(_errs, err)
			}
			if category != nil {
				rowProduct.categoryID = category.ID
			}
		}

		// TODO: create collection
		//
		// {
		// 	if len(rowProduct.Collections) > 0 {
		// 		rowProduct.collectionIDs = make([]dot.ID, len(rowProduct.Collections))
		// 	}
		// 	for i, name := range rowProduct.Collections {
		// 		if debug.FailPercent != 0 && isRandomFail(debug.FailPercent) {
		// 			_errs = append(_errs, imcsv.CellErrorWithCode(idx.indexer, cm.Internal, errors.New("random error"), rowProduct.RowIndex, -1, "Random error for development"))
		// 			continue
		// 		}
		//
		// 		nameNorm := validate.NormalizeSearch(name)
		// 		collection := collections[nameNorm]
		// 		if collection == nil {
		//
		// 			collection = &catalogmodel.ShopCollection{
		// 				ShopID: shop.ID,
		// 				Name:   name,
		// 			}
		// 			createCollectionCmd := &catalogmodelx.CreateShopCollectionCommand{
		// 				Collection: collection,
		// 			}
		// 			if err := bus.Dispatch(ctx, createCollectionCmd); err != nil {
		// 				err = imcsv.CellErrorWithCode(idx.indexer, cm.Internal, err, rowProduct.RowIndex, -1, "Không thể tạo bộ sưu tập \"%v\": %v", name, err)
		// 				_errs = append(_errs, err)
		// 				continue
		// 			}
		//
		// 			msgs = append(msgs, "Đã tạo bộ sưu tập "+name)
		// 			collection = createCollectionCmd.Result
		// 			collections[nameNorm] = collection
		// 		}
		// 		rowProduct.collectionIDs[i] = collection.ID
		// 	}
		// }
	}
	if len(_errs) > 0 {
		return
	}

	now := time.Now()
	// Create new products/variants and add them to corresponding categories/collection
	var createStocktakeCmd = &apishop.CreateStocktakeEndpoint{
		CreateStocktakeRequest: &topintshop.CreateStocktakeRequest{
			Note: "Tạo phiếu quản lý tồn kho theo file import",
		},
	}
	createStocktakeCmd.Context.Shop = shop
	createStocktakeCmd.Context.User = user
	var stocktakeLines []*topintshop.StocktakeLine
	for _, rowProduct := range rowProducts {
		if debug.FailPercent != 0 && isRandomFail(debug.FailPercent) {
			_errs = append(_errs, imcsv.CellErrorWithCode(idx.indexer, cm.Internal, errors.New("random error"), rowProduct.RowIndex, -1, "Random error for development"))
			continue
		}

		variantReq := rowToCreateVariant(rowProduct, now)
		if p := products[rowProduct.GetProductKey()]; p != nil {
			variantReq.ProductId = p.ProductID

		} else {
			productReq := rowToCreateProduct(rowProduct, now)
			createProductCmd := &apishop.CreateProductEndpoint{
				CreateProductRequest: productReq,
			}
			createProductCmd.Context.Shop = shop
			if err := apishop.ProductServiceImpl.CreateProduct(ctx, createProductCmd); err != nil {
				err = imcsv.CellErrorWithCode(idx.indexer, cm.Unknown, err, rowProduct.RowIndex, -1,
					`Không thể tạo sản phẩm "%v": %v`,
					rowProduct.GetProductNameOrCode(), err).
					WithMeta("product_code", rowProduct.ProductCode)
				_errs = append(_errs, err)
				continue
			}
			variantReq.ProductId = createProductCmd.Result.Id
		}

		createVariantCmd := &apishop.CreateVariantEndpoint{
			CreateVariantRequest: variantReq,
		}
		createVariantCmd.Context.Shop = shop
		if err := apishop.ProductServiceImpl.CreateVariant(ctx, createVariantCmd); err != nil {
			err = imcsv.CellErrorWithCode(idx.indexer, cm.Unknown, err, rowProduct.RowIndex, -1,
				`Không thể tạo phiên bản "%v" của sản phẩm "%v": %v`,
				variantReq.Name, rowProduct.GetProductNameOrCode(), err).
				WithMeta("product_code", rowProduct.ProductCode).
				WithMeta("variant_code", rowProduct.VariantCode)
			_errs = append(_errs, err)
			continue
		}
		if rowProduct.CostPrice > 0 {
			updatePriceCmd := &apishop.UpdateInventoryVariantCostPriceEndpoint{
				UpdateInventoryVariantCostPriceRequest: &topintshop.UpdateInventoryVariantCostPriceRequest{
					VariantId: createVariantCmd.Result.Id,
					CostPrice: rowProduct.CostPrice,
				},
			}
			updatePriceCmd.Context.Shop = shop
			if err := apishop.InventoryServiceImpl.UpdateInventoryVariantCostPrice(ctx, updatePriceCmd); err != nil {
				err = imcsv.CellErrorWithCode(idx.indexer, cm.Unknown, err, rowProduct.RowIndex, -1,
					`Không thể cập nhập giá cho phiên bản "%v" của sản phẩm "%v": %v`,
					variantReq.Name, rowProduct.GetProductNameOrCode(), err).
					WithMeta("product_code", rowProduct.ProductCode).
					WithMeta("variant_code", rowProduct.VariantCode)
				_errs = append(_errs, err)
			}
		}

		if rowProduct.QuantityAvail != 0 {
			var imageURL string
			if len(createVariantCmd.CreateVariantRequest.ImageUrls) > 0 {
				imageURL = createVariantCmd.CreateVariantRequest.ImageUrls[0]
			}

			createStocktakeCmd.TotalQuantity += rowProduct.QuantityAvail
			// Prepare create stocktake
			var stocktakeLine = &topintshop.StocktakeLine{
				ProductId:   createVariantCmd.CreateVariantRequest.ProductId,
				ProductName: rowProduct.ProductName,
				VariantName: createVariantCmd.CreateVariantRequest.Name,
				VariantId:   createVariantCmd.Result.Id,
				OldQuantity: 0,
				NewQuantity: rowProduct.QuantityAvail,
				Code:        createVariantCmd.CreateVariantRequest.Code,
				ImageUrl:    imageURL,
				CostPrice:   rowProduct.CostPrice,
				Attributes:  createVariantCmd.CreateVariantRequest.Attributes,
			}
			stocktakeLines = append(stocktakeLines, stocktakeLine)
		}

		// Fake the product, so subsequent create variant requests reuse the created product
		products[rowProduct.GetProductKey()] = &catalog.ShopProduct{
			ProductID: variantReq.ProductId,
		}

		var msg string
		if variantReq.Name != "" {
			msg = fmt.Sprintf("Đã tạo sản phẩm \"%v\" - \"%v\"", rowProduct.ProductName, variantReq.Name)
		} else {
			msg = fmt.Sprintf("Đã tạo sản phẩm \"%v\"", rowProduct.ProductName)
		}
		msgs = append(msgs, msg)

		productIDs := []dot.ID{variantReq.ProductId}
		if rowProduct.categoryID != 0 {
			updateProductsCategoryCmd := &catalogmodelx.UpdateProductsShopCategoryCommand{
				CategoryID: rowProduct.categoryID,
				ProductIDs: productIDs,
				ShopID:     shop.ID,
			}
			if err := bus.Dispatch(ctx, updateProductsCategoryCmd); err != nil {
				err = imcsv.CellErrorWithCode(idx.indexer, cm.Internal, err, rowProduct.RowIndex, -1,
					`Không thể thêm sản phẩm "%v" vào danh mục: %v`,
					rowProduct.GetProductNameOrCode(), err).
					WithMeta("product_code", rowProduct.ProductCode).
					WithMeta("variant_code", rowProduct.VariantCode)
				_errs = append(_errs, err)
			}
		}

		// TODO: add product to collection
		//
		// for _, collectionID := range rowProduct.collectionIDs {
		// 	updateProductsCollectionCmd := &catalogmodelx.AddProductsToShopCollectionCommand{
		// 		ShopID:       shop.ID,
		// 		ProductIDs:   productIDs,
		// 		CollectionID: collectionID,
		// 	}
		// 	if err := bus.Dispatch(ctx, updateProductsCollectionCmd); err != nil {
		// 		err = imcsv.CellErrorWithCode(idx.indexer, cm.Internal, err, rowProduct.RowIndex, -1,
		// 			`Không thể thêm sản phẩm "%v" vào bộ sưu tập: %v`,
		// 			rowProduct.GetProductNameOrCode(), err).
		// 			WithMeta("product_code", rowProduct.ProductCode).
		// 			WithMeta("variant_code", rowProduct.VariantCode)
		// 		_errs = append(_errs, err)
		// 	}
		// }
	}
	createStocktakeCmd.Lines = stocktakeLines
	stocktakeId = 0
	if len(stocktakeLines) > 0 {
		err := apishop.StocktakeServiceImpl.CreateStocktake(ctx, createStocktakeCmd)
		if err != nil {
			_errs = append(_errs, err)
		} else {
			stocktakeId = createStocktakeCmd.Result.Id
		}
	}

	return
}

type Categories struct {
	List []*catalogmodel.ShopCategory
	Map  map[dot.ID]*catalogmodel.ShopCategory
	Sort map[[3]string]*catalogmodel.ShopCategory
}

func normalizeCategory(cc [3]string) (res [3]string) {
	for i := 0; i < len(cc); i++ {
		res[i] = validate.NormalizeSearch(cc[i])
	}
	return res
}

func loadCategories(ctx context.Context, shopID dot.ID) (*Categories, error) {
	query := &catalogmodelx.GetProductSourceCategoriesQuery{
		ShopID: shopID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	categories := query.Result.Categories

	mapCategory := make(map[dot.ID]*catalogmodel.ShopCategory)
	for _, c := range categories {
		mapCategory[c.ID] = c
	}
	return &Categories{
		List: categories,
		Map:  mapCategory,
		Sort: sortCategories(mapCategory),
	}, nil
}

func sortCategories(mapCategory map[dot.ID]*catalogmodel.ShopCategory) map[[3]string]*catalogmodel.ShopCategory {
	categories := make(map[[3]string]*catalogmodel.ShopCategory)
	for _, c := range mapCategory {
		cc, ok := buildCategoryHierarchy(mapCategory, c)
		if ok {
			categories[cc] = c
		}
	}
	return categories
}

func buildCategoryHierarchy(mapCategory map[dot.ID]*catalogmodel.ShopCategory, category *catalogmodel.ShopCategory) (res [3]string, ok bool) {
	i := 0
	res[0] = validate.NormalizeSearch(category.Name)
	for category.ParentID != 0 {
		i++
		if i >= 3 {
			return res, false
		}

		parent := mapCategory[category.ParentID]
		if parent == nil {
			return res, false
		}
		category = parent
		res[i] = validate.NormalizeSearch(category.Name)
	}
	return res, true
}

// Load all collections and sort them into normalized map
// func loadCollections(ctx context.Context, shopID dot.ID) (map[string]*catalogmodel.ShopCollection, error) {
// 	query := &catalogmodelx.GetShopCollectionsQuery{
// 		ShopID: shopID,
// 	}
// 	if err := bus.Dispatch(ctx, query); err != nil {
// 		return nil, err
// 	}
// 	mapCollection := make(map[string]*catalogmodel.ShopCollection)
// 	for _, collection := range query.Result.Collections {
// 		name := validate.NormalizeSearch(collection.Name)
// 		mapCollection[name] = collection
// 	}
// 	return mapCollection, nil
// }

func (im *Import) loadProducts(ctx context.Context, codeMode CodeMode, shopID dot.ID, keys []string) (map[string]*catalog.ShopProduct, error) {
	s := im.shopProductStore(ctx).ShopID(shopID)
	useCode := codeMode == CodeModeUseCode
	if useCode {
		s.Codes(keys...)
	} else {
		// only query products with ed_code is null
		s.Where(s.FtShopProduct.ByCode("").Nullable())
		s.ByNameNormUas(keys...)
	}
	products, err := s.WithPaging(maxPaging).ListShopProductsDB()
	if err != nil {
		return nil, err
	}

	mapProducts := make(map[string]*catalog.ShopProduct)
	for _, p := range products {
		product := convert.Convert_catalogmodel_ShopProduct_catalog_ShopProduct(p, nil)
		if useCode {
			mapProducts[p.Code] = product
		} else {
			// Use p.NameNormUa here instead of p.NameNorm because NameNorm
			// is sorted by Postgres while normalizing keeps the word order.
			mapProducts[p.NameNormUa] = product
		}
	}
	return mapProducts, nil
}

// different to loadProducts, we query variants with both ed_code and
// attr_norm_kv to make sure that there is no duplicate in variant
func (im *Import) loadVariants(
	ctx context.Context,
	codeMode CodeMode,
	shopID dot.ID,
	codes []string,
	attrNorms []interface{},
) (
	variantByCode map[string]*catalogmodel.ShopVariant,
	variantByAttr map[string]*catalogmodel.ShopVariant,
	_ error,
) {
	s := im.shopVariantStore(ctx).ShopID(shopID)
	args := catalogsqlstore.ListShopVariantsForImportArgs{
		Codes:     codes,
		AttrNorms: attrNorms,
	}
	useCode := codeMode == CodeModeUseCode
	if useCode {
		args.Codes = codes
	}
	variants, err := s.WithPaging(maxPaging).ListShopVariantsDB()
	if err != nil {
		return nil, nil, err
	}

	variantByCode = make(map[string]*catalogmodel.ShopVariant)
	variantByAttr = make(map[string]*catalogmodel.ShopVariant)
	for _, v := range variants {
		if useCode && v.Code != "" {
			variantByCode[v.Code] = v
		}
		variantByAttr[v.AttrNormKv] = v
	}
	return
}

func ensureCategory(
	ctx context.Context,
	msgs []string,
	categories map[[3]string]*catalogmodel.ShopCategory,
	shop *identitymodel.Shop,
	names [3]string,
	cc [3]string,
) (*catalogmodel.ShopCategory, []string, error) {
	if cc == [3]string{} {
		return nil, msgs, nil
	}
	category := categories[cc]
	if category == nil {
		ccParent := [3]string{cc[1], cc[2]}
		namesNext := [3]string{names[1], names[2]}

		var parent *catalogmodel.ShopCategory
		var err error
		parent, msgs, err = ensureCategory(ctx, msgs, categories, shop, namesNext, ccParent)
		if err != nil {
			return nil, msgs, err
		}

		cmd := &catalogmodelx.CreateShopCategoryCommand{
			ShopID: shop.ID,
			Name:   names[0],
		}
		if parent != nil {
			cmd.ParentID = parent.ID
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return nil, msgs, err
		}
		msgs = append(msgs, "Đã tạo danh mục "+names[0])
		category = cmd.Result
		categories[cc] = category
	}
	return category, msgs, nil
}
