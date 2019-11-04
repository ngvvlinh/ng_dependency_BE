// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package crm

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	crm "etop.vn/backend/pb/services/crm"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type CrmServiceServer struct {
	CrmAPI
}

func NewCrmServiceServer(svc CrmAPI) Server {
	return &CrmServiceServer{
		CrmAPI: svc,
	}
}

const CrmServicePathPrefix = "/api/crm.Crm/"

func (s *CrmServiceServer) PathPrefix() string {
	return CrmServicePathPrefix
}

func (s *CrmServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CrmServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/crm.Crm/RefreshFulfillmentFromCarrier":
		msg := new(crm.RefreshFulfillmentFromCarrierRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CrmAPI.RefreshFulfillmentFromCarrier(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Crm/SendNotification":
		msg := new(crm.SendNotificationRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CrmAPI.SendNotification(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	MiscAPI
}

func NewMiscServiceServer(svc MiscAPI) Server {
	return &MiscServiceServer{
		MiscAPI: svc,
	}
}

const MiscServicePathPrefix = "/api/crm.Misc/"

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

func (s *MiscServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/crm.Misc/VersionInfo":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MiscAPI.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type VhtServiceServer struct {
	VhtAPI
}

func NewVhtServiceServer(svc VhtAPI) Server {
	return &VhtServiceServer{
		VhtAPI: svc,
	}
}

const VhtServicePathPrefix = "/api/crm.Vht/"

func (s *VhtServiceServer) PathPrefix() string {
	return VhtServicePathPrefix
}

func (s *VhtServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *VhtServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/crm.Vht/CreateOrUpdateCallHistoryByCallID":
		msg := new(crm.VHTCallLog)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VhtAPI.CreateOrUpdateCallHistoryByCallID(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vht/CreateOrUpdateCallHistoryBySDKCallID":
		msg := new(crm.VHTCallLog)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VhtAPI.CreateOrUpdateCallHistoryBySDKCallID(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vht/GetCallHistories":
		msg := new(crm.GetCallHistoriesRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VhtAPI.GetCallHistories(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type VtigerServiceServer struct {
	VtigerAPI
}

func NewVtigerServiceServer(svc VtigerAPI) Server {
	return &VtigerServiceServer{
		VtigerAPI: svc,
	}
}

const VtigerServicePathPrefix = "/api/crm.Vtiger/"

func (s *VtigerServiceServer) PathPrefix() string {
	return VtigerServicePathPrefix
}

func (s *VtigerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *VtigerServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/crm.Vtiger/CountTicketByStatus":
		msg := new(crm.CountTicketByStatusRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.CountTicketByStatus(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/CreateOrUpdateContact":
		msg := new(crm.ContactRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.CreateOrUpdateContact(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/CreateOrUpdateLead":
		msg := new(crm.LeadRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.CreateOrUpdateLead(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/CreateTicket":
		msg := new(crm.CreateOrUpdateTicketRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.CreateTicket(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/GetCategories":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.GetCategories(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/GetContacts":
		msg := new(crm.GetContactsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.GetContacts(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/GetStatus":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.GetStatus(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/GetTicketStatusCount":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.GetTicketStatusCount(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/GetTickets":
		msg := new(crm.GetTicketsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.GetTickets(ctx, msg)
		}
		return msg, fn, nil
	case "/api/crm.Vtiger/UpdateTicket":
		msg := new(crm.CreateOrUpdateTicketRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.VtigerAPI.UpdateTicket(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
