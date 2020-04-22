package fabo

import (
	service "etop.vn/api/top/int/fabo"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/int/fabo
// +gen:wrapper:package=fabo

func NewFaboServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewPageServiceServer(WrapPageService(pageService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
