package affiliate

import (
	service "o.o/api/top/int/affiliate"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/affiliate
// +gen:wrapper:package=affiliate

func NewAffiliateServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
