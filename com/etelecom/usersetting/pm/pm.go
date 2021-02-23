package pm

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/etelecom/usersetting"
	"o.o/api/top/types/etc/charge_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	userSettingQuery usersetting.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	userSettingQ usersetting.QueryBus,
) *ProcessManager {
	p := &ProcessManager{userSettingQuery: userSettingQ}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ExtensionCreating)
}

func (m *ProcessManager) ExtensionCreating(ctx context.Context, event *etelecom.ExtensionCreatingEvent) error {
	query := &usersetting.GetUserSettingQuery{
		UserID: event.OwnerID,
	}
	if err := m.userSettingQuery.Dispatch(ctx, query); err != nil && cm.ErrorCode(err) != cm.NotFound {
		return err
	}
	setting := query.Result
	if setting == nil {
		return nil
	}
	if setting.ExtensionChargeType == charge_type.Prepaid {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn gói dịch vụ để tạo extension.")
	}
	return nil
}
