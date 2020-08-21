package ws

import (
	"o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/product"
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
		Category:  category.PbShopCategory(arg.Category),
	}
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
		Product:      product.PbShopProductWithVariants(arg.Product),
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

func PbWsProducts(arg []*webserver.WsProduct) []*shop.WsProduct {
	var wsProducts []*shop.WsProduct
	for _, value := range arg {
		wsProducts = append(wsProducts, PbWsProduct(value))
	}
	return wsProducts
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
