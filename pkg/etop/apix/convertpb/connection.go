package convertpb

import (
	"o.o/api/main/connectioning"
	extpartnercarrier "o.o/api/top/external/partnercarrier"
)

func PbShipmentConnection(in *connectioning.Connection) *extpartnercarrier.ShipmentConnection {
	if in == nil {
		return nil
	}
	res := &extpartnercarrier.ShipmentConnection{
		ID:        in.ID,
		Name:      in.Name,
		Status:    in.Status,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		ImageURL:  in.ImageURL,
	}
	driverCfg := in.DriverConfig
	if driverCfg != nil {
		res.TrackingURL = driverCfg.TrackingURL
		res.CreateFulfillmentURL = driverCfg.CreateFulfillmentURL
		res.GetFulfillmentURL = driverCfg.GetFulfillmentURL
		res.CancelFulfillmentURL = driverCfg.CancelFulfillmentURL
		res.GetShippingServicesURL = driverCfg.GetShippingServicesURL
		res.SignUpURL = driverCfg.SignUpURL
		res.SignInURL = driverCfg.SignInURL
	}
	return res
}

func PbShipmentConnections(items []*connectioning.Connection) []*extpartnercarrier.ShipmentConnection {
	result := make([]*extpartnercarrier.ShipmentConnection, len(items))
	for i, item := range items {
		result[i] = PbShipmentConnection(item)
	}
	return result
}
