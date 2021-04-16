package types

import (
	"time"

	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type Transaction struct {
	Name         string                           `json:"name"`
	ID           dot.ID                           `json:"id"`
	Amount       int                              `json:"amount"`
	AccountID    dot.ID                           `json:"account_id"`
	Status       status3.Status                   `json:"status"`
	Type         transaction_type.TransactionType `json:"type"`
	Classify     service_classify.ServiceClassify `json:"classify"`
	Note         string                           `json:"note"`
	ReferralType subject_referral.SubjectReferral `json:"referral_type"`
	ReferralIDs  []dot.ID                         `json:"referral_ids"'`
	CreatedAt    time.Time                        `json:"created_at"`
	UpdatedAt    time.Time                        `json:"updated_at"`
}

func (m *Transaction) String() string { return jsonx.MustMarshalToString(m) }

type GetTransactionRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTransactionsRequest struct {
	Paging *common.CursorPaging   `json:"paging"`
	Filter *GetTransactionsFilter `json:"filter"`
}

func (m *GetTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTransactionsResponse struct {
	Transactions []*Transaction         `json:"transactions"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (m *GetTransactionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTransactionsFilter struct {
	RefID    dot.ID                           `json:"ref_id"`
	RefType  subject_referral.SubjectReferral `json:"ref_type"`
	DateFrom time.Time                        `json:"date_from"`
	DateTo   time.Time                        `json:"date_to"`
}

type GetAdminTransactionsRequest struct {
	Paging *common.CursorPaging        `json:"paging"`
	Filter *GetAdminTransactionsFilter `json:"filter"`
}

func (m *GetAdminTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetAdminTransactionsFilter struct {
	AccountID dot.ID                           `json:"account_id"`
	RefID     dot.ID                           `json:"ref_id"`
	RefType   subject_referral.SubjectReferral `json:"ref_type"`
	DateFrom  time.Time                        `json:"date_from"`
	DateTo    time.Time                        `json:"date_to"`
}

func (m *GetAdminTransactionsFilter) String() string { return jsonx.MustMarshalToString(m) }
