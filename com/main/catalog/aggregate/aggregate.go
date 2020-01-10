package aggregate

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ catalog.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                    *cmsql.Database
	shopProduct           sqlstore.ShopProductStoreFactory
	shopVariant           sqlstore.ShopVariantStoreFactory
	shopCategory          sqlstore.ShopCategoryStoreFactory
	shopCollection        sqlstore.ShopCollectionStoreFactory
	shopProductCollection sqlstore.ShopProductCollectionStoreFactory
	shopBrand             sqlstore.ShopBrandStoreFactory
	shopVariantSupplier   sqlstore.ShopVariantSupplierStoreFactory
	eventBus              capi.EventBus
}

func New(eventBus capi.EventBus, db *cmsql.Database) *Aggregate {
	return &Aggregate{
		db:                    db,
		shopProduct:           sqlstore.NewShopProductStore(db),
		shopVariant:           sqlstore.NewShopVariantStore(db),
		shopCategory:          sqlstore.NewShopCategoryStore(db),
		shopCollection:        sqlstore.NewShopCollectionStore(db),
		shopProductCollection: sqlstore.NewShopProductCollectionStore(db),
		shopBrand:             sqlstore.NewShopBrandStore(db),
		shopVariantSupplier:   sqlstore.NewVariantSupplierStore(db),
		eventBus:              eventBus,
	}
}

func (a *Aggregate) MessageBus() catalog.CommandBus {
	b := bus.New()
	return catalog.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShopProduct(ctx context.Context, args *catalog.CreateShopProductArgs) (*catalog.ShopProductWithVariants, error) {
	if args.BrandID != 0 {
		_, err := a.shopBrand(ctx).ShopID(args.ShopID).ID(args.BrandID).GetShopBrand()
		if err != nil {
			return nil, cm.MapError(err).
				Mapf(cm.NotFound, cm.InvalidArgument, "Mã thương hiệu không tồn tại").
				Throw()
		}
	}
	product := &catalog.ShopProduct{
		ExternalID:   args.ExternalID,
		ExternalCode: args.ExternalCode,
		PartnerID:    args.PartnerID,

		ProductID:   cm.NewID(),
		ShopID:      args.ShopID,
		Code:        args.Code,
		Name:        args.Name,
		Unit:        args.Unit,
		ImageURLs:   args.ImageURLs,
		Note:        args.Note,
		ShortDesc:   args.ShortDesc,
		Description: args.Description,
		DescHTML:    args.DescHTML,
		CostPrice:   args.CostPrice,
		ListPrice:   args.ListPrice,
		RetailPrice: args.RetailPrice,
		ProductType: args.ProductType,
		MetaFields:  args.MetaFields,
		BrandID:     args.BrandID,
	}
	if product.Code != "" {
		number, ok := convert.ParseCodeNorm(product.Code)
		if ok {
			product.CodeNorm = number
		}
	}
	if product.Code == "" {
		var maxCodeNorm int
		productTemp, err := a.shopProduct(ctx).ShopID(args.ShopID).IncludeDeleted().GetProductByMaximumCodeNorm()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			maxCodeNorm = productTemp.CodeNorm
		case cm.NotFound:
			// no-op
		default:
			return nil, err
		}
		if maxCodeNorm >= convert.MaxCodeNorm {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
		}
		codeNorm := maxCodeNorm + 1
		product.Code = convert.GenerateCodeProduct(codeNorm)
		product.CodeNorm = codeNorm
	}
	if err := a.shopProduct(ctx).CreateShopProduct(product); err != nil {
		return nil, sqlstore.CheckProductExternalError(err, args.ExternalID, args.ExternalCode)
	}
	result, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(product.ProductID).GetShopProductWithVariants()
	return result, err
}

