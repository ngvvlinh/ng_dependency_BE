// +build wireinject

package shipment_all

import (
	"github.com/google/wire"

	_shipment "o.o/backend/cogs/shipment"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_ntx "o.o/backend/cogs/shipment/ntx"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	shipmentwebhookall "o.o/backend/cogs/shipment/webhook/_all"
	"o.o/backend/com/main/shippingcode"
	"o.o/backend/pkg/etop/logic/money-transaction/dhlimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	imhandlers "o.o/backend/pkg/etop/logic/money-transaction/handlers"
	"o.o/backend/pkg/etop/logic/money-transaction/jtexpressimport"
	"o.o/backend/pkg/etop/logic/money-transaction/njvimport"
	"o.o/backend/pkg/etop/logic/money-transaction/snappyimport"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
	directwebhook "o.o/backend/pkg/integration/shipping/direct/webhook"
)

var WireSet = wire.NewSet(
	_shipment.WireSet,
	shipmentwebhookall.WireSet,
	_ghtk.WireSet,
	_vtpost.WireSet,
	_ntx.WireSet,
	wire.FieldsOf(new(Config), "GHN", "GHNWebhook", "GHTK", "GHTKWebhook", "VTPost", "VTPostWebhook", "NTX", "NTXWebhook"),
	ghnimport.WireSet,
	ghtkimport.WireSet,
	vtpostimport.WireSet,
	jtexpressimport.WireSet,
	njvimport.WireSet,
	dhlimport.WireSet,
	snappyimport.WireSet,
	imhandlers.WireSet,
	directwebhook.WireSet,
	shippingcode.WireSet,
	SupportedCarrierDriver,
	SupportedShipmentServices,
)
