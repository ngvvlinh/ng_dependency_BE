package shop

import (
	"context"

	"o.o/api/summary"
	"o.o/api/top/int/shop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
)

type SummaryService struct {
	SummaryQuery summary.QueryBus
	SummaryOld   *logicsummary.Summary
}

func (s *SummaryService) Clone() *SummaryService { res := *s; return &res }

func (s *SummaryService) SummarizeFulfillments(ctx context.Context, q *SummarizeFulfillmentsEndpoint) error {
	query := &model.SummarizeFulfillmentsRequest{
		ShopID:   q.Context.Shop.ID,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
	}
	if err := s.SummaryOld.SummarizeFulfillments(ctx, query); err != nil {
		return err
	}

	q.Result = &shop.SummarizeFulfillmentsResponse{
		Tables: convertpb.PbSummaryTables(query.Result.Tables),
	}
	return nil
}

func (s *SummaryService) SummarizeTopShip(ctx context.Context, q *SummarizeTopShipEndpoint) error {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return err
	}
	query := &summary.SummaryTopShipQuery{
		ShopID:   q.Context.Shop.ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.SummarizeTopShipResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return nil
}

func (s *SummaryService) SummarizePOS(ctx context.Context, q *SummarizePOSEndpoint) error {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return err
	}
	query := &summary.SummaryPOSQuery{
		ShopID:   q.Context.Shop.ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.SummarizePOSResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return nil
}

func (s *SummaryService) CalcBalanceShop(ctx context.Context, q *CalcBalanceShopEndpoint) error {
	query := &model.GetBalanceShopCommand{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CalcBalanceShopResponse{
		Balance: query.Result.Amount,
	}
	return nil
}
