package shop

import (
	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/purchaserefund"
	"etop.vn/api/main/refund"
	"etop.vn/api/main/stocktaking"
	"etop.vn/api/top/int/shop"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/capi/dot"
)

func PbPurchaseRefunds(args []*purchaserefund.PurchaseRefund) []*shop.PurchaseRefund {
	var result []*shop.PurchaseRefund
	for _, value := range args {
		result = append(result, PbPurchaseRefund(value))
	}
	return result
}

func PbPurchaseRefund(args *purchaserefund.PurchaseRefund) *shop.PurchaseRefund {
	var result = &shop.PurchaseRefund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		PurchaseOrderID: args.PurchaseOrderID,
		Note:            args.Note,
		Code:            args.Code,
		TotalAdjustment: args.TotalAdjustment,
		AdjustmentLines: args.AdjustmentLines,
		Lines:           PbPurchaseRefundLine(args.Lines),
		CreatedAt:       cmapi.PbTime(args.CreatedAt),
		UpdatedAt:       cmapi.PbTime(args.UpdatedAt),
		CancelledAt:     cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:     cmapi.PbTime(args.ConfirmedAt),
		CreatedBy:       args.CreatedBy,
		UpdatedBy:       args.UpdatedBy,
		CancelReason:    args.CancelReason,
		Status:          args.Status,
		TotalAmount:     args.TotalAmount,
		BasketValue:     args.BasketValue,
	}
	return result
}

func PbPurchaseRefundLine(args []*purchaserefund.PurchaseRefundLine) []*shop.PurchaseRefundLine {
	var result []*shop.PurchaseRefundLine
	for _, v := range args {
		result = append(result, &shop.PurchaseRefundLine{
			VariantID:    v.VariantID,
			ProductID:    v.ProductID,
			Quantity:     v.Quantity,
			Code:         v.Code,
			ImageURL:     v.ImageURL,
			Name:         v.ProductName,
			PaymentPrice: v.PaymentPrice,
			Attributes:   v.Attributes,
			Adjustment:   v.Adjustment,
		})
	}
	return result
}

func PbRefunds(args []*refund.Refund) []*shop.Refund {
	var result []*shop.Refund
	for _, value := range args {
		result = append(result, PbRefund(value))
	}
	return result
}

func PbRefund(args *refund.Refund) *shop.Refund {
	var result = &shop.Refund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		OrderID:         args.OrderID,
		Note:            args.Note,
		Code:            args.Code,
		AdjustmentLines: args.AdjustmentLines,
		TotalAdjustment: args.TotalAdjustment,
		Lines:           PbRefundLine(args.Lines),
		CreatedAt:       cmapi.PbTime(args.CreatedAt),
		UpdatedAt:       cmapi.PbTime(args.UpdatedAt),
		CancelledAt:     cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:     cmapi.PbTime(args.ConfirmedAt),
		CreatedBy:       args.CreatedBy,
		UpdatedBy:       args.UpdatedBy,
		CancelReason:    args.CancelReason,
		Status:          args.Status,
		TotalAmount:     args.TotalAmount,
		BasketValue:     args.BasketValue,
	}
	return result
}

func PbRefundLine(args []*refund.RefundLine) []*shop.RefundLine {
	var result []*shop.RefundLine
	for _, v := range args {
		result = append(result, &shop.RefundLine{
			VariantID:   v.VariantID,
			ProductID:   v.ProductID,
			Quantity:    v.Quantity,
			Code:        v.Code,
			ImageURL:    v.ImageURL,
			Name:        v.ProductName,
			RetailPrice: v.RetailPrice,
			Attributes:  v.Attributes,
			Adjustment:  v.Adjustment,
		})
	}
	return result
}

func PbProductsQuantity(shopProducts []*catalog.ShopProductWithVariants, inventoryVariants map[dot.ID]*inventory.InventoryVariant) (res []*shop.ShopProduct) {
	for _, product := range shopProducts {
		productPb := PbProductQuantity(product, inventoryVariants)
		res = append(res, productPb)
	}
	return
}

func PbProductQuantity(shopProduct *catalog.ShopProductWithVariants, inventoryVariants map[dot.ID]*inventory.InventoryVariant) *shop.ShopProduct {
	shopProductPb := PbShopProductWithVariants(shopProduct)
	shopProductPb.Variants = PbVariantsQuantity(shopProduct.Variants, inventoryVariants)
	return shopProductPb
}

func PbVariantsQuantity(shopVariants []*catalog.ShopVariant, inventoryVariants map[dot.ID]*inventory.InventoryVariant) []*shop.ShopVariant {
	var variants []*shop.ShopVariant
	for _, variant := range shopVariants {
		inventoryVariant := inventoryVariants[variant.VariantID]
		valuePb := PbVariantQuantity(variant, inventoryVariant)
		variants = append(variants, valuePb)
	}
	return variants
}

