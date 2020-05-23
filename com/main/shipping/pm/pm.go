package pm

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
)

const MinShopBalance = -200000

type ProcessManager struct {
	eventBus      capi.EventBus
	shippingAggr  shipping.CommandBus
	shippingQuery shipping.QueryBus
	redisStore    redis.Store
}

func New(eventBus capi.EventBus, shippingQ shipping.QueryBus, shippingA shipping.CommandBus, redisS redis.Store) *ProcessManager {
	return &ProcessManager{
		eventBus:      eventBus,
		shippingQuery: shippingQ,
		shippingAggr:  shippingA,
		redisStore:    redisS,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ConnectionUpdated)
	eventBus.AddEventListener(m.ShopConnectionUpdated)
	eventBus.AddEventListener(m.MoneyTxShippingExternalCreated)
	eventBus.AddEventListener(m.MoneyTxShippingExternalLinesDeleted)
	eventBus.AddEventListener(m.MoneyTxShippingExternalDeleted)
	eventBus.AddEventListener(m.MoneyTxShippingExternalsConfirming)
	eventBus.AddEventListener(m.MoneyTxShippingCreated)
	eventBus.AddEventListener(m.MoneyTxShippingConfirmed)
	eventBus.AddEventListener(m.MoneyTxShippingDeleted)
	eventBus.AddEventListener(m.MoneyTxShippingRemoveFfms)
	eventBus.AddEventListener(m.MoneyTxShippingEtopConfirmed)
	eventBus.AddEventListener(m.SingleFulfillmentCreatingEvent)
}

func (m *ProcessManager) MoneyTxShippingExternalCreated(ctx context.Context, event *moneytx.MoneyTxShippingExternalCreatedEvent) error {
	if event.MoneyTxShippingExternalID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Event MoneyTransactionShippingExternalCreated missing ID")
	}
	if len(event.FulfillementIDs) == 0 {
		return nil
	}
	cmd := &shipping.UpdateFulfillmentsMoneyTxIDCommand{
		FulfillmentIDs:            event.FulfillementIDs,
		MoneyTxShippingExternalID: event.MoneyTxShippingExternalID,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) MoneyTxShippingExternalLinesDeleted(ctx context.Context, event *moneytx.MoneyTxShippingExternalLinesDeletedEvent) error {
	if len(event.FulfillmentIDs) == 0 {
		return nil
	}
	cmd := &shipping.RemoveFulfillmentsMoneyTxIDCommand{
		FulfillmentIDs: event.FulfillmentIDs,
	}
	if err := m.shippingAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) MoneyTxShippingExternalDeleted(ctx context.Context, event *moneytx.MoneyTxShippingExternalDeletedEvent) error {
	if event.MoneyTxShippingExternalID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingExternalID").WithMetap("event", "MoneyTxShippingExternalDeleted")
	}

	cmd := &shipping.RemoveFulfillmentsMoneyTxIDCommand{
		MoneyTxShippingExternalID: event.MoneyTxShippingExternalID,
	}
	if err := m.shippingAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) MoneyTxShippingCreated(ctx context.Context, event *moneytx.MoneyTxShippingCreatedEvent) error {
	if event.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingID").WithMetap("event", "MoneyTxShippingCreatedEvent")
	}
	if len(event.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "FulfillmentIDs can not be empty").WithMetap("event", "MoneyTxShippingCreatedEvent")
	}
	cmd := &shipping.UpdateFulfillmentsMoneyTxIDCommand{
		FulfillmentIDs:    event.FulfillmentIDs,
		MoneyTxShippingID: event.MoneyTxShippingID,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) MoneyTxShippingConfirmed(ctx context.Context, event *moneytx.MoneyTxShippingConfirmedEvent) error {
	if event.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingID").WithMetap("event", "MoneyTxShippingConfirmed")
	}
	cmd := &shipping.UpdateFulfillmentsCODTransferedAtCommand{
		MoneyTxShippingIDs: []dot.ID{event.MoneyTxShippingID},
		CODTransferedAt:    event.ConfirmedAt,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) MoneyTxShippingDeleted(ctx context.Context, event *moneytx.MoneyTxShippingDeletedEvent) error {
	if event.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingID").WithMetap("event", "MoneyTxShippingDeleted")
	}
	cmd := &shipping.RemoveFulfillmentsMoneyTxIDCommand{
		MoneyTxShippingID: event.MoneyTxShippingID,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) MoneyTxShippingRemoveFfms(ctx context.Context, event *moneytx.MoneyTxShippingRemovedFfmsEvent) error {
	if event.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingID").WithMetap("event", "MoneyTxShippingRemoveFfms")
	}
	cmd := &shipping.RemoveFulfillmentsMoneyTxIDCommand{
		FulfillmentIDs: event.FulfillmentIDs,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) MoneyTxShippingExternalsConfirming(ctx context.Context, event *moneytx.MoneyTxShippingExternalsConfirmingEvent) error {
	// Lọc tất cả các đơn trả hàng và đang trả hàng của VTPOST
	// Dùng ListFulfillmentsForMoneyTxQuery để lấy những đơn chưa đối soát
	// Cập nhật phí trả hàng nếu chưa có
	query := &shipping.ListFulfillmentsForMoneyTxQuery{
		ShippingProvider: shipping_provider.VTPost,
		ShippingStates:   []shippingstate.State{shippingstate.Returning, shippingstate.Returned},
	}
	if err := m.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var ffms []*shipmodel.Fulfillment

	ffms = shippingconvert.Convert_shipping_Fulfillments_shippingmodel_Fulfillments(query.Result)

	for _, ffm := range ffms {
		if ffm.ShippingState != shippingstate.Returned && ffm.ShippingState != shippingstate.Returning {
			continue
		}
		returnedFee := shippingsharemodel.GetReturnedFee(ffm.ShippingFeeShopLines)
		newReturnedFee := CalcVtpostShippingFeeReturned(ffm)
		if newReturnedFee == 0 || newReturnedFee == returnedFee {
			continue
		}
		lines := ffm.ProviderShippingFeeLines
		providerShippingFeeLines := shippingsharemodel.UpdateShippingFees(lines, newReturnedFee, shipping_fee_type.Return)
		shippingFeeShopLines := shippingsharemodel.GetShippingFeeShopLines(providerShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		update := &shipping.UpdateFulfillmentShippingFeesCommand{
			FulfillmentID:            ffm.ID,
			ProviderShippingFeeLines: shippingconvert.Convert_sharemodel_ShippingFeeLines_shipping_ShippingFeeLines(providerShippingFeeLines),
			ShippingFeeLines:         shippingconvert.Convert_sharemodel_ShippingFeeLines_shipping_ShippingFeeLines(shippingFeeShopLines),
		}
		if err := m.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

// CalcVtpostShippingFeeReturned: Tính cước phí trả hàng vtpost
func CalcVtpostShippingFeeReturned(ffm *shipmodel.Fulfillment) int {
	// Nội tỉnh miễn phí trả hàng
	// Liên tỉnh 50% cước phí chiều đi
	from := ffm.AddressFrom
	to := ffm.AddressTo
	if from.ProvinceCode == to.ProvinceCode {
		return 0
	}

	returnedFee := shippingsharemodel.GetReturnedFee(ffm.ShippingFeeShopLines)
	totalFee := shippingsharemodel.GetTotalShippingFee(ffm.ShippingFeeShopLines)
	newReturnedFee := (totalFee - returnedFee) / 2
	return newReturnedFee
}

func (m *ProcessManager) MoneyTxShippingEtopConfirmed(ctx context.Context, event *moneytx.MoneyTxShippingEtopConfirmedEvent) error {
	if event.MoneyTxShippingEtopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingEtopID").WithMetap("event", "MoneyTxShippingEtopConfirmed")
	}
	if len(event.MoneyTxShippingIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "MoneyTxShippingIDs can not be empty").WithMetap("event", "MoneyTxShippingEtopConfirmed")
	}

	cmd := &shipping.UpdateFulfillmentsCODTransferedAtCommand{
		MoneyTxShippingIDs: event.MoneyTxShippingIDs,
		CODTransferedAt:    event.ConfirmedAt,
	}
	return m.shippingAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) ConnectionUpdated(ctx context.Context, event *connectioning.ConnectionUpdatedEvent) error {
	if event.ConnectionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID").WithMetap("event", "ConnectionUpdatedEvent")
	}
	// Delete cache connection in carrier manager
	key := carrier.GetRedisConnectionKeyByID(event.ConnectionID)
	return m.redisStore.Del(key)
}

