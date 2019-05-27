package sqlstore

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	txmodel "etop.vn/backend/pkg/services/moneytx/model"
	"etop.vn/backend/pkg/services/moneytx/modelx"
	txmodely "etop.vn/backend/pkg/services/moneytx/modely"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
	shipmodelx "etop.vn/backend/pkg/services/shipping/modelx"
	"etop.vn/backend/pkg/services/shipping/modely"
	shipmodely "etop.vn/backend/pkg/services/shipping/modely"
)

func init() {
	bus.AddHandlers("sql",
		CreateMoneyTransactions,
		CreateMoneyTransaction,
		GetMoneyTransaction,
		GetMoneyTransactions,
		RemoveFfmsMoneyTransaction,
		ConfirmMoneyTransaction,
		DeleteMoneyTransaction,
		UpdateMoneyTransaction,

		CreateMoneyTransactionShippingExternal,
		CreateMoneyTransactionShippingExternalLine,
		RemoveMoneyTransactionShippingExternalLines,
		ComfirmMoneyTransactionShippingExternals,
		GetMoneyTransactionShippingExternal,
		GetMoneyTransactionShippingExternals,
		DeleteMoneyTransactionShippingExternal,
		UpdateMoneyTransactionShippingExternal,
		CreateCredit,
		GetCredit,
		GetCredits,
		UpdateCredit,
		ConfirmCredit,
		DeleteCredit,
		CalcBalanceShop,

		CreateMoneyTransactionShippingEtop,
		GetMoneyTransactionShippingEtop,
		GetMoneyTransactionShippingEtops,
		UpdateMoneyTransactionShippingEtop,
		DeleteMoneyTransactionShippingEtop,
		ConfirmMoneyTransactionShippingEtop,
	)
}

var zeroTime = time.Unix(0, 0)

var acceptStates = []string{
	string(model.StateReturned), string(model.StateReturning), string(model.StateDelivered), string(model.StateUndeliverable),
}

var filterMoneyTransactionShippingWhitelist = FilterWhitelist{
	Arrays:   nil,
	Contains: []string{"shop.name"},
	Dates:    []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals: []string{
		"code", "shop.money_transaction_rrule", "shop.name",
		"shop.phone",
	},
	Numbers: []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:  []string{"status"},
	PrefixOrRename: map[string]string{
		"shop.name":                    "s.name",
		"shop.phone":                   "s.phone",
		"shop.money_transaction_rrule": "s.money_transaction_rrule",
		"code":                         "m",
		"created_at":                   "m",
		"updated_at":                   "m",
		"confirmed_at":                 "m",
		"etop_transfered_at":           "m",
		"status":                       "m",
	},
}

var filterMoneyTransactionWhitelist = FilterWhitelist{
	Arrays:         nil,
	Contains:       []string{},
	Dates:          []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals:         []string{"code"},
	Numbers:        []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:         []string{"status"},
	PrefixOrRename: map[string]string{},
}

func CreateMoneyTransactions(ctx context.Context, cmd *modelx.CreateMoneyTransactions) error {
	return createMoneyTransactions(ctx, x, cmd)
}

func createMoneyTransactions(ctx context.Context, x Qx, cmd *modelx.CreateMoneyTransactions) error {
	if len(cmd.ShopIDMapFfms) == 0 {
		return cm.Error(cm.InvalidArgument, "ShopIDMapFfms can not be empty", nil)
	}

	created := 0
	for shopID, fulfillments := range cmd.ShopIDMapFfms {
		totalCOD, totalAmount, totalOrders, _, fulfillmentIDs := CalcFulfillmentsInfo(fulfillments)
		shop := cmd.ShopIDMap[shopID]
		command := &modelx.CreateMoneyTransaction{
			Shop:           shop,
			FulFillmentIDs: fulfillmentIDs,
			TotalCOD:       totalCOD,
			TotalAmount:    totalAmount,
			TotalOrders:    totalOrders,
		}

		if err := createMoneyTransaction(ctx, x, command); err != nil {
			return err
		}
		created++
	}
	cmd.Result.Created = created
	return nil
}

func CalcFulfillmentsInfo(fulfillments []*shipmodel.Fulfillment) (totalCOD int, totalAmount int, totalOrders int, totalShippingFee int, ffmIDs []int64) {
	ffmIDs = make([]int64, len(fulfillments))
	totalCOD = 0
	totalAmount = 0
	totalOrders = 0
	totalShippingFee = 0
	for i, ffm := range fulfillments {
		ffmIDs[i] = ffm.ID
		amount := ffm.TotalCODAmount
		if ffm.ShippingState == model.StateReturned || ffm.ShippingState == model.StateReturning {
			// make sure COD = 0
			amount = 0
		} else if ffm.ShippingState == model.StateUndeliverable {
			// trường hợp đơn bồi hoàn
			amount = ffm.ActualCompensationAmount
		}
		totalAmount = totalAmount + amount - ffm.ShippingFeeShop
		totalCOD += amount
		totalOrders++
		totalShippingFee += ffm.ShippingFeeShop
	}
	return totalCOD, totalAmount, totalOrders, totalShippingFee, ffmIDs
}

func CreateMoneyTransaction(ctx context.Context, cmd *modelx.CreateMoneyTransaction) error {
	return createMoneyTransaction(ctx, x, cmd)
}

func createMoneyTransaction(ctx context.Context, x Qx, cmd *modelx.CreateMoneyTransaction) error {
	// only handler for shop
	if cmd.Shop == nil {
		return cm.Error(cm.InvalidArgument, "Missing Shop ", nil)
	}
	if len(cmd.FulFillmentIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "FulfillmentIDs can not be empty", nil)
	}

	transaction := &txmodel.MoneyTransactionShipping{
		ShopID: cmd.Shop.ID,
		Status: model.S3Zero,
	}

	transaction.ID = cm.NewID()
	// generate order code
	code, errCode := GenerateCode(ctx, x, model.CodeTypeMoneyTransaction, cmd.Shop.Code)
	if errCode != nil {
		return errCode
	}
	transaction.Code = code
	if err := x.Table("money_transaction_shipping").ShouldInsert(transaction); err != nil {
		return err
	}

	var fulfillments []*modely.FulfillmentExtended
	var ffms []*shipmodel.Fulfillment
	if err := x.Table("fulfillment").Where("f.shop_id = ? AND f.type_from = ?",
		cmd.Shop.ID, model.FFShop).
		In("f.id", cmd.FulFillmentIDs).Find((*shipmodely.FulfillmentExtendeds)(&fulfillments)); err != nil {
		return err
	}
	for _, id := range cmd.FulFillmentIDs {
		found := false
		for _, ffm := range fulfillments {
			if id == ffm.ID {
				if err := CheckFulfillmentValid(ffm.Fulfillment); err != nil {
					return err
				}
				ffms = append(ffms, ffm.Fulfillment)
				found = true
				break
			}
		}
		if !found {
			return cm.Errorf(cm.NotFound, nil, "Fulfillment id not found or it has not done or it is belongs to another transaction", id)
		}
	}
	totalCOD, totalAmount, totalOrders, _, _ := CalcFulfillmentsInfo(ffms)
	if totalCOD != cmd.TotalCOD {
		return cm.Error(cm.FailedPrecondition, "Total COD does not match", nil)
	}
	if totalCOD != cmd.TotalCOD {
		return cm.Error(cm.FailedPrecondition, "Total Amount does not match", nil)
	}
	if totalOrders != cmd.TotalOrders {
		return cm.Error(cm.FailedPrecondition, "Total Order does not match", nil)
	}
	if err := x.Table("fulfillment").In("id", cmd.FulFillmentIDs).
		ShouldUpdateMap(M{"money_transaction_id": transaction.ID}); err != nil {
		return err
	}

	transaction.TotalCOD = totalCOD
	transaction.TotalOrders = totalOrders
	transaction.TotalAmount = totalAmount
	if err := x.Table("money_transaction_shipping").Where("id = ?", transaction.ID).ShouldUpdate(transaction); err != nil {
		return err
	}

	cmd.Result = &txmodely.MoneyTransactionExtended{
		MoneyTransactionShipping: transaction,
		Fulfillments:             fulfillments,
	}
	return nil
}

func GetMoneyTransaction(ctx context.Context, query *modelx.GetMoneyTransaction) error {
	return getMoneyTransaction(ctx, x, query)
}

