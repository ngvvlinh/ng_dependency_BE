package webserver

import (
	"net/http"

	"github.com/labstack/echo"

	"etop.vn/backend/com/web/ecom/middlewares"
)

type indexData struct {
	Site  string
	Title string
	Items []map[string]interface{}
}

func (s *Server) Index(c echo.Context) error {
	site := middlewares.GetSite(c)
	if site == "" {
		site = "<empty X-Forwarded-Host>"
	}
	data := &indexData{
		Site:  site,
		Title: "sample page",
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
