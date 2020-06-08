package faboinfo

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/capi/dot"
)

type FaboPages struct {
	FbPageIDs       []dot.ID
	ExternalPageIDs []string
}

type FaboPagesKit struct {
	FBPageQuery fbpaging.QueryBus
	FBUserQuery fbusering.QueryBus
}

func (fi *FaboPagesKit) GetPages(ctx context.Context, shopID dot.ID) (*FaboPages, error) {
	listFbPagesQuery := &fbpaging.ListFbExternalPagesQuery{
		ShopID: shopID,
		Filters: []meta.Filter{
			{
				Name:  "status",
				Op:    "=",
				Value: "P",
			},
		},
	}
	if err := fi.FBPageQuery.Dispatch(ctx, listFbPagesQuery); err != nil {
		return nil, err
	}

	fbPageIDs := make([]dot.ID, 0, len(listFbPagesQuery.Result.FbPages))
	externalPageIDs := make([]string, 0, len(listFbPagesQuery.Result.FbPages))
	for _, fbPage := range listFbPagesQuery.Result.FbPages {
		fbPageIDs = append(fbPageIDs, fbPage.ID)
		externalPageIDs = append(externalPageIDs, fbPage.ExternalID)
	}

	return &FaboPages{
		FbPageIDs:       fbPageIDs,
		ExternalPageIDs: externalPageIDs,
	}, nil
}
