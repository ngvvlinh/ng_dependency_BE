// source: services/affiliate/affiliate.proto

package affiliate

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	status4 "etop.vn/api/pb/etop/etc/status4"
	order "etop.vn/api/pb/etop/order"
	shop "etop.vn/api/pb/etop/shop"
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

type UpdateReferralRequest struct {
	ReferralCode         string   `protobuf:"bytes,2,opt,name=referral_code,json=referralCode" json:"referral_code"`
	SaleReferralCode     string   `protobuf:"bytes,3,opt,name=sale_referral_code,json=saleReferralCode" json:"sale_referral_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateReferralRequest) Reset()         { *m = UpdateReferralRequest{} }
func (m *UpdateReferralRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateReferralRequest) ProtoMessage()    {}
func (*UpdateReferralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{0}
}

var xxx_messageInfo_UpdateReferralRequest proto.InternalMessageInfo

func (m *UpdateReferralRequest) GetReferralCode() string {
	if m != nil {
		return m.ReferralCode
	}
	return ""
}

func (m *UpdateReferralRequest) GetSaleReferralCode() string {
	if m != nil {
		return m.SaleReferralCode
	}
	return ""
}

type UserReferral struct {
	UserId               dot.ID   `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id"`
	ReferralCode         string   `protobuf:"bytes,2,opt,name=referral_code,json=referralCode" json:"referral_code"`
	SaleReferralCode     string   `protobuf:"bytes,3,opt,name=sale_referral_code,json=saleReferralCode" json:"sale_referral_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserReferral) Reset()         { *m = UserReferral{} }
func (m *UserReferral) String() string { return proto.CompactTextString(m) }
func (*UserReferral) ProtoMessage()    {}
func (*UserReferral) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{1}
}

var xxx_messageInfo_UserReferral proto.InternalMessageInfo

func (m *UserReferral) GetUserId() dot.ID {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *UserReferral) GetReferralCode() string {
	if m != nil {
		return m.ReferralCode
	}
	return ""
}

func (m *UserReferral) GetSaleReferralCode() string {
	if m != nil {
		return m.SaleReferralCode
	}
	return ""
}

type SellerCommission struct {
	Id                   dot.ID            `protobuf:"varint,1,opt,name=id" json:"id"`
	Value                int32             `protobuf:"varint,4,opt,name=value" json:"value"`
	Description          string            `protobuf:"bytes,5,opt,name=description" json:"description"`
	Note                 string            `protobuf:"bytes,6,opt,name=note" json:"note"`
	Status               status4.Status    `protobuf:"varint,8,opt,name=status,enum=status4.Status" json:"status"`
	Type                 string            `protobuf:"bytes,9,opt,name=type" json:"type"`
	OValue               int32             `protobuf:"varint,16,opt,name=o_value,json=oValue" json:"o_value"`
	OBaseValue           int32             `protobuf:"varint,17,opt,name=o_base_value,json=oBaseValue" json:"o_base_value"`
	Product              *shop.ShopProduct `protobuf:"bytes,13,opt,name=product" json:"product"`
	Order                *order.Order      `protobuf:"bytes,14,opt,name=order" json:"order"`
	FromSeller           *etop.Affiliate   `protobuf:"bytes,15,opt,name=from_seller,json=fromSeller" json:"from_seller"`
	ValidAt              dot.Time          `protobuf:"bytes,10,opt,name=valid_at,json=validAt" json:"valid_at"`
	CreatedAt            dot.Time          `protobuf:"bytes,11,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time          `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SellerCommission) Reset()         { *m = SellerCommission{} }
func (m *SellerCommission) String() string { return proto.CompactTextString(m) }
func (*SellerCommission) ProtoMessage()    {}
func (*SellerCommission) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{2}
}

var xxx_messageInfo_SellerCommission proto.InternalMessageInfo

func (m *SellerCommission) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SellerCommission) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *SellerCommission) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SellerCommission) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *SellerCommission) GetStatus() status4.Status {
	if m != nil {
		return m.Status
	}
	return status4.Status_Z
}

func (m *SellerCommission) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *SellerCommission) GetOValue() int32 {
	if m != nil {
		return m.OValue
	}
	return 0
}

func (m *SellerCommission) GetOBaseValue() int32 {
	if m != nil {
		return m.OBaseValue
	}
	return 0
}

func (m *SellerCommission) GetProduct() *shop.ShopProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *SellerCommission) GetOrder() *order.Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *SellerCommission) GetFromSeller() *etop.Affiliate {
	if m != nil {
		return m.FromSeller
	}
	return nil
}

type GetCommissionsResponse struct {
	Commissions          []*SellerCommission `protobuf:"bytes,1,rep,name=commissions" json:"commissions"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetCommissionsResponse) Reset()         { *m = GetCommissionsResponse{} }
func (m *GetCommissionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetCommissionsResponse) ProtoMessage()    {}
func (*GetCommissionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{3}
}

var xxx_messageInfo_GetCommissionsResponse proto.InternalMessageInfo

func (m *GetCommissionsResponse) GetCommissions() []*SellerCommission {
	if m != nil {
		return m.Commissions
	}
	return nil
}

type NotifyNewShopPurchaseRequest struct {
	OrderId              dot.ID   `protobuf:"varint,1,opt,name=order_id,json=orderId" json:"order_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyNewShopPurchaseRequest) Reset()         { *m = NotifyNewShopPurchaseRequest{} }
func (m *NotifyNewShopPurchaseRequest) String() string { return proto.CompactTextString(m) }
func (*NotifyNewShopPurchaseRequest) ProtoMessage()    {}
func (*NotifyNewShopPurchaseRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{4}
}

var xxx_messageInfo_NotifyNewShopPurchaseRequest proto.InternalMessageInfo

func (m *NotifyNewShopPurchaseRequest) GetOrderId() dot.ID {
	if m != nil {
		return m.OrderId
	}
	return 0
}

type NotifyNewShopPurchaseResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message" json:"message"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyNewShopPurchaseResponse) Reset()         { *m = NotifyNewShopPurchaseResponse{} }
func (m *NotifyNewShopPurchaseResponse) String() string { return proto.CompactTextString(m) }
func (*NotifyNewShopPurchaseResponse) ProtoMessage()    {}
func (*NotifyNewShopPurchaseResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{5}
}

var xxx_messageInfo_NotifyNewShopPurchaseResponse proto.InternalMessageInfo

func (m *NotifyNewShopPurchaseResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetCouponsResponse struct {
	Coupons              []*Coupon `protobuf:"bytes,1,rep,name=coupons" json:"coupons"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetCouponsResponse) Reset()         { *m = GetCouponsResponse{} }
func (m *GetCouponsResponse) String() string { return proto.CompactTextString(m) }
func (*GetCouponsResponse) ProtoMessage()    {}
func (*GetCouponsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{6}
}

var xxx_messageInfo_GetCouponsResponse proto.InternalMessageInfo

func (m *GetCouponsResponse) GetCoupons() []*Coupon {
	if m != nil {
		return m.Coupons
	}
	return nil
}

type CreateCouponRequest struct {
	Value                int32    `protobuf:"varint,1,opt,name=value" json:"value"`
	Unit                 *string  `protobuf:"bytes,2,opt,name=unit" json:"unit"`
	Description          *string  `protobuf:"bytes,3,opt,name=description" json:"description"`
	ProductId            dot.ID   `protobuf:"varint,4,opt,name=product_id,json=productId" json:"product_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateCouponRequest) Reset()         { *m = CreateCouponRequest{} }
func (m *CreateCouponRequest) String() string { return proto.CompactTextString(m) }
func (*CreateCouponRequest) ProtoMessage()    {}
func (*CreateCouponRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{7}
}

var xxx_messageInfo_CreateCouponRequest proto.InternalMessageInfo

func (m *CreateCouponRequest) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *CreateCouponRequest) GetUnit() string {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return ""
}

func (m *CreateCouponRequest) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *CreateCouponRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

type Coupon struct {
	Id                   dot.ID   `protobuf:"varint,1,opt,name=id" json:"id"`
	Code                 string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Value                int32    `protobuf:"varint,3,opt,name=value" json:"value"`
	Unit                 string   `protobuf:"bytes,4,opt,name=unit" json:"unit"`
	Description          *string  `protobuf:"bytes,5,opt,name=description" json:"description"`
	UserId               dot.ID   `protobuf:"varint,6,opt,name=user_id,json=userId" json:"user_id"`
	StartDate            dot.Time `protobuf:"bytes,7,opt,name=start_date,json=startDate" json:"start_date"`
	EndDate              dot.Time `protobuf:"bytes,8,opt,name=end_date,json=endDate" json:"end_date"`
	ProductId            dot.ID   `protobuf:"varint,9,opt,name=product_id,json=productId" json:"product_id"`
	CreatedAt            dot.Time `protobuf:"bytes,10,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Coupon) Reset()         { *m = Coupon{} }
func (m *Coupon) String() string { return proto.CompactTextString(m) }
func (*Coupon) ProtoMessage()    {}
func (*Coupon) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{8}
}

var xxx_messageInfo_Coupon proto.InternalMessageInfo

func (m *Coupon) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Coupon) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Coupon) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Coupon) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *Coupon) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *Coupon) GetUserId() dot.ID {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *Coupon) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

type GetTransactionsResponse struct {
	Transactions         []*Transaction `protobuf:"bytes,1,rep,name=transactions" json:"transactions"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetTransactionsResponse) Reset()         { *m = GetTransactionsResponse{} }
func (m *GetTransactionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetTransactionsResponse) ProtoMessage()    {}
func (*GetTransactionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{9}
}

var xxx_messageInfo_GetTransactionsResponse proto.InternalMessageInfo

func (m *GetTransactionsResponse) GetTransactions() []*Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type Transaction struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{10}
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

type CommissionSetting struct {
	ProductId            dot.ID   `protobuf:"varint,2,opt,name=product_id,json=productId" json:"product_id"`
	Amount               int32    `protobuf:"varint,3,opt,name=amount" json:"amount"`
	Unit                 string   `protobuf:"bytes,4,opt,name=unit" json:"unit"`
	CreatedAt            dot.Time `protobuf:"bytes,5,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt            dot.Time `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommissionSetting) Reset()         { *m = CommissionSetting{} }
func (m *CommissionSetting) String() string { return proto.CompactTextString(m) }
func (*CommissionSetting) ProtoMessage()    {}
func (*CommissionSetting) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{11}
}

var xxx_messageInfo_CommissionSetting proto.InternalMessageInfo

func (m *CommissionSetting) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CommissionSetting) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *CommissionSetting) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

