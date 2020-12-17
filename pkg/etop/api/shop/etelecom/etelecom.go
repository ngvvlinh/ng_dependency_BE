package etelecom

import (
	"context"

	"o.o/api/etelecom"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ExtensionService struct {
	session.Session

	EtelecomAggr  etelecom.CommandBus
	EtelecomQuery etelecom.QueryBus
}

func (s *ExtensionService) Clone() api.EtelecomService {
	res := *s
	return &res
}

func (s *ExtensionService) GetExtensions(ctx context.Context, r *api.GetExtensionsRequest) (*api.GetExtensionsResponse, error) {
	query := &etelecom.ListExtensionsQuery{
		AccountIDs: []dot.ID{s.SS.Shop().ID},
	}
	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := Convert_etelecom_Extensions_shop_Extensions(query.Result)

	// censor extension password
	for _, ext := range res {
		if ext.UserID != s.SS.User().ID {
			ext.ExtensionPassword = ""
		}
	}
	return &api.GetExtensionsResponse{Extensions: res}, nil
}

func (s *ExtensionService) CreateExtension(ctx context.Context, r *api.CreateExtensionRequest) (*api.Extension, error) {
	cmd := &etelecom.CreateExtensionCommand{
		UserID:    r.UserID,
		AccountID: s.SS.Shop().ID,
		HotlineID: r.HotlineID,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	var res api.Extension
	Convert_etelecom_Extension_shop_Extension(cmd.Result, &res)
	return &res, nil
}

func (s *ExtensionService) GetHotlines(ctx context.Context, _ *pbcm.Empty) (*api.GetHotLinesResponse, error) {
	// list all hotline builtin
	queryBuiltinHotlines := &etelecom.ListBuiltinHotlinesQuery{}
	if err := s.EtelecomQuery.Dispatch(ctx, queryBuiltinHotlines); err != nil {
		return nil, err
	}
	builtinHotlines := queryBuiltinHotlines.Result

	query := &etelecom.ListHotlinesQuery{
		OwnerID: s.SS.Shop().OwnerID,
	}
	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	hotlines := append(builtinHotlines, query.Result...)

	res := Convert_etelecom_Hotlines_shop_Hotlines(hotlines)
	return &api.GetHotLinesResponse{Hotlines: res}, nil
}

func (s *ExtensionService) GetCallLogs(ctx context.Context, r *api.GetCallLogsRequest) (*api.GetCallLogsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	query := &etelecom.ListCallLogsQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
	}
	if r.Filter != nil && (len(r.Filter.ExtensionIDs) > 0 || len(r.Filter.HotlineIDs) > 0) {
		query.HotlineIDs = r.Filter.HotlineIDs
		query.ExtensionIDs = r.Filter.ExtensionIDs
	}

	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := Convert_etelecom_CallLogs_shop_CallLogs(query.Result.CallLogs)
	return &api.GetCallLogsResponse{
		CallLogs: res,
		Paging:   cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}, nil
}
