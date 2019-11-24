// source: etop/order/order.proto

package order

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	spreadsheet "etop.vn/api/pb/common/spreadsheet"
	etop "etop.vn/api/pb/etop"
	fee "etop.vn/api/pb/etop/etc/fee"
	gender "etop.vn/api/pb/etop/etc/gender"
	ghn_note_code "etop.vn/api/pb/etop/etc/ghn_note_code"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	shipping_fee_type "etop.vn/api/pb/etop/etc/shipping_fee_type"
	shipping_provider "etop.vn/api/pb/etop/etc/shipping_provider"
	status3 "etop.vn/api/pb/etop/etc/status3"
	status4 "etop.vn/api/pb/etop/etc/status4"
	status5 "etop.vn/api/pb/etop/etc/status5"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	source "etop.vn/api/pb/etop/order/source"
	"etop.vn/capi/dot"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type OrdersResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Orders               []*Order         `protobuf:"bytes,2,rep,name=orders" json:"orders"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *OrdersResponse) Reset()         { *m = OrdersResponse{} }
func (m *OrdersResponse) String() string { return proto.CompactTextString(m) }
func (*OrdersResponse) ProtoMessage()    {}
func (*OrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{0}
}

var xxx_messageInfo_OrdersResponse proto.InternalMessageInfo

func (m *OrdersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *OrdersResponse) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type Order struct {
	ExportedFields []string `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	Id             dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId         dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	ShopName       string   `protobuf:"bytes,3,opt,name=shop_name,json=shopName" json:"shop_name"`
	Code           string   `protobuf:"bytes,4,opt,name=code" json:"code"`
	// the same as external_code
	EdCode          string         `protobuf:"bytes,5,opt,name=ed_code,json=edCode" json:"ed_code"`
	ExternalCode    string         `protobuf:"bytes,8,opt,name=external_code,json=externalCode" json:"external_code"`
	Source          source.Source  `protobuf:"varint,6,opt,name=source,enum=source.Source" json:"source"`
	PartnerId       dot.ID         `protobuf:"varint,55,opt,name=partner_id,json=partnerId" json:"partner_id"`
	ExternalId      string         `protobuf:"bytes,56,opt,name=external_id,json=externalId" json:"external_id"`
	ExternalUrl     string         `protobuf:"bytes,57,opt,name=external_url,json=externalUrl" json:"external_url"`
	SelfUrl         string         `protobuf:"bytes,58,opt,name=self_url,json=selfUrl" json:"self_url"`
	PaymentMethod   string         `protobuf:"bytes,7,opt,name=payment_method,json=paymentMethod" json:"payment_method"`
	Customer        *OrderCustomer `protobuf:"bytes,9,opt,name=customer" json:"customer"`
	CustomerAddress *OrderAddress  `protobuf:"bytes,10,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	BillingAddress  *OrderAddress  `protobuf:"bytes,11,opt,name=billing_address,json=billingAddress" json:"billing_address"`
	ShippingAddress *OrderAddress  `protobuf:"bytes,12,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	CreatedAt       dot.Time       `protobuf:"bytes,13,opt,name=created_at,json=createdAt" json:"created_at"`
	ProcessedAt     dot.Time       `protobuf:"bytes,14,opt,name=processed_at,json=processedAt" json:"processed_at"`
	UpdatedAt       dot.Time       `protobuf:"bytes,15,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt        dot.Time       `protobuf:"bytes,16,opt,name=closed_at,json=closedAt" json:"closed_at"`
	ConfirmedAt     dot.Time       `protobuf:"bytes,17,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelledAt     dot.Time       `protobuf:"bytes,18,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CancelReason    string         `protobuf:"bytes,19,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	ShConfirm       status3.Status `protobuf:"varint,31,opt,name=sh_confirm,json=shConfirm,enum=status3.Status" json:"sh_confirm"`
	// @deprecated replaced by confirm_status
	Confirm       status3.Status `protobuf:"varint,33,opt,name=confirm,enum=status3.Status" json:"confirm"`
	ConfirmStatus status3.Status `protobuf:"varint,70,opt,name=confirm_status,json=confirmStatus,enum=status3.Status" json:"confirm_status"`
	Status        status5.Status `protobuf:"varint,34,opt,name=status,enum=status5.Status" json:"status"`
	// @deprecated replaced by fulfillment_shipping_status
	FulfillmentStatus         status5.Status   `protobuf:"varint,35,opt,name=fulfillment_status,json=fulfillmentStatus,enum=status5.Status" json:"fulfillment_status"`
	FulfillmentShippingStatus status5.Status   `protobuf:"varint,36,opt,name=fulfillment_shipping_status,json=fulfillmentShippingStatus,enum=status5.Status" json:"fulfillment_shipping_status"`
	CustomerPaymentStatus     status3.Status   `protobuf:"varint,37,opt,name=customer_payment_status,json=customerPaymentStatus,enum=status3.Status" json:"customer_payment_status"`
	EtopPaymentStatus         status4.Status   `protobuf:"varint,69,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	Lines                     []*OrderLine     `protobuf:"bytes,41,rep,name=lines" json:"lines"`
	Discounts                 []*OrderDiscount `protobuf:"bytes,42,rep,name=discounts" json:"discounts"`
	TotalItems                int32            `protobuf:"varint,43,opt,name=total_items,json=totalItems" json:"total_items"`
	BasketValue               int32            `protobuf:"varint,44,opt,name=basket_value,json=basketValue" json:"basket_value"`
	TotalWeight               int32            `protobuf:"varint,45,opt,name=total_weight,json=totalWeight" json:"total_weight"`
	OrderDiscount             int32            `protobuf:"varint,46,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	TotalDiscount             int32            `protobuf:"varint,47,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount               int32            `protobuf:"varint,48,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	OrderNote                 string           `protobuf:"bytes,50,opt,name=order_note,json=orderNote" json:"order_note"`
	// @deprecated use fee_lines.shipping
	ShippingFee int32           `protobuf:"varint,52,opt,name=shipping_fee,json=shippingFee" json:"shipping_fee"`
	TotalFee    int32           `protobuf:"varint,53,opt,name=total_fee,json=totalFee" json:"total_fee"`
	FeeLines    []*OrderFeeLine `protobuf:"bytes,54,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	// @deperecated use fee_lines.shipping instead
	ShopShippingFee int32 `protobuf:"varint,61,opt,name=shop_shipping_fee,json=shopShippingFee" json:"shop_shipping_fee"`
	// @deprecated use shop_shipping.shipping_note instead
	ShippingNote string `protobuf:"bytes,51,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	// @deperecated use shop_shipping.cod_amount instead
	ShopCod      int32           `protobuf:"varint,62,opt,name=shop_cod,json=shopCod" json:"shop_cod"`
	ReferenceUrl string          `protobuf:"bytes,63,opt,name=reference_url,json=referenceUrl" json:"reference_url"`
	Fulfillments []*XFulfillment `protobuf:"bytes,74,rep,name=fulfillments" json:"fulfillments"`
	// @deprecated use shipping instead
	ShopShipping *OrderShipping `protobuf:"bytes,66,opt,name=shop_shipping,json=shopShipping" json:"shop_shipping"`
	Shipping     *OrderShipping `protobuf:"bytes,71,opt,name=shipping" json:"shipping"`
	// @deprecated use try_on_code instead
	GhnNoteCode     string   `protobuf:"bytes,68,opt,name=ghn_note_code,json=ghnNoteCode" json:"ghn_note_code"`
	FulfillmentType string   `protobuf:"bytes,72,opt,name=fulfillment_type,json=fulfillmentType" json:"fulfillment_type"`
	FulfillmentIds  []dot.ID `protobuf:"varint,73,rep,name=fulfillment_ids,json=fulfillmentIds" json:"fulfillment_ids"`
	// received_amount get from receipt
	ReceivedAmount       int            `protobuf:"varint,75,opt,name=received_amount,json=receivedAmount" json:"received_amount"`
	CustomerId           dot.ID         `protobuf:"varint,76,opt,name=customer_id,json=customerId" json:"customer_id"`
	PaymentStatus        status4.Status `protobuf:"varint,77,opt,name=payment_status,json=paymentStatus,enum=status4.Status" json:"payment_status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{1}
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *Order) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Order) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Order) GetShopName() string {
	if m != nil {
		return m.ShopName
	}
	return ""
}

func (m *Order) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Order) GetEdCode() string {
	if m != nil {
		return m.EdCode
	}
	return ""
}

func (m *Order) GetExternalCode() string {
	if m != nil {
		return m.ExternalCode
	}
	return ""
}

func (m *Order) GetSource() source.Source {
	if m != nil {
		return m.Source
	}
	return source.Source_unknown
}

func (m *Order) GetPartnerId() dot.ID {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *Order) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *Order) GetExternalUrl() string {
	if m != nil {
		return m.ExternalUrl
	}
	return ""
}

func (m *Order) GetSelfUrl() string {
	if m != nil {
		return m.SelfUrl
	}
	return ""
}

func (m *Order) GetPaymentMethod() string {
	if m != nil {
		return m.PaymentMethod
	}
	return ""
}

func (m *Order) GetCustomer() *OrderCustomer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *Order) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *Order) GetBillingAddress() *OrderAddress {
	if m != nil {
		return m.BillingAddress
	}
	return nil
}

func (m *Order) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *Order) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *Order) GetShConfirm() status3.Status {
	if m != nil {
		return m.ShConfirm
	}
	return status3.Status_Z
}

func (m *Order) GetConfirm() status3.Status {
	if m != nil {
		return m.Confirm
	}
	return status3.Status_Z
}

func (m *Order) GetConfirmStatus() status3.Status {
	if m != nil {
		return m.ConfirmStatus
	}
	return status3.Status_Z
}

func (m *Order) GetStatus() status5.Status {
	if m != nil {
		return m.Status
	}
	return status5.Status_Z
}

func (m *Order) GetFulfillmentStatus() status5.Status {
	if m != nil {
		return m.FulfillmentStatus
	}
	return status5.Status_Z
}

func (m *Order) GetFulfillmentShippingStatus() status5.Status {
	if m != nil {
		return m.FulfillmentShippingStatus
	}
	return status5.Status_Z
}

func (m *Order) GetCustomerPaymentStatus() status3.Status {
	if m != nil {
		return m.CustomerPaymentStatus
	}
	return status3.Status_Z
}

func (m *Order) GetEtopPaymentStatus() status4.Status {
	if m != nil {
		return m.EtopPaymentStatus
	}
	return status4.Status_Z
}

func (m *Order) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *Order) GetDiscounts() []*OrderDiscount {
	if m != nil {
		return m.Discounts
	}
	return nil
}

func (m *Order) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *Order) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *Order) GetTotalWeight() int32 {
	if m != nil {
		return m.TotalWeight
	}
	return 0
}

func (m *Order) GetOrderDiscount() int32 {
	if m != nil {
		return m.OrderDiscount
	}
	return 0
}

func (m *Order) GetTotalDiscount() int32 {
	if m != nil {
		return m.TotalDiscount
	}
	return 0
}

func (m *Order) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *Order) GetOrderNote() string {
	if m != nil {
		return m.OrderNote
	}
	return ""
}

func (m *Order) GetShippingFee() int32 {
	if m != nil {
		return m.ShippingFee
	}
	return 0
}

func (m *Order) GetTotalFee() int32 {
	if m != nil {
		return m.TotalFee
	}
	return 0
}

func (m *Order) GetFeeLines() []*OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *Order) GetShopShippingFee() int32 {
	if m != nil {
		return m.ShopShippingFee
	}
	return 0
}

func (m *Order) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *Order) GetShopCod() int32 {
	if m != nil {
		return m.ShopCod
	}
	return 0
}

func (m *Order) GetReferenceUrl() string {
	if m != nil {
		return m.ReferenceUrl
	}
	return ""
}

func (m *Order) GetFulfillments() []*XFulfillment {
	if m != nil {
		return m.Fulfillments
	}
	return nil
}

func (m *Order) GetShopShipping() *OrderShipping {
	if m != nil {
		return m.ShopShipping
	}
	return nil
}

func (m *Order) GetShipping() *OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

func (m *Order) GetGhnNoteCode() string {
	if m != nil {
		return m.GhnNoteCode
	}
	return ""
}

func (m *Order) GetFulfillmentType() string {
	if m != nil {
		return m.FulfillmentType
	}
	return ""
}

func (m *Order) GetFulfillmentIds() []dot.ID {
	if m != nil {
		return m.FulfillmentIds
	}
	return nil
}

func (m *Order) GetCustomerId() dot.ID {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

func (m *Order) GetPaymentStatus() status4.Status {
	if m != nil {
		return m.PaymentStatus
	}
	return status4.Status_Z
}

type OrderLineMetaField struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key"`
	Value                string   `protobuf:"bytes,2,opt,name=value" json:"value"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderLineMetaField) Reset()         { *m = OrderLineMetaField{} }
func (m *OrderLineMetaField) String() string { return proto.CompactTextString(m) }
func (*OrderLineMetaField) ProtoMessage()    {}
func (*OrderLineMetaField) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{2}
}

var xxx_messageInfo_OrderLineMetaField proto.InternalMessageInfo

func (m *OrderLineMetaField) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *OrderLineMetaField) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *OrderLineMetaField) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type OrderLine struct {
	ExportedFields       []string              `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	OrderId              dot.ID                `protobuf:"varint,2,opt,name=order_id,json=orderId" json:"order_id"`
	VariantId            dot.ID                `protobuf:"varint,3,opt,name=variant_id,json=variantId" json:"variant_id"`
	ProductName          string                `protobuf:"bytes,7,opt,name=product_name,json=productName" json:"product_name"`
	IsOutsideEtop        bool                  `protobuf:"varint,8,opt,name=is_outside_etop,json=isOutsideEtop" json:"is_outside_etop"`
	Quantity             int32                 `protobuf:"varint,17,opt,name=quantity" json:"quantity"`
	ListPrice            int32                 `protobuf:"varint,18,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          int32                 `protobuf:"varint,19,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	PaymentPrice         int32                 `protobuf:"varint,20,opt,name=payment_price,json=paymentPrice" json:"payment_price"`
	ImageUrl             string                `protobuf:"bytes,21,opt,name=image_url,json=imageUrl" json:"image_url"`
	Attributes           []*Attribute          `protobuf:"bytes,22,rep,name=attributes" json:"attributes"`
	ProductId            dot.ID                `protobuf:"varint,23,opt,name=product_id,json=productId" json:"product_id"`
	TotalDiscount        int32                 `protobuf:"varint,24,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	MetaFields           []*OrderLineMetaField `protobuf:"bytes,25,rep,name=meta_fields,json=metaFields" json:"meta_fields"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *OrderLine) Reset()         { *m = OrderLine{} }
func (m *OrderLine) String() string { return proto.CompactTextString(m) }
func (*OrderLine) ProtoMessage()    {}
func (*OrderLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{3}
}

var xxx_messageInfo_OrderLine proto.InternalMessageInfo

func (m *OrderLine) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *OrderLine) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *OrderLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *OrderLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *OrderLine) GetIsOutsideEtop() bool {
	if m != nil {
		return m.IsOutsideEtop
	}
	return false
}

func (m *OrderLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *OrderLine) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *OrderLine) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *OrderLine) GetPaymentPrice() int32 {
	if m != nil {
		return m.PaymentPrice
	}
	return 0
}

func (m *OrderLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *OrderLine) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *OrderLine) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *OrderLine) GetTotalDiscount() int32 {
	if m != nil {
		return m.TotalDiscount
	}
	return 0
}

func (m *OrderLine) GetMetaFields() []*OrderLineMetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

type OrderFeeLine struct {
	Type fee.FeeType `protobuf:"varint,1,opt,name=type,enum=fee.FeeType" json:"type"`
	// @required
	Name string `protobuf:"bytes,2,opt,name=name" json:"name"`
	Code string `protobuf:"bytes,3,opt,name=code" json:"code"`
	Desc string `protobuf:"bytes,5,opt,name=desc" json:"desc"`
	// @required
	Amount               int32    `protobuf:"varint,6,opt,name=amount" json:"amount"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderFeeLine) Reset()         { *m = OrderFeeLine{} }
func (m *OrderFeeLine) String() string { return proto.CompactTextString(m) }
func (*OrderFeeLine) ProtoMessage()    {}
func (*OrderFeeLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{4}
}

var xxx_messageInfo_OrderFeeLine proto.InternalMessageInfo

func (m *OrderFeeLine) GetType() fee.FeeType {
	if m != nil {
		return m.Type
	}
	return fee.FeeType_other
}

func (m *OrderFeeLine) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *OrderFeeLine) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *OrderFeeLine) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *OrderFeeLine) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type CreateOrderLine struct {
	VariantId            dot.ID                `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	ProductName          string                `protobuf:"bytes,2,opt,name=product_name,json=productName" json:"product_name"`
	Quantity             int32                 `protobuf:"varint,3,opt,name=quantity" json:"quantity"`
	ListPrice            int32                 `protobuf:"varint,7,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          int32                 `protobuf:"varint,4,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	PaymentPrice         int32                 `protobuf:"varint,5,opt,name=payment_price,json=paymentPrice" json:"payment_price"`
	ImageUrl             string                `protobuf:"bytes,21,opt,name=image_url,json=imageUrl" json:"image_url"`
	Attributes           []*Attribute          `protobuf:"bytes,22,rep,name=attributes" json:"attributes"`
	MetaFields           []*OrderLineMetaField `protobuf:"bytes,23,rep,name=meta_fields,json=metaFields" json:"meta_fields"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *CreateOrderLine) Reset()         { *m = CreateOrderLine{} }
func (m *CreateOrderLine) String() string { return proto.CompactTextString(m) }
func (*CreateOrderLine) ProtoMessage()    {}
func (*CreateOrderLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{5}
}

var xxx_messageInfo_CreateOrderLine proto.InternalMessageInfo

func (m *CreateOrderLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *CreateOrderLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *CreateOrderLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *CreateOrderLine) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *CreateOrderLine) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *CreateOrderLine) GetPaymentPrice() int32 {
	if m != nil {
		return m.PaymentPrice
	}
	return 0
}

func (m *CreateOrderLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *CreateOrderLine) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *CreateOrderLine) GetMetaFields() []*OrderLineMetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

type OrderCustomer struct {
	ExportedFields       []string      `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	FirstName            string        `protobuf:"bytes,1,opt,name=first_name,json=firstName" json:"first_name"`
	LastName             string        `protobuf:"bytes,2,opt,name=last_name,json=lastName" json:"last_name"`
	FullName             string        `protobuf:"bytes,3,opt,name=full_name,json=fullName" json:"full_name"`
	Email                string        `protobuf:"bytes,4,opt,name=email" json:"email"`
	Phone                string        `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Gender               gender.Gender `protobuf:"varint,6,opt,name=gender,enum=gender.Gender" json:"gender"`
	Type                 string        `protobuf:"bytes,7,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *OrderCustomer) Reset()         { *m = OrderCustomer{} }
func (m *OrderCustomer) String() string { return proto.CompactTextString(m) }
func (*OrderCustomer) ProtoMessage()    {}
func (*OrderCustomer) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{6}
}

var xxx_messageInfo_OrderCustomer proto.InternalMessageInfo

func (m *OrderCustomer) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *OrderCustomer) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *OrderCustomer) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *OrderCustomer) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *OrderCustomer) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OrderCustomer) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *OrderCustomer) GetGender() gender.Gender {
	if m != nil {
		return m.Gender
	}
	return gender.Gender_unknown
}

func (m *OrderCustomer) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type OrderAddress struct {
	ExportedFields       []string          `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	FullName             string            `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	FirstName            string            `protobuf:"bytes,2,opt,name=first_name,json=firstName" json:"first_name"`
	LastName             string            `protobuf:"bytes,3,opt,name=last_name,json=lastName" json:"last_name"`
	Phone                string            `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	Email                string            `protobuf:"bytes,17,opt,name=email" json:"email"`
	Country              string            `protobuf:"bytes,5,opt,name=country" json:"country"`
	City                 string            `protobuf:"bytes,6,opt,name=city" json:"city"`
	Province             string            `protobuf:"bytes,7,opt,name=province" json:"province"`
	District             string            `protobuf:"bytes,8,opt,name=district" json:"district"`
	Ward                 string            `protobuf:"bytes,9,opt,name=ward" json:"ward"`
	Zip                  string            `protobuf:"bytes,10,opt,name=zip" json:"zip"`
	Company              string            `protobuf:"bytes,11,opt,name=company" json:"company"`
	Address1             string            `protobuf:"bytes,12,opt,name=address1" json:"address1"`
	Address2             string            `protobuf:"bytes,13,opt,name=address2" json:"address2"`
	ProvinceCode         string            `protobuf:"bytes,14,opt,name=province_code,json=provinceCode" json:"province_code"`
	DistrictCode         string            `protobuf:"bytes,15,opt,name=district_code,json=districtCode" json:"district_code"`
	WardCode             string            `protobuf:"bytes,16,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Coordinates          *etop.Coordinates `protobuf:"bytes,18,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *OrderAddress) Reset()         { *m = OrderAddress{} }
func (m *OrderAddress) String() string { return proto.CompactTextString(m) }
func (*OrderAddress) ProtoMessage()    {}
func (*OrderAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{7}
}

var xxx_messageInfo_OrderAddress proto.InternalMessageInfo

func (m *OrderAddress) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *OrderAddress) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *OrderAddress) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *OrderAddress) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *OrderAddress) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *OrderAddress) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OrderAddress) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *OrderAddress) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *OrderAddress) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *OrderAddress) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *OrderAddress) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *OrderAddress) GetZip() string {
	if m != nil {
		return m.Zip
	}
	return ""
}

