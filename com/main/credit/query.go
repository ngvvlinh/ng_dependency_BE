package credit

import (
	"context"
	"database/sql"

	"o.o/api/main/credit"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_type"
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
	argsCredit := &credit.GetTotalCreditArgs{
		ShopIDs:  shopIDs,
		Classify: credit_type.CreditClassifyTelecom,
	}
	totalCredit, err := q.getTotalUserCredit(ctx, argsCredit)
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

func (q *CreditQueryService) getTotalUserCredit(ctx context.Context, args *credit.GetTotalCreditArgs) (int, error) {
	var total sql.NullInt64
	query := q.dbMain.SQL("SELECT SUM(amount) from credit").In("shop_id", args.ShopIDs).
		Where("status = 1 AND paid_at IS NOT NULL")
	if args.Classify == credit_type.CreditClassifyShipping {
		query = query.Where("classify IS NULL OR classify = ?", credit_type.CreditClassifyShipping)
	} else {
		query = query.Where("classify = ?", args.Classify)
	}
	if err := query.Scan(&total); err != nil {
		return 0, err
	}
	return int(total.Int64), nil
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
