package fabo

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
)

var appScopes = map[string]string{
	"public_profile":          "Hiển thị thông tin cơ bản của tài khoản",
	"pages_show_list":         "Hiển thị các trang do tài khoản quản lý",
	"pages_messaging":         "Quản lý và truy cập các cuộc trò chuyện của trang",
	"pages_read_engagement":   "",
	"pages_manage_metadata":   "",
	"pages_read_user_content": "",
	"pages_manage_engagement": "",
	"pages_manage_posts":      "",
}

type PageService struct {
	session.Session

	FaboInfo            *faboinfo.FaboPagesKit
	FBExternalUserQuery fbusering.QueryBus
	FBExternalUserAggr  fbusering.CommandBus
	FBExternalPageQuery fbpaging.QueryBus
	FBExternalPageAggr  fbpaging.CommandBus
	FBClient            *fbclient.FbClient
}

func (s *PageService) Clone() fabo.PageService {
	res := *s
	return &res
}

func (s *PageService) RemovePages(ctx context.Context, r *fabo.RemovePagesRequest) (*common.Empty, error) {
	if len(r.ExternalIDs) == 0 && len(r.NewExternalIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ids must not be null")
	}
	externalIDs := r.ExternalIDs
	if len(externalIDs) == 0 {
		externalIDs = r.NewExternalIDs
	}
	disablePagesByIDsCmd := &fbpaging.DisableFbExternalPagesByExternalIDsCommand{
		ExternalIDs: externalIDs,
		ShopID:      s.SS.Shop().ID,
	}
	if err := s.FBExternalPageAggr.Dispatch(ctx, disablePagesByIDsCmd); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *PageService) ListPages(ctx context.Context, r *fabo.ListPagesRequest) (*fabo.ListPagesResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	listFbExternalPagesQuery := &fbpaging.ListFbExternalPagesQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.FBExternalPageQuery.Dispatch(ctx, listFbExternalPagesQuery); err != nil {
		return nil, err
	}
	resp := &fabo.ListPagesResponse{
		FbPages: convertpb.PbFbPages(listFbExternalPagesQuery.Result.FbPages),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return resp, nil
}

func (s *PageService) ConnectPages(ctx context.Context, r *fabo.ConnectPagesRequest) (*fabo.ConnectPagesResponse, error) {
	shopID := s.SS.Shop().ID

	// Check accessToken is alive
	userToken, err := s.FBClient.CallAPICheckAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}

	// Verify permissions
	if err := verifyScopes(appScopes, userToken.Data.Scopes); err != nil {
		return nil, err
	}

	// Get long lived accessToken from accessToken (above)
	if r.AccessToken == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "access_token must not be null")
	}
	longLivedAccessToken, err := s.FBClient.CallAPIGetLongLivedAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get information of user from accessToken (above)
	me, err := s.FBClient.CallAPIGetMe(longLivedAccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get all accounts of user (above)
	accounts, err := s.FBClient.CallAPIGetAccounts(longLivedAccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get externalIDs
	var externalIDs []string
	for _, account := range accounts.Accounts.Data {
		externalIDs = append(externalIDs, account.Id)
	}

	// Subcribe app (enable webhook messager)
	if contains(userToken.Data.Scopes, "pages_messaging") {
		var wg sync.WaitGroup
		wg.Add(len(accounts.Accounts.Data))
		for _, account := range accounts.Accounts.Data {
			go func(accessToken, externalPageID string) {
				defer wg.Done()
				// TODO: Ngoc handle err
				if _, err := s.FBClient.CallAPICreateSubscribedApps(&fbclient.CreateSubscribedAppsRequest{
					AccessToken: accessToken,
					Fields:      []string{fbclient.MessagesField, fbclient.MessageEchoesField, fbclient.FeedField},
					PageID:      externalPageID,
				}); err != nil {
					return
				}
			}(account.AccessToken, account.Id)
		}
		wg.Wait()
	}

	// Get fbPages active from externalIDs (accounts)
	listFbPagesActiveQuery := &fbpaging.ListFbExternalPagesActiveByExternalIDsQuery{
		ExternalIDs: externalIDs,
	}
	if err := s.FBExternalPageQuery.Dispatch(ctx, listFbPagesActiveQuery); err != nil {
		return nil, err
	}

	// key externalID
	mapFbPageActive := make(map[string]*fbpaging.FbExternalPage)
	for _, fbPage := range listFbPagesActiveQuery.Result {
		mapFbPageActive[fbPage.ExternalID] = fbPage
	}

	createFbUserCombinedCmd := &fbusering.CreateFbExternalUserCombinedCommand{
		FbUser: &fbusering.CreateFbExternalUserArgs{
			ExternalID: me.ID,
			ExternalInfo: &fbusering.FbExternalUserInfo{
				Name:      me.Name,
				FirstName: me.FirstName,
				LastName:  me.LastName,
				ShortName: me.ShortName,
				ImageURL:  me.Picture.Data.Url,
			},
			Status: status3.P,
		},
		FbUserInternal: &fbusering.CreateFbExternalUserInternalArgs{
			ExternalID: me.ID,
			Token:      longLivedAccessToken.AccessToken,
			ExpiresIn:  fbclient.ExpiresInUserToken, // 60 days
		},
	}
	if err := s.FBExternalUserAggr.Dispatch(ctx, createFbUserCombinedCmd); err != nil {
		return nil, err
	}

	var fbErrorPages []*fabo.FbErrorPage

	permissionsGranted := getPermissionsGranted(accounts.Permissions)

	listCreateFbPageCombinedCmd := make([]*fbpaging.CreateFbExternalPageCombinedArgs, 0, len(accounts.Accounts.Data))
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

		if fbPage, ok := mapFbPageActive[account.Id]; ok && fbPage.ShopID != shopID {
			fbErrorPages = append(fbErrorPages, &fabo.FbErrorPage{
				ExternalID:       account.Id,
				ExternalName:     account.Name,
				ExternalImageURL: account.Picture.Data.Url,
				Reason:           "Fanpage đã được kết nối với tài khoản trong hệ thống.",
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
		createFbPageCmd := &fbpaging.CreateFbExternalPageArgs{
			ID:                   fbPageID,
			ExternalID:           account.Id,
			ShopID:               shopID,
			ExternalName:         account.Name,
			ExternalCategory:     account.Category,
			ExternalCategoryList: categories,
			ExternalTasks:        account.Tasks,
			ExternalPermissions:  permissionsGranted,
			ExternalImageURL:     account.Picture.Data.Url,
			Status:               status3.P,
			ConnectionStatus:     status3.P,
		}
		createFbPageInternalCmd := &fbpaging.CreateFbExternalPageInternalArgs{
			ID:         fbPageID,
			ExternalID: account.Id,
			Token:      account.AccessToken,
		}
		listCreateFbPageCombinedCmd = append(listCreateFbPageCombinedCmd, &fbpaging.CreateFbExternalPageCombinedArgs{
			FbPage:         createFbPageCmd,
			FbPageInternal: createFbPageInternalCmd,
		})
	}
	var fbPageCombinedsResult []*fabo.FbPageCombined

	if len(listCreateFbPageCombinedCmd) > 0 {
		createFbExternalPageCombinedsCmd := &fbpaging.CreateFbExternalPageCombinedsCommand{
			FbPageCombineds: listCreateFbPageCombinedCmd,
		}
		if err := s.FBExternalPageAggr.Dispatch(ctx, createFbExternalPageCombinedsCmd); err != nil {
			return nil, err
		}

		fbPageCombinedsResult = convertpb.PbFbPageCombineds(createFbExternalPageCombinedsCmd.Result)
	}

	resp := &fabo.ConnectPagesResponse{
		FbUser:       convertpb.PbFbUserCombined(createFbUserCombinedCmd.Result),
		FbPages:      fbPageCombinedsResult,
		FbErrorPages: fbErrorPages,
	}
	return resp, nil
}

func verifyScopes(appScopes map[string]string, scopes []string) error {
	mapScope := make(map[string]bool)
	for _, scope := range scopes {
		mapScope[scope] = true
	}

	var permissionsMissing []string
	for scope := range appScopes {
		if _, ok := mapScope[scope]; !ok {
			permissionsMissing = append(permissionsMissing, scope)
		}
	}

	if len(permissionsMissing) > 0 {
		listPermissions := strings.Join(permissionsMissing, ",")
		dialogMsg := fmt.Sprintf("You must grant permission (%s) to perform this action.", listPermissions)
		if len(permissionsMissing) > 1 {
			dialogMsg = fmt.Sprintf("You must grant permissions (%s) to perform this action.", listPermissions)
		}
		return cm.Errorf(cm.FacebookPermissionMissing, nil, "Missing permissions").
			WithMeta("require_permissions", strings.Join(permissionsMissing, ",")).
			WithMeta("dialog_msg", dialogMsg)
	}

	return nil
}

func (s *PageService) CheckPermissions(ctx context.Context, req *fabo.CheckPagePermissionsRequest) (*fabo.CheckPagePermissionsResponse, error) {
	pageMissingRoles := &fabo.CheckPagePermissionsResponse{
		PageMissingRoles: map[string][]string{},
	}

	for _, pageID := range req.ExternalPageIDS {
		getAccessTokenQuery := &fbpaging.GetPageAccessTokenQuery{
			ExternalID: pageID,
		}
		if err := s.FBExternalPageQuery.Dispatch(ctx, getAccessTokenQuery); err != nil {
			return nil, err
		}

		pageToken, err := s.FBClient.CallAPICheckAccessToken(getAccessTokenQuery.Result)
		if err != nil {
			return nil, err
		}

		missingRoles := getMissingPermissions(pageToken.Data.Scopes)
		pageMissingRoles.PageMissingRoles[pageID] = missingRoles
	}
	return pageMissingRoles, nil
}

func getMissingPermissions(pagePerms []string) []string {
	mapPagePerms := make(map[string]struct{})
	for _, perm := range pagePerms {
		mapPagePerms[perm] = struct{}{}
	}

	missingRoles := []string{}
	for perm, _ := range appScopes {
		if _, ok := mapPagePerms[perm]; !ok {
			missingRoles = append(missingRoles, perm)
		}
	}

	return missingRoles
}

func getPermissionsGranted(permissionsData model.AccountsPermissions) []string {
	var permissions []string
	for _, permission := range permissionsData.Data {
		if permission.Status == fbclient.PermissionGranted {
			permissions = append(permissions, permission.Permission)
		}
	}
	return permissions
}

func contains(arr []string, str string) bool {
	for _, el := range arr {
		if el == str {
			return true
		}
	}

	return false
}
