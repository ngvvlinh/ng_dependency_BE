package shipmentprice

import (
	"context"

	"o.o/backend/com/main/shipmentpricing/util"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/redis"
)

const (
	ShipmentPricesRedisKey = "active_shipment_prices"
)

func getActiveShipmentPricesRedisKey(ctx context.Context) string {
	// cache riêng từng wl_partner_id
	return ShipmentPricesRedisKey +
		":" + util.VersionCaching +
		":wl" + wl.X(ctx).ID.String()
}

func DeleteRedisCache(ctx context.Context, redisStore redis.Store) error {
	key := getActiveShipmentPricesRedisKey(ctx)
	return redisStore.Del(key)
}
