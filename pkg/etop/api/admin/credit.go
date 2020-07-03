package admin

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type CreditService struct {
	session.Session

	CreditAggr  credit.CommandBus
	CreditQuery credit.QueryBus
}

func (s *CreditService) Clone() admin.CreditService {
	res := *s
	return &res
}

func (s *CreditService) CreateCredit(ctx context.Context, q *admin.CreateCreditRequest) (*etop.Credit, error) {
	cmd := &credit.CreateCreditCommand{
		Amount: q.Amount,
		ShopID: q.ShopId,
		PaidAt: cmapi.PbTimeToModel(q.PaidAt),
		Type:   q.Type,
	}

	if err := s.CreditAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_CreditExtended_to_api_Credit(cmd.Result)
	return result, nil
}

func (s *CreditService) GetCredit(ctx context.Context, q *admin.GetCreditRequest) (*etop.Credit, error) {
	query := &credit.GetCreditQuery{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := s.CreditQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_CreditExtended_to_api_Credit(query.Result)
	return result, nil
}

func (s *CreditService) GetCredits(ctx context.Context, q *admin.GetCreditsRequest) (*etop.CreditsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &credit.ListCreditsQuery{
		ShopID: q.ShopId,
		Paging: paging,
	}
	if err := s.CreditQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &etop.CreditsResponse{
		Credits: convertpb.Convert_core_CreditExtendeds_to_api_Credits(query.Result.Credits),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *CreditService) ConfirmCredit(ctx context.Context, q *admin.ConfirmCreditRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &credit.ConfirmCreditCommand{
		ID:     q.Id,
		ShopID: q.ShopId,
	}
	if err := s.CreditAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *CreditService) DeleteCredit(ctx context.Context, q *pbcm.IDRequest) (*pbcm.RemovedResponse, error) {
	cmd := &credit.DeleteCreditCommand{
		ID: q.Id,
	}
	if err := s.CreditAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return result, nil
}
