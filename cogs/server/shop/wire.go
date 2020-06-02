package server_shop

import (
	"net/http"
	"time"

	"github.com/google/wire"

	"o.o/backend/pkg/common/apifw/httpx"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/logic/orders/imcsv"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
)

var WireSet = wire.NewSet(
	BuildImportHandler,
	BuildEventStreamHandler,
	BuildDownloadHandler,
)

type ImportHandler httpx.Server

func BuildImportHandler(
	orderImport *imcsv.Import,
	productImport *productimcsv.Import,
) ImportHandler {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))
	rt.Use(httpx.Auth(permission.Shop))

	rt.POST("/api/shop.Import/Orders", orderImport.HandleImportOrders)
	rt.POST("/api/shop.Import/Products", productImport.HandleShopImportProducts)
	rt.POST("/api/shop.Import/SampleProducts", productImport.HandleShopImportSampleProducts)
	return httpx.MakeServer("/api/shop.Import/", rt)
}

type EventStreamHandler httpx.Server

func BuildEventStreamHandler(
	eventStreamer *eventstream.EventStream,
) EventStreamHandler {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))
	rt.Use(httpx.Auth(permission.Shop))
	rt.GET("/api/event-stream", eventStreamer.HandleEventStream)

	s := headers.ForwardHeaders(rt, headers.Config{AllowQueryAuthorization: true})
	return httpx.MakeServer("/api/event-stream", s)
}

type DownloadHandler httpx.Server

func BuildDownloadHandler() DownloadHandler {
	mux := http.NewServeMux()

	// change path for clearing browser cache and still keep the old/dl
	// path for backward compatible
	mux.Handle("/dl/imports/shop_orders.v1.xlsx",
		cmservice.ServeAssetsByContentGenerator(
			cmservice.MIMEExcel,
			orderimcsv.AssetShopOrderPath,
			5*time.Minute,
			orderimcsv.GenerateImportFile,
		),
	)
	mux.Handle("/dl/imports/shop_orders.v1b.xlsx",
		cmservice.ServeAssetsByContentGenerator(
			cmservice.MIMEExcel,
			orderimcsv.AssetShopOrderPath,
			5*time.Minute,
			orderimcsv.GenerateImportFile,
		),
	)
	mux.Handle("/dl/imports/shop_products.v1.xlsx",
		cmservice.ServeAssetsByContentGenerator(
			cmservice.MIMEExcel,
			productimcsv.AssetShopProductPath,
			5*time.Minute,
			productimcsv.GenerateImportFile,
		),
	)
	mux.Handle("/dl/imports/shop_products.v1.simplified.xlsx",
		cmservice.ServeAssets(
			productimcsv.AssetShopProductSimplifiedPath,
			cmservice.MIMEExcel,
		),
	)
	return httpx.MakeServer("/dl/imports/", mux)
}
