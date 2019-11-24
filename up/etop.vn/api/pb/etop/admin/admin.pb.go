// source: etop/admin/admin.proto

package admin

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	credit_type "etop.vn/api/pb/etop/etc/credit_type"
	notifier_entity "etop.vn/api/pb/etop/etc/notifier_entity"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	status3 "etop.vn/api/pb/etop/etc/status3"
	_ "etop.vn/api/pb/etop/order"
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

type GetOrdersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetOrdersRequest) Reset()         { *m = GetOrdersRequest{} }
func (m *GetOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*GetOrdersRequest) ProtoMessage()    {}
func (*GetOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{0}
}

var xxx_messageInfo_GetOrdersRequest proto.InternalMessageInfo

type GetFulfillmentsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	ShopId               dot.ID           `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	OrderId              dot.ID           `protobuf:"varint,4,opt,name=order_id,json=orderId" json:"order_id"`
	Status               *status3.Status  `protobuf:"varint,5,opt,name=status,enum=status3.Status" json:"status"`
	Filters              []*common.Filter `protobuf:"bytes,6,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetFulfillmentsRequest) Reset()         { *m = GetFulfillmentsRequest{} }
func (m *GetFulfillmentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetFulfillmentsRequest) ProtoMessage()    {}
func (*GetFulfillmentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{1}
}

var xxx_messageInfo_GetFulfillmentsRequest proto.InternalMessageInfo

type LoginAsAccountRequest struct {
	UserId    dot.ID `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id"`
	AccountId dot.ID `protobuf:"varint,2,opt,name=account_id,json=accountId" json:"account_id"`
	// This is a sensitive API, so admin must provide password before processing!
	Password             string   `protobuf:"bytes,3,opt,name=password" json:"password"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginAsAccountRequest) Reset()         { *m = LoginAsAccountRequest{} }
func (m *LoginAsAccountRequest) String() string { return proto.CompactTextString(m) }
func (*LoginAsAccountRequest) ProtoMessage()    {}
func (*LoginAsAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{2}
}

var xxx_messageInfo_LoginAsAccountRequest proto.InternalMessageInfo

type GetMoneyTransactionsRequest struct {
	Ids                                []dot.ID         `protobuf:"varint,1,rep,name=ids" json:"ids"`
	ShopId                             dot.ID           `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	MoneyTransactionShippingExternalId dot.ID           `protobuf:"varint,3,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	Paging                             *common.Paging   `protobuf:"bytes,4,opt,name=paging" json:"paging"`
	Filters                            []*common.Filter `protobuf:"bytes,5,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral               struct{}         `json:"-"`
	XXX_sizecache                      int32            `json:"-"`
}

func (m *GetMoneyTransactionsRequest) Reset()         { *m = GetMoneyTransactionsRequest{} }
func (m *GetMoneyTransactionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetMoneyTransactionsRequest) ProtoMessage()    {}
func (*GetMoneyTransactionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{3}
}

var xxx_messageInfo_GetMoneyTransactionsRequest proto.InternalMessageInfo

type RemoveFfmsMoneyTransactionRequest struct {
	FulfillmentIds       []dot.ID `protobuf:"varint,1,rep,name=fulfillment_ids,json=fulfillmentIds" json:"fulfillment_ids"`
	MoneyTransactionId   dot.ID   `protobuf:"varint,2,opt,name=money_transaction_id,json=moneyTransactionId" json:"money_transaction_id"`
	ShopId               dot.ID   `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveFfmsMoneyTransactionRequest) Reset()         { *m = RemoveFfmsMoneyTransactionRequest{} }
func (m *RemoveFfmsMoneyTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveFfmsMoneyTransactionRequest) ProtoMessage()    {}
func (*RemoveFfmsMoneyTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{4}
}

var xxx_messageInfo_RemoveFfmsMoneyTransactionRequest proto.InternalMessageInfo

type GetMoneyTransactionShippingExternalsRequest struct {
	Ids                  []dot.ID         `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Paging               *common.Paging   `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,3,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetMoneyTransactionShippingExternalsRequest) Reset() {
	*m = GetMoneyTransactionShippingExternalsRequest{}
}
func (m *GetMoneyTransactionShippingExternalsRequest) String() string {
	return proto.CompactTextString(m)
}
func (*GetMoneyTransactionShippingExternalsRequest) ProtoMessage() {}
func (*GetMoneyTransactionShippingExternalsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{5}
}

var xxx_messageInfo_GetMoneyTransactionShippingExternalsRequest proto.InternalMessageInfo

type RemoveMoneyTransactionShippingExternalLinesRequest struct {
	LineIds                            []dot.ID `protobuf:"varint,1,rep,name=line_ids,json=lineIds" json:"line_ids"`
	MoneyTransactionShippingExternalId dot.ID   `protobuf:"varint,2,opt,name=money_transaction_shipping_external_id,json=moneyTransactionShippingExternalId" json:"money_transaction_shipping_external_id"`
	XXX_NoUnkeyedLiteral               struct{} `json:"-"`
	XXX_sizecache                      int32    `json:"-"`
}

func (m *RemoveMoneyTransactionShippingExternalLinesRequest) Reset() {
	*m = RemoveMoneyTransactionShippingExternalLinesRequest{}
}
func (m *RemoveMoneyTransactionShippingExternalLinesRequest) String() string {
	return proto.CompactTextString(m)
}
func (*RemoveMoneyTransactionShippingExternalLinesRequest) ProtoMessage() {}
func (*RemoveMoneyTransactionShippingExternalLinesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{6}
}

var xxx_messageInfo_RemoveMoneyTransactionShippingExternalLinesRequest proto.InternalMessageInfo

type ConfirmMoneyTransactionRequest struct {
	MoneyTransactionId   dot.ID   `protobuf:"varint,1,opt,name=money_transaction_id,json=moneyTransactionId" json:"money_transaction_id"`
	ShopId               dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	TotalCod             dot.ID   `protobuf:"varint,3,opt,name=total_cod,json=totalCod" json:"total_cod"`
	TotalAmount          dot.ID   `protobuf:"varint,4,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	TotalOrders          dot.ID   `protobuf:"varint,5,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmMoneyTransactionRequest) Reset()         { *m = ConfirmMoneyTransactionRequest{} }
