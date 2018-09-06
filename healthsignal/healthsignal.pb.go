// Code generated by protoc-gen-go. DO NOT EDIT.
// source: healthsignal.proto

package healthsignal

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type EndpointsQuery struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EndpointsQuery) Reset()         { *m = EndpointsQuery{} }
func (m *EndpointsQuery) String() string { return proto.CompactTextString(m) }
func (*EndpointsQuery) ProtoMessage()    {}
func (*EndpointsQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_healthsignal_a1c3dd764c55366c, []int{0}
}
func (m *EndpointsQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndpointsQuery.Unmarshal(m, b)
}
func (m *EndpointsQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndpointsQuery.Marshal(b, m, deterministic)
}
func (dst *EndpointsQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndpointsQuery.Merge(dst, src)
}
func (m *EndpointsQuery) XXX_Size() int {
	return xxx_messageInfo_EndpointsQuery.Size(m)
}
func (m *EndpointsQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_EndpointsQuery.DiscardUnknown(m)
}

var xxx_messageInfo_EndpointsQuery proto.InternalMessageInfo

type Endpoints struct {
	Endpoints            []*Endpoint `protobuf:"bytes,1,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Endpoints) Reset()         { *m = Endpoints{} }
func (m *Endpoints) String() string { return proto.CompactTextString(m) }
func (*Endpoints) ProtoMessage()    {}
func (*Endpoints) Descriptor() ([]byte, []int) {
	return fileDescriptor_healthsignal_a1c3dd764c55366c, []int{1}
}
func (m *Endpoints) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoints.Unmarshal(m, b)
}
func (m *Endpoints) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoints.Marshal(b, m, deterministic)
}
func (dst *Endpoints) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoints.Merge(dst, src)
}
func (m *Endpoints) XXX_Size() int {
	return xxx_messageInfo_Endpoints.Size(m)
}
func (m *Endpoints) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoints.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoints proto.InternalMessageInfo

func (m *Endpoints) GetEndpoints() []*Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

type Endpoint struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url                  string   `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	ExpectedStatus       int32    `protobuf:"varint,3,opt,name=expectedStatus,proto3" json:"expectedStatus,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_healthsignal_a1c3dd764c55366c, []int{2}
}
func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (dst *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(dst, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Endpoint) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Endpoint) GetExpectedStatus() int32 {
	if m != nil {
		return m.ExpectedStatus
	}
	return 0
}

func (m *Endpoint) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_healthsignal_a1c3dd764c55366c, []int{3}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Ping struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	RequestDuration      int32    `protobuf:"varint,2,opt,name=requestDuration,proto3" json:"requestDuration,omitempty"`
	ConnectionDuration   int32    `protobuf:"varint,3,opt,name=connectionDuration,proto3" json:"connectionDuration,omitempty"`
	Location             string   `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	Endpoint             int32    `protobuf:"varint,5,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	TimedOut             bool     `protobuf:"varint,6,opt,name=timedOut,proto3" json:"timedOut,omitempty"`
	OnlineStatus         int32    `protobuf:"varint,7,opt,name=onlineStatus,proto3" json:"onlineStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_healthsignal_a1c3dd764c55366c, []int{4}
}
func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (dst *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(dst, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

func (m *Ping) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Ping) GetRequestDuration() int32 {
	if m != nil {
		return m.RequestDuration
	}
	return 0
}

func (m *Ping) GetConnectionDuration() int32 {
	if m != nil {
		return m.ConnectionDuration
	}
	return 0
}

func (m *Ping) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Ping) GetEndpoint() int32 {
	if m != nil {
		return m.Endpoint
	}
	return 0
}

func (m *Ping) GetTimedOut() bool {
	if m != nil {
		return m.TimedOut
	}
	return false
}

func (m *Ping) GetOnlineStatus() int32 {
	if m != nil {
		return m.OnlineStatus
	}
	return 0
}

