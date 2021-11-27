// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamServiceClient interface {
	EchoStr(ctx context.Context, in *Empty, opts ...grpc.CallOption) (StreamService_EchoStrClient, error)
	EchoReq(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error)
}

type streamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamServiceClient(cc grpc.ClientConnInterface) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) EchoStr(ctx context.Context, in *Empty, opts ...grpc.CallOption) (StreamService_EchoStrClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[0], "/stream.StreamService/EchoStr", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceEchoStrClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_EchoStrClient interface {
	Recv() (*Result, error)
	grpc.ClientStream
}

type streamServiceEchoStrClient struct {
	grpc.ClientStream
}

func (x *streamServiceEchoStrClient) Recv() (*Result, error) {
	m := new(Result)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamServiceClient) EchoReq(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/stream.StreamService/EchoReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamServiceServer is the server API for StreamService service.
// All implementations must embed UnimplementedStreamServiceServer
// for forward compatibility
type StreamServiceServer interface {
	EchoStr(*Empty, StreamService_EchoStrServer) error
	EchoReq(context.Context, *Empty) (*Result, error)
	mustEmbedUnimplementedStreamServiceServer()
}

// UnimplementedStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServiceServer struct {
}

func (UnimplementedStreamServiceServer) EchoStr(*Empty, StreamService_EchoStrServer) error {
	return status.Errorf(codes.Unimplemented, "method EchoStr not implemented")
}
func (UnimplementedStreamServiceServer) EchoReq(context.Context, *Empty) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EchoReq not implemented")
}
func (UnimplementedStreamServiceServer) mustEmbedUnimplementedStreamServiceServer() {}

// UnsafeStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServiceServer will
// result in compilation errors.
type UnsafeStreamServiceServer interface {
	mustEmbedUnimplementedStreamServiceServer()
}

func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {
	s.RegisterService(&StreamService_ServiceDesc, srv)
}

func _StreamService_EchoStr_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).EchoStr(m, &streamServiceEchoStrServer{stream})
}

type StreamService_EchoStrServer interface {
	Send(*Result) error
	grpc.ServerStream
}

type streamServiceEchoStrServer struct {
	grpc.ServerStream
}

func (x *streamServiceEchoStrServer) Send(m *Result) error {
	return x.ServerStream.SendMsg(m)
}

func _StreamService_EchoReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServiceServer).EchoReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stream.StreamService/EchoReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServiceServer).EchoReq(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StreamService_ServiceDesc is the grpc.ServiceDesc for StreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stream.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EchoReq",
			Handler:    _StreamService_EchoReq_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EchoStr",
			Handler:       _StreamService_EchoStr_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/stream.proto",
}
