package etelecom

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/etelecom/summary"
	api "o.o/api/top/int/shop"
	shoptypes "o.o/api/top/int/shop/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ExtensionService struct {
	session.Session

	EtelecomAggr  etelecom.CommandBus
	EtelecomQuery etelecom.QueryBus
	SummaryQuery  summary.QueryBus
}

func (s *ExtensionService) Clone() api.EtelecomService {
	res := *s
	return &res
}

func (s *ExtensionService) GetExtensions(ctx context.Context, r *shoptypes.GetExtensionsRequest) (*shoptypes.GetExtensionsResponse, error) {
	query := &etelecom.ListExtensionsQuery{
		AccountIDs: []dot.ID{s.SS.Shop().ID},
	}
	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := Convert_etelecom_Extensions_shoptypes_Extensions(query.Result)

	// censor extension password
	for _, ext := range res {
		if ext.UserID != s.SS.User().ID {
			ext.ExtensionPassword = ""
		}
	}
	return &shoptypes.GetExtensionsResponse{Extensions: res}, nil
}

func (s *ExtensionService) CreateExtension(ctx context.Context, r *shoptypes.CreateExtensionRequest) (*shoptypes.Extension, error) {
	cmd := &etelecom.CreateExtensionCommand{
		UserID:    r.UserID,
		AccountID: s.SS.Shop().ID,
		HotlineID: r.HotlineID,
		OwnerID:   s.SS.User().ID,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	var res shoptypes.Extension
	Convert_etelecom_Extension_shoptypes_Extension(cmd.Result, &res)
	return &res, nil
}

func (s *ExtensionService) GetHotlines(ctx context.Context, _ *pbcm.Empty) (*shoptypes.GetHotLinesResponse, error) {
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

	res := Convert_etelecom_Hotlines_shoptypes_Hotlines(hotlines)
	return &shoptypes.GetHotLinesResponse{Hotlines: res}, nil
}

func (s *ExtensionService) GetCallLogs(ctx context.Context, r *shoptypes.GetCallLogsRequest) (*shoptypes.GetCallLogsResponse, error) {
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
	res := Convert_etelecom_CallLogs_shoptypes_CallLogs(query.Result.CallLogs)
	return &shoptypes.GetCallLogsResponse{
		CallLogs: res,
		Paging:   cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *ExtensionService) SummaryEtelecom(
	ctx context.Context, req *api.SummaryEtelecomRequest,
) (*api.SummaryEtelecomResponse, error) {
	dateFrom, dateTo, err := cm.ParseDateFromTo(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}

	if dateTo.Before(dateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}

	query := &summary.SummaryQuery{
		ShopID:   s.SS.Shop().ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err := s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return &api.SummaryEtelecomResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}, nil
}
