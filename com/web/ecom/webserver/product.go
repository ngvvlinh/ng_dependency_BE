package webserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"

	"o.o/api/webserver"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var redisStore redis.Store

type DataProduct struct {
	Product  *webserver.WsProduct
	Category *WsCategoriesWithProducstWithVariants
	*IndexData
	Attributes map[string]*AttributeWithOption
	Host       string
}

type AttributeWithOption struct {
	Option []string
}

func (s *Server) Product(c echo.Context) error {
	path := c.Request().URL.Path
	pathSplit := strings.Split(path, "/")
	if len(pathSplit) != 3 {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	productSubPath := pathSplit[2]
	productSubPathSplit := strings.Split(productSubPath, "-")
	if len(productSubPathSplit) == 0 {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	productIDString := productSubPathSplit[len(productSubPathSplit)-1]
	productID, err := dot.ParseID(productIDString)
	if err != nil {
		return err
	}
	shopID, err := GetShopID(c)
	if err != nil {
		return err
	}
	query := &webserver.GetWsProductByIDQuery{
		ShopID: shopID,
		ID:     productID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return err
	}
	product := query.Result
	if product.Slug == "" {
		product.Slug = fmt.Sprintf("product-%v", product.ID)
	}
	if !query.Result.Appear {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	//check product have more than one variant
	if len(query.Result.Product.Variants) == 0 {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	var category = &WsCategoriesWithProducstWithVariants{}
	for _, categoryValue := range indexData.Categories {
		isFound := false
		for _, productValue := range categoryValue.Products {
			if productValue.ID == product.ID {
				category = categoryValue
				isFound = true
				break
			}
		}
		if isFound {
			break
		}
	}
	var concernedProduct []*webserver.WsProduct
	for k, v := range category.Products {
		concernedProduct = append(concernedProduct, v)
		if k > 10 {
			break
		}
	}
	category.Products = concernedProduct

	var mapAttributes = make(map[string]*AttributeWithOption)
	for _, variant := range product.Product.Variants {
		for _, attribute := range variant.Attributes {
			if mapAttributes[attribute.Name] == nil {
				mapAttributes[attribute.Name] = &AttributeWithOption{}
			}
			isExist := false
			for _, attributeOption := range mapAttributes[attribute.Name].Option {
				if attribute.Value == attributeOption {
					isExist = true
					break
				}
			}
			if !isExist {
				mapAttributes[attribute.Name].Option = append(mapAttributes[attribute.Name].Option, attribute.Value)
			}
		}
	}
	data := &DataProduct{
		Attributes: mapAttributes,
		Product:    product,
		IndexData:  indexData,
		Host:       c.Request().Host,
		Category:   category,
	}

	err = c.Render(http.StatusOK, "header.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "product.html", data)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "footer.html", data)
	if err != nil {
		return err
	}
	return nil
}

func GetCookieKey(c echo.Context) (string, error) {
	shopID, err := GetShopID(c)
	if err != nil {
		return "", err
	}
	cookie, err := c.Request().Cookie("session")
	if err != nil {
		id := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "session",
			Value:    id.String() + shopID.String(),
			Expires:  time.Now().Add(30 * 24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
		}
		c.SetCookie(cookie)
	}
	return cookie.Value, nil
}

func GetShopID(c echo.Context) (dot.ID, error) {
	subDoimain, err := getSubDomain(c)
	if err != nil {
		return 0, err
	}
	query := &webserver.GetShopIDBySiteSubdomainQuery{
		SiteSubDoimain: subDoimain,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return 0, c.Render(http.StatusOK, "404.html", nil)
	}
	return query.Result, nil
}
