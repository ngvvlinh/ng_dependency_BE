package admin

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type MoneyTransactionService struct {
	MoneyTxQuery moneytx.QueryBus
	MoneyTxAggr  moneytx.CommandBus
}

func (s *MoneyTransactionService) Clone() *MoneyTransactionService {
	res := *s
	return &res
}

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *GetMoneyTransactionEndpoint) error {
	query := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShipping(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingsQuery{
		MoneyTxShippingIDs: q.Ids,
		ShopID:             q.ShopId,
		Paging:             *paging,
		Filters:            cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippings(query.Result.MoneyTxShippings),
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransaction(ctx context.Context, q *UpdateMoneyTransactionEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingInfoCommand{
		MoneyTxShippingID: q.Id,
		Note:              q.Note,
		InvoiceNumber:     q.InvoiceNumber,
		BankAccount:       convertpb.Convert_api_BankAccount_To_core_BankAccount(q.BankAccount),
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShipping(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransaction(ctx context.Context, q *ConfirmMoneyTransactionEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingCommand{
		MoneyTxShippingID: q.MoneyTransactionId,
		ShopID:            q.ShopId,
		TotalCOD:          q.TotalCod,
		TotalAmount:       q.TotalAmount,
		TotalOrders:       q.TotalOrders,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternal(ctx context.Context, q *GetMoneyTransactionShippingExternalEndpoint) error {
	query := &moneytx.GetMoneyTxShippingExternalQuery{
		ID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternals(ctx context.Context, q *GetMoneyTransactionShippingExternalsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingExternalsQuery{
		MoneyTxShippingExternalIDs: q.Ids,
		Paging:                     *paging,
		Filters:                    cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionShippingExternalsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippingExternalsFtLine(query.Result.MoneyTxShippingExternals),
	}
	return nil
}

func (s *MoneyTransactionService) RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *RemoveMoneyTransactionShippingExternalLinesEndpoint) error {
	cmd := &moneytx.RemoveMoneyTxShippingExternalLinesCommand{
		MoneyTxShippingExternalID: q.MoneyTransactionShippingExternalId,
		LineIDs:                   q.LineIds,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingExternal(ctx context.Context, q *DeleteMoneyTransactionShippingExternalEndpoint) error {
	cmd := &moneytx.DeleteMoneyTxShippingExternalCommand{
		ID: q.Id,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *ConfirmMoneyTransactionShippingExternalsEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingExternalsCommand{
		IDs: q.Ids,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingExternal(ctx context.Context, q *UpdateMoneyTransactionShippingExternalEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingExternalInfoCommand{
		MoneyTxShippingExternalID: q.Id,
		BankAccount:               convertpb.BankAccountToCoreBankAccount(q.BankAccount),
		Note:                      q.Note,
		InvoiceNumber:             q.InvoiceNumber,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtop(ctx context.Context, q *GetMoneyTransactionShippingEtopEndpoint) error {
	query := &moneytx.GetMoneyTxShippingEtopQuery{
		ID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtops(ctx context.Context, q *GetMoneyTransactionShippingEtopsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingEtopsQuery{
		MoneyTxShippingEtopIDs: q.Ids,
		Status:                 q.Status,
		Paging:                 *paging,
		Filter:                 cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.MoneyTransactionShippingEtopsResponse{
		Paging:                        cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactionShippingEtops: convertpb.PbMoneyTxShippingEtops(query.Result.MoneyTxShippingEtops),
	}
	return nil
}

func (s *MoneyTransactionService) CreateMoneyTransactionShippingEtop(ctx context.Context, q *CreateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.CreateMoneyTxShippingEtopCommand{
		MoneyTxShippingIDs: q.Ids,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingEtop(ctx context.Context, q *UpdateMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.UpdateMoneyTxShippingEtopCommand{
		MoneyTxShippingEtopID: q.Id,
		BankAccount:           convertpb.Convert_api_BankAccount_To_core_BankAccount(q.BankAccount),
		Note:                  q.Note,
		InvoiceNumber:         q.InvoiceNumber,
		Adds:                  q.Adds,
		Deletes:               q.Deletes,
		ReplaceAll:            q.ReplaceAll,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingEtop(ctx context.Context, q *DeleteMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.DeleteMoneyTxShippingEtopCommand{
		ID: q.Id,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *ConfirmMoneyTransactionShippingEtopEndpoint) error {
	cmd := &moneytx.ConfirmMoneyTxShippingEtopCommand{
		MoneyTxShippingEtopID: q.Id,
		TotalCOD:              q.TotalCod,
		TotalAmount:           q.TotalAmount,
		TotalOrders:           q.TotalOrders,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}
