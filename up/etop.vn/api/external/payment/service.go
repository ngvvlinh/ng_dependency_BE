package payment

import (
	"context"
	"encoding/json"
	"time"

	"etop.vn/api/main/etop"
)

type Aggregate interface {
	CreateOrUpdatePayment(context.Context, *CreatePaymentArgs) (*Payment, error)

	UpdateExternalPaymentInfo(context.Context, *UpdateExternalPaymentInfoArgs) (*Payment, error)
}

type QueryService interface {
	GetPaymentByID(ctx context.Context, ID int64) (*Payment, error)

	GetPaymentByExternalTransID(ctx context.Context, TransactionID string) (*Payment, error)
}

type CreatePaymentArgs struct {
	Amount          int
	Status          etop.Status4
	State           PaymentState
	PaymentProvider PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}

type UpdateExternalPaymentInfoArgs struct {
	ID              int64
	Amount          int
	Status          etop.Status4
	State           PaymentState
	ExternalData    json.RawMessage
	ExternalTransID string
}