func PbVariantQuantity(shopVariant *catalog.ShopVariant, inventoryVariant *inventory.InventoryVariant) *shop.ShopVariant {
	shopVariantDB := PbShopVariant(shopVariant)
	if inventoryVariant != nil {
		shopVariantDB.InventoryVariant = &shop.InventoryVariantShopVariant{
			QuantityOnHand: inventoryVariant.QuantityOnHand,
			QuantityPicked: inventoryVariant.QuantityPicked,
			Quantity:       inventoryVariant.QuantitySummary,
			CostPrice:      inventoryVariant.CostPrice,
		}
		shopVariantDB.QuantityOnHand = inventoryVariant.QuantityOnHand
		shopVariantDB.QuantityPicked = inventoryVariant.QuantityPicked
		shopVariantDB.Quantity = inventoryVariant.QuantitySummary
		shopVariantDB.CostPrice = inventoryVariant.CostPrice
	}
	return shopVariantDB
}

func PbStocktakes(args []*stocktaking.ShopStocktake) []*shop.Stocktake {
	var stocktakesPb []*shop.Stocktake
	for _, value := range args {
		stocktakesPb = append(stocktakesPb, PbStocktake(value))
	}
	return stocktakesPb
}

func PbStocktake(args *stocktaking.ShopStocktake) *shop.Stocktake {
	return &shop.Stocktake{
		Id:            args.ID,
		ShopId:        args.ShopID,
		TotalQuantity: args.TotalQuantity,
		Note:          args.Note,
		CreatedBy:     args.CreatedBy,
		UpdatedBy:     args.UpdatedBy,
		CancelReason:  args.CancelReason,
		CreatedAt:     cmapi.PbTime(args.CreatedAt),
		UpdatedAt:     cmapi.PbTime(args.UpdatedAt),
		ConfirmedAt:   cmapi.PbTime(args.ConfirmedAt),
		CancelledAt:   cmapi.PbTime(args.CancelledAt),
		Status:        args.Status,
		Type:          args.Type.String(),
		Code:          args.Code,
		Lines:         PbstocktakeLines(args.Lines),
	}
}

func PbstocktakeLines(args []*stocktaking.StocktakeLine) []*shop.StocktakeLine {
	var lines []*shop.StocktakeLine
	for _, value := range args {
		lines = append(lines, &shop.StocktakeLine{
			VariantId:   value.VariantID,
			OldQuantity: value.OldQuantity,
			NewQuantity: value.NewQuantity,
			VariantName: value.VariantName,
			ProductName: value.ProductName,
			CostPrice:   value.CostPrice,
			ProductId:   value.ProductID,
			Code:        value.Code,
			ImageUrl:    value.ImageURL,
			Attributes:  value.Attributes,
		})
	}
	return lines
}

func PbInventory(args *inventory.InventoryVariant) *shop.InventoryVariant {
	return &shop.InventoryVariant{
		ShopId:         args.ShopID,
		VariantId:      args.VariantID,
		QuantityOnHand: args.QuantityOnHand,
		QuantityPicked: args.QuantityPicked,
		Quantity:       args.QuantitySummary,
		CostPrice:      args.CostPrice,
		CreatedAt:      cmapi.PbTime(args.CreatedAt),
		UpdatedAt:      cmapi.PbTime(args.UpdatedAt),
	}
}

func PbInventoryVariants(args []*inventory.InventoryVariant) []*shop.InventoryVariant {
	var inventoryVariants []*shop.InventoryVariant
	for _, value := range args {
		inventoryVariants = append(inventoryVariants, PbInventory(value))
	}
	return inventoryVariants
}

func PbBrand(args *catalog.ShopBrand) *shop.Brand {
	return &shop.Brand{
		ShopId:      args.ShopID,
		Id:          args.ID,
		Name:        args.BrandName,
		Description: args.Description,
		CreatedAt:   cmapi.PbTime(args.CreatedAt),
		UpdatedAt:   cmapi.PbTime(args.UpdatedAt),
	}
}

func PbBrands(args []*catalog.ShopBrand) []*shop.Brand {
	var brands []*shop.Brand
	for _, value := range args {
		brands = append(brands, PbBrand(value))
	}
	return brands
}

func PbShopInventoryVoucher(args *inventory.InventoryVoucher) *shop.InventoryVoucher {
	if args == nil {
		return nil
	}

	var inventoryVoucherItem []*shop.InventoryVoucherLine
	for _, value := range args.Lines {
		inventoryVoucherItem = append(inventoryVoucherItem, &shop.InventoryVoucherLine{
			VariantId:   value.VariantID,
			VariantName: value.VariantName,
			ProductId:   value.ProductID,
			Code:        value.Code,
			ProductName: value.ProductName,
			ImageUrl:    value.ImageURL,
			Attributes:  value.Attributes,
			Price:       value.Price,
			Quantity:    value.Quantity,
		})
	}
	return &shop.InventoryVoucher{
		Title:        args.Title,
		TotalAmount:  args.TotalAmount,
		CreatedBy:    args.CreatedBy,
		UpdatedBy:    args.UpdatedBy,
		Lines:        inventoryVoucherItem,
		RefId:        args.RefID,
		Code:         args.Code,
		RefCode:      args.RefCode,
		RefType:      args.RefType.String(),
		RefName:      args.RefName,
		TraderId:     args.TraderID,
		Trader:       PbShopTrader(args.Trader),
		Status:       args.Status,
		Note:         args.Note,
		Type:         args.Type.String(),
		Id:           args.ID,
		ShopId:       args.ShopID,
		CreatedAt:    cmapi.PbTime(args.CreatedAt),
		UpdatedAt:    cmapi.PbTime(args.UpdatedAt),
		CancelledAt:  cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:  cmapi.PbTime(args.ConfirmedAt),
		CancelReason: args.CancelReason,
		Rollback:     args.Rollback,
	}
}

