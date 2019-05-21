package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type OrderFeeType string

const (
	OrderFeeOther    OrderFeeType = "other"
	OrderFeeShipping OrderFeeType = "shipping"
	OrderFeeTax      OrderFeeType = "tax"
)

var _ = sqlgenOrder(&Order{})

type Order struct {
	ID          int64
	ShopID      int64
	Code        string
	EdCode      string
	ProductIDs  []int64
	VariantIDs  []int64
	SupplierIDs []int64
	PartnerID   int64

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
	ExternalConfirm model.Status3
	ShopConfirm     model.Status3
	ConfirmStatus   model.Status3

	FulfillmentShippingStatus model.Status5
	CustomerPaymentStatus     model.Status3
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
	OrderSourceID   int64
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool
	// Fulfillments    []*Fulfillment `sq:"-"`
	ExternalData *OrderExternal `sq:"-"`

	// @deprecated: use try_on instead
	GhnNoteCode string
	TryOn       model.TryOn

	CustomerNameNorm string
	ProductNameNorm  string
}

func (m *Order) SelfURL(baseURL string, accType int) string {
	switch accType {
	case model.TagEtop, model.TagSupplier:
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

func (s *OrderShipping) GetPtrShippingServiceCode() *string {
	if s == nil || s.ExternalServiceID == "" {
		return nil
	}
	return &s.ExternalServiceID
}

type OrderDiscount struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

var _ = sqlgenOrderExternal(&OrderExternal{})

type OrderExternal struct {
	ID                   int64
	OrderSourceID        int64
	ExternalOrderSource  string
	ExternalProvider     string
	ExternalOrderID      string
	ExternalOrderCode    string
	ExternalUserID       string
	ExternalCustomerID   string
	ExternalCreatedAt    time.Time
	ExternalProcessedAt  time.Time
	ExternalUpdatedAt    time.Time
	ExternalClosedAt     time.Time
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
	ExternalLines        []*ExternalOrderLine
	ExternalData         json.RawMessage
}

type OrderExternalCreate struct {
	ID int64

	OrderExternalCreateOrder
	OrderExternalCreateExternal
}

var _ = sqlgenOrderExternalCreateOrder(&OrderExternalCreateOrder{}, &Order{})

type OrderExternalCreateOrder struct {
	ShopID          int64
	Code            string
	SupplierIDs     []int64 `sq:"'supplier_ids'"`
	Currency        string
	PaymentMethod   string
	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	CustomerPhone   string
	CustomerEmail   string

	CreatedAt time.Time `sq:"create"`

	ProcessedAt  time.Time
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string

	Lines         OrderLinesList
	VariantIDs    []int64 `sq:"'variant_ids'"`
	Discounts     []*OrderDiscount
	TotalItems    int
	BasketValue   int
	TotalWeight   int
	TotalTax      int
	TotalDiscount int
	TotalAmount   int
	ShopCOD       int

	ExternalConfirm model.Status3

	OrderSourceType model.OrderSourceType
	OrderSourceID   int64
	ExternalOrderID string
	IsOutsideEtop   bool
}

var _ = sqlgenOrderExternalCreateExternal(&OrderExternalCreateExternal{}, &OrderExternal{})

type OrderExternalCreateExternal struct {
	OrderSourceID        int64
	ExternalOrderSource  string
	ExternalOrderID      string
	ExternalOrderCode    string
	ExternalUserID       string
	ExternalCustomerID   string
	ExternalCreatedAt    time.Time `sq:"create"`
	ExternalProcessedAt  time.Time
	ExternalUpdatedAt    time.Time `sq:"update"`
	ExternalClosedAt     time.Time
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
	ExternalLines        []*ExternalOrderLine

	ExternalData json.RawMessage
}

type OrderExternalUpdate struct {
	ID int64

	OrderExternalUpdateOrder
	OrderExternalUpdateExternal
}

var _ = sqlgenOrderExternalUpdateOrder(&OrderExternalUpdateOrder{}, &Order{})

type OrderExternalUpdateOrder struct {
	UpdatedAt    time.Time `sq:"update"`
	ClosedAt     time.Time
	CancelledAt  time.Time
	CancelReason string

	ExternalConfirm model.Status3
}

var _ = sqlgenOrderExternalUpdateExternal(&OrderExternalUpdateExternal{}, &OrderExternal{})

type OrderExternalUpdateExternal struct {
	ExternalUpdatedAt    time.Time `sq:"update"`
	ExternalCancelledAt  time.Time
	ExternalCancelReason string
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

var _ = sqlgenOrderLine(&OrderLine{})

type OrderLine struct {
	OrderID     int64  `json:"order_id"`
	VariantID   int64  `json:"variant_id"`
	ProductName string `json:"product_name"`
	ProductID   int64  `json:"product_id"`

	SupplierID int64 `json:"supplier_id"`
	ShopID     int64 `json:"shop_id"`

	ExternalVariantID       string `json:"x_variant_id"`
	ExternalSupplierOrderID string `json:"-"`

	SupplierName string `json:"supplier_name"`

	// Vendor      string `json:"vendor"` // supplier?

	// CreatedAt    time.Time `json:"-"`
	// CreatedBy    int64     `json:"-"`
	UpdatedAt time.Time `json:"-" `
	// UpdatedBy int64     `json:"-"`
	ClosedAt    time.Time `json:"-" `
	ConfirmedAt time.Time `json:"-" `
	// ConfirmedBy  int64     `json:"-"`
	CancelledAt time.Time `json:"-" `
	// CancelledBy  int64     `json:"-"`
	CancelReason string `json:"-"`

	SupplierConfirm model.Status3 `json:"-"`
	// EtopConfirm     int `json:"etop_confirm"`
	// ConfirmStatus int `json:"confirm_status"`

	Status model.Status3 `json:"-"`

	Weight          int `json:"weight"`
	Quantity        int `json:"quantity"`
	WholesalePrice0 int `json:"wholesale_price_0" sq:"'wholesale_price_0'"`
	WholesalePrice  int `json:"wholesale_price"`
	ListPrice       int `json:"list_price"`
	RetailPrice     int `json:"retail_price"`
	PaymentPrice    int `json:"payment_price"`

	LineAmount       int  `json:"line_amount"`
	TotalDiscount    int  `json:"discount"`
	TotalLineAmount  int  `json:"total_line_amount"`
	RequiresShipping bool `json:"requires_shipping"`

	ImageURL      string                   `json:"image_url" `
	Attributes    []model.ProductAttribute `json:"attributes" sq:"-"`
	IsOutsideEtop bool                     `json:"is_outside_etop"`
	Code          string                   `json:"code"`
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

var _ = sqlgenOrderLineExtended(&OrderLineExtended{},
	&OrderLine{}, sq.AS("ol"), sq.LEFT_JOIN,
	&model.Variant{}, sq.AS("v"), "v.id = ol.variant_id",
)

type OrderLineExtended struct {
	*OrderLine
	*model.Variant
}

func (olExtended *OrderLineExtended) ToOrderLine() *OrderLine {
	return &OrderLine{
		OrderID:                 olExtended.OrderID,
		VariantID:               olExtended.VariantID,
		ProductName:             olExtended.ProductName,
		ProductID:               olExtended.OrderLine.ProductID,
		SupplierID:              olExtended.OrderLine.SupplierID,
		ShopID:                  olExtended.OrderLine.ShopID,
		ExternalVariantID:       olExtended.ExternalVariantID,
		ExternalSupplierOrderID: olExtended.ExternalSupplierOrderID,
		SupplierName:            olExtended.SupplierName,
		UpdatedAt:               olExtended.OrderLine.UpdatedAt,
		ClosedAt:                olExtended.OrderLine.ClosedAt,
		ConfirmedAt:             olExtended.OrderLine.ConfirmedAt,
		CancelledAt:             olExtended.OrderLine.CancelledAt,
		CancelReason:            olExtended.OrderLine.CancelReason,
		SupplierConfirm:         olExtended.OrderLine.SupplierConfirm,
		Status:                  olExtended.OrderLine.Status,
		Weight:                  olExtended.OrderLine.Weight,
		Quantity:                olExtended.OrderLine.Quantity,
		WholesalePrice0:         olExtended.OrderLine.WholesalePrice0,
		WholesalePrice:          olExtended.OrderLine.WholesalePrice,
		ListPrice:               olExtended.OrderLine.ListPrice,
		RetailPrice:             olExtended.OrderLine.RetailPrice,
		PaymentPrice:            olExtended.OrderLine.PaymentPrice,
		LineAmount:              olExtended.OrderLine.LineAmount,
		TotalDiscount:           olExtended.OrderLine.TotalDiscount,
		TotalLineAmount:         olExtended.OrderLine.TotalLineAmount,
		RequiresShipping:        olExtended.OrderLine.RequiresShipping,
		ImageURL:                olExtended.OrderLine.ImageURL,
		Attributes:              olExtended.Variant.Attributes,
		IsOutsideEtop:           olExtended.OrderLine.IsOutsideEtop,
		Code:                    olExtended.OrderLine.Code,
	}
}

type ExternalOrderLine struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	SKU       string `json:"sku"`

	Name        string `json:"name"`
	Title       string `json:"title"`
	VariantName string `json:"variant_name"`
	Vendor      string `json:"vendor"`
	Type        string `json:"type"`

	Weight         int `json:"weight"`
	Quantity       int `json:"quantity"`
	Price          int `json:"price"`
	PriceOriginal  int `json:"price_original"`
	PricePromotion int `json:"price_promotion"`
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

func (o *OrderCustomer) UpdateCustomer(fullname string) *OrderCustomer {
	if fullname != "" {
		o.FullName = fullname
	}
	return o
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

	Company  string `json:"company"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
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

func (m *OrderAddress) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

func GetFeeLinesWithFallback(lines []OrderFeeLine, totalFee *int32, shopShippingFee *int32) []OrderFeeLine {
	if len(lines) == 0 &&
		(totalFee == nil || *totalFee == 0) &&
		shopShippingFee != nil && *shopShippingFee != 0 {
		return []OrderFeeLine{
			{
				Amount: int(*shopShippingFee),
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
