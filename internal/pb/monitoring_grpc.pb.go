// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/monitoring.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LoadAvg_LoadAvgGetMon_FullMethodName = "/pb.LoadAvg/LoadAvgGetMon"
)

// LoadAvgClient is the client API for LoadAvg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoadAvgClient interface {
	LoadAvgGetMon(ctx context.Context, in *LoadAvgRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LoadAvgReply], error)
}

type loadAvgClient struct {
	cc grpc.ClientConnInterface
}

func NewLoadAvgClient(cc grpc.ClientConnInterface) LoadAvgClient {
	return &loadAvgClient{cc}
}

func (c *loadAvgClient) LoadAvgGetMon(ctx context.Context, in *LoadAvgRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[LoadAvgReply], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LoadAvg_ServiceDesc.Streams[0], LoadAvg_LoadAvgGetMon_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[LoadAvgRequest, LoadAvgReply]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LoadAvg_LoadAvgGetMonClient = grpc.ServerStreamingClient[LoadAvgReply]

// LoadAvgServer is the server API for LoadAvg service.
// All implementations must embed UnimplementedLoadAvgServer
// for forward compatibility.
type LoadAvgServer interface {
	LoadAvgGetMon(*LoadAvgRequest, grpc.ServerStreamingServer[LoadAvgReply]) error
	mustEmbedUnimplementedLoadAvgServer()
}

// UnimplementedLoadAvgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLoadAvgServer struct{}

func (UnimplementedLoadAvgServer) LoadAvgGetMon(*LoadAvgRequest, grpc.ServerStreamingServer[LoadAvgReply]) error {
	return status.Errorf(codes.Unimplemented, "method LoadAvgGetMon not implemented")
}
func (UnimplementedLoadAvgServer) mustEmbedUnimplementedLoadAvgServer() {}
func (UnimplementedLoadAvgServer) testEmbeddedByValue()                 {}

// UnsafeLoadAvgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoadAvgServer will
// result in compilation errors.
type UnsafeLoadAvgServer interface {
	mustEmbedUnimplementedLoadAvgServer()
}

func RegisterLoadAvgServer(s grpc.ServiceRegistrar, srv LoadAvgServer) {
	// If the following call pancis, it indicates UnimplementedLoadAvgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LoadAvg_ServiceDesc, srv)
}

func _LoadAvg_LoadAvgGetMon_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LoadAvgRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LoadAvgServer).LoadAvgGetMon(m, &grpc.GenericServerStream[LoadAvgRequest, LoadAvgReply]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LoadAvg_LoadAvgGetMonServer = grpc.ServerStreamingServer[LoadAvgReply]

// LoadAvg_ServiceDesc is the grpc.ServiceDesc for LoadAvg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoadAvg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.LoadAvg",
	HandlerType: (*LoadAvgServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LoadAvgGetMon",
			Handler:       _LoadAvg_LoadAvgGetMon_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/monitoring.proto",
}
