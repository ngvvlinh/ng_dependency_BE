// +build wireinject

package carrier

import (
	"github.com/google/wire"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	"o.o/backend/com/main/shipping/carrier/types"
	"o.o/backend/pkg/common/redis"
	"o.o/capi"
)

func MockManager(
	mockBus capi.Bus,
	eventBus capi.EventBus,
	redisStore redis.Store,
	cfg types.Config,
) (*ShipmentManager, error) {
	panic(wire.Build(
		location.NewQueryBus,
		connectioning.NewQueryBus,
		connectioning.NewCommandBus,
		shipmentservice.NewQueryBus,
		shipmentprice.NewQueryBus,
		pricelistpromotion.NewQueryBus,
		NewShipmentManager,
		connectionmanager.NewConnectionManager,
	))
}
