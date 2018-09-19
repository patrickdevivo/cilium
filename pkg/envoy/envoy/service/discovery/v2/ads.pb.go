// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/service/discovery/v2/ads.proto

package v2

import (
	fmt "fmt"
	v2 "github.com/cilium/cilium/pkg/envoy/envoy/api/v2"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

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

// [#not-implemented-hide:] Not configuration. Workaround c++ protobuf issue with importing
// services: https://github.com/google/protobuf/issues/4221
type AdsDummy struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdsDummy) Reset()         { *m = AdsDummy{} }
func (m *AdsDummy) String() string { return proto.CompactTextString(m) }
func (*AdsDummy) ProtoMessage()    {}
func (*AdsDummy) Descriptor() ([]byte, []int) {
	return fileDescriptor_187fd5dcc2dab695, []int{0}
}

func (m *AdsDummy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdsDummy.Unmarshal(m, b)
}
func (m *AdsDummy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdsDummy.Marshal(b, m, deterministic)
}
func (m *AdsDummy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdsDummy.Merge(m, src)
}
func (m *AdsDummy) XXX_Size() int {
	return xxx_messageInfo_AdsDummy.Size(m)
}
func (m *AdsDummy) XXX_DiscardUnknown() {
	xxx_messageInfo_AdsDummy.DiscardUnknown(m)
}

var xxx_messageInfo_AdsDummy proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AdsDummy)(nil), "envoy.service.discovery.v2.AdsDummy")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AggregatedDiscoveryServiceClient is the client API for AggregatedDiscoveryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AggregatedDiscoveryServiceClient interface {
	// This is a gRPC-only API.
	StreamAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_StreamAggregatedResourcesClient, error)
	IncrementalAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_IncrementalAggregatedResourcesClient, error)
}

type aggregatedDiscoveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewAggregatedDiscoveryServiceClient(cc *grpc.ClientConn) AggregatedDiscoveryServiceClient {
	return &aggregatedDiscoveryServiceClient{cc}
}

func (c *aggregatedDiscoveryServiceClient) StreamAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_StreamAggregatedResourcesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AggregatedDiscoveryService_serviceDesc.Streams[0], "/envoy.service.discovery.v2.AggregatedDiscoveryService/StreamAggregatedResources", opts...)
	if err != nil {
		return nil, err
	}
	x := &aggregatedDiscoveryServiceStreamAggregatedResourcesClient{stream}
	return x, nil
}

type AggregatedDiscoveryService_StreamAggregatedResourcesClient interface {
	Send(*v2.DiscoveryRequest) error
	Recv() (*v2.DiscoveryResponse, error)
	grpc.ClientStream
}

type aggregatedDiscoveryServiceStreamAggregatedResourcesClient struct {
	grpc.ClientStream
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesClient) Send(m *v2.DiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesClient) Recv() (*v2.DiscoveryResponse, error) {
	m := new(v2.DiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aggregatedDiscoveryServiceClient) IncrementalAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_IncrementalAggregatedResourcesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AggregatedDiscoveryService_serviceDesc.Streams[1], "/envoy.service.discovery.v2.AggregatedDiscoveryService/IncrementalAggregatedResources", opts...)
	if err != nil {
		return nil, err
	}
	x := &aggregatedDiscoveryServiceIncrementalAggregatedResourcesClient{stream}
	return x, nil
}

type AggregatedDiscoveryService_IncrementalAggregatedResourcesClient interface {
	Send(*v2.IncrementalDiscoveryRequest) error
	Recv() (*v2.IncrementalDiscoveryResponse, error)
	grpc.ClientStream
}

type aggregatedDiscoveryServiceIncrementalAggregatedResourcesClient struct {
	grpc.ClientStream
}

func (x *aggregatedDiscoveryServiceIncrementalAggregatedResourcesClient) Send(m *v2.IncrementalDiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceIncrementalAggregatedResourcesClient) Recv() (*v2.IncrementalDiscoveryResponse, error) {
	m := new(v2.IncrementalDiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AggregatedDiscoveryServiceServer is the server API for AggregatedDiscoveryService service.
type AggregatedDiscoveryServiceServer interface {
	// This is a gRPC-only API.
	StreamAggregatedResources(AggregatedDiscoveryService_StreamAggregatedResourcesServer) error
	IncrementalAggregatedResources(AggregatedDiscoveryService_IncrementalAggregatedResourcesServer) error
}

func RegisterAggregatedDiscoveryServiceServer(s *grpc.Server, srv AggregatedDiscoveryServiceServer) {
	s.RegisterService(&_AggregatedDiscoveryService_serviceDesc, srv)
}

func _AggregatedDiscoveryService_StreamAggregatedResources_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AggregatedDiscoveryServiceServer).StreamAggregatedResources(&aggregatedDiscoveryServiceStreamAggregatedResourcesServer{stream})
}

