package subscriptionbill

import (
	"context"

	"o.o/api/meta"
	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/status4"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSubscriptionBill(context.Context, *CreateSubscriptionBillArgs) (*SubscriptionBillFtLine, error)
	CreateSubscriptionBillBySubrID(context.Context, *CreateSubscriptionBillBySubrIDArgs) (*SubscriptionBillFtLine, error)
	UpdateSubscriptionBillPaymentInfo(context.Context, *UpdateSubscriptionBillPaymentInfoArgs) error
	UpdateSubscriptionBillStatus(context.Context, *UpdateSubscriptionBillStatusArgs) error
	DeleteSubsciptionBill(ctx context.Context, ID dot.ID, AccountID dot.ID) error
	ManualPaymentSubscriptionBill(context.Context, *ManualPaymentSubrBillArgs) error
}

type QueryService interface {
	GetSubscriptionBillByID(ctx context.Context, ID dot.ID, AccountID dot.ID) (*SubscriptionBillFtLine, error)
	ListSubscriptionBills(context.Context, *ListSubscriptionBillsArgs) (*ListSubscriptionBillsResponse, error)
}

// +convert:create=SubscriptionBill
type CreateSubscriptionBillArgs struct {
	AccountID      dot.ID
	SubscriptionID dot.ID
	TotalAmount    int
	Lines          []*SubscriptionBillLine
	Description    string
	Customer       *types.CustomerInfo
}

type UpdateSubscriptionBillPaymentInfoArgs struct {
	ID            dot.ID
	AccountID     dot.ID
	PaymentID     dot.ID
	PaymentStatus status4.Status
}

type UpdateSubscriptionBillStatusArgs struct {
	ID        dot.ID
	AccountID dot.ID
	Status    status4.NullStatus
}

type ManualPaymentSubrBillArgs struct {
	ID          dot.ID
	AccountID   dot.ID
	TotalAmount int
}

type CreateSubscriptionBillBySubrIDArgs struct {
	SubscriptionID dot.ID
	AccountID      dot.ID
	TotalAmount    int
	Customer       *types.CustomerInfo
	Description    string
}

type ListSubscriptionBillsArgs struct {
	AccountID      dot.ID
	SubscriptionID dot.ID
	Paging         meta.Paging
	Filters        meta.Filters
}

type ListSubscriptionBillsResponse struct {
	SubscriptionBills []*SubscriptionBillFtLine
	Paging            meta.PageInfo
}
