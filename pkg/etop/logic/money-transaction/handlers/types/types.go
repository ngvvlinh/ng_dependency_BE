package types

import (
	"context"
	"mime/multipart"

	"o.o/api/main/moneytx"
)

type ImportMoneyTx interface {
	ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) ([]*moneytx.MoneyTransactionShippingExternalLine, error)
}
