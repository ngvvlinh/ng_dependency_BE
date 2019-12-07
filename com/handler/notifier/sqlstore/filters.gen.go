// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/notifier_entity"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/capi/dot"
)

type NotificationFilters struct{ prefix string }

func NewNotificationFilters(prefix string) NotificationFilters {
	return NotificationFilters{prefix}
}

func (ft *NotificationFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft NotificationFilters) Prefix() string {
	return ft.prefix
}

func (ft *NotificationFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *NotificationFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *NotificationFilters) ByTitle(Title string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == "",
	}
}

func (ft *NotificationFilters) ByTitlePtr(Title *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == nil,
		IsZero: Title != nil && (*Title) == "",
	}
}

func (ft *NotificationFilters) ByMessage(Message string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "message",
		Value:  Message,
		IsNil:  Message == "",
	}
}

func (ft *NotificationFilters) ByMessagePtr(Message *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "message",
		Value:  Message,
		IsNil:  Message == nil,
		IsZero: Message != nil && (*Message) == "",
	}
}

func (ft *NotificationFilters) ByIsRead(IsRead bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_read",
		Value:  IsRead,
		IsNil:  bool(!IsRead),
	}
}

func (ft *NotificationFilters) ByIsReadPtr(IsRead *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_read",
		Value:  IsRead,
		IsNil:  IsRead == nil,
		IsZero: IsRead != nil && bool(!(*IsRead)),
	}
}

func (ft *NotificationFilters) ByEntityID(EntityID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "entity_id",
		Value:  EntityID,
		IsNil:  EntityID == 0,
	}
}

func (ft *NotificationFilters) ByEntityIDPtr(EntityID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "entity_id",
		Value:  EntityID,
		IsNil:  EntityID == nil,
		IsZero: EntityID != nil && (*EntityID) == 0,
	}
}

func (ft *NotificationFilters) ByEntity(Entity notifier_entity.NotifierEntity) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "entity",
		Value:  Entity,
		IsNil:  Entity == 0,
	}
}

func (ft *NotificationFilters) ByEntityPtr(Entity *notifier_entity.NotifierEntity) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "entity",
		Value:  Entity,
		IsNil:  Entity == nil,
		IsZero: Entity != nil && (*Entity) == 0,
	}
}

func (ft *NotificationFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *NotificationFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *NotificationFilters) BySyncStatus(SyncStatus status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "sync_status",
		Value:  SyncStatus,
		IsNil:  SyncStatus == 0,
	}
}

func (ft *NotificationFilters) BySyncStatusPtr(SyncStatus *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "sync_status",
		Value:  SyncStatus,
		IsNil:  SyncStatus == nil,
		IsZero: SyncStatus != nil && (*SyncStatus) == 0,
	}
}

func (ft *NotificationFilters) ByExternalServiceID(ExternalServiceID int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_service_id",
		Value:  ExternalServiceID,
		IsNil:  ExternalServiceID == 0,
	}
}

func (ft *NotificationFilters) ByExternalServiceIDPtr(ExternalServiceID *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_service_id",
		Value:  ExternalServiceID,
		IsNil:  ExternalServiceID == nil,
		IsZero: ExternalServiceID != nil && (*ExternalServiceID) == 0,
	}
}

func (ft *NotificationFilters) ByExternalNotiID(ExternalNotiID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_noti_id",
		Value:  ExternalNotiID,
		IsNil:  ExternalNotiID == "",
	}
}

func (ft *NotificationFilters) ByExternalNotiIDPtr(ExternalNotiID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_noti_id",
		Value:  ExternalNotiID,
		IsNil:  ExternalNotiID == nil,
		IsZero: ExternalNotiID != nil && (*ExternalNotiID) == "",
	}
}

func (ft *NotificationFilters) BySendNotification(SendNotification bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "send_notification",
		Value:  SendNotification,
		IsNil:  bool(!SendNotification),
	}
}

