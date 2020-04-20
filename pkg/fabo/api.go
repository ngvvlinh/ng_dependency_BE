package fabo

import (
	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/fabo/fbusering"
)

var (
	fbUserQuery fbusering.QueryBus
	fbUserAggr  fbusering.CommandBus
	fbPageQuery fbpaging.QueryBus
	fbPageAggr  fbpaging.CommandBus
	appScopes   []string
)

func Init(
	fbUserQ fbusering.QueryBus,
	fbuserA fbusering.CommandBus,
	fbPageQ fbpaging.QueryBus,
	fbPageA fbpaging.CommandBus,
	_appScopes []string,
) {
	fbUserQuery = fbUserQ
	fbUserAggr = fbuserA
	fbPageQuery = fbPageQ
	fbPageAggr = fbPageA
	appScopes = _appScopes
}

type SessionService struct{}
type PageService struct{}

var sessionService = &SessionService{}
var pageService = &PageService{}
