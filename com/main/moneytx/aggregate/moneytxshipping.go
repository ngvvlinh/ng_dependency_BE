package aggregate

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	identityconvert "o.o/backend/com/main/identity/convert"
	moneytxmodel "o.o/backend/com/main/moneytx/model"
	moneytxsqlstore "o.o/backend/com/main/moneytx/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ moneytx.Aggregate = &MoneyTxAggregate{}

type MoneyTxAggregate struct {
	db                           *cmsql.Database
	moneyTxShippingStore         moneytxsqlstore.MoneyTxShippingStoreFactory
	moneyTxShippingExternalStore moneytxsqlstore.MoneyTxShippingExternalStoreFactory
	moneyTxShippingEtopStore     moneytxsqlstore.MoneyTxShippingEtopStoreFactory
	shippingQuery                shipping.QueryBus
	identityQuery                identity.QueryBus
	eventBus                     capi.EventBus
}

func NewMoneyTxAggregate(
	db com.MainDB,
	shippingQS shipping.QueryBus,
	identityQS identity.QueryBus,
	eventB capi.EventBus,
) *MoneyTxAggregate {
	return &MoneyTxAggregate{
		db:                           db,
		moneyTxShippingStore:         moneytxsqlstore.NewMoneyTxShippingStore(db),
		moneyTxShippingExternalStore: moneytxsqlstore.NewMoneyTxShippingExternalStore(db),
		moneyTxShippingEtopStore:     moneytxsqlstore.NewMoneyTxShippingEtopStore(db),
		shippingQuery:                shippingQS,
		identityQuery:                identityQS,
		eventBus:                     eventB,
	}
}

func MoneyTxAggregateMessageBus(a *MoneyTxAggregate) moneytx.CommandBus {
	b := bus.New()
	return moneytx.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *MoneyTxAggregate) CreateMoneyTxShipping(ctx context.Context, args *moneytx.CreateMoneyTxShippingArgs) (*moneytx.MoneyTransactionShipping, error) {
	if args.Shop == nil {
		if args.ShopID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
		}
		query := &identity.GetShopByIDQuery{
			ID: args.ShopID,
		}
		if err := a.identityQuery.Dispatch(ctx, query); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Shop ID kh??ng h???p l???")
		}
		args.Shop = query.Result
	}
	if len(args.FulfillmentIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "FulfillmentIDs can not be empty")
	}

	id := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		code, err := sqlstore.GenerateCode(ctx, tx, model.CodeTypeMoneyTransaction, args.Shop.Code)
		if err != nil {
			return err
		}
		query := &shipping.ListFulfillmentsByIDsQuery{
			IDs:    args.FulfillmentIDs,
			ShopID: args.Shop.ID,
		}
		if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		fulfillments := query.Result
		var ffms []*shipping.Fulfillment
		for _, id := range args.FulfillmentIDs {
			found := false
			for _, ffm := range fulfillments {
				if id == ffm.ID {
					if err := checkFulfillmentValid(ffm); err != nil {
						return err
					}
					found = true
					ffms = append(ffms, ffm)
					break
				}
			}
			if !found {
				return cm.Errorf(cm.NotFound, nil, "Fulfillment id (%v) not found", id)
			}
		}
		statistics, err := calcFulfillmentStatistics(ffms)
		if err != nil {
			return err
		}
		if statistics.TotalCOD != args.TotalCOD {
			return cm.Errorf(cm.FailedPrecondition, nil, "Total COD does not match")
		}
		if statistics.TotalAmount != args.TotalAmount {
			return cm.Errorf(cm.FailedPrecondition, nil, "Total Amount does not match")
		}
		if statistics.TotalOrders != args.TotalOrders {
			return cm.Errorf(cm.FailedPrecondition, nil, "Total Orders does not match")
		}

		moneyTx := &moneytx.MoneyTransactionShipping{
			ID:          id,
			ShopID:      args.Shop.ID,
			Status:      status3.Z,
			Code:        code,
			TotalCOD:    statistics.TotalCOD,
			TotalAmount: statistics.TotalAmount,
			TotalOrders: statistics.TotalOrders,
		}
		if err := a.moneyTxShippingStore(ctx).CreateMoneyTxShipping(moneyTx); err != nil {
			return err
		}

		event := &moneytx.MoneyTxShippingCreatedEvent{
			EventMeta:         meta.NewEvent(),
			MoneyTxShippingID: moneyTx.ID,
			ShopID:            moneyTx.ShopID,
			FulfillmentIDs:    statistics.FulfillmentIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingStore(ctx).ID(id).GetMoneyTxShipping()
}

