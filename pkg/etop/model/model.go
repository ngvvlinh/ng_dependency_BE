package model

import (
	"encoding/json"
	"strings"
	"time"

	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/order_source"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
)

type (
	FulfillmentEndpoint string
	ShippingPriceSource string

	RoleType  string
	FfmAction int
	CodeType  string

	ShippingRouteType string

	DBName string
)

// Type constants
const (
	DBMain     DBName = "main"
	DBNotifier DBName = "notifier"

	TypeShippingSourceEtop    ShippingPriceSource = "etop"
	TypeShippingSourceCarrier ShippingPriceSource = "carrier"

	TypeGHNCode  = "GHN"
	TypeGHTKCode = "GHTK"

	ShippingServiceNameStandard = "Chuẩn"
	ShippingServiceNameFaster   = "Nhanh"

	FFShop     FulfillmentEndpoint = "shop"
	FFCustomer FulfillmentEndpoint = "customer"

	// Don't change these values
	TagUser      = 17
	TagPartner   = 21
	TagShop      = 33
	TagAffiliate = 35
	TagEtop      = 101
	TagImport    = 111

	EtopAccountID        = TagEtop
	EtopTradingAccountID = 1000015764575267699
	TopShipID            = 1000030662086749358

	CurrencyVND = "vnd"

	AddressTypeGeneral   = "general"
	AddressTypeWarehouse = "warehouse"
	AddressTypeShipTo    = "shipto"
	AddressTypeShipFrom  = "shipfrom"

	CodeTypeShop  = "shop"
	CodeTypeOrder = "order"

	CodeTypeMoneyTransaction         = "money_transaction_shipping"
	CodeTypeMoneyTransactionEtop     = "money_transaction_shipping_etop"
	CodeTypeMoneyTransactionExternal = "money_transaction_external"
	CodeTypeConnection               = "connection"

	// For import

	FfmActionNothing FfmAction = iota
	FfmActionCancel
	FfmActionReturn

	RouteSameProvince ShippingRouteType = "noi_tinh"
	RouteSameRegion   ShippingRouteType = "noi_mien"
	RouteNationWide   ShippingRouteType = "toan_quoc"
)

type NormalizedPhone = validate.NormalizedPhone
type NormalizedEmail = validate.NormalizedEmail

//        SyncAt: updated when sending data to external service successfully
//     TrySyncAt: updated when sending data to external service (may unsuccessfully)
//    LastSyncAt: updated when querying external data source successfully
// LastTrySypkg/etop/model/model.go:18:ncAt: updated when querying external data source (may unsuccessfully)

func EtopPaymentStatusLabel(s status4.Status) string {
	switch s {
	case status4.N:
		return "Không cần thanh toán"
	case status4.Z:
		return "Chưa đối soát"
	case status4.P:
		return "Đã thanh toán"
	case status4.S:
		return "Cần thanh toán"
	default:
		return "Không xác định"
	}
}

func OrderStatusLabel(s status5.Status) string {
	switch s {
	case status5.NS:
		return "Trả hàng"
	case status5.N:
		return "Huỷ"
	case status5.Z:
		return "Mới"
	case status5.P:
		return "Thành công"
	case status5.S:
		return "Đang xử lý"
	default:
		return "Không xác định"
	}
}

func PaymentStatusLabel(s status4.Status) string {
	switch s {
	case 0:
		return ""
	case 1:
		return "Đã đối soát"
	case 2:
		return "Chưa đối soát"
	default:
		return "Không xác định"
	}
}

func AccountTypeLabel(t account_type.AccountType) string {
	switch t {
	case account_type.Shop:
		return "cửa hàng"
	case account_type.Partner:
		return "tài khoản"
	case account_type.Etop:
		return "tài khoản"
	default:
		return "tài khoản"
	}
}

func (t CodeType) Validate() error {
	switch t {
	case "":
		return cm.Error(cm.InvalidArgument, "Missing Type", nil)
	case CodeTypeShop,
		CodeTypeOrder,
		CodeTypeMoneyTransaction,
		CodeTypeMoneyTransactionExternal,
		CodeTypeMoneyTransactionEtop,
		CodeTypeConnection:
		// no-op
	default:
		return cm.Error(cm.InvalidArgument, "Type is not valid", nil)
	}
	return nil
}

func VerifyPaymentMethod(s payment_method.PaymentMethod) bool {
	switch s {
	case payment_method.COD, payment_method.Bank,
		payment_method.Other, payment_method.VTPay:
		return true
	}
	return false
}

func VerifyOrderSource(s order_source.Source) bool {
	switch s {
	case order_source.EtopPOS, order_source.EtopPXS, order_source.EtopCMX, order_source.TSApp, order_source.API, order_source.Import, order_source.EtopApp:
		return true
	}
	return false
}

func VerifyShippingProvider(s shipping_provider.ShippingProvider) bool {
	if s == 0 {
		return true
	}
	switch s {
	case shipping_provider.GHN, shipping_provider.GHTK, shipping_provider.VTPost, shipping_provider.Manual:
		return true
	}
	return false
}

