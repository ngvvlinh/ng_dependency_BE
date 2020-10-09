// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	substate "o.o/api/top/types/etc/shipping/substate"
	shippingmodel "o.o/backend/com/main/shipping/model"
	conversion "o.o/backend/pkg/common/conversion"
	fulfillmentmodel "o.o/backend/zexp/etl/main/fulfillment/model"
	orderconvert "o.o/backend/zexp/etl/main/order/convert"
	dot "o.o/capi/dot"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*fulfillmentmodel.Fulfillment)(nil), (*shippingmodel.Fulfillment)(nil), func(arg, out interface{}) error {
		Convert_fulfillmentmodel_Fulfillment_shippingmodel_Fulfillment(arg.(*fulfillmentmodel.Fulfillment), out.(*shippingmodel.Fulfillment))
		return nil
	})
	s.Register(([]*fulfillmentmodel.Fulfillment)(nil), (*[]*shippingmodel.Fulfillment)(nil), func(arg, out interface{}) error {
		out0 := Convert_fulfillmentmodel_Fulfillments_shippingmodel_Fulfillments(arg.([]*fulfillmentmodel.Fulfillment))
		*out.(*[]*shippingmodel.Fulfillment) = out0
		return nil
	})
	s.Register((*shippingmodel.Fulfillment)(nil), (*fulfillmentmodel.Fulfillment)(nil), func(arg, out interface{}) error {
		Convert_shippingmodel_Fulfillment_fulfillmentmodel_Fulfillment(arg.(*shippingmodel.Fulfillment), out.(*fulfillmentmodel.Fulfillment))
		return nil
	})
	s.Register(([]*shippingmodel.Fulfillment)(nil), (*[]*fulfillmentmodel.Fulfillment)(nil), func(arg, out interface{}) error {
		out0 := Convert_shippingmodel_Fulfillments_fulfillmentmodel_Fulfillments(arg.([]*shippingmodel.Fulfillment))
		*out.(*[]*fulfillmentmodel.Fulfillment) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/shipping/model.Fulfillment --//

func Convert_fulfillmentmodel_Fulfillment_shippingmodel_Fulfillment(arg *fulfillmentmodel.Fulfillment, out *shippingmodel.Fulfillment) *shippingmodel.Fulfillment {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shippingmodel.Fulfillment{}
	}
	convert_fulfillmentmodel_Fulfillment_shippingmodel_Fulfillment(arg, out)
	return out
}

func convert_fulfillmentmodel_Fulfillment_shippingmodel_Fulfillment(arg *fulfillmentmodel.Fulfillment, out *shippingmodel.Fulfillment) {
	out.ID = arg.ID                                                   // simple assign
	out.OrderID = arg.OrderID                                         // simple assign
	out.ShopID = arg.ShopID                                           // simple assign
	out.PartnerID = arg.PartnerID                                     // simple assign
	out.ShopConfirm = arg.ShopConfirm                                 // simple assign
	out.ConfirmStatus = arg.ConfirmStatus                             // simple assign
	out.TotalItems = arg.TotalItems                                   // simple assign
	out.TotalWeight = arg.TotalWeight                                 // simple assign
	out.BasketValue = arg.BasketValue                                 // simple assign
	out.TotalDiscount = arg.TotalDiscount                             // simple assign
	out.TotalAmount = arg.TotalAmount                                 // simple assign
	out.TotalCODAmount = arg.TotalCODAmount                           // simple assign
	out.OriginalCODAmount = arg.OriginalCODAmount                     // simple assign
	out.ActualCompensationAmount = arg.ActualCompensationAmount       // simple assign
	out.ShippingFeeCustomer = arg.ShippingFeeCustomer                 // simple assign
	out.ShippingFeeShop = arg.ShippingFeeShop                         // simple assign
	out.ShippingFeeShopLines = arg.ShippingFeeShopLines               // simple assign
	out.ShippingServiceFee = arg.ShippingServiceFee                   // simple assign
	out.ExternalShippingFee = arg.ExternalShippingFee                 // simple assign
	out.ProviderShippingFeeLines = arg.ProviderShippingFeeLines       // simple assign
	out.EtopDiscount = arg.EtopDiscount                               // simple assign
	out.EtopFeeAdjustment = arg.EtopFeeAdjustment                     // simple assign
	out.ShippingFeeMain = arg.ShippingFeeMain                         // simple assign
	out.ShippingFeeReturn = arg.ShippingFeeReturn                     // simple assign
	out.ShippingFeeInsurance = arg.ShippingFeeInsurance               // simple assign
	out.ShippingFeeAdjustment = arg.ShippingFeeAdjustment             // simple assign
	out.ShippingFeeCODS = arg.ShippingFeeCODS                         // simple assign
	out.ShippingFeeInfoChange = arg.ShippingFeeInfoChange             // simple assign
	out.ShippingFeeOther = arg.ShippingFeeOther                       // simple assign
	out.UpdatedBy = 0                                                 // zero value
	out.EtopAdjustedShippingFeeMain = arg.EtopAdjustedShippingFeeMain // simple assign
	out.EtopPriceRule = arg.EtopPriceRule                             // simple assign
	out.VariantIDs = arg.VariantIDs                                   // simple assign
	out.Lines = orderconvert.Convert_ordermodel_OrderLines_orderingmodel_OrderLines(arg.Lines)
	out.TypeFrom = ""                                                               // zero value
	out.TypeTo = ""                                                                 // zero value
	out.AddressFrom = arg.AddressFrom                                               // simple assign
	out.AddressTo = arg.AddressTo                                                   // simple assign
	out.AddressReturn = arg.AddressReturn                                           // simple assign
	out.AddressToProvinceCode = arg.AddressToProvinceCode                           // simple assign
	out.AddressToDistrictCode = arg.AddressToDistrictCode                           // simple assign
	out.AddressToWardCode = arg.AddressToWardCode                                   // simple assign
	out.AddressToPhone = ""                                                         // zero value
	out.AddressToFullNameNorm = ""                                                  // zero value
	out.CreatedAt = arg.CreatedAt                                                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                                                   // simple assign
	out.ClosedAt = arg.ClosedAt                                                     // simple assign
	out.ExpectedDeliveryAt = arg.ExpectedDeliveryAt                                 // simple assign
	out.ExpectedPickAt = arg.ExpectedPickAt                                         // simple assign
	out.CODEtopTransferedAt = arg.CODEtopTransferedAt                               // simple assign
	out.ShippingFeeShopTransferedAt = arg.ShippingFeeShopTransferedAt               // simple assign
	out.ShippingCancelledAt = arg.ShippingCancelledAt                               // simple assign
	out.ShippingDeliveredAt = arg.ShippingDeliveredAt                               // simple assign
	out.ShippingReturnedAt = arg.ShippingReturnedAt                                 // simple assign
	out.ShippingCreatedAt = arg.ShippingCreatedAt                                   // simple assign
	out.ShippingPickingAt = arg.ShippingPickingAt                                   // simple assign
	out.ShippingHoldingAt = arg.ShippingHoldingAt                                   // simple assign
	out.ShippingDeliveringAt = arg.ShippingDeliveringAt                             // simple assign
	out.ShippingReturningAt = arg.ShippingReturningAt                               // simple assign
	out.MoneyTransactionID = arg.MoneyTransactionID                                 // simple assign
	out.MoneyTransactionShippingExternalID = arg.MoneyTransactionShippingExternalID // simple assign
	out.CancelReason = arg.CancelReason                                             // simple assign
	out.ShippingProvider = arg.ShippingProvider                                     // simple assign
	out.ProviderServiceID = arg.ProviderServiceID                                   // simple assign
	out.ShippingCode = arg.ShippingCode                                             // simple assign
	out.ShippingNote = arg.ShippingNote                                             // simple assign
	out.TryOn = arg.TryOn                                                           // simple assign
	out.IncludeInsurance = dot.NullBool{}                                           // types do not match
	out.InsuranceValue = dot.NullInt{}                                              // zero value
	out.ShippingType = 0                                                            // zero value
	out.ShippingPaymentType = 0                                                     // zero value
	out.ConnectionID = arg.ConnectionID                                             // simple assign
	out.ConnectionMethod = arg.ConnectionMethod                                     // simple assign
	out.ShopCarrierID = 0                                                           // zero value
	out.ShippingServiceName = arg.ShippingServiceName                               // simple assign
	out.ExternalShippingName = arg.ExternalShippingName                             // simple assign
	out.ExternalShippingID = arg.ExternalShippingID                                 // simple assign
	out.ExternalShippingCode = arg.ExternalShippingCode                             // simple assign
	out.ExternalShippingCreatedAt = arg.ExternalShippingCreatedAt                   // simple assign
	out.ExternalShippingUpdatedAt = arg.ExternalShippingUpdatedAt                   // simple assign
	out.ExternalShippingCancelledAt = arg.ExternalShippingCancelledAt               // simple assign
	out.ExternalShippingDeliveredAt = arg.ExternalShippingDeliveredAt               // simple assign
	out.ExternalShippingReturnedAt = arg.ExternalShippingReturnedAt                 // simple assign
	out.ExternalShippingClosedAt = arg.ExternalShippingClosedAt                     // simple assign
	out.ExternalShippingState = arg.ExternalShippingState                           // simple assign
	out.ExternalShippingStateCode = arg.ExternalShippingStateCode                   // simple assign
	out.ExternalShippingStatus = arg.ExternalShippingStatus                         // simple assign
	out.ExternalShippingNote = dot.NullString{}                                     // types do not match
	out.ExternalShippingSubState = dot.NullString{}                                 // types do not match
	out.ExternalShippingData = arg.ExternalShippingData                             // simple assign
	out.ShippingState = arg.ShippingState                                           // simple assign
	out.ShippingStatus = arg.ShippingStatus                                         // simple assign
	out.EtopPaymentStatus = arg.EtopPaymentStatus                                   // simple assign
	out.Status = arg.Status                                                         // simple assign
	out.SyncStatus = arg.SyncStatus                                                 // simple assign
	out.SyncStates = arg.SyncStates                                                 // simple assign
	out.LastSyncAt = arg.LastSyncAt                                                 // simple assign
	out.ExternalShippingLogs = arg.ExternalShippingLogs                             // simple assign
	out.AdminNote = arg.AdminNote                                                   // simple assign
	out.IsPartialDelivery = arg.IsPartialDelivery                                   // simple assign
	out.CreatedBy = arg.CreatedBy                                                   // simple assign
	out.GrossWeight = arg.GrossWeight                                               // simple assign
	out.ChargeableWeight = arg.ChargeableWeight                                     // simple assign
	out.Length = arg.Length                                                         // simple assign
	out.Width = arg.Width                                                           // simple assign
	out.Height = arg.Height                                                         // simple assign
	out.DeliveryRoute = arg.DeliveryRoute                                           // simple assign
	out.ExternalAffiliateID = ""                                                    // zero value
	out.Coupon = ""                                                                 // zero value
	out.ShipmentPriceInfo = nil                                                     // zero value
	out.LinesContent = ""                                                           // zero value
	out.EdCode = ""                                                                 // zero value
	out.ShippingSubstate = substate.NullSubstate{}                                  // zero value
	out.Rid = arg.Rid                                                               // simple assign
}

func Convert_fulfillmentmodel_Fulfillments_shippingmodel_Fulfillments(args []*fulfillmentmodel.Fulfillment) (outs []*shippingmodel.Fulfillment) {
	if args == nil {
		return nil
	}
	tmps := make([]shippingmodel.Fulfillment, len(args))
	outs = make([]*shippingmodel.Fulfillment, len(args))
	for i := range tmps {
		outs[i] = Convert_fulfillmentmodel_Fulfillment_shippingmodel_Fulfillment(args[i], &tmps[i])
	}
	return outs
}

func Convert_shippingmodel_Fulfillment_fulfillmentmodel_Fulfillment(arg *shippingmodel.Fulfillment, out *fulfillmentmodel.Fulfillment) *fulfillmentmodel.Fulfillment {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &fulfillmentmodel.Fulfillment{}
	}
	convert_shippingmodel_Fulfillment_fulfillmentmodel_Fulfillment(arg, out)
	return out
}

