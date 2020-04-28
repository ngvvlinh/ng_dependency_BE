package shipmentservice

import (
	"context"

	"o.o/backend/com/main/shipmentpricing"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

const (
	ShipmentServiceRedisKey = "shipment_service"
)

func getShipmentServiceRedisKey(ctx context.Context, serviceID string, connID dot.ID) string {
	return ShipmentServiceRedisKey +
		":" + shipmentpricing.VersionCaching +
		":wl" + wl.X(ctx).ID.String() +
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
