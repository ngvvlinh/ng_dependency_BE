// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type ConnectionFilters struct{ prefix string }

func NewConnectionFilters(prefix string) ConnectionFilters {
	return ConnectionFilters{prefix}
}

func (ft *ConnectionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ConnectionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ConnectionFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ConnectionFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ConnectionFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ConnectionFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ConnectionFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ConnectionFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ConnectionFilters) ByPartnerID(PartnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "partner_id",
		Value:  PartnerID,
		IsNil:  PartnerID == 0,
	}
}

func (ft *ConnectionFilters) ByPartnerIDPtr(PartnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "partner_id",
		Value:  PartnerID,
		IsNil:  PartnerID == nil,
		IsZero: PartnerID != nil && (*PartnerID) == 0,
	}
}

func (ft *ConnectionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ConnectionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ConnectionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ConnectionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ConnectionFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ConnectionFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ConnectionFilters) ByDriver(Driver string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "driver",
		Value:  Driver,
		IsNil:  Driver == "",
	}
}

func (ft *ConnectionFilters) ByDriverPtr(Driver *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "driver",
		Value:  Driver,
		IsNil:  Driver == nil,
		IsZero: Driver != nil && (*Driver) == "",
	}
}

func (ft *ConnectionFilters) ByConnectionType(ConnectionType connection_type.ConnectionType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_type",
		Value:  ConnectionType,
		IsNil:  ConnectionType == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionTypePtr(ConnectionType *connection_type.ConnectionType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_type",
		Value:  ConnectionType,
		IsNil:  ConnectionType == nil,
		IsZero: ConnectionType != nil && (*ConnectionType) == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionSubtype(ConnectionSubtype connection_type.ConnectionSubtype) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_subtype",
		Value:  ConnectionSubtype,
		IsNil:  ConnectionSubtype == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionSubtypePtr(ConnectionSubtype *connection_type.ConnectionSubtype) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_subtype",
		Value:  ConnectionSubtype,
		IsNil:  ConnectionSubtype == nil,
		IsZero: ConnectionSubtype != nil && (*ConnectionSubtype) == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionMethod(ConnectionMethod connection_type.ConnectionMethod) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionMethodPtr(ConnectionMethod *connection_type.ConnectionMethod) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == nil,
		IsZero: ConnectionMethod != nil && (*ConnectionMethod) == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionProvider(ConnectionProvider connection_type.ConnectionProvider) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_provider",
		Value:  ConnectionProvider,
		IsNil:  ConnectionProvider == 0,
	}
}

func (ft *ConnectionFilters) ByConnectionProviderPtr(ConnectionProvider *connection_type.ConnectionProvider) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_provider",
		Value:  ConnectionProvider,
		IsNil:  ConnectionProvider == nil,
		IsZero: ConnectionProvider != nil && (*ConnectionProvider) == 0,
	}
}

func (ft *ConnectionFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *ConnectionFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *ConnectionFilters) ByImageURL(ImageURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageURL,
		IsNil:  ImageURL == "",
	}
}

func (ft *ConnectionFilters) ByImageURLPtr(ImageURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageURL,
		IsNil:  ImageURL == nil,
		IsZero: ImageURL != nil && (*ImageURL) == "",
	}
}

type ShopConnectionFilters struct{ prefix string }

func NewShopConnectionFilters(prefix string) ShopConnectionFilters {
	return ShopConnectionFilters{prefix}
}

func (ft *ShopConnectionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopConnectionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopConnectionFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopConnectionFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopConnectionFilters) ByConnectionID(ConnectionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == 0,
	}
}

func (ft *ShopConnectionFilters) ByConnectionIDPtr(ConnectionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == nil,
		IsZero: ConnectionID != nil && (*ConnectionID) == 0,
	}
}

func (ft *ShopConnectionFilters) ByToken(Token string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == "",
	}
}

func (ft *ShopConnectionFilters) ByTokenPtr(Token *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == nil,
		IsZero: Token != nil && (*Token) == "",
	}
}

func (ft *ShopConnectionFilters) ByTokenExpiresAt(TokenExpiresAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "token_expires_at",
		Value:  TokenExpiresAt,
		IsNil:  TokenExpiresAt.IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByTokenExpiresAtPtr(TokenExpiresAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "token_expires_at",
		Value:  TokenExpiresAt,
		IsNil:  TokenExpiresAt == nil,
		IsZero: TokenExpiresAt != nil && (*TokenExpiresAt).IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopConnectionFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ShopConnectionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ShopConnectionFilters) ByIsGlobal(IsGlobal bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_global",
		Value:  IsGlobal,
		IsNil:  bool(!IsGlobal),
	}
}

func (ft *ShopConnectionFilters) ByIsGlobalPtr(IsGlobal *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_global",
		Value:  IsGlobal,
		IsNil:  IsGlobal == nil,
		IsZero: IsGlobal != nil && bool(!(*IsGlobal)),
	}
}