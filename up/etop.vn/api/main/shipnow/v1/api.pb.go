// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: etop.vn/api/main/shipnow/v1/api.proto

package v1

import (
	fmt "fmt"
	math "math"

	types "etop.vn/api/main/order/v1/types"
	types1 "etop.vn/api/main/shipping/v1/types"
	v1 "etop.vn/api/meta/v1"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

func (m *ShipnowFulfillment) Reset()         { *m = ShipnowFulfillment{} }
func (m *ShipnowFulfillment) String() string { return proto.CompactTextString(m) }
func (*ShipnowFulfillment) ProtoMessage()    {}
func (*ShipnowFulfillment) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{0}
}
func (m *ShipnowFulfillment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShipnowFulfillment.Unmarshal(m, b)
}
func (m *ShipnowFulfillment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShipnowFulfillment.Marshal(b, m, deterministic)
}
func (m *ShipnowFulfillment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShipnowFulfillment.Merge(m, src)
}
func (m *ShipnowFulfillment) XXX_Size() int {
	return xxx_messageInfo_ShipnowFulfillment.Size(m)
}
func (m *ShipnowFulfillment) XXX_DiscardUnknown() {
	xxx_messageInfo_ShipnowFulfillment.DiscardUnknown(m)
}

var xxx_messageInfo_ShipnowFulfillment proto.InternalMessageInfo

func (m *ShipnowFulfillment) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShipnowFulfillment) GetShopId() int64 {
	if m != nil {
		return m.ShopId
	}
	return 0
}

func (m *ShipnowFulfillment) GetPartnerId() int64 {
	if m != nil {
		return m.PartnerId
	}
	return 0
}

func (m *ShipnowFulfillment) GetPickupAddress() *types.Address {
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

func (m *ShipnowFulfillment) GetValueInfo() types1.ValueInfo {
	if m != nil {
		return m.ValueInfo
	}
	return types1.ValueInfo{}
}

func (m *ShipnowFulfillment) GetShippingNote() string {
	if m != nil {
		return m.ShippingNote
	}
	return ""
}

func (m *ShipnowFulfillment) GetRequestPickupAt() *v1.Timestamp {
	if m != nil {
		return m.RequestPickupAt
	}
	return nil
}

func (m *DeliveryPoint) Reset()         { *m = DeliveryPoint{} }
func (m *DeliveryPoint) String() string { return proto.CompactTextString(m) }
func (*DeliveryPoint) ProtoMessage()    {}
func (*DeliveryPoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{1}
}
func (m *DeliveryPoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeliveryPoint.Unmarshal(m, b)
}
func (m *DeliveryPoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeliveryPoint.Marshal(b, m, deterministic)
}
func (m *DeliveryPoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeliveryPoint.Merge(m, src)
}
func (m *DeliveryPoint) XXX_Size() int {
	return xxx_messageInfo_DeliveryPoint.Size(m)
}
func (m *DeliveryPoint) XXX_DiscardUnknown() {
	xxx_messageInfo_DeliveryPoint.DiscardUnknown(m)
}

var xxx_messageInfo_DeliveryPoint proto.InternalMessageInfo

func (m *DeliveryPoint) GetShippingAddress() *types.Address {
	if m != nil {
		return m.ShippingAddress
	}
	return nil
}

func (m *DeliveryPoint) GetLines() []*types.ItemLine {
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

func (m *DeliveryPoint) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *DeliveryPoint) GetTryOn() types1.TryOnCode {
	if m != nil {
		return m.TryOn
	}
	return types1.TryOnCode_unknown
}

func (m *CreateShipnowFulfillmentCommand) Reset()         { *m = CreateShipnowFulfillmentCommand{} }
func (m *CreateShipnowFulfillmentCommand) String() string { return proto.CompactTextString(m) }
func (*CreateShipnowFulfillmentCommand) ProtoMessage()    {}
func (*CreateShipnowFulfillmentCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{2}
}
func (m *CreateShipnowFulfillmentCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateShipnowFulfillmentCommand.Unmarshal(m, b)
}
func (m *CreateShipnowFulfillmentCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateShipnowFulfillmentCommand.Marshal(b, m, deterministic)
}
func (m *CreateShipnowFulfillmentCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateShipnowFulfillmentCommand.Merge(m, src)
}
func (m *CreateShipnowFulfillmentCommand) XXX_Size() int {
	return xxx_messageInfo_CreateShipnowFulfillmentCommand.Size(m)
}
func (m *CreateShipnowFulfillmentCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateShipnowFulfillmentCommand.DiscardUnknown(m)
}

var xxx_messageInfo_CreateShipnowFulfillmentCommand proto.InternalMessageInfo

func (m *ConfirmShipnowFulfillmentCommand) Reset()         { *m = ConfirmShipnowFulfillmentCommand{} }
func (m *ConfirmShipnowFulfillmentCommand) String() string { return proto.CompactTextString(m) }
func (*ConfirmShipnowFulfillmentCommand) ProtoMessage()    {}
func (*ConfirmShipnowFulfillmentCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{3}
}
func (m *ConfirmShipnowFulfillmentCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmShipnowFulfillmentCommand.Unmarshal(m, b)
}
func (m *ConfirmShipnowFulfillmentCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmShipnowFulfillmentCommand.Marshal(b, m, deterministic)
}
func (m *ConfirmShipnowFulfillmentCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmShipnowFulfillmentCommand.Merge(m, src)
}
func (m *ConfirmShipnowFulfillmentCommand) XXX_Size() int {
	return xxx_messageInfo_ConfirmShipnowFulfillmentCommand.Size(m)
}
func (m *ConfirmShipnowFulfillmentCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmShipnowFulfillmentCommand.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmShipnowFulfillmentCommand proto.InternalMessageInfo

func (m *ConfirmShipnowFulfillmentCommand) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelShipnowFulfillmentCommand) Reset()         { *m = CancelShipnowFulfillmentCommand{} }
func (m *CancelShipnowFulfillmentCommand) String() string { return proto.CompactTextString(m) }
func (*CancelShipnowFulfillmentCommand) ProtoMessage()    {}
func (*CancelShipnowFulfillmentCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{4}
}
func (m *CancelShipnowFulfillmentCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelShipnowFulfillmentCommand.Unmarshal(m, b)
}
func (m *CancelShipnowFulfillmentCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelShipnowFulfillmentCommand.Marshal(b, m, deterministic)
}
func (m *CancelShipnowFulfillmentCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelShipnowFulfillmentCommand.Merge(m, src)
}
func (m *CancelShipnowFulfillmentCommand) XXX_Size() int {
	return xxx_messageInfo_CancelShipnowFulfillmentCommand.Size(m)
}
func (m *CancelShipnowFulfillmentCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelShipnowFulfillmentCommand.DiscardUnknown(m)
}

var xxx_messageInfo_CancelShipnowFulfillmentCommand proto.InternalMessageInfo

func (m *CancelShipnowFulfillmentCommand) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelShipnowFulfillmentCommand) GetCancelReason() string {
	if m != nil {
		return m.CancelReason
	}
	return ""
}

