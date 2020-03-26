package webserver

import (
	"net/http"

	"github.com/labstack/echo"

	"o.o/backend/pkg/common/redis"
)

type DataSuccessOrder struct {
	OrderCode string
	*IndexData
}

func (s *Server) Checkout(c echo.Context) error {
	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}

	err = c.Render(http.StatusOK, "header.html", indexData)
	if err != nil {
		return err
	}
	if len(indexData.Cart.Products) == 0 {
		err = c.Render(http.StatusOK, "cart-empty.html", indexData)
		if err != nil {
			return err
		}
	} else {
		err = c.Render(http.StatusOK, "checkout.html", indexData)
		if err != nil {
			return err
		}
	}
	err = c.Render(http.StatusOK, "footer.html", indexData)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) CartOrder(c echo.Context) error {
	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	orderCode, err := getSessionOrder(c)
	if err != nil && err != redis.ErrNil {
		return err
	}
	if err != nil && err == redis.ErrNil {
		return c.Redirect(http.StatusFound, "/")
	}
	data := &DataSuccessOrder{
		IndexData: indexData,
		OrderCode: orderCode,
	}
	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "checkout-success.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}
