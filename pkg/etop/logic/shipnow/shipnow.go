package shipnow

import (
	"context"

	"etop.vn/backend/pkg/etop/model"

	"etop.vn/backend/pkg/common/bus"

	cm "etop.vn/backend/pkg/common"

	ordermodelx "etop.vn/backend/pkg/services/ordering/modelx"
	shipnowmodelx "etop.vn/backend/pkg/services/shipnow/modelx"
)

func CreateShipNowFulfillment(ctx context.Context, args shipnowmodelx.CreateShipNowFulfillmentArgs) error {
	orderIDs := args.OrderIDs
	if len(orderIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Order IDs")
	}
	query := ordermodelx.GetOrdersQuery{
		ShopIDs: []int64{args.ShopID},
		IDs:     args.OrderIDs,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	PickupAddressID := args.ShopAddressID
	// get pickup address
	if args.ShopAddressID == 0 {
		q := model.GetShopQuery{
			ShopID: args.ShopID,
		}
		if err := bus.Dispatch(ctx, q); err != nil {
			return err
		}
		PickupAddressID = q.Result.AddressID
	}
	addressQuery := &model.GetAddressQuery{
		AddressID: PickupAddressID,
	}
	if err := bus.Dispatch(ctx, addressQuery); err != nil {
		return cm.Error(cm.Internal, "Lỗi khi kiểm tra thông tin địa chỉ của cửa hàng: "+err.Error(), err)
	}
	// shopAddress := addressQuery.Result
	//
	// cmd := &shipnow.CreateShipNowFulfillmentCommand{
	// 	ShipNowFulfillment: shipnow.ShipNowFulfillment{
	// 		ShopId:              args.ShopID,
	// 		PickupAddress:       nil,
	// 		DeliveryPoints:      nil,
	// 		Carrier:             args.Carrier,
	// 		ShippingServiceCode: args.ShippingServiceCode,
	// 		ShippingServiceFee:  int32(args.ShippingServiceFee),
	// 		WeightInfo:          types.WeightInfo{},
	// 		ValueInfo:           types.ValueInfo{},
	// 		ShippingNote:        args.ShippingNote,
	// 		RequestPickupAt:     nil,
	// 	},
	// }
	return nil
}
