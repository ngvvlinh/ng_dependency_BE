// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package fabo

import (
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/faboinfo"
)

// Injectors from wire.go:

func NewFaboServer(ss session.Session, fbExternalUserQuery fbusering.QueryBus, fbExternalUserAggr fbusering.CommandBus, fbPagingQuery fbpaging.QueryBus, fbPagingAggr fbpaging.CommandBus, fbMessagingQuery fbmessaging.QueryBus, fbMessagingAggr fbmessaging.CommandBus, appScopes map[string]string, fbClient *fbclient.FbClient, customerQ customering.QueryBus) Servers {
	faboInfo := faboinfo.New(fbPagingQuery, fbExternalUserQuery)
	pageService := NewPageService(ss, faboInfo, fbExternalUserQuery, fbExternalUserAggr, fbPagingQuery, fbPagingAggr, appScopes, fbClient)
	customerConversationService := NewCustomerConversationService(ss, faboInfo, fbMessagingQuery, fbMessagingAggr, fbPagingQuery, fbClient)
	customerService := NewCustomerService(customerQ, fbExternalUserQuery, fbExternalUserAggr, ss)
	servers := NewServer(pageService, customerConversationService, customerService)
	return servers
}
