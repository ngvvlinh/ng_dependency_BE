package aggregate

import (
	"context"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	moneytxsqlstore "o.o/backend/com/main/moneytx/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

func (a *MoneyTxAggregate) CreateMoneyTxShippingExternal(ctx context.Context, args *moneytx.CreateMoneyTxShippingExternalArgs) (*moneytx.MoneyTransactionShippingExternalFtLine, error) {
	if len(args.Lines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Danh sách Vận đơn không được rỗng")
	}

	connectionID := shipping.GetConnectionID(args.ConnectionID, args.Provider)
	if connectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển")
	}

	totalCOD := 0
	totalOrders := 0
	for _, line := range args.Lines {
		totalCOD += line.ExternalTotalCOD
		totalOrders++
	}

	externalTxID := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		subCode := args.Provider.String()
		subCode = strings.ToUpper(subCode)
		code, errCode := sqlstore.GenerateCode(ctx, tx, model.CodeTypeMoneyTransactionExternal, subCode)
		if errCode != nil {
			return errCode
		}
		externalTx := &moneytx.MoneyTransactionShippingExternal{
			ID:             externalTxID,
			Code:           code,
			TotalCOD:       totalCOD,
			TotalOrders:    totalOrders,
			ExternalPaidAt: args.ExternalPaidAt,
			Provider:       args.Provider,
			BankAccount:    args.BankAccount,
			Note:           args.Note,
			InvoiceNumber:  args.InvoiceNumber,
			ConnectionID:   connectionID,
		}
		if err := a.moneyTxShippingExternalStore(ctx).CreateMoneyTxShippingExternal(externalTx); err != nil {
			return err
		}

		ffmIDs := make([]dot.ID, 0, len(args.Lines))
		for _, line := range args.Lines {
			createCmd := &moneytx.CreateMoneyTxShippingExternalLineArgs{
				ExternalCode:                       line.ExternalCode,
				ExternalTotalCOD:                   line.ExternalTotalCOD,
				ExternalCreatedAt:                  line.ExternalCreatedAt,
				ExternalClosedAt:                   line.ExternalClosedAt,
				EtopFulfillmentID:                  line.EtopFulfillmentID,
				EtopFulfillmentIDRaw:               line.EtopFulfillmentIDRaw,
				ExternalCustomer:                   line.ExternalCustomer,
				ExternalAddress:                    line.ExternalAddress,
				MoneyTransactionShippingExternalID: externalTx.ID,
				ExternalTotalShippingFee:           line.ExternalTotalShippingFee,
			}
			externalLine, err := a.CreateMoneyTxShippingExternalLine(ctx, createCmd)
			if err != nil {
				return err
			}
			if externalLine.EtopFulfillmentID != 0 && externalLine.ImportError == nil {
				ffmIDs = append(ffmIDs, externalLine.EtopFulfillmentID)
			}
		}

		event := &moneytx.MoneyTxShippingExternalCreatedEvent{
			EventMeta:                 meta.NewEvent(),
			MoneyTxShippingExternalID: externalTxID,
			FulfillementIDs:           ffmIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingExternalStore(ctx).ID(externalTxID).GetMoneyTxShippingExternalFtLine()
}

func (a *MoneyTxAggregate) CreateMoneyTxShippingExternalLine(ctx context.Context, args *moneytx.CreateMoneyTxShippingExternalLineArgs) (*moneytx.MoneyTransactionShippingExternalLine, error) {
	if args.MoneyTransactionShippingExternalID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing MoneyTransactionShippingExternalID", nil)
	}
	line := &moneytx.MoneyTransactionShippingExternalLine{
		ID:                                 cm.NewID(),
		ExternalCode:                       args.ExternalCode,
		ExternalTotalCOD:                   args.ExternalTotalCOD,
		ExternalCreatedAt:                  args.ExternalCreatedAt,
		ExternalClosedAt:                   args.ExternalClosedAt,
		ExternalCustomer:                   args.ExternalCustomer,
		ExternalAddress:                    args.ExternalAddress,
		EtopFulfillmentIDRaw:               args.EtopFulfillmentIDRaw,
		EtopFulfillmentID:                  args.EtopFulfillmentID,
		MoneyTransactionShippingExternalID: args.MoneyTransactionShippingExternalID,
		ExternalTotalShippingFee:           args.ExternalTotalShippingFee,
	}
	if line.ExternalCode == "" && line.EtopFulfillmentID == 0 {
		line.ImportError = &meta.Error{
			Code: "ffm_id_empty",
			Msg:  "Thiếu mã vận đơn",
		}
	} else {
		query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
			ID:           line.EtopFulfillmentID,
			ShippingCode: line.ExternalCode,
			Result:       nil,
		}
		if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
			line.ImportError = &meta.Error{
				Code: "ffm_not_found",
				Msg:  "Không tìm thấy vận đơn trên Etop",
			}
		} else {
			ffm := query.Result
			line.EtopFulfillmentID = ffm.ID
			line.ExternalCode = ffm.ShippingCode
			if ffm.MoneyTransactionShippingExternalID != 0 {
				line.ImportError = &meta.Error{
					Code: "ffm_exist_money_transaction_shipping_external",
					Msg:  "Vận đơn nằm trong phiên thanh toán nhà vận chuyển khác: " + strconv.Itoa(int(ffm.MoneyTransactionShippingExternalID)),
				}
			} else if !cm.StringsContain(moneytx.ShippingAcceptStates, ffm.ShippingState.String()) {
				line.ImportError = &meta.Error{
					Code: "ffm_not_done",
					Msg:  "Vận đơn chưa hoàn thành trên Etop",
				}
			} else if ffm.ShippingState == shippingstate.Delivered && ffm.TotalCODAmount != line.ExternalTotalCOD {
				line.ImportError = &meta.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị vận đơn không đúng",
					Meta: map[string]string{
						"Etop":     strconv.Itoa(ffm.TotalCODAmount),
						"Provider": strconv.Itoa(line.ExternalTotalCOD),
					},
				}
			} else if ffm.ShippingState == shippingstate.Undeliverable && line.ExternalTotalCOD != ffm.ActualCompensationAmount {
				line.ImportError = &meta.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị bồi hoàn không đúng",
				}
			} else if ffm.MoneyTransactionID != 0 {
				line.ImportError = &meta.Error{
					Code: "ffm_exist_money_transaction",
					Msg:  "Vận đơn nằm trong phiên thanh toán khác: " + strconv.Itoa(int(ffm.MoneyTransactionID)),
				}
			} else if line.ExternalTotalShippingFee != 0 && line.ExternalTotalShippingFee != ffm.ShippingFeeShop {
				line.ImportError = &meta.Error{
					Code: "ffm_shipping_fee_not_match",
					Msg:  "Tổng tiền cước không đúng.",
					Meta: map[string]string{
						"Etop":     strconv.Itoa(ffm.ShippingFeeShop),
						"Provider": strconv.Itoa(line.ExternalTotalShippingFee),
					},
				}
			} else if ffm.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
				if ffm.ShippingType == 0 {
					// backward compatible
					// remove later
					// no error
				} else {
					line.ImportError = &meta.Error{
						Code: "ffm_not_in_etop",
						Msg:  "Vận đơn không được đối soát bởi Etop",
					}
				}
			}
		}
	}
	if err := a.moneyTxShippingExternalStore(ctx).CreateMoneyTxShippingExternalLine(line); err != nil {
		return nil, err
	}

	return line, nil
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingExternalInfo(ctx context.Context, args *moneytx.UpdateMoneyTxShippingExternalInfoArgs) (*moneytx.MoneyTransactionShippingExternalFtLine, error) {
	if args.MoneyTxShippingExternalID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing money transaction ID")
	}
	moneyTx, err := a.moneyTxShippingExternalStore(ctx).ID(args.MoneyTxShippingExternalID).
		GetMoneyTxShippingExternal()
	if err != nil {
		return nil, err
	}
	if moneyTx.Status == status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update this money transaction")
	}
	if err := a.moneyTxShippingExternalStore(ctx).UpdateMoneyTxShippingExternalInfo(args); err != nil {
		return nil, err
	}
	return a.moneyTxShippingExternalStore(ctx).ID(args.MoneyTxShippingExternalID).GetMoneyTxShippingExternalFtLine()
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShippingExternal(ctx context.Context, id dot.ID) (updated int, _ error) {
	panic("implement me")
}

