package shop

import (
	"context"

	"o.o/api/main/ledgering"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/ledger_type"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type LedgerService struct {
	session.Session

	LedgerAggr  ledgering.CommandBus
	LedgerQuery ledgering.QueryBus
}

func (s *LedgerService) Clone() api.LedgerService { res := *s; return &res }

func (s *LedgerService) GetLedger(ctx context.Context, r *pbcm.IDRequest) (*api.Ledger, error) {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.LedgerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := convertpb.PbLedger(query.Result)
	return result, nil
}

func (s *LedgerService) GetLedgers(ctx context.Context, r *api.GetLedgersRequest) (*api.LedgersResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &ledgering.ListLedgersQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.LedgerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &api.LedgersResponse{
		Ledgers: convertpb.PbLedgers(query.Result.Ledgers),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *LedgerService) CreateLedger(ctx context.Context, r *api.CreateLedgerRequest) (*api.Ledger, error) {
	cmd := &ledgering.CreateLedgerCommand{
		ShopID:      s.SS.Shop().ID,
		Name:        r.Name,
		BankAccount: convertpb.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        r.Note,
		Type:        ledger_type.LedgerTypeBank,
		CreatedBy:   s.SS.Claim().UserID,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbLedger(cmd.Result)

	return result, nil
}

func (s *LedgerService) UpdateLedger(ctx context.Context, r *api.UpdateLedgerRequest) (*api.Ledger, error) {
	cmd := &ledgering.UpdateLedgerCommand{
		ID:          r.Id,
		ShopID:      s.SS.Shop().ID,
		Name:        r.Name,
		BankAccount: convertpb.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        r.Note,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := convertpb.PbLedger(cmd.Result)
	return result, nil
}

func (s *LedgerService) DeleteLedger(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &ledgering.DeleteLedgerCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
}
