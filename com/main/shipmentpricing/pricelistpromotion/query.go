package pricelistpromotion

import (
	"context"
	"sort"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelist"
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
	priceListQS             pricelist.QueryBus
}

func NewQueryService(db com.MainDB, redisStore redis.Store,
	locationQuery location.QueryBus,
	identityQuery identity.QueryBus,
	shopPriceListQS shopshipmentpricelist.QueryBus,
	priceListQS pricelist.QueryBus,
) *QueryService {
	return &QueryService{
		priceListPromotionStore: sqlstore.NewPriceListStorePromotion(db),
		redisStore:              redisStore,
		locationQS:              locationQuery,
		identityQS:              identityQuery,
		shopPriceListQS:         shopPriceListQS,
		priceListQS:             priceListQS,
	}
}

func QueryServiceMessageBus(q *QueryService) pricelistpromotion.QueryBus {
	b := bus.New()
	return pricelistpromotion.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetPriceListPromotion(ctx context.Context, id dot.ID) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	return q.priceListPromotionStore(ctx).ID(id).GetShipmentPriceListPromotion()
}

func (q *QueryService) ListPriceListPromotions(ctx context.Context,
	args *pricelistpromotion.ListPriceListPromotionArgs) ([]*pricelistpromotion.ShipmentPriceListPromotion, error) {
	query := q.priceListPromotionStore(ctx).OptionalConnectionID(args.ConnectionID).OptionalPriceListID(args.PriceListID).WithPaging(args.Paging)
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
	var res []*pricelistpromotion.ShipmentPriceListPromotion

	var fromCustomRegionIDs []dot.ID
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

	// Kiểm tra ngày tạo shop || user
	if !q.checkShopOrUserCreatedAt(ctx, promotion.AppliedRules, args) {
		return false
	}

	// Kiểm tra bảng giá đang sử dụng
	return q.checkUsingPriceList(ctx, promotion, args)
}

func (q *QueryService) checkShopOrUserCreatedAt(ctx context.Context, rules *pricelistpromotion.AppliedRules, args *CheckMatchingPromotionArgs) bool {
	shop := args.Shop
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
	return true
}

// checkUsingPriceList: Kiểm bảng giá đang sử dụng
// - Nếu là bảng giá thường: shop phải có bảng giá đó mới được apply
// - Nếu là bảng giá mặc định: shop có bảng giá mặc định hoặc shop chưa được gán bảng giá đều được apply bảng giá KM này
// Note: Bảng giá khuyến mãi của GHN có thể áp dụng cho các shop đang sử dụng bảng giá của NVC khác
func (q *QueryService) checkUsingPriceList(ctx context.Context, promotion *pricelistpromotion.ShipmentPriceListPromotion, args *CheckMatchingPromotionArgs) bool {
	rules := promotion.AppliedRules
	if len(rules.UsingPriceListIDs) == 0 {
		return true
	}

	shopPriceListQuery := &shopshipmentpricelist.ListShopShipmentPriceListsQuery{
		ShopID: args.Shop.ID,
	}
	if err := q.shopPriceListQS.Dispatch(ctx, shopPriceListQuery); err != nil {
		return false
	}
	shopPriceLists := shopPriceListQuery.Result.ShopShipmentPriceLists
	var mapConnectionShopPriceList = make(map[dot.ID]*shopshipmentpricelist.ShopShipmentPriceList)
	for _, shopPL := range shopPriceLists {
		priceListID := shopPL.ShipmentPriceListID
		mapConnectionShopPriceList[shopPL.ConnectionID] = shopPL
		if cm.IDsContain(rules.UsingPriceListIDs, priceListID) {
			return true
		}
	}

	// Kiểm tra xem bảng giá mặc định có nằm trong UsingPriceListIDs không
	defaultPriceListQuery := &pricelist.ListShipmentPriceListsQuery{
		IsDefault: dot.Bool(true),
		IDs:       rules.UsingPriceListIDs,
	}
	if err := q.priceListQS.Dispatch(ctx, defaultPriceListQuery); err != nil {
		return false
	}
	defaultPriceLists := defaultPriceListQuery.Result
	if len(defaultPriceLists) == 0 {
		return false
	}
	for _, defaultPL := range defaultPriceLists {
		// trường hợp có bảng giá mặc định (tính theo NVC defaultPL.ConnectionID)
		shopPL, ok := mapConnectionShopPriceList[defaultPL.ConnectionID]
		if !ok || shopPL.ShipmentPriceListID == defaultPL.ID {
			// shop chưa được gán bảng giá nào
			// hoặc shop đã được gán vào bảng giá mặc định
			// => apply
			return true
		}
	}
	return false
}
