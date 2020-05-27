package partnercarrier

import (
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(ss *session.Session) Servers {
	servers := httprpc.MustNewServers(
		nil,
		NewMiscService(ss).Clone,
		NewShipmentConnectionService(ss).Clone,
		NewShipmentService(ss).Clone,
	)
	return servers
}
