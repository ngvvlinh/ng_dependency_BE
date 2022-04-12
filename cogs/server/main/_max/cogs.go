package server_max

import (
	"context"
	"net/http"
	"o.o/backend/pkg/etop/apix/portsip_pbx"
	"strings"

	_main "o.o/backend/cogs/server/main"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
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
	hs = append(hs, affServer...)
	hs = append(hs, saffServer...)
	hs = httprpc.WithHooks(hs, ssHooks, logging)
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
	perms := acl.GetExtACL()
	ssExtHooks, err := session.NewHook(perms)
	if err != nil {
		return nil, err
	}

	hs = append(hs, xshopServers...)
	hs = append(hs, carrierServers...)
	hs = append(hs, partnerServers...)
	hs = append(hs, partnerImportServers...)
	hs = append(hs, vnpostServers...)
	hs = append(hs, vhtServers...)

	// check whitelist IP for partner APIs
	var whitelistIPHooks httprpc.HooksFunc = func() httprpc.Hooks {
		return httprpc.Hooks{
			RequestRouted: func(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
				var ss *session.Session
				var partner *identitymodel.Partner
				if _ss, ok := info.Inner.(session.Sessioner); ok {
					ss = _ss.GetSession()
					partner = ss.SS.Partner()
				}
				if partner != nil {
					permission := perms[info.Route]
					if err := checkWhitelistIPs(GetIP(info.HTTPRequest), partner.WhitelistIPs, permission.RequiredWhitelistIPs); err != nil {
						return nil, err
					}
				}

				return ctx, nil
			},
		}
	}

	hs = httprpc.WithHooks(hs, ssExtHooks, whitelistIPHooks, logging)
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

func BuildPortSipPBXHandler(
	portsipService portsip_pbx.PortsipService,
) _main.PortSipHandler {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))

	rt.GET("/portsip-pbx/v1/cdr", portsipService.GetCallLogs)
	return httpx.MakeServer("/portsip-pbx/v1/", rt)
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func checkWhitelistIPs(currentIP string, whitelistIPs []string, requiredWhitelistIPs bool) error {
	if whitelistIPs == nil || len(whitelistIPs) == 0 {
		if requiredWhitelistIPs {
			return cm.ErrPermissionDenied
		}
		return nil
	}

	for _, ip := range whitelistIPs {
		if strings.Contains(currentIP, ip) {
			return nil
		}
	}
	return cm.ErrPermissionDenied
}
