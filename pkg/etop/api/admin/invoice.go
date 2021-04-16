package admin

import (
	"context"

	"o.o/api/subscripting/invoice"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type InvoiceService struct {
	session.Session

	InvoiceAggr  invoice.CommandBus
	InvoiceQuery invoice.QueryBus
}

func (s *InvoiceService) Clone() admin.InvoiceService {
	res := *s
	return &res
}

func (s *InvoiceService) GetInvoices(ctx context.Context, r *types.GetInvoicesRequest) (*types.GetInvoicesResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	query := &invoice.ListInvoicesQuery{
		AccountID: r.AccountID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := s.InvoiceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpball.PbInvoices(query.Result.Invoices)
	result := &types.GetInvoicesResponse{
		Invoices: res,
		Paging:   cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return result, nil
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, r *types.CreateInvoiceRequest) (*types.Invoice, error) {
	cmd := &invoice.CreateInvoiceBySubrIDCommand{
		SubscriptionID: r.SubscriptionID,
		AccountID:      r.AccountID,
		TotalAmount:    r.TotalAmount,
		Customer:       convertpball.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
		Description:    r.Description,
	}
	if err := s.InvoiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpball.PbInvoice(cmd.Result)
	return result, nil
}

func (s *InvoiceService) ManualPaymentInvoice(ctx context.Context, r *types.ManualPaymentInvoiceRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &invoice.ManualPaymentInvoiceCommand{
		ID:          r.InvoiceID,
		AccountID:   r.AccountID,
		TotalAmount: r.TotalAmount,
	}
	if err := s.InvoiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *InvoiceService) DeleteInvoice(ctx context.Context, r *types.SubscriptionIDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &invoice.DeleteInvoiceCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.InvoiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}
