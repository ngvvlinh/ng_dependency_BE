package aggregate

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"o.o/api/main/address"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	shipstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	com "o.o/backend/com/main"
	addressconvert "o.o/backend/com/main/address/convert"
	orderconvert "o.o/backend/com/main/ordering/convert"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	"o.o/backend/com/main/shipping/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/backend/com/main/shipping/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ shipping.Aggregate = &Aggregate{}
var ll = l.New()
var scheme = conversion.Build(shippingconvert.RegisterConversions)

type Aggregate struct {
	db             *cmsql.Database
	locationQS     location.QueryBus
	orderQS        ordering.QueryBus
	shimentManager *carrier.ShipmentManager
	connectionQS   connectioning.QueryBus
	identityQS     identity.QueryBus
	addressQS      address.QueryBus
	ffmStore       sqlstore.FulfillmentStoreFactory
	eventBus       capi.EventBus
}

func NewAggregate(
	db com.MainDB, eventB capi.EventBus,
	locationQS location.QueryBus, orderQS ordering.QueryBus,
	shipmentManager *carrier.ShipmentManager, connectionQS connectioning.QueryBus,
	identityQS identity.QueryBus, addressQS address.QueryBus,
) *Aggregate {
	return &Aggregate{
		db:             db,
		locationQS:     locationQS,
		orderQS:        orderQS,
		shimentManager: shipmentManager,
		connectionQS:   connectionQS,
		identityQS:     identityQS,
		addressQS:      addressQS,
		ffmStore:       sqlstore.NewFulfillmentStore(db),
		eventBus:       eventB,
	}
}

func AggregateMessageBus(a *Aggregate) shipping.CommandBus {
	b := bus.New()
	return shipping.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateFulfillments(ctx context.Context, args *shipping.CreateFulfillmentsArgs) (fulfillmentIDs []dot.ID, _ error) {
	query := &ordering.GetOrderByIDQuery{
		ID: args.OrderID,
	}
	if err := a.orderQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	order := query.Result
	switch order.Status {
	case status5.N:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hủy")
	case status5.NS:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng")
	}

	if order.ConfirmStatus != status3.P || order.ShopConfirm != status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng xác nhận đơn hàng trước khi tạo đơn giao hàng")
	}

	var ffms []*shipmodel.Fulfillment
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		event := &shipping.FulfillmentsCreatingEvent{
			EventMeta: meta.NewEvent(),
			ShopID:    args.ShopID,
			OrderID:   args.OrderID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		_, _, err := a.getAndVerifyAddress(ctx, args.ShippingAddress)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
		}
		_, _, err = a.getAndVerifyAddress(ctx, args.PickupAddress)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
		}
		if args.ReturnAddress != nil {
			_, _, err = a.getAndVerifyAddress(ctx, args.ReturnAddress)
			if err != nil {
				return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ trả hàng không hợp lệ: %v", err)
			}
		}

		oldFulfillments, err := a.ffmStore(ctx).OrderID(args.OrderID).ListFfmsDB()
		if err != nil {
			return err
		}

		ffm, err := a.prepareFulfillmentFromOrder(ctx, order, args)
		if err != nil {
			return err
		}

		creates, updates, err := CompareFulfillments(oldFulfillments, ffm)
		if err != nil {
			return err
		}
		if creates != nil {
			if err := a.ffmStore(ctx).CreateFulfillmentsDB(creates); err != nil {
				return err
			}
		}
		if updates != nil {
			if err := a.ffmStore(ctx).StatusNotIn(status5.N, status5.NS, status5.P).UpdateFulfillmentsDB(updates); err != nil {
				return err
			}
		}
		ll.S.Infof("Compare fulfillments: create %v update %v", len(creates), len(updates))
		totalChanges := len(creates) + len(updates)
		if totalChanges == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "Không tạo được fulfillment")
		}

		ffms = append(creates, updates...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := a.shimentManager.CreateFulfillments(ctx, ffms); err != nil {
		return nil, err
	}

	var res []dot.ID
	for _, ffm := range ffms {
		res = append(res, ffm.ID)
	}

	// TODO: kiểm tra trường hợp tự giao hoặc giao qua NVC
	// rollback khi gặp lỗi. VD: hủy ffm bên NVC
	event2 := &shipping.FulfillmentsCreatedEvent{
		EventMeta:      meta.NewEvent(),
		FulfillmentIDs: res,
		ShippingType:   args.ShippingType,
		OrderID:        order.ID,
	}
	if err := a.eventBus.Publish(ctx, event2); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *Aggregate) getAndVerifyAddress(ctx context.Context, address *ordertypes.Address) (*location.Province, *location.District, error) {
	if address == nil {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin địa chỉ")
	}
	if address.Phone == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin số điện thoại")
	}
	if _, ok := validate.NormalizePhone(address.Phone); !ok {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ (%v).", address.Phone)
	}

	query := &location.FindOrGetLocationQuery{
		Province:     address.Province,
		District:     address.District,
		Ward:         address.Ward,
		ProvinceCode: address.ProvinceCode,
		DistrictCode: address.DistrictCode,
		WardCode:     address.WardCode,
	}
	if err := a.locationQS.Dispatch(ctx, query); err != nil {
		return nil, nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ không hợp lệ %v", err.Error())
	}
	district, province := query.Result.District, query.Result.Province

	address.Location = ordertypes.Location{
		ProvinceCode: province.Code,
		Province:     province.Name,
		DistrictCode: district.Code,
		District:     district.Name,
		Coordinates:  address.Coordinates,
	}
	ward := query.Result.Ward
	if ward != nil {
		address.Location.Ward = ward.Name
		address.Location.WardCode = ward.Code
	}
	return province, district, nil
}

