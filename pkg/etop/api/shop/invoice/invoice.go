package invoice

import (
	"context"

	"o.o/api/main/invoicing"
	api "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type InvoiceService struct {
	session.Session
	InvoiceAggr  invoicing.CommandBus
	InvoiceQuery invoicing.QueryBus
}

func (s *InvoiceService) Clone() api.InvoiceService { res := *s; return &res }

func (s *InvoiceService) GetInvoices(ctx context.Context, r *types.GetShopInvoicesRequest) (*types.GetShopInvoicesResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	query := &invoicing.ListInvoicesQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
	}

	if r.Filter != nil {
		query.RefID = r.Filter.RefID
		query.RefType = r.Filter.RefType
		query.DateFrom = r.Filter.DateFrom
		query.DateTo = r.Filter.DateTo
		query.Type = r.Filter.Type
	}
	if err = s.InvoiceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	res := convertpball.PbInvoices(query.Result.Invoices)
	result := &types.GetShopInvoicesResponse{
		Invoices: res,
		Paging:   cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return result, nil
}

func (s *InvoiceService) GetInvoice(ctx context.Context, r *pbcm.IDRequest) (*types.Invoice, error) {
	query := &invoicing.GetInvoiceByIDQuery{
		ID:        r.Id,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.InvoiceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpball.PbInvoice(query.Result)
	return res, nil
}
