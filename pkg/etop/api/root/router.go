package root

import (
	"net/http"
	"strings"

	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func ProxyEtop(servers []httprpc.Server) []httprpc.Server {
	// proxy /api/root... to /api/etop
	for _, s := range servers {
		pathPrefix := strings.Replace(s.PathPrefix(), "/etop.", "/root.", 1)
		prx := &Proxy{pathPrefix, s}
		servers = append(servers, prx)
	}
	return servers
}

var _ httprpc.Server = &Proxy{}

type Proxy struct {
	pathPrefix string
	next       httprpc.Server
}

func (p *Proxy) PathPrefix() string {
	return p.pathPrefix
}

func (p *Proxy) WithHooks(builder httprpc.HooksBuilder) httprpc.Server {
	p.next = p.next.WithHooks(builder)
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newPath := strings.Replace(r.URL.Path, "/root.", "/etop.", 1)
	if newPath == r.URL.Path {
		p.next.ServeHTTP(w, r)
		return
	}
	r2 := *r
	u := *r.URL
	r2.URL = &u
	r2.URL.Path = newPath
	p.next.ServeHTTP(w, &r2)
}
