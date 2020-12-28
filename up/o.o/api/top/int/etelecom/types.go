package etelecom

import (
	"o.o/api/top/int/shop"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type SummaryEtelecomRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (r *SummaryEtelecomRequest) String() string { return jsonx.MustMarshalToString(r) }

type SummaryEtelecomResponse struct {
	Tables []*shop.SummaryTable `json:"tables"`
}

func (r *SummaryEtelecomResponse) String() string { return jsonx.MustMarshalToString(r) }

type CreateUserAndAssignExtensionRequest struct {
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	HotlineID dot.ID `json:"hotline_id"`
}

func (r *CreateUserAndAssignExtensionRequest) String() string { return jsonx.MustMarshalToString(r) }
