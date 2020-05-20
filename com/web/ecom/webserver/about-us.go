package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) AboutUs(c echo.Context) error {
	data, err := GetIndexData(c)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "about-us.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}
