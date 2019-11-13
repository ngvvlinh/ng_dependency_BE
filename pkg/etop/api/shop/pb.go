package shop

import (
	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	common "etop.vn/backend/pb/common"
	pbcm "etop.vn/backend/pb/common"
	pbproducttype "etop.vn/backend/pb/etop/etc/product_type"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
)

func PbProductsQuantity(shopProducts []*catalog.ShopProductWithVariants, inventoryVariants map[int64]*inventory.InventoryVariant) (res []*pbshop.ShopProduct) {
	for _, product := range shopProducts {
		productPb := PbProductQuantity(product, inventoryVariants)
		res = append(res, productPb)
	}
	return
}

func PbProductQuantity(shopProduct *catalog.ShopProductWithVariants, inventoryVariants map[int64]*inventory.InventoryVariant) *pbshop.ShopProduct {
	shopProductPb := PbShopProductWithVariants(shopProduct)
	shopProductPb.Variants = PbVariantsQuantity(shopProduct.Variants, inventoryVariants)
	return shopProductPb
}

func PbVariantsQuantity(shopVariants []*catalog.ShopVariant, inventoryVariants map[int64]*inventory.InventoryVariant) []*pbshop.ShopVariant {
	var variants []*pbshop.ShopVariant
	for _, variant := range shopVariants {
		inventoryVariant := inventoryVariants[variant.VariantID]
		valuePb := PbVariantQuantity(variant, inventoryVariant)
		variants = append(variants, valuePb)
	}
	return variants
}

func PbVariantQuantity(shopVariant *catalog.ShopVariant, inventoryVariant *inventory.InventoryVariant) *pbshop.ShopVariant {
	shopVariantDB := PbShopVariant(shopVariant)
	if inventoryVariant != nil {
		shopVariantDB.InventoryVariant = &pbshop.InventoryVariantQuantity{
			QuantityOnHand: inventoryVariant.QuantityOnHand,
			QuantityPicked: inventoryVariant.QuantityPicked,
			Quantity:       inventoryVariant.QuantitySummary,
		}
	}
	return shopVariantDB
}

func PbInventory(args *inventory.InventoryVariant) *pbshop.InventoryVariant {
	return &pbshop.InventoryVariant{
		ShopId:         args.ShopID,
		VariantId:      args.VariantID,
		QuantityOnHand: args.QuantityOnHand,
		QuantityPicked: args.QuantityPicked,
		Quantity:       args.QuantitySummary,
		PurchasePrice:  args.PurchasePrice,
	}
}

func PbInventoryVariants(args []*inventory.InventoryVariant) []*pbshop.InventoryVariant {
	var inventoryVariants []*pbshop.InventoryVariant
	for _, value := range args {
		inventoryVariants = append(inventoryVariants, PbInventory(value))
	}
	return inventoryVariants
}

func PbBrand(args *catalog.ShopBrand) *pbshop.Brand {
	return &pbshop.Brand{
		ShopId:      args.ShopID,
		Id:          args.ID,
		Name:        args.BrandName,
		Description: args.Description,
		CreatedAt:   pbcm.PbTime(args.CreatedAt),
		UpdatedAt:   pbcm.PbTime(args.UpdatedAt),
	}
}

func PbBrands(args []*catalog.ShopBrand) []*pbshop.Brand {
	var brands []*pbshop.Brand
	for _, value := range args {
		brands = append(brands, PbBrand(value))
	}
	return brands
}

