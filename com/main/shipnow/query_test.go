package shipnow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"o.o/api/main/shipnow"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

func TestQueryService(t *testing.T) {
	db, err := cmsql.Connect(cc.DefaultPostgres())
	if err != nil {
		panic(err)
	}
	shipnowQuery := NewQueryService(db)
	shipnowQueryBus := QueryServiceMessageBus(shipnowQuery)
	ctx := context.Background()

	query := &shipnow.GetShipnowFulfillmentQuery{ID: 100}
	err = shipnowQueryBus.Dispatch(ctx, query)
	require.Error(t, err)

	t.Logf("got error: %v", err)
	require.NotContains(t, err.Error(), "Handler not found")
}
