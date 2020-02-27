package aggregate

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/backend/com/main/connectioning/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	. "etop.vn/backend/pkg/common/testing"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll     = l.New()
	db     *cmsql.Database
	connID = dot.ID(123)
	shopID = dot.ID(123456)
)

func init() {
	postgres := cc.DefaultPostgres()
	db = cmsql.MustConnect(postgres)
	sqlstore.Init(db)
	db.MustExec(`
		DROP TABLE IF EXISTS shop_connection, connection CASCADE;
		CREATE TABLE connection (
			id INT8 PRIMARY KEY
			, name TEXT
			, status INT2
			, partner_id INT8
			, driver_config JSONB
			, driver TEXT
			, connection_type TEXT
			, connection_subtype TEXT
			, connection_method TEXT
			, connection_provider TEXT
			, created_at TIMESTAMP WITH TIME ZONE
			, updated_at TIMESTAMP WITH TIME ZONE
			, deleted_at TIMESTAMP WITH TIME ZONE
			, etop_affiliate_account JSONB
			, code TEXT
			, image_url TEXT
		);
		CREATE TABLE shop_connection (
			shop_id INT8
			, connection_id INT8
			, token TEXT
			, token_expires_at TIMESTAMPTZ
			, status INT2
			, is_global BOOLEAN
			, connection_states JSONB
			, created_at TIMESTAMPTZ
			, updated_at TIMESTAMPTZ
			, deleted_at TIMESTAMPTZ
			, external_data JSONB
		);
	`)
}

func TestConnectionAggregate(t *testing.T) {

	Convey("Connection Aggregate", t, func() {
		_conn := &model.Connection{
			ID:     connID,
			Name:   "Connection",
			Status: 0,
		}

		ctx := bus.Ctx()
		Aggr := NewConnectionAggregate(db, bus.New()).MessageBus()
		_, err := db.Insert(_conn)
		So(err, ShouldBeNil)

		Reset(func() {
			db.MustExec("truncate connection CASCADE")
		})

		Convey("Create Connection Success", func() {
			cmd := &connectioning.CreateConnectionCommand{
				Name:               "test create",
				Driver:             "shipping/shipment/topship/ghn",
				ConnectionType:     connection_type.Shipping,
				ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
				ConnectionMethod:   connection_type.ConnectionMethodTopship,
				ConnectionProvider: connection_type.ConnectionProviderGHN,
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			conn := cmd.Result
			So(conn.Name, ShouldEqual, cmd.Name)
		})

		Convey("Update Driver Config", func() {
			cmd := &connectioning.UpdateConnectionCommand{
				ID: connID,
				DriverConfig: &connectioning.ConnectionDriverConfig{
					CreateFulfillmentURL:   "http://create-fulfillment",
					GetFulfillmentURL:      "http://get-fulfillment",
					GetShippingServicesURL: "http://get-shipping-services",
					CancelFulfillmentURL:   "http://cancel-fulfillment",
				},
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			conn := cmd.Result
			So(conn.ID, ShouldEqual, cmd.ID)
			So(conn.DriverConfig, ShouldDeepEqual, cmd.DriverConfig)
		})

		Convey("Confirm Missing Connection ID", func() {
			cmd := &connectioning.ConfirmConnectionCommand{}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldCMError, cm.InvalidArgument, "missing id")
		})

		Convey("Confirm Success", func() {
			cmd := &connectioning.ConfirmConnectionCommand{
				ID: connID,
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result, ShouldEqual, 1)

			Convey("Confirm Fail Precondition: Status = 1", func() {
				cmd := &connectioning.ConfirmConnectionCommand{
					ID: connID,
				}
				err := Aggr.Dispatch(ctx, cmd)
				So(err, ShouldCMError, cm.FailedPrecondition, "Can not confirm this Connection")
			})
		})

		Convey("Delete Success", func() {
			cmd := &connectioning.DeleteConnectionCommand{
				ID: connID,
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result, ShouldEqual, 1)
		})
	})
}

func TestShopConnectionAggregate(t *testing.T) {
	Convey("Connection Aggregate", t, func() {
		_shopConn := &model.ShopConnection{
			ShopID:       shopID,
			ConnectionID: connID,
			Token:        "token",
			Status:       0,
		}

		ctx := context.Background()
		Aggr := NewConnectionAggregate(db, bus.New()).MessageBus()
		_, err := db.Insert(_shopConn)
		So(err, ShouldBeNil)

		Reset(func() {
			db.MustExec("truncate connection, shop_connection CASCADE")
		})

		Convey("Create Success", func() {
			cmd := &connectioning.CreateShopConnectionCommand{
				ShopID:       shopID,
				ConnectionID: dot.ID(124),
				Token:        "token",
			}

			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			conn := cmd.Result
			So(conn.ConnectionID, ShouldEqual, cmd.ConnectionID)
		})

		Convey("Update Token", func() {
			cmd := &connectioning.UpdateShopConnectionTokenCommand{
				ShopID:       shopID,
				ConnectionID: connID,
				Token:        "token update",
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			conn := cmd.Result
			So(conn.ShopID, ShouldEqual, shopID)
			So(conn.Token, ShouldEqual, cmd.Token)
		})

		Convey("Confirm Missing Connection ID", func() {
			cmd := &connectioning.ConfirmShopConnectionCommand{}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldCMError, cm.InvalidArgument, "missing shop_id")
		})

		Convey("Confirm Success", func() {
			cmd := &connectioning.ConfirmShopConnectionCommand{
				ShopID:       shopID,
				ConnectionID: connID,
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result, ShouldEqual, 1)

			Convey("Confirm Fail Precondition: Status = 1", func() {
				cmd := &connectioning.ConfirmShopConnectionCommand{
					ShopID:       shopID,
					ConnectionID: connID,
				}
				err := Aggr.Dispatch(ctx, cmd)
				So(err, ShouldCMError, cm.FailedPrecondition, "Can not confirm this Connection")
			})
		})

		Convey("Delete Success", func() {
			cmd := &connectioning.DeleteShopConnectionCommand{
				ShopID:       shopID,
				ConnectionID: connID,
			}
			err := Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result, ShouldEqual, 1)
		})
	})
}
