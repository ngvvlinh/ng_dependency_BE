package shop

import (
	"context"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/receipting"

	"etop.vn/api/main/ledgering"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		GetLedger,
		GetLedgers,
		CreateLedger,
		UpdateLedger,
		DeleteLedger)
}

func GetLedger(ctx context.Context, r *wrapshop.GetLedgerEndpoint) error {
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

func GetLedgers(ctx context.Context, r *wrapshop.GetLedgersEndpoint) error {
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

func CreateLedger(ctx context.Context, r *wrapshop.CreateLedgerEndpoint) error {
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

func UpdateLedger(ctx context.Context, r *wrapshop.UpdateLedgerEndpoint) error {
	query := &ledgering.GetLedgerByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	err := ledgerQuery.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return cm.Errorf(cm.NotFound, err, "sổ quỹ không tồn tại")
	case cm.NoError:
		if query.Result.Type == string(ledgering.LedgerTypeCash) {
			return cm.Errorf(cm.FailedPrecondition, nil, "bạn không được sửa sổ quỹ mặc định")
		}
	default:
		return err
	}

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

func DeleteLedger(ctx context.Context, r *wrapshop.DeleteLedgerEndpoint) error {
	query := &receipting.ListReceiptsByLedgerIDsQuery{
		ShopID:    r.Context.Shop.ID,
		LedgerIDs: []int64{r.Id},
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result.Receipts) != 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "sổ quỹ đã được được sử dụng, không thể xóa")
	}

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
