package api

import (
	"context"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/catalog"
	servicecatalog "etop.vn/api/main/catalog"
	"etop.vn/api/services/affiliate"
	pbcm "etop.vn/backend/pb/common"
	pbaff "etop.vn/backend/pb/services/affiliate"
	pbshop "etop.vn/backend/pkg/etop/api/shop"
	modeletop "etop.vn/backend/pkg/etop/model"
	wrapaff "etop.vn/backend/wrapper/services/affiliate"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

func init() {
	bus.AddHandlers("",
		GetCommissions,
		NotifyNewShopPurchase,
		CreateOrUpdateTradingCommissionSetting,
		CreateOrUpdateAffiliateCommissionSetting,
		TradingGetProducts,
		AffiliateGetProducts,
	)
}

var ll = l.New()

var (
	catalogQuery   catalog.QueryBus
	affiliateCmd   affiliate.CommandBus
	affiliateQuery affiliate.QueryBus
)

func Init(
	affCmd affiliate.CommandBus,
	affQuery affiliate.QueryBus,
	catQuery catalog.QueryBus,
) {
	affiliateCmd = affCmd
	catalogQuery = catQuery
	affiliateQuery = affQuery
}

func GetCommissions(ctx context.Context, q *wrapaff.GetCommissionsEndpoint) error {
	q.Result = &pbaff.GetCommissionsResponse{Message: "hello"}
	return nil
}

func NotifyNewShopPurchase(ctx context.Context, q *wrapaff.NotifyNewShopPurchaseEndpoint) error {
	return nil
}

func CreateOrUpdateTradingCommissionSetting(ctx context.Context, q *wrapaff.CreateOrUpdateTradingCommissionSettingEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.Unauthenticated, nil, "Unauthorized")
	}

	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: modeletop.EtopTradingAccountID, // TODO test public api
		Amount:    q.Amount,
		Unit:      *q.Unit,
		Type:      "shop",
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbCommissionSetting(cmd.Result)
	return nil
}

func CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, q *wrapaff.CreateOrUpdateAffiliateCommissionSettingEndpoint) error {
	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: q.Context.Affiliate.ID,
		Amount:    q.Amount,
		Unit:      *q.Unit,
		Type:      "affiliate",
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbCommissionSetting(cmd.Result)
	return nil
}

func TradingGetProducts(ctx context.Context, q *wrapaff.TradingGetProductsEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.Unauthenticated, nil, "Unauthenticated")
	}
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	shopCommissionSettingMap := GetShopCommissionSettingsByProducts(ctx, modeletop.EtopTradingAccountID, query.Result.Products)
	var products []*pbaff.ShopProductResponse
	for _, product := range query.Result.Products {
		shopCommissionSetting := shopCommissionSettingMap[product.ProductID]
		var pbShopCommissionSetting *pbaff.CommissionSetting = nil
		if shopCommissionSetting != nil {
			pbShopCommissionSetting = pbaff.PbCommissionSetting(shopCommissionSetting)
		}
		products = append(products, &pbaff.ShopProductResponse{
			Product:               pbshop.PbShopProduct(product),
			ShopCommissionSetting: pbShopCommissionSetting,
		})
	}

	q.Result = &pbaff.ShopGetProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: products,
	}
	return nil
}

func GetShopCommissionSettingsByProducts(ctx context.Context, accountID int64, products []*servicecatalog.ShopProduct) map[int64]*affiliate.CommissionSetting {
	var productIds []int64
	for _, product := range products {
		productIds = append(productIds, product.ProductID)
	}

	getShopCommissionByProductIDsQuery := &affiliate.GetCommissionByProductIDsQuery{
		AccountID:  accountID,
		ProductIDs: productIds,
	}
	if err := affiliateQuery.Dispatch(ctx, getShopCommissionByProductIDsQuery); err != nil {
		return map[int64]*affiliate.CommissionSetting{}
	}

	var interfaceArr []interface{}
	for _, e := range getShopCommissionByProductIDsQuery.Result {
		interfaceArr = append(interfaceArr, e)
	}
	shopCommissionMap := map[int64]*affiliate.CommissionSetting{}
	for _, e := range getShopCommissionByProductIDsQuery.Result {
		shopCommissionMap[e.ProductID] = e
	}

	return shopCommissionMap
}

func AffiliateGetProducts(ctx context.Context, q *wrapaff.AffiliateGetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	tradingCommissionMap := GetShopCommissionSettingsByProducts(ctx, modeletop.EtopTradingAccountID, query.Result.Products)
	affCommissionMap := GetShopCommissionSettingsByProducts(ctx, q.Context.Affiliate.ID, query.Result.Products)
	var products []*pbaff.AffiliateProductResponse
	for _, product := range query.Result.Products {
		tradingCommissionSetting := tradingCommissionMap[product.ProductID]
		affCommissionSetting := affCommissionMap[product.ProductID]
		var pbTradingCommissionSetting *pbaff.CommissionSetting = nil
		if tradingCommissionSetting != nil {
			pbTradingCommissionSetting = pbaff.PbCommissionSetting(tradingCommissionSetting)
		}
		var pbAffCommissionSetting *pbaff.CommissionSetting = nil
		if affCommissionSetting != nil {
			pbAffCommissionSetting = pbaff.PbCommissionSetting(affCommissionSetting)
		}

		products = append(products, &pbaff.AffiliateProductResponse{
			Product:                    pbshop.PbShopProduct(product),
			ShopCommissionSetting:      pbTradingCommissionSetting,
			AffiliateCommissionSetting: pbAffCommissionSetting,
		})
	}

	q.Result = &pbaff.AffiliateGetProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: products,
	}

	return nil
}
