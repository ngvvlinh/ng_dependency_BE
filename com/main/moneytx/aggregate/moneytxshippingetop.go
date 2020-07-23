package aggregate

import (
	"context"
	"time"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	identityconvert "o.o/backend/com/main/identity/convert"
	moneytxmodel "o.o/backend/com/main/moneytx/model"
	moneytxsqlstore "o.o/backend/com/main/moneytx/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

func (a *MoneyTxAggregate) CreateMoneyTxShippingEtop(ctx context.Context, args *moneytx.CreateMoneyTxShippingEtopArgs) (*moneytx.MoneyTransactionShippingEtop, error) {
	if len(args.MoneyTxShippingIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "MoneyTxShippingIDs can not be empty")
	}
	statistics, err := a.prepareMoneyTxShippingEtop(ctx, 0, args.MoneyTxShippingIDs)
	if err != nil {
		return nil, err
	}
	newID := cm.NewID()
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		code, err := sqlstore.GenerateCode(ctx, tx, model.CodeTypeMoneyTransaction, "ETOP")
		if err != nil {
			return err
		}
		moneyTxShippingEtop := &moneytxmodel.MoneyTransactionShippingEtop{
			ID:                    newID,
			TotalCOD:              statistics.TotalCOD,
			TotalAmount:           statistics.TotalAmount,
			TotalOrders:           statistics.TotalOrders,
			TotalFee:              statistics.TotalShippingFee,
			TotalMoneyTransaction: statistics.TotalMoneyTransaction,
			Code:                  code,
			Note:                  args.Note,
			InvoiceNumber:         args.Note,
			BankAccount:           identityconvert.Convert_identitytypes_BankAccount_sharemodel_BankAccount(args.BankAccount, nil),
		}
		if err := a.moneyTxShippingEtopStore(ctx).CreateMoneyTxShippingEtopDB(moneyTxShippingEtop); err != nil {
			return err
		}
		update := &moneytxmodel.MoneyTransactionShipping{
			MoneyTransactionShippingEtopID: newID,
		}
		if err := a.moneyTxShippingStore(ctx).IDs(args.MoneyTxShippingIDs...).UpdateMoneyTxShippingDB(update); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingEtopStore(ctx).ID(newID).GetMoneyTxShippingEtop()
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingEtop(ctx context.Context, args moneytx.UpdateMoneyTxShippingEtopArgs) (*moneytx.MoneyTransactionShippingEtop, error) {
	if args.MoneyTxShippingEtopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTransactionShippingEtopID")
	}
	mtxse, err := a.moneyTxShippingEtopStore(ctx).ID(args.MoneyTxShippingEtopID).GetMoneyTxShippingEtop()
	if err != nil {
		return nil, err
	}
	if mtxse.Status == status3.P {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "This MoneyTransactionShippingEtop was confirmed. Can not update")
	}

	moneyTxShippings, err := a.moneyTxShippingStore(ctx).MoneyTxShippingEtopID(args.MoneyTxShippingEtopID).ListMoneyTxShippings()
	if err != nil {
		return nil, err
	}
	oldIDs := make([]dot.ID, len(moneyTxShippings))
	for i, mtxs := range moneyTxShippings {
		oldIDs[i] = mtxs.ID
	}
	newIDs := patchID(oldIDs, args.Adds, args.Deletes, args.ReplaceAll)

	statistics, err := a.prepareMoneyTxShippingEtop(ctx, args.MoneyTxShippingEtopID, newIDs)
	if err != nil {
		return nil, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if len(oldIDs) > 0 {
			if err := a.moneyTxShippingStore(ctx).IDs(oldIDs...).RemoveMoneyTxShippingMoneyTxShippingEtopID(); err != nil {
				return err
			}
		}

		// update phiên etop
		updateStatistics := &moneytxsqlstore.UpdateMoneyTxShippingEtopStatisticsArgs{
			ID:                    args.MoneyTxShippingEtopID,
			TotalCOD:              dot.Int(statistics.TotalCOD),
			TotalAmount:           dot.Int(statistics.TotalAmount),
			TotalOrders:           dot.Int(statistics.TotalOrders),
			TotalFee:              dot.Int(statistics.TotalShippingFee),
			TotalMoneyTransaction: dot.Int(statistics.TotalOrders),
		}
		if err := a.moneyTxShippingEtopStore(ctx).UpdateMoneyTxShippingEtopStatistics(updateStatistics); err != nil {
			return err
		}

		update := &moneytxmodel.MoneyTransactionShippingEtop{
			ID:            args.MoneyTxShippingEtopID,
			BankAccount:   identityconvert.Convert_identitytypes_BankAccount_sharemodel_BankAccount(args.BankAccount, nil),
			Note:          args.Note,
			InvoiceNumber: args.InvoiceNumber,
		}
		if err := a.moneyTxShippingEtopStore(ctx).ID(args.MoneyTxShippingEtopID).UpdateMoneyTxShippingEtopDB(update); err != nil {
			return err
		}

		if len(newIDs) > 0 {
			updateMoneyTxShipping := &moneytxmodel.MoneyTransactionShipping{
				MoneyTransactionShippingEtopID: args.MoneyTxShippingEtopID,
			}
			if err := a.moneyTxShippingStore(ctx).IDs(newIDs...).UpdateMoneyTxShippingDB(updateMoneyTxShipping); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingEtopStore(ctx).ID(args.MoneyTxShippingEtopID).GetMoneyTxShippingEtop()
}

func patchID(list, adds, deletes, replaceAll []dot.ID) []dot.ID {
	if len(replaceAll) > 0 {
		return replaceAll
	}
	newList := make([]dot.ID, 0, len(list)+len(adds))
	for _, id := range list {
		if !cm.IDsContain(newList, id) && !cm.IDsContain(deletes, id) {
			newList = append(newList, id)
		}
	}
	for _, id := range adds {
		if !cm.IDsContain(newList, id) && !cm.IDsContain(deletes, id) {
			newList = append(newList, id)
		}
	}
	return newList
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShippingEtop(ctx context.Context, args *moneytx.ConfirmMoneyTxShippingEtopArgs) error {
	if args.MoneyTxShippingEtopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingEtopID")
	}
	moneyTxShippingEtop, err := a.moneyTxShippingEtopStore(ctx).ID(args.MoneyTxShippingEtopID).GetMoneyTxShippingEtop()
	if err != nil {
		return err
	}
	if moneyTxShippingEtop.Status != status3.Z {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this MoneyTxShippingEtop")
	}

	moneyTxShippings, err := a.moneyTxShippingStore(ctx).MoneyTxShippingEtopID(args.MoneyTxShippingEtopID).ListMoneyTxShippings()
	if err != nil {
		return err
	}

	var totalCOD, totalAmount, totalOrders int
	var isContainMoneyTxManual bool
	moneyTxShippingIDs := make([]dot.ID, len(moneyTxShippings))
	for i, mtxs := range moneyTxShippings {
		if mtxs.Status != status3.Z {
			return cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this MoneyTxShipping (id = %v)", mtxs.ID)
		}
		moneyTxShippingIDs[i] = mtxs.ID
		totalCOD += mtxs.TotalCOD
		totalAmount += mtxs.TotalAmount
		totalOrders += mtxs.TotalOrders
		if mtxs.Type == "manual" {
			isContainMoneyTxManual = true
		}
	}

	query := &shipping.ListFulfillmentsByMoneyTxQuery{
		MoneyTxShippingIDs: moneyTxShippingIDs,
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	ffms := query.Result
	ffmIDs := make([]dot.ID, len(ffms))
	for i, ffm := range ffms {
		if !cm.StringsContain(moneytx.ShippingAcceptStates, ffm.ShippingState.Text()) {
			return cm.Errorf(cm.FailedPrecondition, nil, "Fulfillment's (#%v) shipping state does not valid.")
		}
		ffmIDs[i] = ffm.ID
	}
	if totalCOD != args.TotalCOD {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total COD does not match. (expected_total_cod = %v)", totalCOD)
	}
	if totalAmount != args.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total Amount does not match. (expected_total_amount = %v", totalAmount)
	}
	if totalOrders != args.TotalOrders {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total Orders does not match. (expected_total_orders = %v", totalOrders)
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		now := time.Now()

		// cập nhật phiên shop
		for _, mtxs := range moneyTxShippings {
			update := &moneytxmodel.MoneyTransactionShipping{
				Status:           status3.P,
				ConfirmedAt:      now,
				EtopTransferedAt: now,
			}
			if err := a.moneyTxShippingStore(ctx).ID(mtxs.ID).UpdateMoneyTxShippingDB(update); err != nil {
				return err
			}
		}

		// cập nhật phiên etop
		mtxseUpdateStatistic := &moneytxsqlstore.UpdateMoneyTxShippingEtopStatisticsArgs{
			ID:                    args.MoneyTxShippingEtopID,
			TotalCOD:              dot.Int(totalCOD),
			TotalAmount:           dot.Int(totalAmount),
			TotalOrders:           dot.Int(totalOrders),
			TotalMoneyTransaction: dot.Int(len(moneyTxShippings)),
		}
		if !isContainMoneyTxManual {
			mtxseUpdateStatistic.TotalFee = dot.Int(totalCOD - totalAmount)
		}
		if err := a.moneyTxShippingEtopStore(ctx).UpdateMoneyTxShippingEtopStatistics(mtxseUpdateStatistic); err != nil {
			return err
		}
		mtxseUpdate := &moneytxmodel.MoneyTransactionShippingEtop{
			ConfirmedAt: now,
			Status:      status3.P,
		}
		if err := a.moneyTxShippingEtopStore(ctx).ID(args.MoneyTxShippingEtopID).UpdateMoneyTxShippingEtopDB(mtxseUpdate); err != nil {
			return err
		}

		// update ffms cod_etop_transfered_at
		event := &moneytx.MoneyTxShippingEtopConfirmedEvent{
			EventMeta:             meta.NewEvent(),
			MoneyTxShippingEtopID: args.MoneyTxShippingEtopID,
			MoneyTxShippingIDs:    moneyTxShippingIDs,
			ConfirmedAt:           now,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
}

func (a *MoneyTxAggregate) DeleteMoneyTxShippingEtop(ctx context.Context, id dot.ID) error {
	if id == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxShippingEtopID")
	}
	moneyTxShippingEtop, err := a.moneyTxShippingEtopStore(ctx).ID(id).GetMoneyTxShippingEtop()
	if err != nil {
		return err
	}
	if moneyTxShippingEtop.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "MoneyTxShippingEtop was confirmed. Can not delete.")
	}
	moneyTxShippings, err := a.moneyTxShippingStore(ctx).MoneyTxShippingEtopID(id).ListMoneyTxShippings()
	if err != nil {
		return err
	}

	var moneyTxShippingIDs = make([]dot.ID, len(moneyTxShippings))
	for i, mtxs := range moneyTxShippings {
		if mtxs.Status == status3.P {
			return cm.Errorf(cm.FailedPrecondition, nil, "Can not delete this MoneyTxShippingEtop. MoneyTxShipping (id = %v) was confirmed", mtxs.ID)
		}
		moneyTxShippingIDs[i] = mtxs.ID
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if len(moneyTxShippingIDs) > 0 {
			if err := a.moneyTxShippingStore(ctx).IDs(moneyTxShippingIDs...).RemoveMoneyTxShippingMoneyTxShippingEtopID(); err != nil {
				return err
			}
		}
		if err := a.moneyTxShippingEtopStore(ctx).DeleteMoneyTxShippingEtop(id); err != nil {
			return err
		}
		return nil
	})
}

type MoneyTxShippingEtopStatistics struct {
	TotalOrders           int
	TotalCOD              int
	TotalAmount           int
	TotalShippingFee      int
	TotalMoneyTransaction int
}

func (a *MoneyTxAggregate) prepareMoneyTxShippingEtop(ctx context.Context, moneyTxShippingEtopID dot.ID, moneyTxShippingIDs []dot.ID) (*MoneyTxShippingEtopStatistics, error) {
	if len(moneyTxShippingIDs) == 0 {
		return &MoneyTxShippingEtopStatistics{}, nil
	}
	moneyTxShippings, err := a.moneyTxShippingStore(ctx).IDs(moneyTxShippingIDs...).ListMoneyTxShippings()
	if err != nil {
		return nil, err
	}
	for _, id := range moneyTxShippingIDs {
		found := false
		for _, mtxs := range moneyTxShippings {
			if id == mtxs.ID {
				found = true
				break
			}
		}
		if !found {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "MoneyTxShipping does not exist (id = %v)", id)
		}
	}
	if len(moneyTxShippingIDs) != len(moneyTxShippings) {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}

	query := &shipping.ListFulfillmentsByMoneyTxQuery{
		MoneyTxShippingIDs: moneyTxShippingIDs,
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	ffms := query.Result

	// Tuan hotfix: 23/07/2020
	totalAmountManual := 0

	for _, mtxs := range moneyTxShippings {
		if mtxs.Status != status3.Z {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "MoneyTxShipping does not valid (id = %v)", mtxs.ID)
		}
		if mtxs.MoneyTransactionShippingEtopID != 0 && mtxs.MoneyTransactionShippingEtopID != moneyTxShippingEtopID {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "MoneyTxShipping belongs to another MoneyTxShippingEtop. (money_tx_shipping_id = %v, money_tx_shipping_etop_id = %v", mtxs.ID, mtxs.MoneyTransactionShippingEtopID)
		}
		// Tuan hotfix: 23/07/2020
		if mtxs.Type == "manual" {
			totalAmountManual += mtxs.TotalAmount
		}
	}

	var totalCOD, totalOrders, totalAmount int
	for _, mtxs := range moneyTxShippings {
		totalCOD += mtxs.TotalCOD
		totalOrders += mtxs.TotalOrders
		totalAmount += mtxs.TotalAmount
	}
	statistics, err := calcFulfillmentStatistics(ffms)
	if err != nil {
		return nil, err
	}

	if statistics.TotalCOD != totalCOD {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Total COD does not match")
	}

	// Tuan hotfix: 23/07/2020
	statistics.TotalAmount += totalAmountManual

	if statistics.TotalAmount != totalAmount {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Total Amount does not match")
	}
	if statistics.TotalOrders != totalOrders {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Total Orders does not match")
	}

	return &MoneyTxShippingEtopStatistics{
		TotalCOD:              statistics.TotalCOD,
		TotalOrders:           statistics.TotalOrders,
		TotalAmount:           statistics.TotalAmount,
		TotalShippingFee:      statistics.TotalShippingFee,
		TotalMoneyTransaction: len(moneyTxShippingIDs),
	}, nil
}