type SupplyCommissionSetting struct {
	ProductId                dot.ID                                 `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	Level1DirectCommission   int32                                  `protobuf:"varint,2,opt,name=level1_direct_commission,json=level1DirectCommission" json:"level1_direct_commission"`
	Level1IndirectCommission int32                                  `protobuf:"varint,3,opt,name=level1_indirect_commission,json=level1IndirectCommission" json:"level1_indirect_commission"`
	Level2DirectCommission   int32                                  `protobuf:"varint,4,opt,name=level2_direct_commission,json=level2DirectCommission" json:"level2_direct_commission"`
	Level2IndirectCommission int32                                  `protobuf:"varint,5,opt,name=level2_indirect_commission,json=level2IndirectCommission" json:"level2_indirect_commission"`
	DependOn                 string                                 `protobuf:"bytes,6,opt,name=depend_on,json=dependOn" json:"depend_on"`
	Level1LimitCount         int32                                  `protobuf:"varint,7,opt,name=level1_limit_count,json=level1LimitCount" json:"level1_limit_count"`
	MLifetimeDuration        *SupplyCommissionSettingDurationObject `protobuf:"bytes,8,opt,name=m_lifetime_duration,json=mLifetimeDuration" json:"m_lifetime_duration"`
	MLevel1LimitDuration     *SupplyCommissionSettingDurationObject `protobuf:"bytes,11,opt,name=m_level1_limit_duration,json=mLevel1LimitDuration" json:"m_level1_limit_duration"`
	CreatedAt                dot.Time                               `protobuf:"bytes,9,opt,name=created_at,json=createdAt" json:"created_at"`
	UpdatedAt                dot.Time                               `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt" json:"updated_at"`
	Group                    string                                 `protobuf:"bytes,12,opt,name=group" json:"group"`
	XXX_NoUnkeyedLiteral     struct{}                               `json:"-"`
	XXX_sizecache            int32                                  `json:"-"`
}

func (m *SupplyCommissionSetting) Reset()         { *m = SupplyCommissionSetting{} }
func (m *SupplyCommissionSetting) String() string { return proto.CompactTextString(m) }
func (*SupplyCommissionSetting) ProtoMessage()    {}
func (*SupplyCommissionSetting) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{12}
}

var xxx_messageInfo_SupplyCommissionSetting proto.InternalMessageInfo

func (m *SupplyCommissionSetting) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *SupplyCommissionSetting) GetLevel1DirectCommission() int32 {
	if m != nil {
		return m.Level1DirectCommission
	}
	return 0
}

func (m *SupplyCommissionSetting) GetLevel1IndirectCommission() int32 {
	if m != nil {
		return m.Level1IndirectCommission
	}
	return 0
}

func (m *SupplyCommissionSetting) GetLevel2DirectCommission() int32 {
	if m != nil {
		return m.Level2DirectCommission
	}
	return 0
}

func (m *SupplyCommissionSetting) GetLevel2IndirectCommission() int32 {
	if m != nil {
		return m.Level2IndirectCommission
	}
	return 0
}

func (m *SupplyCommissionSetting) GetDependOn() string {
	if m != nil {
		return m.DependOn
	}
	return ""
}

func (m *SupplyCommissionSetting) GetLevel1LimitCount() int32 {
	if m != nil {
		return m.Level1LimitCount
	}
	return 0
}

func (m *SupplyCommissionSetting) GetMLifetimeDuration() *SupplyCommissionSettingDurationObject {
	if m != nil {
		return m.MLifetimeDuration
	}
	return nil
}

func (m *SupplyCommissionSetting) GetMLevel1LimitDuration() *SupplyCommissionSettingDurationObject {
	if m != nil {
		return m.MLevel1LimitDuration
	}
	return nil
}

func (m *SupplyCommissionSetting) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

type SupplyCommissionSettingDurationObject struct {
	Duration             int32    `protobuf:"varint,1,opt,name=duration" json:"duration"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SupplyCommissionSettingDurationObject) Reset()         { *m = SupplyCommissionSettingDurationObject{} }
func (m *SupplyCommissionSettingDurationObject) String() string { return proto.CompactTextString(m) }
func (*SupplyCommissionSettingDurationObject) ProtoMessage()    {}
func (*SupplyCommissionSettingDurationObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{13}
}

var xxx_messageInfo_SupplyCommissionSettingDurationObject proto.InternalMessageInfo

func (m *SupplyCommissionSettingDurationObject) GetDuration() int32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *SupplyCommissionSettingDurationObject) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type GetEtopCommissionSettingByProductIDsRequest struct {
	ProductIds           []dot.ID `protobuf:"varint,1,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetEtopCommissionSettingByProductIDsRequest) Reset() {
	*m = GetEtopCommissionSettingByProductIDsRequest{}
}
func (m *GetEtopCommissionSettingByProductIDsRequest) String() string {
	return proto.CompactTextString(m)
}
func (*GetEtopCommissionSettingByProductIDsRequest) ProtoMessage() {}
func (*GetEtopCommissionSettingByProductIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{14}
}

var xxx_messageInfo_GetEtopCommissionSettingByProductIDsRequest proto.InternalMessageInfo

