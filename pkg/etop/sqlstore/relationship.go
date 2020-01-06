package sqlstore

import (
	"context"
	"strings"
	"sync"
	"time"

	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		GetAccountUser,
		GetAccountUserExtended,
		GetAccountUserExtendeds,
		GetAllAccountRoles,
		UpdateRole,
		UpdateInfos,
		CreateAccountUser,
		UpdateAccountUser,
		GetAllAccountUsers,
		DeleteAccountUser,
	)
}

var filterAccountUserWhitelist = FilterWhitelist{
	Arrays:   nil,
	Bools:    nil,
	Contains: nil,
	Dates:    nil,
	Equals:   nil,
	Nullable: []string{"deleted_at"},
	Numbers:  nil,
	Status:   []string{"status"},
	Unaccent: nil,
	PrefixOrRename: map[string]string{
		"status":     `"au".status`,
		"deleted_at": `"au".deleted_at`,
	},
}

func GetAllAccountRoles(ctx context.Context, query *identitymodelx.GetAllAccountRolesQuery) error {
	if query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}

	s := x.Table("account_user").
		Where("au.user_id = ? AND au.deleted_at is NULL", query.UserID)

	if query.Type.Valid {
		s = s.Where("type = ?", query.Type)
	}
	return s.Find((*identitymodel.AccountUserExtendeds)(&query.Result))
}

func UpdateInfos(ctx context.Context, cmd *identitymodelx.UpdateInfosCommand) error {
	return inTransaction(func(s Qx) error {
		return updateInfos(ctx, s, cmd)
	})
}

func updateInfos(ctx context.Context, s Qx, cmd *identitymodelx.UpdateInfosCommand) error {
	mapUpdate := make(map[string]interface{})
	if cmd.ShortName.Valid {
		mapUpdate["short_name"] = cmd.ShortName.String
	}
	if cmd.FullName.Valid {
		mapUpdate["full_name"] = cmd.FullName.String
	}
	if cmd.Position.Valid {
		mapUpdate["position"] = cmd.Position.String
	}
	if _, err := s.Table("account_user").
		Where("account_id = ?", cmd.AccountID).
		Where("user_id = ?", cmd.UserID).
		Where("deleted_at is NULL").
		UpdateMap(mapUpdate); err != nil {
		return err
	}

	cmd.Result = new(identitymodel.AccountUser)
	s = x.
		Where("deleted_at is NULL").
		Where("account_id = ?", cmd.AccountID).
		Where("user_id = ?", cmd.UserID)

	return s.ShouldGet(cmd.Result)
}

func UpdateRole(ctx context.Context, cmd *identitymodelx.UpdateRoleCommand) error {
	return inTransaction(func(s Qx) error {
		return updateRole(ctx, s, cmd)
	})
}

func updateRole(ctx context.Context, s Qx, cmd *identitymodelx.UpdateRoleCommand) error {
	permission := &identitymodel.AccountUser{
		AccountID:  cmd.AccountID,
		UserID:     cmd.UserID,
		Permission: cmd.Permission,
	}
	roles := "{" + strings.Join(cmd.Permission.Roles, ",") + "}"
	permissions := "{" + strings.Join(cmd.Permission.Permissions, ",") + "}"
	_, err := s.Table("account_user").
		Where("account_id = ?", cmd.AccountID).
		Where("user_id = ?", cmd.UserID).
		UpdateMap(map[string]interface{}{
			"roles":       roles,
			"permissions": permissions,
		})
	cmd.Result = permission
	return err
}

func GetAccountUser(ctx context.Context, query *identitymodelx.GetAccountUserQuery) error {
	if query.UserID == 0 && !query.FindByAccountID {
		return cm.Error(cm.InvalidArgument, "Missing UserID", nil)
	}
	if query.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}

	query.Result = new(identitymodel.AccountUser)
	s := x.
		Where("deleted_at is NULL").
		Where("account_id = ?", query.AccountID)
	if query.UserID != 0 && !query.FindByAccountID {
		s = s.Where("user_id = ?", query.UserID)
	}
	return s.ShouldGet(query.Result)
}

func GetAccountUserExtended(ctx context.Context, query *identitymodelx.GetAccountUserExtendedQuery) error {
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

func GetAccountUserExtendeds(ctx context.Context, query *identitymodelx.GetAccountUserExtendedsQuery) error {
	if len(query.AccountIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountIDs", nil)
	}

	s := x.Table("account_user").
		In("au.account_id", query.AccountIDs)
	if !query.IncludeDeleted {
		s = s.Where("au.deleted_at IS NULL")
	}

	s, _, err := Filters(s, query.Filters, filterAccountUserWhitelist)
	if err != nil {
		return err
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, sqlstore.ConvertPaging(query.Paging), Ms{"id": "u.id", "updated_at": "au.updated_at"})
		if err != nil {
			return err
		}
		if err := s2.Find((*identitymodel.AccountUserExtendeds)(&query.Result.AccountUsers)); err != nil {
			return err
		}
	}
	return nil
}

func CreateAccountUser(ctx context.Context, cmd *identitymodelx.CreateAccountUserCommand) error {
	accUser := cmd.AccountUser
	if accUser.UserID == 0 || accUser.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	err := x.Table("account_user").ShouldInsert(accUser)
	if err != nil {
		return err
	}

	cmd.Result = new(identitymodel.AccountUser)
	err = x.Table("account_user").
		Where("account_id = ? AND user_id = ?", accUser.AccountID, accUser.UserID).
		ShouldGet(cmd.Result)
	return err
}

func UpdateAccountUser(ctx context.Context, cmd *identitymodelx.UpdateAccountUserCommand) error {
	accUser := cmd.AccountUser
	if accUser.UserID == 0 || accUser.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	return x.Table("account_user").
		Where("user_id = ?", accUser.UserID).
		Where("account_id = ?", accUser.AccountID).
		ShouldUpdate(accUser)
}

func GetAllAccountUsers(ctx context.Context, query *identitymodelx.GetAllAccountUsersQuery) error {
	if len(query.UserIDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Missing UserIDs", nil)
	}
	var res []*identitymodel.AccountUser
	guard := make(chan int, 8)
	var m sync.Mutex
	for i, userID := range query.UserIDs {
		guard <- i
		go func(uID dot.ID) {
			defer func() {
				<-guard
			}()
			var _res []*identitymodel.AccountUser
			s := x.Table("account_user").
				Where("au.user_id = ? AND au.deleted_at is NULL", uID)
			if query.Type.Valid {
				s = s.Where("type = ?", query.Type)
			}
			if err := s.Find((*identitymodel.AccountUsers)(&_res)); err == nil {
				m.Lock()
				res = append(res, _res...)
				m.Unlock()
			}
		}(userID)
	}

	query.Result = res
	return nil
}

func DeleteAccountUser(ctx context.Context, cmd *identitymodelx.DeleteAccountUserCommand) error {
	if cmd.UserID == 0 || cmd.AccountID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}
	updated, err := x.Table("account_user").
		Where("account_id = ?", cmd.AccountID).
		Where("user_id = ?", cmd.UserID).
		UpdateMap(map[string]interface{}{
			"deleted_at": time.Now(),
			"status":     int(status3.N),
		})
	if err != nil {
		return err
	}

	cmd.Result.Updated = updated
	return nil
}
