// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: TestService.proto

package common

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TestService_HelloWorld_FullMethodName          = "/testservice.TestService/HelloWorld"
	TestService_HelloToUser_FullMethodName         = "/testservice.TestService/HelloToUser"
	TestService_Store_FullMethodName               = "/testservice.TestService/Store"
	TestService_Get_FullMethodName                 = "/testservice.TestService/Get"
	TestService_WaitAndRand_FullMethodName         = "/testservice.TestService/WaitAndRand"
	TestService_ExtractLinksFromURL_FullMethodName = "/testservice.TestService/ExtractLinksFromURL"
	TestService_IsAlive_FullMethodName             = "/testservice.TestService/IsAlive"
)

// TestServiceClient is the client API for TestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestServiceClient interface {
	// returns "Hello World"
	HelloWorld(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	// receives user name, return "Hello [user name]"
	HelloToUser(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	// receives key/value pair and stores it in a map
	Store(ctx context.Context, in *StoreKeyValue, opts ...grpc.CallOption) (*empty.Empty, error)
	// returns value for a given key from the map
	Get(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	// Wait given number of seconds and return random number
	// async function
	WaitAndRand(ctx context.Context, in *wrappers.Int32Value, opts ...grpc.CallOption) (TestService_WaitAndRandClient, error)
	// extracts links from URL using beautiful soup
	ExtractLinksFromURL(ctx context.Context, in *ExtractLinksFromURLParameters, opts ...grpc.CallOption) (*ExtractLinksFromURLReturnedValue, error)
	// returns true
	IsAlive(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
}

type testServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTestServiceClient(cc grpc.ClientConnInterface) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) HelloWorld(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, TestService_HelloWorld_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) HelloToUser(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, TestService_HelloToUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Store(ctx context.Context, in *StoreKeyValue, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, TestService_Store_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Get(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, TestService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) WaitAndRand(ctx context.Context, in *wrappers.Int32Value, opts ...grpc.CallOption) (TestService_WaitAndRandClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestService_ServiceDesc.Streams[0], TestService_WaitAndRand_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &testServiceWaitAndRandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestService_WaitAndRandClient interface {
	Recv() (*wrappers.Int32Value, error)
	grpc.ClientStream
}

type testServiceWaitAndRandClient struct {
	grpc.ClientStream
}

func (x *testServiceWaitAndRandClient) Recv() (*wrappers.Int32Value, error) {
	m := new(wrappers.Int32Value)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testServiceClient) ExtractLinksFromURL(ctx context.Context, in *ExtractLinksFromURLParameters, opts ...grpc.CallOption) (*ExtractLinksFromURLReturnedValue, error) {
	out := new(ExtractLinksFromURLReturnedValue)
	err := c.cc.Invoke(ctx, TestService_ExtractLinksFromURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) IsAlive(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, TestService_IsAlive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestServiceServer is the server API for TestService service.
// All implementations must embed UnimplementedTestServiceServer
// for forward compatibility
type TestServiceServer interface {
	// returns "Hello World"
	HelloWorld(context.Context, *empty.Empty) (*wrappers.StringValue, error)
	// receives user name, return "Hello [user name]"
	HelloToUser(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error)
	// receives key/value pair and stores it in a map
	Store(context.Context, *StoreKeyValue) (*empty.Empty, error)
	// returns value for a given key from the map
	Get(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error)
	// Wait given number of seconds and return random number
	// async function
	WaitAndRand(*wrappers.Int32Value, TestService_WaitAndRandServer) error
	// extracts links from URL using beautiful soup
	ExtractLinksFromURL(context.Context, *ExtractLinksFromURLParameters) (*ExtractLinksFromURLReturnedValue, error)
	// returns true
	IsAlive(context.Context, *empty.Empty) (*wrappers.BoolValue, error)
	mustEmbedUnimplementedTestServiceServer()
}

// UnimplementedTestServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTestServiceServer struct {
}

func (UnimplementedTestServiceServer) HelloWorld(context.Context, *empty.Empty) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloWorld not implemented")
}
func (UnimplementedTestServiceServer) HelloToUser(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloToUser not implemented")
}
func (UnimplementedTestServiceServer) Store(context.Context, *StoreKeyValue) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (UnimplementedTestServiceServer) Get(context.Context, *wrappers.StringValue) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedTestServiceServer) WaitAndRand(*wrappers.Int32Value, TestService_WaitAndRandServer) error {
	return status.Errorf(codes.Unimplemented, "method WaitAndRand not implemented")
}
func (UnimplementedTestServiceServer) ExtractLinksFromURL(context.Context, *ExtractLinksFromURLParameters) (*ExtractLinksFromURLReturnedValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtractLinksFromURL not implemented")
}
func (UnimplementedTestServiceServer) IsAlive(context.Context, *empty.Empty) (*wrappers.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAlive not implemented")
}
func (UnimplementedTestServiceServer) mustEmbedUnimplementedTestServiceServer() {}

// UnsafeTestServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestServiceServer will
// result in compilation errors.
type UnsafeTestServiceServer interface {
	mustEmbedUnimplementedTestServiceServer()
}

func RegisterTestServiceServer(s grpc.ServiceRegistrar, srv TestServiceServer) {
	s.RegisterService(&TestService_ServiceDesc, srv)
}

func _TestService_HelloWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).HelloWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_HelloWorld_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).HelloWorld(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_HelloToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).HelloToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_HelloToUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).HelloToUser(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreKeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_Store_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Store(ctx, req.(*StoreKeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Get(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_WaitAndRand_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(wrappers.Int32Value)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServiceServer).WaitAndRand(m, &testServiceWaitAndRandServer{stream})
}

type TestService_WaitAndRandServer interface {
	Send(*wrappers.Int32Value) error
	grpc.ServerStream
}

type testServiceWaitAndRandServer struct {
	grpc.ServerStream
}

func (x *testServiceWaitAndRandServer) Send(m *wrappers.Int32Value) error {
	return x.ServerStream.SendMsg(m)
}

func _TestService_ExtractLinksFromURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtractLinksFromURLParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).ExtractLinksFromURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_ExtractLinksFromURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).ExtractLinksFromURL(ctx, req.(*ExtractLinksFromURLParameters))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_IsAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).IsAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_IsAlive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).IsAlive(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// TestService_ServiceDesc is the grpc.ServiceDesc for TestService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "testservice.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HelloWorld",
			Handler:    _TestService_HelloWorld_Handler,
		},
		{
			MethodName: "HelloToUser",
			Handler:    _TestService_HelloToUser_Handler,
		},
		{
			MethodName: "Store",
			Handler:    _TestService_Store_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _TestService_Get_Handler,
		},
		{
			MethodName: "ExtractLinksFromURL",
			Handler:    _TestService_ExtractLinksFromURL_Handler,
		},
		{
			MethodName: "IsAlive",
			Handler:    _TestService_IsAlive_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WaitAndRand",
			Handler:       _TestService_WaitAndRand_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "TestService.proto",
}