func getMoneyTransaction(ctx context.Context, x Qx, query *modelx.GetMoneyTransaction) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	s := x.Table("money_transaction_shipping").Where("id = ?", query.ID)
	if query.ShopID != 0 {
		s = s.Where("shop_id = ?", query.ShopID)
	}

	var transaction = new(txmodel.MoneyTransactionShipping)
	if err := s.ShouldGet(transaction); err != nil {
		return err
	}
	var fulfillments []*modely.FulfillmentExtended
	if err := x.Table("fulfillment").Where("f.money_transaction_id = ?", query.ID).Find((*shipmodely.FulfillmentExtendeds)(&fulfillments)); err != nil {
		return err
	}
	query.Result = &txmodely.MoneyTransactionExtended{
		MoneyTransactionShipping: transaction,
		Fulfillments:             fulfillments,
	}
	return nil
}

func GetMoneyTransactions(ctx context.Context, query *modelx.GetMoneyTransactions) error {
	s := x.Table("money_transaction_shipping")
	if query.ShopID != 0 {
		s = s.Where("m.shop_id = ?", query.ShopID)
	}
	if query.MoneyTransactionShippingExternalID != 0 {
		s = s.Where("m.money_transaction_shipping_external_id = ?", query.MoneyTransactionShippingExternalID)
	}
	s, _, err := Filters(s, query.Filters, filterMoneyTransactionShippingWhitelist)
	if err != nil {
		return err
	}
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-m.created_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"m.created_at": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("m.id", query.IDs)
		}
		var moneyTransactions []*txmodel.MoneyTransactionShippingFtShop
		if err := s2.Find((*txmodel.MoneyTransactionShippingFtShops)(&moneyTransactions)); err != nil {
			return err
		}

		result := make([]*txmodely.MoneyTransactionExtended, len(moneyTransactions))
		for i, transaction := range moneyTransactions {
			result[i] = &txmodely.MoneyTransactionExtended{
				MoneyTransactionShipping: transaction.MoneyTransactionShipping,
			}
		}
		if query.IncludeFulfillments {
			moneyTransactionIDs := make([]int64, len(moneyTransactions))
			for i, transaction := range moneyTransactions {
				moneyTransactionIDs[i] = transaction.ID
			}

			var fulfillments []*modely.FulfillmentExtended
			if err := x.Table("fulfillment").In("f.money_transaction_id", moneyTransactionIDs).
				Find((*shipmodely.FulfillmentExtendeds)(&fulfillments)); err != nil {
				return err
			}
			ffmsByMoneyTransactionID := make(map[int64][]*modely.FulfillmentExtended)
			for _, ffm := range fulfillments {
				ffmsByMoneyTransactionID[ffm.MoneyTransactionID] = append(ffmsByMoneyTransactionID[ffm.MoneyTransactionID], ffm)
			}
			for _, transaction := range result {
				transaction.Fulfillments = ffmsByMoneyTransactionID[transaction.ID]
			}
		}
		query.Result.MoneyTransactions = result
	}
	{
		total, err := s.Count(&txmodel.MoneyTransactionShippingFtShop{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func UpdateMoneyTransaction(ctx context.Context, cmd *modelx.UpdateMoneyTransaction) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	s := x.Table("money_transaction_shipping").Where("id = ?", cmd.ID)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	{
		s1 := s.Clone()
		var transaction = new(txmodel.MoneyTransactionShipping)
		if err := s1.ShouldGet(transaction); err != nil {
			return err
		}
		if transaction.Status == model.S3Positive {
			return cm.Error(cm.InvalidArgument, "This money transaction was confirm. Can not update!", nil)
		}
	}
	{
		s2 := s.Clone()
		m := &txmodel.MoneyTransactionShipping{
			ID:            cmd.ID,
			Note:          cmd.Note,
			InvoiceNumber: cmd.InvoiceNumber,
			BankAccount:   cmd.BankAccount,
		}
		if err := s2.ShouldUpdate(m); err != nil {
			return err
		}
	}

	query := &modelx.GetMoneyTransaction{
		ID: cmd.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func RemoveFfmsMoneyTransaction(ctx context.Context, cmd *modelx.RemoveFfmsMoneyTransaction) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Mising AccountID", nil)
	}
	if cmd.MoneyTransactionID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money Transaction ID", nil)
	}
	if len(cmd.FulfillmentIDs) == 0 {
		return cm.Error(cm.FailedPrecondition, "FulfillmentIDs can not be empty", nil)
	}

	var transaction = new(txmodel.MoneyTransactionShipping)
	if err := x.Table("money_transaction_shipping").Where("id = ? AND shop_id = ?", cmd.MoneyTransactionID, cmd.ShopID).
		ShouldGet(transaction); err != nil {
		return err
	}

	var fulfillments []*shipmodel.Fulfillment
	if err := x.Table("fulfillment").
		Where("shop_id = ? AND money_transaction_id = ?", cmd.ShopID, cmd.MoneyTransactionID).
		In("id", cmd.FulfillmentIDs).Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return err
	}

	for _, id := range cmd.FulfillmentIDs {
		found := false
		for _, ffm := range fulfillments {
			if id == ffm.ID {
				found = true
				break
			}
		}
		if !found {
			return cm.Error(cm.NotFound, "Fulfillment #"+strconv.Itoa(int(id))+" does not exist in this money transaction", nil)
		}
	}
	return inTransaction(func(s Qx) error {
		s2 := s.Table("fulfillment").Where("shop_id = ?", cmd.ShopID).In("id", cmd.FulfillmentIDs)
		if _, err := s2.UpdateMap(M{
			"money_transaction_id":                   nil,
			"money_transaction_shipping_external_id": nil,
		}); err != nil {
			return err
		}
		var fulfillments []*shipmodel.Fulfillment
		if err := s.Table("fulfillment").Where("money_transaction_id = ?", cmd.MoneyTransactionID).Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
			return err
		}

		// update money_transaction
		totalCOD, totalAmount, totalOrders, _, _ := CalcFulfillmentsInfo(fulfillments)
		m := map[string]interface{}{
			"total_cod":    totalCOD,
			"total_orders": totalOrders,
			"total_amount": totalAmount,
		}
		if totalOrders == 0 && totalCOD == 0 {
			// money transaction does not have any ffm => update status = -1
			m["status"] = model.S3Negative
		}

		if err := s.Table("money_transaction_shipping").Where("id = ? AND shop_id = ?", cmd.MoneyTransactionID, cmd.ShopID).ShouldUpdateMap(m); err != nil {
			return err
		}

		// reGet money_transaction
		query := &modelx.GetMoneyTransaction{
			ID:     cmd.MoneyTransactionID,
			ShopID: cmd.ShopID,
		}
		if err := getMoneyTransaction(ctx, s, query); err != nil {
			return err
		}

		cmd.Result = query.Result
		return nil
	})

	return nil
}

func ConfirmMoneyTransaction(ctx context.Context, cmd *modelx.ConfirmMoneyTransaction) error {
	if cmd.MoneyTransactionID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money transaction ID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Shop ID", nil)
	}

	query := &modelx.GetMoneyTransaction{
		ID:     cmd.MoneyTransactionID,
		ShopID: cmd.ShopID,
	}
	if err := GetMoneyTransaction(ctx, query); err != nil {
		return err
	}
	fulfillments := query.Result.Fulfillments
	transaction := query.Result.MoneyTransactionShipping
	if transaction.Status != model.S3Zero {
		return cm.Error(cm.FailedPrecondition, "Can not confirm this money transaction", nil)
	}
	var ffms = make([]*shipmodel.Fulfillment, len(fulfillments))
	for i, ffm := range fulfillments {
		if !cm.StringsContain(acceptStates, string(ffm.ShippingState)) {
			return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" does not valid. Status must be delivered or returning or returned.", nil)
		}
		ffms[i] = ffm.Fulfillment
	}
	totalCOD, totalAmount, totalOrders, _, ffmIDs := CalcFulfillmentsInfo(ffms)

	if totalCOD != cmd.TotalCOD {
		return cm.Error(cm.FailedPrecondition, "Total COD does not match", nil)
	}
	if totalAmount != cmd.TotalAmount {
		return cm.Error(cm.FailedPrecondition, "Total Amount does not match", nil)
	}
	if totalOrders != cmd.TotalOrders {
		return cm.Error(cm.FailedPrecondition, "Total Order does not match", nil)
	}

	return inTransaction(func(s Qx) error {
		now := time.Now()
		if err := s.Table("money_transaction_shipping").Where("id = ?", cmd.MoneyTransactionID).
			ShouldUpdateMap(M{
				"total_orders":       totalOrders,
				"total_cod":          totalCOD,
				"total_amount":       totalAmount,
				"status":             model.S3Positive,
				"confirmed_at":       now,
				"etop_transfered_at": now,
			}); err != nil {
			return err
		}

		if err := s.Table("fulfillment").In("id", ffmIDs).
			ShouldUpdateMap(M{
				"cod_etop_transfered_at": now,
			}); err != nil {
			return err
		}
		cmd.Result.Updated = 1

		return nil
	})
}

