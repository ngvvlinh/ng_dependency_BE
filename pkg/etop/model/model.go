package model

import (
	"time"

	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

type (
	FulfillmentEndpoint string
	ShippingPriceSource string

	RoleType  string
	FfmAction int
	CodeType  string

	ShippingRouteType string
)

// Type constants
const (
	TypeShippingSourceEtop    ShippingPriceSource = "etop"
	TypeShippingSourceCarrier ShippingPriceSource = "carrier"

	TypeGHNCode  = "GHN"
	TypeGHTKCode = "GHTK"

	ShippingServiceNameStandard = "Chuẩn"
	ShippingServiceNameFaster   = "Nhanh"

	FFShop     FulfillmentEndpoint = "shop"
	FFCustomer FulfillmentEndpoint = "customer"

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
	case order_source.EtopPOS, order_source.EtopPXS, order_source.EtopCMX, order_source.TSApp, order_source.API, order_source.Import, order_source.EtopApp, order_source.Ecomify:
		return true
	}
	return false
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