func (m *OrderAddress) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *OrderAddress) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *OrderAddress) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *OrderAddress) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *OrderAddress) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *OrderAddress) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *OrderAddress) GetCoordinates() *etop.Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type OrderDiscount struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type"`
	Amount               int32    `protobuf:"varint,3,opt,name=amount" json:"amount"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderDiscount) Reset()         { *m = OrderDiscount{} }
func (m *OrderDiscount) String() string { return proto.CompactTextString(m) }
func (*OrderDiscount) ProtoMessage()    {}
func (*OrderDiscount) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{8}
}

var xxx_messageInfo_OrderDiscount proto.InternalMessageInfo

func (m *OrderDiscount) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *OrderDiscount) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *OrderDiscount) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type CreateOrderRequest struct {
	Source        source.Source `protobuf:"varint,1,opt,name=source,enum=source.Source" json:"source"`
	ExternalId    string        `protobuf:"bytes,2,opt,name=external_id,json=externalId" json:"external_id"`
	ExternalCode  string        `protobuf:"bytes,3,opt,name=external_code,json=externalCode" json:"external_code"`
	ExternalUrl   string        `protobuf:"bytes,5,opt,name=external_url,json=externalUrl" json:"external_url"`
	PaymentMethod string        `protobuf:"bytes,7,opt,name=payment_method,json=paymentMethod" json:"payment_method"`
	// If order_source is self, customer must be shop information
	// and customer_address, shipping_address must be shop address.
	Customer        *OrderCustomer `protobuf:"bytes,9,opt,name=customer" json:"customer"`
	CustomerAddress *OrderAddress  `protobuf:"bytes,10,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	BillingAddress  *OrderAddress  `protobuf:"bytes,11,opt,name=billing_address,json=billingAddress" json:"billing_address"`
	ShippingAddress *OrderAddress  `protobuf:"bytes,12,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	// If there are products from shop, this field should be set.
	// Otherwise, use shop default address.
	ShopAddress *OrderAddress      `protobuf:"bytes,13,opt,name=shop_address,json=shopAddress" json:"shop_address"`
	ShConfirm   *status3.Status    `protobuf:"varint,31,opt,name=sh_confirm,json=shConfirm,enum=status3.Status" json:"sh_confirm"`
	Lines       []*CreateOrderLine `protobuf:"bytes,41,rep,name=lines" json:"lines"`
	Discounts   []*OrderDiscount   `protobuf:"bytes,42,rep,name=discounts" json:"discounts"`
	TotalItems  int32              `protobuf:"varint,43,opt,name=total_items,json=totalItems" json:"total_items"`
	BasketValue int32              `protobuf:"varint,44,opt,name=basket_value,json=basketValue" json:"basket_value"`
	// @deprecated use shipping.gross_weight, shipping.chargeable_weight
	TotalWeight   int32           `protobuf:"varint,45,opt,name=total_weight,json=totalWeight" json:"total_weight"`
	OrderDiscount int32           `protobuf:"varint,46,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	TotalFee      int32           `protobuf:"varint,60,opt,name=total_fee,json=totalFee" json:"total_fee"`
	FeeLines      []*OrderFeeLine `protobuf:"bytes,61,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	TotalDiscount *int32          `protobuf:"varint,47,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount   int32           `protobuf:"varint,48,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	OrderNote     string          `protobuf:"bytes,50,opt,name=order_note,json=orderNote" json:"order_note"`
	ShippingNote  string          `protobuf:"bytes,51,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	// @deprecated use fee_lines.shipping instead
	ShopShippingFee int32 `protobuf:"varint,52,opt,name=shop_shipping_fee,json=shopShippingFee" json:"shop_shipping_fee"`
	// @deprecated use shipping.cod_amount instead
	ShopCod      int32  `protobuf:"varint,53,opt,name=shop_cod,json=shopCod" json:"shop_cod"`
	ReferenceUrl string `protobuf:"bytes,54,opt,name=reference_url,json=referenceUrl" json:"reference_url"`
	// @deprecated use shipping instead
	ShopShipping *OrderShipping `protobuf:"bytes,55,opt,name=shop_shipping,json=shopShipping" json:"shop_shipping"`
	Shipping     *OrderShipping `protobuf:"bytes,57,opt,name=shipping" json:"shipping"`
	// @deprecated use shop_shipping.try_on instead
	GhnNoteCode          ghn_note_code.GHNNoteCode `protobuf:"varint,56,opt,name=ghn_note_code,json=ghnNoteCode,enum=ghn_note_code.GHNNoteCode" json:"ghn_note_code"`
	ExternalMeta         map[string]string         `protobuf:"bytes,62,rep,name=external_meta,json=externalMeta" json:"external_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ReferralMeta         map[string]string         `protobuf:"bytes,63,rep,name=referral_meta,json=referralMeta" json:"referral_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CustomerId           dot.ID                    `protobuf:"varint,64,opt,name=customer_id,json=customerId" json:"customer_id"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CreateOrderRequest) Reset()         { *m = CreateOrderRequest{} }
func (m *CreateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrderRequest) ProtoMessage()    {}
func (*CreateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{9}
}

var xxx_messageInfo_CreateOrderRequest proto.InternalMessageInfo

func (m *CreateOrderRequest) GetSource() source.Source {
	if m != nil {
		return m.Source
	}
	return source.Source_unknown
}

func (m *CreateOrderRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *CreateOrderRequest) GetExternalCode() string {
	if m != nil {
		return m.ExternalCode
	}
	return ""
}

func (m *CreateOrderRequest) GetExternalUrl() string {
	if m != nil {
		return m.ExternalUrl
	}
	return ""
}

func (m *CreateOrderRequest) GetPaymentMethod() string {
	if m != nil {
		return m.PaymentMethod
	}
	return ""
}

func (m *CreateOrderRequest) GetCustomer() *OrderCustomer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *CreateOrderRequest) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetBillingAddress() *OrderAddress {
	if m != nil {
		return m.BillingAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetShopAddress() *OrderAddress {
	if m != nil {
		return m.ShopAddress
	}
	return nil
}

func (m *CreateOrderRequest) GetShConfirm() status3.Status {
	if m != nil && m.ShConfirm != nil {
		return *m.ShConfirm
	}
	return status3.Status_Z
}

func (m *CreateOrderRequest) GetLines() []*CreateOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *CreateOrderRequest) GetDiscounts() []*OrderDiscount {
	if m != nil {
		return m.Discounts
	}
	return nil
}

func (m *CreateOrderRequest) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *CreateOrderRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *CreateOrderRequest) GetTotalWeight() int32 {
	if m != nil {
		return m.TotalWeight
	}
	return 0
}

func (m *CreateOrderRequest) GetOrderDiscount() int32 {
	if m != nil {
		return m.OrderDiscount
	}
	return 0
}

func (m *CreateOrderRequest) GetTotalFee() int32 {
	if m != nil {
		return m.TotalFee
	}
	return 0
}

func (m *CreateOrderRequest) GetFeeLines() []*OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *CreateOrderRequest) GetTotalDiscount() int32 {
	if m != nil && m.TotalDiscount != nil {
		return *m.TotalDiscount
	}
	return 0
}

func (m *CreateOrderRequest) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *CreateOrderRequest) GetOrderNote() string {
	if m != nil {
		return m.OrderNote
	}
	return ""
}

func (m *CreateOrderRequest) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *CreateOrderRequest) GetShopShippingFee() int32 {
	if m != nil {
		return m.ShopShippingFee
	}
	return 0
}

func (m *CreateOrderRequest) GetShopCod() int32 {
	if m != nil {
		return m.ShopCod
	}
	return 0
}

func (m *CreateOrderRequest) GetReferenceUrl() string {
	if m != nil {
		return m.ReferenceUrl
	}
	return ""
}

func (m *CreateOrderRequest) GetShopShipping() *OrderShipping {
	if m != nil {
		return m.ShopShipping
	}
	return nil
}

func (m *CreateOrderRequest) GetShipping() *OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

func (m *CreateOrderRequest) GetGhnNoteCode() ghn_note_code.GHNNoteCode {
	if m != nil {
		return m.GhnNoteCode
	}
	return ghn_note_code.GHNNoteCode_unknown
}

func (m *CreateOrderRequest) GetExternalMeta() map[string]string {
	if m != nil {
		return m.ExternalMeta
	}
	return nil
}

func (m *CreateOrderRequest) GetReferralMeta() map[string]string {
	if m != nil {
		return m.ReferralMeta
	}
	return nil
}

func (m *CreateOrderRequest) GetCustomerId() dot.ID {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

type UpdateOrderRequest struct {
	// @required
	Id              dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Customer        *OrderCustomer `protobuf:"bytes,2,opt,name=customer" json:"customer"`
	CustomerAddress *OrderAddress  `protobuf:"bytes,3,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	BillingAddress  *OrderAddress  `protobuf:"bytes,4,opt,name=billing_address,json=billingAddress" json:"billing_address"`
	ShippingAddress *OrderAddress  `protobuf:"bytes,5,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	ShopAddress     *OrderAddress  `protobuf:"bytes,9,opt,name=shop_address,json=shopAddress" json:"shop_address"`
	OrderNote       string         `protobuf:"bytes,6,opt,name=order_note,json=orderNote" json:"order_note"`
	ShippingNote    string         `protobuf:"bytes,7,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	// @deprecated use fee_lines instead
	ShopShippingFee *int32 `protobuf:"varint,8,opt,name=shop_shipping_fee,json=shopShippingFee" json:"shop_shipping_fee"`
	// @deprecated use shipping.cod_amount instead
	ShopCod *int32 `protobuf:"varint,11,opt,name=shop_cod,json=shopCod" json:"shop_cod"`
	// @deprecated use shipping instead
	ShopShipping  *OrderShipping  `protobuf:"bytes,10,opt,name=shop_shipping,json=shopShipping" json:"shop_shipping"`
	Shipping      *OrderShipping  `protobuf:"bytes,15,opt,name=shipping" json:"shipping"`
	FeeLines      []*OrderFeeLine `protobuf:"bytes,54,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	OrderDiscount *int32          `protobuf:"varint,13,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	// @deprecated
	TotalWeight          int32              `protobuf:"varint,14,opt,name=total_weight,json=totalWeight" json:"total_weight"`
	ChargeableWeight     int32              `protobuf:"varint,19,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	Lines                []*CreateOrderLine `protobuf:"bytes,20,rep,name=lines" json:"lines"`
	BasketValue          int32              `protobuf:"varint,21,opt,name=basket_value,json=basketValue" json:"basket_value"`
	TotalAmount          int32              `protobuf:"varint,22,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	TotalItems           int32              `protobuf:"varint,23,opt,name=total_items,json=totalItems" json:"total_items"`
	TotalFee             *int32             `protobuf:"varint,24,opt,name=total_fee,json=totalFee" json:"total_fee"`
	CustomerId           dot.ID             `protobuf:"varint,25,opt,name=customer_id,json=customerId" json:"customer_id"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *UpdateOrderRequest) Reset()         { *m = UpdateOrderRequest{} }
func (m *UpdateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateOrderRequest) ProtoMessage()    {}
func (*UpdateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{10}
}

var xxx_messageInfo_UpdateOrderRequest proto.InternalMessageInfo

func (m *UpdateOrderRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateOrderRequest) GetCustomer() *OrderCustomer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *UpdateOrderRequest) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *UpdateOrderRequest) GetBillingAddress() *OrderAddress {
	if m != nil {
		return m.BillingAddress
	}
	return nil
}

func (m *UpdateOrderRequest) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *UpdateOrderRequest) GetShopAddress() *OrderAddress {
	if m != nil {
		return m.ShopAddress
	}
	return nil
}

func (m *UpdateOrderRequest) GetOrderNote() string {
	if m != nil {
		return m.OrderNote
	}
	return ""
}

func (m *UpdateOrderRequest) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *UpdateOrderRequest) GetShopShippingFee() int32 {
	if m != nil && m.ShopShippingFee != nil {
		return *m.ShopShippingFee
	}
	return 0
}

func (m *UpdateOrderRequest) GetShopCod() int32 {
	if m != nil && m.ShopCod != nil {
		return *m.ShopCod
	}
	return 0
}

func (m *UpdateOrderRequest) GetShopShipping() *OrderShipping {
	if m != nil {
		return m.ShopShipping
	}
	return nil
}

func (m *UpdateOrderRequest) GetShipping() *OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

func (m *UpdateOrderRequest) GetFeeLines() []*OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *UpdateOrderRequest) GetOrderDiscount() int32 {
	if m != nil && m.OrderDiscount != nil {
		return *m.OrderDiscount
	}
	return 0
}

func (m *UpdateOrderRequest) GetTotalWeight() int32 {
	if m != nil {
		return m.TotalWeight
	}
	return 0
}

func (m *UpdateOrderRequest) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *UpdateOrderRequest) GetLines() []*CreateOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *UpdateOrderRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *UpdateOrderRequest) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *UpdateOrderRequest) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *UpdateOrderRequest) GetTotalFee() int32 {
	if m != nil && m.TotalFee != nil {
		return *m.TotalFee
	}
	return 0
}

func (m *UpdateOrderRequest) GetCustomerId() dot.ID {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

type OrderShipping struct {
	ExportedFields []string `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	// @deprecated use pickup_address
	ShAddress *OrderAddress `protobuf:"bytes,1,opt,name=sh_address,json=shAddress" json:"sh_address"`
	// @deprecated
	XServiceId string `protobuf:"bytes,2,opt,name=x_service_id,json=xServiceId" json:"x_service_id"`
	// @deprecated
	XShippingFee int32 `protobuf:"varint,3,opt,name=x_shipping_fee,json=xShippingFee" json:"x_shipping_fee"`
	// @deprecated
	XServiceName        string        `protobuf:"bytes,4,opt,name=x_service_name,json=xServiceName" json:"x_service_name"`
	PickupAddress       *OrderAddress `protobuf:"bytes,21,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ReturnAddress       *OrderAddress `protobuf:"bytes,26,opt,name=return_address,json=returnAddress" json:"return_address"`
	ShippingServiceName string        `protobuf:"bytes,24,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceCode string        `protobuf:"bytes,22,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee  int32         `protobuf:"varint,23,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	// @deprecated use carrier
	ShippingProvider shipping_provider.ShippingProvider `protobuf:"varint,5,opt,name=shipping_provider,json=shippingProvider,enum=shipping_provider.ShippingProvider" json:"shipping_provider"`
	Carrier          shipping_provider.ShippingProvider `protobuf:"varint,14,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	IncludeInsurance bool                               `protobuf:"varint,6,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	TryOn            try_on.TryOnCode                   `protobuf:"varint,7,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	ShippingNote     string                             `protobuf:"bytes,8,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	CodAmount        *int32                             `protobuf:"varint,9,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	// @deprecated
	Weight               *int32   `protobuf:"varint,25,opt,name=weight" json:"weight"`
	GrossWeight          *int32   `protobuf:"varint,10,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	Length               *int32   `protobuf:"varint,11,opt,name=length" json:"length"`
	Width                *int32   `protobuf:"varint,12,opt,name=width" json:"width"`
	Height               *int32   `protobuf:"varint,13,opt,name=height" json:"height"`
	ChargeableWeight     *int32   `protobuf:"varint,15,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderShipping) Reset()         { *m = OrderShipping{} }
func (m *OrderShipping) String() string { return proto.CompactTextString(m) }
func (*OrderShipping) ProtoMessage()    {}
func (*OrderShipping) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{11}
}

var xxx_messageInfo_OrderShipping proto.InternalMessageInfo

func (m *OrderShipping) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *OrderShipping) GetShAddress() *OrderAddress {
	if m != nil {
		return m.ShAddress
	}
	return nil
}

func (m *OrderShipping) GetXServiceId() string {
	if m != nil {
		return m.XServiceId
	}
	return ""
}

func (m *OrderShipping) GetXShippingFee() int32 {
	if m != nil {
		return m.XShippingFee
	}
	return 0
}

func (m *OrderShipping) GetXServiceName() string {
	if m != nil {
		return m.XServiceName
	}
	return ""
}

func (m *OrderShipping) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *OrderShipping) GetReturnAddress() *OrderAddress {
	if m != nil {
		return m.ReturnAddress
	}
	return nil
}

func (m *OrderShipping) GetShippingServiceName() string {
	if m != nil {
		return m.ShippingServiceName
	}
	return ""
}

func (m *OrderShipping) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *OrderShipping) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *OrderShipping) GetShippingProvider() shipping_provider.ShippingProvider {
	if m != nil {
		return m.ShippingProvider
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *OrderShipping) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *OrderShipping) GetIncludeInsurance() bool {
	if m != nil {
		return m.IncludeInsurance
	}
	return false
}

func (m *OrderShipping) GetTryOn() try_on.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *OrderShipping) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *OrderShipping) GetCodAmount() int32 {
	if m != nil && m.CodAmount != nil {
		return *m.CodAmount
	}
	return 0
}

func (m *OrderShipping) GetWeight() int32 {
	if m != nil && m.Weight != nil {
		return *m.Weight
	}
	return 0
}

func (m *OrderShipping) GetGrossWeight() int32 {
	if m != nil && m.GrossWeight != nil {
		return *m.GrossWeight
	}
	return 0
}

func (m *OrderShipping) GetLength() int32 {
	if m != nil && m.Length != nil {
		return *m.Length
	}
	return 0
}

func (m *OrderShipping) GetWidth() int32 {
	if m != nil && m.Width != nil {
		return *m.Width
	}
	return 0
}

func (m *OrderShipping) GetHeight() int32 {
	if m != nil && m.Height != nil {
		return *m.Height
	}
	return 0
}

func (m *OrderShipping) GetChargeableWeight() int32 {
	if m != nil && m.ChargeableWeight != nil {
		return *m.ChargeableWeight
	}
	return 0
}

type Fulfillment struct {
	ExportedFields []string     `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	Id             dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	OrderId        dot.ID       `protobuf:"varint,2,opt,name=order_id,json=orderId" json:"order_id"`
	ShopId         dot.ID       `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	PartnerId      dot.ID       `protobuf:"varint,55,opt,name=partner_id,json=partnerId" json:"partner_id"`
	SelfUrl        string       `protobuf:"bytes,59,opt,name=self_url,json=selfUrl" json:"self_url"`
	Lines          []*OrderLine `protobuf:"bytes,5,rep,name=lines" json:"lines"`
	TotalItems     int32        `protobuf:"varint,6,opt,name=total_items,json=totalItems" json:"total_items"`
	// @deprecated use chargeable_weight
	TotalWeight int32 `protobuf:"varint,7,opt,name=total_weight,json=totalWeight" json:"total_weight"`
	BasketValue int32 `protobuf:"varint,10,opt,name=basket_value,json=basketValue" json:"basket_value"`
	// @deprecated use cod_amount
	TotalCodAmount int32 `protobuf:"varint,9,opt,name=total_cod_amount,json=totalCodAmount" json:"total_cod_amount"`
	CodAmount      int32 `protobuf:"varint,60,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	// @deprecated
	TotalAmount      int32    `protobuf:"varint,8,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	ChargeableWeight int32    `protobuf:"varint,72,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	CreatedAt        dot.Time `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt        dot.Time `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt         dot.Time `protobuf:"bytes,14,opt,name=closed_at,json=closedAt" json:"closed_at"`
	CancelledAt      dot.Time `protobuf:"bytes,13,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CancelReason     string   `protobuf:"bytes,15,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	// @deprecated use carrier instead
	ShippingProvider     string                             `protobuf:"bytes,17,opt,name=shipping_provider,json=shippingProvider" json:"shipping_provider"`
	Carrier              shipping_provider.ShippingProvider `protobuf:"varint,71,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	ShippingServiceName  string                             `protobuf:"bytes,70,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceFee   int32                              `protobuf:"varint,68,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingServiceCode  string                             `protobuf:"bytes,69,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingCode         string                             `protobuf:"bytes,18,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	ShippingNote         string                             `protobuf:"bytes,19,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	TryOn                try_on.TryOnCode                   `protobuf:"varint,20,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	IncludeInsurance     bool                               `protobuf:"varint,63,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	ShConfirm            status3.Status                     `protobuf:"varint,23,opt,name=sh_confirm,json=shConfirm,enum=status3.Status" json:"sh_confirm"`
	ShippingState        shipping.State                     `protobuf:"varint,25,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	Status               status5.Status                     `protobuf:"varint,26,opt,name=status,enum=status5.Status" json:"status"`
	ShippingStatus       status5.Status                     `protobuf:"varint,64,opt,name=shipping_status,json=shippingStatus,enum=status5.Status" json:"shipping_status"`
	EtopPaymentStatus    status4.Status                     `protobuf:"varint,22,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	ShippingFeeCustomer  int32                              `protobuf:"varint,27,opt,name=shipping_fee_customer,json=shippingFeeCustomer" json:"shipping_fee_customer"`
	ShippingFeeShop      int32                              `protobuf:"varint,28,opt,name=shipping_fee_shop,json=shippingFeeShop" json:"shipping_fee_shop"`
	XShippingFee         int32                              `protobuf:"varint,29,opt,name=x_shipping_fee,json=xShippingFee" json:"x_shipping_fee"`
	XShippingId          string                             `protobuf:"bytes,30,opt,name=x_shipping_id,json=xShippingId" json:"x_shipping_id"`
	XShippingCode        string                             `protobuf:"bytes,31,opt,name=x_shipping_code,json=xShippingCode" json:"x_shipping_code"`
	XShippingCreatedAt   dot.Time                           `protobuf:"bytes,32,opt,name=x_shipping_created_at,json=xShippingCreatedAt" json:"x_shipping_created_at"`
	XShippingUpdatedAt   dot.Time                           `protobuf:"bytes,33,opt,name=x_shipping_updated_at,json=xShippingUpdatedAt" json:"x_shipping_updated_at"`
	XShippingCancelledAt dot.Time                           `protobuf:"bytes,34,opt,name=x_shipping_cancelled_at,json=xShippingCancelledAt" json:"x_shipping_cancelled_at"`
	XShippingDeliveredAt dot.Time                           `protobuf:"bytes,35,opt,name=x_shipping_delivered_at,json=xShippingDeliveredAt" json:"x_shipping_delivered_at"`
	XShippingReturnedAt  dot.Time                           `protobuf:"bytes,36,opt,name=x_shipping_returned_at,json=xShippingReturnedAt" json:"x_shipping_returned_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt dot.Time `protobuf:"bytes,41,opt,name=expected_delivery_at,json=expectedDeliveryAt" json:"expected_delivery_at"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt              dot.Time               `protobuf:"bytes,51,opt,name=expected_pick_at,json=expectedPickAt" json:"expected_pick_at"`
	EstimatedDeliveryAt         dot.Time               `protobuf:"bytes,74,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	EstimatedPickupAt           dot.Time               `protobuf:"bytes,75,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	CodEtopTransferedAt         dot.Time               `protobuf:"bytes,43,opt,name=cod_etop_transfered_at,json=codEtopTransferedAt" json:"cod_etop_transfered_at"`
	ShippingFeeShopTransferedAt dot.Time               `protobuf:"bytes,44,opt,name=shipping_fee_shop_transfered_at,json=shippingFeeShopTransferedAt" json:"shipping_fee_shop_transfered_at"`
	XShippingState              string                 `protobuf:"bytes,37,opt,name=x_shipping_state,json=xShippingState" json:"x_shipping_state"`
	XShippingStatus             status5.Status         `protobuf:"varint,38,opt,name=x_shipping_status,json=xShippingStatus,enum=status5.Status" json:"x_shipping_status"`
	XSyncStatus                 status4.Status         `protobuf:"varint,39,opt,name=x_sync_status,json=xSyncStatus,enum=status4.Status" json:"x_sync_status"`
	XSyncStates                 *FulfillmentSyncStates `protobuf:"bytes,40,opt,name=x_sync_states,json=xSyncStates" json:"x_sync_states"`
	// @deprecated use shipping_address instead
	AddressTo *etop.Address `protobuf:"bytes,42,opt,name=address_to,json=addressTo" json:"address_to"`
	// @deprecated use pickup_address instead
	AddressFrom                        *etop.Address          `protobuf:"bytes,50,opt,name=address_from,json=addressFrom" json:"address_from"`
	PickupAddress                      *OrderAddress          `protobuf:"bytes,65,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ReturnAddress                      *OrderAddress          `protobuf:"bytes,66,opt,name=return_address,json=returnAddress" json:"return_address"`
	ShippingAddress                    *OrderAddress          `protobuf:"bytes,67,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	Shop                               *etop.Shop             `protobuf:"bytes,45,opt,name=shop" json:"shop"`
	Order                              *Order                 `protobuf:"bytes,46,opt,name=order" json:"order"`
	ProviderShippingFeeLines           []*ShippingFeeLine     `protobuf:"bytes,47,rep,name=provider_shipping_fee_lines,json=providerShippingFeeLines" json:"provider_shipping_fee_lines"`
	ShippingFeeShopLines               []*ShippingFeeLine     `protobuf:"bytes,48,rep,name=shipping_fee_shop_lines,json=shippingFeeShopLines" json:"shipping_fee_shop_lines"`
	EtopDiscount                       int32                  `protobuf:"varint,49,opt,name=etop_discount,json=etopDiscount" json:"etop_discount"`
	MoneyTransactionShippingId         dot.ID                 `protobuf:"varint,52,opt,name=money_transaction_shipping_id,json=moneyTransactionShippingId" json:"money_transaction_shipping_id"`
	MoneyTransactionShippingExternalId dot.ID                 `protobuf:"varint,53,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	XShippingLogs                      []*ExternalShippingLog `protobuf:"bytes,54,rep,name=x_shipping_logs,json=xShippingLogs" json:"x_shipping_logs"`
	XShippingNote                      string                 `protobuf:"bytes,56,opt,name=x_shipping_note,json=xShippingNote" json:"x_shipping_note"`
	XShippingSubState                  string                 `protobuf:"bytes,57,opt,name=x_shipping_sub_state,json=xShippingSubState" json:"x_shipping_sub_state"`
	Code                               string                 `protobuf:"bytes,58,opt,name=code" json:"code"`
	ActualCompensationAmount           int32                  `protobuf:"varint,76,opt,name=actual_compensation_amount,json=actualCompensationAmount" json:"actual_compensation_amount"`
	XXX_NoUnkeyedLiteral               struct{}               `json:"-"`
	XXX_sizecache                      int32                  `json:"-"`
}

func (m *Fulfillment) Reset()         { *m = Fulfillment{} }
func (m *Fulfillment) String() string { return proto.CompactTextString(m) }
func (*Fulfillment) ProtoMessage()    {}
func (*Fulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{12}
}

var xxx_messageInfo_Fulfillment proto.InternalMessageInfo

func (m *Fulfillment) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *Fulfillment) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Fulfillment) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *Fulfillment) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Fulfillment) GetPartnerId() dot.ID {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *Fulfillment) GetSelfUrl() string {
	if m != nil {
		return m.SelfUrl
	}
	return ""
}

func (m *Fulfillment) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *Fulfillment) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *Fulfillment) GetTotalWeight() int32 {
	if m != nil {
		return m.TotalWeight
	}
	return 0
}

func (m *Fulfillment) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *Fulfillment) GetTotalCodAmount() int32 {
	if m != nil {
		return m.TotalCodAmount
	}
	return 0
}

func (m *Fulfillment) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *Fulfillment) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *Fulfillment) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *Fulfillment) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *Fulfillment) GetShippingProvider() string {
	if m != nil {
		return m.ShippingProvider
	}
	return ""
}

func (m *Fulfillment) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *Fulfillment) GetShippingServiceName() string {
	if m != nil {
		return m.ShippingServiceName
	}
	return ""
}

func (m *Fulfillment) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *Fulfillment) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *Fulfillment) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

func (m *Fulfillment) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *Fulfillment) GetTryOn() try_on.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *Fulfillment) GetIncludeInsurance() bool {
	if m != nil {
		return m.IncludeInsurance
	}
	return false
}

func (m *Fulfillment) GetShConfirm() status3.Status {
	if m != nil {
		return m.ShConfirm
	}
	return status3.Status_Z
}

func (m *Fulfillment) GetShippingState() shipping.State {
	if m != nil {
		return m.ShippingState
	}
	return shipping.State_default
}

func (m *Fulfillment) GetStatus() status5.Status {
	if m != nil {
		return m.Status
	}
	return status5.Status_Z
}

func (m *Fulfillment) GetShippingStatus() status5.Status {
	if m != nil {
		return m.ShippingStatus
	}
	return status5.Status_Z
}

func (m *Fulfillment) GetEtopPaymentStatus() status4.Status {
	if m != nil {
		return m.EtopPaymentStatus
	}
	return status4.Status_Z
}

func (m *Fulfillment) GetShippingFeeCustomer() int32 {
	if m != nil {
		return m.ShippingFeeCustomer
	}
	return 0
}

func (m *Fulfillment) GetShippingFeeShop() int32 {
	if m != nil {
		return m.ShippingFeeShop
	}
	return 0
}

func (m *Fulfillment) GetXShippingFee() int32 {
	if m != nil {
		return m.XShippingFee
	}
	return 0
}

func (m *Fulfillment) GetXShippingId() string {
	if m != nil {
		return m.XShippingId
	}
	return ""
}

func (m *Fulfillment) GetXShippingCode() string {
	if m != nil {
		return m.XShippingCode
	}
	return ""
}

func (m *Fulfillment) GetXShippingState() string {
	if m != nil {
		return m.XShippingState
	}
	return ""
}

func (m *Fulfillment) GetXShippingStatus() status5.Status {
	if m != nil {
		return m.XShippingStatus
	}
	return status5.Status_Z
}

func (m *Fulfillment) GetXSyncStatus() status4.Status {
	if m != nil {
		return m.XSyncStatus
	}
	return status4.Status_Z
}

func (m *Fulfillment) GetXSyncStates() *FulfillmentSyncStates {
	if m != nil {
		return m.XSyncStates
	}
	return nil
}

func (m *Fulfillment) GetAddressTo() *etop.Address {
	if m != nil {
		return m.AddressTo
	}
	return nil
}

func (m *Fulfillment) GetAddressFrom() *etop.Address {
	if m != nil {
		return m.AddressFrom
	}
	return nil
}

func (m *Fulfillment) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *Fulfillment) GetReturnAddress() *OrderAddress {
	if m != nil {
		return m.ReturnAddress
	}
	return nil
}

func (m *Fulfillment) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *Fulfillment) GetShop() *etop.Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *Fulfillment) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *Fulfillment) GetProviderShippingFeeLines() []*ShippingFeeLine {
	if m != nil {
		return m.ProviderShippingFeeLines
	}
	return nil
}

func (m *Fulfillment) GetShippingFeeShopLines() []*ShippingFeeLine {
	if m != nil {
		return m.ShippingFeeShopLines
	}
	return nil
}

func (m *Fulfillment) GetEtopDiscount() int32 {
	if m != nil {
		return m.EtopDiscount
	}
	return 0
}

func (m *Fulfillment) GetMoneyTransactionShippingId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingId
	}
	return 0
}

func (m *Fulfillment) GetMoneyTransactionShippingExternalId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingExternalId
	}
	return 0
}

func (m *Fulfillment) GetXShippingLogs() []*ExternalShippingLog {
	if m != nil {
		return m.XShippingLogs
	}
	return nil
}

func (m *Fulfillment) GetXShippingNote() string {
	if m != nil {
		return m.XShippingNote
	}
	return ""
}

func (m *Fulfillment) GetXShippingSubState() string {
	if m != nil {
		return m.XShippingSubState
	}
	return ""
}

func (m *Fulfillment) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Fulfillment) GetActualCompensationAmount() int32 {
	if m != nil {
		return m.ActualCompensationAmount
	}
	return 0
}

type ShippingFeeLine struct {
	ShippingFeeType          shipping_fee_type.ShippingFeeType `protobuf:"varint,1,opt,name=shipping_fee_type,json=shippingFeeType,enum=shipping_fee_type.ShippingFeeType" json:"shipping_fee_type"`
	Cost                     int32                             `protobuf:"varint,2,opt,name=cost" json:"cost"`
	ExternalServiceId        string                            `protobuf:"bytes,3,opt,name=external_service_id,json=externalServiceId" json:"external_service_id"`
	ExternalServiceName      string                            `protobuf:"bytes,4,opt,name=external_service_name,json=externalServiceName" json:"external_service_name"`
	ExternalServiceType      string                            `protobuf:"bytes,5,opt,name=external_service_type,json=externalServiceType" json:"external_service_type"`
	ExternalShippingOrderId  string                            `protobuf:"bytes,6,opt,name=external_shipping_order_id,json=externalShippingOrderId" json:"external_shipping_order_id"`
	ExternalPaymentChannelId string                            `protobuf:"bytes,7,opt,name=external_payment_channel_id,json=externalPaymentChannelId" json:"external_payment_channel_id"`
	XXX_NoUnkeyedLiteral     struct{}                          `json:"-"`
	XXX_sizecache            int32                             `json:"-"`
}

func (m *ShippingFeeLine) Reset()         { *m = ShippingFeeLine{} }
func (m *ShippingFeeLine) String() string { return proto.CompactTextString(m) }
func (*ShippingFeeLine) ProtoMessage()    {}
func (*ShippingFeeLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{13}
}

var xxx_messageInfo_ShippingFeeLine proto.InternalMessageInfo

func (m *ShippingFeeLine) GetShippingFeeType() shipping_fee_type.ShippingFeeType {
	if m != nil {
		return m.ShippingFeeType
	}
	return shipping_fee_type.ShippingFeeType_main
}

func (m *ShippingFeeLine) GetCost() int32 {
	if m != nil {
		return m.Cost
	}
	return 0
}

func (m *ShippingFeeLine) GetExternalServiceId() string {
	if m != nil {
		return m.ExternalServiceId
	}
	return ""
}

func (m *ShippingFeeLine) GetExternalServiceName() string {
	if m != nil {
		return m.ExternalServiceName
	}
	return ""
}

func (m *ShippingFeeLine) GetExternalServiceType() string {
	if m != nil {
		return m.ExternalServiceType
	}
	return ""
}

func (m *ShippingFeeLine) GetExternalShippingOrderId() string {
	if m != nil {
		return m.ExternalShippingOrderId
	}
	return ""
}

func (m *ShippingFeeLine) GetExternalPaymentChannelId() string {
	if m != nil {
		return m.ExternalPaymentChannelId
	}
	return ""
}

type ExternalShippingLog struct {
	StateText            string   `protobuf:"bytes,1,opt,name=state_text,json=stateText" json:"state_text"`
	Time                 string   `protobuf:"bytes,2,opt,name=time" json:"time"`
	Message              string   `protobuf:"bytes,3,opt,name=message" json:"message"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExternalShippingLog) Reset()         { *m = ExternalShippingLog{} }