func (a *MoneyTxAggregate) CreateMoneyTxShippings(ctx context.Context, args *moneytx.CreateMoneyTxShippingsArgs) (created int, _ error) {
	shopIDs := []dot.ID{}
	for shopID, _ := range args.ShopIDMapFfms {
		shopIDs = append(shopIDs, shopID)
	}

	query := &identity.ListShopsByIDsQuery{
		IDs: shopIDs,
	}
	if err := a.identityQuery.Dispatch(ctx, query); err != nil {
		return 0, err
	}
	if len(query.Result) != len(args.ShopIDMapFfms) {
		return 0, cm.Errorf(cm.Internal, nil, "ShopIDs does not have the expected length").WithMetap("func", "CreateMoneyTxShippings")
	}
	shopIDMapShop := make(map[dot.ID]*identity.Shop)
	for _, shop := range query.Result {
		shopIDMapShop[shop.ID] = shop
	}

	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for shopID, ffms := range args.ShopIDMapFfms {
			statistics, err := calcFulfillmentStatistics(ffms)
			if err != nil {
				return err
			}
			cmd := &moneytx.CreateMoneyTxShippingArgs{
				Shop:           shopIDMapShop[shopID],
				FulfillmentIDs: statistics.FulfillmentIDs,
				TotalCOD:       statistics.TotalCOD,
				TotalAmount:    statistics.TotalAmount,
				TotalOrders:    statistics.TotalOrders,
			}
			if _, err := a.CreateMoneyTxShipping(ctx, cmd); err != nil {
				return err
			}
			created++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingInfo(ctx context.Context, args *moneytx.UpdateMoneyTxShippingInfoArgs) (*moneytx.MoneyTransactionShipping, error) {
	moneyTx, err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).OptionalShopID(args.ShopID).GetMoneyTxShipping()
	if err != nil {
		return nil, err
	}
	if moneyTx.Status == status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "This money transaction was confirmed. Can not update.")
	}
	update := &moneytxmodel.MoneyTransactionShipping{
		ID:            args.MoneyTxShippingID,
		Note:          args.Note,
		InvoiceNumber: args.InvoiceNumber,
		BankAccount:   identityconvert.Convert_identitytypes_BankAccount_sharemodel_BankAccount(args.BankAccount, nil),
	}
	if err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).UpdateMoneyTxShippingDB(update); err != nil {
		return nil, err
	}
	return a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).GetMoneyTxShipping()
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShipping(ctx context.Context, args *moneytx.ConfirmMoneyTxShippingArgs) error {
	if args.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing money transaction id")
	}
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop id")
	}
	moneyTx, err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).ShopID(args.ShopID).GetMoneyTxShipping()
	if err != nil {
		return err
	}
	if moneyTx.Status != status3.Z {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this money transaction")
	}

	query := &shipping.ListFulfillmentsByMoneyTxQuery{
		MoneyTxShippingIDs: []dot.ID{args.MoneyTxShippingID},
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffms := query.Result
	statistics, err := calcFulfillmentStatistics(ffms)
	if err != nil {
		return err
	}
	if args.TotalCOD != statistics.TotalCOD {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total COD does not match")
	}
	if args.TotalAmount != statistics.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match")
	}
	if args.TotalOrders != statistics.TotalOrders {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total orders does not match")
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		now := time.Now()
		update1 := &moneytxsqlstore.UpdateMoneyTxShippingStatisticsArgs{
			ID:          args.MoneyTxShippingID,
			TotalOrders: dot.Int(statistics.TotalOrders),
			TotalCOD:    dot.Int(statistics.TotalCOD),
			TotalAmount: dot.Int(statistics.TotalAmount),
		}
		if err := a.moneyTxShippingStore(ctx).UpdateMoneyTxShippingStatistics(update1); err != nil {
			return err
		}

		update2 := &moneytxmodel.MoneyTransactionShipping{
			Status:           status3.P,
			ConfirmedAt:      now,
			EtopTransferedAt: now,
		}
		if err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).UpdateMoneyTxShippingDB(update2); err != nil {
			return err
		}

		event := &moneytx.MoneyTxShippingConfirmedEvent{
			EventMeta:         meta.NewEvent(),
			MoneyTxShippingID: args.MoneyTxShippingID,
			ShopID:            args.ShopID,
			ConfirmedAt:       now,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
}

