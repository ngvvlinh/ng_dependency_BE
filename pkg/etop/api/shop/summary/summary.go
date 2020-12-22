package summary

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/summary"
	api "o.o/api/top/int/shop"
	"o.o/api/top/types/etc/credit_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	logicsummary "o.o/backend/pkg/etop/logic/summary"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
)

type SummaryService struct {
	session.Session

	SummaryQuery summary.QueryBus
	SummaryOld   *logicsummary.Summary
	CreditQuery  credit.QueryBus

	MoneyTxStore sqlstore.MoneyTxStoreInterface
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

	// creditClassify: default is shipping
	creditClassify := q.CreditClassify
	result := &api.CalcBalanceUserResponse{}
	if !creditClassify.Valid || creditClassify.Enum == credit_type.CreditClassifyShipping {
		queryActual := &model.GetActualUserBalanceCommand{
			UserID: shop.OwnerID,
		}
		if err := s.MoneyTxStore.CalcActualUserBalance(ctx, queryActual); err != nil {
			return nil, err
		}

		queryAvailable := &model.GetAvailableUserBalanceCommand{
			UserID: shop.OwnerID,
		}
		if err := s.MoneyTxStore.CalcAvailableUserBalance(ctx, queryAvailable); err != nil {
			return nil, err
		}
		result.AvailableBalance = queryAvailable.Result.Amount
		result.ActualBalance = queryActual.Result.Amount
	}

	if !creditClassify.Valid || creditClassify.Enum == credit_type.CreditClassifyTelecom {
		query := &credit.GetTelecomUserBalanceQuery{
			UserID: shop.OwnerID,
		}
		if err := s.CreditQuery.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		result.TelecomBalance = query.Result
	}

	return result, nil
}
