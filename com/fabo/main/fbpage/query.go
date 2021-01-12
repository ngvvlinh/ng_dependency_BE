package fbpage

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ fbpaging.QueryService = &FbPageQuery{}

type FbPageQuery struct {
	db                  *cmsql.Database
	fbPageStore         sqlstore.FbExternalPageStoreFactory
	fbPageInternalStore sqlstore.FbExternalPageInternalStoreFactory
	fbPageUtil          *FbPageUtil
}

func NewFbPageQuery(database com.MainDB, fbPageUtil *FbPageUtil) *FbPageQuery {
	return &FbPageQuery{
		db:                  database,
		fbPageStore:         sqlstore.NewFbExternalPageStore(database),
		fbPageInternalStore: sqlstore.NewFbExternalPageInternalStore(database),
		fbPageUtil:          fbPageUtil,
	}
}

func FbPageQueryMessageBus(q *FbPageQuery) fbpaging.QueryBus {
	b := bus.New()
	return fbpaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbPageQuery) GetFbExternalPageActiveByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPage, error) {
	fbExternalPageCached, err := q.fbPageUtil.GetFbPage(externalID)
	if err != nil {
		return nil, err
	}
	if fbExternalPageCached != nil {
		return fbExternalPageCached, nil
	}

	return q.fbPageStore(ctx).ExternalID(externalID).Status(status3.P).GetFbExternalPage()
}

func (q *FbPageQuery) GetFbExternalPageByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPage, error) {
	fbExternalPageCached, err := q.fbPageUtil.GetFbPage(externalID)
	if err != nil {
		return nil, err
	}
	if fbExternalPageCached != nil {
		return fbExternalPageCached, nil
	}

	return q.fbPageStore(ctx).ExternalID(externalID).GetFbExternalPage()
}

func (q *FbPageQuery) GetFbExternalPageInternalByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPageInternal, error) {
	fbExternalPageInternalCached, err := q.fbPageUtil.GetFbPageInternal(externalID)
	if err != nil {
		return nil, err
	}
	if fbExternalPageInternalCached != nil {
		return fbExternalPageInternalCached, nil
	}

	return q.fbPageInternalStore(ctx).ExternalID(externalID).GetFbExternalPageInternal()
}

func (q *FbPageQuery) ListFbExternalPagesByExternalIDs(
	ctx context.Context, externalIDs []string,
) ([]*fbpaging.FbExternalPage, error) {
	fbPages, err := q.fbPageStore(ctx).ExternalIDs(externalIDs).ListFbPages()
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

func (q *FbPageQuery) GetFbExternalPageInternalActiveByExternalID(
	ctx context.Context, externalID string,
) (*fbpaging.FbExternalPageInternal, error) {
	fbExternalPageCached, err := q.fbPageUtil.GetFbPage(externalID)
	if err != nil {
		return nil, err
	}

	if fbExternalPageCached != nil {
		if fbExternalPageCached.Status != status3.P {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "FbExternalPage (%s) isn't active", externalID)
		}
	} else {
		if _, err := q.fbPageStore(ctx).ExternalID(externalID).Status(status3.P).GetFbExternalPage(); err != nil {
			return nil, err
		}
	}

	fbExternalPageInternalCached, err := q.fbPageUtil.GetFbPageInternal(externalID)
	if err != nil {
		return nil, err
	}
	if fbExternalPageInternalCached != nil {
		return fbExternalPageInternalCached, nil
	}

	fbExternalPageInternal, err := q.fbPageInternalStore(ctx).ExternalID(externalID).GetFbExternalPageInternal()
	if err != nil {
		return nil, err
	}
	return fbExternalPageInternal, nil
}

func (q *FbPageQuery) ListFbPagesByShop(ctx context.Context, shopIDs []dot.ID) ([]*fbpaging.FbExternalPage, error) {
	return q.fbPageStore(ctx).ShopIDs(shopIDs...).ListFbPages()
}

func (q *FbPageQuery) ListActiveFbPagesByShopIDs(ctx context.Context, shopIDs []dot.ID) ([]*fbpaging.FbExternalPage, error) {
	return q.fbPageStore(ctx).Status(status3.P).ShopIDs(shopIDs...).ListFbPages()
}

func (q *FbPageQuery) GetPageAccessToken(ctx context.Context, externalID string) (string, error) {
	fbExternalPageInternalCached, err := q.fbPageUtil.GetFbPageInternal(externalID)
	if err != nil {
		return "", err
	}
	if fbExternalPageInternalCached != nil {
		return fbExternalPageInternalCached.Token, nil
	}

	return q.fbPageInternalStore(ctx).ExternalID(externalID).GetAccessToken()
}
