package pm

import (
	"context"
	"time"

	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/models/shipping"
)

var ll = l.New()
var _ shipping.ProcessManager = &Impl{}

type Impl struct {
	self     shipping.ProcessManagerBus
	location location.Bus

	db cmsql.Database
}

func New(
	locationBus location.Bus,
) *Impl {
	im := &Impl{
		location: locationBus,
	}
	im.self = shipping.ProcessManagerBus{bus.New()}
	im.RegisterHandlers(im.self)
	return im
}

func (im *Impl) MessageBus() shipping.ProcessManagerBus {
	return im.self
}

func (im *Impl) RegisterHandlers(bus bus.Bus) {
	bus.AddHandlers(
		im.CreateFulfillmentHandler,
		im.ConfirmFulfillmentHandler,
		im.CancelFulfillmentHandler,
	)
}

func (im *Impl) CreateFulfillmentHandler(ctx context.Context, cmd *shipping.CreateFulfillmentCommand) error {
	return im.CreateFulfillment(ctx, cmd)
}

func (im *Impl) CreateFulfillment(ctx context.Context, cmd *shipping.CreateFulfillmentCommand) error {

	ffm := &model.Fulfillment{
		ID:                                 cm.NewID(),
		OrderID:                            0,
		ShopID:                             0,
		SupplierID:                         0,
		PartnerID:                          0,
		SupplierConfirm:                    0,
		ShopConfirm:                        0,
		ConfirmStatus:                      0,
		TotalItems:                         0,
		TotalWeight:                        0,
		BasketValue:                        0,
		TotalDiscount:                      0,
		TotalAmount:                        0,
		TotalCODAmount:                     0,
		OriginalCODAmount:                  0,
		ActualCompensationAmount:           0,
		ShippingFeeCustomer:                0,
		ShippingFeeShop:                    0,
		ShippingFeeShopLines:               nil,
		ShippingServiceFee:                 0,
		ExternalShippingFee:                0,
		ProviderShippingFeeLines:           nil,
		EtopDiscount:                       0,
		EtopFeeAdjustment:                  0,
		ShippingFeeMain:                    0,
		ShippingFeeReturn:                  0,
		ShippingFeeInsurance:               0,
		ShippingFeeAdjustment:              0,
		ShippingFeeCODS:                    0,
		ShippingFeeInfoChange:              0,
		ShippingFeeOther:                   0,
		EtopAdjustedShippingFeeMain:        0,
		EtopPriceRule:                      false,
		VariantIDs:                         nil,
		Lines:                              nil,
		TypeFrom:                           "",
		TypeTo:                             "",
		AddressFrom:                        nil,
		AddressTo:                          nil,
		AddressReturn:                      nil,
		AddressToProvinceCode:              "",
		AddressToDistrictCode:              "",
		AddressToWardCode:                  "",
		CreatedAt:                          time.Time{},
		UpdatedAt:                          time.Time{},
		ClosedAt:                           time.Time{},
		ExpectedDeliveryAt:                 time.Time{},
		ExpectedPickAt:                     time.Time{},
		CODEtopTransferedAt:                time.Time{},
		ShippingFeeShopTransferedAt:        time.Time{},
		ShippingCancelledAt:                time.Time{},
		ShippingDeliveredAt:                time.Time{},
		ShippingReturnedAt:                 time.Time{},
		ShippingCreatedAt:                  time.Time{},
		ShippingPickingAt:                  time.Time{},
		ShippingHoldingAt:                  time.Time{},
		ShippingDeliveringAt:               time.Time{},
		ShippingReturningAt:                time.Time{},
		MoneyTransactionID:                 0,
		MoneyTransactionShippingExternalID: 0,
		CancelReason:                       "",
		ShippingProvider:                   "",
		ProviderServiceID:                  "",
		ShippingCode:                       "",
		ShippingNote:                       "",
		TryOn:                              "",
		IncludeInsurance:                   false,
		ExternalShippingName:               "",
		ExternalShippingID:                 "",
		ExternalShippingCode:               "",
		ExternalShippingCreatedAt:          time.Time{},
		ExternalShippingUpdatedAt:          time.Time{},
		ExternalShippingCancelledAt:        time.Time{},
		ExternalShippingDeliveredAt:        time.Time{},
		ExternalShippingReturnedAt:         time.Time{},
		ExternalShippingClosedAt:           time.Time{},
		ExternalShippingState:              "",
		ExternalShippingStateCode:          "",
		ExternalShippingStatus:             0,
		ExternalShippingNote:               "",
		ExternalShippingSubState:           "",
		ExternalShippingData:               nil,
		ShippingState:                      "",
		ShippingStatus:                     0,
		EtopPaymentStatus:                  0,
		Status:                             0,
		SyncStatus:                         0,
		SyncStates:                         nil,
		LastSyncAt:                         time.Time{},
		ExternalShippingLogs:               nil,
		AdminNote:                          "",
		IsPartialDelivery:                  false,
	}
	return cm.ErrTODO
}

func (im *Impl) ConfirmFulfillmentHandler(ctx context.Context, cmd *shipping.ConfirmFulfillmentCommand) error {
	return im.ConfirmFulfillment(ctx, cmd)
}

func (im *Impl) ConfirmFulfillment(ctx context.Context, cmd *shipping.ConfirmFulfillmentCommand) error {
	return cm.ErrTODO
}

func (im *Impl) CancelFulfillmentHandler(ctx context.Context, cmd *shipping.CancelFulfillmentCommand) error {
	return im.CancelFulfillment(ctx, cmd)
}

func (im *Impl) CancelFulfillment(ctx context.Context, cmd *shipping.CancelFulfillmentCommand) error {
	return cm.ErrTODO
}

func (im *Impl) GetFulfillmentByIDHandler(ctx context.Context, query *shipping.GetFulfillmentByIDQuery) error {
	result, err := im.GetFulfillmentByID(ctx, query)
	query.Result = result
	return err
}

func (im *Impl) GetFulfillmentByID(ctx context.Context, query *shipping.GetFulfillmentByIDQuery) (*model.Fulfillment, error) {
	return nil, cm.ErrTODO
}