func (a *Aggregate) UpdateShopProductInfo(ctx context.Context, args *catalog.UpdateShopProductInfoArgs) (*catalog.ShopProductWithVariants, error) {
	if args.BrandID.Valid {
		_, err := a.shopBrand(ctx).ShopID(args.ShopID).ID(args.BrandID.ID).GetShopBrand()
		if err != nil {
			return nil, cm.MapError(err).
				Map(cm.NotFound, cm.InvalidArgument, "Mã thương hiệu không tồn tại").
				Throw()
		}
	}
	product, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProduct()
	if err != nil {
		return nil, err
	}
	product = convert.Apply_catalog_UpdateShopProductInfoArgs_catalog_ShopProduct(args, product)
	productModel := &model.ShopProduct{}
	if err := scheme.Convert(product, productModel); err != nil {
		return nil, err
	}
	if err = a.shopProduct(ctx).UpdateShopProduct(productModel); err != nil {
		return nil, err
	}
	result, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProductWithVariants()
	return result, err
}

func (a *Aggregate) UpdateShopProductStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (int, error) {
	count, err := a.shopProduct(ctx).ShopID(args.ShopID).IDs(args.IDs...).UpdateStatusShopProducts(args.Status)
	return count, err
}

func (a *Aggregate) UpdateShopProductImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopProductWithVariants, error) {
	productDB, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ID).GetShopProduct()
	if err != nil {
		return nil, err
	}
	var newImageURLs []string
	newImageURLs = productDB.ImageURLs
	newImageURLs, err = Patch(newImageURLs, args.Updates)
	if err != nil {
		return nil, err
	}
	productDB.ImageURLs = newImageURLs
	if err = a.shopProduct(ctx).ShopID(productDB.ShopID).ID(productDB.ProductID).UpdateImageShopProduct(productDB); err != nil {
		return nil, err
	}
	result, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ID).GetShopProductWithVariants()
	return result, nil
}

func (a *Aggregate) DeleteShopProducts(ctx context.Context, args *shopping.IDsQueryShopArgs) (int, error) {
	deletedProduct, err := a.shopProduct(ctx).ShopID(args.ShopID).IDs(args.IDs...).SoftDelete()
	if err != nil {
		return 0, err
	}
	_, err = a.shopVariant(ctx).ShopID(args.ShopID).ProductIDs(args.IDs...).SoftDelete()
	if err != nil {
		return 0, err
	}
	return deletedProduct, nil
}

func (a *Aggregate) CreateShopVariant(ctx context.Context, args *catalog.CreateShopVariantArgs) (*catalog.ShopVariant, error) {
	prodcut, err := a.shopProduct(ctx).
		ShopID(args.ShopID).
		ID(args.ProductID).
		GetShopProductDB()
	if err != nil {
		return nil, err
	}

	variant := &catalog.ShopVariant{
		ExternalID:   args.ExternalID,
		ExternalCode: args.ExternalCode,
		PartnerID:    args.PartnerID,
		ShopID:       args.ShopID,
		ProductID:    args.ProductID,
		VariantID:    cm.NewID(),
		Code:         args.Code,
		Name:         args.Name,
		ShortDesc:    args.ShortDesc,
		Description:  args.Description,
		DescHTML:     args.DescHTML,
		ImageURLs:    args.ImageURLs,
		Status:       0,
		Attributes:   args.Attributes,
		CostPrice:    args.CostPrice,
		ListPrice:    args.ListPrice,
		RetailPrice:  args.RetailPrice,
		Note:         args.Note,
	}
	if variant.Code != "" {
		ss := strings.Split(variant.Code, "-")
		if len(ss) == 2 {
			_, ok := convert.ParseCodeNorm(ss[0])
			if ok {
				log.Println(ss[1])
				codeNorm, err := strconv.Atoi(ss[1])
				if err == nil {
					variant.CodeNorm = codeNorm
				}
			}
		}
	}
	if variant.Code == "" {
		var maxCodeNorm int
		variantTemp, err := a.shopVariant(ctx).ShopID(args.ShopID).IncludeDeleted().GetVariantByMaximumCodeNorm(variant.ProductID)
		switch cm.ErrorCode(err) {
		case cm.NoError:
			maxCodeNorm = variantTemp.CodeNorm
		case cm.NotFound:
			// no-op
		default:
			return nil, err
		}
		if maxCodeNorm >= convert.MaxCodeNormVariant {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
		}
		codeNorm := maxCodeNorm + 1
		variant.Code = convert.GenerateCodeVariant(prodcut.Code, codeNorm)
		variant.CodeNorm = codeNorm
	}

	if err = a.shopVariant(ctx).CreateShopVariant(variant); err != nil {
		return nil, sqlstore.CheckShopVariantExternalError(err, args.ExternalID, args.ExternalCode)
	}
	return variant, nil
}

