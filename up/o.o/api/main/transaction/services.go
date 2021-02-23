package transaction

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateTransaction(context.Context, *CreateTransactionArgs) (*Transaction, error)

	ConfirmTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)

	CancelTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)

	DeleteTransaction(ctx context.Context, trxnID, accountID dot.ID) error
}

type QueryService interface {
	GetTransactionByID(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)

	GetTransactionByReferral(context.Context, *GetTrxnByReferralArgs) (*Transaction, error)

	ListTransactions(context.Context, *GetTransactionsArgs) (*TransactionResponse, error)

	GetBalanceUser(context.Context, *GetBalanceUserArgs) (*GetBalanceUserResponse, error)
}

// +convert:create=Transaction
type CreateTransactionArgs struct {
	ID           dot.ID
	Name         string
	Amount       int
	AccountID    dot.ID
	Status       status3.Status
	Type         transaction_type.TransactionType
	Classify     service_classify.ServiceClassify
	Note         string
	ReferralType subject_referral.SubjectReferral
	ReferralIDs  []dot.ID
}

func (a *CreateTransactionArgs) Validate() error {
	if a.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing AccountID")
	}
	if a.Amount == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing amount")
	}
	return nil
}

type GetBalanceUserArgs struct {
	UserID   dot.ID
	Classify service_classify.ServiceClassify
}

type GetBalanceUserResponse struct {
	AvailableBalance int
	ActualBalance    int
	Classify         service_classify.ServiceClassify
}

type GetTransactionsArgs struct {
	AccountID dot.ID
	Paging    meta.Paging
}

type TransactionResponse struct {
	Count        int
	Paging       meta.PageInfo
	Transactions []*Transaction
}

type GetTrxnByReferralArgs struct {
	ReferralType subject_referral.SubjectReferral
	ReferralID   dot.ID
}