func (ft *NotificationFilters) BySendNotificationPtr(SendNotification *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "send_notification",
		Value:  SendNotification,
		IsNil:  SendNotification == nil,
		IsZero: SendNotification != nil && bool(!(*SendNotification)),
	}
}

func (ft *NotificationFilters) BySyncedAt(SyncedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "synced_at",
		Value:  SyncedAt,
		IsNil:  SyncedAt.IsZero(),
	}
}

func (ft *NotificationFilters) BySyncedAtPtr(SyncedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "synced_at",
		Value:  SyncedAt,
		IsNil:  SyncedAt == nil,
		IsZero: SyncedAt != nil && (*SyncedAt).IsZero(),
	}
}

func (ft *NotificationFilters) BySeenAt(SeenAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "seen_at",
		Value:  SeenAt,
		IsNil:  SeenAt.IsZero(),
	}
}

func (ft *NotificationFilters) BySeenAtPtr(SeenAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "seen_at",
		Value:  SeenAt,
		IsNil:  SeenAt == nil,
		IsZero: SeenAt != nil && (*SeenAt).IsZero(),
	}
}

func (ft *NotificationFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *NotificationFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *NotificationFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *NotificationFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type DeviceFilters struct{ prefix string }

func NewDeviceFilters(prefix string) DeviceFilters {
	return DeviceFilters{prefix}
}

func (ft *DeviceFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft DeviceFilters) Prefix() string {
	return ft.prefix
}

func (ft *DeviceFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *DeviceFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *DeviceFilters) ByDeviceID(DeviceID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "device_id",
		Value:  DeviceID,
		IsNil:  DeviceID == "",
	}
}

func (ft *DeviceFilters) ByDeviceIDPtr(DeviceID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "device_id",
		Value:  DeviceID,
		IsNil:  DeviceID == nil,
		IsZero: DeviceID != nil && (*DeviceID) == "",
	}
}

func (ft *DeviceFilters) ByDeviceName(DeviceName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "device_name",
		Value:  DeviceName,
		IsNil:  DeviceName == "",
	}
}

func (ft *DeviceFilters) ByDeviceNamePtr(DeviceName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "device_name",
		Value:  DeviceName,
		IsNil:  DeviceName == nil,
		IsZero: DeviceName != nil && (*DeviceName) == "",
	}
}

func (ft *DeviceFilters) ByExternalDeviceID(ExternalDeviceID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_device_id",
		Value:  ExternalDeviceID,
		IsNil:  ExternalDeviceID == "",
	}
}

func (ft *DeviceFilters) ByExternalDeviceIDPtr(ExternalDeviceID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_device_id",
		Value:  ExternalDeviceID,
		IsNil:  ExternalDeviceID == nil,
		IsZero: ExternalDeviceID != nil && (*ExternalDeviceID) == "",
	}
}

func (ft *DeviceFilters) ByExternalServiceID(ExternalServiceID int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_service_id",
		Value:  ExternalServiceID,
		IsNil:  ExternalServiceID == 0,
	}
}

func (ft *DeviceFilters) ByExternalServiceIDPtr(ExternalServiceID *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_service_id",
		Value:  ExternalServiceID,
		IsNil:  ExternalServiceID == nil,
		IsZero: ExternalServiceID != nil && (*ExternalServiceID) == 0,
	}
}

func (ft *DeviceFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *DeviceFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *DeviceFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *DeviceFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *DeviceFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *DeviceFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *DeviceFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *DeviceFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *DeviceFilters) ByDeactivatedAt(DeactivatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deactivated_at",
		Value:  DeactivatedAt,
		IsNil:  DeactivatedAt.IsZero(),
	}
}

func (ft *DeviceFilters) ByDeactivatedAtPtr(DeactivatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deactivated_at",
		Value:  DeactivatedAt,
		IsNil:  DeactivatedAt == nil,
		IsZero: DeactivatedAt != nil && (*DeactivatedAt).IsZero(),
	}
}