func (a *Aggregate) UpdateShopVariantInfo(ctx context.Context, args *catalog.UpdateShopVariantInfoArgs) (*catalog.ShopVariant, error) {
	variant, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariant()
	if err != nil {
		return nil, err
	}
	variant = convert.Apply_catalog_UpdateShopVariantInfoArgs_catalog_ShopVariant(args, variant)
	variantModel := &model.ShopVariant{}
	if err := scheme.Convert(variant, variantModel); err != nil {
		return nil, err
	}
	if err = a.shopVariant(ctx).UpdateShopVariant(variantModel); err != nil {
		return nil, err
	}
	result, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariant()
	return result, err
}

func (a *Aggregate) DeleteShopVariants(ctx context.Context, args *shopping.IDsQueryShopArgs) (int, error) {
	deleted, err := a.shopVariant(ctx).ShopID(args.ShopID).IDs(args.IDs...).SoftDelete()
	if err != nil {
		return 0, err
	}
	variants, err := a.shopVariantSupplier(ctx).ShopID(args.ShopID).VariantIDs(args.IDs...).ListVariantSupplier()
	if err != nil {
		return deleted, err
	}
	if len(variants) != 0 {
		variantIDs := make([]dot.ID, 0, len(variants))
		for _, variant := range variants {
			variantIDs = append(variantIDs, variant.VariantID)
		}
		if err := a.deleteVariantsSupplier(ctx, variantIDs, args.ShopID); err != nil {
			return deleted, err
		}
	}
	return deleted, err
}

func (a *Aggregate) UpdateShopVariantStatus(ctx context.Context, args *catalog.UpdateStatusArgs) (int, error) {
	update, err := a.shopVariant(ctx).ShopID(args.ShopID).IDs(args.IDs...).UpdateStatusShopVariant(args.Status)
	return update, err
}

func (a *Aggregate) UpdateShopVariantImages(ctx context.Context, args *catalog.UpdateImagesArgs) (*catalog.ShopVariant, error) {
	variantDB, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.ID).GetShopVariant()
	if err != nil {
		return nil, err
	}

	var newImageURLs []string
	newImageURLs = variantDB.ImageURLs
	newImageURLs, err = Patch(newImageURLs, args.Updates)
	if err != nil {
		return nil, err
	}
	variantDB.ImageURLs = newImageURLs
	if err = a.shopVariant(ctx).ShopID(variantDB.ShopID).ID(args.ID).UpdateImageShopVariant(variantDB); err != nil {
		return nil, err
	}
	return variantDB, nil
}

func (a *Aggregate) UpdateShopVariantAttributes(ctx context.Context, args *catalog.UpdateShopVariantAttributes) (*catalog.ShopVariant, error) {
	variantDB, err := a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	if len(args.Attributes) <= 0 {
		return nil, cm.Error(cm.InvalidArgument, "Atributes is empty", nil)
	}
	var attributesUpdate model.ProductAttributes
	for _, value := range args.Attributes {
		attributesUpdate = append(attributesUpdate, &model.ProductAttribute{
			Name:  value.Value,
			Value: value.Name,
		})
	}
	variantDB.Attributes = attributesUpdate
	if err = a.shopVariant(ctx).UpdateShopVariant(variantDB); err != nil {
		return nil, err
	}
	return a.shopVariant(ctx).ShopID(args.ShopID).ID(args.VariantID).GetShopVariant()
}

func Patch(source []string, update []*meta.UpdateSet) ([]string, error) {
	if OptionValue(update, meta.OpDeleteAll) != nil {
		return []string{}, nil
	}
	if replace := OptionValue(update, meta.OpReplaceAll); replace != nil {
		arr, _, err := replace.Update(source)
		if err != nil {
			return []string{}, err
		}
		return arr, nil
	}
	var arrResult []string
	arrResult = source
	if add := OptionValue(update, meta.OpAdd); add != nil {
		var err error
		arrResult, _, err = add.Update(arrResult)
		if err != nil {
			return []string{}, err
		}
	}
	if remove := OptionValue(update, meta.OpRemove); remove != nil {
		var err error
		arrResult, _, err = remove.Update(arrResult)
		if err != nil {
			return []string{}, err
		}
	}
	return arrResult, nil
}

