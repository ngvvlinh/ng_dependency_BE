package pm

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/service_classify"
	shippingstate "o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
)

const MinShopBalance = -100000

type ProcessManager struct {
	eventBus        capi.EventBus
	shippingAggr    shipping.CommandBus
	shippingQuery   shipping.QueryBus
	redisStore      redis.Store
	connectionQuery connectioning.QueryBus

	ShopStore        sqlstore.ShopStoreInterface
	TransactionQuery transaction.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	shippingQ shipping.QueryBus,
	shippingA shipping.CommandBus,
	redisS redis.Store,
	connectionQ connectioning.QueryBus,
	ShopStore sqlstore.ShopStoreInterface,
	TransactionQ transaction.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:         eventBus,
		shippingQuery:    shippingQ,
		shippingAggr:     shippingA,
		redisStore:       redisS,
		connectionQuery:  connectionQ,
		ShopStore:        ShopStore,
		TransactionQuery: TransactionQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
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
	eventBus.AddEventListener(m.HandleDHLFulfillmentCancelledEvent)
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

	for _, ffm := range query.Result {
		if ffm.ShippingState != shippingstate.Returned && ffm.ShippingState != shippingstate.Returning {
			continue
		}
		returnedFee := shippingtypes.GetShippingFee(ffm.ShippingFeeShopLines, shipping_fee_type.Return)
		newReturnedFee := CalcVtpostShippingFeeReturned(ffm)
		if newReturnedFee == 0 || newReturnedFee == returnedFee {
			continue
		}
		providerShippingFeeLines := shippingtypes.UpdateShippingFees(ffm.ProviderShippingFeeLines, newReturnedFee, shipping_fee_type.Return)
		shippingFeeShopLines := shippingtypes.UpdateShippingFees(ffm.ShippingFeeShopLines, newReturnedFee, shipping_fee_type.Return)
		update := &shipping.UpdateFulfillmentShippingFeesCommand{
			FulfillmentID:            ffm.ID,
			ProviderShippingFeeLines: providerShippingFeeLines,
			ShippingFeeLines:         shippingFeeShopLines,
		}
		if err := m.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

// CalcVtpostShippingFeeReturned: Tính cước phí trả hàng vtpost
func CalcVtpostShippingFeeReturned(ffm *shipping.Fulfillment) int {
	// Nội tỉnh miễn phí trả hàng
	// Liên tỉnh 50% cước phí chiều đi
	from := ffm.AddressFrom
	to := ffm.AddressTo
	if from.ProvinceCode == to.ProvinceCode {
		return 0
	}

	returnedFee := shippingtypes.GetShippingFee(ffm.ShippingFeeShopLines, shipping_fee_type.Return)
	totalFee := shippingtypes.GetTotalShippingFee(ffm.ShippingFeeShopLines)
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
	key := connectionmanager.GetRedisConnectionKeyByID(event.ConnectionID)
	return m.redisStore.Del(key)
}

func (m *ProcessManager) ShopConnectionUpdated(ctx context.Context, event *connectioning.ShopConnectionUpdatedEvent) error {
	if event.ConnectionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID").WithMetap("event", "ShopConnectionUpdatedEvent")
	}
	// Delete cache shop connection in carrier manager
	availableRedisKeys := []string{}
	key1 := connectionmanager.GetRedisShopConnectionKey(connectionmanager.GetShopConnectionArgs{ConnectionID: event.ConnectionID, ShopID: event.ShopID})
	availableRedisKeys = append(availableRedisKeys, key1)

	key2 := connectionmanager.GetRedisShopConnectionKey(connectionmanager.GetShopConnectionArgs{ConnectionID: event.ConnectionID, OwnerID: event.OwnerID})
	if !cm.StringsContain(availableRedisKeys, key2) {
		availableRedisKeys = append(availableRedisKeys, key2)
	}

	if err := m.redisStore.Del(availableRedisKeys...); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) SingleFulfillmentCreatingEvent(ctx context.Context, event *shipping.SingleFulfillmentCreatingEvent) error {
	if event.ConnectionID != 0 {
		queryConn := &connectioning.GetConnectionByIDQuery{
			ID: event.ConnectionID,
		}
		err := m.connectionQuery.Dispatch(ctx, queryConn)
		if err != nil {
			return err
		}

		if queryConn.Result.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
			return nil
		}
	}

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

	// Trường hợp địa chỉ lấy hàng nằm ngoài HCM, HN
	// Tính số dư user: GetBalanceUser
	// Số dư còn lại (1) = số dư User - shippingFee
	// Nếu (1) < 0 => không cho tạo đơn giao hàng
	queryShop := &identitymodelx.GetShopQuery{
		ShopID: event.ShopID,
	}
	if err := m.ShopStore.GetShop(ctx, queryShop); err != nil {
		return err
	}

	query := &transaction.GetBalanceUserQuery{
		UserID:   queryShop.Result.OwnerID,
		Classify: service_classify.Shipping,
	}
	if err := m.TransactionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	actualBalance := query.Result.ActualBalance

	// HCM, HN
	if cm.StringsContain(provinces, fromAddress.ProvinceCode) {
		if actualBalance-event.ShippingFee < MinShopBalance {
			return cm.Errorf(cm.FailedPrecondition, nil, "Số dư của bạn không đủ để tạo đơn. Vui lòng nạp thêm tiền.")
		}
		return nil
	}

	if actualBalance-event.ShippingFee < 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số dư của bạn không đủ để tạo đơn. Vui lòng nạp thêm tiền.")
	}
	return nil
}

func (m *ProcessManager) HandleDHLFulfillmentCancelledEvent(ctx context.Context, event *shipping.DHLFulfillmentCancelledEvent) error {
	updateFfmShippingSubstateCmd := &shipping.UpdateFulfillmentShippingSubstateCommand{
		FulfillmentID:    event.FulfillmentID,
		ShippingSubstate: shippingsubstate.Cancelling,
	}
	return m.shippingAggr.Dispatch(ctx, updateFfmShippingSubstateCmd)
}