func (m *ConfirmMoneyTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmMoneyTransactionRequest) ProtoMessage()    {}
func (*ConfirmMoneyTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{7}
}

var xxx_messageInfo_ConfirmMoneyTransactionRequest proto.InternalMessageInfo

type DeleteMoneyTransactionRequest struct {
	MoneyTransactionId   dot.ID   `protobuf:"varint,1,opt,name=money_transaction_id,json=moneyTransactionId" json:"money_transaction_id"`
	ShopId               dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteMoneyTransactionRequest) Reset()         { *m = DeleteMoneyTransactionRequest{} }
func (m *DeleteMoneyTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteMoneyTransactionRequest) ProtoMessage()    {}
func (*DeleteMoneyTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{8}
}

var xxx_messageInfo_DeleteMoneyTransactionRequest proto.InternalMessageInfo

type GetShopsRequest struct {
	Paging               *common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetShopsRequest) Reset()         { *m = GetShopsRequest{} }
func (m *GetShopsRequest) String() string { return proto.CompactTextString(m) }
func (*GetShopsRequest) ProtoMessage()    {}
func (*GetShopsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{9}
}

var xxx_messageInfo_GetShopsRequest proto.InternalMessageInfo

type GetShopsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Shops                []*etop.Shop     `protobuf:"bytes,2,rep,name=shops" json:"shops"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetShopsResponse) Reset()         { *m = GetShopsResponse{} }
func (m *GetShopsResponse) String() string { return proto.CompactTextString(m) }
func (*GetShopsResponse) ProtoMessage()    {}
func (*GetShopsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{10}
}

var xxx_messageInfo_GetShopsResponse proto.InternalMessageInfo

type GetCreditRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCreditRequest) Reset()         { *m = GetCreditRequest{} }
func (m *GetCreditRequest) String() string { return proto.CompactTextString(m) }
func (*GetCreditRequest) ProtoMessage()    {}
func (*GetCreditRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{11}
}

var xxx_messageInfo_GetCreditRequest proto.InternalMessageInfo

type GetCreditsRequest struct {
	ShopId               dot.ID         `protobuf:"varint,1,opt,name=shop_id,json=shopId" json:"shop_id"`
	Paging               *common.Paging `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetCreditsRequest) Reset()         { *m = GetCreditsRequest{} }
func (m *GetCreditsRequest) String() string { return proto.CompactTextString(m) }
func (*GetCreditsRequest) ProtoMessage()    {}
func (*GetCreditsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{12}
}

var xxx_messageInfo_GetCreditsRequest proto.InternalMessageInfo

type CreateCreditRequest struct {
	Amount               dot.ID                  `protobuf:"varint,1,opt,name=amount" json:"amount"`
	ShopId               dot.ID                  `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Type                 *credit_type.CreditType `protobuf:"varint,3,opt,name=type,enum=credit_type.CreditType" json:"type"`
	PaidAt               dot.Time                `protobuf:"bytes,4,opt,name=paid_at,json=paidAt" json:"paid_at"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *CreateCreditRequest) Reset()         { *m = CreateCreditRequest{} }
func (m *CreateCreditRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCreditRequest) ProtoMessage()    {}
func (*CreateCreditRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{13}
}

var xxx_messageInfo_CreateCreditRequest proto.InternalMessageInfo

type UpdateCreditRequest struct {
	Id                   dot.ID                  `protobuf:"varint,1,opt,name=id" json:"id"`
	Amount               dot.ID                  `protobuf:"varint,2,opt,name=amount" json:"amount"`
	ShopId               dot.ID                  `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	Type                 *credit_type.CreditType `protobuf:"varint,4,opt,name=type,enum=credit_type.CreditType" json:"type"`
	PaidAt               dot.Time                `protobuf:"bytes,5,opt,name=paid_at,json=paidAt" json:"paid_at"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *UpdateCreditRequest) Reset()         { *m = UpdateCreditRequest{} }
func (m *UpdateCreditRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCreditRequest) ProtoMessage()    {}
func (*UpdateCreditRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{14}
}

var xxx_messageInfo_UpdateCreditRequest proto.InternalMessageInfo

type ConfirmCreditRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmCreditRequest) Reset()         { *m = ConfirmCreditRequest{} }
func (m *ConfirmCreditRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmCreditRequest) ProtoMessage()    {}
func (*ConfirmCreditRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{15}
}

var xxx_messageInfo_ConfirmCreditRequest proto.InternalMessageInfo