func DeleteMoneyTransaction(ctx context.Context, cmd *modelx.DeleteMoneyTransaction) error {
	if cmd.MoneyTransactionID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money Transaction ID", nil)
	}
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Shop ID", nil)
	}
	var transaction = new(txmodel.MoneyTransactionShipping)
	if err := x.Table("money_transaction_shipping").Where("id = ? AND shop_id = ?", cmd.MoneyTransactionID, cmd.ShopID).
		ShouldGet(transaction); err != nil {
		return err
	}
	if transaction.Status == model.S3Positive {
		return cm.Error(cm.FailedPrecondition, "Can not delete this money transaction", nil)
	}
	return inTransaction(func(s Qx) error {
		if _, err := s.Table("fulfillment").Where("money_transaction_id = ?", cmd.MoneyTransactionID).
			UpdateMap(M{
				"money_transaction_id":                   nil,
				"money_transaction_shipping_external_id": nil,
			}); err != nil {
			return err
		}

		if deleted, err := s.Table("money_transaction_shipping").Where("id = ? AND shop_id = ?", cmd.MoneyTransactionID, cmd.ShopID).
			Delete(&txmodel.MoneyTransactionShipping{}); err != nil {
			return err
		} else if deleted == 0 {
			return cm.Error(cm.NotFound, "", nil)
		}
		cmd.Result.Deleted = 1
		return nil
	})
}

func CreateMoneyTransactionShippingExternal(ctx context.Context, cmd *modelx.CreateMoneyTransactionShippingExternal) error {
	if cmd.Provider == "" {
		return cm.Error(cm.InvalidArgument, "Vui lòng chọn nhà vận chuyển", nil)
	}
	if len(cmd.Lines) == 0 {
		return cm.Error(cm.InvalidArgument, "Vận đơn không được rỗng", nil)
	}
	totalCOD := 0
	totalOrders := 0
	for _, line := range cmd.Lines {
		totalCOD += line.ExternalTotalCOD
		totalOrders++
	}
	return inTransaction(func(s Qx) error {
		code, errCode := GenerateCode(ctx, s, model.CodeTypeMoneyTransactionExternal, model.TypeGHNCode)
		if errCode != nil {
			return errCode
		}
		externalTransaction := &txmodel.MoneyTransactionShippingExternal{
			ID:             cm.NewID(),
			Code:           code,
			TotalCOD:       totalCOD,
			TotalOrders:    totalOrders,
			ExternalPaidAt: cmd.ExternalPaidAt,
			Provider:       cmd.Provider,
			Note:           cmd.Note,
			InvoiceNumber:  cmd.InvoiceNumber,
			BankAccount:    cmd.BankAccount,
		}
		if err := s.Table("money_transaction_shipping_external").ShouldInsert(externalTransaction); err != nil {
			return err
		}

		ffmIDs := make([]int64, 0, len(cmd.Lines))
		for _, line := range cmd.Lines {
			createCmd := &modelx.CreateMoneyTransactionShippingExternalLine{
				ExternalCode:                       line.ExternalCode,
				ExternalTotalCOD:                   line.ExternalTotalCOD,
				ExternalCreatedAt:                  line.ExternalCreatedAt,
				ExternalClosedAt:                   line.ExternalClosedAt,
				EtopFulfillmentIdRaw:               line.EtopFulfillmentIdRaw,
				ExternalCustomer:                   line.ExternalCustomer,
				ExternalAddress:                    line.ExternalAddress,
				MoneyTransactionShippingExternalID: externalTransaction.ID,
				ExternalTotalShippingFee:           line.ExternalTotalShippingFee,
			}
			if err := createMoneyTransactionShippingExternalLine(ctx, s, createCmd); err != nil {
				return err
			}
			if createCmd.Result.EtopFulfillmentID != 0 && createCmd.Result.ImportError == nil {
				ffmIDs = append(ffmIDs, createCmd.Result.EtopFulfillmentID)
			}
		}
		if len(ffmIDs) > 0 {
			if err := s.Table("fulfillment").In("id", ffmIDs).ShouldUpdateMap(M{
				"money_transaction_shipping_external_id": externalTransaction.ID,
			}); err != nil {
				return err
			}
		}

		query := &modelx.GetMoneyTransactionShippingExternal{
			ID: externalTransaction.ID,
		}
		if err := getMoneyTransactionShippingExternal(ctx, s, query); err != nil {
			return err
		}
		cmd.Result = query.Result
		return nil
	})
}

func CreateMoneyTransactionShippingExternalLine(ctx context.Context, cmd *modelx.CreateMoneyTransactionShippingExternalLine) error {
	return createMoneyTransactionShippingExternalLine(ctx, x, cmd)
}

func createMoneyTransactionShippingExternalLine(ctx context.Context, x Qx, cmd *modelx.CreateMoneyTransactionShippingExternalLine) error {
	if cmd.MoneyTransactionShippingExternalID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing MoneyTransactionShippingExternalID", nil)
	}
	line := &txmodel.MoneyTransactionShippingExternalLine{
		ID:                                 cm.NewID(),
		ExternalCode:                       cmd.ExternalCode,
		ExternalTotalCOD:                   cmd.ExternalTotalCOD,
		ExternalCreatedAt:                  cmd.ExternalCreatedAt,
		ExternalClosedAt:                   cmd.ExternalClosedAt,
		ExternalCustomer:                   cmd.ExternalCustomer,
		ExternalAddress:                    cmd.ExternalAddress,
		EtopFulfillmentIdRaw:               cmd.EtopFulfillmentIdRaw,
		MoneyTransactionShippingExternalID: cmd.MoneyTransactionShippingExternalID,
		ExternalTotalShippingFee:           cmd.ExternalTotalShippingFee,
	}
	if line.ExternalCode == "" {
		line.ImportError = &model.Error{
			Code: "ffm_id_empty",
			Msg:  "Thiếu mã vận đơn",
		}
	} else {
		var ffm = new(shipmodel.Fulfillment)
		if has, err := x.Table("fulfillment").Where("shipping_code = ?", line.ExternalCode).Get(ffm); err != nil || !has {
			line.ImportError = &model.Error{
				Code: "ffm_not_found",
				Msg:  "Không tìm thấy vận đơn trên Etop",
			}
		} else {
			line.EtopFulfillmentID = ffm.ID
			if ffm.MoneyTransactionShippingExternalID != 0 {
				line.ImportError = &model.Error{
					Code: "ffm_exist_money_transaction_shipping_external",
					Msg:  "Vận đơn nằm trong phiên thanh toán nhà vận chuyển khác: " + strconv.Itoa(int(ffm.MoneyTransactionShippingExternalID)),
				}
			} else if !cm.StringsContain(acceptStates, string(ffm.ShippingState)) {
				line.ImportError = &model.Error{
					Code: "ffm_not_done",
					Msg:  "Vận đơn chưa hoàn thành trên Etop",
				}
			} else if ffm.ShippingState == model.StateDelivered && ffm.TotalCODAmount != line.ExternalTotalCOD {
				line.ImportError = &model.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị vận đơn không đúng",
				}
			} else if ffm.ShippingState == model.StateUndeliverable && line.ExternalTotalCOD != ffm.ActualCompensationAmount {
				line.ImportError = &model.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị bồi hoàn không đúng",
				}
			} else if ffm.MoneyTransactionID != 0 {
				line.ImportError = &model.Error{
					Code: "ffm_exist_money_transaction",
					Msg:  "Vận đơn nằm trong phiên thanh toán khác: " + strconv.Itoa(int(ffm.MoneyTransactionID)),
				}
			} else if line.ExternalTotalShippingFee != 0 && line.ExternalTotalShippingFee != ffm.ShippingFeeShop {
				line.ImportError = &model.Error{
					Code: "ffm_shipping_fee_not_match",
					Msg:  "Tổng tiền cước không đúng.",
					Meta: map[string]string{
						"Etop":     strconv.Itoa(ffm.ShippingFeeShop),
						"Provider": strconv.Itoa(line.ExternalTotalShippingFee),
					},
				}
			}
		}

	}
	if err := x.Table("money_transaction_shipping_external_line").ShouldInsert(line); err != nil {
		return err
	}
	cmd.Result = line

	return nil
}

