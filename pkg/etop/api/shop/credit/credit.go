package credit

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/top/int/etop"
	api "o.o/api/top/int/shop"
	"o.o/api/top/types/etc/credit_type"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type CreditService struct {
	session.Session
	CreditAggr credit.CommandBus
}

func (s *CreditService) Clone() api.CreditService { res := *s; return &res }

func (s *CreditService) CreateCredit(ctx context.Context, req *api.CreateCreditRequest) (res *etop.Credit, err error) {
	cmd := credit.CreateCreditCommand{
		Amount:   req.Amount,
		ShopID:   s.SS.Shop().ID,
		Type:     credit_type.Shop,
		Classify: req.Classify,
	}
	if err := s.CreditAggr.Dispatch(ctx, &cmd); err != nil {
		return nil, err
	}
	extendedCredit := cmd.Result
	return convertpball.Convert_core_CreditExtended_to_api_Credit(extendedCredit), nil
}
