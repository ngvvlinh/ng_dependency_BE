package faboinfo

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/capi/dot"
)

type Info struct {
	FbUserID  dot.ID
	FbPageIDs []dot.ID
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

func (fi *FaboInfo) GetFaboInfo(ctx context.Context, shopID, userID dot.ID) (*Info, error) {
	getFbUserByIDQuery := &fbusering.GetFbExternalUserByUserIDQuery{
		UserID: userID,
	}
	if err := fi.fbUserQuery.Dispatch(ctx, getFbUserByIDQuery); err != nil {
		return nil, err
	}
	fbUserID := getFbUserByIDQuery.Result.ID

	listFbPagesQuery := &fbpaging.ListFbExternalPagesQuery{
		ShopID:   shopID,
		UserID:   userID,
		FbUserID: dot.WrapID(fbUserID),
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
	for _, fbPage := range listFbPagesQuery.Result.FbPages {
		fbPageIDs = append(fbPageIDs, fbPage.ID)
	}

	return &Info{
		FbUserID:  fbUserID,
		FbPageIDs: fbPageIDs,
	}, nil
}
