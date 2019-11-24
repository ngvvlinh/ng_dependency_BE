package permission

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("session", VerifyAdmin)
}

type VerifyAdminQuery struct {
	UserID dot.ID
	Result *model.Permission
}

func VerifyAdmin(ctx context.Context, query *VerifyAdminQuery) error {
	adminQuery := &model.GetAccountRolesQuery{
		AccountID: model.EtopAccountID,
		UserID:    query.UserID,
	}
	if err := bus.Dispatch(ctx, adminQuery); err != nil {
		return cm.MapError(err).
			Mapf(cm.NotFound, cm.PermissionDenied, "Người dùng phải là quản trị viên để sử dụng tính năng này").
			Throw()
	}

	query.Result = &adminQuery.Result.AccountUser.Permission
	return nil
}

func VerifyRoleLevel(requiredRole, userRole RoleType) bool {
	return RoleLevel(userRole) >= RoleLevel(requiredRole)
}