func init() {
	proto.RegisterType((*EndpointsQuery)(nil), "EndpointsQuery")
	proto.RegisterType((*Endpoints)(nil), "Endpoints")
	proto.RegisterType((*Endpoint)(nil), "Endpoint")
	proto.RegisterType((*Empty)(nil), "Empty")
	proto.RegisterType((*Ping)(nil), "Ping")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HealthSignalClient is the client API for HealthSignal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HealthSignalClient interface {
	SavePing(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Empty, error)
	GetEndpoints(ctx context.Context, in *EndpointsQuery, opts ...grpc.CallOption) (*Endpoints, error)
}

type healthSignalClient struct {
	cc *grpc.ClientConn
}

func NewHealthSignalClient(cc *grpc.ClientConn) HealthSignalClient {
	return &healthSignalClient{cc}
}

func (c *healthSignalClient) SavePing(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/HealthSignal/SavePing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthSignalClient) GetEndpoints(ctx context.Context, in *EndpointsQuery, opts ...grpc.CallOption) (*Endpoints, error) {
	out := new(Endpoints)
	err := c.cc.Invoke(ctx, "/HealthSignal/GetEndpoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthSignalServer is the server API for HealthSignal service.
type HealthSignalServer interface {
	SavePing(context.Context, *Ping) (*Empty, error)
	GetEndpoints(context.Context, *EndpointsQuery) (*Endpoints, error)
}

func RegisterHealthSignalServer(s *grpc.Server, srv HealthSignalServer) {
	s.RegisterService(&_HealthSignal_serviceDesc, srv)
}

func _HealthSignal_SavePing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthSignalServer).SavePing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HealthSignal/SavePing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthSignalServer).SavePing(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

func _HealthSignal_GetEndpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndpointsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthSignalServer).GetEndpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HealthSignal/GetEndpoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthSignalServer).GetEndpoints(ctx, req.(*EndpointsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _HealthSignal_serviceDesc = grpc.ServiceDesc{
	ServiceName: "HealthSignal",
	HandlerType: (*HealthSignalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SavePing",
			Handler:    _HealthSignal_SavePing_Handler,
		},
		{
			MethodName: "GetEndpoints",
			Handler:    _HealthSignal_GetEndpoints_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "healthsignal.proto",
}

func init() { proto.RegisterFile("healthsignal.proto", fileDescriptor_healthsignal_a1c3dd764c55366c) }

var fileDescriptor_healthsignal_a1c3dd764c55366c = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xc1, 0x4e, 0x83, 0x40,
	0x10, 0x2d, 0x6d, 0xa1, 0x30, 0x36, 0x6d, 0x33, 0x07, 0xb3, 0xa9, 0x17, 0xb2, 0x07, 0xe5, 0x22,
	0x87, 0xea, 0x27, 0xd8, 0xe8, 0x4d, 0xa5, 0x37, 0x6f, 0x08, 0x93, 0x76, 0x13, 0xba, 0x8b, 0xb0,
	0x18, 0xfb, 0xd7, 0x7e, 0x82, 0x61, 0xb7, 0xd0, 0xb4, 0xf1, 0xc4, 0xbc, 0xf7, 0x66, 0x86, 0xc9,
	0x7b, 0x0b, 0xb8, 0xa3, 0xb4, 0xd0, 0xbb, 0x5a, 0x6c, 0x65, 0x5a, 0xc4, 0x65, 0xa5, 0xb4, 0xe2,
	0x0b, 0x98, 0xad, 0x65, 0x5e, 0x2a, 0x21, 0x75, 0xfd, 0xde, 0x50, 0x75, 0xe0, 0x8f, 0x10, 0xf4,
	0x0c, 0xde, 0x41, 0x40, 0x1d, 0x60, 0x4e, 0x38, 0x8a, 0xae, 0x56, 0x41, 0xdc, 0xc9, 0xc9, 0x49,
	0xe3, 0x3b, 0xf0, 0x3b, 0x1a, 0x67, 0x30, 0x14, 0x39, 0x73, 0x42, 0x27, 0x72, 0x93, 0xa1, 0xc8,
	0x71, 0x01, 0xa3, 0xa6, 0x2a, 0xd8, 0x30, 0x74, 0xa2, 0x20, 0x69, 0x4b, 0xbc, 0x85, 0x19, 0xfd,
	0x94, 0x94, 0x69, 0xca, 0x37, 0x3a, 0xd5, 0x4d, 0xcd, 0x46, 0xa6, 0xfb, 0x82, 0x45, 0x84, 0xb1,
	0x4c, 0xf7, 0xc4, 0xc6, 0x66, 0xd4, 0xd4, 0x7c, 0x02, 0xee, 0x7a, 0x5f, 0xea, 0x03, 0xff, 0x75,
	0x60, 0xfc, 0x26, 0xe4, 0x16, 0xaf, 0xc1, 0xab, 0xed, 0x16, 0xfb, 0xcf, 0x23, 0xc2, 0x08, 0xe6,
	0x15, 0x7d, 0x35, 0x54, 0xeb, 0xa7, 0xa6, 0x4a, 0xb5, 0x50, 0xd2, 0xdc, 0xe0, 0x26, 0x97, 0x34,
	0xc6, 0x80, 0x99, 0x92, 0x92, 0xb2, 0x16, 0xf5, 0xcd, 0xf6, 0xa6, 0x7f, 0x14, 0x5c, 0x82, 0x5f,
	0xa8, 0xcc, 0x76, 0xd9, 0xdb, 0x7a, 0xdc, 0x6a, 0x9d, 0x2d, 0xcc, 0x35, 0x1b, 0x7a, 0xdc, 0x6a,
	0x5a, 0xec, 0x29, 0x7f, 0x6d, 0x34, 0xf3, 0x42, 0x27, 0xf2, 0x93, 0x1e, 0x23, 0x87, 0xa9, 0x92,
	0x85, 0x90, 0x74, 0x74, 0x64, 0x62, 0x66, 0xcf, 0xb8, 0xd5, 0x07, 0x4c, 0x5f, 0x4c, 0x86, 0x1b,
	0x93, 0x21, 0xde, 0x80, 0xbf, 0x49, 0xbf, 0xc9, 0xb8, 0xe0, 0xc6, 0xed, 0x67, 0xe9, 0xc5, 0xd6,
	0x9d, 0x01, 0xde, 0xc3, 0xf4, 0x99, 0xf4, 0x29, 0xcb, 0x79, 0x7c, 0x9e, 0xf4, 0x12, 0x4e, 0x04,
	0x1f, 0x7c, 0x7a, 0xe6, 0x41, 0x3c, 0xfc, 0x05, 0x00, 0x00, 0xff, 0xff, 0x17, 0x79, 0xf1, 0x6a,
	0x26, 0x02, 0x00, 0x00,
}
