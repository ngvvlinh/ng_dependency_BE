package admin

import (
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

func NewSadminServer(ss *session.Session, hooks ...httprpc.HooksBuilder) []httprpc.Server {
	servers := httprpc.MustNewServers(
		httprpc.ChainHooks(hooks...),
		NewMiscService(ss).Clone,
		NewUserService(ss).Clone,
	)
	return servers
}