func (m *ExternalShippingLog) String() string { return proto.CompactTextString(m) }
func (*ExternalShippingLog) ProtoMessage()    {}
func (*ExternalShippingLog) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{14}
}

var xxx_messageInfo_ExternalShippingLog proto.InternalMessageInfo

func (m *ExternalShippingLog) GetStateText() string {
	if m != nil {
		return m.StateText
	}
	return ""
}

func (m *ExternalShippingLog) GetTime() string {
	if m != nil {
		return m.Time
	}
	return ""
}

func (m *ExternalShippingLog) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type FulfillmentsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Fulfillments         []*Fulfillment   `protobuf:"bytes,2,rep,name=fulfillments" json:"fulfillments"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *FulfillmentsResponse) Reset()         { *m = FulfillmentsResponse{} }
func (m *FulfillmentsResponse) String() string { return proto.CompactTextString(m) }
func (*FulfillmentsResponse) ProtoMessage()    {}
func (*FulfillmentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{15}
}

var xxx_messageInfo_FulfillmentsResponse proto.InternalMessageInfo

func (m *FulfillmentsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *FulfillmentsResponse) GetFulfillments() []*Fulfillment {
	if m != nil {
		return m.Fulfillments
	}
	return nil
}

type Attribute struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Value                string   `protobuf:"bytes,2,opt,name=value" json:"value"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Attribute) Reset()         { *m = Attribute{} }
func (m *Attribute) String() string { return proto.CompactTextString(m) }
func (*Attribute) ProtoMessage()    {}
func (*Attribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{16}
}

var xxx_messageInfo_Attribute proto.InternalMessageInfo

func (m *Attribute) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type OrderWithErrorsResponse struct {
	// @deprecated
	Errors               []*common.Error `protobuf:"bytes,1,rep,name=errors" json:"errors"`
	Order                *Order          `protobuf:"bytes,2,opt,name=order" json:"order"`
	FulfillmentErrors    []*common.Error `protobuf:"bytes,3,rep,name=fulfillment_errors,json=fulfillmentErrors" json:"fulfillment_errors"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *OrderWithErrorsResponse) Reset()         { *m = OrderWithErrorsResponse{} }
func (m *OrderWithErrorsResponse) String() string { return proto.CompactTextString(m) }
func (*OrderWithErrorsResponse) ProtoMessage()    {}
func (*OrderWithErrorsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{17}
}

var xxx_messageInfo_OrderWithErrorsResponse proto.InternalMessageInfo

func (m *OrderWithErrorsResponse) GetErrors() []*common.Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

func (m *OrderWithErrorsResponse) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *OrderWithErrorsResponse) GetFulfillmentErrors() []*common.Error {
	if m != nil {
		return m.FulfillmentErrors
	}
	return nil
}

type GetExternalShippingServicesRequest struct {
	// @deprecated use carrier instead
	Provider         shipping_provider.ShippingProvider `protobuf:"varint,1,opt,name=provider,enum=shipping_provider.ShippingProvider" json:"provider"`
	Carrier          shipping_provider.ShippingProvider `protobuf:"varint,18,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	FromDistrictCode string                             `protobuf:"bytes,2,opt,name=from_district_code,json=fromDistrictCode" json:"from_district_code"`
	FromProvinceCode string                             `protobuf:"bytes,8,opt,name=from_province_code,json=fromProvinceCode" json:"from_province_code"`
	ToDistrictCode   string                             `protobuf:"bytes,3,opt,name=to_district_code,json=toDistrictCode" json:"to_district_code"`
	ToProvinceCode   string                             `protobuf:"bytes,9,opt,name=to_province_code,json=toProvinceCode" json:"to_province_code"`
	FromProvince     string                             `protobuf:"bytes,12,opt,name=from_province,json=fromProvince" json:"from_province"`
	FromDistrict     string                             `protobuf:"bytes,13,opt,name=from_district,json=fromDistrict" json:"from_district"`
	ToProvince       string                             `protobuf:"bytes,14,opt,name=to_province,json=toProvince" json:"to_province"`
	ToDistrict       string                             `protobuf:"bytes,15,opt,name=to_district,json=toDistrict" json:"to_district"`
	// @deprecated use gross_weight instead
	Weight           int32 `protobuf:"varint,4,opt,name=weight" json:"weight"`
	GrossWeight      int32 `protobuf:"varint,19,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	ChargeableWeight int32 `protobuf:"varint,17,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	Length           int32 `protobuf:"varint,5,opt,name=length" json:"length"`
	Width            int32 `protobuf:"varint,6,opt,name=width" json:"width"`
	Height           int32 `protobuf:"varint,7,opt,name=height" json:"height"`
	// @deprecated use basket_value instead
	Value int32 `protobuf:"varint,10,opt,name=value" json:"value"`
	// @deprecated use cod_amount instead
	TotalCodAmount       int32    `protobuf:"varint,11,opt,name=total_cod_amount,json=totalCodAmount" json:"total_cod_amount"`
	CodAmount            int32    `protobuf:"varint,16,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	BasketValue          int32    `protobuf:"varint,21,opt,name=basket_value,json=basketValue" json:"basket_value"`
	IncludeInsurance     *bool    `protobuf:"varint,20,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetExternalShippingServicesRequest) Reset()         { *m = GetExternalShippingServicesRequest{} }
func (m *GetExternalShippingServicesRequest) String() string { return proto.CompactTextString(m) }
func (*GetExternalShippingServicesRequest) ProtoMessage()    {}
func (*GetExternalShippingServicesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{18}
}

var xxx_messageInfo_GetExternalShippingServicesRequest proto.InternalMessageInfo

func (m *GetExternalShippingServicesRequest) GetProvider() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Provider
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *GetExternalShippingServicesRequest) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *GetExternalShippingServicesRequest) GetFromDistrictCode() string {
	if m != nil {
		return m.FromDistrictCode
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetFromProvinceCode() string {
	if m != nil {
		return m.FromProvinceCode
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetToDistrictCode() string {
	if m != nil {
		return m.ToDistrictCode
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetToProvinceCode() string {
	if m != nil {
		return m.ToProvinceCode
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetFromProvince() string {
	if m != nil {
		return m.FromProvince
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetFromDistrict() string {
	if m != nil {
		return m.FromDistrict
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetToProvince() string {
	if m != nil {
		return m.ToProvince
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetToDistrict() string {
	if m != nil {
		return m.ToDistrict
	}
	return ""
}

func (m *GetExternalShippingServicesRequest) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetGrossWeight() int32 {
	if m != nil {
		return m.GrossWeight
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetTotalCodAmount() int32 {
	if m != nil {
		return m.TotalCodAmount
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *GetExternalShippingServicesRequest) GetIncludeInsurance() bool {
	if m != nil && m.IncludeInsurance != nil {
		return *m.IncludeInsurance
	}
	return false
}

type GetExternalShippingServicesResponse struct {
	Services             []*ExternalShippingService `protobuf:"bytes,1,rep,name=services" json:"services"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *GetExternalShippingServicesResponse) Reset()         { *m = GetExternalShippingServicesResponse{} }
func (m *GetExternalShippingServicesResponse) String() string { return proto.CompactTextString(m) }
func (*GetExternalShippingServicesResponse) ProtoMessage()    {}
func (*GetExternalShippingServicesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{19}
}

var xxx_messageInfo_GetExternalShippingServicesResponse proto.InternalMessageInfo

func (m *GetExternalShippingServicesResponse) GetServices() []*ExternalShippingService {
	if m != nil {
		return m.Services
	}
	return nil
}

type ExternalShippingService struct {
	ExportedFields []string `protobuf:"bytes,100,rep,name=exported_fields,json=exportedFields" json:"exported_fields"`
	// @deprecated use code
	ExternalId string `protobuf:"bytes,1,opt,name=external_id,json=externalId" json:"external_id"`
	// @deprecated use fee
	ServiceFee int32 `protobuf:"varint,3,opt,name=service_fee,json=serviceFee" json:"service_fee"`
	// @deprecated use carier
	Provider shipping_provider.ShippingProvider `protobuf:"varint,4,opt,name=provider,enum=shipping_provider.ShippingProvider" json:"provider"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt dot.Time `protobuf:"bytes,5,opt,name=expected_pick_at,json=expectedPickAt" json:"expected_pick_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt   dot.Time                           `protobuf:"bytes,6,opt,name=expected_delivery_at,json=expectedDeliveryAt" json:"expected_delivery_at"`
	Name                 string                             `protobuf:"bytes,2,opt,name=name" json:"name"`
	Code                 string                             `protobuf:"bytes,11,opt,name=code" json:"code"`
	Fee                  int32                              `protobuf:"varint,13,opt,name=fee" json:"fee"`
	Carrier              shipping_provider.ShippingProvider `protobuf:"varint,14,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	EstimatedPickupAt    dot.Time                           `protobuf:"bytes,15,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	EstimatedDeliveryAt  dot.Time                           `protobuf:"bytes,16,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *ExternalShippingService) Reset()         { *m = ExternalShippingService{} }
func (m *ExternalShippingService) String() string { return proto.CompactTextString(m) }
func (*ExternalShippingService) ProtoMessage()    {}
func (*ExternalShippingService) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{20}
}

var xxx_messageInfo_ExternalShippingService proto.InternalMessageInfo

func (m *ExternalShippingService) GetExportedFields() []string {
	if m != nil {
		return m.ExportedFields
	}
	return nil
}

func (m *ExternalShippingService) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *ExternalShippingService) GetServiceFee() int32 {
	if m != nil {
		return m.ServiceFee
	}
	return 0
}

func (m *ExternalShippingService) GetProvider() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Provider
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *ExternalShippingService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ExternalShippingService) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ExternalShippingService) GetFee() int32 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func (m *ExternalShippingService) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

type FulfillmentSyncStates struct {
	SyncAt               dot.Time      `protobuf:"bytes,1,opt,name=sync_at,json=syncAt" json:"sync_at"`
	NextShippingState    string        `protobuf:"bytes,2,opt,name=next_shipping_state,json=nextShippingState" json:"next_shipping_state"`
	Error                *common.Error `protobuf:"bytes,3,opt,name=error" json:"error"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *FulfillmentSyncStates) Reset()         { *m = FulfillmentSyncStates{} }
func (m *FulfillmentSyncStates) String() string { return proto.CompactTextString(m) }
func (*FulfillmentSyncStates) ProtoMessage()    {}
func (*FulfillmentSyncStates) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{21}
}

var xxx_messageInfo_FulfillmentSyncStates proto.InternalMessageInfo

func (m *FulfillmentSyncStates) GetNextShippingState() string {
	if m != nil {
		return m.NextShippingState
	}
	return ""
}

func (m *FulfillmentSyncStates) GetError() *common.Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type MoneyTransaction struct {
	Id                                 dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId                             dot.ID            `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Status                             status3.Status    `protobuf:"varint,4,opt,name=status,enum=status3.Status" json:"status"`
	TotalCod                           int               `protobuf:"varint,5,opt,name=total_cod,json=totalCod" json:"total_cod"`
	TotalOrders                        int               `protobuf:"varint,6,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	Code                               string            `protobuf:"bytes,7,opt,name=code" json:"code"`
	Provider                           string            `protobuf:"bytes,9,opt,name=provider" json:"provider"`
	MoneyTransactionShippingExternalId dot.ID            `protobuf:"varint,14,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	MoneyTransactionShippingEtopId     dot.ID            `protobuf:"varint,17,opt,name=money_transaction_shipping_etop_id,json=moneyTransactionShippingEtopId" json:"money_transaction_shipping_etop_id"`
	TotalAmount                        int               `protobuf:"varint,15,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	CreatedAt                          dot.Time          `protobuf:"bytes,10,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                          dot.Time          `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt                           dot.Time          `protobuf:"bytes,12,opt,name=closed_at,json=closedAt" json:"closed_at"`
	ConfirmedAt                        dot.Time          `protobuf:"bytes,13,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	EtopTransferedAt                   dot.Time          `protobuf:"bytes,16,opt,name=etop_transfered_at,json=etopTransferedAt" json:"etop_transfered_at"`
	Note                               string            `protobuf:"bytes,18,opt,name=note" json:"note"`
	InvoiceNumber                      string            `protobuf:"bytes,19,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount                        *etop.BankAccount `protobuf:"bytes,20,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral               struct{}          `json:"-"`
	XXX_sizecache                      int32             `json:"-"`
}

func (m *MoneyTransaction) Reset()         { *m = MoneyTransaction{} }
func (m *MoneyTransaction) String() string { return proto.CompactTextString(m) }
func (*MoneyTransaction) ProtoMessage()    {}
func (*MoneyTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{22}
}

var xxx_messageInfo_MoneyTransaction proto.InternalMessageInfo

func (m *MoneyTransaction) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MoneyTransaction) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *MoneyTransaction) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *MoneyTransaction) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *MoneyTransaction) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *MoneyTransaction) GetMoneyTransactionShippingExternalId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingExternalId
	}
	return 0
}

func (m *MoneyTransaction) GetMoneyTransactionShippingEtopId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingEtopId
	}
	return 0
}

func (m *MoneyTransaction) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *MoneyTransaction) GetInvoiceNumber() string {
	if m != nil {
		return m.InvoiceNumber
	}
	return ""
}

func (m *MoneyTransaction) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

type MoneyTransactionsResponse struct {
	Paging               *common.PageInfo    `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	MoneyTransactions    []*MoneyTransaction `protobuf:"bytes,2,rep,name=money_transactions,json=moneyTransactions" json:"money_transactions"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *MoneyTransactionsResponse) Reset()         { *m = MoneyTransactionsResponse{} }
func (m *MoneyTransactionsResponse) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionsResponse) ProtoMessage()    {}
func (*MoneyTransactionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{23}
}

var xxx_messageInfo_MoneyTransactionsResponse proto.InternalMessageInfo

func (m *MoneyTransactionsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *MoneyTransactionsResponse) GetMoneyTransactions() []*MoneyTransaction {
	if m != nil {
		return m.MoneyTransactions
	}
	return nil
}

type MoneyTransactionShippingExternalLine struct {
	Id                                 dot.ID        `protobuf:"varint,1,opt,name=id" json:"id"`
	ExternalCode                       string        `protobuf:"bytes,2,opt,name=external_code,json=externalCode" json:"external_code"`
	ExternalCustomer                   string        `protobuf:"bytes,3,opt,name=external_customer,json=externalCustomer" json:"external_customer"`
	ExternalAddress                    string        `protobuf:"bytes,4,opt,name=external_address,json=externalAddress" json:"external_address"`
	ExternalTotalCod                   int           `protobuf:"varint,5,opt,name=external_total_cod,json=externalTotalCod" json:"external_total_cod"`
	ExternalTotalShippingFee           int           `protobuf:"varint,15,opt,name=external_total_shipping_fee,json=externalTotalShippingFee" json:"external_total_shipping_fee"`
	EtopFulfillmentId                  dot.ID        `protobuf:"varint,6,opt,name=etop_fulfillment_id,json=etopFulfillmentId" json:"etop_fulfillment_id"`
	EtopFulfillmentIdRaw               string        `protobuf:"bytes,7,opt,name=etop_fulfillment_id_raw,json=etopFulfillmentIdRaw" json:"etop_fulfillment_id_raw"`
	Note                               string        `protobuf:"bytes,8,opt,name=note" json:"note"`
	MoneyTransactionShippingExternalId dot.ID        `protobuf:"varint,9,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	ImportError                        *common.Error `protobuf:"bytes,10,opt,name=import_error,json=importError" json:"import_error"`
	CreatedAt                          dot.Time      `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                          dot.Time      `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ExternalCreatedAt                  dot.Time      `protobuf:"bytes,13,opt,name=external_created_at,json=externalCreatedAt" json:"external_created_at"`
	ExternalClosedAt                   dot.Time      `protobuf:"bytes,14,opt,name=external_closed_at,json=externalClosedAt" json:"external_closed_at"`
	Fulfillment                        *Fulfillment  `protobuf:"bytes,16,opt,name=fulfillment" json:"fulfillment"`
	XXX_NoUnkeyedLiteral               struct{}      `json:"-"`
	XXX_sizecache                      int32         `json:"-"`
}

func (m *MoneyTransactionShippingExternalLine) Reset()         { *m = MoneyTransactionShippingExternalLine{} }
func (m *MoneyTransactionShippingExternalLine) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionShippingExternalLine) ProtoMessage()    {}
func (*MoneyTransactionShippingExternalLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{24}
}

var xxx_messageInfo_MoneyTransactionShippingExternalLine proto.InternalMessageInfo

func (m *MoneyTransactionShippingExternalLine) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MoneyTransactionShippingExternalLine) GetExternalCode() string {
	if m != nil {
		return m.ExternalCode
	}
	return ""
}

func (m *MoneyTransactionShippingExternalLine) GetExternalCustomer() string {
	if m != nil {
		return m.ExternalCustomer
	}
	return ""
}

func (m *MoneyTransactionShippingExternalLine) GetExternalAddress() string {
	if m != nil {
		return m.ExternalAddress
	}
	return ""
}

func (m *MoneyTransactionShippingExternalLine) GetEtopFulfillmentId() dot.ID {
	if m != nil {
		return m.EtopFulfillmentId
	}
	return 0
}

func (m *MoneyTransactionShippingExternalLine) GetEtopFulfillmentIdRaw() string {
	if m != nil {
		return m.EtopFulfillmentIdRaw
	}
	return ""
}

func (m *MoneyTransactionShippingExternalLine) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *MoneyTransactionShippingExternalLine) GetMoneyTransactionShippingExternalId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingExternalId
	}
	return 0
}

func (m *MoneyTransactionShippingExternalLine) GetImportError() *common.Error {
	if m != nil {
		return m.ImportError
	}
	return nil
}

func (m *MoneyTransactionShippingExternalLine) GetFulfillment() *Fulfillment {
	if m != nil {
		return m.Fulfillment
	}
	return nil
}

type MoneyTransactionShippingExternal struct {
	Id                   dot.ID                                  `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string                                  `protobuf:"bytes,2,opt,name=code" json:"code"`
	TotalCod             int                                     `protobuf:"varint,3,opt,name=total_cod,json=totalCod" json:"total_cod"`
	TotalOrders          int                                     `protobuf:"varint,4,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	Status               status3.Status                          `protobuf:"varint,5,opt,name=status,enum=status3.Status" json:"status"`
	Provider             string                                  `protobuf:"bytes,6,opt,name=provider" json:"provider"`
	Lines                []*MoneyTransactionShippingExternalLine `protobuf:"bytes,7,rep,name=lines" json:"lines"`
	CreatedAt            dot.Time                                `protobuf:"bytes,8,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time                                `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ExternalPaidAt       dot.Time                                `protobuf:"bytes,10,opt,name=external_paid_at,json=externalPaidAt" json:"external_paid_at"`
	Note                 string                                  `protobuf:"bytes,11,opt,name=note" json:"note"`
	InvoiceNumber        string                                  `protobuf:"bytes,12,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount          *etop.BankAccount                       `protobuf:"bytes,13,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *MoneyTransactionShippingExternal) Reset()         { *m = MoneyTransactionShippingExternal{} }
func (m *MoneyTransactionShippingExternal) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionShippingExternal) ProtoMessage()    {}
func (*MoneyTransactionShippingExternal) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{25}
}

var xxx_messageInfo_MoneyTransactionShippingExternal proto.InternalMessageInfo

func (m *MoneyTransactionShippingExternal) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MoneyTransactionShippingExternal) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *MoneyTransactionShippingExternal) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *MoneyTransactionShippingExternal) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *MoneyTransactionShippingExternal) GetLines() []*MoneyTransactionShippingExternalLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *MoneyTransactionShippingExternal) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *MoneyTransactionShippingExternal) GetInvoiceNumber() string {
	if m != nil {
		return m.InvoiceNumber
	}
	return ""
}

func (m *MoneyTransactionShippingExternal) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

type MoneyTransactionShippingExternalsResponse struct {
	Paging               *common.PageInfo                    `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	MoneyTransactions    []*MoneyTransactionShippingExternal `protobuf:"bytes,2,rep,name=money_transactions,json=moneyTransactions" json:"money_transactions"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *MoneyTransactionShippingExternalsResponse) Reset() {
	*m = MoneyTransactionShippingExternalsResponse{}
}
func (m *MoneyTransactionShippingExternalsResponse) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionShippingExternalsResponse) ProtoMessage()    {}
func (*MoneyTransactionShippingExternalsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{26}
}

var xxx_messageInfo_MoneyTransactionShippingExternalsResponse proto.InternalMessageInfo

func (m *MoneyTransactionShippingExternalsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *MoneyTransactionShippingExternalsResponse) GetMoneyTransactions() []*MoneyTransactionShippingExternal {
	if m != nil {
		return m.MoneyTransactions
	}
	return nil
}