/*
# Khi tạo phiên thanh toán cho Shop (ConfirmMoneyTxShippingExternals)
	TopShip sẽ tự động đối soát với tất cả NVC trường hợp:
	- returned
	- returning
	- COD = 0 (state: delivered & total_cod_amount = 0)
*/

func (a *MoneyTxAggregate) ConfirmMoneyTxShippingExternals(ctx context.Context, ids []dot.ID) (updated int, _ error) {
	if len(ids) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing money transaction shipping external IDs")
	}
	moneyTxs, err := a.moneyTxShippingExternalStore(ctx).IDs(ids...).ListMoneyTxShippingExternalsFtLine()
	if err != nil {
		return 0, err
	}

	var moneyTxExternalIDs []dot.ID
	shopFfmsMap := make(map[dot.ID][]*shipping.Fulfillment)
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// raise event confirming
		// cập nhật phí trả hàng vtpost (requirement: Lọc tất cả đơn trả hàng, đang trả hàng của VTPOST, tính toán & thêm phí trả hàng vào)
		// Công thức tính cụ thể xem ở shipping pm
		event := &moneytx.MoneyTxShippingExternalsConfirmingEvent{
			EventMeta:                  meta.NewEvent(),
			MoneyTxShippingExternalIDs: moneyTxExternalIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		for _, moneyTx := range moneyTxs {
			_shopFfmsMap, err := a.preprocessConfirmMoneyTxExternal(ctx, moneyTx)
			if err != nil {
				return err
			}
			moneyTxExternalIDs = append(moneyTxExternalIDs, moneyTx.ID)
			for shopID, ffms := range _shopFfmsMap {
				shopFfmsMap[shopID] = mergeFulfillments(shopFfmsMap[shopID], ffms)
			}
		}

		shopFfmsAdditionMap, err := a.combineWithExtraFfms(ctx)
		if err != nil {
			return err
		}
		// make sure do not dupplicate ffm
		for shopID, ffms := range shopFfmsAdditionMap {
			shopFfmsMap[shopID] = mergeFulfillments(shopFfmsMap[shopID], ffms)
		}

		cmd := &moneytx.CreateMoneyTxShippingsArgs{
			ShopIDMapFfms: shopFfmsMap,
		}
		if _, err := a.CreateMoneyTxShippings(ctx, cmd); err != nil {
			return err
		}
		return a.moneyTxShippingExternalStore(ctx).ConfirmMoneyTxShippingExternals(moneyTxExternalIDs)
	})
	if err != nil {
		return 0, err
	}
	return len(moneyTxExternalIDs), nil
}

