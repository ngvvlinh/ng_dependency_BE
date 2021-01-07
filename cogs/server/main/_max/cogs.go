package server_max

import (
	_main "o.o/backend/cogs/server/main"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/api/admin"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/integration"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/authx"
	"o.o/backend/pkg/etop/apix/mc/vht"
	"o.o/backend/pkg/etop/apix/mc/vnp"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/middlewares"
	serviceaffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
)

func BuildIntHandlers(
	rootServers apiroot.Servers,
	shopServers shop.Servers,
	adminServers admin.Servers,
	sadminServers sadmin.Servers,
	integrationServers integration.Servers,
	affServer affapi.Servers,
	saffServer serviceaffapi.Servers,
	c *captcha.Captcha,
) (hs _main.IntHandlers, _ error) {
	logging := middlewares.NewLogging()
	ssHooks, err := session.NewHook(acl.GetACL(), session.OptCaptcha(c))
	if err != nil {
		return nil, err
	}

	hs = append(hs, rootServers...)
	hs = append(hs, shopServers...)
	hs = append(hs, adminServers...)
	hs = append(hs, sadminServers...)
	hs = append(hs, integrationServers...)
	hs = httprpc.WithHooks(hs, ssHooks, logging)

	hs = append(hs, affServer...)
	hs = append(hs, saffServer...)
	return
}

func BuildExtHandlers(
	partnerServers partner.Servers,
	xshopServers xshop.Servers,
	carrierServers partnercarrier.Servers,
	partnerImportServers partnerimport.Servers,
	vnpostServers vnp.Servers,
	vhtServers vht.Servers,
) (hs _main.ExtHandlers, _ error) {
	logging := middlewares.NewLogging()
	ssExtHooks, err := session.NewHook(acl.GetExtACL())
	if err != nil {
		return nil, err
	}

	hs = append(hs, xshopServers...)
	hs = append(hs, carrierServers...)
	hs = append(hs, partnerServers...)
	hs = append(hs, partnerImportServers...)
	hs = append(hs, vnpostServers...)
	hs = append(hs, vhtServers...)
	hs = httprpc.WithHooks(hs, ssExtHooks, logging)
	return
}

func BuildAuthxHandler(
	authxService authx.AuthxService,
) _main.AuthxHandler {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))

	rt.POST("/v1/authx/AuthUser", authxService.AuthUser)
	return httpx.MakeServer("/v1/authx/", rt)
}
