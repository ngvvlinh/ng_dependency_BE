package haravan

import (
	"context"
	"encoding/json"

	"etop.vn/api/external/haravan"
	"etop.vn/api/external/haravan/identity"
	"etop.vn/api/meta"
	identitysqlstore "etop.vn/backend/com/external/haravan/identity/sqlstore"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/external/haravan/convert"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
	"etop.vn/common/bus"
)

var _ haravan.Aggregate = &Aggregate{}
var haravanPartner *model.Partner

type Aggregate struct {
	db                   cmsql.Transactioner
	haravanClient        *haravanclient.Client
	xAccountHaravanStore identitysqlstore.XAccountHaravanStoreFactory
}

func NewAggregate(db cmsql.Database, cfg haravanclient.Config) *Aggregate {
	ctx := context.Background()
	if partner, err := sqlstore.Partner(ctx).ID(haravan.HaravanPartnerID).Get(); err == nil {
		haravanPartner = partner
	}
	return &Aggregate{
		db:                   db,
		haravanClient:        haravanclient.New(cfg),
		xAccountHaravanStore: identitysqlstore.NewXAccountHaravanStore(db),
	}
}

func (a *Aggregate) MessageBus() haravan.CommandBus {
	b := bus.New()
	return haravan.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) SendUpdateExternalFulfillmentState(ctx context.Context, args *haravan.SendUpdateExternalFulfillmentStateArgs) (*meta.Empty, error) {
	xAccountHaravan, externalOrderID, externalFulfillmentID, err := a.validateUpdateExternalFulfillment(ctx, args.FulfillmentID)
	if err != nil {
		return nil, err
	}
	state, err := getExternalFulfillmentState(ctx, args.FulfillmentID)
	if err != nil {
		return nil, err
	}

	cmd := &haravanclient.UpdateShippingStateRequest{
		Connection: haravanclient.Connection{
			Subdomain: xAccountHaravan.Subdomain,
			TokenStr:  xAccountHaravan.AccessToken,
		},
		OrderID:   externalOrderID,
		FulfillID: externalFulfillmentID,
		State:     state,
	}
	if err = a.haravanClient.UpdateShippingState(ctx, cmd); err != nil {
		return nil, err
	}

	return &meta.Empty{}, nil
}

func (a *Aggregate) SendUpdateExternalPaymentStatus(ctx context.Context, args *haravan.SendUpdateExternalPaymentStatusArgs) (*meta.Empty, error) {
	xAccountHaravan, externalOrderID, externalFulfillmentID, err := a.validateUpdateExternalFulfillment(ctx, args.FulfillmentID)
	if err != nil {
		return nil, err
	}
	paymentStatus, err := getExternalPaymentStatus(ctx, args.FulfillmentID)
	if err != nil {
		return nil, err
	}

	cmd := &haravanclient.UpdatePaymentStatusRequest{
		Connection: haravanclient.Connection{
			Subdomain: xAccountHaravan.Subdomain,
			TokenStr:  xAccountHaravan.AccessToken,
		},
		OrderID:   externalOrderID,
		FulfillID: externalFulfillmentID,
		Status:    paymentStatus,
	}
	if err = a.haravanClient.UpdatePaymentStatus(ctx, cmd); err != nil {
		return nil, err
	}

	return &meta.Empty{}, nil
}

func (a *Aggregate) validateUpdateExternalFulfillment(ctx context.Context, fulfillmentID int64) (xAccountHaravan *identity.ExternalAccountHaravan, externalOrderID string, externalFulfillmentID string, _err error) {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: fulfillmentID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, "", "", err
	}
	ffm := query.Result.Fulfillment
	order := query.Result.Order

	if haravanPartner == nil || haravanPartner.ID != ffm.PartnerID {
		return nil, "", "", cm.Errorf(cm.InvalidArgument, nil, "This ffm does not come from Haravan")
	}

	var externalMeta haravan.ExternalMeta
	if err := json.Unmarshal(order.ExternalMeta, &externalMeta); err != nil {
		return nil, "", "", nil
	}

	if externalMeta.ExternalOrderID == "" || externalMeta.ExternalFulfillmentID == "" {
		return nil, "", "", cm.Errorf(cm.FailedPrecondition, nil, "Missing ExternalOrderID and ExternalFulfillmentID")
	}
	xAccountHaravan, _err = a.getExternalAccountHaravan(ctx, ffm.ShopID)
	if _err != nil {
		return nil, "", "", _err
	}
	return xAccountHaravan, externalMeta.ExternalOrderID, externalMeta.ExternalFulfillmentID, nil
}

func (a *Aggregate) getExternalAccountHaravan(ctx context.Context, shopID int64) (*identity.ExternalAccountHaravan, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(shopID).GetXAccountHaravan()
	if err != nil {
		return nil, err
	}
	if account == nil || account.AccessToken == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Missing Accesstoken External Account Haravan.")
	}
	return account, nil
}

func getExternalFulfillmentState(ctx context.Context, ffmID int64) (haravanclient.FulfillmentState, error) {
	query := &shipmodelx.GetFulfillmentQuery{
		FulfillmentID: ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return "", nil
	}
	ffm := query.Result
	state := convert.ToFulfillmentState(ffm.ShippingState)
	if state == "" {
		return "", cm.Errorf(cm.FailedPrecondition, nil, "Missing fulfillment state")
	}
	return state, nil
}

func getExternalPaymentStatus(ctx context.Context, ffmID int64) (haravanclient.PaymentStatus, error) {
	query := &shipmodelx.GetFulfillmentQuery{
		FulfillmentID: ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return "", nil
	}
	ffm := query.Result
	status := convert.ToCODStatus(ffm.EtopPaymentStatus)
	if status == "" {
		return "", cm.Errorf(cm.FailedPrecondition, nil, "Missing payment status")
	}
	return status, nil
}
