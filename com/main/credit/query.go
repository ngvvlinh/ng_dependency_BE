package credit

import (
	"context"
	"database/sql"

	"o.o/api/main/credit"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/credit_type"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/credit/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ credit.QueryService = &CreditQueryService{}

type CreditQueryService struct {
	dbMain        *cmsql.Database
	dbEtelecom    *cmsql.Database
	CreditStore   sqlstore.CreditFactory
	eventBus      capi.EventBus
	identityQuery identity.QueryBus
}

func NewQueryCredit(
	bus capi.EventBus,
	dbMain com.MainDB,
	dbEtelecom com.EtelecomDB,
	identityQ identity.QueryBus,
) *CreditQueryService {
	return &CreditQueryService{
		dbMain:        dbMain,
		dbEtelecom:    dbEtelecom,
		identityQuery: identityQ,
		eventBus:      bus,
		CreditStore:   sqlstore.NewCreditStore(dbMain),
	}
}

func CreditQueryServiceMessageBus(q *CreditQueryService) credit.QueryBus {
	b := bus.New()
	return credit.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *CreditQueryService) GetCredit(ctx context.Context, args *credit.GetCreditArgs) (*credit.CreditExtended, error) {
	if args.ID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	query := q.CreditStore(ctx).ID(args.ID)
	if args.ShopID != 0 {
		query = query.ShopID(args.ShopID)
	}
	creditValue, err := query.Get()
	if err != nil {
		return nil, err
	}
	getShopQuery := &identity.GetShopByIDQuery{
		ID: creditValue.ShopID,
	}
	if err = q.identityQuery.Dispatch(ctx, getShopQuery); err != nil {
		return nil, err
	}
	return &credit.CreditExtended{
		Credit: creditValue,
		Shop:   getShopQuery.Result,
	}, nil
}

func (q *CreditQueryService) ListCredits(ctx context.Context, args *credit.ListCreditsArgs) (*credit.ListCreditsResponse, error) {
	var creditsResult []*credit.CreditExtended
	query := q.CreditStore(ctx)
	if args.ShopID != 0 {
		query = query.ShopID(args.ShopID)
	}
	creditValues, err := query.WithPaging(args.Paging).ListCredit()
	if err != nil {
		return nil, err
	}
	if len(creditValues) > 0 {
		var shopIDs []dot.ID
		var mapShop = make(map[dot.ID]*identity.Shop)
		for _, v := range creditValues {
			shopIDs = append(shopIDs, v.ShopID)
		}

		getShopQuery := &identity.ListShopsByIDsQuery{
			IDs: shopIDs,
		}
		if err = q.identityQuery.Dispatch(ctx, getShopQuery); err != nil {
			return nil, err
		}
		for _, v := range getShopQuery.Result {
			mapShop[v.ID] = v
		}
		for _, v := range creditValues {
			creditsResult = append(creditsResult, &credit.CreditExtended{
				Credit: v,
				Shop:   mapShop[v.ShopID],
			})
		}
	}
	return &credit.ListCreditsResponse{
		Credits: creditsResult,
	}, nil
}

func (q *CreditQueryService) GetTelecomUserBalance(ctx context.Context, userID dot.ID) (int, error) {
	if userID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing user ID")
	}
	queryAccounts := &identity.GetAllAccountsByUsersQuery{
		UserIDs: []dot.ID{userID},
		Type:    account_type.Shop.Wrap(),
	}
	if err := q.identityQuery.Dispatch(ctx, queryAccounts); err != nil {
		return 0, err
	}

	shopIDs := []dot.ID{}
	for _, accountUser := range queryAccounts.Result {
		shopIDs = append(shopIDs, accountUser.AccountID)
	}

	// get total user credit
	totalCredit, err := q.CreditStore(ctx).ShopIDs(shopIDs...).Classify(credit_type.CreditClassifyTelecom).SumCredit()
	if err != nil {
		return 0, err
	}

	// get total postage
	totalPostage, err := q.getTotalPostage(ctx, shopIDs)
	if err != nil {
		return 0, err
	}
	res := totalCredit - totalPostage
	return res, nil
}

// calc from call logs
func (q *CreditQueryService) getTotalPostage(ctx context.Context, shopIDs []dot.ID) (int, error) {
	var total sql.NullInt64
	query := q.dbEtelecom.SQL("SELECT SUM(postage) from call_log").In("account_id", shopIDs).
		Where("call_status = 1")
	if err := query.Scan(&total); err != nil {
		return 0, err
	}
	return int(total.Int64), nil
}