func (m *GetShipnowFulfillmentQueryArgs) Reset()         { *m = GetShipnowFulfillmentQueryArgs{} }
func (m *GetShipnowFulfillmentQueryArgs) String() string { return proto.CompactTextString(m) }
func (*GetShipnowFulfillmentQueryArgs) ProtoMessage()    {}
func (*GetShipnowFulfillmentQueryArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{5}
}
func (m *GetShipnowFulfillmentQueryArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetShipnowFulfillmentQueryArgs.Unmarshal(m, b)
}
func (m *GetShipnowFulfillmentQueryArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetShipnowFulfillmentQueryArgs.Marshal(b, m, deterministic)
}
func (m *GetShipnowFulfillmentQueryArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetShipnowFulfillmentQueryArgs.Merge(m, src)
}
func (m *GetShipnowFulfillmentQueryArgs) XXX_Size() int {
	return xxx_messageInfo_GetShipnowFulfillmentQueryArgs.Size(m)
}
func (m *GetShipnowFulfillmentQueryArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GetShipnowFulfillmentQueryArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GetShipnowFulfillmentQueryArgs proto.InternalMessageInfo

func (m *GetShipnowFulfillmentQueryArgs) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetShipnowFulfillmentQueryResult) Reset()         { *m = GetShipnowFulfillmentQueryResult{} }
func (m *GetShipnowFulfillmentQueryResult) String() string { return proto.CompactTextString(m) }
func (*GetShipnowFulfillmentQueryResult) ProtoMessage()    {}
func (*GetShipnowFulfillmentQueryResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{6}
}
func (m *GetShipnowFulfillmentQueryResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetShipnowFulfillmentQueryResult.Unmarshal(m, b)
}
func (m *GetShipnowFulfillmentQueryResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetShipnowFulfillmentQueryResult.Marshal(b, m, deterministic)
}
func (m *GetShipnowFulfillmentQueryResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetShipnowFulfillmentQueryResult.Merge(m, src)
}
func (m *GetShipnowFulfillmentQueryResult) XXX_Size() int {
	return xxx_messageInfo_GetShipnowFulfillmentQueryResult.Size(m)
}
func (m *GetShipnowFulfillmentQueryResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GetShipnowFulfillmentQueryResult.DiscardUnknown(m)
}

var xxx_messageInfo_GetShipnowFulfillmentQueryResult proto.InternalMessageInfo

func (m *ShipnowEvent) Reset()         { *m = ShipnowEvent{} }
func (m *ShipnowEvent) String() string { return proto.CompactTextString(m) }
func (*ShipnowEvent) ProtoMessage()    {}
func (*ShipnowEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{7}
}
func (m *ShipnowEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShipnowEvent.Unmarshal(m, b)
}
func (m *ShipnowEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShipnowEvent.Marshal(b, m, deterministic)
}
func (m *ShipnowEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShipnowEvent.Merge(m, src)
}
func (m *ShipnowEvent) XXX_Size() int {
	return xxx_messageInfo_ShipnowEvent.Size(m)
}
func (m *ShipnowEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ShipnowEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ShipnowEvent proto.InternalMessageInfo

func (m *ShipnowEvent) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ShipnowEvent) GetUuid() v1.UUID {
	if m != nil {
		return m.Uuid
	}
	return v1.UUID{}
}

