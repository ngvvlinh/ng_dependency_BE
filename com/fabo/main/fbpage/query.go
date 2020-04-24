package fbpage

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/filter"
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
	ctx context.Context, ID dot.ID,
) (*fbpaging.FbPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbPageByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbPageInternalByID(
	ctx context.Context, ID dot.ID,
) (*fbpaging.FbPageInternal, error) {
	panic("implement me")
}

func (q *FbPageQuery) ListFbPagesByIDs(
	ctx context.Context, IDs filter.IDs,
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

func (q *FbPageQuery) ListFbPagesActiveByExternalIDs(
	ctx context.Context, externalIDs []string,
) ([]*fbpaging.FbPage, error) {
	fbPages, err := q.fbPageStore(ctx).ExternalIDs(externalIDs).Status(status3.P).ListFbPages()
	if err != nil {
		return nil, err
	}
	return fbPages, nil
}
