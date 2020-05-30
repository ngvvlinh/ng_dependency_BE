package api

import (
	"context"
	"net/http"
	"strings"

	service "o.o/api/top/int/etop"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/integration/sms"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

// +gen:wrapper=o.o/api/top/int/etop
// +gen:wrapper:package=etop

type Servers []httprpc.Server

func NewServers(
	miscService *MiscService,
	userService *UserService,
	accountService *AccountService,
	locationService *LocationService,
	bankService *BankService,
	addressService *AddressService,
	accountRelationshipService *AccountRelationshipService,
	userRelationshipService *UserRelationshipService,
	ecomService *EcomService,

	rd redis.Store,
	_cfgEmail EmailConfig,
	_cfgSMS sms.Config,

) (Servers, func()) {
	UserServiceImpl = userService // MUSTDO: remove this
	enabledEmail = _cfgEmail.Enabled
	enabledSMS = _cfgSMS.Enabled
	cfgEmail = _cfgEmail
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempUser, 0)
	if enabledEmail {
		if _, err := validate.ValidateStruct(cfgEmail); err != nil {
			ll.Fatal("Can not validate config", l.Error(err))
		}
	}

	var cookieHooks httprpc.HooksFunc = func() httprpc.Hooks {
		return httprpc.Hooks{
			BeforeResponse: func(ctx context.Context, info httprpc.HookInfo, respHeaders http.Header) (context.Context, error) {
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

	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewUserServiceServer(WrapUserService(userService.Clone), cookieHooks),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
		service.NewLocationServiceServer(WrapLocationService(locationService.Clone)),
		service.NewBankServiceServer(WrapBankService(bankService.Clone)),
		service.NewAddressServiceServer(WrapAddressService(addressService.Clone)),
		service.NewAccountRelationshipServiceServer(WrapAccountRelationshipService(accountRelationshipService.Clone)),
		service.NewUserRelationshipServiceServer(WrapUserRelationshipService(userRelationshipService.Clone)),
		service.NewEcomServiceServer(WrapEcomService(ecomService.Clone)),
	}

	var result []httprpc.Server
	result = append(result, servers...)

	// proxy /api/root... to /api/etop
	for _, s := range servers {
		pathPrefix := strings.Replace(s.PathPrefix(), "/etop.", "/root.", 1)
		prx := &Proxy{pathPrefix, s}
		result = append(result, prx)
	}
	return result, idempgroup.Shutdown
}

var _ httprpc.Server = &Proxy{}

type Proxy struct {
	pathPrefix string
	next       httprpc.Server
}

func (p *Proxy) PathPrefix() string {
	return p.pathPrefix
}

func (p *Proxy) WithHooks(builder httprpc.HooksBuilder) httprpc.Server {
	p.next = p.next.WithHooks(builder)
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newPath := strings.Replace(r.URL.Path, "/root.", "/etop.", 1)
	if newPath == r.URL.Path {
		p.next.ServeHTTP(w, r)
		return
	}
	r2 := *r
	u := *r.URL
	r2.URL = &u
	r2.URL.Path = newPath
	p.next.ServeHTTP(w, &r2)
}
