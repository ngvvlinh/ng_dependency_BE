// source: etop/shop/shop.proto

package shop

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	spreadsheet "etop.vn/api/pb/common/spreadsheet"
	etop "etop.vn/api/pb/etop"
	ghn_note_code "etop.vn/api/pb/etop/etc/ghn_note_code"
	payment_provider "etop.vn/api/pb/etop/etc/payment_provider"
	product_type "etop.vn/api/pb/etop/etc/product_type"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	status3 "etop.vn/api/pb/etop/etc/status3"
	status4 "etop.vn/api/pb/etop/etc/status4"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	order "etop.vn/api/pb/etop/order"
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

type InvitationsResponse struct {
	Invitations          []*etop.Invitation `protobuf:"bytes,1,rep,name=invitations" json:"invitations"`
	Paging               *common.PageInfo   `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *InvitationsResponse) Reset()         { *m = InvitationsResponse{} }
func (m *InvitationsResponse) String() string { return proto.CompactTextString(m) }
func (*InvitationsResponse) ProtoMessage()    {}
func (*InvitationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{0}
}

var xxx_messageInfo_InvitationsResponse proto.InternalMessageInfo

func (m *InvitationsResponse) GetInvitations() []*etop.Invitation {
	if m != nil {
		return m.Invitations
	}
	return nil
}

func (m *InvitationsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type GetInvitationsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetInvitationsRequest) Reset()         { *m = GetInvitationsRequest{} }
func (m *GetInvitationsRequest) String() string { return proto.CompactTextString(m) }
func (*GetInvitationsRequest) ProtoMessage()    {}
func (*GetInvitationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{1}
}

var xxx_messageInfo_GetInvitationsRequest proto.InternalMessageInfo

func (m *GetInvitationsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetInvitationsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type CreateInvitationRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email" json:"email"`
	Roles                []string `protobuf:"bytes,2,rep,name=roles" json:"roles"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateInvitationRequest) Reset()         { *m = CreateInvitationRequest{} }
func (m *CreateInvitationRequest) String() string { return proto.CompactTextString(m) }
func (*CreateInvitationRequest) ProtoMessage()    {}
func (*CreateInvitationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{2}
}

var xxx_messageInfo_CreateInvitationRequest proto.InternalMessageInfo

func (m *CreateInvitationRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateInvitationRequest) GetRoles() []string {
	if m != nil {
		return m.Roles
	}
	return nil
}

type PurchaseOrder struct {
	Id                   dot.ID                 `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID                 `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	SupplierId           dot.ID                 `protobuf:"varint,3,opt,name=supplier_id,json=supplierId" json:"supplier_id"`
	Supplier             *PurchaseOrderSupplier `protobuf:"bytes,4,opt,name=supplier" json:"supplier"`
	BasketValue          int                    `protobuf:"varint,5,opt,name=basket_value,json=basketValue" json:"basket_value"`
	TotalDiscount        int                    `protobuf:"varint,6,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount          int                    `protobuf:"varint,7,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	Code                 string                 `protobuf:"bytes,8,opt,name=code" json:"code"`
	Note                 string                 `protobuf:"bytes,9,opt,name=note" json:"note"`
	Status               status3.Status         `protobuf:"varint,10,opt,name=status,enum=status3.Status" json:"status"`
	Lines                []*PurchaseOrderLine   `protobuf:"bytes,11,rep,name=lines" json:"lines"`
	CreatedBy            dot.ID                 `protobuf:"varint,12,opt,name=created_by,json=createdBy" json:"created_by"`
	CancelledReason      string                 `protobuf:"bytes,13,opt,name=cancelled_reason,json=cancelledReason" json:"cancelled_reason"`
	ConfirmedAt          dot.Time               `protobuf:"bytes,14,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelledAt          dot.Time               `protobuf:"bytes,15,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	CreatedAt            dot.Time               `protobuf:"bytes,16,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time               `protobuf:"bytes,17,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	DeletedAt            dot.Time               `protobuf:"bytes,18,opt,name=deleted_at,json=deletedAt" json:"deleted_at"`
	InventoryVoucher     *InventoryVoucher      `protobuf:"bytes,19,opt,name=inventory_voucher,json=inventoryVoucher" json:"inventory_voucher"`
	PaidAmount           int                    `protobuf:"varint,20,opt,name=paid_amount,json=paidAmount" json:"paid_amount"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *PurchaseOrder) Reset()         { *m = PurchaseOrder{} }
func (m *PurchaseOrder) String() string { return proto.CompactTextString(m) }
func (*PurchaseOrder) ProtoMessage()    {}
func (*PurchaseOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{3}
}

var xxx_messageInfo_PurchaseOrder proto.InternalMessageInfo

func (m *PurchaseOrder) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PurchaseOrder) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *PurchaseOrder) GetSupplierId() dot.ID {
	if m != nil {
		return m.SupplierId
	}
	return 0
}

func (m *PurchaseOrder) GetSupplier() *PurchaseOrderSupplier {
	if m != nil {
		return m.Supplier
	}
	return nil
}

func (m *PurchaseOrder) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *PurchaseOrder) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *PurchaseOrder) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *PurchaseOrder) GetLines() []*PurchaseOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *PurchaseOrder) GetCreatedBy() dot.ID {
	if m != nil {
		return m.CreatedBy
	}
	return 0
}

func (m *PurchaseOrder) GetCancelledReason() string {
	if m != nil {
		return m.CancelledReason
	}
	return ""
}

func (m *PurchaseOrder) GetInventoryVoucher() *InventoryVoucher {
	if m != nil {
		return m.InventoryVoucher
	}
	return nil
}

type PurchaseOrderSupplier struct {
	FullName             string   `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,3,opt,name=email" json:"email"`
	CompanyName          string   `protobuf:"bytes,4,opt,name=company_name,json=companyName" json:"company_name"`
	TaxNumber            string   `protobuf:"bytes,5,opt,name=tax_number,json=taxNumber" json:"tax_number"`
	HeadquarterAddress   string   `protobuf:"bytes,6,opt,name=headquarter_address,json=headquarterAddress" json:"headquarter_address"`
	Deleted              bool     `protobuf:"varint,7,opt,name=deleted" json:"deleted"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PurchaseOrderSupplier) Reset()         { *m = PurchaseOrderSupplier{} }
func (m *PurchaseOrderSupplier) String() string { return proto.CompactTextString(m) }
func (*PurchaseOrderSupplier) ProtoMessage()    {}
func (*PurchaseOrderSupplier) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{4}
}

var xxx_messageInfo_PurchaseOrderSupplier proto.InternalMessageInfo

func (m *PurchaseOrderSupplier) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetTaxNumber() string {
	if m != nil {
		return m.TaxNumber
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetHeadquarterAddress() string {
	if m != nil {
		return m.HeadquarterAddress
	}
	return ""
}

func (m *PurchaseOrderSupplier) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

type GetPurchaseOrdersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetPurchaseOrdersRequest) Reset()         { *m = GetPurchaseOrdersRequest{} }
func (m *GetPurchaseOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*GetPurchaseOrdersRequest) ProtoMessage()    {}
func (*GetPurchaseOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{5}
}

var xxx_messageInfo_GetPurchaseOrdersRequest proto.InternalMessageInfo

func (m *GetPurchaseOrdersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetPurchaseOrdersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type PurchaseOrdersResponse struct {
	PurchaseOrders       []*PurchaseOrder `protobuf:"bytes,1,rep,name=purchase_orders,json=purchaseOrders" json:"purchase_orders"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PurchaseOrdersResponse) Reset()         { *m = PurchaseOrdersResponse{} }
func (m *PurchaseOrdersResponse) String() string { return proto.CompactTextString(m) }
func (*PurchaseOrdersResponse) ProtoMessage()    {}
func (*PurchaseOrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{6}
}

var xxx_messageInfo_PurchaseOrdersResponse proto.InternalMessageInfo

func (m *PurchaseOrdersResponse) GetPurchaseOrders() []*PurchaseOrder {
	if m != nil {
		return m.PurchaseOrders
	}
	return nil
}

func (m *PurchaseOrdersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type CreatePurchaseOrderRequest struct {
	SupplierId           dot.ID               `protobuf:"varint,1,opt,name=supplier_id,json=supplierId" json:"supplier_id"`
	BasketValue          int                  `protobuf:"varint,2,opt,name=basket_value,json=basketValue" json:"basket_value"`
	TotalDiscount        int                  `protobuf:"varint,3,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount          int                  `protobuf:"varint,4,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	Note                 string               `protobuf:"bytes,5,opt,name=note" json:"note"`
	Lines                []*PurchaseOrderLine `protobuf:"bytes,6,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CreatePurchaseOrderRequest) Reset()         { *m = CreatePurchaseOrderRequest{} }
func (m *CreatePurchaseOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePurchaseOrderRequest) ProtoMessage()    {}
func (*CreatePurchaseOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{7}
}

var xxx_messageInfo_CreatePurchaseOrderRequest proto.InternalMessageInfo

func (m *CreatePurchaseOrderRequest) GetSupplierId() dot.ID {
	if m != nil {
		return m.SupplierId
	}
	return 0
}

func (m *CreatePurchaseOrderRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreatePurchaseOrderRequest) GetLines() []*PurchaseOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type UpdatePurchaseOrderRequest struct {
	Id                   dot.ID               `protobuf:"varint,1,opt,name=id" json:"id"`
	SupplierId           *dot.ID              `protobuf:"varint,2,opt,name=supplier_id,json=supplierId" json:"supplier_id"`
	BasketValue          *int                 `protobuf:"varint,3,opt,name=basket_value,json=basketValue" json:"basket_value"`
	TotalDiscount        *int                 `protobuf:"varint,4,opt,name=total_discount,json=totalDiscount" json:"total_discount"`
	TotalAmount          *int                 `protobuf:"varint,5,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	Note                 *string              `protobuf:"bytes,6,opt,name=note" json:"note"`
	Lines                []*PurchaseOrderLine `protobuf:"bytes,7,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *UpdatePurchaseOrderRequest) Reset()         { *m = UpdatePurchaseOrderRequest{} }
func (m *UpdatePurchaseOrderRequest) String() string { return proto.CompactTextString(m) }
func (*UpdatePurchaseOrderRequest) ProtoMessage()    {}
func (*UpdatePurchaseOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{8}
}

var xxx_messageInfo_UpdatePurchaseOrderRequest proto.InternalMessageInfo

func (m *UpdatePurchaseOrderRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdatePurchaseOrderRequest) GetSupplierId() dot.ID {
	if m != nil && m.SupplierId != nil {
		return *m.SupplierId
	}
	return 0
}

func (m *UpdatePurchaseOrderRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdatePurchaseOrderRequest) GetLines() []*PurchaseOrderLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type PurchaseOrderLine struct {
	VariantId            dot.ID       `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	Quantity             int32        `protobuf:"varint,2,opt,name=quantity" json:"quantity"`
	PaymentPrice         int32        `protobuf:"varint,3,opt,name=payment_price,json=paymentPrice" json:"payment_price"`
	ProductId            dot.ID       `protobuf:"varint,4,opt,name=product_id,json=productId" json:"product_idint32"`
	ProductName          string       `protobuf:"bytes,7,opt,name=product_name,json=productName" json:"product_name"`
	ImageUrl             string       `protobuf:"bytes,5,opt,name=image_url,json=imageUrl" json:"image_url"`
	Code                 string       `protobuf:"bytes,9,opt,name=code" json:"code"`
	Attributes           []*Attribute `protobuf:"bytes,6,rep,name=attributes" json:"attributes"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PurchaseOrderLine) Reset()         { *m = PurchaseOrderLine{} }
func (m *PurchaseOrderLine) String() string { return proto.CompactTextString(m) }
func (*PurchaseOrderLine) ProtoMessage()    {}
func (*PurchaseOrderLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{9}
}

var xxx_messageInfo_PurchaseOrderLine proto.InternalMessageInfo

func (m *PurchaseOrderLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *PurchaseOrderLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *PurchaseOrderLine) GetPaymentPrice() int32 {
	if m != nil {
		return m.PaymentPrice
	}
	return 0
}

func (m *PurchaseOrderLine) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *PurchaseOrderLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *PurchaseOrderLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *PurchaseOrderLine) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *PurchaseOrderLine) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type CancelPurchaseOrderRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason" json:"reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelPurchaseOrderRequest) Reset()         { *m = CancelPurchaseOrderRequest{} }
func (m *CancelPurchaseOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CancelPurchaseOrderRequest) ProtoMessage()    {}
func (*CancelPurchaseOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{10}
}

var xxx_messageInfo_CancelPurchaseOrderRequest proto.InternalMessageInfo

func (m *CancelPurchaseOrderRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelPurchaseOrderRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ConfirmPurchaseOrderRequest struct {
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// enum create, confirm
	AutoInventoryVoucher string   `protobuf:"bytes,2,opt,name=auto_inventory_voucher,json=autoInventoryVoucher" json:"auto_inventory_voucher"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmPurchaseOrderRequest) Reset()         { *m = ConfirmPurchaseOrderRequest{} }
func (m *ConfirmPurchaseOrderRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmPurchaseOrderRequest) ProtoMessage()    {}
func (*ConfirmPurchaseOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{11}
}

var xxx_messageInfo_ConfirmPurchaseOrderRequest proto.InternalMessageInfo

func (m *ConfirmPurchaseOrderRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ConfirmPurchaseOrderRequest) GetAutoInventoryVoucher() string {
	if m != nil {
		return m.AutoInventoryVoucher
	}
	return ""
}

type GetLedgersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetLedgersRequest) Reset()         { *m = GetLedgersRequest{} }
func (m *GetLedgersRequest) String() string { return proto.CompactTextString(m) }
func (*GetLedgersRequest) ProtoMessage()    {}
func (*GetLedgersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{12}
}

var xxx_messageInfo_GetLedgersRequest proto.InternalMessageInfo

func (m *GetLedgersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetLedgersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type CreateLedgerRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name"`
	// required name, account_number, account_name
	BankAccount          *etop.BankAccount `protobuf:"bytes,2,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	Note                 string            `protobuf:"bytes,3,opt,name=note" json:"note"`
	CreatedBy            string            `protobuf:"bytes,4,opt,name=created_by,json=createdBy" json:"created_by"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CreateLedgerRequest) Reset()         { *m = CreateLedgerRequest{} }
func (m *CreateLedgerRequest) String() string { return proto.CompactTextString(m) }
func (*CreateLedgerRequest) ProtoMessage()    {}
func (*CreateLedgerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{13}
}

var xxx_messageInfo_CreateLedgerRequest proto.InternalMessageInfo

func (m *CreateLedgerRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateLedgerRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *CreateLedgerRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateLedgerRequest) GetCreatedBy() string {
	if m != nil {
		return m.CreatedBy
	}
	return ""
}

type UpdateLedgerRequest struct {
	Id   dot.ID  `protobuf:"varint,1,opt,name=id" json:"id"`
	Name *string `protobuf:"bytes,2,opt,name=name" json:"name"`
	// required name, account_number, account_name
	BankAccount          *etop.BankAccount `protobuf:"bytes,3,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	Note                 *string           `protobuf:"bytes,4,opt,name=note" json:"note"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateLedgerRequest) Reset()         { *m = UpdateLedgerRequest{} }
func (m *UpdateLedgerRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateLedgerRequest) ProtoMessage()    {}
func (*UpdateLedgerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{14}
}

var xxx_messageInfo_UpdateLedgerRequest proto.InternalMessageInfo

func (m *UpdateLedgerRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateLedgerRequest) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpdateLedgerRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *UpdateLedgerRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

type LedgersResponse struct {
	Ledgers              []*Ledger        `protobuf:"bytes,1,rep,name=ledgers" json:"ledgers"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *LedgersResponse) Reset()         { *m = LedgersResponse{} }
func (m *LedgersResponse) String() string { return proto.CompactTextString(m) }
func (*LedgersResponse) ProtoMessage()    {}
func (*LedgersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{15}
}

var xxx_messageInfo_LedgersResponse proto.InternalMessageInfo

func (m *LedgersResponse) GetLedgers() []*Ledger {
	if m != nil {
		return m.Ledgers
	}
	return nil
}

func (m *LedgersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type Ledger struct {
	Id          dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Name        string            `protobuf:"bytes,2,opt,name=name" json:"name"`
	BankAccount *etop.BankAccount `protobuf:"bytes,3,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	Note        string            `protobuf:"bytes,4,opt,name=note" json:"note"`
	// enum: cash, bank
	Type                 string   `protobuf:"bytes,5,opt,name=type" json:"type"`
	CreatedBy            dot.ID   `protobuf:"varint,6,opt,name=created_by,json=createdBy" json:"created_by"`
	CreatedAt            dot.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ledger) Reset()         { *m = Ledger{} }
func (m *Ledger) String() string { return proto.CompactTextString(m) }
func (*Ledger) ProtoMessage()    {}
func (*Ledger) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{16}
}

var xxx_messageInfo_Ledger proto.InternalMessageInfo

func (m *Ledger) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Ledger) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Ledger) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *Ledger) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Ledger) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Ledger) GetCreatedBy() dot.ID {
	if m != nil {
		return m.CreatedBy
	}
	return 0
}

type RegisterShopRequest struct {
	// @required
	Name        string            `protobuf:"bytes,1,opt,name=name" json:"name"`
	Address     *etop.Address     `protobuf:"bytes,2,opt,name=address" json:"address"`
	Phone       string            `protobuf:"bytes,7,opt,name=phone" json:"phone"`
	BankAccount *etop.BankAccount `protobuf:"bytes,8,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	WebsiteUrl  string            `protobuf:"bytes,12,opt,name=website_url,json=websiteUrl" json:"website_url"`
	ImageUrl    string            `protobuf:"bytes,13,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email       string            `protobuf:"bytes,14,opt,name=email" json:"email"`
	UrlSlug     string            `protobuf:"bytes,16,opt,name=url_slug,json=urlSlug" json:"url_slug"`
	CompanyInfo *etop.CompanyInfo `protobuf:"bytes,17,opt,name=company_info,json=companyInfo" json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `protobuf:"bytes,18,opt,name=money_transaction_rrule,json=moneyTransactionRrule" json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `protobuf:"bytes,19,rep,name=survey_info,json=surveyInfo" json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `protobuf:"bytes,20,rep,name=shipping_service_select_strategy,json=shippingServiceSelectStrategy" json:"shipping_service_select_strategy"`
	XXX_NoUnkeyedLiteral          struct{}                                  `json:"-"`
	XXX_sizecache                 int32                                     `json:"-"`
}

func (m *RegisterShopRequest) Reset()         { *m = RegisterShopRequest{} }
func (m *RegisterShopRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterShopRequest) ProtoMessage()    {}
func (*RegisterShopRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{17}
}

var xxx_messageInfo_RegisterShopRequest proto.InternalMessageInfo

func (m *RegisterShopRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RegisterShopRequest) GetAddress() *etop.Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *RegisterShopRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterShopRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *RegisterShopRequest) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *RegisterShopRequest) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *RegisterShopRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *RegisterShopRequest) GetUrlSlug() string {
	if m != nil {
		return m.UrlSlug
	}
	return ""
}

func (m *RegisterShopRequest) GetCompanyInfo() *etop.CompanyInfo {
	if m != nil {
		return m.CompanyInfo
	}
	return nil
}

func (m *RegisterShopRequest) GetMoneyTransactionRrule() string {
	if m != nil {
		return m.MoneyTransactionRrule
	}
	return ""
}

func (m *RegisterShopRequest) GetSurveyInfo() []*etop.SurveyInfo {
	if m != nil {
		return m.SurveyInfo
	}
	return nil
}

func (m *RegisterShopRequest) GetShippingServiceSelectStrategy() []*etop.ShippingServiceSelectStrategyItem {
	if m != nil {
		return m.ShippingServiceSelectStrategy
	}
	return nil
}

type RegisterShopResponse struct {
	// @required
	Shop                 *etop.Shop `protobuf:"bytes,1,opt,name=shop" json:"shop"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RegisterShopResponse) Reset()         { *m = RegisterShopResponse{} }
func (m *RegisterShopResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterShopResponse) ProtoMessage()    {}
func (*RegisterShopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{18}
}

var xxx_messageInfo_RegisterShopResponse proto.InternalMessageInfo

func (m *RegisterShopResponse) GetShop() *etop.Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

type UpdateShopRequest struct {
	InventoryOverstock *bool             `protobuf:"varint,4,opt,name=inventory_overstock,json=inventoryOverstock" json:"inventory_overstock"`
	Name               string            `protobuf:"bytes,1,opt,name=name" json:"name"`
	Address            *etop.Address     `protobuf:"bytes,2,opt,name=address" json:"address"`
	Phone              string            `protobuf:"bytes,7,opt,name=phone" json:"phone"`
	BankAccount        *etop.BankAccount `protobuf:"bytes,8,opt,name=bank_account,json=bankAccount" json:"bank_account"`
	WebsiteUrl         string            `protobuf:"bytes,12,opt,name=website_url,json=websiteUrl" json:"website_url"`
	ImageUrl           string            `protobuf:"bytes,13,opt,name=image_url,json=imageUrl" json:"image_url"`
	Email              string            `protobuf:"bytes,14,opt,name=email" json:"email"`
	AutoCreateFfm      *bool             `protobuf:"varint,15,opt,name=auto_create_ffm,json=autoCreateFfm" json:"auto_create_ffm"`
	// @deprecated use try_on instead
	GhnNoteCode *ghn_note_code.GHNNoteCode `protobuf:"varint,16,opt,name=ghn_note_code,json=ghnNoteCode,enum=ghn_note_code.GHNNoteCode" json:"ghn_note_code"`
	TryOn       *try_on.TryOnCode          `protobuf:"varint,17,opt,name=try_on,json=tryOn,enum=try_on.TryOnCode" json:"try_on"`
	CompanyInfo *etop.CompanyInfo          `protobuf:"bytes,18,opt,name=company_info,json=companyInfo" json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `protobuf:"bytes,19,opt,name=money_transaction_rrule,json=moneyTransactionRrule" json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `protobuf:"bytes,20,rep,name=survey_info,json=surveyInfo" json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `protobuf:"bytes,21,rep,name=shipping_service_select_strategy,json=shippingServiceSelectStrategy" json:"shipping_service_select_strategy"`
	XXX_NoUnkeyedLiteral          struct{}                                  `json:"-"`
	XXX_sizecache                 int32                                     `json:"-"`
}

func (m *UpdateShopRequest) Reset()         { *m = UpdateShopRequest{} }
func (m *UpdateShopRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateShopRequest) ProtoMessage()    {}
func (*UpdateShopRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{19}
}

var xxx_messageInfo_UpdateShopRequest proto.InternalMessageInfo

func (m *UpdateShopRequest) GetInventoryOverstock() bool {
	if m != nil && m.InventoryOverstock != nil {
		return *m.InventoryOverstock
	}
	return false
}

func (m *UpdateShopRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateShopRequest) GetAddress() *etop.Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *UpdateShopRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *UpdateShopRequest) GetBankAccount() *etop.BankAccount {
	if m != nil {
		return m.BankAccount
	}
	return nil
}

func (m *UpdateShopRequest) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *UpdateShopRequest) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *UpdateShopRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UpdateShopRequest) GetAutoCreateFfm() bool {
	if m != nil && m.AutoCreateFfm != nil {
		return *m.AutoCreateFfm
	}
	return false
}

func (m *UpdateShopRequest) GetGhnNoteCode() ghn_note_code.GHNNoteCode {
	if m != nil && m.GhnNoteCode != nil {
		return *m.GhnNoteCode
	}
	return ghn_note_code.GHNNoteCode_unknown
}

func (m *UpdateShopRequest) GetTryOn() try_on.TryOnCode {
	if m != nil && m.TryOn != nil {
		return *m.TryOn
	}
	return try_on.TryOnCode_unknown
}

func (m *UpdateShopRequest) GetCompanyInfo() *etop.CompanyInfo {
	if m != nil {
		return m.CompanyInfo
	}
	return nil
}

func (m *UpdateShopRequest) GetMoneyTransactionRrule() string {
	if m != nil {
		return m.MoneyTransactionRrule
	}
	return ""
}

func (m *UpdateShopRequest) GetSurveyInfo() []*etop.SurveyInfo {
	if m != nil {
		return m.SurveyInfo
	}
	return nil
}

func (m *UpdateShopRequest) GetShippingServiceSelectStrategy() []*etop.ShippingServiceSelectStrategyItem {
	if m != nil {
		return m.ShippingServiceSelectStrategy
	}
	return nil
}

type UpdateShopResponse struct {
	Shop                 *etop.Shop `protobuf:"bytes,1,opt,name=shop" json:"shop"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UpdateShopResponse) Reset()         { *m = UpdateShopResponse{} }
func (m *UpdateShopResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateShopResponse) ProtoMessage()    {}
func (*UpdateShopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{20}
}

var xxx_messageInfo_UpdateShopResponse proto.InternalMessageInfo

func (m *UpdateShopResponse) GetShop() *etop.Shop {
	if m != nil {
		return m.Shop
	}
	return nil
}

type Collection struct {
	// @required
	Id     dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId dot.ID `protobuf:"varint,8,opt,name=shop_id,json=shopId" json:"shop_id"`
	// @required
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name"`
	Description string `protobuf:"bytes,3,opt,name=description" json:"description"`
	ShortDesc   string `protobuf:"bytes,4,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml    string `protobuf:"bytes,5,opt,name=desc_html,json=descHtml" json:"desc_html"`
	// @required
	CreatedAt dot.Time `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	// @required
	UpdatedAt            dot.Time `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Collection) Reset()         { *m = Collection{} }
func (m *Collection) String() string { return proto.CompactTextString(m) }
func (*Collection) ProtoMessage()    {}
func (*Collection) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{21}
}

var xxx_messageInfo_Collection proto.InternalMessageInfo

func (m *Collection) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Collection) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Collection) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Collection) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Collection) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *Collection) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

type CreateCollectionRequest struct {
	// @required
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Description          string   `protobuf:"bytes,3,opt,name=description" json:"description"`
	ShortDesc            string   `protobuf:"bytes,4,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             string   `protobuf:"bytes,5,opt,name=desc_html,json=descHtml" json:"desc_html"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCollectionRequest) Reset()         { *m = CreateCollectionRequest{} }
func (m *CreateCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCollectionRequest) ProtoMessage()    {}
func (*CreateCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{22}
}

var xxx_messageInfo_CreateCollectionRequest proto.InternalMessageInfo

func (m *CreateCollectionRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateCollectionRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateCollectionRequest) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *CreateCollectionRequest) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

type UpdateProductCategoryRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	CategoryId           dot.ID   `protobuf:"varint,2,opt,name=category_id,json=categoryId" json:"category_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductCategoryRequest) Reset()         { *m = UpdateProductCategoryRequest{} }
func (m *UpdateProductCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductCategoryRequest) ProtoMessage()    {}
func (*UpdateProductCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{23}
}

var xxx_messageInfo_UpdateProductCategoryRequest proto.InternalMessageInfo

func (m *UpdateProductCategoryRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *UpdateProductCategoryRequest) GetCategoryId() dot.ID {
	if m != nil {
		return m.CategoryId
	}
	return 0
}

type CollectionsResponse struct {
	Collections          []*ShopCollection `protobuf:"bytes,1,rep,name=collections" json:"collections"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CollectionsResponse) Reset()         { *m = CollectionsResponse{} }
func (m *CollectionsResponse) String() string { return proto.CompactTextString(m) }
func (*CollectionsResponse) ProtoMessage()    {}
func (*CollectionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{24}
}

var xxx_messageInfo_CollectionsResponse proto.InternalMessageInfo

func (m *CollectionsResponse) GetCollections() []*ShopCollection {
	if m != nil {
		return m.Collections
	}
	return nil
}

type UpdateCollectionRequest struct {
	// @required
	Id dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	// @required
	Name                 *string  `protobuf:"bytes,2,opt,name=name" json:"name"`
	Description          *string  `protobuf:"bytes,3,opt,name=description" json:"description"`
	ShortDesc            *string  `protobuf:"bytes,4,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             *string  `protobuf:"bytes,5,opt,name=desc_html,json=descHtml" json:"desc_html"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateCollectionRequest) Reset()         { *m = UpdateCollectionRequest{} }
func (m *UpdateCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCollectionRequest) ProtoMessage()    {}
func (*UpdateCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{25}
}

var xxx_messageInfo_UpdateCollectionRequest proto.InternalMessageInfo

func (m *UpdateCollectionRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateCollectionRequest) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpdateCollectionRequest) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *UpdateCollectionRequest) GetShortDesc() string {
	if m != nil && m.ShortDesc != nil {
		return *m.ShortDesc
	}
	return ""
}

func (m *UpdateCollectionRequest) GetDescHtml() string {
	if m != nil && m.DescHtml != nil {
		return *m.DescHtml
	}
	return ""
}

type UpdateProductsCollectionRequest struct {
	// @required
	CollectionId         dot.ID   `protobuf:"varint,1,opt,name=collection_id,json=collectionId" json:"collection_id"`
	ProductIds           []dot.ID `protobuf:"varint,2,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductsCollectionRequest) Reset()         { *m = UpdateProductsCollectionRequest{} }
func (m *UpdateProductsCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductsCollectionRequest) ProtoMessage()    {}
func (*UpdateProductsCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{26}
}

var xxx_messageInfo_UpdateProductsCollectionRequest proto.InternalMessageInfo

func (m *UpdateProductsCollectionRequest) GetCollectionId() dot.ID {
	if m != nil {
		return m.CollectionId
	}
	return 0
}

func (m *UpdateProductsCollectionRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type RemoveProductsCollectionRequest struct {
	// @required
	CollectionId         dot.ID   `protobuf:"varint,1,opt,name=collection_id,json=collectionId" json:"collection_id"`
	ProductIds           []dot.ID `protobuf:"varint,2,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveProductsCollectionRequest) Reset()         { *m = RemoveProductsCollectionRequest{} }
func (m *RemoveProductsCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveProductsCollectionRequest) ProtoMessage()    {}
func (*RemoveProductsCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{27}
}

var xxx_messageInfo_RemoveProductsCollectionRequest proto.InternalMessageInfo

func (m *RemoveProductsCollectionRequest) GetCollectionId() dot.ID {
	if m != nil {
		return m.CollectionId
	}
	return 0
}

func (m *RemoveProductsCollectionRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type EtopVariant struct {
	Id                   dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string       `protobuf:"bytes,2,opt,name=code" json:"code"`
	Name                 string       `protobuf:"bytes,4,opt,name=name" json:"name"`
	Description          string       `protobuf:"bytes,6,opt,name=description" json:"description"`
	ShortDesc            string       `protobuf:"bytes,7,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             string       `protobuf:"bytes,8,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls            []string     `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	ListPrice            int32        `protobuf:"varint,15,opt,name=list_price,json=listPrice" json:"list_price"`
	CostPrice            int32        `protobuf:"varint,16,opt,name=cost_price,json=costPrice" json:"cost_price"`
	Attributes           []*Attribute `protobuf:"bytes,18,rep,name=attributes" json:"attributes"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *EtopVariant) Reset()         { *m = EtopVariant{} }
func (m *EtopVariant) String() string { return proto.CompactTextString(m) }
func (*EtopVariant) ProtoMessage()    {}
func (*EtopVariant) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{28}
}

var xxx_messageInfo_EtopVariant proto.InternalMessageInfo

func (m *EtopVariant) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EtopVariant) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *EtopVariant) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EtopVariant) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *EtopVariant) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *EtopVariant) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *EtopVariant) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *EtopVariant) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *EtopVariant) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *EtopVariant) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type EtopProduct struct {
	Id          dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Code        string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Name        string   `protobuf:"bytes,4,opt,name=name" json:"name"`
	Description string   `protobuf:"bytes,6,opt,name=description" json:"description"`
	ShortDesc   string   `protobuf:"bytes,7,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml    string   `protobuf:"bytes,8,opt,name=desc_html,json=descHtml" json:"desc_html"`
	Unit        string   `protobuf:"bytes,12,opt,name=unit" json:"unit"`
	ImageUrls   []string `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	ListPrice   int32    `protobuf:"varint,15,opt,name=list_price,json=listPrice" json:"list_price"`
	CostPrice   int32    `protobuf:"varint,16,opt,name=cost_price,json=costPrice" json:"cost_price"`
	CategoryId  dot.ID   `protobuf:"varint,10,opt,name=category_id,json=categoryId" json:"category_id"`
	// @deprecated
	ProductSourceCategoryId dot.ID   `protobuf:"varint,11,opt,name=product_source_category_id,json=productSourceCategoryId" json:"product_source_category_id"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *EtopProduct) Reset()         { *m = EtopProduct{} }
func (m *EtopProduct) String() string { return proto.CompactTextString(m) }
func (*EtopProduct) ProtoMessage()    {}
func (*EtopProduct) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{29}
}

var xxx_messageInfo_EtopProduct proto.InternalMessageInfo

func (m *EtopProduct) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EtopProduct) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *EtopProduct) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EtopProduct) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *EtopProduct) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *EtopProduct) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *EtopProduct) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *EtopProduct) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *EtopProduct) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *EtopProduct) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *EtopProduct) GetCategoryId() dot.ID {
	if m != nil {
		return m.CategoryId
	}
	return 0
}

func (m *EtopProduct) GetProductSourceCategoryId() dot.ID {
	if m != nil {
		return m.ProductSourceCategoryId
	}
	return 0
}

type ShopVariant struct {
	// @required
	Id   dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	Info *EtopVariant `protobuf:"bytes,8,opt,name=info" json:"info"`
	Code string       `protobuf:"bytes,2,opt,name=code" json:"code"`
	// @deprecated use code instead
	EdCode           string                       `protobuf:"bytes,3,opt,name=ed_code,json=edCode" json:"ed_code"`
	Name             string                       `protobuf:"bytes,4,opt,name=name" json:"name"`
	Description      string                       `protobuf:"bytes,5,opt,name=description" json:"description"`
	ShortDesc        string                       `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml         string                       `protobuf:"bytes,7,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls        []string                     `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	ListPrice        int32                        `protobuf:"varint,14,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice      int32                        `protobuf:"varint,15,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	Note             string                       `protobuf:"bytes,11,opt,name=note" json:"note"`
	Status           status3.Status               `protobuf:"varint,12,opt,name=status,enum=status3.Status" json:"status"`
	IsAvailable      bool                         `protobuf:"varint,13,opt,name=is_available,json=isAvailable" json:"is_available"`
	InventoryVariant *InventoryVariantShopVariant `protobuf:"bytes,17,opt,name=inventory_variant,json=inventoryVariant" json:"inventory_variant"`
	// @deprecated use stags instead
	Tags                 []string          `protobuf:"bytes,10,rep,name=tags" json:"tags"`
	Stags                []*Tag            `protobuf:"bytes,22,rep,name=stags" json:"stags"`
	Attributes           []*Attribute      `protobuf:"bytes,18,rep,name=attributes" json:"attributes"`
	Product              *ShopShortProduct `protobuf:"bytes,23,opt,name=product" json:"product"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ShopVariant) Reset()         { *m = ShopVariant{} }
func (m *ShopVariant) String() string { return proto.CompactTextString(m) }
func (*ShopVariant) ProtoMessage()    {}
func (*ShopVariant) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{30}
}

var xxx_messageInfo_ShopVariant proto.InternalMessageInfo

func (m *ShopVariant) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShopVariant) GetInfo() *EtopVariant {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *ShopVariant) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ShopVariant) GetEdCode() string {
	if m != nil {
		return m.EdCode
	}
	return ""
}

func (m *ShopVariant) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShopVariant) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ShopVariant) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *ShopVariant) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *ShopVariant) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *ShopVariant) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *ShopVariant) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *ShopVariant) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *ShopVariant) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *ShopVariant) GetIsAvailable() bool {
	if m != nil {
		return m.IsAvailable
	}
	return false
}

func (m *ShopVariant) GetInventoryVariant() *InventoryVariantShopVariant {
	if m != nil {
		return m.InventoryVariant
	}
	return nil
}

func (m *ShopVariant) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ShopVariant) GetStags() []*Tag {
	if m != nil {
		return m.Stags
	}
	return nil
}

func (m *ShopVariant) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *ShopVariant) GetProduct() *ShopShortProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

type InventoryVariantShopVariant struct {
	QuantityOnHand       int32    `protobuf:"varint,3,opt,name=quantity_on_hand,json=quantityOnHand" json:"quantity_on_hand"`
	QuantityPicked       int32    `protobuf:"varint,4,opt,name=quantity_picked,json=quantityPicked" json:"quantity_picked"`
	CostPrice            int32    `protobuf:"varint,5,opt,name=cost_price,json=costPrice" json:"cost_price"`
	Quantity             int32    `protobuf:"varint,6,opt,name=quantity" json:"quantity"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InventoryVariantShopVariant) Reset()         { *m = InventoryVariantShopVariant{} }
func (m *InventoryVariantShopVariant) String() string { return proto.CompactTextString(m) }
func (*InventoryVariantShopVariant) ProtoMessage()    {}
func (*InventoryVariantShopVariant) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{31}
}

var xxx_messageInfo_InventoryVariantShopVariant proto.InternalMessageInfo

func (m *InventoryVariantShopVariant) GetQuantityOnHand() int32 {
	if m != nil {
		return m.QuantityOnHand
	}
	return 0
}

func (m *InventoryVariantShopVariant) GetQuantityPicked() int32 {
	if m != nil {
		return m.QuantityPicked
	}
	return 0
}

func (m *InventoryVariantShopVariant) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *InventoryVariantShopVariant) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type ShopProduct struct {
	// @required
	Id   dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	Info *EtopProduct `protobuf:"bytes,8,opt,name=info" json:"info"`
	Code string       `protobuf:"bytes,2,opt,name=code" json:"code"`
	// @deprecated use code instead
	EdCode      string   `protobuf:"bytes,3,opt,name=ed_code,json=edCode" json:"ed_code"`
	Name        string   `protobuf:"bytes,4,opt,name=name" json:"name"`
	Description string   `protobuf:"bytes,5,opt,name=description" json:"description"`
	ShortDesc   string   `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml    string   `protobuf:"bytes,7,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls   []string `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	CategoryId  dot.ID   `protobuf:"varint,24,opt,name=category_id,json=categoryId" json:"category_id"`
	// @deprecated use stags instead
	Tags                 []string                  `protobuf:"bytes,10,rep,name=tags" json:"tags"`
	Stags                []*Tag                    `protobuf:"bytes,22,rep,name=stags" json:"stags"`
	Note                 string                    `protobuf:"bytes,11,opt,name=note" json:"note"`
	Status               status3.Status            `protobuf:"varint,12,opt,name=status,enum=status3.Status" json:"status"`
	IsAvailable          bool                      `protobuf:"varint,13,opt,name=is_available,json=isAvailable" json:"is_available"`
	ListPrice            int32                     `protobuf:"varint,14,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          int32                     `protobuf:"varint,15,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	CollectionIds        []dot.ID                  `protobuf:"varint,17,rep,name=collection_ids,json=collectionIds" json:"collection_ids"`
	Variants             []*ShopVariant            `protobuf:"bytes,18,rep,name=variants" json:"variants"`
	ProductSourceId      dot.ID                    `protobuf:"varint,19,opt,name=product_source_id,json=productSourceId" json:"product_source_id"`
	CreatedAt            dot.Time                  `protobuf:"bytes,21,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time                  `protobuf:"bytes,20,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ProductType          *product_type.ProductType `protobuf:"varint,23,opt,name=product_type,json=productType,enum=product_type.ProductType" json:"product_type"`
	MetaFields           []*common.MetaField       `protobuf:"bytes,26,rep,name=meta_fields,json=metaFields" json:"meta_fields"`
	BrandId              dot.ID                    `protobuf:"varint,27,opt,name=brand_id,json=brandId" json:"brand_id"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ShopProduct) Reset()         { *m = ShopProduct{} }
func (m *ShopProduct) String() string { return proto.CompactTextString(m) }
func (*ShopProduct) ProtoMessage()    {}
func (*ShopProduct) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{32}
}

var xxx_messageInfo_ShopProduct proto.InternalMessageInfo

func (m *ShopProduct) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShopProduct) GetInfo() *EtopProduct {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *ShopProduct) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ShopProduct) GetEdCode() string {
	if m != nil {
		return m.EdCode
	}
	return ""
}

func (m *ShopProduct) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShopProduct) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ShopProduct) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *ShopProduct) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *ShopProduct) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *ShopProduct) GetCategoryId() dot.ID {
	if m != nil {
		return m.CategoryId
	}
	return 0
}

func (m *ShopProduct) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *ShopProduct) GetStags() []*Tag {
	if m != nil {
		return m.Stags
	}
	return nil
}

func (m *ShopProduct) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *ShopProduct) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *ShopProduct) GetIsAvailable() bool {
	if m != nil {
		return m.IsAvailable
	}
	return false
}

func (m *ShopProduct) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *ShopProduct) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *ShopProduct) GetCollectionIds() []dot.ID {
	if m != nil {
		return m.CollectionIds
	}
	return nil
}

func (m *ShopProduct) GetVariants() []*ShopVariant {
	if m != nil {
		return m.Variants
	}
	return nil
}

func (m *ShopProduct) GetProductSourceId() dot.ID {
	if m != nil {
		return m.ProductSourceId
	}
	return 0
}

func (m *ShopProduct) GetProductType() product_type.ProductType {
	if m != nil && m.ProductType != nil {
		return *m.ProductType
	}
	return product_type.ProductType_unknown
}

func (m *ShopProduct) GetMetaFields() []*common.MetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

func (m *ShopProduct) GetBrandId() dot.ID {
	if m != nil {
		return m.BrandId
	}
	return 0
}

type ShopShortProduct struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopShortProduct) Reset()         { *m = ShopShortProduct{} }
func (m *ShopShortProduct) String() string { return proto.CompactTextString(m) }
func (*ShopShortProduct) ProtoMessage()    {}
func (*ShopShortProduct) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{33}
}

var xxx_messageInfo_ShopShortProduct proto.InternalMessageInfo

func (m *ShopShortProduct) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShopShortProduct) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ShopCollection struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	Description          string   `protobuf:"bytes,4,opt,name=description" json:"description"`
	DescHtml             string   `protobuf:"bytes,5,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ShortDesc            string   `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopCollection) Reset()         { *m = ShopCollection{} }
func (m *ShopCollection) String() string { return proto.CompactTextString(m) }
func (*ShopCollection) ProtoMessage()    {}
func (*ShopCollection) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{34}
}

var xxx_messageInfo_ShopCollection proto.InternalMessageInfo

func (m *ShopCollection) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShopCollection) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *ShopCollection) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShopCollection) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ShopCollection) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *ShopCollection) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

type GetVariantsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetVariantsRequest) Reset()         { *m = GetVariantsRequest{} }
func (m *GetVariantsRequest) String() string { return proto.CompactTextString(m) }
func (*GetVariantsRequest) ProtoMessage()    {}
func (*GetVariantsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{35}
}

var xxx_messageInfo_GetVariantsRequest proto.InternalMessageInfo

func (m *GetVariantsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetVariantsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetCategoriesRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetCategoriesRequest) Reset()         { *m = GetCategoriesRequest{} }
func (m *GetCategoriesRequest) String() string { return proto.CompactTextString(m) }
func (*GetCategoriesRequest) ProtoMessage()    {}
func (*GetCategoriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{36}
}

var xxx_messageInfo_GetCategoriesRequest proto.InternalMessageInfo

func (m *GetCategoriesRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCategoriesRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type ShopVariantsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Variants             []*ShopVariant   `protobuf:"bytes,2,rep,name=variants" json:"variants"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ShopVariantsResponse) Reset()         { *m = ShopVariantsResponse{} }
func (m *ShopVariantsResponse) String() string { return proto.CompactTextString(m) }
func (*ShopVariantsResponse) ProtoMessage()    {}
func (*ShopVariantsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{37}
}

var xxx_messageInfo_ShopVariantsResponse proto.InternalMessageInfo

func (m *ShopVariantsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShopVariantsResponse) GetVariants() []*ShopVariant {
	if m != nil {
		return m.Variants
	}
	return nil
}

type ShopProductsResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Products             []*ShopProduct   `protobuf:"bytes,2,rep,name=products" json:"products"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ShopProductsResponse) Reset()         { *m = ShopProductsResponse{} }
func (m *ShopProductsResponse) String() string { return proto.CompactTextString(m) }
func (*ShopProductsResponse) ProtoMessage()    {}
func (*ShopProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{38}
}

var xxx_messageInfo_ShopProductsResponse proto.InternalMessageInfo

func (m *ShopProductsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShopProductsResponse) GetProducts() []*ShopProduct {
	if m != nil {
		return m.Products
	}
	return nil
}

type ShopCategoriesResponse struct {
	Paging               *common.PageInfo `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Categories           []*ShopCategory  `protobuf:"bytes,2,rep,name=categories" json:"categories"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ShopCategoriesResponse) Reset()         { *m = ShopCategoriesResponse{} }
func (m *ShopCategoriesResponse) String() string { return proto.CompactTextString(m) }
func (*ShopCategoriesResponse) ProtoMessage()    {}
func (*ShopCategoriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{39}
}

var xxx_messageInfo_ShopCategoriesResponse proto.InternalMessageInfo

func (m *ShopCategoriesResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShopCategoriesResponse) GetCategories() []*ShopCategory {
	if m != nil {
		return m.Categories
	}
	return nil
}

type UpdateVariantRequest struct {
	// @required
	Id          dot.ID       `protobuf:"varint,1,opt,name=id" json:"id"`
	Name        *string      `protobuf:"bytes,2,opt,name=name" json:"name"`
	Note        *string      `protobuf:"bytes,3,opt,name=note" json:"note"`
	Code        *string      `protobuf:"bytes,11,opt,name=code" json:"code"`
	CostPrice   *int32       `protobuf:"varint,4,opt,name=cost_price,json=costPrice" json:"cost_price"`
	ListPrice   *int32       `protobuf:"varint,5,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice *int32       `protobuf:"varint,6,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	Description *string      `protobuf:"bytes,7,opt,name=description" json:"description"`
	ShortDesc   *string      `protobuf:"bytes,8,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml    *string      `protobuf:"bytes,9,opt,name=desc_html,json=descHtml" json:"desc_html"`
	Attributes  []*Attribute `protobuf:"bytes,10,rep,name=attributes" json:"attributes"`
	// deprecated
	Sku                  string   `protobuf:"bytes,16,opt,name=sku" json:"sku"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateVariantRequest) Reset()         { *m = UpdateVariantRequest{} }
func (m *UpdateVariantRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVariantRequest) ProtoMessage()    {}
func (*UpdateVariantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{40}
}

var xxx_messageInfo_UpdateVariantRequest proto.InternalMessageInfo

func (m *UpdateVariantRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateVariantRequest) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpdateVariantRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateVariantRequest) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *UpdateVariantRequest) GetCostPrice() int32 {
	if m != nil && m.CostPrice != nil {
		return *m.CostPrice
	}
	return 0
}

func (m *UpdateVariantRequest) GetListPrice() int32 {
	if m != nil && m.ListPrice != nil {
		return *m.ListPrice
	}
	return 0
}

func (m *UpdateVariantRequest) GetRetailPrice() int32 {
	if m != nil && m.RetailPrice != nil {
		return *m.RetailPrice
	}
	return 0
}

func (m *UpdateVariantRequest) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *UpdateVariantRequest) GetShortDesc() string {
	if m != nil && m.ShortDesc != nil {
		return *m.ShortDesc
	}
	return ""
}

func (m *UpdateVariantRequest) GetDescHtml() string {
	if m != nil && m.DescHtml != nil {
		return *m.DescHtml
	}
	return ""
}

func (m *UpdateVariantRequest) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *UpdateVariantRequest) GetSku() string {
	if m != nil {
		return m.Sku
	}
	return ""
}

type UpdateProductRequest struct {
	// @required
	Id                   dot.ID                    `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 *string                   `protobuf:"bytes,2,opt,name=name" json:"name"`
	Code                 *string                   `protobuf:"bytes,3,opt,name=code" json:"code"`
	Note                 *string                   `protobuf:"bytes,4,opt,name=note" json:"note"`
	Unit                 *string                   `protobuf:"bytes,5,opt,name=unit" json:"unit"`
	Description          *string                   `protobuf:"bytes,6,opt,name=description" json:"description"`
	ShortDesc            *string                   `protobuf:"bytes,7,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             *string                   `protobuf:"bytes,8,opt,name=desc_html,json=descHtml" json:"desc_html"`
	CostPrice            *int32                    `protobuf:"varint,9,opt,name=cost_price,json=costPrice" json:"cost_price"`
	ListPrice            *int32                    `protobuf:"varint,10,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          *int32                    `protobuf:"varint,11,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	ProductType          *product_type.ProductType `protobuf:"varint,12,opt,name=product_type,json=productType,enum=product_type.ProductType" json:"product_type"`
	MetaFields           *common.MetaField         `protobuf:"bytes,14,opt,name=meta_fields,json=metaFields" json:"meta_fields"`
	BrandId              *dot.ID                   `protobuf:"varint,15,opt,name=brand_id,json=brandId" json:"brand_id"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *UpdateProductRequest) Reset()         { *m = UpdateProductRequest{} }
func (m *UpdateProductRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductRequest) ProtoMessage()    {}
func (*UpdateProductRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{41}
}

var xxx_messageInfo_UpdateProductRequest proto.InternalMessageInfo

func (m *UpdateProductRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateProductRequest) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpdateProductRequest) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *UpdateProductRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateProductRequest) GetUnit() string {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return ""
}

func (m *UpdateProductRequest) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *UpdateProductRequest) GetShortDesc() string {
	if m != nil && m.ShortDesc != nil {
		return *m.ShortDesc
	}
	return ""
}

func (m *UpdateProductRequest) GetDescHtml() string {
	if m != nil && m.DescHtml != nil {
		return *m.DescHtml
	}
	return ""
}

func (m *UpdateProductRequest) GetCostPrice() int32 {
	if m != nil && m.CostPrice != nil {
		return *m.CostPrice
	}
	return 0
}

func (m *UpdateProductRequest) GetListPrice() int32 {
	if m != nil && m.ListPrice != nil {
		return *m.ListPrice
	}
	return 0
}

func (m *UpdateProductRequest) GetRetailPrice() int32 {
	if m != nil && m.RetailPrice != nil {
		return *m.RetailPrice
	}
	return 0
}

func (m *UpdateProductRequest) GetProductType() product_type.ProductType {
	if m != nil && m.ProductType != nil {
		return *m.ProductType
	}
	return product_type.ProductType_unknown
}

func (m *UpdateProductRequest) GetMetaFields() *common.MetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

func (m *UpdateProductRequest) GetBrandId() dot.ID {
	if m != nil && m.BrandId != nil {
		return *m.BrandId
	}
	return 0
}

type UpdateCategoryRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 *string  `protobuf:"bytes,2,opt,name=name" json:"name"`
	ParentId             dot.ID   `protobuf:"varint,3,opt,name=parent_id,json=parentId" json:"parent_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateCategoryRequest) Reset()         { *m = UpdateCategoryRequest{} }
func (m *UpdateCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCategoryRequest) ProtoMessage()    {}
func (*UpdateCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{42}
}

var xxx_messageInfo_UpdateCategoryRequest proto.InternalMessageInfo

func (m *UpdateCategoryRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateCategoryRequest) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UpdateCategoryRequest) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

type UpdateVariantsRequest struct {
	Updates              []*UpdateVariantRequest `protobuf:"bytes,1,rep,name=updates" json:"updates"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *UpdateVariantsRequest) Reset()         { *m = UpdateVariantsRequest{} }
func (m *UpdateVariantsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVariantsRequest) ProtoMessage()    {}
func (*UpdateVariantsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{43}
}

var xxx_messageInfo_UpdateVariantsRequest proto.InternalMessageInfo

func (m *UpdateVariantsRequest) GetUpdates() []*UpdateVariantRequest {
	if m != nil {
		return m.Updates
	}
	return nil
}

type UpdateProductsTagsRequest struct {
	// @required
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Adds                 []string `protobuf:"bytes,2,rep,name=adds" json:"adds"`
	Deletes              []string `protobuf:"bytes,3,rep,name=deletes" json:"deletes"`
	ReplaceAll           []string `protobuf:"bytes,4,rep,name=replace_all,json=replaceAll" json:"replace_all"`
	DeleteAll            bool     `protobuf:"varint,5,opt,name=delete_all,json=deleteAll" json:"delete_all"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductsTagsRequest) Reset()         { *m = UpdateProductsTagsRequest{} }
func (m *UpdateProductsTagsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductsTagsRequest) ProtoMessage()    {}
func (*UpdateProductsTagsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{44}
}

var xxx_messageInfo_UpdateProductsTagsRequest proto.InternalMessageInfo

func (m *UpdateProductsTagsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *UpdateProductsTagsRequest) GetAdds() []string {
	if m != nil {
		return m.Adds
	}
	return nil
}

func (m *UpdateProductsTagsRequest) GetDeletes() []string {
	if m != nil {
		return m.Deletes
	}
	return nil
}

func (m *UpdateProductsTagsRequest) GetReplaceAll() []string {
	if m != nil {
		return m.ReplaceAll
	}
	return nil
}

func (m *UpdateProductsTagsRequest) GetDeleteAll() bool {
	if m != nil {
		return m.DeleteAll
	}
	return false
}

type UpdateVariantsResponse struct {
	Variants             []*ShopVariant  `protobuf:"bytes,2,rep,name=variants" json:"variants"`
	Errors               []*common.Error `protobuf:"bytes,3,rep,name=errors" json:"errors"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateVariantsResponse) Reset()         { *m = UpdateVariantsResponse{} }
func (m *UpdateVariantsResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateVariantsResponse) ProtoMessage()    {}
func (*UpdateVariantsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{45}
}

var xxx_messageInfo_UpdateVariantsResponse proto.InternalMessageInfo

func (m *UpdateVariantsResponse) GetVariants() []*ShopVariant {
	if m != nil {
		return m.Variants
	}
	return nil
}

func (m *UpdateVariantsResponse) GetErrors() []*common.Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

type AddVariantsRequest struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Tags                 []string `protobuf:"bytes,2,rep,name=tags" json:"tags"`
	CollectionId         dot.ID   `protobuf:"varint,3,opt,name=collection_id,json=collectionId" json:"collection_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddVariantsRequest) Reset()         { *m = AddVariantsRequest{} }
func (m *AddVariantsRequest) String() string { return proto.CompactTextString(m) }
func (*AddVariantsRequest) ProtoMessage()    {}
func (*AddVariantsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{46}
}

var xxx_messageInfo_AddVariantsRequest proto.InternalMessageInfo

func (m *AddVariantsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *AddVariantsRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *AddVariantsRequest) GetCollectionId() dot.ID {
	if m != nil {
		return m.CollectionId
	}
	return 0
}

type AddVariantsResponse struct {
	Variants             []*ShopVariant  `protobuf:"bytes,1,rep,name=variants" json:"variants"`
	Errors               []*common.Error `protobuf:"bytes,2,rep,name=errors" json:"errors"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AddVariantsResponse) Reset()         { *m = AddVariantsResponse{} }
func (m *AddVariantsResponse) String() string { return proto.CompactTextString(m) }
func (*AddVariantsResponse) ProtoMessage()    {}
func (*AddVariantsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{47}
}

var xxx_messageInfo_AddVariantsResponse proto.InternalMessageInfo

func (m *AddVariantsResponse) GetVariants() []*ShopVariant {
	if m != nil {
		return m.Variants
	}
	return nil
}

func (m *AddVariantsResponse) GetErrors() []*common.Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

type RemoveVariantsRequest struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveVariantsRequest) Reset()         { *m = RemoveVariantsRequest{} }
func (m *RemoveVariantsRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveVariantsRequest) ProtoMessage()    {}
func (*RemoveVariantsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{48}
}

var xxx_messageInfo_RemoveVariantsRequest proto.InternalMessageInfo

func (m *RemoveVariantsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

type GetOrdersRequest struct {
	Paging               *common.Paging     `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter   `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	Mixed                *etop.MixedAccount `protobuf:"bytes,5,opt,name=mixed" json:"mixed"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetOrdersRequest) Reset()         { *m = GetOrdersRequest{} }
func (m *GetOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*GetOrdersRequest) ProtoMessage()    {}
func (*GetOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{49}
}

var xxx_messageInfo_GetOrdersRequest proto.InternalMessageInfo

func (m *GetOrdersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetOrdersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *GetOrdersRequest) GetMixed() *etop.MixedAccount {
	if m != nil {
		return m.Mixed
	}
	return nil
}

type UpdateOrdersStatusRequest struct {
	// @required
	Ids []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	// @required
	Confirm              *status3.Status `protobuf:"varint,2,opt,name=confirm,enum=status3.Status" json:"confirm"`
	CancelReason         string          `protobuf:"bytes,3,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	Status               status4.Status  `protobuf:"varint,4,opt,name=status,enum=status4.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateOrdersStatusRequest) Reset()         { *m = UpdateOrdersStatusRequest{} }
func (m *UpdateOrdersStatusRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateOrdersStatusRequest) ProtoMessage()    {}
func (*UpdateOrdersStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{50}
}

var xxx_messageInfo_UpdateOrdersStatusRequest proto.InternalMessageInfo

func (m *UpdateOrdersStatusRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *UpdateOrdersStatusRequest) GetConfirm() status3.Status {
	if m != nil && m.Confirm != nil {
		return *m.Confirm
	}
	return status3.Status_Z
}

func (m *UpdateOrdersStatusRequest) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *UpdateOrdersStatusRequest) GetStatus() status4.Status {
	if m != nil {
		return m.Status
	}
	return status4.Status_Z
}

type ConfirmOrderRequest struct {
	OrderId dot.ID `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	// enum ('create', 'create')
	AutoInventoryVoucher  *string  `protobuf:"bytes,2,opt,name=auto_inventory_voucher,json=autoInventoryVoucher" json:"auto_inventory_voucher"`
	AutoCreateFulfillment bool     `protobuf:"varint,3,opt,name=auto_create_fulfillment,json=autoCreateFulfillment" json:"auto_create_fulfillment"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *ConfirmOrderRequest) Reset()         { *m = ConfirmOrderRequest{} }
func (m *ConfirmOrderRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmOrderRequest) ProtoMessage()    {}
func (*ConfirmOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{51}
}

var xxx_messageInfo_ConfirmOrderRequest proto.InternalMessageInfo

func (m *ConfirmOrderRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *ConfirmOrderRequest) GetAutoInventoryVoucher() string {
	if m != nil && m.AutoInventoryVoucher != nil {
		return *m.AutoInventoryVoucher
	}
	return ""
}

func (m *ConfirmOrderRequest) GetAutoCreateFulfillment() bool {
	if m != nil {
		return m.AutoCreateFulfillment
	}
	return false
}

type OrderIDRequest struct {
	OrderId              dot.ID   `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderIDRequest) Reset()         { *m = OrderIDRequest{} }
func (m *OrderIDRequest) String() string { return proto.CompactTextString(m) }
func (*OrderIDRequest) ProtoMessage()    {}
func (*OrderIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{52}
}

var xxx_messageInfo_OrderIDRequest proto.InternalMessageInfo

func (m *OrderIDRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

type OrderIDsRequest struct {
	OrderIds             []dot.ID `protobuf:"varint,1,rep,name=order_ids,json=orderIds" json:"order_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderIDsRequest) Reset()         { *m = OrderIDsRequest{} }
func (m *OrderIDsRequest) String() string { return proto.CompactTextString(m) }
func (*OrderIDsRequest) ProtoMessage()    {}
func (*OrderIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{53}
}

var xxx_messageInfo_OrderIDsRequest proto.InternalMessageInfo

func (m *OrderIDsRequest) GetOrderIds() []dot.ID {
	if m != nil {
		return m.OrderIds
	}
	return nil
}

type CreateFulfillmentsForOrderRequest struct {
	OrderId              dot.ID   `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	VariantIds           []dot.ID `protobuf:"varint,3,rep,name=variant_ids,json=variantIds" json:"variant_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateFulfillmentsForOrderRequest) Reset()         { *m = CreateFulfillmentsForOrderRequest{} }
func (m *CreateFulfillmentsForOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CreateFulfillmentsForOrderRequest) ProtoMessage()    {}
func (*CreateFulfillmentsForOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{54}
}

var xxx_messageInfo_CreateFulfillmentsForOrderRequest proto.InternalMessageInfo

func (m *CreateFulfillmentsForOrderRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *CreateFulfillmentsForOrderRequest) GetVariantIds() []dot.ID {
	if m != nil {
		return m.VariantIds
	}
	return nil
}

type CancelOrderRequest struct {
	OrderId              dot.ID   `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	CancelReason         string   `protobuf:"bytes,2,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	AutoInventoryVoucher string   `protobuf:"bytes,3,opt,name=auto_inventory_voucher,json=autoInventoryVoucher" json:"auto_inventory_voucher"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelOrderRequest) Reset()         { *m = CancelOrderRequest{} }
func (m *CancelOrderRequest) String() string { return proto.CompactTextString(m) }
func (*CancelOrderRequest) ProtoMessage()    {}
func (*CancelOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{55}
}

var xxx_messageInfo_CancelOrderRequest proto.InternalMessageInfo

func (m *CancelOrderRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *CancelOrderRequest) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *CancelOrderRequest) GetAutoInventoryVoucher() string {
	if m != nil {
		return m.AutoInventoryVoucher
	}
	return ""
}

type CancelOrdersRequest struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason" json:"reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelOrdersRequest) Reset()         { *m = CancelOrdersRequest{} }
func (m *CancelOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*CancelOrdersRequest) ProtoMessage()    {}
func (*CancelOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{56}
}

var xxx_messageInfo_CancelOrdersRequest proto.InternalMessageInfo

func (m *CancelOrdersRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *CancelOrdersRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ProductSource struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	Type                 string         `protobuf:"bytes,2,opt,name=type" json:"type"`
	Name                 string         `protobuf:"bytes,3,opt,name=name" json:"name"`
	Status               status3.Status `protobuf:"varint,4,opt,name=status,enum=status3.Status" json:"status"`
	UpdatedAt            dot.Time       `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	CreatedAt            dot.Time       `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ProductSource) Reset()         { *m = ProductSource{} }
func (m *ProductSource) String() string { return proto.CompactTextString(m) }
func (*ProductSource) ProtoMessage()    {}
func (*ProductSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{57}
}

var xxx_messageInfo_ProductSource proto.InternalMessageInfo

func (m *ProductSource) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProductSource) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ProductSource) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProductSource) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

// deprecated
type CreateProductSourceRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateProductSourceRequest) Reset()         { *m = CreateProductSourceRequest{} }
func (m *CreateProductSourceRequest) String() string { return proto.CompactTextString(m) }
func (*CreateProductSourceRequest) ProtoMessage()    {}
func (*CreateProductSourceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{58}
}

var xxx_messageInfo_CreateProductSourceRequest proto.InternalMessageInfo

type ProductSourcesResponse struct {
	ProductSources       []*ProductSource `protobuf:"bytes,1,rep,name=product_sources,json=productSources" json:"product_sources"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ProductSourcesResponse) Reset()         { *m = ProductSourcesResponse{} }
func (m *ProductSourcesResponse) String() string { return proto.CompactTextString(m) }
func (*ProductSourcesResponse) ProtoMessage()    {}
func (*ProductSourcesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{59}
}

var xxx_messageInfo_ProductSourcesResponse proto.InternalMessageInfo

func (m *ProductSourcesResponse) GetProductSources() []*ProductSource {
	if m != nil {
		return m.ProductSources
	}
	return nil
}

type CreateCategoryRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	ParentId             dot.ID   `protobuf:"varint,2,opt,name=parent_id,json=parentId" json:"parent_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCategoryRequest) Reset()         { *m = CreateCategoryRequest{} }
func (m *CreateCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCategoryRequest) ProtoMessage()    {}
func (*CreateCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{60}
}

var xxx_messageInfo_CreateCategoryRequest proto.InternalMessageInfo

func (m *CreateCategoryRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateCategoryRequest) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

type CreateProductRequest struct {
	Code                 string                    `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string                    `protobuf:"bytes,2,opt,name=name" json:"name"`
	Unit                 string                    `protobuf:"bytes,3,opt,name=unit" json:"unit"`
	Note                 string                    `protobuf:"bytes,4,opt,name=note" json:"note"`
	Description          string                    `protobuf:"bytes,5,opt,name=description" json:"description"`
	ShortDesc            string                    `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             string                    `protobuf:"bytes,7,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls            []string                  `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	CostPrice            int32                     `protobuf:"varint,16,opt,name=cost_price,json=costPrice" json:"cost_price"`
	ListPrice            int32                     `protobuf:"varint,15,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          int32                     `protobuf:"varint,26,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	ProductType          *product_type.ProductType `protobuf:"varint,27,opt,name=product_type,json=productType,enum=product_type.ProductType" json:"product_type"`
	BrandId              dot.ID                    `protobuf:"varint,29,opt,name=brand_id,json=brandId" json:"brand_id"`
	MetaFields           []*common.MetaField       `protobuf:"bytes,28,rep,name=meta_fields,json=metaFields" json:"meta_fields"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CreateProductRequest) Reset()         { *m = CreateProductRequest{} }
func (m *CreateProductRequest) String() string { return proto.CompactTextString(m) }
func (*CreateProductRequest) ProtoMessage()    {}
func (*CreateProductRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{61}
}

var xxx_messageInfo_CreateProductRequest proto.InternalMessageInfo

func (m *CreateProductRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CreateProductRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateProductRequest) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *CreateProductRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateProductRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateProductRequest) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *CreateProductRequest) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *CreateProductRequest) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *CreateProductRequest) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *CreateProductRequest) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *CreateProductRequest) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *CreateProductRequest) GetProductType() product_type.ProductType {
	if m != nil && m.ProductType != nil {
		return *m.ProductType
	}
	return product_type.ProductType_unknown
}

func (m *CreateProductRequest) GetBrandId() dot.ID {
	if m != nil {
		return m.BrandId
	}
	return 0
}

func (m *CreateProductRequest) GetMetaFields() []*common.MetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

type CreateVariantRequest struct {
	Code                 string       `protobuf:"bytes,1,opt,name=code" json:"code"`
	Name                 string       `protobuf:"bytes,2,opt,name=name" json:"name"`
	ProductId            dot.ID       `protobuf:"varint,3,opt,name=product_id,json=productId" json:"product_id"`
	Note                 string       `protobuf:"bytes,4,opt,name=note" json:"note"`
	Description          string       `protobuf:"bytes,5,opt,name=description" json:"description"`
	ShortDesc            string       `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml             string       `protobuf:"bytes,7,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls            []string     `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	Attributes           []*Attribute `protobuf:"bytes,14,rep,name=attributes" json:"attributes"`
	CostPrice            int32        `protobuf:"varint,16,opt,name=cost_price,json=costPrice" json:"cost_price"`
	ListPrice            int32        `protobuf:"varint,15,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice          int32        `protobuf:"varint,26,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CreateVariantRequest) Reset()         { *m = CreateVariantRequest{} }
func (m *CreateVariantRequest) String() string { return proto.CompactTextString(m) }
func (*CreateVariantRequest) ProtoMessage()    {}
func (*CreateVariantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{62}
}

var xxx_messageInfo_CreateVariantRequest proto.InternalMessageInfo

func (m *CreateVariantRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CreateVariantRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateVariantRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CreateVariantRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateVariantRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateVariantRequest) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *CreateVariantRequest) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *CreateVariantRequest) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *CreateVariantRequest) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *CreateVariantRequest) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *CreateVariantRequest) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *CreateVariantRequest) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

type DeprecatedCreateVariantRequest struct {
	// required
	ProductSourceId dot.ID `protobuf:"varint,23,opt,name=product_source_id,json=productSourceId" json:"product_source_id"`
	ProductId       dot.ID `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string `protobuf:"bytes,2,opt,name=product_name,json=productName" json:"product_name"`
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name              string         `protobuf:"bytes,3,opt,name=name" json:"name"`
	Description       string         `protobuf:"bytes,5,opt,name=description" json:"description"`
	ShortDesc         string         `protobuf:"bytes,6,opt,name=short_desc,json=shortDesc" json:"short_desc"`
	DescHtml          string         `protobuf:"bytes,7,opt,name=desc_html,json=descHtml" json:"desc_html"`
	ImageUrls         []string       `protobuf:"bytes,9,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	Tags              []string       `protobuf:"bytes,10,rep,name=tags" json:"tags"`
	Status            status3.Status `protobuf:"varint,12,opt,name=status,enum=status3.Status" json:"status"`
	CostPrice         int32          `protobuf:"varint,16,opt,name=cost_price,json=costPrice" json:"cost_price"`
	ListPrice         int32          `protobuf:"varint,15,opt,name=list_price,json=listPrice" json:"list_price"`
	RetailPrice       int32          `protobuf:"varint,26,opt,name=retail_price,json=retailPrice" json:"retail_price"`
	Code              string         `protobuf:"bytes,19,opt,name=code" json:"code"`
	QuantityAvailable int32          `protobuf:"varint,20,opt,name=quantity_available,json=quantityAvailable" json:"quantity_available"`
	QuantityOnHand    int32          `protobuf:"varint,21,opt,name=quantity_on_hand,json=quantityOnHand" json:"quantity_on_hand"`
	QuantityReserved  int32          `protobuf:"varint,22,opt,name=quantity_reserved,json=quantityReserved" json:"quantity_reserved"`
	Attributes        []*Attribute   `protobuf:"bytes,24,rep,name=attributes" json:"attributes"`
	Unit              string         `protobuf:"bytes,25,opt,name=unit" json:"unit"`
	// deprecated: use code instead
	Sku                  string   `protobuf:"bytes,18,opt,name=sku" json:"sku"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeprecatedCreateVariantRequest) Reset()         { *m = DeprecatedCreateVariantRequest{} }
func (m *DeprecatedCreateVariantRequest) String() string { return proto.CompactTextString(m) }
func (*DeprecatedCreateVariantRequest) ProtoMessage()    {}
func (*DeprecatedCreateVariantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{63}
}

var xxx_messageInfo_DeprecatedCreateVariantRequest proto.InternalMessageInfo

func (m *DeprecatedCreateVariantRequest) GetProductSourceId() dot.ID {
	if m != nil {
		return m.ProductSourceId
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetShortDesc() string {
	if m != nil {
		return m.ShortDesc
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetDescHtml() string {
	if m != nil {
		return m.DescHtml
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

func (m *DeprecatedCreateVariantRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *DeprecatedCreateVariantRequest) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *DeprecatedCreateVariantRequest) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetListPrice() int32 {
	if m != nil {
		return m.ListPrice
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetRetailPrice() int32 {
	if m != nil {
		return m.RetailPrice
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetQuantityAvailable() int32 {
	if m != nil {
		return m.QuantityAvailable
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetQuantityOnHand() int32 {
	if m != nil {
		return m.QuantityOnHand
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetQuantityReserved() int32 {
	if m != nil {
		return m.QuantityReserved
	}
	return 0
}

func (m *DeprecatedCreateVariantRequest) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *DeprecatedCreateVariantRequest) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *DeprecatedCreateVariantRequest) GetSku() string {
	if m != nil {
		return m.Sku
	}
	return ""
}

type ConnectProductSourceResquest struct {
	ProductSourceId      dot.ID   `protobuf:"varint,1,opt,name=product_source_id,json=productSourceId" json:"product_source_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectProductSourceResquest) Reset()         { *m = ConnectProductSourceResquest{} }
func (m *ConnectProductSourceResquest) String() string { return proto.CompactTextString(m) }
func (*ConnectProductSourceResquest) ProtoMessage()    {}
func (*ConnectProductSourceResquest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{64}
}

var xxx_messageInfo_ConnectProductSourceResquest proto.InternalMessageInfo

func (m *ConnectProductSourceResquest) GetProductSourceId() dot.ID {
	if m != nil {
		return m.ProductSourceId
	}
	return 0
}

// deprecated
type CreatePSCategoryRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	ParentId             dot.ID   `protobuf:"varint,2,opt,name=parent_id,json=parentId" json:"parent_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreatePSCategoryRequest) Reset()         { *m = CreatePSCategoryRequest{} }
func (m *CreatePSCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePSCategoryRequest) ProtoMessage()    {}
func (*CreatePSCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{65}
}

var xxx_messageInfo_CreatePSCategoryRequest proto.InternalMessageInfo

func (m *CreatePSCategoryRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreatePSCategoryRequest) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

type UpdateProductsPSCategoryRequest struct {
	CategoryId           dot.ID   `protobuf:"varint,1,opt,name=category_id,json=categoryId" json:"category_id"`
	ProductIds           []dot.ID `protobuf:"varint,2,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductsPSCategoryRequest) Reset()         { *m = UpdateProductsPSCategoryRequest{} }
func (m *UpdateProductsPSCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductsPSCategoryRequest) ProtoMessage()    {}
func (*UpdateProductsPSCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{66}
}

var xxx_messageInfo_UpdateProductsPSCategoryRequest proto.InternalMessageInfo

func (m *UpdateProductsPSCategoryRequest) GetCategoryId() dot.ID {
	if m != nil {
		return m.CategoryId
	}
	return 0
}

func (m *UpdateProductsPSCategoryRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type UpdateProductsCollectionResponse struct {
	Updated              int32           `protobuf:"varint,1,opt,name=updated" json:"updated"`
	Errors               []*common.Error `protobuf:"bytes,2,rep,name=errors" json:"errors"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateProductsCollectionResponse) Reset()         { *m = UpdateProductsCollectionResponse{} }
func (m *UpdateProductsCollectionResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateProductsCollectionResponse) ProtoMessage()    {}
func (*UpdateProductsCollectionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{67}
}

var xxx_messageInfo_UpdateProductsCollectionResponse proto.InternalMessageInfo

func (m *UpdateProductsCollectionResponse) GetUpdated() int32 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func (m *UpdateProductsCollectionResponse) GetErrors() []*common.Error {
	if m != nil {
		return m.Errors
	}
	return nil
}

type UpdateProductSourceCategoryRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ParentId             dot.ID   `protobuf:"varint,2,opt,name=parent_id,json=parentId" json:"parent_id"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductSourceCategoryRequest) Reset()         { *m = UpdateProductSourceCategoryRequest{} }
func (m *UpdateProductSourceCategoryRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductSourceCategoryRequest) ProtoMessage()    {}
func (*UpdateProductSourceCategoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{68}
}

var xxx_messageInfo_UpdateProductSourceCategoryRequest proto.InternalMessageInfo

func (m *UpdateProductSourceCategoryRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateProductSourceCategoryRequest) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

func (m *UpdateProductSourceCategoryRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// deprecated
type GetProductSourceCategoriesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProductSourceCategoriesRequest) Reset()         { *m = GetProductSourceCategoriesRequest{} }
func (m *GetProductSourceCategoriesRequest) String() string { return proto.CompactTextString(m) }
func (*GetProductSourceCategoriesRequest) ProtoMessage()    {}
func (*GetProductSourceCategoriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{69}
}

var xxx_messageInfo_GetProductSourceCategoriesRequest proto.InternalMessageInfo

type GetFulfillmentsRequest struct {
	Paging               *common.Paging     `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter   `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	Mixed                *etop.MixedAccount `protobuf:"bytes,3,opt,name=mixed" json:"mixed"`
	OrderId              dot.ID             `protobuf:"varint,4,opt,name=order_id,json=orderId" json:"order_id"`
	Status               *status3.Status    `protobuf:"varint,5,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetFulfillmentsRequest) Reset()         { *m = GetFulfillmentsRequest{} }
func (m *GetFulfillmentsRequest) String() string { return proto.CompactTextString(m) }
func (*GetFulfillmentsRequest) ProtoMessage()    {}
func (*GetFulfillmentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{70}
}

var xxx_messageInfo_GetFulfillmentsRequest proto.InternalMessageInfo

func (m *GetFulfillmentsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetFulfillmentsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *GetFulfillmentsRequest) GetMixed() *etop.MixedAccount {
	if m != nil {
		return m.Mixed
	}
	return nil
}

func (m *GetFulfillmentsRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *GetFulfillmentsRequest) GetStatus() status3.Status {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return status3.Status_Z
}

type GetFulfillmentHistoryRequest struct {
	Paging               *common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	All                  bool           `protobuf:"varint,2,opt,name=all" json:"all"`
	Id                   dot.ID         `protobuf:"varint,3,opt,name=id" json:"id"`
	OrderId              dot.ID         `protobuf:"varint,4,opt,name=order_id,json=orderId" json:"order_id"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetFulfillmentHistoryRequest) Reset()         { *m = GetFulfillmentHistoryRequest{} }
func (m *GetFulfillmentHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*GetFulfillmentHistoryRequest) ProtoMessage()    {}
func (*GetFulfillmentHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{71}
}

var xxx_messageInfo_GetFulfillmentHistoryRequest proto.InternalMessageInfo

func (m *GetFulfillmentHistoryRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetFulfillmentHistoryRequest) GetAll() bool {
	if m != nil {
		return m.All
	}
	return false
}

func (m *GetFulfillmentHistoryRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetFulfillmentHistoryRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

type GetBalanceShopResponse struct {
	Amount               int32    `protobuf:"varint,1,opt,name=amount" json:"amount"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBalanceShopResponse) Reset()         { *m = GetBalanceShopResponse{} }
func (m *GetBalanceShopResponse) String() string { return proto.CompactTextString(m) }
func (*GetBalanceShopResponse) ProtoMessage()    {}
func (*GetBalanceShopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{72}
}

var xxx_messageInfo_GetBalanceShopResponse proto.InternalMessageInfo

func (m *GetBalanceShopResponse) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type GetMoneyTransactionsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetMoneyTransactionsRequest) Reset()         { *m = GetMoneyTransactionsRequest{} }
func (m *GetMoneyTransactionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetMoneyTransactionsRequest) ProtoMessage()    {}
func (*GetMoneyTransactionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{73}
}

var xxx_messageInfo_GetMoneyTransactionsRequest proto.InternalMessageInfo

func (m *GetMoneyTransactionsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetMoneyTransactionsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetPublicFulfillmentRequest struct {
	// @Required
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPublicFulfillmentRequest) Reset()         { *m = GetPublicFulfillmentRequest{} }
func (m *GetPublicFulfillmentRequest) String() string { return proto.CompactTextString(m) }
func (*GetPublicFulfillmentRequest) ProtoMessage()    {}
func (*GetPublicFulfillmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{74}
}

var xxx_messageInfo_GetPublicFulfillmentRequest proto.InternalMessageInfo

func (m *GetPublicFulfillmentRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type UpdateFulfillmentsShippingStateRequest struct {
	// Only support for manual order
	Ids []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	// @required
	ShippingState        shipping.State `protobuf:"varint,3,opt,name=shipping_state,json=shippingState,enum=shipping.State" json:"shipping_state"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UpdateFulfillmentsShippingStateRequest) Reset() {
	*m = UpdateFulfillmentsShippingStateRequest{}
}
func (m *UpdateFulfillmentsShippingStateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateFulfillmentsShippingStateRequest) ProtoMessage()    {}
func (*UpdateFulfillmentsShippingStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{75}
}

var xxx_messageInfo_UpdateFulfillmentsShippingStateRequest proto.InternalMessageInfo

func (m *UpdateFulfillmentsShippingStateRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *UpdateFulfillmentsShippingStateRequest) GetShippingState() shipping.State {
	if m != nil {
		return m.ShippingState
	}
	return shipping.State_default
}

type UpdateOrderPaymentStatusRequest struct {
	OrderId              dot.ID          `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	Status               *status3.Status `protobuf:"varint,2,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UpdateOrderPaymentStatusRequest) Reset()         { *m = UpdateOrderPaymentStatusRequest{} }
func (m *UpdateOrderPaymentStatusRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateOrderPaymentStatusRequest) ProtoMessage()    {}
func (*UpdateOrderPaymentStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{76}
}

var xxx_messageInfo_UpdateOrderPaymentStatusRequest proto.InternalMessageInfo

func (m *UpdateOrderPaymentStatusRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *UpdateOrderPaymentStatusRequest) GetStatus() status3.Status {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return status3.Status_Z
}

type SummarizeFulfillmentsRequest struct {
	DateFrom             string   `protobuf:"bytes,1,opt,name=date_from,json=dateFrom" json:"date_from"`
	DateTo               string   `protobuf:"bytes,2,opt,name=date_to,json=dateTo" json:"date_to"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SummarizeFulfillmentsRequest) Reset()         { *m = SummarizeFulfillmentsRequest{} }
func (m *SummarizeFulfillmentsRequest) String() string { return proto.CompactTextString(m) }
func (*SummarizeFulfillmentsRequest) ProtoMessage()    {}
func (*SummarizeFulfillmentsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{77}
}

var xxx_messageInfo_SummarizeFulfillmentsRequest proto.InternalMessageInfo

func (m *SummarizeFulfillmentsRequest) GetDateFrom() string {
	if m != nil {
		return m.DateFrom
	}
	return ""
}

func (m *SummarizeFulfillmentsRequest) GetDateTo() string {
	if m != nil {
		return m.DateTo
	}
	return ""
}

type SummarizeFulfillmentsResponse struct {
	Tables               []*SummaryTable `protobuf:"bytes,1,rep,name=tables" json:"tables"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SummarizeFulfillmentsResponse) Reset()         { *m = SummarizeFulfillmentsResponse{} }
func (m *SummarizeFulfillmentsResponse) String() string { return proto.CompactTextString(m) }
func (*SummarizeFulfillmentsResponse) ProtoMessage()    {}
func (*SummarizeFulfillmentsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{78}
}

var xxx_messageInfo_SummarizeFulfillmentsResponse proto.InternalMessageInfo

func (m *SummarizeFulfillmentsResponse) GetTables() []*SummaryTable {
	if m != nil {
		return m.Tables
	}
	return nil
}

type SummarizePOSResponse struct {
	Tables               []*SummaryTable `protobuf:"bytes,1,rep,name=tables" json:"tables"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SummarizePOSResponse) Reset()         { *m = SummarizePOSResponse{} }
func (m *SummarizePOSResponse) String() string { return proto.CompactTextString(m) }
func (*SummarizePOSResponse) ProtoMessage()    {}
func (*SummarizePOSResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{79}
}

var xxx_messageInfo_SummarizePOSResponse proto.InternalMessageInfo

func (m *SummarizePOSResponse) GetTables() []*SummaryTable {
	if m != nil {
		return m.Tables
	}
	return nil
}

type SummarizePOSRequest struct {
	DateFrom             string   `protobuf:"bytes,1,opt,name=date_from,json=dateFrom" json:"date_from"`
	DateTo               string   `protobuf:"bytes,2,opt,name=date_to,json=dateTo" json:"date_to"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SummarizePOSRequest) Reset()         { *m = SummarizePOSRequest{} }
func (m *SummarizePOSRequest) String() string { return proto.CompactTextString(m) }
func (*SummarizePOSRequest) ProtoMessage()    {}
func (*SummarizePOSRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{80}
}

var xxx_messageInfo_SummarizePOSRequest proto.InternalMessageInfo

func (m *SummarizePOSRequest) GetDateFrom() string {
	if m != nil {
		return m.DateFrom
	}
	return ""
}

func (m *SummarizePOSRequest) GetDateTo() string {
	if m != nil {
		return m.DateTo
	}
	return ""
}

type SummaryTable struct {
	Label                string          `protobuf:"bytes,1,opt,name=label" json:"label"`
	Tags                 []string        `protobuf:"bytes,2,rep,name=tags" json:"tags"`
	Columns              []SummaryColRow `protobuf:"bytes,3,rep,name=columns" json:"columns"`
	Rows                 []SummaryColRow `protobuf:"bytes,4,rep,name=rows" json:"rows"`
	Data                 []SummaryItem   `protobuf:"bytes,5,rep,name=data" json:"data"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SummaryTable) Reset()         { *m = SummaryTable{} }
func (m *SummaryTable) String() string { return proto.CompactTextString(m) }
func (*SummaryTable) ProtoMessage()    {}
func (*SummaryTable) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{81}
}

var xxx_messageInfo_SummaryTable proto.InternalMessageInfo

func (m *SummaryTable) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *SummaryTable) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *SummaryTable) GetColumns() []SummaryColRow {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *SummaryTable) GetRows() []SummaryColRow {
	if m != nil {
		return m.Rows
	}
	return nil
}

func (m *SummaryTable) GetData() []SummaryItem {
	if m != nil {
		return m.Data
	}
	return nil
}

type SummaryColRow struct {
	Label                string   `protobuf:"bytes,1,opt,name=label" json:"label"`
	Spec                 string   `protobuf:"bytes,2,opt,name=spec" json:"spec"`
	Unit                 string   `protobuf:"bytes,3,opt,name=unit" json:"unit"`
	Indent               int32    `protobuf:"varint,4,opt,name=indent" json:"indent"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SummaryColRow) Reset()         { *m = SummaryColRow{} }
func (m *SummaryColRow) String() string { return proto.CompactTextString(m) }
func (*SummaryColRow) ProtoMessage()    {}
func (*SummaryColRow) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{82}
}

var xxx_messageInfo_SummaryColRow proto.InternalMessageInfo

func (m *SummaryColRow) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *SummaryColRow) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *SummaryColRow) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *SummaryColRow) GetIndent() int32 {
	if m != nil {
		return m.Indent
	}
	return 0
}

type SummaryItem struct {
	Label                string   `protobuf:"bytes,1,opt,name=label" json:"label"`
	Spec                 string   `protobuf:"bytes,2,opt,name=spec" json:"spec"`
	Value                int32    `protobuf:"varint,3,opt,name=value" json:"value"`
	Unit                 string   `protobuf:"bytes,4,opt,name=unit" json:"unit"`
	ImageUrls            []string `protobuf:"bytes,5,rep,name=image_urls,json=imageUrls" json:"image_urls"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SummaryItem) Reset()         { *m = SummaryItem{} }
func (m *SummaryItem) String() string { return proto.CompactTextString(m) }
func (*SummaryItem) ProtoMessage()    {}
func (*SummaryItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{83}
}

var xxx_messageInfo_SummaryItem proto.InternalMessageInfo

func (m *SummaryItem) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *SummaryItem) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *SummaryItem) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *SummaryItem) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *SummaryItem) GetImageUrls() []string {
	if m != nil {
		return m.ImageUrls
	}
	return nil
}

type ImportProductsResponse struct {
	Data                 *spreadsheet.SpreadsheetData `protobuf:"bytes,1,opt,name=data" json:"data"`
	ImportErrors         []*common.Error              `protobuf:"bytes,3,rep,name=import_errors,json=importErrors" json:"import_errors"`
	CellErrors           []*common.Error              `protobuf:"bytes,4,rep,name=cell_errors,json=cellErrors" json:"cell_errors"`
	ImportId             dot.ID                       `protobuf:"varint,5,opt,name=import_id,json=importId" json:"import_id"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *ImportProductsResponse) Reset()         { *m = ImportProductsResponse{} }
func (m *ImportProductsResponse) String() string { return proto.CompactTextString(m) }
func (*ImportProductsResponse) ProtoMessage()    {}
func (*ImportProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{84}
}

var xxx_messageInfo_ImportProductsResponse proto.InternalMessageInfo

func (m *ImportProductsResponse) GetData() *spreadsheet.SpreadsheetData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ImportProductsResponse) GetImportErrors() []*common.Error {
	if m != nil {
		return m.ImportErrors
	}
	return nil
}

func (m *ImportProductsResponse) GetCellErrors() []*common.Error {
	if m != nil {
		return m.CellErrors
	}
	return nil
}

func (m *ImportProductsResponse) GetImportId() dot.ID {
	if m != nil {
		return m.ImportId
	}
	return 0
}

type CalcBalanceShopResponse struct {
	Balance              int32    `protobuf:"varint,1,opt,name=balance" json:"balance"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CalcBalanceShopResponse) Reset()         { *m = CalcBalanceShopResponse{} }
func (m *CalcBalanceShopResponse) String() string { return proto.CompactTextString(m) }
func (*CalcBalanceShopResponse) ProtoMessage()    {}
func (*CalcBalanceShopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{85}
}

var xxx_messageInfo_CalcBalanceShopResponse proto.InternalMessageInfo

func (m *CalcBalanceShopResponse) GetBalance() int32 {
	if m != nil {
		return m.Balance
	}
	return 0
}

type RequestExportRequest struct {
	ExportType string           `protobuf:"bytes,1,opt,name=export_type,json=exportType" json:"export_type"`
	Filters    []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	DateFrom   string           `protobuf:"bytes,3,opt,name=date_from,json=dateFrom" json:"date_from"`
	DateTo     string           `protobuf:"bytes,4,opt,name=date_to,json=dateTo" json:"date_to"`
	// Accept '\t', ',' or ';'. Default to ','.
	Delimiter string `protobuf:"bytes,5,opt,name=delimiter" json:"delimiter"`
	// For exporting csv compatible with Excel
	ExcelCompatibleMode bool `protobuf:"varint,6,opt,name=excel_compatible_mode,json=excelCompatibleMode" json:"excel_compatible_mode"`
	// Export specific ids
	Ids                  []dot.ID `protobuf:"varint,7,rep,name=ids" json:"ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RequestExportRequest) Reset()         { *m = RequestExportRequest{} }
func (m *RequestExportRequest) String() string { return proto.CompactTextString(m) }
func (*RequestExportRequest) ProtoMessage()    {}
func (*RequestExportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{86}
}

var xxx_messageInfo_RequestExportRequest proto.InternalMessageInfo

func (m *RequestExportRequest) GetExportType() string {
	if m != nil {
		return m.ExportType
	}
	return ""
}

func (m *RequestExportRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *RequestExportRequest) GetDateFrom() string {
	if m != nil {
		return m.DateFrom
	}
	return ""
}

func (m *RequestExportRequest) GetDateTo() string {
	if m != nil {
		return m.DateTo
	}
	return ""
}

func (m *RequestExportRequest) GetDelimiter() string {
	if m != nil {
		return m.Delimiter
	}
	return ""
}

func (m *RequestExportRequest) GetExcelCompatibleMode() bool {
	if m != nil {
		return m.ExcelCompatibleMode
	}
	return false
}

func (m *RequestExportRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

type RequestExportResponse struct {
	Id                   string         `protobuf:"bytes,1,opt,name=id" json:"id"`
	Filename             string         `protobuf:"bytes,2,opt,name=filename" json:"filename"`
	ExportType           string         `protobuf:"bytes,3,opt,name=export_type,json=exportType" json:"export_type"`
	Status               status4.Status `protobuf:"varint,4,opt,name=status,enum=status4.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RequestExportResponse) Reset()         { *m = RequestExportResponse{} }
func (m *RequestExportResponse) String() string { return proto.CompactTextString(m) }
func (*RequestExportResponse) ProtoMessage()    {}
func (*RequestExportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{87}
}

var xxx_messageInfo_RequestExportResponse proto.InternalMessageInfo

func (m *RequestExportResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RequestExportResponse) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *RequestExportResponse) GetExportType() string {
	if m != nil {
		return m.ExportType
	}
	return ""
}

func (m *RequestExportResponse) GetStatus() status4.Status {
	if m != nil {
		return m.Status
	}
	return status4.Status_Z
}

type GetExportsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetExportsRequest) Reset()         { *m = GetExportsRequest{} }
func (m *GetExportsRequest) String() string { return proto.CompactTextString(m) }
func (*GetExportsRequest) ProtoMessage()    {}
func (*GetExportsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{88}
}

var xxx_messageInfo_GetExportsRequest proto.InternalMessageInfo

type GetExportsResponse struct {
	ExportItems          []*ExportItem `protobuf:"bytes,1,rep,name=export_items,json=exportItems" json:"export_items"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetExportsResponse) Reset()         { *m = GetExportsResponse{} }
func (m *GetExportsResponse) String() string { return proto.CompactTextString(m) }
func (*GetExportsResponse) ProtoMessage()    {}
func (*GetExportsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{89}
}

var xxx_messageInfo_GetExportsResponse proto.InternalMessageInfo

func (m *GetExportsResponse) GetExportItems() []*ExportItem {
	if m != nil {
		return m.ExportItems
	}
	return nil
}

type ExportItem struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id"`
	Filename string `protobuf:"bytes,2,opt,name=filename" json:"filename"`
	// example: shop/fulfillments, admin/orders
	ExportType           string          `protobuf:"bytes,3,opt,name=export_type,json=exportType" json:"export_type"`
	DownloadUrl          string          `protobuf:"bytes,4,opt,name=download_url,json=downloadUrl" json:"download_url"`
	AccountId            dot.ID          `protobuf:"varint,5,opt,name=account_id,json=accountId" json:"account_id"`
	UserId               dot.ID          `protobuf:"varint,6,opt,name=user_id,json=userId" json:"user_id"`
	CreatedAt            dot.Time        `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	DeletedAt            dot.Time        `protobuf:"bytes,8,opt,name=deleted_at,json=deletedAt" json:"deleted_at"`
	RequestQuery         string          `protobuf:"bytes,9,opt,name=request_query,json=requestQuery" json:"request_query"`
	MimeType             string          `protobuf:"bytes,10,opt,name=mime_type,json=mimeType" json:"mime_type"`
	Status               status4.Status  `protobuf:"varint,11,opt,name=status,enum=status4.Status" json:"status"`
	ExportErrors         []*common.Error `protobuf:"bytes,12,rep,name=export_errors,json=exportErrors" json:"export_errors"`
	Error                *common.Error   `protobuf:"bytes,13,opt,name=error" json:"error"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ExportItem) Reset()         { *m = ExportItem{} }
func (m *ExportItem) String() string { return proto.CompactTextString(m) }
func (*ExportItem) ProtoMessage()    {}
func (*ExportItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{90}
}

var xxx_messageInfo_ExportItem proto.InternalMessageInfo

func (m *ExportItem) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExportItem) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *ExportItem) GetExportType() string {
	if m != nil {
		return m.ExportType
	}
	return ""
}

func (m *ExportItem) GetDownloadUrl() string {
	if m != nil {
		return m.DownloadUrl
	}
	return ""
}

func (m *ExportItem) GetAccountId() dot.ID {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *ExportItem) GetUserId() dot.ID {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *ExportItem) GetRequestQuery() string {
	if m != nil {
		return m.RequestQuery
	}
	return ""
}

func (m *ExportItem) GetMimeType() string {
	if m != nil {
		return m.MimeType
	}
	return ""
}

func (m *ExportItem) GetStatus() status4.Status {
	if m != nil {
		return m.Status
	}
	return status4.Status_Z
}

func (m *ExportItem) GetExportErrors() []*common.Error {
	if m != nil {
		return m.ExportErrors
	}
	return nil
}

func (m *ExportItem) GetError() *common.Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type GetExportsStatusRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetExportsStatusRequest) Reset()         { *m = GetExportsStatusRequest{} }
func (m *GetExportsStatusRequest) String() string { return proto.CompactTextString(m) }
func (*GetExportsStatusRequest) ProtoMessage()    {}
func (*GetExportsStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{91}
}

var xxx_messageInfo_GetExportsStatusRequest proto.InternalMessageInfo

type ExportStatusItem struct {
	Id                   string        `protobuf:"bytes,1,opt,name=id" json:"id"`
	ProgressMax          int32         `protobuf:"varint,2,opt,name=progress_max,json=progressMax" json:"progress_max"`
	ProgressValue        int32         `protobuf:"varint,3,opt,name=progress_value,json=progressValue" json:"progress_value"`
	ProgressError        int32         `protobuf:"varint,4,opt,name=progress_error,json=progressError" json:"progress_error"`
	Error                *common.Error `protobuf:"bytes,5,opt,name=error" json:"error"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ExportStatusItem) Reset()         { *m = ExportStatusItem{} }
func (m *ExportStatusItem) String() string { return proto.CompactTextString(m) }
func (*ExportStatusItem) ProtoMessage()    {}
func (*ExportStatusItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{92}
}

var xxx_messageInfo_ExportStatusItem proto.InternalMessageInfo

func (m *ExportStatusItem) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ExportStatusItem) GetProgressMax() int32 {
	if m != nil {
		return m.ProgressMax
	}
	return 0
}

func (m *ExportStatusItem) GetProgressValue() int32 {
	if m != nil {
		return m.ProgressValue
	}
	return 0
}

func (m *ExportStatusItem) GetProgressError() int32 {
	if m != nil {
		return m.ProgressError
	}
	return 0
}

func (m *ExportStatusItem) GetError() *common.Error {
	if m != nil {
		return m.Error
	}
	return nil
}

type AuthorizePartnerRequest struct {
	PartnerId            dot.ID   `protobuf:"varint,1,opt,name=partner_id,json=partnerId" json:"partner_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizePartnerRequest) Reset()         { *m = AuthorizePartnerRequest{} }
func (m *AuthorizePartnerRequest) String() string { return proto.CompactTextString(m) }
func (*AuthorizePartnerRequest) ProtoMessage()    {}
func (*AuthorizePartnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{93}
}

var xxx_messageInfo_AuthorizePartnerRequest proto.InternalMessageInfo

func (m *AuthorizePartnerRequest) GetPartnerId() dot.ID {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

type GetPartnersResponse struct {
	Partners             []*etop.PublicAccountInfo `protobuf:"bytes,1,rep,name=partners" json:"partners"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *GetPartnersResponse) Reset()         { *m = GetPartnersResponse{} }
func (m *GetPartnersResponse) String() string { return proto.CompactTextString(m) }
func (*GetPartnersResponse) ProtoMessage()    {}
func (*GetPartnersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{94}
}

var xxx_messageInfo_GetPartnersResponse proto.InternalMessageInfo

func (m *GetPartnersResponse) GetPartners() []*etop.PublicAccountInfo {
	if m != nil {
		return m.Partners
	}
	return nil
}

type AuthorizedPartnerResponse struct {
	Partner              *etop.PublicAccountInfo `protobuf:"bytes,1,opt,name=partner" json:"partner"`
	RedirectUrl          string                  `protobuf:"bytes,2,opt,name=redirect_url,json=redirectUrl" json:"redirect_url"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *AuthorizedPartnerResponse) Reset()         { *m = AuthorizedPartnerResponse{} }
func (m *AuthorizedPartnerResponse) String() string { return proto.CompactTextString(m) }
func (*AuthorizedPartnerResponse) ProtoMessage()    {}
func (*AuthorizedPartnerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{95}
}

var xxx_messageInfo_AuthorizedPartnerResponse proto.InternalMessageInfo

func (m *AuthorizedPartnerResponse) GetPartner() *etop.PublicAccountInfo {
	if m != nil {
		return m.Partner
	}
	return nil
}

func (m *AuthorizedPartnerResponse) GetRedirectUrl() string {
	if m != nil {
		return m.RedirectUrl
	}
	return ""
}

type GetAuthorizedPartnersResponse struct {
	Partners             []*AuthorizedPartnerResponse `protobuf:"bytes,1,rep,name=partners" json:"partners"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *GetAuthorizedPartnersResponse) Reset()         { *m = GetAuthorizedPartnersResponse{} }
func (m *GetAuthorizedPartnersResponse) String() string { return proto.CompactTextString(m) }
func (*GetAuthorizedPartnersResponse) ProtoMessage()    {}
func (*GetAuthorizedPartnersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{96}
}

var xxx_messageInfo_GetAuthorizedPartnersResponse proto.InternalMessageInfo

func (m *GetAuthorizedPartnersResponse) GetPartners() []*AuthorizedPartnerResponse {
	if m != nil {
		return m.Partners
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
	return fileDescriptor_1e3a4043938e9392, []int{97}
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

type UpdateVariantImagesRequest struct {
	// @required
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Adds                 []string `protobuf:"bytes,2,rep,name=adds" json:"adds"`
	Deletes              []string `protobuf:"bytes,3,rep,name=deletes" json:"deletes"`
	ReplaceAll           []string `protobuf:"bytes,4,rep,name=replace_all,json=replaceAll" json:"replace_all"`
	DeleteAll            bool     `protobuf:"varint,5,opt,name=delete_all,json=deleteAll" json:"delete_all"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateVariantImagesRequest) Reset()         { *m = UpdateVariantImagesRequest{} }
func (m *UpdateVariantImagesRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVariantImagesRequest) ProtoMessage()    {}
func (*UpdateVariantImagesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{98}
}

var xxx_messageInfo_UpdateVariantImagesRequest proto.InternalMessageInfo

func (m *UpdateVariantImagesRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateVariantImagesRequest) GetAdds() []string {
	if m != nil {
		return m.Adds
	}
	return nil
}

func (m *UpdateVariantImagesRequest) GetDeletes() []string {
	if m != nil {
		return m.Deletes
	}
	return nil
}

func (m *UpdateVariantImagesRequest) GetReplaceAll() []string {
	if m != nil {
		return m.ReplaceAll
	}
	return nil
}

func (m *UpdateVariantImagesRequest) GetDeleteAll() bool {
	if m != nil {
		return m.DeleteAll
	}
	return false
}

type UpdateProductMetaFieldsRequest struct {
	// @required
	Id                   dot.ID              `protobuf:"varint,1,opt,name=id" json:"id"`
	MetaFields           []*common.MetaField `protobuf:"bytes,2,rep,name=meta_fields,json=metaFields" json:"meta_fields"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *UpdateProductMetaFieldsRequest) Reset()         { *m = UpdateProductMetaFieldsRequest{} }
func (m *UpdateProductMetaFieldsRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductMetaFieldsRequest) ProtoMessage()    {}
func (*UpdateProductMetaFieldsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{99}
}

var xxx_messageInfo_UpdateProductMetaFieldsRequest proto.InternalMessageInfo

func (m *UpdateProductMetaFieldsRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateProductMetaFieldsRequest) GetMetaFields() []*common.MetaField {
	if m != nil {
		return m.MetaFields
	}
	return nil
}

type CategoriesResponse struct {
	Categories           []*Category `protobuf:"bytes,1,rep,name=categories" json:"categories"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *CategoriesResponse) Reset()         { *m = CategoriesResponse{} }
func (m *CategoriesResponse) String() string { return proto.CompactTextString(m) }
func (*CategoriesResponse) ProtoMessage()    {}
func (*CategoriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{100}
}

var xxx_messageInfo_CategoriesResponse proto.InternalMessageInfo

func (m *CategoriesResponse) GetCategories() []*Category {
	if m != nil {
		return m.Categories
	}
	return nil
}

type Category struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	ProductSourceId      dot.ID   `protobuf:"varint,3,opt,name=product_source_id,json=productSourceId" json:"product_source_id"`
	ParentId             dot.ID   `protobuf:"varint,6,opt,name=parent_id,json=parentId" json:"parent_id"`
	ShopId               dot.ID   `protobuf:"varint,7,opt,name=shop_id,json=shopId" json:"shop_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Category) Reset()         { *m = Category{} }
func (m *Category) String() string { return proto.CompactTextString(m) }
func (*Category) ProtoMessage()    {}
func (*Category) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{101}
}

var xxx_messageInfo_Category proto.InternalMessageInfo

func (m *Category) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Category) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Category) GetProductSourceId() dot.ID {
	if m != nil {
		return m.ProductSourceId
	}
	return 0
}

func (m *Category) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

func (m *Category) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

type Tag struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Label                string   `protobuf:"bytes,2,opt,name=label" json:"label"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{102}
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Tag) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

type ExternalAccountAhamove struct {
	Id               dot.ID `protobuf:"varint,1,opt,name=id" json:"id"`
	Phone            string `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Name             string `protobuf:"bytes,3,opt,name=name" json:"name"`
	ExternalVerified bool   `protobuf:"varint,4,opt,name=external_verified,json=externalVerified" json:"external_verified"`
	//    optional string external_token = 5 [(gogoproto.nullable) = false];
	ExternalCreatedAt    dot.Time `protobuf:"bytes,6,opt,name=external_created_at,json=externalCreatedAt" json:"external_created_at"`
	CreatedAt            dot.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	LastSendVerifyAt     dot.Time `protobuf:"bytes,9,opt,name=last_send_verify_at,json=lastSendVerifyAt" json:"last_send_verify_at"`
	ExternalTicketId     string   `protobuf:"bytes,10,opt,name=external_ticket_id,json=externalTicketId" json:"external_ticket_id"`
	IdCardFrontImg       string   `protobuf:"bytes,11,opt,name=id_card_front_img,json=idCardFrontImg" json:"id_card_front_img"`
	IdCardBackImg        string   `protobuf:"bytes,12,opt,name=id_card_back_img,json=idCardBackImg" json:"id_card_back_img"`
	PortraitImg          string   `protobuf:"bytes,13,opt,name=portrait_img,json=portraitImg" json:"portrait_img"`
	UploadedAt           dot.Time `protobuf:"bytes,14,opt,name=uploaded_at,json=uploadedAt" json:"uploaded_at"`
	FanpageUrl           string   `protobuf:"bytes,15,opt,name=fanpage_url,json=fanpageUrl" json:"fanpage_url"`
	WebsiteUrl           string   `protobuf:"bytes,16,opt,name=website_url,json=websiteUrl" json:"website_url"`
	CompanyImgs          []string `protobuf:"bytes,17,rep,name=company_imgs,json=companyImgs" json:"company_imgs"`
	BusinessLicenseImgs  []string `protobuf:"bytes,18,rep,name=business_license_imgs,json=businessLicenseImgs" json:"business_license_imgs"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExternalAccountAhamove) Reset()         { *m = ExternalAccountAhamove{} }
func (m *ExternalAccountAhamove) String() string { return proto.CompactTextString(m) }
func (*ExternalAccountAhamove) ProtoMessage()    {}
func (*ExternalAccountAhamove) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{103}
}

var xxx_messageInfo_ExternalAccountAhamove proto.InternalMessageInfo

func (m *ExternalAccountAhamove) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ExternalAccountAhamove) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ExternalAccountAhamove) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ExternalAccountAhamove) GetExternalVerified() bool {
	if m != nil {
		return m.ExternalVerified
	}
	return false
}

func (m *ExternalAccountAhamove) GetExternalTicketId() string {
	if m != nil {
		return m.ExternalTicketId
	}
	return ""
}

func (m *ExternalAccountAhamove) GetIdCardFrontImg() string {
	if m != nil {
		return m.IdCardFrontImg
	}
	return ""
}

func (m *ExternalAccountAhamove) GetIdCardBackImg() string {
	if m != nil {
		return m.IdCardBackImg
	}
	return ""
}

func (m *ExternalAccountAhamove) GetPortraitImg() string {
	if m != nil {
		return m.PortraitImg
	}
	return ""
}

func (m *ExternalAccountAhamove) GetFanpageUrl() string {
	if m != nil {
		return m.FanpageUrl
	}
	return ""
}

func (m *ExternalAccountAhamove) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *ExternalAccountAhamove) GetCompanyImgs() []string {
	if m != nil {
		return m.CompanyImgs
	}
	return nil
}

func (m *ExternalAccountAhamove) GetBusinessLicenseImgs() []string {
	if m != nil {
		return m.BusinessLicenseImgs
	}
	return nil
}

type UpdateXAccountAhamoveVerificationRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	IdCardFrontImg       string   `protobuf:"bytes,2,opt,name=id_card_front_img,json=idCardFrontImg" json:"id_card_front_img"`
	IdCardBackImg        string   `protobuf:"bytes,3,opt,name=id_card_back_img,json=idCardBackImg" json:"id_card_back_img"`
	PortraitImg          string   `protobuf:"bytes,4,opt,name=portrait_img,json=portraitImg" json:"portrait_img"`
	FanpageUrl           string   `protobuf:"bytes,5,opt,name=fanpage_url,json=fanpageUrl" json:"fanpage_url"`
	WebsiteUrl           string   `protobuf:"bytes,6,opt,name=website_url,json=websiteUrl" json:"website_url"`
	CompanyImgs          []string `protobuf:"bytes,7,rep,name=company_imgs,json=companyImgs" json:"company_imgs"`
	BusinessLicenseImgs  []string `protobuf:"bytes,8,rep,name=business_license_imgs,json=businessLicenseImgs" json:"business_license_imgs"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateXAccountAhamoveVerificationRequest) Reset() {
	*m = UpdateXAccountAhamoveVerificationRequest{}
}
func (m *UpdateXAccountAhamoveVerificationRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateXAccountAhamoveVerificationRequest) ProtoMessage()    {}
func (*UpdateXAccountAhamoveVerificationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{104}
}

var xxx_messageInfo_UpdateXAccountAhamoveVerificationRequest proto.InternalMessageInfo

func (m *UpdateXAccountAhamoveVerificationRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetIdCardFrontImg() string {
	if m != nil {
		return m.IdCardFrontImg
	}
	return ""
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetIdCardBackImg() string {
	if m != nil {
		return m.IdCardBackImg
	}
	return ""
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetPortraitImg() string {
	if m != nil {
		return m.PortraitImg
	}
	return ""
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetFanpageUrl() string {
	if m != nil {
		return m.FanpageUrl
	}
	return ""
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetWebsiteUrl() string {
	if m != nil {
		return m.WebsiteUrl
	}
	return ""
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetCompanyImgs() []string {
	if m != nil {
		return m.CompanyImgs
	}
	return nil
}

func (m *UpdateXAccountAhamoveVerificationRequest) GetBusinessLicenseImgs() []string {
	if m != nil {
		return m.BusinessLicenseImgs
	}
	return nil
}

type ExternalAccountHaravan struct {
	Id                                dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId                            dot.ID   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	Subdomain                         string   `protobuf:"bytes,3,opt,name=subdomain" json:"subdomain"`
	ExternalCarrierServiceId          int32    `protobuf:"varint,8,opt,name=external_carrier_service_id,json=externalCarrierServiceId" json:"external_carrier_service_id"`
	ExternalConnectedCarrierServiceAt dot.Time `protobuf:"bytes,9,opt,name=external_connected_carrier_service_at,json=externalConnectedCarrierServiceAt" json:"external_connected_carrier_service_at"`
	ExpiresAt                         dot.Time `protobuf:"bytes,5,opt,name=expires_at,json=expiresAt" json:"expires_at"`
	CreatedAt                         dot.Time `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                         dot.Time `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral              struct{} `json:"-"`
	XXX_sizecache                     int32    `json:"-"`
}

func (m *ExternalAccountHaravan) Reset()         { *m = ExternalAccountHaravan{} }
func (m *ExternalAccountHaravan) String() string { return proto.CompactTextString(m) }
func (*ExternalAccountHaravan) ProtoMessage()    {}
func (*ExternalAccountHaravan) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{105}
}

var xxx_messageInfo_ExternalAccountHaravan proto.InternalMessageInfo

func (m *ExternalAccountHaravan) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ExternalAccountHaravan) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *ExternalAccountHaravan) GetSubdomain() string {
	if m != nil {
		return m.Subdomain
	}
	return ""
}

func (m *ExternalAccountHaravan) GetExternalCarrierServiceId() int32 {
	if m != nil {
		return m.ExternalCarrierServiceId
	}
	return 0
}

type ExternalAccountHaravanRequest struct {
	// @required
	Subdomain string `protobuf:"bytes,1,opt,name=subdomain" json:"subdomain"`
	// @required OAuth code
	Code string `protobuf:"bytes,2,opt,name=code" json:"code"`
	// @required
	RedirectUri          string   `protobuf:"bytes,3,opt,name=redirect_uri,json=redirectUri" json:"redirect_uri"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExternalAccountHaravanRequest) Reset()         { *m = ExternalAccountHaravanRequest{} }
func (m *ExternalAccountHaravanRequest) String() string { return proto.CompactTextString(m) }
func (*ExternalAccountHaravanRequest) ProtoMessage()    {}
func (*ExternalAccountHaravanRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{106}
}

var xxx_messageInfo_ExternalAccountHaravanRequest proto.InternalMessageInfo

func (m *ExternalAccountHaravanRequest) GetSubdomain() string {
	if m != nil {
		return m.Subdomain
	}
	return ""
}

func (m *ExternalAccountHaravanRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *ExternalAccountHaravanRequest) GetRedirectUri() string {
	if m != nil {
		return m.RedirectUri
	}
	return ""
}

type CustomerLiability struct {
	TotalOrders          int32    `protobuf:"varint,1,opt,name=total_orders,json=totalOrders" json:"total_orders"`
	TotalAmount          int      `protobuf:"varint,2,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	ReceivedAmount       int      `protobuf:"varint,3,opt,name=received_amount,json=receivedAmount" json:"received_amount"`
	Liability            int      `protobuf:"varint,4,opt,name=liability" json:"liability"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CustomerLiability) Reset()         { *m = CustomerLiability{} }
func (m *CustomerLiability) String() string { return proto.CompactTextString(m) }
func (*CustomerLiability) ProtoMessage()    {}
func (*CustomerLiability) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{107}
}

var xxx_messageInfo_CustomerLiability proto.InternalMessageInfo

func (m *CustomerLiability) GetTotalOrders() int32 {
	if m != nil {
		return m.TotalOrders
	}
	return 0
}

type Customer struct {
	Id                   dot.ID             `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID             `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	FullName             string             `protobuf:"bytes,3,opt,name=full_name,json=fullName" json:"full_name"`
	Code                 string             `protobuf:"bytes,13,opt,name=code" json:"code"`
	Note                 string             `protobuf:"bytes,4,opt,name=note" json:"note"`
	Phone                string             `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Email                string             `protobuf:"bytes,6,opt,name=email" json:"email"`
	Gender               string             `protobuf:"bytes,7,opt,name=gender" json:"gender"`
	Type                 string             `protobuf:"bytes,8,opt,name=type" json:"type"`
	Birthday             string             `protobuf:"bytes,12,opt,name=birthday" json:"birthday"`
	CreatedAt            dot.Time           `protobuf:"bytes,9,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time           `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	Status               status3.Status     `protobuf:"varint,11,opt,name=status,enum=status3.Status" json:"status"`
	GroupIds             []dot.ID           `protobuf:"varint,14,rep,name=group_ids,json=groupIds" json:"group_ids"`
	Liability            *CustomerLiability `protobuf:"bytes,15,opt,name=liability" json:"liability"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Customer) Reset()         { *m = Customer{} }
func (m *Customer) String() string { return proto.CompactTextString(m) }
func (*Customer) ProtoMessage()    {}
func (*Customer) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{108}
}

var xxx_messageInfo_Customer proto.InternalMessageInfo

func (m *Customer) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Customer) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Customer) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Customer) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Customer) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Customer) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Customer) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Customer) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

func (m *Customer) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Customer) GetBirthday() string {
	if m != nil {
		return m.Birthday
	}
	return ""
}

func (m *Customer) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Customer) GetGroupIds() []dot.ID {
	if m != nil {
		return m.GroupIds
	}
	return nil
}

func (m *Customer) GetLiability() *CustomerLiability {
	if m != nil {
		return m.Liability
	}
	return nil
}

type CreateCustomerRequest struct {
	// @required
	FullName string `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Gender   string `protobuf:"bytes,12,opt,name=gender" json:"gender"`
	Birthday string `protobuf:"bytes,13,opt,name=birthday" json:"birthday"`
	// enum ('individual', 'organization')
	Type string `protobuf:"bytes,14,opt,name=type" json:"type"`
	Note string `protobuf:"bytes,2,opt,name=note" json:"note"`
	// @required
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCustomerRequest) Reset()         { *m = CreateCustomerRequest{} }
func (m *CreateCustomerRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCustomerRequest) ProtoMessage()    {}
func (*CreateCustomerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{109}
}

var xxx_messageInfo_CreateCustomerRequest proto.InternalMessageInfo

func (m *CreateCustomerRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateCustomerRequest) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

func (m *CreateCustomerRequest) GetBirthday() string {
	if m != nil {
		return m.Birthday
	}
	return ""
}

func (m *CreateCustomerRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *CreateCustomerRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateCustomerRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CreateCustomerRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UpdateCustomerRequest struct {
	Id       dot.ID  `protobuf:"varint,5,opt,name=id" json:"id"`
	FullName *string `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Gender   *string `protobuf:"bytes,12,opt,name=gender" json:"gender"`
	Birthday *string `protobuf:"bytes,13,opt,name=birthday" json:"birthday"`
	// enum ('individual', 'organization','independent')
	Type                 *string  `protobuf:"bytes,14,opt,name=type" json:"type"`
	Note                 *string  `protobuf:"bytes,2,opt,name=note" json:"note"`
	Phone                *string  `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                *string  `protobuf:"bytes,4,opt,name=email" json:"email"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateCustomerRequest) Reset()         { *m = UpdateCustomerRequest{} }
func (m *UpdateCustomerRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCustomerRequest) ProtoMessage()    {}
func (*UpdateCustomerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{110}
}

var xxx_messageInfo_UpdateCustomerRequest proto.InternalMessageInfo

func (m *UpdateCustomerRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateCustomerRequest) GetFullName() string {
	if m != nil && m.FullName != nil {
		return *m.FullName
	}
	return ""
}

func (m *UpdateCustomerRequest) GetGender() string {
	if m != nil && m.Gender != nil {
		return *m.Gender
	}
	return ""
}

func (m *UpdateCustomerRequest) GetBirthday() string {
	if m != nil && m.Birthday != nil {
		return *m.Birthday
	}
	return ""
}

func (m *UpdateCustomerRequest) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *UpdateCustomerRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateCustomerRequest) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

func (m *UpdateCustomerRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

type GetCustomersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetCustomersRequest) Reset()         { *m = GetCustomersRequest{} }
func (m *GetCustomersRequest) String() string { return proto.CompactTextString(m) }
func (*GetCustomersRequest) ProtoMessage()    {}
func (*GetCustomersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{111}
}

var xxx_messageInfo_GetCustomersRequest proto.InternalMessageInfo

func (m *GetCustomersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCustomersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type CustomersResponse struct {
	Customers            []*Customer      `protobuf:"bytes,1,rep,name=customers" json:"customers"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CustomersResponse) Reset()         { *m = CustomersResponse{} }
func (m *CustomersResponse) String() string { return proto.CompactTextString(m) }
func (*CustomersResponse) ProtoMessage()    {}
func (*CustomersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{112}
}

var xxx_messageInfo_CustomersResponse proto.InternalMessageInfo

func (m *CustomersResponse) GetCustomers() []*Customer {
	if m != nil {
		return m.Customers
	}
	return nil
}

func (m *CustomersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type SetCustomersStatusRequest struct {
	Ids                  []dot.ID       `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Status               status3.Status `protobuf:"varint,2,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SetCustomersStatusRequest) Reset()         { *m = SetCustomersStatusRequest{} }
func (m *SetCustomersStatusRequest) String() string { return proto.CompactTextString(m) }
func (*SetCustomersStatusRequest) ProtoMessage()    {}
func (*SetCustomersStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{113}
}

var xxx_messageInfo_SetCustomersStatusRequest proto.InternalMessageInfo

func (m *SetCustomersStatusRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *SetCustomersStatusRequest) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

type CustomerDetailsResponse struct {
	Customer             *Customer                 `protobuf:"bytes,1,opt,name=customer" json:"customer"`
	SummaryItems         []*IndependentSummaryItem `protobuf:"bytes,2,rep,name=summary_items,json=summaryItems" json:"summary_items"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *CustomerDetailsResponse) Reset()         { *m = CustomerDetailsResponse{} }
func (m *CustomerDetailsResponse) String() string { return proto.CompactTextString(m) }
func (*CustomerDetailsResponse) ProtoMessage()    {}
func (*CustomerDetailsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{114}
}

var xxx_messageInfo_CustomerDetailsResponse proto.InternalMessageInfo

func (m *CustomerDetailsResponse) GetCustomer() *Customer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *CustomerDetailsResponse) GetSummaryItems() []*IndependentSummaryItem {
	if m != nil {
		return m.SummaryItems
	}
	return nil
}

type IndependentSummaryItem struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Label                string   `protobuf:"bytes,2,opt,name=label" json:"label"`
	Spec                 string   `protobuf:"bytes,3,opt,name=spec" json:"spec"`
	Value                int32    `protobuf:"varint,4,opt,name=value" json:"value"`
	Unit                 string   `protobuf:"bytes,5,opt,name=unit" json:"unit"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndependentSummaryItem) Reset()         { *m = IndependentSummaryItem{} }
func (m *IndependentSummaryItem) String() string { return proto.CompactTextString(m) }
func (*IndependentSummaryItem) ProtoMessage()    {}
func (*IndependentSummaryItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{115}
}

var xxx_messageInfo_IndependentSummaryItem proto.InternalMessageInfo

func (m *IndependentSummaryItem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *IndependentSummaryItem) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *IndependentSummaryItem) GetSpec() string {
	if m != nil {
		return m.Spec
	}
	return ""
}

func (m *IndependentSummaryItem) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *IndependentSummaryItem) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

type GetCustomerAddressesRequest struct {
	CustomerId           dot.ID   `protobuf:"varint,1,opt,name=customer_id,json=customerId" json:"customer_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCustomerAddressesRequest) Reset()         { *m = GetCustomerAddressesRequest{} }
func (m *GetCustomerAddressesRequest) String() string { return proto.CompactTextString(m) }
func (*GetCustomerAddressesRequest) ProtoMessage()    {}
func (*GetCustomerAddressesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{116}
}

var xxx_messageInfo_GetCustomerAddressesRequest proto.InternalMessageInfo

func (m *GetCustomerAddressesRequest) GetCustomerId() dot.ID {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

type CustomerAddress struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Province             string            `protobuf:"bytes,2,opt,name=province" json:"province"`
	ProvinceCode         string            `protobuf:"bytes,6,opt,name=province_code,json=provinceCode" json:"province_code"`
	District             string            `protobuf:"bytes,3,opt,name=district" json:"district"`
	DistrictCode         string            `protobuf:"bytes,7,opt,name=district_code,json=districtCode" json:"district_code"`
	Ward                 string            `protobuf:"bytes,4,opt,name=ward" json:"ward"`
	WardCode             string            `protobuf:"bytes,8,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1             string            `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2             string            `protobuf:"bytes,9,opt,name=address2" json:"address2"`
	Country              string            `protobuf:"bytes,11,opt,name=country" json:"country"`
	FullName             string            `protobuf:"bytes,12,opt,name=full_name,json=fullName" json:"full_name"`
	Company              string            `protobuf:"bytes,10,opt,name=company" json:"company"`
	Phone                string            `protobuf:"bytes,15,opt,name=phone" json:"phone"`
	Email                string            `protobuf:"bytes,16,opt,name=email" json:"email"`
	Position             string            `protobuf:"bytes,17,opt,name=position" json:"position"`
	Coordinates          *etop.Coordinates `protobuf:"bytes,20,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CustomerAddress) Reset()         { *m = CustomerAddress{} }
func (m *CustomerAddress) String() string { return proto.CompactTextString(m) }
func (*CustomerAddress) ProtoMessage()    {}
func (*CustomerAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{117}
}

var xxx_messageInfo_CustomerAddress proto.InternalMessageInfo

func (m *CustomerAddress) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CustomerAddress) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *CustomerAddress) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *CustomerAddress) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *CustomerAddress) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *CustomerAddress) GetWard() string {
	if m != nil {
		return m.Ward
	}
	return ""
}

func (m *CustomerAddress) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *CustomerAddress) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *CustomerAddress) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *CustomerAddress) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *CustomerAddress) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CustomerAddress) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *CustomerAddress) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CustomerAddress) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CustomerAddress) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *CustomerAddress) GetCoordinates() *etop.Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type CreateCustomerAddressRequest struct {
	CustomerId           dot.ID            `protobuf:"varint,1,opt,name=customer_id,json=customerId" json:"customer_id"`
	ProvinceCode         string            `protobuf:"bytes,2,opt,name=province_code,json=provinceCode" json:"province_code"`
	DistrictCode         string            `protobuf:"bytes,3,opt,name=district_code,json=districtCode" json:"district_code"`
	WardCode             string            `protobuf:"bytes,4,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1             string            `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2             string            `protobuf:"bytes,6,opt,name=address2" json:"address2"`
	Country              string            `protobuf:"bytes,7,opt,name=country" json:"country"`
	FullName             string            `protobuf:"bytes,8,opt,name=full_name,json=fullName" json:"full_name"`
	Company              string            `protobuf:"bytes,12,opt,name=company" json:"company"`
	Phone                string            `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	Email                string            `protobuf:"bytes,10,opt,name=email" json:"email"`
	Position             string            `protobuf:"bytes,11,opt,name=position" json:"position"`
	Coordinates          *etop.Coordinates `protobuf:"bytes,18,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CreateCustomerAddressRequest) Reset()         { *m = CreateCustomerAddressRequest{} }
func (m *CreateCustomerAddressRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCustomerAddressRequest) ProtoMessage()    {}
func (*CreateCustomerAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{118}
}

var xxx_messageInfo_CreateCustomerAddressRequest proto.InternalMessageInfo

func (m *CreateCustomerAddressRequest) GetCustomerId() dot.ID {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

func (m *CreateCustomerAddressRequest) GetProvinceCode() string {
	if m != nil {
		return m.ProvinceCode
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetDistrictCode() string {
	if m != nil {
		return m.DistrictCode
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetWardCode() string {
	if m != nil {
		return m.WardCode
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetAddress1() string {
	if m != nil {
		return m.Address1
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetCompany() string {
	if m != nil {
		return m.Company
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetPosition() string {
	if m != nil {
		return m.Position
	}
	return ""
}

func (m *CreateCustomerAddressRequest) GetCoordinates() *etop.Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type UpdateCustomerAddressRequest struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	ProvinceCode         *string           `protobuf:"bytes,2,opt,name=province_code,json=provinceCode" json:"province_code"`
	DistrictCode         *string           `protobuf:"bytes,3,opt,name=district_code,json=districtCode" json:"district_code"`
	WardCode             *string           `protobuf:"bytes,4,opt,name=ward_code,json=wardCode" json:"ward_code"`
	Address1             *string           `protobuf:"bytes,5,opt,name=address1" json:"address1"`
	Address2             *string           `protobuf:"bytes,6,opt,name=address2" json:"address2"`
	Country              *string           `protobuf:"bytes,7,opt,name=country" json:"country"`
	FullName             *string           `protobuf:"bytes,8,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                *string           `protobuf:"bytes,9,opt,name=phone" json:"phone"`
	Email                *string           `protobuf:"bytes,10,opt,name=email" json:"email"`
	Position             *string           `protobuf:"bytes,11,opt,name=position" json:"position"`
	Company              *string           `protobuf:"bytes,12,opt,name=company" json:"company"`
	Coordinates          *etop.Coordinates `protobuf:"bytes,18,opt,name=coordinates" json:"coordinates"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateCustomerAddressRequest) Reset()         { *m = UpdateCustomerAddressRequest{} }
func (m *UpdateCustomerAddressRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCustomerAddressRequest) ProtoMessage()    {}
func (*UpdateCustomerAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{119}
}

var xxx_messageInfo_UpdateCustomerAddressRequest proto.InternalMessageInfo

func (m *UpdateCustomerAddressRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateCustomerAddressRequest) GetProvinceCode() string {
	if m != nil && m.ProvinceCode != nil {
		return *m.ProvinceCode
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetDistrictCode() string {
	if m != nil && m.DistrictCode != nil {
		return *m.DistrictCode
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetWardCode() string {
	if m != nil && m.WardCode != nil {
		return *m.WardCode
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetAddress1() string {
	if m != nil && m.Address1 != nil {
		return *m.Address1
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetAddress2() string {
	if m != nil && m.Address2 != nil {
		return *m.Address2
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetCountry() string {
	if m != nil && m.Country != nil {
		return *m.Country
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetFullName() string {
	if m != nil && m.FullName != nil {
		return *m.FullName
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetPosition() string {
	if m != nil && m.Position != nil {
		return *m.Position
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetCompany() string {
	if m != nil && m.Company != nil {
		return *m.Company
	}
	return ""
}

func (m *UpdateCustomerAddressRequest) GetCoordinates() *etop.Coordinates {
	if m != nil {
		return m.Coordinates
	}
	return nil
}

type CustomerAddressesResponse struct {
	Addresses            []*CustomerAddress `protobuf:"bytes,1,rep,name=addresses" json:"addresses"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CustomerAddressesResponse) Reset()         { *m = CustomerAddressesResponse{} }
func (m *CustomerAddressesResponse) String() string { return proto.CompactTextString(m) }
func (*CustomerAddressesResponse) ProtoMessage()    {}
func (*CustomerAddressesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{120}
}

var xxx_messageInfo_CustomerAddressesResponse proto.InternalMessageInfo

func (m *CustomerAddressesResponse) GetAddresses() []*CustomerAddress {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type UpdateProductStatusRequest struct {
	Ids                  []dot.ID       `protobuf:"varint,1,rep,name=ids" json:"ids"`
	Status               status3.Status `protobuf:"varint,3,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UpdateProductStatusRequest) Reset()         { *m = UpdateProductStatusRequest{} }
func (m *UpdateProductStatusRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProductStatusRequest) ProtoMessage()    {}
func (*UpdateProductStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{121}
}

var xxx_messageInfo_UpdateProductStatusRequest proto.InternalMessageInfo

func (m *UpdateProductStatusRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

func (m *UpdateProductStatusRequest) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

type UpdateProductStatusResponse struct {
	Updated              int32    `protobuf:"varint,2,opt,name=updated" json:"updated"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateProductStatusResponse) Reset()         { *m = UpdateProductStatusResponse{} }
func (m *UpdateProductStatusResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateProductStatusResponse) ProtoMessage()    {}
func (*UpdateProductStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{122}
}

var xxx_messageInfo_UpdateProductStatusResponse proto.InternalMessageInfo

func (m *UpdateProductStatusResponse) GetUpdated() int32 {
	if m != nil {
		return m.Updated
	}
	return 0
}

type PaymentTradingOrderRequest struct {
	OrderId              dot.ID                           `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	Desc                 string                           `protobuf:"bytes,2,opt,name=desc" json:"desc"`
	ReturnUrl            string                           `protobuf:"bytes,3,opt,name=return_url,json=returnUrl" json:"return_url"`
	Amount               int32                            `protobuf:"varint,4,opt,name=amount" json:"amount"`
	PaymentProvider      payment_provider.PaymentProvider `protobuf:"varint,5,opt,name=payment_provider,json=paymentProvider,enum=payment_provider.PaymentProvider" json:"payment_provider"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *PaymentTradingOrderRequest) Reset()         { *m = PaymentTradingOrderRequest{} }
func (m *PaymentTradingOrderRequest) String() string { return proto.CompactTextString(m) }
func (*PaymentTradingOrderRequest) ProtoMessage()    {}
func (*PaymentTradingOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{123}
}

var xxx_messageInfo_PaymentTradingOrderRequest proto.InternalMessageInfo

func (m *PaymentTradingOrderRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *PaymentTradingOrderRequest) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *PaymentTradingOrderRequest) GetReturnUrl() string {
	if m != nil {
		return m.ReturnUrl
	}
	return ""
}

func (m *PaymentTradingOrderRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *PaymentTradingOrderRequest) GetPaymentProvider() payment_provider.PaymentProvider {
	if m != nil {
		return m.PaymentProvider
	}
	return payment_provider.PaymentProvider_unknown
}

type PaymentTradingOrderResponse struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url" json:"url"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PaymentTradingOrderResponse) Reset()         { *m = PaymentTradingOrderResponse{} }
func (m *PaymentTradingOrderResponse) String() string { return proto.CompactTextString(m) }
func (*PaymentTradingOrderResponse) ProtoMessage()    {}
func (*PaymentTradingOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{124}
}

var xxx_messageInfo_PaymentTradingOrderResponse proto.InternalMessageInfo

func (m *PaymentTradingOrderResponse) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type UpdateVariantAttributesRequest struct {
	// @required
	VariantId            dot.ID       `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	Attributes           []*Attribute `protobuf:"bytes,2,rep,name=attributes" json:"attributes"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *UpdateVariantAttributesRequest) Reset()         { *m = UpdateVariantAttributesRequest{} }
func (m *UpdateVariantAttributesRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVariantAttributesRequest) ProtoMessage()    {}
func (*UpdateVariantAttributesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{125}
}

var xxx_messageInfo_UpdateVariantAttributesRequest proto.InternalMessageInfo

func (m *UpdateVariantAttributesRequest) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *UpdateVariantAttributesRequest) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type PaymentCheckReturnDataRequest struct {
	Id                    string                           `protobuf:"bytes,1,opt,name=id" json:"id"`
	Code                  string                           `protobuf:"bytes,2,opt,name=code" json:"code"`
	PaymentStatus         string                           `protobuf:"bytes,3,opt,name=payment_status,json=paymentStatus" json:"payment_status"`
	Amount                int32                            `protobuf:"varint,4,opt,name=amount" json:"amount"`
	ExternalTransactionId string                           `protobuf:"bytes,5,opt,name=external_transaction_id,json=externalTransactionId" json:"external_transaction_id"`
	PaymentProvider       payment_provider.PaymentProvider `protobuf:"varint,6,opt,name=payment_provider,json=paymentProvider,enum=payment_provider.PaymentProvider" json:"payment_provider"`
	XXX_NoUnkeyedLiteral  struct{}                         `json:"-"`
	XXX_sizecache         int32                            `json:"-"`
}

func (m *PaymentCheckReturnDataRequest) Reset()         { *m = PaymentCheckReturnDataRequest{} }
func (m *PaymentCheckReturnDataRequest) String() string { return proto.CompactTextString(m) }
func (*PaymentCheckReturnDataRequest) ProtoMessage()    {}
func (*PaymentCheckReturnDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{126}
}

var xxx_messageInfo_PaymentCheckReturnDataRequest proto.InternalMessageInfo

func (m *PaymentCheckReturnDataRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PaymentCheckReturnDataRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *PaymentCheckReturnDataRequest) GetPaymentStatus() string {
	if m != nil {
		return m.PaymentStatus
	}
	return ""
}

func (m *PaymentCheckReturnDataRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *PaymentCheckReturnDataRequest) GetExternalTransactionId() string {
	if m != nil {
		return m.ExternalTransactionId
	}
	return ""
}

func (m *PaymentCheckReturnDataRequest) GetPaymentProvider() payment_provider.PaymentProvider {
	if m != nil {
		return m.PaymentProvider
	}
	return payment_provider.PaymentProvider_unknown
}

type ShopCategory struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	ParentId             dot.ID   `protobuf:"varint,3,opt,name=parent_id,json=parentId" json:"parent_id"`
	ShopId               dot.ID   `protobuf:"varint,4,opt,name=shop_id,json=shopId" json:"shop_id"`
	Status               dot.ID   `protobuf:"varint,5,opt,name=status" json:"status"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopCategory) Reset()         { *m = ShopCategory{} }
func (m *ShopCategory) String() string { return proto.CompactTextString(m) }
func (*ShopCategory) ProtoMessage()    {}
func (*ShopCategory) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{127}
}

var xxx_messageInfo_ShopCategory proto.InternalMessageInfo

func (m *ShopCategory) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShopCategory) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShopCategory) GetParentId() dot.ID {
	if m != nil {
		return m.ParentId
	}
	return 0
}

func (m *ShopCategory) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *ShopCategory) GetStatus() dot.ID {
	if m != nil {
		return m.Status
	}
	return 0
}

type GetCollectionsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,4,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetCollectionsRequest) Reset()         { *m = GetCollectionsRequest{} }
func (m *GetCollectionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetCollectionsRequest) ProtoMessage()    {}
func (*GetCollectionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{128}
}

var xxx_messageInfo_GetCollectionsRequest proto.InternalMessageInfo

func (m *GetCollectionsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCollectionsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type ShopCollectionsResponse struct {
	Paging               *common.PageInfo  `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Collections          []*ShopCollection `protobuf:"bytes,2,rep,name=collections" json:"collections"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ShopCollectionsResponse) Reset()         { *m = ShopCollectionsResponse{} }
func (m *ShopCollectionsResponse) String() string { return proto.CompactTextString(m) }
func (*ShopCollectionsResponse) ProtoMessage()    {}
func (*ShopCollectionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{129}
}

var xxx_messageInfo_ShopCollectionsResponse proto.InternalMessageInfo

func (m *ShopCollectionsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShopCollectionsResponse) GetCollections() []*ShopCollection {
	if m != nil {
		return m.Collections
	}
	return nil
}

type AddShopProductCollectionRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	CollectionIds        []dot.ID `protobuf:"varint,2,rep,name=collection_ids,json=collectionIds" json:"collection_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddShopProductCollectionRequest) Reset()         { *m = AddShopProductCollectionRequest{} }
func (m *AddShopProductCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*AddShopProductCollectionRequest) ProtoMessage()    {}
func (*AddShopProductCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{130}
}

var xxx_messageInfo_AddShopProductCollectionRequest proto.InternalMessageInfo

func (m *AddShopProductCollectionRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *AddShopProductCollectionRequest) GetCollectionIds() []dot.ID {
	if m != nil {
		return m.CollectionIds
	}
	return nil
}

type RemoveShopProductCollectionRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	CollectionIds        []dot.ID `protobuf:"varint,2,rep,name=collection_ids,json=collectionIds" json:"collection_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveShopProductCollectionRequest) Reset()         { *m = RemoveShopProductCollectionRequest{} }
func (m *RemoveShopProductCollectionRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveShopProductCollectionRequest) ProtoMessage()    {}
func (*RemoveShopProductCollectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{131}
}

var xxx_messageInfo_RemoveShopProductCollectionRequest proto.InternalMessageInfo

func (m *RemoveShopProductCollectionRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *RemoveShopProductCollectionRequest) GetCollectionIds() []dot.ID {
	if m != nil {
		return m.CollectionIds
	}
	return nil
}

type AddCustomerToGroupRequest struct {
	CustomerIds          []dot.ID `protobuf:"varint,1,rep,name=customer_ids,json=customerIds" json:"customer_ids"`
	GroupId              dot.ID   `protobuf:"varint,2,opt,name=group_id,json=groupId" json:"group_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddCustomerToGroupRequest) Reset()         { *m = AddCustomerToGroupRequest{} }
func (m *AddCustomerToGroupRequest) String() string { return proto.CompactTextString(m) }
func (*AddCustomerToGroupRequest) ProtoMessage()    {}
func (*AddCustomerToGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{132}
}

var xxx_messageInfo_AddCustomerToGroupRequest proto.InternalMessageInfo

func (m *AddCustomerToGroupRequest) GetCustomerIds() []dot.ID {
	if m != nil {
		return m.CustomerIds
	}
	return nil
}

func (m *AddCustomerToGroupRequest) GetGroupId() dot.ID {
	if m != nil {
		return m.GroupId
	}
	return 0
}

type RemoveCustomerOutOfGroupRequest struct {
	CustomerIds          []dot.ID `protobuf:"varint,1,rep,name=customer_ids,json=customerIds" json:"customer_ids"`
	GroupId              dot.ID   `protobuf:"varint,2,opt,name=group_id,json=groupId" json:"group_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveCustomerOutOfGroupRequest) Reset()         { *m = RemoveCustomerOutOfGroupRequest{} }
func (m *RemoveCustomerOutOfGroupRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveCustomerOutOfGroupRequest) ProtoMessage()    {}
func (*RemoveCustomerOutOfGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{133}
}

var xxx_messageInfo_RemoveCustomerOutOfGroupRequest proto.InternalMessageInfo

func (m *RemoveCustomerOutOfGroupRequest) GetCustomerIds() []dot.ID {
	if m != nil {
		return m.CustomerIds
	}
	return nil
}

func (m *RemoveCustomerOutOfGroupRequest) GetGroupId() dot.ID {
	if m != nil {
		return m.GroupId
	}
	return 0
}

type SupplierLiability struct {
	TotalPurchaseOrders  int      `protobuf:"varint,1,opt,name=total_purchase_orders,json=totalPurchaseOrders" json:"total_purchase_orders"`
	TotalAmount          int      `protobuf:"varint,2,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	PaidAmount           int      `protobuf:"varint,3,opt,name=paid_amount,json=paidAmount" json:"paid_amount"`
	Liability            int      `protobuf:"varint,4,opt,name=liability" json:"liability"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SupplierLiability) Reset()         { *m = SupplierLiability{} }
func (m *SupplierLiability) String() string { return proto.CompactTextString(m) }
func (*SupplierLiability) ProtoMessage()    {}
func (*SupplierLiability) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{134}
}

var xxx_messageInfo_SupplierLiability proto.InternalMessageInfo

type Supplier struct {
	Id                   dot.ID             `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID             `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	FullName             string             `protobuf:"bytes,3,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 string             `protobuf:"bytes,4,opt,name=note" json:"note"`
	Phone                string             `protobuf:"bytes,5,opt,name=phone" json:"phone"`
	Email                string             `protobuf:"bytes,6,opt,name=email" json:"email"`
	CompanyName          string             `protobuf:"bytes,7,opt,name=company_name,json=companyName" json:"company_name"`
	TaxNumber            string             `protobuf:"bytes,8,opt,name=tax_number,json=taxNumber" json:"tax_number"`
	HeadquaterAddress    string             `protobuf:"bytes,9,opt,name=headquater_address,json=headquaterAddress" json:"headquater_address"`
	Code                 string             `protobuf:"bytes,10,opt,name=code" json:"code"`
	Status               status3.Status     `protobuf:"varint,11,opt,name=status,enum=status3.Status" json:"status"`
	CreatedAt            dot.Time           `protobuf:"bytes,12,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time           `protobuf:"bytes,13,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	Liability            *SupplierLiability `protobuf:"bytes,14,opt,name=liability" json:"liability"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Supplier) Reset()         { *m = Supplier{} }
func (m *Supplier) String() string { return proto.CompactTextString(m) }
func (*Supplier) ProtoMessage()    {}
func (*Supplier) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{135}
}

var xxx_messageInfo_Supplier proto.InternalMessageInfo

func (m *Supplier) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Supplier) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Supplier) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Supplier) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Supplier) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Supplier) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Supplier) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *Supplier) GetTaxNumber() string {
	if m != nil {
		return m.TaxNumber
	}
	return ""
}

func (m *Supplier) GetHeadquaterAddress() string {
	if m != nil {
		return m.HeadquaterAddress
	}
	return ""
}

func (m *Supplier) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Supplier) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Supplier) GetLiability() *SupplierLiability {
	if m != nil {
		return m.Liability
	}
	return nil
}

type CreateSupplierRequest struct {
	FullName             string   `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 string   `protobuf:"bytes,2,opt,name=note" json:"note"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,4,opt,name=email" json:"email"`
	CompanyName          string   `protobuf:"bytes,5,opt,name=company_name,json=companyName" json:"company_name"`
	TaxNumber            string   `protobuf:"bytes,6,opt,name=tax_number,json=taxNumber" json:"tax_number"`
	HeadquaterAddress    string   `protobuf:"bytes,7,opt,name=headquater_address,json=headquaterAddress" json:"headquater_address"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateSupplierRequest) Reset()         { *m = CreateSupplierRequest{} }
func (m *CreateSupplierRequest) String() string { return proto.CompactTextString(m) }
func (*CreateSupplierRequest) ProtoMessage()    {}
func (*CreateSupplierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{136}
}

var xxx_messageInfo_CreateSupplierRequest proto.InternalMessageInfo

func (m *CreateSupplierRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateSupplierRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateSupplierRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *CreateSupplierRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateSupplierRequest) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *CreateSupplierRequest) GetTaxNumber() string {
	if m != nil {
		return m.TaxNumber
	}
	return ""
}

func (m *CreateSupplierRequest) GetHeadquaterAddress() string {
	if m != nil {
		return m.HeadquaterAddress
	}
	return ""
}

type UpdateSupplierRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	FullName             *string  `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 *string  `protobuf:"bytes,3,opt,name=note" json:"note"`
	Phone                *string  `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	Email                *string  `protobuf:"bytes,5,opt,name=email" json:"email"`
	CompanyName          *string  `protobuf:"bytes,6,opt,name=company_name,json=companyName" json:"company_name"`
	TaxNumber            *string  `protobuf:"bytes,7,opt,name=tax_number,json=taxNumber" json:"tax_number"`
	HeadquaterAddress    *string  `protobuf:"bytes,8,opt,name=headquater_address,json=headquaterAddress" json:"headquater_address"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateSupplierRequest) Reset()         { *m = UpdateSupplierRequest{} }
func (m *UpdateSupplierRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateSupplierRequest) ProtoMessage()    {}
func (*UpdateSupplierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{137}
}

var xxx_messageInfo_UpdateSupplierRequest proto.InternalMessageInfo

func (m *UpdateSupplierRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateSupplierRequest) GetFullName() string {
	if m != nil && m.FullName != nil {
		return *m.FullName
	}
	return ""
}

func (m *UpdateSupplierRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateSupplierRequest) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

func (m *UpdateSupplierRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *UpdateSupplierRequest) GetCompanyName() string {
	if m != nil && m.CompanyName != nil {
		return *m.CompanyName
	}
	return ""
}

func (m *UpdateSupplierRequest) GetTaxNumber() string {
	if m != nil && m.TaxNumber != nil {
		return *m.TaxNumber
	}
	return ""
}

func (m *UpdateSupplierRequest) GetHeadquaterAddress() string {
	if m != nil && m.HeadquaterAddress != nil {
		return *m.HeadquaterAddress
	}
	return ""
}

type GetSuppliersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetSuppliersRequest) Reset()         { *m = GetSuppliersRequest{} }
func (m *GetSuppliersRequest) String() string { return proto.CompactTextString(m) }
func (*GetSuppliersRequest) ProtoMessage()    {}
func (*GetSuppliersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{138}
}

var xxx_messageInfo_GetSuppliersRequest proto.InternalMessageInfo

func (m *GetSuppliersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetSuppliersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type SuppliersResponse struct {
	Suppliers            []*Supplier      `protobuf:"bytes,1,rep,name=suppliers" json:"suppliers"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SuppliersResponse) Reset()         { *m = SuppliersResponse{} }
func (m *SuppliersResponse) String() string { return proto.CompactTextString(m) }
func (*SuppliersResponse) ProtoMessage()    {}
func (*SuppliersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{139}
}

var xxx_messageInfo_SuppliersResponse proto.InternalMessageInfo

func (m *SuppliersResponse) GetSuppliers() []*Supplier {
	if m != nil {
		return m.Suppliers
	}
	return nil
}

func (m *SuppliersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type Carrier struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID         `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	FullName             string         `protobuf:"bytes,3,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 string         `protobuf:"bytes,4,opt,name=note" json:"note"`
	Status               status3.Status `protobuf:"varint,5,opt,name=status,enum=status3.Status" json:"status"`
	CreatedAt            dot.Time       `protobuf:"bytes,6,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time       `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Carrier) Reset()         { *m = Carrier{} }
func (m *Carrier) String() string { return proto.CompactTextString(m) }
func (*Carrier) ProtoMessage()    {}
func (*Carrier) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{140}
}

var xxx_messageInfo_Carrier proto.InternalMessageInfo

func (m *Carrier) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Carrier) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Carrier) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Carrier) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Carrier) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

type CreateCarrierRequest struct {
	FullName             string   `protobuf:"bytes,1,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 string   `protobuf:"bytes,2,opt,name=note" json:"note"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCarrierRequest) Reset()         { *m = CreateCarrierRequest{} }
func (m *CreateCarrierRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCarrierRequest) ProtoMessage()    {}
func (*CreateCarrierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{141}
}

var xxx_messageInfo_CreateCarrierRequest proto.InternalMessageInfo

func (m *CreateCarrierRequest) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *CreateCarrierRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type UpdateCarrierRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	FullName             *string  `protobuf:"bytes,2,opt,name=full_name,json=fullName" json:"full_name"`
	Note                 *string  `protobuf:"bytes,3,opt,name=note" json:"note"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateCarrierRequest) Reset()         { *m = UpdateCarrierRequest{} }
func (m *UpdateCarrierRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCarrierRequest) ProtoMessage()    {}
func (*UpdateCarrierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{142}
}

var xxx_messageInfo_UpdateCarrierRequest proto.InternalMessageInfo

func (m *UpdateCarrierRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateCarrierRequest) GetFullName() string {
	if m != nil && m.FullName != nil {
		return *m.FullName
	}
	return ""
}

func (m *UpdateCarrierRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

type GetCarriersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetCarriersRequest) Reset()         { *m = GetCarriersRequest{} }
func (m *GetCarriersRequest) String() string { return proto.CompactTextString(m) }
func (*GetCarriersRequest) ProtoMessage()    {}
func (*GetCarriersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{143}
}

var xxx_messageInfo_GetCarriersRequest proto.InternalMessageInfo

func (m *GetCarriersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCarriersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type CarriersResponse struct {
	Carriers             []*Carrier       `protobuf:"bytes,1,rep,name=carriers" json:"carriers"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CarriersResponse) Reset()         { *m = CarriersResponse{} }
func (m *CarriersResponse) String() string { return proto.CompactTextString(m) }
func (*CarriersResponse) ProtoMessage()    {}
func (*CarriersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{144}
}

var xxx_messageInfo_CarriersResponse proto.InternalMessageInfo

func (m *CarriersResponse) GetCarriers() []*Carrier {
	if m != nil {
		return m.Carriers
	}
	return nil
}

func (m *CarriersResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type ReceiptLine struct {
	RefId                dot.ID   `protobuf:"varint,1,opt,name=ref_id,json=refId" json:"ref_id"`
	Title                string   `protobuf:"bytes,2,opt,name=title" json:"title"`
	Amount               int      `protobuf:"varint,3,req,name=amount" json:"amount"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReceiptLine) Reset()         { *m = ReceiptLine{} }
func (m *ReceiptLine) String() string { return proto.CompactTextString(m) }
func (*ReceiptLine) ProtoMessage()    {}
func (*ReceiptLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{145}
}

var xxx_messageInfo_ReceiptLine proto.InternalMessageInfo

func (m *ReceiptLine) GetRefId() dot.ID {
	if m != nil {
		return m.RefId
	}
	return 0
}

func (m *ReceiptLine) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

type Trader struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type"`
	FullName             string   `protobuf:"bytes,3,opt,name=full_name,json=fullName" json:"full_name"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone" json:"phone"`
	Deleted              bool     `protobuf:"varint,5,opt,name=deleted" json:"deleted"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Trader) Reset()         { *m = Trader{} }
func (m *Trader) String() string { return proto.CompactTextString(m) }
func (*Trader) ProtoMessage()    {}
func (*Trader) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{146}
}

var xxx_messageInfo_Trader proto.InternalMessageInfo

func (m *Trader) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Trader) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Trader) GetFullName() string {
	if m != nil {
		return m.FullName
	}
	return ""
}

func (m *Trader) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Trader) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

type Receipt struct {
	Id                   dot.ID         `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID         `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	TraderId             dot.ID         `protobuf:"varint,3,opt,name=trader_id,json=traderId" json:"trader_id"`
	CreatedBy            dot.ID         `protobuf:"varint,4,opt,name=created_by,json=createdBy" json:"created_by"`
	CreatedType          string         `protobuf:"bytes,5,opt,name=created_type,json=createdType" json:"created_type"`
	Code                 string         `protobuf:"bytes,6,opt,name=code" json:"code"`
	Title                string         `protobuf:"bytes,7,opt,name=title" json:"title"`
	Type                 string         `protobuf:"bytes,8,opt,name=type" json:"type"`
	Description          string         `protobuf:"bytes,9,opt,name=description" json:"description"`
	Amount               int            `protobuf:"varint,10,opt,name=amount" json:"amount"`
	LedgerId             dot.ID         `protobuf:"varint,11,opt,name=ledger_id,json=ledgerId" json:"ledger_id"`
	RefType              string         `protobuf:"bytes,23,opt,name=ref_type,json=refType" json:"ref_type"`
	Lines                []*ReceiptLine `protobuf:"bytes,12,rep,name=lines" json:"lines"`
	Status               status3.Status `protobuf:"varint,13,opt,name=status,enum=status3.Status" json:"status"`
	Reason               string         `protobuf:"bytes,14,opt,name=reason" json:"reason"`
	PaidAt               dot.Time       `protobuf:"bytes,15,opt,name=paid_at,json=paidAt" json:"paid_at"`
	CreatedAt            dot.Time       `protobuf:"bytes,16,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time       `protobuf:"bytes,17,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ConfirmedAt          dot.Time       `protobuf:"bytes,18,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelledAt          dot.Time       `protobuf:"bytes,19,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	User                 *etop.User     `protobuf:"bytes,20,opt,name=user" json:"user"`
	Trader               *Trader        `protobuf:"bytes,21,opt,name=trader" json:"trader"`
	Ledger               *Ledger        `protobuf:"bytes,22,opt,name=ledger" json:"ledger"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{147}
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Receipt) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Receipt) GetTraderId() dot.ID {
	if m != nil {
		return m.TraderId
	}
	return 0
}

func (m *Receipt) GetCreatedBy() dot.ID {
	if m != nil {
		return m.CreatedBy
	}
	return 0
}

func (m *Receipt) GetCreatedType() string {
	if m != nil {
		return m.CreatedType
	}
	return ""
}

func (m *Receipt) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Receipt) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Receipt) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Receipt) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Receipt) GetLedgerId() dot.ID {
	if m != nil {
		return m.LedgerId
	}
	return 0
}

func (m *Receipt) GetRefType() string {
	if m != nil {
		return m.RefType
	}
	return ""
}

func (m *Receipt) GetLines() []*ReceiptLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *Receipt) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Receipt) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Receipt) GetUser() *etop.User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Receipt) GetTrader() *Trader {
	if m != nil {
		return m.Trader
	}
	return nil
}

func (m *Receipt) GetLedger() *Ledger {
	if m != nil {
		return m.Ledger
	}
	return nil
}

type CreateReceiptRequest struct {
	TraderId dot.ID `protobuf:"varint,1,opt,name=trader_id,json=traderId" json:"trader_id"`
	Title    string `protobuf:"bytes,2,opt,name=title" json:"title"`
	// enum('receipt', 'payment')
	Type        string `protobuf:"bytes,3,opt,name=type" json:"type"`
	Description string `protobuf:"bytes,4,opt,name=description" json:"description"`
	Amount      int    `protobuf:"varint,5,opt,name=amount" json:"amount"`
	LedgerId    dot.ID `protobuf:"varint,6,opt,name=ledger_id,json=ledgerId" json:"ledger_id"`
	// enum ('order', 'fulfillment', 'inventory voucher'
	RefType              string         `protobuf:"bytes,9,opt,name=ref_type,json=refType" json:"ref_type"`
	PaidAt               dot.Time       `protobuf:"bytes,7,req,name=paid_at,json=paidAt" json:"paid_at"`
	Lines                []*ReceiptLine `protobuf:"bytes,8,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateReceiptRequest) Reset()         { *m = CreateReceiptRequest{} }
func (m *CreateReceiptRequest) String() string { return proto.CompactTextString(m) }
func (*CreateReceiptRequest) ProtoMessage()    {}
func (*CreateReceiptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{148}
}

var xxx_messageInfo_CreateReceiptRequest proto.InternalMessageInfo

func (m *CreateReceiptRequest) GetTraderId() dot.ID {
	if m != nil {
		return m.TraderId
	}
	return 0
}

func (m *CreateReceiptRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreateReceiptRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *CreateReceiptRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateReceiptRequest) GetLedgerId() dot.ID {
	if m != nil {
		return m.LedgerId
	}
	return 0
}

func (m *CreateReceiptRequest) GetRefType() string {
	if m != nil {
		return m.RefType
	}
	return ""
}

func (m *CreateReceiptRequest) GetLines() []*ReceiptLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type UpdateReceiptRequest struct {
	Id          dot.ID  `protobuf:"varint,1,opt,name=id" json:"id"`
	TraderId    *dot.ID `protobuf:"varint,2,opt,name=trader_id,json=traderId" json:"trader_id"`
	Title       *string `protobuf:"bytes,3,opt,name=title" json:"title"`
	Description *string `protobuf:"bytes,4,opt,name=description" json:"description"`
	Amount      *int    `protobuf:"varint,5,opt,name=amount" json:"amount"`
	LedgerId    *dot.ID `protobuf:"varint,6,opt,name=ledger_id,json=ledgerId" json:"ledger_id"`
	// enum ('order', 'fulfillment', 'inventory voucher'
	RefType              *string        `protobuf:"bytes,9,opt,name=ref_type,json=refType" json:"ref_type"`
	PaidAt               dot.Time       `protobuf:"bytes,7,opt,name=paid_at,json=paidAt" json:"paid_at"`
	Lines                []*ReceiptLine `protobuf:"bytes,8,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UpdateReceiptRequest) Reset()         { *m = UpdateReceiptRequest{} }
func (m *UpdateReceiptRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateReceiptRequest) ProtoMessage()    {}
func (*UpdateReceiptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{149}
}

var xxx_messageInfo_UpdateReceiptRequest proto.InternalMessageInfo

func (m *UpdateReceiptRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateReceiptRequest) GetTraderId() dot.ID {
	if m != nil && m.TraderId != nil {
		return *m.TraderId
	}
	return 0
}

func (m *UpdateReceiptRequest) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *UpdateReceiptRequest) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *UpdateReceiptRequest) GetLedgerId() dot.ID {
	if m != nil && m.LedgerId != nil {
		return *m.LedgerId
	}
	return 0
}

func (m *UpdateReceiptRequest) GetRefType() string {
	if m != nil && m.RefType != nil {
		return *m.RefType
	}
	return ""
}

func (m *UpdateReceiptRequest) GetPaidAt() dot.Time {
	if m != nil {
		return m.PaidAt
	}
	return dot.Time{}
}

func (m *UpdateReceiptRequest) GetLines() []*ReceiptLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type CancelReceiptRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason" json:"reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelReceiptRequest) Reset()         { *m = CancelReceiptRequest{} }
func (m *CancelReceiptRequest) String() string { return proto.CompactTextString(m) }
func (*CancelReceiptRequest) ProtoMessage()    {}
func (*CancelReceiptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{150}
}

var xxx_messageInfo_CancelReceiptRequest proto.InternalMessageInfo

func (m *CancelReceiptRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelReceiptRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type GetReceiptsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetReceiptsRequest) Reset()         { *m = GetReceiptsRequest{} }
func (m *GetReceiptsRequest) String() string { return proto.CompactTextString(m) }
func (*GetReceiptsRequest) ProtoMessage()    {}
func (*GetReceiptsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{151}
}

var xxx_messageInfo_GetReceiptsRequest proto.InternalMessageInfo

func (m *GetReceiptsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetReceiptsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetReceiptsByLedgerTypeRequest struct {
	Type                 string           `protobuf:"bytes,1,opt,name=type" json:"type"`
	Paging               *common.Paging   `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,3,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetReceiptsByLedgerTypeRequest) Reset()         { *m = GetReceiptsByLedgerTypeRequest{} }
func (m *GetReceiptsByLedgerTypeRequest) String() string { return proto.CompactTextString(m) }
func (*GetReceiptsByLedgerTypeRequest) ProtoMessage()    {}
func (*GetReceiptsByLedgerTypeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{152}
}

var xxx_messageInfo_GetReceiptsByLedgerTypeRequest proto.InternalMessageInfo

func (m *GetReceiptsByLedgerTypeRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *GetReceiptsByLedgerTypeRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetReceiptsByLedgerTypeRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type ReceiptsResponse struct {
	TotalAmountConfirmedReceipt int              `protobuf:"varint,1,opt,name=total_amount_confirmed_receipt,json=totalAmountConfirmedReceipt" json:"total_amount_confirmed_receipt"`
	TotalAmountConfirmedPayment int              `protobuf:"varint,2,opt,name=total_amount_confirmed_payment,json=totalAmountConfirmedPayment" json:"total_amount_confirmed_payment"`
	Receipts                    []*Receipt       `protobuf:"bytes,3,rep,name=receipts" json:"receipts"`
	Paging                      *common.PageInfo `protobuf:"bytes,4,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral        struct{}         `json:"-"`
	XXX_sizecache               int32            `json:"-"`
}

func (m *ReceiptsResponse) Reset()         { *m = ReceiptsResponse{} }
func (m *ReceiptsResponse) String() string { return proto.CompactTextString(m) }
func (*ReceiptsResponse) ProtoMessage()    {}
func (*ReceiptsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{153}
}

var xxx_messageInfo_ReceiptsResponse proto.InternalMessageInfo

func (m *ReceiptsResponse) GetReceipts() []*Receipt {
	if m != nil {
		return m.Receipts
	}
	return nil
}

func (m *ReceiptsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type GetShopCollectionsByProductIDRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetShopCollectionsByProductIDRequest) Reset()         { *m = GetShopCollectionsByProductIDRequest{} }
func (m *GetShopCollectionsByProductIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetShopCollectionsByProductIDRequest) ProtoMessage()    {}
func (*GetShopCollectionsByProductIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{154}
}

var xxx_messageInfo_GetShopCollectionsByProductIDRequest proto.InternalMessageInfo

func (m *GetShopCollectionsByProductIDRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

type CreateInventoryVoucherRequest struct {
	RefId   dot.ID `protobuf:"varint,16,opt,name=ref_id,json=refId" json:"ref_id"`
	RefType string `protobuf:"bytes,17,opt,name=ref_type,json=refType" json:"ref_type"`
	//enum "in" or "out" only for ref_type = "order"
	Type                 string   `protobuf:"bytes,15,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateInventoryVoucherRequest) Reset()         { *m = CreateInventoryVoucherRequest{} }
func (m *CreateInventoryVoucherRequest) String() string { return proto.CompactTextString(m) }
func (*CreateInventoryVoucherRequest) ProtoMessage()    {}
func (*CreateInventoryVoucherRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{155}
}

var xxx_messageInfo_CreateInventoryVoucherRequest proto.InternalMessageInfo

func (m *CreateInventoryVoucherRequest) GetRefId() dot.ID {
	if m != nil {
		return m.RefId
	}
	return 0
}

func (m *CreateInventoryVoucherRequest) GetRefType() string {
	if m != nil {
		return m.RefType
	}
	return ""
}

func (m *CreateInventoryVoucherRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type InventoryVoucherLine struct {
	VariantId            dot.ID      `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	Code                 string      `protobuf:"bytes,3,opt,name=code" json:"code"`
	VariantName          string      `protobuf:"bytes,6,opt,name=variant_name,json=variantName" json:"variant_name"`
	ProductId            dot.ID      `protobuf:"varint,7,opt,name=product_id,json=productId" json:"product_id"`
	ProductName          string      `protobuf:"bytes,8,opt,name=product_name,json=productName" json:"product_name"`
	ImageUrl             string      `protobuf:"bytes,9,opt,name=image_url,json=imageUrl" json:"image_url"`
	Attributes           []Attribute `protobuf:"bytes,10,rep,name=attributes" json:"attributes"`
	Price                int32       `protobuf:"varint,2,opt,name=price" json:"price"`
	Quantity             int32       `protobuf:"varint,5,opt,name=quantity" json:"quantity"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *InventoryVoucherLine) Reset()         { *m = InventoryVoucherLine{} }
func (m *InventoryVoucherLine) String() string { return proto.CompactTextString(m) }
func (*InventoryVoucherLine) ProtoMessage()    {}
func (*InventoryVoucherLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{156}
}

var xxx_messageInfo_InventoryVoucherLine proto.InternalMessageInfo

func (m *InventoryVoucherLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *InventoryVoucherLine) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *InventoryVoucherLine) GetVariantName() string {
	if m != nil {
		return m.VariantName
	}
	return ""
}

func (m *InventoryVoucherLine) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *InventoryVoucherLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *InventoryVoucherLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *InventoryVoucherLine) GetAttributes() []Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *InventoryVoucherLine) GetPrice() int32 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *InventoryVoucherLine) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type CreateInventoryVoucherResponse struct {
	InventoryVouchers    []*InventoryVoucher `protobuf:"bytes,1,rep,name=inventory_vouchers,json=inventoryVouchers" json:"inventory_vouchers"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CreateInventoryVoucherResponse) Reset()         { *m = CreateInventoryVoucherResponse{} }
func (m *CreateInventoryVoucherResponse) String() string { return proto.CompactTextString(m) }
func (*CreateInventoryVoucherResponse) ProtoMessage()    {}
func (*CreateInventoryVoucherResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{157}
}

var xxx_messageInfo_CreateInventoryVoucherResponse proto.InternalMessageInfo

func (m *CreateInventoryVoucherResponse) GetInventoryVouchers() []*InventoryVoucher {
	if m != nil {
		return m.InventoryVouchers
	}
	return nil
}

type ConfirmInventoryVoucherRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmInventoryVoucherRequest) Reset()         { *m = ConfirmInventoryVoucherRequest{} }
func (m *ConfirmInventoryVoucherRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmInventoryVoucherRequest) ProtoMessage()    {}
func (*ConfirmInventoryVoucherRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{158}
}

var xxx_messageInfo_ConfirmInventoryVoucherRequest proto.InternalMessageInfo

func (m *ConfirmInventoryVoucherRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

type ConfirmInventoryVoucherResponse struct {
	InventoryVoucher     *InventoryVoucher `protobuf:"bytes,1,opt,name=inventory_voucher,json=inventoryVoucher" json:"inventory_voucher"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ConfirmInventoryVoucherResponse) Reset()         { *m = ConfirmInventoryVoucherResponse{} }
func (m *ConfirmInventoryVoucherResponse) String() string { return proto.CompactTextString(m) }
func (*ConfirmInventoryVoucherResponse) ProtoMessage()    {}
func (*ConfirmInventoryVoucherResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{159}
}

var xxx_messageInfo_ConfirmInventoryVoucherResponse proto.InternalMessageInfo

func (m *ConfirmInventoryVoucherResponse) GetInventoryVoucher() *InventoryVoucher {
	if m != nil {
		return m.InventoryVoucher
	}
	return nil
}

type CancelInventoryVoucherRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Reason               string   `protobuf:"bytes,2,opt,name=reason" json:"reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelInventoryVoucherRequest) Reset()         { *m = CancelInventoryVoucherRequest{} }
func (m *CancelInventoryVoucherRequest) String() string { return proto.CompactTextString(m) }
func (*CancelInventoryVoucherRequest) ProtoMessage()    {}
func (*CancelInventoryVoucherRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{160}
}

var xxx_messageInfo_CancelInventoryVoucherRequest proto.InternalMessageInfo

func (m *CancelInventoryVoucherRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelInventoryVoucherRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type CancelInventoryVoucherResponse struct {
	Inventory            *InventoryVoucher `protobuf:"bytes,1,opt,name=inventory" json:"inventory"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CancelInventoryVoucherResponse) Reset()         { *m = CancelInventoryVoucherResponse{} }
func (m *CancelInventoryVoucherResponse) String() string { return proto.CompactTextString(m) }
func (*CancelInventoryVoucherResponse) ProtoMessage()    {}
func (*CancelInventoryVoucherResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{161}
}

var xxx_messageInfo_CancelInventoryVoucherResponse proto.InternalMessageInfo

func (m *CancelInventoryVoucherResponse) GetInventory() *InventoryVoucher {
	if m != nil {
		return m.Inventory
	}
	return nil
}

type UpdateInventoryVoucherRequest struct {
	Id                   dot.ID                 `protobuf:"varint,1,opt,name=id" json:"id"`
	TraderId             *dot.ID                `protobuf:"varint,3,opt,name=trader_id,json=traderId" json:"trader_id"`
	Lines                []InventoryVoucherLine `protobuf:"bytes,6,rep,name=lines" json:"lines"`
	Note                 *string                `protobuf:"bytes,2,opt,name=note" json:"note"`
	Type                 string                 `protobuf:"bytes,4,opt,name=type" json:"type"`
	Title                *string                `protobuf:"bytes,7,opt,name=title" json:"title"`
	TotalAmount          int32                  `protobuf:"varint,8,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *UpdateInventoryVoucherRequest) Reset()         { *m = UpdateInventoryVoucherRequest{} }
func (m *UpdateInventoryVoucherRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateInventoryVoucherRequest) ProtoMessage()    {}
func (*UpdateInventoryVoucherRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{162}
}

var xxx_messageInfo_UpdateInventoryVoucherRequest proto.InternalMessageInfo

func (m *UpdateInventoryVoucherRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateInventoryVoucherRequest) GetTraderId() dot.ID {
	if m != nil && m.TraderId != nil {
		return *m.TraderId
	}
	return 0
}

func (m *UpdateInventoryVoucherRequest) GetLines() []InventoryVoucherLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *UpdateInventoryVoucherRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateInventoryVoucherRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *UpdateInventoryVoucherRequest) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *UpdateInventoryVoucherRequest) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

type UpdateInventoryVoucherResponse struct {
	InventoryVoucher     *InventoryVoucher `protobuf:"bytes,1,opt,name=inventory_voucher,json=inventoryVoucher" json:"inventory_voucher"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateInventoryVoucherResponse) Reset()         { *m = UpdateInventoryVoucherResponse{} }
func (m *UpdateInventoryVoucherResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateInventoryVoucherResponse) ProtoMessage()    {}
func (*UpdateInventoryVoucherResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{163}
}

var xxx_messageInfo_UpdateInventoryVoucherResponse proto.InternalMessageInfo

func (m *UpdateInventoryVoucherResponse) GetInventoryVoucher() *InventoryVoucher {
	if m != nil {
		return m.InventoryVoucher
	}
	return nil
}

type AdjustInventoryQuantityRequest struct {
	InventoryVariants    []*InventoryVariant `protobuf:"bytes,1,rep,name=inventory_variants,json=inventoryVariants" json:"inventory_variants"`
	Note                 string              `protobuf:"bytes,2,opt,name=note" json:"note"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *AdjustInventoryQuantityRequest) Reset()         { *m = AdjustInventoryQuantityRequest{} }
func (m *AdjustInventoryQuantityRequest) String() string { return proto.CompactTextString(m) }
func (*AdjustInventoryQuantityRequest) ProtoMessage()    {}
func (*AdjustInventoryQuantityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{164}
}

var xxx_messageInfo_AdjustInventoryQuantityRequest proto.InternalMessageInfo

func (m *AdjustInventoryQuantityRequest) GetInventoryVariants() []*InventoryVariant {
	if m != nil {
		return m.InventoryVariants
	}
	return nil
}

func (m *AdjustInventoryQuantityRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type AdjustInventoryQuantityResponse struct {
	InventoryVariants    []*InventoryVariant `protobuf:"bytes,1,rep,name=inventory_variants,json=inventoryVariants" json:"inventory_variants"`
	InventoryVouchers    []*InventoryVoucher `protobuf:"bytes,2,rep,name=inventory_vouchers,json=inventoryVouchers" json:"inventory_vouchers"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *AdjustInventoryQuantityResponse) Reset()         { *m = AdjustInventoryQuantityResponse{} }
func (m *AdjustInventoryQuantityResponse) String() string { return proto.CompactTextString(m) }
func (*AdjustInventoryQuantityResponse) ProtoMessage()    {}
func (*AdjustInventoryQuantityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{165}
}

var xxx_messageInfo_AdjustInventoryQuantityResponse proto.InternalMessageInfo

func (m *AdjustInventoryQuantityResponse) GetInventoryVariants() []*InventoryVariant {
	if m != nil {
		return m.InventoryVariants
	}
	return nil
}

func (m *AdjustInventoryQuantityResponse) GetInventoryVouchers() []*InventoryVoucher {
	if m != nil {
		return m.InventoryVouchers
	}
	return nil
}

type GetInventoryVariantsRequest struct {
	Paging               common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetInventoryVariantsRequest) Reset()         { *m = GetInventoryVariantsRequest{} }
func (m *GetInventoryVariantsRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVariantsRequest) ProtoMessage()    {}
func (*GetInventoryVariantsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{166}
}

var xxx_messageInfo_GetInventoryVariantsRequest proto.InternalMessageInfo

func (m *GetInventoryVariantsRequest) GetPaging() common.Paging {
	if m != nil {
		return m.Paging
	}
	return common.Paging{}
}

type GetInventoryVariantsResponse struct {
	InventoryVariants    []*InventoryVariant `protobuf:"bytes,1,rep,name=inventory_variants,json=inventoryVariants" json:"inventory_variants"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetInventoryVariantsResponse) Reset()         { *m = GetInventoryVariantsResponse{} }
func (m *GetInventoryVariantsResponse) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVariantsResponse) ProtoMessage()    {}
func (*GetInventoryVariantsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{167}
}

var xxx_messageInfo_GetInventoryVariantsResponse proto.InternalMessageInfo

func (m *GetInventoryVariantsResponse) GetInventoryVariants() []*InventoryVariant {
	if m != nil {
		return m.InventoryVariants
	}
	return nil
}

type GetInventoryVariantsByVariantIDsRequest struct {
	VariantIds           []dot.ID `protobuf:"varint,1,rep,name=variant_ids,json=variantIds" json:"variant_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInventoryVariantsByVariantIDsRequest) Reset() {
	*m = GetInventoryVariantsByVariantIDsRequest{}
}
func (m *GetInventoryVariantsByVariantIDsRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVariantsByVariantIDsRequest) ProtoMessage()    {}
func (*GetInventoryVariantsByVariantIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{168}
}

var xxx_messageInfo_GetInventoryVariantsByVariantIDsRequest proto.InternalMessageInfo

func (m *GetInventoryVariantsByVariantIDsRequest) GetVariantIds() []dot.ID {
	if m != nil {
		return m.VariantIds
	}
	return nil
}

type InventoryVariant struct {
	ShopId               dot.ID   `protobuf:"varint,1,opt,name=shop_id,json=shopId" json:"shop_id"`
	VariantId            dot.ID   `protobuf:"varint,2,opt,name=variant_id,json=variantId" json:"variant_id"`
	QuantityOnHand       int32    `protobuf:"varint,3,opt,name=quantity_on_hand,json=quantityOnHand" json:"quantity_on_hand"`
	QuantityPicked       int32    `protobuf:"varint,4,opt,name=quantity_picked,json=quantityPicked" json:"quantity_picked"`
	CostPrice            int32    `protobuf:"varint,5,opt,name=cost_price,json=costPrice" json:"cost_price"`
	Quantity             int32    `protobuf:"varint,6,opt,name=quantity" json:"quantity"`
	CreatedAt            dot.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InventoryVariant) Reset()         { *m = InventoryVariant{} }
func (m *InventoryVariant) String() string { return proto.CompactTextString(m) }
func (*InventoryVariant) ProtoMessage()    {}
func (*InventoryVariant) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{169}
}

var xxx_messageInfo_InventoryVariant proto.InternalMessageInfo

func (m *InventoryVariant) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *InventoryVariant) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *InventoryVariant) GetQuantityOnHand() int32 {
	if m != nil {
		return m.QuantityOnHand
	}
	return 0
}

func (m *InventoryVariant) GetQuantityPicked() int32 {
	if m != nil {
		return m.QuantityPicked
	}
	return 0
}

func (m *InventoryVariant) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *InventoryVariant) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type InventoryVariantQuantity struct {
	QuantityOnHand       int32    `protobuf:"varint,1,opt,name=quantity_on_hand,json=quantityOnHand" json:"quantity_on_hand"`
	QuantityPicked       int32    `protobuf:"varint,2,opt,name=quantity_picked,json=quantityPicked" json:"quantity_picked"`
	Quantity             int32    `protobuf:"varint,4,opt,name=quantity" json:"quantity"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InventoryVariantQuantity) Reset()         { *m = InventoryVariantQuantity{} }
func (m *InventoryVariantQuantity) String() string { return proto.CompactTextString(m) }
func (*InventoryVariantQuantity) ProtoMessage()    {}
func (*InventoryVariantQuantity) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{170}
}

var xxx_messageInfo_InventoryVariantQuantity proto.InternalMessageInfo

func (m *InventoryVariantQuantity) GetQuantityOnHand() int32 {
	if m != nil {
		return m.QuantityOnHand
	}
	return 0
}

func (m *InventoryVariantQuantity) GetQuantityPicked() int32 {
	if m != nil {
		return m.QuantityPicked
	}
	return 0
}

func (m *InventoryVariantQuantity) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type InventoryVoucher struct {
	TotalAmount          int32                   `protobuf:"varint,7,opt,name=total_amount,json=totalAmount" json:"total_amount"`
	CreatedBy            dot.ID                  `protobuf:"varint,8,opt,name=created_by,json=createdBy" json:"created_by"`
	UpdatedBy            dot.ID                  `protobuf:"varint,9,opt,name=updated_by,json=updatedBy" json:"updated_by"`
	Lines                []*InventoryVoucherLine `protobuf:"bytes,1,rep,name=lines" json:"lines"`
	TraderId             dot.ID                  `protobuf:"varint,2,opt,name=trader_id,json=traderId" json:"trader_id"`
	Note                 string                  `protobuf:"bytes,3,opt,name=note" json:"note"`
	Type                 string                  `protobuf:"bytes,4,opt,name=type" json:"type"`
	Id                   dot.ID                  `protobuf:"varint,5,opt,name=id" json:"id"`
	ShopId               dot.ID                  `protobuf:"varint,6,opt,name=shop_id,json=shopId" json:"shop_id"`
	Title                string                  `protobuf:"bytes,15,opt,name=title" json:"title"`
	RefId                dot.ID                  `protobuf:"varint,16,opt,name=ref_id,json=refId" json:"ref_id"`
	RefType              string                  `protobuf:"bytes,17,opt,name=ref_type,json=refType" json:"ref_type"`
	RefName              string                  `protobuf:"bytes,18,opt,name=ref_name,json=refName" json:"ref_name"`
	RefCode              string                  `protobuf:"bytes,21,opt,name=ref_code,json=refCode" json:"ref_code"`
	Code                 string                  `protobuf:"bytes,19,opt,name=code" json:"code"`
	CreatedAt            dot.Time                `protobuf:"bytes,10,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time                `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	CancelledAt          dot.Time                `protobuf:"bytes,12,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	ConfirmedAt          dot.Time                `protobuf:"bytes,13,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelReason         string                  `protobuf:"bytes,14,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	Trader               *Trader                 `protobuf:"bytes,20,opt,name=trader" json:"trader"`
	Status               status3.Status          `protobuf:"varint,22,opt,name=status,enum=status3.Status" json:"status"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *InventoryVoucher) Reset()         { *m = InventoryVoucher{} }
func (m *InventoryVoucher) String() string { return proto.CompactTextString(m) }
func (*InventoryVoucher) ProtoMessage()    {}
func (*InventoryVoucher) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{171}
}

var xxx_messageInfo_InventoryVoucher proto.InternalMessageInfo

func (m *InventoryVoucher) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *InventoryVoucher) GetCreatedBy() dot.ID {
	if m != nil {
		return m.CreatedBy
	}
	return 0
}

func (m *InventoryVoucher) GetUpdatedBy() dot.ID {
	if m != nil {
		return m.UpdatedBy
	}
	return 0
}

func (m *InventoryVoucher) GetLines() []*InventoryVoucherLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

func (m *InventoryVoucher) GetTraderId() dot.ID {
	if m != nil {
		return m.TraderId
	}
	return 0
}

func (m *InventoryVoucher) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *InventoryVoucher) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *InventoryVoucher) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *InventoryVoucher) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *InventoryVoucher) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *InventoryVoucher) GetRefId() dot.ID {
	if m != nil {
		return m.RefId
	}
	return 0
}

func (m *InventoryVoucher) GetRefType() string {
	if m != nil {
		return m.RefType
	}
	return ""
}

func (m *InventoryVoucher) GetRefName() string {
	if m != nil {
		return m.RefName
	}
	return ""
}

func (m *InventoryVoucher) GetRefCode() string {
	if m != nil {
		return m.RefCode
	}
	return ""
}

func (m *InventoryVoucher) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *InventoryVoucher) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *InventoryVoucher) GetTrader() *Trader {
	if m != nil {
		return m.Trader
	}
	return nil
}

func (m *InventoryVoucher) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

type GetInventoryVouchersRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetInventoryVouchersRequest) Reset()         { *m = GetInventoryVouchersRequest{} }
func (m *GetInventoryVouchersRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVouchersRequest) ProtoMessage()    {}
func (*GetInventoryVouchersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{172}
}

var xxx_messageInfo_GetInventoryVouchersRequest proto.InternalMessageInfo

func (m *GetInventoryVouchersRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetInventoryVouchersRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetInventoryVouchersByIDsRequest struct {
	Ids                  []dot.ID `protobuf:"varint,1,rep,name=ids" json:"ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInventoryVouchersByIDsRequest) Reset()         { *m = GetInventoryVouchersByIDsRequest{} }
func (m *GetInventoryVouchersByIDsRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVouchersByIDsRequest) ProtoMessage()    {}
func (*GetInventoryVouchersByIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{173}
}

var xxx_messageInfo_GetInventoryVouchersByIDsRequest proto.InternalMessageInfo

func (m *GetInventoryVouchersByIDsRequest) GetIds() []dot.ID {
	if m != nil {
		return m.Ids
	}
	return nil
}

type GetInventoryVouchersResponse struct {
	InventoryVouchers    []*InventoryVoucher `protobuf:"bytes,1,rep,name=inventory_vouchers,json=inventoryVouchers" json:"inventory_vouchers"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetInventoryVouchersResponse) Reset()         { *m = GetInventoryVouchersResponse{} }
func (m *GetInventoryVouchersResponse) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVouchersResponse) ProtoMessage()    {}
func (*GetInventoryVouchersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{174}
}

var xxx_messageInfo_GetInventoryVouchersResponse proto.InternalMessageInfo

func (m *GetInventoryVouchersResponse) GetInventoryVouchers() []*InventoryVoucher {
	if m != nil {
		return m.InventoryVouchers
	}
	return nil
}

type CustomerGroup struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CustomerGroup) Reset()         { *m = CustomerGroup{} }
func (m *CustomerGroup) String() string { return proto.CompactTextString(m) }
func (*CustomerGroup) ProtoMessage()    {}
func (*CustomerGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{175}
}

var xxx_messageInfo_CustomerGroup proto.InternalMessageInfo

func (m *CustomerGroup) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CustomerGroup) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateCustomerGroupRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCustomerGroupRequest) Reset()         { *m = CreateCustomerGroupRequest{} }
func (m *CreateCustomerGroupRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCustomerGroupRequest) ProtoMessage()    {}
func (*CreateCustomerGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{176}
}

var xxx_messageInfo_CreateCustomerGroupRequest proto.InternalMessageInfo

func (m *CreateCustomerGroupRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type UpdateCustomerGroupRequest struct {
	GroupId              dot.ID   `protobuf:"varint,1,opt,name=group_id,json=groupId" json:"group_id"`
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateCustomerGroupRequest) Reset()         { *m = UpdateCustomerGroupRequest{} }
func (m *UpdateCustomerGroupRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateCustomerGroupRequest) ProtoMessage()    {}
func (*UpdateCustomerGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{177}
}

var xxx_messageInfo_UpdateCustomerGroupRequest proto.InternalMessageInfo

func (m *UpdateCustomerGroupRequest) GetGroupId() dot.ID {
	if m != nil {
		return m.GroupId
	}
	return 0
}

func (m *UpdateCustomerGroupRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetCustomerGroupsRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetCustomerGroupsRequest) Reset()         { *m = GetCustomerGroupsRequest{} }
func (m *GetCustomerGroupsRequest) String() string { return proto.CompactTextString(m) }
func (*GetCustomerGroupsRequest) ProtoMessage()    {}
func (*GetCustomerGroupsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{178}
}

var xxx_messageInfo_GetCustomerGroupsRequest proto.InternalMessageInfo

func (m *GetCustomerGroupsRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetCustomerGroupsRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type CustomerGroupsResponse struct {
	CustomerGroups       []*CustomerGroup `protobuf:"bytes,1,rep,name=customer_groups,json=customerGroups" json:"customer_groups"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CustomerGroupsResponse) Reset()         { *m = CustomerGroupsResponse{} }
func (m *CustomerGroupsResponse) String() string { return proto.CompactTextString(m) }
func (*CustomerGroupsResponse) ProtoMessage()    {}
func (*CustomerGroupsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{179}
}

var xxx_messageInfo_CustomerGroupsResponse proto.InternalMessageInfo

func (m *CustomerGroupsResponse) GetCustomerGroups() []*CustomerGroup {
	if m != nil {
		return m.CustomerGroups
	}
	return nil
}

func (m *CustomerGroupsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type GetOrdersByReceiptIDRequest struct {
	ReceiptId            dot.ID   `protobuf:"varint,1,opt,name=receipt_id,json=receiptId" json:"receipt_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetOrdersByReceiptIDRequest) Reset()         { *m = GetOrdersByReceiptIDRequest{} }
func (m *GetOrdersByReceiptIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetOrdersByReceiptIDRequest) ProtoMessage()    {}
func (*GetOrdersByReceiptIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{180}
}

var xxx_messageInfo_GetOrdersByReceiptIDRequest proto.InternalMessageInfo

func (m *GetOrdersByReceiptIDRequest) GetReceiptId() dot.ID {
	if m != nil {
		return m.ReceiptId
	}
	return 0
}

type GetInventoryVariantRequest struct {
	VariantId            dot.ID   `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInventoryVariantRequest) Reset()         { *m = GetInventoryVariantRequest{} }
func (m *GetInventoryVariantRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVariantRequest) ProtoMessage()    {}
func (*GetInventoryVariantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{181}
}

var xxx_messageInfo_GetInventoryVariantRequest proto.InternalMessageInfo

func (m *GetInventoryVariantRequest) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

type CreateBrandRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Description          string   `protobuf:"bytes,2,opt,name=description" json:"description"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateBrandRequest) Reset()         { *m = CreateBrandRequest{} }
func (m *CreateBrandRequest) String() string { return proto.CompactTextString(m) }
func (*CreateBrandRequest) ProtoMessage()    {}
func (*CreateBrandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{182}
}

var xxx_messageInfo_CreateBrandRequest proto.InternalMessageInfo

func (m *CreateBrandRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateBrandRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type UpdateBrandRequest struct {
	Id                   dot.ID   `protobuf:"varint,3,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Description          string   `protobuf:"bytes,2,opt,name=description" json:"description"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateBrandRequest) Reset()         { *m = UpdateBrandRequest{} }
func (m *UpdateBrandRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateBrandRequest) ProtoMessage()    {}
func (*UpdateBrandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{183}
}

var xxx_messageInfo_UpdateBrandRequest proto.InternalMessageInfo

func (m *UpdateBrandRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateBrandRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateBrandRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type DeleteBrandResponse struct {
	Count                int32    `protobuf:"varint,1,opt,name=count" json:"count"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteBrandResponse) Reset()         { *m = DeleteBrandResponse{} }
func (m *DeleteBrandResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteBrandResponse) ProtoMessage()    {}
func (*DeleteBrandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{184}
}

var xxx_messageInfo_DeleteBrandResponse proto.InternalMessageInfo

func (m *DeleteBrandResponse) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Brand struct {
	ShopId               dot.ID   `protobuf:"varint,1,opt,name=shop_id,json=shopId" json:"shop_id"`
	Id                   dot.ID   `protobuf:"varint,2,opt,name=id" json:"id"`
	Name                 string   `protobuf:"bytes,3,opt,name=name" json:"name"`
	Description          string   `protobuf:"bytes,4,opt,name=description" json:"description"`
	CreatedAt            dot.Time `protobuf:"bytes,5,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Brand) Reset()         { *m = Brand{} }
func (m *Brand) String() string { return proto.CompactTextString(m) }
func (*Brand) ProtoMessage()    {}
func (*Brand) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{185}
}

var xxx_messageInfo_Brand proto.InternalMessageInfo

func (m *Brand) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Brand) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Brand) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Brand) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type GetBrandsByIDsResponse struct {
	Brands               []*Brand `protobuf:"bytes,1,rep,name=brands" json:"brands"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBrandsByIDsResponse) Reset()         { *m = GetBrandsByIDsResponse{} }
func (m *GetBrandsByIDsResponse) String() string { return proto.CompactTextString(m) }
func (*GetBrandsByIDsResponse) ProtoMessage()    {}
func (*GetBrandsByIDsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{186}
}

var xxx_messageInfo_GetBrandsByIDsResponse proto.InternalMessageInfo

func (m *GetBrandsByIDsResponse) GetBrands() []*Brand {
	if m != nil {
		return m.Brands
	}
	return nil
}

type GetBrandsRequest struct {
	Paging               common.Paging `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetBrandsRequest) Reset()         { *m = GetBrandsRequest{} }
func (m *GetBrandsRequest) String() string { return proto.CompactTextString(m) }
func (*GetBrandsRequest) ProtoMessage()    {}
func (*GetBrandsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{187}
}

var xxx_messageInfo_GetBrandsRequest proto.InternalMessageInfo

func (m *GetBrandsRequest) GetPaging() common.Paging {
	if m != nil {
		return m.Paging
	}
	return common.Paging{}
}

type GetBrandsResponse struct {
	Brands               []*Brand         `protobuf:"bytes,1,rep,name=brands" json:"brands"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetBrandsResponse) Reset()         { *m = GetBrandsResponse{} }
func (m *GetBrandsResponse) String() string { return proto.CompactTextString(m) }
func (*GetBrandsResponse) ProtoMessage()    {}
func (*GetBrandsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{188}
}

var xxx_messageInfo_GetBrandsResponse proto.InternalMessageInfo

func (m *GetBrandsResponse) GetBrands() []*Brand {
	if m != nil {
		return m.Brands
	}
	return nil
}

func (m *GetBrandsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type GetInventoryVouchersByReferenceRequest struct {
	RefId dot.ID `protobuf:"varint,1,opt,name=ref_id,json=refId" json:"ref_id"`
	// enum ('order', 'purchase_order', 'return', 'purchase_order')
	RefType              string   `protobuf:"bytes,2,opt,name=ref_type,json=refType" json:"ref_type"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInventoryVouchersByReferenceRequest) Reset() {
	*m = GetInventoryVouchersByReferenceRequest{}
}
func (m *GetInventoryVouchersByReferenceRequest) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVouchersByReferenceRequest) ProtoMessage()    {}
func (*GetInventoryVouchersByReferenceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{189}
}

var xxx_messageInfo_GetInventoryVouchersByReferenceRequest proto.InternalMessageInfo

func (m *GetInventoryVouchersByReferenceRequest) GetRefId() dot.ID {
	if m != nil {
		return m.RefId
	}
	return 0
}

func (m *GetInventoryVouchersByReferenceRequest) GetRefType() string {
	if m != nil {
		return m.RefType
	}
	return ""
}

type GetInventoryVouchersByReferenceResponse struct {
	InventoryVouchers    []*InventoryVoucher `protobuf:"bytes,1,rep,name=inventory_vouchers,json=inventoryVouchers" json:"inventory_vouchers"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetInventoryVouchersByReferenceResponse) Reset() {
	*m = GetInventoryVouchersByReferenceResponse{}
}
func (m *GetInventoryVouchersByReferenceResponse) String() string { return proto.CompactTextString(m) }
func (*GetInventoryVouchersByReferenceResponse) ProtoMessage()    {}
func (*GetInventoryVouchersByReferenceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{190}
}

var xxx_messageInfo_GetInventoryVouchersByReferenceResponse proto.InternalMessageInfo

func (m *GetInventoryVouchersByReferenceResponse) GetInventoryVouchers() []*InventoryVoucher {
	if m != nil {
		return m.InventoryVouchers
	}
	return nil
}

type UpdateOrderShippingInfoRequest struct {
	OrderId              dot.ID               `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	Shipping             *order.OrderShipping `protobuf:"bytes,2,opt,name=shipping" json:"shipping"`
	ShippingAddress      *order.OrderAddress  `protobuf:"bytes,3,opt,name=shipping_address,json=shippingAddress" json:"shipping_address"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *UpdateOrderShippingInfoRequest) Reset()         { *m = UpdateOrderShippingInfoRequest{} }
func (m *UpdateOrderShippingInfoRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateOrderShippingInfoRequest) ProtoMessage()    {}
func (*UpdateOrderShippingInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{191}
}

var xxx_messageInfo_UpdateOrderShippingInfoRequest proto.InternalMessageInfo

func (m *UpdateOrderShippingInfoRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *UpdateOrderShippingInfoRequest) GetShipping() *order.OrderShipping {
	if m != nil {
		return m.Shipping
	}
	return nil
}

func (m *UpdateOrderShippingInfoRequest) GetShippingAddress() *order.OrderAddress {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

type GetStocktakesByIDsResponse struct {
	Stocktakes           []*Stocktake `protobuf:"bytes,1,rep,name=stocktakes" json:"stocktakes"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetStocktakesByIDsResponse) Reset()         { *m = GetStocktakesByIDsResponse{} }
func (m *GetStocktakesByIDsResponse) String() string { return proto.CompactTextString(m) }
func (*GetStocktakesByIDsResponse) ProtoMessage()    {}
func (*GetStocktakesByIDsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{192}
}

var xxx_messageInfo_GetStocktakesByIDsResponse proto.InternalMessageInfo

func (m *GetStocktakesByIDsResponse) GetStocktakes() []*Stocktake {
	if m != nil {
		return m.Stocktakes
	}
	return nil
}

type CreateStocktakeRequest struct {
	TotalQuantity int32  `protobuf:"varint,1,opt,name=total_quantity,json=totalQuantity" json:"total_quantity"`
	Note          string `protobuf:"bytes,2,opt,name=note" json:"note"`
	//  length more than one
	Lines                []*StocktakeLine `protobuf:"bytes,3,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CreateStocktakeRequest) Reset()         { *m = CreateStocktakeRequest{} }
func (m *CreateStocktakeRequest) String() string { return proto.CompactTextString(m) }
func (*CreateStocktakeRequest) ProtoMessage()    {}
func (*CreateStocktakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{193}
}

var xxx_messageInfo_CreateStocktakeRequest proto.InternalMessageInfo

func (m *CreateStocktakeRequest) GetTotalQuantity() int32 {
	if m != nil {
		return m.TotalQuantity
	}
	return 0
}

func (m *CreateStocktakeRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateStocktakeRequest) GetLines() []*StocktakeLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type UpdateStocktakeRequest struct {
	Id            dot.ID `protobuf:"varint,4,opt,name=id" json:"id"`
	TotalQuantity int32  `protobuf:"varint,1,opt,name=total_quantity,json=totalQuantity" json:"total_quantity"`
	Note          string `protobuf:"bytes,2,opt,name=note" json:"note"`
	//  length more than one
	Lines                []*StocktakeLine `protobuf:"bytes,3,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *UpdateStocktakeRequest) Reset()         { *m = UpdateStocktakeRequest{} }
func (m *UpdateStocktakeRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateStocktakeRequest) ProtoMessage()    {}
func (*UpdateStocktakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{194}
}

var xxx_messageInfo_UpdateStocktakeRequest proto.InternalMessageInfo

func (m *UpdateStocktakeRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateStocktakeRequest) GetTotalQuantity() int32 {
	if m != nil {
		return m.TotalQuantity
	}
	return 0
}

func (m *UpdateStocktakeRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *UpdateStocktakeRequest) GetLines() []*StocktakeLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type Stocktake struct {
	Id                   dot.ID           `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId               dot.ID           `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	TotalQuantity        int32            `protobuf:"varint,3,opt,name=total_quantity,json=totalQuantity" json:"total_quantity"`
	Note                 string           `protobuf:"bytes,4,opt,name=note" json:"note"`
	Code                 string           `protobuf:"bytes,13,opt,name=code" json:"code"`
	CancelReason         string           `protobuf:"bytes,14,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	CreatedBy            dot.ID           `protobuf:"varint,5,opt,name=created_by,json=createdBy" json:"created_by"`
	UpdatedBy            dot.ID           `protobuf:"varint,6,opt,name=updated_by,json=updatedBy" json:"updated_by"`
	CreatedAt            dot.Time         `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time         `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	ConfirmedAt          dot.Time         `protobuf:"bytes,9,opt,name=confirmed_at,json=confirmedAt" json:"confirmed_at"`
	CancelledAt          dot.Time         `protobuf:"bytes,10,opt,name=cancelled_at,json=cancelledAt" json:"cancelled_at"`
	Status               status3.Status   `protobuf:"varint,12,opt,name=status,enum=status3.Status" json:"status"`
	Lines                []*StocktakeLine `protobuf:"bytes,11,rep,name=lines" json:"lines"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Stocktake) Reset()         { *m = Stocktake{} }
func (m *Stocktake) String() string { return proto.CompactTextString(m) }
func (*Stocktake) ProtoMessage()    {}
func (*Stocktake) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{195}
}

var xxx_messageInfo_Stocktake proto.InternalMessageInfo

func (m *Stocktake) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Stocktake) GetShopId() dot.ID {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *Stocktake) GetTotalQuantity() int32 {
	if m != nil {
		return m.TotalQuantity
	}
	return 0
}

func (m *Stocktake) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *Stocktake) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Stocktake) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *Stocktake) GetCreatedBy() dot.ID {
	if m != nil {
		return m.CreatedBy
	}
	return 0
}

func (m *Stocktake) GetUpdatedBy() dot.ID {
	if m != nil {
		return m.UpdatedBy
	}
	return 0
}

func (m *Stocktake) GetStatus() status3.Status {
	if m != nil {
		return m.Status
	}
	return status3.Status_Z
}

func (m *Stocktake) GetLines() []*StocktakeLine {
	if m != nil {
		return m.Lines
	}
	return nil
}

type GetStocktakesRequest struct {
	Paging               *common.Paging   `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Filters              []*common.Filter `protobuf:"bytes,2,rep,name=filters" json:"filters"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetStocktakesRequest) Reset()         { *m = GetStocktakesRequest{} }
func (m *GetStocktakesRequest) String() string { return proto.CompactTextString(m) }
func (*GetStocktakesRequest) ProtoMessage()    {}
func (*GetStocktakesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{196}
}

var xxx_messageInfo_GetStocktakesRequest proto.InternalMessageInfo

func (m *GetStocktakesRequest) GetPaging() *common.Paging {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetStocktakesRequest) GetFilters() []*common.Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type GetStocktakesResponse struct {
	Stocktakes           []*Stocktake     `protobuf:"bytes,1,rep,name=stocktakes" json:"stocktakes"`
	Paging               *common.PageInfo `protobuf:"bytes,2,opt,name=paging" json:"paging"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetStocktakesResponse) Reset()         { *m = GetStocktakesResponse{} }
func (m *GetStocktakesResponse) String() string { return proto.CompactTextString(m) }
func (*GetStocktakesResponse) ProtoMessage()    {}
func (*GetStocktakesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{197}
}

var xxx_messageInfo_GetStocktakesResponse proto.InternalMessageInfo

func (m *GetStocktakesResponse) GetStocktakes() []*Stocktake {
	if m != nil {
		return m.Stocktakes
	}
	return nil
}

func (m *GetStocktakesResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

type StocktakeLine struct {
	ProductId            dot.ID       `protobuf:"varint,9,opt,name=product_id,json=productId" json:"product_id"`
	ProductName          string       `protobuf:"bytes,10,opt,name=product_name,json=productName" json:"product_name"`
	VariantName          string       `protobuf:"bytes,4,opt,name=variant_name,json=variantName" json:"variant_name"`
	VariantId            dot.ID       `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	OldQuantity          int32        `protobuf:"varint,2,opt,name=old_quantity,json=oldQuantity" json:"old_quantity"`
	NewQuantity          int32        `protobuf:"varint,3,opt,name=new_quantity,json=newQuantity" json:"new_quantity"`
	Code                 string       `protobuf:"bytes,5,opt,name=code" json:"code"`
	ImageUrl             string       `protobuf:"bytes,6,opt,name=image_url,json=imageUrl" json:"image_url"`
	CostPrice            int32        `protobuf:"varint,8,opt,name=cost_price,json=costPrice" json:"cost_price"`
	Attributes           []*Attribute `protobuf:"bytes,7,rep,name=attributes" json:"attributes"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *StocktakeLine) Reset()         { *m = StocktakeLine{} }
func (m *StocktakeLine) String() string { return proto.CompactTextString(m) }
func (*StocktakeLine) ProtoMessage()    {}
func (*StocktakeLine) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{198}
}

var xxx_messageInfo_StocktakeLine proto.InternalMessageInfo

func (m *StocktakeLine) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *StocktakeLine) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *StocktakeLine) GetVariantName() string {
	if m != nil {
		return m.VariantName
	}
	return ""
}

func (m *StocktakeLine) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *StocktakeLine) GetOldQuantity() int32 {
	if m != nil {
		return m.OldQuantity
	}
	return 0
}

func (m *StocktakeLine) GetNewQuantity() int32 {
	if m != nil {
		return m.NewQuantity
	}
	return 0
}

func (m *StocktakeLine) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *StocktakeLine) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func (m *StocktakeLine) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func (m *StocktakeLine) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type ConfirmStocktakeRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	AutoInventoryVoucher string   `protobuf:"bytes,2,opt,name=auto_inventory_voucher,json=autoInventoryVoucher" json:"auto_inventory_voucher"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfirmStocktakeRequest) Reset()         { *m = ConfirmStocktakeRequest{} }
func (m *ConfirmStocktakeRequest) String() string { return proto.CompactTextString(m) }
func (*ConfirmStocktakeRequest) ProtoMessage()    {}
func (*ConfirmStocktakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{199}
}

var xxx_messageInfo_ConfirmStocktakeRequest proto.InternalMessageInfo

func (m *ConfirmStocktakeRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ConfirmStocktakeRequest) GetAutoInventoryVoucher() string {
	if m != nil {
		return m.AutoInventoryVoucher
	}
	return ""
}

type GetVariantsBySupplierIDRequest struct {
	SupplierId           dot.ID   `protobuf:"varint,1,opt,name=supplier_id,json=supplierId" json:"supplier_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetVariantsBySupplierIDRequest) Reset()         { *m = GetVariantsBySupplierIDRequest{} }
func (m *GetVariantsBySupplierIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetVariantsBySupplierIDRequest) ProtoMessage()    {}
func (*GetVariantsBySupplierIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{200}
}

var xxx_messageInfo_GetVariantsBySupplierIDRequest proto.InternalMessageInfo

func (m *GetVariantsBySupplierIDRequest) GetSupplierId() dot.ID {
	if m != nil {
		return m.SupplierId
	}
	return 0
}

type GetSuppliersByVariantIDRequest struct {
	VariantId            dot.ID   `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSuppliersByVariantIDRequest) Reset()         { *m = GetSuppliersByVariantIDRequest{} }
func (m *GetSuppliersByVariantIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetSuppliersByVariantIDRequest) ProtoMessage()    {}
func (*GetSuppliersByVariantIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{201}
}

var xxx_messageInfo_GetSuppliersByVariantIDRequest proto.InternalMessageInfo

func (m *GetSuppliersByVariantIDRequest) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

type CancelStocktakeRequest struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	CancelReason         string   `protobuf:"bytes,2,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelStocktakeRequest) Reset()         { *m = CancelStocktakeRequest{} }
func (m *CancelStocktakeRequest) String() string { return proto.CompactTextString(m) }
func (*CancelStocktakeRequest) ProtoMessage()    {}
func (*CancelStocktakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{202}
}

var xxx_messageInfo_CancelStocktakeRequest proto.InternalMessageInfo

func (m *CancelStocktakeRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelStocktakeRequest) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

type UpdateInventoryVariantCostPriceResponse struct {
	InventoryVariant     *InventoryVariant `protobuf:"bytes,2,opt,name=inventory_variant,json=inventoryVariant" json:"inventory_variant"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UpdateInventoryVariantCostPriceResponse) Reset() {
	*m = UpdateInventoryVariantCostPriceResponse{}
}
func (m *UpdateInventoryVariantCostPriceResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateInventoryVariantCostPriceResponse) ProtoMessage()    {}
func (*UpdateInventoryVariantCostPriceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{203}
}

var xxx_messageInfo_UpdateInventoryVariantCostPriceResponse proto.InternalMessageInfo

func (m *UpdateInventoryVariantCostPriceResponse) GetInventoryVariant() *InventoryVariant {
	if m != nil {
		return m.InventoryVariant
	}
	return nil
}

type UpdateInventoryVariantCostPriceRequest struct {
	VariantId            dot.ID   `protobuf:"varint,1,opt,name=variant_id,json=variantId" json:"variant_id"`
	CostPrice            int32    `protobuf:"varint,2,opt,name=cost_price,json=costPrice" json:"cost_price"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateInventoryVariantCostPriceRequest) Reset() {
	*m = UpdateInventoryVariantCostPriceRequest{}
}
func (m *UpdateInventoryVariantCostPriceRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateInventoryVariantCostPriceRequest) ProtoMessage()    {}
func (*UpdateInventoryVariantCostPriceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e3a4043938e9392, []int{204}
}

var xxx_messageInfo_UpdateInventoryVariantCostPriceRequest proto.InternalMessageInfo

func (m *UpdateInventoryVariantCostPriceRequest) GetVariantId() dot.ID {
	if m != nil {
		return m.VariantId
	}
	return 0
}

func (m *UpdateInventoryVariantCostPriceRequest) GetCostPrice() int32 {
	if m != nil {
		return m.CostPrice
	}
	return 0
}

func init() {
	proto.RegisterType((*InvitationsResponse)(nil), "shop.InvitationsResponse")
	proto.RegisterType((*GetInvitationsRequest)(nil), "shop.GetInvitationsRequest")
	proto.RegisterType((*CreateInvitationRequest)(nil), "shop.CreateInvitationRequest")
	proto.RegisterType((*PurchaseOrder)(nil), "shop.PurchaseOrder")
	proto.RegisterType((*PurchaseOrderSupplier)(nil), "shop.PurchaseOrderSupplier")
	proto.RegisterType((*GetPurchaseOrdersRequest)(nil), "shop.GetPurchaseOrdersRequest")
	proto.RegisterType((*PurchaseOrdersResponse)(nil), "shop.PurchaseOrdersResponse")
	proto.RegisterType((*CreatePurchaseOrderRequest)(nil), "shop.CreatePurchaseOrderRequest")
	proto.RegisterType((*UpdatePurchaseOrderRequest)(nil), "shop.UpdatePurchaseOrderRequest")
	proto.RegisterType((*PurchaseOrderLine)(nil), "shop.PurchaseOrderLine")
	proto.RegisterType((*CancelPurchaseOrderRequest)(nil), "shop.CancelPurchaseOrderRequest")
	proto.RegisterType((*ConfirmPurchaseOrderRequest)(nil), "shop.ConfirmPurchaseOrderRequest")
	proto.RegisterType((*GetLedgersRequest)(nil), "shop.GetLedgersRequest")
	proto.RegisterType((*CreateLedgerRequest)(nil), "shop.CreateLedgerRequest")
	proto.RegisterType((*UpdateLedgerRequest)(nil), "shop.UpdateLedgerRequest")
	proto.RegisterType((*LedgersResponse)(nil), "shop.LedgersResponse")
	proto.RegisterType((*Ledger)(nil), "shop.Ledger")
	proto.RegisterType((*RegisterShopRequest)(nil), "shop.RegisterShopRequest")
	proto.RegisterType((*RegisterShopResponse)(nil), "shop.RegisterShopResponse")
	proto.RegisterType((*UpdateShopRequest)(nil), "shop.UpdateShopRequest")
	proto.RegisterType((*UpdateShopResponse)(nil), "shop.UpdateShopResponse")
	proto.RegisterType((*Collection)(nil), "shop.Collection")
	proto.RegisterType((*CreateCollectionRequest)(nil), "shop.CreateCollectionRequest")
	proto.RegisterType((*UpdateProductCategoryRequest)(nil), "shop.UpdateProductCategoryRequest")
	proto.RegisterType((*CollectionsResponse)(nil), "shop.CollectionsResponse")
	proto.RegisterType((*UpdateCollectionRequest)(nil), "shop.UpdateCollectionRequest")
	proto.RegisterType((*UpdateProductsCollectionRequest)(nil), "shop.UpdateProductsCollectionRequest")
	proto.RegisterType((*RemoveProductsCollectionRequest)(nil), "shop.RemoveProductsCollectionRequest")
	proto.RegisterType((*EtopVariant)(nil), "shop.EtopVariant")
	proto.RegisterType((*EtopProduct)(nil), "shop.EtopProduct")
	proto.RegisterType((*ShopVariant)(nil), "shop.ShopVariant")
	proto.RegisterType((*InventoryVariantShopVariant)(nil), "shop.InventoryVariantShopVariant")
	proto.RegisterType((*ShopProduct)(nil), "shop.ShopProduct")
	proto.RegisterType((*ShopShortProduct)(nil), "shop.ShopShortProduct")
	proto.RegisterType((*ShopCollection)(nil), "shop.ShopCollection")
	proto.RegisterType((*GetVariantsRequest)(nil), "shop.GetVariantsRequest")
	proto.RegisterType((*GetCategoriesRequest)(nil), "shop.GetCategoriesRequest")
	proto.RegisterType((*ShopVariantsResponse)(nil), "shop.ShopVariantsResponse")
	proto.RegisterType((*ShopProductsResponse)(nil), "shop.ShopProductsResponse")
	proto.RegisterType((*ShopCategoriesResponse)(nil), "shop.ShopCategoriesResponse")
	proto.RegisterType((*UpdateVariantRequest)(nil), "shop.UpdateVariantRequest")
	proto.RegisterType((*UpdateProductRequest)(nil), "shop.UpdateProductRequest")
	proto.RegisterType((*UpdateCategoryRequest)(nil), "shop.UpdateCategoryRequest")
	proto.RegisterType((*UpdateVariantsRequest)(nil), "shop.UpdateVariantsRequest")
	proto.RegisterType((*UpdateProductsTagsRequest)(nil), "shop.UpdateProductsTagsRequest")
	proto.RegisterType((*UpdateVariantsResponse)(nil), "shop.UpdateVariantsResponse")
	proto.RegisterType((*AddVariantsRequest)(nil), "shop.AddVariantsRequest")
	proto.RegisterType((*AddVariantsResponse)(nil), "shop.AddVariantsResponse")
	proto.RegisterType((*RemoveVariantsRequest)(nil), "shop.RemoveVariantsRequest")
	proto.RegisterType((*GetOrdersRequest)(nil), "shop.GetOrdersRequest")
	proto.RegisterType((*UpdateOrdersStatusRequest)(nil), "shop.UpdateOrdersStatusRequest")
	proto.RegisterType((*ConfirmOrderRequest)(nil), "shop.ConfirmOrderRequest")
	proto.RegisterType((*OrderIDRequest)(nil), "shop.OrderIDRequest")
	proto.RegisterType((*OrderIDsRequest)(nil), "shop.OrderIDsRequest")
	proto.RegisterType((*CreateFulfillmentsForOrderRequest)(nil), "shop.CreateFulfillmentsForOrderRequest")
	proto.RegisterType((*CancelOrderRequest)(nil), "shop.CancelOrderRequest")
	proto.RegisterType((*CancelOrdersRequest)(nil), "shop.CancelOrdersRequest")
	proto.RegisterType((*ProductSource)(nil), "shop.ProductSource")
	proto.RegisterType((*CreateProductSourceRequest)(nil), "shop.CreateProductSourceRequest")
	proto.RegisterType((*ProductSourcesResponse)(nil), "shop.ProductSourcesResponse")
	proto.RegisterType((*CreateCategoryRequest)(nil), "shop.CreateCategoryRequest")
	proto.RegisterType((*CreateProductRequest)(nil), "shop.CreateProductRequest")
	proto.RegisterType((*CreateVariantRequest)(nil), "shop.CreateVariantRequest")
	proto.RegisterType((*DeprecatedCreateVariantRequest)(nil), "shop.DeprecatedCreateVariantRequest")
	proto.RegisterType((*ConnectProductSourceResquest)(nil), "shop.ConnectProductSourceResquest")
	proto.RegisterType((*CreatePSCategoryRequest)(nil), "shop.CreatePSCategoryRequest")
	proto.RegisterType((*UpdateProductsPSCategoryRequest)(nil), "shop.UpdateProductsPSCategoryRequest")
	proto.RegisterType((*UpdateProductsCollectionResponse)(nil), "shop.UpdateProductsCollectionResponse")
	proto.RegisterType((*UpdateProductSourceCategoryRequest)(nil), "shop.UpdateProductSourceCategoryRequest")
	proto.RegisterType((*GetProductSourceCategoriesRequest)(nil), "shop.GetProductSourceCategoriesRequest")
	proto.RegisterType((*GetFulfillmentsRequest)(nil), "shop.GetFulfillmentsRequest")
	proto.RegisterType((*GetFulfillmentHistoryRequest)(nil), "shop.GetFulfillmentHistoryRequest")
	proto.RegisterType((*GetBalanceShopResponse)(nil), "shop.GetBalanceShopResponse")
	proto.RegisterType((*GetMoneyTransactionsRequest)(nil), "shop.GetMoneyTransactionsRequest")
	proto.RegisterType((*GetPublicFulfillmentRequest)(nil), "shop.GetPublicFulfillmentRequest")
	proto.RegisterType((*UpdateFulfillmentsShippingStateRequest)(nil), "shop.UpdateFulfillmentsShippingStateRequest")
	proto.RegisterType((*UpdateOrderPaymentStatusRequest)(nil), "shop.UpdateOrderPaymentStatusRequest")
	proto.RegisterType((*SummarizeFulfillmentsRequest)(nil), "shop.SummarizeFulfillmentsRequest")
	proto.RegisterType((*SummarizeFulfillmentsResponse)(nil), "shop.SummarizeFulfillmentsResponse")
	proto.RegisterType((*SummarizePOSResponse)(nil), "shop.SummarizePOSResponse")
	proto.RegisterType((*SummarizePOSRequest)(nil), "shop.SummarizePOSRequest")
	proto.RegisterType((*SummaryTable)(nil), "shop.SummaryTable")
	proto.RegisterType((*SummaryColRow)(nil), "shop.SummaryColRow")
	proto.RegisterType((*SummaryItem)(nil), "shop.SummaryItem")
	proto.RegisterType((*ImportProductsResponse)(nil), "shop.ImportProductsResponse")
	proto.RegisterType((*CalcBalanceShopResponse)(nil), "shop.CalcBalanceShopResponse")
	proto.RegisterType((*RequestExportRequest)(nil), "shop.RequestExportRequest")
	proto.RegisterType((*RequestExportResponse)(nil), "shop.RequestExportResponse")
	proto.RegisterType((*GetExportsRequest)(nil), "shop.GetExportsRequest")
	proto.RegisterType((*GetExportsResponse)(nil), "shop.GetExportsResponse")
	proto.RegisterType((*ExportItem)(nil), "shop.ExportItem")
	proto.RegisterType((*GetExportsStatusRequest)(nil), "shop.GetExportsStatusRequest")
	proto.RegisterType((*ExportStatusItem)(nil), "shop.ExportStatusItem")
	proto.RegisterType((*AuthorizePartnerRequest)(nil), "shop.AuthorizePartnerRequest")
	proto.RegisterType((*GetPartnersResponse)(nil), "shop.GetPartnersResponse")
	proto.RegisterType((*AuthorizedPartnerResponse)(nil), "shop.AuthorizedPartnerResponse")
	proto.RegisterType((*GetAuthorizedPartnersResponse)(nil), "shop.GetAuthorizedPartnersResponse")
	proto.RegisterType((*Attribute)(nil), "shop.Attribute")
	proto.RegisterType((*UpdateVariantImagesRequest)(nil), "shop.UpdateVariantImagesRequest")
	proto.RegisterType((*UpdateProductMetaFieldsRequest)(nil), "shop.UpdateProductMetaFieldsRequest")
	proto.RegisterType((*CategoriesResponse)(nil), "shop.CategoriesResponse")
	proto.RegisterType((*Category)(nil), "shop.Category")
	proto.RegisterType((*Tag)(nil), "shop.Tag")
	proto.RegisterType((*ExternalAccountAhamove)(nil), "shop.ExternalAccountAhamove")
	proto.RegisterType((*UpdateXAccountAhamoveVerificationRequest)(nil), "shop.UpdateXAccountAhamoveVerificationRequest")
	proto.RegisterType((*ExternalAccountHaravan)(nil), "shop.ExternalAccountHaravan")
	proto.RegisterType((*ExternalAccountHaravanRequest)(nil), "shop.ExternalAccountHaravanRequest")
	proto.RegisterType((*CustomerLiability)(nil), "shop.CustomerLiability")
	proto.RegisterType((*Customer)(nil), "shop.Customer")
	proto.RegisterType((*CreateCustomerRequest)(nil), "shop.CreateCustomerRequest")
	proto.RegisterType((*UpdateCustomerRequest)(nil), "shop.UpdateCustomerRequest")
	proto.RegisterType((*GetCustomersRequest)(nil), "shop.GetCustomersRequest")
	proto.RegisterType((*CustomersResponse)(nil), "shop.CustomersResponse")
	proto.RegisterType((*SetCustomersStatusRequest)(nil), "shop.SetCustomersStatusRequest")
	proto.RegisterType((*CustomerDetailsResponse)(nil), "shop.CustomerDetailsResponse")
	proto.RegisterType((*IndependentSummaryItem)(nil), "shop.IndependentSummaryItem")
	proto.RegisterType((*GetCustomerAddressesRequest)(nil), "shop.GetCustomerAddressesRequest")
	proto.RegisterType((*CustomerAddress)(nil), "shop.CustomerAddress")
	proto.RegisterType((*CreateCustomerAddressRequest)(nil), "shop.CreateCustomerAddressRequest")
	proto.RegisterType((*UpdateCustomerAddressRequest)(nil), "shop.UpdateCustomerAddressRequest")
	proto.RegisterType((*CustomerAddressesResponse)(nil), "shop.CustomerAddressesResponse")
	proto.RegisterType((*UpdateProductStatusRequest)(nil), "shop.UpdateProductStatusRequest")
	proto.RegisterType((*UpdateProductStatusResponse)(nil), "shop.UpdateProductStatusResponse")
	proto.RegisterType((*PaymentTradingOrderRequest)(nil), "shop.PaymentTradingOrderRequest")
	proto.RegisterType((*PaymentTradingOrderResponse)(nil), "shop.PaymentTradingOrderResponse")
	proto.RegisterType((*UpdateVariantAttributesRequest)(nil), "shop.UpdateVariantAttributesRequest")
	proto.RegisterType((*PaymentCheckReturnDataRequest)(nil), "shop.PaymentCheckReturnDataRequest")
	proto.RegisterType((*ShopCategory)(nil), "shop.ShopCategory")
	proto.RegisterType((*GetCollectionsRequest)(nil), "shop.GetCollectionsRequest")
	proto.RegisterType((*ShopCollectionsResponse)(nil), "shop.ShopCollectionsResponse")
	proto.RegisterType((*AddShopProductCollectionRequest)(nil), "shop.AddShopProductCollectionRequest")
	proto.RegisterType((*RemoveShopProductCollectionRequest)(nil), "shop.RemoveShopProductCollectionRequest")
	proto.RegisterType((*AddCustomerToGroupRequest)(nil), "shop.AddCustomerToGroupRequest")
	proto.RegisterType((*RemoveCustomerOutOfGroupRequest)(nil), "shop.RemoveCustomerOutOfGroupRequest")
	proto.RegisterType((*SupplierLiability)(nil), "shop.SupplierLiability")
	proto.RegisterType((*Supplier)(nil), "shop.Supplier")
	proto.RegisterType((*CreateSupplierRequest)(nil), "shop.CreateSupplierRequest")
	proto.RegisterType((*UpdateSupplierRequest)(nil), "shop.UpdateSupplierRequest")
	proto.RegisterType((*GetSuppliersRequest)(nil), "shop.GetSuppliersRequest")
	proto.RegisterType((*SuppliersResponse)(nil), "shop.SuppliersResponse")
	proto.RegisterType((*Carrier)(nil), "shop.Carrier")
	proto.RegisterType((*CreateCarrierRequest)(nil), "shop.CreateCarrierRequest")
	proto.RegisterType((*UpdateCarrierRequest)(nil), "shop.UpdateCarrierRequest")
	proto.RegisterType((*GetCarriersRequest)(nil), "shop.GetCarriersRequest")
	proto.RegisterType((*CarriersResponse)(nil), "shop.CarriersResponse")
	proto.RegisterType((*ReceiptLine)(nil), "shop.ReceiptLine")
	proto.RegisterType((*Trader)(nil), "shop.Trader")
	proto.RegisterType((*Receipt)(nil), "shop.Receipt")
	proto.RegisterType((*CreateReceiptRequest)(nil), "shop.CreateReceiptRequest")
	proto.RegisterType((*UpdateReceiptRequest)(nil), "shop.UpdateReceiptRequest")
	proto.RegisterType((*CancelReceiptRequest)(nil), "shop.CancelReceiptRequest")
	proto.RegisterType((*GetReceiptsRequest)(nil), "shop.GetReceiptsRequest")
	proto.RegisterType((*GetReceiptsByLedgerTypeRequest)(nil), "shop.GetReceiptsByLedgerTypeRequest")
	proto.RegisterType((*ReceiptsResponse)(nil), "shop.ReceiptsResponse")
	proto.RegisterType((*GetShopCollectionsByProductIDRequest)(nil), "shop.GetShopCollectionsByProductIDRequest")
	proto.RegisterType((*CreateInventoryVoucherRequest)(nil), "shop.CreateInventoryVoucherRequest")
	proto.RegisterType((*InventoryVoucherLine)(nil), "shop.InventoryVoucherLine")
	proto.RegisterType((*CreateInventoryVoucherResponse)(nil), "shop.CreateInventoryVoucherResponse")
	proto.RegisterType((*ConfirmInventoryVoucherRequest)(nil), "shop.ConfirmInventoryVoucherRequest")
	proto.RegisterType((*ConfirmInventoryVoucherResponse)(nil), "shop.ConfirmInventoryVoucherResponse")
	proto.RegisterType((*CancelInventoryVoucherRequest)(nil), "shop.CancelInventoryVoucherRequest")
	proto.RegisterType((*CancelInventoryVoucherResponse)(nil), "shop.CancelInventoryVoucherResponse")
	proto.RegisterType((*UpdateInventoryVoucherRequest)(nil), "shop.UpdateInventoryVoucherRequest")
	proto.RegisterType((*UpdateInventoryVoucherResponse)(nil), "shop.UpdateInventoryVoucherResponse")
	proto.RegisterType((*AdjustInventoryQuantityRequest)(nil), "shop.AdjustInventoryQuantityRequest")
	proto.RegisterType((*AdjustInventoryQuantityResponse)(nil), "shop.AdjustInventoryQuantityResponse")
	proto.RegisterType((*GetInventoryVariantsRequest)(nil), "shop.GetInventoryVariantsRequest")
	proto.RegisterType((*GetInventoryVariantsResponse)(nil), "shop.GetInventoryVariantsResponse")
	proto.RegisterType((*GetInventoryVariantsByVariantIDsRequest)(nil), "shop.GetInventoryVariantsByVariantIDsRequest")
	proto.RegisterType((*InventoryVariant)(nil), "shop.InventoryVariant")
	proto.RegisterType((*InventoryVariantQuantity)(nil), "shop.InventoryVariantQuantity")
	proto.RegisterType((*InventoryVoucher)(nil), "shop.InventoryVoucher")
	proto.RegisterType((*GetInventoryVouchersRequest)(nil), "shop.GetInventoryVouchersRequest")
	proto.RegisterType((*GetInventoryVouchersByIDsRequest)(nil), "shop.GetInventoryVouchersByIDsRequest")
	proto.RegisterType((*GetInventoryVouchersResponse)(nil), "shop.GetInventoryVouchersResponse")
	proto.RegisterType((*CustomerGroup)(nil), "shop.CustomerGroup")
	proto.RegisterType((*CreateCustomerGroupRequest)(nil), "shop.CreateCustomerGroupRequest")
	proto.RegisterType((*UpdateCustomerGroupRequest)(nil), "shop.UpdateCustomerGroupRequest")
	proto.RegisterType((*GetCustomerGroupsRequest)(nil), "shop.GetCustomerGroupsRequest")
	proto.RegisterType((*CustomerGroupsResponse)(nil), "shop.CustomerGroupsResponse")
	proto.RegisterType((*GetOrdersByReceiptIDRequest)(nil), "shop.GetOrdersByReceiptIDRequest")
	proto.RegisterType((*GetInventoryVariantRequest)(nil), "shop.GetInventoryVariantRequest")
	proto.RegisterType((*CreateBrandRequest)(nil), "shop.CreateBrandRequest")
	proto.RegisterType((*UpdateBrandRequest)(nil), "shop.UpdateBrandRequest")
	proto.RegisterType((*DeleteBrandResponse)(nil), "shop.DeleteBrandResponse")
	proto.RegisterType((*Brand)(nil), "shop.Brand")
	proto.RegisterType((*GetBrandsByIDsResponse)(nil), "shop.GetBrandsByIDsResponse")
	proto.RegisterType((*GetBrandsRequest)(nil), "shop.GetBrandsRequest")
	proto.RegisterType((*GetBrandsResponse)(nil), "shop.GetBrandsResponse")
	proto.RegisterType((*GetInventoryVouchersByReferenceRequest)(nil), "shop.GetInventoryVouchersByReferenceRequest")
	proto.RegisterType((*GetInventoryVouchersByReferenceResponse)(nil), "shop.GetInventoryVouchersByReferenceResponse")
	proto.RegisterType((*UpdateOrderShippingInfoRequest)(nil), "shop.UpdateOrderShippingInfoRequest")
	proto.RegisterType((*GetStocktakesByIDsResponse)(nil), "shop.GetStocktakesByIDsResponse")
	proto.RegisterType((*CreateStocktakeRequest)(nil), "shop.CreateStocktakeRequest")
	proto.RegisterType((*UpdateStocktakeRequest)(nil), "shop.UpdateStocktakeRequest")
	proto.RegisterType((*Stocktake)(nil), "shop.Stocktake")
	proto.RegisterType((*GetStocktakesRequest)(nil), "shop.GetStocktakesRequest")
	proto.RegisterType((*GetStocktakesResponse)(nil), "shop.GetStocktakesResponse")
	proto.RegisterType((*StocktakeLine)(nil), "shop.StocktakeLine")
	proto.RegisterType((*ConfirmStocktakeRequest)(nil), "shop.ConfirmStocktakeRequest")
	proto.RegisterType((*GetVariantsBySupplierIDRequest)(nil), "shop.GetVariantsBySupplierIDRequest")
	proto.RegisterType((*GetSuppliersByVariantIDRequest)(nil), "shop.GetSuppliersByVariantIDRequest")
	proto.RegisterType((*CancelStocktakeRequest)(nil), "shop.CancelStocktakeRequest")
	proto.RegisterType((*UpdateInventoryVariantCostPriceResponse)(nil), "shop.UpdateInventoryVariantCostPriceResponse")
	proto.RegisterType((*UpdateInventoryVariantCostPriceRequest)(nil), "shop.UpdateInventoryVariantCostPriceRequest")
}

func init() { proto.RegisterFile("etop/shop/shop.proto", fileDescriptor_1e3a4043938e9392) }

var fileDescriptor_1e3a4043938e9392 = []byte{
	// 8224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe4, 0x7d, 0x5b, 0x6c, 0x24, 0xc9,
	0x91, 0xd8, 0xf4, 0xbb, 0x3b, 0x9a, 0xcf, 0xe2, 0xab, 0x87, 0x1c, 0x72, 0xc8, 0x9a, 0x9d, 0x9d,
	0x59, 0xad, 0x87, 0xdc, 0x9d, 0x9d, 0x5d, 0xbd, 0x56, 0x32, 0x48, 0xce, 0xee, 0x2c, 0x57, 0xfb,
	0x18, 0x35, 0x67, 0x47, 0xb6, 0x21, 0xa9, 0x55, 0xec, 0x4a, 0x36, 0x4b, 0xac, 0xee, 0xea, 0xad,
	0xaa, 0xe6, 0x0c, 0x0d, 0xf9, 0xc3, 0xb0, 0x21, 0xc8, 0xf6, 0x87, 0x64, 0x40, 0x86, 0x04, 0x58,
	0x10, 0x24, 0xd8, 0x96, 0x0d, 0xc3, 0x0f, 0xc0, 0x96, 0x01, 0x7d, 0xd8, 0x1f, 0x82, 0xfc, 0xa1,
	0x3b, 0x1c, 0x20, 0xe1, 0x70, 0xba, 0xfb, 0x3b, 0x68, 0x57, 0xc2, 0x7d, 0x1c, 0x70, 0x82, 0x70,
	0x10, 0x70, 0x38, 0x1c, 0x70, 0x38, 0xe4, 0xb3, 0x32, 0xb3, 0xaa, 0xba, 0xab, 0x39, 0x9c, 0xd9,
	0x13, 0xf4, 0x33, 0xc3, 0xce, 0x8c, 0xca, 0x47, 0x64, 0x44, 0x64, 0x44, 0x64, 0x64, 0x24, 0xcc,
	0xa3, 0xd0, 0xeb, 0x6f, 0x05, 0x47, 0xec, 0x9f, 0xcd, 0xbe, 0xef, 0x85, 0x9e, 0x51, 0xc4, 0x7f,
	0x2f, 0xcf, 0x77, 0xbc, 0x8e, 0x47, 0x0a, 0xb6, 0xf0, 0x5f, 0xb4, 0x6e, 0xf9, 0x72, 0xc7, 0xf3,
	0x3a, 0x2e, 0xda, 0x22, 0xbf, 0x0e, 0x06, 0x87, 0x5b, 0xa1, 0xd3, 0x45, 0x41, 0x68, 0x75, 0xd9,
	0xc7, 0xcb, 0x73, 0x6d, 0xaf, 0xdb, 0xf5, 0x7a, 0x5b, 0xf4, 0x3f, 0x56, 0xf8, 0x14, 0x2b, 0x0c,
	0xfa, 0x3e, 0xb2, 0xec, 0xe0, 0x08, 0xa1, 0x50, 0xfe, 0x9b, 0x41, 0x6d, 0x90, 0xd1, 0xa0, 0xb0,
	0xbd, 0xd5, 0x39, 0xea, 0xb5, 0x7a, 0x5e, 0x88, 0x5a, 0x6d, 0xcf, 0x46, 0x5b, 0xf8, 0x1f, 0x06,
	0xb2, 0x2c, 0x40, 0x42, 0xff, 0xb4, 0x45, 0xba, 0x11, 0x75, 0xab, 0xa2, 0x2e, 0x38, 0x72, 0xfa,
	0x7d, 0xa7, 0xd7, 0xd9, 0x0a, 0x42, 0x2b, 0x4c, 0xa8, 0x0e, 0xad, 0x70, 0x10, 0xbc, 0xc0, 0xfe,
	0x4f, 0xa9, 0xbe, 0xa5, 0x56, 0x4f, 0xb3, 0x6a, 0x8e, 0xa4, 0xe5, 0x45, 0x52, 0xe0, 0xf9, 0x36,
	0xf2, 0xe9, 0xbf, 0xac, 0xfc, 0x39, 0xd1, 0x4e, 0xdf, 0x3a, 0xed, 0xa2, 0x5e, 0xd8, 0xea, 0xfb,
	0xde, 0x89, 0x83, 0xc1, 0xf4, 0x02, 0xf6, 0xc5, 0xf5, 0xe8, 0x0b, 0xdf, 0xb3, 0x07, 0xed, 0xb0,
	0x15, 0x9e, 0xf6, 0x91, 0xf2, 0x83, 0x42, 0x9a, 0x1e, 0xcc, 0xed, 0xf5, 0x4e, 0x9c, 0xd0, 0x0a,
	0x1d, 0xaf, 0x17, 0x34, 0x51, 0xd0, 0xf7, 0x7a, 0x01, 0x32, 0x6e, 0x42, 0xdd, 0x89, 0x8a, 0x1b,
	0xb9, 0xf5, 0xc2, 0xf5, 0xfa, 0xcd, 0x99, 0x4d, 0x32, 0xd8, 0x08, 0xbe, 0x29, 0x03, 0x19, 0x4f,
	0x41, 0xb9, 0x6f, 0x75, 0x9c, 0x5e, 0xa7, 0x91, 0x5f, 0xcf, 0x5d, 0xaf, 0xdf, 0x9c, 0xd8, 0x6c,
	0x77, 0x37, 0xef, 0x5a, 0x1d, 0xb4, 0xd7, 0x3b, 0xf4, 0x9a, 0xac, 0xce, 0xb4, 0x60, 0xe1, 0x0e,
	0x0a, 0x95, 0x3e, 0xdf, 0x1d, 0xa0, 0x20, 0x34, 0x4c, 0xf1, 0x79, 0x8e, 0x7c, 0x0e, 0xec, 0x73,
	0xa7, 0xd7, 0xe1, 0x1f, 0x1b, 0x4f, 0x41, 0xe5, 0xd0, 0x71, 0x43, 0xe4, 0x07, 0x8d, 0x3c, 0x19,
	0x12, 0x01, 0x7a, 0x95, 0x14, 0x35, 0x79, 0x95, 0xf9, 0x69, 0x58, 0xda, 0xf5, 0x91, 0x15, 0x22,
	0x69, 0xa4, 0xac, 0x93, 0x65, 0x28, 0xa1, 0xae, 0xe5, 0xb8, 0xa4, 0x8f, 0xda, 0x4e, 0xf1, 0xc7,
	0x7f, 0x7a, 0xf9, 0x42, 0x93, 0x16, 0xe1, 0x3a, 0xdf, 0x73, 0x11, 0x6d, 0x5a, 0xd4, 0x91, 0x22,
	0xf3, 0x7b, 0x15, 0x98, 0xbc, 0x3b, 0xf0, 0xdb, 0x47, 0x56, 0x80, 0xde, 0xc6, 0x4b, 0x63, 0xcc,
	0x43, 0xde, 0xb1, 0x49, 0x33, 0x05, 0x06, 0x9a, 0x77, 0x6c, 0x63, 0x15, 0x2a, 0x98, 0xd2, 0x5b,
	0x8e, 0x4d, 0x90, 0xc0, 0xab, 0xca, 0xb8, 0x70, 0xcf, 0x36, 0xae, 0x42, 0x3d, 0x18, 0xf4, 0xfb,
	0xae, 0x83, 0x7c, 0x0c, 0x52, 0x90, 0x40, 0x80, 0x57, 0xec, 0xd9, 0xc6, 0x87, 0xa1, 0xca, 0x7f,
	0x35, 0x8a, 0x04, 0x19, 0x2b, 0x9b, 0x84, 0x99, 0x94, 0x21, 0xec, 0x33, 0x90, 0xa6, 0x00, 0x36,
	0xae, 0xc1, 0xc4, 0x81, 0x15, 0x1c, 0xa3, 0xb0, 0x75, 0x62, 0xb9, 0x03, 0xd4, 0x28, 0x49, 0x1d,
	0xd4, 0x69, 0xcd, 0x7d, 0x5c, 0x61, 0x3c, 0x0b, 0x53, 0xa1, 0x17, 0x5a, 0x6e, 0xcb, 0x76, 0x82,
	0xb6, 0x37, 0xe8, 0x85, 0x8d, 0xb2, 0x04, 0x3a, 0x49, 0xea, 0x6e, 0xb3, 0x2a, 0xdc, 0x2a, 0x05,
	0xb6, 0xba, 0x04, 0xb4, 0x22, 0xb7, 0x4a, 0x6a, 0xb6, 0x49, 0x85, 0xd1, 0x80, 0x22, 0x66, 0x9e,
	0x46, 0x55, 0x42, 0x2e, 0x29, 0xc1, 0x35, 0x98, 0xf9, 0x1a, 0x35, 0xb9, 0x06, 0x97, 0x18, 0x37,
	0xa0, 0x4c, 0xb9, 0xa2, 0x01, 0xeb, 0xb9, 0xeb, 0x53, 0x37, 0xa7, 0x37, 0x19, 0x2f, 0x6d, 0xee,
	0x93, 0xff, 0x05, 0x06, 0xc9, 0x2f, 0xe3, 0x06, 0x94, 0x5c, 0xa7, 0x87, 0x82, 0x46, 0x9d, 0xac,
	0xff, 0x52, 0x02, 0x5e, 0xde, 0x70, 0x7a, 0xa8, 0x49, 0xa1, 0x8c, 0x2b, 0x00, 0x6d, 0x42, 0x0a,
	0x76, 0xeb, 0xe0, 0xb4, 0x31, 0x21, 0x0d, 0xbc, 0xc6, 0xca, 0x77, 0x4e, 0x8d, 0x2d, 0x98, 0x69,
	0x5b, 0xbd, 0x36, 0x72, 0x5d, 0x64, 0xb7, 0x7c, 0x64, 0x05, 0x5e, 0xaf, 0x31, 0x29, 0x0d, 0x74,
	0x5a, 0xd4, 0x36, 0x49, 0xa5, 0xf1, 0x09, 0x98, 0x68, 0x7b, 0xbd, 0x43, 0xc7, 0xef, 0x22, 0xbb,
	0x65, 0x85, 0x8d, 0x29, 0xb2, 0x46, 0xcb, 0x9b, 0x54, 0x90, 0x6d, 0x72, 0x41, 0xb6, 0x79, 0x8f,
	0x0b, 0xb2, 0x66, 0x5d, 0xc0, 0x6f, 0x87, 0xe4, 0x73, 0xd1, 0x9f, 0x15, 0x36, 0xa6, 0x33, 0x7c,
	0xce, 0xe1, 0xb7, 0x43, 0xe3, 0xa3, 0xd1, 0x9c, 0xac, 0xb0, 0x31, 0x33, 0xf2, 0x63, 0x3e, 0x53,
	0xfa, 0xe9, 0xa0, 0x6f, 0xf3, 0x4f, 0x67, 0x47, 0x7f, 0xca, 0xa0, 0xe9, 0xa7, 0x36, 0x72, 0x11,
	0xfb, 0xd4, 0x18, 0xfd, 0x29, 0x83, 0xde, 0x0e, 0x8d, 0x5d, 0x98, 0x75, 0x7a, 0x27, 0xa8, 0x17,
	0x7a, 0xfe, 0x69, 0xeb, 0xc4, 0x1b, 0xb4, 0x8f, 0x90, 0xdf, 0x98, 0x23, 0x2d, 0x2c, 0xd2, 0xf5,
	0xdb, 0xe3, 0xd5, 0xf7, 0x69, 0x6d, 0x73, 0xc6, 0xd1, 0x4a, 0x30, 0xeb, 0xf4, 0x2d, 0xc7, 0xe6,
	0x34, 0x38, 0x2f, 0xb3, 0x0e, 0xae, 0xa0, 0x24, 0x68, 0x7e, 0x23, 0x0f, 0x0b, 0x89, 0x5c, 0x62,
	0x6c, 0x40, 0xed, 0x70, 0xe0, 0xba, 0xad, 0x9e, 0xd5, 0x45, 0x0a, 0xfb, 0x57, 0x71, 0xf1, 0x5b,
	0x56, 0x17, 0x61, 0x09, 0xd0, 0x3f, 0xf2, 0x7a, 0x88, 0xf0, 0xae, 0x90, 0x00, 0xa4, 0x28, 0x92,
	0x1c, 0x85, 0xb8, 0xe4, 0xb8, 0x86, 0xe9, 0xa1, 0xdb, 0xb7, 0x7a, 0xa7, 0xb4, 0xf5, 0xa2, 0x04,
	0x52, 0x67, 0x35, 0xa4, 0x83, 0x2b, 0x00, 0xa1, 0xf5, 0xb0, 0xd5, 0x1b, 0x74, 0x0f, 0x90, 0x4f,
	0xb8, 0x93, 0x83, 0xd5, 0x42, 0xeb, 0xe1, 0x5b, 0xa4, 0xd8, 0x78, 0x11, 0xe6, 0x8e, 0x90, 0x65,
	0xbf, 0x3b, 0xb0, 0xfc, 0x10, 0xf9, 0x2d, 0xcb, 0xb6, 0x7d, 0x14, 0x04, 0x84, 0x41, 0x39, 0xb4,
	0x21, 0x01, 0x6c, 0xd3, 0x7a, 0x63, 0x0d, 0x2a, 0x0c, 0xe5, 0x84, 0x41, 0xab, 0x0c, 0x94, 0x17,
	0x9a, 0x36, 0x34, 0xee, 0xa0, 0x50, 0xc1, 0xcd, 0x63, 0x90, 0xbd, 0x5f, 0x82, 0x45, 0xbd, 0x0b,
	0xb6, 0xa5, 0xbc, 0x0c, 0xd3, 0x7d, 0x56, 0xd3, 0x22, 0xbb, 0x1b, 0xdf, 0x56, 0xe6, 0x12, 0x78,
	0xb8, 0x39, 0xd5, 0x57, 0x5a, 0xc9, 0xb8, 0xb9, 0xfc, 0x9b, 0x3c, 0x2c, 0x53, 0xd1, 0xaf, 0xb6,
	0xc6, 0xa6, 0xa9, 0x89, 0xdf, 0x5c, 0x8a, 0xf8, 0xd5, 0xa5, 0x68, 0x3e, 0xbb, 0x14, 0x2d, 0x64,
	0x97, 0xa2, 0xc5, 0x21, 0x52, 0x94, 0xc8, 0xca, 0x52, 0x82, 0xac, 0x64, 0xc2, 0xaf, 0x9c, 0x45,
	0xf8, 0x99, 0xff, 0x22, 0x0f, 0xcb, 0xef, 0x10, 0x06, 0x4e, 0xc4, 0x46, 0xf2, 0x0e, 0x76, 0x59,
	0xc5, 0x11, 0x99, 0xbb, 0x82, 0x9d, 0x0d, 0x0d, 0x3b, 0x64, 0xca, 0x2a, 0x5e, 0xae, 0xc6, 0xf0,
	0x42, 0x26, 0xab, 0x63, 0x64, 0x43, 0xc3, 0x48, 0x89, 0xb6, 0x24, 0xe3, 0xc2, 0x60, 0xb8, 0x20,
	0xc4, 0xaf, 0x63, 0xa1, 0x92, 0x09, 0x0b, 0x3f, 0xcc, 0xc3, 0x6c, 0xac, 0x12, 0x73, 0xe2, 0x89,
	0xe5, 0x3b, 0x56, 0x2f, 0xd4, 0x29, 0xa1, 0xc6, 0xca, 0xf7, 0x6c, 0x63, 0x1d, 0xaa, 0xef, 0x0e,
	0xac, 0x5e, 0xe8, 0x84, 0xa7, 0x0a, 0x11, 0x88, 0x52, 0xe3, 0x19, 0x98, 0x8c, 0x54, 0x30, 0xa7,
	0x8d, 0x14, 0x02, 0x98, 0x60, 0x55, 0x77, 0x71, 0x0d, 0xee, 0x91, 0xeb, 0x5f, 0x8e, 0xad, 0xac,
	0x7e, 0x8d, 0x95, 0x53, 0xd2, 0xe3, 0x40, 0x44, 0x92, 0x54, 0x64, 0x49, 0xc2, 0x6a, 0x88, 0x24,
	0xd9, 0x80, 0x9a, 0xd3, 0xb5, 0x3a, 0xa8, 0x35, 0xf0, 0x5d, 0x85, 0x52, 0xaa, 0xa4, 0xf8, 0x1d,
	0xdf, 0x15, 0xbb, 0x71, 0x2d, 0xb6, 0x1b, 0x6f, 0x01, 0x58, 0x61, 0xe8, 0x3b, 0x07, 0x83, 0x50,
	0x10, 0xd3, 0x34, 0x45, 0xe3, 0x36, 0x2f, 0x6f, 0x4a, 0x20, 0xe6, 0x5d, 0x58, 0xde, 0x25, 0x3b,
	0xd0, 0x18, 0x84, 0x74, 0x09, 0xca, 0x6c, 0x2f, 0x95, 0xa5, 0x29, 0x2b, 0x33, 0x3d, 0x58, 0xd9,
	0xa5, 0x5b, 0xe2, 0x18, 0x4d, 0x7e, 0x0c, 0x16, 0xad, 0x41, 0xe8, 0xb5, 0xe2, 0xbb, 0x89, 0xdc,
	0xc5, 0x3c, 0x86, 0xd1, 0x77, 0x14, 0xf3, 0x73, 0x30, 0x7b, 0x07, 0x85, 0x6f, 0x20, 0xbb, 0xf3,
	0x58, 0xe4, 0xde, 0x7f, 0xcc, 0xc1, 0x1c, 0x95, 0x3c, 0xb4, 0x0b, 0xde, 0x03, 0x66, 0x66, 0x7d,
	0xc3, 0x21, 0x25, 0xc6, 0x2d, 0xcc, 0x47, 0xbd, 0xe3, 0x96, 0xd5, 0xa6, 0x2c, 0x42, 0xe5, 0xda,
	0x2c, 0xd5, 0xb1, 0x77, 0xac, 0xde, 0xf1, 0x36, 0xad, 0xc0, 0xac, 0x25, 0x7e, 0x08, 0xe1, 0x50,
	0x88, 0x09, 0x07, 0x55, 0xd5, 0x91, 0xb7, 0xa0, 0x48, 0xd5, 0x31, 0xff, 0x55, 0x0e, 0xe6, 0xa8,
	0x48, 0x50, 0x87, 0x99, 0x8c, 0x6f, 0x83, 0x0d, 0x3e, 0xcf, 0xb8, 0x2f, 0x69, 0xd8, 0x85, 0x4c,
	0xc3, 0xe6, 0x7c, 0x5c, 0x8c, 0xf8, 0xd8, 0x6c, 0xc1, 0xb4, 0x58, 0x0e, 0xb6, 0x47, 0x3c, 0x0d,
	0x15, 0x97, 0x16, 0xb1, 0xbd, 0x61, 0x82, 0x52, 0x25, 0x1b, 0x2c, 0xaf, 0xcc, 0xb8, 0x1b, 0xfc,
	0xbf, 0x3c, 0x94, 0xe9, 0x97, 0x29, 0xf3, 0x6b, 0xc8, 0xf3, 0x1b, 0xba, 0x38, 0x85, 0xb1, 0x16,
	0xa7, 0x18, 0x5b, 0x9c, 0x06, 0x14, 0xb1, 0xd1, 0xa5, 0xca, 0x74, 0x5c, 0xa2, 0x2d, 0x5b, 0x39,
	0x59, 0x43, 0x55, 0x55, 0xbe, 0xca, 0xd9, 0x55, 0xbe, 0xea, 0x18, 0x2a, 0x9f, 0xf9, 0xfb, 0x45,
	0x98, 0x6b, 0xa2, 0x8e, 0x13, 0x84, 0xc8, 0xdf, 0x3f, 0xf2, 0xfa, 0xa3, 0x69, 0xfa, 0x1a, 0x54,
	0xb8, 0xba, 0x42, 0x17, 0x66, 0x92, 0x62, 0x8c, 0xe9, 0x28, 0x4d, 0x5e, 0x1b, 0x69, 0x5a, 0x95,
	0xb8, 0xa6, 0xa5, 0xe3, 0xbe, 0x9a, 0x09, 0xf7, 0x57, 0xa1, 0xfe, 0x00, 0x1d, 0x04, 0x4e, 0x48,
	0x45, 0xe2, 0x84, 0xd4, 0x2e, 0xb0, 0x0a, 0x2c, 0x14, 0x15, 0xb9, 0x39, 0x99, 0x28, 0x37, 0x85,
	0xa6, 0x37, 0x15, 0xd7, 0xf4, 0x2e, 0x43, 0x75, 0xe0, 0xbb, 0xad, 0xc0, 0x1d, 0x74, 0x88, 0xe6,
	0xcd, 0xab, 0x2b, 0x03, 0xdf, 0xdd, 0x77, 0x07, 0x1d, 0x3c, 0x78, 0xae, 0x0a, 0x3a, 0xbd, 0x43,
	0x8f, 0xe9, 0xd8, 0x6c, 0xf0, 0xbb, 0xb4, 0x86, 0x10, 0x29, 0xd7, 0x0b, 0xf1, 0x0f, 0xe3, 0x65,
	0x58, 0xea, 0x7a, 0x3d, 0x74, 0xda, 0x0a, 0x7d, 0xab, 0x17, 0x58, 0x6d, 0x6c, 0xb2, 0xb6, 0x7c,
	0x7f, 0xe0, 0x22, 0xa2, 0x69, 0xf3, 0x5e, 0x16, 0x08, 0xd0, 0xbd, 0x08, 0xa6, 0x89, 0x41, 0x8c,
	0xe7, 0xf1, 0x96, 0xed, 0x9f, 0x20, 0xd6, 0xe5, 0x9c, 0x6c, 0xac, 0xef, 0x93, 0x0a, 0xd2, 0x23,
	0x04, 0xe2, 0x6f, 0xa3, 0x0f, 0xeb, 0xdc, 0xa3, 0xd1, 0x0a, 0x90, 0x7f, 0xe2, 0xb4, 0x51, 0x2b,
	0x40, 0x2e, 0x6a, 0x87, 0xad, 0x20, 0xf4, 0xad, 0x10, 0x75, 0x4e, 0x1b, 0xf3, 0xa4, 0x9d, 0x6b,
	0xac, 0x1d, 0x06, 0xbd, 0x4f, 0x81, 0xf7, 0x09, 0xec, 0x3e, 0x03, 0xdd, 0x0b, 0x51, 0xb7, 0xb9,
	0x1a, 0x0c, 0x03, 0x31, 0x5f, 0x82, 0x79, 0x95, 0x96, 0x18, 0xcb, 0xaf, 0x01, 0xf1, 0x0d, 0x09,
	0x01, 0xcc, 0x7a, 0xf3, 0xfa, 0x4d, 0x52, 0x6e, 0xfe, 0x79, 0x09, 0x66, 0xa9, 0xc4, 0x92, 0x49,
	0x70, 0x0b, 0xe6, 0xa2, 0x4d, 0xc0, 0x3b, 0x41, 0x7e, 0x10, 0x7a, 0xed, 0x63, 0xc2, 0x78, 0xd5,
	0xa6, 0x21, 0xaa, 0xde, 0xe6, 0x35, 0xbf, 0xcb, 0x34, 0xfb, 0x34, 0x4c, 0x93, 0x5d, 0x93, 0xca,
	0x84, 0xd6, 0xe1, 0x61, 0x97, 0x58, 0x9c, 0xd5, 0xe6, 0x24, 0x2e, 0xa6, 0x9b, 0xd6, 0xab, 0x87,
	0x5d, 0xe3, 0x93, 0x30, 0xa9, 0x38, 0xc9, 0x08, 0x81, 0x4f, 0x61, 0x61, 0x21, 0x97, 0x6e, 0xde,
	0x79, 0xed, 0xad, 0xb7, 0xbc, 0x10, 0xed, 0x7a, 0x36, 0x6a, 0xd6, 0x3b, 0x47, 0x3d, 0xfe, 0xc3,
	0xb8, 0x0e, 0x65, 0xea, 0x41, 0x23, 0x44, 0x3f, 0x75, 0x73, 0x76, 0x93, 0xfe, 0xdc, 0xbc, 0xe7,
	0x9f, 0xbe, 0xdd, 0x23, 0xf0, 0xa5, 0x10, 0xff, 0x19, 0x63, 0x12, 0xe3, 0x51, 0x99, 0x64, 0x6e,
	0x6c, 0x26, 0x99, 0x3f, 0x27, 0x26, 0x59, 0x38, 0x57, 0x26, 0xb9, 0x05, 0x86, 0x4c, 0xeb, 0x19,
	0x59, 0xe4, 0x87, 0x79, 0x80, 0x5d, 0xcf, 0xc5, 0x2d, 0x39, 0x5e, 0x6f, 0xb4, 0x67, 0xaa, 0x9a,
	0xe0, 0x99, 0x4a, 0xdf, 0x0a, 0x9f, 0x86, 0xba, 0x8d, 0x82, 0xb6, 0xef, 0xf4, 0x71, 0xeb, 0x8a,
	0xe2, 0x21, 0x57, 0xe0, 0x8d, 0x2c, 0x38, 0xf2, 0xfc, 0xb0, 0x85, 0x0b, 0x55, 0xfd, 0x83, 0x94,
	0xdf, 0x46, 0x41, 0x1b, 0x93, 0x32, 0xae, 0x6e, 0x1d, 0x85, 0x5d, 0x4d, 0x6d, 0xc5, 0xc5, 0xaf,
	0x85, 0x5d, 0x57, 0xdb, 0xeb, 0xca, 0x67, 0xdf, 0xeb, 0x2a, 0xe3, 0xec, 0x75, 0xff, 0x21, 0xc7,
	0x9d, 0x86, 0x11, 0x26, 0xf5, 0xfd, 0xee, 0x03, 0xc7, 0x8d, 0xf9, 0x45, 0xb8, 0xc4, 0x0c, 0x3a,
	0x6a, 0x0a, 0xec, 0x62, 0xb2, 0xf1, 0xfc, 0x53, 0x3e, 0x52, 0xd5, 0xc6, 0xc8, 0x25, 0xdb, 0x18,
	0x57, 0xa1, 0xde, 0x66, 0xdf, 0xe9, 0x7e, 0x4a, 0xe0, 0x15, 0x7b, 0xb6, 0xf9, 0x26, 0xcc, 0x45,
	0xa8, 0x88, 0x54, 0xb4, 0x97, 0xa0, 0xde, 0x8e, 0x8a, 0x99, 0x9a, 0x36, 0x4f, 0xd5, 0x34, 0x4c,
	0x93, 0x12, 0xfa, 0x64, 0x40, 0xf3, 0xbb, 0x39, 0x58, 0xa2, 0x63, 0x8f, 0x23, 0x38, 0xbb, 0xf6,
	0xb9, 0x9e, 0x80, 0x70, 0x15, 0xd5, 0xab, 0x71, 0x54, 0xcb, 0x48, 0x5e, 0x89, 0x21, 0x59, 0x42,
	0x6f, 0x17, 0x2e, 0x2b, 0xe8, 0x0d, 0xe2, 0x43, 0x7d, 0x06, 0x26, 0xa3, 0x59, 0xe9, 0x48, 0x9e,
	0x88, 0xaa, 0xf6, 0x88, 0x25, 0x1d, 0x2d, 0x06, 0x35, 0x1e, 0x0a, 0x4d, 0x10, 0xeb, 0x10, 0xe0,
	0xee, 0x9a, 0xa8, 0xeb, 0x9d, 0x3c, 0xa1, 0xee, 0xde, 0xcb, 0x43, 0xfd, 0x95, 0xd0, 0xeb, 0xdf,
	0xa7, 0xf6, 0x6d, 0xba, 0x4e, 0x4c, 0x84, 0x7f, 0x3e, 0xd1, 0x87, 0xab, 0x7b, 0xb7, 0x12, 0xd9,
	0xa0, 0x9c, 0x8d, 0x0d, 0x2a, 0x19, 0xd8, 0xa0, 0x9a, 0x28, 0x22, 0x56, 0x01, 0xc4, 0x86, 0x18,
	0x34, 0x6a, 0xeb, 0x05, 0xbc, 0xc6, 0x7c, 0x2f, 0x24, 0x4e, 0x5f, 0xd7, 0x09, 0xb8, 0x45, 0x8e,
	0xf7, 0xba, 0x12, 0xef, 0x06, 0x97, 0x0b, 0x73, 0xbc, 0xed, 0x09, 0xa0, 0x19, 0x19, 0x08, 0x97,
	0x53, 0x20, 0xd5, 0x50, 0x36, 0x46, 0x1b, 0xca, 0xdf, 0x2f, 0x50, 0x1c, 0xb3, 0x15, 0xfd, 0x2d,
	0xc6, 0x71, 0x03, 0x8a, 0x83, 0x9e, 0x13, 0x2a, 0x4a, 0x09, 0x29, 0x79, 0x72, 0xd8, 0xd7, 0x04,
	0x15, 0x24, 0x0b, 0x2a, 0x63, 0x1b, 0x96, 0x39, 0xe1, 0x07, 0xde, 0xc0, 0x6f, 0xa3, 0x96, 0xfc,
	0x55, 0x5d, 0xfa, 0x6a, 0x89, 0xc1, 0xed, 0x13, 0xb0, 0xdd, 0x48, 0xd6, 0xfd, 0x51, 0x09, 0xea,
	0x58, 0x78, 0x0d, 0x67, 0x8d, 0xab, 0x50, 0x24, 0xba, 0x03, 0x57, 0xee, 0x08, 0x1d, 0x48, 0x1c,
	0xd5, 0x24, 0xd5, 0x43, 0x56, 0x77, 0x15, 0x2a, 0xc8, 0xa6, 0xba, 0x95, 0xbc, 0x55, 0x94, 0x91,
	0xbd, 0x3b, 0xd6, 0xe2, 0x97, 0xb2, 0x2d, 0x7e, 0x39, 0xc3, 0xe2, 0x57, 0x1e, 0x9d, 0xc1, 0xa6,
	0x92, 0x97, 0xf8, 0x1a, 0x4c, 0xf8, 0x28, 0xb4, 0x1c, 0x37, 0x81, 0x12, 0xea, 0xb4, 0x86, 0x02,
	0x72, 0xab, 0xb9, 0x3e, 0xe4, 0x6c, 0x68, 0x22, 0xcb, 0xd9, 0xd0, 0x35, 0x98, 0x70, 0x82, 0x96,
	0x75, 0x62, 0x39, 0xae, 0x75, 0xe0, 0x22, 0xa2, 0x2a, 0x73, 0x37, 0x78, 0xdd, 0x09, 0xb6, 0x79,
	0x85, 0xf1, 0x96, 0x72, 0x20, 0x41, 0x17, 0x8f, 0x59, 0x6a, 0x1b, 0xfa, 0x81, 0x04, 0xad, 0x95,
	0x88, 0x43, 0x3e, 0x9b, 0x60, 0xe4, 0x62, 0x40, 0x31, 0xb4, 0x3a, 0x41, 0x03, 0x08, 0xa2, 0xc8,
	0xdf, 0xc6, 0x65, 0x28, 0x05, 0xa4, 0x70, 0x91, 0x48, 0x8d, 0x1a, 0x6d, 0xf7, 0x9e, 0xd5, 0x69,
	0xd2, 0xf2, 0xb1, 0x65, 0x8b, 0xf1, 0x1c, 0x54, 0x18, 0xfd, 0x36, 0x96, 0xe4, 0xc3, 0x13, 0x3c,
	0xb6, 0x7d, 0xbc, 0xba, 0x4c, 0xe8, 0x34, 0x39, 0x98, 0xf9, 0xa3, 0x1c, 0xac, 0x0c, 0x99, 0x89,
	0xb1, 0x09, 0x33, 0xdc, 0x93, 0xd9, 0xf2, 0x7a, 0xad, 0x23, 0xab, 0x47, 0xcf, 0x24, 0xf9, 0x32,
	0x4d, 0xf1, 0xda, 0xb7, 0x7b, 0xaf, 0x59, 0x3d, 0xdb, 0xb8, 0x01, 0xd3, 0x02, 0xbe, 0xef, 0xb4,
	0x8f, 0x11, 0xf5, 0x63, 0xc6, 0xc0, 0xef, 0x92, 0x3a, 0x8d, 0xc9, 0x4b, 0xc9, 0x4c, 0x2e, 0xfb,
	0x58, 0xcb, 0x12, 0x88, 0x28, 0x35, 0x7f, 0x56, 0xa1, 0xcc, 0x39, 0x5c, 0xa6, 0xa6, 0x32, 0x27,
	0xc7, 0xca, 0xef, 0x0a, 0x73, 0x6a, 0x52, 0xb3, 0x91, 0x22, 0x35, 0xcf, 0x44, 0xb3, 0x4f, 0x9e,
	0x55, 0xcf, 0x57, 0xd4, 0x5c, 0x85, 0x29, 0x45, 0x9d, 0x0a, 0x1a, 0xb3, 0x44, 0x4d, 0x9a, 0x94,
	0x35, 0xa9, 0xc0, 0xb8, 0x01, 0x55, 0x26, 0x15, 0x38, 0x63, 0xce, 0x46, 0xac, 0xc6, 0xc5, 0x80,
	0x00, 0x31, 0x9e, 0x83, 0x59, 0x6d, 0x03, 0x72, 0x6c, 0x62, 0x92, 0x72, 0xbc, 0x4f, 0x2b, 0xfb,
	0xce, 0x9e, 0xad, 0xd9, 0x38, 0x0b, 0x67, 0xb7, 0x71, 0xe6, 0xc7, 0x39, 0xc2, 0x7d, 0x39, 0x3a,
	0x5c, 0x20, 0xce, 0xc8, 0x25, 0xb2, 0x52, 0x17, 0x37, 0x95, 0xb0, 0x10, 0xc6, 0x2e, 0xf7, 0x4e,
	0xfb, 0x48, 0x9c, 0x38, 0xe0, 0x1f, 0xc6, 0x26, 0xd4, 0xbb, 0x28, 0xb4, 0x5a, 0x87, 0x0e, 0x72,
	0xed, 0xa0, 0xb1, 0x4c, 0xf0, 0x32, 0xb9, 0xd9, 0xee, 0x6e, 0xbe, 0x89, 0x42, 0xeb, 0x55, 0x5c,
	0xda, 0x84, 0x2e, 0xff, 0x13, 0x13, 0x53, 0xf5, 0xc0, 0xb7, 0x7a, 0x36, 0x46, 0xc6, 0x8a, 0x84,
	0x8c, 0x0a, 0x29, 0xdd, 0xb3, 0xcd, 0x1d, 0x98, 0xd1, 0x45, 0xd7, 0xb8, 0x7e, 0x5a, 0xf3, 0x27,
	0x39, 0x98, 0x52, 0xad, 0x8e, 0xb3, 0x05, 0x66, 0xf0, 0x1e, 0x0a, 0xa3, 0xb8, 0xbb, 0x98, 0xc6,
	0xdd, 0x19, 0x2c, 0xdb, 0x2c, 0x02, 0xc0, 0xfc, 0x3c, 0x18, 0x77, 0x50, 0xc8, 0x88, 0xec, 0xac,
	0x07, 0x15, 0xc5, 0xf4, 0x83, 0x8a, 0x2f, 0xc0, 0xfc, 0x1d, 0xc4, 0x0d, 0x47, 0x07, 0x3d, 0x86,
	0x1e, 0x8e, 0x61, 0x5e, 0xe2, 0x93, 0xc8, 0x72, 0x7c, 0x4a, 0xeb, 0x21, 0xd1, 0x69, 0xaf, 0xf0,
	0x5e, 0x7e, 0x24, 0xef, 0xf1, 0xce, 0xb8, 0x05, 0x35, 0x7e, 0x67, 0x8c, 0xc4, 0x13, 0x3a, 0xe3,
	0x1b, 0x87, 0x00, 0x31, 0x7d, 0x58, 0x24, 0xc4, 0x26, 0x21, 0x6f, 0xac, 0xee, 0x6e, 0x02, 0x97,
	0xc0, 0x0e, 0xe2, 0x1d, 0x1a, 0x92, 0xe9, 0xcc, 0xad, 0x79, 0x09, 0xca, 0xfc, 0x75, 0x1e, 0xe6,
	0xa9, 0x51, 0xca, 0x27, 0x3f, 0xb6, 0xd1, 0x6c, 0xc8, 0x67, 0x46, 0x4c, 0x5e, 0x1b, 0x6c, 0x1f,
	0xac, 0xd3, 0x32, 0xb6, 0x03, 0xca, 0xfb, 0x35, 0xd9, 0xd9, 0xe5, 0x9d, 0x7a, 0x55, 0x11, 0xc5,
	0x25, 0x5a, 0x1d, 0x09, 0xe1, 0x0d, 0x4d, 0x08, 0x93, 0xcd, 0x5c, 0x15, 0xbf, 0x9a, 0xf5, 0x5e,
	0x19, 0x65, 0xbd, 0x57, 0x87, 0x5a, 0xef, 0x35, 0xd5, 0x7a, 0xd7, 0x14, 0x2a, 0x18, 0xad, 0x50,
	0x2d, 0x42, 0x21, 0x38, 0x1e, 0x28, 0x7e, 0x7c, 0x5c, 0x60, 0xfe, 0xff, 0x02, 0x47, 0x39, 0x27,
	0x81, 0xb3, 0xa0, 0x3c, 0xd2, 0x24, 0x18, 0x7a, 0x13, 0xce, 0xc0, 0x70, 0x19, 0xb1, 0xb2, 0xa8,
	0x27, 0x82, 0xda, 0x57, 0xeb, 0x09, 0x96, 0xde, 0x30, 0x2c, 0x55, 0x86, 0x62, 0xa9, 0xaa, 0x61,
	0x49, 0x5d, 0xe4, 0xda, 0xf0, 0x45, 0x86, 0x51, 0x8b, 0x5c, 0x8f, 0x2f, 0xb2, 0xbe, 0xcb, 0x4c,
	0x3c, 0xca, 0x2e, 0x33, 0xc5, 0x3c, 0xf2, 0xa9, 0xbb, 0xcc, 0x45, 0x69, 0x97, 0x99, 0x26, 0xf1,
	0x03, 0x62, 0x7f, 0xb1, 0x61, 0x81, 0x39, 0x9c, 0x34, 0x2f, 0x59, 0xf6, 0x65, 0xdc, 0x80, 0x5a,
	0xdf, 0xf2, 0x11, 0x0d, 0x12, 0x90, 0x8f, 0xf6, 0xab, 0xb4, 0x98, 0xb8, 0xc9, 0x16, 0x14, 0xf6,
	0x14, 0x02, 0xf5, 0x16, 0x54, 0xe8, 0xd6, 0xcb, 0x9d, 0x64, 0xcb, 0x94, 0x16, 0x93, 0x98, 0xb9,
	0xc9, 0x41, 0xcd, 0xef, 0xe6, 0xe0, 0xa2, 0xea, 0x83, 0xba, 0x67, 0x75, 0x44, 0x9b, 0x33, 0x50,
	0xc0, 0x4a, 0x4b, 0x8e, 0x28, 0x2d, 0xf8, 0x4f, 0x3c, 0x6a, 0xcb, 0x66, 0xee, 0x9e, 0x5a, 0x93,
	0xfc, 0x6d, 0x34, 0x78, 0x24, 0x50, 0xd0, 0x28, 0x90, 0x62, 0xfe, 0xd3, 0xb8, 0x0c, 0x75, 0x1f,
	0xf5, 0x5d, 0xab, 0x8d, 0x5a, 0x96, 0xeb, 0x12, 0x21, 0x5e, 0x6b, 0x02, 0x2b, 0xda, 0x76, 0xc9,
	0x16, 0x45, 0x61, 0x49, 0x7d, 0x49, 0xd2, 0xca, 0x58, 0x3c, 0xd7, 0xb6, 0xeb, 0x9a, 0x5f, 0x84,
	0x45, 0x7d, 0xca, 0x4c, 0x0c, 0x8e, 0x27, 0xbc, 0x8d, 0x0d, 0x28, 0x23, 0xdf, 0xf7, 0x7c, 0x3a,
	0x4e, 0xac, 0x70, 0xb6, 0xbb, 0x9b, 0xaf, 0xe0, 0x92, 0x26, 0xab, 0x30, 0x11, 0x18, 0xdb, 0xb6,
	0xad, 0xe3, 0x36, 0x11, 0x0f, 0x44, 0x73, 0xcd, 0x4b, 0xea, 0x6c, 0xcc, 0x79, 0x56, 0x48, 0x73,
	0x9e, 0x99, 0x1d, 0x98, 0x53, 0xba, 0x49, 0x98, 0x4f, 0x6e, 0x9c, 0xf9, 0xe4, 0xd3, 0xe6, 0xf3,
	0x0c, 0x2c, 0x50, 0x9f, 0xdf, 0xc8, 0x29, 0x99, 0x5f, 0xce, 0xc1, 0xcc, 0x1d, 0x14, 0x3e, 0x52,
	0xa4, 0x56, 0xfa, 0x36, 0x6d, 0x5c, 0x87, 0x52, 0xd7, 0x79, 0x88, 0x6c, 0xb2, 0xca, 0x78, 0x1f,
	0x22, 0xc7, 0x0a, 0x6f, 0xe2, 0x22, 0x7e, 0x58, 0x45, 0x01, 0xcc, 0xff, 0x25, 0x68, 0x92, 0x8e,
	0x85, 0x6a, 0xf4, 0xe9, 0x6b, 0xf1, 0x0c, 0x54, 0x58, 0xb8, 0x23, 0x61, 0xa6, 0xb8, 0x31, 0xd0,
	0xe4, 0xf5, 0x64, 0x89, 0x48, 0x60, 0x09, 0x8f, 0xbb, 0x94, 0x15, 0x30, 0x16, 0x25, 0xc9, 0x82,
	0x2e, 0x23, 0x0b, 0xa3, 0xa8, 0x34, 0x7a, 0x2b, 0xd1, 0xc2, 0x30, 0xff, 0x6b, 0x0e, 0xe6, 0x58,
	0x84, 0x89, 0x12, 0x59, 0x72, 0x19, 0xaa, 0x24, 0xfa, 0x4c, 0x77, 0xa6, 0x56, 0x48, 0xe9, 0x9e,
	0x6d, 0xdc, 0x1a, 0x1e, 0x64, 0x92, 0x1c, 0x5e, 0x62, 0xbc, 0x0c, 0x4b, 0xca, 0x21, 0xdb, 0xc0,
	0x3d, 0x74, 0x5c, 0xb7, 0x8b, 0x58, 0xec, 0x00, 0xe7, 0xa2, 0x05, 0xe9, 0xc8, 0x2d, 0x02, 0x31,
	0x9f, 0x87, 0x29, 0x32, 0xc8, 0xbd, 0xdb, 0x59, 0x87, 0x69, 0x6e, 0xc2, 0x34, 0xfb, 0x44, 0xac,
	0xc4, 0x0a, 0xd4, 0xf8, 0x37, 0x7c, 0x3d, 0xaa, 0x0c, 0x1c, 0x33, 0xd2, 0x46, 0xac, 0xdf, 0xe0,
	0x55, 0xcf, 0x1f, 0x0f, 0x39, 0x97, 0xa1, 0x1e, 0x85, 0x4d, 0x51, 0xb6, 0x2d, 0x34, 0x41, 0x44,
	0x4c, 0x05, 0xe6, 0xb7, 0x72, 0x60, 0xd0, 0x50, 0xa1, 0xf1, 0x1a, 0x8e, 0x11, 0x42, 0x3e, 0x95,
	0x10, 0xd2, 0xa3, 0x80, 0x0a, 0x23, 0xa3, 0x80, 0x5e, 0x81, 0x39, 0x69, 0x74, 0x43, 0x68, 0x78,
	0x78, 0xf4, 0xd2, 0xdf, 0xe6, 0x60, 0xf2, 0xae, 0x6c, 0xd3, 0xa5, 0x1b, 0x2e, 0x64, 0x0f, 0xcc,
	0xc7, 0xc2, 0x3e, 0xd2, 0x0d, 0x8e, 0x34, 0x3a, 0x4f, 0xb1, 0xa4, 0x55, 0x7b, 0xb0, 0x34, 0x66,
	0x48, 0xef, 0x19, 0x4f, 0xda, 0xcc, 0x4b, 0x22, 0xce, 0x52, 0xc6, 0x02, 0x43, 0xa7, 0x79, 0x1f,
	0x16, 0x95, 0x72, 0x35, 0x08, 0x54, 0x31, 0x95, 0xf5, 0x20, 0x50, 0xa5, 0xb9, 0x29, 0xc5, 0x6e,
	0x0e, 0xcc, 0x7b, 0xb0, 0xc0, 0xce, 0xe8, 0xb4, 0x1d, 0x3d, 0xfd, 0x74, 0x5f, 0xd9, 0xc1, 0xf3,
	0x89, 0x3b, 0xf8, 0x57, 0x8b, 0x30, 0xaf, 0x4c, 0x46, 0x6a, 0x95, 0x28, 0x71, 0xb9, 0x54, 0x37,
	0x7d, 0xfc, 0x44, 0x90, 0xbb, 0xcd, 0x0b, 0x31, 0xb7, 0x79, 0x7a, 0x70, 0xd0, 0xdf, 0x3f, 0xef,
	0xee, 0x68, 0xdf, 0x7c, 0x26, 0x2f, 0xbf, 0xee, 0x97, 0x59, 0x4e, 0xf3, 0xcb, 0xe8, 0x3a, 0xe3,
	0xca, 0x58, 0x3a, 0xa3, 0xec, 0x69, 0x58, 0x4d, 0xf0, 0x34, 0xe8, 0x4a, 0xe5, 0xa5, 0x11, 0xae,
	0x0b, 0xf3, 0xff, 0x16, 0x38, 0x45, 0x68, 0x36, 0xd7, 0x59, 0x28, 0x42, 0x3d, 0x93, 0x2d, 0x24,
	0x9f, 0xc9, 0xfe, 0xd6, 0x10, 0x87, 0x6a, 0x64, 0x4d, 0x8d, 0x36, 0xb2, 0x9e, 0x3c, 0x35, 0x99,
	0x5f, 0x2f, 0xc3, 0xda, 0x6d, 0xd4, 0xf7, 0x11, 0x36, 0xa3, 0xed, 0xc4, 0x85, 0x4c, 0x74, 0xd9,
	0x2d, 0x0d, 0x73, 0xd9, 0x65, 0x3a, 0x5a, 0xd7, 0xc3, 0x77, 0xf3, 0x69, 0xe1, 0xbb, 0x99, 0xfd,
	0x4d, 0x1f, 0xf4, 0x7a, 0x27, 0xb9, 0x89, 0xc7, 0xf4, 0xf5, 0x7e, 0x00, 0xf2, 0x84, 0xf3, 0xe9,
	0x5c, 0x8c, 0x4f, 0x5f, 0x00, 0x43, 0x1c, 0x61, 0x44, 0xee, 0xe7, 0x79, 0xa9, 0xa1, 0x59, 0x5e,
	0x1f, 0x39, 0xa1, 0x93, 0xce, 0x49, 0x16, 0x86, 0x9c, 0x93, 0x3c, 0x0f, 0xa2, 0x91, 0x96, 0x8f,
	0x02, 0xe4, 0x9f, 0x20, 0xbb, 0xb1, 0x28, 0x7d, 0x20, 0x9a, 0x6b, 0xb2, 0x5a, 0x8d, 0xaf, 0x1a,
	0xa3, 0xf9, 0x8a, 0x6f, 0x34, 0x17, 0x63, 0x1b, 0x0d, 0x73, 0x6b, 0x18, 0xba, 0x5b, 0xe3, 0x2e,
	0x5c, 0xda, 0xf5, 0x7a, 0x3d, 0xd4, 0x0e, 0xb5, 0x4d, 0x3b, 0x18, 0xc2, 0x13, 0xb9, 0x21, 0x3c,
	0x61, 0xde, 0xe7, 0x31, 0x33, 0x77, 0xf7, 0xcf, 0x75, 0x47, 0x76, 0xf4, 0x38, 0x8c, 0x78, 0xfb,
	0xda, 0x29, 0x47, 0x2e, 0xe5, 0x94, 0x63, 0x64, 0x50, 0x04, 0x82, 0xf5, 0xf4, 0x90, 0x0f, 0x11,
	0x7f, 0xc5, 0xcc, 0x73, 0xda, 0x4f, 0x49, 0xc4, 0x7c, 0xd2, 0xc2, 0x2c, 0x66, 0xdf, 0x03, 0x30,
	0x95, 0x6e, 0xd4, 0x13, 0xe8, 0xe1, 0x8e, 0x89, 0xd1, 0x08, 0x4b, 0x17, 0x27, 0xe6, 0x15, 0xd8,
	0xb8, 0x83, 0xc2, 0xa4, 0x5e, 0x23, 0xdf, 0xaf, 0xf9, 0xc7, 0x39, 0x58, 0xbc, 0x83, 0x42, 0xd9,
	0x32, 0x38, 0xf7, 0x08, 0xf9, 0xc8, 0xde, 0x2c, 0x8c, 0xb0, 0x37, 0x15, 0x63, 0xa1, 0x98, 0x64,
	0x2c, 0x5c, 0x13, 0x02, 0xa8, 0x94, 0x6c, 0x5f, 0x72, 0x23, 0xf0, 0xdf, 0xe6, 0xe0, 0x92, 0x3a,
	0xb1, 0xd7, 0x9c, 0x20, 0x94, 0x30, 0x9e, 0x65, 0x7a, 0x8b, 0x50, 0xb0, 0x5c, 0x97, 0x60, 0x9e,
	0x9b, 0x71, 0xb8, 0x80, 0xad, 0x56, 0x21, 0x76, 0x7f, 0x66, 0xf8, 0xe0, 0xcd, 0x97, 0x08, 0xae,
	0x77, 0x2c, 0x17, 0xdb, 0x21, 0x4a, 0x9c, 0xdf, 0x25, 0x28, 0xb3, 0x9b, 0x30, 0x32, 0x99, 0xb1,
	0x32, 0xb3, 0x03, 0x2b, 0x77, 0x50, 0xf8, 0xa6, 0x16, 0xdc, 0xf8, 0x18, 0xae, 0x32, 0x7c, 0x98,
	0x74, 0x74, 0x77, 0x70, 0xe0, 0x3a, 0x6d, 0x09, 0x73, 0x23, 0x75, 0x20, 0xf3, 0x21, 0x3c, 0x4d,
	0x89, 0x5c, 0x26, 0x24, 0x11, 0x13, 0x19, 0x5a, 0x21, 0x4a, 0xb7, 0xb7, 0x5e, 0x86, 0xa9, 0x28,
	0xd6, 0x12, 0x83, 0x12, 0xc4, 0x92, 0xa5, 0x65, 0xc5, 0x64, 0x6d, 0x11, 0xbf, 0x5b, 0x15, 0xc8,
	0xcd, 0x9a, 0xc7, 0x5c, 0x60, 0x10, 0xb3, 0xee, 0x2e, 0xbd, 0x76, 0xa3, 0xba, 0x29, 0x46, 0x5a,
	0xa0, 0x11, 0x51, 0xe5, 0x87, 0x13, 0xd5, 0x17, 0xe0, 0xd2, 0xfe, 0xa0, 0xdb, 0xb5, 0x7c, 0xe7,
	0x9f, 0xa2, 0x24, 0x96, 0xc1, 0x9b, 0x2e, 0xf1, 0x01, 0xf8, 0x5e, 0x57, 0xbd, 0x68, 0x48, 0x70,
	0xe3, 0x7b, 0x5d, 0x63, 0x15, 0x2a, 0x04, 0x24, 0xf4, 0x54, 0xf3, 0x12, 0x17, 0xde, 0xf3, 0xcc,
	0x4f, 0xc1, 0x6a, 0x4a, 0x0f, 0x8c, 0x52, 0x3e, 0x04, 0xe5, 0x10, 0xef, 0x4c, 0xdc, 0x7a, 0xe2,
	0x87, 0x08, 0xe4, 0xa3, 0xd3, 0x7b, 0xb8, 0xaa, 0xc9, 0x20, 0xcc, 0x1d, 0x98, 0x17, 0x8d, 0xdd,
	0x7d, 0x7b, 0xff, 0x4c, 0x6d, 0x7c, 0x06, 0xe6, 0xd4, 0x36, 0xce, 0x6b, 0xa6, 0x3f, 0xc9, 0xc1,
	0x84, 0xdc, 0xa3, 0xb1, 0x0c, 0x25, 0xd7, 0x3a, 0x40, 0xda, 0x05, 0x6d, 0x52, 0x94, 0xe8, 0xe3,
	0x7b, 0x01, 0x2a, 0x6d, 0xcf, 0x1d, 0x74, 0x7b, 0xdc, 0x87, 0x38, 0xa7, 0x4c, 0x63, 0xd7, 0x73,
	0x9b, 0xde, 0x03, 0xbe, 0xd4, 0x0c, 0xd2, 0xb8, 0x01, 0x45, 0xdf, 0x7b, 0xc0, 0xbd, 0x63, 0x43,
	0xbe, 0x20, 0x60, 0xc6, 0xb3, 0x50, 0xb4, 0xad, 0xd0, 0x6a, 0x94, 0x14, 0x0f, 0x20, 0x05, 0xdf,
	0x0b, 0x51, 0x97, 0x03, 0x63, 0x20, 0xf3, 0x9f, 0xc1, 0xa4, 0xd2, 0xd2, 0xd0, 0x19, 0x35, 0xa0,
	0x18, 0xf4, 0x51, 0x5b, 0xb5, 0x1a, 0x70, 0xc9, 0x10, 0x3b, 0xf2, 0x12, 0x94, 0x9d, 0x9e, 0x8d,
	0xd8, 0xa5, 0x3a, 0x21, 0x25, 0x68, 0x99, 0xf9, 0xcd, 0x1c, 0xd4, 0xa5, 0xa1, 0x9d, 0xb1, 0xf7,
	0x65, 0x28, 0x45, 0x97, 0xfb, 0x78, 0x17, 0xb4, 0x48, 0x8c, 0xac, 0x38, 0x22, 0x30, 0xac, 0xa4,
	0xa9, 0x92, 0xe6, 0xef, 0xe5, 0x60, 0x71, 0xaf, 0xdb, 0x8f, 0x4e, 0x7b, 0x23, 0x7a, 0x7e, 0x8e,
	0x61, 0x98, 0x8a, 0xae, 0x4b, 0x9b, 0x72, 0x22, 0x87, 0xfd, 0xe8, 0xef, 0xdb, 0x56, 0x68, 0x51,
	0x34, 0x1b, 0x9b, 0x30, 0xe9, 0x90, 0xb6, 0x5a, 0x69, 0x1e, 0xe4, 0x09, 0x5a, 0x4f, 0x7e, 0x04,
	0xc6, 0x87, 0xa0, 0xde, 0x46, 0xae, 0xcb, 0xa1, 0x8b, 0x3a, 0x34, 0xe0, 0x5a, 0x06, 0x4b, 0xe2,
	0xed, 0x49, 0xdb, 0x8e, 0xad, 0x5c, 0xa1, 0xaf, 0xd2, 0xe2, 0x3d, 0xdb, 0xfc, 0x28, 0x2c, 0xed,
	0x5a, 0x6e, 0x3b, 0x49, 0x8a, 0xaf, 0x41, 0xe5, 0x80, 0x16, 0xab, 0xda, 0x02, 0x2b, 0x34, 0xff,
	0x5d, 0x1e, 0xe6, 0x19, 0x03, 0xbd, 0xf2, 0x10, 0x37, 0x27, 0xa9, 0x34, 0x88, 0x14, 0x50, 0x1b,
	0x58, 0x5e, 0x30, 0xa0, 0x15, 0xc4, 0xda, 0xcd, 0xb6, 0xdb, 0x2a, 0xac, 0x59, 0x18, 0xc5, 0x9a,
	0xc5, 0x38, 0x6b, 0x1a, 0x26, 0xb6, 0x1d, 0x5c, 0xa7, 0xeb, 0x84, 0xfa, 0x55, 0x65, 0x51, 0x6c,
	0x7c, 0x04, 0x16, 0xd0, 0xc3, 0x36, 0x72, 0x5b, 0x24, 0x4e, 0x3f, 0x74, 0x0e, 0x5c, 0xd4, 0xea,
	0xe2, 0xcd, 0xa1, 0x2c, 0x6d, 0x96, 0x73, 0x04, 0x64, 0x57, 0x40, 0xbc, 0x89, 0xf5, 0x70, 0xb6,
	0x03, 0x54, 0x22, 0x77, 0xf7, 0x7f, 0xca, 0xc1, 0x82, 0x86, 0x17, 0x86, 0xd1, 0x48, 0x2d, 0xaa,
	0x49, 0x1b, 0xed, 0x3a, 0x54, 0x0f, 0x1d, 0x17, 0xc5, 0xec, 0x2c, 0x51, 0xaa, 0x23, 0xb4, 0x90,
	0x82, 0xd0, 0x31, 0x1d, 0xcb, 0x73, 0xe4, 0x22, 0x21, 0x1d, 0xa2, 0xd0, 0xa0, 0xf6, 0xc8, 0xa9,
	0xbd, 0x28, 0x64, 0x03, 0x7f, 0x01, 0x26, 0xd8, 0x00, 0x9c, 0x10, 0x75, 0xa3, 0x34, 0x1a, 0x34,
	0x36, 0x89, 0xd4, 0x90, 0x5b, 0x01, 0x6c, 0x98, 0xf8, 0xef, 0xc0, 0xfc, 0x56, 0x11, 0x20, 0xaa,
	0x7b, 0xdc, 0x93, 0xbf, 0x06, 0x13, 0xb6, 0xf7, 0xa0, 0xe7, 0x7a, 0x96, 0x4d, 0xae, 0x97, 0xa8,
	0xf1, 0x0d, 0xac, 0xe6, 0x1d, 0x9f, 0x9c, 0x0c, 0xb1, 0xcb, 0x2d, 0x3a, 0x57, 0xd4, 0x58, 0xf9,
	0x1e, 0x89, 0xb2, 0x18, 0x04, 0x74, 0x8f, 0x95, 0x6f, 0xb2, 0x95, 0x71, 0x61, 0x2c, 0xec, 0x65,
	0xdc, 0x6b, 0x6c, 0x52, 0xfa, 0x81, 0xea, 0x38, 0xe9, 0x07, 0x9e, 0x81, 0x49, 0x9f, 0x2e, 0x53,
	0xeb, 0xdd, 0x01, 0xf2, 0x4f, 0x95, 0x0b, 0xb1, 0x13, 0xac, 0xea, 0xd3, 0xb8, 0x06, 0x73, 0x4d,
	0xd7, 0xe9, 0x22, 0x8a, 0x32, 0x90, 0xf1, 0x8a, 0x8b, 0x35, 0x6a, 0xa9, 0x67, 0xa0, 0x16, 0x2c,
	0xa7, 0xd8, 0x32, 0x30, 0xc9, 0x33, 0x11, 0x93, 0x53, 0xb4, 0x9e, 0xc9, 0x9e, 0xcb, 0x50, 0x22,
	0x80, 0x24, 0x22, 0x4a, 0x81, 0xa3, 0xe5, 0xe6, 0x45, 0x58, 0x8a, 0x28, 0x4d, 0x51, 0x71, 0xcc,
	0x3f, 0xc8, 0xc1, 0x0c, 0xad, 0xa0, 0xe5, 0x43, 0xe8, 0x87, 0x3a, 0x2a, 0x3a, 0x3e, 0x0a, 0x82,
	0x56, 0xd7, 0x7a, 0x48, 0x68, 0xa8, 0x24, 0x39, 0x2a, 0x48, 0xcd, 0x9b, 0xd6, 0x43, 0xe3, 0x59,
	0x98, 0x12, 0x80, 0xf1, 0x2d, 0x61, 0x92, 0xd7, 0x89, 0xfb, 0xf0, 0x02, 0x98, 0xce, 0xa2, 0x98,
	0x04, 0x4c, 0xe6, 0x13, 0xcd, 0xb4, 0x94, 0x32, 0xd3, 0x4f, 0xc2, 0xd2, 0xf6, 0x20, 0x3c, 0xf2,
	0x88, 0xd2, 0x61, 0xf9, 0x61, 0x2f, 0x3a, 0x4e, 0xb8, 0x02, 0xd0, 0xa7, 0x25, 0x71, 0x67, 0x0c,
	0x2d, 0xdf, 0xb3, 0xcd, 0xd7, 0x61, 0x0e, 0xeb, 0xb1, 0xf4, 0xb7, 0xcc, 0x94, 0x55, 0x06, 0xc3,
	0x19, 0x72, 0x89, 0x9a, 0x22, 0x54, 0xe3, 0x65, 0xb6, 0x08, 0x89, 0xdb, 0x10, 0x80, 0xe6, 0x03,
	0xb8, 0x28, 0xc6, 0x62, 0x8b, 0xc1, 0xb0, 0x16, 0x9f, 0x87, 0x0a, 0x03, 0x64, 0x1b, 0x58, 0x6a,
	0x83, 0x1c, 0x8e, 0x7a, 0x32, 0x6c, 0xc7, 0x47, 0xed, 0x90, 0xb0, 0x9d, 0xe2, 0x28, 0xe2, 0x35,
	0xef, 0xf8, 0xae, 0xf9, 0x59, 0x58, 0xbd, 0x83, 0xc2, 0x58, 0xdf, 0xd1, 0x74, 0x3e, 0x1e, 0x9b,
	0xce, 0x65, 0xe6, 0x36, 0x48, 0x1b, 0xaf, 0x34, 0xad, 0x6d, 0xa8, 0x09, 0xef, 0xc2, 0x10, 0x93,
	0x5d, 0xa8, 0x03, 0x4a, 0x5e, 0x0c, 0x52, 0x64, 0x7e, 0x2f, 0xc7, 0x93, 0x0c, 0x30, 0x17, 0xdb,
	0x1e, 0xde, 0xf0, 0x83, 0x91, 0x67, 0xed, 0x4f, 0xfc, 0xd4, 0xfa, 0x10, 0xd6, 0x14, 0x13, 0x5c,
	0xb8, 0x7e, 0x47, 0x8c, 0x55, 0x73, 0x1e, 0xe7, 0x47, 0x39, 0x8f, 0x6f, 0x83, 0x91, 0x10, 0x20,
	0xb4, 0xa9, 0x84, 0xfe, 0xd0, 0x85, 0x9a, 0xa2, 0x0b, 0x95, 0x18, 0xf6, 0xf3, 0x3f, 0x72, 0x50,
	0xe5, 0x15, 0x63, 0xdf, 0x5e, 0x4e, 0xf4, 0xe4, 0x14, 0x86, 0x79, 0x37, 0x15, 0x1f, 0x43, 0x39,
	0xd1, 0xc7, 0x20, 0x45, 0xd0, 0x55, 0xe2, 0x11, 0x74, 0xe6, 0x87, 0xa1, 0x70, 0xcf, 0xea, 0xa4,
	0x0c, 0x55, 0x68, 0xa1, 0xf9, 0x98, 0x16, 0x6a, 0xfe, 0xa8, 0x0c, 0x8b, 0xaf, 0x3c, 0x0c, 0x91,
	0xdf, 0xb3, 0x5c, 0xc6, 0x2b, 0xdb, 0x47, 0x56, 0xd7, 0x3b, 0x41, 0xe9, 0x8d, 0xa5, 0x66, 0x69,
	0x49, 0xf7, 0xab, 0x3e, 0x0f, 0xb3, 0x88, 0xf5, 0xd2, 0x3a, 0x41, 0xbe, 0x73, 0xe8, 0xb0, 0xe8,
	0x65, 0x4e, 0x2a, 0x33, 0xbc, 0xfa, 0x3e, 0xab, 0x35, 0x5e, 0x87, 0x39, 0xf1, 0xc9, 0x58, 0x07,
	0x65, 0xa2, 0xa7, 0x5d, 0x79, 0xff, 0x7a, 0xf2, 0x37, 0xb8, 0x8d, 0x3d, 0x98, 0x73, 0xad, 0x20,
	0x6c, 0x05, 0xa8, 0x67, 0xd3, 0x59, 0x9f, 0xe2, 0x36, 0x6a, 0x23, 0xdb, 0x98, 0xc1, 0x9f, 0xed,
	0xa3, 0x9e, 0x4d, 0x90, 0x71, 0xba, 0x1d, 0x1a, 0x37, 0xc1, 0x10, 0xc8, 0x08, 0x9d, 0xf6, 0x31,
	0x0a, 0xf9, 0x9d, 0x8c, 0x9a, 0x8e, 0xc0, 0x7b, 0xa4, 0x7a, 0xcf, 0x36, 0xb6, 0x60, 0xd6, 0xb1,
	0x5b, 0x6d, 0xcb, 0xb7, 0xb1, 0x1e, 0x8a, 0x89, 0xab, 0xdb, 0x51, 0x62, 0x87, 0xa7, 0x1c, 0x7b,
	0xd7, 0xf2, 0xed, 0x57, 0x71, 0xe5, 0x5e, 0xb7, 0x63, 0xdc, 0x80, 0x19, 0xfe, 0xc1, 0x81, 0xd5,
	0x3e, 0x26, 0xf0, 0xf2, 0x05, 0x94, 0x49, 0x0a, 0xbf, 0x63, 0xb5, 0x8f, 0x31, 0x38, 0xde, 0xc5,
	0x3c, 0x3f, 0xf4, 0x2d, 0x87, 0x36, 0x3d, 0xa9, 0xb8, 0xdb, 0x59, 0x0d, 0x06, 0xfc, 0x38, 0xd4,
	0x07, 0x7d, 0xac, 0xc9, 0x64, 0xcd, 0xd7, 0x04, 0x1c, 0x7c, 0x9b, 0xe8, 0xe5, 0x87, 0x56, 0xaf,
	0xcf, 0x2f, 0xe0, 0x4e, 0xcb, 0x9a, 0x14, 0xab, 0xc0, 0x0a, 0x92, 0x76, 0x99, 0x77, 0x26, 0xf5,
	0x32, 0x6f, 0x74, 0xf7, 0xb5, 0xdb, 0xa1, 0x01, 0xc8, 0xb5, 0xe8, 0xa2, 0x6b, 0xb7, 0x13, 0x18,
	0x37, 0x61, 0xe1, 0x60, 0x10, 0x38, 0x3d, 0xbc, 0x8d, 0xba, 0x4e, 0x1b, 0xf5, 0x02, 0x44, 0x61,
	0x0d, 0x02, 0x3b, 0xc7, 0x2b, 0xdf, 0xa0, 0x75, 0xf8, 0x1b, 0xf3, 0x57, 0x79, 0xb8, 0x4e, 0xc5,
	0xdb, 0x3f, 0x52, 0x99, 0x88, 0x52, 0x73, 0xdb, 0x1a, 0x7d, 0xdf, 0x2e, 0x71, 0xb5, 0xf2, 0x63,
	0xae, 0x56, 0x21, 0xfb, 0x6a, 0x15, 0xd3, 0x56, 0x4b, 0x43, 0x78, 0x29, 0x1b, 0xc2, 0xcb, 0x19,
	0x11, 0x5e, 0x19, 0x03, 0xe1, 0xd5, 0x74, 0x84, 0xff, 0xac, 0x10, 0x13, 0x5b, 0xaf, 0x59, 0xbe,
	0x75, 0x62, 0x9d, 0x31, 0x02, 0xd9, 0x84, 0x5a, 0x30, 0x38, 0xb0, 0xbd, 0xae, 0xe5, 0xa8, 0x51,
	0x30, 0x51, 0xb1, 0xb1, 0x0b, 0x2b, 0x91, 0x40, 0xb2, 0x7c, 0xdf, 0x41, 0xbe, 0xb8, 0x98, 0xcc,
	0xee, 0xf5, 0x72, 0x65, 0xab, 0x21, 0x44, 0x10, 0x85, 0x63, 0xf7, 0x8d, 0xf7, 0x6c, 0xc3, 0x85,
	0xab, 0x51, 0x23, 0xf4, 0x3c, 0x00, 0xd9, 0xb1, 0xe6, 0x32, 0x49, 0x89, 0x0d, 0xd1, 0x09, 0x6f,
	0x47, 0xed, 0x8d, 0x0a, 0x2f, 0xf4, 0xb0, 0xef, 0xf8, 0x28, 0xc8, 0x18, 0x9e, 0xc0, 0xa0, 0x1f,
	0x29, 0x3c, 0xe1, 0x51, 0x2e, 0x02, 0x7f, 0x39, 0x07, 0xab, 0xc9, 0xeb, 0x1a, 0x79, 0x5a, 0xa5,
	0x95, 0xca, 0x25, 0xaf, 0x54, 0xfa, 0x35, 0x13, 0x55, 0xf3, 0x73, 0xd4, 0x3b, 0xc3, 0x91, 0xe6,
	0xe7, 0x98, 0x3f, 0xc8, 0xc1, 0xec, 0xee, 0x20, 0x08, 0xbd, 0x2e, 0xf2, 0xdf, 0x70, 0xac, 0x03,
	0xc7, 0x75, 0xc2, 0xd3, 0x28, 0x8b, 0x94, 0x48, 0xa1, 0x25, 0x29, 0xee, 0xa4, 0x86, 0x25, 0xcc,
	0xd2, 0xd3, 0x4d, 0xe5, 0xd3, 0xd2, 0x4d, 0xdd, 0x80, 0x69, 0x1f, 0xb5, 0x91, 0x73, 0x82, 0x44,
	0x72, 0x35, 0x59, 0x55, 0x98, 0xe2, 0x95, 0x0c, 0xdc, 0x84, 0x9a, 0xcb, 0x47, 0xa3, 0x66, 0x31,
	0x12, 0xc5, 0xe6, 0x7f, 0x2e, 0x42, 0x95, 0x0f, 0xfd, 0x6c, 0xdc, 0xa0, 0x24, 0x6b, 0x2b, 0x24,
	0x26, 0x6b, 0xe3, 0x28, 0x9e, 0x4c, 0x4d, 0x36, 0x18, 0x3f, 0x4c, 0x17, 0xaa, 0x43, 0x69, 0x48,
	0x82, 0xb7, 0x72, 0x3c, 0x85, 0xc2, 0x25, 0x28, 0x77, 0x50, 0xcf, 0x46, 0xbe, 0x72, 0x88, 0xca,
	0xca, 0x44, 0x94, 0x4f, 0x35, 0x16, 0xe5, 0xb3, 0x0e, 0xd5, 0x03, 0xc7, 0x0f, 0x8f, 0x6c, 0xeb,
	0x54, 0xd9, 0xc7, 0x44, 0xa9, 0x46, 0xe4, 0xb5, 0xb3, 0x13, 0x39, 0x8c, 0xa3, 0x17, 0xa4, 0x19,
	0xb1, 0x29, 0x27, 0xb8, 0x1b, 0x50, 0xeb, 0xf8, 0xde, 0xa0, 0x4f, 0xce, 0xd0, 0xa6, 0xd6, 0x0b,
	0x91, 0x76, 0x48, 0x8a, 0xf7, 0xec, 0xc0, 0x78, 0x51, 0x26, 0x8b, 0x69, 0x66, 0x05, 0x51, 0xf5,
	0x56, 0xa7, 0x61, 0x99, 0x52, 0xfe, 0x2c, 0x27, 0x42, 0x7a, 0x18, 0x98, 0xe4, 0x5b, 0x1e, 0x95,
	0xae, 0x2f, 0x5a, 0x95, 0x89, 0x84, 0x55, 0x91, 0x71, 0x3f, 0x99, 0x88, 0x7b, 0xbe, 0x6e, 0x53,
	0x89, 0xd1, 0x59, 0x98, 0x82, 0xf2, 0xe9, 0x14, 0x54, 0x18, 0x42, 0x41, 0xc5, 0x18, 0x05, 0x99,
	0x7f, 0x98, 0x13, 0xd1, 0xc8, 0xda, 0x44, 0x29, 0x7f, 0x94, 0x34, 0xfe, 0x58, 0x89, 0x4d, 0x5f,
	0x9a, 0xf8, 0xa2, 0x3a, 0x71, 0x31, 0xe5, 0x65, 0x7d, 0xca, 0xd2, 0x64, 0x0d, 0x79, 0xb2, 0x6c,
	0x9a, 0x86, 0x3c, 0x4d, 0x36, 0xc1, 0x79, 0x65, 0x82, 0x7c, 0x6a, 0xf3, 0xca, 0xd4, 0xf8, 0xa4,
	0x5a, 0xc4, 0xc2, 0xe6, 0x13, 0x7a, 0x0c, 0x47, 0x51, 0x9d, 0x48, 0x04, 0x46, 0xa6, 0xd4, 0x3f,
	0x80, 0x5a, 0x9b, 0x17, 0x6a, 0x96, 0x14, 0x47, 0x6d, 0x04, 0x90, 0x31, 0x55, 0xd4, 0x67, 0xe1,
	0xe2, 0xbe, 0x34, 0x93, 0x51, 0x11, 0xae, 0x37, 0x46, 0x9c, 0x15, 0x69, 0x2e, 0xc3, 0xaf, 0xe4,
	0x60, 0x89, 0xb7, 0x7d, 0x9b, 0x84, 0x29, 0xc8, 0x47, 0x39, 0x55, 0x3e, 0x58, 0x86, 0x2e, 0x7d,
	0x32, 0xa2, 0xde, 0xd8, 0x86, 0xc9, 0x80, 0xfa, 0xf6, 0x99, 0x43, 0x91, 0xa2, 0xee, 0x12, 0xbf,
	0xb3, 0x6a, 0xa3, 0x3e, 0x22, 0xc7, 0x00, 0xd2, 0x09, 0x40, 0x73, 0x22, 0x88, 0x7e, 0x04, 0xe6,
	0xb7, 0x73, 0xb0, 0x98, 0x0c, 0x38, 0xdc, 0xfe, 0x4f, 0x33, 0xdf, 0xc4, 0x21, 0x42, 0x21, 0xfd,
	0x10, 0xa1, 0x98, 0x7e, 0x88, 0x50, 0xd2, 0x0f, 0x11, 0xcc, 0xdb, 0xe4, 0xf4, 0x91, 0x4f, 0x9e,
	0x25, 0xe1, 0x89, 0xfc, 0x09, 0x57, 0xa1, 0xce, 0xd1, 0x11, 0x3f, 0xf7, 0x67, 0x15, 0x7b, 0xb6,
	0xf9, 0xfd, 0x22, 0x4c, 0x6b, 0x6d, 0xa4, 0x6c, 0x44, 0xeb, 0xe4, 0x0a, 0xd0, 0x89, 0xd3, 0x6b,
	0x6b, 0x9e, 0x54, 0x5e, 0x4a, 0x72, 0xfc, 0xb1, 0xbf, 0xe9, 0xf5, 0x50, 0x79, 0x83, 0x98, 0xe0,
	0x55, 0xe4, 0x92, 0xe8, 0x3a, 0x54, 0x6d, 0x27, 0x08, 0x7d, 0xa7, 0x1d, 0x6a, 0x4e, 0x77, 0x56,
	0x8a, 0x1b, 0xe3, 0x7f, 0xd3, 0xc6, 0xe4, 0x0d, 0x65, 0x82, 0x57, 0xf1, 0x1b, 0xa7, 0x0f, 0x2c,
	0xdf, 0x56, 0xb7, 0x31, 0x5c, 0x82, 0x65, 0x23, 0xfe, 0xbf, 0x15, 0x4b, 0xb6, 0x5b, 0xc5, 0xc5,
	0x7c, 0x24, 0x2c, 0x6f, 0xd1, 0xf3, 0xea, 0x6d, 0x34, 0x5e, 0x2a, 0x41, 0xdc, 0x54, 0x3c, 0xa2,
	0xa2, 0xd4, 0x58, 0x83, 0x0a, 0xd1, 0x7f, 0xfc, 0x53, 0xc5, 0x68, 0xe3, 0x85, 0xaa, 0x88, 0x9e,
	0x48, 0x14, 0xd1, 0xa4, 0x09, 0xa2, 0x68, 0x2b, 0xa6, 0x22, 0x2f, 0x8c, 0xc4, 0xe9, 0xf4, 0x10,
	0x71, 0x3a, 0x13, 0xdf, 0x90, 0xf1, 0xaa, 0x79, 0x81, 0x43, 0x42, 0xa4, 0x66, 0x95, 0x55, 0x63,
	0xa5, 0xc6, 0x0b, 0x50, 0x6f, 0x7b, 0x9e, 0x6f, 0x3b, 0x3d, 0x72, 0x05, 0x63, 0x5e, 0x4d, 0x31,
	0x24, 0x2a, 0x9a, 0x32, 0x94, 0xf9, 0x9b, 0x02, 0x5c, 0x52, 0xb7, 0x23, 0x9e, 0x05, 0x6a, 0x2c,
	0xf2, 0x8b, 0x93, 0x4c, 0x3e, 0x95, 0x64, 0x62, 0x04, 0x51, 0x48, 0x25, 0x08, 0x65, 0xd9, 0x8b,
	0xe7, 0xb0, 0xec, 0xe5, 0x51, 0xcb, 0x5e, 0x19, 0xb9, 0xec, 0xd5, 0x51, 0xcb, 0x3e, 0x31, 0x74,
	0xd9, 0x6b, 0x43, 0x96, 0x1d, 0x86, 0x2f, 0x7b, 0x3d, 0xcb, 0xb2, 0x1b, 0x99, 0x96, 0xfd, 0xbb,
	0x05, 0x9e, 0x57, 0x27, 0x65, 0xd9, 0x93, 0x45, 0xc7, 0x95, 0xc4, 0x55, 0xd6, 0xd6, 0xf7, 0x4a,
	0xe2, 0xfa, 0x6a, 0x2b, 0xbb, 0x12, 0x5b, 0x59, 0x69, 0x4d, 0x97, 0xf5, 0x35, 0x95, 0x56, 0x73,
	0x59, 0x5f, 0x4d, 0x69, 0x1d, 0x1b, 0xda, 0x3a, 0x46, 0x2b, 0xb8, 0x12, 0x5b, 0x41, 0x69, 0xed,
	0xe6, 0x95, 0xb5, 0x89, 0x29, 0x00, 0x20, 0x29, 0x00, 0xb8, 0x7b, 0x75, 0x3d, 0xa4, 0x95, 0x68,
	0x68, 0x34, 0x10, 0xad, 0xfe, 0x99, 0xd6, 0xe8, 0x2e, 0x5c, 0x4c, 0xd8, 0x14, 0x84, 0x4f, 0xbf,
	0x66, 0xf1, 0x42, 0xa6, 0x12, 0x2c, 0xa8, 0xbb, 0x28, 0x5f, 0xd0, 0x08, 0xce, 0xfc, 0x9c, 0xc8,
	0x8e, 0xcb, 0x9c, 0xa1, 0x99, 0x37, 0xfd, 0x42, 0x96, 0x4d, 0xff, 0x13, 0xb0, 0x92, 0xd8, 0x7c,
	0x3c, 0xa8, 0x2c, 0x9f, 0x10, 0x54, 0x66, 0xfe, 0x2a, 0x07, 0xcb, 0x2c, 0x90, 0xe5, 0x9e, 0x6f,
	0xd9, 0x4e, 0xaf, 0x33, 0xde, 0x85, 0x8a, 0x06, 0x14, 0x49, 0x64, 0xa8, 0xa2, 0xc2, 0xe2, 0x12,
	0xe3, 0x0a, 0x80, 0x8f, 0xc2, 0x81, 0xdf, 0x23, 0x6e, 0x13, 0xc5, 0xd5, 0x40, 0xcb, 0xdf, 0xf1,
	0x5d, 0x29, 0x16, 0xa9, 0x18, 0x8f, 0x45, 0x32, 0x9a, 0x30, 0xa3, 0xbf, 0x3c, 0xc0, 0x42, 0xb1,
	0x36, 0x36, 0x63, 0x4f, 0x12, 0xdc, 0xe5, 0x59, 0x70, 0xe9, 0x6f, 0xe1, 0x82, 0x56, 0x8b, 0xcd,
	0x17, 0x61, 0x25, 0x71, 0xbe, 0x0c, 0x5f, 0x8b, 0x50, 0xc0, 0xc3, 0x95, 0x95, 0x13, 0x5c, 0x60,
	0x9e, 0x70, 0xb7, 0x3e, 0x3b, 0x7e, 0x10, 0xe7, 0x19, 0x81, 0x74, 0x58, 0x34, 0x3a, 0xd5, 0xaf,
	0x1a, 0x7f, 0x99, 0x1f, 0x9d, 0xe9, 0xe7, 0xbf, 0xe5, 0x61, 0x95, 0x8d, 0x77, 0xf7, 0x08, 0xb5,
	0x8f, 0x9b, 0x04, 0x75, 0x24, 0x42, 0x21, 0x26, 0x34, 0x6a, 0x99, 0x72, 0xff, 0x3c, 0x0b, 0x53,
	0x1c, 0x77, 0x12, 0x9d, 0x09, 0x67, 0x5a, 0x5f, 0x8e, 0x6a, 0x1a, 0xb1, 0x3e, 0x2f, 0xc3, 0x52,
	0xe4, 0xac, 0x95, 0xb2, 0xe5, 0x31, 0x43, 0x43, 0xa4, 0xca, 0x13, 0x1e, 0xdb, 0x08, 0x66, 0xcf,
	0x4e, 0x5c, 0xdd, 0xf2, 0x23, 0xae, 0xee, 0xb7, 0x73, 0x30, 0x21, 0xdf, 0x71, 0x1e, 0xfb, 0x4c,
	0x63, 0xf4, 0x55, 0x4c, 0xd9, 0xa7, 0x50, 0x4c, 0xf0, 0x29, 0x5c, 0x52, 0xa2, 0x06, 0x0b, 0x1a,
	0xbb, 0xd2, 0x77, 0x29, 0x94, 0x8c, 0x67, 0xe7, 0x7d, 0x31, 0xfe, 0x01, 0x2c, 0xa9, 0xb9, 0x0a,
	0xc6, 0xbd, 0x3f, 0xae, 0xe5, 0x5e, 0xcb, 0x67, 0xcd, 0xbd, 0x16, 0xc0, 0xe5, 0x6d, 0xdb, 0x96,
	0xee, 0xb4, 0xc7, 0x13, 0x8d, 0x65, 0x0a, 0x6f, 0x7f, 0x36, 0x96, 0x3e, 0x23, 0x2f, 0x39, 0x03,
	0xd4, 0x24, 0x1a, 0xe6, 0x09, 0x98, 0xf4, 0xa6, 0xe3, 0x13, 0xee, 0x17, 0xc1, 0xc5, 0x6d, 0xdb,
	0xe6, 0x72, 0xff, 0x9e, 0x77, 0xc7, 0xf7, 0x06, 0x22, 0x6f, 0xe8, 0x35, 0x98, 0x90, 0xf4, 0x37,
	0x26, 0xde, 0x45, 0xa6, 0x7e, 0xa1, 0xc0, 0x91, 0xec, 0x15, 0xdc, 0xe5, 0xa1, 0x38, 0xa8, 0x2a,
	0xcc, 0xe3, 0x61, 0x1e, 0xf3, 0xe4, 0x6d, 0xbc, 0xa7, 0xb7, 0x07, 0xe1, 0xdb, 0x87, 0x8f, 0xa9,
	0xb3, 0x1f, 0xe6, 0x60, 0x96, 0x3f, 0x64, 0x10, 0xf9, 0x02, 0x3f, 0x02, 0x0b, 0xd4, 0xc5, 0x17,
	0xcf, 0xab, 0x1f, 0x89, 0x84, 0x39, 0x02, 0xa2, 0xe6, 0xe4, 0xcf, 0xee, 0x1c, 0xd4, 0x5e, 0x5d,
	0x28, 0x24, 0xbf, 0xba, 0x90, 0xc9, 0x29, 0xf8, 0xdf, 0x8b, 0x50, 0x15, 0x8f, 0x31, 0x3c, 0x46,
	0xa7, 0xe0, 0x39, 0xbb, 0xfe, 0xf4, 0xb7, 0x1d, 0x2a, 0xd9, 0xde, 0x76, 0xa8, 0x26, 0xbf, 0xed,
	0xf0, 0x02, 0xf0, 0xa7, 0x1b, 0xe4, 0xa7, 0x1d, 0x64, 0x2d, 0x78, 0x36, 0xaa, 0xe7, 0x86, 0x2b,
	0xdf, 0x32, 0x20, 0xb6, 0x65, 0x8c, 0xe9, 0xc7, 0x53, 0x9d, 0x8d, 0x13, 0x67, 0x77, 0x36, 0x4e,
	0x8e, 0xe3, 0x6c, 0x54, 0x5c, 0x83, 0x53, 0xb2, 0x6b, 0x30, 0x46, 0xd2, 0x32, 0xbd, 0x7c, 0x35,
	0xcf, 0x5d, 0x83, 0xe2, 0xa1, 0x9b, 0xec, 0xae, 0xc1, 0x73, 0x77, 0xe0, 0xc5, 0xe8, 0xa0, 0x94,
	0x8d, 0x0e, 0xca, 0xe3, 0xd0, 0x41, 0x65, 0x28, 0x1d, 0x98, 0x7f, 0x23, 0x7c, 0x88, 0x3a, 0x46,
	0x92, 0xd9, 0x49, 0x51, 0xf3, 0xf3, 0x9a, 0x9a, 0x9f, 0x94, 0x14, 0x44, 0xa8, 0xfe, 0xc5, 0x44,
	0xd5, 0xbf, 0x24, 0xab, 0xfe, 0x1b, 0x1a, 0x3e, 0x58, 0x9a, 0x0a, 0x19, 0x13, 0xab, 0x0a, 0x26,
	0x58, 0x9a, 0x8a, 0x08, 0x07, 0x37, 0x12, 0x71, 0x40, 0xcd, 0x91, 0x84, 0xd9, 0x53, 0x67, 0x23,
	0x9f, 0xf9, 0xe3, 0x71, 0x36, 0x4a, 0xad, 0x47, 0xce, 0x46, 0xfe, 0xf6, 0x85, 0xe6, 0x6c, 0x14,
	0x6b, 0x10, 0x01, 0x64, 0x74, 0x36, 0x7e, 0x27, 0x0f, 0x15, 0x76, 0x50, 0xf6, 0x01, 0x08, 0xc2,
	0x1b, 0x23, 0x2e, 0x4f, 0x0c, 0x95, 0x19, 0x4f, 0xea, 0x14, 0x6e, 0x9f, 0x5f, 0xc0, 0x64, 0x78,
	0x3a, 0x0f, 0xd6, 0x37, 0x3f, 0xc7, 0xd3, 0xba, 0x68, 0x8d, 0x9e, 0x0f, 0xf7, 0xb0, 0xcc, 0x4d,
	0xac, 0xed, 0xc7, 0x40, 0x9f, 0x6d, 0x98, 0x89, 0x1a, 0x67, 0xe4, 0xf9, 0x0c, 0x54, 0xd9, 0xc9,
	0x2d, 0xa7, 0xce, 0x49, 0x1e, 0x54, 0x44, 0xa7, 0x28, 0xaa, 0x33, 0xd2, 0xa6, 0x0d, 0xf5, 0x26,
	0x6a, 0x23, 0xa7, 0x1f, 0x92, 0x67, 0x52, 0x56, 0xa0, 0xec, 0xa3, 0x43, 0x5d, 0x35, 0x2b, 0xf9,
	0xe8, 0x70, 0x8f, 0x04, 0xe2, 0x84, 0x4e, 0xe8, 0x6a, 0x81, 0x38, 0xa4, 0x48, 0xb2, 0x4f, 0x0a,
	0xeb, 0xf9, 0xd8, 0x5d, 0x96, 0x6f, 0xe4, 0xa0, 0x8c, 0xad, 0xbc, 0x61, 0x2f, 0x33, 0xa4, 0x5c,
	0x9c, 0xcf, 0x40, 0xfb, 0xcb, 0x8a, 0x18, 0x53, 0x45, 0xbc, 0xf4, 0x4a, 0x52, 0x29, 0xe9, 0x95,
	0xa4, 0xaf, 0x57, 0xa0, 0xc2, 0x10, 0x70, 0x66, 0xde, 0x0c, 0xc9, 0xcc, 0x62, 0x76, 0x0a, 0x2d,
	0xa6, 0x57, 0x49, 0xb5, 0x97, 0x3a, 0x12, 0x9e, 0x7c, 0xc0, 0xfb, 0x0d, 0x03, 0x8a, 0xbd, 0x1c,
	0x51, 0x67, 0x35, 0xf7, 0xd8, 0x59, 0x55, 0xcc, 0x2b, 0x4d, 0xb5, 0x03, 0xb1, 0x3e, 0x95, 0xf8,
	0xfa, 0xa4, 0x9f, 0x59, 0x6a, 0x17, 0x50, 0x6b, 0x69, 0x17, 0x50, 0xa3, 0x15, 0x86, 0x04, 0x0b,
	0x74, 0x03, 0x6a, 0xf4, 0xb1, 0x0e, 0x3d, 0x07, 0x6b, 0x95, 0x16, 0x93, 0xbb, 0x79, 0x55, 0x4c,
	0x5b, 0x22, 0x15, 0x9d, 0xf0, 0x12, 0xfa, 0xe8, 0x90, 0xc5, 0x26, 0xb3, 0x87, 0x7e, 0x26, 0xe4,
	0x8b, 0x17, 0x12, 0x79, 0xf2, 0x57, 0xde, 0x22, 0x91, 0x36, 0x99, 0x45, 0xa4, 0x45, 0xb9, 0x1d,
	0xa6, 0xe2, 0xb9, 0x1d, 0x8c, 0x17, 0xa0, 0x42, 0x55, 0xde, 0x2c, 0x0f, 0xb3, 0x95, 0x89, 0x12,
	0xfc, 0x41, 0xbd, 0xc9, 0xa6, 0xbf, 0x43, 0x67, 0x3c, 0xda, 0x3b, 0x74, 0x73, 0xe3, 0xbd, 0x43,
	0xb7, 0x06, 0xc5, 0x41, 0x80, 0x7c, 0xe6, 0x5a, 0x67, 0x69, 0xe9, 0xdf, 0x09, 0x90, 0xdf, 0x24,
	0xe5, 0x58, 0xe0, 0x50, 0x92, 0x67, 0x09, 0x0e, 0xd9, 0x5b, 0x2e, 0x94, 0xef, 0x9b, 0xac, 0x0e,
	0x43, 0x51, 0x8a, 0x20, 0x17, 0x64, 0xf5, 0x17, 0x5f, 0x58, 0x9d, 0xf9, 0x27, 0x79, 0xbe, 0x21,
	0xb0, 0xe5, 0x97, 0x36, 0x84, 0x88, 0xdd, 0x72, 0x89, 0xec, 0x36, 0x4c, 0x4c, 0x71, 0x36, 0x28,
	0x8c, 0x62, 0x83, 0xe2, 0x68, 0x36, 0x28, 0x8d, 0x62, 0x83, 0xf2, 0x48, 0x36, 0xa8, 0x25, 0xb1,
	0x81, 0x44, 0x90, 0x95, 0xf5, 0x7c, 0x46, 0x82, 0x14, 0xbc, 0x53, 0x1d, 0xce, 0x3b, 0xe6, 0xff,
	0x11, 0xf9, 0xe5, 0x34, 0xcc, 0xa6, 0xee, 0x8a, 0x11, 0xbe, 0xe9, 0xe3, 0x60, 0x11, 0xa6, 0xe7,
	0x39, 0xa6, 0xd9, 0xd9, 0x31, 0xc5, 0xf1, 0x7a, 0x02, 0x26, 0x55, 0x1c, 0x2e, 0xaa, 0x38, 0x14,
	0xd8, 0x5b, 0x89, 0x61, 0x4f, 0xc2, 0xdb, 0x45, 0x1d, 0x6f, 0x11, 0xc6, 0x3e, 0x2a, 0x63, 0x6c,
	0x04, 0x4d, 0xf3, 0x05, 0x1b, 0x17, 0x6f, 0xaf, 0xc3, 0xfc, 0x2e, 0xcb, 0x4a, 0x93, 0x01, 0x6d,
	0xc3, 0xd3, 0xc9, 0x50, 0xcd, 0x81, 0x35, 0xf4, 0x18, 0x34, 0x87, 0x7f, 0x99, 0x83, 0x35, 0xa9,
	0x83, 0x9d, 0x53, 0xca, 0x5d, 0x24, 0x93, 0x46, 0x74, 0xab, 0x33, 0x76, 0xeb, 0x88, 0x32, 0x82,
	0xa9, 0xe9, 0x0d, 0x23, 0x86, 0x51, 0x48, 0x1f, 0xc6, 0x5f, 0xe7, 0x60, 0x26, 0x9a, 0x24, 0xd3,
	0x60, 0xf6, 0x60, 0x4d, 0x76, 0x45, 0xb4, 0x22, 0x81, 0xe6, 0x53, 0x50, 0x05, 0x97, 0x2b, 0x92,
	0x73, 0x62, 0x97, 0x43, 0xf2, 0xfd, 0x3a, 0xbd, 0x29, 0xe6, 0x8d, 0x54, 0x36, 0xec, 0xc4, 0xa6,
	0x98, 0x37, 0x13, 0xeb, 0x55, 0xac, 0x7b, 0x3e, 0xa3, 0x49, 0x85, 0x12, 0x9a, 0xa2, 0x5a, 0xd2,
	0xab, 0x8a, 0x43, 0xf4, 0xaa, 0x4f, 0xc1, 0x53, 0xd8, 0x7a, 0x51, 0xdd, 0x7f, 0x3b, 0xa7, 0xcc,
	0x2f, 0x16, 0xe5, 0x7d, 0xca, 0xe2, 0x0f, 0x33, 0x1f, 0xc0, 0xaa, 0x78, 0xe0, 0x56, 0x7d, 0x37,
	0x53, 0x64, 0x82, 0xe2, 0x6a, 0xdb, 0x4c, 0x5c, 0x6d, 0x93, 0x05, 0xce, 0x6c, 0x92, 0xc0, 0xe1,
	0xb4, 0x30, 0xad, 0xd3, 0x82, 0xf9, 0xcb, 0x3c, 0xcc, 0xeb, 0x7d, 0x66, 0x7f, 0x4e, 0xaf, 0x21,
	0x27, 0x45, 0xd4, 0x83, 0xe2, 0xf8, 0xe7, 0x91, 0x31, 0xc9, 0xa5, 0x2d, 0xab, 0x79, 0x2b, 0x9e,
	0x4c, 0xa5, 0x92, 0x2d, 0x0b, 0x47, 0x35, 0xd3, 0x23, 0x7a, 0xb5, 0xc4, 0x87, 0x75, 0x5e, 0xcc,
	0x90, 0x54, 0x92, 0x3b, 0xc4, 0xa4, 0xec, 0x0c, 0x58, 0x05, 0x25, 0x29, 0x2a, 0xe4, 0xb3, 0x1d,
	0x5a, 0xa4, 0x64, 0xbc, 0x2e, 0x25, 0x66, 0xbc, 0xee, 0xc0, 0x5a, 0xda, 0xfa, 0x32, 0xae, 0x79,
	0x05, 0x8c, 0x58, 0xfa, 0x2b, 0x6e, 0x01, 0xa4, 0xbd, 0xa9, 0x3a, 0xab, 0xbf, 0xa9, 0x1a, 0x98,
	0x2f, 0xc1, 0x1a, 0x23, 0xfd, 0x34, 0x4a, 0x4a, 0x14, 0x67, 0xe6, 0x21, 0x5c, 0x4e, 0xfd, 0x8e,
	0x8d, 0x30, 0xf1, 0xd1, 0xd7, 0xdc, 0x78, 0x8f, 0xbe, 0x9a, 0xfb, 0xb0, 0x4a, 0x85, 0xec, 0x58,
	0xc3, 0x1b, 0x21, 0x6d, 0xef, 0xc3, 0x5a, 0x5a, 0xa3, 0x6c, 0xec, 0xb7, 0xa0, 0x26, 0x86, 0x32,
	0x62, 0xcc, 0x11, 0xa0, 0xf9, 0x57, 0x39, 0x58, 0xa5, 0x3b, 0xe9, 0x78, 0xa3, 0x5d, 0x89, 0x59,
	0x0c, 0xd2, 0x96, 0xfa, 0x92, 0xfa, 0xe4, 0xe7, 0x72, 0xf2, 0x30, 0x30, 0x0f, 0x8a, 0xb0, 0x1c,
	0xa2, 0x12, 0x27, 0x85, 0x76, 0x71, 0xbe, 0x2e, 0xc6, 0x64, 0xfc, 0xbc, 0x62, 0x29, 0xf0, 0x8d,
	0x5b, 0xf7, 0x12, 0x57, 0x63, 0xb1, 0xa6, 0xec, 0xd1, 0x5d, 0xc4, 0xcf, 0xe0, 0x1e, 0x2f, 0x35,
	0xfc, 0xf3, 0x1c, 0xac, 0x6d, 0xdb, 0x5f, 0x1c, 0x04, 0xa1, 0x00, 0xfe, 0xb4, 0xc8, 0xa3, 0x42,
	0x31, 0xac, 0xf2, 0x85, 0x9a, 0xba, 0x71, 0x31, 0x39, 0xb5, 0xbf, 0xcc, 0x17, 0x3c, 0x91, 0x63,
	0xba, 0x0f, 0xe1, 0x7f, 0xe6, 0xe0, 0x72, 0xea, 0x18, 0x12, 0x99, 0xf3, 0xcc, 0x83, 0x48, 0xe6,
	0xf1, 0xfc, 0xb8, 0x3c, 0x7e, 0x87, 0x04, 0x54, 0xe9, 0x1d, 0x0a, 0x2d, 0xe3, 0x7a, 0xba, 0x96,
	0x11, 0xa9, 0x46, 0x64, 0x0b, 0x43, 0x24, 0x97, 0x46, 0x42, 0x43, 0xe7, 0x3a, 0x6d, 0xf3, 0x75,
	0xb8, 0x96, 0xd4, 0xcd, 0x0e, 0xff, 0x4b, 0x4a, 0x78, 0xa8, 0x65, 0x23, 0xcc, 0xc5, 0xb2, 0x11,
	0xfe, 0x65, 0x1e, 0x66, 0xf4, 0x96, 0x64, 0x03, 0x3e, 0x97, 0x60, 0xc0, 0xab, 0x5b, 0x59, 0x3e,
	0x79, 0x2b, 0xfb, 0xad, 0x78, 0x39, 0xe1, 0x03, 0x7a, 0x36, 0xf2, 0xdf, 0xe7, 0xa0, 0xa1, 0x23,
	0x9d, 0xf3, 0x48, 0x22, 0xe2, 0x72, 0xe3, 0x21, 0x2e, 0x3f, 0x04, 0x71, 0x32, 0x4e, 0x8a, 0x89,
	0x7b, 0xeb, 0x5f, 0x94, 0x65, 0x92, 0x60, 0xd9, 0x3b, 0x93, 0x5e, 0xb8, 0x4f, 0x92, 0x74, 0x9a,
	0xeb, 0xa6, 0x9a, 0xec, 0xba, 0xb9, 0x12, 0xe1, 0xee, 0x80, 0xde, 0x36, 0x16, 0x40, 0xac, 0x7c,
	0xe7, 0xd4, 0x78, 0x8e, 0x0b, 0xf6, 0xdc, 0x28, 0xc1, 0xce, 0x45, 0xfa, 0x46, 0xcc, 0xf4, 0x8a,
	0x99, 0xba, 0xe9, 0xaf, 0xc3, 0x0e, 0x93, 0xfd, 0x49, 0xf1, 0xc7, 0x12, 0x93, 0x94, 0x13, 0x98,
	0x44, 0xd8, 0xd4, 0xd3, 0x71, 0x9b, 0xfa, 0xd1, 0x94, 0x4f, 0x06, 0x40, 0x14, 0x37, 0x43, 0x03,
	0x20, 0x4a, 0x1b, 0x03, 0x20, 0x9a, 0xe4, 0x82, 0x06, 0xc0, 0xe3, 0x26, 0x53, 0x92, 0x7f, 0xa9,
	0xfc, 0x00, 0x67, 0xe7, 0x87, 0xfa, 0xb8, 0x5e, 0x1a, 0xd9, 0xcd, 0x32, 0x31, 0x9e, 0x9b, 0x45,
	0x77, 0xf2, 0x4c, 0x8e, 0xe7, 0xe4, 0x89, 0x25, 0x56, 0x9d, 0x4a, 0x4d, 0xac, 0x1a, 0x39, 0x6c,
	0xe6, 0x87, 0x38, 0x6c, 0x22, 0x67, 0xdb, 0x62, 0x96, 0x30, 0xa8, 0x8e, 0xb6, 0xfd, 0xb0, 0x6d,
	0xe9, 0xfc, 0x8d, 0xdc, 0x5b, 0xb0, 0x9e, 0xd4, 0xd1, 0xce, 0xa9, 0xb4, 0x61, 0xc4, 0x93, 0x2c,
	0xeb, 0x9b, 0x9a, 0x18, 0xde, 0xf9, 0x2a, 0xda, 0xff, 0x10, 0x26, 0x79, 0x9c, 0x00, 0x09, 0x11,
	0x18, 0xfb, 0xa1, 0x8b, 0x97, 0x78, 0xc2, 0x55, 0xa5, 0x99, 0x91, 0xd9, 0xd6, 0xcc, 0xcf, 0xf0,
	0x20, 0xb7, 0xc4, 0xef, 0xe4, 0xc0, 0x83, 0x5c, 0x42, 0xe0, 0xc1, 0x90, 0x01, 0xd9, 0xd0, 0x90,
	0xe2, 0xb4, 0x49, 0xab, 0x8f, 0x61, 0x51, 0xbf, 0x04, 0x8b, 0x7a, 0x17, 0x51, 0x26, 0x59, 0x11,
	0x5c, 0x41, 0x46, 0xab, 0x65, 0x92, 0x55, 0xe7, 0x3b, 0xd5, 0x56, 0x5a, 0xc9, 0x78, 0x18, 0xb2,
	0x43, 0x68, 0x97, 0xc6, 0x4c, 0xec, 0x9c, 0x32, 0xd3, 0x5f, 0xb1, 0xd5, 0x99, 0x17, 0x20, 0x66,
	0xf4, 0xb2, 0xf2, 0x3d, 0xdb, 0xdc, 0x86, 0xe5, 0x04, 0x75, 0x66, 0x9c, 0xd8, 0x34, 0xf3, 0x3e,
	0x18, 0x74, 0xed, 0x77, 0x7c, 0xab, 0x67, 0x8f, 0xce, 0xb0, 0xa7, 0xb9, 0x2e, 0xf3, 0x29, 0xae,
	0x4b, 0xd3, 0xe5, 0xaf, 0x8d, 0x2a, 0xed, 0x26, 0xa7, 0x35, 0x7b, 0xf4, 0xde, 0x9e, 0x87, 0xb9,
	0xdb, 0xe4, 0x90, 0x85, 0xf5, 0xc6, 0xd6, 0x71, 0x19, 0x4a, 0xed, 0x58, 0xce, 0x33, 0x5a, 0x64,
	0xfe, 0x26, 0x07, 0x25, 0x02, 0x3d, 0x4a, 0x67, 0xa3, 0x63, 0xce, 0xa7, 0x8c, 0xf9, 0xec, 0x8f,
	0xba, 0xa8, 0x1b, 0x46, 0xe9, 0xec, 0x1b, 0x46, 0x79, 0x1c, 0x05, 0xea, 0x13, 0x34, 0x43, 0x1c,
	0x9e, 0x38, 0x17, 0x5f, 0x0c, 0x59, 0x57, 0xa0, 0x4c, 0x72, 0xda, 0x72, 0x5a, 0xaf, 0x53, 0x5a,
	0xa7, 0x18, 0x65, 0x55, 0xe6, 0xcb, 0x24, 0x6d, 0x3c, 0xfd, 0x7c, 0x7c, 0x2d, 0xff, 0xf3, 0x24,
	0xbd, 0x0d, 0xff, 0x7a, 0x8c, 0x7e, 0x33, 0xf2, 0xd4, 0x21, 0x3c, 0x9d, 0x2c, 0xa6, 0x9b, 0xe8,
	0x10, 0xf9, 0xa8, 0x27, 0xb2, 0x48, 0x0f, 0x3f, 0x7b, 0x94, 0xf5, 0x88, 0x7c, 0x82, 0x1e, 0x61,
	0xf6, 0x35, 0x33, 0x22, 0xa9, 0x9f, 0xf3, 0x95, 0xf1, 0xff, 0x3b, 0xc7, 0xcd, 0x60, 0x22, 0x31,
	0x78, 0xe2, 0x3b, 0x32, 0xfb, 0xac, 0x51, 0xbb, 0xcf, 0x41, 0x95, 0x67, 0xb6, 0x63, 0x58, 0x9c,
	0xdf, 0x24, 0x75, 0x9b, 0x4a, 0x9b, 0x4d, 0x01, 0x65, 0x7c, 0x12, 0x66, 0x44, 0xe2, 0x3c, 0x1e,
	0x43, 0x41, 0x33, 0x2c, 0xce, 0xc9, 0x5f, 0xf2, 0xf8, 0xe7, 0x69, 0x0e, 0xcc, 0xc3, 0x2a, 0xde,
	0x24, 0xf2, 0x69, 0x3f, 0xf4, 0xda, 0xc7, 0xa1, 0x75, 0x8c, 0x34, 0x82, 0xdb, 0x02, 0x08, 0x44,
	0x15, 0x43, 0x09, 0x73, 0x7f, 0x89, 0x4f, 0x9a, 0x12, 0x88, 0xf9, 0xaf, 0x73, 0xb0, 0xc8, 0xa2,
	0x76, 0x44, 0x3d, 0x9b, 0xfc, 0xb3, 0x30, 0x45, 0x95, 0x6c, 0xa1, 0xa1, 0xcb, 0x2c, 0x3f, 0x49,
	0xea, 0x84, 0x9d, 0x90, 0x1e, 0xc0, 0xf3, 0x0c, 0x57, 0x9c, 0xd5, 0x7c, 0x6f, 0xbc, 0x37, 0xd9,
	0x47, 0xff, 0x9d, 0x1c, 0x7f, 0xa9, 0x22, 0x36, 0x18, 0x2a, 0x31, 0x8a, 0x9a, 0xc4, 0x78, 0xf2,
	0x43, 0xfc, 0x65, 0x11, 0x6a, 0xa2, 0xe2, 0x6c, 0x27, 0xce, 0xf1, 0x41, 0x17, 0x46, 0x0f, 0xba,
	0x98, 0x64, 0x21, 0xa4, 0xdc, 0xa7, 0x1d, 0x43, 0x85, 0x54, 0xed, 0xa3, 0x52, 0x16, 0xfb, 0xa8,
	0x9c, 0x6c, 0x1f, 0x7d, 0x30, 0x09, 0x33, 0x74, 0x65, 0xbb, 0xf6, 0x68, 0x27, 0xaa, 0x30, 0x9e,
	0xaa, 0x3f, 0x66, 0x62, 0x65, 0x41, 0x66, 0xf5, 0x91, 0x64, 0x46, 0x5f, 0xfd, 0x8a, 0xb8, 0xfc,
	0xfc, 0x35, 0xb5, 0x1e, 0x89, 0x9f, 0x96, 0x7b, 0x38, 0xa3, 0x08, 0xc9, 0xb8, 0x8f, 0x7c, 0xad,
	0x00, 0x93, 0xca, 0x54, 0xb5, 0xb3, 0x81, 0x5a, 0xb6, 0xb3, 0x01, 0x48, 0x3b, 0x1b, 0xd0, 0x8f,
	0x24, 0x8a, 0x43, 0x8e, 0x24, 0x46, 0x1f, 0x7d, 0x5c, 0x83, 0x09, 0xcf, 0xb5, 0x23, 0x0e, 0x55,
	0xf2, 0x6d, 0x79, 0xae, 0x2d, 0xf8, 0xf3, 0x1a, 0x4c, 0xf4, 0xd0, 0x83, 0x64, 0x56, 0xae, 0xf7,
	0xd0, 0x03, 0x99, 0x91, 0x09, 0xbb, 0x96, 0x62, 0xec, 0xaa, 0x9c, 0x6a, 0x94, 0x13, 0x4f, 0x35,
	0x54, 0xff, 0x52, 0x35, 0xcb, 0xe3, 0xc7, 0x95, 0xd1, 0x57, 0x22, 0x8e, 0x61, 0x89, 0x9d, 0x0a,
	0xa4, 0x88, 0x5b, 0x5d, 0xb0, 0x7d, 0x6c, 0xf8, 0x53, 0x2b, 0x43, 0x5f, 0xf2, 0xb8, 0x43, 0x8e,
	0x34, 0x23, 0xef, 0x20, 0x0f, 0xc7, 0x8b, 0xd4, 0xf3, 0xab, 0x50, 0xe7, 0x91, 0x79, 0xb1, 0xbb,
	0x7a, 0xbc, 0x62, 0xcf, 0x36, 0x5f, 0x21, 0x0d, 0x89, 0xc8, 0x3f, 0xc9, 0xcf, 0x38, 0x96, 0x92,
	0xfe, 0x8f, 0x61, 0x91, 0x9e, 0x2a, 0x64, 0x9c, 0x7b, 0xf6, 0x07, 0x4f, 0xcc, 0x1e, 0x5c, 0xd3,
	0xdd, 0xeb, 0xb4, 0xdb, 0x5d, 0xbe, 0x56, 0x29, 0x7e, 0x76, 0xf6, 0xb2, 0x6d, 0x3e, 0xd9, 0xcf,
	0x9e, 0xf6, 0x9c, 0xad, 0xe9, 0xf3, 0x3c, 0xbe, 0x43, 0xfa, 0x1b, 0xe3, 0x6a, 0x8d, 0x4a, 0x6c,
	0xf9, 0x44, 0x62, 0xdb, 0x79, 0xfd, 0xc7, 0xef, 0xad, 0xe5, 0x7e, 0xfa, 0xde, 0x5a, 0xee, 0xe7,
	0xef, 0xad, 0x5d, 0xf8, 0xf5, 0x7b, 0x6b, 0x17, 0xbe, 0xf2, 0xfe, 0xda, 0x85, 0xff, 0xf2, 0xfe,
	0xda, 0x85, 0x1f, 0xbc, 0xbf, 0x76, 0xe1, 0xc7, 0xef, 0xaf, 0x5d, 0xf8, 0xe9, 0xfb, 0x6b, 0x17,
	0x7e, 0xfe, 0xfe, 0xda, 0x85, 0xaf, 0xfd, 0x62, 0xed, 0xc2, 0x37, 0x7f, 0xb1, 0x76, 0xe1, 0x9f,
	0x34, 0x48, 0xbc, 0xc9, 0x49, 0x6f, 0xcb, 0xea, 0x3b, 0x5b, 0xfd, 0x83, 0x2d, 0xfc, 0x73, 0x0b,
	0x4f, 0xf1, 0xef, 0x02, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x15, 0xfe, 0xe2, 0x83, 0x96, 0x00, 0x00,
}
