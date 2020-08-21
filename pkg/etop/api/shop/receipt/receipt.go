package receipt

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/ledgering"
	"o.o/api/main/receipting"
	"o.o/api/shopping/carrying"
	"o.o/api/shopping/customering"
	"o.o/api/shopping/suppliering"
	"o.o/api/shopping/tradering"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	shop2 "o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
	"o.o/capi/util"
)

type ReceiptService struct {
	session.Session

	CarrierQuery  carrying.QueryBus
	CustomerQuery customering.QueryBus
	LedgerQuery   ledgering.QueryBus
	ReceiptAggr   receipting.CommandBus
	ReceiptQuery  receipting.QueryBus
	SupplierQuery suppliering.QueryBus
	TraderQuery   tradering.QueryBus
}

func (s *ReceiptService) Clone() api.ReceiptService { res := *s; return &res }

func (s *ReceiptService) CreateReceipt(ctx context.Context, q *api.CreateReceiptRequest) (_ *api.Receipt, _err error) {
	key := fmt.Sprintf("Create receipt %v-%v-%v-%v-%v-%v-%v-%v",
		s.SS.Shop().ID, s.SS.Claim().UserID, q.TraderId, q.LedgerId, q.Title, q.Description, q.Amount, q.Type)
	resp, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 15*time.Second, "Create receipt",
		func() (interface{}, error) { return s.createReceipt(ctx, q) })
	if err != nil {
		return nil, err
	}
	result := convertpb.PbReceipt(resp.(*receipting.Receipt))
	return result, nil
}

func (s *ReceiptService) createReceipt(ctx context.Context, q *api.CreateReceiptRequest) (_ *receipting.Receipt, err error) {
	cmd := &receipting.CreateReceiptCommand{
		ShopID:      s.SS.Shop().ID,
		CreatedBy:   s.SS.Claim().UserID,
		TraderID:    q.TraderId,
		Title:       q.Title,
		Description: q.Description,
		Amount:      q.Amount,
		LedgerID:    q.LedgerId,
		RefType:     q.RefType,
		Type:        q.Type,
		Mode:        receipt_mode.Manual,
		Status:      int(status3.Z),
		Lines:       convertpb.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
		PaidAt:      q.PaidAt.ToTime(),
		Note:        q.Note,
	}
	if err := s.ReceiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

func (s *ReceiptService) UpdateReceipt(ctx context.Context, q *api.UpdateReceiptRequest) (_ *api.Receipt, err error) {
	cmd := &receipting.UpdateReceiptCommand{
		ID:          q.Id,
		ShopID:      s.SS.Shop().ID,
		TraderID:    q.TraderId,
		Title:       q.Title,
		Description: q.Description,
		Amount:      q.Amount,
		LedgerID:    q.LedgerId,
		RefIDs:      nil,
		RefType:     q.RefType,
		Lines:       convertpb.Convert_api_ReceiptLines_To_core_ReceiptLines(q.Lines),
		Trader:      nil,
		PaidAt:      q.PaidAt.ToTime(),
		Note:        q.Note,
	}
	err = s.ReceiptAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	result := convertpb.PbReceipt(cmd.Result)
	return result, nil
}

func (s *ReceiptService) ConfirmReceipt(ctx context.Context, q *pbcm.IDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &receipting.ConfirmReceiptCommand{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.ReceiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}

	return result, nil
}

func (s *ReceiptService) CancelReceipt(ctx context.Context, q *api.CancelReceiptRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &receipting.CancelReceiptCommand{
		ID:           q.Id,
		ShopID:       s.SS.Shop().ID,
		CancelReason: util.CoalesceString(q.CancelReason, q.Reason),
	}
	if err := s.ReceiptAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}

	return result, nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, q *pbcm.IDRequest) (*api.Receipt, error) {
	// Check receipt is exist
	getReceiptQuery := &receipting.GetReceiptByIDQuery{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.ReceiptQuery.Dispatch(ctx, getReceiptQuery); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy phiếu").
			Throw()
	}
	receipts, err := s.getInfosForReceipts(ctx, s.SS.Shop().ID, []*receipting.Receipt{getReceiptQuery.Result})
	if err != nil {
		return nil, err
	}
	result := receipts[0]
	return result, nil
}

