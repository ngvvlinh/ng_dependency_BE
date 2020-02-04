package types

import (
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type Partner struct {
	Id         dot.ID                   `json:"id"`
	Name       string                   `json:"name"`
	PublicName string                   `json:"public_name"`
	Type       account_type.AccountType `json:"type"`
	Phone      string                   `json:"phone"`
	// only domain, no scheme
	Website         string   `json:"website"`
	WebsiteUrl      string   `json:"website_url"`
	ImageUrl        string   `json:"image_url"`
	Email           string   `json:"email"`
	RecognizedHosts []string `json:"recognized_hosts"`
	RedirectUrls    []string `json:"redirect_urls"`

	Meta map[string]string `json:"meta,omitempty"`
}

func (m *Partner) String() string { return jsonx.MustMarshalToString(m) }
