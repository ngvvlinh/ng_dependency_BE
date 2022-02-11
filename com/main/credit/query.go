package credit

import (
	"context"

	"o.o/api/main/credit"
	"o.o/api/main/identity"
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
	if args.Classify.Valid {
		query = query.Classify(args.Classify.Enum)
	}
	if args.DateTo.Before(args.DateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !args.DateFrom.IsZero() {
		query = query.BetweenDateFromAndDateTo(args.DateFrom, args.DateTo)
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