func RemoveMoneyTransactionShippingExternalLines(ctx context.Context, cmd *modelx.RemoveMoneyTransactionShippingExternalLines) error {
	if cmd.MoneyTransactionShippingExternalID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money Transaction Shipping External ID", nil)
	}
	if len(cmd.LineIDs) == 0 {
		return cm.Error(cm.FailedPrecondition, "LineIDs can not be empty", nil)
	}

	var transaction = new(txmodel.MoneyTransactionShippingExternal)
	if err := x.Table("money_transaction_shipping_external").Where("id = ?", cmd.MoneyTransactionShippingExternalID).
		ShouldGet(transaction); err != nil {
		return err
	}

	var lines []*txmodel.MoneyTransactionShippingExternalLine
	if err := x.Table("money_transaction_shipping_external_line").Where("money_transaction_shipping_external_id = ?", cmd.MoneyTransactionShippingExternalID).
		In("id", cmd.LineIDs).
		Find((*txmodel.MoneyTransactionShippingExternalLines)(&lines)); err != nil {
		return err
	}
	ffmIDs := make([]int64, 0, len(cmd.LineIDs))
	for _, id := range cmd.LineIDs {
		found := false
		for _, line := range lines {
			if id == line.ID {
				found = true
				if line.EtopFulfillmentID != 0 {
					ffmIDs = append(ffmIDs, line.EtopFulfillmentID)
				}
				break
			}
		}
		if !found {
			return cm.Error(cm.NotFound, "Line #"+strconv.Itoa(int(id))+" does not exist in this money transaction", nil)
		}
	}
	return inTransaction(func(s Qx) error {
		s2 := s.Table("money_transaction_shipping_external_line").
			Where("money_transaction_shipping_external_id = ?", cmd.MoneyTransactionShippingExternalID).
			In("id", cmd.LineIDs)
		if _, err := s2.Delete(&txmodel.MoneyTransactionShippingExternalLine{}); err != nil {
			return err
		}

		if len(ffmIDs) > 0 {
			if err := s.Table("fulfillment").In("id", ffmIDs).ShouldUpdateMap(M{
				"money_transaction_shipping_external_id": nil,
			}); err != nil {
				return err
			}
		}

		var lines []*txmodel.MoneyTransactionShippingExternalLine
		if err := s.Table("money_transaction_shipping_external_line").
			Where("money_transaction_shipping_external_id = ?", cmd.MoneyTransactionShippingExternalID).
			Find((*txmodel.MoneyTransactionShippingExternalLines)(&lines)); err != nil {
			return err
		}

		// update money_transaction
		totalOrders := len(lines)
		totalCOD := 0
		for _, line := range lines {
			totalCOD += line.ExternalTotalCOD
		}
		m := map[string]interface{}{
			"total_cod":    totalCOD,
			"total_orders": totalOrders,
		}
		if totalOrders == 0 && totalCOD == 0 {
			// money transaction does not have any ffm => update status = -1
			m["status"] = model.S3Negative
		}

		if err := s.Table("money_transaction_shipping_external").Where("id = ?", cmd.MoneyTransactionShippingExternalID).
			ShouldUpdateMap(m); err != nil {
			return err
		}

		// reGet money_transaction_shipping_external
		query := &modelx.GetMoneyTransactionShippingExternal{
			ID: cmd.MoneyTransactionShippingExternalID,
		}
		if err := getMoneyTransactionShippingExternal(ctx, s, query); err != nil {
			return err
		}

		cmd.Result = query.Result
		return nil
	})

	return nil
}

func DeleteMoneyTransactionShippingExternal(ctx context.Context, cmd *modelx.DeleteMoneyTransactionShippingExternal) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money Transaction ID", nil)
	}

	var transaction = new(txmodel.MoneyTransactionShippingExternal)
	if err := x.Table("money_transaction_shipping_external").Where("id = ?", cmd.ID).
		ShouldGet(transaction); err != nil {
		return err
	}
	if transaction.Status == model.S3Positive {
		return cm.Error(cm.FailedPrecondition, "Can not delete this money transaction", nil)
	}
	return inTransaction(func(s Qx) error {
		if _, err := s.Table("money_transaction_shipping_external_line").Where("money_transaction_shipping_external_id = ?", cmd.ID).
			Delete(&txmodel.MoneyTransactionShippingExternalLine{}); err != nil {
			return err
		}

		if _, err := s.Table("fulfillment").Where("money_transaction_shipping_external_id = ?", cmd.ID).UpdateMap(M{
			"money_transaction_shipping_external_id": nil,
		}); err != nil {
			return err
		}
		if deleted, err := s.Table("money_transaction_shipping_external").Where("id = ?", cmd.ID).
			Delete(&txmodel.MoneyTransactionShippingExternal{}); err != nil {
			return err
		} else if deleted == 0 {
			return cm.Error(cm.NotFound, "", nil)
		}
		cmd.Result.Deleted = 1
		return nil
	})
}

func UpdateMoneyTransactionShippingExternal(ctx context.Context, cmd *modelx.UpdateMoneyTransactionShippingExternal) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	s := x.Table("money_transaction_shipping_external").Where("id = ?", cmd.ID)
	var transaction = new(txmodel.MoneyTransactionShippingExternal)
	{
		s1 := s.Clone()
		if err := s1.ShouldGet(transaction); err != nil {
			return err
		}
		if transaction.Status == model.S3Positive {
			return cm.Error(cm.FailedPrecondition, "Can not update this money transaction", nil)
		}
	}
	{
		s2 := s.Clone()
		m := &txmodel.MoneyTransactionShippingExternal{
			ID:            cmd.ID,
			Note:          cmd.Note,
			InvoiceNumber: cmd.InvoiceNumber,
			BankAccount:   cmd.BankAccount,
		}

		if err := s2.ShouldUpdate(m); err != nil {
			return err
		}
	}
	query := &modelx.GetMoneyTransactionShippingExternal{
		ID: cmd.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

/*
# Khi tạo phiên thanh toán của etop cho Shop
	+ Lọc đơn của phiên thanh toán ghn cho etop
	+ Thêm các ffms có trạng thái "returned"
    Điều kiện:
    	* Đảm bảo 1 phiên chuyển tiền của etop cho shop >= 0
    	* Cho các ffms returned vào lần lượt theo thứ tự từ cũ tới mới cho tới khi total >= 0
    	* Các ffms returned còn lại để nguyên, cho vào phiên sau
*/

func PreprocessMoneyTransactionExternal(ctx context.Context, externalMoneyTransactionExtended *txmodel.MoneyTransactionShippingExternalExtended) (shopFfmMap map[int64][]*shipmodel.Fulfillment, _err error) {
	shopFfmMap = make(map[int64][]*shipmodel.Fulfillment)
	if externalMoneyTransactionExtended.Status != model.S3Zero {
		_err = cm.Error(cm.FailedPrecondition, "Can not confirm this money transaction", nil).WithMetap("id", externalMoneyTransactionExtended.ID)
		return shopFfmMap, _err
	}
	lines := externalMoneyTransactionExtended.Lines
	if len(lines) == 0 {
		_err = cm.Error(cm.FailedPrecondition, "There are no lines in this money transaction", nil).WithMetap("id", externalMoneyTransactionExtended.ID)
		return shopFfmMap, _err
	}

	ffmCodes := make([]string, len(lines))
	for i, line := range lines {
		if line.ImportError != nil && line.ImportError.Code != "" {
			_err = cm.Error(cm.FailedPrecondition, "Vui lòng xử lý lỗi trước khi xác nhận phiên", nil)
			return shopFfmMap, _err
		}
		ffmCodes[i] = line.ExternalCode
	}
	var fulfillments []*shipmodel.Fulfillment
	if err := x.Table("fulfillment").
		In("shipping_code", ffmCodes).
		Find((*shipmodel.Fulfillments)(&fulfillments)); err != nil {
		return shopFfmMap, err
	}

	for _, line := range lines {
		found := false
		for _, ffm := range fulfillments {
			if line.ExternalCode == ffm.ShippingCode {
				found = true
				break
			}
		}
		if !found {
			_err = cm.Errorf(cm.NotFound, nil, "Không tìm thấy vận đơn %v", line.ExternalCode)
			return shopFfmMap, _err
		}
	}
	for _, ffm := range fulfillments {
		if shopFfmMap[ffm.ShopID] == nil {
			shopFfmMap[ffm.ShopID] = make([]*shipmodel.Fulfillment, 0, len(fulfillments))
		}
		found := false
		for _, _ffm := range shopFfmMap[ffm.ShopID] {
			if ffm.ID == _ffm.ID {
				found = true
				break
			}
		}
		if !found {
			shopFfmMap[ffm.ShopID] = append(shopFfmMap[ffm.ShopID], ffm)
		}
	}
	return shopFfmMap, nil
}

/*
* Confirm multiplelity money transaction Shipping Externals
* Collect lines to create money transaction for shop.
 */
func ComfirmMoneyTransactionShippingExternals(ctx context.Context, cmd *modelx.ConfirmMoneyTransactionShippingExternals) error {
	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Money transaction shipping external ID", nil)
	}
	query := &modelx.GetMoneyTransactionShippingExternals{
		IDs: cmd.IDs,
	}
	if err := GetMoneyTransactionShippingExternals(ctx, query); err != nil {
		return err
	}
	externalTransactions := query.Result.MoneyTransactionShippingExternals

	var externalTransactionIDs []int64
	var shopIDs []int64
	shopFfmMap := make(map[int64][]*shipmodel.Fulfillment)
	for _, externalTransaction := range externalTransactions {
		_shopFfmMap, err := PreprocessMoneyTransactionExternal(ctx, externalTransaction)
		if err != nil {
			return err
		}
		externalTransactionIDs = append(externalTransactionIDs, externalTransaction.ID)
		for shopId, ffms := range _shopFfmMap {
			shopFfmMap[shopId] = mergeFulfillments(shopFfmMap[shopId], ffms)
			if !cm.ContainInt64(shopIDs, shopId) {
				shopIDs = append(shopIDs, shopId)
			}
		}
	}
	{
		_shopFfmMap := combineWithExtraFfms()
		// make sure do not dupplicate ffm
		for shopId, ffms := range _shopFfmMap {
			shopFfmMap[shopId] = mergeFulfillments(shopFfmMap[shopId], ffms)
			if !cm.ContainInt64(shopIDs, shopId) {
				shopIDs = append(shopIDs, shopId)
			}
		}
	}
	shopsQuery := &model.GetShopsQuery{
		ShopIDs: shopIDs,
	}
	if err := bus.Dispatch(ctx, shopsQuery); err != nil {
		return err
	}

	shopsMap := make(map[int64]*model.Shop)
	for _, shop := range shopsQuery.Result.Shops {
		shopsMap[shop.ID] = shop
	}

	if len(shopIDs) != len(shopsMap) {
		return cm.Error(cm.Internal, "", nil).WithMeta("reason", "ShopIDs does not have the expected length")
	}

	cmdCreateMoneyTransaction := &modelx.CreateMoneyTransactions{
		ShopIDMapFfms: shopFfmMap,
		ShopIDMap:     shopsMap,
	}

	return inTransaction(func(s Qx) error {
		if err := createMoneyTransactions(ctx, s, cmdCreateMoneyTransaction); err != nil {
			return err
		}
		if err := s.Table("money_transaction_shipping_external").In("id", externalTransactionIDs).
			ShouldUpdateMap(M{
				"status": model.S3Positive,
			}); err != nil {
			return err
		}
		cmd.Result.Updated = 1
		return nil
	})
}