type MoneyTransactionShippingEtop struct {
	Id                   dot.ID              `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string              `protobuf:"bytes,2,opt,name=code" json:"code"`
	TotalCod             int                 `protobuf:"varint,3,opt,name=total_cod,json=totalCod" json:"total_cod"`
	TotalOrders          int                 `protobuf:"varint,4,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	TotalAmount          int                 `protobuf:"varint,5,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	TotalFee             int                 `protobuf:"varint,6,opt,name=total_fee,json=totalFee" json:"total_fee"`
	Status               status3.Status      `protobuf:"varint,7,opt,name=status,enum=status3.Status" json:"status"`
	MoneyTransactions    []*MoneyTransaction `protobuf:"bytes,8,rep,name=money_transactions,json=moneyTransactions" json:"money_transactions"`
	CreatedAt            dot.Time            `protobuf:"bytes,9,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time            `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ConfirmedAt          dot.Time            `protobuf:"bytes,11,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	Note                 string              `protobuf:"bytes,12,opt,name=note" json:"note"`
	InvoiceNumber        string              `protobuf:"bytes,13,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount          *etop.BankAccount   `protobuf:"bytes,14,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *MoneyTransactionShippingEtop) Reset()         { *m = MoneyTransactionShippingEtop{} }
func (m *MoneyTransactionShippingEtop) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionShippingEtop) ProtoMessage()    {}
func (*MoneyTransactionShippingEtop) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{27}
}

var xxx_messageInfo_MoneyTransactionShippingEtop proto.InternalMessageInfo

func (m *MoneyTransactionShippingEtop) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MoneyTransactionShippingEtop) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *MoneyTransactionShippingEtop) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *MoneyTransactionShippingEtop) GetMoneyTransactions() []*MoneyTransaction {
	if m != nil {
		return m.MoneyTransactions
	}
	return nil
}

func (m *MoneyTransactionShippingEtop) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *MoneyTransactionShippingEtop) GetInvoiceNumber() string {
	if m != nil {
		return m.InvoiceNumber
	}
	return ""
}

func (m *MoneyTransactionShippingEtop) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

type MoneyTransactionShippingEtopsResponse struct {
	Paging                        *common.PageInfo                `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	MoneyTransactionShippingEtops []*MoneyTransactionShippingEtop `protobuf:"bytes,2,rep,name=money_transaction_shipping_etops,json=moneyTransactionShippingEtops" json:"money_transaction_shipping_etops"`
	XXX_NoUnkeyedLiteral          struct{}                        `json:"-"`
	XXX_sizecache                 int32                           `json:"-"`
}

func (m *MoneyTransactionShippingEtopsResponse) Reset()         { *m = MoneyTransactionShippingEtopsResponse{} }
func (m *MoneyTransactionShippingEtopsResponse) String() string { return proto.CompactTextString(m) }
func (*MoneyTransactionShippingEtopsResponse) ProtoMessage()    {}
func (*MoneyTransactionShippingEtopsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{28}
}

var xxx_messageInfo_MoneyTransactionShippingEtopsResponse proto.InternalMessageInfo

func (m *MoneyTransactionShippingEtopsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *MoneyTransactionShippingEtopsResponse) GetMoneyTransactionShippingEtops() []*MoneyTransactionShippingEtop {
	if m != nil {
		return m.MoneyTransactionShippingEtops
	}
	return nil
}

type ImportOrdersResponse struct {
	Data                 *spreadsheet.SpreadsheetData `protobuf:"bytes,1,opt,name=data" json:"data"`
	Orders               []*Order                     `protobuf:"bytes,2,rep,name=orders" json:"orders"`
	ImportErrors         []*common.Error              `protobuf:"bytes,3,rep,name=import_errors,json=importErrors" json:"import_errors"`
	CellErrors           []*common.Error              `protobuf:"bytes,4,rep,name=cell_errors,json=cellErrors" json:"cell_errors"`
	ImportId             dot.ID                       `protobuf:"varint,5,opt,name=import_id,json=importId" json:"import_id"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *ImportOrdersResponse) Reset()         { *m = ImportOrdersResponse{} }
func (m *ImportOrdersResponse) String() string { return proto.CompactTextString(m) }
func (*ImportOrdersResponse) ProtoMessage()    {}
func (*ImportOrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{29}
}

var xxx_messageInfo_ImportOrdersResponse proto.InternalMessageInfo

func (m *ImportOrdersResponse) GetData() *spreadsheet.SpreadsheetData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ImportOrdersResponse) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

func (m *ImportOrdersResponse) GetImportErrors() []*common.Error {
	if m != nil {
		return m.ImportErrors
	}
	return nil
}

func (m *ImportOrdersResponse) GetCellErrors() []*common.Error {
	if m != nil {
		return m.CellErrors
	}
	return nil
}

func (m *ImportOrdersResponse) GetImportId() dot.ID {
	if m != nil {
		return m.ImportId
	}
	return 0
}

// Public API for using with ManyChat
type PublicFulfillment struct {
	Id            dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	ShippingState shipping.State `protobuf:"varint,2,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	Status        status5.Status `protobuf:"varint,3,opt,name=status,enum=status5.Status" json:"status"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt  dot.Time `protobuf:"bytes,4,opt,name=expected_delivery_at,json=expectedDeliveryAt" json:"expected_delivery_at"`
	DeliveredAt         dot.Time `protobuf:"bytes,5,opt,name=delivered_at,json=deliveredAt" json:"delivered_at"`
	EstimatedPickupAt   dot.Time `protobuf:"bytes,15,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time `protobuf:"bytes,16,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	ShippingCode        string   `protobuf:"bytes,6,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	OrderId             dot.ID   `protobuf:"varint,8,opt,name=order_id,json=orderId" json:"order_id"`
	// For using with ManyChat
	DeliveredAtText      string   `protobuf:"bytes,9,opt,name=delivered_at_text,json=deliveredAtText" json:"delivered_at_text"`
	ShippingStateText    string   `protobuf:"bytes,10,opt,name=shipping_state_text,json=shippingStateText" json:"shipping_state_text"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicFulfillment) Reset()         { *m = PublicFulfillment{} }
func (m *PublicFulfillment) String() string { return proto.CompactTextString(m) }
func (*PublicFulfillment) ProtoMessage()    {}
func (*PublicFulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{30}
}

var xxx_messageInfo_PublicFulfillment proto.InternalMessageInfo

func (m *PublicFulfillment) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PublicFulfillment) GetShippingState() shipping.State {
	if m != nil {
		return m.ShippingState
	}
	return shipping.State_default
}

func (m *PublicFulfillment) GetStatus() status5.Status {
	if m != nil {
		return m.Status
	}
	return status5.Status_Z
}

func (m *PublicFulfillment) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

func (m *PublicFulfillment) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *PublicFulfillment) GetDeliveredAtText() string {
	if m != nil {
		return m.DeliveredAtText
	}
	return ""
}

func (m *PublicFulfillment) GetShippingStateText() string {
	if m != nil {
		return m.ShippingStateText
	}
	return ""
}

type ShipnowFulfillment struct {
	Id                         dot.ID           `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId                     dot.ID           `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	PartnerId                  dot.ID           `protobuf:"varint,3,opt,name=partner_id,json=partnerId" json:"partner_id"`
	PickupAddress              *OrderAddress    `protobuf:"bytes,4,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	DeliveryPoints             []*DeliveryPoint `protobuf:"bytes,5,rep,name=delivery_points,json=deliveryPoints" json:"delivery_points"`
	Carrier                    string           `protobuf:"bytes,6,opt,name=carrier" json:"carrier"`
	ShippingServiceCode        string           `protobuf:"bytes,7,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee         int32            `protobuf:"varint,8,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingServiceName        string           `protobuf:"bytes,28,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceDescription string           `protobuf:"bytes,31,opt,name=shipping_service_description,json=shippingServiceDescription" json:"shipping_service_description"`
	WeightInfo                 `protobuf:"bytes,9,opt,name=weight_info,json=weightInfo,embedded=weight_info" json:"weight_info"`
	ValueInfo                  ValueInfo      `protobuf:"bytes,10,opt,name=value_info,json=valueInfo" json:"value_info"`
	ShippingNote               string         `protobuf:"bytes,11,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt            dot.Time       `protobuf:"bytes,12,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at"`
	CreatedAt                  dot.Time       `protobuf:"bytes,13,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                  dot.Time       `protobuf:"bytes,14,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	Status                     status5.Status `protobuf:"varint,15,opt,name=status,enum=status5.Status" json:"status"`
	ShippingStatus             status5.Status `protobuf:"varint,16,opt,name=shipping_status,json=shippingStatus,enum=status5.Status" json:"shipping_status"`
	ShippingState              string         `protobuf:"bytes,17,opt,name=shipping_state,json=shippingState" json:"shipping_state"`
	ConfirmStatus              status3.Status `protobuf:"varint,18,opt,name=confirm_status,json=confirmStatus,enum=status3.Status" json:"confirm_status"`
	OrderIds                   []dot.ID       `protobuf:"varint,19,rep,name=order_ids,json=orderIds" json:"order_ids"`
	ShippingCreatedAt          dot.Time       `protobuf:"bytes,20,opt,name=shipping_created_at,json=shippingCreatedAt" json:"shipping_created_at"`
	ShippingCode               string         `protobuf:"bytes,21,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	EtopPaymentStatus          status4.Status `protobuf:"varint,22,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	CodEtopTransferedAt        dot.Time       `protobuf:"bytes,23,opt,name=cod_etop_transfered_at,json=codEtopTransferedAt" json:"cod_etop_transfered_at"`
	ShippingPickingAt          dot.Time       `protobuf:"bytes,24,opt,name=shipping_picking_at,json=shippingPickingAt" json:"shipping_picking_at"`
	ShippingDeliveringAt       dot.Time       `protobuf:"bytes,25,opt,name=shipping_delivering_at,json=shippingDeliveringAt" json:"shipping_delivering_at"`
	ShippingDeliveredAt        dot.Time       `protobuf:"bytes,26,opt,name=shipping_delivered_at,json=shippingDeliveredAt" json:"shipping_delivered_at"`
	ShippingCancelledAt        dot.Time       `protobuf:"bytes,27,opt,name=shipping_cancelled_at,json=shippingCancelledAt" json:"shipping_cancelled_at"`
	ShippingSharedLink         string         `protobuf:"bytes,29,opt,name=shipping_shared_link,json=shippingSharedLink" json:"shipping_shared_link"`
	CancelReason               string         `protobuf:"bytes,30,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	XXX_NoUnkeyedLiteral       struct{}       `json:"-"`
	XXX_sizecache              int32          `json:"-"`
}

func (m *ShipnowFulfillment) Reset()         { *m = ShipnowFulfillment{} }
func (m *ShipnowFulfillment) String() string { return proto.CompactTextString(m) }
func (*ShipnowFulfillment) ProtoMessage()    {}
func (*ShipnowFulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{31}
}

var xxx_messageInfo_ShipnowFulfillment proto.InternalMessageInfo

func (m *ShipnowFulfillment) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShipnowFulfillment) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *ShipnowFulfillment) GetPartnerId() dot.ID {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *ShipnowFulfillment) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *ShipnowFulfillment) GetDeliveryPoints() []*DeliveryPoint {
	if m != nil {
		return m.DeliveryPoints
	}
	return nil
}

func (m *ShipnowFulfillment) GetCarrier() string {
	if m != nil {
		return m.Carrier
	}
	return ""
}

func (m *ShipnowFulfillment) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *ShipnowFulfillment) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *ShipnowFulfillment) GetShippingServiceName() string {
	if m != nil {
		return m.ShippingServiceName
	}
	return ""
}

func (m *ShipnowFulfillment) GetShippingServiceDescription() string {
	if m != nil {
		return m.ShippingServiceDescription
	}
	return ""
}

func (m *ShipnowFulfillment) GetValueInfo() ValueInfo {
	if m != nil {
		return m.ValueInfo
	}
	return ValueInfo{}
}

func (m *ShipnowFulfillment) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *ShipnowFulfillment) GetStatus() status5.Status {
	if m != nil {
		return m.Status
	}
	return status5.Status_Z
}

func (m *ShipnowFulfillment) GetShippingStatus() status5.Status {
	if m != nil {
		return m.ShippingStatus
	}
	return status5.Status_Z
}

func (m *ShipnowFulfillment) GetShippingState() string {
	if m != nil {
		return m.ShippingState
	}
	return ""
}

func (m *ShipnowFulfillment) GetConfirmStatus() status3.Status {
	if m != nil {
		return m.ConfirmStatus
	}
	return status3.Status_Z
}

func (m *ShipnowFulfillment) GetOrderIds() []dot.ID {
	if m != nil {
		return m.OrderIds
	}
	return nil
}

func (m *ShipnowFulfillment) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

func (m *ShipnowFulfillment) GetEtopPaymentStatus() status4.Status {
	if m != nil {
		return m.EtopPaymentStatus
	}
	return status4.Status_Z
}

func (m *ShipnowFulfillment) GetShippingSharedLink() string {
	if m != nil {
		return m.ShippingSharedLink
	}
	return ""
}

func (m *ShipnowFulfillment) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

type GetShipnowFulfillmentsRequest struct {
	Paging               *common.Paging     `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter   `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	Mixed                *etop.MixedAccount `protobuf:"bytes,3,opt,name=mixed" json:"mixed"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetShipnowFulfillmentsRequest) Reset()         { *m = GetShipnowFulfillmentsRequest{} }
func (m *GetShipnowFulfillmentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetShipnowFulfillmentsRequest) ProtoMessage()    {}
func (*GetShipnowFulfillmentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{32}
}

var xxx_messageInfo_GetShipnowFulfillmentsRequest proto.InternalMessageInfo

func (m *GetShipnowFulfillmentsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetShipnowFulfillmentsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *GetShipnowFulfillmentsRequest) GetMixed() *etop.MixedAccount {
	if m != nil {
		return m.Mixed
	}
	return nil
}

type ShipnowFulfillments struct {
	Paging               *common.PageInfo      `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	ShipnowFulfillments  []*ShipnowFulfillment `protobuf:"bytes,2,rep,name=shipnow_fulfillments,json=shipnowFulfillments" json:"shipnow_fulfillments"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ShipnowFulfillments) Reset()         { *m = ShipnowFulfillments{} }
func (m *ShipnowFulfillments) String() string { return proto.CompactTextString(m) }
func (*ShipnowFulfillments) ProtoMessage()    {}
func (*ShipnowFulfillments) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{33}
}

var xxx_messageInfo_ShipnowFulfillments proto.InternalMessageInfo

func (m *ShipnowFulfillments) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShipnowFulfillments) GetShipnowFulfillments() []*ShipnowFulfillment {
	if m != nil {
		return m.ShipnowFulfillments
	}
	return nil
}

type DeliveryPoint struct {
	ShippingAddress      *OrderAddress `protobuf:"bytes,1,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	Lines                []*OrderLine  `protobuf:"bytes,2,rep,name=lines" json:"lines"`
	ShippingNote         string        `protobuf:"bytes,3,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	OrderId              dot.ID        `protobuf:"varint,4,opt,name=order_id,json=orderId" json:"order_id"`
	WeightInfo           `protobuf:"bytes,5,opt,name=weight_info,json=weightInfo,embedded=weight_info" json:"weight_info"`
	ValueInfo            `protobuf:"bytes,6,opt,name=value_info,json=valueInfo,embedded=value_info" json:"value_info"`
	TryOn                try_on.TryOnCode `protobuf:"varint,7,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *DeliveryPoint) Reset()         { *m = DeliveryPoint{} }
func (m *DeliveryPoint) String() string { return proto.CompactTextString(m) }
func (*DeliveryPoint) ProtoMessage()    {}
func (*DeliveryPoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{34}
}

var xxx_messageInfo_DeliveryPoint proto.InternalMessageInfo

func (m *DeliveryPoint) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *DeliveryPoint) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *DeliveryPoint) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *DeliveryPoint) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *DeliveryPoint) GetTryOn() try_on.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return try_on.TryOnCode_unknown
}

type WeightInfo struct {
	GrossWeight          int32    `protobuf:"varint,1,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	ChargeableWeight     int32    `protobuf:"varint,2,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	Length               int32    `protobuf:"varint,3,opt,name=length" json:"length"`
	Width                int32    `protobuf:"varint,4,opt,name=width" json:"width"`
	Height               int32    `protobuf:"varint,5,opt,name=height" json:"height"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeightInfo) Reset()         { *m = WeightInfo{} }
func (m *WeightInfo) String() string { return proto.CompactTextString(m) }
func (*WeightInfo) ProtoMessage()    {}
func (*WeightInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{35}
}

var xxx_messageInfo_WeightInfo proto.InternalMessageInfo

func (m *WeightInfo) GetGrossWeight() int32 {
	if m != nil {
		return m.GrossWeight
	}
	return 0
}

func (m *WeightInfo) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *WeightInfo) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *WeightInfo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *WeightInfo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type ValueInfo struct {
	BasketValue          int32    `protobuf:"varint,1,opt,name=basket_value,json=basketValue" json:"basket_value"`
	CodAmount            int32    `protobuf:"varint,2,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	IncludeInsurance     bool     `protobuf:"varint,3,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValueInfo) Reset()         { *m = ValueInfo{} }
func (m *ValueInfo) String() string { return proto.CompactTextString(m) }
func (*ValueInfo) ProtoMessage()    {}
func (*ValueInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{36}
}

var xxx_messageInfo_ValueInfo proto.InternalMessageInfo

func (m *ValueInfo) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *ValueInfo) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *ValueInfo) GetIncludeInsurance() bool {
	if m != nil {
		return m.IncludeInsurance
	}
	return false
}

type CreateShipnowFulfillmentRequest struct {
	OrderIds             []dot.ID      `protobuf:"varint,1,rep,name=order_ids,json=orderIds" json:"order_ids"`
	Carrier              string        `protobuf:"bytes,2,opt,name=carrier" json:"carrier"`
	ShippingServiceCode  string        `protobuf:"bytes,3,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee   int32         `protobuf:"varint,4,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingNote         string        `protobuf:"bytes,5,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt      dot.Time      `protobuf:"bytes,6,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at"`
	PickupAddress        *OrderAddress `protobuf:"bytes,7,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CreateShipnowFulfillmentRequest) Reset()         { *m = CreateShipnowFulfillmentRequest{} }
func (m *CreateShipnowFulfillmentRequest) String() string { return proto.CompactTextString(m) }
func (*CreateShipnowFulfillmentRequest) ProtoMessage()    {}
func (*CreateShipnowFulfillmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{37}
}

var xxx_messageInfo_CreateShipnowFulfillmentRequest proto.InternalMessageInfo

func (m *CreateShipnowFulfillmentRequest) GetOrderIds() []dot.ID {
	if m != nil {
		return m.OrderIds
	}
	return nil
}

func (m *CreateShipnowFulfillmentRequest) GetCarrier() string {
	if m != nil {
		return m.Carrier
	}
	return ""
}

func (m *CreateShipnowFulfillmentRequest) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *CreateShipnowFulfillmentRequest) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *CreateShipnowFulfillmentRequest) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *CreateShipnowFulfillmentRequest) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

type UpdateShipnowFulfillmentRequest struct {
	Id                   dot.ID        `protobuf:"varint,8,opt,name=id" json:"id"`
	OrderIds             []dot.ID      `protobuf:"varint,1,rep,name=order_ids,json=orderIds" json:"order_ids"`
	Carrier              string        `protobuf:"bytes,2,opt,name=carrier" json:"carrier"`
	ShippingServiceCode  string        `protobuf:"bytes,3,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee   int32         `protobuf:"varint,4,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingNote         string        `protobuf:"bytes,5,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt      dot.Time      `protobuf:"bytes,6,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at"`
	PickupAddress        *OrderAddress `protobuf:"bytes,7,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *UpdateShipnowFulfillmentRequest) Reset()         { *m = UpdateShipnowFulfillmentRequest{} }
func (m *UpdateShipnowFulfillmentRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateShipnowFulfillmentRequest) ProtoMessage()    {}
func (*UpdateShipnowFulfillmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{38}
}

var xxx_messageInfo_UpdateShipnowFulfillmentRequest proto.InternalMessageInfo

func (m *UpdateShipnowFulfillmentRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateShipnowFulfillmentRequest) GetOrderIds() []dot.ID {
	if m != nil {
		return m.OrderIds
	}
	return nil
}

func (m *UpdateShipnowFulfillmentRequest) GetCarrier() string {
	if m != nil {
		return m.Carrier
	}
	return ""
}

func (m *UpdateShipnowFulfillmentRequest) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *UpdateShipnowFulfillmentRequest) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *UpdateShipnowFulfillmentRequest) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *UpdateShipnowFulfillmentRequest) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

type CancelShipnowFulfillmentRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	CancelReason         string   `protobuf:"bytes,2,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelShipnowFulfillmentRequest) Reset()         { *m = CancelShipnowFulfillmentRequest{} }
func (m *CancelShipnowFulfillmentRequest) String() string { return proto.CompactTextString(m) }
func (*CancelShipnowFulfillmentRequest) ProtoMessage()    {}
func (*CancelShipnowFulfillmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{39}
}

var xxx_messageInfo_CancelShipnowFulfillmentRequest proto.InternalMessageInfo

func (m *CancelShipnowFulfillmentRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelShipnowFulfillmentRequest) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

type GetShipnowServicesRequest struct {
	OrderIds             []dot.ID                `protobuf:"varint,1,rep,name=order_ids,json=orderIds" json:"order_ids"`
	PickupAddress        *OrderAddress           `protobuf:"bytes,2,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	DeliveryPoints       []*DeliveryPointRequest `protobuf:"bytes,3,rep,name=delivery_points,json=deliveryPoints" json:"delivery_points"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *GetShipnowServicesRequest) Reset()         { *m = GetShipnowServicesRequest{} }
func (m *GetShipnowServicesRequest) String() string { return proto.CompactTextString(m) }
func (*GetShipnowServicesRequest) ProtoMessage()    {}
func (*GetShipnowServicesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{40}
}

var xxx_messageInfo_GetShipnowServicesRequest proto.InternalMessageInfo

func (m *GetShipnowServicesRequest) GetOrderIds() []dot.ID {
	if m != nil {
		return m.OrderIds
	}
	return nil
}

func (m *GetShipnowServicesRequest) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *GetShipnowServicesRequest) GetDeliveryPoints() []*DeliveryPointRequest {
	if m != nil {
		return m.DeliveryPoints
	}
	return nil
}

type DeliveryPointRequest struct {
	ShippingAddress      *OrderAddress `protobuf:"bytes,1,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	CodAmount            int32         `protobuf:"varint,2,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *DeliveryPointRequest) Reset()         { *m = DeliveryPointRequest{} }
func (m *DeliveryPointRequest) String() string { return proto.CompactTextString(m) }
func (*DeliveryPointRequest) ProtoMessage()    {}
func (*DeliveryPointRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{41}
}

var xxx_messageInfo_DeliveryPointRequest proto.InternalMessageInfo

func (m *DeliveryPointRequest) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *DeliveryPointRequest) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

