package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	"o.o/api/top/types/etc/account_type"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

type UserService struct {
	IdentityQuery identity.QueryBus
}

func (s *UserService) Clone() *UserService {
	res := *s
	return &res
}

func (s *UserService) GetUsers(ctx context.Context, q *GetUsersEndpoint) error {
	paging, err := cmapi.CMCursorPaging(q.Paging)
	if err != nil {
		return err
	}
	query := &identity.GetUsersQuery{
		Name:      q.Filters.Name,
		Phone:     q.Filters.Phone,
		Email:     q.Filters.Email,
		CreatedAt: q.Filters.CreatedAt,
		Paging:    *paging,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	Users := query.Result
	var UserIDs []dot.ID
	for _, user := range query.Result.ListUsers {
		UserIDs = append(UserIDs, user.ID)
	}
	queryAccount := &identity.GetAllAccountsByUsersQuery{
		UserIDs: UserIDs,
		Type:    account_type.Shop.Wrap(),
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryAccount); err != nil {
		return err
	}
	q.Result = &admin.UserResponse{
		Paging: cmapi.PbCursorPageInfo(paging, &Users.Paging),
		Users:  convertpb.PbUsers(Users.ListUsers),
	}
	populateShopCount(q.Result.Users, queryAccount.Result)
	return nil
}

func (s *UserService) GetUser(ctx context.Context, q *GetUserEndpoint) error {
	query := &identity.GetUserByIDQuery{
		UserID: q.Id,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	queryAccount := &identity.GetAllAccountsByUsersQuery{
		UserIDs: []dot.ID{query.Result.ID},
		Type:    account_type.Shop.Wrap(),
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryAccount); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_User_To_api_User(query.Result)
	populateShopCount([]*etop.User{q.Result}, queryAccount.Result)
	return nil
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

func (s *UserService) GetUsersByIDs(ctx context.Context, q *GetUsersByIDsEndpoint) error {
	query := &identity.GetUsersByIDsQuery{
		IDs: q.Ids,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
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
		return err
	}
	q.Result = &admin.UserResponse{
		Users: convertpb.PbUsers(query.Result),
	}
	populateShopCount(q.Result.Users, queryAccount.Result)
	return nil
}