func (m *ShipnowEvent) GetCorrelationId() v1.UUID {
	if m != nil {
		return m.CorrelationId
	}
	return v1.UUID{}
}

func (m *ShipnowEvent) GetShipnowFulfillmentId() int64 {
	if m != nil {
		return m.ShipnowFulfillmentId
	}
	return 0
}

func (m *ShipnowEvent) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *ShipnowEvent) GetData() *EventData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *EventData) Reset()         { *m = EventData{} }
func (m *EventData) String() string { return proto.CompactTextString(m) }
func (*EventData) ProtoMessage()    {}
func (*EventData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{8}
}
func (m *EventData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventData.Unmarshal(m, b)
}
func (m *EventData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventData.Marshal(b, m, deterministic)
}
func (m *EventData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventData.Merge(m, src)
}
func (m *EventData) XXX_Size() int {
	return xxx_messageInfo_EventData.Size(m)
}
func (m *EventData) XXX_DiscardUnknown() {
	xxx_messageInfo_EventData.DiscardUnknown(m)
}

var xxx_messageInfo_EventData proto.InternalMessageInfo

func (*EventData_Created) isEventData_Data()               {}
func (*EventData_ConfirmationRequested) isEventData_Data() {}
func (*EventData_ConfirmationAccepted) isEventData_Data()  {}
func (*EventData_ConfirmationRejected) isEventData_Data()  {}
func (*EventData_CancellationRequested) isEventData_Data() {}
func (*EventData_CancellationAccepted) isEventData_Data()  {}
func (*EventData_CancellationRejected) isEventData_Data()  {}

func (m *EventData) GetData() isEventData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *EventData) GetCreated() *CreatedData {
	if x, ok := m.GetData().(*EventData_Created); ok {
		return x.Created
	}
	return nil
}

func (m *EventData) GetConfirmationRequested() *ConfirmationRequestedData {
	if x, ok := m.GetData().(*EventData_ConfirmationRequested); ok {
		return x.ConfirmationRequested
	}
	return nil
}

func (m *EventData) GetConfirmationAccepted() *ConfirmationAcceptedData {
	if x, ok := m.GetData().(*EventData_ConfirmationAccepted); ok {
		return x.ConfirmationAccepted
	}
	return nil
}

func (m *EventData) GetConfirmationRejected() *ConfirmationRejectedData {
	if x, ok := m.GetData().(*EventData_ConfirmationRejected); ok {
		return x.ConfirmationRejected
	}
	return nil
}

func (m *EventData) GetCancellationRequested() *CancellationRequestedData {
	if x, ok := m.GetData().(*EventData_CancellationRequested); ok {
		return x.CancellationRequested
	}
	return nil
}

func (m *EventData) GetCancellationAccepted() *CancellationAcceptedData {
	if x, ok := m.GetData().(*EventData_CancellationAccepted); ok {
		return x.CancellationAccepted
	}
	return nil
}

func (m *EventData) GetCancellationRejected() *CancellationRejectedData {
	if x, ok := m.GetData().(*EventData_CancellationRejected); ok {
		return x.CancellationRejected
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*EventData) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _EventData_OneofMarshaler, _EventData_OneofUnmarshaler, _EventData_OneofSizer, []interface{}{
		(*EventData_Created)(nil),
		(*EventData_ConfirmationRequested)(nil),
		(*EventData_ConfirmationAccepted)(nil),
		(*EventData_ConfirmationRejected)(nil),
		(*EventData_CancellationRequested)(nil),
		(*EventData_CancellationAccepted)(nil),
		(*EventData_CancellationRejected)(nil),
	}
}

func _EventData_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*EventData)
	// data
	switch x := m.Data.(type) {
	case *EventData_Created:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Created); err != nil {
			return err
		}
	case *EventData_ConfirmationRequested:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConfirmationRequested); err != nil {
			return err
		}
	case *EventData_ConfirmationAccepted:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConfirmationAccepted); err != nil {
			return err
		}
	case *EventData_ConfirmationRejected:
		_ = b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ConfirmationRejected); err != nil {
			return err
		}
	case *EventData_CancellationRequested:
		_ = b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CancellationRequested); err != nil {
			return err
		}
	case *EventData_CancellationAccepted:
		_ = b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CancellationAccepted); err != nil {
			return err
		}
	case *EventData_CancellationRejected:
		_ = b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CancellationRejected); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("EventData.Data has unexpected type %T", x)
	}
	return nil
}

func _EventData_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*EventData)
	switch tag {
	case 1: // data.Created
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CreatedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_Created{msg}
		return true, err
	case 4: // data.ConfirmationRequested
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConfirmationRequestedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_ConfirmationRequested{msg}
		return true, err
	case 5: // data.ConfirmationAccepted
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConfirmationAcceptedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_ConfirmationAccepted{msg}
		return true, err
	case 6: // data.ConfirmationRejected
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ConfirmationRejectedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_ConfirmationRejected{msg}
		return true, err
	case 7: // data.CancellationRequested
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CancellationRequestedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_CancellationRequested{msg}
		return true, err
	case 8: // data.CancellationAccepted
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CancellationAcceptedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_CancellationAccepted{msg}
		return true, err
	case 9: // data.CancellationRejected
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CancellationRejectedData)
		err := b.DecodeMessage(msg)
		m.Data = &EventData_CancellationRejected{msg}
		return true, err
	default:
		return false, nil
	}
}

