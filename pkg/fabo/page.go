package fabo

import (
	"context"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/top/int/fabo"
	"etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/fabo/convertpb"
	"etop.vn/capi/dot"
)

func (s *PageService) RemoveFbPages(ctx context.Context, r *RemoveFbPagesEndpoint) error {
	if len(r.IDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "ids must not be null")
	}
	disablePagesByIDsCmd := &fbpaging.DisableFbPagesByIDsCommand{
		IDs:    r.IDs,
		ShopID: r.Context.Shop.ID,
		UserID: r.Context.UserID,
	}
	if err := fbPageAggr.Dispatch(ctx, disablePagesByIDsCmd); err != nil {
		return err
	}

	r.Result = &common.Empty{}
	return nil
}

func (s *PageService) ListFbPages(ctx context.Context, r *ListFbPagesEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	listFbPagesQuery := &fbpaging.ListFbPagesQuery{
		ShopID:   r.Context.Shop.ID,
		UserID:   r.Context.UserID,
		FbUserID: dot.WrapID(r.Context.FaboInfo.FbUserID),
		Paging:   *paging,
		Filters:  cmapi.ToFilters(r.Filters),
	}
	if err := fbPageQuery.Dispatch(ctx, listFbPagesQuery); err != nil {
		return err
	}
	r.Result = &fabo.FbPagesResponse{
		FbPages: convertpb.PbFbPages(listFbPagesQuery.Result.FbPages),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}