type GetShipnowServicesResponse struct {
	Services             []*ShippnowService `protobuf:"bytes,1,rep,name=services" json:"services"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetShipnowServicesResponse) Reset()         { *m = GetShipnowServicesResponse{} }
func (m *GetShipnowServicesResponse) String() string { return proto.CompactTextString(m) }
func (*GetShipnowServicesResponse) ProtoMessage()    {}
func (*GetShipnowServicesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{42}
}

var xxx_messageInfo_GetShipnowServicesResponse proto.InternalMessageInfo

func (m *GetShipnowServicesResponse) GetServices() []*ShippnowService {
	if m != nil {
		return m.Services
	}
	return nil
}

type ShippnowService struct {
	Carrier              string   `protobuf:"bytes,1,opt,name=carrier" json:"carrier"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Code                 string   `protobuf:"bytes,3,opt,name=code" json:"code"`
	Fee                  int32    `protobuf:"varint,4,opt,name=fee" json:"fee"`
	Description          string   `protobuf:"bytes,5,opt,name=description" json:"description"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShippnowService) Reset()         { *m = ShippnowService{} }
func (m *ShippnowService) String() string { return proto.CompactTextString(m) }
func (*ShippnowService) ProtoMessage()    {}
func (*ShippnowService) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{43}
}

var xxx_messageInfo_ShippnowService proto.InternalMessageInfo

func (m *ShippnowService) GetCarrier() string {
	if m != nil {
		return m.Carrier
	}
	return ""
}

func (m *ShippnowService) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShippnowService) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ShippnowService) GetFee() int32 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func (m *ShippnowService) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type XFulfillment struct {
	Shipnow  *ShipnowFulfillment `protobuf:"bytes,77,opt,name=shipnow" json:"shipnow"`
	Shipment *Fulfillment        `protobuf:"bytes,78,opt,name=shipment" json:"shipment"`
	// backward-compatible fields from shipment
	Id         dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	OrderId    dot.ID       `protobuf:"varint,2,opt,name=order_id,json=orderId" json:"order_id"`
	ShopId     dot.ID       `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	PartnerId  dot.ID       `protobuf:"varint,55,opt,name=partner_id,json=partnerId" json:"partner_id"`
	SelfUrl    string       `protobuf:"bytes,59,opt,name=self_url,json=selfUrl" json:"self_url"`
	Lines      []*OrderLine `protobuf:"bytes,5,rep,name=lines" json:"lines"`
	TotalItems int32        `protobuf:"varint,6,opt,name=total_items,json=totalItems" json:"total_items"`
	// @deprecated use chargeable_weight
	TotalWeight int32 `protobuf:"varint,7,opt,name=total_weight,json=totalWeight" json:"total_weight"`
	BasketValue int32 `protobuf:"varint,10,opt,name=basket_value,json=basketValue" json:"basket_value"`
	// @deprecated use cod_amount
	TotalCodAmount int32 `protobuf:"varint,9,opt,name=total_cod_amount,json=totalCodAmount" json:"total_cod_amount"`
	CodAmount      int32 `protobuf:"varint,60,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	// @deprecated
	TotalAmount      int32    `protobuf:"varint,8,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	ChargeableWeight int32    `protobuf:"varint,72,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	CreatedAt        dot.Time `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt        dot.Time `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ClosedAt         dot.Time `protobuf:"bytes,14,opt,name=closed_at,json=closedAt" json:"closed_at"`
	CancelledAt      dot.Time `protobuf:"bytes,13,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CancelReason     string   `protobuf:"bytes,15,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	// @deprecated use carrier instead
	ShippingProvider     string                             `protobuf:"bytes,17,opt,name=shipping_provider,json=shippingProvider" json:"shipping_provider"`
	Carrier              shipping_provider.ShippingProvider `protobuf:"varint,71,opt,name=carrier,enum=shipping_provider.ShippingProvider" json:"carrier"`
	ShippingServiceName  string                             `protobuf:"bytes,70,opt,name=shipping_service_name,json=shippingServiceName" json:"shipping_service_name"`
	ShippingServiceFee   int32                              `protobuf:"varint,68,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingServiceCode  string                             `protobuf:"bytes,69,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingCode         string                             `protobuf:"bytes,18,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	ShippingNote         string                             `protobuf:"bytes,19,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	TryOn                try_on.TryOnCode                   `protobuf:"varint,20,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	IncludeInsurance     bool                               `protobuf:"varint,63,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
	ShConfirm            status3.Status                     `protobuf:"varint,23,opt,name=sh_confirm,json=shConfirm,enum=status3.Status" json:"sh_confirm"`
	ShippingState        shipping.State                     `protobuf:"varint,25,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	Status               status5.Status                     `protobuf:"varint,26,opt,name=status,enum=status5.Status" json:"status"`
	ShippingStatus       status5.Status                     `protobuf:"varint,64,opt,name=shipping_status,json=shippingStatus,enum=status5.Status" json:"shipping_status"`
	EtopPaymentStatus    status4.Status                     `protobuf:"varint,22,opt,name=etop_payment_status,json=etopPaymentStatus,enum=status4.Status" json:"etop_payment_status"`
	ShippingFeeCustomer  int32                              `protobuf:"varint,27,opt,name=shipping_fee_customer,json=shippingFeeCustomer" json:"shipping_fee_customer"`
	ShippingFeeShop      int32                              `protobuf:"varint,28,opt,name=shipping_fee_shop,json=shippingFeeShop" json:"shipping_fee_shop"`
	XShippingFee         int32                              `protobuf:"varint,29,opt,name=x_shipping_fee,json=xShippingFee" json:"x_shipping_fee"`
	XShippingId          string                             `protobuf:"bytes,30,opt,name=x_shipping_id,json=xShippingId" json:"x_shipping_id"`
	XShippingCode        string                             `protobuf:"bytes,31,opt,name=x_shipping_code,json=xShippingCode" json:"x_shipping_code"`
	XShippingCreatedAt   dot.Time                           `protobuf:"bytes,32,opt,name=x_shipping_created_at,json=xShippingCreatedAt" json:"x_shipping_created_at"`
	XShippingUpdatedAt   dot.Time                           `protobuf:"bytes,33,opt,name=x_shipping_updated_at,json=xShippingUpdatedAt" json:"x_shipping_updated_at"`
	XShippingCancelledAt dot.Time                           `protobuf:"bytes,34,opt,name=x_shipping_cancelled_at,json=xShippingCancelledAt" json:"x_shipping_cancelled_at"`
	XShippingDeliveredAt dot.Time                           `protobuf:"bytes,35,opt,name=x_shipping_delivered_at,json=xShippingDeliveredAt" json:"x_shipping_delivered_at"`
	XShippingReturnedAt  dot.Time                           `protobuf:"bytes,36,opt,name=x_shipping_returned_at,json=xShippingReturnedAt" json:"x_shipping_returned_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt dot.Time `protobuf:"bytes,41,opt,name=expected_delivery_at,json=expectedDeliveryAt" json:"expected_delivery_at"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt              dot.Time               `protobuf:"bytes,51,opt,name=expected_pick_at,json=expectedPickAt" json:"expected_pick_at"`
	EstimatedDeliveryAt         dot.Time               `protobuf:"bytes,74,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at"`
	EstimatedPickupAt           dot.Time               `protobuf:"bytes,75,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at"`
	CodEtopTransferedAt         dot.Time               `protobuf:"bytes,43,opt,name=cod_etop_transfered_at,json=codEtopTransferedAt" json:"cod_etop_transfered_at"`
	ShippingFeeShopTransferedAt dot.Time               `protobuf:"bytes,44,opt,name=shipping_fee_shop_transfered_at,json=shippingFeeShopTransferedAt" json:"shipping_fee_shop_transfered_at"`
	XShippingState              string                 `protobuf:"bytes,37,opt,name=x_shipping_state,json=xShippingState" json:"x_shipping_state"`
	XShippingStatus             status5.Status         `protobuf:"varint,38,opt,name=x_shipping_status,json=xShippingStatus,enum=status5.Status" json:"x_shipping_status"`
	XSyncStatus                 status4.Status         `protobuf:"varint,39,opt,name=x_sync_status,json=xSyncStatus,enum=status4.Status" json:"x_sync_status"`
	XSyncStates                 *FulfillmentSyncStates `protobuf:"bytes,40,opt,name=x_sync_states,json=xSyncStates" json:"x_sync_states"`
	// @deprecated use shipping_address instead
	AddressTo *etop.Address `protobuf:"bytes,42,opt,name=address_to,json=addressTo" json:"address_to"`
	// @deprecated use pickup_address instead
	AddressFrom                        *etop.Address          `protobuf:"bytes,50,opt,name=address_from,json=addressFrom" json:"address_from"`
	PickupAddress                      *OrderAddress          `protobuf:"bytes,65,opt,name=pickup_address,json=pickupAddress" json:"pickup_address"`
	ReturnAddress                      *OrderAddress          `protobuf:"bytes,66,opt,name=return_address,json=returnAddress" json:"return_address"`
	ShippingAddress                    *OrderAddress          `protobuf:"bytes,67,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	Shop                               *etop.Shop             `protobuf:"bytes,45,opt,name=shop" json:"shop"`
	Order                              *Order                 `protobuf:"bytes,46,opt,name=order" json:"order"`
	ProviderShippingFeeLines           []*ShippingFeeLine     `protobuf:"bytes,47,rep,name=provider_shipping_fee_lines,json=providerShippingFeeLines" json:"provider_shipping_fee_lines"`
	ShippingFeeShopLines               []*ShippingFeeLine     `protobuf:"bytes,48,rep,name=shipping_fee_shop_lines,json=shippingFeeShopLines" json:"shipping_fee_shop_lines"`
	EtopDiscount                       int32                  `protobuf:"varint,49,opt,name=etop_discount,json=etopDiscount" json:"etop_discount"`
	MoneyTransactionShippingId         dot.ID                 `protobuf:"varint,52,opt,name=money_transaction_shipping_id,json=moneyTransactionShippingId" json:"money_transaction_shipping_id"`
	MoneyTransactionShippingExternalId dot.ID                 `protobuf:"varint,53,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	XShippingLogs                      []*ExternalShippingLog `protobuf:"bytes,54,rep,name=x_shipping_logs,json=xShippingLogs" json:"x_shipping_logs"`
	XShippingNote                      string                 `protobuf:"bytes,56,opt,name=x_shipping_note,json=xShippingNote" json:"x_shipping_note"`
	XShippingSubState                  string                 `protobuf:"bytes,57,opt,name=x_shipping_sub_state,json=xShippingSubState" json:"x_shipping_sub_state"`
	Code                               string                 `protobuf:"bytes,58,opt,name=code" json:"code"`
	ActualCompensationAmount           int32                  `protobuf:"varint,76,opt,name=actual_compensation_amount,json=actualCompensationAmount" json:"actual_compensation_amount"`
	XXX_NoUnkeyedLiteral               struct{}               `json:"-"`
	XXX_sizecache                      int32                  `json:"-"`
}

func (m *XFulfillment) Reset()         { *m = XFulfillment{} }
func (m *XFulfillment) String() string { return proto.CompactTextString(m) }
func (*XFulfillment) ProtoMessage()    {}
func (*XFulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{44}
}

var xxx_messageInfo_XFulfillment proto.InternalMessageInfo

func (m *XFulfillment) GetShipnow() *ShipnowFulfillment {
	if m != nil {
		return m.Shipnow
	}
	return nil
}

func (m *XFulfillment) GetShipment() *Fulfillment {
	if m != nil {
		return m.Shipment
	}
	return nil
}

func (m *XFulfillment) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *XFulfillment) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *XFulfillment) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *XFulfillment) GetPartnerId() dot.ID {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *XFulfillment) GetSelfUrl() string {
	if m != nil {
		return m.SelfUrl
	}
	return ""
}

func (m *XFulfillment) GetLines() []*OrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *XFulfillment) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *XFulfillment) GetTotalWeight() int32 {
	if m != nil {
		return m.TotalWeight
	}
	return 0
}

func (m *XFulfillment) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *XFulfillment) GetTotalCodAmount() int32 {
	if m != nil {
		return m.TotalCodAmount
	}
	return 0
}

func (m *XFulfillment) GetCodAmount() int32 {
	if m != nil {
		return m.CodAmount
	}
	return 0
}

func (m *XFulfillment) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *XFulfillment) GetChargeableWeight() int32 {
	if m != nil {
		return m.ChargeableWeight
	}
	return 0
}

func (m *XFulfillment) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *XFulfillment) GetShippingProvider() string {
	if m != nil {
		return m.ShippingProvider
	}
	return ""
}

func (m *XFulfillment) GetCarrier() shipping_provider.ShippingProvider {
	if m != nil {
		return m.Carrier
	}
	return shipping_provider.ShippingProvider_unknown
}

func (m *XFulfillment) GetShippingServiceName() string {
	if m != nil {
		return m.ShippingServiceName
	}
	return ""
}

func (m *XFulfillment) GetShippingServiceFee() int32 {
	if m != nil {
		return m.ShippingServiceFee
	}
	return 0
}

func (m *XFulfillment) GetShippingServiceCode() string {
	if m != nil {
		return m.ShippingServiceCode
	}
	return ""
}

func (m *XFulfillment) GetShippingCode() string {
	if m != nil {
		return m.ShippingCode
	}
	return ""
}

func (m *XFulfillment) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *XFulfillment) GetTryOn() try_on.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *XFulfillment) GetIncludeInsurance() bool {
	if m != nil {
		return m.IncludeInsurance
	}
	return false
}

func (m *XFulfillment) GetShConfirm() status3.Status {
	if m != nil {
		return m.ShConfirm
	}
	return status3.Status_Z
}

func (m *XFulfillment) GetShippingState() shipping.State {
	if m != nil {
		return m.ShippingState
	}
	return shipping.State_default
}

func (m *XFulfillment) GetStatus() status5.Status {
	if m != nil {
		return m.Status
	}
	return status5.Status_Z
}

func (m *XFulfillment) GetShippingStatus() status5.Status {
	if m != nil {
		return m.ShippingStatus
	}
	return status5.Status_Z
}

func (m *XFulfillment) GetEtopPaymentStatus() status4.Status {
	if m != nil {
		return m.EtopPaymentStatus
	}
	return status4.Status_Z
}

func (m *XFulfillment) GetShippingFeeCustomer() int32 {
	if m != nil {
		return m.ShippingFeeCustomer
	}
	return 0
}

func (m *XFulfillment) GetShippingFeeShop() int32 {
	if m != nil {
		return m.ShippingFeeShop
	}
	return 0
}

func (m *XFulfillment) GetXShippingFee() int32 {
	if m != nil {
		return m.XShippingFee
	}
	return 0
}

func (m *XFulfillment) GetXShippingId() string {
	if m != nil {
		return m.XShippingId
	}
	return ""
}

func (m *XFulfillment) GetXShippingCode() string {
	if m != nil {
		return m.XShippingCode
	}
	return ""
}

func (m *XFulfillment) GetXShippingState() string {
	if m != nil {
		return m.XShippingState
	}
	return ""
}

func (m *XFulfillment) GetXShippingStatus() status5.Status {
	if m != nil {
		return m.XShippingStatus
	}
	return status5.Status_Z
}

func (m *XFulfillment) GetXSyncStatus() status4.Status {
	if m != nil {
		return m.XSyncStatus
	}
	return status4.Status_Z
}

func (m *XFulfillment) GetXSyncStates() *FulfillmentSyncStates {
	if m != nil {
		return m.XSyncStates
	}
	return nil
}

func (m *XFulfillment) GetAddressTo() *etop.Address {
	if m != nil {
		return m.AddressTo
	}
	return nil
}

func (m *XFulfillment) GetAddressFrom() *etop.Address {
	if m != nil {
		return m.AddressFrom
	}
	return nil
}

func (m *XFulfillment) GetPickupAddress() *OrderAddress {
	if m != nil {
		return m.PickupAddress
	}
	return nil
}

func (m *XFulfillment) GetReturnAddress() *OrderAddress {
	if m != nil {
		return m.ReturnAddress
	}
	return nil
}

func (m *XFulfillment) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *XFulfillment) GetShop() *etop.Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

func (m *XFulfillment) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *XFulfillment) GetProviderShippingFeeLines() []*ShippingFeeLine {
	if m != nil {
		return m.ProviderShippingFeeLines
	}
	return nil
}

func (m *XFulfillment) GetShippingFeeShopLines() []*ShippingFeeLine {
	if m != nil {
		return m.ShippingFeeShopLines
	}
	return nil
}

func (m *XFulfillment) GetEtopDiscount() int32 {
	if m != nil {
		return m.EtopDiscount
	}
	return 0
}

func (m *XFulfillment) GetMoneyTransactionShippingId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingId
	}
	return 0
}

func (m *XFulfillment) GetMoneyTransactionShippingExternalId() dot.ID {
	if m != nil {
		return m.MoneyTransactionShippingExternalId
	}
	return 0
}

func (m *XFulfillment) GetXShippingLogs() []*ExternalShippingLog {
	if m != nil {
		return m.XShippingLogs
	}
	return nil
}

func (m *XFulfillment) GetXShippingNote() string {
	if m != nil {
		return m.XShippingNote
	}
	return ""
}

func (m *XFulfillment) GetXShippingSubState() string {
	if m != nil {
		return m.XShippingSubState
	}
	return ""
}

func (m *XFulfillment) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *XFulfillment) GetActualCompensationAmount() int32 {
	if m != nil {
		return m.ActualCompensationAmount
	}
	return 0
}

