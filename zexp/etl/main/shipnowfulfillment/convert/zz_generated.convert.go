// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	shipnowmodel "o.o/backend/com/main/shipnow/model"
	conversion "o.o/backend/pkg/common/conversion"
	shipnowfulfillmentmodel "o.o/backend/zexp/etl/main/shipnowfulfillment/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shipnowfulfillmentmodel.DeliveryPoint)(nil), (*shipnowmodel.DeliveryPoint)(nil), func(arg, out interface{}) error {
		Convert_shipnowfulfillmentmodel_DeliveryPoint_shipnowmodel_DeliveryPoint(arg.(*shipnowfulfillmentmodel.DeliveryPoint), out.(*shipnowmodel.DeliveryPoint))
		return nil
	})
	s.Register(([]*shipnowfulfillmentmodel.DeliveryPoint)(nil), (*[]*shipnowmodel.DeliveryPoint)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipnowfulfillmentmodel_DeliveryPoints_shipnowmodel_DeliveryPoints(arg.([]*shipnowfulfillmentmodel.DeliveryPoint))
		*out.(*[]*shipnowmodel.DeliveryPoint) = out0
		return nil
	})
	s.Register((*shipnowmodel.DeliveryPoint)(nil), (*shipnowfulfillmentmodel.DeliveryPoint)(nil), func(arg, out interface{}) error {
		Convert_shipnowmodel_DeliveryPoint_shipnowfulfillmentmodel_DeliveryPoint(arg.(*shipnowmodel.DeliveryPoint), out.(*shipnowfulfillmentmodel.DeliveryPoint))
		return nil
	})
	s.Register(([]*shipnowmodel.DeliveryPoint)(nil), (*[]*shipnowfulfillmentmodel.DeliveryPoint)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipnowmodel_DeliveryPoints_shipnowfulfillmentmodel_DeliveryPoints(arg.([]*shipnowmodel.DeliveryPoint))
		*out.(*[]*shipnowfulfillmentmodel.DeliveryPoint) = out0
		return nil
	})
	s.Register((*shipnowfulfillmentmodel.ShipnowFulfillment)(nil), (*shipnowmodel.ShipnowFulfillment)(nil), func(arg, out interface{}) error {
		Convert_shipnowfulfillmentmodel_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(arg.(*shipnowfulfillmentmodel.ShipnowFulfillment), out.(*shipnowmodel.ShipnowFulfillment))
		return nil
	})
	s.Register(([]*shipnowfulfillmentmodel.ShipnowFulfillment)(nil), (*[]*shipnowmodel.ShipnowFulfillment)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipnowfulfillmentmodel_ShipnowFulfillments_shipnowmodel_ShipnowFulfillments(arg.([]*shipnowfulfillmentmodel.ShipnowFulfillment))
		*out.(*[]*shipnowmodel.ShipnowFulfillment) = out0
		return nil
	})
	s.Register((*shipnowmodel.ShipnowFulfillment)(nil), (*shipnowfulfillmentmodel.ShipnowFulfillment)(nil), func(arg, out interface{}) error {
		Convert_shipnowmodel_ShipnowFulfillment_shipnowfulfillmentmodel_ShipnowFulfillment(arg.(*shipnowmodel.ShipnowFulfillment), out.(*shipnowfulfillmentmodel.ShipnowFulfillment))
		return nil
	})
	s.Register(([]*shipnowmodel.ShipnowFulfillment)(nil), (*[]*shipnowfulfillmentmodel.ShipnowFulfillment)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipnowmodel_ShipnowFulfillments_shipnowfulfillmentmodel_ShipnowFulfillments(arg.([]*shipnowmodel.ShipnowFulfillment))
		*out.(*[]*shipnowfulfillmentmodel.ShipnowFulfillment) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/shipnow/model.DeliveryPoint --//

func Convert_shipnowfulfillmentmodel_DeliveryPoint_shipnowmodel_DeliveryPoint(arg *shipnowfulfillmentmodel.DeliveryPoint, out *shipnowmodel.DeliveryPoint) *shipnowmodel.DeliveryPoint {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipnowmodel.DeliveryPoint{}
	}
	convert_shipnowfulfillmentmodel_DeliveryPoint_shipnowmodel_DeliveryPoint(arg, out)
	return out
}