func (a *Aggregate) prepareFulfillmentFromOrder(ctx context.Context, order *ordering.Order, args *shipping.CreateFulfillmentsArgs) (*shipmodel.Fulfillment, error) {
	shippingType := args.ShippingType
	var connectionMethod connection_type.ConnectionMethod
	var conn *connectioning.Connection

	switch shippingType {
	case ordertypes.ShippingTypeManual:
		if args.ShopCarrierID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển tự quản lý")
		}
	case ordertypes.ShippingTypeShipment:
		if args.ConnectionID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển (connection_id)")
		}
		queryConn := &connectioning.GetConnectionByIDQuery{
			ID: args.ConnectionID,
		}
		if err := a.connectionQS.Dispatch(ctx, queryConn); err != nil {
			return nil, err
		}
		conn = queryConn.Result

		if conn.ConnectionProvider == connection_type.ConnectionProviderGHN {
			if args.TryOn.String() == "" {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng chọn ghi chú xem hàng!")
			}
		}
		connectionMethod = conn.ConnectionMethod

		if err := checkValidShippingPaymentType(connectionMethod, args.ShippingPaymentType); err != nil {
			return nil, err
		}
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phương thức vận chuyển không hợp lệ. Chỉ hỗ trợ: manual và shipment.")
	}

	tryOn := args.TryOn
	if tryOn == 0 {
		tryOn = try_on.None
	}

	variantIDs := []dot.ID{}
	totalItems, basketValue, totalAmount := 0, 0, 0

	if len(order.Lines) != 0 {
		for _, line := range order.Lines {
			variantIDs = append(variantIDs, line.VariantID)
		}
	}
	totalItems = order.TotalItems
	basketValue = order.BasketValue
	totalAmount = order.TotalAmount
	variantIDs = variantIDs

	typeFrom := etopmodel.FFShop
	typeTo := etopmodel.FFCustomer

	ffm := &shipmodel.Fulfillment{
		ID:                cm.NewID(),
		OrderID:           order.ID,
		ShopID:            order.ShopID,
		PartnerID:         order.PartnerID,
		ShopConfirm:       status3.P, // Always set shop_confirm to 1
		ConfirmStatus:     0,
		TotalItems:        totalItems,
		TotalWeight:       cm.CoalesceInt(args.ChargeableWeight, args.GrossWeight, order.TotalWeight),
		BasketValue:       basketValue,
		TotalDiscount:     0,
		TotalAmount:       totalAmount,
		TotalCODAmount:    args.CODAmount,
		OriginalCODAmount: order.ShopCOD,
		EtopDiscount:      0,
		EtopFeeAdjustment: 0,
		VariantIDs:        variantIDs,
		Lines:             orderconvert.OrderLinesToModel(order.Lines),
		// after this
		TypeFrom:            typeFrom,
		TypeTo:              typeTo,
		AddressFrom:         addressconvert.OrderAddressToModel(args.PickupAddress),
		AddressTo:           addressconvert.OrderAddressToModel(args.ShippingAddress),
		AddressReturn:       addressconvert.OrderAddressToModel(args.ReturnAddress),
		ShippingFeeCustomer: order.ShopShippingFee,
		ProviderServiceID:   args.ShippingServiceCode,
		ShippingServiceFee:  args.ShippingServiceFee,
		ShippingServiceName: args.ShippingServiceName,
		ShippingNote:        args.ShippingNote,
		TryOn:               tryOn,
		IncludeInsurance:    dot.Bool(args.IncludeInsurance),
		InsuranceValue:      args.InsuranceValue,
		ConnectionID:        args.ConnectionID,
		ConnectionMethod:    connectionMethod,
		ShopCarrierID:       args.ShopCarrierID,
		GrossWeight:         args.GrossWeight,
		ChargeableWeight:    args.ChargeableWeight,
		Width:               args.Width,
		Height:              args.Height,
		Length:              args.Length,
		ShippingState:       shipstate.Default,
		ShippingType:        args.ShippingType,
		Coupon:              args.Coupon,
		ShippingPaymentType: args.ShippingPaymentType.Apply(shipping_payment_type.Seller),
		CreatedBy:           order.CreatedBy,
	}
	if conn != nil {
		// backward compatible
		if shippingProvider, ok := shipping_provider.ParseShippingProvider(conn.ConnectionProvider.String()); ok {
			ffm.ShippingProvider = shippingProvider
		}
	}
	return ffm, nil
}

// Compare the old fulfillments and expected fulfillments
// - Missing/cancelled fulfillments: create
// - Error fulfillments: update
// - Processing fulfillments: ignore
func CompareFulfillments(olds []*shipmodel.Fulfillment, ffm *shipmodel.Fulfillment) (creates, updates []*shipmodel.Fulfillment, err error) {
	// active ffm: Those which are not cancelled
	var old *shipmodel.Fulfillment
	for _, oldFfm := range olds {
		if oldFfm.Status != status5.N && oldFfm.Status != status5.NS &&
			oldFfm.ShopConfirm != status3.N {
			old = oldFfm
		}
	}

	if old == nil {
		return []*shipmodel.Fulfillment{ffm}, nil, nil
	}
	// update error fulfillments
	if old.Status == status5.Z {
		ffm.ID = old.ID
		return nil, []*shipmodel.Fulfillment{ffm}, nil
	}
	// ignore
	return nil, nil, nil
}