func (a *MoneyTxAggregate) DeleteMoneyTxShipping(ctx context.Context, args *moneytx.DeleteMoneyTxShippingArgs) error {
	if args.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing money_tx_shipping_id")
	}
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}
	moneyTx, err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).ShopID(args.ShopID).GetMoneyTxShipping()
	if err != nil {
		return err
	}
	if moneyTx.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Cannot delete this money transaction")
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		event := &moneytx.MoneyTxShippingDeletedEvent{
			EventMeta:         meta.NewEvent(),
			MoneyTxShippingID: args.MoneyTxShippingID,
			ShopID:            args.ShopID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return a.moneyTxShippingStore(ctx).DeleteMoneyTxShipping(args.MoneyTxShippingID)
	})
}

func (a *MoneyTxAggregate) AddFulfillmentsMoneyTxShipping(context.Context, *moneytx.FfmsMoneyTxShippingArgs) error {
	panic("implement me")
}

func (a *MoneyTxAggregate) RemoveFulfillmentsMoneyTxShipping(ctx context.Context, args *moneytx.FfmsMoneyTxShippingArgs) error {
	if args.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing money_tx_shipping_id")
	}
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing shop_id")
	}
	if len(args.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "FulfillmentIDs can not be empty")
	}
	moneyTx, err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).ShopID(args.ShopID).GetMoneyTxShipping()
	if err != nil {
		return err
	}
	if moneyTx.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Cannot update this money transaction")
	}
	query := &shipping.ListFulfillmentsByMoneyTxQuery{
		MoneyTxShippingIDs: []dot.ID{args.MoneyTxShippingID},
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffms := query.Result
	ffmsMap := make(map[dot.ID]*shipping.Fulfillment)
	var retainFfms []*shipping.Fulfillment
	for _, ffm := range ffms {
		ffmsMap[ffm.ID] = ffm
		if !cm.IDsContain(args.FulfillmentIDs, ffm.ID) {
			retainFfms = append(retainFfms, ffm)
		}
	}

	for _, id := range args.FulfillmentIDs {
		_, ok := ffmsMap[id]
		if !ok {
			return cm.Errorf(cm.NotFound, nil, "Fulfillment #%v does not exist in this money transaction", id)
		}
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		event := &moneytx.MoneyTxShippingRemovedFfmsEvent{
			EventMeta:         meta.NewEvent(),
			MoneyTxShippingID: args.MoneyTxShippingID,
			FulfillmentIDs:    args.FulfillmentIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		statistics, err := calcFulfillmentStatistics(retainFfms)
		if err != nil {
			return err
		}
		update := &moneytxsqlstore.UpdateMoneyTxShippingStatisticsArgs{
			ID:          args.MoneyTxShippingID,
			TotalOrders: dot.Int(statistics.TotalOrders),
			TotalCOD:    dot.Int(statistics.TotalCOD),
			TotalAmount: dot.Int(statistics.TotalAmount),
		}
		if err := a.moneyTxShippingStore(ctx).UpdateMoneyTxShippingStatistics(update); err != nil {
			return err
		}
		return nil
	})
}

