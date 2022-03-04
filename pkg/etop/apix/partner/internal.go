package partner

import (
	"context"

	"o.o/api/main/bankstatement"
	api "o.o/api/top/external/partner"
	"o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type InternalService struct {
	session.Session

	BankStatementAggr bankstatement.CommandBus
}

func (s *InternalService) Clone() api.InternalService { res := *s; return &res }

func (s *InternalService) CreateBankStatement(ctx context.Context, r *api.CreateBankStatementRequest) (*common.Empty, error) {
	cmd := &bankstatement.CreateBankStatementCommand{
		Amount:                r.Amount,
		Description:           r.Description,
		TransferedAt:          r.TransferedAt,
		ExternalTransactionID: r.TransactionID,
		SenderName:            r.SenderName,
		SenderBankAccount:     r.SenderBankAccount,
		OtherInfo:             r.OtherInfo,
	}
	if err := s.BankStatementAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
