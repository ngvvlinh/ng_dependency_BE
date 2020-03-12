package shipmentprice

import (
	"context"

	"etop.vn/backend/com/main/shipmentpricing"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/redis"
)

const (
	ShipmentPricesRedisKey = "active_shipment_prices"
)

func getActiveShipmentPricesRedisKey(ctx context.Context) string {
	// cache riêng từng wl_partner_id
	return ShipmentPricesRedisKey +
		":" + shipmentpricing.VersionCaching +
		":wl" + wl.X(ctx).ID.String()
}

func DeleteRedisCache(ctx context.Context, redisStore redis.Store) error {
	key := getActiveShipmentPricesRedisKey(ctx)
	return redisStore.Del(key)
}
