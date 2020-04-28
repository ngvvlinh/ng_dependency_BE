package convert

import (
	"o.o/api/main/shipmentpricing/shipmentservice"
)

// +gen:convert: o.o/backend/com/main/shipmentpricing/shipmentservice/model -> o.o/api/main/shipmentpricing/shipmentservice
// +gen:convert:  o.o/api/main/shipmentpricing/shipmentservice

func updateShipmentService(args *shipmentservice.UpdateShipmentServiceArgs, out *shipmentservice.ShipmentService) {
	apply_shipmentservice_UpdateShipmentServiceArgs_shipmentservice_ShipmentService(args, out)

	out.AvailableLocations = nil
	out.BlacklistLocations = nil
	return
}