func _EventData_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*EventData)
	// data
	switch x := m.Data.(type) {
	case *EventData_Created:
		s := proto.Size(x.Created)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_ConfirmationRequested:
		s := proto.Size(x.ConfirmationRequested)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_ConfirmationAccepted:
		s := proto.Size(x.ConfirmationAccepted)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_ConfirmationRejected:
		s := proto.Size(x.ConfirmationRejected)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_CancellationRequested:
		s := proto.Size(x.CancellationRequested)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_CancellationAccepted:
		s := proto.Size(x.CancellationAccepted)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EventData_CancellationRejected:
		s := proto.Size(x.CancellationRejected)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func (m *CreatedData) Reset()         { *m = CreatedData{} }
func (m *CreatedData) String() string { return proto.CompactTextString(m) }
func (*CreatedData) ProtoMessage()    {}
func (*CreatedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{9}
}
func (m *CreatedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatedData.Unmarshal(m, b)
}
func (m *CreatedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatedData.Marshal(b, m, deterministic)
}
func (m *CreatedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatedData.Merge(m, src)
}
func (m *CreatedData) XXX_Size() int {
	return xxx_messageInfo_CreatedData.Size(m)
}
func (m *CreatedData) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatedData.DiscardUnknown(m)
}

var xxx_messageInfo_CreatedData proto.InternalMessageInfo

func (m *ConfirmationRequestedData) Reset()         { *m = ConfirmationRequestedData{} }
func (m *ConfirmationRequestedData) String() string { return proto.CompactTextString(m) }
func (*ConfirmationRequestedData) ProtoMessage()    {}
func (*ConfirmationRequestedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{10}
}
func (m *ConfirmationRequestedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmationRequestedData.Unmarshal(m, b)
}
func (m *ConfirmationRequestedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmationRequestedData.Marshal(b, m, deterministic)
}
func (m *ConfirmationRequestedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmationRequestedData.Merge(m, src)
}
func (m *ConfirmationRequestedData) XXX_Size() int {
	return xxx_messageInfo_ConfirmationRequestedData.Size(m)
}
func (m *ConfirmationRequestedData) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmationRequestedData.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmationRequestedData proto.InternalMessageInfo

func (m *ConfirmationAcceptedData) Reset()         { *m = ConfirmationAcceptedData{} }
func (m *ConfirmationAcceptedData) String() string { return proto.CompactTextString(m) }
func (*ConfirmationAcceptedData) ProtoMessage()    {}
func (*ConfirmationAcceptedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{11}
}
func (m *ConfirmationAcceptedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmationAcceptedData.Unmarshal(m, b)
}
func (m *ConfirmationAcceptedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmationAcceptedData.Marshal(b, m, deterministic)
}
func (m *ConfirmationAcceptedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmationAcceptedData.Merge(m, src)
}
func (m *ConfirmationAcceptedData) XXX_Size() int {
	return xxx_messageInfo_ConfirmationAcceptedData.Size(m)
}
func (m *ConfirmationAcceptedData) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmationAcceptedData.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmationAcceptedData proto.InternalMessageInfo

func (m *ConfirmationRejectedData) Reset()         { *m = ConfirmationRejectedData{} }
func (m *ConfirmationRejectedData) String() string { return proto.CompactTextString(m) }
func (*ConfirmationRejectedData) ProtoMessage()    {}
func (*ConfirmationRejectedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{12}
}
func (m *ConfirmationRejectedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfirmationRejectedData.Unmarshal(m, b)
}
func (m *ConfirmationRejectedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfirmationRejectedData.Marshal(b, m, deterministic)
}
func (m *ConfirmationRejectedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfirmationRejectedData.Merge(m, src)
}
func (m *ConfirmationRejectedData) XXX_Size() int {
	return xxx_messageInfo_ConfirmationRejectedData.Size(m)
}
func (m *ConfirmationRejectedData) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfirmationRejectedData.DiscardUnknown(m)
}

var xxx_messageInfo_ConfirmationRejectedData proto.InternalMessageInfo

func (m *CancellationRequestedData) Reset()         { *m = CancellationRequestedData{} }
func (m *CancellationRequestedData) String() string { return proto.CompactTextString(m) }
func (*CancellationRequestedData) ProtoMessage()    {}
func (*CancellationRequestedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{13}
}
func (m *CancellationRequestedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancellationRequestedData.Unmarshal(m, b)
}
func (m *CancellationRequestedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancellationRequestedData.Marshal(b, m, deterministic)
}
func (m *CancellationRequestedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancellationRequestedData.Merge(m, src)
}
func (m *CancellationRequestedData) XXX_Size() int {
	return xxx_messageInfo_CancellationRequestedData.Size(m)
}
func (m *CancellationRequestedData) XXX_DiscardUnknown() {
	xxx_messageInfo_CancellationRequestedData.DiscardUnknown(m)
}

