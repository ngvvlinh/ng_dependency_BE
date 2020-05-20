package webserver

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"o.o/api/webserver"
	"o.o/capi/dot"
)

type DataPage struct {
	*IndexData
	Page        *webserver.WsPage
	ContentPage template.HTML
}

func (s *Server) Page(c echo.Context) error {
	path := c.Request().URL.Path
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) != 3 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	pageSubPath := pathSplit[2]
	pageSubPathSplit := strings.Split(pageSubPath, "-")
	if len(pageSubPathSplit) == 0 {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	pageIDString := pageSubPathSplit[len(pageSubPathSplit)-1]
	pageID, err := dot.ParseID(pageIDString)
	if err != nil {
		return err
	}
	shopID, err := GetShopID(c)
	if err != nil {
		return err
	}
	query := &webserver.GetWsPageByIDQuery{
		ShopID: shopID,
		ID:     pageID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return err
	}
	page := query.Result
	if page.Slug == "" {
		page.Slug = fmt.Sprintf("page-%v", page.ID)
	}
	if !query.Result.Appear {
		return c.Render(http.StatusOK, "404.html", nil)
	}
	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	contentPage := template.HTML(page.DescHTML)

	data := &DataPage{
		Page:        page,
		IndexData:   indexData,
		ContentPage: contentPage,
	}
	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "post-without-sidebar.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}
