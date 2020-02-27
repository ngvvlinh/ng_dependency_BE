package convertpb

import (
	"etop.vn/api/main/connectioning"
	extpartner "etop.vn/api/top/external/partner"
)

func PbShipmentConnection(in *connectioning.Connection) *extpartner.ShipmentConnection {
	if in == nil {
		return nil
	}
	res := &extpartner.ShipmentConnection{
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

func PbShipmentConnections(items []*connectioning.Connection) []*extpartner.ShipmentConnection {
	result := make([]*extpartner.ShipmentConnection, len(items))
	for i, item := range items {
		result[i] = PbShipmentConnection(item)
	}
	return result
}
