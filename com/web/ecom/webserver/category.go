package webserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"o.o/api/meta"
	"o.o/api/webserver"
	"o.o/capi/dot"
)

type DataCategory struct {
	Website  *webserver.WsWebsite
	Category *WsCategoriesWithProducstWithVariants
	*IndexData
	Paging       *CategoryPaging
	Host         string
	ProductCount int
}

type CategoryPaging struct {
	Max       int
	Selection int
	Limit     int
	ListPage  []int
}

func (s *Server) CategoryWithPagePaging(c echo.Context) error {
	path := c.Request().URL.Path
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) != 5 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	categorySubPath := pathSplit[2]
	categorySubPathSplit := strings.Split(categorySubPath, "-")
	if len(categorySubPathSplit) == 0 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	categoryIDString := categorySubPathSplit[len(categorySubPathSplit)-1]
	categoryID, err := dot.ParseID(categoryIDString)
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(pathSplit[3])
	if err != nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	countProduct, err := strconv.Atoi(pathSplit[4])
	if err != nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}

	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	var category *WsCategoriesWithProducstWithVariants
	for _, v := range indexData.Categories {
		if v.ID == categoryID {
			category = v
		}
	}
	if category == nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}

	var productIDs []dot.ID
	var mapProducts = make(map[dot.ID]*webserver.WsProduct)
	for _, v := range category.Products {
		if !contains(productIDs, v.ID) {
			productIDs = append(productIDs, v.ID)
			mapProducts[v.ID] = v
		}
	}
	shopID, err := GetShopID(c)
	if err != nil {
		return err
	}
	queryProduct := &webserver.ListWsProductsByIDsWithPagingQuery{
		ShopID: shopID,
		IDs:    productIDs,
		Paging: meta.Paging{
			Offset: (page - 1) * countProduct,
			Limit:  countProduct,
			Sort:   []string{"name"},
		},
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryProduct)
	if err != nil {
		return err
	}
	category.Products = queryProduct.Result.WsProducts
	var paging = &CategoryPaging{
		Max:       len(productIDs)/countProduct + 1,
		Selection: page,
		Limit:     countProduct,
	}
	for i := 1; i <= paging.Max; i++ {
		paging.ListPage = append(paging.ListPage, i)
	}
	var data = &DataCategory{
		Category:     category,
		IndexData:    indexData,
		Paging:       paging,
		ProductCount: len(productIDs),
		Host:         c.Request().Host,
	}

	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "category.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Category(c echo.Context) error {
	path := c.Request().URL.Path
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) != 3 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	categorySubPath := pathSplit[2]
	categorySubPathSplit := strings.Split(categorySubPath, "-")
	if len(categorySubPathSplit) == 0 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	categoryIDString := categorySubPathSplit[len(categorySubPathSplit)-1]
	categoryID, err := dot.ParseID(categoryIDString)
	if err != nil {
		return err
	}

	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	var category *WsCategoriesWithProducstWithVariants
	for _, v := range indexData.Categories {
		if v.ID == categoryID {
			category = v
		}
	}
	if category == nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	var data = &DataCategory{
		Category:  category,
		IndexData: indexData,
		Host:      c.Request().Host,
	}

	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "category.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}

func contains(s []dot.ID, e dot.ID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
