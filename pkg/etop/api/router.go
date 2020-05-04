package api

import (
	"net/http"
	"strings"

	service "o.o/api/top/int/etop"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/etop
// +gen:wrapper:package=etop

func NewEtopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewUserServiceServer(WrapUserService(userService.Clone)),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
		service.NewLocationServiceServer(WrapLocationService(locationService.Clone)),
		service.NewBankServiceServer(WrapBankService(bankService.Clone)),
		service.NewAddressServiceServer(WrapAddressService(addressService.Clone)),
		service.NewAccountRelationshipServiceServer(WrapAccountRelationshipService(accountRelationshipService.Clone)),
		service.NewUserRelationshipServiceServer(WrapUserRelationshipService(userRelationshipService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}

	// proxy /api/root... to /api/etop
	for _, s := range servers {
		pathPrefix := strings.Replace(s.PathPrefix(), "/etop.", "/root.", 1)
		m.Handle(pathPrefix, proxy(s))
	}
}

func proxy(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newPath := strings.Replace(r.URL.Path, "/root.", "/etop.", 1)
		if newPath == r.URL.Path {
			next.ServeHTTP(w, r)
			return
		}
		r2 := *r
		u := *r.URL
		r2.URL = &u
		r2.URL.Path = newPath
		next.ServeHTTP(w, &r2)
	}
}
