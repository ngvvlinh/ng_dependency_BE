package invoicing

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateInvoice(context.Context, *CreateInvoiceArgs) (*InvoiceFtLine, error)
	CreateInvoiceBySubrID(context.Context, *CreateInvoiceBySubrIDArgs) (*InvoiceFtLine, error)
	UpdateInvoicePaymentInfo(context.Context, *UpdateInvoicePaymentInfoArgs) error
	DeleteInvoice(context.Context, *DeleteInvoiceArgs) error

	PaymentInvoice(context.Context, *PaymentInvoiceArgs) error
}

type QueryService interface {
	GetInvoiceByID(ctx context.Context, ID dot.ID, AccountID dot.ID) (*InvoiceFtLine, error)
	GetInvoiceByPaymentID(ctx context.Context, paymentID dot.ID) (*InvoiceFtLine, error)
	GetInvoiceByReferral(ctx context.Context, _ *GetInvoiceByReferralArgs) (*InvoiceFtLine, error)
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
	Classify     service_classify.ServiceClassify
	Type         invoice_type.InvoiceType
}

func (args *CreateInvoiceArgs) Validate() error {
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if args.Customer == nil || args.Customer.FullName == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing customer")
	}
	if len(args.Lines) == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing lines")
	}
	if args.Type == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing invoice type")
	}
	if args.TotalAmount <= 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Total amount must be >= 0")
	}
	return nil
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
	Classify       service_classify.ServiceClassify
}

type GetInvoiceByReferralArgs struct {
	ReferralType subject_referral.SubjectReferral
	ReferralIDs  []dot.ID
}

type ListInvoicesArgs struct {
	AccountID dot.ID
	Paging    meta.Paging
	Filters   meta.Filters
	RefID     dot.ID
	RefType   subject_referral.SubjectReferral
	DateFrom  time.Time
	DateTo    time.Time
	Type      invoice_type.InvoiceType
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
