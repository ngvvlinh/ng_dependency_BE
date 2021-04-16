package query

import (
	"context"
	"database/sql"

	"o.o/api/main/identity"
	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/service_classify"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/transaction/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ transaction.QueryService = &QueryService{}

type QueryService struct {
	dbMain        *cmsql.Database
	dbEtelecom    *cmsql.Database
	store         sqlstore.TransactionStoreFactory
	identityQuery identity.QueryBus
}

func NewQueryService(
	db com.MainDB,
	dbEtelecom com.EtelecomDB,
	identityQ identity.QueryBus,
) *QueryService {
	return &QueryService{
		dbMain:        db,
		dbEtelecom:    dbEtelecom,
		store:         sqlstore.NewTransactionStore(db),
		identityQuery: identityQ,
	}
}

func QueryServiceMessageBus(q *QueryService) transaction.QueryBus {
	b := bus.New()
	return transaction.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetTransactionByID(ctx context.Context, tranID dot.ID, userID dot.ID) (*transaction.Transaction, error) {
	return q.store(ctx).ID(tranID).AccountID(userID).GetTransaction()
}

func (q *QueryService) ListTransactions(ctx context.Context, args *transaction.GetTransactionsArgs) (*transaction.TransactionResponse, error) {
	query := q.store(ctx).OptionalAccountID(args.AccountID)

	if args.DateTo.Before(args.DateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !args.DateFrom.IsZero() {
		query = query.BetweenDateFromAndDateTo(args.DateFrom, args.DateTo)
	}

	if args.RefID != 0 {
		query = query.ReferralID(args.RefID)
	}
	if args.RefType != 0 {
		query = query.ReferralType(args.RefType)
	}

	transactions, err := query.WithPaging(args.Paging).ListTransactions()
	if err != nil {
		return nil, err
	}
	return &transaction.TransactionResponse{
		Paging:       query.GetPaging(),
		Transactions: transactions,
	}, nil
}

func (q *QueryService) GetTransactionByReferral(ctx context.Context, args *transaction.GetTrxnByReferralArgs) (*transaction.Transaction, error) {
	return q.store(ctx).ReferralType(args.ReferralType).ReferralID(args.ReferralID).GetTransaction()
}

func (q *QueryService) GetBalanceUser(ctx context.Context, args *transaction.GetBalanceUserArgs) (*transaction.GetBalanceUserResponse, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}
	switch args.Classify {
	case service_classify.Shipping:
		return q.GetShippingUserBalance(ctx, args.UserID)
	case service_classify.Telecom:
		return q.GetTelecomUserBalance(ctx, args.UserID)
	case service_classify.All:
		shippingBalance, err := q.GetShippingUserBalance(ctx, args.UserID)
		if err != nil {
			return nil, err
		}
		telecomBalance, err := q.GetTelecomUserBalance(ctx, args.UserID)
		if err != nil {
			return nil, err
		}
		return &transaction.GetBalanceUserResponse{
			AvailableBalance: shippingBalance.AvailableBalance + telecomBalance.AvailableBalance,
			ActualBalance:    shippingBalance.ActualBalance + telecomBalance.ActualBalance,
			Classify:         service_classify.All,
		}, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Service classify does not valid")
	}
}

func (q *QueryService) GetTelecomUserBalance(ctx context.Context, userID dot.ID) (*transaction.GetBalanceUserResponse, error) {
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

	// get total user transaction
	totalCredit, err := q.store(ctx).AccountIDs(shopIDs...).Classify(service_classify.Telecom).GetBalance()
	if err != nil {
		return nil, err
	}

	// get total postage
	totalPostage, err := q.getTotalPostage(ctx, shopIDs)
	if err != nil {
		return nil, err
	}
	res := totalCredit - totalPostage
	return &transaction.GetBalanceUserResponse{
		AvailableBalance: res,
		ActualBalance:    res,
		Classify:         service_classify.Telecom,
	}, nil
}

// calc from call logs
func (q *QueryService) getTotalPostage(ctx context.Context, shopIDs []dot.ID) (int, error) {
	var total sql.NullInt64
	query := q.dbEtelecom.SQL("SELECT SUM(postage) from call_log").In("account_id", shopIDs).
		Where("call_status = 1")
	if err := query.Scan(&total); err != nil {
		return 0, err
	}
	return int(total.Int64), nil
}

func (q *QueryService) GetShippingUserBalance(ctx context.Context, userID dot.ID) (*transaction.GetBalanceUserResponse, error) {
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
	totalCredit, err := q.store(ctx).AccountIDs(shopIDs...).Classify(service_classify.Shipping).GetBalance()
	if err != nil {
		return nil, err
	}

	return &transaction.GetBalanceUserResponse{
		AvailableBalance: totalCredit + shippingAvailableFee,
		ActualBalance:    totalCredit + shippingActualFee,
		Classify:         service_classify.Shipping,
	}, nil
}

/*
	calcShippingAvailableFee: tính số dư dự kiến user (không tính những ffm đã thanh toán)

	- Tính theo user (dùng shopIDs của user đó)
	- COD: của tất cả Shop (ffm) khác trạng thái hủy và không phải là đơn trả hàng (status != -1 AND status != 0 AND shipping_status != -2 AND etop_payment_status != 1)
	- Cước phí: đơn có trạng thái khác hủy (chưa đối soát)
*/
func (q *QueryService) calcShippingAvailableFee(shopIDs []dot.ID) (int, error) {
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
func (q *QueryService) calcShippingActualFee(shopIDs []dot.ID) (int, error) {
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
