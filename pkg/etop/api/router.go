package api

import (
	service "etop.vn/api/root/int/etop"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/root/int/etop
// +gen:wrapper:package=etop

func NewEtopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewUserServiceServer(WrapUserService(userService)),
		service.NewAccountServiceServer(WrapAccountService(accountService)),
		service.NewRelationshipServiceServer(WrapRelationshipService(relationshipService)),
		service.NewLocationServiceServer(WrapLocationService(locationService)),
		service.NewBankServiceServer(WrapBankService(bankService)),
		service.NewAddressServiceServer(WrapAddressService(addressService)),
		service.NewInvitationServiceServer(WrapInvitationService(invitationService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
