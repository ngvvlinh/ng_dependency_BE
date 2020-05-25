package shop

import (
	"context"

	"o.o/api/main/ledgering"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/ledger_type"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type LedgerService struct {
	LedgerAggr  ledgering.CommandBus
	LedgerQuery ledgering.QueryBus
}

func (s *LedgerService) Clone() *LedgerService { res := *s; return &res }

func (s *LedgerService) GetLedger(ctx context.Context, r *GetLedgerEndpoint) error {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := s.LedgerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = convertpb.PbLedger(query.Result)
	return nil
}

func (s *LedgerService) GetLedgers(ctx context.Context, r *GetLedgersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &ledgering.ListLedgersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.LedgerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &shop.LedgersResponse{
		Ledgers: convertpb.PbLedgers(query.Result.Ledgers),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *LedgerService) CreateLedger(ctx context.Context, r *CreateLedgerEndpoint) error {
	cmd := &ledgering.CreateLedgerCommand{
		ShopID:      r.Context.Shop.ID,
		Name:        r.Name,
		BankAccount: convertpb.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        r.Note,
		Type:        ledger_type.LedgerTypeBank,
		CreatedBy:   r.Context.UserID,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbLedger(cmd.Result)

	return nil
}

func (s *LedgerService) UpdateLedger(ctx context.Context, r *UpdateLedgerEndpoint) error {
	cmd := &ledgering.UpdateLedgerCommand{
		ID:          r.Id,
		ShopID:      r.Context.Shop.ID,
		Name:        r.Name,
		BankAccount: convertpb.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        r.Note,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = convertpb.PbLedger(cmd.Result)
	return nil
}

func (s *LedgerService) DeleteLedger(ctx context.Context, r *DeleteLedgerEndpoint) error {
	cmd := &ledgering.DeleteLedgerCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := s.LedgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}
