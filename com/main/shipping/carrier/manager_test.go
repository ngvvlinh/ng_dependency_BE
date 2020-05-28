package carrier

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi/dot"
)

var (
	connID = dot.ID(123)
	shopID = dot.ID(456)
)

func TestShipmentManager(t *testing.T) {
	ctx := bus.Ctx()
	if cmenv.Env() == 0 {
		cmenv.SetEnvironment(cmenv.EnvDev.String())
	}
	mockBus, eventBus := bus.New(), bus.New()
	mockBus.MockHandler(func(query *connectioning.GetConnectionByIDQuery) error {
		query.Result = &connectioning.Connection{
			ID:                 connID,
			Name:               "topship-ghn",
			Status:             1,
			Driver:             "shipping/shipment/builtin/ghn",
			ConnectionType:     connection_type.Shipping,
			ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
			ConnectionMethod:   connection_type.ConnectionMethodBuiltin,
			ConnectionProvider: connection_type.ConnectionProviderGHN,
		}
		return nil
	})
	mockBus.MockHandler(func(query *connectioning.GetShopConnectionByIDQuery) error {
		query.Result = &connectioning.ShopConnection{
			ShopID:       shopID,
			ConnectionID: connID,
			Token:        "token",
			Status:       1,
			ExternalData: &connectioning.ShopConnectionExternalData{
				UserID: "1158799",
				Email:  "tuan@etop.vn",
			},
		}
		return nil
	})
	cfg := Config{
		Endpoints: []ConfigEndpoint{
			{
				connection_type.ConnectionProviderGHN,
				"/callback/ghn",
			},
		},
	}
	shipmentManager, err := MockManager(mockBus, eventBus, nil, false, cfg)
	if err != nil {
		panic(err)
	}

	Convey("Get Shipment driver", t, func() {
		shipmentType, err := shipmentManager.getShipmentDriver(ctx, connID, 0)
		So(err, ShouldBeNil)
		So(shipmentType, ShouldNotBeNil)
	})
}