func PbShopInventoryVoucher(args *inventory.InventoryVoucher) *pbshop.InventoryVoucher {
	if args == nil {
		return nil
	}

	var inventoryVoucherItem []*pbshop.InventoryVoucherLine
	for _, value := range args.Lines {
		inventoryVoucherItem = append(inventoryVoucherItem, &pbshop.InventoryVoucherLine{
			VariantId: value.VariantID,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}
	return &pbshop.InventoryVoucher{
		Title:        args.Title,
		TotalAmount:  args.TotalAmount,
		CreatedBy:    args.CreatedBy,
		UpdatedBy:    args.UpdatedBy,
		Lines:        inventoryVoucherItem,
		RefId:        args.RefID,
		Code:         args.Code,
		RefType:      string(args.RefType),
		RefName:      string(args.RefName),
		TraderId:     args.TraderID,
		Trader:       PbShopTrader(args.Trader),
		Note:         args.Note,
		Type:         string(args.Type),
		Id:           args.ID,
		ShopId:       args.ShopID,
		CreatedAt:    pbcm.PbTime(args.CreatedAt),
		UpdatedAt:    pbcm.PbTime(args.UpdatedAt),
		CancelledAt:  pbcm.PbTime(args.CancelledAt),
		ConfirmedAt:  pbcm.PbTime(args.ConfirmedAt),
		CancelReason: args.CancelReason,
	}
}

func PbShopTrader(args *inventory.Trader) *pbshop.Trader {
	if args == nil {
		return nil
	}
	return &pbshop.Trader{
		Id:       args.ID,
		Type:     args.Type,
		FullName: args.FullName,
		Phone:    args.Phone,
		Deleted:  false,
	}
}

func PbShopInventoryVouchers(inventory []*inventory.InventoryVoucher) []*pbshop.InventoryVoucher {
	var inventoryVouchers []*pbshop.InventoryVoucher
	for _, value := range inventory {
		inventoryVouchers = append(inventoryVouchers, PbShopInventoryVoucher(value))
	}
	return inventoryVouchers
}

func PbShopVariants(items []*catalog.ShopVariant) []*pbshop.ShopVariant {
	res := make([]*pbshop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalog.ShopVariant) *pbshop.ShopVariant {
	res := &pbshop.ShopVariant{
		Id: m.VariantID,
		Info: &pbshop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			CostPrice:   m.CostPrice,
			Attributes:  convertpb.PbAttributes(m.Attributes),
		},
		Code:        m.Code,
		EdCode:      m.Code,
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		DescHtml:    m.DescHTML,
		ImageUrls:   m.ImageURLs,
		Tags:        nil,
		Note:        m.Note,
		Status:      pbs3.Pb(model.Status3(m.Status)),
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt32(m.RetailPrice, m.ListPrice),
		CostPrice:   m.CostPrice,
		Attributes:  convertpb.PbAttributes(m.Attributes),
	}
	return res
}

func PbShopCategory(m *catalog.ShopCategory) *pbshop.ShopCategory {
	res := &pbshop.ShopCategory{
		Id:       m.ID,
		ShopId:   m.ShopID,
		Status:   0,
		ParentId: m.ParentID,
		Name:     m.Name,
	}
	return res
}

func PbShopCollection(m *catalog.ShopCollection) *pbshop.ShopCollection {
	res := &pbshop.ShopCollection{
		Id:          m.ID,
		ShopId:      m.ShopID,
		Description: m.Description,
		DescHtml:    m.DescHTML,
		Name:        m.Name,
		ShortDesc:   m.ShortDesc,
	}
	return res
}

func PbShopCollections(items []*catalog.ShopCollection) []*pbshop.ShopCollection {
	res := make([]*pbshop.ShopCollection, len(items))
	for i, item := range items {
		res[i] = PbShopCollection(item)
	}
	return res
}

func PbShopProductsWithVariants(items []*catalog.ShopProductWithVariants) []*pbshop.ShopProduct {
	res := make([]*pbshop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProductWithVariants(item)
	}
	return res
}

func PbShopCategories(items []*catalog.ShopCategory) []*pbshop.ShopCategory {
	res := make([]*pbshop.ShopCategory, len(items))
	for i, item := range items {
		res[i] = PbShopCategory(item)
	}
	return res
}

func PbShopVariantWithProduct(m *catalog.ShopVariantWithProduct) *pbshop.ShopVariant {
	if m == nil {
		return nil
	}
	res := &pbshop.ShopVariant{
		Id: m.VariantID,
		Info: &pbshop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			CostPrice:   m.CostPrice,
			Attributes:  convertpb.PbAttributes(m.Attributes),
		},
		Code:        m.Code,
		EdCode:      m.Code,
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		DescHtml:    m.DescHTML,
		ImageUrls:   m.ImageURLs,
		Tags:        nil,
		Note:        m.Note,
		Status:      pbs3.Pb(model.Status3(m.Status)),
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt32(m.RetailPrice, m.ListPrice),
		CostPrice:   m.CostPrice,
		Attributes:  convertpb.PbAttributes(m.Attributes),
	}
	if m.ShopProduct != nil {
		res.Product = &pbshop.ShopShortProduct{
			Id:   m.ShopProduct.ProductID,
			Name: m.ShopProduct.Name,
		}
	}

	return res
}

func PbShopVariantsWithProducts(items []*catalog.ShopVariantWithProduct) []*pbshop.ShopVariant {
	res := make([]*pbshop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariantWithProduct(item)
	}
	return res
}

func PbShopProductWithVariants(m *catalog.ShopProductWithVariants) *pbshop.ShopProduct {
	shopID := m.ShopProduct.ShopID
	metaFields := []*common.MetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &common.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	res := &pbshop.ShopProduct{
		Id: m.ShopProduct.ProductID,
		Info: &pbshop.EtopProduct{
			Id:          m.ShopProduct.ProductID,
			Code:        m.ShopProduct.Code,
			Name:        m.ShopProduct.Name,
			Description: m.ShopProduct.Description,
			ShortDesc:   m.ShopProduct.ShortDesc,
			DescHtml:    m.ShopProduct.DescHTML,
			Unit:        m.ShopProduct.Unit,
			ImageUrls:   m.ShopProduct.ImageURLs,
			ListPrice:   m.ShopProduct.ListPrice,
			CostPrice:   m.ShopProduct.CostPrice,
			CategoryId:  m.ShopProduct.CategoryID,
			// deprecated
			ProductSourceCategoryId: m.ShopProduct.CategoryID,
		},
		CategoryId:      m.ShopProduct.CategoryID,
		Code:            m.ShopProduct.Code,
		EdCode:          m.ShopProduct.Code,
		Name:            m.ShopProduct.Name,
		Description:     m.ShopProduct.Description,
		ShortDesc:       m.ShopProduct.ShortDesc,
		DescHtml:        m.ShopProduct.DescHTML,
		ImageUrls:       m.ShopProduct.ImageURLs,
		Tags:            m.ShopProduct.Tags,
		Stags:           nil,
		Note:            m.Note,
		Status:          pbs3.Pb(model.Status3(m.ShopProduct.Status)),
		IsAvailable:     false,
		ListPrice:       m.ShopProduct.ListPrice,
		RetailPrice:     coalesceInt32(m.ShopProduct.RetailPrice, m.ShopProduct.ListPrice),
		CostPrice:       m.ShopProduct.CostPrice,
		CollectionIds:   m.ShopProduct.CollectionIDs,
		Variants:        PbShopVariants(m.Variants),
		ProductSourceId: shopID, // backward-compatible: use shop_id in place of product_source_id
		CreatedAt:       pbcm.PbTime(m.CreatedAt),
		UpdatedAt:       pbcm.PbTime(m.UpdatedAt),
		ProductType:     pbproducttype.PbProductType(string(m.ProductType)),
		MetaFields:      metaFields,
		BrandId:         m.BrandID,
	}
	return res
}

func coalesce(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func coalesceInt32(is ...int32) int32 {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func coalesceStrings(sss ...[]string) []string {
	for _, ss := range sss {
		if len(ss) != 0 {
			return ss
		}
	}
	return nil
}

func merge(sss ...[]string) []string {
	s0 := sss[0]
	for _, ss := range sss[1:] {
		for _, s := range ss {
			if !contain(s0, s) {
				s0 = append(s0, s)
			}
		}
	}
	return s0
}

func contain(ss []string, s string) bool {
	for _, S := range ss {
		if S == s {
			return true
		}
	}
	return false
}
