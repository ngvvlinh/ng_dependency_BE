package fabo

import (
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/capi/httprpc"
)

func NewFaboServer(
	hooks *httprpc.Hooks,
	ss *session.Session,
	fbUserQuery fbusering.QueryBus,
	fbUserAggr fbusering.CommandBus,
	fbPageQuery fbpaging.QueryBus,
	fbPageAggr fbpaging.CommandBus,
	appScopes map[string]string,
	fbClient *fbclient.FbClient,
) []httprpc.Server {
	faboInfo := faboinfo.New(fbPageQuery, fbUserQuery)
	pageService := NewPageService(
		ss, faboInfo,
		fbUserQuery, fbUserAggr, fbPageQuery,
		fbPageAggr, appScopes, fbClient,
	)
	servers := httprpc.MustNewServers(
		hooks,
		pageService.Clone,
	)
	return servers
}
