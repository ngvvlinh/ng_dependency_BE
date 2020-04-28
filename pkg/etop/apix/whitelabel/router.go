package whitelabel

import (
	"o.o/api/top/external/whitelabel"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/external/whitelabel
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