func convert_shippingmodel_Fulfillment_fulfillmentmodel_Fulfillment(arg *shippingmodel.Fulfillment, out *fulfillmentmodel.Fulfillment) {
	out.ID = arg.ID                                                   // simple assign
	out.OrderID = arg.OrderID                                         // simple assign
	out.ShopID = arg.ShopID                                           // simple assign
	out.PartnerID = arg.PartnerID                                     // simple assign
	out.ShopConfirm = arg.ShopConfirm                                 // simple assign
	out.ConfirmStatus = arg.ConfirmStatus                             // simple assign
	out.TotalItems = arg.TotalItems                                   // simple assign
	out.TotalWeight = arg.TotalWeight                                 // simple assign
	out.BasketValue = arg.BasketValue                                 // simple assign
	out.TotalDiscount = arg.TotalDiscount                             // simple assign
	out.TotalAmount = arg.TotalAmount                                 // simple assign
	out.TotalCODAmount = arg.TotalCODAmount                           // simple assign
	out.OriginalCODAmount = arg.OriginalCODAmount                     // simple assign
	out.ActualCompensationAmount = arg.ActualCompensationAmount       // simple assign
	out.ShippingFeeCustomer = arg.ShippingFeeCustomer                 // simple assign
	out.ShippingFeeShop = arg.ShippingFeeShop                         // simple assign
	out.ShippingFeeShopLines = arg.ShippingFeeShopLines               // simple assign
	out.ShippingServiceFee = arg.ShippingServiceFee                   // simple assign
	out.ExternalShippingFee = arg.ExternalShippingFee                 // simple assign
	out.ProviderShippingFeeLines = arg.ProviderShippingFeeLines       // simple assign
	out.EtopDiscount = arg.EtopDiscount                               // simple assign
	out.EtopFeeAdjustment = arg.EtopFeeAdjustment                     // simple assign
	out.ShippingFeeMain = arg.ShippingFeeMain                         // simple assign
	out.ShippingFeeReturn = arg.ShippingFeeReturn                     // simple assign
	out.ShippingFeeInsurance = arg.ShippingFeeInsurance               // simple assign
	out.ShippingFeeAdjustment = arg.ShippingFeeAdjustment             // simple assign
	out.ShippingFeeCODS = arg.ShippingFeeCODS                         // simple assign
	out.ShippingFeeInfoChange = arg.ShippingFeeInfoChange             // simple assign
	out.ShippingFeeOther = arg.ShippingFeeOther                       // simple assign
	out.EtopAdjustedShippingFeeMain = arg.EtopAdjustedShippingFeeMain // simple assign
	out.EtopPriceRule = arg.EtopPriceRule                             // simple assign
	out.VariantIDs = arg.VariantIDs                                   // simple assign
	out.Lines = orderconvert.Convert_orderingmodel_OrderLines_ordermodel_OrderLines(arg.Lines)
	out.AddressFrom = arg.AddressFrom                                               // simple assign
	out.AddressTo = arg.AddressTo                                                   // simple assign
	out.AddressReturn = arg.AddressReturn                                           // simple assign
	out.AddressToProvinceCode = arg.AddressToProvinceCode                           // simple assign
	out.AddressToDistrictCode = arg.AddressToDistrictCode                           // simple assign
	out.AddressToWardCode = arg.AddressToWardCode                                   // simple assign
	out.CreatedAt = arg.CreatedAt                                                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                                                   // simple assign
	out.ClosedAt = arg.ClosedAt                                                     // simple assign
	out.ExpectedDeliveryAt = arg.ExpectedDeliveryAt                                 // simple assign
	out.ExpectedPickAt = arg.ExpectedPickAt                                         // simple assign
	out.CODEtopTransferedAt = arg.CODEtopTransferedAt                               // simple assign
	out.ShippingFeeShopTransferedAt = arg.ShippingFeeShopTransferedAt               // simple assign
	out.ShippingCancelledAt = arg.ShippingCancelledAt                               // simple assign
	out.ShippingDeliveredAt = arg.ShippingDeliveredAt                               // simple assign
	out.ShippingReturnedAt = arg.ShippingReturnedAt                                 // simple assign
	out.ShippingCreatedAt = arg.ShippingCreatedAt                                   // simple assign
	out.ShippingPickingAt = arg.ShippingPickingAt                                   // simple assign
	out.ShippingHoldingAt = arg.ShippingHoldingAt                                   // simple assign
	out.ShippingDeliveringAt = arg.ShippingDeliveringAt                             // simple assign
	out.ShippingReturningAt = arg.ShippingReturningAt                               // simple assign
	out.MoneyTransactionID = arg.MoneyTransactionID                                 // simple assign
	out.MoneyTransactionShippingExternalID = arg.MoneyTransactionShippingExternalID // simple assign
	out.CancelReason = arg.CancelReason                                             // simple assign
	out.ShippingProvider = arg.ShippingProvider                                     // simple assign
	out.ProviderServiceID = arg.ProviderServiceID                                   // simple assign
	out.ShippingCode = arg.ShippingCode                                             // simple assign
	out.ShippingNote = arg.ShippingNote                                             // simple assign
	out.TryOn = arg.TryOn                                                           // simple assign
	out.IncludeInsurance = false                                                    // types do not match
	out.ConnectionID = arg.ConnectionID                                             // simple assign
	out.ConnectionMethod = arg.ConnectionMethod                                     // simple assign
	out.ShippingServiceName = arg.ShippingServiceName                               // simple assign
	out.ExternalShippingName = arg.ExternalShippingName                             // simple assign
	out.ExternalShippingID = arg.ExternalShippingID                                 // simple assign
	out.ExternalShippingCode = arg.ExternalShippingCode                             // simple assign
	out.ExternalShippingCreatedAt = arg.ExternalShippingCreatedAt                   // simple assign
	out.ExternalShippingUpdatedAt = arg.ExternalShippingUpdatedAt                   // simple assign
	out.ExternalShippingCancelledAt = arg.ExternalShippingCancelledAt               // simple assign
	out.ExternalShippingDeliveredAt = arg.ExternalShippingDeliveredAt               // simple assign
	out.ExternalShippingReturnedAt = arg.ExternalShippingReturnedAt                 // simple assign
	out.ExternalShippingClosedAt = arg.ExternalShippingClosedAt                     // simple assign
	out.ExternalShippingState = arg.ExternalShippingState                           // simple assign
	out.ExternalShippingStateCode = arg.ExternalShippingStateCode                   // simple assign
	out.ExternalShippingStatus = arg.ExternalShippingStatus                         // simple assign
	out.ExternalShippingNote = ""                                                   // types do not match
	out.ExternalShippingSubState = ""                                               // types do not match
	out.ExternalShippingData = arg.ExternalShippingData                             // simple assign
	out.ShippingState = arg.ShippingState                                           // simple assign
	out.ShippingStatus = arg.ShippingStatus                                         // simple assign
	out.EtopPaymentStatus = arg.EtopPaymentStatus                                   // simple assign
	out.Status = arg.Status                                                         // simple assign
	out.SyncStatus = arg.SyncStatus                                                 // simple assign
	out.SyncStates = arg.SyncStates                                                 // simple assign
	out.LastSyncAt = arg.LastSyncAt                                                 // simple assign
	out.ExternalShippingLogs = arg.ExternalShippingLogs                             // simple assign
	out.AdminNote = arg.AdminNote                                                   // simple assign
	out.IsPartialDelivery = arg.IsPartialDelivery                                   // simple assign
	out.CreatedBy = arg.CreatedBy                                                   // simple assign
	out.GrossWeight = arg.GrossWeight                                               // simple assign
	out.ChargeableWeight = arg.ChargeableWeight                                     // simple assign
	out.Length = arg.Length                                                         // simple assign
	out.Width = arg.Width                                                           // simple assign
	out.Height = arg.Height                                                         // simple assign
	out.DeliveryRoute = arg.DeliveryRoute                                           // simple assign
	out.Rid = arg.Rid                                                               // simple assign
}

func Convert_shippingmodel_Fulfillments_fulfillmentmodel_Fulfillments(args []*shippingmodel.Fulfillment) (outs []*fulfillmentmodel.Fulfillment) {
	if args == nil {
		return nil
	}
	tmps := make([]fulfillmentmodel.Fulfillment, len(args))
	outs = make([]*fulfillmentmodel.Fulfillment, len(args))
	for i := range tmps {
		outs[i] = Convert_shippingmodel_Fulfillment_fulfillmentmodel_Fulfillment(args[i], &tmps[i])
	}
	return outs
}