type TradingCreateOrderRequest struct {
	// Customer should be shop's information
	// and customer_address, shipping_address must be shop address.
	Customer             *OrderCustomer     `protobuf:"bytes,1,opt,name=customer" json:"customer"`
	CustomerAddress      *OrderAddress      `protobuf:"bytes,2,opt,name=customer_address,json=customerAddress" json:"customer_address"`
	BillingAddress       *OrderAddress      `protobuf:"bytes,3,opt,name=billing_address,json=billingAddress" json:"billing_address"`
	ShippingAddress      *OrderAddress      `protobuf:"bytes,4,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	Lines                []*CreateOrderLine `protobuf:"bytes,5,rep,name=lines" json:"lines"`
	Discounts            []*OrderDiscount   `protobuf:"bytes,6,rep,name=discounts" json:"discounts"`
	TotalItems           int32              `protobuf:"varint,7,opt,name=total_items,json=totalItems" json:"total_items"`
	BasketValue          int32              `protobuf:"varint,8,opt,name=basket_value,json=basketValue" json:"basket_value"`
	OrderDiscount        int32              `protobuf:"varint,9,opt,name=order_discount,json=orderDiscount" json:"order_discount"`
	TotalFee             int32              `protobuf:"varint,10,opt,name=total_fee,json=totalFee" json:"total_fee"`
	FeeLines             []*OrderFeeLine    `protobuf:"bytes,11,rep,name=fee_lines,json=feeLines" json:"fee_lines"`
	TotalDiscount        *int32             `protobuf:"varint,12,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount          int32              `protobuf:"varint,13,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	OrderNote            string             `protobuf:"bytes,14,opt,name=order_note,json=orderNote" json:"order_note"`
	PaymentMethod        string             `protobuf:"bytes,15,opt,name=payment_method,json=paymentMethod" json:"payment_method"`
	ReferralMeta         map[string]string  `protobuf:"bytes,16,rep,name=referral_meta,json=referralMeta" json:"referral_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TradingCreateOrderRequest) Reset()         { *m = TradingCreateOrderRequest{} }
func (m *TradingCreateOrderRequest) String() string { return proto.CompactTextString(m) }
func (*TradingCreateOrderRequest) ProtoMessage()    {}
func (*TradingCreateOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1df1ee62235eeec9, []int{45}
}

var xxx_messageInfo_TradingCreateOrderRequest proto.InternalMessageInfo

func (m *TradingCreateOrderRequest) GetCustomer() *OrderCustomer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetCustomerAddress() *OrderAddress {
	if m != nil {
		return m.CustomerAddress
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetBillingAddress() *OrderAddress {
	if m != nil {
		return m.BillingAddress
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetShippingAddress() *OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetLines() []*CreateOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetDiscounts() []*OrderDiscount {
	if m != nil {
		return m.Discounts
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetBasketValue() int32 {
	if m != nil {
		return m.BasketValue
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetOrderDiscount() int32 {
	if m != nil {
		return m.OrderDiscount
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetTotalFee() int32 {
	if m != nil {
		return m.TotalFee
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetFeeLines() []*OrderFeeLine {
	if m != nil {
		return m.FeeLines
	}
	return nil
}

func (m *TradingCreateOrderRequest) GetTotalDiscount() int32 {
	if m != nil && m.TotalDiscount != nil {
		return *m.TotalDiscount
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *TradingCreateOrderRequest) GetOrderNote() string {
	if m != nil {
		return m.OrderNote
	}
	return ""
}

func (m *TradingCreateOrderRequest) GetPaymentMethod() string {
	if m != nil {
		return m.PaymentMethod
	}
	return ""
}

func (m *TradingCreateOrderRequest) GetReferralMeta() map[string]string {
	if m != nil {
		return m.ReferralMeta
	}
	return nil
}

func init() {
	proto.RegisterType((*OrdersResponse)(nil), "order.OrdersResponse")
	proto.RegisterType((*Order)(nil), "order.Order")
	proto.RegisterType((*OrderLineMetaField)(nil), "order.OrderLineMetaField")
	proto.RegisterType((*OrderLine)(nil), "order.OrderLine")
	proto.RegisterType((*OrderFeeLine)(nil), "order.OrderFeeLine")
	proto.RegisterType((*CreateOrderLine)(nil), "order.CreateOrderLine")
	proto.RegisterType((*OrderCustomer)(nil), "order.OrderCustomer")
	proto.RegisterType((*OrderAddress)(nil), "order.OrderAddress")
	proto.RegisterType((*OrderDiscount)(nil), "order.OrderDiscount")
	proto.RegisterType((*CreateOrderRequest)(nil), "order.CreateOrderRequest")
	proto.RegisterMapType((map[string]string)(nil), "order.CreateOrderRequest.ExternalMetaEntry")
	proto.RegisterMapType((map[string]string)(nil), "order.CreateOrderRequest.ReferralMetaEntry")
	proto.RegisterType((*UpdateOrderRequest)(nil), "order.UpdateOrderRequest")
	proto.RegisterType((*OrderShipping)(nil), "order.OrderShipping")
	proto.RegisterType((*Fulfillment)(nil), "order.Fulfillment")
	proto.RegisterType((*ShippingFeeLine)(nil), "order.ShippingFeeLine")
	proto.RegisterType((*ExternalShippingLog)(nil), "order.ExternalShippingLog")
	proto.RegisterType((*FulfillmentsResponse)(nil), "order.FulfillmentsResponse")
	proto.RegisterType((*Attribute)(nil), "order.Attribute")
	proto.RegisterType((*OrderWithErrorsResponse)(nil), "order.OrderWithErrorsResponse")
	proto.RegisterType((*GetExternalShippingServicesRequest)(nil), "order.GetExternalShippingServicesRequest")
	proto.RegisterType((*GetExternalShippingServicesResponse)(nil), "order.GetExternalShippingServicesResponse")
	proto.RegisterType((*ExternalShippingService)(nil), "order.ExternalShippingService")
	proto.RegisterType((*FulfillmentSyncStates)(nil), "order.FulfillmentSyncStates")
	proto.RegisterType((*MoneyTransaction)(nil), "order.MoneyTransaction")
	proto.RegisterType((*MoneyTransactionsResponse)(nil), "order.MoneyTransactionsResponse")
	proto.RegisterType((*MoneyTransactionShippingExternalLine)(nil), "order.MoneyTransactionShippingExternalLine")
	proto.RegisterType((*MoneyTransactionShippingExternal)(nil), "order.MoneyTransactionShippingExternal")
	proto.RegisterType((*MoneyTransactionShippingExternalsResponse)(nil), "order.MoneyTransactionShippingExternalsResponse")
	proto.RegisterType((*MoneyTransactionShippingEtop)(nil), "order.MoneyTransactionShippingEtop")
	proto.RegisterType((*MoneyTransactionShippingEtopsResponse)(nil), "order.MoneyTransactionShippingEtopsResponse")
	proto.RegisterType((*ImportOrdersResponse)(nil), "order.ImportOrdersResponse")
	proto.RegisterType((*PublicFulfillment)(nil), "order.PublicFulfillment")
	proto.RegisterType((*ShipnowFulfillment)(nil), "order.ShipnowFulfillment")
	proto.RegisterType((*GetShipnowFulfillmentsRequest)(nil), "order.GetShipnowFulfillmentsRequest")
	proto.RegisterType((*ShipnowFulfillments)(nil), "order.ShipnowFulfillments")
	proto.RegisterType((*DeliveryPoint)(nil), "order.DeliveryPoint")
	proto.RegisterType((*WeightInfo)(nil), "order.WeightInfo")
	proto.RegisterType((*ValueInfo)(nil), "order.ValueInfo")
	proto.RegisterType((*CreateShipnowFulfillmentRequest)(nil), "order.CreateShipnowFulfillmentRequest")
	proto.RegisterType((*UpdateShipnowFulfillmentRequest)(nil), "order.UpdateShipnowFulfillmentRequest")
	proto.RegisterType((*CancelShipnowFulfillmentRequest)(nil), "order.CancelShipnowFulfillmentRequest")
	proto.RegisterType((*GetShipnowServicesRequest)(nil), "order.GetShipnowServicesRequest")
	proto.RegisterType((*DeliveryPointRequest)(nil), "order.DeliveryPointRequest")
	proto.RegisterType((*GetShipnowServicesResponse)(nil), "order.GetShipnowServicesResponse")
	proto.RegisterType((*ShippnowService)(nil), "order.ShippnowService")
	proto.RegisterType((*XFulfillment)(nil), "order._Fulfillment")
	proto.RegisterType((*TradingCreateOrderRequest)(nil), "order.TradingCreateOrderRequest")
	proto.RegisterMapType((map[string]string)(nil), "order.TradingCreateOrderRequest.ReferralMetaEntry")
}

func init() { proto.RegisterFile("etop/order/order.proto", fileDescriptor_1df1ee62235eeec9) }

var fileDescriptor_1df1ee62235eeec9 = []byte{
	// 5932 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x7d, 0xdb, 0x6f, 0x1b, 0xd9,
	0x79, 0xb8, 0x79, 0x15, 0xf9, 0x51, 0x24, 0xa5, 0x91, 0x6c, 0x8d, 0x65, 0x5b, 0x92, 0xe9, 0xf5,
	0xae, 0xbd, 0x76, 0x28, 0xaf, 0xd6, 0xde, 0xf5, 0x3a, 0x7b, 0x89, 0x2c, 0x59, 0x5e, 0x6d, 0x6c,
	0xaf, 0x7e, 0xb4, 0x37, 0xf9, 0xa1, 0x28, 0xc0, 0x8e, 0x38, 0x47, 0xd2, 0xc0, 0xe4, 0x0c, 0x33,
	0x33, 0xb4, 0xa5, 0xf4, 0xad, 0x4f, 0x09, 0x5a, 0xa0, 0x05, 0x8a, 0x02, 0x41, 0xdb, 0xa7, 0x06,
	0x05, 0xda, 0xb7, 0x02, 0x45, 0x81, 0xa2, 0x45, 0x81, 0x3e, 0xe6, 0xa5, 0xe8, 0xfe, 0x05, 0x41,
	0x76, 0xf3, 0xd8, 0x16, 0x28, 0x10, 0xa0, 0x2f, 0x7d, 0x68, 0x71, 0x6e, 0xc3, 0xef, 0xcc, 0x85,
	0x1c, 0x52, 0x49, 0x53, 0x14, 0x7a, 0x59, 0x8b, 0xe7, 0xbb, 0xcc, 0x99, 0x33, 0xdf, 0xfd, 0x7c,
	0xe7, 0x2c, 0x5c, 0x20, 0xbe, 0xd3, 0x5f, 0x77, 0x5c, 0x93, 0xb8, 0xfc, 0xbf, 0xcd, 0xbe, 0xeb,
	0xf8, 0x8e, 0x56, 0x60, 0x3f, 0x96, 0x17, 0x0f, 0x9d, 0x43, 0x87, 0x8d, 0xac, 0xd3, 0xbf, 0x38,
	0x70, 0x79, 0xf5, 0xd0, 0x71, 0x0e, 0xbb, 0x64, 0x9d, 0xfd, 0xda, 0x1f, 0x1c, 0xac, 0xfb, 0x56,
	0x8f, 0x78, 0xbe, 0xd1, 0xeb, 0x0b, 0x84, 0x85, 0x8e, 0xd3, 0xeb, 0x39, 0xf6, 0x3a, 0xff, 0x47,
	0x0c, 0xbe, 0x21, 0x06, 0xbd, 0xbe, 0x4b, 0x0c, 0xd3, 0x3b, 0x22, 0xc4, 0xc7, 0x7f, 0x0b, 0xac,
	0x25, 0x36, 0x21, 0xe2, 0x77, 0xd6, 0x0f, 0x08, 0x59, 0xf7, 0x4f, 0xfa, 0x44, 0x00, 0x2e, 0x07,
	0x80, 0x43, 0x62, 0xd3, 0xe9, 0xf2, 0x7f, 0x04, 0xf4, 0xea, 0x10, 0x7a, 0x64, 0xb7, 0x6d, 0xc7,
	0x27, 0xed, 0x8e, 0x63, 0x92, 0x75, 0xfa, 0x1f, 0x81, 0x72, 0x25, 0x40, 0xf1, 0x8e, 0xac, 0x7e,
	0xdf, 0xb2, 0x0f, 0xd7, 0x3d, 0xdf, 0xf0, 0x25, 0xf8, 0x7a, 0x04, 0xdc, 0x3e, 0x20, 0xa4, 0x4d,
	0x67, 0x80, 0xa7, 0x81, 0xb8, 0xf8, 0x86, 0x3f, 0xf0, 0xde, 0x15, 0xff, 0x26, 0x80, 0xef, 0x8e,
	0x06, 0xdf, 0x53, 0xc1, 0xcb, 0x01, 0xd8, 0x77, 0x4f, 0xda, 0x6c, 0x01, 0x83, 0xe9, 0xdf, 0x8c,
	0xce, 0xaf, 0xef, 0x3a, 0xaf, 0x2c, 0xba, 0x14, 0xf2, 0x0f, 0x81, 0x5a, 0x17, 0xa8, 0x8e, 0xfc,
	0x1e, 0x2b, 0xe8, 0x2b, 0x7b, 0xce, 0xc0, 0xed, 0x10, 0xf1, 0x0f, 0x87, 0x37, 0x7e, 0x13, 0x6a,
	0x9f, 0x53, 0xa0, 0xd7, 0x22, 0x5e, 0xdf, 0xb1, 0x3d, 0xa2, 0xbd, 0x01, 0xc5, 0xbe, 0x71, 0x68,
	0xd9, 0x87, 0x7a, 0x66, 0x2d, 0x73, 0xa3, 0xb2, 0x31, 0xdb, 0xec, 0xf4, 0x9a, 0x7b, 0xc6, 0x21,
	0xd9, 0xb5, 0x0f, 0x9c, 0x96, 0x80, 0x51, 0x2c, 0xc6, 0xd4, 0xd3, 0xb3, 0x6b, 0x39, 0x86, 0xc5,
	0x65, 0x88, 0x31, 0x6b, 0x09, 0x58, 0xe3, 0x9f, 0xce, 0x43, 0x81, 0x8d, 0x68, 0x6f, 0x41, 0x9d,
	0x1c, 0xf7, 0x1d, 0xd7, 0x27, 0x66, 0xfb, 0xc0, 0x22, 0x5d, 0xd3, 0xd3, 0xcd, 0xb5, 0xdc, 0x8d,
	0x72, 0xab, 0x26, 0x87, 0x77, 0xd8, 0xa8, 0xb6, 0x08, 0x59, 0xcb, 0x64, 0x8f, 0xce, 0x3d, 0xcc,
	0xff, 0xe4, 0xa7, 0xab, 0xe7, 0x5a, 0x59, 0xcb, 0xd4, 0xae, 0xc0, 0x8c, 0x77, 0xe4, 0xf4, 0xdb,
	0x96, 0xa9, 0x67, 0x11, 0xa8, 0x48, 0x07, 0x77, 0x4d, 0xed, 0x2a, 0x94, 0x19, 0xd8, 0x36, 0x7a,
	0x44, 0xcf, 0xad, 0x65, 0x6e, 0x94, 0x05, 0x42, 0x89, 0x0e, 0x3f, 0x33, 0x7a, 0x44, 0xd3, 0x21,
	0x4f, 0x97, 0x54, 0xcf, 0x23, 0x28, 0x1b, 0xa1, 0xbc, 0x89, 0xc9, 0x64, 0x46, 0x2f, 0x20, 0x60,
	0x91, 0x98, 0x5b, 0x14, 0x7c, 0x13, 0xaa, 0xe4, 0xd8, 0x27, 0xae, 0x6d, 0x74, 0x39, 0x52, 0x09,
	0x21, 0xcd, 0x4a, 0x10, 0x43, 0xbd, 0x0d, 0x45, 0xbe, 0xb8, 0x7a, 0x71, 0x2d, 0x73, 0xa3, 0xb6,
	0x51, 0x6b, 0x8a, 0xb5, 0x7e, 0xce, 0xfe, 0x09, 0x26, 0xcd, 0x7e, 0x69, 0xd7, 0x00, 0xfa, 0x86,
	0xeb, 0xdb, 0xc4, 0xa5, 0xaf, 0xf5, 0x3e, 0x7a, 0xad, 0xb2, 0x18, 0xdf, 0x35, 0xb5, 0xeb, 0x50,
	0x09, 0x9e, 0x6e, 0x99, 0xfa, 0x7d, 0xf4, 0x6c, 0x90, 0x80, 0x5d, 0x53, 0x7b, 0x0b, 0x82, 0x99,
	0xb4, 0x07, 0x6e, 0x57, 0xff, 0x00, 0xe1, 0x05, 0x0c, 0xbe, 0x70, 0xbb, 0xda, 0x2a, 0x94, 0x3c,
	0xd2, 0x3d, 0x60, 0x48, 0x0f, 0x10, 0xd2, 0x0c, 0x1d, 0xa5, 0x08, 0xb7, 0xa0, 0xd6, 0x37, 0x4e,
	0x7a, 0xc4, 0xf6, 0xdb, 0x3d, 0xe2, 0x1f, 0x39, 0xa6, 0x3e, 0x83, 0xd0, 0xaa, 0x02, 0xf6, 0x94,
	0x81, 0xb4, 0x3b, 0x50, 0xea, 0x0c, 0x3c, 0xdf, 0xe9, 0x11, 0x57, 0x2f, 0x33, 0x69, 0x59, 0xc4,
	0x72, 0xb0, 0x25, 0x60, 0xad, 0x00, 0x4b, 0xfb, 0x18, 0xe6, 0xe4, 0xdf, 0x6d, 0xc3, 0x34, 0x5d,
	0xe2, 0x79, 0x3a, 0x30, 0xca, 0x05, 0x4c, 0xb9, 0xc9, 0x41, 0xad, 0xba, 0x44, 0x16, 0x03, 0xda,
	0x87, 0x50, 0xdf, 0xb7, 0xba, 0x5d, 0xaa, 0x03, 0x92, 0xbc, 0x92, 0x4c, 0x5e, 0x13, 0xb8, 0x92,
	0xfa, 0x63, 0x98, 0x0b, 0x54, 0x48, 0x92, 0xcf, 0x8e, 0x78, 0xba, 0x44, 0x96, 0xf4, 0x1f, 0x00,
	0x74, 0x5c, 0x62, 0x50, 0x21, 0x36, 0x7c, 0xbd, 0xca, 0x28, 0x97, 0x9b, 0xdc, 0x26, 0x36, 0xa5,
	0x4d, 0x6c, 0xbe, 0x90, 0x36, 0xb1, 0x55, 0x16, 0xd8, 0x9b, 0xbe, 0xf6, 0x11, 0xcc, 0xf6, 0x5d,
	0xa7, 0x43, 0x3c, 0x8f, 0x13, 0xd7, 0xc6, 0x12, 0x57, 0x02, 0xfc, 0x4d, 0x9f, 0x3e, 0x79, 0xd0,
	0x37, 0xe5, 0x93, 0xeb, 0xe3, 0x9f, 0x2c, 0xb0, 0x37, 0x7d, 0xed, 0x7d, 0x28, 0x77, 0xba, 0x8e,
	0x78, 0xec, 0xdc, 0x58, 0xca, 0x12, 0x47, 0xe6, 0x53, 0xee, 0x38, 0xf6, 0x81, 0xe5, 0xf6, 0x38,
	0xed, 0xfc, 0xf8, 0x29, 0x07, 0xf8, 0x82, 0xdc, 0xb0, 0x3b, 0xa4, 0xdb, 0xe5, 0xe4, 0x5a, 0x0a,
	0x72, 0x89, 0xbf, 0xe9, 0x53, 0xbd, 0xe3, 0x3f, 0xdb, 0x2e, 0x31, 0x3c, 0xc7, 0xd6, 0x17, 0xb0,
	0xde, 0x71, 0x50, 0x8b, 0x41, 0xb4, 0xbb, 0x00, 0xde, 0x51, 0x5b, 0x3c, 0x5b, 0x5f, 0x65, 0xba,
	0x57, 0x6f, 0x0a, 0x2b, 0xdd, 0x7c, 0xce, 0xfe, 0x95, 0xaa, 0xe5, 0x1d, 0x6d, 0x71, 0x3c, 0x6d,
	0x1d, 0x66, 0x24, 0xc9, 0xd5, 0x51, 0x24, 0x12, 0x4b, 0xfb, 0x10, 0x6a, 0xe2, 0xcf, 0x36, 0x47,
	0xd4, 0x77, 0x46, 0xd1, 0x55, 0x05, 0x32, 0x1f, 0xd4, 0xbe, 0x01, 0x45, 0x41, 0xd5, 0x50, 0xa8,
	0xee, 0xa9, 0x54, 0x02, 0x49, 0xdb, 0x06, 0xed, 0x60, 0xd0, 0x3d, 0xb0, 0xba, 0x5d, 0xa6, 0x8b,
	0x82, 0xf4, 0xda, 0x28, 0xd2, 0x79, 0x44, 0x20, 0x1e, 0xfa, 0x05, 0x5c, 0x52, 0xb8, 0x48, 0xe1,
	0x17, 0xec, 0xde, 0x18, 0xc5, 0xee, 0x22, 0x66, 0x27, 0x08, 0x05, 0xdb, 0xa7, 0xb0, 0x14, 0x68,
	0xb1, 0xb4, 0x16, 0x82, 0xe5, 0xf5, 0x51, 0x4b, 0x72, 0x5e, 0x52, 0xed, 0x71, 0x22, 0xc1, 0xee,
	0x11, 0x2c, 0x50, 0x37, 0x15, 0x66, 0xf5, 0x48, 0x61, 0x75, 0x37, 0xf4, 0xb2, 0x94, 0x42, 0x65,
	0xf3, 0x26, 0x14, 0xba, 0x96, 0x4d, 0x3c, 0xfd, 0x26, 0x73, 0x49, 0x73, 0x58, 0xa5, 0x9f, 0x58,
	0x36, 0x69, 0x71, 0xb0, 0xb6, 0x01, 0x65, 0xd3, 0xf2, 0x3a, 0xce, 0xc0, 0xf6, 0x3d, 0xfd, 0x6d,
	0x86, 0xab, 0x98, 0xad, 0x6d, 0x01, 0x6c, 0x0d, 0xd1, 0xa8, 0x1d, 0xf6, 0x1d, 0x9f, 0x1a, 0x61,
	0x9f, 0xf4, 0x3c, 0xfd, 0xd6, 0x5a, 0xe6, 0x46, 0x41, 0xda, 0x61, 0x06, 0xd8, 0xa5, 0xe3, 0xd4,
	0x0e, 0xef, 0x1b, 0xde, 0x4b, 0xe2, 0xb7, 0x5f, 0x19, 0xdd, 0x01, 0xd1, 0x6f, 0x23, 0xbc, 0x0a,
	0x87, 0x7c, 0x87, 0x02, 0x28, 0x22, 0xe7, 0xf7, 0x9a, 0x58, 0x87, 0x47, 0xbe, 0xfe, 0x0d, 0x8c,
	0xc8, 0x20, 0xdf, 0x65, 0x00, 0x6a, 0x8f, 0xd9, 0xd4, 0xda, 0x72, 0x2e, 0x7a, 0x13, 0xa1, 0x56,
	0x1d, 0x3c, 0x61, 0x8a, 0xcc, 0xb9, 0x06, 0xc8, 0xeb, 0x18, 0x99, 0xc1, 0x02, 0xe4, 0x60, 0x0a,
	0x46, 0x8f, 0xa1, 0xde, 0x89, 0x4c, 0x61, 0x93, 0x01, 0xa8, 0xa3, 0xe2, 0x53, 0xa0, 0xc1, 0x95,
	0xbe, 0x81, 0xd4, 0xb0, 0xcc, 0xc6, 0x9f, 0x39, 0x3e, 0x7b, 0x21, 0x1c, 0x3d, 0xe9, 0x77, 0x31,
	0x37, 0x09, 0xd9, 0x21, 0x84, 0xfa, 0x6a, 0xfe, 0x58, 0x8a, 0x75, 0x0f, 0x61, 0x95, 0xd8, 0x30,
	0x45, 0xb9, 0x03, 0x65, 0x1a, 0x80, 0xf1, 0x8f, 0xf9, 0x1e, 0xfb, 0x40, 0x8a, 0x7d, 0xde, 0x21,
	0x84, 0x7d, 0xcf, 0xd2, 0x01, 0xff, 0xc3, 0xd3, 0xee, 0xc0, 0x3c, 0x0b, 0x00, 0x94, 0x29, 0x7c,
	0x84, 0x98, 0xd7, 0x29, 0xf8, 0x39, 0x9a, 0xc6, 0x4d, 0xa8, 0x06, 0xc8, 0xec, 0xbd, 0xde, 0xc5,
	0xe6, 0x45, 0x82, 0xd8, 0xab, 0x51, 0x9f, 0x49, 0x99, 0x77, 0x1c, 0x53, 0xff, 0x18, 0xf1, 0x64,
	0x21, 0xc9, 0x96, 0x63, 0x52, 0x5e, 0x2e, 0x39, 0x20, 0x2e, 0xb1, 0x3b, 0x84, 0x79, 0xd6, 0x4f,
	0x30, 0xaf, 0x00, 0x44, 0xdd, 0xeb, 0xfb, 0x30, 0x8b, 0xd4, 0xca, 0xd3, 0x3f, 0x53, 0xde, 0xae,
	0xbd, 0x33, 0x84, 0xb5, 0x14, 0x44, 0xed, 0x03, 0x3a, 0x5f, 0xf4, 0x86, 0xfa, 0xc3, 0xa8, 0xbf,
	0x95, 0xef, 0x47, 0xe7, 0x3f, 0x7c, 0x5b, 0xea, 0xa5, 0x03, 0xaa, 0xc7, 0x23, 0xa8, 0x02, 0x2c,
	0xed, 0x06, 0x54, 0x95, 0x60, 0x5a, 0xdf, 0xc6, 0xf1, 0xc4, 0xe1, 0x91, 0x4d, 0xd7, 0x85, 0x85,
	0x3c, 0xeb, 0x30, 0x87, 0x0d, 0x0c, 0x0d, 0x97, 0xf5, 0x4f, 0x11, 0x72, 0x1d, 0x41, 0x5f, 0x9c,
	0xf4, 0xa9, 0x9c, 0xe0, 0xa1, 0xb6, 0x65, 0x7a, 0xfa, 0xee, 0x5a, 0xee, 0x46, 0xae, 0x55, 0x43,
	0xc3, 0xbb, 0x26, 0xb5, 0x97, 0x75, 0x97, 0x74, 0x88, 0xf5, 0x8a, 0x7a, 0x0f, 0x2e, 0xa1, 0xdf,
	0x46, 0x8b, 0x5f, 0x93, 0x40, 0x21, 0xa4, 0xd7, 0xa1, 0x12, 0x98, 0x24, 0xcb, 0xd4, 0x9f, 0xa0,
	0x70, 0x0a, 0x24, 0x60, 0xd7, 0xa4, 0x36, 0x3c, 0x64, 0x65, 0x9e, 0x8e, 0xb2, 0x32, 0x32, 0xde,
	0xe1, 0x83, 0x8d, 0x7d, 0xd0, 0x02, 0x6b, 0xf2, 0x94, 0xf8, 0x06, 0x8b, 0x59, 0xb5, 0x0b, 0x90,
	0x7b, 0x49, 0x4e, 0x58, 0xcc, 0x2a, 0x5f, 0x9b, 0x0e, 0x68, 0xcb, 0x50, 0xe0, 0x56, 0x20, 0x8b,
	0x20, 0x7c, 0x88, 0x86, 0xa3, 0x91, 0x60, 0x95, 0x8d, 0x34, 0xfe, 0x39, 0x0f, 0xe5, 0xe0, 0x21,
	0xe9, 0xe3, 0xe6, 0x55, 0x28, 0x71, 0x25, 0x0d, 0x85, 0xc8, 0x33, 0x6c, 0x74, 0xd7, 0xa4, 0x5a,
	0xfc, 0xca, 0x70, 0x2d, 0x83, 0x2d, 0x3a, 0x7b, 0x6e, 0x10, 0x6e, 0x8a, 0x71, 0x1e, 0x47, 0xf6,
	0x5d, 0xc7, 0x1c, 0x74, 0x7c, 0x1e, 0x4b, 0xe3, 0xd8, 0xaf, 0x22, 0x20, 0x2c, 0x9c, 0xbe, 0x0d,
	0x75, 0xcb, 0x6b, 0x3b, 0x03, 0xdf, 0xb3, 0x4c, 0xd2, 0xa6, 0xb6, 0x98, 0xc5, 0xc5, 0x25, 0xb9,
	0x6e, 0x96, 0xf7, 0x39, 0x87, 0x3d, 0xf2, 0x9d, 0xbe, 0xb6, 0x06, 0xa5, 0xef, 0x0d, 0x0c, 0xdb,
	0xb7, 0xfc, 0x13, 0x16, 0x45, 0x04, 0x2a, 0x2f, 0x47, 0xe9, 0xec, 0xba, 0x96, 0xe7, 0xb7, 0xfb,
	0xae, 0xd5, 0x21, 0x2c, 0x54, 0x90, 0x38, 0x65, 0x3a, 0xbe, 0x47, 0x87, 0xe9, 0xec, 0x5c, 0xe2,
	0x1b, 0x56, 0x57, 0xa0, 0x2d, 0x60, 0x1b, 0xc3, 0x21, 0x1c, 0xf1, 0x26, 0xc8, 0x0f, 0x27, 0x30,
	0x17, 0x11, 0xe6, 0xac, 0x00, 0x71, 0xd4, 0xab, 0x50, 0xb6, 0x7a, 0xc6, 0x21, 0xd7, 0xdb, 0xf3,
	0x38, 0x75, 0x60, 0xc3, 0x54, 0x67, 0xef, 0x00, 0x18, 0xbe, 0xef, 0x5a, 0xfb, 0x03, 0x9f, 0x78,
	0xfa, 0x05, 0xc5, 0xb9, 0x6c, 0x4a, 0x40, 0x0b, 0xe1, 0xb0, 0xd0, 0x5e, 0x2c, 0xa3, 0x65, 0xea,
	0x4b, 0x4a, 0x68, 0xcf, 0xc7, 0x77, 0xcd, 0x18, 0x63, 0xad, 0x27, 0x1b, 0xeb, 0x07, 0x50, 0xe9,
	0x11, 0xdf, 0x90, 0x32, 0x70, 0x91, 0x4d, 0xe2, 0x62, 0xd8, 0xc3, 0x05, 0x32, 0xd9, 0x82, 0x9e,
	0xfc, 0xd3, 0x6b, 0xfc, 0x59, 0x06, 0x66, 0xb1, 0xdd, 0xd4, 0xde, 0x84, 0x3c, 0x53, 0xd4, 0x0c,
	0x13, 0xfd, 0xd9, 0xe6, 0x01, 0x21, 0xcd, 0x1d, 0x42, 0xa8, 0x7e, 0x4a, 0x51, 0xa4, 0xf0, 0x40,
	0x48, 0xb3, 0x61, 0x21, 0x0d, 0xb2, 0xa9, 0x5c, 0x24, 0x9b, 0xd2, 0x21, 0x6f, 0x12, 0xaf, 0xa3,
	0xa4, 0x52, 0x6c, 0x44, 0xbb, 0x0c, 0x45, 0xa1, 0xc7, 0x45, 0xf4, 0x9e, 0x62, 0xac, 0xf1, 0xbb,
	0x39, 0xa8, 0x6f, 0xb1, 0x68, 0x79, 0x28, 0xfc, 0xaa, 0xc8, 0x66, 0xd2, 0x89, 0x6c, 0x36, 0x49,
	0x64, 0xb1, 0x10, 0xe6, 0x52, 0x08, 0xe1, 0x4c, 0x3a, 0x21, 0xcc, 0xa7, 0x16, 0xc2, 0xc2, 0xff,
	0xac, 0x10, 0x86, 0x44, 0x66, 0x69, 0x12, 0x91, 0xf9, 0x71, 0x16, 0xaa, 0x4a, 0x0a, 0x97, 0xde,
	0x10, 0x5d, 0x03, 0x38, 0xb0, 0x5c, 0x4f, 0x7c, 0x0d, 0x6c, 0x14, 0xcb, 0x6c, 0x9c, 0x7d, 0x8b,
	0xab, 0x50, 0xee, 0x1a, 0x5e, 0xcc, 0x17, 0x2b, 0xd1, 0x61, 0x89, 0x72, 0x30, 0xe8, 0x76, 0x63,
	0x72, 0x7a, 0x3a, 0xcc, 0x50, 0x96, 0xa1, 0x40, 0x7a, 0x86, 0xd5, 0x55, 0x92, 0x7a, 0x3e, 0x44,
	0x61, 0xfd, 0x23, 0xc7, 0x56, 0x73, 0x7a, 0x3e, 0x44, 0xf3, 0x74, 0x5e, 0x42, 0x0a, 0xf2, 0x74,
	0x51, 0x51, 0x7a, 0xcc, 0xfe, 0x91, 0x92, 0xc9, 0x07, 0xa9, 0x44, 0x33, 0x6d, 0xc1, 0xb6, 0x90,
	0x8d, 0x34, 0xfe, 0xb4, 0x20, 0x14, 0x4b, 0xe6, 0x87, 0xa9, 0x17, 0x49, 0x79, 0xb9, 0x4c, 0xec,
	0xcb, 0xa9, 0xeb, 0x98, 0x4d, 0xb1, 0x8e, 0xb9, 0xd8, 0x75, 0x0c, 0x16, 0x22, 0x1f, 0x5d, 0x88,
	0x60, 0x01, 0xe7, 0xa3, 0x0b, 0xb8, 0x42, 0xd3, 0xa3, 0x81, 0xed, 0xbb, 0x27, 0xca, 0x12, 0xca,
	0x41, 0x66, 0x02, 0xa8, 0x2a, 0x15, 0x15, 0x13, 0x40, 0xd5, 0x68, 0x0d, 0x4a, 0xac, 0x2c, 0x65,
	0x77, 0xd4, 0x45, 0x0b, 0x46, 0x29, 0x86, 0x69, 0x79, 0xbe, 0x6b, 0x75, 0x7c, 0xa5, 0x9c, 0x12,
	0x8c, 0x52, 0xee, 0xaf, 0x0d, 0xd7, 0x64, 0x55, 0x85, 0x80, 0x3b, 0x1d, 0xa1, 0xde, 0xf6, 0xfb,
	0x56, 0x9f, 0x15, 0x0d, 0x02, 0x6f, 0xfb, 0x7d, 0xab, 0xcf, 0xe7, 0xdb, 0xeb, 0x1b, 0xf6, 0x09,
	0xab, 0x08, 0xa0, 0xf9, 0xb2, 0x41, 0xfa, 0x4c, 0x91, 0xf2, 0xbf, 0xc3, 0x72, 0xfe, 0xe0, 0x99,
	0x72, 0x14, 0x61, 0x6c, 0xb0, 0xdc, 0x3e, 0x8c, 0xb1, 0xc1, 0x54, 0x5a, 0xbc, 0x03, 0x8f, 0x8b,
	0x6a, 0x38, 0xd0, 0x93, 0x20, 0x59, 0x36, 0x92, 0x2f, 0xc3, 0x51, 0xeb, 0x18, 0x55, 0x82, 0x18,
	0xea, 0x55, 0x28, 0xd3, 0x37, 0xe3, 0x68, 0x73, 0xf8, 0xc1, 0x74, 0x98, 0xa1, 0xbc, 0x0b, 0x95,
	0x8e, 0xe3, 0xb8, 0xa6, 0x65, 0x1b, 0x54, 0xfd, 0x79, 0x2a, 0x3d, 0xdf, 0x64, 0x85, 0xbe, 0xad,
	0x21, 0xa0, 0x85, 0xb1, 0x1a, 0x86, 0xd0, 0xe1, 0xc0, 0x89, 0x48, 0xab, 0x9d, 0x89, 0xb3, 0xda,
	0x4c, 0xc6, 0xb3, 0x61, 0x19, 0x47, 0x56, 0x3b, 0x17, 0x63, 0xb5, 0xff, 0xa1, 0x0a, 0x1a, 0xb2,
	0xda, 0x2d, 0xf2, 0xbd, 0x01, 0xf1, 0x7c, 0x54, 0x08, 0xcb, 0xa4, 0x28, 0x84, 0x85, 0x6a, 0x5c,
	0xd9, 0x84, 0x1a, 0x57, 0xa4, 0x10, 0x97, 0x4b, 0x2c, 0xc4, 0x85, 0xcb, 0x61, 0x85, 0xa4, 0x72,
	0xd8, 0x59, 0xb5, 0x6b, 0x44, 0xb5, 0xeb, 0x3d, 0x60, 0x79, 0x44, 0x40, 0x5b, 0x4d, 0xa6, 0xad,
	0x50, 0x44, 0x49, 0xd7, 0x4c, 0x51, 0x8e, 0xc1, 0x85, 0x98, 0xdb, 0x6a, 0xde, 0x7e, 0x41, 0x3c,
	0x20, 0x14, 0x0d, 0xfc, 0xdf, 0xcf, 0xde, 0x95, 0xcc, 0xf8, 0xc3, 0xf1, 0x99, 0xf1, 0x47, 0x69,
	0x32, 0xe3, 0xeb, 0xf1, 0x25, 0x81, 0x5f, 0x6d, 0x31, 0x60, 0x82, 0xe4, 0x3a, 0x36, 0x73, 0xbf,
	0x3b, 0x2a, 0x73, 0xc7, 0xe9, 0xf8, 0xbd, 0x54, 0xe9, 0xf8, 0x7b, 0x89, 0xe9, 0x78, 0x24, 0xab,
	0x7e, 0x7f, 0xaa, 0xac, 0xfa, 0x83, 0x54, 0x59, 0xf5, 0x76, 0x38, 0xab, 0xbe, 0xcf, 0x54, 0x63,
	0xb9, 0xa9, 0x8c, 0x36, 0x1f, 0x7f, 0xfa, 0x4c, 0xa6, 0xd7, 0x71, 0x19, 0xf7, 0x1e, 0x32, 0x83,
	0x34, 0x62, 0xd3, 0x3f, 0x66, 0x62, 0x70, 0x2b, 0xaa, 0x35, 0xc2, 0x1a, 0x37, 0x1f, 0x09, 0x74,
	0x1a, 0xea, 0x3d, 0xa2, 0xbe, 0x7b, 0x68, 0x2d, 0xe9, 0x10, 0xe5, 0xc8, 0x16, 0xc5, 0x95, 0x1c,
	0x3f, 0x19, 0xc7, 0xb1, 0x25, 0xd0, 0x11, 0x47, 0x17, 0x0d, 0x85, 0x93, 0xf1, 0x6f, 0xc5, 0x27,
	0xe3, 0xcb, 0x9f, 0xc0, 0x7c, 0x64, 0x6e, 0xda, 0x1c, 0xca, 0xa6, 0x79, 0x1e, 0xbd, 0xa8, 0xe4,
	0xd1, 0x22, 0x83, 0x7e, 0x90, 0xbd, 0x9f, 0xa1, 0x0c, 0x22, 0x53, 0x99, 0x84, 0x41, 0xe3, 0x3f,
	0x66, 0x40, 0xfb, 0x82, 0x55, 0xca, 0x15, 0xff, 0x15, 0xbf, 0x09, 0x85, 0xed, 0x7f, 0x76, 0x6a,
	0xfb, 0x9f, 0x3b, 0x9d, 0xfd, 0xcf, 0x9f, 0xce, 0xfe, 0x17, 0x4e, 0x61, 0xff, 0xcb, 0x29, 0xed,
	0xbf, 0x6a, 0x22, 0x8a, 0x29, 0x4d, 0xc4, 0x4c, 0xa2, 0x89, 0x78, 0x3b, 0xce, 0x44, 0x94, 0x98,
	0x15, 0x8b, 0x18, 0x87, 0x8b, 0xc8, 0x38, 0x54, 0x18, 0x4a, 0x60, 0x16, 0x22, 0xba, 0x0e, 0x53,
	0xe9, 0x7a, 0x3d, 0x95, 0xae, 0x4f, 0x5e, 0xc2, 0xbc, 0x1e, 0x71, 0x15, 0x55, 0x6e, 0xa8, 0x55,
	0x27, 0x11, 0x76, 0x3d, 0xb5, 0x24, 0xd7, 0xf3, 0x0e, 0xcc, 0x77, 0x8e, 0x0c, 0xf7, 0x90, 0x18,
	0xfb, 0x5d, 0x22, 0xb1, 0x71, 0xc5, 0x64, 0x6e, 0x08, 0x16, 0x24, 0x81, 0x23, 0x5e, 0x4c, 0xe3,
	0x88, 0xc3, 0xde, 0xf2, 0xfc, 0x58, 0x6f, 0x29, 0x7c, 0xcb, 0x85, 0x24, 0xdf, 0x12, 0x72, 0xd3,
	0x4b, 0x09, 0x6e, 0xfa, 0x12, 0xf6, 0x93, 0xac, 0x66, 0x82, 0x3c, 0x64, 0xc8, 0xf4, 0x5c, 0x8c,
	0x37, 0x3d, 0x8d, 0x1f, 0x94, 0x44, 0x70, 0x1c, 0x7c, 0xe3, 0xd4, 0xb9, 0xdb, 0x06, 0x0b, 0x6f,
	0xa4, 0x52, 0x64, 0x92, 0x95, 0xa2, 0xec, 0x1d, 0x49, 0x95, 0x78, 0x13, 0x66, 0x8f, 0xdb, 0x1e,
	0x71, 0x5f, 0x59, 0x1d, 0x12, 0x89, 0x71, 0x8f, 0x9f, 0x73, 0xc0, 0xae, 0xa9, 0xbd, 0x0d, 0xb5,
	0x63, 0x55, 0xce, 0x71, 0xd4, 0x3d, 0x7b, 0x8c, 0x45, 0x9d, 0xe3, 0x0a, 0x9e, 0x2c, 0x01, 0xc4,
	0x19, 0xde, 0xac, 0xe4, 0xca, 0x92, 0xc0, 0x07, 0x50, 0xeb, 0x5b, 0x9d, 0x97, 0x83, 0xa1, 0x32,
	0x9f, 0x4f, 0x9e, 0x77, 0x95, 0xa3, 0xca, 0xb9, 0x3f, 0x80, 0x9a, 0x4b, 0xfc, 0x81, 0x6b, 0x07,
	0xb4, 0xcb, 0x23, 0x68, 0x39, 0xaa, 0xa4, 0xbd, 0x0f, 0xe7, 0x87, 0x7b, 0x4e, 0x78, 0xaa, 0x3a,
	0x9a, 0xea, 0x82, 0x44, 0xc1, 0x33, 0x8e, 0xa3, 0x64, 0x4e, 0xf3, 0xc2, 0x08, 0x4a, 0xe6, 0x20,
	0xdf, 0x83, 0xc5, 0x08, 0x25, 0x5d, 0x49, 0x2c, 0x4e, 0x5a, 0x88, 0x90, 0xae, 0xe7, 0x77, 0xa8,
	0x99, 0x09, 0xf5, 0x57, 0x30, 0x7b, 0x59, 0xdb, 0xb8, 0xd6, 0x8c, 0x40, 0x9a, 0xf2, 0x53, 0xec,
	0x89, 0x01, 0xa9, 0x55, 0x5e, 0x68, 0x5c, 0xdb, 0x82, 0x99, 0x8e, 0xe1, 0xba, 0x16, 0x71, 0x99,
	0xb2, 0x4e, 0xc4, 0x4d, 0x52, 0x52, 0x6d, 0xb6, 0xec, 0x4e, 0x77, 0x60, 0x92, 0xb6, 0x65, 0x7b,
	0x03, 0xd7, 0xb0, 0x45, 0x97, 0x81, 0xac, 0xb8, 0xce, 0x09, 0xf0, 0xae, 0x84, 0x6a, 0x4d, 0x28,
	0xf2, 0x5e, 0x12, 0x66, 0x5a, 0x6b, 0x1b, 0xf3, 0x4d, 0xfe, 0xb3, 0xf9, 0xc2, 0x3d, 0xf9, 0xdc,
	0x46, 0xe1, 0x45, 0xc1, 0xa7, 0x03, 0x51, 0x8b, 0x5c, 0x4a, 0xb4, 0xc8, 0x57, 0x00, 0x3a, 0x4e,
	0x50, 0x96, 0x2f, 0x33, 0x15, 0x2c, 0x77, 0x1c, 0x59, 0x8b, 0xbf, 0x00, 0x45, 0x61, 0x6f, 0x2e,
	0x32, 0x90, 0xf8, 0xa5, 0x5d, 0x85, 0xd9, 0x43, 0xd7, 0xf1, 0x3c, 0x69, 0x8d, 0x80, 0x41, 0x2b,
	0x6c, 0x4c, 0x98, 0xa0, 0x0b, 0x50, 0xec, 0x12, 0xfb, 0xd0, 0x3f, 0x12, 0xd6, 0x5b, 0xfc, 0xa2,
	0x2e, 0xfc, 0xb5, 0x65, 0xfa, 0x47, 0x2c, 0x81, 0x29, 0xb4, 0xf8, 0x0f, 0x8a, 0x7d, 0xc4, 0x59,
	0x71, 0x5b, 0x29, 0x7e, 0x69, 0xb7, 0xe2, 0x6c, 0x5f, 0x9d, 0xa1, 0x44, 0xac, 0x5e, 0xe3, 0x17,
	0x2b, 0x50, 0x41, 0xfb, 0x2e, 0xa7, 0x6d, 0x55, 0x19, 0x5b, 0x88, 0x47, 0xbd, 0x2c, 0xb9, 0x98,
	0x5e, 0x96, 0x54, 0x6d, 0x21, 0xb8, 0x8d, 0xe3, 0x9b, 0x71, 0x6d, 0x1c, 0xc1, 0x5e, 0x68, 0x61,
	0xf4, 0x5e, 0x68, 0xc8, 0xe4, 0x16, 0x93, 0x33, 0x23, 0xc5, 0xeb, 0xcc, 0x24, 0x79, 0x9d, 0xb0,
	0x53, 0x80, 0x24, 0xa7, 0xd0, 0x84, 0x39, 0xce, 0x31, 0x2c, 0x48, 0x72, 0x7f, 0x87, 0x41, 0xb7,
	0x02, 0x99, 0xba, 0xa6, 0x88, 0x1c, 0xce, 0x8e, 0x90, 0xe0, 0x85, 0x3d, 0x4d, 0x29, 0xc9, 0xd3,
	0xc4, 0x3a, 0xc7, 0x4f, 0x47, 0x3a, 0x47, 0xb5, 0xf7, 0xa3, 0x32, 0x49, 0xef, 0x87, 0xda, 0xbc,
	0x31, 0x3b, 0x75, 0xf3, 0x46, 0x6d, 0xc2, 0xe6, 0x0d, 0xdc, 0x7d, 0x51, 0x3d, 0x65, 0xf7, 0x45,
	0x3d, 0xb1, 0xfb, 0xe2, 0x9d, 0x38, 0xbb, 0x89, 0x0b, 0x8a, 0x23, 0x4d, 0xe2, 0xe3, 0xa9, 0x4d,
	0x62, 0xa2, 0x6f, 0xd9, 0x19, 0xe7, 0x5b, 0x92, 0x3c, 0xc4, 0xf6, 0x18, 0x0f, 0x91, 0xe8, 0x93,
	0x1e, 0x8d, 0xf3, 0x49, 0xd8, 0xb6, 0x32, 0x0a, 0x2d, 0xce, 0xb6, 0x46, 0x50, 0x99, 0x19, 0x5e,
	0x48, 0x34, 0xc3, 0x43, 0x0b, 0xbf, 0x98, 0xca, 0xc2, 0xc7, 0x3a, 0x91, 0x4f, 0x46, 0x3a, 0x11,
	0xb5, 0xb5, 0x66, 0x29, 0x65, 0x6b, 0xcd, 0x87, 0x50, 0x53, 0x5a, 0x4d, 0x08, 0x73, 0x04, 0x8c,
	0x52, 0x0c, 0x33, 0x52, 0x39, 0xbd, 0xe0, 0x85, 0xd9, 0x20, 0xea, 0x94, 0x59, 0x4e, 0xd3, 0x29,
	0xf3, 0x31, 0xd4, 0xc3, 0x7d, 0x2d, 0xdf, 0x1a, 0x45, 0x57, 0xf3, 0xd4, 0x66, 0x96, 0x84, 0xee,
	0x93, 0x0b, 0x13, 0x76, 0x9f, 0x60, 0xe1, 0xa0, 0xa1, 0x7f, 0x90, 0x2a, 0x5e, 0x42, 0x52, 0xb5,
	0x80, 0x3a, 0x21, 0x82, 0xad, 0x95, 0x3b, 0x48, 0x81, 0x28, 0x25, 0x75, 0x04, 0xfa, 0x65, 0xb5,
	0x04, 0x12, 0x50, 0x3d, 0x3f, 0x72, 0xfa, 0x31, 0x61, 0xe2, 0x95, 0xc4, 0x30, 0xf1, 0x06, 0x54,
	0x11, 0xae, 0x65, 0xea, 0x2b, 0xb8, 0x18, 0x1a, 0xa0, 0xee, 0x9a, 0xda, 0x6d, 0xa8, 0x23, 0x4c,
	0x26, 0xa6, 0xab, 0xb8, 0x1a, 0x1a, 0xe0, 0x32, 0x39, 0x7d, 0x0a, 0xe7, 0x31, 0xf6, 0xd0, 0x34,
	0xae, 0x8d, 0xb5, 0x34, 0xda, 0x90, 0x53, 0x60, 0x23, 0x55, 0x76, 0xc8, 0x5c, 0x5e, 0x9d, 0x80,
	0xdd, 0x17, 0x81, 0xdd, 0xfc, 0x7f, 0xb0, 0x84, 0x67, 0x87, 0x2d, 0x61, 0x63, 0x2c, 0xc3, 0xc5,
	0xe1, 0xfc, 0x90, 0x49, 0x54, 0x59, 0x9a, 0xa4, 0x6b, 0xbd, 0x22, 0x2e, 0x67, 0x79, 0x6d, 0x02,
	0x96, 0xdb, 0x92, 0x70, 0xd3, 0xd7, 0x3e, 0x87, 0x0b, 0x88, 0x25, 0x0f, 0x9d, 0x39, 0xc7, 0x37,
	0xc6, 0x72, 0x5c, 0x08, 0x38, 0xb6, 0x04, 0xdd, 0xa6, 0xaf, 0x3d, 0x81, 0x45, 0x72, 0xdc, 0x27,
	0x1d, 0xba, 0x76, 0x62, 0x86, 0x27, 0x94, 0xdd, 0xcd, 0xf1, 0x8b, 0x28, 0xe9, 0xc4, 0xfc, 0x4e,
	0x36, 0x7d, 0x6d, 0x1b, 0xe6, 0x02, 0x6e, 0x34, 0x27, 0xa0, 0x9c, 0xde, 0x1d, 0xcb, 0xa9, 0x26,
	0x69, 0xf6, 0xac, 0xce, 0xcb, 0x4d, 0x5f, 0x7b, 0x06, 0xe7, 0x89, 0xe7, 0x5b, 0x3d, 0x23, 0x3c,
	0xa9, 0xcf, 0xc6, 0xbf, 0x63, 0x40, 0x88, 0x66, 0xf5, 0x19, 0x0c, 0x87, 0xdb, 0x32, 0xab, 0xe1,
	0xcd, 0x21, 0xa3, 0xb9, 0xcd, 0x07, 0x64, 0x7b, 0x3c, 0xc1, 0x61, 0x1f, 0x80, 0x46, 0x15, 0x4c,
	0xff, 0x7d, 0xd7, 0xb0, 0xbd, 0x03, 0xf9, 0x49, 0x6f, 0x8d, 0x9f, 0x5c, 0xc7, 0x31, 0x1f, 0xf9,
	0x4e, 0xff, 0x45, 0x40, 0xb7, 0xe9, 0x6b, 0xbf, 0x05, 0xab, 0x11, 0x5d, 0x0e, 0x71, 0xbe, 0x3d,
	0x96, 0xf3, 0xa5, 0x90, 0xbe, 0x2b, 0x4f, 0x68, 0xc2, 0xdc, 0x71, 0x3b, 0x64, 0x5d, 0xaf, 0x23,
	0x35, 0xad, 0x1d, 0x3f, 0x57, 0xac, 0xe9, 0x26, 0xcc, 0x1f, 0x47, 0x1a, 0xff, 0xde, 0x1c, 0x65,
	0x20, 0xeb, 0xc7, 0xa1, 0x76, 0xbf, 0x0f, 0x98, 0x09, 0x39, 0xb1, 0x3b, 0x92, 0xfc, 0xad, 0x51,
	0xb6, 0xb1, 0x72, 0xfc, 0xfc, 0xc4, 0xee, 0x08, 0xd2, 0x6f, 0x29, 0xa4, 0xc4, 0xd3, 0x6f, 0xb0,
	0xb7, 0xbf, 0x2c, 0xe2, 0x51, 0x14, 0x77, 0x4b, 0x22, 0xe2, 0x21, 0x0e, 0xc4, 0xd3, 0x6e, 0x03,
	0x88, 0xbc, 0xb3, 0xed, 0x3b, 0xfa, 0xdb, 0x8c, 0xbc, 0xca, 0x77, 0xbe, 0x82, 0x44, 0x5b, 0x20,
	0xbc, 0x70, 0xb4, 0x3b, 0x30, 0x2b, 0xb1, 0x0f, 0x5c, 0xa7, 0xc7, 0x0a, 0xd4, 0x11, 0xfc, 0x8a,
	0x40, 0xd9, 0x71, 0x9d, 0x5e, 0x4c, 0x6a, 0xbc, 0x79, 0x8a, 0xd4, 0xf8, 0x61, 0xea, 0xd4, 0x38,
	0xae, 0x3a, 0xb7, 0x35, 0x41, 0x75, 0x6e, 0x05, 0xf2, 0xcc, 0x51, 0x7c, 0x83, 0xd1, 0x00, 0x7f,
	0x43, 0x2a, 0x2d, 0x2d, 0x36, 0xae, 0x35, 0x80, 0x9f, 0xe4, 0x60, 0x3b, 0x0e, 0xe1, 0x06, 0x7d,
	0x0e, 0xd2, 0xbe, 0x80, 0x4b, 0x32, 0xdc, 0x52, 0xdc, 0x89, 0xa8, 0x5b, 0xad, 0x2b, 0x65, 0x20,
	0xe4, 0x54, 0x58, 0x06, 0xa1, 0x4b, 0xd2, 0x10, 0x80, 0xb5, 0x87, 0x46, 0x95, 0x80, 0xb3, 0xbc,
	0x33, 0x92, 0xe5, 0x62, 0x48, 0xf0, 0x39, 0xbb, 0x9b, 0x50, 0x65, 0x0a, 0x1a, 0x14, 0xc6, 0xde,
	0xc1, 0xce, 0x8e, 0x82, 0x82, 0xea, 0xd8, 0x63, 0xb8, 0xd2, 0x73, 0x6c, 0x72, 0xc2, 0x55, 0xce,
	0xe8, 0xf8, 0x96, 0x63, 0x2b, 0xce, 0xef, 0x2e, 0xca, 0xa7, 0x96, 0x19, 0xea, 0x8b, 0x21, 0x26,
	0xf2, 0x85, 0xff, 0x1f, 0xde, 0x1c, 0xc1, 0x08, 0x6f, 0x57, 0xde, 0x43, 0x1c, 0x1b, 0x49, 0x1c,
	0x1f, 0x0d, 0xb7, 0x31, 0x1f, 0x2a, 0x5e, 0xb6, 0xeb, 0x1c, 0xca, 0xfa, 0xe0, 0xb2, 0x58, 0x14,
	0x89, 0x2b, 0x69, 0x9f, 0x38, 0x87, 0xc8, 0xf7, 0x3e, 0x71, 0x0e, 0xbd, 0x90, 0xa7, 0x66, 0x51,
	0xe2, 0xfd, 0x58, 0x4f, 0xcd, 0xc2, 0xc4, 0x7b, 0xb0, 0x88, 0x2d, 0xc0, 0x60, 0x5f, 0x58, 0x0d,
	0x7c, 0x48, 0x60, 0x7e, 0xa8, 0xf3, 0x83, 0x7d, 0x6e, 0x38, 0xe4, 0x6e, 0xf1, 0x83, 0xc8, 0x6e,
	0xf1, 0x43, 0x58, 0x36, 0x3a, 0xfe, 0x80, 0x25, 0x6f, 0xbd, 0x3e, 0xb1, 0x3d, 0x83, 0x2d, 0x8f,
	0x48, 0xba, 0x9e, 0xa0, 0xaf, 0xa3, 0x73, 0xbc, 0x2d, 0x84, 0xc6, 0x33, 0xb0, 0xc6, 0x5f, 0xe7,
	0xa0, 0x1e, 0xfa, 0xfc, 0xda, 0x8b, 0x50, 0x20, 0x84, 0x9a, 0x94, 0x1a, 0xcd, 0x08, 0x04, 0x4b,
	0x0f, 0x6a, 0x5d, 0xc2, 0xc1, 0xd2, 0x0b, 0xd1, 0xc5, 0xd4, 0x71, 0x3c, 0x9f, 0x25, 0xe3, 0x85,
	0xe1, 0x7b, 0x78, 0xbe, 0x76, 0x17, 0x16, 0x82, 0x2f, 0x89, 0x8a, 0x73, 0x78, 0x5f, 0x79, 0x5e,
	0x22, 0x0c, 0x6b, 0x74, 0xf7, 0xe1, 0x7c, 0x84, 0x2a, 0x52, 0x7e, 0x5b, 0x08, 0xd1, 0xc9, 0x9a,
	0x56, 0x84, 0x92, 0xbd, 0x63, 0x61, 0x04, 0x25, 0x7b, 0x87, 0x4d, 0x58, 0x1e, 0x52, 0xca, 0x85,
	0x08, 0xca, 0x0c, 0xb8, 0xc4, 0xbe, 0x44, 0x42, 0xe2, 0xf3, 0xb9, 0x28, 0x3b, 0x6c, 0xc1, 0xa5,
	0x80, 0x85, 0x0c, 0x75, 0x3b, 0x47, 0x86, 0x6d, 0x13, 0x26, 0xc6, 0xb8, 0xfc, 0xae, 0x4b, 0x44,
	0x11, 0xe1, 0x6e, 0x71, 0xb4, 0x5d, 0xb3, 0xe1, 0xc3, 0x42, 0x8c, 0x78, 0xd2, 0xe4, 0x9c, 0x89,
	0x54, 0xdb, 0x27, 0xc7, 0xbe, 0xda, 0xf3, 0xc3, 0xc6, 0x5f, 0x90, 0x63, 0xd6, 0x7d, 0xe0, 0x5b,
	0xe1, 0x6e, 0x32, 0x3a, 0xa2, 0xad, 0xc0, 0x4c, 0x8f, 0x78, 0x9e, 0x71, 0xa8, 0xee, 0xe9, 0xcb,
	0xc1, 0x86, 0x0f, 0x8b, 0xc8, 0x51, 0x4c, 0x7a, 0x54, 0xe9, 0xbd, 0x50, 0xcb, 0x2d, 0x3f, 0xb0,
	0xa4, 0x45, 0x3d, 0x90, 0xda, 0x71, 0xdb, 0xd8, 0x84, 0x72, 0xd0, 0x58, 0x15, 0xb4, 0xc2, 0x65,
	0x22, 0xad, 0x70, 0x23, 0xba, 0x3c, 0x1b, 0x7f, 0x9c, 0x81, 0x25, 0xb6, 0xfe, 0xdf, 0xb5, 0xfc,
	0xa3, 0x47, 0xae, 0xeb, 0xa0, 0x73, 0x56, 0x57, 0xa1, 0x48, 0xd8, 0x88, 0x9e, 0x61, 0x13, 0x2a,
	0xd3, 0xc9, 0x33, 0x9c, 0x96, 0x00, 0x0c, 0x4d, 0x78, 0x36, 0xd9, 0x84, 0xdf, 0x57, 0xcf, 0x09,
	0x08, 0x96, 0xb9, 0x30, 0x4b, 0x7c, 0x36, 0x80, 0x4f, 0xa4, 0xf1, 0x37, 0x33, 0xd0, 0x78, 0x4c,
	0xfc, 0xf0, 0xf7, 0x14, 0x62, 0xe7, 0xc9, 0xbd, 0xb0, 0x47, 0xa2, 0x9b, 0x87, 0xce, 0x23, 0x33,
	0x69, 0xb2, 0x1e, 0x90, 0xe2, 0x94, 0x5f, 0x9b, 0x3a, 0xe5, 0xdf, 0x00, 0x8d, 0x7a, 0xf5, 0xb6,
	0xda, 0x59, 0x83, 0x17, 0x7e, 0x8e, 0xc2, 0xb7, 0x71, 0x77, 0x8d, 0xa4, 0x51, 0x1b, 0x77, 0x4a,
	0x61, 0x9a, 0x3d, 0xdc, 0xbc, 0xc3, 0x8a, 0x53, 0xa1, 0xa7, 0x60, 0xc9, 0xac, 0xf9, 0x8e, 0xf2,
	0x0c, 0x8e, 0xaf, 0x3e, 0xa1, 0xac, 0xe2, 0xef, 0x85, 0x9a, 0x83, 0x94, 0x39, 0x29, 0x0d, 0x49,
	0xb3, 0x78, 0x3a, 0x01, 0x6a, 0xd0, 0x2f, 0x55, 0x0d, 0xa3, 0xca, 0x99, 0xf0, 0x5a, 0xde, 0x90,
	0x27, 0xee, 0x4d, 0x82, 0xe1, 0x04, 0x04, 0x5a, 0xc0, 0xaf, 0xae, 0xa2, 0x05, 0xdc, 0x2e, 0x07,
	0x45, 0x5c, 0xdc, 0xe1, 0x28, 0x4b, 0xb9, 0x6f, 0x85, 0x4a, 0xb9, 0x4a, 0x2b, 0x2e, 0x2e, 0xe8,
	0xc6, 0x56, 0xda, 0xe6, 0x47, 0x56, 0xda, 0x2e, 0x07, 0x35, 0x60, 0xdc, 0x31, 0x29, 0x2b, 0xc1,
	0xcb, 0xb2, 0x12, 0x8c, 0x6b, 0x95, 0xa2, 0x1e, 0x7c, 0x39, 0xa8, 0x07, 0xe3, 0x02, 0xa5, 0xac,
	0x0a, 0x07, 0x9a, 0x8a, 0x8b, 0x92, 0xa2, 0x1f, 0x3b, 0xae, 0x1c, 0x59, 0x49, 0x5d, 0x8e, 0x9c,
	0x4b, 0x2c, 0x47, 0xa6, 0xdb, 0x21, 0xbb, 0x15, 0x57, 0x98, 0x59, 0x5c, 0xcb, 0xdc, 0x28, 0x45,
	0x4b, 0x32, 0x0d, 0x03, 0xae, 0x8d, 0x54, 0x5b, 0x61, 0x5f, 0x1e, 0x40, 0x49, 0xf8, 0x18, 0x69,
	0x61, 0x56, 0x12, 0x02, 0x0c, 0x41, 0xda, 0x0a, 0xf0, 0x1b, 0x3f, 0x2a, 0xc0, 0x52, 0x02, 0x56,
	0xfa, 0xf2, 0x78, 0xa8, 0xad, 0x2b, 0x93, 0xd0, 0xd6, 0x75, 0x1d, 0x2a, 0xb8, 0x06, 0x87, 0xf7,
	0xbb, 0xc0, 0x1b, 0xd6, 0xde, 0xb0, 0x19, 0xca, 0x4f, 0x6f, 0x86, 0xe2, 0x52, 0xda, 0xc2, 0xc4,
	0x29, 0x6d, 0x52, 0x9a, 0x5d, 0x9c, 0x2a, 0xcd, 0x1e, 0xdf, 0x66, 0x5d, 0x89, 0x84, 0x60, 0x17,
	0x20, 0x47, 0x57, 0xab, 0x8a, 0x56, 0x8b, 0x0e, 0xfc, 0x72, 0x36, 0x9b, 0x12, 0x32, 0xec, 0xfa,
	0x34, 0x19, 0x76, 0x62, 0xf6, 0x3f, 0x37, 0x55, 0xf6, 0xdf, 0xf8, 0x71, 0x06, 0xce, 0xc7, 0x66,
	0x8d, 0xda, 0xbb, 0x30, 0xc3, 0x12, 0x4d, 0xc3, 0x17, 0xe1, 0xc0, 0x28, 0xde, 0x45, 0x8a, 0xba,
	0xc9, 0x42, 0x40, 0x9b, 0x1c, 0xfb, 0xe1, 0x84, 0x1a, 0x7f, 0x8a, 0x79, 0x8a, 0xa0, 0xe6, 0xd4,
	0xab, 0x50, 0x60, 0x8e, 0x56, 0x34, 0x73, 0x20, 0x3f, 0xcb, 0xc7, 0x1b, 0x3f, 0x98, 0x81, 0xb9,
	0xa7, 0xa1, 0x5c, 0x60, 0xba, 0xa3, 0xcd, 0xc3, 0x62, 0x68, 0x7e, 0x54, 0xf1, 0x55, 0x16, 0x43,
	0x83, 0x1e, 0xb2, 0x8e, 0x63, 0x32, 0xc1, 0xce, 0x29, 0x3d, 0x64, 0x5b, 0x8e, 0x39, 0xdc, 0x24,
	0x11, 0x07, 0xb8, 0x8b, 0x08, 0x8b, 0x6f, 0x92, 0xf0, 0x13, 0xe1, 0x81, 0xf4, 0xcd, 0x44, 0xa4,
	0x6f, 0x0d, 0x29, 0x63, 0x39, 0xd2, 0xe1, 0x4b, 0xf5, 0x2c, 0x7d, 0xfe, 0x54, 0x9b, 0x30, 0x7f,
	0xda, 0x83, 0xc6, 0x28, 0xce, 0x3e, 0x5f, 0xca, 0x79, 0xc4, 0x75, 0x25, 0x91, 0xab, 0xcf, 0x96,
	0x38, 0xbc, 0x6b, 0x54, 0x8f, 0x2c, 0x88, 0xb0, 0xe7, 0xea, 0x16, 0x10, 0x4c, 0xbf, 0x05, 0x54,
	0x99, 0x7a, 0x0b, 0x68, 0xf6, 0x14, 0xe7, 0x77, 0xab, 0x93, 0x9d, 0xdf, 0xfd, 0x14, 0xb4, 0x98,
	0xba, 0xd8, 0x78, 0xb5, 0x9d, 0x23, 0xe1, 0xa2, 0x18, 0x35, 0x70, 0x34, 0x47, 0xd5, 0x14, 0x03,
	0x47, 0x53, 0xd3, 0x5b, 0x50, 0xb3, 0xec, 0x57, 0x0e, 0x4b, 0xa1, 0x06, 0xbd, 0x7d, 0xe2, 0x2a,
	0xbb, 0x1d, 0x55, 0x01, 0x7b, 0xc6, 0x40, 0xda, 0x5d, 0xea, 0x4e, 0xed, 0x97, 0x6d, 0xa3, 0xc3,
	0xcb, 0x00, 0x8b, 0xb8, 0x0b, 0xfa, 0xa1, 0x61, 0xbf, 0xdc, 0xe4, 0x00, 0xea, 0x5b, 0x83, 0x1f,
	0x8d, 0x1f, 0x66, 0xe0, 0x62, 0x58, 0x15, 0x27, 0x4d, 0x21, 0x76, 0x40, 0x8b, 0xc8, 0x9c, 0x4c,
	0x24, 0x96, 0x84, 0x57, 0x0d, 0x3f, 0xa3, 0x35, 0x1f, 0x16, 0x3b, 0xaf, 0xf1, 0x27, 0x33, 0xf0,
	0xc6, 0xd3, 0x31, 0x22, 0xce, 0x32, 0xe1, 0x78, 0x53, 0x11, 0xe9, 0x80, 0xce, 0x26, 0x76, 0x40,
	0xbf, 0x03, 0xf3, 0x43, 0x54, 0xb9, 0x13, 0x81, 0x43, 0xd8, 0xb9, 0x00, 0x5d, 0x6e, 0x43, 0xac,
	0x43, 0x30, 0xa6, 0x74, 0x9b, 0x05, 0x47, 0xf9, 0x24, 0x54, 0x56, 0xa0, 0x36, 0x40, 0x0b, 0x08,
	0xe2, 0x8d, 0x4e, 0xc0, 0xf0, 0x85, 0x34, 0x3e, 0x38, 0x0b, 0xe5, 0x34, 0xca, 0x36, 0x06, 0x56,
	0x3d, 0x5d, 0x21, 0xc6, 0x5b, 0x1a, 0x77, 0xc5, 0x8e, 0x8d, 0x7a, 0x90, 0x50, 0x31, 0x64, 0x6c,
	0x83, 0x66, 0x07, 0x9f, 0x28, 0xd4, 0xbe, 0x09, 0x4b, 0x31, 0x54, 0x6d, 0xd7, 0x78, 0xad, 0x58,
	0xb8, 0xc5, 0x08, 0x65, 0xcb, 0x78, 0x1d, 0x88, 0x70, 0x29, 0x22, 0xc2, 0xe9, 0x2d, 0x5d, 0x79,
	0x42, 0x4b, 0x77, 0x1b, 0x66, 0xad, 0x1e, 0x0d, 0xa9, 0x78, 0x56, 0x27, 0x0c, 0x0e, 0x72, 0x36,
	0x15, 0x0e, 0x66, 0x3f, 0x7e, 0x4d, 0xfb, 0xd3, 0x9f, 0xa1, 0x12, 0xca, 0x44, 0x57, 0x23, 0x04,
	0xe2, 0x39, 0xdc, 0x02, 0xfa, 0x14, 0xc9, 0xd3, 0x24, 0x9b, 0xde, 0x43, 0x51, 0x96, 0x96, 0xef,
	0x2e, 0x54, 0xd0, 0x57, 0x16, 0x36, 0x2b, 0x2e, 0xe3, 0xc7, 0x68, 0x8d, 0x7f, 0xc9, 0xc3, 0xda,
	0x38, 0xed, 0x4c, 0xd0, 0x4c, 0xe9, 0x2a, 0xb3, 0x11, 0x57, 0xa9, 0x38, 0xe4, 0x5c, 0x2a, 0x87,
	0x9c, 0x4f, 0x72, 0xc8, 0xc3, 0x58, 0xa0, 0x90, 0x26, 0x16, 0xc0, 0x5e, 0xba, 0x18, 0xeb, 0xa5,
	0x37, 0x65, 0x97, 0xc8, 0x8c, 0xd2, 0xf1, 0x9b, 0xc6, 0x44, 0xc9, 0x06, 0x12, 0x55, 0xec, 0x4a,
	0xd3, 0x8b, 0x5d, 0x79, 0x12, 0xb1, 0xdb, 0x46, 0xb6, 0xaa, 0x6f, 0x58, 0x29, 0xfd, 0x71, 0x6d,
	0x58, 0xd7, 0xb2, 0xb0, 0x5f, 0xaa, 0xa4, 0xf0, 0x4b, 0xb3, 0xe9, 0xfd, 0x52, 0x35, 0x95, 0x5f,
	0xfa, 0xcb, 0x0c, 0xdc, 0x1c, 0xb7, 0xd0, 0x93, 0xfa, 0xa9, 0xef, 0x8c, 0xf0, 0x53, 0x6f, 0xa5,
	0xfc, 0xb8, 0x71, 0x7e, 0xeb, 0xf7, 0x0b, 0x70, 0xf9, 0xe9, 0x88, 0x20, 0xea, 0xd7, 0xaa, 0x15,
	0xe1, 0xf0, 0xad, 0x90, 0x14, 0xbe, 0x29, 0xe7, 0x2b, 0x8a, 0x91, 0x87, 0x52, 0xcf, 0x32, 0xd4,
	0xb0, 0x99, 0x34, 0x1a, 0x16, 0x1f, 0x17, 0x94, 0x26, 0x8d, 0x0b, 0x42, 0x4a, 0x54, 0x9e, 0x5e,
	0x89, 0x60, 0x12, 0x25, 0x0a, 0xc7, 0x87, 0x95, 0xc9, 0xe2, 0x43, 0xa9, 0x3d, 0xb3, 0x29, 0xb4,
	0xa7, 0x9a, 0x5e, 0x7b, 0x6a, 0xa9, 0xb4, 0xe7, 0xef, 0x33, 0x70, 0x7d, 0x94, 0x44, 0x4e, 0xaa,
	0x39, 0x5d, 0x58, 0x1b, 0x93, 0x55, 0x48, 0x3d, 0xba, 0x36, 0x4e, 0x8f, 0x7c, 0xa7, 0xdf, 0xba,
	0x32, 0x2a, 0xe5, 0xf0, 0x1a, 0xff, 0x96, 0x81, 0xc5, 0x5d, 0xe6, 0xbb, 0x43, 0x97, 0x6f, 0xdd,
	0x81, 0xbc, 0x69, 0xf8, 0x86, 0x98, 0xea, 0xe5, 0x26, 0xbe, 0x25, 0xed, 0xf9, 0xf0, 0xef, 0x6d,
	0xc3, 0x37, 0x5a, 0x0c, 0x33, 0xdd, 0x45, 0x5c, 0x5a, 0x13, 0xaa, 0x38, 0x94, 0x88, 0x29, 0x10,
	0xcf, 0xa2, 0x58, 0xc2, 0xd3, 0xde, 0x86, 0x4a, 0x87, 0x74, 0xbb, 0x12, 0x3b, 0x1f, 0xc6, 0x06,
	0x0a, 0x15, 0xb8, 0xec, 0xf0, 0x32, 0xe3, 0x6d, 0x85, 0x52, 0x4e, 0x3e, 0xbc, 0x6b, 0x36, 0x7e,
	0x91, 0x87, 0xf9, 0xbd, 0xc1, 0x7e, 0xd7, 0xea, 0xe0, 0x46, 0xcb, 0x78, 0xa3, 0x11, 0xed, 0x1d,
	0xca, 0x4e, 0xd5, 0x3b, 0x94, 0x4b, 0xd3, 0x3b, 0x94, 0x54, 0xc8, 0xc9, 0x4f, 0x55, 0xc8, 0xf9,
	0x08, 0x66, 0x95, 0xb6, 0x90, 0xf1, 0x85, 0xa5, 0x8a, 0x89, 0xba, 0x41, 0xfe, 0x17, 0x97, 0x5d,
	0xa2, 0x0d, 0x6c, 0xc5, 0xc4, 0x06, 0x36, 0xdc, 0x00, 0x5b, 0x8a, 0x6b, 0x80, 0xbd, 0x03, 0xf3,
	0x78, 0x99, 0xf8, 0xa6, 0x11, 0x2e, 0x23, 0xd4, 0xd1, 0xaa, 0xb0, 0xad, 0xa3, 0xbb, 0xb0, 0xa0,
	0xca, 0x04, 0xa7, 0xc1, 0x67, 0x80, 0xe7, 0x15, 0x39, 0xa0, 0x54, 0x8d, 0xdf, 0xa9, 0x81, 0x46,
	0xf5, 0xce, 0x76, 0x5e, 0x8f, 0x17, 0xbb, 0x31, 0x65, 0x18, 0xb5, 0x2b, 0x37, 0x17, 0xdf, 0x95,
	0x1b, 0x6d, 0x25, 0xc8, 0xa7, 0x6e, 0x25, 0xf8, 0x08, 0xea, 0xc1, 0x67, 0xea, 0x3b, 0x96, 0xed,
	0xcb, 0xd6, 0x5d, 0x79, 0xd2, 0x44, 0x7e, 0x8c, 0x3d, 0x0a, 0x6c, 0xd5, 0x4c, 0xfc, 0xd3, 0x63,
	0xa7, 0x9f, 0x45, 0xdd, 0xaf, 0xa8, 0x9c, 0x7e, 0x1e, 0xd1, 0x2c, 0x19, 0x29, 0xee, 0x4c, 0xd4,
	0x4e, 0x5f, 0x9a, 0xa2, 0x59, 0x92, 0x95, 0x39, 0x2f, 0x8f, 0x6b, 0xcf, 0xdc, 0x81, 0xcb, 0x11,
	0x4a, 0x93, 0x78, 0x1d, 0xd7, 0xea, 0x53, 0x53, 0xaa, 0x34, 0xa5, 0x2d, 0x87, 0x18, 0x6c, 0x0f,
	0xf1, 0xb4, 0x0f, 0xa1, 0xc2, 0xf7, 0x1b, 0xda, 0x96, 0x7d, 0xe0, 0x08, 0xb7, 0x3a, 0x2f, 0x96,
	0x93, 0xef, 0x35, 0x50, 0x1f, 0xf0, 0xb0, 0x44, 0x39, 0x7d, 0xf9, 0xd3, 0xd5, 0x4c, 0x0b, 0x5e,
	0x07, 0xa3, 0xda, 0x3d, 0x00, 0x56, 0xb5, 0xe7, 0xc4, 0xdc, 0xb1, 0xca, 0x36, 0x6a, 0x56, 0xb5,
	0x67, 0xb4, 0xc1, 0x75, 0x14, 0x62, 0x20, 0xda, 0xbe, 0x59, 0x49, 0x6c, 0xdf, 0xdc, 0x81, 0x79,
	0x97, 0x6f, 0xb2, 0x21, 0x6d, 0x1f, 0x9f, 0x7d, 0xd5, 0x05, 0x51, 0xa0, 0xeb, 0xa7, 0xb8, 0x95,
	0x4e, 0x8d, 0x1e, 0x6a, 0x93, 0x44, 0x0f, 0x43, 0x4b, 0x5b, 0x9f, 0xb2, 0x4b, 0x73, 0x6e, 0x92,
	0x2e, 0xcd, 0x5b, 0x11, 0xb7, 0x80, 0x5b, 0x8c, 0x43, 0x5e, 0x20, 0x7a, 0x53, 0x9b, 0x36, 0xc1,
	0x4d, 0x6d, 0x97, 0xa0, 0x2c, 0x0d, 0x98, 0xa7, 0x2f, 0xb0, 0xcb, 0x89, 0x4a, 0xc2, 0x76, 0x79,
	0xd4, 0x48, 0xc7, 0x35, 0x3d, 0x2e, 0x8e, 0x37, 0xd2, 0x5e, 0xa4, 0xe7, 0x31, 0x62, 0x54, 0xcf,
	0x27, 0x1a, 0xd5, 0x5f, 0x52, 0x93, 0x6a, 0x72, 0xbf, 0xdb, 0xd2, 0x74, 0xfd, 0x6e, 0x78, 0x39,
	0xa8, 0x10, 0xb3, 0x6e, 0x26, 0x7e, 0x93, 0x4d, 0xca, 0xe5, 0xd8, 0xe3, 0x54, 0x9b, 0xbe, 0xb6,
	0x07, 0x17, 0xc2, 0xed, 0x95, 0x82, 0xdd, 0xc5, 0xf1, 0xfd, 0x95, 0x9e, 0xda, 0x5e, 0xc9, 0x39,
	0x3e, 0x43, 0x36, 0x48, 0xf1, 0xcc, 0xcb, 0xe3, 0xdf, 0xd6, 0x8b, 0xe9, 0xd7, 0xc4, 0xfc, 0x94,
	0x9e, 0xd2, 0x4b, 0xe9, 0xf9, 0xe1, 0x96, 0x52, 0xc5, 0xb6, 0x1e, 0x19, 0x74, 0x72, 0x5d, 0xcb,
	0x7e, 0xc9, 0xba, 0x79, 0xcb, 0x11, 0xdb, 0xca, 0x10, 0x9e, 0x58, 0xf6, 0xcb, 0x68, 0x77, 0xfe,
	0x4a, 0x52, 0x77, 0x7e, 0xe3, 0x0f, 0x33, 0x70, 0xe5, 0x31, 0xf1, 0xa3, 0x7e, 0x30, 0xd8, 0xe0,
	0x6f, 0x84, 0x02, 0x64, 0x10, 0x01, 0xb2, 0x65, 0x1f, 0xa2, 0xeb, 0x5e, 0x67, 0x0e, 0xac, 0xae,
	0x3f, 0x0c, 0x33, 0x19, 0xd2, 0x0e, 0x1b, 0x6a, 0x49, 0x90, 0x76, 0x03, 0x0a, 0x3d, 0xeb, 0x98,
	0x98, 0x62, 0x5b, 0x44, 0xe3, 0x31, 0xfc, 0x53, 0x3a, 0x24, 0x83, 0x78, 0x8e, 0xd0, 0xf8, 0x61,
	0x06, 0x16, 0x62, 0xa6, 0x94, 0x32, 0x58, 0x7f, 0xc2, 0x97, 0xcd, 0x76, 0x5e, 0xb7, 0x63, 0x3a,
	0x3b, 0x2e, 0xa2, 0xe6, 0x32, 0x95, 0x3f, 0xff, 0x08, 0xa1, 0x67, 0x36, 0xfe, 0x35, 0x0b, 0x55,
	0xc5, 0xb9, 0xc6, 0xb6, 0xe6, 0x65, 0x26, 0x68, 0xcd, 0x0b, 0x0e, 0xdf, 0x64, 0x47, 0x1f, 0xbe,
	0x89, 0xf8, 0x8a, 0xdc, 0xa8, 0x3b, 0xe8, 0x82, 0xa0, 0x2a, 0x1f, 0x17, 0x54, 0x85, 0x9c, 0x5d,
	0x61, 0x32, 0x67, 0xf7, 0x81, 0xe2, 0xec, 0x8a, 0x09, 0xce, 0x6e, 0x48, 0x8b, 0x1c, 0xde, 0x84,
	0xc7, 0xcc, 0x1a, 0xff, 0x98, 0x01, 0x18, 0xce, 0x27, 0xd2, 0x48, 0x90, 0x99, 0xa8, 0x91, 0x20,
	0x9b, 0xb2, 0x91, 0x20, 0x37, 0xaa, 0x91, 0x20, 0x3f, 0xaa, 0x91, 0xa0, 0x10, 0x6d, 0x24, 0x68,
	0xfc, 0x5e, 0x06, 0xca, 0xc1, 0xaa, 0x44, 0x76, 0xf9, 0x33, 0x49, 0xbb, 0xfc, 0x6a, 0xcf, 0x40,
	0x36, 0xbe, 0x67, 0x20, 0xf6, 0x8c, 0x46, 0x6e, 0xd4, 0x19, 0x8d, 0xc6, 0x7f, 0x65, 0x61, 0x95,
	0x3b, 0x95, 0x18, 0x91, 0x17, 0x4a, 0xae, 0xf8, 0xb4, 0x4c, 0xc8, 0xa7, 0xa1, 0xe0, 0x31, 0x3b,
	0x51, 0xf0, 0x98, 0x9b, 0x36, 0x78, 0xcc, 0x8f, 0x09, 0x1e, 0x23, 0x9a, 0x51, 0x98, 0x2c, 0x8a,
	0x2a, 0x4e, 0x1e, 0x45, 0x45, 0x83, 0xf7, 0x99, 0xb4, 0xc1, 0x7b, 0xe3, 0x8f, 0x72, 0xb0, 0xca,
	0xcf, 0x1e, 0x24, 0x7f, 0x01, 0x9e, 0x76, 0x94, 0x42, 0x69, 0xc7, 0xd9, 0x77, 0xf9, 0x55, 0x7e,
	0x97, 0x7d, 0x58, 0xe5, 0xee, 0x76, 0xdc, 0x67, 0x89, 0xd9, 0x69, 0x53, 0x1d, 0x6c, 0x36, 0xd1,
	0xc1, 0xfe, 0x5d, 0x06, 0x2e, 0x0e, 0x1d, 0x6c, 0xb8, 0x7b, 0x6e, 0xe4, 0xf7, 0x8d, 0xbe, 0x5a,
	0x36, 0x75, 0xbe, 0xb8, 0x1d, 0xcd, 0x17, 0x79, 0x4d, 0xe7, 0x52, 0x6c, 0xbe, 0xc8, 0xa7, 0x13,
	0x4e, 0x1b, 0x1b, 0xbf, 0x0d, 0x8b, 0x71, 0x78, 0xa7, 0xf6, 0x80, 0x69, 0x4c, 0x5d, 0x63, 0x0f,
	0x96, 0xe3, 0x16, 0x4e, 0x94, 0xc2, 0x36, 0x22, 0xfd, 0x4b, 0x4a, 0xd7, 0xf8, 0x90, 0x04, 0xf5,
	0x2d, 0xfd, 0x79, 0x46, 0x34, 0x15, 0x0f, 0xa1, 0x58, 0x89, 0x32, 0x71, 0x4a, 0x34, 0xcd, 0x25,
	0x87, 0xa2, 0xfb, 0x26, 0x1f, 0xee, 0xbe, 0x79, 0x13, 0x2a, 0x38, 0x51, 0x55, 0xae, 0x9d, 0x42,
	0x80, 0xc6, 0x5f, 0xad, 0xc2, 0x2c, 0xbe, 0xeb, 0x95, 0xf5, 0xae, 0xf0, 0x75, 0x60, 0xf7, 0x91,
	0x8e, 0x0c, 0x62, 0x24, 0xa6, 0xd6, 0xe4, 0xb7, 0x52, 0xb0, 0x2d, 0xae, 0x67, 0x89, 0x5b, 0x5c,
	0x01, 0xce, 0xd9, 0x79, 0xe5, 0xb3, 0xf3, 0xca, 0x67, 0xe7, 0x95, 0xcf, 0xce, 0x2b, 0x9f, 0x9d,
	0x57, 0x3e, 0x3b, 0xaf, 0x7c, 0x76, 0x5e, 0xf9, 0xec, 0xbc, 0xf2, 0xd9, 0x79, 0xe5, 0xb3, 0xf3,
	0xca, 0x67, 0xe7, 0x95, 0xcf, 0xce, 0x2b, 0x9f, 0x9d, 0x57, 0x3e, 0x3b, 0xaf, 0x7c, 0x76, 0x5e,
	0x79, 0xdc, 0x79, 0xe5, 0xff, 0x2c, 0xc2, 0xc5, 0x17, 0xae, 0x61, 0x06, 0x41, 0x80, 0x72, 0x61,
	0x28, 0xbe, 0x1a, 0x34, 0x33, 0xf5, 0xd5, 0xa0, 0xd9, 0xd3, 0x5d, 0x0d, 0x9a, 0x3b, 0xdd, 0xd5,
	0xa0, 0xf9, 0x09, 0x94, 0xf9, 0xb6, 0x9a, 0xae, 0x4f, 0x72, 0x65, 0x73, 0x71, 0xaa, 0x2b, 0x9b,
	0x67, 0x52, 0x5e, 0xd9, 0x5c, 0x4a, 0x3e, 0x62, 0x17, 0xbe, 0x5e, 0xb3, 0x9c, 0xf2, 0x26, 0x66,
	0x18, 0x7f, 0x13, 0x73, 0x65, 0xba, 0x9b, 0x98, 0x67, 0xd3, 0xdc, 0xc4, 0x5c, 0x4d, 0x77, 0x13,
	0x73, 0x2d, 0xfe, 0x9a, 0xd5, 0xe8, 0x05, 0xe7, 0xf5, 0xe4, 0x0b, 0xce, 0xbf, 0x1b, 0xbe, 0x08,
	0x78, 0x8e, 0xbd, 0xd7, 0x86, 0x78, 0xaf, 0x44, 0xf1, 0x1f, 0x77, 0x1f, 0xf0, 0xa9, 0xef, 0xe9,
	0x7d, 0xf8, 0xed, 0x9f, 0x7c, 0xb5, 0x92, 0xf9, 0xf2, 0xab, 0x95, 0xcc, 0xcf, 0xbe, 0x5a, 0x39,
	0xf7, 0xef, 0x5f, 0xad, 0x9c, 0xfb, 0xc1, 0xd7, 0x2b, 0xe7, 0xfe, 0xe2, 0xeb, 0x95, 0x73, 0x7f,
	0xfb, 0xf5, 0xca, 0xb9, 0x9f, 0x7c, 0xbd, 0x72, 0xee, 0xcb, 0xaf, 0x57, 0xce, 0xfd, 0xec, 0xeb,
	0x95, 0x73, 0x7f, 0xf0, 0xf3, 0x95, 0x73, 0x3f, 0xfa, 0xf9, 0xca, 0xb9, 0xdf, 0xb8, 0xc8, 0xfc,
	0xcb, 0x2b, 0x7b, 0xdd, 0xe8, 0x5b, 0xeb, 0xfd, 0xfd, 0xf5, 0xe1, 0xff, 0x07, 0xf3, 0xbf, 0x03,
	0x00, 0x00, 0xff, 0xff, 0xc2, 0x5b, 0x1c, 0x8c, 0xfa, 0x74, 0x00, 0x00,
}
