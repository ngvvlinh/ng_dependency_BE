package pm

import (
	"context"
	"time"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/address"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	etopconvert "o.o/backend/com/main/etop/convert"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type ProcessManager struct {
	eventBus            capi.EventBus
	shipnowQuery        shipnow.QueryBus
	shipnow             shipnow.CommandBus
	order               ordering.CommandBus
	shipnowManager      carrier.Manager
	identityQuery       identity.QueryBus
	identityAggr        identity.CommandBus
	addressQuery        address.QueryBus
	accountshipnowQuery accountshipnow.QueryBus
	accountshipnowAggr  accountshipnow.CommandBus
	connectionAggr      connectioning.CommandBus
}

func New(
	eventBus bus.EventRegistry,
	shipnowQuery shipnow.QueryBus,
	shipnowAggrBus shipnow.CommandBus,
	orderAggrBus ordering.CommandBus,
	carrierManager carrier.Manager,
	identityQS identity.QueryBus,
	addressQS address.QueryBus,
	accountshipnowQS accountshipnow.QueryBus,
	accountshipnowA accountshipnow.CommandBus,
	connectionA connectioning.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:            eventBus,
		shipnowQuery:        shipnowQuery,
		shipnow:             shipnowAggrBus,
		order:               orderAggrBus,
		shipnowManager:      carrierManager,
		identityQuery:       identityQS,
		addressQuery:        addressQS,
		accountshipnowQuery: accountshipnowQS,
		accountshipnowAggr:  accountshipnowA,
		connectionAggr:      connectionA,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShipnowOrderReservation)
	eventBus.AddEventListener(m.ShipnowOrderChanged)
	eventBus.AddEventListener(m.ShipnowCancelled)
	eventBus.AddEventListener(m.ValidateConfirmed)
	eventBus.AddEventListener(m.ShipnowExternalCreated)
	eventBus.AddEventListener(m.ExternalAccountAhamoveCreated)
	eventBus.AddEventListener(m.ExternalAccountAhamoveVerifyRequested)
	eventBus.AddEventListener(m.ExternalAccountShipnowUpdateVerificationInfo)
}

