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
	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
)

type PageService struct {
	session.Sessioner
	ss *session.Session

	faboInfo            *faboinfo.FaboInfo
	fbExternalUserQuery fbusering.QueryBus
	fbExternalUserAggr  fbusering.CommandBus
	fbExternalPageQuery fbpaging.QueryBus
	fbExternalPageAggr  fbpaging.CommandBus
	appScopes           map[string]string
	fbClient            *fbclient.FbClient
}

func NewPageService(
	ss *session.Session,
	faboInfo *faboinfo.FaboInfo,
	fbUserQuery fbusering.QueryBus,
	fbUserAggr fbusering.CommandBus,
	fbPageQuery fbpaging.QueryBus,
	fbPageAggr fbpaging.CommandBus,
	appScopes map[string]string,
	fbClient *fbclient.FbClient,
) *PageService {
	s := &PageService{
		ss:                  ss,
		faboInfo:            faboInfo,
		fbExternalUserQuery: fbUserQuery,
		fbExternalUserAggr:  fbUserAggr,
		fbExternalPageQuery: fbPageQuery,
		fbExternalPageAggr:  fbPageAggr,
		appScopes:           appScopes,
		fbClient:            fbClient,
	}
	return s
}

func (s *PageService) Clone() fabo.PageService {
	res := *s
	res.Sessioner, res.ss = s.ss.Split()
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
		ShopID:      s.ss.Shop().ID,
	}
	if err := s.fbExternalPageAggr.Dispatch(ctx, disablePagesByIDsCmd); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *PageService) ListPages(ctx context.Context, r *fabo.ListPagesRequest) (*fabo.ListPagesResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	listFbExternalPagesQuery := &fbpaging.ListFbExternalPagesQuery{
		ShopID:  s.ss.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.fbExternalPageQuery.Dispatch(ctx, listFbExternalPagesQuery); err != nil {
		return nil, err
	}
	resp := &fabo.ListPagesResponse{
		FbPages: convertpb.PbFbPages(listFbExternalPagesQuery.Result.FbPages),
		Paging:  cmapi.PbPageInfo(paging),
	}
	return resp, nil
}

func (s *PageService) ConnectPages(ctx context.Context, r *fabo.ConnectPagesRequest) (*fabo.ConnectPagesResponse, error) {
	shopID := s.ss.Shop().ID

	// Check accessToken is alive
	_, err := s.fbClient.CallAPICheckAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}

	// verify permissions
	//if err := verifyScopes(s.appScopes, userToken.Data.Scopes); err != nil {
	//	return nil, err
	//}

	// Get long lived accessToken from accessToken (above)
	if r.AccessToken == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "access_token must not be null")
	}
	longLivedAccessToken, err := s.fbClient.CallAPIGetLongLivedAccessToken(r.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get information of user from accessToken (above)
	me, err := s.fbClient.CallAPIGetMe(longLivedAccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get all accounts of user (above)
	accounts, err := s.fbClient.CallAPIGetAccounts(longLivedAccessToken.AccessToken)
	if err != nil {
		return nil, err
	}

	var externalIDs []string
	for _, account := range accounts.Accounts.Data {
		externalIDs = append(externalIDs, account.Id)
	}

	// Get fbPages active from externalIDs (accounts)
	listFbPagesActiveQuery := &fbpaging.ListFbExternalPagesActiveByExternalIDsQuery{
		ExternalIDs: externalIDs,
	}
	if err := s.fbExternalPageQuery.Dispatch(ctx, listFbPagesActiveQuery); err != nil {
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
			Token:  longLivedAccessToken.AccessToken,
			Status: status3.P,
		},
		FbUserInternal: &fbusering.CreateFbExternalUserInternalArgs{
			ExternalID: me.ID,
			Token:      longLivedAccessToken.AccessToken,
			ExpiresIn:  fbclient.ExpiresInUserToken, // 60 days
		},
	}
	if err := s.fbExternalUserAggr.Dispatch(ctx, createFbUserCombinedCmd); err != nil {
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
		if err := s.fbExternalPageAggr.Dispatch(ctx, createFbExternalPageCombinedsCmd); err != nil {
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

	for scope, messageScope := range appScopes {
		if _, ok := mapScope[scope]; !ok {
			return cm.Errorf(cm.FacebookPermissionDenied, nil, "Bạn chưa cấp đủ quyền để tiếp tục").
				WithMeta(fmt.Sprintf("scope.%s", scope), messageScope)
		}
	}
	return nil
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
