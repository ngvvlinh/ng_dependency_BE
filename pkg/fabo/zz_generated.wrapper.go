// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package fabo

import (
	"context"
	"time"

	api "etop.vn/api/top/int/fabo"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	cmwrapper "etop.vn/backend/pkg/common/apifw/wrapper"
	bus "etop.vn/backend/pkg/common/bus"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapPageService(s *PageService) api.PageService {
	return wrapPageService{s: s}
}

type wrapPageService struct {
	s *PageService
}

type ConnectPagesEndpoint struct {
	*api.ConnectPagesRequest
	Result  *api.ConnectPagesResponse
	Context claims.ShopClaim
}

func (s wrapPageService) ConnectPages(ctx context.Context, req *api.ConnectPagesRequest) (resp *api.ConnectPagesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "fabo.Page/ConnectPages"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ConnectPagesEndpoint{ConnectPagesRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.ConnectPages(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type ListPagesEndpoint struct {
	*api.ListPagesRequest
	Result  *api.ListPagesResponse
	Context claims.ShopClaim
}

func (s wrapPageService) ListPages(ctx context.Context, req *api.ListPagesRequest) (resp *api.ListPagesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "fabo.Page/ListPages"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ListPagesEndpoint{ListPagesRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	// Include Fabo's information
	if query.Context.Shop != nil {
		getFaboInfoQuery := &middleware.GetFaboInfoQuery{
			ShopID: query.Context.Shop.ID,
			UserID: query.Context.UserID,
		}
		faboInfo, err := middleware.GetFaboInfo(ctx, getFaboInfoQuery)
		if err != nil {
			return nil, err
		}
		query.Context.FaboInfo = faboInfo
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.ListPages(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type RemovePagesEndpoint struct {
	*api.RemovePagesRequest
	Result  *cm.Empty
	Context claims.ShopClaim
}

func (s wrapPageService) RemovePages(ctx context.Context, req *api.RemovePagesRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "fabo.Page/RemovePages"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &RemovePagesEndpoint{RemovePagesRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	// Include Fabo's information
	if query.Context.Shop != nil {
		getFaboInfoQuery := &middleware.GetFaboInfoQuery{
			ShopID: query.Context.Shop.ID,
			UserID: query.Context.UserID,
		}
		faboInfo, err := middleware.GetFaboInfo(ctx, getFaboInfoQuery)
		if err != nil {
			return nil, err
		}
		query.Context.FaboInfo = faboInfo
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.RemovePages(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}
