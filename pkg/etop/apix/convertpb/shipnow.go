package convertpb

import (
	"o.o/api/main/connectioning"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	shipnowtypes "o.o/api/main/shipnow/types"
	typesx "o.o/api/top/external/types"
	typesint "o.o/api/top/int/types"
	"o.o/capi/dot"
)

func Convert_core_ShipnowService_To_apix_ShipnowService(in *shipnowtypes.ShipnowService) *typesx.ShipnowService {
	if in == nil {
		return nil
	}
	res := &typesx.ShipnowService{
		Code:        in.Code,
		Name:        in.Name,
		Fee:         in.Fee,
		Description: in.Description,
	}
	if in.ConnectionInfo != nil {
		res.CarrierInfo = &typesx.CarrierInfo{
			Name:     in.ConnectionInfo.Name,
			ImageURL: in.ConnectionInfo.ImageURL,
		}
	}
	return res
}

func Convert_core_ShipnowServices_To_apix_ShipnowServices(items []*shipnowtypes.ShipnowService) []*typesx.ShipnowService {
	result := make([]*typesx.ShipnowService, len(items))
	for i, item := range items {
		result[i] = Convert_core_ShipnowService_To_apix_ShipnowService(item)
	}
	return result
}

func Convert_apix_ShipnowLocationAddressShortVersion_To_core_OrderAddress(in *typesx.ShipnowAddressShortVersion) *ordertypes.Address {
	if in == nil {
		return nil
	}
	res := &ordertypes.Address{
		Address1: in.Address,
		Location: ordertypes.Location{
			Province: in.Province,
			District: in.District,
			Ward:     in.Ward,
		},
	}
	if in.Coordinates != nil {
		res.Coordinates = &ordertypes.Coordinates{
			Latitude:  in.Coordinates.Latitude,
			Longitude: in.Coordinates.Longitude,
		}
	}
	return res
}

func Convert_apix_ShipnowAddress_To_api_OrderAddress(in *typesx.ShipnowAddress) *typesint.OrderAddress {
	if in == nil {
		return nil
	}
	return &typesint.OrderAddress{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Province: in.Province,
		District: in.District,
		Ward:     in.Ward,
		Company:  in.Company,
		Address1: in.Address,
	}
}

func Convert_apix_ShipnowAddress_To_core_OrderAddress(in *typesx.ShipnowAddress) *ordertypes.Address {
	if in == nil {
		return nil
	}
	res := &ordertypes.Address{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Company:  in.Company,
		Address1: in.Address,
		Location: ordertypes.Location{
			Province: in.Province,
			District: in.District,
			Ward:     in.Ward,
		},
	}
	if in.Coordinates != nil {
		res.Location.Coordinates = &ordertypes.Coordinates{
			Latitude:  in.Coordinates.Latitude,
			Longitude: in.Coordinates.Longitude,
		}
	}
	return res
}

func Convert_core_OrderAddress_To_apix_ShipnowAddress(in *ordertypes.Address) *typesx.ShipnowAddress {
	if in == nil {
		return nil
	}
	res := &typesx.ShipnowAddress{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Province: in.Province,
		District: in.District,
		Ward:     in.Ward,
		Address:  in.Address1,
		Company:  in.Company,
	}
	if in.Coordinates != nil {
		res.Coordinates = &typesx.Coordinates{
			Latitude:  in.Coordinates.Latitude,
			Longitude: in.Coordinates.Longitude,
		}
	}
	return res
}

func Convert_core_OrderLine_To_apix_OrderLine(in *ordertypes.ItemLine) *typesx.OrderLine {
	if in == nil {
		return nil
	}
	return &typesx.OrderLine{
		VariantId:    in.VariantID,
		ProductId:    in.ProductID,
		ProductName:  in.ProductInfo.ProductName,
		Quantity:     in.Quantity,
		ListPrice:    dot.Int(in.ProductInfo.ListPrice),
		RetailPrice:  dot.Int(in.ProductInfo.RetailPrice),
		PaymentPrice: dot.Int(in.ProductInfo.PaymentPrice),
		ImageUrl:     in.ProductInfo.ImageURL,
		Attributes:   in.ProductInfo.Attributes,
	}
}

