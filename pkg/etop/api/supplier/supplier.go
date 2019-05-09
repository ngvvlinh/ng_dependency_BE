package supplier

import (
	"context"

	"github.com/asaskevich/govalidator"

	cmP "etop.vn/backend/pb/common"
	supplierP "etop.vn/backend/pb/etop/supplier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	supplierW "etop.vn/backend/wrapper/etop/supplier"
)

var ll = l.New()

func init() {
	bus.AddHandler("api", GetCategories)
	bus.AddHandler("api", GetCategoriesByIDs)
	bus.AddHandler("api", GetCategory)
	bus.AddHandler("api", GetPriceRules)
	bus.AddHandler("api", GetVariant)
	bus.AddHandler("api", GetVariants)
	bus.AddHandler("api", GetVariantsByIDs)
	bus.AddHandler("api", UpdatePriceRules)
	bus.AddHandler("api", UpdateVariant)
	bus.AddHandler("api", UpdateVariantImages)
	bus.AddHandler("api", UpdateVariants)
	bus.AddHandler("api", UpdateVariantsStatus)
	bus.AddHandler("api", VersionInfo)
	bus.AddHandler("api", GetBrand)
	bus.AddHandler("api", GetBrands)
	bus.AddHandler("api", GetBrandsByIDs)
	bus.AddHandler("api", CreateBrand)
	bus.AddHandler("api", UpdateBrand)
	bus.AddHandler("api", DeleteBrand)
	bus.AddHandler("api", UpdateBrandImages)
	bus.AddHandler("api", UpdateSupplier)
	bus.AddHandler("api", GetProductEndpoint)
	bus.AddHandler("api", GetProductsEndpoint)
	bus.AddHandler("api", GetProductsByIDsEndpoint)
	bus.AddHandler("api", UpdateProductEndpoint)
	bus.AddHandler("api", UpdateProductImagesEndpoint)
	bus.AddHandler("api", SetDefaultAddress)
}

func VersionInfo(ctx context.Context, q *supplierW.VersionInfoEndpoint) error {
	q.Result = &cmP.VersionInfoResponse{
		Service: "etop.Supplier",
		Version: "0.1",
	}
	return nil
}

