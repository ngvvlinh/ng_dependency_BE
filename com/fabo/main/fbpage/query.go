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
	fbPageStore         sqlstore.FbExternalPageStoreFactory
	fbPageInternalStore sqlstore.FbExternalPageInternalStoreFactory
}

func NewFbPageQuery(database *cmsql.Database) *FbPageQuery {
	return &FbPageQuery{
		db:                  database,
		fbPageStore:         sqlstore.NewFbExternalPageStore(database),
		fbPageInternalStore: sqlstore.NewFbExternalPageInternalStore(database),
	}
}

func (q *FbPageQuery) MessageBus() fbpaging.QueryBus {
	b := bus.New()
	return fbpaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (f FbPageQuery) GetFbExternalPageByID(
	ctx context.Context, ID dot.ID,
) (*fbpaging.FbExternalPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbExternalPageByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPage, error) {
	panic("implement me")
}

func (f FbPageQuery) GetFbExternalPageInternalByID(
	ctx context.Context, ID dot.ID,
) (*fbpaging.FbExternalPageInternal, error) {
	panic("implement me")
}

func (q *FbPageQuery) ListFbExternalPagesByIDs(
	ctx context.Context, IDs filter.IDs,
) ([]*fbpaging.FbExternalPage, error) {
	panic("implement me")
}

func (q *FbPageQuery) ListFbExternalPages(
	ctx context.Context, args *fbpaging.ListFbExternalPagesArgs,
) (*fbpaging.FbPagesResponse, error) {
	query := q.fbPageStore(ctx).OptionalShopID(args.ShopID).UserID(args.UserID).
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

func (q *FbPageQuery) ListFbExternalPagesActiveByExternalIDs(
	ctx context.Context, externalIDs []string,
) ([]*fbpaging.FbExternalPage, error) {
	fbPages, err := q.fbPageStore(ctx).ExternalIDs(externalIDs).Status(status3.P).ListFbPages()
	if err != nil {
		return nil, err
	}
	return fbPages, nil
}