func (a *Aggregate) UpdateFulfillmentShippingState(ctx context.Context, args *shipping.UpdateFulfillmentShippingStateArgs) (updated int, _ error) {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	if args.ShippingState == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "shipping_state không được để trống.")
	}

	query := a.ffmStore(ctx).OptionalPartnerID(args.PartnerID).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode)
	ffm, err := query.GetFulfillment()
	if err != nil {
		return 0, err
	}

	if err := canUpdateFulfillment(ffm); err != nil {
		return 0, err
	}
	if args.ActualCompensationAmount.Valid {
		if ffm.ShippingState != shipstate.Undeliverable && args.ShippingState != shipstate.Undeliverable {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Chỉ cập nhật giá trị hoàn hàng khi đơn vận chuyển không giao được hàng")
		}
	} else {
		if args.ShippingState == shipstate.Undeliverable {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập giá trị hoàn hàng.")
		}
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		args.FulfillmentID = ffm.ID
		if err := a.ffmStore(ctx).UpdateFulfillmentShippingState(args, ffm.CODAmount); err != nil {
			return err
		}

		if shipping.IsStateReturn(args.ShippingState) {
			// Nếu là đơn trả hàng
			// Cộng thêm phí trả hàng
			// Nếu đã có phí trả hàng thì chỉ update giá theo bảng giá hiện tại
			calcShippingFeesArgs := &carrier.CalcMakeupShippingFeesByFfmArgs{
				Fulfillment: ffm,
				State:       args.ShippingState,
			}
			calcFeeResp, err := a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, calcShippingFeesArgs)
			if err != nil {
				return err
			}
			returnFeeLine := shippingtypes.GetShippingFeeLine(calcFeeResp.ShippingFeeLines, shipping_fee_type.Return)
			shippingFeeShopLines := shippingtypes.ApplyShippingFeeLine(ffm.ShippingFeeShopLines, returnFeeLine)

			update := &shipping.UpdateFulfillmentShippingFeesArgs{
				FulfillmentID:    ffm.ID,
				ShippingCode:     ffm.ShippingCode,
				ShippingFeeLines: shippingFeeShopLines,
				UpdatedBy:        args.UpdatedBy,
				ShipmentPriceInfo: &shipping.ShipmentPriceInfo{
					ShipmentPriceID:     calcFeeResp.ShipmentPriceID,
					ShipmentPriceListID: calcFeeResp.ShipmentPriceListID,
				},
			}
			if _, err := a.UpdateFulfillmentShippingFees(ctx, update); err != nil {
				return err
			}
		} else {
			// this event is called in if statement already (in func UpdateFulfillmentShippingFees)
			event := &shipping.FulfillmentUpdatedEvent{
				FulfillmentID:     ffm.ID,
				MoneyTxShippingID: ffm.MoneyTransactionID,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *Aggregate) UpdateFulfillmentShippingSubstate(ctx context.Context, args *shipping.UpdateFulfillmentShippingSubstateArgs) (updated int, _ error) {
	return a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentShippingSubstate(args.ShippingSubstate)
}

func (a *Aggregate) UpdateFulfillmentShippingFees(ctx context.Context, args *shipping.UpdateFulfillmentShippingFeesArgs) (updated int, _ error) {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFulfillment()
	if err != nil {
		return 0, err
	}

	if err := canUpdateFulfillment(ffm); err != nil {
		return 0, err
	}
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var lines []*shippingsharemodel.ShippingFeeLine
		var providerShippingFeeLines []*shippingsharemodel.ShippingFeeLine
		if err := scheme.Convert(args.ProviderShippingFeeLines, &providerShippingFeeLines); err != nil {
			return err
		}
		// Khi cước phí ffm ko có bảng giá riêng
		// => phí ship tính với shop được tính từ phí NVC
		if !ffm.EtopPriceRule && providerShippingFeeLines != nil {
			lines = shippingsharemodel.GetShippingFeeShopLines(providerShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		} else {
			if err := scheme.Convert(args.ShippingFeeLines, &lines); err != nil {
				return err
			}
		}

		totalShippingFeeShop := shippingsharemodel.GetTotalShippingFee(lines)
		update := &shipmodel.Fulfillment{
			ProviderShippingFeeLines: providerShippingFeeLines,
			ShippingFeeShopLines:     lines,
			ShippingFeeShop:          shipping.CalcShopShippingFee(totalShippingFeeShop, ffm),
			UpdatedBy:                args.UpdatedBy,
			AdminNote:                args.AdminNote,
		}
		if args.ShipmentPriceInfo != nil {
			update.ShipmentPriceInfo = shippingconvert.
				Convert_shipping_ShipmentPriceInfo_sharemodel_ShipmentPriceInfo(args.ShipmentPriceInfo, nil)
		}
		if err := a.ffmStore(ctx).ID(ffm.ID).UpdateFulfillmentDB(update); err != nil {
			return err
		}

		if args.TotalCODAmount.Valid {
			if _, err := a.ffmStore(ctx).ID(ffm.ID).UpdateFulfillmentCOD(&shipping.UpdateFulfillmentCODAmountArgs{
				TotalCODAmount: args.TotalCODAmount,
				UpdatedBy:      args.UpdatedBy,
			}); err != nil {
				return err
			}
		}

		eventChanged := &shipping.FulfillmentUpdatedEvent{
			FulfillmentID:     ffm.ID,
			MoneyTxShippingID: ffm.MoneyTransactionID,
		}
		if err := a.eventBus.Publish(ctx, eventChanged); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *Aggregate) UpdateFulfillmentCODAmount(ctx context.Context, args *shipping.UpdateFulfillmentCODAmountArgs) error {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFulfillment()
	if err != nil {
		return err
	}
	if err := canUpdateFulfillment(ffm); err != nil {
		return err
	}
	if !args.TotalCODAmount.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing total COD amount")
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if _, err := a.ffmStore(ctx).ID(ffm.ID).UpdateFulfillmentCOD(args); err != nil {
			return err
		}
		eventChanged := &shipping.FulfillmentUpdatedEvent{
			FulfillmentID:     ffm.ID,
			MoneyTxShippingID: ffm.MoneyTransactionID,
		}
		if err := a.eventBus.Publish(ctx, eventChanged); err != nil {
			return err
		}
		return nil
	})
}

func (a *Aggregate) ShopUpdateFulfillmentCOD(ctx context.Context, args *shipping.ShopUpdateFulfillmentCODArgs) (updated int, _ error) {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	if !args.TotalCODAmount.Valid {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing COD Amount")
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFfmDB()
	if err != nil {
		return 0, err
	}
	if ffm.MoneyTransactionID != 0 || ffm.MoneyTransactionShippingExternalID != 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Đơn đã nằm trong phiên, không được chỉnh sửa")
	}

	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm.ConnectionID = shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)

		ffm.TotalCODAmount = args.TotalCODAmount.Int

		// case shipment: update ffm from carrier
		if ffm.ConnectionID != 0 {
			if err := a.shimentManager.UpdateFulfillmentCOD(ctx, ffm); err != nil {
				return err
			}
		}

		update := &shipping.UpdateFulfillmentCODAmountArgs{
			TotalCODAmount: args.TotalCODAmount,
			UpdatedBy:      args.UpdatedBy,
		}
		if _, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).UpdateFulfillmentCOD(update); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *Aggregate) UpdateFulfillmentsMoneyTxID(ctx context.Context, args *shipping.UpdateFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	if len(args.FulfillmentIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs").WithMetap("function", "UpdateFulfillmentsMoneyTxID")
	}

	if args.MoneyTxShippingExternalID == 0 && args.MoneyTxShippingID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxID").WithMetap("function", "UpdateFulfillmentsMoneyTxID")
	}

	return a.ffmStore(ctx).UpdateFulfillmentsMoneyTxID(args)
}

