package bankstatement

import (
	"context"
	"time"

	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateBankStatement(context.Context, *CreateBankStatementArgs) (*BankStatement, error)
}

type QueryService interface {
	GetBankStatement(context.Context, *GetBankStatementArgs) (*BankStatement, error)
}

// +convert:create=BankStatement
type CreateBankStatementArgs struct {
	ID                    dot.ID
	Amount                int
	Description           string // format: {mã shop} {số điện thoại}
	AccountID             dot.ID
	TransferedAt          time.Time
	ExternalTransactionID string
	SenderName            string
	SenderBankAccount     string
	OtherInfo             map[string]string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (a *CreateBankStatementArgs) Validate() error {
	if a.Amount == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing amount")
	}
	if a.Description == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing description")
	}
	if a.TransferedAt.IsZero() {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing transferd_at")
	}
	if a.SenderName == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing sender name")
	}
	return nil
}

type GetBankStatementArgs struct {
	ID                    dot.ID
	ExternalTransactionID string
}
