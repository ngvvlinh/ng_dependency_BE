package pm

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/customer_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	CustomerQS   customering.QueryBus
	CustomerAggr customering.CommandBus
	ShippingQS   shipping.QueryBus
}

func New(
	eventBus bus.EventRegistry, customerQS customering.QueryBus,
	customerAggr customering.CommandBus, shippingQS shipping.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:     eventBus,
		CustomerQS:   customerQS,
		CustomerAggr: customerAggr,
		ShippingQS:   shippingQS,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.FulfillmentFromImportCreated)
}

// create customer depends on phone
// if customer exists then ignore
func (m *ProcessManager) FulfillmentFromImportCreated(ctx context.Context, event *shipping.FulfillmentFromImportCreatedEvent) error {
	getFfmQuery := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
		ID: event.FulfillmentID,
	}
	if err := m.ShippingQS.Dispatch(ctx, getFfmQuery); err != nil {
		return err
	}

	ffm := getFfmQuery.Result
	addressTo := ffm.AddressTo
	if addressTo == nil {
		return nil
	}

	phoneNorm, ok := validate.NormalizePhone(ffm.AddressTo.Phone)
	if !ok {
		return nil
	}

	getCustomersByPhonesQuery := &customering.GetCustomerByPhoneQuery{
		ShopID: event.ShopID,
		Phone:  phoneNorm.String(),
	}
	if err := m.CustomerQS.Dispatch(ctx, getCustomersByPhonesQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return err
	}
	customer := getCustomersByPhonesQuery.Result
	// ignore if customer exists
	if customer != nil {
		return nil
	}

	createCustomerCmd := &customering.CreateCustomerCommand{
		ShopID:   event.ShopID,
		FullName: addressTo.FullName,
		Type:     customer_type.Individual,
		Phone:    phoneNorm.String(),
		Email:    addressTo.Email,
	}
	if err := m.CustomerAggr.Dispatch(ctx, createCustomerCmd); err != nil {
		return err
	}

	return nil
}
