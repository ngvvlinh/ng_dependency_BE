package affiliate

import (
	service "etop.vn/api/top/int/affiliate"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/int/affiliate
// +gen:wrapper:package=affiliate

func NewAffiliateServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewAccountServiceServer(WrapAccountService(accountService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
