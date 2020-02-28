package webserver

import (
	"net/http"

	"github.com/labstack/echo"

	"etop.vn/api/webserver"
	"etop.vn/backend/com/web/ecom/middlewares"
	"etop.vn/capi/dot"
)

type indexData struct {
	Site       string
	Title      string
	Categories []map[string]interface{}
	Items      []map[string]interface{}
}

func (s *Server) Index(c echo.Context) error {
	// cmd := &webserver.ListWsProductsQuery{
	// 	ShopID: 0,
	// }
	// err := webserverQueryBus.Dispatch(c.Request().Context(), cmd)
	// if err != nil {
	// 	return err
	// }
	//TODO shop_id
	shopID := dot.ID(1048069941737578021)
	query := &webserver.ListWsCategoriesQuery{
		ShopID: shopID,
		Result: nil,
	}
	err := webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return err
	}
	site := middlewares.GetSite(c)
	if site == "" {
		site = "<empty X-Forwarded-Host>"
	}
	var categories = make([]map[string]interface{}, len(query.Result.WsCategories))
	for _, v := range query.Result.WsCategories {
		categories = append(categories, map[string]interface{}{
			"Name": v.Category.Name,
			"Slug": v.Slug,
		})
	}
	data := &indexData{
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
	}
	return c.Render(http.StatusOK, "index.html", data)
}
