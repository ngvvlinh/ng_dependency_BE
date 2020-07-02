package server_admin

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/session"
	imcsvghtk "o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
)

var WireSet = wire.NewSet(
	BuildImportHandlers,
)

type ImportServer httpx.Server

func BuildImportHandlers(
	ghnIm imcsvghn.Import,
	ghtkIm imcsvghtk.Import,
	vtpostIm vtpostimxlsx.Import,
	ss session.Session,
) ImportServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))

	perm := permission.Decl{Type: permission.EtopAdmin}
	rt.Use(httpx.Auth(perm, ss))
	rt.POST("/api/admin.Import/ghn/MoneyTransactions", ghnIm.HandleImportMoneyTransactions)
	rt.POST("/api/admin.Import/ghtk/MoneyTransactions", ghtkIm.HandleImportMoneyTransactions)
	rt.POST("/api/admin.Import/vtpost/MoneyTransactions", vtpostIm.HandleImportMoneyTransactions)
	return httpx.MakeServer("/api/admin.Import/", rt)
}