func (a *Aggregate) UpdateFulfillmentsCODTransferedAt(ctx context.Context, args *shipping.UpdateFulfillmentsCODTransferedAtArgs) error {
	if len(args.MoneyTxShippingIDs) == 0 &&
		len(args.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required fields").WithMetap("func", "UpdateFulfillmentsCODTransferedAt")
	}
	if args.CODTransferedAt.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing cod_transfered_at").WithMetap("func", "UpdateFulfillmentsCODTransferedAt")
	}
	query := a.ffmStore(ctx)
	if len(args.MoneyTxShippingIDs) > 0 {
		query = query.MoneyTxShippingIDs(args.MoneyTxShippingIDs...)
	}
	if len(args.FulfillmentIDs) > 0 {
		query = query.IDs(args.FulfillmentIDs...)
	}
	update := &shipmodel.Fulfillment{
		CODEtopTransferedAt: args.CODTransferedAt,
	}
	return query.UpdateFfmDB(update)
}

func (a *Aggregate) RemoveFulfillmentsMoneyTxID(ctx context.Context, args *shipping.RemoveFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	// make sure only remove fulfillment moneyTxID if etop not payment for shop
	return a.ffmStore(ctx).EtopNotPayment().RemoveFulfillmentsMoneyTxID(args)
}

func (a *Aggregate) UpdateFulfillmentsStatus(ctx context.Context, args *shipping.UpdateFulfillmentsStatusArgs) error {
	if len(args.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs").WithMetap("function", "UpdateFulfillmentsStatus")
	}

	return a.ffmStore(ctx).UpdateFulfillmentsStatus(args)
}

func (a *Aggregate) CancelFulfillment(ctx context.Context, args *shipping.CancelFulfillmentArgs) error {
	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFulfillment()
	if err != nil {
		return err
	}

	switch ffm.ShippingState {
	case shipstate.Unknown,
		shipstate.Default,
		shipstate.Created,
		shipstate.Picking:
	// continue
	default:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn đang ở trạng thái '%v'. Không thể hủy.", ffm.ShippingState.GetLabelRefName())
	}

	switch ffm.Status {
	case status5.P, status5.NS:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hoàn thành. Không thể hủy")
	case status5.N:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy.")
	}

	// backward compatible
	ffm.ConnectionID = shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)

	var ffmDB shipmodel.Fulfillment
	if err := scheme.Convert(ffm, &ffmDB); err != nil {
		return err
	}
	// case shipment: cancel ffm from carrier
	if ffm.ConnectionID != 0 {
		if err := a.shimentManager.CancelFulfillment(ctx, &ffmDB); err != nil {
			return err
		}
	}
	if err := a.ffmStore(ctx).CancelFulfillment(args); err != nil {
		return err
	}
	return err
}

func (a *Aggregate) UpdateFulfillmentExternalShippingInfo(ctx context.Context, args *shipping.UpdateFfmExternalShippingInfoArgs) (updated int, err error) {
	if args.FulfillmentID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Fulfillment ID")
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		update := &shipmodel.Fulfillment{
			ID:                        args.FulfillmentID,
			ShippingState:             args.ShippingState,
			ShippingSubstate:          args.ShippingSubstate,
			ShippingStatus:            args.ShippingStatus,
			ExternalShippingData:      args.ExternalShippingData,
			ExternalShippingState:     args.ExternalShippingState,
			ExternalShippingStatus:    args.ExternalShippingStatus,
			ExternalShippingUpdatedAt: args.ExternalShippingUpdatedAt,
			ExternalShippingStateCode: args.ExternalShippingStateCode,
			ClosedAt:                  args.ClosedAt,
			LastSyncAt:                args.LastSyncAt,
			ShippingCreatedAt:         args.ShippingCreatedAt,
			ShippingPickingAt:         args.ShippingPickingAt,
			ShippingHoldingAt:         args.ShippingHoldingAt,
			ShippingDeliveringAt:      args.ShippingDeliveringAt,
			ShippingDeliveredAt:       args.ShippingDeliveredAt,
			ShippingReturningAt:       args.ShippingReturningAt,
			ShippingReturnedAt:        args.ShippingReturnedAt,
			ShippingCancelledAt:       args.ShippingCancelledAt,
			ExternalShippingNote:      args.ExternalShippingNote,
			ExternalShippingSubState:  args.ExternalShippingSubState,
		}
		if args.Weight != 0 {
			// Deprecated total_weight
			update.TotalWeight = args.Weight
			update.GrossWeight = args.Weight
			update.ChargeableWeight = args.Weight
		}
		if args.ExternalShippingLogs != nil {
			logs := shippingconvert.Convert_shipping_ExternalShippingLogs_shippingmodel_ExternalShippingLogs(args.ExternalShippingLogs)
			update.ExternalShippingLogs = logs
		}

		if err := a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentDB(update); err != nil {
			return err
		}
		updated = 1
		return nil
	})
	return
}