func (a *MoneyTxAggregate) ReCalcMoneyTxShippingEtop(ctx context.Context, MoneyTxShippingEtopID dot.ID) error {
	moneyTxShippingEtop, err := a.moneyTxShippingEtopStore(ctx).ID(MoneyTxShippingEtopID).GetMoneyTxShippingEtop()
	if err != nil {
		return err
	}
	if moneyTxShippingEtop.Status == status3.P || !moneyTxShippingEtop.ConfirmedAt.IsZero() {
		return cm.Errorf(cm.FailedPrecondition, nil, "phiên etop đã hoàn thành").WithMetap("money_transaction_etop_id", moneyTxShippingEtop.ID)
	}
	moneyTxShippings, err := a.moneyTxShippingStore(ctx).MoneyTxShippingEtopID(MoneyTxShippingEtopID).ListMoneyTxShippings()
	if err != nil {
		return err
	}
	var moneyTxShippingIDs []dot.ID
	for _, item := range moneyTxShippings {
		moneyTxShippingIDs = append(moneyTxShippingIDs, item.ID)
	}
	statistics, err := a.prepareMoneyTxShippingEtop(ctx, MoneyTxShippingEtopID, moneyTxShippingIDs)
	if err != nil {
		return err
	}
	updateStatistics := &moneytxsqlstore.UpdateMoneyTxShippingEtopStatisticsArgs{
		ID:                    MoneyTxShippingEtopID,
		TotalCOD:              dot.Int(statistics.TotalCOD),
		TotalAmount:           dot.Int(statistics.TotalAmount),
		TotalOrders:           dot.Int(statistics.TotalOrders),
		TotalFee:              dot.Int(statistics.TotalShippingFee),
		TotalMoneyTransaction: dot.Int(statistics.TotalOrders),
	}
	if err := a.moneyTxShippingEtopStore(ctx).UpdateMoneyTxShippingEtopStatistics(updateStatistics); err != nil {
		return err
	}
	return nil
}