type CreatePartnerRequest struct {
	Partner              etop.Partner `protobuf:"bytes,1,opt,name=partner" json:"partner"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CreatePartnerRequest) Reset()         { *m = CreatePartnerRequest{} }
func (m *CreatePartnerRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePartnerRequest) ProtoMessage()    {}
func (*CreatePartnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{16}
}

var xxx_messageInfo_CreatePartnerRequest proto.InternalMessageInfo

type UpdateFulfillmentRequest struct {
	Id                       dot.ID          `protobuf:"varint,1,opt,name=id" json:"id"`
	FullName                 string          `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                    string          `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	TotalCodAmount           *int32          `protobuf:"varint,5,opt,name=total_cod_amount,json=totalCodAmount" json:"total_cod_amount"`
	IsPartialDelivery        bool            `protobuf:"varint,6,opt,name=is_partial_delivery,json=isPartialDelivery" json:"is_partial_delivery"`
	AdminNote                string          `protobuf:"bytes,7,opt,name=admin_note,json=adminNote" json:"admin_note"`
	ActualCompensationAmount int32           `protobuf:"varint,8,opt,name=actual_compensation_amount,json=actualCompensationAmount" json:"actual_compensation_amount"`
	ShippingState            *shipping.State `protobuf:"varint,9,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	XXX_NoUnkeyedLiteral     struct{}        `json:"-"`
	XXX_sizecache            int32           `json:"-"`
}

func (m *UpdateFulfillmentRequest) Reset()         { *m = UpdateFulfillmentRequest{} }
func (m *UpdateFulfillmentRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateFulfillmentRequest) ProtoMessage()    {}
func (*UpdateFulfillmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{17}
}

var xxx_messageInfo_UpdateFulfillmentRequest proto.InternalMessageInfo

type GenerateAPIKeyRequest struct {
	AccountId            dot.ID   `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateAPIKeyRequest) Reset()         { *m = GenerateAPIKeyRequest{} }
func (m *GenerateAPIKeyRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateAPIKeyRequest) ProtoMessage()    {}
func (*GenerateAPIKeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{18}
}

var xxx_messageInfo_GenerateAPIKeyRequest proto.InternalMessageInfo

type GenerateAPIKeyResponse struct {
	AccountId            dot.ID   `protobuf:"varint,1,opt,name=account_id,json=accountId" json:"account_id"`
	ApiKey               string   `protobuf:"bytes,2,opt,name=api_key,json=apiKey" json:"api_key"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateAPIKeyResponse) Reset()         { *m = GenerateAPIKeyResponse{} }
func (m *GenerateAPIKeyResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateAPIKeyResponse) ProtoMessage()    {}
func (*GenerateAPIKeyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{19}
}

var xxx_messageInfo_GenerateAPIKeyResponse proto.InternalMessageInfo

type UpdateMoneyTransactionShippingEtopRequest struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Adds                 []dot.ID          `protobuf:"varint,2,rep,name=adds" json:"adds"`
	Deletes              []dot.ID          `protobuf:"varint,3,rep,name=deletes" json:"deletes"`
	ReplaceAll           []dot.ID          `protobuf:"varint,4,rep,name=replace_all,json=replaceAll" json:"replace_all"`
	Note                 string            `protobuf:"bytes,5,opt,name=note" json:"note"`
	InvoiceNumber        string            `protobuf:"bytes,6,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount          *etop.BankAccount `protobuf:"bytes,7,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateMoneyTransactionShippingEtopRequest) Reset() {
	*m = UpdateMoneyTransactionShippingEtopRequest{}
}
func (m *UpdateMoneyTransactionShippingEtopRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateMoneyTransactionShippingEtopRequest) ProtoMessage()    {}
func (*UpdateMoneyTransactionShippingEtopRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{20}
}

var xxx_messageInfo_UpdateMoneyTransactionShippingEtopRequest proto.InternalMessageInfo

type GetMoneyTransactionShippingEtopsRequest struct {
	Ids                  []dot.ID         `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Status               *status3.Status  `protobuf:"varint,2,opt,name=status,enum=status3.Status" json:"status"`
	Paging               *common.Paging   `protobuf:"bytes,3,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetMoneyTransactionShippingEtopsRequest) Reset() {
	*m = GetMoneyTransactionShippingEtopsRequest{}
}
func (m *GetMoneyTransactionShippingEtopsRequest) String() string { return proto.CompactTextString(m) }
func (*GetMoneyTransactionShippingEtopsRequest) ProtoMessage()    {}
func (*GetMoneyTransactionShippingEtopsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{21}
}

var xxx_messageInfo_GetMoneyTransactionShippingEtopsRequest proto.InternalMessageInfo

type ConfirmMoneyTransactionShippingEtopRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	TotalCod             dot.ID   `protobuf:"varint,3,opt,name=total_cod,json=totalCod" json:"total_cod"`
	TotalAmount          dot.ID   `protobuf:"varint,4,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	TotalOrders          dot.ID   `protobuf:"varint,5,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmMoneyTransactionShippingEtopRequest) Reset() {
	*m = ConfirmMoneyTransactionShippingEtopRequest{}
}
func (m *ConfirmMoneyTransactionShippingEtopRequest) String() string {
	return proto.CompactTextString(m)
}
func (*ConfirmMoneyTransactionShippingEtopRequest) ProtoMessage() {}
func (*ConfirmMoneyTransactionShippingEtopRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{22}
}

var xxx_messageInfo_ConfirmMoneyTransactionShippingEtopRequest proto.InternalMessageInfo

