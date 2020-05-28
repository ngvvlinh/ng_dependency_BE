package imports

import (
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
)

type ShopImport struct {
	*httpx.Router
}

func New(
	orderImport *imcsv.Import,
	productImport *productimcsv.Import,
) ShopImport {
	rt := httpx.New()
	rt.Use(httpx.Auth(permission.Shop))

	rt.POST("/api/shop.Import/Orders", orderImport.HandleImportOrders)
	rt.POST("/api/shop.Import/Products", productImport.HandleShopImportProducts)
	rt.POST("/api/shop.Import/SampleProducts", productImport.HandleShopImportSampleProducts)
	return ShopImport{rt}
}
