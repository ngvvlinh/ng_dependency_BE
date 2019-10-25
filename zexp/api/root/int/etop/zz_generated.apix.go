// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package etop

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	etop "etop.vn/backend/pb/etop"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type UserServiceServer struct {
	UserAPI
}

func NewUserServiceServer(svc UserAPI) Server {
	return &UserServiceServer{
		UserAPI: svc,
	}
}

const UserServicePathPrefix = "/api/etop.User/"

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

func (s *UserServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/etop.User/ChangePassword":
		msg := new(etop.ChangePasswordRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.ChangePassword(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/ChangePasswordUsingToken":
		msg := new(etop.ChangePasswordUsingTokenRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.ChangePasswordUsingToken(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/Login":
		msg := new(etop.LoginRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.Login(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/Register":
		msg := new(etop.CreateUserRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.Register(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/ResetPassword":
		msg := new(etop.ResetPasswordRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.ResetPassword(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/SendEmailVerification":
		msg := new(etop.SendEmailVerificationRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.SendEmailVerification(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/SendPhoneVerification":
		msg := new(etop.SendPhoneVerificationRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.SendPhoneVerification(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/SendSTokenEmail":
		msg := new(etop.SendSTokenEmailRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.SendSTokenEmail(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/SessionInfo":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.SessionInfo(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/SwitchAccount":
		msg := new(etop.SwitchAccountRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.SwitchAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/UpdatePermission":
		msg := new(etop.UpdatePermissionRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.UpdatePermission(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/UpdateReferenceSale":
		msg := new(etop.UpdateReferenceSaleRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.UpdateReferenceSale(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/UpdateReferenceUser":
		msg := new(etop.UpdateReferenceUserRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.UpdateReferenceUser(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/UpgradeAccessToken":
		msg := new(etop.UpgradeAccessTokenRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.UpgradeAccessToken(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/VerifyEmailUsingToken":
		msg := new(etop.VerifyEmailUsingTokenRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.VerifyEmailUsingToken(ctx, msg)
		}
		return msg, fn, nil
	case "/api/etop.User/VerifyPhoneUsingToken":
		msg := new(etop.VerifyPhoneUsingTokenRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.VerifyPhoneUsingToken(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
