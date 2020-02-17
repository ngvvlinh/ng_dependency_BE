package whitelabel

import (
	"text/template"

	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/connection_type"
)

type WL struct {
	identity.Partner
	Config
	Driver
}

func (w *WL) Clone() *WL {
	w2 := *w
	return &w2
}

func (w *WL) IsWhiteLabel() bool {
	return w.ID != 0
}

type Config struct {
	Key string

	Host string

	RootURL string

	AuthURL string

	SiteName string

	CompanyName string

	CompanyFullName string

	CSEmail string

	Shipment *ShipmentConfig

	Templates *Templates
}

type Driver interface {
}

type ShipmentConfig struct {
	// Các NVC topship hỗ trợ cho Partner WL này
	Topship []connection_type.ConnectionProvider
}

type Templates struct {
	RequestLoginSmsTpl         *template.Template
	NewAccountViaPartnerSmsTpl *template.Template
}

func MustParseTemplate(name, tpl string) *template.Template {
	return template.Must(template.New(name).Parse(tpl))
}