func (m *ProcessManager) ShopConnectionUpdated(ctx context.Context, event *connectioning.ShopConnectionUpdatedEvent) error {
	if event.ConnectionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID").WithMetap("event", "ShopConnectionUpdatedEvent")
	}
	// Delete cache connection in carrier manager
	key := carrier.GetRedisShopConnectionKey(event.ConnectionID, event.ShopID)
	return m.redisStore.Del(key)
}

func (m *ProcessManager) SingleFulfillmentCreatingEvent(ctx context.Context, event *shipping.SingleFulfillmentCreatingEvent) error {
	fromAddress := event.FromAddress
	if fromAddress == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing from address").WithMeta("event", "SingleFulfillmentCreatingEvent")
	}
	if event.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID").WithMeta("event", "SingleFulfillmentCreatingEvent")
	}

	provinces := []string{
		"01", // HN
		"79", // HCM
	}
	if cm.StringsContain(provinces, fromAddress.ProvinceCode) {
		query := &etopmodel.GetBalanceShopCommand{
			ShopID: event.ShopID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			return err
		}
		balance := query.Result.Amount
		if balance-event.ShippingFee < MinShopBalance {
			return cm.Errorf(cm.FailedPrecondition, nil, "Số dư của bạn không đủ để tạo đơn. Vui lòng nạp thêm tiền.")
		}
		return nil
	}

	// Trường hợp địa chỉ lấy hàng nằm ngoài HCM, HN
	// Tính số dư user: GetBalanceUser
	// Số dư còn lại (1) = số dư User - shippingFee
	// Nếu (1) < 0 => không cho tạo đơn giao hàng
	queryShop := &identitymodelx.GetShopQuery{
		ShopID: event.ShopID,
	}
	if err := bus.Dispatch(ctx, queryShop); err != nil {
		return err
	}
	query := &etopmodel.GetBalanceUserCommand{
		UserID: queryShop.Result.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	balance := query.Result.Amount
	if balance-event.ShippingFee < 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số dư của bạn không đủ để tạo đơn. Vui lòng nạp thêm tiền.")
	}
	return nil
}