func GetCategory(ctx context.Context, q *supplierW.GetCategoryEndpoint) error {
	query := &model.GetProductSourceCategoryQuery{
		SupplierID: q.Context.AccountID,
		CategoryID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = supplierP.PbCategory(query.Result)
	return nil
}

func GetCategories(ctx context.Context, q *supplierW.GetCategoriesEndpoint) error {
	query := &model.GetProductSourceCategoriesExtendedQuery{SupplierID: q.Context.AccountID}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.CategoriesResponse{
		Categories: supplierP.PbCategories(query.Result.Categories),
	}
	return nil
}

func GetCategoriesByIDs(ctx context.Context, q *supplierW.GetCategoriesByIDsEndpoint) error {
	query := &model.GetProductSourceCategoriesExtendedQuery{
		SupplierID: q.Context.AccountID,
		IDs:        q.Ids,
	}
	if q.Ids == nil {
		query.IDs = make([]int64, 0)
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.CategoriesResponse{
		Categories: supplierP.PbCategories(query.Result.Categories),
	}
	return nil
}

func GetVariant(ctx context.Context, q *supplierW.GetVariantEndpoint) error {
	query := &model.GetVariantQuery{
		SupplierID: q.Context.AccountID,
		VariantID:  q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = supplierP.PbVariant(query.Result)
	return nil
}

func GetVariants(ctx context.Context, q *supplierW.GetVariantsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetVariantsExtendedQuery{
		SupplierID: q.Context.AccountID,
		Paging:     paging,
		Filters:    cmP.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.VariantsResponse{
		Variants: supplierP.PbVariants(query.Result.Variants),
		Paging:   cmP.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}

func GetVariantsByIDs(ctx context.Context, q *supplierW.GetVariantsByIDsEndpoint) error {
	query := &model.GetVariantsExtendedQuery{
		SupplierID: q.Context.AccountID,
		IDs:        q.Ids,
	}
	if q.Ids == nil {
		query.IDs = make([]int64, 0)
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.VariantsResponse{
		Variants: supplierP.PbVariants(query.Result.Variants),
	}
	return nil
}

func UpdateVariant(ctx context.Context, q *supplierW.UpdateVariantEndpoint) error {
	cmd := &model.UpdateVariantCommand{
		SupplierID: q.Context.AccountID,
		Variant:    supplierP.PbUpdateVariantToModel(q.UpdateVariantRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = supplierP.PbVariant(cmd.Result)
	return nil
}

func UpdateVariantImages(ctx context.Context, q *supplierW.UpdateVariantImagesEndpoint) error {
	query := &model.GetVariantQuery{
		SupplierID: q.Context.AccountID,
		VariantID:  q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	imageURLs, err := patchVariantImages(query.Result, q.UpdateVariantImagesRequest)
	if err != nil {
		return err
	}

	cmd := &model.UpdateVariantImagesCommand{
		SupplierID: q.Context.AccountID,
		VariantID:  q.Id,
		ImageURLs:  imageURLs,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = supplierP.PbVariant(cmd.Result)
	return nil
}

func patchVariantImages(variant *model.VariantExtended, req *supplierP.UpdateVariantImagesRequest) ([]string, error) {
	r := &model.UpdateListRequest{
		Adds:       req.Adds,
		Deletes:    req.Deletes,
		ReplaceAll: req.ReplaceAll,
		DeleteAll:  req.DeleteAll,
	}
	if err := r.Verify(); err != nil {
		return nil, err
	}

	for _, imgURL := range req.Adds {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "Invalid url: "+imgURL, nil)
		}
	}
	for _, imgURL := range req.ReplaceAll {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "Invalid url: "+imgURL, nil)
		}
	}

	return r.Patch(variant.ImageURLs), nil
}

func UpdateVariants(ctx context.Context, q *supplierW.UpdateVariantsEndpoint) error {
	cmd := &model.UpdateVariantsCommand{
		SupplierID: q.Context.AccountID,
		Variants:   supplierP.PbUpdateVariantsToModel(q.Updates),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &supplierP.UpdateVariantsResponse{
		Variants: supplierP.PbVariants(cmd.Result.Variants),
		Errors:   cmP.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func UpdateVariantsStatus(ctx context.Context, q *supplierW.UpdateVariantsStatusEndpoint) error {
	if q.SStatus == nil {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	cmd := &model.UpdateVariantsStatusCommand{
		SupplierID: q.Context.AccountID,
		IDs:        q.Ids,
	}
	cmd.Update.SupplierStatus = q.SStatus.ToModel()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func GetPriceRules(ctx context.Context, q *supplierW.GetPriceRulesEndpoint) error {
	priceRules, ok := q.Context.Supplier.GetPriceRules()
	q.Result = supplierP.PbPriceRules(priceRules, ok)
	return nil
}

func UpdatePriceRules(ctx context.Context, q *supplierW.UpdatePriceRulesEndpoint) error {
	priceRules, _ := q.Context.Supplier.GetPriceRules()
	rules := priceRules.Rules

	r := q.UpdatePriceRulesRequest
	if r.General != nil {
		priceRules.General = supplierP.PbPriceRuleToModel(r.General)
		if !priceRules.General.IsZeroIdentifier() {
			return cm.Error(cm.InvalidArgument, "General rule must has no identifier", nil)
		}
		if err := priceRules.General.Validate(); err != nil {
			return err
		}
	}

	updateErrors := make([]error, len(q.Updates))
	for i, update := range q.Updates {
		if update == nil {
			updateErrors[i] = cm.Error(cm.InvalidArgument, "", nil)
			continue
		}

		updateRule := supplierP.PbPriceRuleToModel(update)
		if err := updateRule.Validate(); err != nil {
			updateErrors[i] = err
			continue
		}

		if index, err := findRule(rules, updateRule); err != nil {
			updateErrors[i] = err
		} else if index >= 0 {
			rules[index] = updateRule
		} else {
			rules = append(rules, updateRule)
		}
	}

	deleteCount := 0
	deleteErrors := make([]error, len(q.Deletes))
	for i, delete := range q.Deletes {
		if delete == nil {
			deleteErrors[i] = cm.Error(cm.InvalidArgument, "", nil)
		}

		deleteRule := supplierP.PbPriceRuleToModel(delete)
		if index, err := findRule(rules, deleteRule); err != nil {
			deleteErrors[i] = err
		} else if index >= 0 {
			rules[index] = nil
			deleteCount++
		} else {
			deleteErrors[i] = cm.Error(cm.NotFound, "", nil)
		}
	}

	if deleteCount > 0 {
		newRules := make([]*model.SupplierPriceRule, 0, len(rules)-deleteCount)
		for _, rule := range rules {
			if rule != nil {
				newRules = append(newRules, rule)
			}
		}
		rules = newRules
	}

	// validate the last time
	priceRules.Rules = rules
	if err := priceRules.Validate(); err != nil {
		return err
	}

	cmd := &model.UpdatePriceRulesCommand{
		SupplierID: q.Context.AccountID,
		PriceRules: priceRules,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	pbPriceRules := supplierP.PbPriceRules(priceRules, true)
	q.Result = &supplierP.UpdatePriceRulesResponse{
		General:      pbPriceRules.General,
		Rules:        pbPriceRules.Rules,
		UpdateErrors: cmP.PbErrors(updateErrors),
		DeleteErrors: cmP.PbErrors(deleteErrors),
	}
	return nil
}

func findRule(rules []*model.SupplierPriceRule, rule *model.SupplierPriceRule) (int, error) {
	if rule.IsZeroIdentifier() {
		return 0, cm.Error(cm.InvalidArgument, "", nil)
	}
	for i, r := range rules {
		if r.ExternalCategoryID == rule.ExternalCategoryID &&
			r.SupplierCategoryID == rule.SupplierCategoryID &&
			r.Tag == rule.Tag {
			return i, nil
		}
	}
	return -1, nil
}

func GetBrand(ctx context.Context, q *supplierW.GetBrandEndpoint) error {
	query := &model.GetProductBrandQuery{
		SupplierID: q.Context.AccountID,
		ID:         q.Id,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = supplierP.PbBrandExt(query.Result)
	return nil
}

func GetBrands(ctx context.Context, q *supplierW.GetBrandsEndpoint) error {
	if q.Context.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}

	query := &model.GetProductBrandsQuery{
		SupplierID: q.Context.AccountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &supplierP.BrandsResponse{
		Brands: supplierP.PbBrandsExt(query.Result.Brands),
	}
	return nil
}

func GetBrandsByIDs(ctx context.Context, q *supplierW.GetBrandsByIDsEndpoint) error {
	query := &model.GetProductBrandsQuery{
		SupplierID: q.Context.AccountID,
		Ids:        q.Ids,
	}

	if q.Ids == nil {
		query.Ids = make([]int64, 0)
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &supplierP.BrandsResponse{
		Brands: supplierP.PbBrandsExt(query.Result.Brands),
	}

	return nil
}

func CreateBrand(ctx context.Context, q *supplierW.CreateBrandEndpoint) error {
	if len(q.ImageUrls) > 0 {
		for _, imgURL := range q.ImageUrls {
			if !govalidator.IsURL(imgURL) {
				return cm.Error(cm.InvalidArgument, "Invalid url: "+imgURL, nil)
			}
		}
	}

	cmd := &model.CreateProductBrandCommand{
		Brand: &model.ProductBrand{
			SupplierID:  q.Context.AccountID,
			Name:        q.Name,
			Description: q.Description,
			Policy:      q.Policy,
			ImageURLs:   q.ImageUrls,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	query := &model.GetProductBrandQuery{
		SupplierID: q.Context.Supplier.ID,
		ID:         cmd.Result.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = supplierP.PbBrandExt(query.Result)
	return nil
}

func UpdateBrand(ctx context.Context, q *supplierW.UpdateBrandEndpoint) error {
	cmd := &model.UpdateProductBrandCommand{
		Brand: &model.ProductBrand{
			ID:          q.Id,
			SupplierID:  q.Context.Supplier.ID,
			Name:        q.Name,
			Description: q.Description,
			Policy:      q.Policy,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil
	}

	query := &model.GetProductBrandQuery{
		SupplierID: q.Context.Supplier.ID,
		ID:         q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = supplierP.PbBrandExt(query.Result)
	return nil
}

func DeleteBrand(ctx context.Context, q *supplierW.DeleteBrandEndpoint) error {
	cmd := &model.DeleteProductBrandCommand{
		SupplierID: q.Context.AccountID,
		ID:         q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &cmP.Empty{}
	return nil
}

func UpdateBrandImages(ctx context.Context, q *supplierW.UpdateBrandImagesEndpoint) error {
	query := &model.GetProductBrandQuery{
		SupplierID: q.Context.AccountID,
		ID:         q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	imageURLs, err := patchBrandImages(query.Result.ProductBrand, q.UpdateBrandImagesRequest)
	if err != nil {
		return err
	}

	cmd := &model.UpdateProductBrandCommand{
		Brand: &model.ProductBrand{
			ID:         q.Id,
			SupplierID: q.Context.AccountID,
			ImageURLs:  imageURLs,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil
	}

	q.Result = supplierP.PbBrand(cmd.Result)
	return nil
}

func patchBrandImages(brand *model.ProductBrand, req *supplierP.UpdateBrandImagesRequest) ([]string, error) {
	r := &model.UpdateListRequest{
		Adds:       req.Adds,
		Deletes:    req.Deletes,
		ReplaceAll: req.ReplaceAll,
		DeleteAll:  req.DeleteAll,
	}
	if err := r.Verify(); err != nil {
		return nil, err
	}

	for _, imgURL := range req.Adds {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "Invalid url: "+imgURL, nil)
		}
	}
	for _, imgURL := range req.ReplaceAll {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "Invalid url: "+imgURL, nil)
		}
	}

	return r.Patch(brand.ImageURLs), nil
}

func UpdateSupplier(ctx context.Context, q *supplierW.UpdateSupplierEndpoint) error {
	var warehouseAddressID int64
	if q.WarehouseAddress != nil {
		warehouseAddress := q.WarehouseAddress
		addressObj := &model.Address{
			Province:     warehouseAddress.Province,
			ProvinceCode: warehouseAddress.ProvinceCode,
			District:     warehouseAddress.District,
			DistrictCode: warehouseAddress.DistrictCode,
			Ward:         warehouseAddress.Ward,
			WardCode:     warehouseAddress.WardCode,
			Zip:          warehouseAddress.Zip,
			Address1:     warehouseAddress.Address1,
			Address2:     warehouseAddress.Address2,
			FullName:     warehouseAddress.FullName,
			FirstName:    warehouseAddress.FirstName,
			LastName:     warehouseAddress.LastName,
			Email:        warehouseAddress.Email,
			Position:     warehouseAddress.Position,
			Phone:        warehouseAddress.Phone,
		}
		if warehouseAddress.Id != 0 {
			// update warehouse address
			addressObj.ID = warehouseAddress.Id
			addressCmd := &model.UpdateAddressCommand{
				Address: addressObj,
			}
			if err := bus.Dispatch(ctx, addressCmd); err != nil {
				return err
			}
			warehouseAddressID = addressCmd.Result.ID
		} else {
			// create new warehouse address
			addressObj.AccountID = q.Context.AccountID
			addressObj.Type = model.AddressTypeWarehouse
			addressCmd := &model.CreateAddressCommand{
				Address: addressObj,
			}
			if err := bus.Dispatch(ctx, addressCmd); err != nil {
				return err
			}
			warehouseAddressID = addressCmd.Result.ID
		}
	}
	cmd := &model.UpdateSupplierCommand{
		Supplier: &model.Supplier{
			ID:                 q.Context.AccountID,
			Name:               q.Name,
			CompanyInfo:        q.CompanyInfo.ToModel(),
			WarehouseAddressID: warehouseAddressID,
			BankAccount:        q.BankAccount.ToModel(),
			ContactPersons:     supplierP.PbContactPersonsToModel(q.ContactPersons),
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = supplierP.PbSupplierExt(cmd.Result)
	return nil
}

func GetProductEndpoint(ctx context.Context, q *supplierW.GetProductEndpoint) error {
	query := &model.GetProductQuery{
		SupplierID: q.Context.AccountID,
		ProductID:  q.Id,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = supplierP.PbProduct(query.Result)
	return nil
}

func GetProductsEndpoint(ctx context.Context, q *supplierW.GetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetProductsExtendedQuery{
		SupplierID:        q.Context.AccountID,
		Paging:            paging,
		Filters:           cmP.ToFilters(q.Filters),
		ProductSourceType: model.ProductSourceKiotViet,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.ProductsResponse{
		Products: supplierP.PbProducts(query.Result.Products),
		Paging:   cmP.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}

func GetProductsByIDsEndpoint(ctx context.Context, q *supplierW.GetProductsByIDsEndpoint) error {
	query := &model.GetProductsExtendedQuery{
		SupplierID:        q.Context.AccountID,
		IDs:               q.Ids,
		ProductSourceType: model.ProductSourceKiotViet,
	}
	if q.Ids == nil {
		query.IDs = make([]int64, 0)
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.ProductsResponse{
		Products: supplierP.PbProducts(query.Result.Products),
	}
	return nil
}

func UpdateProductEndpoint(ctx context.Context, q *supplierW.UpdateProductEndpoint) error {
	cmd := &model.UpdateProductCommand{
		SupplierID: q.Context.AccountID,
		Product:    supplierP.PbUpdateProductToModel(q.UpdateProductRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = supplierP.PbProduct(cmd.Result)
	return nil
}

func UpdateProductImagesEndpoint(ctx context.Context, q *supplierW.UpdateProductImagesEndpoint) error {
	query := &model.GetProductQuery{
		SupplierID: q.Context.AccountID,
		ProductID:  q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	updateRequest := &model.UpdateListRequest{
		Adds:       q.Adds,
		Deletes:    q.Deletes,
		ReplaceAll: q.ReplaceAll,
		DeleteAll:  q.DeleteAll,
	}
	if err := updateRequest.Verify(); err != nil {
		return err
	}

	imageURLs, err := cmP.PatchImage(query.Result.Product.ImageURLs, updateRequest)
	if err != nil {
		return err
	}

	cmd := &model.UpdateProductImagesCommand{
		SupplierID: q.Context.AccountID,
		ProductID:  q.Id,
		ImageURLs:  imageURLs,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = supplierP.PbProduct(cmd.Result)
	return nil
}

func SetDefaultAddress(ctx context.Context, q *supplierW.SetDefaultAddressEndpoint) error {
	cmd := &model.SetDefaultAddressSupplierCommand{
		SupplierID: q.Context.Supplier.ID,
		Type:       q.Type.ToModel(),
		AddressID:  q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}

	return nil
}
