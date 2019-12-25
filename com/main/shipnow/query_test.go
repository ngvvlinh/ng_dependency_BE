package shipnow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"etop.vn/api/main/shipnow"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
)

func TestQueryService(t *testing.T) {
	db, err := cmsql.Connect(cc.DefaultPostgres())
	if err != nil {
		panic(err)
	}
	shipnowQuery := NewQueryService(db)
	shipnowQueryBus := shipnowQuery.MessageBus()
	ctx := context.Background()

	query := &shipnow.GetShipnowFulfillmentQuery{Id: 100}
	err = shipnowQueryBus.Dispatch(ctx, query)
	require.Error(t, err)

	t.Logf("got error: %v", err)
	require.NotContains(t, err.Error(), "Handler not found")
}
