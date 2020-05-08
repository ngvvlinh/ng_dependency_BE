package partnercarrier

import (
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

func NewPartnerCarrierServer(m httprpc.Muxer, ss *session.Session, hooks ...*httprpc.Hooks) {
	servers := httprpc.MustNewServers(
		httprpc.ChainHooks(hooks...),
		NewMiscService(ss).Clone,
		NewShipmentConnectionService(ss).Clone,
		NewShipmentService(ss).Clone,
	)
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
