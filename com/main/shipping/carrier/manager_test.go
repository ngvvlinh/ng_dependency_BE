package carrier

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/location"
	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/capi/dot"
)

var (
	shipmentManager *ShipmentManager
	locationQS      location.QueryBus
	connID          = dot.ID(123)
	shopID          = dot.ID(456)
)

func init() {

}

func TestShipmentManager(t *testing.T) {
	ctx := bus.Ctx()
	if cmenv.Env() == 0 {
		cmenv.SetEnvironment(cmenv.EnvDev.String())
	}
	mockBus := bus.New()
	mockBus.MockHandler(func(query *connectioning.GetConnectionByIDQuery) error {
		query.Result = &connectioning.Connection{
			ID:                 connID,
			Name:               "topship-ghn",
			Status:             1,
			Driver:             "shipping/shipment/topship/ghn",
			ConnectionType:     connection_type.Shipping,
			ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
			ConnectionMethod:   connection_type.ConnectionMethodTopship,
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
	connectionQS := connectioning.NewQueryBus(mockBus)
	connectionAggr := connectioning.NewCommandBus(mockBus)
	shipmentServiceQS := shipmentservice.NewQueryBus(mockBus)
	shipmentPriceQS := shipmentprice.NewQueryBus(mockBus)
	shipmentManager = NewShipmentManager(locationQS, connectionQS, connectionAggr, nil, shipmentServiceQS, shipmentPriceQS)
	Convey("Get Shipment driver", t, func() {
		shipmentType, err := shipmentManager.getShipmentDriver(ctx, connID, 0)
		So(err, ShouldBeNil)
		So(shipmentType, ShouldNotBeNil)
	})
}
