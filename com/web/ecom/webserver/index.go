package webserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
	"o.o/backend/com/web/ecom/middlewares"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

type IndexData struct {
	Site               string
	Title              string
	Categories         []*WsCategoriesWithProducstWithVariants
	Pages              []*webserver.WsPage
	Items              []map[string]interface{}
	FaviconImage       string
	MainColor          string
	Banner             *webserver.Banner
	OutstandingProduct *webserver.SpecialProduct
	NewProduct         *webserver.SpecialProduct
	SEOConfig          *webserver.WsGeneralSEO
	Facebook           *webserver.Facebook
	GoogleAnalyticsID  string
	DomainName         string
	OverStock          bool
	ShopInfo           *webserver.ShopInfo
	Description        string
	LogoImage          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Cart               *SessionCart
	SearchKey          string
}

type WsCategoriesWithProducstWithVariants struct {
	Products []*webserver.WsProduct
	*webserver.WsCategory
	CoupleProducts []*CoupleProduct
}

type CoupleProduct struct {
	ProductOne *webserver.WsProduct
	ProductTwo *webserver.WsProduct
}

func (s *Server) Index(c echo.Context) error {
	data, err := GetIndexData(c)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "index.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}

func GetIndexData(c echo.Context) (*IndexData, error) {
	shopID, err := GetShopID(c)
	if err != nil {
		return nil, err
	}
	// menu Catalog
	queryWsCategories := &webserver.ListWsCategoriesQuery{
		ShopID: shopID,
		Result: nil,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryWsCategories)
	if err != nil {
		return nil, err
	}

	site := middlewares.GetSite(c)
	if site == "" {
		site = "<empty X-Forwarded-Host>"
	}
	var categoriesIDs []dot.ID
	var categories []*WsCategoriesWithProducstWithVariants
	for k, v := range queryWsCategories.Result.WsCategories {
		if v.Appear {
			var category = &WsCategoriesWithProducstWithVariants{
				Products:   []*webserver.WsProduct{},
				WsCategory: v,
			}
			categoriesIDs = append(categoriesIDs, v.ID)
			name := v.Category.Name
			slug := fmt.Sprintf("%v-%v", v.Slug, v.ID)
			if v.Slug == "" {
				slug = fmt.Sprintf("category-%v-%v", k, v.ID)
			}
			if name == "" {
				name = fmt.Sprintf("category %v", k)
			}
			category.Slug = slug
			category.Category.Name = name
			categories = append(categories, category)
		}
	}

	var wsProductIDs []dot.ID
	//TODO lay san phan cua tung catalog 1.0
	queryShopProductByCategory := &catalog.ListShopProductWithVariantByCategoriesIDsQuery{
		ShopID:        shopID,
		CategoriesIds: categoriesIDs,
	}
	err = catelogQueryBus.Dispatch(c.Request().Context(), queryShopProductByCategory)
	if err != nil {
		return nil, err
	}
	var mapCategoriesWithProductWithVariant = make(map[dot.ID][]*catalog.ShopProductWithVariants)
	for _, v := range queryShopProductByCategory.Result.Products {
		wsProductIDs = append(wsProductIDs, v.ProductID)
		if mapCategoriesWithProductWithVariant[v.CategoryID] == nil {
			mapCategoriesWithProductWithVariant[v.CategoryID] = []*catalog.ShopProductWithVariants{}
		}
		var listProduct = mapCategoriesWithProductWithVariant[v.CategoryID]
		listProduct = append(listProduct, v)
		mapCategoriesWithProductWithVariant[v.CategoryID] = listProduct
	}
	// lấy sản phẩm 1.1
	queryWsProduct := &webserver.ListWsProductsByIDsQuery{
		IDs:    wsProductIDs,
		ShopID: shopID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryWsProduct)
	if err != nil {
		return nil, err
	}
	var mapWsProduct = make(map[dot.ID]*webserver.WsProduct)
	for _, v := range queryWsProduct.Result {
		mapWsProduct[v.ID] = v
	}
	// lấy sản phẩm 1.2 làm map danh mục map sản phẩm
	categoryWithProductWithVariant := removeProductDisAppear(mapCategoriesWithProductWithVariant, mapWsProduct)
	// main site
	queryWevsites := &webserver.ListWsWebsitesQuery{
		ShopID: shopID,
		Result: nil,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryWevsites)
	if err != nil {
		return nil, err
	}
	var wsWebSiteDefault *webserver.WsWebsite
	if len(queryWevsites.Result.WsWebsites) > 0 {
		wsWebSiteDefault = queryWevsites.Result.WsWebsites[0]
	}
	// if null get default
	if wsWebSiteDefault == nil {
		wsWebSiteDefault = &webserver.WsWebsite{}
	}
	for k, v := range categories {
		categories[k].Products = categoryWithProductWithVariant[v.ID]
	}
	categories = checkCatagoryContainProduct(categories)
	for key, v := range categories {
		// if len(v.Products)%2 == 1 {
		// 	v.Products = append(v.Products, v.Products...)
		// }
		for i := 0; i < len(v.Products); i = i + 2 {
			coupleProduct := &CoupleProduct{}
			coupleProduct.ProductOne = v.Products[i]
			if i+1 < len(v.Products) {
				coupleProduct.ProductTwo = v.Products[i+1]
			}
			categories[key].CoupleProducts = append(categories[key].CoupleProducts, coupleProduct)
		}
	}

	wsPagesQuery := &webserver.ListWsPagesQuery{
		ShopID: shopID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), wsPagesQuery)
	if err != nil {
		return nil, err
	}
	var pages []*webserver.WsPage
	for _, v := range wsPagesQuery.Result.WsPages {
		if v.Appear {
			pages = append(pages, v)
		}
	}
	cart, err := getSessionCart(c)
	if err != nil {
		return nil, err
	}

	data := &IndexData{
		Cart:       cart,
		Pages:      pages,
		Site:       site,
		Title:      "sample page",
		Categories: categories,
		Items: []map[string]interface{}{
			{
				"Label": "Alice",
			},
			{
				"Label": "Bob",
			},
		},
		FaviconImage:       wsWebSiteDefault.FaviconImage,
		MainColor:          wsWebSiteDefault.MainColor,
		Banner:             wsWebSiteDefault.Banner,
		OutstandingProduct: wsWebSiteDefault.OutstandingProduct,
		NewProduct:         wsWebSiteDefault.NewProduct,
		SEOConfig:          wsWebSiteDefault.SEOConfig,
		Facebook:           wsWebSiteDefault.Facebook,
		GoogleAnalyticsID:  wsWebSiteDefault.GoogleAnalyticsID,
		DomainName:         wsWebSiteDefault.DomainName,
		OverStock:          wsWebSiteDefault.OverStock,
		ShopInfo:           wsWebSiteDefault.ShopInfo,
		Description:        wsWebSiteDefault.Description,
		LogoImage:          wsWebSiteDefault.LogoImage,
	}
	return data, nil
}

