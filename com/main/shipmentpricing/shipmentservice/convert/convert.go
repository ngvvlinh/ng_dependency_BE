package convert

import (
	"etop.vn/api/main/shipmentpricing/shipmentservice"
)

// +gen:convert: etop.vn/backend/com/main/shipmentpricing/shipmentservice/model->etop.vn/api/main/shipmentpricing/shipmentservice
// +gen:convert: etop.vn/api/main/shipmentpricing/shipmentservice

func updateShipmentService(args *shipmentservice.UpdateShipmentServiceArgs, out *shipmentservice.ShipmentService) {
	apply_shipmentservice_UpdateShipmentServiceArgs_shipmentservice_ShipmentService(args, out)

	out.AvailableLocations = nil
	out.BlacklistLocations = nil
	return
}
