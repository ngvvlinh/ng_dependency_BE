// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package etop

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	"etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type AccountServiceServer struct {
	inner AccountService
}

func NewAccountServiceServer(svc AccountService) Server {
	return &AccountServiceServer{
		inner: svc,
	}
}

const AccountServicePathPrefix = "/etop.Account/"

func (s *AccountServiceServer) PathPrefix() string {
	return AccountServicePathPrefix
}

func (s *AccountServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *AccountServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Account/GetPublicPartnerInfo":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetPublicPartnerInfo(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Account/GetPublicPartners":
		msg := &common.IDsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetPublicPartners(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Account/UpdateURLSlug":
		msg := &etop.UpdateURLSlugRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateURLSlug(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type AddressServiceServer struct {
	inner AddressService
}

func NewAddressServiceServer(svc AddressService) Server {
	return &AddressServiceServer{
		inner: svc,
	}
}

const AddressServicePathPrefix = "/etop.Address/"

func (s *AddressServiceServer) PathPrefix() string {
	return AddressServicePathPrefix
}

func (s *AddressServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *AddressServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Address/CreateAddress":
		msg := &etop.CreateAddressRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateAddress(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Address/GetAddresses":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetAddresses(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Address/RemoveAddress":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.RemoveAddress(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Address/UpdateAddress":
		msg := &etop.UpdateAddressRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateAddress(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type BankServiceServer struct {
	inner BankService
}

func NewBankServiceServer(svc BankService) Server {
	return &BankServiceServer{
		inner: svc,
	}
}

const BankServicePathPrefix = "/etop.Bank/"

func (s *BankServiceServer) PathPrefix() string {
	return BankServicePathPrefix
}

func (s *BankServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *BankServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Bank/GetBanks":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetBanks(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Bank/GetBranchesByBankProvince":
		msg := &etop.GetBranchesByBankProvinceResquest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetBranchesByBankProvince(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Bank/GetProvincesByBank":
		msg := &etop.GetProvincesByBankResquest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetProvincesByBank(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type InvitationServiceServer struct {
	inner InvitationService
}

func NewInvitationServiceServer(svc InvitationService) Server {
	return &InvitationServiceServer{
		inner: svc,
	}
}

const InvitationServicePathPrefix = "/etop.Invitation/"

func (s *InvitationServiceServer) PathPrefix() string {
	return InvitationServicePathPrefix
}

func (s *InvitationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *InvitationServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Invitation/AcceptInvitation":
		msg := &etop.AcceptInvitationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.AcceptInvitation(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Invitation/GetInvitationByToken":
		msg := &etop.GetInvitationByTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetInvitationByToken(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Invitation/GetInvitations":
		msg := &etop.GetInvitationsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetInvitations(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Invitation/RejectInvitation":
		msg := &etop.RejectInvitationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.RejectInvitation(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type LocationServiceServer struct {
	inner LocationService
}

func NewLocationServiceServer(svc LocationService) Server {
	return &LocationServiceServer{
		inner: svc,
	}
}

const LocationServicePathPrefix = "/etop.Location/"

func (s *LocationServiceServer) PathPrefix() string {
	return LocationServicePathPrefix
}

func (s *LocationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *LocationServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Location/GetDistricts":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetDistricts(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Location/GetDistrictsByProvince":
		msg := &etop.GetDistrictsByProvinceRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetDistrictsByProvince(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Location/GetProvinces":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetProvinces(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Location/GetWards":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetWards(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Location/GetWardsByDistrict":
		msg := &etop.GetWardsByDistrictRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetWardsByDistrict(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Location/ParseLocation":
		msg := &etop.ParseLocationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ParseLocation(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	inner MiscService
}

func NewMiscServiceServer(svc MiscService) Server {
	return &MiscServiceServer{
		inner: svc,
	}
}

const MiscServicePathPrefix = "/etop.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *MiscServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type RelationshipServiceServer struct {
	inner RelationshipService
}

func NewRelationshipServiceServer(svc RelationshipService) Server {
	return &RelationshipServiceServer{
		inner: svc,
	}
}

const RelationshipServicePathPrefix = "/etop.Relationship/"

func (s *RelationshipServiceServer) PathPrefix() string {
	return RelationshipServicePathPrefix
}

func (s *RelationshipServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *RelationshipServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.Relationship/AnswerInvitation":
		msg := &etop.AnswerInvitationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.AnswerInvitation(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Relationship/GetUsersInCurrentAccounts":
		msg := &etop.GetUsersInCurrentAccountsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetUsersInCurrentAccounts(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Relationship/InviteUserToAccount":
		msg := &etop.InviteUserToAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.InviteUserToAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Relationship/LeaveAccount":
		msg := &etop.LeaveAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.LeaveAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.Relationship/RemoveUserFromCurrentAccount":
		msg := &etop.RemoveUserFromCurrentAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.RemoveUserFromCurrentAccount(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	inner UserService
}

func NewUserServiceServer(svc UserService) Server {
	return &UserServiceServer{
		inner: svc,
	}
}

const UserServicePathPrefix = "/etop.User/"

func (s *UserServiceServer) PathPrefix() string {
	return UserServicePathPrefix
}

func (s *UserServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *UserServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etop.User/ChangePassword":
		msg := &etop.ChangePasswordRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ChangePassword(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/ChangePasswordUsingToken":
		msg := &etop.ChangePasswordUsingTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ChangePasswordUsingToken(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/Login":
		msg := &etop.LoginRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.Login(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/Register":
		msg := &etop.CreateUserRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.Register(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/ResetPassword":
		msg := &etop.ResetPasswordRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ResetPassword(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/SendEmailVerification":
		msg := &etop.SendEmailVerificationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SendEmailVerification(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/SendPhoneVerification":
		msg := &etop.SendPhoneVerificationRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SendPhoneVerification(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/SendSTokenEmail":
		msg := &etop.SendSTokenEmailRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SendSTokenEmail(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/SessionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SessionInfo(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/SwitchAccount":
		msg := &etop.SwitchAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SwitchAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/UpdatePermission":
		msg := &etop.UpdatePermissionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdatePermission(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/UpdateReferenceSale":
		msg := &etop.UpdateReferenceSaleRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateReferenceSale(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/UpdateReferenceUser":
		msg := &etop.UpdateReferenceUserRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateReferenceUser(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/UpgradeAccessToken":
		msg := &etop.UpgradeAccessTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpgradeAccessToken(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/VerifyEmailUsingToken":
		msg := &etop.VerifyEmailUsingTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.VerifyEmailUsingToken(ctx, msg)
		}
		return msg, fn, nil
	case "/etop.User/VerifyPhoneUsingToken":
		msg := &etop.VerifyPhoneUsingTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.VerifyPhoneUsingToken(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
