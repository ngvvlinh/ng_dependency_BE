package fbpage

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/backend/com/fabo/main/fbpage/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ fbpaging.QueryService = &FbPageQuery{}

type FbPageQuery struct {
	db                  *cmsql.Database
	fbPageStore         sqlstore.FbPageStoreFactory
	fbPageInternalStore sqlstore.FbPageInternalStoreFactory
}

func NewFbPageQuery(database *cmsql.Database) *FbPageQuery {
	return &FbPageQuery{
		db:                  database,
		fbPageStore:         sqlstore.NewFbPageStore(database),
		fbPageInternalStore: sqlstore.NewFbPageInternalStore(database),
	}
}

func (q *FbPageQuery) MessageBus() fbpaging.QueryBus {
	b := bus.New()
	return fbpaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (f FbPageQuery) GetFbPageByID(
	ctx context.Context, args *fbpaging.GetFbPageByIDArgs,
) (*fbpaging.FbPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbPageByExternalID(
	ctx context.Context, args *fbpaging.GetFbPageByExternalIDArgs,
) (*fbpaging.FbPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbPageInternalByID(
	ctx context.Context, args *fbpaging.GetFbPageInternalByIDArgs,
) (*fbpaging.FbPageInternal, error) {
	panic("implement me")
}

func (q *FbPageQuery) ListFbPagesByIDs(
	ctx context.Context, args *fbpaging.ListFbPagesByIDsArgs,
) ([]*fbpaging.FbPage, error) {
	panic("implement me")
}

func (q *FbPageQuery) ListFbPages(
	ctx context.Context, args *fbpaging.ListFbPagesArgs,
) (*fbpaging.FbPagesResponse, error) {
	query := q.fbPageStore(ctx).ShopID(args.ShopID).UserID(args.UserID).
		WithPaging(args.Paging).Filters(args.Filters)
	if args.FbUserID.Valid {
		query = query.FbUserID(args.FbUserID.ID)
	}
	fbPages, err := query.ListFbPages()
	if err != nil {
		return nil, err
	}
	return &fbpaging.FbPagesResponse{
		FbPages: fbPages,
		Paging:  query.GetPaging(),
	}, nil
}
