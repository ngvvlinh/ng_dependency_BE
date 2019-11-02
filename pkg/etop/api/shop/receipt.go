package shop

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/shopping/vendoring"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		receiptService.CreateReceipt,
		receiptService.UpdateReceipt,
		receiptService.ConfirmReceipt,
		receiptService.CancelReceipt,
		receiptService.GetReceipt,
		receiptService.GetReceipts,
		receiptService.GetReceiptsByLedgerType)
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (_err error) {
	key := fmt.Sprintf("Create receipt %v-%v-%v-%v-%v-%v-%v-%v",
		q.Context.Shop.ID, q.Context.UserID, q.TraderId, q.LedgerId, q.Title, q.Description, q.Amount, q.Type)
	result, err := idempgroup.DoAndWrap(
		key, 15*time.Second,
		func() (interface{}, error) { return s.createReceipt(ctx, q) },
		"Create receipt")
	if err != nil {
		return err
	}
	q.Result = pbshop.PbReceipt(result.(*receipting.CreateReceiptCommand).Result)

	return nil
}

func (s *ReceiptService) createReceipt(ctx context.Context, q *wrapshop.CreateReceiptEndpoint) (_ *receipting.CreateReceiptCommand, err error) {
	var paidAt time.Time
	var checkHavePaidAt bool
	if q.PaidAt.Seconds != 0 {
		paidAt, err = ptypes.Timestamp(q.PaidAt)
		if err != nil {
			return nil, err
		}
		checkHavePaidAt = true
	}

	cmd := &receipting.CreateReceiptCommand{
		ShopID:      q.Context.Shop.ID,
		CreatedBy:   q.Context.UserID,
		TraderID:    q.TraderId,
		Title:       q.Title,
		Description: q.Description,
		Amount:      q.Amount,
		LedgerID:    q.LedgerId,
		Type:        receipting.ReceiptType(q.Type),
		CreatedType: receipting.ReceiptCreatedTypeManual,
		Status:      int32(model.S3Zero),
		Lines:       pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
	}
	if checkHavePaidAt {
		cmd.PaidAt = paidAt
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (s *ReceiptService) UpdateReceipt(ctx context.Context, q *wrapshop.UpdateReceiptEndpoint) (err error) {
	var paidAt time.Time
	var checkHavePaidAt bool
	if q.PaidAt.Seconds != 0 {
		paidAt, err = ptypes.Timestamp(&q.PaidAt)
		if err != nil {
			return err
		}
		checkHavePaidAt = true
	}

	cmd := &receipting.UpdateReceiptCommand{
		ID:          q.Id,
		ShopID:      q.Context.Shop.ID,
		Title:       PString(q.Title),
		Description: PString(q.Description),
		LedgerID:    PInt64(q.LedgerId),
		TraderID:    PInt64(q.TraderId),
		Amount:      PInt32(q.Amount),
		Lines:       pbshop.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
	}
	if checkHavePaidAt {
		cmd.PaidAt = paidAt
	}
	err = receiptAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}

	q.Result = pbshop.PbReceipt(cmd.Result)
	return nil
}

func (s *ReceiptService) ConfirmReceipt(ctx context.Context, q *wrapshop.ConfirmReceiptEndpoint) error {
	cmd := &receipting.ConfirmReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}

	return nil
}

func (s *ReceiptService) CancelReceipt(ctx context.Context, q *wrapshop.CancelReceiptEndpoint) error {
	cmd := &receipting.CancelReceiptCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
		Reason: q.Reason,
	}
	if err := receiptAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}

	return nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, q *wrapshop.GetReceiptEndpoint) error {
	// Check receipt is exist
	getReceiptQuery := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := receiptQuery.Dispatch(ctx, getReceiptQuery); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}

	if receipts, err := s.getInfosForReceipts(ctx, q.Context.Shop.ID, []*receipting.Receipt{getReceiptQuery.Result}); err != nil {
		return err
	} else {
		q.Result = receipts[0]
	}

	return nil
}

func (s *ReceiptService) GetReceipts(ctx context.Context, q *wrapshop.GetReceiptsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &receipting.ListReceiptsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	if receipts, err := s.getInfosForReceipts(ctx, q.Context.Shop.ID, query.Result.Receipts); err != nil {
		return err
	} else {
		q.Result = &pbshop.ReceiptsResponse{

			TotalAmountConfirmedReceipt: query.Result.TotalAmountConfirmedReceipt,
			TotalAmountConfirmedPayment: query.Result.TotalAmountConfirmedPayment,
			Receipts:                    receipts,
			Paging:                      pbcm.PbPageInfo(paging, query.Result.Count),
		}
	}

	return nil
}

func (s *ReceiptService) GetReceiptsByLedgerType(ctx context.Context, q *wrapshop.GetReceiptsByLedgerTypeEndpoint) error {
	paging := q.Paging.CMPaging()
	listLedgersByType := &ledgering.ListLedgersByTypeQuery{
		LedgerType: ledgering.LedgerType(q.Type),
		ShopID:     q.Context.Shop.ID,
	}
	if err := ledgerQuery.Dispatch(ctx, listLedgersByType); err != nil {
		return err
	}

	var ledgerIDs []int64
	for _, ledger := range listLedgersByType.Result.Ledgers {
		ledgerIDs = append(ledgerIDs, ledger.ID)
	}

	query := &receipting.ListReceiptsByLedgerIDsQuery{
		ShopID:    q.Context.Shop.ID,
		LedgerIDs: ledgerIDs,
		Paging:    *paging,
		Filters:   pbcm.ToFilters(q.Filters),
	}
	if err := receiptQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	if receipts, err := s.getInfosForReceipts(ctx, q.Context.Shop.ID, query.Result.Receipts); err != nil {
		return err
	} else {
		q.Result = &pbshop.ReceiptsResponse{
			TotalAmountConfirmedReceipt: query.Result.TotalAmountConfirmedReceipt,
			TotalAmountConfirmedPayment: query.Result.TotalAmountConfirmedPayment,
			Receipts:                    receipts,
			Paging:                      pbcm.PbPageInfo(paging, query.Result.Count),
		}
	}
	return nil
}

