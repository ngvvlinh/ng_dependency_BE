package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/top/types/etc/fee"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/order_source"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	addressmodel "etop.vn/backend/com/main/address/model"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

const (
	OrderFeeOther    = fee.Other
	OrderFeeShipping = fee.Shipping
	OrderFeeTax      = fee.Tax
)

var _ = sqlgenOrder(&Order{})

// +convert:type=ordering.Order
type Order struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	EdCode     string
	ProductIDs []dot.ID
	VariantIDs []dot.ID
	PartnerID  dot.ID

	Currency      string
	PaymentMethod payment_method.PaymentMethod

	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	CustomerName    string
	CustomerPhone   string
	CustomerEmail   string

	CreatedAt    time.Time
	ProcessedAt  time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string

	CustomerConfirm status3.Status
	ShopConfirm     status3.Status
	ConfirmStatus   status3.Status

	FulfillmentShippingStatus status5.Status
	EtopPaymentStatus         status4.Status

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status status5.Status

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int
	FulfillmentStatuses        []int

	Lines           OrderLinesList
	Discounts       []*OrderDiscount
	TotalItems      int
	BasketValue     int
	TotalWeight     int
	TotalTax        int
	OrderDiscount   int
	TotalDiscount   int
	ShopShippingFee int
	TotalFee        int
	FeeLines        OrderFeeLines
	ShopCOD         int
	TotalAmount     int

	OrderNote       string
	ShopNote        string
	ShippingNote    string
	OrderSourceType order_source.Source
	OrderSourceID   dot.ID
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool

	// @deprecated: use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode
	TryOn       try_on.TryOnCode

	CustomerNameNorm string
	ProductNameNorm  string
	FulfillmentType  ordertypes.ShippingType
	FulfillmentIDs   []dot.ID
	ExternalMeta     json.RawMessage
	TradingShopID    dot.ID

	// payment
	PaymentStatus status4.Status
	PaymentID     dot.ID

	ReferralMeta json.RawMessage

	CustomerID dot.ID
	CreatedBy  dot.ID
}