func convert_shipnowfulfillmentmodel_DeliveryPoint_shipnowmodel_DeliveryPoint(arg *shipnowfulfillmentmodel.DeliveryPoint, out *shipnowmodel.DeliveryPoint) {
	out.ShippingAddress = arg.ShippingAddress   // simple assign
	out.Items = arg.Items                       // simple assign
	out.OrderID = arg.OrderID                   // simple assign
	out.OrderCode = arg.OrderCode               // simple assign
	out.GrossWeight = arg.GrossWeight           // simple assign
	out.ChargeableWeight = arg.ChargeableWeight // simple assign
	out.Length = arg.Length                     // simple assign
	out.Width = arg.Width                       // simple assign
	out.Height = arg.Height                     // simple assign
	out.BasketValue = arg.BasketValue           // simple assign
	out.CODAmount = arg.CODAmount               // simple assign
	out.TryOn = arg.TryOn                       // simple assign
	out.ShippingNote = arg.ShippingNote         // simple assign
	out.ShippingState = 0                       // zero value
}

func Convert_shipnowfulfillmentmodel_DeliveryPoints_shipnowmodel_DeliveryPoints(args []*shipnowfulfillmentmodel.DeliveryPoint) (outs []*shipnowmodel.DeliveryPoint) {
	if args == nil {
		return nil
	}
	tmps := make([]shipnowmodel.DeliveryPoint, len(args))
	outs = make([]*shipnowmodel.DeliveryPoint, len(args))
	for i := range tmps {
		outs[i] = Convert_shipnowfulfillmentmodel_DeliveryPoint_shipnowmodel_DeliveryPoint(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipnowmodel_DeliveryPoint_shipnowfulfillmentmodel_DeliveryPoint(arg *shipnowmodel.DeliveryPoint, out *shipnowfulfillmentmodel.DeliveryPoint) *shipnowfulfillmentmodel.DeliveryPoint {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipnowfulfillmentmodel.DeliveryPoint{}
	}
	convert_shipnowmodel_DeliveryPoint_shipnowfulfillmentmodel_DeliveryPoint(arg, out)
	return out
}

func convert_shipnowmodel_DeliveryPoint_shipnowfulfillmentmodel_DeliveryPoint(arg *shipnowmodel.DeliveryPoint, out *shipnowfulfillmentmodel.DeliveryPoint) {
	out.ShippingAddress = arg.ShippingAddress   // simple assign
	out.Items = arg.Items                       // simple assign
	out.OrderID = arg.OrderID                   // simple assign
	out.OrderCode = arg.OrderCode               // simple assign
	out.GrossWeight = arg.GrossWeight           // simple assign
	out.ChargeableWeight = arg.ChargeableWeight // simple assign
	out.Length = arg.Length                     // simple assign
	out.Width = arg.Width                       // simple assign
	out.Height = arg.Height                     // simple assign
	out.BasketValue = arg.BasketValue           // simple assign
	out.CODAmount = arg.CODAmount               // simple assign
	out.TryOn = arg.TryOn                       // simple assign
	out.ShippingNote = arg.ShippingNote         // simple assign
}

func Convert_shipnowmodel_DeliveryPoints_shipnowfulfillmentmodel_DeliveryPoints(args []*shipnowmodel.DeliveryPoint) (outs []*shipnowfulfillmentmodel.DeliveryPoint) {
	if args == nil {
		return nil
	}
	tmps := make([]shipnowfulfillmentmodel.DeliveryPoint, len(args))
	outs = make([]*shipnowfulfillmentmodel.DeliveryPoint, len(args))
	for i := range tmps {
		outs[i] = Convert_shipnowmodel_DeliveryPoint_shipnowfulfillmentmodel_DeliveryPoint(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/shipnow/model.ShipnowFulfillment --//

func Convert_shipnowfulfillmentmodel_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(arg *shipnowfulfillmentmodel.ShipnowFulfillment, out *shipnowmodel.ShipnowFulfillment) *shipnowmodel.ShipnowFulfillment {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipnowmodel.ShipnowFulfillment{}
	}
	convert_shipnowfulfillmentmodel_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(arg, out)
	return out
}

func convert_shipnowfulfillmentmodel_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(arg *shipnowfulfillmentmodel.ShipnowFulfillment, out *shipnowmodel.ShipnowFulfillment) {
	out.ID = arg.ID                                                 // simple assign
	out.ShopID = arg.ShopID                                         // simple assign
	out.PartnerID = arg.PartnerID                                   // simple assign
	out.OrderIDs = arg.OrderIDs                                     // simple assign
	out.PickupAddress = arg.PickupAddress                           // simple assign
	out.Carrier = 0                                                 // types do not match
	out.ShippingServiceCode = arg.ShippingServiceCode               // simple assign
	out.ShippingServiceFee = arg.ShippingServiceFee                 // simple assign
	out.ShippingServiceName = arg.ShippingServiceName               // simple assign
	out.ShippingServiceDescription = arg.ShippingServiceDescription // simple assign
	out.ChargeableWeight = arg.ChargeableWeight                     // simple assign
	out.GrossWeight = arg.GrossWeight                               // simple assign
	out.BasketValue = arg.BasketValue                               // simple assign
	out.CODAmount = arg.CODAmount                                   // simple assign
	out.ShippingNote = arg.ShippingNote                             // simple assign
	out.RequestPickupAt = arg.RequestPickupAt                       // simple assign
	out.DeliveryPoints = Convert_shipnowfulfillmentmodel_DeliveryPoints_shipnowmodel_DeliveryPoints(arg.DeliveryPoints)
	out.CancelReason = arg.CancelReason                   // simple assign
	out.Status = arg.Status                               // simple assign
	out.ConfirmStatus = arg.ConfirmStatus                 // simple assign
	out.ShippingStatus = arg.ShippingStatus               // simple assign
	out.EtopPaymentStatus = arg.EtopPaymentStatus         // simple assign
	out.ShippingState = arg.ShippingState                 // simple assign
	out.ShippingCode = arg.ShippingCode                   // simple assign
	out.FeeLines = arg.FeeLines                           // simple assign
	out.CarrierFeeLines = arg.CarrierFeeLines             // simple assign
	out.TotalFee = arg.TotalFee                           // simple assign
	out.ShippingCreatedAt = arg.ShippingCreatedAt         // simple assign
	out.ShippingPickingAt = arg.ShippingPickingAt         // simple assign
	out.ShippingDeliveringAt = arg.ShippingDeliveringAt   // simple assign
	out.ShippingDeliveredAt = arg.ShippingDeliveredAt     // simple assign
	out.ShippingCancelledAt = arg.ShippingCancelledAt     // simple assign
	out.SyncStatus = arg.SyncStatus                       // simple assign
	out.SyncStates = arg.SyncStates                       // simple assign
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
	out.CODEtopTransferedAt = arg.CODEtopTransferedAt     // simple assign
	out.ShippingSharedLink = arg.ShippingSharedLink       // simple assign
	out.AddressToProvinceCode = arg.AddressToProvinceCode // simple assign
	out.AddressToDistrictCode = arg.AddressToDistrictCode // simple assign
	out.ConnectionID = 0                                  // zero value
	out.ConnectionMethod = 0                              // zero value
	out.ExternalID = ""                                   // zero value
	out.Rid = arg.Rid                                     // simple assign
}

func Convert_shipnowfulfillmentmodel_ShipnowFulfillments_shipnowmodel_ShipnowFulfillments(args []*shipnowfulfillmentmodel.ShipnowFulfillment) (outs []*shipnowmodel.ShipnowFulfillment) {
	if args == nil {
		return nil
	}
	tmps := make([]shipnowmodel.ShipnowFulfillment, len(args))
	outs = make([]*shipnowmodel.ShipnowFulfillment, len(args))
	for i := range tmps {
		outs[i] = Convert_shipnowfulfillmentmodel_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipnowmodel_ShipnowFulfillment_shipnowfulfillmentmodel_ShipnowFulfillment(arg *shipnowmodel.ShipnowFulfillment, out *shipnowfulfillmentmodel.ShipnowFulfillment) *shipnowfulfillmentmodel.ShipnowFulfillment {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipnowfulfillmentmodel.ShipnowFulfillment{}
	}
	convert_shipnowmodel_ShipnowFulfillment_shipnowfulfillmentmodel_ShipnowFulfillment(arg, out)
	return out
}

func convert_shipnowmodel_ShipnowFulfillment_shipnowfulfillmentmodel_ShipnowFulfillment(arg *shipnowmodel.ShipnowFulfillment, out *shipnowfulfillmentmodel.ShipnowFulfillment) {
	out.ID = arg.ID                                                 // simple assign
	out.ShopID = arg.ShopID                                         // simple assign
	out.PartnerID = arg.PartnerID                                   // simple assign
	out.OrderIDs = arg.OrderIDs                                     // simple assign
	out.PickupAddress = arg.PickupAddress                           // simple assign
	out.Carrier = ""                                                // types do not match
	out.ShippingServiceCode = arg.ShippingServiceCode               // simple assign
	out.ShippingServiceFee = arg.ShippingServiceFee                 // simple assign
	out.ShippingServiceName = arg.ShippingServiceName               // simple assign
	out.ShippingServiceDescription = arg.ShippingServiceDescription // simple assign
	out.ChargeableWeight = arg.ChargeableWeight                     // simple assign
	out.GrossWeight = arg.GrossWeight                               // simple assign
	out.BasketValue = arg.BasketValue                               // simple assign
	out.CODAmount = arg.CODAmount                                   // simple assign
	out.ShippingNote = arg.ShippingNote                             // simple assign
	out.RequestPickupAt = arg.RequestPickupAt                       // simple assign
	out.DeliveryPoints = Convert_shipnowmodel_DeliveryPoints_shipnowfulfillmentmodel_DeliveryPoints(arg.DeliveryPoints)
	out.CancelReason = arg.CancelReason                   // simple assign
	out.Status = arg.Status                               // simple assign
	out.ConfirmStatus = arg.ConfirmStatus                 // simple assign
	out.ShippingStatus = arg.ShippingStatus               // simple assign
	out.EtopPaymentStatus = arg.EtopPaymentStatus         // simple assign
	out.ShippingState = arg.ShippingState                 // simple assign
	out.ShippingCode = arg.ShippingCode                   // simple assign
	out.FeeLines = arg.FeeLines                           // simple assign
	out.CarrierFeeLines = arg.CarrierFeeLines             // simple assign
	out.TotalFee = arg.TotalFee                           // simple assign
	out.ShippingCreatedAt = arg.ShippingCreatedAt         // simple assign
	out.ShippingPickingAt = arg.ShippingPickingAt         // simple assign
	out.ShippingDeliveringAt = arg.ShippingDeliveringAt   // simple assign
	out.ShippingDeliveredAt = arg.ShippingDeliveredAt     // simple assign
	out.ShippingCancelledAt = arg.ShippingCancelledAt     // simple assign
	out.SyncStatus = arg.SyncStatus                       // simple assign
	out.SyncStates = arg.SyncStates                       // simple assign
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
	out.CODEtopTransferedAt = arg.CODEtopTransferedAt     // simple assign
	out.ShippingSharedLink = arg.ShippingSharedLink       // simple assign
	out.AddressToProvinceCode = arg.AddressToProvinceCode // simple assign
	out.AddressToDistrictCode = arg.AddressToDistrictCode // simple assign
	out.Rid = arg.Rid                                     // simple assign
}

func Convert_shipnowmodel_ShipnowFulfillments_shipnowfulfillmentmodel_ShipnowFulfillments(args []*shipnowmodel.ShipnowFulfillment) (outs []*shipnowfulfillmentmodel.ShipnowFulfillment) {
	if args == nil {
		return nil
	}
	tmps := make([]shipnowfulfillmentmodel.ShipnowFulfillment, len(args))
	outs = make([]*shipnowfulfillmentmodel.ShipnowFulfillment, len(args))
	for i := range tmps {
		outs[i] = Convert_shipnowmodel_ShipnowFulfillment_shipnowfulfillmentmodel_ShipnowFulfillment(args[i], &tmps[i])
	}
	return outs
}
