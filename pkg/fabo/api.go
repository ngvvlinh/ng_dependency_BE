package fabo

import (
	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/fabo/fbusering"
	"etop.vn/backend/com/fabo/pkg/fbclient"
)

var (
	fbUserQuery fbusering.QueryBus
	fbUserAggr  fbusering.CommandBus
	fbPageQuery fbpaging.QueryBus
	fbPageAggr  fbpaging.CommandBus
	appScopes   map[string]string
	fbClient    *fbclient.FbClient
)

func Init(
	fbUserQ fbusering.QueryBus,
	fbuserA fbusering.CommandBus,
	fbPageQ fbpaging.QueryBus,
	fbPageA fbpaging.CommandBus,
	_fbClient *fbclient.FbClient,
	_appScopes map[string]string,
) {
	fbUserQuery = fbUserQ
	fbUserAggr = fbuserA
	fbPageQuery = fbPageQ
	fbPageAggr = fbPageA
	fbClient = _fbClient
	appScopes = _appScopes
}

type PageService struct{}

var pageService = &PageService{}
