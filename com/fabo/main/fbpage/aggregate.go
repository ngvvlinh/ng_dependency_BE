package fbpage

import (
	"context"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/fabo/main/fbpage/convert"
	"etop.vn/backend/com/fabo/main/fbpage/model"
	"etop.vn/backend/com/fabo/main/fbpage/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbPageAggregate struct {
	db                  *cmsql.Database
	fbPageStore         sqlstore.FbPageStoreFactory
	fbPageInternalStore sqlstore.FbPageInternalStoreFactory
}

func NewFbPageAggregate(db *cmsql.Database) *FbPageAggregate {
	return &FbPageAggregate{
		db:                  db,
		fbPageStore:         sqlstore.NewFbPageStore(db),
		fbPageInternalStore: sqlstore.NewFbPageInternalStore(db),
	}
}

func (a *FbPageAggregate) MessageBus() fbpaging.CommandBus {
	b := bus.New()
	return fbpaging.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbPageAggregate) CreateFbPage(
	ctx context.Context, args *fbpaging.CreateFbPageArgs,
) (*fbpaging.FbPage, error) {
	fbPageResult := new(fbpaging.FbPage)
	if err := scheme.Convert(args, fbPageResult); err != nil {
		return nil, err
	}
	if err := a.fbPageStore(ctx).CreateFbPage(fbPageResult); err != nil {
		return nil, err
	}
	return fbPageResult, nil
}

func (a *FbPageAggregate) CreateFbPageInternal(
	ctx context.Context, args *fbpaging.CreateFbPageInternalArgs,
) (*fbpaging.FbPageInternal, error) {
	fbPageInternalResult := new(fbpaging.FbPageInternal)
	if err := scheme.Convert(args, fbPageInternalResult); err != nil {
		return nil, err
	}
	if err := a.fbPageInternalStore(ctx).CreateFbPageInternal(fbPageInternalResult); err != nil {
		return nil, err
	}
	return fbPageInternalResult, nil
}

func (a *FbPageAggregate) CreateFbPageCombined(
	ctx context.Context, args *fbpaging.CreateFbPageCombinedArgs,
) (*fbpaging.FbPageCombined, error) {
	panic("implement me")
}

func (a *FbPageAggregate) CreateFbPageCombineds(
	ctx context.Context, args *fbpaging.CreateFbPageCombinedsArgs,
) ([]*fbpaging.FbPageCombined, error) {
	// create map arguments with external_id
	mapExternalIDAndFbPageCombined := make(map[string]*fbpaging.CreateFbPageCombinedArgs)

	for _, fbPageCombined := range args.FbPageCombineds {
		mapExternalIDAndFbPageCombined[fbPageCombined.FbPage.ExternalID] = fbPageCombined
	}

	// get all oldFbPages by (user_id, shop_id) from DB
	// create map fbPageDisabled (oldFbPages don't appear into arg)
	oldFbPages, err := a.fbPageStore(ctx).UserID(args.UserID).ShopID(args.ShopID).ListFbPagesDB()
	if err != nil {
		return nil, err
	}

	mapFbPageDisabled := make(map[string]*model.FbPage)
	for _, oldFbPage := range oldFbPages {
		if _, ok := mapExternalIDAndFbPageCombined[oldFbPage.ExternalID]; !ok {
			mapFbPageDisabled[oldFbPage.ExternalID] = oldFbPage
		} else {
			// Depend on (ON CONFLICT), we change the IDs of args to same with oldFbPages
			// Then when create args, db will update elements with exists IDs
			mapExternalIDAndFbPageCombined[oldFbPage.ExternalID].FbPage.ID = oldFbPage.ID
			mapExternalIDAndFbPageCombined[oldFbPage.ExternalID].FbPageInternal.ID = oldFbPage.ID
		}
	}

	mapFbPages := make(map[dot.ID]*fbpaging.FbPage)
	mapFbPageInternals := make(map[dot.ID]*fbpaging.FbPageInternal)
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// Disable fbPages (mapFbPagesDisabled)
		{
			externalIDs := make([]string, 0, len(mapFbPageDisabled))
			for externalID := range mapFbPageDisabled {
				externalIDs = append(externalIDs, externalID)
			}

			if len(externalIDs) > 0 {
				if _, err := a.fbPageStore(ctx).ExternalIDs(externalIDs).UpdateStatus(int(status3.N)); err != nil {
					return err
				}
			}
		}

		// Create newfbPages (mapFbPagesEnabled)
		{
			newFbPageModels := make([]*fbpaging.FbPage, 0, len(args.FbPageCombineds))
			newFbPageInternalModels := make([]*fbpaging.FbPageInternal, 0, len(args.FbPageCombineds))
			for _, fbPageCombined := range args.FbPageCombineds {
				newFbPage := new(fbpaging.FbPage)
				newFbPageInternal := new(fbpaging.FbPageInternal)
				if err := scheme.Convert(fbPageCombined.FbPage, newFbPage); err != nil {
					return err
				}
				if err := scheme.Convert(fbPageCombined.FbPageInternal, newFbPageInternal); err != nil {
					return err
				}
				newFbPageModels = append(newFbPageModels, newFbPage)
				newFbPageInternalModels = append(newFbPageInternalModels, newFbPageInternal)
			}

			if err := a.fbPageStore(ctx).CreateFbPages(newFbPageModels); err != nil {
				return err
			}
			for _, newFbPage := range newFbPageModels {
				mapFbPages[newFbPage.ID] = newFbPage
			}

			if err := a.fbPageInternalStore(ctx).CreateFbPageInternals(newFbPageInternalModels); err != nil {
				return err
			}
			for _, newFbPageInternal := range newFbPageInternalModels {
				mapFbPageInternals[newFbPageInternal.ID] = newFbPageInternal
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	fbPageCombineds := make([]*fbpaging.FbPageCombined, 0, len(args.FbPageCombineds))
	for fbPageID := range mapFbPages {
		fbPageCombineds = append(fbPageCombineds, &fbpaging.FbPageCombined{
			FbPage:         mapFbPages[fbPageID],
			FbPageInternal: mapFbPageInternals[fbPageID],
		})
	}
	return fbPageCombineds, err
}

func (a *FbPageAggregate) DisableFbPagesByIDs(
	ctx context.Context, args *fbpaging.DisableFbPagesByIDsArgs,
) (int, error) {
	return a.fbPageStore(ctx).ShopID(args.ShopID).UserID(args.UserID).IDs(args.IDs).UpdateStatus(int(status3.N))
}

func (a *FbPageAggregate) DisableAllFbPages(
	ctx context.Context, args *fbpaging.DisableAllFbPagesArgs,
) (int, error) {
	return a.fbPageStore(ctx).ShopID(args.ShopID).UserID(args.UserID).Status(status3.P).UpdateStatus(int(status3.N))
}
