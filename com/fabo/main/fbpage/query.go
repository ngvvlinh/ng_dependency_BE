package fbpage

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
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

func FbPageQueryMessageBus(q *FbPageQuery) fbpaging.QueryBus {
	b := bus.New()
	return fbpaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbPageQuery) GetFbExternalPageByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPage, error) {
	return q.fbPageStore(ctx).ExternalID(externalID).GetFbExternalPage()
}

func (q *FbPageQuery) ListFbExternalPagesByIDs(
	ctx context.Context, IDs []dot.ID,
) ([]*fbpaging.FbExternalPage, error) {
	return q.fbPageStore(ctx).IDs(IDs).ListFbPages()
}

func (q *FbPageQuery) ListFbExternalPagesByExternalIDs(
	ctx context.Context, external_IDs []string,
) ([]*fbpaging.FbExternalPage, error) {
	fbPages, err := q.fbPageStore(ctx).ExternalIDs(external_IDs).ListFbPages()
	if err != nil {
		return nil, err
	}
	return fbPages, nil
}

func (q *FbPageQuery) ListFbExternalPages(
	ctx context.Context, args *fbpaging.ListFbExternalPagesArgs,
) (*fbpaging.FbPagesResponse, error) {
	query := q.fbPageStore(ctx).OptionalShopID(args.ShopID).
		WithPaging(args.Paging).Filters(args.Filters)
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
