package admin

import (
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

func NewSadminServer(m httprpc.Muxer, ss *session.Session, hooks ...httprpc.HooksBuilder) {
	servers := httprpc.MustNewServers(
		httprpc.ChainHooks(hooks...),
		NewMiscService(ss).Clone,
		NewUserService(ss).Clone,
	)
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
