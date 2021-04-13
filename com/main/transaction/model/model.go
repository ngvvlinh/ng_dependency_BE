package model

import (
	"time"

	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	"o.o/capi/dot"
)

// +sqlgen
type Transaction struct {
	Name         string
	ID           dot.ID `paging:"id"`
	Amount       int
	AccountID    dot.ID
	Status       status3.Status
	Type         transaction_type.TransactionType
	Classify     service_classify.ServiceClassify
	Note         string
	ReferralType subject_referral.SubjectReferral
	ReferralIDs  []dot.ID
	CreatedAt    time.Time `sq:"create" paging:"created_at"`
	UpdatedAt    time.Time `sq:"update" paging:"updated_at"`
}