func Convert_core_OrderLines_To_apix_OrderLines(items []*ordertypes.ItemLine) []*typesx.OrderLine {
	result := make([]*typesx.OrderLine, len(items))
	for i, item := range items {
		result[i] = Convert_core_OrderLine_To_apix_OrderLine(item)
	}
	return result
}

func Convert_core_DeliveryPoint_To_apix_ShipnowDeliveryPoint(in *shipnow.DeliveryPoint) *typesx.ShipnowDeliveryPoint {
	if in == nil {
		return nil
	}
	return &typesx.ShipnowDeliveryPoint{
		ChargeableWeight: typesint.Int(in.WeightInfo.ChargeableWeight),
		GrossWeight:      typesint.Int(in.WeightInfo.GrossWeight),
		CODAmount:        typesint.Int(in.ValueInfo.CODAmount),
		ShippingNote:     in.ShippingNote,
		ShippingAddress:  Convert_core_OrderAddress_To_apix_ShipnowAddress(in.ShippingAddress),
		BasketValue:      typesint.Int(in.BasketValue),
		Lines:            Convert_core_OrderLines_To_apix_OrderLines(in.Lines),
		ShippingState:    in.ShippingState,
	}
}

func Convert_core_DeliveryPoints_To_apix_ShipnowDeliveryPoints(items []*shipnow.DeliveryPoint) []*typesx.ShipnowDeliveryPoint {
	result := make([]*typesx.ShipnowDeliveryPoint, len(items))
	for i, item := range items {
		result[i] = Convert_core_DeliveryPoint_To_apix_ShipnowDeliveryPoint(item)
	}
	return result
}

func Convert_core_ShipnowFulfillment_To_apix_ShipnowFulfillment(in *shipnow.ShipnowFulfillment, conn *connectioning.Connection) *typesx.ShipnowFulfillment {
	if in == nil {
		return nil
	}
	res := &typesx.ShipnowFulfillment{
		ID:                         in.ID,
		ShopID:                     in.ShopID,
		PickupAddress:              Convert_core_OrderAddress_To_apix_ShipnowAddress(in.PickupAddress),
		DeliveryPoints:             Convert_core_DeliveryPoints_To_apix_ShipnowDeliveryPoints(in.DeliveryPoints),
		ShippingServiceCode:        dot.String(in.ShippingServiceCode),
		ShippingServiceFee:         dot.Int(in.ShippingServiceFee),
		ActualShippingServiceFee:   dot.Int(in.TotalFee),
		ShippingServiceName:        dot.String(in.ShippingServiceName),
		ShippingServiceDescription: dot.String(in.ShippingServiceDescription),
		GrossWeight:                dot.Int(in.GrossWeight),
		ChargeableWeight:           dot.Int(in.ChargeableWeight),
		ShippingNote:               dot.String(in.ShippingNote),
		Status:                     in.Status.Wrap(),
		ShippingStatus:             in.ShippingStatus.Wrap(),
		ShippingCode:               dot.String(in.ShippingCode),
		ShippingState:              in.ShippingState.Wrap(),
		ConfirmStatus:              in.ConfirmStatus,
		OrderIDs:                   in.OrderIDs,
		CreatedAt:                  in.CreatedAt,
		UpdatedAt:                  in.UpdatedAt,
		ShippingSharedLink:         dot.String(in.ShippingSharedLink),
		CancelReason:               dot.String(in.CancelReason),
		CODAmount:                  dot.Int(in.ValueInfo.CODAmount),
		BasketValue:                dot.Int(in.ValueInfo.BasketValue),
		ExternalID:                 dot.String(in.ExternalID),
	}
	if conn != nil {
		res.CarrierInfo = &typesx.CarrierInfo{
			Name:     conn.Name,
			ImageURL: conn.ImageURL,
		}
	}
	return res
}