func (q *CreditQueryService) GetShippingUserBalance(ctx context.Context, userID dot.ID) (*credit.GetShippingUserBalanceResponse, error) {
	if userID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing user ID")
	}
	queryAccounts := &identity.GetAllAccountsByUsersQuery{
		UserIDs: []dot.ID{userID},
		Type:    account_type.Shop.Wrap(),
	}
	if err := q.identityQuery.Dispatch(ctx, queryAccounts); err != nil {
		return nil, err
	}

	shopIDs := []dot.ID{}
	for _, accountUser := range queryAccounts.Result {
		shopIDs = append(shopIDs, accountUser.AccountID)
	}

	shippingAvailableFee, err := q.calcShippingAvailableFee(shopIDs)
	if err != nil {
		return nil, err
	}
	shippingActualFee, err := q.calcShippingActualFee(shopIDs)
	if err != nil {
		return nil, err
	}

	// get total user credit
	totalCredit, err := q.CreditStore(ctx).ShopIDs(shopIDs...).Classify(credit_type.CreditClassifyShipping).SumCredit()
	if err != nil {
		return nil, err
	}

	res := &credit.GetShippingUserBalanceResponse{
		ShippingActualUserBalance:    totalCredit + shippingActualFee,
		ShippingAvailableUserBalance: totalCredit + shippingAvailableFee,
	}
	return res, nil
}

/*
	calcShippingAvailableFee: tính số dư dự kiến user (không tính những ffm đã thanh toán)

	- Tính theo user (dùng shopIDs của user đó)
	- COD: của tất cả Shop (ffm) khác trạng thái hủy và không phải là đơn trả hàng (status != -1 AND status != 0 AND shipping_status != -2 AND etop_payment_status != 1)
	- Cước phí: đơn có trạng thái khác hủy (chưa đối soát)
*/
func (q *CreditQueryService) calcShippingAvailableFee(shopIDs []dot.ID) (int, error) {
	var totalCODAmount, totalShippingFee sql.NullInt64
	if err := q.dbMain.SQL("SELECT SUM(total_cod_amount) from fulfillment").
		In("shop_id", shopIDs).
		Where("status not in (-1, 0) AND etop_payment_status != 1").
		Where("shipping_status != -2").
		Where("connection_method IS NULL OR connection_method = ?", connection_type.ConnectionMethodBuiltin).
		Scan(&totalCODAmount); err != nil {
		return 0, err
	}

	if err := q.dbMain.SQL("SELECT SUM(shipping_fee_shop) from fulfillment").
		In("shop_id", shopIDs).
		Where("status not in (-1, 0) AND etop_payment_status != 1").
		Where("connection_method IS NULL OR connection_method = ?", connection_type.ConnectionMethodBuiltin).
		Scan(&totalShippingFee); err != nil {
		return 0, err
	}

	res := int(totalCODAmount.Int64 - totalShippingFee.Int64)
	return res, nil
}

/*
	calcShippingActualFee: tính số dư thực tế user (không tính những ffm đã thanh toán)

	- Tính theo user (dùng shopIDs của user đó)
	- COD: của tất cả Shop (ffm), chỉ tính đơn giao thành công và chưa đối soát
	- Cước phí: đơn có trạng thái khác hủy (chưa đối soát)
*/
func (q *CreditQueryService) calcShippingActualFee(shopIDs []dot.ID) (int, error) {
	var totalCODAmount, totalShippingFee sql.NullInt64
	if err := q.dbMain.SQL("SELECT SUM(total_cod_amount) from fulfillment").
		In("shop_id", shopIDs).
		Where("status not in (-1, 0) AND etop_payment_status != 1").
		Where("shipping_status = 1").
		Where("connection_method IS NULL OR connection_method = ?", connection_type.ConnectionMethodBuiltin).
		Scan(&totalCODAmount); err != nil {
		return 0, err
	}

	if err := q.dbMain.SQL("SELECT SUM(shipping_fee_shop) from fulfillment").
		In("shop_id", shopIDs).
		Where("status not in (-1, 0) AND etop_payment_status != 1").
		Where("connection_method IS NULL OR connection_method = ?", connection_type.ConnectionMethodBuiltin).
		Scan(&totalShippingFee); err != nil {
		return 0, err
	}

	res := int(totalCODAmount.Int64 - totalShippingFee.Int64)
	return res, nil
}
