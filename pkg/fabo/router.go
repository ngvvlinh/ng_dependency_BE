package fabo

import (
	service "o.o/api/top/int/fabo"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/fabo
// +gen:wrapper:package=fabo

func NewFaboServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewPageServiceServer(WrapPageService(pageService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
