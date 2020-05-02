package fabo

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/pkg/fbclient"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/capi/dot"
)

func (s *PageService) Clone() *PageService {
	res := *s
	return &res
}

func (s *PageService) RemovePages(ctx context.Context, r *RemovePagesEndpoint) error {
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

func (s *PageService) ListPages(ctx context.Context, r *ListPagesEndpoint) error {
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
	r.Result = &fabo.ListPagesResponse{
		FbPages: convertpb.PbFbPages(listFbPagesQuery.Result.FbPages),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *PageService) ConnectPages(ctx context.Context, r *ConnectPagesEndpoint) error {
	shopID := r.Context.Shop.ID
	userID := r.Context.UserID

	userToken, err := fbClient.CallAPICheckAccessToken(r.AccessToken)
	if err != nil {
		return err
	}

	// verify permissions
	if err := verifyScopes(userToken.Data.Scopes); err != nil {
		return err
	}

	if r.AccessToken == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "access_token must not be null")
	}
	longLivedAccessToken, err := fbClient.CallAPIGetLongLivedAccessToken(r.AccessToken)
	if err != nil {
		return err
	}

	me, err := fbClient.CallAPIGetMe(longLivedAccessToken.AccessToken)
	if err != nil {
		return err
	}

	accounts, err := fbClient.CallAPIGetAccounts(longLivedAccessToken.AccessToken)
	if err != nil {
		return err
	}

	fbUserID := cm.NewID()
	createFbUserCombinedCmd := &fbusering.CreateFbUserCombinedCommand{
		UserID: userID,
		ShopID: shopID,
		FbUser: &fbusering.CreateFbUserArgs{
			ID:         fbUserID,
			ExternalID: me.ID,
			UserID:     userID,
			ExternalInfo: &fbusering.ExternalFBUserInfo{
				Name:      me.Name,
				FirstName: me.FirstName,
				LastName:  me.LastName,
				ShortName: me.ShortName,
				ImageURL:  me.Picture.Data.Url,
			},
			Token:  longLivedAccessToken.AccessToken,
			Status: status3.P,
		},
		FbUserInternal: &fbusering.CreateFbUserInternalArgs{
			ID:        fbUserID,
			Token:     longLivedAccessToken.AccessToken,
			ExpiresIn: fbclient.ExpiresInUserToken, // 60 days
		},
	}
	if err := fbUserAggr.Dispatch(ctx, createFbUserCombinedCmd); err != nil {
		return err
	}
	fbUserID = createFbUserCombinedCmd.Result.FbUser.ID

	var fbErrorPages []*fabo.FbErrorPage

	listCreateFbPageCombinedCmd := make([]*fbpaging.CreateFbPageCombinedArgs, 0, len(accounts.Accounts.Data))
	for _, account := range accounts.Accounts.Data {
		// Verify role (Admin)
		if fbclient.GetRole(account.Tasks) != fbclient.ADMIN {
			fbErrorPages = append(fbErrorPages, &fabo.FbErrorPage{
				ExternalID:       account.Id,
				ExternalName:     account.Name,
				ExternalImageURL: account.Picture.Data.Url,
				Reason:           "Tài khoản Facebook cần có quyền Admin trên Fanpage để kết nối.",
			})
			continue
		}

		fbPageID := cm.NewID()
		categories := make([]*fbpaging.ExternalCategory, 0, len(account.CategoryList))
		for _, category := range account.CategoryList {
			categories = append(categories, &fbpaging.ExternalCategory{
				ID:   category.ID,
				Name: category.Name,
			})
		}
		createFbPageCmd := &fbpaging.CreateFbPageArgs{
			ID:                   fbPageID,
			ExternalID:           account.Id,
			FbUserID:             fbUserID,
			ShopID:               shopID,
			UserID:               userID,
			ExternalName:         account.Name,
			ExternalCategory:     account.Category,
			ExternalCategoryList: categories,
			ExternalTasks:        account.Tasks,
			ExternalImageURL:     account.Picture.Data.Url,
			Status:               status3.P,
			ConnectionStatus:     status3.P,
		}
		createFbPageInternalCmd := &fbpaging.CreateFbPageInternalArgs{
			ID:    fbPageID,
			Token: account.AccessToken,
		}
		listCreateFbPageCombinedCmd = append(listCreateFbPageCombinedCmd, &fbpaging.CreateFbPageCombinedArgs{
			FbPage:         createFbPageCmd,
			FbPageInternal: createFbPageInternalCmd,
		})
	}
	var fbPageCombinedsResult []*fabo.FbPageCombined

	if len(listCreateFbPageCombinedCmd) > 0 {
		createFbPageCombinedsCmd := &fbpaging.CreateFbPageCombinedsCommand{
			ShopID:          shopID,
			UserID:          userID,
			FbPageCombineds: listCreateFbPageCombinedCmd,
			Result:          nil,
		}
		if err := fbPageAggr.Dispatch(ctx, createFbPageCombinedsCmd); err != nil {
			return err
		}

		fbPageCombinedsResult = convertpb.PbFbPageCombineds(createFbPageCombinedsCmd.Result)
	}

	r.Result = &fabo.ConnectPagesResponse{
		FbUser:       convertpb.PbFbUserCombined(createFbUserCombinedCmd.Result),
		FbPages:      fbPageCombinedsResult,
		FbErrorPages: fbErrorPages,
	}
	return nil
}

func verifyScopes(scopes []string) error {
	mapScope := make(map[string]bool)
	for _, scope := range scopes {
		mapScope[scope] = true
	}

	for scope, messageScope := range appScopes {
		if _, ok := mapScope[scope]; !ok {
			return cm.Errorf(cm.FacebookPermissionDenied, nil, "Bạn chưa cấp đủ quyền để tiếp tục").
				WithMeta(fmt.Sprintf("scope.%s", scope), messageScope)
		}
	}
	return nil
}
