package shop

import (
	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/purchaserefund"
	"o.o/api/main/refund"
	"o.o/api/main/stocktaking"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func PbWsPages(arg []*webserver.WsPage) []*shop.WsPage {
	var wsPages []*shop.WsPage
	for _, value := range arg {
		wsPages = append(wsPages, PbWsPage(value))
	}
	return wsPages
}

func PbWsPage(arg *webserver.WsPage) *shop.WsPage {
	if arg == nil {
		return nil
	}
	return &shop.WsPage{
		ShopID:    arg.ShopID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		Name:      arg.Name,
		ID:        arg.ID,
		SEOConfig: PbWsSEOConfig(arg.SEOConfig),
		Slug:      arg.Slug,
		Appear:    arg.Appear,
		DescHTML:  arg.DescHTML,
	}
}

func PbWsCategories(arg []*webserver.WsCategory) []*shop.WsCategory {
	var wsCategories []*shop.WsCategory
	for _, value := range arg {
		wsCategories = append(wsCategories, PbWsCategory(value))
	}
	return wsCategories
}

func PbWsCategory(arg *webserver.WsCategory) *shop.WsCategory {
	if arg == nil {
		return nil
	}
	return &shop.WsCategory{
		Image:     arg.Image,
		ShopID:    arg.ShopID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		ID:        arg.ID,
		SEOConfig: PbWsSEOConfig(arg.SEOConfig),
		Slug:      arg.Slug,
		Appear:    arg.Appear,
		Category:  PbShopCategory(arg.Category),
	}
}

func PbWsProducts(arg []*webserver.WsProduct) []*shop.WsProduct {
	var wsProducts []*shop.WsProduct
	for _, value := range arg {
		wsProducts = append(wsProducts, PbWsProduct(value))
	}
	return wsProducts
}

func PbWsProduct(arg *webserver.WsProduct) *shop.WsProduct {
	if arg == nil {
		return nil
	}
	return &shop.WsProduct{
		ShopID:       arg.ShopID,
		CreatedAt:    arg.CreatedAt,
		UpdatedAt:    arg.UpdatedAt,
		ID:           arg.ID,
		SEOConfig:    PbWsSEOConfig(arg.SEOConfig),
		Slug:         arg.Slug,
		Appear:       arg.Appear,
		ComparePrice: PbComparePrice(arg.ComparePrice),
		DescHTML:     arg.DescHTML,
		Product:      PbShopProductWithVariants(arg.Product),
		Sale:         arg.IsSale,
	}
}

func PbComparePrice(arg []*webserver.ComparePrice) []*shop.ComparePrice {
	if arg == nil {
		return nil
	}
	var listComparePrice []*shop.ComparePrice
	for _, value := range arg {
		listComparePrice = append(listComparePrice, &shop.ComparePrice{
			VariantID:    value.VariantID,
			ComparePrice: value.ComparePrice,
		})
	}
	return listComparePrice
}

func PbWsSEOConfig(arg *webserver.WsSEOConfig) *shop.WsSEOConfig {
	if arg == nil {
		return nil
	}
	return &shop.WsSEOConfig{
		Content:     arg.Content,
		Keyword:     arg.Keyword,
		Description: arg.Description,
	}
}
func PbWsWebsites(arg []*webserver.WsWebsite) []*shop.WsWebsite {
	var wsWebsites []*shop.WsWebsite
	for _, value := range arg {
		wsWebsites = append(wsWebsites, PbWsWebsite(value))
	}
	return wsWebsites
}

func PbWsWebsite(arg *webserver.WsWebsite) *shop.WsWebsite {
	return &shop.WsWebsite{
		ShopID:             arg.ShopID,
		ID:                 arg.ID,
		MainColor:          arg.MainColor,
		Banner:             PbWsBanner(arg.Banner),
		OutstandingProduct: PbWsSpecialProduct(arg.OutstandingProduct),
		NewProduct:         PbWsSpecialProduct(arg.NewProduct),
		SEOConfig:          PbWsGeneralSEO(arg.SEOConfig),
		Facebook:           PbWsFacebook(arg.Facebook),
		GoogleAnalyticsID:  arg.GoogleAnalyticsID,
		DomainName:         arg.DomainName,
		OverStock:          arg.OverStock,
		ShopInfo:           PbWsShopInfo(arg.ShopInfo),
		Description:        arg.Description,
		LogoImage:          arg.LogoImage,
		FaviconImage:       arg.FaviconImage,
		UpdatedAt:          arg.UpdatedAt,
		CreatedAt:          arg.CreatedAt,
		SiteSubdomain:      arg.SiteSubdomain,
	}
}

func PbWsShopInfo(arg *webserver.ShopInfo) *shop.ShopInfo {
	if arg == nil {
		return nil
	}
	var address *shop.AddressShopInfo
	if arg.Address != nil {
		address = &shop.AddressShopInfo{
			Province:     arg.Address.Province,
			District:     arg.Address.District,
			Ward:         arg.Address.Ward,
			DistrictCode: arg.Address.DistrictCode,
			ProvinceCode: arg.Address.ProvinceCode,
			WardCode:     arg.Address.WardCode,
			Address:      arg.Address.Address,
		}
	}
	return &shop.ShopInfo{
		Email:           arg.Email,
		Name:            arg.Name,
		Phone:           arg.Phone,
		Address:         address,
		FacebookFanpage: arg.FacebookFanpage,
	}
}
func PbWsFacebook(arg *webserver.Facebook) *shop.Facebook {
	if arg == nil {
		return nil
	}
	return &shop.Facebook{
		FacebookID:     arg.FacebookID,
		WelcomeMessage: arg.WelcomeMessage,
	}
}

func PbWsGeneralSEO(arg *webserver.WsGeneralSEO) *shop.WsGeneralSEO {
	if arg == nil {
		return nil
	}
	return &shop.WsGeneralSEO{
		Title:               arg.Title,
		SiteContent:         arg.SiteContent,
		SiteMetaKeyword:     arg.SiteMetaKeyword,
		SiteMetaDescription: arg.SiteMetaDescription,
	}
}

func PbWsSpecialProduct(arg *webserver.SpecialProduct) *shop.SpecialProduct {
	if arg == nil {
		return nil
	}
	return &shop.SpecialProduct{
		ProductIDs: arg.ProductIDs,
		Style:      arg.Style,
	}
}

func PbWsBanners(arg []*webserver.Banner) []*shop.Banner {
	if len(arg) == 0 {
		return nil
	}
	var result []*shop.Banner
	for _, v := range arg {
		result = append(result, PbWsBanner(v))
	}
	return result
}

func PbWsBanner(arg *webserver.Banner) *shop.Banner {
	if arg == nil {
		return nil
	}
	return &shop.Banner{
		BannerItems: PbWsBannerItems(arg.BannerItems),
		Style:       arg.Style,
	}
}

func PbWsBannerItems(arg []*webserver.BannerItem) []*shop.BannerItem {
	if arg == nil {
		return nil
	}
	var result []*shop.BannerItem
	for _, v := range arg {
		result = append(result, &shop.BannerItem{
			Alt:   v.Alt,
			Url:   v.Url,
			Image: v.Image,
		})
	}
	return result
}

func ConvertShopInfo(args *shop.ShopInfo) *webserver.ShopInfo {
	var shopInfo *webserver.ShopInfo
	if args != nil {
		shopInfo = &webserver.ShopInfo{}
		shopInfo.Phone = args.Phone
		shopInfo.FacebookFanpage = args.FacebookFanpage
		shopInfo.Name = args.Name
		shopInfo.Email = args.Email
		if args.Address != nil {
			shopInfo.Address = &webserver.AddressShopInfo{
				Province:     args.Address.Province,
				District:     args.Address.District,
				Ward:         args.Address.Ward,
				DistrictCode: args.Address.DistrictCode,
				ProvinceCode: args.Address.ProvinceCode,
				WardCode:     args.Address.WardCode,
				Address:      args.Address.Address,
			}
		}

	}
	return shopInfo
}

func ConvertBanner(args *shop.Banner) *webserver.Banner {
	var banner *webserver.Banner
	if args != nil {
		banner = &webserver.Banner{}
		banner.Style = args.Style
		var bannerItems = []*webserver.BannerItem{}
		for _, item := range args.BannerItems {
			bannerItems = append(bannerItems, &webserver.BannerItem{
				Alt:   item.Alt,
				Url:   item.Url,
				Image: item.Image,
			})
		}
		banner.BannerItems = bannerItems
	}
	return banner
}

func ConvertFacebook(args *shop.Facebook) *webserver.Facebook {
	var facebook *webserver.Facebook
	if args != nil {
		facebook = &webserver.Facebook{}
		facebook.FacebookID = args.FacebookID
		facebook.WelcomeMessage = args.WelcomeMessage
	}
	return facebook
}

func ConvertSpecialProduct(args *shop.SpecialProduct) *webserver.SpecialProduct {
	var specialProduct *webserver.SpecialProduct
	if args != nil {
		specialProduct = &webserver.SpecialProduct{}
		specialProduct.Style = args.Style
		specialProduct.ProductIDs = args.ProductIDs
	}
	return specialProduct
}

func ConvertSEOConfig(args *shop.WsSEOConfig) *webserver.WsSEOConfig {
	var wsSEOConfig *webserver.WsSEOConfig
	if args != nil {
		wsSEOConfig = &webserver.WsSEOConfig{}
		wsSEOConfig.Content = args.Content
		wsSEOConfig.Description = args.Description
		wsSEOConfig.Keyword = args.Keyword
	}
	return wsSEOConfig
}

func ConvertComparePrice(args []*shop.ComparePrice) []*webserver.ComparePrice {
	var listComparePrice []*webserver.ComparePrice
	for _, value := range args {
		comparePriceItem := &webserver.ComparePrice{
			VariantID:    value.VariantID,
			ComparePrice: value.ComparePrice,
		}
		listComparePrice = append(listComparePrice, comparePriceItem)
	}
	return listComparePrice
}

func ConvertWsGeneralSEO(args *shop.WsGeneralSEO) *webserver.WsGeneralSEO {
	var wsGeneralSEO *webserver.WsGeneralSEO
	if args != nil {
		wsGeneralSEO = &webserver.WsGeneralSEO{}
		wsGeneralSEO.Title = args.Title
		wsGeneralSEO.SiteContent = args.SiteContent
		wsGeneralSEO.SiteMetaDescription = args.SiteMetaDescription
		wsGeneralSEO.SiteMetaKeyword = args.SiteMetaKeyword
	}
	return wsGeneralSEO
}

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
	if args == nil {
		return nil
	}
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
		RefAction:    args.RefAction,
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
