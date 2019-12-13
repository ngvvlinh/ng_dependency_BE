package transaction

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateTransaction(context.Context, *CreateTransactionArgs) (*Transaction, error)

	ConfirmTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)

	CancelTransaction(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)
}

type QueryService interface {
	GetTransactionByID(ctx context.Context, trxnID dot.ID, accountID dot.ID) (*Transaction, error)

	ListTransactions(context.Context, *GetTransactionsArgs) (*TransactionResponse, error)

	GetBalance(context.Context, *GetBalanceArgs) (int, error)
}

type CreateTransactionArgs struct {
	ID        dot.ID
	Amount    int
	AccountID dot.ID
	Status    status3.Status
	Type      TransactionType
	Note      string
	Metadata  *TransactionMetadata
}

type GetBalanceArgs struct {
	AccountID       dot.ID
	TransactionType TransactionType
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