func (m *ProcessManager) ShipnowOrderReservation(ctx context.Context, event *shipnow.ShipnowOrderReservationEvent) error {
	// Call orderAggr for ReserveOrdersForFfm
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIDs,
		Fulfill:    ordertypes.ShippingTypeShipnow,
		FulfillIDs: []dot.ID{event.ShipnowFulfillmentID},
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowOrderChanged(ctx context.Context, event *shipnow.ShipnowOrderChangedEvent) error {
	// release old orderIDs and reserve new orderIDs
	cmd := &ordering.ReleaseOrdersForFfmCommand{
		OrderIDs: event.OldOrderIDs,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	cmd2 := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   event.OrderIDs,
		Fulfill:    ordertypes.ShippingTypeShipnow,
		FulfillIDs: []dot.ID{event.ShipnowFulfillmentID},
	}
	if err := m.order.Dispatch(ctx, cmd2); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowCancelled(ctx context.Context, event *shipnow.ShipnowCancelledEvent) error {
	// release orderIDs
	cmd := &ordering.ReleaseOrdersForFfmCommand{
		OrderIDs: event.OrderIDs,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	query := &shipnow.GetShipnowFulfillmentQuery{
		ID: event.ShipnowFulfillmentID,
	}
	if err := m.shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result.ShipnowFulfillment
	if ffm.ShippingCode != "" {
		cmd2 := &carrier.CancelExternalShipnowCommand{
			ShopID:               ffm.ShopID,
			ShipnowFulfillmentID: ffm.ID,
			ExternalShipnowID:    ffm.ShippingCode,
			CarrierServiceCode:   ffm.ShippingServiceCode,
			CancelReason:         event.CancelReason,
			Carrier:              ffm.Carrier,
			ConnectionID:         ffm.ConnectionID,
		}
		if err := m.shipnowManager.CancelExternalShipping(ctx, cmd2); err != nil {
			return err
		}
	}

	return nil
}

func (m *ProcessManager) ValidateConfirmed(ctx context.Context, event *shipnow.ShipnowValidateConfirmedEvent) error {
	cmd := &ordering.ValidateOrdersForShippingCommand{
		OrderIDs: event.OrderIDs,
	}
	if err := m.order.Dispatch(ctx, cmd); err != nil {
		return err
	}

	// update order confirm status
	cmd2 := &ordering.UpdateOrdersConfirmStatusCommand{
		IDs:           event.OrderIDs,
		ShopConfirm:   status3.P,
		ConfirmStatus: status3.P,
	}
	if err := m.order.Dispatch(ctx, cmd2); err != nil {
		return err
	}
	return nil
}

func (m *ProcessManager) ShipnowExternalCreated(ctx context.Context, event *shipnow.ShipnowExternalCreatedEvent) (_err error) {
	query := &shipnow.GetShipnowFulfillmentQuery{
		ID: event.ShipnowFulfillmentID,
	}
	if err := m.shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result.ShipnowFulfillment
	{
		// update sync status
		update := &shipnow.UpdateShipnowFulfillmentStateCommand{
			Id:         ffm.ID,
			SyncStatus: status4.S,
			SyncStates: &shipnow.SyncStates{
				TrySyncAt: time.Now(),
			},
		}
		if err := m.shipnow.Dispatch(ctx, update); err != nil {
			return err
		}
	}

	defer func() {
		if _err == nil {
			return
		}
		update := &shipnow.UpdateShipnowFulfillmentStateCommand{
			Id:         ffm.ID,
			SyncStatus: status4.N,
			SyncStates: &shipnow.SyncStates{
				TrySyncAt: time.Now(),
				Error:     etopconvert.Error(model.ToError(_err)),
			},
		}
		// Keep the original error
		_ = m.shipnow.Dispatch(ctx, update)
	}()

	cmd := &carrier.CreateExternalShipnowCommand{
		ShopID:               ffm.ShopID,
		ShipnowFulfillmentID: ffm.ID,
		PickupAddress:        ffm.PickupAddress,
		DeliveryPoints:       ffm.DeliveryPoints,
		ShippingNote:         ffm.ShippingNote,
		Coupon:               ffm.Coupon,
	}
	xShipnow, err := m.shipnowManager.CreateExternalShipnow(ctx, cmd)
	if err != nil {
		return err
	}

	cmd2 := &shipnow.UpdateShipnowFulfillmentCarrierInfoCommand{
		ID:                         ffm.ID,
		ShippingCode:               xShipnow.ID,
		ShippingState:              xShipnow.State,
		TotalFee:                   xShipnow.TotalFee,
		FeeLines:                   xShipnow.FeeLines,
		CarrierFeeLines:            xShipnow.FeeLines,
		ShippingCreatedAt:          xShipnow.CreatedAt,
		ShippingServiceName:        xShipnow.Service.Name,
		ShippingServiceDescription: xShipnow.Service.Description,
		ShippingSharedLink:         xShipnow.SharedLink,
	}
	if err := m.shipnow.Dispatch(ctx, cmd2); err != nil {
		return nil
	}
	return nil
}

func (m *ProcessManager) ExternalAccountAhamoveCreated(ctx context.Context, event *accountshipnow.ExternalAccountAhamoveCreatedEvent) error {
	query := &accountshipnow.GetExternalAccountAhamoveQuery{
		OwnerID: event.OwnerID,
		Phone:   event.Phone,
	}
	if err := m.accountshipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	accountAhamove := query.Result

	// Ahamove register account
	queryShop := &identity.GetShopByIDQuery{
		ID: event.ShopID,
	}
	if err := m.identityQuery.Dispatch(ctx, queryShop); err != nil {
		return err
	}
	queryAddress := &address.GetAddressByIDQuery{
		ID: queryShop.Result.AddressID,
	}
	if err := m.addressQuery.Dispatch(ctx, queryAddress); err != nil {
		return err
	}

	args2 := &carrier.RegisterExternalAccountCommand{
		Phone:        event.Phone,
		Name:         accountAhamove.Name,
		Address:      queryAddress.Result.GetFullAddress(),
		Carrier:      carriertypes.Ahamove,
		ConnectionID: event.ConnectionID,
		OwnerID:      event.OwnerID,
		ShopID:       event.ShopID,
	}
	regisResult, err := m.shipnowManager.RegisterExternalAccount(ctx, args2)
	if err != nil {
		return err
	}

	// create shop_connection
	shopConnCmd := &connectioning.CreateOrUpdateShopConnectionCommand{
		OwnerID:      event.OwnerID,
		ConnectionID: event.ConnectionID,
		Token:        regisResult.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			Identifier: event.Phone,
		},
	}
	if err := m.connectionAggr.Dispatch(ctx, shopConnCmd); err != nil {
		return err
	}

	// Re-get external account and Update external info
	xAccount, err := m.shipnowManager.GetExternalAccount(ctx, &carrier.GetExternalAccountCommand{
		OwnerID:      event.OwnerID,
		ConnectionID: event.ConnectionID,
	})
	if err != nil {
		return err
	}
	// update another info
	args3 := &accountshipnow.UpdateExternalAccountAhamoveExternalInfoCommand{
		ID:                accountAhamove.ID,
		ExternalID:        xAccount.ID,
		ExternalCreatedAt: xAccount.CreatedAt,
		ExternalVerified:  xAccount.Verified,
		ExternalToken:     regisResult.Token,
	}
	return m.accountshipnowAggr.Dispatch(ctx, args3)
}

func (m *ProcessManager) ExternalAccountAhamoveVerifyRequested(ctx context.Context, event *accountshipnow.ExternalAccountShipnowVerifyRequestedEvent) error {
	getXAccountArgs := &carrier.GetExternalAccountCommand{
		OwnerID:      event.OwnerID,
		ConnectionID: event.ConnectionID,
		ShopID:       event.ShopID,
	}
	xAccount, err := m.shipnowManager.GetExternalAccount(ctx, getXAccountArgs)
	if err != nil {
		return err
	}
	if xAccount.Verified {
		update := &accountshipnow.UpdateExternalAccountAhamoveExternalInfoCommand{
			ID:               event.ID,
			ExternalID:       xAccount.ID,
			ExternalVerified: xAccount.Verified,
		}
		if err := m.accountshipnowAggr.Dispatch(ctx, update); err != nil {
			return err
		}
		return nil
	}

	// send verify request to Ahamove
	cmd := &carrier.VerifyExternalAccountCommand{
		OwnerID:      event.OwnerID,
		ConnectionID: event.ConnectionID,
		ShopID:       event.ShopID,
	}
	res, err := m.shipnowManager.VerifyExternalAccount(ctx, cmd)
	if err != nil {
		return err
	}
	// update external_ticket_id
	externalData, _ := jsonx.Marshal(res)
	update := &accountshipnow.UpdateExternalAccountAhamoveExternalInfoCommand{
		ID:                   event.ID,
		ExternalTicketID:     res.TicketID,
		LastSendVerifiedAt:   time.Now(),
		ExternalDataVerified: externalData,
	}
	return m.accountshipnowAggr.Dispatch(ctx, update)
}

func (m *ProcessManager) ExternalAccountShipnowUpdateVerificationInfo(ctx context.Context, event *accountshipnow.ExternalAccountShipnowUpdateVerificationInfoEvent) error {
	query := &carrier.GetExternalAccountCommand{
		OwnerID:      event.OwnerID,
		ConnectionID: event.ConnectionID,
	}
	xAccount, err := m.shipnowManager.GetExternalAccount(ctx, query)
	if err != nil {
		return err
	}
	if !xAccount.Verified {
		return nil
	}

	update := &accountshipnow.UpdateExternalAccountAhamoveExternalInfoCommand{
		ID:               event.ID,
		ExternalVerified: xAccount.Verified,
	}
	return m.accountshipnowAggr.Dispatch(ctx, update)
}