type AggregatedDiscoveryService_StreamAggregatedResourcesServer interface {
	Send(*v2.DiscoveryResponse) error
	Recv() (*v2.DiscoveryRequest, error)
	grpc.ServerStream
}

type aggregatedDiscoveryServiceStreamAggregatedResourcesServer struct {
	grpc.ServerStream
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesServer) Send(m *v2.DiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesServer) Recv() (*v2.DiscoveryRequest, error) {
	m := new(v2.DiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _AggregatedDiscoveryService_IncrementalAggregatedResources_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AggregatedDiscoveryServiceServer).IncrementalAggregatedResources(&aggregatedDiscoveryServiceIncrementalAggregatedResourcesServer{stream})
}

type AggregatedDiscoveryService_IncrementalAggregatedResourcesServer interface {
	Send(*v2.IncrementalDiscoveryResponse) error
	Recv() (*v2.IncrementalDiscoveryRequest, error)
	grpc.ServerStream
}

type aggregatedDiscoveryServiceIncrementalAggregatedResourcesServer struct {
	grpc.ServerStream
}

func (x *aggregatedDiscoveryServiceIncrementalAggregatedResourcesServer) Send(m *v2.IncrementalDiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceIncrementalAggregatedResourcesServer) Recv() (*v2.IncrementalDiscoveryRequest, error) {
	m := new(v2.IncrementalDiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AggregatedDiscoveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "envoy.service.discovery.v2.AggregatedDiscoveryService",
	HandlerType: (*AggregatedDiscoveryServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamAggregatedResources",
			Handler:       _AggregatedDiscoveryService_StreamAggregatedResources_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "IncrementalAggregatedResources",
			Handler:       _AggregatedDiscoveryService_IncrementalAggregatedResources_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "envoy/service/discovery/v2/ads.proto",
}

func init() {
	proto.RegisterFile("envoy/service/discovery/v2/ads.proto", fileDescriptor_187fd5dcc2dab695)
}

var fileDescriptor_187fd5dcc2dab695 = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0xd0, 0xbd, 0x4a, 0x04, 0x31,
	0x14, 0x05, 0x60, 0x63, 0xa1, 0x92, 0x32, 0x9d, 0x41, 0x56, 0x58, 0x2c, 0xd4, 0x22, 0x91, 0xf8,
	0x04, 0x2b, 0xdb, 0xd8, 0xee, 0x76, 0x76, 0x99, 0xe4, 0x32, 0x04, 0xcc, 0x8f, 0xb9, 0x99, 0xc0,
	0x14, 0xf6, 0xbe, 0xb5, 0xe2, 0x8c, 0xce, 0x8c, 0xa8, 0xb0, 0xf5, 0xfd, 0xce, 0x39, 0x21, 0xf4,
	0x0a, 0x42, 0x8d, 0xbd, 0x44, 0xc8, 0xd5, 0x19, 0x90, 0xd6, 0xa1, 0x89, 0x15, 0x72, 0x2f, 0xab,
	0x92, 0xda, 0xa2, 0x48, 0x39, 0x96, 0xc8, 0xf8, 0xa0, 0xc4, 0x97, 0x12, 0x93, 0x12, 0x55, 0xf1,
	0x8b, 0xb1, 0x41, 0x27, 0xf7, 0x99, 0x99, 0x4f, 0x43, 0x72, 0x4d, 0xe9, 0xd9, 0xc6, 0xe2, 0xb6,
	0xf3, 0xbe, 0x57, 0xef, 0x84, 0xf2, 0x4d, 0xdb, 0x66, 0x68, 0x75, 0x01, 0xbb, 0xfd, 0x96, 0xfb,
	0xb1, 0x95, 0x35, 0xf4, 0x7c, 0x5f, 0x32, 0x68, 0x3f, 0x9b, 0x1d, 0x60, 0xec, 0xb2, 0x01, 0x64,
	0x2b, 0x31, 0x3e, 0x41, 0x27, 0x27, 0xaa, 0x12, 0x53, 0x78, 0x07, 0x2f, 0x1d, 0x60, 0xe1, 0x97,
	0xff, 0xde, 0x31, 0xc5, 0x80, 0xb0, 0x3e, 0xba, 0x26, 0x77, 0x84, 0xbd, 0xd2, 0xd5, 0x63, 0x30,
	0x19, 0x3c, 0x84, 0xa2, 0x9f, 0xff, 0x1a, 0xba, 0xf9, 0x59, 0xb4, 0xd0, 0xbf, 0x36, 0x6f, 0x0f,
	0xa1, 0xcb, 0xf9, 0x87, 0xd3, 0xa7, 0xe3, 0xaa, 0xde, 0x08, 0x69, 0x4e, 0x86, 0xdf, 0xb9, 0xff,
	0x08, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x78, 0xe4, 0x1b, 0x7f, 0x01, 0x00, 0x00,
}
