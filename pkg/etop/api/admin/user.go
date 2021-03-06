package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type UserService struct {
	session.Session

	IdentityQuery identity.QueryBus
	IdentityAggr  identity.CommandBus
	RedisStore    redis.Store
}

func (s *UserService) Clone() admin.UserService {
	res := *s
	return &res
}

func (s *UserService) BlockUser(ctx context.Context, q *admin.BlockUserRequest) (*etop.User, error) {
	blockBy := s.SS.User().ID
	cmd := &identity.BlockUserCommand{
		UserID:      q.UserID,
		BlockBy:     blockBy,
		BlockReason: q.BlockReason,
	}
	err := s.IdentityAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_User_To_api_User(cmd.Result), nil
}

func (s *UserService) UnblockUser(ctx context.Context, q *admin.UnblockUserRequest) (*etop.User, error) {
	cmd := &identity.UnblockUserCommand{
		UserID: q.UserID,
	}
	err := s.IdentityAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_User_To_api_User(cmd.Result), nil
}

func (s *UserService) GetUsers(ctx context.Context, q *admin.GetUsersRequest) (*admin.UserResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	if q.Filters == nil {
		q.Filters = &admin.UsersFilter{}
	}

	var fullTextSearch filter.FullTextSearch = ""
	if q.Filters != nil {
		fullTextSearch = q.Filters.Name
	}

	query := &identity.GetUserFtRefSaffsQuery{
		Name:      fullTextSearch,
		Phone:     q.Filters.Phone,
		Email:     q.Filters.Email,
		CreatedAt: q.Filters.CreatedAt,
		RefSale:   q.Filters.RefSale,
		RefAff:    q.Filters.RefAff,

		Paging: *paging,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	Users := query.Result
	if len(Users.ListUsers) > 0 {
		var UserIDs []dot.ID
		for _, user := range query.Result.ListUsers {
			UserIDs = append(UserIDs, user.ID)
		}
		queryAccount := &identity.GetAllAccountsByUsersQuery{
			UserIDs: UserIDs,
			Type:    account_type.Shop.Wrap(),
		}
		if err := s.IdentityQuery.Dispatch(ctx, queryAccount); err != nil {
			return nil, err
		}
		result := &admin.UserResponse{
			Paging: cmapi.PbCursorPageInfo(paging, &Users.Paging),
			Users:  convertpb.PbUserFtRefSaffs(Users.ListUsers),
		}
		populateShopCount(result.Users, queryAccount.Result)
		return result, nil
	}
	return &admin.UserResponse{
		Paging: cmapi.PbCursorPageInfo(paging, &Users.Paging),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, q *pbcm.IDRequest) (*etop.User, error) {
	query := &identity.GetUserFtRefSaffByIDQuery{
		UserID: q.Id,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	queryAccount := &identity.GetAllAccountsByUsersQuery{
		UserIDs: []dot.ID{query.Result.ID},
		Type:    account_type.Shop.Wrap(),
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryAccount); err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_UserFtRefSaff_To_api_User(query.Result)
	populateShopCount([]*etop.User{result}, queryAccount.Result)
	return result, nil
}

func populateShopCount(Users []*etop.User, Accounts []*identity.AccountUser) {
	result := make(map[dot.ID][]*identity.AccountUser)
	for _, account := range Accounts {
		result[account.UserID] = append(result[account.UserID], account)
	}
	for _, user := range Users {
		if result[user.Id] != nil {
			user.TotalShop = len(result[user.Id])
		}
	}
}

func (s *UserService) GetUsersByIDs(ctx context.Context, q *pbcm.IDsRequest) (*admin.UserResponse, error) {
	query := &identity.GetUsersByIDsQuery{
		IDs: q.Ids,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var UserIDs []dot.ID
	for _, user := range query.Result {
		UserIDs = append(UserIDs, user.ID)
	}
	queryAccount := &identity.GetAllAccountsByUsersQuery{
		UserIDs: UserIDs,
		Type:    account_type.Shop.Wrap(),
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryAccount); err != nil {
		return nil, err
	}
	result := &admin.UserResponse{
		Users: convertpb.PbUsers(query.Result),
	}
	populateShopCount(result.Users, queryAccount.Result)
	return result, nil
}

func (s *UserService) UpdateUserRef(ctx context.Context, r *admin.UpdateUserRefRequest) (*pbcm.Empty, error) {
	cmdUpdate := &identity.UpdateUserRefCommand{
		UserID:  r.UserID,
		RefAff:  r.RefAff,
		RefSale: r.RefSale,
	}
	err := s.IdentityAggr.Dispatch(ctx, cmdUpdate)
	if err != nil {
		return nil, err
	}
	return &pbcm.Empty{}, nil
}

func (s *UserService) ChangeUserCredential(ctx context.Context, r *admin.ChangeUserCredentialRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &identity.ChangeUserCredentialCommand{
		UserID:   r.UserID,
		Email:    r.Email,
		Phone:    r.Phone,
		Password: r.Password,
	}
	err := s.IdentityAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *UserService) GetLatestUserOTP(ctx context.Context, r *admin.GetLatestUserOTPRequest) (*admin.GetLatestUserOTPResponse, error) {
	if r.UserID == 0 && r.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing user ID or phone number")
	}
	if r.UserID != 0 && r.Phone != "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Only provide user ID or phone number")
	}
	redisKey := identity.GetLatestUserOTPRedisKey(r.UserID, r.Phone)
	res := identity.LatestUserOTPData{}
	if err := s.RedisStore.Get(redisKey, &res); err != nil && err != redis.ErrNil {
		return nil, err
	}
	return &admin.GetLatestUserOTPResponse{
		OTP:    res.OTP,
		Action: res.Action,
		Label:  res.Action.GetLabelRefName(),
	}, nil
}
