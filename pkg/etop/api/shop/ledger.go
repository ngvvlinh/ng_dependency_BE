package shop

import (
	"context"

	"etop.vn/api/main/ledgering"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		ledgerService.GetLedger,
		ledgerService.GetLedgers,
		ledgerService.CreateLedger,
		ledgerService.UpdateLedger,
		ledgerService.DeleteLedger)
}

func (s LedgerService) GetLedger(ctx context.Context, r *GetLedgerEndpoint) error {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := ledgerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = pbshop.PbLedger(query.Result)
	return nil
}

func (s LedgerService) GetLedgers(ctx context.Context, r *GetLedgersEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &ledgering.ListLedgersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := ledgerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &pbshop.LedgersResponse{
		Ledgers: pbshop.PbLedgers(query.Result.Ledgers),
		Paging:  pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func (s LedgerService) CreateLedger(ctx context.Context, r *CreateLedgerEndpoint) error {
	cmd := &ledgering.CreateLedgerCommand{
		ShopID:      r.Context.Shop.ID,
		Name:        r.Name,
		BankAccount: pbshop.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        r.Note,
		Type:        ledgering.LedgerTypeBank,
		CreatedBy:   r.Context.UserID,
	}
	if err := ledgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbLedger(cmd.Result)

	return nil
}

func (s LedgerService) UpdateLedger(ctx context.Context, r *UpdateLedgerEndpoint) error {
	cmd := &ledgering.UpdateLedgerCommand{
		ID:          r.Id,
		ShopID:      r.Context.Shop.ID,
		Name:        PString(r.Name),
		BankAccount: pbshop.Convert_api_BankAccount_To_core_BankAccount(r.BankAccount),
		Note:        PString(r.Note),
	}
	if err := ledgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = pbshop.PbLedger(cmd.Result)
	return nil
}

func (s LedgerService) DeleteLedger(ctx context.Context, r *DeleteLedgerEndpoint) error {
	cmd := &ledgering.DeleteLedgerCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := ledgerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}
