package middleware

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/capi/dot"
)

var (
	fbPageQuery fbpaging.QueryBus
	fbUserQuery fbusering.QueryBus
)

func NewFabo(fbPageQ fbpaging.QueryBus, fbUserQ fbusering.QueryBus) {
	fbPageQuery = fbPageQ
	fbUserQuery = fbUserQ
}

type GetFaboInfoQuery struct {
	ShopID dot.ID
	UserID dot.ID
}

func GetFaboInfo(ctx context.Context, r *GetFaboInfoQuery) (*claims.FaboInfo, error) {
	getFbUserByIDQuery := &fbusering.GetFbUserByUserIDQuery{
		UserID: r.UserID,
	}
	if err := fbUserQuery.Dispatch(ctx, getFbUserByIDQuery); err != nil {
		return nil, err
	}
	fbUserID := getFbUserByIDQuery.Result.ID

	listFbPagesQuery := &fbpaging.ListFbPagesQuery{
		ShopID:   r.ShopID,
		UserID:   r.UserID,
		FbUserID: dot.WrapID(fbUserID),
		Filters: []meta.Filter{
			{
				Name:  "status",
				Op:    "=",
				Value: "P",
			},
		},
	}
	if err := fbPageQuery.Dispatch(ctx, listFbPagesQuery); err != nil {
		return nil, err
	}

	fbPageIDs := make([]dot.ID, 0, len(listFbPagesQuery.Result.FbPages))

	for _, fbPage := range listFbPagesQuery.Result.FbPages {
		fbPageIDs = append(fbPageIDs, fbPage.ID)
	}

	return &claims.FaboInfo{
		FbUserID:  fbUserID,
		FbPageIDs: fbPageIDs,
	}, nil
}
