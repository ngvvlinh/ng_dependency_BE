package carrier

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/com/main/shipping/sharemodel"
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
		cmenv.SetEnvironment("test", cmenv.EnvDev.String())
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
				UserID:     "1158799",
				Identifier: "tuan@etop.vn",
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
	shipmentManager, err := MockManager(mockBus, eventBus, nil, cfg)
	if err != nil {
		panic(err)
	}

	Convey("Get Shipment driver", t, func() {
		shipmentType, err := shipmentManager.getShipmentDriver(ctx, connID, 0)
		So(err, ShouldBeNil)
		So(shipmentType, ShouldNotBeNil)
	})
}

func TestFilterShipmentServices(t *testing.T) {
	services1 := []*sharemodel.AvailableShippingService{
		{
			Name:              "ABCDE",
			ServiceFee:        26000,
			ProviderServiceID: "53321",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 37000,
				MakeupFee: 26000,
			},
		},
		{
			Name:              "EFGHI",
			ServiceFee:        26000,
			ProviderServiceID: "53320",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 59000,
				MakeupFee: 26000,
			},
		},
	}
	t.Run("All services are valid", func(t *testing.T) {
		ss := filterShipmentServicesByEdCode(services1)
		assert.EqualValues(t, 1, len(ss))
		assert.EqualValues(t, "ABCDE", ss[0].Name)
		assert.EqualValues(t, "53321", ss[0].ProviderServiceID)
	})

	services2 := []*sharemodel.AvailableShippingService{
		{
			Name:              "ABCDE",
			ServiceFee:        26000,
			ProviderServiceID: "53321",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 37000,
				MakeupFee: 26000,
			},
		},
		{
			Name:              "EFGHI",
			ServiceFee:        26000,
			ProviderServiceID: "53320",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 59000,
				MakeupFee: 26000,
			},
		},
		{
			Name:              "IJKLM",
			ServiceFee:        29000,
			ProviderServiceID: "53320",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 59000,
				MakeupFee: 29000,
			},
		},
	}
	t.Run("1 service don't have code", func(t *testing.T) {
		ss := filterShipmentServicesByEdCode(services2)
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Name < ss[j].Name
		})
		assert.EqualValues(t, 2, len(ss))
		assert.EqualValues(t, "ABCDE", ss[0].Name)
		assert.EqualValues(t, "IJKLM", ss[1].Name)
	})

	services3 := []*sharemodel.AvailableShippingService{
		{
			Name:              "ABCDE",
			ServiceFee:        26000,
			ProviderServiceID: "53321",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: false,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 37000,
				MakeupFee: 26000,
			},
		},
		{
			Name:              "EFGHI",
			ServiceFee:        26000,
			ProviderServiceID: "53320",
			ShipmentServiceInfo: &sharemodel.ShipmentServiceInfo{
				Code:        "GHN.TOPSHIP",
				IsAvailable: true,
			},
			ShipmentPriceInfo: &sharemodel.ShipmentPriceInfo{
				OriginFee: 59000,
				MakeupFee: 26000,
			},
		},
	}
	t.Run("1 services are not valid", func(t *testing.T) {
		ss := filterShipmentServicesByEdCode(services3)
		assert.EqualValues(t, 1, len(ss))
		assert.EqualValues(t, "EFGHI", ss[0].Name)
		assert.EqualValues(t, "53320", ss[0].ProviderServiceID)
	})
}
