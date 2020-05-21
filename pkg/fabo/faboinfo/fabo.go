package faboinfo

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/capi/dot"
)

type Info struct {
	FbPageIDs       []dot.ID
	ExternalPageIDs []string
}

type FaboInfo struct {
	fbPageQuery fbpaging.QueryBus
	fbUserQuery fbusering.QueryBus
}

func New(fbPageQuery fbpaging.QueryBus, fbUserQuery fbusering.QueryBus) *FaboInfo {
	fi := &FaboInfo{
		fbPageQuery: fbPageQuery,
		fbUserQuery: fbUserQuery,
	}
	return fi
}

type GetFaboInfoQuery struct {
	ShopID dot.ID
	UserID dot.ID
}

func (fi *FaboInfo) GetFaboInfo(ctx context.Context, shopID dot.ID) (*Info, error) {
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
	if err := fi.fbPageQuery.Dispatch(ctx, listFbPagesQuery); err != nil {
		return nil, err
	}

	fbPageIDs := make([]dot.ID, 0, len(listFbPagesQuery.Result.FbPages))
	externalPageIDs := make([]string, 0, len(listFbPagesQuery.Result.FbPages))
	for _, fbPage := range listFbPagesQuery.Result.FbPages {
		fbPageIDs = append(fbPageIDs, fbPage.ID)
		externalPageIDs = append(externalPageIDs, fbPage.ExternalID)
	}

	return &Info{
		FbPageIDs:       fbPageIDs,
		ExternalPageIDs: externalPageIDs,
	}, nil
}
