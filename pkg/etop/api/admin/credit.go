package admin

import (
	"context"

	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	creditmodelx "o.o/backend/com/main/credit/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
)

type CreditService struct{}

func (s *CreditService) Clone() *CreditService {
	res := *s
	return &res
}

func (s *CreditService) CreateCredit(ctx context.Context, q *CreateCreditEndpoint) error {
	cmd := &creditmodelx.CreateCreditCommand{
		Amount: q.Amount,
		ShopID: q.ShopId,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
		Type:   q.Type,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) GetCredit(ctx context.Context, q *GetCreditEndpoint) error {
	query := &creditmodelx.GetCreditQuery{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(query.Result)
	return nil
}

func (s *CreditService) GetCredits(ctx context.Context, q *GetCreditsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &creditmodelx.GetCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &etop.CreditsResponse{
		Credits: convertpb.PbCreditExtendeds(query.Result.Credits),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *CreditService) UpdateCredit(ctx context.Context, q *UpdateCreditEndpoint) error {
	cmd := &creditmodelx.UpdateCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
		Amount: q.Amount,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCreditExtended(cmd.Result)
	return nil
}

func (s *CreditService) ConfirmCredit(ctx context.Context, q *ConfirmCreditEndpoint) error {
	cmd := &creditmodelx.ConfirmCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *CreditService) DeleteCredit(ctx context.Context, q *DeleteCreditEndpoint) error {
	cmd := &creditmodelx.DeleteCreditCommand{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result.Deleted,
	}
	return nil
}
