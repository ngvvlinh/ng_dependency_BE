package convertpb

import (
	"o.o/api/main/accountshipnow"
	"o.o/api/main/shipnow"
	shipnowtypes "o.o/api/main/shipnow/types"
	"o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func ShipnowNullState(s dot.NullString) shipnow_state.NullState {
	if s.Apply("") == "" {
		return shipnow_state.NullState{}
	}
	st, ok := shipnow_state.ParseState(s.String)
	if !ok {
		return shipnow_state.StateDefault.Wrap()
	}
	return st.Wrap()
}

func Convert_core_XAccountAhamove_To_api_XAccountAhamove(in *accountshipnow.ExternalAccountAhamove, hideInfo bool) *shop.ExternalAccountAhamove {
	if in == nil {
		return nil
	}
	if hideInfo {
		return &shop.ExternalAccountAhamove{
			Id: in.ID,
		}
	}
	return &shop.ExternalAccountAhamove{
		Id:                  in.ID,
		Phone:               in.Phone,
		Name:                in.Name,
		ExternalVerified:    in.ExternalVerified,
		ExternalCreatedAt:   cmapi.PbTime(in.ExternalCreatedAt),
		CreatedAt:           cmapi.PbTime(in.CreatedAt),
		UpdatedAt:           cmapi.PbTime(in.UpdatedAt),
		LastSendVerifyAt:    cmapi.PbTime(in.LastSendVerifiedAt),
		ExternalTicketId:    in.ExternalTicketID,
		IdCardFrontImg:      in.IDCardFrontImg,
		IdCardBackImg:       in.IDCardBackImg,
		PortraitImg:         in.PortraitImg,
		UploadedAt:          cmapi.PbTime(in.UploadedAt),
		FanpageUrl:          in.FanpageURL,
		WebsiteUrl:          in.WebsiteURL,
		CompanyImgs:         in.CompanyImgs,
		BusinessLicenseImgs: in.BusinessLicenseImgs,
		OwnerID:             in.OwnerID,
		ConnectionID:        in.ConnectionID,
	}
}

func Convert_core_ShipnowService_To_api_ShipnowService(in *shipnowtypes.ShipnowService) *types.ShipnowService {
	if in == nil {
		return nil
	}
	return &types.ShipnowService{
		Carrier:        in.Carrier.String(),
		Name:           in.Name,
		Code:           in.Code,
		Fee:            in.Fee,
		Description:    in.Description,
		ConnectionInfo: Convert_core_ConnectionInfo_To_api_ConnectionInfo(in.ConnectionInfo),
	}
}

func Convert_core_ShipnowServices_To_api_ShipnowServices(ins []*shipnowtypes.ShipnowService) (outs []*types.ShipnowService) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowService_To_api_ShipnowService(in))
	}
	return
}

func Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in *shipnow.ShipnowFulfillment) *types.ShipnowFulfillment {
	if in == nil {
		return nil
	}
	return &types.ShipnowFulfillment{
		Id:                         in.ID,
		ShopId:                     in.ShopID,
		PartnerId:                  in.PartnerID,
		PickupAddress:              Convert_core_OrderAddress_To_api_OrderAddress(in.PickupAddress),
		DeliveryPoints:             Convert_core_DeliveryPoints_To_api_DeliveryPoints(in.DeliveryPoints),
		Carrier:                    in.Carrier.String(),
		ShippingServiceCode:        in.ShippingServiceCode,
		ShippingServiceFee:         in.ShippingServiceFee,
		ShippingServiceName:        in.ShippingServiceName,
		ShippingServiceDescription: in.ShippingServiceDescription,
		WeightInfo:                 Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:                  Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		ShippingNote:               in.ShippingNote,
		RequestPickupAt:            cmapi.PbTime(in.RequestPickupAt),
		CreatedAt:                  cmapi.PbTime(in.CreatedAt),
		UpdatedAt:                  cmapi.PbTime(in.UpdatedAt),
		Status:                     in.Status,
		ShippingStatus:             in.ShippingStatus,
		ShippingState:              in.ShippingState,
		ConfirmStatus:              in.ConfirmStatus,
		OrderIds:                   in.OrderIDs,
		ShippingCreatedAt:          cmapi.PbTime(in.ShippingCreatedAt),
		ShippingCode:               in.ShippingCode,
		EtopPaymentStatus:          in.EtopPaymentStatus,
		CodEtopTransferedAt:        cmapi.PbTime(in.CODEtopTransferedAt),
		ShippingPickingAt:          cmapi.PbTime(in.ShippingPickingAt),
		ShippingDeliveringAt:       cmapi.PbTime(in.ShippingDeliveringAt),
		ShippingDeliveredAt:        cmapi.PbTime(in.ShippingDeliveredAt),
		ShippingCancelledAt:        cmapi.PbTime(in.ShippingCancelledAt),
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
		ConnectionID:               in.ConnectionID,
		Coupon:                     in.Coupon,
	}
}

func Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ins []*shipnow.ShipnowFulfillment) (outs []*types.ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in))
	}
	return
}

func Convert_api_DeliveryPoint_To_core_DeliveryPoint(in *types.DeliveryPoint) *shipnow.DeliveryPoint {
	return &shipnow.DeliveryPoint{
		ShippingAddress: Convert_api_OrderAddress_To_core_OrderAddress(in.ShippingAddress),
		Lines:           Convert_api_OrderLines_To_core_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_api_WeightInfo_To_core_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_api_ValueInfo_To_core_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Conver_api_DeliveryPoints_To_core_DeliveryPoints(ins []*types.DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_api_DeliveryPoint_To_core_DeliveryPoint(in))
	}
	return
}

func Convert_core_DeliveryPoint_To_api_DeliveryPoint(in *shipnow.DeliveryPoint) *types.DeliveryPoint {
	return &types.DeliveryPoint{
		ShippingAddress: Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		Lines:           Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Convert_core_DeliveryPoints_To_api_DeliveryPoints(ins []*shipnow.DeliveryPoint) (outs []*types.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_core_DeliveryPoint_To_api_DeliveryPoint(in))
	}
	return
}
