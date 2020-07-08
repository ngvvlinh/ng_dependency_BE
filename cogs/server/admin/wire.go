package server_admin

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	moneytxhandlers "o.o/backend/pkg/etop/logic/money-transaction/handlers"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
)

var WireSet = wire.NewSet(
	BuildImportHandlers,
)

type ImportServer httpx.Server

func BuildImportHandlers(
	ghnIm ghnimport.Import,
	ghtkIm ghtkimport.Import,
	vtpostIm vtpostimport.Import,
	importer moneytxhandlers.ImportService,
	ss session.Session,
) ImportServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))

	perm := permission.Decl{Type: permission.EtopAdmin}
	rt.Use(httpx.Auth(perm, ss))
	rt.POST("/api/admin.Import/ghn/MoneyTransactions", ghnIm.HandleImportMoneyTransactions)
	rt.POST("/api/admin.Import/ghtk/MoneyTransactions", ghtkIm.HandleImportMoneyTransactions)
	rt.POST("/api/admin.Import/vtpost/MoneyTransactions", vtpostIm.HandleImportMoneyTransactions)

	rt.POST("/api/admin.Import/MoneyTxShippingExternal", importer.HandleImportMoneyTxs)
	return httpx.MakeServer("/api/admin.Import/", rt)
}