func OptionValue(update []*meta.UpdateSet, op meta.UpdateOp) *meta.UpdateSet {
	for i := 0; i < len(update); i++ {
		if update[i].Op == op {
			return update[i]
		}
	}
	return nil
}

func (a *Aggregate) UpdateShopProductCategory(ctx context.Context, args *catalog.UpdateShopProductCategoryArgs) (*catalog.ShopProductWithVariants, error) {
	product, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProduct()
	if err != nil {
		return nil, err
	}
	_, err = a.shopCategory(ctx).ShopID(args.ShopID).ID(args.CategoryID).GetShopCategoryDB()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Mã danh mục không không tồn tại")
	}
	product = convert.Apply_catalog_UpdateShopProductCategoryArgs_catalog_ShopProduct(args, product)
	productModel := &model.ShopProduct{}
	if err := scheme.Convert(product, productModel); err != nil {
		return nil, err
	}
	if err = a.shopProduct(ctx).UpdateShopProductCategory(productModel); err != nil {
		return nil, err
	}
	result, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProductWithVariants()
	return result, err
}

func (a *Aggregate) CreateShopCategory(ctx context.Context, args *catalog.CreateShopCategoryArgs) (*catalog.ShopCategory, error) {
	category := &catalog.ShopCategory{
		ID:       cm.NewID(),
		ShopID:   args.ShopID,
		Name:     args.Name,
		ParentID: args.ParentID,
		Status:   args.Status,
	}
	if args.ParentID != 0 {
		if _, err := a.shopCategory(ctx).ID(args.ParentID).GetShopCategory(); err != nil {
			return nil, err
		}
	}
	if err := a.shopCategory(ctx).CreateShopCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (a *Aggregate) UpdateShopCategory(ctx context.Context, args *catalog.UpdateShopCategoryArgs) (*catalog.ShopCategory, error) {
	category, err := a.shopCategory(ctx).ShopID(args.ShopID).ID(args.ID).GetShopCategory()
	if err != nil {
		return nil, err
	}
	if args.ParentID != 0 {
		if _, err = a.shopCategory(ctx).ID(args.ParentID).GetShopCategory(); err != nil {
			return nil, err
		}
	}
	category = convert.Apply_catalog_UpdateShopCategoryArgs_catalog_ShopCategory(args, category)
	categoryModel := &model.ShopCategory{}
	if err := scheme.Convert(category, categoryModel); err != nil {
		return nil, err
	}
	if err = a.shopCategory(ctx).UpdateShopCategory(categoryModel); err != nil {
		return nil, err
	}
	result, err := a.shopCategory(ctx).ShopID(args.ShopID).ID(args.ID).GetShopCategory()
	return result, err
}

func (a *Aggregate) DeleteShopCategory(ctx context.Context, args *catalog.DeleteShopCategoryArgs) (deleted int, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var err error
		deleted, err = a.shopProduct(ctx).ShopID(args.ShopID).RemoveShopProductCategory()
		if err != nil {
			return err
		}
		deleted, err = a.shopCategory(ctx).ID(args.ID).ShopID(args.ShopID).SoftDelete()
		return err
	})

	return deleted, err
}

func (a *Aggregate) RemoveShopProductCategory(ctx context.Context, args *catalog.RemoveShopProductCategoryArgs) (*catalog.ShopProductWithVariants, error) {
	_, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).RemoveShopProductCategory()
	if err != nil {
		return nil, err
	}
	product, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProductWithVariants()
	if err != nil {
		return nil, err
	}
	return product, err
}

func (a *Aggregate) CreateShopCollection(ctx context.Context, args *catalog.CreateShopCollectionArgs) (*catalog.ShopCollection, error) {
	collection := &catalog.ShopCollection{
		ID:          cm.NewID(),
		ShopID:      args.ShopID,
		Name:        args.Name,
		DescHTML:    args.DescHTML,
		Description: args.Description,
		ShortDesc:   args.ShortDesc,
	}
	if err := a.shopCollection(ctx).CreateShopCollection(collection); err != nil {
		return nil, err
	}
	return collection, nil
}

