package admin

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type MoneyTransactionService struct {
	session.Session

	MoneyTxQuery moneytx.QueryBus
	MoneyTxAggr  moneytx.CommandBus
}

func (s *MoneyTransactionService) Clone() admin.MoneyTransactionService {
	res := *s
	return &res
}

func (s *MoneyTransactionService) CreateMoneyTransaction(ctx context.Context, q *admin.CreateMoneyTransactionRequest) (*types.MoneyTransaction, error) {
	cmd := &moneytx.CreateMoneyTxShippingCommand{
		ShopID:         q.ShopID,
		FulfillmentIDs: q.FulfillmentIDs,
		TotalCOD:       q.TotalCOD,
		TotalAmount:    q.TotalAmount,
		TotalOrders:    q.TotalOrders,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbMoneyTxShipping(cmd.Result), nil
}

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *pbcm.IDRequest) (*types.MoneyTransaction, error) {
	query := &moneytx.GetMoneyTxShippingByIDQuery{
		MoneyTxShippingID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShipping(query.Result)
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *admin.GetMoneyTransactionsRequest) (*types.MoneyTransactionsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingsQuery{
		MoneyTxShippingIDs: q.Ids,
		ShopID:             q.ShopId,
		Paging:             *paging,
		Filters:            cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.MoneyTransactionsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippings(query.Result.MoneyTxShippings),
	}
	return result, nil
}

func (s *MoneyTransactionService) UpdateMoneyTransaction(ctx context.Context, q *admin.UpdateMoneyTransactionRequest) (*types.MoneyTransaction, error) {
	cmd := &moneytx.UpdateMoneyTxShippingInfoCommand{
		MoneyTxShippingID: q.Id,
		Note:              q.Note,
		InvoiceNumber:     q.InvoiceNumber,
		BankAccount:       convertpb.Convert_api_BankAccount_To_core_BankAccount(q.BankAccount),
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShipping(cmd.Result)
	return result, nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransaction(ctx context.Context, q *admin.ConfirmMoneyTransactionRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &moneytx.ConfirmMoneyTxShippingCommand{
		MoneyTxShippingID: q.MoneyTransactionId,
		ShopID:            q.ShopId,
		TotalCOD:          q.TotalCod,
		TotalAmount:       q.TotalAmount,
		TotalOrders:       q.TotalOrders,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternal(ctx context.Context, q *pbcm.IDRequest) (*types.MoneyTransactionShippingExternal, error) {
	query := &moneytx.GetMoneyTxShippingExternalQuery{
		ID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingExternalFtLine(query.Result)
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingExternals(ctx context.Context, q *admin.GetMoneyTransactionShippingExternalsRequest) (*types.MoneyTransactionShippingExternalsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingExternalsQuery{
		MoneyTxShippingExternalIDs: q.Ids,
		Paging:                     *paging,
		Filters:                    cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.MoneyTransactionShippingExternalsResponse{
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactions: convertpb.PbMoneyTxShippingExternalsFtLine(query.Result.MoneyTxShippingExternals),
	}
	return result, nil
}

func (s *MoneyTransactionService) RemoveMoneyTransactionShippingExternalLines(ctx context.Context, q *admin.RemoveMoneyTransactionShippingExternalLinesRequest) (*types.MoneyTransactionShippingExternal, error) {
	cmd := &moneytx.RemoveMoneyTxShippingExternalLinesCommand{
		MoneyTxShippingExternalID: q.MoneyTransactionShippingExternalId,
		LineIDs:                   q.LineIds,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return result, nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingExternal(ctx context.Context, q *pbcm.IDRequest) (*pbcm.RemovedResponse, error) {
	cmd := &moneytx.DeleteMoneyTxShippingExternalCommand{
		ID: q.Id,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return result, nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingExternals(ctx context.Context, q *pbcm.IDsRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &moneytx.ConfirmMoneyTxShippingExternalsCommand{
		IDs: q.Ids,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result,
	}
	return result, nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingExternal(ctx context.Context, q *admin.UpdateMoneyTransactionShippingExternalRequest) (*types.MoneyTransactionShippingExternal, error) {
	cmd := &moneytx.UpdateMoneyTxShippingExternalInfoCommand{
		MoneyTxShippingExternalID: q.Id,
		BankAccount:               convertpb.BankAccountToCoreBankAccount(q.BankAccount),
		Note:                      q.Note,
		InvoiceNumber:             q.InvoiceNumber,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result)
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtop(ctx context.Context, q *pbcm.IDRequest) (*types.MoneyTransactionShippingEtop, error) {
	query := &moneytx.GetMoneyTxShippingEtopQuery{
		ID: q.Id,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingEtop(query.Result)
	return result, nil
}

func (s *MoneyTransactionService) GetMoneyTransactionShippingEtops(ctx context.Context, q *admin.GetMoneyTransactionShippingEtopsRequest) (*types.MoneyTransactionShippingEtopsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneytx.ListMoneyTxShippingEtopsQuery{
		MoneyTxShippingEtopIDs: q.Ids,
		Status:                 q.Status,
		Paging:                 *paging,
		Filter:                 cmapi.ToFilters(q.Filters),
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.MoneyTransactionShippingEtopsResponse{
		Paging:                        cmapi.PbMetaPageInfo(query.Result.Paging),
		MoneyTransactionShippingEtops: convertpb.PbMoneyTxShippingEtops(query.Result.MoneyTxShippingEtops),
	}
	return result, nil
}

func (s *MoneyTransactionService) CreateMoneyTransactionShippingEtop(ctx context.Context, q *pbcm.IDsRequest) (*types.MoneyTransactionShippingEtop, error) {
	cmd := &moneytx.CreateMoneyTxShippingEtopCommand{
		MoneyTxShippingIDs: q.Ids,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return result, nil
}

func (s *MoneyTransactionService) UpdateMoneyTransactionShippingEtop(ctx context.Context, q *admin.UpdateMoneyTransactionShippingEtopRequest) (*types.MoneyTransactionShippingEtop, error) {
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
		return nil, err
	}
	result := convertpb.PbMoneyTxShippingEtop(cmd.Result)
	return result, nil
}

func (s *MoneyTransactionService) DeleteMoneyTransactionShippingEtop(ctx context.Context, q *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &moneytx.DeleteMoneyTxShippingEtopCommand{
		ID: q.Id,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

func (s *MoneyTransactionService) ConfirmMoneyTransactionShippingEtop(ctx context.Context, q *admin.ConfirmMoneyTransactionShippingEtopRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &moneytx.ConfirmMoneyTxShippingEtopCommand{
		MoneyTxShippingEtopID: q.Id,
		TotalCOD:              q.TotalCod,
		TotalAmount:           q.TotalAmount,
		TotalOrders:           q.TotalOrders,
	}
	if err := s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}
