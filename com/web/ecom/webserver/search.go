package webserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"

	"o.o/api/meta"
	"o.o/api/webserver"
	"o.o/capi/dot"
)

type DataSearch struct {
	Website *webserver.WsWebsite
	Search  []*webserver.WsProduct
	*IndexData
	Paging            *CategoryPaging
	SearchNameProduct string
	Host              string
	ProductCount      int
}

func (s *Server) Search(c echo.Context) error {
	path := c.Request().URL.Path
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) != 5 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	nameProduct := pathSplit[2]
	shopID, err := GetShopID(c)
	if err != nil {
		return err
	}
	queryWsProduct := &webserver.SearchProductByNameQuery{
		ShopID: shopID,
		Name:   nameProduct,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryWsProduct)
	if err != nil {
		return err
	}

	var productIDs []dot.ID
	for _, v := range queryWsProduct.Result.WsProducts {
		if v.Appear && len(v.Product.Variants) > 0 {
			productIDs = append(productIDs, v.ID)
		}
	}

	// paging
	page, err := strconv.Atoi(pathSplit[3])
	if err != nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	countProduct, err := strconv.Atoi(pathSplit[4])
	if err != nil {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	if countProduct != 12 && countProduct != 24 {
		return c.Render(http.StatusOK, "404.html", nil)
	}

	queryProductWithPaging := &webserver.ListWsProductsByIDsWithPagingQuery{
		ShopID: shopID,
		IDs:    productIDs,
		Paging: meta.Paging{
			Offset: (page - 1) * countProduct,
			Limit:  countProduct,
			Sort:   []string{"name"},
		},
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), queryProductWithPaging)
	if err != nil {
		return err
	}

	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	indexData.SearchKey = nameProduct
	var paging = &CategoryPaging{
		Max:       len(productIDs)/countProduct + 1,
		Selection: page,
		Limit:     countProduct,
	}
	for i := 1; i <= paging.Max; i++ {
		paging.ListPage = append(paging.ListPage, i)
	}
	var data = &DataSearch{
		Search:            queryProductWithPaging.Result.WsProducts,
		IndexData:         indexData,
		Paging:            paging,
		SearchNameProduct: nameProduct,
		ProductCount:      len(productIDs),
		Host:              c.Request().Host,
	}
	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "search.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}
