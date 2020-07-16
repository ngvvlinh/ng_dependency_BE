package shop

import (
	"context"

	"o.o/api/summary"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
)

type SummaryService struct {
	session.Session

	SummaryQuery summary.QueryBus
	SummaryOld   *logicsummary.Summary
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

func (s *SummaryService) CalcBalanceUser(ctx context.Context, q *pbcm.Empty) (*api.CalcBalanceUserResponse, error) {
	shop := s.SS.Shop()
	queryActual := &model.GetActualUserBalanceCommand{
		UserID: shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, queryActual); err != nil {
		return nil, err
	}

	queryAvailable := &model.GetAvailableUserBalanceCommand{
		UserID: shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, queryAvailable); err != nil {
		return nil, err
	}

	result := &api.CalcBalanceUserResponse{
		AvailableBalance: queryAvailable.Result.Amount,
		ActualBalance:    queryActual.Result.Amount,
	}
	return result, nil
}