func (m *GetEtopCommissionSettingByProductIDsRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type GetEtopCommissionSettingByProductIDsResponse struct {
	EtopCommissionSettings []*CommissionSetting `protobuf:"bytes,1,rep,name=etop_commission_settings,json=etopCommissionSettings" json:"etop_commission_settings"`
	XXX_NoUnkeyedLiteral   struct{}             `json:"-"`
	XXX_sizecache          int32                `json:"-"`
}

func (m *GetEtopCommissionSettingByProductIDsResponse) Reset() {
	*m = GetEtopCommissionSettingByProductIDsResponse{}
}
func (m *GetEtopCommissionSettingByProductIDsResponse) String() string {
	return proto.CompactTextString(m)
}
func (*GetEtopCommissionSettingByProductIDsResponse) ProtoMessage() {}
func (*GetEtopCommissionSettingByProductIDsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{15}
}

var xxx_messageInfo_GetEtopCommissionSettingByProductIDsResponse proto.InternalMessageInfo

func (m *GetEtopCommissionSettingByProductIDsResponse) GetEtopCommissionSettings() []*CommissionSetting {
	if m != nil {
		return m.EtopCommissionSettings
	}
	return nil
}

type GetCommissionSettingByProductIDsRequest struct {
	ProductIds           []dot.ID `protobuf:"varint,1,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCommissionSettingByProductIDsRequest) Reset() {
	*m = GetCommissionSettingByProductIDsRequest{}
}
func (m *GetCommissionSettingByProductIDsRequest) String() string { return proto.CompactTextString(m) }
func (*GetCommissionSettingByProductIDsRequest) ProtoMessage()    {}
func (*GetCommissionSettingByProductIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{16}
}

var xxx_messageInfo_GetCommissionSettingByProductIDsRequest proto.InternalMessageInfo

func (m *GetCommissionSettingByProductIDsRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type GetCommissionSettingByProductIDsResponse struct {
	CommissionSettings   []*CommissionSetting `protobuf:"bytes,1,rep,name=commission_settings,json=commissionSettings" json:"commission_settings"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GetCommissionSettingByProductIDsResponse) Reset() {
	*m = GetCommissionSettingByProductIDsResponse{}
}
func (m *GetCommissionSettingByProductIDsResponse) String() string { return proto.CompactTextString(m) }
func (*GetCommissionSettingByProductIDsResponse) ProtoMessage()    {}
func (*GetCommissionSettingByProductIDsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{17}
}

var xxx_messageInfo_GetCommissionSettingByProductIDsResponse proto.InternalMessageInfo

func (m *GetCommissionSettingByProductIDsResponse) GetCommissionSettings() []*CommissionSetting {
	if m != nil {
		return m.CommissionSettings
	}
	return nil
}

type CreateOrUpdateCommissionSettingRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	Amount               int32    `protobuf:"varint,2,opt,name=amount" json:"amount"`
	Unit                 *string  `protobuf:"bytes,3,opt,name=unit" json:"unit"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrUpdateCommissionSettingRequest) Reset() {
	*m = CreateOrUpdateCommissionSettingRequest{}
}
func (m *CreateOrUpdateCommissionSettingRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrUpdateCommissionSettingRequest) ProtoMessage()    {}
func (*CreateOrUpdateCommissionSettingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{18}
}

var xxx_messageInfo_CreateOrUpdateCommissionSettingRequest proto.InternalMessageInfo

func (m *CreateOrUpdateCommissionSettingRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CreateOrUpdateCommissionSettingRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *CreateOrUpdateCommissionSettingRequest) GetUnit() string {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return ""
}

type CreateOrUpdateTradingCommissionSettingRequest struct {
	ProductId                dot.ID `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	Level1DirectCommission   int32  `protobuf:"varint,2,opt,name=level1_direct_commission,json=level1DirectCommission" json:"level1_direct_commission"`
	Level1IndirectCommission int32  `protobuf:"varint,3,opt,name=level1_indirect_commission,json=level1IndirectCommission" json:"level1_indirect_commission"`
	Level2DirectCommission   int32  `protobuf:"varint,4,opt,name=level2_direct_commission,json=level2DirectCommission" json:"level2_direct_commission"`
	Level2IndirectCommission int32  `protobuf:"varint,5,opt,name=level2_indirect_commission,json=level2IndirectCommission" json:"level2_indirect_commission"`
	// product, customer
	DependOn         string `protobuf:"bytes,6,opt,name=depend_on,json=dependOn" json:"depend_on"`
	Level1LimitCount int32  `protobuf:"varint,7,opt,name=level1_limit_count,json=level1LimitCount" json:"level1_limit_count"`
	// day, month
	Level1LimitDurationType string `protobuf:"bytes,9,opt,name=level1_limit_duration_type,json=level1LimitDurationType" json:"level1_limit_duration_type"`
	Level1LimitDuration     int32  `protobuf:"varint,10,opt,name=level1_limit_duration,json=level1LimitDuration" json:"level1_limit_duration"`
	// day, month
	LifetimeDurationType string   `protobuf:"bytes,11,opt,name=lifetime_duration_type,json=lifetimeDurationType" json:"lifetime_duration_type"`
	LifetimeDuration     int32    `protobuf:"varint,12,opt,name=lifetime_duration,json=lifetimeDuration" json:"lifetime_duration"`
	Group                string   `protobuf:"bytes,13,opt,name=group" json:"group"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) Reset() {
	*m = CreateOrUpdateTradingCommissionSettingRequest{}
}
func (m *CreateOrUpdateTradingCommissionSettingRequest) String() string {
	return proto.CompactTextString(m)
}
func (*CreateOrUpdateTradingCommissionSettingRequest) ProtoMessage() {}
func (*CreateOrUpdateTradingCommissionSettingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{19}
}

var xxx_messageInfo_CreateOrUpdateTradingCommissionSettingRequest proto.InternalMessageInfo

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel1DirectCommission() int32 {
	if m != nil {
		return m.Level1DirectCommission
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel1IndirectCommission() int32 {
	if m != nil {
		return m.Level1IndirectCommission
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel2DirectCommission() int32 {
	if m != nil {
		return m.Level2DirectCommission
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel2IndirectCommission() int32 {
	if m != nil {
		return m.Level2IndirectCommission
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetDependOn() string {
	if m != nil {
		return m.DependOn
	}
	return ""
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel1LimitCount() int32 {
	if m != nil {
		return m.Level1LimitCount
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel1LimitDurationType() string {
	if m != nil {
		return m.Level1LimitDurationType
	}
	return ""
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLevel1LimitDuration() int32 {
	if m != nil {
		return m.Level1LimitDuration
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLifetimeDurationType() string {
	if m != nil {
		return m.LifetimeDurationType
	}
	return ""
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetLifetimeDuration() int32 {
	if m != nil {
		return m.LifetimeDuration
	}
	return 0
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

type ProductPromotion struct {
	Product              *shop.ShopProduct `protobuf:"bytes,1,opt,name=product" json:"product"`
	Id                   dot.ID            `protobuf:"varint,2,opt,name=id" json:"id"`
	ProductId            dot.ID            `protobuf:"varint,3,opt,name=product_id,json=productId" json:"product_id"`
	Amount               int32             `protobuf:"varint,4,opt,name=amount" json:"amount"`
	Unit                 string            `protobuf:"bytes,5,opt,name=unit" json:"unit"`
	Type                 string            `protobuf:"bytes,6,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ProductPromotion) Reset()         { *m = ProductPromotion{} }
func (m *ProductPromotion) String() string { return proto.CompactTextString(m) }
func (*ProductPromotion) ProtoMessage()    {}
func (*ProductPromotion) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{20}
}

var xxx_messageInfo_ProductPromotion proto.InternalMessageInfo

func (m *ProductPromotion) GetProduct() *shop.ShopProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *ProductPromotion) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProductPromotion) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *ProductPromotion) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *ProductPromotion) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *ProductPromotion) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type CreateOrUpdateProductPromotionRequest struct {
	Id                   dot.ID   `protobuf:"varint,8,opt,name=id" json:"id"`
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	Amount               int32    `protobuf:"varint,2,opt,name=amount" json:"amount"`
	Unit                 string   `protobuf:"bytes,3,opt,name=unit" json:"unit"`
	Code                 string   `protobuf:"bytes,4,opt,name=code" json:"code"`
	Description          string   `protobuf:"bytes,5,opt,name=description" json:"description"`
	Note                 string   `protobuf:"bytes,6,opt,name=note" json:"note"`
	Type                 string   `protobuf:"bytes,7,opt,name=type" json:"type"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateOrUpdateProductPromotionRequest) Reset()         { *m = CreateOrUpdateProductPromotionRequest{} }
func (m *CreateOrUpdateProductPromotionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateOrUpdateProductPromotionRequest) ProtoMessage()    {}
func (*CreateOrUpdateProductPromotionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{21}
}

var xxx_messageInfo_CreateOrUpdateProductPromotionRequest proto.InternalMessageInfo

func (m *CreateOrUpdateProductPromotionRequest) GetId() dot.ID {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CreateOrUpdateProductPromotionRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CreateOrUpdateProductPromotionRequest) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *CreateOrUpdateProductPromotionRequest) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *CreateOrUpdateProductPromotionRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *CreateOrUpdateProductPromotionRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateOrUpdateProductPromotionRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

func (m *CreateOrUpdateProductPromotionRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type GetProductPromotionByProductIDRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProductPromotionByProductIDRequest) Reset()         { *m = GetProductPromotionByProductIDRequest{} }
func (m *GetProductPromotionByProductIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetProductPromotionByProductIDRequest) ProtoMessage()    {}
func (*GetProductPromotionByProductIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{22}
}

var xxx_messageInfo_GetProductPromotionByProductIDRequest proto.InternalMessageInfo

type GetProductPromotionByProductIDResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProductPromotionByProductIDResponse) Reset() {
	*m = GetProductPromotionByProductIDResponse{}
}
func (m *GetProductPromotionByProductIDResponse) String() string { return proto.CompactTextString(m) }
func (*GetProductPromotionByProductIDResponse) ProtoMessage()    {}
func (*GetProductPromotionByProductIDResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{23}
}

var xxx_messageInfo_GetProductPromotionByProductIDResponse proto.InternalMessageInfo

type SupplyProductResponse struct {
	Product                 *shop.ShopProduct        `protobuf:"bytes,1,opt,name=product" json:"product"`
	SupplyCommissionSetting *SupplyCommissionSetting `protobuf:"bytes,2,opt,name=supply_commission_setting,json=supplyCommissionSetting" json:"supply_commission_setting"`
	Promotion               *ProductPromotion        `protobuf:"bytes,3,opt,name=promotion" json:"promotion"`
	XXX_NoUnkeyedLiteral    struct{}                 `json:"-"`
	XXX_sizecache           int32                    `json:"-"`
}

func (m *SupplyProductResponse) Reset()         { *m = SupplyProductResponse{} }
func (m *SupplyProductResponse) String() string { return proto.CompactTextString(m) }
func (*SupplyProductResponse) ProtoMessage()    {}
func (*SupplyProductResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{24}
}

var xxx_messageInfo_SupplyProductResponse proto.InternalMessageInfo

func (m *SupplyProductResponse) GetProduct() *shop.ShopProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *SupplyProductResponse) GetSupplyCommissionSetting() *SupplyCommissionSetting {
	if m != nil {
		return m.SupplyCommissionSetting
	}
	return nil
}

func (m *SupplyProductResponse) GetPromotion() *ProductPromotion {
	if m != nil {
		return m.Promotion
	}
	return nil
}

type ShopProductResponse struct {
	Product              *shop.ShopProduct `protobuf:"bytes,1,opt,name=product" json:"product"`
	Promotion            *ProductPromotion `protobuf:"bytes,3,opt,name=promotion" json:"promotion"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ShopProductResponse) Reset()         { *m = ShopProductResponse{} }
