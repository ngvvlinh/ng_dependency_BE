package server_admin

import (
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/session"
	hotfixextension "o.o/backend/pkg/etop/logic/hotfix/extension"
	hotfixmoneytx "o.o/backend/pkg/etop/logic/hotfix/moneytx"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	moneytxhandlers "o.o/backend/pkg/etop/logic/money-transaction/handlers"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
)

type ImportServer httpx.Server

func BuildImportHandlers(
	ghnIm ghnimport.Import,
	ghtkIm ghtkimport.Import,
	vtpostIm vtpostimport.Import,
	importer moneytxhandlers.ImportService,
	hotfixMoneyTx *hotfixmoneytx.HotFixMoneyTxService,
	hotfixExtension *hotfixextension.ExtensionService,
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
	// hot fix
	//
	// create money tx shipping
	rt.POST("/api/admin.Import/CreateMoneyTransactionShipping", hotfixMoneyTx.HandleImportMoneyTransactionManual)
	// import extension portsip
	rt.POST("/api/admin.Import/Extensions", hotfixExtension.HandleImportExtension)
	return httpx.MakeServer("/api/admin.Import/", rt)
}