// UpdateFulfillmentShippingFeesFromWebhook
//
// Cập nhật giá vận chuyển từ webhook
// Luôn cập nhật ProviderShippingFeeLines
// Kiểm tra nếu có thay đổi về khối lượng, sẽ tính lại giá theo bảng giá hiện tại của TOPSHIP. Nếu không có giá TOPSHIP sẽ giữ nguyên giá từ NVC => cập nhật lại giá mới hoặc thông báo qua telegram nếu đơn đã nằm trong phiên thanh toán với shop
func (a *Aggregate) UpdateFulfillmentShippingFeesFromWebhook(ctx context.Context, args *shipping.UpdateFulfillmentShippingFeesFromWebhookArgs) (_err error) {
	providerFeeLines := args.ProviderFeeLines
	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFulfillment()
	if err != nil {
		return err
	}

	update := &shipping.Fulfillment{
		ID:                       ffm.ID,
		ProviderShippingFeeLines: args.ProviderFeeLines,
	}
	var calcFeeResp *carrier.CalcMakeupShippingFeesByFfmResponse
	defer func() error {

		update.ExternalShippingFee = shippingtypes.GetTotalShippingFee(args.ProviderFeeLines)

		// always update provider shipping fee even if error occurred
		if update.ShippingFeeShopLines != nil {
			totalFee := shippingtypes.GetTotalShippingFee(update.ShippingFeeShopLines)
			shippingFeeShop := shipping.CalcShopShippingFee(totalFee, ffm)
			update.ShippingFeeShop = shippingFeeShop
			if shippingFeeShop != ffm.ShippingFeeShop {
				// Giá thay đổi
				// check money_transaction_shipping
				if ffm.MoneyTransactionID != 0 || ffm.MoneyTransactionShippingExternalID != 0 {
					// Đơn đã nằm trong phiên
					// Giá cước đơn thay đổi
					// Không cập nhật + bắn noti telegram để follow
					connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
					if connectionID == 0 {
						return cm.Errorf(cm.FailedPrecondition, nil, "ConnectionID can not be empty")
					}
					connection, _ := a.shimentManager.ConnectionManager.GetConnectionByID(ctx, connectionID)
					str := "–––\n👹 %v: đơn %v có thay đổi về giá nhưng đã nằm trong phiên thanh toán. Không thể cập nhật, vui lòng kiểm tra lại. 👹 \n- Giá hiện tại: %v \n- Giá mới: %v\n–––"
					ll.WithChannel(meta.ChannelShipmentCarrier).SendMessage(fmt.Sprintf(str, connection.Name, ffm.ShippingCode, ffm.ShippingFeeShop, shippingFeeShop))
					// shop shipping fee does not change
					update.ShippingFeeShopLines = nil
					update.ShippingFeeShop = 0
				}
			}
		}

		if update.ShippingFeeShopLines != nil && calcFeeResp != nil {
			update.ShipmentPriceInfo = &shipping.ShipmentPriceInfo{
				ShipmentPriceID:     calcFeeResp.ShipmentPriceID,
				ShipmentPriceListID: calcFeeResp.ShipmentPriceListID,
			}
		}

		if update.ShippingFeeShopLines != nil || update.ProviderShippingFeeLines != nil {
			if err := a.ffmStore(ctx).ID(ffm.ID).UpdateFulfillment(update); err != nil {
				return err
			}
		}
		// keep origin error
		return _err
	}()

	if !ffm.EtopPriceRule {
		update.ShippingFeeShopLines = shippingtypes.GetShippingFeeShopLines(providerFeeLines, false, dot.Int(0))
		return nil
	}

	// Trường hợp có áp dụng bảng giá
	// Các trường hợp cần tính lại cước phí đơn
	//   - Thay đổi khối lượng
	//   - Đơn trả hàng => tính phí trả hàng
	var feeLines []*shippingtypes.ShippingFeeLine
	shippingFeeShopLines := ffm.ShippingFeeShopLines
	calcShippingFeesArgs := &carrier.CalcMakeupShippingFeesByFfmArgs{
		Fulfillment:        ffm,
		Weight:             args.NewWeight,
		State:              args.NewState,
		AdditionalFeeTypes: nil,
	}

	// Đơn giao hàng 1 phần chỉ cần tính lại phí trả hàng nên bỏ qua phí chính
	if args.NewWeight != 0 && args.NewWeight != ffm.ChargeableWeight && !ffm.IsPartialDelivery {
		calcFeeResp, err = a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, calcShippingFeesArgs)
		if err != nil {
			return err
		}
		feeLines = calcFeeResp.ShippingFeeLines
		mainFeeLine := shippingtypes.GetShippingFeeLine(feeLines, shipping_fee_type.Main)
		shippingFeeShopLines = shippingtypes.ApplyShippingFeeLine(shippingFeeShopLines, mainFeeLine)

		// Remove field EtopAdjustedShippingFeeMain if not use
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.EtopAdjustedShippingFeeMain = mainFeeLine.Cost
	}

	if shipping.IsStateReturn(args.NewState) {
		if feeLines == nil {
			calcFeeResp, err = a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, calcShippingFeesArgs)
			if err != nil {
				return err
			}
			feeLines = calcFeeResp.ShippingFeeLines
		}
		returnFeeLine := shippingtypes.GetShippingFeeLine(feeLines, shipping_fee_type.Return)
		shippingFeeShopLines = shippingtypes.ApplyShippingFeeLine(shippingFeeShopLines, returnFeeLine)
		update.ShippingFeeShopLines = shippingFeeShopLines
	}

	return nil
}