func (s *ReceiptService) GetReceipts(ctx context.Context, q *api.GetReceiptsRequest) (*api.ReceiptsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &receipting.ListReceiptsQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.ReceiptQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	receipts, err := s.getInfosForReceipts(ctx, s.SS.Shop().ID, query.Result.Receipts)
	if err != nil {
		return nil, err
	}
	result := &api.ReceiptsResponse{
		TotalAmountConfirmedReceipt: query.Result.TotalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: query.Result.TotalAmountConfirmedPayment,
		Receipts:                    receipts,
		Paging:                      cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *ReceiptService) GetReceiptsByLedgerType(ctx context.Context, q *api.GetReceiptsByLedgerTypeRequest) (*api.ReceiptsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	listLedgersByType := &ledgering.ListLedgersByTypeQuery{
		LedgerType: q.Type,
		ShopID:     s.SS.Shop().ID,
	}
	if err := s.LedgerQuery.Dispatch(ctx, listLedgersByType); err != nil {
		return nil, err
	}

	var ledgerIDs []dot.ID
	for _, ledger := range listLedgersByType.Result.Ledgers {
		ledgerIDs = append(ledgerIDs, ledger.ID)
	}

	query := &receipting.ListReceiptsByLedgerIDsQuery{
		ShopID:    s.SS.Shop().ID,
		LedgerIDs: ledgerIDs,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(q.Filters),
	}
	if err := s.ReceiptQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	receipts, err := s.getInfosForReceipts(ctx, s.SS.Shop().ID, query.Result.Receipts)
	if err != nil {
		return nil, err
	}
	result := &api.ReceiptsResponse{
		TotalAmountConfirmedReceipt: query.Result.TotalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: query.Result.TotalAmountConfirmedPayment,
		Receipts:                    receipts,
		Paging:                      cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *ReceiptService) getInfosForReceipts(ctx context.Context, shopID dot.ID, receipts []*receipting.Receipt) (receiptsResult []*shop.Receipt, _ error) {
	mapOrderIDAndReceivedAmount := make(map[dot.ID]int)
	mapLedger := make(map[dot.ID]*ledgering.ShopLedger)
	var refIDs, userIDs, traderIDs, ledgerIDs []dot.ID

	receiptsResult = convertpb.PbReceipts(receipts)

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
	getUsersOfCurrAccount := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs:     []dot.ID{shopID},
		IncludeDeleted: true,
	}
	if err := bus.Dispatch(ctx, getUsersOfCurrAccount); err != nil {
		return nil, err
	}
	mapUserIDAndUser := make(map[dot.ID]*identitymodel.User)
	for _, accountUser := range getUsersOfCurrAccount.Result.AccountUsers {
		mapUserIDAndUser[accountUser.User.ID] = accountUser.User
	}
	for _, receipt := range receiptsResult {
		if receipt.CreatedBy != 0 {
			receipt.User = convertpb.PbUser(mapUserIDAndUser[receipt.CreatedBy])
		}
	}

	// List traders
	if err := s.listTraders(ctx, shopID, traderIDs, receiptsResult); err != nil {
		return nil, err
	}

	// List ledgers
	getLedgersByIDs := &ledgering.ListLedgersByIDsQuery{
		ShopID: shopID,
		IDs:    ledgerIDs,
	}
	if err := s.LedgerQuery.Dispatch(ctx, getLedgersByIDs); err != nil {
		return nil, err
	}
	for _, ledger := range getLedgersByIDs.Result.Ledgers {
		mapLedger[ledger.ID] = ledger
	}
	for _, receipt := range receiptsResult {
		receipt.Ledger = convertpb.PbLedger(mapLedger[receipt.LedgerId])
	}
	return receiptsResult, nil
}

func (s *ReceiptService) listTraders(
	ctx context.Context, shopID dot.ID,
	traderIDs []dot.ID, receiptsResult []*shop.Receipt,
) error {
	mapSupplier := make(map[dot.ID]*suppliering.ShopSupplier)
	mapCustomer := make(map[dot.ID]*customering.ShopCustomer)
	mapCarrier := make(map[dot.ID]*carrying.ShopCarrier)
	var supplierIDs, customerIDs, carrierIDs []dot.ID
	mapTraderID := make(map[dot.ID]bool)
	for _, traderID := range traderIDs {
		if traderID == idutil.TopShipID {
			carrierIDs = append(carrierIDs, traderID)
		}
	}
	getTradersByIDsQuery := &tradering.ListTradersByIDsQuery{
		ShopID: shopID,
		IDs:    traderIDs,
	}
	if err := s.TraderQuery.Dispatch(ctx, getTradersByIDsQuery); err != nil {
		return err
	}

	for _, trader := range getTradersByIDsQuery.Result.Traders {
		switch trader.Type {
		case tradering.CarrierType:
			carrierIDs = append(carrierIDs, trader.ID)
		case tradering.CustomerType:
			customerIDs = append(customerIDs, trader.ID)
		case tradering.SupplierType:
			supplierIDs = append(supplierIDs, trader.ID)
		}
	}
	for _, traderID := range traderIDs {
		if traderID == customering.CustomerAnonymous {
			customerIDs = append(customerIDs, traderID)
		}
	}
	// Get elements for each of type
	if supplierIDs != nil && len(supplierIDs) > 0 {
		query := &suppliering.ListSuppliersByIDsQuery{
			ShopID: shopID,
			IDs:    supplierIDs,
		}
		if err := s.SupplierQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, supplier := range query.Result.Suppliers {
			mapSupplier[supplier.ID] = supplier
			mapTraderID[supplier.ID] = true
		}
	}
	if customerIDs != nil && len(customerIDs) > 0 {
		query := &customering.ListCustomersByIDsQuery{
			ShopID: shopID,
			IDs:    customerIDs,
		}
		if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, customer := range query.Result.Customers {
			mapCustomer[customer.ID] = customer
			mapTraderID[customer.ID] = true
		}

		getIndependentCustomerQuery := &customering.GetCustomerIndependentQuery{}
		if err := s.CustomerQuery.Dispatch(ctx, getIndependentCustomerQuery); err != nil {
			return err
		}
		anonymousCustomer := getIndependentCustomerQuery.Result
		mapCustomer[anonymousCustomer.ID] = getIndependentCustomerQuery.Result
		mapTraderID[anonymousCustomer.ID] = true
	}
	if carrierIDs != nil && len(carrierIDs) > 0 {
		query := &carrying.ListCarriersByIDsQuery{
			ShopID: shopID,
			IDs:    carrierIDs,
		}
		if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
			return err
		}
		for _, carrier := range query.Result.Carriers {
			mapCarrier[carrier.ID] = carrier
			mapTraderID[carrier.ID] = true
		}

		var hasTopShipID bool
		for _, carrierID := range carrierIDs {
			if carrierID == idutil.TopShipID {
				hasTopShipID = true
				break
			}
		}
		if hasTopShipID {
			getTopShipCarrierQuery := &carrying.GetCarrierByIDQuery{
				ID: idutil.TopShipID,
			}
			if err := s.CarrierQuery.Dispatch(ctx, getTopShipCarrierQuery); err != nil {
				return err
			}
			mapCarrier[idutil.TopShipID] = getTopShipCarrierQuery.Result
			mapTraderID[idutil.TopShipID] = true
		}
	}
	for _, receipt := range receiptsResult {
		if _, ok := mapTraderID[receipt.TraderId]; !ok && receipt.Trader != nil {
			receipt.Trader.Deleted = true
		}
	}
	return nil
}