func (a *Aggregate) UpdateShopCollection(ctx context.Context, args *catalog.UpdateShopCollectionArgs) (*catalog.ShopCollection, error) {
	collection, err := a.shopCollection(ctx).ShopID(args.ShopID).ID(args.ID).GetShopCollection()
	if err != nil {
		return nil, err
	}
	collection = convert.Apply_catalog_UpdateShopCollectionArgs_catalog_ShopCollection(args, collection)
	collectionModel := &model.ShopCollection{}
	if err := scheme.Convert(collection, collectionModel); err != nil {
		return nil, err
	}

	if err = a.shopCollection(ctx).UpdateShopCollection(collectionModel); err != nil {
		return nil, err
	}
	result, err := a.shopCollection(ctx).ShopID(args.ShopID).ID(args.ID).GetShopCollection()
	return result, err
}

func (a *Aggregate) AddShopProductCollection(ctx context.Context, args *catalog.AddShopProductCollectionArgs) (created int, _ error) {
	var err error
	if len(args.CollectionIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Mã bộ sưu tập không được để trống")
	}
	if args.ProductID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Mã sản phẩm không được để trống")
	}
	_, err = a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ProductID).GetShopProduct()
	if err != nil {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Mã sản phẩm không không tồn tại")
	}
	for _, collectionID := range args.CollectionIDs {
		if collectionID == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, err, "Mã bộ sưu tập không được để trống")
		}
	}
	collections, err := a.shopCollection(ctx).ShopID(args.ShopID).IDs(args.CollectionIDs).ListShopCollections()
	if err != nil {
		return 0, err
	}
	if len(collections) != len(args.CollectionIDs) {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Mã bộ sưu tập không tồn tại")
	}

	err = a.db.InTransaction(ctx, func(q cmsql.QueryInterface) error {
		for _, collectionID := range args.CollectionIDs {
			productCollection := &catalog.ShopProductCollection{
				ProductID:    args.ProductID,
				ShopID:       args.ShopID,
				CollectionID: collectionID,
			}
			lineCreated, err := a.shopProductCollection(ctx).AddProductToCollection(productCollection)
			if err != nil {
				return err
			}
			created += lineCreated
		}
		return nil
	})
	return created, err
}

func (a *Aggregate) RemoveShopProductCollection(ctx context.Context, args *catalog.RemoveShopProductColelctionArgs) (deleted int, _ error) {
	var err error
	var removedProduct int
	if len(args.CollectionIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Mã bộ sưu tập không được để trống")
	}
	if args.ProductID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Mã sản phẩm không được để trống")
	}
	for i := 0; i < len(args.CollectionIDs); i++ {
		if args.CollectionIDs[i] == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, err, "Mã bộ sưu tập không được để trống")
		}
	}
	removedProduct, err = a.shopProductCollection(ctx).ShopID(args.ShopID).ProductID(args.ProductID).IDs(args.CollectionIDs).RemoveProductFromCollection()
	return removedProduct, err
}

func (a *Aggregate) UpdateShopProductMetaFields(ctx context.Context, args *catalog.UpdateShopProductMetaFieldsArgs) (*catalog.ShopProductWithVariants, error) {
	productDB, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ID).GetShopProduct()
	if err != nil {
		return nil, err
	}
	productDB.MetaFields = args.MetaFields
	if err := a.shopProduct(ctx).ShopID(args.ShopID).ID(productDB.ProductID).UpdateMetaFieldsShopProduct(productDB); err != nil {
		return nil, err
	}
	result, err := a.shopProduct(ctx).ShopID(args.ShopID).ID(args.ID).GetShopProductWithVariants()
	return result, nil
}

func (a *Aggregate) CreateBrand(ctx context.Context, brand *catalog.CreateBrandArgs) (*catalog.ShopBrand, error) {
	brandCreate := convert.Apply_catalog_CreateBrandArgs_catalog_ShopBrand(brand, nil)
	err := a.shopBrand(ctx).CreateShopBrand(brandCreate)
	if err != nil {
		return nil, err
	}
	resultDB, err := a.shopBrand(ctx).ShopID(brandCreate.ShopID).ID(brandCreate.ID).GetShopBrand()
	return resultDB, err
}

