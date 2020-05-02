package api

import (
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
}