func (a *Aggregate) UpdateFulfillmentInfo(ctx context.Context, args *shipping.UpdateFulfillmentInfoByAdminArgs) (updated int, err error) {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	if args.AdminNote == "" {
		return 0, cm.Error(cm.InvalidArgument, "Ghi chú chỉnh sửa không được để trống", nil)
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFulfillment()
	if err != nil {
		return 0, err
	}
	if err := canUpdateFulfillment(ffm); err != nil {
		return 0, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).UpdateFulfillmentInfo(args, ffm.AddressTo)
		if err != nil {
			return err
		}
		event := &shipping.FulfillmentUpdatedInfoEvent{
			OrderID:  ffm.OrderID,
			Phone:    args.Phone,
			FullName: args.FullName,
		}
		if err = a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
	return 1, err
}

func (a *Aggregate) ShopUpdateFulfillmentInfo(ctx context.Context, args *shipping.UpdateFulfillmentInfoArgs) (updated int, _ error) {
	if args.GrossWeight.Valid && args.GrossWeight.Int <= 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "gross_weight không hợp lệ")
	}

	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFulfillment()
	if err != nil {
		return 0, err
	}

	if err := checkValidShippingPaymentType(ffm.ConnectionMethod, args.ShippingPaymentType); err != nil {
		return 0, err
	}

	if _, _, err = a.getAndVerifyAddress(ctx, args.AddressFrom); err != nil {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	if _, _, err = a.getAndVerifyAddress(ctx, args.AddressTo); err != nil {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}

	switch ffm.Status {
	case status5.P, status5.NS:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hoàn thành. Không thể hủy")
	case status5.N:
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy.")
	}

	if err := canUpdateFulfillment(ffm); err != nil {
		return 0, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// backward compatible
		ffm.ConnectionID = shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)

		ffmUpdate := &shipmodel.Fulfillment{
			ShopID:              ffm.ShopID,
			ConnectionID:        ffm.ConnectionID,
			ShippingCode:        ffm.ShippingCode,
			AddressFrom:         addressconvert.Convert_orderingtypes_Address_addressmodel_Address(args.AddressFrom, nil),
			AddressTo:           addressconvert.Convert_orderingtypes_Address_addressmodel_Address(args.AddressTo, nil),
			TryOn:               args.TryOn.Apply(ffm.TryOn),
			GrossWeight:         args.GrossWeight.Apply(ffm.GrossWeight),
			ChargeableWeight:    args.GrossWeight.Apply(ffm.ChargeableWeight),
			ShippingNote:        args.ShippingNote.Apply(""),
			ShippingPaymentType: args.ShippingPaymentType.Apply(ffm.ShippingPaymentType),
		}

		// set new insuranceValue
		if args.IncludeInsurance.Valid {
			ffmUpdate.IncludeInsurance = args.IncludeInsurance
			if args.IncludeInsurance.Bool {
				ffmUpdate.InsuranceValue = args.InsuranceValue
			}
		}

		// case shipment: update ffm from carrier
		if ffm.ConnectionID != 0 {
			if err := a.shimentManager.UpdateFulfillmentInfo(ctx, ffmUpdate); err != nil {
				return err
			}
		}
		if err := a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentDB(ffmUpdate); err != nil {
			return err
		}

		return nil
	})

	return 1, err
}

func canUpdateFulfillment(ffm *shipping.Fulfillment) error {
	if !ffm.CODEtopTransferedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã đối soát").WithMetap("money_transaction_id", ffm.MoneyTransactionID)
	}
	return nil
}

func checkValidShippingPaymentType(connectionMethod connection_type.ConnectionMethod, shippingPaymentType shipping_payment_type.NullShippingPaymentType) error {
	if connectionMethod == connection_type.ConnectionMethodBuiltin &&
		shippingPaymentType.Apply(shipping_payment_type.Seller) != shipping_payment_type.Seller {
		return cm.Errorf(cm.InvalidArgument, nil, "Hình thức thanh toán phí giao hàng không hợp lệ.")
	}
	return nil
}

// AddFulfillmentShippingFee
//
// Sử dụng khi thêm phí kích hoạt giao lại/phí điều chỉnh thông tin đơn hàng
// Các phí trên có thể thêm nhiều lần
var allowAddShippingFeeTypes = []shipping_fee_type.ShippingFeeType{shipping_fee_type.Redelivery, shipping_fee_type.Adjustment}

func (a *Aggregate) AddFulfillmentShippingFee(ctx context.Context, args *shipping.AddFulfillmentShippingFeeArgs) error {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFulfillment()
	if err != nil {
		return err
	}

	if !shipping_fee_type.Contain(allowAddShippingFeeTypes, args.ShippingFeeType) {
		return cm.Errorf(cm.InvalidArgument, nil, "Chỉ hỗ trợ thêm phí kích hoạt giao lại, phí điều chỉnh thông tin đơn hàng.")
	}
	if err := canUpdateFulfillment(ffm); err != nil {
		return err
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		calcShippingFeesArgs := &carrier.CalcMakeupShippingFeesByFfmArgs{
			Fulfillment:        ffm,
			AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{args.ShippingFeeType},
		}
		calcFeeResp, err := a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, calcShippingFeesArgs)
		if err != nil {
			return err
		}
		feeLine := shippingtypes.GetShippingFeeLine(calcFeeResp.ShippingFeeLines, args.ShippingFeeType)
		if feeLine == nil || feeLine.Cost == 0 {
			return nil
		}
		update := &shipping.UpdateFulfillmentShippingFeesArgs{
			FulfillmentID:    ffm.ID,
			ShippingCode:     ffm.ShippingCode,
			ShippingFeeLines: append(ffm.ShippingFeeShopLines, feeLine),
			UpdatedBy:        args.UpdatedBy,
			ShipmentPriceInfo: &shipping.ShipmentPriceInfo{
				ShipmentPriceID:     calcFeeResp.ShipmentPriceID,
				ShipmentPriceListID: calcFeeResp.ShipmentPriceListID,
			},
		}
		if _, err := a.UpdateFulfillmentShippingFees(ctx, update); err != nil {
			return err
		}
		return nil
	})
}

func (a *Aggregate) CreateFulfillmentsFromImport(
	ctx context.Context, args *shipping.CreateFulfillmentsFromImportArgs,
) ([]*shipping.CreateFullfillmentsFromImportResult, error) {
	if len(args.Fulfillments) == 0 {
		return nil, nil
	}

	result := make([]*shipping.CreateFullfillmentsFromImportResult, len(args.Fulfillments))

	var wg sync.WaitGroup
	var m sync.Mutex

	validationErrors, err := a.validateFulfillmentsFromImport(ctx, args.Fulfillments)
	if err != nil {
		return nil, err
	}

	for idx := range validationErrors {
		// ignore row validation fail
		if validationErrors[idx] != nil {
			result[idx] = &shipping.CreateFullfillmentsFromImportResult{
				Error: validationErrors[idx],
			}
			continue
		}
		importFulfillmentArgs := args.Fulfillments[idx]
		wg.Add(1)
		go func(_idx int, args *shipping.CreateFulfillmentFromImportArgs) {
			defer func() {
				wg.Done()
				m.Unlock()
			}()

			fulfillmentID := cm.NewID()
			args.ID = fulfillmentID
			err := a.createFulfillmentFromImport(ctx, args)
			m.Lock()

			result[_idx] = &shipping.CreateFullfillmentsFromImportResult{
				FulfillmentID: fulfillmentID,
				Error:         err,
			}
		}(idx, importFulfillmentArgs)
	}

	wg.Wait()

	return result, nil
}