func (m *ShopProductResponse) String() string { return proto.CompactTextString(m) }
func (*ShopProductResponse) ProtoMessage()    {}
func (*ShopProductResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{25}
}

var xxx_messageInfo_ShopProductResponse proto.InternalMessageInfo

func (m *ShopProductResponse) GetProduct() *shop.ShopProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *ShopProductResponse) GetPromotion() *ProductPromotion {
	if m != nil {
		return m.Promotion
	}
	return nil
}

type AffiliateProductResponse struct {
	Product                    *shop.ShopProduct  `protobuf:"bytes,1,opt,name=product" json:"product"`
	ShopCommissionSetting      *CommissionSetting `protobuf:"bytes,2,opt,name=shop_commission_setting,json=shopCommissionSetting" json:"shop_commission_setting"`
	AffiliateCommissionSetting *CommissionSetting `protobuf:"bytes,3,opt,name=affiliate_commission_setting,json=affiliateCommissionSetting" json:"affiliate_commission_setting"`
	Promotion                  *ProductPromotion  `protobuf:"bytes,4,opt,name=promotion" json:"promotion"`
	XXX_NoUnkeyedLiteral       struct{}           `json:"-"`
	XXX_sizecache              int32              `json:"-"`
}

func (m *AffiliateProductResponse) Reset()         { *m = AffiliateProductResponse{} }
func (m *AffiliateProductResponse) String() string { return proto.CompactTextString(m) }
func (*AffiliateProductResponse) ProtoMessage()    {}
func (*AffiliateProductResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{26}
}

var xxx_messageInfo_AffiliateProductResponse proto.InternalMessageInfo

func (m *AffiliateProductResponse) GetProduct() *shop.ShopProduct {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *AffiliateProductResponse) GetShopCommissionSetting() *CommissionSetting {
	if m != nil {
		return m.ShopCommissionSetting
	}
	return nil
}

func (m *AffiliateProductResponse) GetAffiliateCommissionSetting() *CommissionSetting {
	if m != nil {
		return m.AffiliateCommissionSetting
	}
	return nil
}

func (m *AffiliateProductResponse) GetPromotion() *ProductPromotion {
	if m != nil {
		return m.Promotion
	}
	return nil
}

type SupplyGetProductsResponse struct {
	Paging               *common.PageInfo         `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Products             []*SupplyProductResponse `protobuf:"bytes,2,rep,name=products" json:"products"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *SupplyGetProductsResponse) Reset()         { *m = SupplyGetProductsResponse{} }
func (m *SupplyGetProductsResponse) String() string { return proto.CompactTextString(m) }
func (*SupplyGetProductsResponse) ProtoMessage()    {}
func (*SupplyGetProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{27}
}

var xxx_messageInfo_SupplyGetProductsResponse proto.InternalMessageInfo

func (m *SupplyGetProductsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *SupplyGetProductsResponse) GetProducts() []*SupplyProductResponse {
	if m != nil {
		return m.Products
	}
	return nil
}

type ShopGetProductsResponse struct {
	Paging               *common.PageInfo       `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Products             []*ShopProductResponse `protobuf:"bytes,2,rep,name=products" json:"products"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ShopGetProductsResponse) Reset()         { *m = ShopGetProductsResponse{} }
func (m *ShopGetProductsResponse) String() string { return proto.CompactTextString(m) }
func (*ShopGetProductsResponse) ProtoMessage()    {}
func (*ShopGetProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{28}
}

var xxx_messageInfo_ShopGetProductsResponse proto.InternalMessageInfo

func (m *ShopGetProductsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *ShopGetProductsResponse) GetProducts() []*ShopProductResponse {
	if m != nil {
		return m.Products
	}
	return nil
}

type CheckReferralCodeValidRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	ReferralCode         string   `protobuf:"bytes,2,opt,name=referral_code,json=referralCode" json:"referral_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckReferralCodeValidRequest) Reset()         { *m = CheckReferralCodeValidRequest{} }
func (m *CheckReferralCodeValidRequest) String() string { return proto.CompactTextString(m) }
func (*CheckReferralCodeValidRequest) ProtoMessage()    {}
func (*CheckReferralCodeValidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{29}
}

var xxx_messageInfo_CheckReferralCodeValidRequest proto.InternalMessageInfo

func (m *CheckReferralCodeValidRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *CheckReferralCodeValidRequest) GetReferralCode() string {
	if m != nil {
		return m.ReferralCode
	}
	return ""
}

type AffiliateGetProductsResponse struct {
	Paging               *common.PageInfo            `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Products             []*AffiliateProductResponse `protobuf:"bytes,2,rep,name=products" json:"products"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *AffiliateGetProductsResponse) Reset()         { *m = AffiliateGetProductsResponse{} }
func (m *AffiliateGetProductsResponse) String() string { return proto.CompactTextString(m) }
func (*AffiliateGetProductsResponse) ProtoMessage()    {}
func (*AffiliateGetProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{30}
}

var xxx_messageInfo_AffiliateGetProductsResponse proto.InternalMessageInfo

func (m *AffiliateGetProductsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *AffiliateGetProductsResponse) GetProducts() []*AffiliateProductResponse {
	if m != nil {
		return m.Products
	}
	return nil
}

type GetProductPromotionRequest struct {
	ProductId            dot.ID   `protobuf:"varint,1,opt,name=product_id,json=productId" json:"product_id"`
	ReferralCode         *string  `protobuf:"bytes,2,opt,name=referral_code,json=referralCode" json:"referral_code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProductPromotionRequest) Reset()         { *m = GetProductPromotionRequest{} }
func (m *GetProductPromotionRequest) String() string { return proto.CompactTextString(m) }
func (*GetProductPromotionRequest) ProtoMessage()    {}
func (*GetProductPromotionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{31}
}

var xxx_messageInfo_GetProductPromotionRequest proto.InternalMessageInfo

func (m *GetProductPromotionRequest) GetProductId() dot.ID {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *GetProductPromotionRequest) GetReferralCode() string {
	if m != nil && m.ReferralCode != nil {
		return *m.ReferralCode
	}
	return ""
}

type GetProductPromotionResponse struct {
	Promotion            *ProductPromotion  `protobuf:"bytes,1,opt,name=promotion" json:"promotion"`
	ReferralDiscount     *CommissionSetting `protobuf:"bytes,2,opt,name=referral_discount,json=referralDiscount" json:"referral_discount"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GetProductPromotionResponse) Reset()         { *m = GetProductPromotionResponse{} }
func (m *GetProductPromotionResponse) String() string { return proto.CompactTextString(m) }
func (*GetProductPromotionResponse) ProtoMessage()    {}
func (*GetProductPromotionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{32}
}

var xxx_messageInfo_GetProductPromotionResponse proto.InternalMessageInfo

func (m *GetProductPromotionResponse) GetPromotion() *ProductPromotion {
	if m != nil {
		return m.Promotion
	}
	return nil
}

func (m *GetProductPromotionResponse) GetReferralDiscount() *CommissionSetting {
	if m != nil {
		return m.ReferralDiscount
	}
	return nil
}

type GetProductPromotionsResponse struct {
	Paging               *common.PageInfo    `protobuf:"bytes,1,opt,name=paging" json:"paging"`
	Promotions           []*ProductPromotion `protobuf:"bytes,2,rep,name=promotions" json:"promotions"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetProductPromotionsResponse) Reset()         { *m = GetProductPromotionsResponse{} }
func (m *GetProductPromotionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetProductPromotionsResponse) ProtoMessage()    {}
func (*GetProductPromotionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{33}
}

var xxx_messageInfo_GetProductPromotionsResponse proto.InternalMessageInfo

func (m *GetProductPromotionsResponse) GetPaging() *common.PageInfo {
	if m != nil {
		return m.Paging
	}
	return nil
}

func (m *GetProductPromotionsResponse) GetPromotions() []*ProductPromotion {
	if m != nil {
		return m.Promotions
	}
	return nil
}

type GetTradingProductPromotionByIDsRequest struct {
	ProductIds           []dot.ID `protobuf:"varint,1,rep,name=product_ids,json=productIds" json:"product_ids"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTradingProductPromotionByIDsRequest) Reset() {
	*m = GetTradingProductPromotionByIDsRequest{}
}
func (m *GetTradingProductPromotionByIDsRequest) String() string { return proto.CompactTextString(m) }
func (*GetTradingProductPromotionByIDsRequest) ProtoMessage()    {}
func (*GetTradingProductPromotionByIDsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{34}
}

var xxx_messageInfo_GetTradingProductPromotionByIDsRequest proto.InternalMessageInfo

func (m *GetTradingProductPromotionByIDsRequest) GetProductIds() []dot.ID {
	if m != nil {
		return m.ProductIds
	}
	return nil
}

type GetTradingProductPromotionByIDsResponse struct {
	Promotions           []*ProductPromotion `protobuf:"bytes,1,rep,name=promotions" json:"promotions"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GetTradingProductPromotionByIDsResponse) Reset() {
	*m = GetTradingProductPromotionByIDsResponse{}
}
func (m *GetTradingProductPromotionByIDsResponse) String() string { return proto.CompactTextString(m) }
func (*GetTradingProductPromotionByIDsResponse) ProtoMessage()    {}
func (*GetTradingProductPromotionByIDsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{35}
}

