package whitelabel

import (
	"text/template"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/wl_type"
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

func (w *WL) IsWLPartnerPOS() bool {
	return w.IsWhiteLabel() && w.WLType == wl_type.POS
}

type Config struct {
	Key string

	Host string

	RootURL string

	AuthURL string

	InviteUserURLByEmail string

	InviteUserURLByPhone string

	SiteName string

	CompanyName string

	CompanyFullName string

	CSEmail string

	Shipment *ShipmentConfig

	Templates *Templates

	DatabaseName string

	WLType wl_type.WhiteLabelType

	// define that this whitelabel partner can not get from host (x-forwarded-header)
	IgnoreParseFromHost bool
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
