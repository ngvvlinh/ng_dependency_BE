// package sadmin generated by wrapper_gen. DO NOT EDIT.
package sadminW

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"

	cm "etop.vn/backend/pb/common"
	etop "etop.vn/backend/pb/etop"
	sadmin "etop.vn/backend/pb/etop/sadmin"
	common "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/metrics"
	cmWrapper "etop.vn/backend/pkg/common/wrapper"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/middleware"
)

var ll = l.New()

var Client Sadmin

type Sadmin interface {
	sadmin.MiscService
	sadmin.UserService
}

type SadminClient struct {
	_MiscService sadmin.MiscService
	_UserService sadmin.UserService
}

func NewSadminClient(addr string, client *http.Client) Sadmin {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	addr = "http://" + addr
	return &SadminClient{
		_MiscService: sadmin.NewMiscServiceProtobufClient(addr, client),
		_UserService: sadmin.NewUserServiceProtobufClient(addr, client),
	}
}

func ConnectSadminService(addr string, client *http.Client) error {
	Client = NewSadminClient(addr, client)
	bus.AddHandler("client", func(ctx context.Context, q *VersionInfoEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *CreateUserEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *LoginAsAccountEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *ResetPasswordEndpoint) error { panic("Unexpected") })
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := Client.VersionInfo(ctx, &cm.Empty{})
	if err == nil {
		ll.S.Infof("Connected to SadminService at %v", addr)
	}
	return err
}

func MustConnectSadminService(addr string, client *http.Client) {
	err := ConnectSadminService(addr, client)
	if err != nil {
		ll.Fatal("Unable to connect Sadmin", l.Error(err))
	}
}

type (
	EmptyClaim   = claims.EmptyClaim
	UserClaim    = claims.UserClaim
	AdminClaim   = claims.AdminClaim
	PartnerClaim = claims.PartnerClaim
	ShopClaim    = claims.ShopClaim
)

func (c *SadminClient) VersionInfo(ctx context.Context, in *cm.Empty) (*cm.VersionInfoResponse, error) {
	resp, err := c._MiscService.VersionInfo(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *SadminClient) CreateUser(ctx context.Context, in *sadmin.SAdminCreateUserRequest) (*etop.RegisterResponse, error) {
	resp, err := c._UserService.CreateUser(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *SadminClient) LoginAsAccount(ctx context.Context, in *sadmin.LoginAsAccountRequest) (*etop.LoginResponse, error) {
	resp, err := c._UserService.LoginAsAccount(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *SadminClient) ResetPassword(ctx context.Context, in *sadmin.SAdminResetPasswordRequest) (*cm.Empty, error) {
	resp, err := c._UserService.ResetPassword(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}

type Muxer interface {
	Handle(string, http.Handler)
}

func NewSadminServer(mux Muxer, hooks *twirp.ServerHooks) {
	bus.Expect(&VersionInfoEndpoint{})
	bus.Expect(&CreateUserEndpoint{})
	bus.Expect(&LoginAsAccountEndpoint{})
	bus.Expect(&ResetPasswordEndpoint{})
	mux.Handle(sadmin.MiscServicePathPrefix, sadmin.NewMiscServiceServer(MiscService{}, hooks))
	mux.Handle(sadmin.UserServicePathPrefix, sadmin.NewUserServiceServer(UserService{}, hooks))
}

type SadminImpl struct {
	MiscService
	UserService
}

func NewSadmin() Sadmin {
	return SadminImpl{}
}

type MiscService struct{}

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context EmptyClaim
}

func (s MiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmWrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &VersionInfoEndpoint{Empty: req}
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type UserService struct{}

type CreateUserEndpoint struct {
	*sadmin.SAdminCreateUserRequest
	Result  *etop.RegisterResponse
	Context EmptyClaim
}

func (s UserService) CreateUser(ctx context.Context, req *sadmin.SAdminCreateUserRequest) (resp *etop.RegisterResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/CreateUser"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmWrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateUserEndpoint{SAdminCreateUserRequest: req}
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type LoginAsAccountEndpoint struct {
	*sadmin.LoginAsAccountRequest
	Result  *etop.LoginResponse
	Context EmptyClaim
}

func (s UserService) LoginAsAccount(ctx context.Context, req *sadmin.LoginAsAccountRequest) (resp *etop.LoginResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/LoginAsAccount"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmWrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &LoginAsAccountEndpoint{LoginAsAccountRequest: req}
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type ResetPasswordEndpoint struct {
	*sadmin.SAdminResetPasswordRequest
	Result  *cm.Empty
	Context EmptyClaim
}

func (s UserService) ResetPassword(ctx context.Context, req *sadmin.SAdminResetPasswordRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/ResetPassword"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmWrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ResetPasswordEndpoint{SAdminResetPasswordRequest: req}
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}
