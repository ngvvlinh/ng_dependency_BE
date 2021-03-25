package _all

import (
	"context"
	"net/http"
	
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/integration/sms"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

func NewServers(
	miscService *apiroot.MiscService,
	userService *apiroot.UserService,
	accountService *apiroot.AccountService,
	locationService *apiroot.LocationService,
	bankService *apiroot.BankService,
	addressService *apiroot.AddressService,
	accountRelationshipService *apiroot.AccountRelationshipService,
	userRelationshipService *apiroot.UserRelationshipService,
	ticketService *TicketService,
	ecomService *apiroot.EcomService,

	rd redis.Store,
	_cfgEmail cc.EmailConfig,
	_cfgSMS sms.Config,

) (apiroot.Servers, func()) {
	apiroot.UserServiceImpl = userService // MUSTDO: remove this
	apiroot.EnabledEmail = _cfgEmail.Enabled
	apiroot.EnabledSMS = _cfgSMS.Enabled
	apiroot.CfgEmail = _cfgEmail
	apiroot.Idempgroup = idemp.NewRedisGroup(rd, apiroot.PrefixIdempUser, 0)
	if apiroot.EnabledEmail {
		if _, err := validate.ValidateStruct(apiroot.CfgEmail); err != nil {
			ll.Fatal("Can not validate config", l.Error(err))
		}
	}

	var cookieHooks httprpc.HooksFunc = func() httprpc.Hooks {
		return httprpc.Hooks{
			ResponsePrepared: func(ctx context.Context, info httprpc.HookInfo, respHeaders http.Header) (context.Context, error) {
				_ctx := bus.GetContext(ctx)
				if _ctx == nil {
					return ctx, nil
				}

				cookieData := _ctx.Value(headers.CookieKey{})
				if cookieData == nil {
					return ctx, nil
				}
				cookies, ok := cookieData.([]*http.Cookie)
				if !ok {
					return ctx, nil
				}
				for _, cookie := range cookies {
					if v := cookie.String(); v != "" {
						respHeaders.Add("Set-Cookie", v)
					}
				}
				return ctx, nil
			},
		}
	}
	
	// Get Origin URL for generate invitation url
	var accountRelationshipHooks httprpc.HooksFunc = func() httprpc.Hooks {
		return httprpc.Hooks{
			RequestReceived: func(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
				reqHeaders := info.HTTPRequest.Header
				origin := reqHeaders.Get("Origin")
				if origin == "" {
					protocol := reqHeaders.Get("X-Forwarded-Proto")
					host := reqHeaders.Get("X-Forwarded-Host")
					if protocol != "" && host != "" {
						origin = protocol + "://" + host
					}
				}
				if origin != "" {
					ctx = context.WithValue(ctx, headers.OriginKey{}, origin)
				}
				return ctx, nil
			},
		}
	}

	servers := httprpc.MustNewServers(
		ticketService.Clone,
		accountService.Clone,
		addressService.Clone,
		bankService.Clone,
		ecomService.Clone,
		locationService.Clone,
		miscService.Clone,
		userRelationshipService.Clone,
	)
	servers = append(servers,
		httprpc.MustNewServer(userService.Clone, cookieHooks),
		httprpc.MustNewServer(accountRelationshipService.Clone, accountRelationshipHooks),
	)

	var result []httprpc.Server
	result = append(result, servers...)

	// proxy /api/root... to /api/etop
	result = apiroot.ProxyEtop(result)
	return result, apiroot.Idempgroup.Shutdown
}
