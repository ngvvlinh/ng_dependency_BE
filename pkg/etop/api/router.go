package api

import (
	service "o.o/api/top/int/etop"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/etop
// +gen:wrapper:package=etop

func NewEtopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewUserServiceServer(WrapUserService(userService)),
		service.NewAccountServiceServer(WrapAccountService(accountService)),
		service.NewLocationServiceServer(WrapLocationService(locationService)),
		service.NewBankServiceServer(WrapBankService(bankService)),
		service.NewAddressServiceServer(WrapAddressService(addressService)),
		service.NewAccountRelationshipServiceServer(WrapAccountRelationshipService(accountRelationshipService)),
		service.NewUserRelationshipServiceServer(WrapUserRelationshipService(userRelationshipService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