var xxx_messageInfo_GetTradingProductPromotionByIDsResponse proto.InternalMessageInfo

func (m *GetTradingProductPromotionByIDsResponse) GetPromotions() []*ProductPromotion {
	if m != nil {
		return m.Promotions
	}
	return nil
}

type CreateReferralCodeRequest struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateReferralCodeRequest) Reset()         { *m = CreateReferralCodeRequest{} }
func (m *CreateReferralCodeRequest) String() string { return proto.CompactTextString(m) }
func (*CreateReferralCodeRequest) ProtoMessage()    {}
func (*CreateReferralCodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{36}
}

var xxx_messageInfo_CreateReferralCodeRequest proto.InternalMessageInfo

func (m *CreateReferralCodeRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type ReferralCode struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code" json:"code"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReferralCode) Reset()         { *m = ReferralCode{} }
func (m *ReferralCode) String() string { return proto.CompactTextString(m) }
func (*ReferralCode) ProtoMessage()    {}
func (*ReferralCode) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{37}
}

var xxx_messageInfo_ReferralCode proto.InternalMessageInfo

func (m *ReferralCode) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type GetReferralCodesResponse struct {
	ReferralCodes        []*ReferralCode `protobuf:"bytes,1,rep,name=referral_codes,json=referralCodes" json:"referral_codes"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetReferralCodesResponse) Reset()         { *m = GetReferralCodesResponse{} }
func (m *GetReferralCodesResponse) String() string { return proto.CompactTextString(m) }
func (*GetReferralCodesResponse) ProtoMessage()    {}
func (*GetReferralCodesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{38}
}

var xxx_messageInfo_GetReferralCodesResponse proto.InternalMessageInfo

func (m *GetReferralCodesResponse) GetReferralCodes() []*ReferralCode {
	if m != nil {
		return m.ReferralCodes
	}
	return nil
}

type Referral struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone" json:"phone"`
	Email                string   `protobuf:"bytes,3,opt,name=email" json:"email"`
	OrderCount           int32    `protobuf:"varint,4,opt,name=order_count,json=orderCount" json:"order_count"`
	TotalRevenue         int32    `protobuf:"varint,5,opt,name=total_revenue,json=totalRevenue" json:"total_revenue"`
	TotalCommission      int32    `protobuf:"varint,6,opt,name=total_commission,json=totalCommission" json:"total_commission"`
	CreatedAt            dot.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt" json:"created_at"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Referral) Reset()         { *m = Referral{} }
func (m *Referral) String() string { return proto.CompactTextString(m) }
func (*Referral) ProtoMessage()    {}
func (*Referral) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{39}
}

var xxx_messageInfo_Referral proto.InternalMessageInfo

func (m *Referral) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Referral) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Referral) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Referral) GetOrderCount() int32 {
	if m != nil {
		return m.OrderCount
	}
	return 0
}

func (m *Referral) GetTotalRevenue() int32 {
	if m != nil {
		return m.TotalRevenue
	}
	return 0
}

func (m *Referral) GetTotalCommission() int32 {
	if m != nil {
		return m.TotalCommission
	}
	return 0
}

