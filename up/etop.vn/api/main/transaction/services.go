package transaction

import (
	"context"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateTransaction(context.Context, *CreateTransactionArgs) (*Transaction, error)

	ConfirmTransaction(ctx context.Context, trxnID int64, accountID int64) (*Transaction, error)

	CancelTransaction(ctx context.Context, trxnID int64, accountID int64) (*Transaction, error)
}

type QueryService interface {
	GetTransactionByID(ctx context.Context, trxnID int64, accountID int64) (*Transaction, error)

	ListTransactions(context.Context, *GetTransactionsArgs) (*TransactionResponse, error)

	GetBalance(context.Context, *GetBalanceArgs) (int, error)
}

type CreateTransactionArgs struct {
	ID        int64
	Amount    int
	AccountID int64
	Status    etoptypes.Status3
	Type      TransactionType
	Note      string
	Metadata  *TransactionMetadata
}

type GetBalanceArgs struct {
	AccountID       int64
	TransactionType TransactionType
}

type GetTransactionsArgs struct {
	AccountID int64
	Paging    meta.Paging
}

type TransactionResponse struct {
	Count        int
	Paging       meta.PageInfo
	Transactions []*Transaction
}
