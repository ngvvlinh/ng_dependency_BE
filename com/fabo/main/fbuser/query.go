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
	fbUserStore         sqlstore.FbUserStoreFactory
	fbUserInternalStore sqlstore.FbUserInternalFactory
}

func (q *FbUserQuery) GetFbUserInternalByID(
	ctx context.Context, ID dot.ID,
) (*fbusering.FbUserInternal, error) {
	panic("implement me")
}

func NewFbUserQuery(database *cmsql.Database) *FbUserQuery {
	return &FbUserQuery{
		db:          database,
		fbUserStore: sqlstore.NewFbUserStore(database),
	}
}

func (q *FbUserQuery) MessageBus() fbusering.QueryBus {
	b := bus.New()
	return fbusering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (f FbUserQuery) GetFbUserByID(
	ctx context.Context, ID dot.ID,
) (*fbusering.FbUser, error) {
	panic("implement me")
}

func (f FbUserQuery) GetFbUserByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbUser, error) {
	panic("implement me")
}

func (q *FbUserQuery) GetFbUserByUserID(
	ctx context.Context, userID dot.ID,
) (*fbusering.FbUser, error) {
	return q.fbUserStore(ctx).UserID(userID).Status(status3.P).GetFbUser()
}
