package invoice

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateInvoice(context.Context, *CreateInvoiceArgs) (*InvoiceFtLine, error)
	CreateInvoiceBySubrID(context.Context, *CreateInvoiceBySubrIDArgs) (*InvoiceFtLine, error)
	UpdateInvoicePaymentInfo(context.Context, *UpdateInvoicePaymentInfoArgs) error
	UpdateInvoiceStatus(context.Context, *UpdateInvoiceStatusArgs) error
	DeleteInvoice(context.Context, *DeleteInvoiceArgs) error
	ManualPaymentInvoice(context.Context, *ManualPaymentInvoiceArgs) error

	PaymentInvoice(context.Context, *PaymentInvoiceArgs) error
}

type QueryService interface {
	GetInvoiceByID(ctx context.Context, ID dot.ID, AccountID dot.ID) (*InvoiceFtLine, error)
	ListInvoices(context.Context, *ListInvoicesArgs) (*ListInvoicesResponse, error)
}

// +convert:create=Invoice
type CreateInvoiceArgs struct {
	AccountID    dot.ID
	TotalAmount  int
	Lines        []*InvoiceLine
	Description  string
	Customer     *types.CustomerInfo
	ReferralType subject_referral.SubjectReferral
}

type UpdateInvoicePaymentInfoArgs struct {
	ID            dot.ID
	AccountID     dot.ID
	PaymentID     dot.ID
	PaymentStatus status4.Status
}

type UpdateInvoiceStatusArgs struct {
	ID        dot.ID
	AccountID dot.ID
	Status    status4.NullStatus
}

type ManualPaymentInvoiceArgs struct {
	ID          dot.ID
	AccountID   dot.ID
	TotalAmount int
}

type CreateInvoiceBySubrIDArgs struct {
	SubscriptionID dot.ID
	AccountID      dot.ID
	TotalAmount    int
	Customer       *types.CustomerInfo
	Description    string
}

type ListInvoicesArgs struct {
	AccountID dot.ID
	Paging    meta.Paging
	Filters   meta.Filters
	RefID     dot.ID
	RefType   subject_referral.SubjectReferral
	DateFrom  time.Time
	DateTo    time.Time
}

type ListInvoicesResponse struct {
	Invoices []*InvoiceFtLine
	Paging   meta.PageInfo
}

type PaymentInvoiceArgs struct {
	InvoiceID       dot.ID
	AccountID       dot.ID
	OwnerID         dot.ID
	TotalAmount     int
	PaymentMethod   payment_method.PaymentMethod
	ServiceClassify service_classify.NullServiceClassify
}

type DeleteInvoiceArgs struct {
	ID          dot.ID
	AccountID   dot.ID
	ForceDelete bool
}