func mergeFulfillments(ffms []*shipmodel.Fulfillment, subFfms []*shipmodel.Fulfillment) []*shipmodel.Fulfillment {
	mergeFfms := append(ffms, subFfms...)
	ffmsMap := make(map[int64]*shipmodel.Fulfillment)
	for _, _ffm := range mergeFfms {
		ffmsMap[_ffm.ID] = _ffm
	}
	var res []*shipmodel.Fulfillment
	for _, _ffm := range ffmsMap {
		res = append(res, _ffm)
	}
	return res
}

func getExtraFfms(provider model.ShippingProvider, isNoneCOD bool, isReturned bool) ([]*shipmodel.Fulfillment, error) {
	// find all ffms has state: "returned" and cod_etop_transfered_at is NULL of this provider
	// find all ffms has state "delivered" and total_cod_amount = 0
	var ffms []*shipmodel.Fulfillment
	s := x.Table("fulfillment").
		Where("shipping_provider = ? AND cod_etop_transfered_at is NULL AND money_transaction_id is NULL AND money_transaction_shipping_external_id is NULL", string(provider))

	if isNoneCOD && isReturned {
		s = s.Where("shipping_state = ? OR (shipping_state = ? AND total_cod_amount = 0)", model.StateReturned, model.StateDelivered)
	} else if isReturned {
		s = s.Where("shipping_state = ?", model.StateReturned)
	} else if isNoneCOD {
		s = s.Where("shipping_state = ? AND total_cod_amount = 0", model.StateDelivered)
	}

	if err := s.Find((*shipmodel.Fulfillments)(&ffms)); err != nil {
		return []*shipmodel.Fulfillment{}, err
	}
	return ffms, nil
}

func combineWithExtraFfms() map[int64][]*shipmodel.Fulfillment {
	var ffmAdditionals []*shipmodel.Fulfillment
	shopFfmMap := make(map[int64][]*shipmodel.Fulfillment)
	// merge with GHN's ffms returned or (ffm delivered and total_cod_amount = 0)
	GHNFfms, _ := getExtraFfms(model.TypeGHN, true, true)
	ffmAdditionals = append(ffmAdditionals, GHNFfms...)

	// merge with VTPOST's ffms
	VtpostFfms := GetVtpostExtraFfms()
	ffmAdditionals = append(ffmAdditionals, VtpostFfms...)

	for _, ffm := range ffmAdditionals {
		if ffm.ID == 0 {
			continue
		}
		shopFfmMap[ffm.ShopID] = append(shopFfmMap[ffm.ShopID], ffm)
	}
	return shopFfmMap
}

func GetVtpostExtraFfms() []*shipmodel.Fulfillment {
	// find all ffms has state: "returned" & "returning" and cod_etop_transfered_at is NULL of vtpost
	// find all ffms has state "delivered" and total_cod_amount = 0
	var ffms []*shipmodel.Fulfillment
	{
		s := x.Table("fulfillment").
			Where("shipping_provider = ? AND cod_etop_transfered_at is NULL AND money_transaction_id is NULL AND money_transaction_shipping_external_id is NULL", string(model.TypeVTPost)).
			Where("shipping_state in (?, ?)", model.StateReturned, model.StateReturning)
		if err := s.Find((*shipmodel.Fulfillments)(&ffms)); err == nil {
			UpdateVtpostShippingFeeReturned(ffms)
		}
	}
	{
		// merge with VTPOST's ffms delivered and total_cod_amount = 0
		VTPOSTFfms, _ := getExtraFfms(model.TypeVTPost, true, false)
		ffms = append(ffms, VTPOSTFfms...)
	}
	return ffms
}

func GetMoneyTransactionShippingExternal(ctx context.Context, query *modelx.GetMoneyTransactionShippingExternal) error {
	return getMoneyTransactionShippingExternal(ctx, x, query)
}

func getMoneyTransactionShippingExternal(ctx context.Context, x Qx, query *modelx.GetMoneyTransactionShippingExternal) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	s := x.Table("money_transaction_shipping_external").Where("id = ?", query.ID)

	var transaction = new(txmodel.MoneyTransactionShippingExternal)
	if err := s.ShouldGet(transaction); err != nil {
		return err
	}
	var lines []*txmodel.MoneyTransactionShippingExternalLineExtended
	if err := x.Table("money_transaction_shipping_external_line").Where("m.money_transaction_shipping_external_id = ?", query.ID).
		Find((*txmodel.MoneyTransactionShippingExternalLineExtendeds)(&lines)); err != nil {
		return err
	}
	query.Result = &txmodel.MoneyTransactionShippingExternalExtended{
		MoneyTransactionShippingExternal: transaction,
		Lines:                            lines,
	}
	return nil
}