func (a *Aggregate) UpdateBrandInfo(ctx context.Context, brand *catalog.UpdateBrandArgs) (*catalog.ShopBrand, error) {
	if brand.ShopID == 0 && brand.ID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Invalid Argument ShopId or ID ", nil)
	}
	brandDb, err := a.shopBrand(ctx).ShopID(brand.ShopID).ID(brand.ID).GetShopBrand()
	if err != nil {
		return nil, err
	}
	brandDb = convert.Apply_catalog_UpdateBrandArgs_catalog_ShopBrand(brand, brandDb)
	brandDb.UpdatedAt = time.Now()
	err = a.shopBrand(ctx).ShopID(brand.ShopID).ID(brand.ID).UpdateShopBrand(brandDb)
	if err != nil {
		return nil, err
	}
	resultDB, err := a.shopBrand(ctx).ShopID(brandDb.ShopID).ID(brandDb.ID).GetShopBrand()
	return resultDB, err
}

func (a *Aggregate) DeleteShopBrand(ctx context.Context, ids []dot.ID, shopID dot.ID) (int, error) {
	var count = 0
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		countRecord, errTrans := a.shopBrand(ctx).ShopID(shopID).IDs(ids...).SoftDelete()
		if errTrans != nil {
			return errTrans
		}
		count = countRecord
		products, errTrans := a.shopProduct(ctx).BrandIDs(ids...).ListShopProductsDB()
		if errTrans != nil {
			return errTrans
		}
		var productIDs []dot.ID
		for _, value := range products {
			productIDs = append(productIDs, value.ProductID)
		}
		errTrans = a.shopProduct(ctx).ShopID(shopID).IDs(productIDs...).RemoveBrands()
		if errTrans != nil {
			return errTrans
		}
		return nil
	})
	return count, err
}

func (a *Aggregate) CreateVariantSupplier(ctx context.Context, sv *catalog.CreateVariantSupplier) (*catalog.ShopVariantSupplier, error) {
	if sv.ShopID == 0 || sv.VariantID == 0 || sv.SupplierID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id, variant_id or supplier_id in request", nil)
	}
	variantSupplier := convert.Apply_catalog_CreateVariantSupplier_catalog_ShopVariantSupplier(sv, nil)
	err := a.shopVariantSupplier(ctx).CreateVariantSupplier(variantSupplier)
	if err != nil {
		return nil, err
	}
	return variantSupplier, nil
}
func (a *Aggregate) CreateVariantsSupplier(ctx context.Context, vs *catalog.CreateVariantsSupplier) (int, error) {
	if vs.VariantIDs == nil {
		return 0, cm.Error(cm.InvalidArgument, "Miss variant_ids ", nil)
	}
	var lineCreate int
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, variant := range vs.VariantIDs {
			create := &catalog.CreateVariantSupplier{
				ShopID:     vs.ShopID,
				SupplierID: vs.SupplierID,
				VariantID:  variant,
			}
			_, errTrans := a.CreateVariantSupplier(ctx, create)
			if errTrans != nil {
				return errTrans
			}
			lineCreate = lineCreate + 1
		}
		return nil
	})
	return lineCreate, err
}

func (a *Aggregate) DeleteVariantSupplier(ctx context.Context, variantID dot.ID, supplierID dot.ID, shopID dot.ID) error {
	if shopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing shop_id in request", nil)
	}
	if supplierID == 0 && variantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing variant_id or supplier_id in request", nil)
	}
	var query = a.shopVariantSupplier(ctx).ShopID(shopID)
	if supplierID != 0 {
		query = query.SupplierID(supplierID)
	}
	if variantID != 0 {
		query = query.VariantID(variantID)
	}
	err := query.DeleteVariantSupplier()
	return err
}

func (a *Aggregate) deleteVariantsSupplier(ctx context.Context, variantIDs []dot.ID, shopID dot.ID) error {
	if shopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing shop_id in request", nil)
	}
	if variantIDs == nil {
		return cm.Error(cm.InvalidArgument, "Missing varianIDs in request", nil)
	}
	return a.shopVariantSupplier(ctx).ShopID(shopID).VariantIDs(variantIDs...).DeleteVariantSupplier()
}
