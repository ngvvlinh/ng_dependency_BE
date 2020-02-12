package whitelabel

import (
	"etop.vn/api/top/external/whitelabel"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/external/whitelabel
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

func NewWhiteLabelServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		whitelabel.NewImportServiceServer(WrapImportService(importService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
