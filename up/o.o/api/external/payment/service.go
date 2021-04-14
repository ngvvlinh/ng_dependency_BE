package payment

import (
	"context"
	"encoding/json"

	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreatePayment(context.Context, *CreatePaymentArgs) (*Payment, error)

	CreateOrUpdatePayment(context.Context, *CreatePaymentArgs) (*Payment, error)

	UpdateExternalPaymentInfo(context.Context, *UpdateExternalPaymentInfoArgs) (*Payment, error)
}

type QueryService interface {
	GetPaymentByID(ctx context.Context, ID dot.ID) (*Payment, error)

	GetPaymentByExternalTransID(ctx context.Context, TransactionID string) (*Payment, error)
}

type CreatePaymentArgs struct {
	ShopID          dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	PaymentProvider payment_provider.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
}

type UpdateExternalPaymentInfoArgs struct {
	ID              dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	ExternalData    json.RawMessage
	ExternalTransID string
}
