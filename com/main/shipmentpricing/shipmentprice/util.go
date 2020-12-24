package shipmentprice

import (
	"context"

	"o.o/backend/com/main/shipmentpricing/util"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

const (
	ShipmentPricesRedisKey = "active_shipment_prices"
)

func getActiveShipmentPricesRedisKey(ctx context.Context, shipmentPriceListID dot.ID) string {
	// cache riêng từng wl_partner_id
	// riêng trường hợp wl partner POS, sử dụng chung cache với TopShip (wlPartnerID = 0)
	wlPartnerID := dot.ID(0)
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWLPartnerPOS() {
		wlPartnerID = wlPartner.ID
	}
	key := ShipmentPricesRedisKey +
		":" + util.VersionCaching +
		":wl" + wlPartnerID.String()
	if shipmentPriceListID != 0 {
		key += ":pricelistid" + shipmentPriceListID.String()
	}
	return concatWithEnvKey(key)
}

func DeleteRedisCache(ctx context.Context, redisStore redis.Store, shipmentPriceListID dot.ID) error {
	// key1: bảng giá mặc định default
	key1 := getActiveShipmentPricesRedisKey(ctx, 0)
	if shipmentPriceListID != 0 {
		key2 := getActiveShipmentPricesRedisKey(ctx, shipmentPriceListID)
		if err := redisStore.Del(key2); err != nil {
			return err
		}
	}
	return redisStore.Del(key1)
}

func concatWithEnvKey(key string) string {
	if cmenv.Env() != cmenv.EnvProd {
		key += ":" + cmenv.Env().String()
	}
	return key
}