func GetMoneyTransactionShippingExternals(ctx context.Context, query *modelx.GetMoneyTransactionShippingExternals) error {
	s := x.Table("money_transaction_shipping_external")

	s, _, err := Filters(s, query.Filters, filterMoneyTransactionWhitelist)
	if err != nil {
		return err
	}
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"updated_at": "", "created_at": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("id", query.IDs)
		}

		var moneyTransactions []*txmodel.MoneyTransactionShippingExternal
		if err := s2.Find((*txmodel.MoneyTransactionShippingExternals)(&moneyTransactions)); err != nil {
			return err
		}

		moneyTransactionIDs := make([]int64, len(moneyTransactions))
		for i, transaction := range moneyTransactions {
			moneyTransactionIDs[i] = transaction.ID
		}

		var lines []*txmodel.MoneyTransactionShippingExternalLineExtended
		if err := x.Table("money_transaction_shipping_external_line").In("m.money_transaction_shipping_external_id", moneyTransactionIDs).
			Find((*txmodel.MoneyTransactionShippingExternalLineExtendeds)(&lines)); err != nil {
			return err
		}
		linesMoneyTransactionHash := make(map[int64][]*txmodel.MoneyTransactionShippingExternalLineExtended)
		for _, line := range lines {
			linesMoneyTransactionHash[line.MoneyTransactionShippingExternalID] = append(linesMoneyTransactionHash[line.MoneyTransactionShippingExternalID], line)
		}

		result := make([]*txmodel.MoneyTransactionShippingExternalExtended, len(moneyTransactions))
		for i, transaction := range moneyTransactions {
			result[i] = &txmodel.MoneyTransactionShippingExternalExtended{
				MoneyTransactionShippingExternal: transaction,
				Lines:                            linesMoneyTransactionHash[transaction.ID],
			}
		}
		query.Result.MoneyTransactionShippingExternals = result
	}
	{
		total, err := s.Count(&txmodel.MoneyTransactionShippingExternal{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func CheckFulfillmentValid(ffm *shipmodel.Fulfillment) error {
	if !cm.StringsContain(acceptStates, string(ffm.ShippingState)) {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" does not valid. Status must be delivered or returning or returned.", nil)
	}
	if ffm.MoneyTransactionID != 0 {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" in another money transaction.", nil)
	}
	if !ffm.CODEtopTransferedAt.IsZero() {
		return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" has paid.", nil)
	}
	return nil
}

func CreateCredit(ctx context.Context, cmd *model.CreateCreditCommand) error {
	if cmd.Type == "" {
		return cm.Error(cm.InvalidArgument, "Missing credit type", nil)
	}
	switch cmd.Type {
	case model.TypeShop:
		if cmd.ShopID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
		}

	default:
		return cm.Error(cm.InvalidArgument, "Type does not support", nil)
	}
	if cmd.Amount == 0 {
		return cm.Error(cm.InvalidArgument, "Missing amount", nil)
	}
	credit := &model.Credit{
		ID:     cm.NewID(),
		Amount: cmd.Amount,
		ShopID: cmd.ShopID,
		Type:   string(cmd.Type),
		PaidAt: cmd.PaidAt,
	}
	if err := x.Table("credit").ShouldInsert(credit); err != nil {
		return err
	}
	query := &model.GetCreditQuery{
		ID: credit.ID,
	}
	if err := GetCredit(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func GetCredit(ctx context.Context, query *model.GetCreditQuery) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}

	s := x.Table("credit").Where("c.id = ?", query.ID)
	if query.ShopID != 0 {
		s = s.Where("c.shop_id = ?", query.ShopID)
	}
	credit := new(model.CreditExtended)
	if err := s.ShouldGet(credit); err != nil {
		return err
	}
	query.Result = credit
	return nil
}