func (a *Aggregate) validateFulfillmentsFromImport(ctx context.Context, createFfmsFromImportArgs []*shipping.CreateFulfillmentFromImportArgs) (validationErrors []error, err error) {
	shopID := createFfmsFromImportArgs[0].ShopID

	// Check ed_code
	{
		validationErrors = make([]error, len(createFfmsFromImportArgs))

		mapEdCodeAndRowIdxs := make(map[string][]int)
		for i, createFfmFromImportArgs := range createFfmsFromImportArgs {
			if createFfmFromImportArgs.EdCode == "" {
				continue
			}
			if _, ok := mapEdCodeAndRowIdxs[createFfmFromImportArgs.EdCode]; ok {
				mapEdCodeAndRowIdxs[createFfmFromImportArgs.EdCode] = append(mapEdCodeAndRowIdxs[createFfmFromImportArgs.EdCode], i)
			} else {
				mapEdCodeAndRowIdxs[createFfmFromImportArgs.EdCode] = []int{i}
			}
		}

		var edCodes []string
		for edCode, rowIdxs := range mapEdCodeAndRowIdxs {
			if len(rowIdxs) == 1 {
				edCodes = append(edCodes, edCode)
				continue
			}
			for _, rowIdx := range rowIdxs {
				validationErrors[rowIdx] = cm.Errorf(cm.InvalidArgument, nil, "Mã nội bộ đã tồn tại")
			}
		}

		if len(edCodes) > 0 {
			ffms, err := a.ffmStore(ctx).ShopID(shopID).StatusNotIn([]status5.Status{status5.N}...).EdCodes(edCodes).ListFfms()
			if err != nil {
				return nil, err
			}
			for _, ffm := range ffms {
				for _, rowIdx := range mapEdCodeAndRowIdxs[ffm.EdCode] {
					validationErrors[rowIdx] = cm.Errorf(cm.InvalidArgument, nil, "Mã nội bộ đã tồn tại trong hệ thống")
				}
			}
		}
	}

	for i := range validationErrors {
		if validationErrors[i] == nil {
			validationErrors[i] = a.validateFulfillmentFromImport(ctx, createFfmsFromImportArgs[i])
		}
	}

	return validationErrors, nil
}

func (a *Aggregate) createFulfillmentFromImport(ctx context.Context, args *shipping.CreateFulfillmentFromImportArgs) (err error) {
	err = a.validateFulfillmentFromImport(ctx, args)
	if err != nil {
		return err
	}

	var ffm *shipmodel.Fulfillment
	ffm, err = a.prepareFulfillmentImport(ctx, args)
	if err != nil {
		return err
	}

	if err = a.ffmStore(ctx).CreateFulfillmentsDB([]*shipmodel.Fulfillment{ffm}); err != nil {
		return err
	}

	// Ignore error when create customers to reduce impact to the operation
	fulfillmentFromImportCreatedEvent := &shipping.FulfillmentFromImportCreatedEvent{
		EventMeta:     meta.NewEvent(),
		ShopID:        ffm.ShopID,
		FulfillmentID: ffm.ID,
	}
	_ = a.eventBus.Publish(ctx, fulfillmentFromImportCreatedEvent)

	defer func() {
		// rollback when get error
		if err != nil {
			cancelFulfillmentArgs := &shipping.CancelFulfillmentArgs{
				CancelReason:  "cancel import fulfillment",
				FulfillmentID: ffm.ID,
			}
			if _err := a.ffmStore(ctx).CancelFulfillment(cancelFulfillmentArgs); _err != nil {
				panic(_err)
			}
		}
	}()

	if err = a.shimentManager.CreateFulfillments(ctx, []*shipmodel.Fulfillment{ffm}); err != nil {
		return err
	}

	return
}

func (a *Aggregate) validateFulfillmentFromImport(ctx context.Context, args *shipping.CreateFulfillmentFromImportArgs) error {
	if strings.TrimSpace(args.ShippingAddress.FullName) == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tên người nhận không được bỏ trống.")
	}
	if strings.TrimSpace(args.ProductDescription) == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mô tả sản phẩm không được bỏ trống")
	}
	if _, _, err := a.getAndVerifyAddress(ctx, args.ShippingAddress); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
	}
	if _, _, err := a.getAndVerifyAddress(ctx, args.PickupAddress); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}

	if args.BasketValue < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị hàng hoá không hợp lệ: %v", args.BasketValue)
	}
	if args.TotalWeight < 50 {
		return cm.Errorf(cm.InvalidArgument, nil, "Khối lượng tối thiểu 50 gram", args.TotalWeight)
	}
	if args.CODAmount != 0 && args.CODAmount < 5000 {
		return cm.Errorf(cm.InvalidArgument, nil, "Thu hộ bằng 0đ hoặc từ 5000đ trở lên", args.CODAmount)
	}

	return nil
}

