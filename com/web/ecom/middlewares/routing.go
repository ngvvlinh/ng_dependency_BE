package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const keySite = "::site"

func GetSite(c echo.Context) string {
	return c.Get(keySite).(string)
}

func SiteRouter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host := c.Request().Header.Get("X-Forwarded-Host")
		site := strings.SplitN(host, ".", 2)[0]

		c.Set(keySite, site)
		return next(c)
	}
}
