package fabo

import (
	"context"
	"time"

	"o.o/api/fabo/summary"
	"o.o/api/top/int/fabo"
	cm "o.o/backend/pkg/common"
	convertpb2 "o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

const (
	// maximum days between dateFrom and dateTo
	MaxDaysBetween = 400 * 24 * time.Hour //days
)

type SummaryService struct {
	session.Session

	SummaryQuery summary.QueryBus
}

func (s *SummaryService) Clone() fabo.SummaryService { res := *s; return &res }

func (s *SummaryService) SummaryShop(ctx context.Context, req *fabo.SummaryShopRequest) (*fabo.SummaryShopResponse, error) {
	dateFrom, dateTo, err := cm.ParseDateFromTo(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	if err := validateDates(dateFrom, dateTo); err != nil {
		return nil, err
	}

	query := &summary.SummaryShopQuery{
		ShopID:   s.SS.Shop().ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err := s.SummaryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return &fabo.SummaryShopResponse{
		Tables: convertpb2.PbSummaryTablesNew(query.Result.ListTable),
	}, nil
}

func validateDates(dateFrom, dateTo time.Time) error {
	if dateTo.Before(dateFrom) {
		return cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}

	if dateTo.Sub(dateFrom) > MaxDaysBetween {
		return cm.Errorf(cm.InvalidArgument, nil, "maximum days between date_from and date_to is 400 days")
	}

	return nil
}
