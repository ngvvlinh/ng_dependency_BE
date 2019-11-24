package pgrid

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"etop.vn/backend/com/handler/pgevent"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/common/xerrors"
)

var db *cmsql.Database

func init() {
	cfg := cc.DefaultPostgres()
	db = cmsql.MustConnect(cfg)
}

func TestQuery(t *testing.T) {
	{
		var foo Foo
		s, args, err := db.
			From("fulfillment AS f").
			SQL("JOIN history.fulfillment AS hf").
			Where("hf.rid = ?", 1000).
			BuildGet(&foo)

		require.NoError(t, err)
		require.Equal(t, len(args), 1)
		require.Equal(t, args[0], 1000)

		expect := `SELECT f.id, f.shop_id, hf._time FROM fulfillment AS f JOIN history.fulfillment AS hf WHERE (hf.rid = $1)`
		require.Equal(t, expect, s)
	}
	{
		var foo Foo
		sql := `FROM fulfillment AS f JOIN history.fulfillment AS hf`
		s, args, err := db.SQL(sql).
			Where("hf.rid = ?", 1000).
			BuildGet(&foo)

		require.NoError(t, err)
		require.Equal(t, len(args), 1)
		require.Equal(t, args[0], 1000)

		expect := `SELECT f.id, f.shop_id, hf._time FROM fulfillment AS f JOIN history.fulfillment AS hf WHERE (hf.rid = $1)`
		require.Equal(t, expect, s)
	}
}

func TestModel(t *testing.T) {
	cases := []struct {
		Model        IModel
		ExpectedCode xerrors.Code
	}{
		// {&UserEvent{}, cm.NotFound},
		// {&ShopEvent{}, cm.NotFound},
		{&ShopProductEvent{}, cm.NotFound},
		{&OrderEvent{}, cm.NotFound},
		{&FulfillmentEvent{}, cm.NotFound},
	}
	for _, c := range cases {
		t.Run(reflect.TypeOf(c.Model).Name(), func(t *testing.T) {
			event := &pgevent.PgEvent{RID: -1}
			err := cmsql.ShouldGet(c.Model.Query(db, event))
			if cm.ErrorCode(err) != c.ExpectedCode {
				t.Errorf("Expect err code: %s, Get: %v", c.ExpectedCode, err)
			}
			assert.Equal(t, "-1", c.Model._meta().RID)
		})
	}
}
