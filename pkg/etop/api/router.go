package api

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/etop"
)

// +gen:wrapper=etop.vn/backend/pb/etop
// +gen:wrapper:package=etop

func NewEtopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewUserServiceServer(NewUserService(userService)),
		service.NewAccountServiceServer(NewAccountService(accountService)),
		service.NewRelationshipServiceServer(NewRelationshipService(relationshipService)),
		service.NewLocationServiceServer(NewLocationService(locationService)),
		service.NewBankServiceServer(NewBankService(bankService)),
		service.NewAddressServiceServer(NewAddressService(addressService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
