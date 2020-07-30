package pricelistpromotion

import (
	"context"
	"sort"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var _ pricelistpromotion.QueryService = &QueryService{}

type QueryService struct {
	redisStore              redis.Store
	priceListPromotionStore sqlstore.PriceListStorePromotionFactory
	locationQS              location.QueryBus
	identityQS              identity.QueryBus
	shopPriceListQS         shopshipmentpricelist.QueryBus
}

func NewQueryService(db com.MainDB, redisStore redis.Store, locationQuery location.QueryBus, identityQuery identity.QueryBus, shopPriceListQS shopshipmentpricelist.QueryBus) *QueryService {
	return &QueryService{
		priceListPromotionStore: sqlstore.NewPriceListStorePromotion(db),
		redisStore:              redisStore,
		locationQS:              locationQuery,
		identityQS:              identityQuery,
		shopPriceListQS:         shopPriceListQS,
	}
}

func QueryServiceMessageBus(q *QueryService) pricelistpromotion.QueryBus {
	b := bus.New()
	return pricelistpromotion.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetPriceListPromotion(ctx context.Context, id dot.ID) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	return q.priceListPromotionStore(ctx).ID(id).GetShipmentPriceListPromotion()
}

func (q *QueryService) ListPriceListPromotion(ctx context.Context,
	args *pricelistpromotion.ListPriceListPromotionArgs) ([]*pricelistpromotion.ShipmentPriceListPromotion, error) {
	query := q.priceListPromotionStore(ctx).OptionalConnectionID(args.ConnectionID)
	return query.ListShipmentPriceListPromotions()
}

func (q *QueryService) GetValidPriceListPromotion(ctx context.Context, args *pricelistpromotion.GetValidPriceListPromotionArgs) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID")
	}
	query := q.priceListPromotionStore(ctx).Status(status3.P).ConnectionID(args.ConnectionID).ActiveDate(time.Now())
	promotions, err := query.ListShipmentPriceListPromotions()
	if err != nil {
		return nil, err
	}

	return q.getMatchingPromotion(ctx, promotions, args)
}

func (q *QueryService) getMatchingPromotion(ctx context.Context, promotions []*pricelistpromotion.ShipmentPriceListPromotion, args *pricelistpromotion.GetValidPriceListPromotionArgs) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	var res = []*pricelistpromotion.ShipmentPriceListPromotion{}

	fromCustomRegionIDs := []dot.ID{}
	if args.FromProvinceCode != "" {
		queryFrom := &location.ListCustomRegionsByCodeQuery{
			ProvinceCode: args.FromProvinceCode,
		}
		if err := q.locationQS.Dispatch(ctx, queryFrom); err == nil {
			for _, region := range queryFrom.Result {
				fromCustomRegionIDs = append(fromCustomRegionIDs, region.ID)
			}
		}
	}

	argCheck := &CheckMatchingPromotionArgs{
		FromCustomRegionIDs: fromCustomRegionIDs,
		ConnectionID:        args.ConnectionID,
	}
	if args.ShopID != 0 {
		queryShop := &identity.GetShopByIDQuery{
			ID: args.ShopID,
		}
		if err := q.identityQS.Dispatch(ctx, queryShop); err != nil {
			return nil, err
		}
		argCheck.Shop = queryShop.Result
	}
	for _, promotion := range promotions {
		if q.isMatchingPromotion(ctx, promotion, argCheck) {
			res = append(res, promotion)
		}
	}

	if len(res) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Không có bảng giá khuyến mãi")
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].PriorityPoint > res[j].PriorityPoint
	})
	return res[0], nil
}

type CheckMatchingPromotionArgs struct {
	FromCustomRegionIDs []dot.ID
	Shop                *identity.Shop
	ConnectionID        dot.ID
}

func (q *QueryService) isMatchingPromotion(ctx context.Context, promotion *pricelistpromotion.ShipmentPriceListPromotion, args *CheckMatchingPromotionArgs) bool {
	if promotion.Status != status3.P {
		return false
	}
	if promotion.AppliedRules == nil {
		return false
	}
	if promotion.ConnectionID != args.ConnectionID {
		return false
	}
	rules := promotion.AppliedRules

	// Kiểm tra vùng tự định nghĩa điểm lấy hàng
	if len(rules.FromCustomRegionIDs) > 0 {
		checkRegion := false
		for _, fromRegionID := range args.FromCustomRegionIDs {
			if cm.IDsContain(rules.FromCustomRegionIDs, fromRegionID) {
				checkRegion = true
				break
			}
		}
		if !checkRegion {
			return false
		}
	}

	if args.Shop == nil {
		return true
	}
	shop := args.Shop
	// Kiểm tra ngày tạo shop || user
	if !rules.ShopCreatedDate.From.IsZero() {
		fromDate := rules.ShopCreatedDate.From.ToTime()
		toDate := rules.ShopCreatedDate.To.ToTime()

		if shop.CreatedAt.Before(fromDate) {
			return false
		}
		if shop.CreatedAt.After(toDate) {
			return false
		}
	}

	if !rules.UserCreatedDate.From.IsZero() {
		queryUser := &identity.GetUserByIDQuery{
			UserID: shop.OwnerID,
		}
		if err := q.identityQS.Dispatch(ctx, queryUser); err != nil {
			return false
		}
		user := queryUser.Result
		fromDate := rules.UserCreatedDate.From.ToTime()
		toDate := rules.UserCreatedDate.To.ToTime()

		if user.CreatedAt.Before(fromDate) {
			return false
		}
		if user.CreatedAt.After(toDate) {
			return false
		}
	}

	// check price list
	if len(rules.UsingPriceListIDs) == 0 {
		return true
	}
	queryShopPriceList := &shopshipmentpricelist.ListShopShipmentPriceListsQuery{
		ShopID: args.Shop.ID,
	}
	if err := q.shopPriceListQS.Dispatch(ctx, queryShopPriceList); err != nil {
		return false
	}
	for _, shopPL := range queryShopPriceList.Result.ShopShipmentPriceLists {
		priceListID := shopPL.ShipmentPriceListID
		if cm.IDsContain(rules.UsingPriceListIDs, priceListID) {
			return true
		}
	}
	return false
}
