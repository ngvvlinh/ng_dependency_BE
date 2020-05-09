package fabo

import (
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/capi/httprpc"
)

func NewFaboServer(
	hooks httprpc.HooksBuilder,
	ss *session.Session,
	fbExternalUserQuery fbusering.QueryBus,
	fbExternalUserAggr fbusering.CommandBus,
	fbPagingQuery fbpaging.QueryBus,
	fbPagingAggr fbpaging.CommandBus,
	fbMessagingQuery fbmessaging.QueryBus,
	fbMessagingAggr fbmessaging.CommandBus,
	appScopes map[string]string,
	fbClient *fbclient.FbClient,
) []httprpc.Server {
	faboInfo := faboinfo.New(fbPagingQuery, fbExternalUserQuery)
	pageService := NewPageService(
		ss, faboInfo,
		fbExternalUserQuery, fbExternalUserAggr, fbPagingQuery,
		fbPagingAggr, appScopes, fbClient,
	)
	customerConversationService := NewCustomerConversationService(
		ss, faboInfo,
		fbMessagingQuery, fbMessagingAggr, fbPagingQuery,
	)
	servers := httprpc.MustNewServers(
		hooks,
		pageService.Clone,
		customerConversationService.Clone,
	)
	return servers
}
