package query

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/api/main/connectioning"
	"o.o/backend/com/main/connectioning/model"
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var (
	db     *cmsql.Database
	connID = dot.ID(1234)
	shopID = dot.ID(4567)
	QS     connectioning.QueryBus
	ctx    context.Context
)

func init() {
	postgres := cc.DefaultPostgres()
	db = cmsql.MustConnect(postgres)
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
			, services JSON
			, wl_partner_id INT8
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
	_conn := &model.Connection{
		ID:          connID,
		Name:        "Connection",
		Status:      1,
		WLPartnerID: drivers.ITopXID,
	}
	_shopConn := &model.ShopConnection{
		ShopID:       shopID,
		ConnectionID: connID,
		Token:        "token",
		Status:       1,
	}
	wl.Init(cmenv.EnvDev)
	ctx = wl.WrapContext(bus.Ctx(), drivers.ITopXID)
	ctx = bus.NewRootContext(ctx)
	QS = NewConnectionQuery(db).MessageBus()

	db.Insert(_conn, _shopConn)
}

func TestConnectionQueryService(t *testing.T) {
	Convey("Connection QueryService", t, func() {
		Convey("Get Connection Success", func() {
			query := &connectioning.GetConnectionByIDQuery{
				ID: connID,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			conn := query.Result
			So(conn.ID, ShouldEqual, connID)
		})

		Convey("List Connections Success", func() {
			query := &connectioning.ListConnectionsQuery{}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			conns := query.Result
			So(conns[0].ID, ShouldEqual, connID)
		})

		Convey("Get ShopConnection Success", func() {
			query := &connectioning.GetShopConnectionByIDQuery{
				ShopID:       shopID,
				ConnectionID: connID,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			shopConn := query.Result
			So(shopConn.ConnectionID, ShouldEqual, connID)
		})

		Convey("List ShopConnections Success", func() {
			query := &connectioning.ListShopConnectionsQuery{
				ShopID: shopID,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			shopConns := query.Result
			So(shopConns[0].ConnectionID, ShouldEqual, connID)
		})
	})
}
