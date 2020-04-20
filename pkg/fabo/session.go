package fabo

import (
	"context"
	"fmt"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/fabo/fbusering"
	"etop.vn/api/top/int/fabo"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/fabo/util"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/fabo/convertpb"
)

func (s *SessionService) InitSession(ctx context.Context, r *InitSessionEndpoint) error {
	shopID := r.Context.Shop.ID
	userID := r.Context.UserID

	// TODO: verify token and permissions
	userToken, err := util.CallAPICheckAccessToken(r.AccessToken)
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
	longLivedAccessToken, err := util.CallAPIGetLongLivedAccessToken(r.AccessToken)
	if err != nil {
		return err
	}

	me, err := util.CallAPIGetMe(longLivedAccessToken.AccessToken)
	if err != nil {
		return err
	}

	accounts, err := util.CallAPIGetAccounts(longLivedAccessToken.AccessToken)
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
			ExpiresIn: util.ExpiresInUserToken, // 60 days
		},
	}
	if err := fbUserAggr.Dispatch(ctx, createFbUserCombinedCmd); err != nil {
		return err
	}
	fbUserID = createFbUserCombinedCmd.Result.FbUser.ID

	listCreateFbPageCombinedCmd := make([]*fbpaging.CreateFbPageCombinedArgs, 0, len(accounts.Accounts.Data))
	for _, account := range accounts.Accounts.Data {
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
			Status:               status3.P,
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
	createFbPageCombinedsCmd := &fbpaging.CreateFbPageCombinedsCommand{
		ShopID:          shopID,
		UserID:          userID,
		FbPageCombineds: listCreateFbPageCombinedCmd,
		Result:          nil,
	}
	if err := fbPageAggr.Dispatch(ctx, createFbPageCombinedsCmd); err != nil {
		return err
	}

	r.Result = &fabo.InitSessionResponse{
		FbUser:  convertpb.PbFbUserCombined(createFbUserCombinedCmd.Result),
		FbPages: convertpb.PbFbPageCombineds(createFbPageCombinedsCmd.Result),
	}
	return nil
}

func verifyScopes(scopes []string) error {
	{
		mapScope := make(map[string]bool)
		for _, scope := range scopes {
			mapScope[scope] = true
		}

		for _, scope := range appScopes {
			if _, ok := mapScope[scope]; !ok {
				return cm.Errorf(cm.PermissionDenied, nil, fmt.Sprintf("Bạn chưa cấp đủ quyền (%s) để tiếp tục", scope))
			}
		}
	}
	return nil
}
