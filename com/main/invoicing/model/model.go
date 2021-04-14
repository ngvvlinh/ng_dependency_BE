package model

import (
	"time"

	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	subcriptingsharemodel "o.o/backend/com/subscripting/sharemodel"
	"o.o/capi/dot"
)

// +sqlgen
type Invoice struct {
	ID            dot.ID `paging:"id"`
	AccountID     dot.ID
	TotalAmount   int
	Description   string
	PaymentID     dot.ID
	PaymentStatus status4.Status
	Status        status4.Status
	Customer      *subcriptingsharemodel.CustomerInfo
	CreatedAt     time.Time `sq:"create" paging:"created_at"`
	UpdatedAt     time.Time `sq:"create" paging:"updated_at"`
	DeletedAt     time.Time
	WLPartnerID   dot.ID
	ReferralType  subject_referral.SubjectReferral
	ReferralIDs   []dot.ID
	Classify      service_classify.ServiceClassify
	Type          invoice_type.InvoiceType
}

// +sqlgen
type InvoiceLine struct {
	ID           dot.ID
	LineAmount   int
	Price        int
	Quantity     int
	Description  string
	InvoiceID    dot.ID
	ReferralType subject_referral.SubjectReferral
	ReferralID   dot.ID
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

type InvoiceFtLine struct {
	*Invoice
	Lines []*InvoiceLine
}