type UpdateMoneyTransactionRequest struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Note                 string            `protobuf:"bytes,2,opt,name=note" json:"note"`
	InvoiceNumber        string            `protobuf:"bytes,3,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount          *etop.BankAccount `protobuf:"bytes,4,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateMoneyTransactionRequest) Reset()         { *m = UpdateMoneyTransactionRequest{} }
func (m *UpdateMoneyTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateMoneyTransactionRequest) ProtoMessage()    {}
func (*UpdateMoneyTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{23}
}

var xxx_messageInfo_UpdateMoneyTransactionRequest proto.InternalMessageInfo

type UpdateMoneyTransactionShippingExternalRequest struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Note                 string            `protobuf:"bytes,2,opt,name=note" json:"note"`
	InvoiceNumber        string            `protobuf:"bytes,3,opt,name=invoice_number,json=invoiceNumber" json:"invoice_number"`
	BankAccount          *etop.BankAccount `protobuf:"bytes,4,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateMoneyTransactionShippingExternalRequest) Reset() {
	*m = UpdateMoneyTransactionShippingExternalRequest{}
}
func (m *UpdateMoneyTransactionShippingExternalRequest) String() string {
	return proto.CompactTextString(m)
}
func (*UpdateMoneyTransactionShippingExternalRequest) ProtoMessage() {}
func (*UpdateMoneyTransactionShippingExternalRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{24}
}

var xxx_messageInfo_UpdateMoneyTransactionShippingExternalRequest proto.InternalMessageInfo

type CreateNotificationsRequest struct {
	Title      string                         `protobuf:"bytes,1,opt,name=title" json:"title"`
	Message    string                         `protobuf:"bytes,2,opt,name=message" json:"message"`
	Entity     notifier_entity.NotifierEntity `protobuf:"varint,3,opt,name=entity,enum=notifier_entity.NotifierEntity" json:"entity"`
	EntityId   dot.ID                         `protobuf:"varint,4,opt,name=entity_id,json=entityId" json:"entity_id"`
	AccountIds []dot.ID                       `protobuf:"varint,5,rep,name=account_ids,json=accountIds" json:"account_ids"`
	// Send to all subscribers
	SendAll              bool     `protobuf:"varint,6,opt,name=send_all,json=sendAll" json:"send_all"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateNotificationsRequest) Reset()         { *m = CreateNotificationsRequest{} }
func (m *CreateNotificationsRequest) String() string { return proto.CompactTextString(m) }
func (*CreateNotificationsRequest) ProtoMessage()    {}
func (*CreateNotificationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{25}
}

var xxx_messageInfo_CreateNotificationsRequest proto.InternalMessageInfo

