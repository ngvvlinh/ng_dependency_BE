package summary

import (
	"context"

	"o.o/api/main/transaction"
	"o.o/api/summary"
	api "o.o/api/top/int/shop"
	"o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/service_classify"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
)

type SummaryService struct {
	session.Session

	SummaryQuery     summary.QueryBus
	SummaryOld       *logicsummary.Summary
	TransactionQuery transaction.QueryBus
}

func (s *SummaryService) Clone() api.SummaryService { res := *s; return &res }

func (s *SummaryService) SummarizeFulfillments(ctx context.Context, q *api.SummarizeFulfillmentsRequest) (*api.SummarizeFulfillmentsResponse, error) {
	query := &model.SummarizeFulfillmentsRequest{
		ShopID:   s.SS.Shop().ID,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
	}
	if err := s.SummaryOld.SummarizeFulfillments(ctx, query); err != nil {
		return nil, err
	}

	result := &api.SummarizeFulfillmentsResponse{
		Tables: convertpb.PbSummaryTables(query.Result.Tables),
	}
	return result, nil
}

func (s *SummaryService) SummarizeTopShip(ctx context.Context, q *api.SummarizeTopShipRequest) (*api.SummarizeTopShipResponse, error) {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return nil, err
	}
	query := &summary.SummaryTopShipQuery{
		ShopID:   s.SS.Shop().ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.SummarizeTopShipResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return result, nil
}

func (s *SummaryService) SummarizePOS(ctx context.Context, q *api.SummarizePOSRequest) (*api.SummarizePOSResponse, error) {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return nil, err
	}
	query := &summary.SummaryPOSQuery{
		ShopID:   s.SS.Shop().ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.SummarizePOSResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return result, nil
}

func (s *SummaryService) CalcBalanceUser(ctx context.Context, q *api.CalcBalanceUserRequest) (*api.CalcBalanceUserResponse, error) {
	shop := s.SS.Shop()
	// classify: default is shipping
	classify := q.ServiceClassify.Apply(service_classify.ServiceClassify(q.CreditClassify.Apply(credit_type.CreditClassifyShipping)))
	query := &transaction.GetBalanceUserQuery{
		UserID:   shop.OwnerID,
		Classify: classify,
	}
	if err := s.TransactionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	balance := query.Result
	var result = &api.CalcBalanceUserResponse{}
	switch classify {
	case service_classify.Shipping:
		result.ActualBalance = balance.ActualBalance
		result.AvailableBalance = balance.AvailableBalance
	case service_classify.Telecom:
		result.TelecomBalance = balance.ActualBalance
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Classify does not support")
	}

	return result, nil
}
