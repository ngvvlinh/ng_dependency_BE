package invoice

import (
	"time"

	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/capi/dot"
)

// +gen:event:topic=event/invoice

type Invoice struct {
	ID            dot.ID
	AccountID     dot.ID
	TotalAmount   int
	Description   string
	PaymentID     dot.ID
	PaymentStatus status4.Status
	Status        status4.Status
	Customer      *types.CustomerInfo
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	WLPartnerID   dot.ID
	ReferralType  subject_referral.SubjectReferral
	ReferralIDs   []dot.ID
}

type InvoiceLine struct {
	ID           dot.ID
	LineAmount   int
	Price        int
	Quantity     int
	Description  string
	InvoiceID    dot.ID
	ReferralType subject_referral.SubjectReferral
	ReferralID   dot.ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type InvoiceFtLine struct {
	*Invoice
	Lines []*InvoiceLine
}

type InvoicePaidEvent struct {
	ID              dot.ID
	AccountID       dot.ID
	PaymentMethod   payment_method.PaymentMethod
	ServiceClassify service_classify.NullServiceClassify
}

type InvoicePayingEvent struct {
	PaymentMethod   payment_method.PaymentMethod
	ServiceClassify service_classify.NullServiceClassify
	OwnerID         dot.ID
	TotalAmount     int
}

type InvoiceDeletedEvent struct {
	InvoinceID dot.ID
}