func (m *Order) SelfURL(baseURL string, accType int) string {
	switch accType {
	case model.TagEtop:
		return ""

	case model.TagShop:
		if baseURL == "" || m.ShopID == 0 || m.ID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/orders/%v", baseURL, m.ShopID, m.ID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

func (m *Order) GetTotalFee() int {
	if m.TotalFee == 0 && m.ShopShippingFee != 0 {
		return m.ShopShippingFee
	}
	return m.TotalFee
}

func (m *Order) GetFeeLines() []OrderFeeLine {
	if m.TotalFee == 0 && m.ShopShippingFee != 0 {
		return []OrderFeeLine{
			{
				Amount: m.ShopShippingFee,
				Desc:   "Phí giao hàng (shop tính cho khách)",
				Name:   "Phí giao hàng",
				Type:   OrderFeeShipping,
			},
		}
	}
	return m.FeeLines
}

func (m *Order) GetTryOn() try_on.TryOnCode {
	if m.TryOn != 0 {
		return m.TryOn
	}
	return model.TryOnFromGHNNoteCode(m.GhnNoteCode)
}

func (m *Order) BeforeInsert() error {
	if m.TryOn == 0 && m.GhnNoteCode != 0 {
		m.TryOn = model.TryOnFromGHNNoteCode(m.GhnNoteCode)
	}
	if m.ShopShipping != nil {
		if err := m.ShopShipping.Validate(); err != nil {
			return err
		}
	}

	m.CustomerName = m.Customer.GetFullName()
	m.CustomerNameNorm = validate.NormalizeSearch(m.CustomerName)

	var b strings.Builder
	for _, line := range m.Lines {
		b.WriteString(line.ProductName)
		b.WriteString(" ")
	}

	m.ProductNameNorm = validate.NormalizeSearch(b.String())
	return nil
}

func (m *Order) BeforeUpdate() error {
	if m.TryOn == 0 && m.GhnNoteCode != 0 {
		m.TryOn = model.TryOnFromGHNNoteCode(m.GhnNoteCode)
	}
	if m.ShopShipping != nil && m.ShopShipping.ShippingProvider != 0 {
		if err := m.ShopShipping.Validate(); err != nil {
			return err
		}
	}

	if m.Customer != nil {
		m.CustomerName = m.Customer.GetFullName()
		m.CustomerNameNorm = validate.NormalizeSearch(m.CustomerName)
		m.CustomerEmail = m.Customer.Email
	}
	if m.ShippingAddress != nil {
		m.CustomerPhone = m.ShippingAddress.Phone
	}
	return nil
}

func (m *Order) GetChargeableWeight() int {
	if m.ShopShipping != nil && m.ShopShipping.ChargeableWeight != 0 {
		return m.ShopShipping.ChargeableWeight
	}
	return m.TotalWeight
}

type OrderShipping struct {
	ShopAddress         *OrderAddress                      `json:"shop_address"`
	ReturnAddress       *OrderAddress                      `json:"return_address"`
	ExternalServiceID   string                             `json:"external_service_id"`
	ExternalShippingFee int                                `json:"external_shipping_fee"`
	ExternalServiceName string                             `json:"external_service_name"`
	ShippingProvider    shipping_provider.ShippingProvider `json:"shipping_provider"`
	ProviderServiceID   string                             `json:"provider_service_id"`
	IncludeInsurance    bool                               `json:"include_insurance"`

	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`

	GrossWeight      int `json:"gross_weight"`
	ChargeableWeight int `json:"chargeable_weight"`
}

func (s *OrderShipping) Validate() error {
	if !model.VerifyShippingProvider(s.ShippingProvider) {
		return cm.Errorf(cm.InvalidArgument, nil, "Nhà vận chuyển không hợp lệ")
	}
	return nil
}

func (s *OrderShipping) GetPickupAddress() *OrderAddress {
	if s == nil {
		return nil
	}
	return s.ShopAddress
}

func (s *OrderShipping) GetShippingProvider() shipping_provider.ShippingProvider {
	if s == nil {
		return 0
	}
	return s.ShippingProvider
}

func (s *OrderShipping) GetShippingServiceCode() string {
	if s == nil {
		return ""
	}
	return cm.Coalesce(s.ProviderServiceID, s.ExternalServiceID)
}

func (s *OrderShipping) GetPtrShippingServiceCode() dot.NullString {
	if s == nil || s.ExternalServiceID == "" {
		return dot.NullString{}
	}
	return dot.String(s.ExternalServiceID)
}

// +convert:type=ordering.OrderDiscount
type OrderDiscount struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type OrderLinesList []*OrderLine

func (lines OrderLinesList) GetTotalItems() int {
	s := 0
	for _, line := range lines {
		s += line.Quantity
	}
	return s
}

func (lines OrderLinesList) GetTotalWeight() int {
	s := 0
	for _, line := range lines {
		s += line.Weight
	}
	return s
}

func (lines OrderLinesList) GetTotalRetailAmount() int {
	s := 0
	for _, line := range lines {
		s += line.RetailPrice
	}
	return s
}

func (lines OrderLinesList) GetTotalPaymentAmount() int {
	s := 0
	for _, line := range lines {
		s += line.PaymentPrice
	}
	return s
}

func (lines OrderLinesList) GetSummary() string {
	var b strings.Builder
	for _, line := range lines {
		fprintf(&b, "%2d x %v", line.Quantity, line.ProductName)
		if len(line.Attributes) > 0 {
			fprintf(&b, " (")
		}
		for i, attr := range line.Attributes {
			if i > 0 {
				fprintf(&b, ", ")
			}
			fprintf(&b, "%v: %v", attr.Name, attr.Value)
		}
		if len(line.Attributes) > 0 {
			fprintf(&b, ")")
		}
		fprintf(&b, "\n")
	}
	return b.String()
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

var _ = sqlgenOrderLine(&OrderLine{})

// +convert:type=ordering/types.ItemLine
type OrderLine struct {
	OrderID     dot.ID `json:"order_id"`
	VariantID   dot.ID `json:"variant_id"`
	ProductName string `json:"product_name"`
	ProductID   dot.ID `json:"product_id"`
	ShopID      dot.ID `json:"shop_id"`

	Weight       int `json:"weight"`
	Quantity     int `json:"quantity"`
	ListPrice    int `json:"list_price"`
	RetailPrice  int `json:"retail_price"`
	PaymentPrice int `json:"payment_price"`

	LineAmount      int `json:"line_amount"`
	TotalDiscount   int `json:"discount"`
	TotalLineAmount int `json:"total_line_amount"`

	ImageURL      string                           `json:"image_url" `
	Attributes    []*catalogmodel.ProductAttribute `json:"attributes" sq:"-"`
	IsOutsideEtop bool                             `json:"is_outside_etop"`
	Code          string                           `json:"code"`
	MetaFields    []*MetaField                     `json:"meta_fields"`
}

func (l *OrderLine) GetRetailAmount() int {
	if l == nil {
		return 0
	}
	return l.RetailPrice * l.Quantity
}

func (l *OrderLine) GetPaymentAmount() int {
	if l == nil {
		return 0
	}
	return l.PaymentPrice * l.Quantity
}

func (l *OrderLine) GetTotalDiscount() int {
	if l == nil {
		return 0
	}
	return (l.RetailPrice - l.PaymentPrice) * l.Quantity
}

// +convert:type=ordering.OrderFeeLine
type OrderFeeLine struct {
	Amount int         `json:"amount"`
	Desc   string      `json:"desc"`
	Code   string      `json:"code"`
	Name   string      `json:"name"`
	Type   fee.FeeType `json:"type"`
}

type OrderFeeLines []OrderFeeLine

func (feeLines OrderFeeLines) GetTotalFee() int {
	s := 0
	for _, line := range feeLines {
		s += line.Amount
	}
	return s
}

func (feeLines OrderFeeLines) GetShippingFee() int {
	s := 0
	for _, line := range feeLines {
		if line.Type == OrderFeeShipping {
			s += line.Amount
		}
	}
	return s
}

func (feeLines OrderFeeLines) GetTax() int {
	s := 0
	for _, line := range feeLines {
		if line.Type == OrderFeeTax {
			s += line.Amount
		}
	}
	return s
}

func (feeLines OrderFeeLines) GetOtherFee() int {
	s := 0
	for _, line := range feeLines {
		if line.Type != OrderFeeShipping && line.Type != OrderFeeTax {
			s += line.Amount
		}
	}
	return s
}

// +convert:type=ordering.OrderCustomer
type OrderCustomer struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	VerifiedEmail bool   `json:"verified_email"`

	ExternalID string `json:"external_id"`
}

func (m *OrderCustomer) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

func (m *OrderCustomer) UpdateCustomer(fullname string) *OrderCustomer {
	if fullname != "" {
		m.FullName = fullname
	}
	return m
}

// +convert:type=ordering/types.Address
type OrderAddress struct {
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company     string                    `json:"company"`
	Address1    string                    `json:"address1"`
	Address2    string                    `json:"address2"`
	Coordinates *addressmodel.Coordinates `json:"coordinates"`
}

func (o *OrderAddress) UpdateAddress(phone string, fullname string) *OrderAddress {
	if phone != "" {
		o.Phone = phone
	}
	if fullname != "" {
		o.FullName = fullname
	}
	return o
}

func (o *OrderAddress) GetFullName() string {
	if o.FullName != "" {
		return o.FullName
	}
	return o.FirstName + " " + o.LastName
}

func GetFeeLinesWithFallback(lines []OrderFeeLine, totalFee dot.NullInt, shopShippingFee dot.NullInt) []OrderFeeLine {
	if len(lines) == 0 &&
		totalFee.Apply(0) == 0 &&
		shopShippingFee.Apply(0) == 0 {
		return []OrderFeeLine{
			{
				Amount: shopShippingFee.Apply(0),
				Desc:   "Phí giao hàng tính cho khách",
				Name:   "Phí giao hàng tính cho khách",
				Type:   OrderFeeShipping,
			},
		}
	}
	return lines
}

func GetShippingFeeFromFeeLines(lines []OrderFeeLine) int {
	s := 0
	for _, line := range lines {
		if line.Type == OrderFeeShipping {
			s += line.Amount
		}
	}
	return s
}

func CalcTotalFee(lines []OrderFeeLine) int {
	s := 0
	for _, line := range lines {
		s += line.Amount
	}
	return s
}