func TryOnFromGHNNoteCode(c ghn_note_code.GHNNoteCode) try_on.TryOnCode {
	switch c {
	case ghn_note_code.KHONGCHOXEMHANG:
		return try_on.None
	case ghn_note_code.CHOXEMHANGKHONGTHU:
		return try_on.Open
	case ghn_note_code.CHOTHUHANG:
		return try_on.Try
	default:
		return 0
	}
}

func GHNNoteCodeFromTryOn(to try_on.TryOnCode) ghn_note_code.GHNNoteCode {
	switch to {
	case try_on.None:
		return ghn_note_code.KHONGCHOXEMHANG
	case try_on.Open:
		return ghn_note_code.CHOXEMHANGKHONGTHU
	case try_on.Try:
		return ghn_note_code.CHOTHUHANG
	default:
		return 0
	}
}

type Unit struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Unit     string  `json:"unit"`
	UnitConv float64 `json:"unit_conv"`
	Price    int     `json:"price"`
}

type StatusQuery struct {
	Status *status3.Status
}

type ExternalShippingLog struct {
	StateText string `json:"StateText"`
	Time      string `json:"Time"`
	Message   string `json:"Message"`
}

func (f FfmAction) ToShippingState() shipping.State {
	switch f {
	case FfmActionCancel:
		return shipping.Cancelled
	case FfmActionReturn:
		return shipping.Returning
	default:
		return shipping.Unknown
	}
}

// +sqlgen
type Code struct {
	Code      string
	Type      CodeType
	CreatedAt time.Time `sq:"create"`
}

func (code *Code) Validate() error {
	if code.Code == "" {
		return cm.Error(cm.InvalidArgument, "Missing Code", nil)
	}
	if err := code.Type.Validate(); err != nil {
		return err
	}
	return nil
}

type CreateCodeCommand struct {
	Code *Code
}

// +sqlgen
type ShippingSource struct {
	ID        dot.ID
	Name      string
	Username  string
	Type      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

// +sqlgen
type ShippingSourceInternal struct {
	ID          dot.ID
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	LastSyncAt  time.Time
	AccessToken string
	ExpiresAt   time.Time `json:"expires_at"`
	Secret      *ShippingSourceSecret
}

type ShippingSourceSecret struct {
	CustomerID int    `json:"CustomerID"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`

	// Use for VTPost
	GroupAddressID int `json:"GroupAddressID"`
}

type AvailableShippingService struct {
	Name string

	// ServiceFee: Tổng phí giao hàng (đã bao gồm phí chính + các phụ phí khác)
	ServiceFee int

	// ShippingFeeMain: Phí chính giao hàng
	ShippingFeeMain   int
	Provider          shipping_provider.ShippingProvider
	ProviderServiceID string

	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
	Source             ShippingPriceSource
	ConnectionInfo     *ConnectionInfo
}

type ConnectionInfo struct {
	ID       dot.ID
	Name     string
	ImageURL string
}

func (service *AvailableShippingService) ApplyFeeMain(feeMain int) {
	service.ServiceFee = service.ServiceFee - service.ShippingFeeMain + feeMain
	service.ShippingFeeMain = feeMain
}

// +sqlgen
type Webhook struct {
	ID        dot.ID
	AccountID dot.ID
	Entities  []string
	Fields    []string
	URL       string
	Metadata  string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

func (m *Webhook) BeforeInsert() error {
	if m == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "empty data")
	}
	if m.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing Name")
	}
	if !validate.URL(m.URL) {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url không hợp lệ")
	}
	if cmenv.IsProd() && !strings.HasPrefix(m.URL, "https://") {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url phải là https://")
	}
	if len(m.Entities) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing entity")
	}
	if len(m.Fields) > 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Thông tin fields chưa được hỗ trợ, vui lòng để trống")
	}

	mp := make(map[string]bool)
	for _, item := range m.Entities {
		if !validate.LowercaseID(item) {
			return cm.Errorf(cm.InvalidArgument, nil, `invalid entity: "%v"`, item)
		}
		switch item {
		case "order", "fulfillment", "product", "variant", "customer", "inventory_level", "customer_address", "customer_group", "customer_group_relationship", "product_collection", "product_collection_relationship":
			if mp[item] {
				return cm.Errorf(cm.InvalidArgument, nil, `duplicated entity: "%v"`, item)
			}
			mp[item] = true
		default:
			return cm.Errorf(cm.InvalidArgument, nil, `unknown entity: "%v"`, item)
		}
	}

	m.ID = cm.NewID()
	return nil
}

// +sqlgen
type Callback struct {
	ID        dot.ID
	WebhookID dot.ID
	AccountID dot.ID
	CreatedAt time.Time `sq:"create"`
	Changes   json.RawMessage
	Result    json.RawMessage // WebhookStatesError
}
