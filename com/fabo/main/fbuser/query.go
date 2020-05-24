package fbuser

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/filter"
)

var _ fbusering.QueryService = &FbUserQuery{}

type FbUserQuery struct {
	db                  *cmsql.Database
	fbUserStore         sqlstore.FbExternalUserStoreFactory
	fbUserInternalStore sqlstore.FbExternalUserInternalFactory
}

func NewFbUserQuery(database *cmsql.Database) *FbUserQuery {
	return &FbUserQuery{
		db:          database,
		fbUserStore: sqlstore.NewFbExternalUserStore(database),
	}
}

func FbUserQueryMessageBus(q *FbUserQuery) fbusering.QueryBus {
	b := bus.New()
	return fbusering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbUserQuery) GetFbExternalUserInternalByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbExternalUserInternal, error) {
	panic("implement me")
}

func (f FbUserQuery) GetFbExternalUserByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbExternalUser, error) {
	return f.fbUserStore(ctx).ExternalID(externalID).GetFbExternalUser()
}

func (q *FbUserQuery) ListFbExternalUsersByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbusering.FbExternalUser, error) {
	return q.fbUserStore(ctx).ExternalIDs(externalIDs).ListFbExternalUsers()
}
