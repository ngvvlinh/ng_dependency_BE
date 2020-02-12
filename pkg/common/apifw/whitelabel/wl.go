package whitelabel

import (
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
	return w.Key != "etop"
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
}

type Driver interface {
}

type ShipmentConfig struct {
	// Các NVC topship hỗ trợ cho Partner WL này
	Topship []connection_type.ConnectionProvider
}