func PbShopTrader(args *inventory.Trader) *shop.Trader {
	if args == nil {
		return nil
	}
	return &shop.Trader{
		Id:       args.ID,
		Type:     args.Type,
		FullName: args.FullName,
		Phone:    args.Phone,
		Deleted:  false,
	}
}

func PbShopInventoryVouchers(inventory []*inventory.InventoryVoucher) []*shop.InventoryVoucher {
	var inventoryVouchers []*shop.InventoryVoucher
	for _, value := range inventory {
		inventoryVouchers = append(inventoryVouchers, PbShopInventoryVoucher(value))
	}
	return inventoryVouchers
}

func PbShopVariants(items []*catalog.ShopVariant) []*shop.ShopVariant {
	res := make([]*shop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalog.ShopVariant) *shop.ShopVariant {
	res := &shop.ShopVariant{
		Id: m.VariantID,
		Info: &shop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			Attributes:  m.Attributes,
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
		Status:      m.Status,
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt(m.RetailPrice, m.ListPrice),
		Attributes:  m.Attributes,
		ProductId:   m.ProductID,
	}
	return res
}

func PbShopCategory(m *catalog.ShopCategory) *shop.ShopCategory {
	res := &shop.ShopCategory{
		Id:       m.ID,
		ShopId:   m.ShopID,
		Status:   0,
		ParentId: m.ParentID,
		Name:     m.Name,
	}
	return res
}

func PbShopCollection(m *catalog.ShopCollection) *shop.ShopCollection {
	res := &shop.ShopCollection{
		Id:          m.ID,
		ShopId:      m.ShopID,
		Description: m.Description,
		DescHtml:    m.DescHTML,
		Name:        m.Name,
		ShortDesc:   m.ShortDesc,
	}
	return res
}

func PbShopCollections(items []*catalog.ShopCollection) []*shop.ShopCollection {
	res := make([]*shop.ShopCollection, len(items))
	for i, item := range items {
		res[i] = PbShopCollection(item)
	}
	return res
}

func PbShopProductsWithVariants(items []*catalog.ShopProductWithVariants) []*shop.ShopProduct {
	res := make([]*shop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProductWithVariants(item)
	}
	return res
}

func PbShopCategories(items []*catalog.ShopCategory) []*shop.ShopCategory {
	res := make([]*shop.ShopCategory, len(items))
	for i, item := range items {
		res[i] = PbShopCategory(item)
	}
	return res
}

func PbShopVariantWithProduct(m *catalog.ShopVariantWithProduct) *shop.ShopVariant {
	if m == nil {
		return nil
	}
	res := &shop.ShopVariant{
		Id: m.VariantID,
		Info: &shop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			CostPrice:   m.CostPrice,
			Attributes:  m.Attributes,
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
		Status:      m.Status,
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt(m.RetailPrice, m.ListPrice),
		Attributes:  m.Attributes,
		ProductId:   m.ShopProduct.ProductID,
	}
	if m.ShopProduct != nil {
		res.Product = &shop.ShopShortProduct{
			Id:   m.ShopProduct.ProductID,
			Name: m.ShopProduct.Name,
		}
	}

	return res
}

func PbShopVariantsWithProducts(items []*catalog.ShopVariantWithProduct) []*shop.ShopVariant {
	res := make([]*shop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariantWithProduct(item)
	}
	return res
}

func PbShopProductWithVariants(m *catalog.ShopProductWithVariants) *shop.ShopProduct {
	shopID := m.ShopProduct.ShopID
	metaFields := []*pbcm.MetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &pbcm.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	res := &shop.ShopProduct{
		Id: m.ShopProduct.ProductID,
		Info: &shop.EtopProduct{
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
		Status:          m.ShopProduct.Status,
		IsAvailable:     false,
		ListPrice:       m.ShopProduct.ListPrice,
		RetailPrice:     coalesceInt(m.ShopProduct.RetailPrice, m.ShopProduct.ListPrice),
		CollectionIds:   m.ShopProduct.CollectionIDs,
		Variants:        PbShopVariants(m.Variants),
		ProductSourceId: shopID, // backward-compatible: use shop_id in place of product_source_id
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
		UpdatedAt:       cmapi.PbTime(m.UpdatedAt),
		ProductType:     m.ProductType.Wrap(),
		MetaFields:      metaFields,
		BrandId:         m.BrandID,
	}
	return res
}

func coalesceInt(is ...int) int {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}
