package fbuser

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ fbusering.QueryService = &FbUserQuery{}

type FbUserQuery struct {
	db                  *cmsql.Database
	fbUserStore         sqlstore.FbExternalUserStoreFactory
	fbUserInternalStore sqlstore.FbExternalUserInternalFactory
}

func (q *FbUserQuery) GetFbExternalUserInternalByID(
	ctx context.Context, ID dot.ID,
) (*fbusering.FbExternalUserInternal, error) {
	panic("implement me")
}

func NewFbUserQuery(database *cmsql.Database) *FbUserQuery {
	return &FbUserQuery{
		db:          database,
		fbUserStore: sqlstore.NewFbExternalUserStore(database),
	}
}

func (q *FbUserQuery) MessageBus() fbusering.QueryBus {
	b := bus.New()
	return fbusering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (f FbUserQuery) GetFbExternalUserByID(
	ctx context.Context, ID dot.ID,
) (*fbusering.FbExternalUser, error) {
	panic("implement me")
}

func (f FbUserQuery) GetFbExternalUserByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbExternalUser, error) {
	panic("implement me")
}

func (q *FbUserQuery) GetFbExternalUserByUserID(
	ctx context.Context, userID dot.ID,
) (*fbusering.FbExternalUser, error) {
	return q.fbUserStore(ctx).UserID(userID).Status(status3.P).GetFbExternalUser()
}
