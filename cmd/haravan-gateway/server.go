package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"

	"etop.vn/backend/com/external/haravan/gateway/aggregate"
	haravanserver "etop.vn/backend/com/external/haravan/gateway/server"
	haravanidentity "etop.vn/backend/com/external/haravan/identity"
	catalogquery "etop.vn/backend/com/main/catalog/query"
	"etop.vn/backend/com/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/metrics"
	orderS "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/common/l"
)

func startServers() *http.Server {
	identityQuery := identity.NewQueryService(db).MessageBus()
	haravanIdentityQuery := haravanidentity.NewQueryService(db).MessageBus()
	shippingManager := shipping_provider.NewCtrl(locationBus, ghnCarrier, ghtkCarrier, vtpostCarrier)
	haravan := aggregate.NewAggregate(db, shippingManager, locationBus, identityQuery).MessageBus()

	catalogQueryService := catalogquery.New(db).MessageBus()
	orderS.Init(shippingManager, catalogQueryService)
	haravanServer := haravanserver.New(haravan, haravanIdentityQuery)

	mux := http.NewServeMux()
	rt := httpx.New()
	mux.Handle("/haravan/", rt)

	rt.Use(httpx.RecoverAndLog(bot, false))
	rt.Use(authMiddleware)

	buildRoute := haravanidentity.BuildGatewayRoute
	rt.GET("/haravan/gateway/__test", test)
	rt.POST(buildRoute(haravanidentity.PathGetShippingRates), haravanServer.GetShippingRates)
	rt.POST(buildRoute(haravanidentity.PathCreateOrder), haravanServer.CreateOrder)
	rt.POST(buildRoute(haravanidentity.PathGetOrder), haravanServer.GetOrder)
	rt.DELETE(buildRoute(haravanidentity.PathCancelOrder), haravanServer.CancelOrder)

	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func authMiddleware(next httpx.Handler) httpx.Handler {
	return func(c *httpx.Context) error {
		h := c.Req.Header
		haravanHMAC := h.Get("X-Haravan-Hmac-SHA256")
		if haravanHMAC == "" {
			return cm.Errorf(cm.Unauthenticated, nil, "Xác thực không hợp lệ. Vui lòng kiểm tra lại.")
		}

		body, err := ioutil.ReadAll(c.Req.Body)
		if err != nil {
			return cm.Error(cm.ExternalServiceError, "failed to read request body", err)
		}

		// Restore the io.ReadCloser to its original state
		c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		hash := generateAuthorizeCode(string(body), cfg.Haravan.Secret)
		ll.Info("hash ::", l.String("author", hash))
		if hash != haravanHMAC {
			return cm.Errorf(cm.Unauthenticated, nil, "Xác thực không hợp lệ. Vui lòng kiểm tra lại.")
		}

		return next(c)
	}
}

func generateAuthorizeCode(data string, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(hash, data)
	if err != nil {
		panic(err)
	}
	macSum := hash.Sum(nil)
	dd := base64.StdEncoding.EncodeToString(macSum)
	return dd
}

func test(c *httpx.Context) error {
	c.SetResult(map[string]string{
		"code": "ok",
	})
	return nil
}
