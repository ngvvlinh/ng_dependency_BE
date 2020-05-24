package fbpage

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/convert"
	"o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/com/fabo/main/fbpage/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbExternalPageAggregate struct {
	db                          *cmsql.Database
	fbExternalPageStore         sqlstore.FbExternalPageStoreFactory
	fbExternalPageInternalStore sqlstore.FbExternalPageInternalStoreFactory
}

func NewFbPageAggregate(db *cmsql.Database) *FbExternalPageAggregate {
	return &FbExternalPageAggregate{
		db:                          db,
		fbExternalPageStore:         sqlstore.NewFbExternalPageStore(db),
		fbExternalPageInternalStore: sqlstore.NewFbExternalPageInternalStore(db),
	}
}

func FbExternalPageAggregateMessageBus(a *FbExternalPageAggregate) fbpaging.CommandBus {
	b := bus.New()
	return fbpaging.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbExternalPageAggregate) CreateFbExternalPage(
	ctx context.Context, args *fbpaging.CreateFbExternalPageArgs,
) (*fbpaging.FbExternalPage, error) {
	fbPageResult := new(fbpaging.FbExternalPage)
	if err := scheme.Convert(args, fbPageResult); err != nil {
		return nil, err
	}
	if err := a.fbExternalPageStore(ctx).CreateFbExternalPage(fbPageResult); err != nil {
		return nil, err
	}
	return fbPageResult, nil
}

func (a *FbExternalPageAggregate) CreateFbExternalPageInternal(
	ctx context.Context, args *fbpaging.CreateFbExternalPageInternalArgs,
) (*fbpaging.FbExternalPageInternal, error) {
	fbPageInternalResult := new(fbpaging.FbExternalPageInternal)
	if err := scheme.Convert(args, fbPageInternalResult); err != nil {
		return nil, err
	}
	if err := a.fbExternalPageInternalStore(ctx).CreateFbExternalPageInternal(fbPageInternalResult); err != nil {
		return nil, err
	}
	return fbPageInternalResult, nil
}

func (a *FbExternalPageAggregate) CreateFbExternalPageCombined(
	ctx context.Context, args *fbpaging.CreateFbExternalPageCombinedArgs,
) (*fbpaging.FbExternalPageCombined, error) {
	panic("implement me")
}

func (a *FbExternalPageAggregate) CreateFbExternalPageCombineds(
	ctx context.Context, args *fbpaging.CreateFbExternalPageCombinedsArgs,
) ([]*fbpaging.FbExternalPageCombined, error) {
	// create map arguments with external_id
	mapExternalIDAndFbPageCombined := make(map[string]*fbpaging.CreateFbExternalPageCombinedArgs)
	var externalIDs []string

	for _, fbPageCombined := range args.FbPageCombineds {
		if _, ok := mapExternalIDAndFbPageCombined[fbPageCombined.FbPage.ExternalID]; !ok {
			mapExternalIDAndFbPageCombined[fbPageCombined.FbPage.ExternalID] = fbPageCombined
			externalIDs = append(externalIDs, fbPageCombined.FbPage.ExternalID)
		}
	}

	// get all oldFbPages by (shop_id) from DB
	// create map fbPageDisabled (oldFbPages don't appear into arg)
	oldFbPages, err := a.fbExternalPageStore(ctx).ExternalIDs(externalIDs).ListFbExternalPagesDB()
	if err != nil {
		return nil, err
	}

	mapFbPageDisabled := make(map[string]*model.FbExternalPage)
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

	mapFbPages := make(map[dot.ID]*fbpaging.FbExternalPage)
	mapFbPageInternals := make(map[dot.ID]*fbpaging.FbExternalPageInternal)
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		//// Disable fbPages (mapFbPagesDisabled)
		//{
		//	externalIDs := make([]string, 0, len(mapFbPageDisabled))
		//	for externalID := range mapFbPageDisabled {
		//		externalIDs = append(externalIDs, externalID)
		//	}
		//
		//	if len(externalIDs) > 0 {
		//		if _, err := a.fbExternalPageStore(ctx).ExternalIDs(externalIDs).UpdateConnectionStatus(int(status3.N)); err != nil {
		//			return err
		//		}
		//	}
		//}

		// Create newfbPages (mapFbPagesEnabled)
		{
			newFbPageModels := make([]*fbpaging.FbExternalPage, 0, len(args.FbPageCombineds))
			newFbPageInternalModels := make([]*fbpaging.FbExternalPageInternal, 0, len(args.FbPageCombineds))
			for _, fbPageCombined := range args.FbPageCombineds {
				newFbPage := new(fbpaging.FbExternalPage)
				newFbPageInternal := new(fbpaging.FbExternalPageInternal)
				if err := scheme.Convert(fbPageCombined.FbPage, newFbPage); err != nil {
					return err
				}
				if err := scheme.Convert(fbPageCombined.FbPageInternal, newFbPageInternal); err != nil {
					return err
				}
				newFbPageModels = append(newFbPageModels, newFbPage)
				newFbPageInternalModels = append(newFbPageInternalModels, newFbPageInternal)
			}

			if err := a.fbExternalPageStore(ctx).CreateFbExternalPages(newFbPageModels); err != nil {
				return err
			}
			for _, newFbPage := range newFbPageModels {
				mapFbPages[newFbPage.ID] = newFbPage
			}

			if err := a.fbExternalPageInternalStore(ctx).CreateFbExternalPageInternals(newFbPageInternalModels); err != nil {
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

	fbPageCombineds := make([]*fbpaging.FbExternalPageCombined, 0, len(args.FbPageCombineds))
	for fbPageID := range mapFbPages {
		fbPageCombineds = append(fbPageCombineds, &fbpaging.FbExternalPageCombined{
			FbExternalPage:         mapFbPages[fbPageID],
			FbExternalPageInternal: mapFbPageInternals[fbPageID],
		})
	}
	return fbPageCombineds, err
}

func (a *FbExternalPageAggregate) DisableFbExternalPagesByExternalIDs(
	ctx context.Context, args *fbpaging.DisableFbExternalPagesByIDsArgs,
) (int, error) {
	return a.fbExternalPageStore(ctx).ShopID(args.ShopID).ExternalIDs(args.ExternalIDs).UpdateStatus(int(status3.N))
}

func (a *FbExternalPageAggregate) DisableAllFbExternalPages(
	ctx context.Context, args *fbpaging.DisableAllFbExternalPagesArgs,
) (int, error) {
	return a.fbExternalPageStore(ctx).ShopID(args.ShopID).Status(status3.P).UpdateStatus(int(status3.N))
}
