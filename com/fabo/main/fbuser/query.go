package fbuser

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ fbusering.QueryService = &FbUserQuery{}

type FbUserQuery struct {
	db                  *cmsql.Database
	fbUserStore         sqlstore.FbUserStoreFactory
	fbUserInternalStore sqlstore.FbUserInternalFactory
}

func (q *FbUserQuery) GetFbUserInternalByID(
	ctx context.Context, args *fbusering.GetFbUserInternalByIDArgs,
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
	ctx context.Context, args *fbusering.GetFbUserByIDArgs,
) (*fbusering.FbUser, error) {
	panic("implement me")
}

func (f FbUserQuery) GetFbUserByExternalID(
	ctx context.Context, args *fbusering.GetFbUserByExternalIDArgs,
) (*fbusering.FbUser, error) {
	panic("implement me")
}

func (q *FbUserQuery) GetFbUserByUserID(
	ctx context.Context, args *fbusering.GetFbUserByUserIDArgs,
) (*fbusering.FbUser, error) {
	return q.fbUserStore(ctx).UserID(args.UserID).Status(status3.P).GetFbUser()
}