type CreateNotificationsResponse struct {
	Created              int32    `protobuf:"varint,1,opt,name=created" json:"created"`
	Errored              int32    `protobuf:"varint,2,opt,name=errored" json:"errored"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateNotificationsResponse) Reset()         { *m = CreateNotificationsResponse{} }
func (m *CreateNotificationsResponse) String() string { return proto.CompactTextString(m) }
func (*CreateNotificationsResponse) ProtoMessage()    {}
func (*CreateNotificationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ede47291efc1dd4, []int{26}
}

var xxx_messageInfo_CreateNotificationsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetOrdersRequest)(nil), "admin.GetOrdersRequest")
	proto.RegisterType((*GetFulfillmentsRequest)(nil), "admin.GetFulfillmentsRequest")
	proto.RegisterType((*LoginAsAccountRequest)(nil), "admin.LoginAsAccountRequest")
	proto.RegisterType((*GetMoneyTransactionsRequest)(nil), "admin.GetMoneyTransactionsRequest")
	proto.RegisterType((*RemoveFfmsMoneyTransactionRequest)(nil), "admin.RemoveFfmsMoneyTransactionRequest")
	proto.RegisterType((*GetMoneyTransactionShippingExternalsRequest)(nil), "admin.GetMoneyTransactionShippingExternalsRequest")
	proto.RegisterType((*RemoveMoneyTransactionShippingExternalLinesRequest)(nil), "admin.RemoveMoneyTransactionShippingExternalLinesRequest")
	proto.RegisterType((*ConfirmMoneyTransactionRequest)(nil), "admin.ConfirmMoneyTransactionRequest")
	proto.RegisterType((*DeleteMoneyTransactionRequest)(nil), "admin.DeleteMoneyTransactionRequest")
	proto.RegisterType((*GetShopsRequest)(nil), "admin.GetShopsRequest")
	proto.RegisterType((*GetShopsResponse)(nil), "admin.GetShopsResponse")
	proto.RegisterType((*GetCreditRequest)(nil), "admin.GetCreditRequest")
	proto.RegisterType((*GetCreditsRequest)(nil), "admin.GetCreditsRequest")
	proto.RegisterType((*CreateCreditRequest)(nil), "admin.CreateCreditRequest")
	proto.RegisterType((*UpdateCreditRequest)(nil), "admin.UpdateCreditRequest")
	proto.RegisterType((*ConfirmCreditRequest)(nil), "admin.ConfirmCreditRequest")
	proto.RegisterType((*CreatePartnerRequest)(nil), "admin.CreatePartnerRequest")
	proto.RegisterType((*UpdateFulfillmentRequest)(nil), "admin.UpdateFulfillmentRequest")
	proto.RegisterType((*GenerateAPIKeyRequest)(nil), "admin.GenerateAPIKeyRequest")
	proto.RegisterType((*GenerateAPIKeyResponse)(nil), "admin.GenerateAPIKeyResponse")
	proto.RegisterType((*UpdateMoneyTransactionShippingEtopRequest)(nil), "admin.UpdateMoneyTransactionShippingEtopRequest")
	proto.RegisterType((*GetMoneyTransactionShippingEtopsRequest)(nil), "admin.GetMoneyTransactionShippingEtopsRequest")
	proto.RegisterType((*ConfirmMoneyTransactionShippingEtopRequest)(nil), "admin.ConfirmMoneyTransactionShippingEtopRequest")
	proto.RegisterType((*UpdateMoneyTransactionRequest)(nil), "admin.UpdateMoneyTransactionRequest")
	proto.RegisterType((*UpdateMoneyTransactionShippingExternalRequest)(nil), "admin.UpdateMoneyTransactionShippingExternalRequest")
	proto.RegisterType((*CreateNotificationsRequest)(nil), "admin.CreateNotificationsRequest")
	proto.RegisterType((*CreateNotificationsResponse)(nil), "admin.CreateNotificationsResponse")
}

func init() { proto.RegisterFile("etop/admin/admin.proto", fileDescriptor_7ede47291efc1dd4) }

var fileDescriptor_7ede47291efc1dd4 = []byte{
	// 1490 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x58, 0xcf, 0x6f, 0x1b, 0xc5,
	0x17, 0xf7, 0x7a, 0xfd, 0x2b, 0x2f, 0x6d, 0xd2, 0x6e, 0xd2, 0x7e, 0xdd, 0xf4, 0x9b, 0x8d, 0xbb,
	0xdf, 0x7e, 0x69, 0xa0, 0xaa, 0x2d, 0xa5, 0xa5, 0x27, 0x38, 0x24, 0x69, 0x1b, 0x59, 0x29, 0x21,
	0xda, 0x16, 0x84, 0xaa, 0xa2, 0xd5, 0xc4, 0x3b, 0x76, 0x46, 0xdd, 0x9d, 0x59, 0x76, 0xc7, 0x01,
	0x1f, 0xb8, 0xc3, 0x01, 0x89, 0x23, 0x77, 0x2e, 0x70, 0xe1, 0x80, 0x90, 0xe0, 0xcc, 0xa9, 0x82,
	0x4b, 0xff, 0x02, 0xd4, 0xa4, 0x57, 0x0e, 0x15, 0x67, 0x0e, 0x68, 0x7e, 0xac, 0xbd, 0x76, 0x12,
	0xc7, 0x11, 0x08, 0xb8, 0x78, 0x77, 0x3e, 0xef, 0xcd, 0xe4, 0xfd, 0xf8, 0xec, 0x7b, 0x6f, 0x02,
	0x17, 0x31, 0x67, 0x51, 0x03, 0xf9, 0x21, 0xa1, 0xea, 0xb7, 0x1e, 0xc5, 0x8c, 0x33, 0xab, 0x28,
	0x17, 0x0b, 0xf3, 0x1d, 0xd6, 0x61, 0x12, 0x69, 0x88, 0x37, 0x25, 0x5c, 0x98, 0x6b, 0xb1, 0x30,
	0x64, 0xb4, 0xa1, 0x1e, 0x1a, 0x9c, 0x95, 0x27, 0x89, 0x1f, 0x0d, 0xa8, 0xa3, 0x59, 0xec, 0xe3,
	0x58, 0xfd, 0x6a, 0x7c, 0x51, 0x2b, 0xb6, 0x1a, 0x09, 0x47, 0xbc, 0x9b, 0xdc, 0xd4, 0xcf, 0xc3,
	0xe2, 0x5d, 0x12, 0x45, 0x84, 0x76, 0xa4, 0x1c, 0x6b, 0xf1, 0x52, 0x5f, 0xdc, 0x8a, 0xb1, 0x4f,
	0xb8, 0xc7, 0x7b, 0x11, 0x6e, 0x88, 0x1f, 0xad, 0xf0, 0xff, 0xbe, 0x02, 0x65, 0x9c, 0xb4, 0x09,
	0x8e, 0x3d, 0x4c, 0x39, 0xe1, 0xbd, 0x86, 0x7a, 0xa4, 0xe7, 0x74, 0x18, 0xeb, 0x04, 0xb8, 0x21,
	0x57, 0x3b, 0xdd, 0x76, 0x83, 0x93, 0x10, 0x27, 0x1c, 0x85, 0xda, 0x7c, 0xe7, 0x31, 0x9c, 0xdb,
	0xc0, 0xfc, 0x6d, 0x61, 0x78, 0xe2, 0xe2, 0x0f, 0xba, 0x38, 0xe1, 0x96, 0x03, 0xa5, 0x08, 0x75,
	0x08, 0xed, 0x54, 0x8d, 0x9a, 0xb1, 0x3c, 0xbd, 0x02, 0xf5, 0x56, 0x58, 0xdf, 0x96, 0x88, 0xab,
	0x25, 0xd6, 0x55, 0x28, 0xb7, 0x49, 0xc0, 0x71, 0x9c, 0x54, 0x0b, 0x35, 0x33, 0x55, 0xba, 0x27,
	0x21, 0x37, 0x15, 0x39, 0x3f, 0x1b, 0x70, 0x71, 0x03, 0xf3, 0x7b, 0xdd, 0xa0, 0x4d, 0x82, 0x20,
	0xc4, 0x94, 0x9f, 0xea, 0x8f, 0x2c, 0x42, 0x39, 0xd9, 0x65, 0x91, 0x47, 0xfc, 0x6a, 0xbe, 0x66,
	0x2c, 0x9b, 0x6b, 0x85, 0xa7, 0xbf, 0x2c, 0xe5, 0xdc, 0x92, 0x00, 0x9b, 0xbe, 0xb5, 0x04, 0x15,
	0x19, 0x71, 0x21, 0x2f, 0x64, 0xe4, 0x65, 0x89, 0x36, 0x7d, 0xeb, 0x1a, 0x94, 0x54, 0xd0, 0xab,
	0xc5, 0x9a, 0xb1, 0x3c, 0xb3, 0x32, 0x5b, 0xd7, 0xb9, 0xa8, 0x3f, 0x90, 0x4f, 0x57, 0x8b, 0xb3,
	0xde, 0x94, 0x8e, 0xf7, 0xe6, 0x63, 0xb8, 0x70, 0x9f, 0x75, 0x08, 0x5d, 0x4d, 0x56, 0x5b, 0x2d,
	0xd6, 0xa5, 0x3c, 0xf5, 0x65, 0x11, 0xca, 0xdd, 0x44, 0xd9, 0x61, 0x64, 0xed, 0x14, 0x60, 0xd3,
	0xb7, 0xfe, 0x07, 0x80, 0xd4, 0x86, 0x51, 0x4f, 0xa6, 0x34, 0xde, 0xf4, 0xad, 0x1a, 0x54, 0x22,
	0x94, 0x24, 0x1f, 0xb2, 0xd8, 0xaf, 0x9a, 0x35, 0x63, 0x79, 0x4a, 0xab, 0xf4, 0x51, 0xe7, 0x37,
	0x03, 0x2e, 0x6f, 0x60, 0xfe, 0x16, 0xa3, 0xb8, 0xf7, 0x30, 0x46, 0x34, 0x41, 0x2d, 0x4e, 0x18,
	0xed, 0x47, 0xf4, 0x1c, 0x98, 0xc4, 0x4f, 0xaa, 0x46, 0xcd, 0x5c, 0x36, 0x5d, 0xf1, 0x7a, 0x52,
	0xfc, 0xde, 0x83, 0x57, 0x42, 0x71, 0x98, 0xc7, 0x07, 0xa7, 0x79, 0x29, 0x1d, 0x3d, 0xfc, 0x11,
	0xc7, 0x31, 0x45, 0x81, 0xd8, 0x6d, 0x66, 0x76, 0x3b, 0xe1, 0x88, 0x01, 0x0f, 0xf4, 0x8e, 0xbb,
	0x7a, 0x43, 0xd3, 0xcf, 0x24, 0xb7, 0x30, 0x09, 0x83, 0x8a, 0xc7, 0xc7, 0xfc, 0x4b, 0x03, 0xae,
	0xb8, 0x38, 0x64, 0x7b, 0xf8, 0x5e, 0x3b, 0x4c, 0x46, 0x7d, 0x4f, 0x5d, 0xbf, 0x06, 0xb3, 0xed,
	0x01, 0xc7, 0xbc, 0x41, 0x18, 0x66, 0x32, 0x70, 0xd3, 0x4f, 0xac, 0xdb, 0x30, 0x7f, 0xd8, 0xe5,
	0x91, 0xf0, 0x58, 0xa3, 0x0e, 0x36, 0xfd, 0x6c, 0x24, 0xcd, 0xc3, 0x91, 0x74, 0x3e, 0x35, 0xe0,
	0xfa, 0x11, 0xa9, 0x19, 0x8d, 0xcc, 0x98, 0x54, 0x0d, 0x22, 0x96, 0x9f, 0x24, 0x62, 0xe6, 0xf1,
	0x11, 0xfb, 0xda, 0x80, 0x15, 0x15, 0xb1, 0x93, 0xcc, 0xb9, 0x4f, 0x28, 0xee, 0x9b, 0x74, 0x09,
	0x2a, 0x01, 0xa1, 0x38, 0x13, 0xbb, 0xb2, 0x58, 0x8b, 0xa0, 0x4d, 0xce, 0x93, 0xfc, 0xe9, 0x78,
	0xe2, 0xfc, 0x6a, 0x80, 0xbd, 0xce, 0x68, 0x9b, 0xc4, 0xe1, 0x71, 0xa9, 0x3d, 0x2e, 0x63, 0xc6,
	0xe4, 0x19, 0x3b, 0x8a, 0xfb, 0x57, 0x60, 0x8a, 0x33, 0x8e, 0x02, 0xaf, 0xc5, 0x86, 0x53, 0x5a,
	0x91, 0xf0, 0x3a, 0x13, 0xd5, 0xe3, 0x8c, 0x52, 0x41, 0xa1, 0xf8, 0x46, 0x87, 0x4a, 0xcc, 0xb4,
	0x94, 0xac, 0x4a, 0xc1, 0x40, 0x51, 0xd6, 0x1d, 0x55, 0x6c, 0x86, 0x15, 0x55, 0x7d, 0x75, 0xf6,
	0x60, 0xf1, 0x0e, 0x0e, 0x30, 0xc7, 0x7f, 0xaf, 0xb3, 0xce, 0xeb, 0x30, 0xbb, 0x81, 0xf9, 0x83,
	0x5d, 0x16, 0x9d, 0xa6, 0xfc, 0x3a, 0x8f, 0x64, 0x6f, 0xd0, 0xdb, 0x92, 0x88, 0xd1, 0x04, 0x5b,
	0x57, 0x47, 0xf6, 0x9d, 0xd1, 0xfb, 0x70, 0x93, 0xb6, 0x59, 0x9f, 0xa9, 0x35, 0x28, 0x8a, 0x3f,
	0x9d, 0x54, 0xf3, 0x9a, 0xa7, 0xb2, 0x61, 0x8a, 0x93, 0x5c, 0x25, 0x70, 0x36, 0xe4, 0xd9, 0xeb,
	0xb2, 0xb9, 0xa5, 0x36, 0xcd, 0x43, 0x7e, 0xc4, 0xd7, 0x3c, 0x39, 0xd1, 0xb7, 0x77, 0xe1, 0x7c,
	0xff, 0xa0, 0x24, 0x53, 0x90, 0xd3, 0x3d, 0xc6, 0x11, 0xc9, 0x9f, 0xe0, 0x63, 0x73, 0xbe, 0x33,
	0x60, 0x6e, 0x3d, 0xc6, 0x88, 0xe3, 0x61, 0x23, 0xff, 0x0b, 0x25, 0xcd, 0x87, 0xa1, 0x93, 0x15,
	0x76, 0x12, 0xeb, 0xae, 0x43, 0x41, 0xf4, 0x70, 0x49, 0xb8, 0x99, 0x95, 0xff, 0xd4, 0x33, 0xcd,
	0xbd, 0xae, 0xfe, 0xcc, 0xc3, 0x5e, 0x84, 0x5d, 0xa9, 0x64, 0xdd, 0x84, 0x72, 0x84, 0x88, 0xef,
	0x21, 0xae, 0xab, 0xe8, 0x42, 0x5d, 0x75, 0xf3, 0x7a, 0xda, 0xcd, 0xeb, 0x0f, 0xd3, 0x6e, 0x2e,
	0xcc, 0x26, 0xfe, 0x2a, 0x77, 0x7e, 0x32, 0x60, 0xee, 0x9d, 0xc8, 0x3f, 0x64, 0xf6, 0xd1, 0xb1,
	0x1d, 0x38, 0x93, 0x1f, 0xef, 0x8c, 0x39, 0xc6, 0x99, 0xc2, 0x29, 0x9d, 0x29, 0x4e, 0xec, 0xcc,
	0x26, 0xcc, 0xeb, 0xea, 0xf0, 0x17, 0x10, 0xe5, 0x2e, 0xcc, 0xab, 0x7c, 0x6e, 0xa3, 0x98, 0x53,
	0x1c, 0xa7, 0x87, 0xdd, 0x10, 0x96, 0x49, 0x44, 0x53, 0xfa, 0xac, 0x62, 0xab, 0x56, 0x4b, 0x67,
	0x0a, 0xad, 0xe3, 0xbc, 0xcc, 0x43, 0x55, 0x05, 0x38, 0x33, 0xd5, 0x8c, 0x37, 0xec, 0x0a, 0x4c,
	0xb5, 0xbb, 0x41, 0xe0, 0x51, 0x14, 0x62, 0x69, 0x5a, 0xbf, 0xb7, 0x0b, 0x78, 0x0b, 0x85, 0xd8,
	0x5a, 0x80, 0x62, 0xb4, 0xcb, 0x28, 0x1e, 0x6a, 0xfd, 0x0a, 0xb2, 0x96, 0xe1, 0x5c, 0xbf, 0x54,
	0xa5, 0xb5, 0x48, 0xc4, 0xb0, 0xe8, 0xce, 0xa4, 0xb5, 0x4a, 0x17, 0xa2, 0x5b, 0x30, 0x47, 0x12,
	0x4f, 0x58, 0x4a, 0x50, 0xe0, 0xf9, 0x38, 0x20, 0x7b, 0x38, 0xee, 0x55, 0x4b, 0x35, 0x63, 0xb9,
	0xa2, 0xcf, 0x3c, 0x4f, 0x92, 0x6d, 0x25, 0xbf, 0xa3, 0xc5, 0x72, 0x3c, 0x11, 0x63, 0xb0, 0x47,
	0x19, 0xc7, 0xd5, 0x72, 0xc6, 0x80, 0x29, 0x89, 0x6f, 0x31, 0x8e, 0xad, 0x35, 0x58, 0x40, 0x2d,
	0xde, 0x95, 0x56, 0x84, 0x11, 0xa6, 0x09, 0x92, 0xb5, 0x49, 0x9b, 0x53, 0x11, 0xe6, 0xe8, 0x4d,
	0x55, 0xa5, 0xb7, 0x9e, 0x51, 0xd3, 0xe6, 0xdd, 0x86, 0x99, 0x7e, 0xd7, 0x90, 0xc3, 0x6e, 0x75,
	0x2a, 0x1d, 0xcb, 0x34, 0x2c, 0xe7, 0x32, 0xec, 0x9e, 0x4d, 0xd7, 0x72, 0xe9, 0xbc, 0x01, 0x17,
	0x36, 0x30, 0xc5, 0x31, 0xe2, 0x78, 0x75, 0xbb, 0xb9, 0x89, 0x7b, 0x69, 0xb8, 0x87, 0x07, 0x2b,
	0xe3, 0xc8, 0xc1, 0xca, 0x79, 0x2c, 0x46, 0xd0, 0xe1, 0xdd, 0xba, 0x96, 0x4d, 0xb2, 0x5d, 0xb0,
	0x0a, 0x45, 0xc4, 0x7b, 0x82, 0x7b, 0x43, 0xa9, 0x2b, 0xa1, 0x88, 0x6c, 0xe2, 0x9e, 0xf3, 0x59,
	0x1e, 0x5e, 0x55, 0x74, 0x38, 0xb6, 0xdb, 0x72, 0x16, 0x8d, 0xe7, 0x87, 0x05, 0x05, 0xe4, 0xfb,
	0xaa, 0x58, 0x9a, 0xae, 0x7c, 0xb7, 0xaa, 0x50, 0xf6, 0x65, 0xab, 0x50, 0xbd, 0xde, 0x74, 0xd3,
	0xa5, 0xb5, 0x04, 0xd3, 0x31, 0x8e, 0x02, 0xd4, 0xc2, 0x1e, 0x0a, 0x02, 0x39, 0x7d, 0x9b, 0x2e,
	0x68, 0x68, 0x35, 0x08, 0xac, 0x2a, 0x14, 0x64, 0x26, 0x8b, 0x19, 0x73, 0x25, 0x62, 0x5d, 0x87,
	0x19, 0x42, 0xf7, 0x18, 0x69, 0x61, 0x8f, 0x76, 0xc3, 0x1d, 0x1c, 0x4b, 0x6a, 0xa4, 0x3a, 0x67,
	0xb5, 0x6c, 0x4b, 0x8a, 0xac, 0x5b, 0x70, 0x66, 0x07, 0xd1, 0x27, 0x9e, 0x0e, 0x85, 0x24, 0xc6,
	0xf4, 0xca, 0x79, 0xf5, 0x71, 0xac, 0x21, 0xfa, 0x24, 0x1d, 0x82, 0xa7, 0x77, 0x06, 0x0b, 0xe7,
	0x1b, 0x03, 0xae, 0x8d, 0x9b, 0x84, 0x78, 0xa6, 0x07, 0x1d, 0x9e, 0x82, 0x06, 0x03, 0x7b, 0x7e,
	0xfc, 0xc0, 0x3e, 0xa8, 0xe0, 0xe6, 0x9f, 0xbc, 0xa2, 0x7c, 0x6f, 0xc0, 0x6b, 0xc7, 0x8c, 0x20,
	0x93, 0x67, 0xf0, 0x1f, 0x99, 0x26, 0xbe, 0x35, 0x60, 0xf1, 0x68, 0xea, 0x8d, 0x37, 0x36, 0xe5,
	0x47, 0x7e, 0x02, 0x7e, 0x98, 0x93, 0xf3, 0xa3, 0x30, 0x11, 0x3f, 0x7e, 0x34, 0xe0, 0xc6, 0x09,
	0xdf, 0x8b, 0x1e, 0x0f, 0xff, 0xc5, 0x4e, 0xfc, 0x6e, 0xc0, 0x82, 0xea, 0x25, 0x5b, 0xf2, 0xf2,
	0xdd, 0x42, 0x43, 0x17, 0xb1, 0x05, 0x28, 0x72, 0xc2, 0x03, 0x2c, 0x8d, 0xee, 0x17, 0x73, 0x09,
	0x59, 0x36, 0x94, 0x43, 0x9c, 0x24, 0xa8, 0x33, 0x6c, 0x7a, 0x0a, 0x5a, 0x6f, 0x42, 0x49, 0x5d,
	0xe0, 0xf5, 0x8c, 0xb0, 0x54, 0x1f, 0xb9, 0xdf, 0xd7, 0xb7, 0xf4, 0xfa, 0xae, 0x5c, 0xa6, 0xe5,
	0x48, 0x09, 0x05, 0x11, 0xd5, 0xdb, 0xe8, 0x9d, 0xb8, 0xa2, 0x60, 0x79, 0x6b, 0x9e, 0x1e, 0x54,
	0x3d, 0x75, 0xf7, 0x32, 0x5d, 0xe8, 0x17, 0x3c, 0x51, 0x60, 0x2a, 0x09, 0xa6, 0xbe, 0xac, 0x2e,
	0xd9, 0xd6, 0x51, 0x16, 0xe8, 0x6a, 0x10, 0x38, 0xef, 0xc3, 0xe5, 0x23, 0xbd, 0xd7, 0x65, 0xd5,
	0x86, 0x72, 0x4b, 0x8a, 0x55, 0xd6, 0xd2, 0xbe, 0x90, 0x82, 0x42, 0x8e, 0xe3, 0x98, 0xc5, 0x58,
	0xf5, 0xe9, 0xbe, 0x5c, 0x83, 0x6b, 0x9b, 0x4f, 0xf7, 0xed, 0xdc, 0xb3, 0x7d, 0xdb, 0x78, 0xbe,
	0x6f, 0xe7, 0x5e, 0xee, 0xdb, 0xb9, 0x4f, 0x0e, 0xec, 0xdc, 0x57, 0x07, 0x76, 0xee, 0x87, 0x03,
	0x3b, 0xf7, 0xf4, 0xc0, 0xce, 0x3d, 0x3b, 0xb0, 0x73, 0xcf, 0x0f, 0xec, 0xdc, 0xe7, 0x2f, 0xec,
	0xdc, 0x17, 0x2f, 0xec, 0xdc, 0xa3, 0x4b, 0x32, 0x5d, 0x7b, 0xb4, 0x81, 0x22, 0xd2, 0x88, 0x76,
	0x1a, 0x83, 0xff, 0xf6, 0xfc, 0x11, 0x00, 0x00, 0xff, 0xff, 0xdf, 0x2c, 0x71, 0xa5, 0xfa, 0x11,
	0x00, 0x00,
}
