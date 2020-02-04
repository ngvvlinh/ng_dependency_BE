package whitelabel

import "etop.vn/api/main/identity"

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
}

type Driver interface {
}
