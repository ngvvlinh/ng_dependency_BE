package pm

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	cm "o.o/backend/pkg/common"
)

func (m *ProcessManager) ExtensionCreatingEvent(ctx context.Context, event *etelecom.ExtensionCreatingEvent) error {
	query := &identity.GetAccountUserQuery{
		UserID:    event.UserID,
		AccountID: event.AccountID,
	}
	if err := m.identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	roles := query.Result.Roles
	acceptRoles := []string{
		authorization.RoleShopOwner.String(),
		authorization.RoleTelecomCustomerService.String(),
	}
	isPermision := false
	for _, role := range roles {
		if cm.StringsContain(acceptRoles, role) {
			isPermision = true
			break
		}
	}
	if !isPermision {
		return cm.Errorf(cm.FailedPrecondition, nil, "Chỉ được tạo extension cho nhân viên có quyền CSKH")
	}
	return nil
}