var xxx_messageInfo_CancellationRequestedData proto.InternalMessageInfo

func (m *CancellationAcceptedData) Reset()         { *m = CancellationAcceptedData{} }
func (m *CancellationAcceptedData) String() string { return proto.CompactTextString(m) }
func (*CancellationAcceptedData) ProtoMessage()    {}
func (*CancellationAcceptedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{14}
}
func (m *CancellationAcceptedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancellationAcceptedData.Unmarshal(m, b)
}
func (m *CancellationAcceptedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancellationAcceptedData.Marshal(b, m, deterministic)
}
func (m *CancellationAcceptedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancellationAcceptedData.Merge(m, src)
}
func (m *CancellationAcceptedData) XXX_Size() int {
	return xxx_messageInfo_CancellationAcceptedData.Size(m)
}
func (m *CancellationAcceptedData) XXX_DiscardUnknown() {
	xxx_messageInfo_CancellationAcceptedData.DiscardUnknown(m)
}

var xxx_messageInfo_CancellationAcceptedData proto.InternalMessageInfo

func (m *CancellationRejectedData) Reset()         { *m = CancellationRejectedData{} }
func (m *CancellationRejectedData) String() string { return proto.CompactTextString(m) }
func (*CancellationRejectedData) ProtoMessage()    {}
func (*CancellationRejectedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_76c324fa88396454, []int{15}
}
func (m *CancellationRejectedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancellationRejectedData.Unmarshal(m, b)
}
func (m *CancellationRejectedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancellationRejectedData.Marshal(b, m, deterministic)
}
func (m *CancellationRejectedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancellationRejectedData.Merge(m, src)
}
func (m *CancellationRejectedData) XXX_Size() int {
	return xxx_messageInfo_CancellationRejectedData.Size(m)
}
func (m *CancellationRejectedData) XXX_DiscardUnknown() {
	xxx_messageInfo_CancellationRejectedData.DiscardUnknown(m)
}

var xxx_messageInfo_CancellationRejectedData proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ShipnowFulfillment)(nil), "etop.vn.api.main.shipnow.v1.ShipnowFulfillment")
	proto.RegisterType((*DeliveryPoint)(nil), "etop.vn.api.main.shipnow.v1.DeliveryPoint")
	proto.RegisterType((*CreateShipnowFulfillmentCommand)(nil), "etop.vn.api.main.shipnow.v1.CreateShipnowFulfillmentCommand")
	proto.RegisterType((*ConfirmShipnowFulfillmentCommand)(nil), "etop.vn.api.main.shipnow.v1.ConfirmShipnowFulfillmentCommand")
	proto.RegisterType((*CancelShipnowFulfillmentCommand)(nil), "etop.vn.api.main.shipnow.v1.CancelShipnowFulfillmentCommand")
	proto.RegisterType((*GetShipnowFulfillmentQueryArgs)(nil), "etop.vn.api.main.shipnow.v1.GetShipnowFulfillmentQueryArgs")
	proto.RegisterType((*GetShipnowFulfillmentQueryResult)(nil), "etop.vn.api.main.shipnow.v1.GetShipnowFulfillmentQueryResult")
	proto.RegisterType((*ShipnowEvent)(nil), "etop.vn.api.main.shipnow.v1.ShipnowEvent")
	proto.RegisterType((*EventData)(nil), "etop.vn.api.main.shipnow.v1.EventData")
	proto.RegisterType((*CreatedData)(nil), "etop.vn.api.main.shipnow.v1.CreatedData")
	proto.RegisterType((*ConfirmationRequestedData)(nil), "etop.vn.api.main.shipnow.v1.ConfirmationRequestedData")
	proto.RegisterType((*ConfirmationAcceptedData)(nil), "etop.vn.api.main.shipnow.v1.ConfirmationAcceptedData")
	proto.RegisterType((*ConfirmationRejectedData)(nil), "etop.vn.api.main.shipnow.v1.ConfirmationRejectedData")
	proto.RegisterType((*CancellationRequestedData)(nil), "etop.vn.api.main.shipnow.v1.CancellationRequestedData")
	proto.RegisterType((*CancellationAcceptedData)(nil), "etop.vn.api.main.shipnow.v1.CancellationAcceptedData")
	proto.RegisterType((*CancellationRejectedData)(nil), "etop.vn.api.main.shipnow.v1.CancellationRejectedData")
}

func init() {
	proto.RegisterFile("etop.vn/api/main/shipnow/v1/api.proto", fileDescriptor_76c324fa88396454)
}

