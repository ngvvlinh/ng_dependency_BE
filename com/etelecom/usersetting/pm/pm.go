package pm

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/etelecom/usersetting"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/charge_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	userSettingQuery usersetting.QueryBus
	identityQuery    identity.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	userSettingQ usersetting.QueryBus,
	identityQ identity.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		userSettingQuery: userSettingQ,
		identityQuery:    identityQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ExtensionCreating)
}

func (m *ProcessManager) ExtensionCreating(ctx context.Context, event *etelecom.ExtensionCreatingEvent) error {
	ownerID := event.OwnerID
	if ownerID == 0 {
		shopQuery := &identity.GetShopByIDQuery{
			ID: event.AccountID,
		}
		if err := m.identityQuery.Dispatch(ctx, shopQuery); err != nil {
			return err
		}
		ownerID = shopQuery.Result.OwnerID
	}

	// Luon cho phép tạo ext cho chủ shop
	if ownerID == event.UserID {
		return nil
	}

	query := &usersetting.GetUserSettingQuery{
		UserID: event.OwnerID,
	}
	if err := m.userSettingQuery.Dispatch(ctx, query); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Không thể sử dụng api này để tạo extension")
	}
	setting := query.Result

	// Cho phép tạo ext đối với user có setting: miễn phí hoặc trả sau
	// Trường hợp trả trước phải tạo thông qua subscription
	if setting.ExtensionChargeType == charge_type.Prepaid {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn gói dịch vụ để tạo extension")
	}
	return nil
}
