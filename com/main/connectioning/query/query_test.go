package query

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"etop.vn/api/main/connectioning"
	"etop.vn/backend/com/main/connectioning/model"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var (
	db     *cmsql.Database
	connID = dot.ID(1234)
	shopID = dot.ID(4567)
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

func TestConnectionQueryService(t *testing.T) {
	Convey("Connection QueryService", t, func() {
		Reset(func() {
			db.MustExec("truncate connection, shop_connection")
		})
		_conn := &model.Connection{
			ID:     connID,
			Name:   "Connection",
			Status: 1,
		}
		_shopConn := &model.ShopConnection{
			ShopID:       shopID,
			ConnectionID: connID,
			Token:        "token",
			Status:       1,
		}

		QS := NewConnectionQuery(db).MessageBus()

		ctx := context.Background()
		_, err := db.Insert(_conn, _shopConn)
		So(err, ShouldBeNil)

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
			query := &connectioning.ListShopConnectionsQuery{}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			shopConns := query.Result
			So(shopConns[0].ConnectionID, ShouldEqual, connID)
		})
	})
}