var fileDescriptor_76c324fa88396454 = []byte{
	// 1152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xcd, 0x6e, 0xdb, 0x46,
	0x17, 0x25, 0xf5, 0x6b, 0x5f, 0xc5, 0xce, 0xf7, 0x4d, 0xec, 0x82, 0x51, 0x10, 0x5a, 0x60, 0xd1,
	0x42, 0x69, 0x51, 0x0a, 0x76, 0x51, 0x23, 0x48, 0x13, 0x14, 0xfe, 0x89, 0x5b, 0x05, 0x4d, 0x93,
	0xd2, 0x49, 0x8a, 0x76, 0x23, 0xb0, 0xe2, 0x48, 0x66, 0x23, 0xcd, 0xb0, 0xc3, 0x11, 0x0d, 0xbd,
	0x41, 0xdb, 0x55, 0x81, 0x6e, 0xfa, 0x08, 0x05, 0xba, 0xec, 0xa6, 0x8f, 0xe0, 0x45, 0x17, 0x5e,
	0x76, 0x15, 0xc4, 0xf2, 0xba, 0x40, 0x1e, 0xa1, 0xe0, 0x70, 0x28, 0xd3, 0x12, 0xa9, 0x48, 0x41,
	0x37, 0x41, 0x34, 0xf7, 0x9e, 0x73, 0xe7, 0x1e, 0x9e, 0xb9, 0xd7, 0xf0, 0x0e, 0xe6, 0xd4, 0x33,
	0x03, 0xd2, 0xb0, 0x3d, 0xb7, 0xd1, 0xb7, 0x5d, 0xd2, 0xf0, 0x8f, 0x5c, 0x8f, 0xd0, 0xe3, 0x46,
	0xb0, 0x19, 0x9e, 0x99, 0x1e, 0xa3, 0x9c, 0xa2, 0x1b, 0x32, 0xcd, 0x0c, 0x8f, 0xc2, 0x34, 0x53,
	0xa6, 0x99, 0xc1, 0x66, 0x75, 0xad, 0x4b, 0xbb, 0x54, 0xe4, 0x35, 0xc2, 0xff, 0x45, 0x90, 0xea,
	0xfb, 0x53, 0xcc, 0x94, 0x39, 0x98, 0x85, 0xbc, 0x7c, 0xe8, 0x61, 0x3f, 0xfa, 0x57, 0x26, 0x9b,
	0xa9, 0xd7, 0xf0, 0x5c, 0xd2, 0x4d, 0xe4, 0xb3, 0x21, 0x25, 0x8b, 0xe4, 0x27, 0xf8, 0x6f, 0x5e,
	0xca, 0xc7, 0xdc, 0xbe, 0xd4, 0x9e, 0xf1, 0x57, 0x11, 0xd0, 0x61, 0xd4, 0xd0, 0xc1, 0xa0, 0xd7,
	0x71, 0x7b, 0xbd, 0x3e, 0x26, 0x1c, 0xad, 0x41, 0xce, 0x75, 0x34, 0xb5, 0xa6, 0xd6, 0xf3, 0xbb,
	0x85, 0x93, 0x17, 0x1b, 0x8a, 0x95, 0x73, 0x1d, 0x74, 0x13, 0xca, 0xfe, 0x11, 0xf5, 0x5a, 0xae,
	0xa3, 0xe5, 0x12, 0xa1, 0x52, 0x78, 0xd8, 0x74, 0xd0, 0xdb, 0x00, 0x9e, 0xcd, 0x38, 0xc1, 0x2c,
	0xcc, 0xc8, 0x27, 0x32, 0x96, 0xe5, 0x79, 0xd3, 0x41, 0x8f, 0x60, 0xd5, 0x73, 0xdb, 0xcf, 0x07,
	0x5e, 0xcb, 0x76, 0x1c, 0x86, 0x7d, 0x5f, 0x2b, 0xd4, 0xd4, 0x7a, 0x65, 0xab, 0x6e, 0x4e, 0x09,
	0x2d, 0x54, 0x33, 0x83, 0x4d, 0x33, 0xea, 0x67, 0x27, 0xca, 0xb7, 0x56, 0x22, 0xbc, 0xfc, 0x89,
	0x0e, 0xe1, 0xaa, 0x83, 0x7b, 0x6e, 0x80, 0xd9, 0xb0, 0xe5, 0x51, 0x97, 0x70, 0x5f, 0x2b, 0xd6,
	0xf2, 0xf5, 0xca, 0xd6, 0x7b, 0xe6, 0x8c, 0x4f, 0x67, 0xee, 0x4b, 0xcc, 0xe3, 0x10, 0x62, 0xad,
	0x3a, 0xc9, 0x9f, 0x3e, 0xd2, 0xa1, 0xdc, 0xb6, 0x19, 0x73, 0x31, 0xd3, 0x4a, 0x35, 0xb5, 0xbe,
	0x2c, 0xfb, 0x88, 0x0f, 0xd1, 0x6d, 0x58, 0x8f, 0x65, 0x6f, 0xf9, 0x98, 0x05, 0x6e, 0x1b, 0xb7,
	0xda, 0xd4, 0xc1, 0x5a, 0x39, 0x91, 0x7d, 0x2d, 0x4e, 0x39, 0x8c, 0x32, 0xf6, 0xa8, 0x83, 0xd1,
	0x36, 0xac, 0x4d, 0x21, 0x3b, 0x18, 0x6b, 0x4b, 0x35, 0xb5, 0x5e, 0x94, 0x40, 0x34, 0x01, 0x3c,
	0xc0, 0x18, 0x7d, 0x0d, 0x95, 0x63, 0xec, 0x76, 0x8f, 0x78, 0xcb, 0x25, 0x1d, 0xaa, 0x2d, 0x0b,
	0xd1, 0xcc, 0xf4, 0x16, 0x43, 0xe8, 0x85, 0x6e, 0x5f, 0x09, 0x58, 0x93, 0x74, 0xe8, 0xee, 0x52,
	0x48, 0x7f, 0xfa, 0x62, 0x43, 0xb5, 0xe0, 0x78, 0x7c, 0x8a, 0x2c, 0x80, 0xc0, 0xee, 0x0d, 0x70,
	0xc4, 0x0c, 0x82, 0xf9, 0x83, 0x79, 0x98, 0x9f, 0x85, 0x28, 0x41, 0x2c, 0x3f, 0x73, 0x10, 0x1f,
	0xa0, 0x5b, 0xb0, 0x32, 0x6e, 0x93, 0x50, 0x8e, 0xb5, 0x4a, 0x42, 0x98, 0x2b, 0x71, 0xe8, 0x0b,
	0xca, 0x31, 0x7a, 0x00, 0xff, 0x67, 0xf8, 0xfb, 0x01, 0xf6, 0x79, 0x2b, 0x76, 0x06, 0xd7, 0xae,
	0x88, 0x5b, 0xe8, 0x97, 0x6f, 0x81, 0xb9, 0x1d, 0x16, 0x7f, 0xe2, 0xf6, 0xb1, 0xcf, 0xed, 0xbe,
	0x67, 0x5d, 0x95, 0xc0, 0xc7, 0x91, 0x23, 0xb8, 0xf1, 0x2a, 0x0f, 0x2b, 0x97, 0xbe, 0x2c, 0x3a,
	0x84, 0xff, 0x8d, 0x2f, 0x12, 0x3b, 0x4e, 0x5d, 0xd0, 0x71, 0x57, 0x63, 0x86, 0xd8, 0x73, 0x9f,
	0x40, 0xb1, 0xe7, 0x12, 0xec, 0x6b, 0x39, 0xe1, 0xb4, 0x5b, 0xaf, 0x65, 0x6a, 0x72, 0xdc, 0xff,
	0xdc, 0x25, 0xd8, 0x8a, 0x70, 0xd3, 0xf2, 0xe4, 0x33, 0xe5, 0xd9, 0x80, 0x25, 0x41, 0x16, 0xbe,
	0xa9, 0x42, 0xe2, 0x4d, 0x95, 0xc5, 0x69, 0xd3, 0x99, 0x74, 0x46, 0xf1, 0x3f, 0x74, 0xc6, 0xb3,
	0x4b, 0xce, 0x28, 0xbd, 0x89, 0x33, 0x2e, 0x88, 0x13, 0xee, 0x78, 0x00, 0x25, 0xce, 0x86, 0x2d,
	0x4a, 0x84, 0x2d, 0x56, 0x5f, 0xcb, 0x29, 0xe6, 0xdf, 0x13, 0x36, 0x7c, 0x44, 0xc2, 0x37, 0x24,
	0x05, 0x28, 0xf2, 0xf0, 0xc0, 0xf0, 0x60, 0x63, 0x8f, 0x61, 0x9b, 0xe3, 0xe9, 0x31, 0xb6, 0x47,
	0xfb, 0x7d, 0x9b, 0x38, 0xe8, 0x21, 0x14, 0x1c, 0x9b, 0xdb, 0xf2, 0xbb, 0x37, 0x66, 0xce, 0x85,
	0x69, 0x96, 0x44, 0x0b, 0x82, 0xc6, 0xb8, 0x0d, 0xb5, 0x3d, 0x4a, 0x3a, 0x2e, 0xeb, 0x67, 0x97,
	0x4c, 0x1d, 0xa0, 0xc6, 0xb7, 0xb0, 0xb1, 0x67, 0x93, 0x36, 0xee, 0x2d, 0x08, 0x0c, 0xfd, 0xd2,
	0x16, 0xc0, 0x16, 0xc3, 0xb6, 0x4f, 0x89, 0x98, 0xbf, 0x63, 0xbf, 0x44, 0x21, 0x4b, 0x44, 0x8c,
	0x6d, 0xd0, 0x3f, 0xc5, 0x7c, 0xba, 0xc0, 0x97, 0x03, 0xcc, 0x86, 0x3b, 0xac, 0xeb, 0x67, 0xdc,
	0xcd, 0x80, 0x5a, 0x36, 0xce, 0xc2, 0xfe, 0xa0, 0xc7, 0x8d, 0x3f, 0x72, 0x70, 0x45, 0x66, 0xdc,
	0x0f, 0xb2, 0xf7, 0xc4, 0x87, 0x50, 0x18, 0x0c, 0xe4, 0x92, 0xa8, 0x6c, 0x5d, 0x4f, 0x7d, 0xc4,
	0x4f, 0x9f, 0x36, 0xf7, 0x25, 0x44, 0x24, 0xa3, 0x03, 0x58, 0x6d, 0x53, 0xc6, 0x70, 0xcf, 0xe6,
	0x2e, 0x25, 0xf1, 0x06, 0x99, 0x03, 0xbe, 0x92, 0x80, 0x35, 0x1d, 0x74, 0x07, 0xde, 0x92, 0x9f,
	0xb3, 0xd5, 0xb9, 0xe8, 0x62, 0xf2, 0xf5, 0xac, 0xf9, 0x53, 0x8d, 0x36, 0x1d, 0xa4, 0x41, 0x21,
	0xf4, 0xaf, 0x98, 0x81, 0xf1, 0x30, 0x16, 0x27, 0xe8, 0x8e, 0xb4, 0x50, 0x45, 0xdc, 0xe9, 0xdd,
	0x99, 0x16, 0x12, 0xd2, 0xec, 0xdb, 0xdc, 0x96, 0x7e, 0xf9, 0xa7, 0x08, 0xcb, 0xe3, 0x33, 0xb4,
	0x0f, 0xe5, 0xc8, 0xaf, 0x4e, 0xf6, 0x1c, 0x4a, 0x90, 0xc9, 0xdc, 0x10, 0xfa, 0x99, 0x62, 0xc5,
	0x50, 0x44, 0x60, 0x5d, 0x7a, 0x50, 0xf4, 0x6d, 0x45, 0x73, 0x10, 0x3b, 0x72, 0x9b, 0x6e, 0xcf,
	0xe6, 0x4c, 0x43, 0xca, 0x0a, 0xe9, 0xb4, 0xe8, 0x39, 0xac, 0x25, 0x03, 0x3b, 0xed, 0x36, 0xf6,
	0xc2, 0x72, 0xd1, 0xb4, 0xf9, 0x68, 0xee, 0x72, 0x31, 0x50, 0x56, 0x4b, 0x25, 0x9d, 0x2c, 0x66,
	0xe1, 0xef, 0x70, 0x3b, 0x2c, 0x56, 0x5a, 0xb0, 0x58, 0x0c, 0x4c, 0x2b, 0x16, 0xc7, 0x84, 0x92,
	0xe2, 0xfd, 0xf4, 0x26, 0x94, 0x2c, 0xcf, 0xa3, 0x64, 0x1a, 0x72, 0xac, 0x64, 0x5a, 0x50, 0x34,
	0x97, 0x08, 0x8c, 0x95, 0x5c, 0x9a, 0xa7, 0xb9, 0x14, 0xe0, 0xb8, 0xb9, 0x94, 0xd8, 0x64, 0xb1,
	0xb1, 0x92, 0xcb, 0x0b, 0x16, 0x9b, 0x52, 0x32, 0x25, 0xb6, 0x5b, 0x8a, 0xde, 0x88, 0xb1, 0x02,
	0x95, 0x84, 0x6b, 0x8d, 0x1b, 0x70, 0x3d, 0xd3, 0x70, 0x46, 0x15, 0xb4, 0x2c, 0x7b, 0x4c, 0xc6,
	0x92, 0x77, 0x10, 0xa4, 0x59, 0xda, 0x0b, 0x60, 0x86, 0x52, 0x93, 0xb1, 0x24, 0xe9, 0xd6, 0x4f,
	0x79, 0x58, 0x95, 0xe3, 0xed, 0xa1, 0x4d, 0xec, 0x2e, 0x66, 0xe8, 0x47, 0x15, 0xb4, 0xac, 0xf5,
	0x82, 0xee, 0xce, 0xf1, 0x72, 0x33, 0x27, 0x7d, 0x75, 0xd1, 0x3d, 0x84, 0x82, 0xb1, 0x90, 0x29,
	0xc1, 0x7b, 0xf3, 0xbc, 0x8a, 0xec, 0xcb, 0x54, 0x53, 0xa7, 0xec, 0xfd, 0xbe, 0xc7, 0x87, 0x88,
	0xc7, 0x92, 0x2d, 0x2e, 0xc1, 0xec, 0x65, 0x37, 0xab, 0xea, 0xd6, 0xef, 0x2a, 0x5c, 0x93, 0x48,
	0xb1, 0x82, 0xe4, 0xdf, 0xc2, 0xe8, 0x17, 0x15, 0xd6, 0x53, 0x17, 0x15, 0xfa, 0x78, 0xe6, 0x5d,
	0x66, 0x2f, 0xc5, 0xea, 0xbd, 0x37, 0x04, 0x47, 0x9b, 0x71, 0xf7, 0xee, 0xc9, 0x99, 0xae, 0x9e,
	0x9e, 0xe9, 0xca, 0xcb, 0x33, 0x5d, 0x79, 0x75, 0xa6, 0x2b, 0x3f, 0x8c, 0x74, 0xe5, 0xb7, 0x91,
	0xae, 0xfc, 0x39, 0xd2, 0x95, 0x93, 0x91, 0xae, 0x9c, 0x8e, 0x74, 0xe5, 0xe5, 0x48, 0x57, 0x7e,
	0x3e, 0xd7, 0x95, 0x5f, 0xcf, 0x75, 0xe5, 0xf4, 0x5c, 0x57, 0xfe, 0x3e, 0xd7, 0x95, 0x6f, 0x72,
	0xc1, 0xe6, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xfe, 0xa6, 0x39, 0x8c, 0x0e, 0x00, 0x00,
}
