package affiliate

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/affiliate"
)

// +gen:wrapper=etop.vn/backend/zexp/api/root/int/affiliate
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