type GetReferralsResponse struct {
	Referrals            []*Referral `protobuf:"bytes,1,rep,name=referrals" json:"referrals"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetReferralsResponse) Reset()         { *m = GetReferralsResponse{} }
func (m *GetReferralsResponse) String() string { return proto.CompactTextString(m) }
func (*GetReferralsResponse) ProtoMessage()    {}
func (*GetReferralsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_831e3e6df584a8a9, []int{40}
}

var xxx_messageInfo_GetReferralsResponse proto.InternalMessageInfo

func (m *GetReferralsResponse) GetReferrals() []*Referral {
	if m != nil {
		return m.Referrals
	}
	return nil
}

func init() {
	proto.RegisterType((*UpdateReferralRequest)(nil), "affiliate.UpdateReferralRequest")
	proto.RegisterType((*UserReferral)(nil), "affiliate.UserReferral")
	proto.RegisterType((*SellerCommission)(nil), "affiliate.SellerCommission")
	proto.RegisterType((*GetCommissionsResponse)(nil), "affiliate.GetCommissionsResponse")
	proto.RegisterType((*NotifyNewShopPurchaseRequest)(nil), "affiliate.NotifyNewShopPurchaseRequest")
	proto.RegisterType((*NotifyNewShopPurchaseResponse)(nil), "affiliate.NotifyNewShopPurchaseResponse")
	proto.RegisterType((*GetCouponsResponse)(nil), "affiliate.GetCouponsResponse")
	proto.RegisterType((*CreateCouponRequest)(nil), "affiliate.CreateCouponRequest")
	proto.RegisterType((*Coupon)(nil), "affiliate.Coupon")
	proto.RegisterType((*GetTransactionsResponse)(nil), "affiliate.GetTransactionsResponse")
	proto.RegisterType((*Transaction)(nil), "affiliate.Transaction")
	proto.RegisterType((*CommissionSetting)(nil), "affiliate.CommissionSetting")
	proto.RegisterType((*SupplyCommissionSetting)(nil), "affiliate.SupplyCommissionSetting")
	proto.RegisterType((*SupplyCommissionSettingDurationObject)(nil), "affiliate.SupplyCommissionSettingDurationObject")
	proto.RegisterType((*GetEtopCommissionSettingByProductIDsRequest)(nil), "affiliate.GetEtopCommissionSettingByProductIDsRequest")
	proto.RegisterType((*GetEtopCommissionSettingByProductIDsResponse)(nil), "affiliate.GetEtopCommissionSettingByProductIDsResponse")
	proto.RegisterType((*GetCommissionSettingByProductIDsRequest)(nil), "affiliate.GetCommissionSettingByProductIDsRequest")
	proto.RegisterType((*GetCommissionSettingByProductIDsResponse)(nil), "affiliate.GetCommissionSettingByProductIDsResponse")
	proto.RegisterType((*CreateOrUpdateCommissionSettingRequest)(nil), "affiliate.CreateOrUpdateCommissionSettingRequest")
	proto.RegisterType((*CreateOrUpdateTradingCommissionSettingRequest)(nil), "affiliate.CreateOrUpdateTradingCommissionSettingRequest")
	proto.RegisterType((*ProductPromotion)(nil), "affiliate.ProductPromotion")
	proto.RegisterType((*CreateOrUpdateProductPromotionRequest)(nil), "affiliate.CreateOrUpdateProductPromotionRequest")
	proto.RegisterType((*GetProductPromotionByProductIDRequest)(nil), "affiliate.GetProductPromotionByProductIDRequest")
	proto.RegisterType((*GetProductPromotionByProductIDResponse)(nil), "affiliate.GetProductPromotionByProductIDResponse")
	proto.RegisterType((*SupplyProductResponse)(nil), "affiliate.SupplyProductResponse")
	proto.RegisterType((*ShopProductResponse)(nil), "affiliate.ShopProductResponse")
	proto.RegisterType((*AffiliateProductResponse)(nil), "affiliate.AffiliateProductResponse")
	proto.RegisterType((*SupplyGetProductsResponse)(nil), "affiliate.SupplyGetProductsResponse")
	proto.RegisterType((*ShopGetProductsResponse)(nil), "affiliate.ShopGetProductsResponse")
	proto.RegisterType((*CheckReferralCodeValidRequest)(nil), "affiliate.CheckReferralCodeValidRequest")
	proto.RegisterType((*AffiliateGetProductsResponse)(nil), "affiliate.AffiliateGetProductsResponse")
	proto.RegisterType((*GetProductPromotionRequest)(nil), "affiliate.GetProductPromotionRequest")
	proto.RegisterType((*GetProductPromotionResponse)(nil), "affiliate.GetProductPromotionResponse")
	proto.RegisterType((*GetProductPromotionsResponse)(nil), "affiliate.GetProductPromotionsResponse")
	proto.RegisterType((*GetTradingProductPromotionByIDsRequest)(nil), "affiliate.GetTradingProductPromotionByIDsRequest")
	proto.RegisterType((*GetTradingProductPromotionByIDsResponse)(nil), "affiliate.GetTradingProductPromotionByIDsResponse")
	proto.RegisterType((*CreateReferralCodeRequest)(nil), "affiliate.CreateReferralCodeRequest")
	proto.RegisterType((*ReferralCode)(nil), "affiliate.ReferralCode")
	proto.RegisterType((*GetReferralCodesResponse)(nil), "affiliate.GetReferralCodesResponse")
	proto.RegisterType((*Referral)(nil), "affiliate.Referral")
	proto.RegisterType((*GetReferralsResponse)(nil), "affiliate.GetReferralsResponse")
}

func init() { proto.RegisterFile("services/affiliate/affiliate.proto", fileDescriptor_831e3e6df584a8a9) }

var fileDescriptor_831e3e6df584a8a9 = []byte{
	// 1959 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x59, 0xcb, 0x6f, 0xdc, 0xc6,
	0x19, 0x5f, 0xee, 0x7b, 0xbf, 0x5d, 0xd9, 0x12, 0x25, 0x4b, 0xb4, 0x22, 0xaf, 0x15, 0xa6, 0x76,
	0x36, 0x70, 0xb3, 0x5b, 0x2f, 0x9a, 0xa2, 0x76, 0xda, 0x18, 0x7a, 0x14, 0x82, 0x0a, 0xd7, 0x32,
	0x28, 0xd9, 0x05, 0x72, 0xc8, 0x82, 0x26, 0x67, 0x57, 0x6c, 0x49, 0x0e, 0x4b, 0xce, 0xaa, 0x10,
	0x8a, 0xa2, 0xe8, 0xa1, 0x45, 0x0f, 0x3d, 0x04, 0x3d, 0xe5, 0x58, 0xa0, 0x97, 0xfe, 0x09, 0x01,
	0x7a, 0xeb, 0xc9, 0x87, 0x1e, 0xf2, 0x17, 0xb4, 0xb1, 0x02, 0xf4, 0xd0, 0x53, 0xff, 0x84, 0x82,
	0x33, 0x43, 0x72, 0xf8, 0x58, 0xad, 0x1e, 0xc9, 0xad, 0x97, 0x5d, 0xe9, 0x7b, 0xcc, 0xfc, 0xbe,
	0x6f, 0xbe, 0xd7, 0xcc, 0x82, 0x1a, 0x20, 0xff, 0xc4, 0x32, 0x50, 0x30, 0xd0, 0xc7, 0x63, 0xcb,
	0xb6, 0x74, 0x82, 0x92, 0xbf, 0xfa, 0x9e, 0x8f, 0x09, 0x96, 0x5b, 0x31, 0x61, 0x7d, 0x65, 0x82,
	0x27, 0x98, 0x52, 0x07, 0xe1, 0x5f, 0x4c, 0x60, 0x7d, 0xd9, 0xc0, 0x8e, 0x83, 0xdd, 0x01, 0xfb,
	0xe2, 0xc4, 0xbb, 0x13, 0x8c, 0x27, 0x36, 0x1a, 0xd0, 0xff, 0x5e, 0x4d, 0xc7, 0x03, 0x62, 0x39,
	0x28, 0x20, 0xba, 0xe3, 0x71, 0x81, 0x15, 0x44, 0xb0, 0x37, 0x08, 0x8e, 0xf9, 0x07, 0xa7, 0xae,
	0x52, 0x2a, 0xf6, 0x4d, 0xe4, 0xb3, 0x4f, 0x4e, 0xbf, 0x49, 0xe9, 0xe1, 0x07, 0x27, 0xdc, 0xe1,
	0x04, 0x63, 0x10, 0x10, 0x9d, 0x4c, 0x83, 0xef, 0xf2, 0x6f, 0xc6, 0x56, 0x4f, 0xe0, 0xd6, 0x0b,
	0xcf, 0xd4, 0x09, 0xd2, 0xd0, 0x18, 0xf9, 0xbe, 0x6e, 0x6b, 0xe8, 0x17, 0x53, 0x14, 0x10, 0xf9,
	0x3d, 0x58, 0xf0, 0x39, 0x69, 0x64, 0x60, 0x13, 0x29, 0xe5, 0x4d, 0xa9, 0xd7, 0xda, 0xae, 0xbe,
	0xfe, 0xe7, 0xdd, 0x92, 0xd6, 0x89, 0x58, 0x3b, 0xd8, 0x44, 0xf2, 0x10, 0xe4, 0x40, 0xb7, 0xd1,
	0x28, 0x2d, 0x5f, 0x11, 0xe4, 0x17, 0x43, 0xbe, 0x26, 0xe8, 0xa8, 0x7f, 0x94, 0xa0, 0xf3, 0x22,
	0x40, 0x7e, 0x44, 0x94, 0xef, 0x40, 0x63, 0x1a, 0x20, 0x7f, 0x64, 0x99, 0x8a, 0xb4, 0x29, 0xf5,
	0x2a, 0x5c, 0xb3, 0x1e, 0x12, 0xf7, 0xcd, 0x6f, 0x1a, 0xce, 0x7f, 0xaa, 0xb0, 0x78, 0x88, 0x6c,
	0x1b, 0xf9, 0x3b, 0xd8, 0x71, 0xac, 0x20, 0xb0, 0xb0, 0x2b, 0xaf, 0x40, 0x39, 0x83, 0xa6, 0x6c,
	0x99, 0xf2, 0x3a, 0xd4, 0x4e, 0x74, 0x7b, 0x8a, 0x94, 0xea, 0xa6, 0xd4, 0xab, 0x71, 0x06, 0x23,
	0xc9, 0xf7, 0xa1, 0x6d, 0xa2, 0xc0, 0xf0, 0x2d, 0x8f, 0x58, 0xd8, 0x55, 0x6a, 0xc2, 0x9e, 0x22,
	0x43, 0x56, 0xa0, 0xea, 0x62, 0x82, 0x94, 0xba, 0x20, 0x40, 0x29, 0xf2, 0xfb, 0x50, 0x67, 0xe7,
	0xa3, 0x34, 0x37, 0xa5, 0xde, 0x8d, 0xe1, 0xcd, 0x3e, 0x3f, 0xb6, 0xfe, 0x21, 0xfd, 0x8e, 0xdc,
	0xc2, 0xa8, 0xe1, 0x42, 0xe4, 0xd4, 0x43, 0x4a, 0x4b, 0x5c, 0x28, 0xa4, 0x84, 0xfe, 0xc4, 0x23,
	0x06, 0x74, 0x51, 0x00, 0x5a, 0xc7, 0x2f, 0x39, 0xd2, 0x0e, 0x1e, 0xbd, 0xd2, 0x03, 0xc4, 0x65,
	0x96, 0x04, 0x19, 0xc0, 0xdb, 0x7a, 0x80, 0x98, 0xdc, 0x43, 0x68, 0x78, 0x3e, 0x36, 0xa7, 0x06,
	0x51, 0x16, 0x36, 0xa5, 0x5e, 0x7b, 0xb8, 0xd4, 0xa7, 0x51, 0x78, 0x78, 0x8c, 0xbd, 0xe7, 0x8c,
	0x41, 0xb5, 0x24, 0x2d, 0x92, 0x93, 0x7b, 0x50, 0xa3, 0x11, 0xa9, 0xdc, 0xa0, 0x0a, 0x9d, 0x3e,
	0x8b, 0xcf, 0x83, 0xf0, 0x93, 0xcb, 0x32, 0x01, 0xf9, 0x7b, 0xd0, 0x1e, 0xfb, 0xd8, 0x19, 0x05,
	0xd4, 0xf3, 0xca, 0x4d, 0x2a, 0x7f, 0xb3, 0x4f, 0xa3, 0x77, 0x2b, 0x4a, 0x26, 0xae, 0x02, 0xa1,
	0x24, 0x3b, 0x22, 0xf9, 0x43, 0x68, 0x9e, 0xe8, 0xb6, 0x65, 0x8e, 0x74, 0xa2, 0x00, 0x55, 0x5a,
	0xef, 0xb3, 0x34, 0xea, 0x47, 0x69, 0xd4, 0x3f, 0x8a, 0xd2, 0x28, 0x82, 0x47, 0x35, 0xb6, 0x88,
	0xfc, 0x04, 0xc0, 0xf0, 0x91, 0x4e, 0x10, 0x55, 0x6f, 0x5f, 0x50, 0xbd, 0xc5, 0x75, 0xd8, 0x02,
	0x53, 0x9a, 0x32, 0x74, 0x81, 0xce, 0x45, 0x17, 0xe0, 0x3a, 0x5b, 0x44, 0xfd, 0x29, 0xac, 0xee,
	0x21, 0x92, 0x04, 0x5a, 0xa0, 0xa1, 0xc0, 0xc3, 0x6e, 0x80, 0xe4, 0x1f, 0x42, 0xdb, 0x48, 0xc8,
	0x8a, 0xb4, 0x59, 0xe9, 0xb5, 0x87, 0x6f, 0xf5, 0x93, 0x4a, 0x93, 0x8d, 0x51, 0x4d, 0x94, 0x57,
	0x9f, 0xc0, 0xc6, 0x33, 0x4c, 0xac, 0xf1, 0xe9, 0x33, 0xf4, 0x4b, 0x7a, 0x40, 0x53, 0xdf, 0x38,
	0xd6, 0x03, 0x14, 0xe5, 0xf4, 0x5d, 0x68, 0x52, 0xc7, 0x67, 0x93, 0xac, 0x41, 0xa9, 0xfb, 0xa6,
	0xfa, 0x04, 0xee, 0xcc, 0x58, 0x80, 0x03, 0xec, 0x42, 0xc3, 0x41, 0x41, 0xa0, 0x4f, 0x10, 0x5d,
	0x20, 0x0a, 0xb9, 0x88, 0xa8, 0x6e, 0x81, 0x4c, 0x4d, 0x9b, 0x7a, 0xa2, 0x59, 0x0f, 0xa0, 0x61,
	0x30, 0x12, 0x37, 0x69, 0x49, 0x30, 0x89, 0x09, 0x6b, 0x91, 0x84, 0xfa, 0x99, 0x04, 0xcb, 0x3b,
	0xd4, 0xd9, 0x9c, 0xc3, 0xc1, 0xc7, 0x79, 0x27, 0xe5, 0xf3, 0x4e, 0x81, 0xea, 0xd4, 0xb5, 0x88,
	0x50, 0x14, 0x24, 0x8d, 0x52, 0xb2, 0x19, 0x59, 0x11, 0x04, 0x52, 0x19, 0xf9, 0x0e, 0x00, 0x8f,
	0xdf, 0xd0, 0x39, 0x55, 0xc1, 0x39, 0x2d, 0x4e, 0xdf, 0x37, 0xd5, 0xbf, 0x57, 0xa0, 0xce, 0x40,
	0xcd, 0xa8, 0x0d, 0x0a, 0x54, 0x73, 0xc5, 0x89, 0x52, 0x12, 0xf4, 0x95, 0xd9, 0xe8, 0xab, 0xa2,
	0x56, 0x11, 0xfa, 0xda, 0x2c, 0xf4, 0x42, 0xf1, 0xac, 0x17, 0x14, 0xcf, 0x47, 0x00, 0x01, 0xd1,
	0x7d, 0x32, 0x0a, 0x23, 0x50, 0x69, 0xcc, 0x8b, 0x58, 0xad, 0x45, 0xa5, 0x77, 0x75, 0x82, 0xe4,
	0x0f, 0xa0, 0x89, 0x5c, 0x93, 0x29, 0x36, 0xe7, 0x2a, 0x36, 0x90, 0x6b, 0x52, 0xb5, 0xb4, 0x3b,
	0x5b, 0x85, 0xee, 0x0c, 0x61, 0x09, 0x99, 0x08, 0xf3, 0x61, 0x25, 0x39, 0xf8, 0x28, 0x95, 0x83,
	0xed, 0xf9, 0xaa, 0x49, 0xf6, 0xbd, 0x80, 0xb5, 0x3d, 0x44, 0x8e, 0x7c, 0xdd, 0x0d, 0x74, 0x83,
	0xa4, 0xd2, 0xef, 0x31, 0x74, 0x88, 0x40, 0xe7, 0xc1, 0xba, 0x2a, 0x04, 0xab, 0xa0, 0xa6, 0xa5,
	0x64, 0xd5, 0x05, 0x68, 0x0b, 0x4c, 0xf5, 0xdf, 0x12, 0x2c, 0x25, 0x69, 0x7a, 0x88, 0x08, 0xb1,
	0xdc, 0x49, 0xc6, 0x2d, 0xe5, 0x62, 0xb7, 0x6c, 0x40, 0x5d, 0x77, 0xf0, 0xd4, 0x25, 0xa9, 0x58,
	0xe1, 0xb4, 0x73, 0x82, 0x25, 0xed, 0xce, 0xda, 0xd5, 0xdd, 0x59, 0xbf, 0x8c, 0x3b, 0xff, 0x54,
	0x87, 0xb5, 0xc3, 0xa9, 0xe7, 0xd9, 0xa7, 0xf3, 0xcc, 0x95, 0x8a, 0xcd, 0xfd, 0x08, 0x14, 0x1b,
	0x9d, 0x20, 0xfb, 0xe1, 0xc8, 0xb4, 0x7c, 0x64, 0x90, 0x51, 0x52, 0xd1, 0xa8, 0x87, 0x22, 0x07,
	0xac, 0x32, 0xa9, 0x5d, 0x2a, 0x24, 0x74, 0xe9, 0x6d, 0x58, 0xe7, 0xfa, 0x96, 0x9b, 0x5f, 0x41,
	0x74, 0x21, 0xdf, 0x67, 0x9f, 0x8b, 0x09, 0x6b, 0x44, 0x18, 0x86, 0x05, 0x18, 0xaa, 0x39, 0x0c,
	0xc3, 0x99, 0x18, 0x86, 0x85, 0x18, 0x6a, 0x39, 0x0c, 0xc3, 0x02, 0x0c, 0x6f, 0x43, 0xcb, 0x44,
	0x5e, 0x98, 0x6c, 0xd8, 0x4d, 0x0d, 0x06, 0x4d, 0x46, 0x3e, 0x70, 0xc3, 0xc9, 0x86, 0x9b, 0x6a,
	0x5b, 0x8e, 0x15, 0xee, 0x11, 0x46, 0x49, 0x43, 0x58, 0x7e, 0x91, 0xf1, 0x9f, 0x86, 0xec, 0x1d,
	0x1a, 0x2f, 0x63, 0x58, 0x76, 0x46, 0xb6, 0x35, 0x46, 0xe1, 0x5c, 0x39, 0x32, 0xa7, 0xbe, 0x4e,
	0x4b, 0x09, 0xcb, 0xe5, 0xef, 0x88, 0xad, 0xa5, 0xf8, 0x10, 0x77, 0xb9, 0xc6, 0xc1, 0xab, 0x9f,
	0xa1, 0xb8, 0xd7, 0x2f, 0x39, 0x4f, 0xf9, 0x8a, 0x11, 0x5b, 0x76, 0x60, 0xcd, 0x19, 0xa5, 0xd0,
	0xc5, 0x7b, 0xb5, 0xaf, 0xb5, 0xd7, 0x8a, 0xf3, 0x34, 0xb1, 0x29, 0xde, 0x2e, 0xdd, 0xc5, 0x5b,
	0xd7, 0xed, 0xe2, 0x70, 0xe9, 0x2e, 0x1e, 0x56, 0xf4, 0x89, 0x8f, 0xa7, 0x1e, 0x9d, 0x00, 0xa2,
	0xb3, 0x62, 0x24, 0xd5, 0x80, 0x7b, 0x17, 0x32, 0x51, 0xde, 0x84, 0x66, 0xec, 0x26, 0xb1, 0xaf,
	0xc5, 0xd4, 0x78, 0xc2, 0x2b, 0x67, 0x27, 0x3c, 0xf5, 0x19, 0x3c, 0xd8, 0x43, 0xe4, 0x47, 0x04,
	0x7b, 0xb9, 0x5d, 0xb6, 0x4f, 0xf9, 0x6c, 0xb6, 0xbf, 0x1b, 0x24, 0xcd, 0xbf, 0x9d, 0x24, 0x23,
	0xab, 0x6d, 0x15, 0x0d, 0xe2, 0x3c, 0x0c, 0xd4, 0xdf, 0x4b, 0xf0, 0xed, 0x8b, 0x2d, 0xc8, 0xcb,
	0xe5, 0x4b, 0x50, 0xc2, 0x51, 0x4d, 0x08, 0xf5, 0x51, 0xc0, 0xc4, 0xa3, 0xd2, 0xb9, 0x91, 0xea,
	0xf3, 0x99, 0x35, 0x35, 0x7a, 0x87, 0xc9, 0x91, 0x03, 0xf5, 0xc7, 0xf0, 0x6e, 0x6a, 0x3e, 0xba,
	0x8e, 0x51, 0xa7, 0xd0, 0x9b, 0xbf, 0x16, 0xb7, 0xe7, 0x27, 0xb0, 0x7c, 0x55, 0x53, 0x64, 0x23,
	0x6f, 0xc6, 0xef, 0x24, 0xb8, 0xcf, 0x06, 0x99, 0x03, 0x9f, 0xdd, 0xb1, 0xf2, 0x7a, 0xdc, 0x8c,
	0x0b, 0x15, 0xca, 0xa4, 0x2f, 0x94, 0xcf, 0xe9, 0x0b, 0x95, 0xec, 0x08, 0xa4, 0xfe, 0xad, 0x06,
	0xef, 0xa7, 0x71, 0x1c, 0xf9, 0xba, 0x69, 0xb9, 0x93, 0xeb, 0xc1, 0xf9, 0x7f, 0xdd, 0xfe, 0x46,
	0xeb, 0xf6, 0x56, 0xec, 0x9e, 0x74, 0x35, 0x1d, 0xe5, 0xee, 0x7b, 0x6b, 0x76, 0xbe, 0x3e, 0x1e,
	0x85, 0x57, 0xc0, 0xef, 0xc3, 0xad, 0xe2, 0x82, 0x0c, 0xc2, 0xce, 0xcb, 0x05, 0xda, 0xf2, 0x63,
	0x58, 0xcd, 0xb5, 0x0c, 0xb6, 0x71, 0x5b, 0xd8, 0x78, 0xc5, 0xce, 0x34, 0x01, 0xba, 0xeb, 0x43,
	0x58, 0xca, 0xb7, 0x9b, 0x4e, 0xca, 0xd6, 0x6c, 0xef, 0x88, 0x4b, 0xe9, 0x42, 0xbe, 0x94, 0xfe,
	0x43, 0x82, 0x45, 0x9e, 0xab, 0xcf, 0x7d, 0xec, 0x60, 0xaa, 0xf0, 0x20, 0xb9, 0x95, 0x4a, 0x33,
	0x6e, 0xa5, 0xc9, 0x7d, 0x94, 0x8d, 0xea, 0xe5, 0xcc, 0xa8, 0x9e, 0x8e, 0xf1, 0xca, 0xbc, 0x94,
	0xab, 0x9e, 0x93, 0x72, 0xb5, 0xdc, 0x28, 0x16, 0x15, 0xed, 0x7a, 0xae, 0x68, 0x7f, 0x5a, 0x86,
	0x7b, 0xe9, 0x64, 0xcc, 0x1a, 0x17, 0x25, 0x21, 0x83, 0xdd, 0x3c, 0x17, 0xf6, 0xd7, 0x54, 0x29,
	0x04, 0xd8, 0xf4, 0xfa, 0x52, 0xcd, 0x5d, 0x5f, 0xae, 0xff, 0xb0, 0x11, 0xb9, 0xa4, 0x91, 0x73,
	0xc9, 0xbb, 0x70, 0x6f, 0x0f, 0x91, 0xac, 0x1b, 0x84, 0x0a, 0xcd, 0x3d, 0xa2, 0xf6, 0xe0, 0xfe,
	0x3c, 0x41, 0x56, 0xc9, 0xd5, 0x7f, 0x49, 0x70, 0x8b, 0x35, 0xe0, 0x28, 0x1a, 0x84, 0xab, 0xe8,
	0xc5, 0x23, 0xe7, 0x13, 0xb8, 0x1d, 0xd0, 0x55, 0x0a, 0x5a, 0x1c, 0x75, 0x6d, 0x7b, 0xa8, 0xce,
	0x9f, 0x6a, 0xb4, 0xb5, 0x60, 0xc6, 0x7c, 0xfc, 0x08, 0xc2, 0x43, 0x63, 0x66, 0xd0, 0xe3, 0x48,
	0x5f, 0xf6, 0x73, 0x91, 0x91, 0x48, 0xab, 0xbf, 0x86, 0x65, 0x11, 0xf2, 0x95, 0xcc, 0xbb, 0xc6,
	0xf6, 0x9f, 0x97, 0x41, 0x89, 0x5f, 0x68, 0xae, 0x05, 0xe2, 0x08, 0xd6, 0x42, 0xe6, 0x6c, 0x0f,
	0x9f, 0xdf, 0x78, 0x6f, 0x85, 0xca, 0x79, 0xcf, 0x7e, 0x02, 0x1b, 0xb1, 0x56, 0xd1, 0xd2, 0x95,
	0x0b, 0x2c, 0xbd, 0x1e, 0x33, 0xe7, 0x9c, 0x5c, 0xf5, 0x52, 0xae, 0xfb, 0x0d, 0xdc, 0x66, 0x81,
	0x92, 0xc4, 0x72, 0x32, 0x82, 0x7c, 0x0b, 0xea, 0x9e, 0x3e, 0x09, 0x11, 0x4a, 0xfc, 0xf1, 0xcc,
	0x70, 0xfa, 0xcf, 0xf5, 0x09, 0xda, 0x77, 0xc7, 0x58, 0xe3, 0x3c, 0xf9, 0x07, 0xd0, 0xe4, 0xee,
	0x0b, 0x94, 0x32, 0x9d, 0x4e, 0x36, 0x73, 0x61, 0x98, 0x39, 0x14, 0x2d, 0xd6, 0x50, 0x7f, 0x05,
	0x6b, 0xe1, 0x49, 0x5c, 0x7d, 0xfb, 0xc7, 0xb9, 0xed, 0xbb, 0xe2, 0xf6, 0xf9, 0xb0, 0x14, 0x36,
	0xc7, 0x70, 0x67, 0xe7, 0x18, 0x19, 0x3f, 0x17, 0x5f, 0x5f, 0x5f, 0xea, 0xb6, 0x65, 0x5e, 0x6a,
	0xf6, 0xb8, 0xf8, 0x6b, 0x70, 0x38, 0x85, 0x6d, 0xc4, 0x91, 0x7a, 0x75, 0x9b, 0x9f, 0xe4, 0x6c,
	0x7e, 0x47, 0xb0, 0x79, 0x56, 0x2a, 0x08, 0x86, 0xdb, 0xb0, 0x5e, 0x50, 0xbc, 0xbe, 0x2e, 0xab,
	0xa5, 0x8c, 0xd5, 0x7f, 0x91, 0xe0, 0xad, 0xc2, 0xed, 0xb8, 0xd1, 0xa9, 0xf8, 0x95, 0x2e, 0x13,
	0xbf, 0xf2, 0x3e, 0x2c, 0xc5, 0x28, 0x4c, 0x2b, 0x30, 0xe2, 0x3e, 0x33, 0x2f, 0x9f, 0x16, 0x23,
	0xb5, 0x5d, 0xae, 0xa5, 0xfe, 0x56, 0x82, 0x8d, 0x02, 0x94, 0x97, 0x3d, 0x9b, 0x0f, 0xa9, 0xf3,
	0xb8, 0x2e, 0x3f, 0x9d, 0x73, 0xad, 0x11, 0xc4, 0xd5, 0x03, 0xda, 0x54, 0xf8, 0x44, 0x9c, 0xef,
	0x2d, 0xc2, 0x5d, 0xe3, 0x5e, 0xc1, 0x5d, 0x23, 0x7a, 0x31, 0x17, 0x6e, 0x1c, 0x63, 0x7a, 0x7b,
	0x39, 0x7f, 0x41, 0x6e, 0x5e, 0x1a, 0xb8, 0x74, 0x39, 0xe0, 0x1f, 0xc0, 0x6d, 0x36, 0x48, 0x88,
	0xa9, 0x14, 0x61, 0x8d, 0x3a, 0xb9, 0x94, 0xed, 0xe4, 0x6a, 0x0f, 0x3a, 0xa2, 0xc2, 0x39, 0x92,
	0x1f, 0x83, 0xb2, 0x87, 0x88, 0x28, 0x9c, 0x20, 0xff, 0x08, 0x6e, 0xa4, 0x42, 0x31, 0x42, 0xbf,
	0x26, 0xa0, 0x4f, 0xe1, 0x5a, 0x10, 0xc3, 0x33, 0x50, 0xff, 0x5c, 0x86, 0x66, 0xfc, 0xd3, 0x4f,
	0x38, 0x34, 0xe8, 0x4e, 0x06, 0x42, 0x48, 0x09, 0x07, 0x43, 0xef, 0x18, 0xbb, 0xe9, 0xfc, 0x66,
	0xa4, 0x90, 0x87, 0x1c, 0xdd, 0xb2, 0x53, 0x73, 0x0c, 0x23, 0x85, 0x47, 0xc5, 0x1e, 0xba, 0x8d,
	0xdc, 0xf0, 0x06, 0x94, 0xc1, 0x66, 0xec, 0xf7, 0x60, 0x81, 0x60, 0xa2, 0xdb, 0x23, 0x1f, 0x9d,
	0x20, 0x77, 0x8a, 0x52, 0x13, 0x7f, 0x87, 0xb2, 0x34, 0xc6, 0x91, 0x07, 0xb0, 0xc8, 0x44, 0x85,
	0xfb, 0x41, 0x5d, 0x90, 0xbe, 0x49, 0xb9, 0xc2, 0xb5, 0x20, 0xfd, 0x40, 0xd1, 0xb8, 0xf4, 0x03,
	0x85, 0xba, 0x0f, 0x2b, 0x82, 0xfb, 0x13, 0xd7, 0x3f, 0x84, 0x56, 0xe4, 0xcb, 0xc8, 0xeb, 0xcb,
	0x05, 0x5e, 0xd7, 0x12, 0xa9, 0xed, 0xc3, 0xd7, 0x6f, 0xba, 0xd2, 0x17, 0x6f, 0xba, 0xd2, 0x97,
	0x6f, 0xba, 0xa5, 0xff, 0xbe, 0xe9, 0x96, 0xfe, 0x70, 0xd6, 0x2d, 0xfd, 0xf5, 0xac, 0x5b, 0xfa,
	0xfc, 0xac, 0x5b, 0x7a, 0x7d, 0xd6, 0x2d, 0x7d, 0x71, 0xd6, 0x2d, 0x7d, 0x79, 0xd6, 0x2d, 0x7d,
	0xfa, 0x55, 0xb7, 0xf4, 0xd9, 0x57, 0xdd, 0xd2, 0xc7, 0x6f, 0xd3, 0xdf, 0x60, 0x4e, 0xdc, 0x81,
	0xee, 0x59, 0x03, 0xef, 0xd5, 0x20, 0xff, 0xf3, 0xe7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x0a,
	0xf5, 0x38, 0xeb, 0x13, 0x1d, 0x00, 0x00,
}
