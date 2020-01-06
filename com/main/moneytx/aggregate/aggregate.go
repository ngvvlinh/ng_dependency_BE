package aggregate

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipping"
	"etop.vn/backend/com/main/moneytx/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ moneytx.Aggregate = &MoneyTxAggregate{}

type MoneyTxAggregate struct {
	db                           cmsql.Transactioner
	moneyTxShippingStore         sqlstore.MoneyTxShippingStoreFactory
	moneyTxShippingExternalStore sqlstore.MoneyTxShippingExternalStoreFactory
	shippingQuery                shipping.QueryBus
	eventBus                     capi.EventBus
}

func (a *MoneyTxAggregate) CreateMoneyTxShipping(context.Context, moneytx.CreateMoneyTxShippingArgs) (*moneytx.MoneyTransactionShippingExtended, error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) CreateMoneyTxShippings(context.Context, *moneytx.CreateMoneyTxShippingsArgs) (created int, _ error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingInfo(context.Context, *moneytx.UpdateMoneyTxShippingInfoArgs) (*moneytx.MoneyTransactionShippingExtended, error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShipping(context.Context, *moneytx.ConfirmMoneyTxShippingArgs) (updated int, _ error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) DeleteMoneyTxShipping(context.Context, *moneytx.DeleteMoneyTxShippingArgs) (deleted int, _ error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) CreateMoneyTxShippingEtop(context.Context, *moneytx.CreateMoneyTxShippingEtopArgs) (*moneytx.MoneyTransactionShippingEtopExtended, error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingEtop(context.Context, moneytx.UpdateMoneyTxShippingEtopArgs) (*moneytx.MoneyTransactionShippingEtopExtended, error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShippingEtop(context.Context, *moneytx.ConfirmMoneyTxShippingEtopArgs) (updated int, _ error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) DeleteMoneyTxShippingEtop(ctx context.Context, ID dot.ID) (deleted int, _ error) {
	panic("implement me")
}

func NewMoneyTxAggregate(
	db *cmsql.Database,
	shippingQS shipping.QueryBus,
	eventB capi.EventBus,
) *MoneyTxAggregate {
	return &MoneyTxAggregate{
		db:                           db,
		moneyTxShippingStore:         sqlstore.NewMoneyTxShippingStore(db),
		moneyTxShippingExternalStore: sqlstore.NewMoneyTxShippingExternalStore(db),
		shippingQuery:                shippingQS,
		eventBus:                     eventB,
	}
}

func (a *MoneyTxAggregate) MessageBus() moneytx.CommandBus {
	b := bus.New()
	return moneytx.NewAggregateHandler(a).RegisterHandlers(b)
}

func (m *MoneyTxAggregate) AddFulfillmentMoneyTxShipping(context.Context, *moneytx.FfmMoneyTxShippingArgs) (updated int, _ error) {
	panic("implement me")
}

func (m *MoneyTxAggregate) RemoveFulfillmentMoneyTxShipping(context.Context, *moneytx.FfmMoneyTxShippingArgs) (removed int, _ error) {
	panic("implement me")
}

func (m *MoneyTxAggregate) ReCalcMoneyTxShipping(ctx context.Context, moneyTxShippingID dot.ID) error {
	if moneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Money Transaction ID")
	}
	// moneyTx, err := m.moneyTxShippingStore(ctx).ID(moneyTxShippingID).GetMoneyTxShipping()
	// if err != nil {
	// 	return err
	// }

	panic("implement me")
}
