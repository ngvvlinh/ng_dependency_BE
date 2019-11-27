package sqlstore

import (
	"context"
	"sync"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		GetAccountUser,
		GetAccountUserExtended,
		GetAccountUserExtendeds,
		GetAllAccountRoles,
		UpdateRole,
		CreateAccountUser,
		UpdateAccountUser,
		GetAllAccountUsers,
	)
}

var filterAccountUserWhitelist = FilterWhitelist{}

func GetAllAccountRoles(ctx context.Context, query *model.GetAllAccountRolesQuery) error {
	if query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	s := x.Table("account_user").
		Where("au.user_id = ? AND au.deleted_at is NULL", query.UserID)

	if query.Type != "" {
		s = s.Where("type = ?", query.Type)
	}
	return s.Find((*model.AccountUserExtendeds)(&query.Result))
}

func UpdateRole(ctx context.Context, cmd *model.UpdateRoleCommand) error {
	return inTransaction(func(s Qx) error {
		return updateRole(ctx, s, cmd)
	})
}

func updateRole(ctx context.Context, s Qx, cmd *model.UpdateRoleCommand) error {
	permission := &model.AccountUser{
		AccountID:  cmd.AccountID,
		UserID:     cmd.UserID,
		Permission: cmd.Permission,
	}
	_, err := s.Insert(permission)
	cmd.Result = permission
	return err
}

func GetAccountUser(ctx context.Context, query *model.GetAccountUserQuery) error {
	if query.UserID == 0 && !query.FindByAccountID {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}
	if query.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	query.Result = new(model.AccountUser)
	s := x.
		Where("deleted_at is NULL").
		Where("account_id = ?", query.AccountID)
	if query.UserID != 0 && !query.FindByAccountID {
		s = s.Where("user_id = ?", query.UserID)
	}
	return s.ShouldGet(query.Result)
}

func GetAccountUserExtended(ctx context.Context, query *model.GetAccountUserExtendedQuery) error {
	if query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}
	if query.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	return x.
		Where("au.deleted_at is NULL").
		Where("au.account_id = ?", query.AccountID).
		Where("au.user_id = ?", query.UserID).
		ShouldGet(&query.Result)
}

func GetAccountUserExtendeds(ctx context.Context, query *model.GetAccountUserExtendedsQuery) error {
	if len(query.AccountIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountIDs", nil)
	}

	s := x.Table("account_user").
		In("au.account_id", query.AccountIDs).
		Where("au.deleted_at IS NULL")

	s, _, err := Filters(s, query.Filters, filterAccountUserWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "u.id", "updated_at": "au.updated_at"})
		if err != nil {
			return err
		}
		if err := s2.Find((*model.AccountUserExtendeds)(&query.Result.AccountUsers)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.AccountUserExtendeds{})
		if err != nil {
			return nil
		}
		query.Result.Total = int(total)
	}
	return nil
}

func CreateAccountUser(ctx context.Context, cmd *model.CreateAccountUserCommand) error {
	accUser := cmd.AccountUser
	if accUser.UserID == 0 || accUser.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	err := x.Table("account_user").ShouldInsert(accUser)
	if err != nil {
		return err
	}

	cmd.Result = new(model.AccountUser)
	err = x.Table("account_user").
		Where("account_id = ? AND user_id = ?", accUser.AccountID, accUser.UserID).
		ShouldGet(cmd.Result)
	return err
}

func UpdateAccountUser(ctx context.Context, cmd *model.UpdateAccountUserCommand) error {
	accUser := cmd.AccountUser
	if accUser.UserID == 0 || accUser.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	return x.Table("account_user").
		Where("user_id = ?", accUser.UserID).
		Where("account_id = ?", accUser.AccountID).
		ShouldUpdate(accUser)
}

func GetAllAccountUsers(ctx context.Context, query *model.GetAllAccountUsersQuery) error {
	if len(query.UserIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserIDs", nil)
	}
	var res []*model.AccountUser
	guard := make(chan int, 8)
	var m sync.Mutex
	for i, userID := range query.UserIDs {
		guard <- i
		go func(uID dot.ID) {
			defer func() {
				<-guard
			}()
			var _res []*model.AccountUser
			s := x.Table("account_user").
				Where("au.user_id = ? AND au.deleted_at is NULL", uID)
			if query.Type != "" {
				s = s.Where("type = ?", query.Type)
			}
			if err := s.Find((*model.AccountUsers)(&_res)); err == nil {
				m.Lock()
				res = append(res, _res...)
				m.Unlock()
			}
		}(userID)
	}

	query.Result = res
	return nil
}
