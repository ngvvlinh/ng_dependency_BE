package shipmentservice

import (
	"context"

	"o.o/backend/com/main/shipmentpricing/util"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

const (
	ShipmentServiceRedisKey = "shipment_service"
)

func getShipmentServiceRedisKey(ctx context.Context, serviceID string, connID dot.ID) string {
	// cache riêng từng wl_partner_id
	// riêng trường hợp wl partner POS, sử dụng chung cache với TopShip
	wlPartnerID := dot.ID(0)
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWLPartnerPOS() {
		wlPartnerID = wlPartner.ID
	}
	return ShipmentServiceRedisKey +
		":" + util.VersionCaching +
		":wl" + wlPartnerID.String() +
		":sid" + serviceID +
		":conn" + connID.String()
}

func DeleteRedisCache(ctx context.Context, redisStore redis.Store, connID dot.ID, serviceIDs []string) error {
	for _, sID := range serviceIDs {
		key := getShipmentServiceRedisKey(ctx, sID, connID)
		if err := redisStore.Del(key); err != nil {
			return err
		}
	}
	return nil
}
