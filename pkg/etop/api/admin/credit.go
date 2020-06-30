package admin

import (
	"context"

	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	creditmodelx "o.o/backend/com/main/credit/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type CreditService struct {
	session.Session
}

func (s *CreditService) Clone() admin.CreditService {
	res := *s
	return &res
}

func (s *CreditService) CreateCredit(ctx context.Context, q *admin.CreateCreditRequest) (*etop.Credit, error) {
	cmd := &creditmodelx.CreateCreditCommand{
		Amount: q.Amount,
		ShopID: q.ShopId,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
		Type:   q.Type,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCreditExtended(cmd.Result)
	return result, nil
}

func (s *CreditService) GetCredit(ctx context.Context, q *admin.GetCreditRequest) (*etop.Credit, error) {
	query := &creditmodelx.GetCreditQuery{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbCreditExtended(query.Result)
	return result, nil
}

func (s *CreditService) GetCredits(ctx context.Context, q *admin.GetCreditsRequest) (*etop.CreditsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &creditmodelx.GetCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &etop.CreditsResponse{
		Credits: convertpb.PbCreditExtendeds(query.Result.Credits),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *CreditService) UpdateCredit(ctx context.Context, q *admin.UpdateCreditRequest) (*etop.Credit, error) {
	cmd := &creditmodelx.UpdateCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
		Amount: q.Amount,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCreditExtended(cmd.Result)
	return result, nil
}

func (s *CreditService) ConfirmCredit(ctx context.Context, q *admin.ConfirmCreditRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &creditmodelx.ConfirmCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}

func (s *CreditService) DeleteCredit(ctx context.Context, q *pbcm.IDRequest) (*pbcm.RemovedResponse, error) {
	cmd := &creditmodelx.DeleteCreditCommand{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result.Deleted,
	}
	return result, nil
}