func (a *MoneyTxAggregate) ReCalcMoneyTxShipping(ctx context.Context, args *moneytx.ReCalcMoneyTxShippingArgs) error {
	if args.MoneyTxShippingID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Money Transaction ID")
	}

	moneyTxShipping, err := a.moneyTxShippingStore(ctx).ID(args.MoneyTxShippingID).GetMoneyTxShipping()
	if err != nil {
		return err
	}
	if moneyTxShipping.Status == status3.P || !moneyTxShipping.ConfirmedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Phi??n ?????i so??t ???? x??c nh???n.").WithMetap("money_transaction_id", moneyTxShipping.ID)
	}
	query := &shipping.ListFulfillmentsByMoneyTxQuery{
		MoneyTxShippingIDs: []dot.ID{args.MoneyTxShippingID},
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffms := query.Result
	statistics, err := calcFulfillmentStatistics(ffms)
	if err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		update := &moneytxsqlstore.UpdateMoneyTxShippingStatisticsArgs{
			ID:          args.MoneyTxShippingID,
			TotalOrders: dot.Int(statistics.TotalOrders),
			TotalCOD:    dot.Int(statistics.TotalCOD),
			TotalAmount: dot.Int(statistics.TotalAmount),
		}
		if err := a.moneyTxShippingStore(ctx).UpdateMoneyTxShippingStatistics(update); err != nil {
			return err
		}
		if moneyTxShipping.MoneyTransactionShippingEtopID != 0 {
			return a.ReCalcMoneyTxShippingEtop(ctx, moneyTxShipping.MoneyTransactionShippingEtopID)
		}
		return nil
	})
}

func calcFulfillmentStatistics(fulfillments []*shipping.Fulfillment) (*moneytx.FulfilmentStatistics, error) {
	var totalCOD, totalAmount, totalOrders, totalShippingFee int
	var ffmIDs []dot.ID
	var ffmIDMap = make(map[dot.ID]bool)
	for _, ffm := range fulfillments {
		if !cm.StringsContain(moneytx.ShippingAcceptStates, ffm.ShippingState.String()) {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Fulfillment #%v's status does not valid.", ffm.ShippingCode)
		}

		if _, ok := ffmIDMap[ffm.ID]; ok {
			return nil, cm.Errorf(cm.Internal, nil, "Dupplicate fulfillment when calculate statistics").WithMetap("ffmID", ffm.ID)
		}
		ffmIDs = append(ffmIDs, ffm.ID)
		amount := ffm.TotalCODAmount
		if ffm.ShippingState == shippingstate.Returned || ffm.ShippingState == shippingstate.Returning {
			// make sure COD = 0
			amount = 0
		} else if ffm.ShippingState == shippingstate.Undeliverable {
			// tr?????ng h???p ????n b???i ho??n
			amount = ffm.ActualCompensationAmount
		}
		totalAmount = totalAmount + amount - ffm.ShippingFeeShop
		totalCOD += amount
		totalOrders++
		totalShippingFee += ffm.ShippingFeeShop
	}
	return &moneytx.FulfilmentStatistics{
		TotalCOD:         totalCOD,
		TotalAmount:      totalAmount,
		TotalOrders:      totalOrders,
		TotalShippingFee: totalShippingFee,
		FulfillmentIDs:   ffmIDs,
	}, nil
}

func checkFulfillmentValid(ffm *shipping.Fulfillment) error {
	if !cm.StringsContain(moneytx.ShippingAcceptStates, ffm.ShippingState.String()) {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" does not valid. Status must be delivered or returning or returned.", nil)
	}
	if ffm.MoneyTransactionID != 0 {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" in another money transaction.", nil)
	}
	if !ffm.CODEtopTransferedAt.IsZero() {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" has paid.", nil)
	}
	// backward compatible
	// remove later
	if ffm.ShippingType == 0 && ffm.ConnectionMethod == 0 {
		return nil
	}
	// -- end backward compatible

	if ffm.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		return cm.Errorf(cm.FailedPrecondition, nil, "Fulfillment #%v can not be paid by Etop", ffm.ShippingCode)
	}
	return nil
}
