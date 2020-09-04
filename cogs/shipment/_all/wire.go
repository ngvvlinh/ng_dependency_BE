// +build wireinject

package shipment_all

import (
	"github.com/google/wire"

	_shipment "o.o/backend/cogs/shipment"
	_ghn "o.o/backend/cogs/shipment/ghn/_all"
	_ghtk "o.o/backend/cogs/shipment/ghtk"
	_vtpost "o.o/backend/cogs/shipment/vtpost"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	imhandlers "o.o/backend/pkg/etop/logic/money-transaction/handlers"
	"o.o/backend/pkg/etop/logic/money-transaction/jtexpressimport"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
	directwebhook "o.o/backend/pkg/integration/shipping/direct/webhook"
)

var WireSet = wire.NewSet(
	_shipment.WireSet,
	_ghn.WireSet,
	_ghtk.WireSet,
	_vtpost.WireSet,
	wire.FieldsOf(new(Config), "GHN", "GHNWebhook", "GHTK", "GHTKWebhook", "VTPost", "VTPostWebhook"),
	ghnimport.WireSet,
	ghtkimport.WireSet,
	vtpostimport.WireSet,
	jtexpressimport.WireSet,
	imhandlers.WireSet,
	directwebhook.WireSet,
	SupportedShippingCarrierConfig,
	SupportedCarrierDriver,
	SupportedShipmentServices,
)