func getSubDomain(c echo.Context) (string, error) {
	domainParts := strings.Split(c.Request().Host, ".")
	if len(domainParts) < 2 {
		return "", cm.Errorf(cm.NotFound, nil, "Site not found")
	}
	return domainParts[0], nil
}

func removeProductDisAppear(mapValue map[dot.ID][]*catalog.ShopProductWithVariants, products map[dot.ID]*webserver.WsProduct) map[dot.ID][]*webserver.WsProduct {
	var result = make(map[dot.ID][]*webserver.WsProduct)
	for categoryID, listWsProducts := range mapValue {
		if result[categoryID] == nil {
			result[categoryID] = []*webserver.WsProduct{}
		}
		if len(listWsProducts) > 0 {
			for _, product := range listWsProducts {
				if products[product.ProductID].Appear && len(product.Variants) > 0 {
					result[categoryID] = append(result[categoryID], products[product.ProductID])
				}
			}
		}
	}
	return result
}

func checkCatagoryContainProduct(args []*WsCategoriesWithProducstWithVariants) []*WsCategoriesWithProducstWithVariants {
	var result []*WsCategoriesWithProducstWithVariants
	for _, v := range args {
		if len(v.Products) > 0 {
			result = append(result, &WsCategoriesWithProducstWithVariants{
				Products:   v.Products,
				WsCategory: v.WsCategory,
			})
		}
	}
	return result
}
