// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package crm

import (
	"context"
	"time"

	api "etop.vn/api/top/services/crm"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	cmwrapper "etop.vn/backend/pkg/common/apifw/wrapper"
	bus "etop.vn/backend/pkg/common/bus"
	headers "etop.vn/backend/pkg/common/headers"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapCrmService(s *CrmService, secret string) api.CrmService {
	return wrapCrmService{s: s, secret: secret}
}

type wrapCrmService struct {
	s      *CrmService
	secret string
}

type RefreshFulfillmentFromCarrierEndpoint struct {
	*api.RefreshFulfillmentFromCarrierRequest
	Result  *cm.UpdatedResponse
	Context claims.EmptyClaim
}

func (s wrapCrmService) RefreshFulfillmentFromCarrier(ctx context.Context, req *api.RefreshFulfillmentFromCarrierRequest) (resp *cm.UpdatedResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Crm/RefreshFulfillmentFromCarrier"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &RefreshFulfillmentFromCarrierEndpoint{RefreshFulfillmentFromCarrierRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := headers.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.RefreshFulfillmentFromCarrier(ctx, query)
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

type SendNotificationEndpoint struct {
	*api.SendNotificationRequest
	Result  *cm.MessageResponse
	Context claims.EmptyClaim
}

func (s wrapCrmService) SendNotification(ctx context.Context, req *api.SendNotificationRequest) (resp *cm.MessageResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Crm/SendNotification"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &SendNotificationEndpoint{SendNotificationRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := headers.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.SendNotification(ctx, query)
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

func WrapMiscService(s *MiscService, secret string) api.MiscService {
	return wrapMiscService{s: s, secret: secret}
}

type wrapMiscService struct {
	s      *MiscService
	secret string
}

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context claims.EmptyClaim
}

func (s wrapMiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &VersionInfoEndpoint{Empty: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := headers.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.VersionInfo(ctx, query)
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

func WrapVhtService(s *VhtService) api.VhtService {
	return wrapVhtService{s: s}
}

type wrapVhtService struct {
	s *VhtService
}

type CreateOrUpdateCallHistoryByCallIDEndpoint struct {
	*api.VHTCallLog
	Result  *api.VHTCallLog
	Context claims.AdminClaim
}

func (s wrapVhtService) CreateOrUpdateCallHistoryByCallID(ctx context.Context, req *api.VHTCallLog) (resp *api.VHTCallLog, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vht/CreateOrUpdateCallHistoryByCallID"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireEtopAdmin: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateCallHistoryByCallIDEndpoint{VHTCallLog: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateOrUpdateCallHistoryByCallID(ctx, query)
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

type CreateOrUpdateCallHistoryBySDKCallIDEndpoint struct {
	*api.VHTCallLog
	Result  *api.VHTCallLog
	Context claims.AdminClaim
}

func (s wrapVhtService) CreateOrUpdateCallHistoryBySDKCallID(ctx context.Context, req *api.VHTCallLog) (resp *api.VHTCallLog, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vht/CreateOrUpdateCallHistoryBySDKCallID"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireEtopAdmin: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateCallHistoryBySDKCallIDEndpoint{VHTCallLog: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateOrUpdateCallHistoryBySDKCallID(ctx, query)
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

type GetCallHistoriesEndpoint struct {
	*api.GetCallHistoriesRequest
	Result  *api.GetCallHistoriesResponse
	Context claims.AdminClaim
}

func (s wrapVhtService) GetCallHistories(ctx context.Context, req *api.GetCallHistoriesRequest) (resp *api.GetCallHistoriesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vht/GetCallHistories"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireEtopAdmin: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetCallHistoriesEndpoint{GetCallHistoriesRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetCallHistories(ctx, query)
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

func WrapVtigerService(s *VtigerService) api.VtigerService {
	return wrapVtigerService{s: s}
}

type wrapVtigerService struct {
	s *VtigerService
}

type CreateOrUpdateContactEndpoint struct {
	*api.ContactRequest
	Result  *api.ContactResponse
	Context claims.ShopClaim
}

func (s wrapVtigerService) CreateOrUpdateContact(ctx context.Context, req *api.ContactRequest) (resp *api.ContactResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/CreateOrUpdateContact"
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
	query := &CreateOrUpdateContactEndpoint{ContactRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateOrUpdateContact(ctx, query)
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

type CreateOrUpdateLeadEndpoint struct {
	*api.LeadRequest
	Result  *api.LeadResponse
	Context claims.ShopClaim
}

func (s wrapVtigerService) CreateOrUpdateLead(ctx context.Context, req *api.LeadRequest) (resp *api.LeadResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/CreateOrUpdateLead"
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
	query := &CreateOrUpdateLeadEndpoint{LeadRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateOrUpdateLead(ctx, query)
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

type CreateTicketEndpoint struct {
	*api.CreateOrUpdateTicketRequest
	Result  *api.Ticket
	Context claims.ShopClaim
}

func (s wrapVtigerService) CreateTicket(ctx context.Context, req *api.CreateOrUpdateTicketRequest) (resp *api.Ticket, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/CreateTicket"
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
	query := &CreateTicketEndpoint{CreateOrUpdateTicketRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateTicket(ctx, query)
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

type GetCategoriesEndpoint struct {
	*cm.Empty
	Result  *api.GetCategoriesResponse
	Context claims.ShopClaim
}

func (s wrapVtigerService) GetCategories(ctx context.Context, req *cm.Empty) (resp *api.GetCategoriesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/GetCategories"
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
	query := &GetCategoriesEndpoint{Empty: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetCategories(ctx, query)
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

type GetContactsEndpoint struct {
	*api.GetContactsRequest
	Result  *api.GetContactsResponse
	Context claims.AdminClaim
}

func (s wrapVtigerService) GetContacts(ctx context.Context, req *api.GetContactsRequest) (resp *api.GetContactsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/GetContacts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireEtopAdmin: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetContactsEndpoint{GetContactsRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetContacts(ctx, query)
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

type GetTicketStatusCountEndpoint struct {
	*cm.Empty
	Result  *api.GetTicketStatusCountResponse
	Context claims.AdminClaim
}

func (s wrapVtigerService) GetTicketStatusCount(ctx context.Context, req *cm.Empty) (resp *api.GetTicketStatusCountResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/GetTicketStatusCount"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireEtopAdmin: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetTicketStatusCountEndpoint{Empty: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetTicketStatusCount(ctx, query)
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

type GetTicketsEndpoint struct {
	*api.GetTicketsRequest
	Result  *api.GetTicketsResponse
	Context claims.ShopClaim
}

func (s wrapVtigerService) GetTickets(ctx context.Context, req *api.GetTicketsRequest) (resp *api.GetTicketsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/GetTickets"
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
	query := &GetTicketsEndpoint{GetTicketsRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetTickets(ctx, query)
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

type UpdateTicketEndpoint struct {
	*api.CreateOrUpdateTicketRequest
	Result  *api.Ticket
	Context claims.ShopClaim
}

func (s wrapVtigerService) UpdateTicket(ctx context.Context, req *api.CreateOrUpdateTicketRequest) (resp *api.Ticket, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "crm.Vtiger/UpdateTicket"
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
	query := &UpdateTicketEndpoint{CreateOrUpdateTicketRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.UpdateTicket(ctx, query)
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