func (a *Aggregate) prepareFulfillmentImport(ctx context.Context, args *shipping.CreateFulfillmentFromImportArgs) (*shipmodel.Fulfillment, error) {
	shippingType := ordertypes.ShippingTypeShipment
	var connectionMethod connection_type.ConnectionMethod
	var conn *connectioning.Connection

	if args.ID == 0 {
		return nil, cm.Errorf(cm.Internal, nil, "id is missing")
	}

	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển (connection_id)")
	}
	conn, err := a.shimentManager.ConnectionManager.GetConnectionByID(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}

	if conn.ConnectionProvider == connection_type.ConnectionProviderGHN {
		if args.TryOn.String() == "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng chọn ghi chú xem hàng!")
		}
	}
	connectionMethod = conn.ConnectionMethod

	tryOn := args.TryOn
	if tryOn == 0 {
		tryOn = try_on.None
	}

	typeFrom := etopmodel.FFShop
	typeTo := etopmodel.FFCustomer

	ffm := &shipmodel.Fulfillment{
		ID:                  args.ID,
		ShopID:              args.ShopID,
		EdCode:              args.EdCode,
		ShopConfirm:         status3.P, // Always set shop_confirm to 1
		TotalItems:          1,         // hardcode
		BasketValue:         args.BasketValue,
		TotalAmount:         args.BasketValue,
		TotalCODAmount:      args.CODAmount,
		TypeFrom:            typeFrom,
		TypeTo:              typeTo,
		AddressFrom:         addressconvert.OrderAddressToModel(args.PickupAddress),
		AddressTo:           addressconvert.OrderAddressToModel(args.ShippingAddress),
		AddressReturn:       addressconvert.OrderAddressToModel(args.PickupAddress),
		ProviderServiceID:   args.ShippingServiceCode,
		ShippingServiceFee:  args.ShippingServiceFee,
		ShippingServiceName: args.ShippingServiceName,
		ShippingNote:        args.ShippingNote,
		TryOn:               tryOn,
		IncludeInsurance:    dot.Bool(args.IncludeInsurance),
		ConnectionID:        args.ConnectionID,
		ConnectionMethod:    connectionMethod,
		TotalWeight:         args.TotalWeight,
		ChargeableWeight:    args.TotalWeight,
		GrossWeight:         args.TotalWeight,
		ShippingState:       shipstate.Default,
		ShippingType:        shippingType,
		ShippingPaymentType: shipping_payment_type.Seller,
		LinesContent:        args.ProductDescription,
		CreatedBy:           args.CreatedBy,
	}

	if conn != nil {
		// backward compatible
		if shippingProvider, ok := shipping_provider.ParseShippingProvider(conn.ConnectionProvider.String()); ok {
			ffm.ShippingProvider = shippingProvider
		}
	}
	return ffm, nil
}

func (a *Aggregate) UpdateFulfillmentShippingCode(ctx context.Context, args *shipping.UpdateFulfillmentShippingCodeArgs) error {
	if args.FulfillmentID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing fulfillment_id")
	}
	if args.ShippingCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing new shipping_code")
	}
	update := &shipmodel.Fulfillment{ShippingCode: args.ShippingCode}
	return a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentDB(update)
}

func (a *Aggregate) CreatePartialFulfillment(ctx context.Context, args *shipping.CreatePartialFulfillmentArgs) (fulfillmentID dot.ID, _ error) {
	oldFfmModel, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFfmDB()
	if err != nil {
		return 0, err
	}

	t0 := time.Now()

	newFfm := &model.Fulfillment{
		ID:                       cm.NewID(),
		OrderID:                  oldFfmModel.OrderID,
		ShopID:                   oldFfmModel.ShopID,
		PartnerID:                oldFfmModel.PartnerID,
		ShopConfirm:              status3.P,
		ConfirmStatus:            status3.P,
		TotalItems:               oldFfmModel.TotalItems,
		TotalWeight:              0,
		BasketValue:              oldFfmModel.BasketValue,
		TotalDiscount:            0,
		TotalAmount:              oldFfmModel.TotalAmount,
		TotalCODAmount:           0,
		OriginalCODAmount:        0,
		ActualCompensationAmount: oldFfmModel.ActualCompensationAmount,
		EtopPriceRule:            oldFfmModel.EtopPriceRule,
		VariantIDs:               oldFfmModel.VariantIDs,
		Lines:                    oldFfmModel.Lines,
		TypeFrom:                 oldFfmModel.TypeFrom,
		TypeTo:                   oldFfmModel.TypeTo,
		AddressFrom:              oldFfmModel.AddressFrom,
		AddressTo:                oldFfmModel.AddressTo,
		AddressReturn:            oldFfmModel.AddressReturn,
		AddressToProvinceCode:    oldFfmModel.AddressToFullNameNorm,
		AddressToDistrictCode:    oldFfmModel.AddressToDistrictCode,
		AddressToWardCode:        oldFfmModel.AddressToWardCode,
		AddressToPhone:           oldFfmModel.AddressToPhone,
		AddressToFullNameNorm:    oldFfmModel.AddressToFullNameNorm,
		CreatedAt:                t0,
		UpdatedAt:                t0,
		ShippingProvider:         oldFfmModel.ShippingProvider,
		ProviderServiceID:        oldFfmModel.ProviderServiceID,
		ShippingCode:             "",
		ShippingNote:             fmt.Sprintf("Đơn giao hàng một phần được tạo từ đơn giao hàng %s", oldFfmModel.ShippingCode),
		TryOn:                    oldFfmModel.TryOn,
		ShippingType:             oldFfmModel.ShippingType,
		ShippingPaymentType:      oldFfmModel.ShippingPaymentType,
		ConnectionID:             oldFfmModel.ConnectionID,
		ConnectionMethod:         oldFfmModel.ConnectionMethod,
		ShippingServiceName:      oldFfmModel.ShippingServiceName,
		ExternalShippingName:     oldFfmModel.ExternalShippingName,
		ExternalShippingID:       oldFfmModel.ExternalShippingID,
		ShippingState:            shipstate.Returning,
		ShippingStatus:           shipstate.Returning.ToStatus5(),
		Status:                   status5.P,
		IsPartialDelivery:        true,
		CreatedBy:                oldFfmModel.CreatedBy,
		DeliveryRoute:            oldFfmModel.DeliveryRoute,
	}

	if args.InfoChanges != nil {
		newFfm.ShippingCode = args.InfoChanges.ShippingCode.String
		newFfm.ExternalShippingCode = args.InfoChanges.ShippingCode.String
		newFfm.GrossWeight = args.InfoChanges.Weight.Int
		newFfm.ChargeableWeight = args.InfoChanges.Weight.Int
		newFfm.Length = args.InfoChanges.Length.Int
		newFfm.Height = args.InfoChanges.Height.Int
		newFfm.Width = args.InfoChanges.Height.Int
	}

	ffm, err := a.ffmStore(ctx).CreateFulfillmentDB(newFfm)
	if err != nil {
		return 0, err
	}

	return ffm.ID, nil
}