func (s *ReceiptService) getInfosForReceipts(ctx context.Context, shopID int64, receipts []*receipting.Receipt) (receiptsResult []*pbshop.Receipt, _ error) {
	mapOrderIDAndReceivedAmount := make(map[int64]int32)
	mapLedger := make(map[int64]*ledgering.ShopLedger)
	var refIDs, userIDs, traderIDs, ledgerIDs []int64

	receiptsResult = pbshop.PbReceipts(receipts)

	// Get ref_ids into receiptLines
	for _, receipt := range receipts {
		if receipt.CreatedBy != 0 {
			userIDs = append(userIDs, receipt.CreatedBy)
		}
		traderIDs = append(traderIDs, receipt.TraderID)
		ledgerIDs = append(ledgerIDs, receipt.LedgerID)

		for _, receiptLine := range receipt.Lines {
			if receiptLine.RefID == 0 {
				continue
			}
			if _, ok := mapOrderIDAndReceivedAmount[receiptLine.RefID]; !ok {
				refIDs = append(refIDs, receiptLine.RefID)
				mapOrderIDAndReceivedAmount[receiptLine.RefID] = 0
			}
		}
	}

	// Get all users into receipts
	getUsersOfCurrAccount := &model.GetAccountUserExtendedsQuery{
		AccountIDs: []int64{shopID},
	}
	if err := bus.Dispatch(ctx, getUsersOfCurrAccount); err != nil {
		return nil, err
	}
	mapUserIDAndUser := make(map[int64]*model.User)
	for _, accountUser := range getUsersOfCurrAccount.Result.AccountUsers {
		mapUserIDAndUser[accountUser.User.ID] = accountUser.User
	}
	for _, receipt := range receiptsResult {
		if receipt.CreatedBy != 0 {
			receipt.User = pbetop.PbUser(mapUserIDAndUser[receipt.CreatedBy])
		}
	}

	// List traders
	if err := listTraders(ctx, shopID, traderIDs, receiptsResult); err != nil {
		return nil, err
	}

	// List ledgers
	getLedgersByIDs := &ledgering.ListLedgersByIDsQuery{
		ShopID: shopID,
		IDs:    ledgerIDs,
	}
	if err := ledgerQuery.Dispatch(ctx, getLedgersByIDs); err != nil {
		return nil, err
	}
	for _, ledger := range getLedgersByIDs.Result.Ledgers {
		mapLedger[ledger.ID] = ledger
	}
	for _, receipt := range receiptsResult {
		receipt.Ledger = pbshop.PbLedger(mapLedger[receipt.LedgerId])
	}
	return receiptsResult, nil
}

func listTraders(
	ctx context.Context, shopID int64,
	traderIDs []int64, receiptsResult []*pbshop.Receipt,
) error {
	mapVendor := make(map[int64]*vendoring.ShopVendor)
	var vendorIDs, customerIDs, carrierIDs []int64
	mapCustomer := make(map[int64]*customering.ShopCustomer)
	mapCarrier := make(map[int64]*carrying.ShopCarrier)
	getTradersByIDsQuery := &tradering.ListTradersByIDsQuery{
		ShopID: shopID,
		IDs:    traderIDs,
	}
	if err := traderQuery.Dispatch(ctx, getTradersByIDsQuery); err != nil {
		return err
	}
	for _, trader := range getTradersByIDsQuery.Result.Traders {
		switch trader.Type {
		case tradering.CarrierType:
			carrierIDs = append(carrierIDs, trader.ID)
		case tradering.CustomerType:
			customerIDs = append(customerIDs, trader.ID)
		case tradering.VendorType:
			vendorIDs = append(vendorIDs, trader.ID)
		}
	}
	// Get elements for each of type
	if vendorIDs != nil && len(vendorIDs) > 0 {
		query := &vendoring.ListVendorsByIDsQuery{
			ShopID: shopID,
			IDs:    vendorIDs,
		}
		if err := vendorQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, vendor := range query.Result.Vendors {
			mapVendor[vendor.ID] = vendor
		}
	}
	if customerIDs != nil && len(customerIDs) > 0 {
		query := &customering.ListCustomersByIDsQuery{
			ShopID: shopID,
			IDs:    customerIDs,
		}
		if err := customerQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, customer := range query.Result.Customers {
			mapCustomer[customer.ID] = customer
		}
	}
	if carrierIDs != nil && len(carrierIDs) > 0 {
		query := &carrying.ListCarriersByIDsQuery{
			ShopID: shopID,
			IDs:    carrierIDs,
		}
		if err := carrierQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, carrier := range query.Result.Carriers {
			mapCarrier[carrier.ID] = carrier
		}
	}
	for _, receipt := range receiptsResult {
		if vendor, ok := mapVendor[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       vendor.ID,
				Type:     tradering.VendorType,
				FullName: vendor.FullName,
			}
		}
		if customer, ok := mapCustomer[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       customer.ID,
				Type:     tradering.CustomerType,
				FullName: customer.FullName,
				Phone:    customer.Phone,
			}
		}
		if carrier, ok := mapCarrier[receipt.TraderId]; ok {
			receipt.Partner = &pbshop.Partner{
				Id:       carrier.ID,
				Type:     tradering.CarrierType,
				FullName: carrier.FullName,
			}
		}
	}
	return nil
}