func (a *MoneyTxAggregate) RemoveMoneyTxShippingExternalLines(ctx context.Context, args *moneytx.RemoveMoneyTxShippingExternalLinesArgs) (*moneytx.MoneyTransactionShippingExternalFtLine, error) {
	if args.MoneyTxShippingExternalID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing money transaction shipping external ID")
	}
	moneyTx, err := a.moneyTxShippingExternalStore(ctx).ID(args.MoneyTxShippingExternalID).GetMoneyTxShippingExternal()
	if err != nil {
		return nil, err
	}
	if moneyTx.Status == status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update money transaction shipping external")
	}
	if len(args.LineIDs) == 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "LineIDs can not be empty")
	}

	lines, err := a.moneyTxShippingExternalStore(ctx).Line_by_MoneyTxShippingExternalID(args.MoneyTxShippingExternalID).ListMoneyTxShippingExternalLinesDB()
	if err != nil {
		return nil, err
	}

	ffmIDs := make([]dot.ID, 0, len(args.LineIDs))
	for _, id := range args.LineIDs {
		found := false
		for _, line := range lines {
			if id == line.ID {
				found = true
				if line.EtopFulfillmentID != 0 && line.ImportError == nil {
					ffmIDs = append(ffmIDs, line.EtopFulfillmentID)
				}
				break
			}
		}
		if !found {
			return nil, cm.Errorf(cm.NotFound, nil, "Line #%v does not exist in this money transaction", id)
		}
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		totalCOD, totalOrders := 0, 0
		for _, line := range lines {
			if cm.IDsContain(args.LineIDs, line.ID) {
				continue
			}
			totalCOD += line.ExternalTotalCOD
			totalOrders++
		}

		if err := a.moneyTxShippingExternalStore(ctx).Line_by_LineIDs(args.LineIDs...).DeleteMoneyTxShippingExternalLines(); err != nil {
			return err
		}
		event := &moneytx.MoneyTxShippingExternalLinesDeletedEvent{
			EventMeta:                 meta.NewEvent(),
			MoneyTxShippingExternalID: args.MoneyTxShippingExternalID,
			FulfillmentIDs:            ffmIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		update := &moneytxsqlstore.UpdateMoneyTxShippingExternalStatisticsArgs{
			MoneyTxShippingExternalID: args.MoneyTxShippingExternalID,
			TotalCOD:                  dot.Int(totalCOD),
			TotalOrders:               dot.Int(totalOrders),
		}
		if err := a.moneyTxShippingExternalStore(ctx).UpdateMoneyTxShippingExternalStatistics(update); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingExternalStore(ctx).ID(args.MoneyTxShippingExternalID).GetMoneyTxShippingExternalFtLine()
}

func (a *MoneyTxAggregate) DeleteMoneyTxShippingExternal(ctx context.Context, id dot.ID) (deleted int, _ error) {
	if id == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing money transaction ID").WithMetap("aggregate", "DeleteMoneyTxShippingExternal")
	}
	moneyTx, err := a.moneyTxShippingExternalStore(ctx).ID(id).GetMoneyTxShippingExternal()
	if err != nil {
		return 0, err
	}
	if moneyTx.Status == status3.P {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Cannot delete this money transaction external").WithMetap("aggregate", "DeleteMoneyTxShippingExternal")
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// raise event:
		// deleted money transaction shipping external line
		// remove money_transaction_shipping_id in fulfillment
		event := &moneytx.MoneyTxShippingExternalDeletedEvent{
			EventMeta:                 meta.NewEvent(),
			MoneyTxShippingExternalID: id,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return a.moneyTxShippingExternalStore(ctx).DeleteMoneyTxShippingExternal(id)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *MoneyTxAggregate) DeleteMoneyTxShippingExternalLines(ctx context.Context, moneyTxShippingExternalID dot.ID) error {
	if moneyTxShippingExternalID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing money_tx_shipping_external_id")
	}
	return a.moneyTxShippingExternalStore(ctx).Line_by_MoneyTxShippingExternalID(moneyTxShippingExternalID).DeleteMoneyTxShippingExternalLines()
}

func (a *MoneyTxAggregate) preprocessConfirmMoneyTxExternal(ctx context.Context, moneyTx *moneytx.MoneyTransactionShippingExternalFtLine) (shopFfmIDMap map[dot.ID][]*shipping.Fulfillment, _err error) {
	shopFfmIDMap = make(map[dot.ID][]*shipping.Fulfillment)
	if moneyTx.Status != status3.Z {
		return shopFfmIDMap, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this money transaction").WithMetap("id", moneyTx.ID)
	}
	lines := moneyTx.Lines
	if len(lines) == 0 {
		return shopFfmIDMap, cm.Errorf(cm.FailedPrecondition, nil, "There are no lines in this money transaction").WithMetap("id", moneyTx.ID)
	}

	ffmCodes := make([]string, len(lines))
	for i, line := range lines {
		if line.ImportError != nil && line.ImportError.Code != "" {
			return shopFfmIDMap, cm.Errorf(cm.FailedPrecondition, nil, "Please handle error before confirm money transaction").WithMetap("id", moneyTx.ID)
		}
		ffmCodes[i] = line.ExternalCode
	}
	ffmQuery := &shipping.ListFulfillmentsByShippingCodesQuery{
		Codes: ffmCodes,
	}
	if err := a.shippingQuery.Dispatch(ctx, ffmQuery); err != nil {
		return nil, err
	}
	fulfillments := ffmQuery.Result

	for _, line := range lines {
		found := false
		for _, ffm := range fulfillments {
			if line.ExternalCode == ffm.ShippingCode {
				found = true
				break
			}
		}
		if !found {
			return shopFfmIDMap, cm.Errorf(cm.NotFound, nil, "Fulfillment not found %v", line.ExternalCode)
		}
	}
	for _, ffm := range fulfillments {
		shopFfmIDMap[ffm.ShopID] = append(shopFfmIDMap[ffm.ShopID], ffm)
	}
	return shopFfmIDMap, nil
}

func mergeFulfillments(ffms []*shipping.Fulfillment, subFfms []*shipping.Fulfillment) []*shipping.Fulfillment {
	mergeFfms := append(ffms, subFfms...)
	ffmMap := make(map[dot.ID]*shipping.Fulfillment)
	for _, ffm := range mergeFfms {
		ffmMap[ffm.ID] = ffm
	}
	var res []*shipping.Fulfillment
	for _, ffm := range ffmMap {
		res = append(res, ffm)
	}
	return res
}

/*
	+ GHN
		- returned
		- COD = 0 (state: delivered & total_cod_amount = 0)
	+ Vtpost
		- returned
		- returning
		- COD = 0 (state: delivered & total_cod_amount = 0)

	Update 06/10/2020 (tuan)
		TopShip sẽ tự động đối soát với tất cả NVC trường hợp:
		- returned
		- returning
		- COD = 0 (state: delivered & total_cod_amount = 0)
*/
func (a *MoneyTxAggregate) combineWithExtraFfms(ctx context.Context) (shopFfmsMap map[dot.ID][]*shipping.Fulfillment, _ error) {
	shopFfmsMap = make(map[dot.ID][]*shipping.Fulfillment)
	query := &shipping.ListFulfillmentsForMoneyTxQuery{
		ShippingStates: []shippingstate.State{shippingstate.Returning, shippingstate.Returned},
		IsNoneCOD:      dot.Bool(true),
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	ffms := query.Result
	for _, ffm := range ffms {
		// Trường hợp đơn returning
		// Chỉ cho vào phiên nếu đơn đã chuyển sang trạng thái returning tối thiểu 5 ngày
		now := time.Now()
		if ffm.ShippingState == shippingstate.Returning && now.Sub(ffm.ExternalShippingReturningAt) < 5*24*time.Hour {
			continue
		}
		shopFfmsMap[ffm.ShopID] = append(shopFfmsMap[ffm.ShopID], ffm)
	}
	return shopFfmsMap, nil
}

const (
	_PartnerShop  = "partner_shop"
	_PriorityShop = "priority_shop"
	_NewShop      = "new_shop"
)

/* SplitMoneyTransactionShippingExternal

	Chia phiên NVC thành nhiều phiên:
	1. Phiên chứa các đơn của shop từ đối tác (ưu tiên cao nhất)
    2. Phiên chứa các đơn của các:
        shop ưu tiên đối soát trước (shop.is_prior_money_transaction)
	3. Phiên chứa các đơn của các:
		shop mới (shop có số phiên đối soát <= args.MaxMoneyTxShippingCount)
	4. Phiên chứa các đơn còn lại
	Xóa bỏ phiên cũ sau khi tách xong
*/
func (a *MoneyTxAggregate) SplitMoneyTxShippingExternal(ctx context.Context, args *moneytx.SplitMoneyTxShippingExternalArgs) (*pbcm.UpdatedResponse, error) {
	if args.MoneyTxShippingExternalID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing money tx shipping external ID")
	}
	if !args.IsSplitByShopPriority && args.MaxMoneyTxShippingCount == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing splitting method")
	}

	moneyTxExternal, err := a.moneyTxShippingExternalStore(ctx).ID(args.MoneyTxShippingExternalID).GetMoneyTxShippingExternalFtLine()
	if err != nil {
		return nil, err
	}
	if moneyTxExternal.Status != status3.Z {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not split this money transaction shipping external")
	}

	ffmIDs := []dot.ID{}
	for _, line := range moneyTxExternal.Lines {
		if line.EtopFulfillmentID == 0 || line.ImportError != nil {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Please fix errors before continuing process")
		}
		ffmIDs = append(ffmIDs, line.EtopFulfillmentID)
	}

	query := &shipping.ListFulfillmentsByIDsQuery{
		IDs: ffmIDs,
	}
	if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	if len(moneyTxExternal.Lines) != len(query.Result) {
		return nil, cm.Errorf(cm.Internal, nil, "The quantity of fulfillments is not enough")
	}

	type moneyTxShippingExternal struct {
		name                  string
		mapShopIDs            map[dot.ID]bool
		shippingExternalLines []*moneytx.MoneyTransactionShippingExternalLine
	}
	moneyTxShippingExternals := map[string]*moneyTxShippingExternal{
		_PartnerShop: {
			name:                  "Phiên shop partner",
			mapShopIDs:            make(map[dot.ID]bool),
			shippingExternalLines: []*moneytx.MoneyTransactionShippingExternalLine{},
		},
		_PriorityShop: {
			name:                  "Phiên shop ưu tiên",
			mapShopIDs:            make(map[dot.ID]bool),
			shippingExternalLines: []*moneytx.MoneyTransactionShippingExternalLine{},
		},
		_NewShop: {
			name:                  "Phiên shop mới",
			mapShopIDs:            make(map[dot.ID]bool),
			shippingExternalLines: []*moneytx.MoneyTransactionShippingExternalLine{},
		},
	}

	mapShopIDs := make(map[dot.ID]bool)
	shopIDs := []dot.ID{}
	for _, ffm := range query.Result {
		if mapShopIDs[ffm.ShopID] {
			continue
		}
		mapShopIDs[ffm.ShopID] = true
		shopIDs = append(shopIDs, ffm.ShopID)
	}

	mapPartnerShopIDs, err := a.buildMapPartnerShopID(ctx, shopIDs)
	if err != nil {
		return nil, err
	}
	if partnerShop, ok := moneyTxShippingExternals[_PartnerShop]; ok {
		partnerShop.mapShopIDs = mapPartnerShopIDs
	}

	if args.IsSplitByShopPriority {
		queryShop := &identity.ListShopsByIDsQuery{
			IDs:                     shopIDs,
			IncludeWLPartnerShop:    true,
			IsPriorMoneyTransaction: true,
		}
		if err := a.identityQuery.Dispatch(ctx, queryShop); err != nil {
			return nil, err
		}
		for _, shop := range queryShop.Result {
			moneyTxShippingExternals[_PriorityShop].mapShopIDs[shop.ID] = true
		}
	}
	if args.MaxMoneyTxShippingCount > 0 {
		queryCount := &moneytx.CountMoneyTxShippingByShopIDsArgs{
			ShopIDs: shopIDs,
		}
		shopCounts, err := a.moneyTxShippingStore(ctx).CountMoneyTxShippingByShopIDs(queryCount)
		if err != nil {
			return nil, err
		}
		mapShopFtMoneyTxCount := make(map[dot.ID]int)
		for _, shop := range shopCounts {
			mapShopFtMoneyTxCount[shop.ShopID] = shop.MoneyTxShippingCount
		}
		for _, shopID := range shopIDs {
			count, ok := mapShopFtMoneyTxCount[shopID]
			// !ok <=> shop chưa có phiên nào => là shop mới
			if !ok || count <= args.MaxMoneyTxShippingCount {
				moneyTxShippingExternals[_NewShop].mapShopIDs[shopID] = true
			}
		}
	}

	remainlines := []*moneytx.MoneyTransactionShippingExternalLine{}
	for _, ffm := range query.Result {
		moneyTxExternalLine := &moneytx.MoneyTransactionShippingExternalLine{
			ExternalCode:     ffm.ShippingCode,
			ExternalTotalCOD: ffm.TotalCODAmount,
		}
		partnerMoneyTxExternal, _ := moneyTxShippingExternals[_PartnerShop]
		priorityMoneyTxExternal := moneyTxShippingExternals[_PriorityShop]
		newMoneyTxExternal := moneyTxShippingExternals[_NewShop]
		if partnerMoneyTxExternal.mapShopIDs[ffm.ShopID] {
			partnerMoneyTxExternal.shippingExternalLines = append(partnerMoneyTxExternal.shippingExternalLines, moneyTxExternalLine)
		} else if priorityMoneyTxExternal.mapShopIDs[ffm.ShopID] {
			priorityMoneyTxExternal.shippingExternalLines = append(priorityMoneyTxExternal.shippingExternalLines, moneyTxExternalLine)
		} else if newMoneyTxExternal.mapShopIDs[ffm.ShopID] {
			newMoneyTxExternal.shippingExternalLines = append(newMoneyTxExternal.shippingExternalLines, moneyTxExternalLine)
		} else {
			remainlines = append(remainlines, moneyTxExternalLine)
		}
	}

	numMoneyTxExternalCreated := 0
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if len(remainlines) == len(moneyTxExternal.Lines) {
			return cm.Errorf(cm.FailedPrecondition, nil, "It doesn't need to split this money transaction shipping external.")
		}

		// xóa phiên cũ
		if _, err := a.DeleteMoneyTxShippingExternal(ctx, args.MoneyTxShippingExternalID); err != nil {
			return err
		}

		cmd := &moneytx.CreateMoneyTxShippingExternalArgs{
			Provider:       moneyTxExternal.Provider,
			ConnectionID:   moneyTxExternal.ConnectionID,
			ExternalPaidAt: moneyTxExternal.ExternalPaidAt,
			BankAccount:    moneyTxExternal.BankAccount,
			InvoiceNumber:  moneyTxExternal.InvoiceNumber,
		}
		for _, _moneyTx := range moneyTxShippingExternals {
			if len(_moneyTx.shippingExternalLines) == 0 {
				continue
			}
			cmdCreate := *cmd
			cmdCreate.Lines = _moneyTx.shippingExternalLines
			if moneyTxExternal.Note == "" {
				cmdCreate.Note = _moneyTx.name
			} else {
				cmdCreate.Note = _moneyTx.name + " - " + moneyTxExternal.Note
			}
			if _, err := a.CreateMoneyTxShippingExternal(ctx, &cmdCreate); err != nil {
				return err
			}
			numMoneyTxExternalCreated++
		}

		// tạo phiên chứa các đơn còn lại
		if len(remainlines) > 0 {
			cmdCreate := *cmd
			cmdCreate.Lines = remainlines
			if moneyTxExternal.Note == "" {
				cmdCreate.Note = "Phiên shop bình thường"
			} else {
				cmdCreate.Note = "Phiên shop bình thường" + " - " + moneyTxExternal.Note
			}
			if _, err := a.CreateMoneyTxShippingExternal(ctx, &cmdCreate); err != nil {
				return err
			}
			numMoneyTxExternalCreated++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: numMoneyTxExternalCreated}, nil
}

func (a *MoneyTxAggregate) buildMapPartnerShopID(ctx context.Context, shopIDs []dot.ID) (map[dot.ID]bool, error) {
	query := &identity.ListPartnerRelationsBySubjectIDsQuery{
		SubjectIDs:  shopIDs,
		SubjectType: identity.SubjectTypeAccount,
	}
	if err := a.identityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var res = make(map[dot.ID]bool)
	for _, relation := range query.Result {
		res[relation.SubjectID] = true
	}
	return res, nil
}