func GetCredits(ctx context.Context, query *model.GetCreditsQuery) error {
	s := x.Table("credit")
	if query.ShopID != 0 {
		s = s.Where("c.shop_id = ?", query.ShopID)
	}
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-c.updated_at"}
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"updated_at": "c.updated_at", "created_at": "c.created_at"})
		if err != nil {
			return err
		}
		var credits []*model.CreditExtended
		if err := s2.Find((*model.CreditExtendeds)(&credits)); err != nil {
			return err
		}
		query.Result.Credits = credits
	}
	{
		total, err := s.Count(&model.Credit{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func UpdateCredit(ctx context.Context, cmd *model.UpdateCreditCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	s := x.Table("credit").Where("id = ?", cmd.ID)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	credit := &model.Credit{
		PaidAt: cmd.PaidAt,
		Amount: cmd.Amount,
	}
	if err := s.ShouldUpdate(credit); err != nil {
		return err
	}
	query := &model.GetCreditQuery{
		ID: cmd.ID,
	}
	if err := GetCredit(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func ConfirmCredit(ctx context.Context, cmd *model.ConfirmCreditCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	s := x.Table("credit").Where("id = ?", cmd.ID)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	{
		s2 := s.Clone()
		credit := new(model.Credit)
		if err := s2.ShouldGet(credit); err != nil {
			return nil
		}
		if credit.Status == model.S3Positive {
			return cm.Error(cm.FailedPrecondition, "This credit has already confirmed", nil)
		}
		if credit.Status != model.S3Zero {
			return cm.Error(cm.FailedPrecondition, "Can not confirm this credit", nil)
		}
		if credit.PaidAt.IsZero() || credit.PaidAt.Equal(zeroTime) {
			return cm.Error(cm.FailedPrecondition, "Missing paid at", nil)
		}
	}
	if err := s.ShouldUpdateMap(M{
		"status": model.S3Positive,
	}); err != nil {
		return err
	}

	cmd.Result.Updated = 1
	return nil
}

func DeleteCredit(ctx context.Context, cmd *model.DeleteCreditCommand) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	s := x.Table("credit").Where("id = ?", cmd.ID)
	if cmd.ShopID != 0 {
		s = s.Where("shop_id = ?", cmd.ShopID)
	}
	{
		s2 := s.Clone()
		credit := new(model.Credit)
		if err := s2.ShouldGet(credit); err != nil {
			return nil
		}
		if credit.Status == model.S3Positive {
			return cm.Error(cm.FailedPrecondition, "This credit has already confirmed", nil)
		}
		if credit.Status != model.S3Zero {
			return cm.Error(cm.FailedPrecondition, "Can not delete this credit", nil)
		}
	}
	if deleted, err := s.Delete(&model.Credit{}); err != nil {
		return err
	} else if deleted == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	cmd.Result.Deleted = 1
	return nil
}

/*
	Số dư shop (không tính những ffm đã thanh toán)
	(1) COD Shop (ffm) khác trạng thái hủy và không phải là đơn trả hàng (status != -1 AND status != 0 AND shipping_status != -2 AND etop_payment_status != 1)
	(2) Phí giao hàng (ffm) khác trạng thái hủy (đã bao gồm đơn trả hàng)
	(3) Credit

	số dư = (1) + (3) - (2)
*/

func CalcBalanceShop(ctx context.Context, cmd *model.GetBalanceShopCommand) error {
	shopID := cmd.ShopID
	if shopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing shop ID", nil)
	}

	var totalCODAmount, totalShippingFee, totalCredit sql.NullInt64
	return inTransaction(func(s Qx) error {
		if err := s.SQL("SELECT SUM(total_cod_amount) from fulfillment").Where("shop_id = ? AND status != -1 AND status != 0 AND shipping_status != -2 AND etop_payment_status != 1", shopID).
			Scan(&totalCODAmount); err != nil {
			return err
		}
		if err := s.SQL("SELECT SUM(shipping_fee_shop) from fulfillment").Where("shop_id = ? AND status != -1 AND status != 0 AND etop_payment_status != 1", shopID).
			Scan(&totalShippingFee); err != nil {
			return err
		}
		if err := s.SQL("SELECT SUM(amount) from credit").Where("shop_id = ? AND status = ? AND paid_at is not NULL", shopID, model.S3Positive).
			Scan(&totalCredit); err != nil {
			return err
		}
		cmd.Result.Amount = int(totalCODAmount.Int64 - totalShippingFee.Int64 + totalCredit.Int64)
		return nil
	})
}

// MoneyTransactionShippingEtop
// Include group of MoneyTransactionShipping

func CreateMoneyTransactionShippingEtop(ctx context.Context, cmd *modelx.CreateMoneyTransactionShippingEtop) error {
	ids := cmd.MoneyTransactionShippingIDs
	if len(ids) == 0 {
		return cm.Error(cm.InvalidArgument, "MoneyTransactionIDs can not be empty", nil)
	}

	mtse, err := prepairMoneyTransactionShippingEtop(ctx, 0, ids)
	if err != nil {
		return err
	}
	newID := cm.NewID()
	_err := inTransaction(func(s Qx) error {
		// generate order code
		code, errCode := GenerateCode(ctx, s, model.CodeTypeMoneyTransactionEtop, "ETOP")
		if errCode != nil {
			return errCode
		}
		mtShippingEtop := &txmodel.MoneyTransactionShippingEtop{
			ID:                    newID,
			TotalCOD:              mtse.TotalCOD,
			TotalOrders:           mtse.TotalOrders,
			TotalAmount:           mtse.TotalAmount,
			TotalFee:              mtse.TotalFee,
			TotalMoneyTransaction: mtse.TotalMoneyTransaction,
			Code:                  code,
		}
		if err := s.Table("money_transaction_shipping_etop").ShouldInsert(mtShippingEtop); err != nil {
			return err
		}
		return nil
	})
	if _err != nil {
		return _err
	}
	if err := x.Table("money_transaction_shipping").In("id", ids).
		ShouldUpdateMap(M{"money_transaction_shipping_etop_id": newID}); err != nil {
		return err
	}

	_query := &modelx.GetMoneyTransactionShippingEtop{
		ID: newID,
	}
	if err := bus.Dispatch(ctx, _query); err != nil {
		return err
	}
	cmd.Result = _query.Result
	return nil
}

func GetMoneyTransactionShippingEtop(ctx context.Context, query *modelx.GetMoneyTransactionShippingEtop) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	s := x.Table("money_transaction_shipping_etop").Where("id = ?", query.ID)
	var transaction = new(txmodel.MoneyTransactionShippingEtop)
	if err := s.Where("id = ?", query.ID).ShouldGet(transaction); err != nil {
		return err
	}
	var moneyTransactionShippings []*txmodel.MoneyTransactionShipping
	if err := x.Table("money_transaction_shipping").Where("money_transaction_shipping_etop_id = ?", query.ID).Find((*txmodel.MoneyTransactionShippings)(&moneyTransactionShippings)); err != nil {
		return err
	}
	moneyTransactionShippingIDs := make([]int64, len(moneyTransactionShippings))
	for i, mt := range moneyTransactionShippings {
		moneyTransactionShippingIDs[i] = mt.ID
	}

	mtQuery := &modelx.GetMoneyTransactions{
		IDs:                 moneyTransactionShippingIDs,
		IncludeFulfillments: true,
	}
	if err := bus.Dispatch(ctx, mtQuery); err != nil {
		return err
	}

	query.Result = &txmodely.MoneyTransactionShippingEtopExtended{
		MoneyTransactionShippingEtop: transaction,
		MoneyTransactions:            mtQuery.Result.MoneyTransactions,
	}
	return nil
}

func GetMoneyTransactionShippingEtops(ctx context.Context, query *modelx.GetMoneyTransactionShippingEtops) error {
	s := x.Table("money_transaction_shipping_etop")
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-updated_at"}
	}
	if query.Status != nil {
		s = s.Where("status = ?", query.Status)
	}
	if len(query.IDs) > 0 {
		s = s.In("id", query.IDs)
	}
	s, _, err := Filters(s, query.Filters, filterMoneyTransactionWhitelist)
	if err != nil {
		return err
	}
	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"updated_at": "", "creted_at": ""})
		if err != nil {
			return err
		}
		var moneyTransactionShippingEtops []*txmodel.MoneyTransactionShippingEtop
		if err := s2.Find((*txmodel.MoneyTransactionShippingEtops)(&moneyTransactionShippingEtops)); err != nil {
			return err
		}
		mtseIDs := make([]int64, len(moneyTransactionShippingEtops))
		for i, mtse := range moneyTransactionShippingEtops {
			mtseIDs[i] = mtse.ID
		}
		var moneyTransactionShippings []*txmodel.MoneyTransactionShipping
		if err := x.Table("money_transaction_shipping").In("money_transaction_shipping_etop_id", mtseIDs).Find((*txmodel.MoneyTransactionShippings)(&moneyTransactionShippings)); err != nil {
			return err
		}
		mtsIDs := make([]int64, len(moneyTransactionShippings))
		for i, mt := range moneyTransactionShippings {
			mtsIDs[i] = mt.ID
		}
		mtQuery := &modelx.GetMoneyTransactions{
			IDs:                 mtsIDs,
			IncludeFulfillments: true,
		}
		if err := bus.Dispatch(ctx, mtQuery); err != nil {
			return err
		}
		moneyTransactionsMap := make(map[int64][]*txmodely.MoneyTransactionExtended)
		for _, mt := range mtQuery.Result.MoneyTransactions {
			moneyTransactionsMap[mt.MoneyTransactionShippingEtopID] = append(moneyTransactionsMap[mt.MoneyTransactionShippingEtopID], mt)
		}
		result := make([]*txmodely.MoneyTransactionShippingEtopExtended, len(mtseIDs))
		for i, mtse := range moneyTransactionShippingEtops {
			result[i] = &txmodely.MoneyTransactionShippingEtopExtended{
				MoneyTransactionShippingEtop: mtse,
				MoneyTransactions:            moneyTransactionsMap[mtse.ID],
			}
		}
		query.Result.MoneyTransactionShippingEtops = result
	}
	{
		total, err := s.Count(&txmodel.MoneyTransactionShippingEtop{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func UpdateMoneyTransactionShippingEtop(ctx context.Context, cmd *modelx.UpdateMoneyTransactionShippingEtop) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing transaction ID", nil)
	}
	{
		s1 := x.Table("money_transaction_shipping_etop").Where("id = ?", cmd.ID)
		var transaction = new(txmodel.MoneyTransactionShippingEtop)
		if err := s1.ShouldGet(transaction); err != nil {
			return err
		}
		if transaction.Status == model.S3Positive {
			return cm.Error(cm.InvalidArgument, "This money transaction was confirmed. Can not update!", nil)
		}
	}
	s2 := x.Table("money_transaction_shipping").Where("money_transaction_shipping_etop_id = ?", cmd.ID)
	var moneyTransactionShippings []*txmodel.MoneyTransactionShipping
	if err := s2.Find((*txmodel.MoneyTransactionShippings)(&moneyTransactionShippings)); err != nil {
		return err
	}
	oldIDs := make([]int64, len(moneyTransactionShippings))
	for i, mt := range moneyTransactionShippings {
		oldIDs[i] = mt.ID
	}
	newIDs := PatchID(oldIDs, cmd.Adds, cmd.Deletes, cmd.ReplaceAll)
	mtse, err := prepairMoneyTransactionShippingEtop(ctx, cmd.ID, newIDs)
	if err != nil {
		return err
	}
	mtse.Note = cmd.Note
	mtse.InvoiceNumber = cmd.InvoiceNumber
	mtse.BankAccount = cmd.BankAccount

	err = inTransaction(func(s Qx) error {
		if len(oldIDs) > 0 {
			if err := s.Table("money_transaction_shipping").In("id", oldIDs).ShouldUpdateMap(M{
				"money_transaction_shipping_etop_id": nil,
			}); err != nil {
				return err
			}
		}
		if err := s.Table("money_transaction_shipping_etop").Where("id = ?", cmd.ID).
			ShouldUpdate(mtse); err != nil {
			return err
		}
		if len(newIDs) > 0 {
			if err := s.Table("money_transaction_shipping").In("id", newIDs).
				ShouldUpdateMap(M{"money_transaction_shipping_etop_id": cmd.ID}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	query := &modelx.GetMoneyTransactionShippingEtop{
		ID: cmd.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	cmd.Result = query.Result
	return nil
}

func PatchID(list []int64, adds, deletes, replaceAll []int64) []int64 {
	if len(replaceAll) > 0 {
		return replaceAll
	}
	newList := make([]int64, 0, len(list)+len(adds))
	for _, id := range list {
		if !cm.ContainInt64(newList, id) && !cm.ContainInt64(deletes, id) {
			newList = append(newList, id)
		}
	}
	for _, id := range adds {
		if !cm.ContainInt64(newList, id) && !cm.ContainInt64(deletes, id) {
			newList = append(newList, id)
		}
	}
	return newList
}

func prepairMoneyTransactionShippingEtop(ctx context.Context, mtseID int64, mtIDs []int64) (*txmodel.MoneyTransactionShippingEtop, error) {
	query := &modelx.GetMoneyTransactions{
		IDs:                 mtIDs,
		IncludeFulfillments: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	moneyTransactions := query.Result.MoneyTransactions
	if len(mtIDs) != len(moneyTransactions) {
		var errID int64
		stop := false
		for _, id := range mtIDs {
			for _, mt := range moneyTransactions {
				if id != mt.ID {
					errID = id
					stop = true
					break
				}
			}
			if stop {
				break
			}
		}
		return nil, cm.Errorf(cm.InvalidArgument, nil, "MoneyTransactionShipping does not exist. (money_transaction_shiping_id = %v)", errID)
	}
	var fulfillments []*shipmodel.Fulfillment
	for _, mt := range moneyTransactions {
		if mt.Status != model.S3Zero {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "MoneyTransactionShipping does not valid. (money_transaction_shipping_id = %v)", mt.ID)
		}
		if mt.MoneyTransactionShippingEtopID != 0 && mt.MoneyTransactionShippingEtopID != mtseID {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "MoneyTransactionShipping belongs to another MoneyTransactionShippingEtop. (money_transaction_shipping_id = %v, money_transaction_shipping_etop_id = %v)", mt.ID, mt.MoneyTransactionShippingEtopID)
		}
		for _, ffm := range mt.Fulfillments {
			fulfillments = append(fulfillments, ffm.Fulfillment)
		}
	}
	totalCOD, totalAmount, totalOrders, totalShippingFee, _ := CalcFulfillmentsInfo(fulfillments)
	mtse := &txmodel.MoneyTransactionShippingEtop{
		TotalCOD:              totalCOD,
		TotalOrders:           totalOrders,
		TotalAmount:           totalAmount,
		TotalFee:              totalShippingFee,
		TotalMoneyTransaction: len(mtIDs),
	}
	return mtse, nil
}

func DeleteMoneyTransactionShippingEtop(ctx context.Context, cmd *modelx.DeleteMoneyTransactionShippingEtop) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing MoneyTransaction ID", nil)
	}
	var mtse = new(txmodel.MoneyTransactionShippingEtop)
	if err := x.Table("money_transaction_shipping_etop").Where("id = ?", cmd.ID).ShouldGet(mtse); err != nil {
		return err
	}
	if mtse.Status == model.S3Positive {
		return cm.Error(cm.FailedPrecondition, "MoneyTransaction was confirmed. Can not delete.", nil)
	}
	var moneyTransactionShippings []*txmodel.MoneyTransactionShipping
	if err := x.Table("money_transaction_shipping").Where("money_transaction_shipping_etop_id = ?", cmd.ID).Find((*txmodel.MoneyTransactionShippings)(&moneyTransactionShippings)); err != nil {
		return err
	}
	var mtIDs = make([]int64, len(moneyTransactionShippings))
	for i, mt := range moneyTransactionShippings {
		if mt.Status == model.S3Positive {
			return cm.Errorf(cm.FailedPrecondition, nil, "Can not delete this MoneyTransactionShippingEtop. This MoneyTransactionShipping (id = %v) was confirmed", mt.ID)
		}
		mtIDs[i] = mt.ID
	}
	return inTransaction(func(s Qx) error {
		if len(mtIDs) > 0 {
			if err := s.Table("money_transaction_shipping").In("id", mtIDs).ShouldUpdateMap(M{
				"money_transaction_shipping_etop_id": nil,
			}); err != nil {
				return err
			}
		}
		if err := s.Table("money_transaction_shipping_etop").Where("id = ?", cmd.ID).
			ShouldDelete(&txmodel.MoneyTransactionShippingEtop{}); err != nil {
			return err
		}
		cmd.Result.Deleted = 1
		return nil
	})
}

func ConfirmMoneyTransactionShippingEtop(ctx context.Context, cmd *modelx.ConfirmMoneyTransactionShippingEtop) error {
	if cmd.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing MoneyTransaction ID", nil)
	}
	query := &modelx.GetMoneyTransactionShippingEtop{
		ID: cmd.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	mtse := query.Result.MoneyTransactionShippingEtop
	moneyTransactionShippings := query.Result.MoneyTransactions
	if mtse.Status != model.S3Zero {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this MoneyTransactionShippingEtop")
	}

	var _moneyTransactions = make([]*txmodel.MoneyTransactionShipping, len(moneyTransactionShippings))
	totalCOD, totalAmount, totalOrders, totalFee := 0, 0, 0, 0
	var ffmIDs []int64
	for i, mt := range moneyTransactionShippings {
		if mt.Status != model.S3Zero {
			return cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this MoneyTransactionShipping (money_transaction_shipping_id = %v).", mt.ID)
		}
		var _ffms = make([]*shipmodel.Fulfillment, len(mt.Fulfillments))
		for j, ffm := range mt.Fulfillments {
			if !cm.StringsContain(acceptStates, string(ffm.ShippingState)) {
				return cm.Error(cm.FailedPrecondition, "Fulfillment #"+ffm.ShippingCode+" does not valid. Status must be delivered or returning or returned.", nil)
			}
			_ffms[j] = ffm.Fulfillment
		}
		_totalCOD, _totalAmount, _totalOrders, _totalFee, _ffmIDs := CalcFulfillmentsInfo(_ffms)
		_moneyTransactions[i] = &txmodel.MoneyTransactionShipping{
			ID:          mt.ID,
			TotalCOD:    _totalCOD,
			TotalAmount: _totalAmount,
			TotalOrders: _totalOrders,
		}
		totalCOD += _totalCOD
		totalAmount += _totalAmount
		totalOrders += _totalOrders
		totalFee += _totalFee
		ffmIDs = append(ffmIDs, _ffmIDs...)
	}
	if totalCOD != cmd.TotalCOD {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total COD does not match. (expected_total_cod = %v)", totalCOD)
	}
	if totalAmount != cmd.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total Amount does not match. (expected_total_amount = %v)", totalAmount)
	}
	if totalOrders != cmd.TotalOrders {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total Orders does not match. (expected_total_order = %v)", totalOrders)
	}

	return inTransaction(func(s Qx) error {
		now := time.Now()
		for _, mt := range _moneyTransactions {
			if err := s.Table("money_transaction_shipping").Where("id = ?", mt.ID).
				ShouldUpdateMap(M{
					"total_orders":       mt.TotalOrders,
					"total_cod":          mt.TotalCOD,
					"total_amount":       mt.TotalAmount,
					"status":             model.S3Positive,
					"confirmed_at":       now,
					"etop_transfered_at": now,
				}); err != nil {
				return err
			}
		}
		if err := s.Table("fulfillment").In("id", ffmIDs).
			ShouldUpdateMap(M{
				"cod_etop_transfered_at": now,
			}); err != nil {
			return err
		}
		if err := s.Table("money_transaction_shipping_etop").Where("id = ?", cmd.ID).
			ShouldUpdateMap(M{
				"total_orders":            totalOrders,
				"total_cod":               totalCOD,
				"total_amount":            totalAmount,
				"status":                  model.S3Positive,
				"total_fee":               totalFee,
				"total_money_transaction": len(_moneyTransactions),
				"confirmed_at":            now,
			}); err != nil {
			return err
		}
		cmd.Result.Updated = 1

		return nil
	})
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

	returnedFee := model.GetReturnedFee(ffm.ShippingFeeShopLines)
	totalFee := model.GetTotalShippingFee(ffm.ShippingFeeShopLines)
	newReturnedFee := (totalFee - returnedFee) / 2
	return newReturnedFee
}

func UpdateVtpostShippingFeeReturned(ffms []*shipmodel.Fulfillment) error {
	var updateFFms []*shipmodel.Fulfillment
	for _, ffm := range ffms {
		if ffm.ShippingState != model.StateReturned && ffm.ShippingState != model.StateReturning {
			continue
		}
		returnedFee := model.GetReturnedFee(ffm.ShippingFeeShopLines)
		newReturnedFee := CalcVtpostShippingFeeReturned(ffm)
		if newReturnedFee == 0 || newReturnedFee == returnedFee {
			continue
		}
		lines := ffm.ProviderShippingFeeLines
		ffm.ProviderShippingFeeLines = model.UpdateShippingFees(lines, newReturnedFee, model.ShippingFeeTypeReturn)
		ffm.ShippingFeeShopLines = model.GetShippingFeeShopLines(ffm.ProviderShippingFeeLines, ffm.EtopPriceRule, &ffm.EtopAdjustedShippingFeeMain)
		updateFFms = append(updateFFms, &shipmodel.Fulfillment{
			ID:                       ffm.ID,
			ProviderShippingFeeLines: ffm.ProviderShippingFeeLines,
			ShippingFeeShopLines:     ffm.ShippingFeeShopLines,
		})
	}

	if len(updateFFms) == 0 {
		return nil
	}
	cmd := &shipmodelx.UpdateFulfillmentsCommand{
		Fulfillments: updateFFms,
	}
	ctx := context.Background()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
