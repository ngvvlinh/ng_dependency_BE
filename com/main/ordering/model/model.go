package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	catalogmodel "etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type OrderFeeType string
type FulfillType int

const (
	OrderFeeOther    OrderFeeType = "other"
	OrderFeeShipping OrderFeeType = "shipping"
	OrderFeeTax      OrderFeeType = "tax"

	// FulfillNone: Tự quản lý trên đơn hàng
	FulfillNone     FulfillType = 0  // none
	FulfillManual   FulfillType = 1  // manual
	FulfillShipment FulfillType = 10 // shipment
	FulfillShipnow  FulfillType = 11 // shipnow
)

var _ = sqlgenOrder(&Order{})

type Order struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	EdCode     string
	ProductIDs []dot.ID
	VariantIDs []dot.ID
	PartnerID  dot.ID

	Currency      string
	PaymentMethod string

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

	CustomerConfirm model.Status3
	ShopConfirm     model.Status3
	ConfirmStatus   model.Status3

	FulfillmentShippingStatus model.Status5
	EtopPaymentStatus         model.Status4

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status model.Status5

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int

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
	OrderSourceType model.OrderSourceType
	OrderSourceID   dot.ID
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool

	// @deprecated: use try_on instead
	GhnNoteCode string
	TryOn       model.TryOn

	CustomerNameNorm string
	ProductNameNorm  string
	FulfillmentType  FulfillType
	FulfillmentIDs   []dot.ID
	ExternalMeta     json.RawMessage
	TradingShopID    dot.ID

	// payment
	PaymentStatus model.Status4
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

func (m *Order) GetTryOn() model.TryOn {
	if m.TryOn != "" {
		return m.TryOn
	}
	return model.TryOnFromGHNNoteCode(m.GhnNoteCode)
}

func (m *Order) BeforeInsert() error {
	if (m.TryOn == "" || m.TryOn == "unknown") && m.GhnNoteCode != "" {
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
	if (m.TryOn == "" || m.TryOn == "unknown") && m.GhnNoteCode != "" {
		m.TryOn = model.TryOnFromGHNNoteCode(m.GhnNoteCode)
	}
	if m.ShopShipping != nil && m.ShopShipping.ShippingProvider != "" {
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
	ShopAddress         *OrderAddress          `json:"shop_address"`
	ReturnAddress       *OrderAddress          `json:"return_address"`
	ExternalServiceID   string                 `json:"external_service_id"`
	ExternalShippingFee int                    `json:"external_shipping_fee"`
	ExternalServiceName string                 `json:"external_service_name"`
	ShippingProvider    model.ShippingProvider `json:"shipping_provider"`
	ProviderServiceID   string                 `json:"provider_service_id"`
	IncludeInsurance    bool                   `json:"include_insurance"`

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

func (s *OrderShipping) GetShippingProvider() model.ShippingProvider {
	if s == nil {
		return ""
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
	MetaFields    []*MetaField
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

type OrderFeeLine struct {
	Amount int          `json:"amount"`
	Desc   string       `json:"desc"`
	Code   string       `json:"code"`
	Name   string       `json:"name"`
	Type   OrderFeeType `json:"type"`
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

	Company     string       `json:"company"`
	Address1    string       `json:"address1"`
	Address2    string       `json:"address2"`
	Coordinates *Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
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
