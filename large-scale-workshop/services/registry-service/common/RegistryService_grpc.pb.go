// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
// source: RegistryService.proto

package RegistryService

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
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	RegistryService_Register_FullMethodName   = "/registryservice.RegistryService/Register"
	RegistryService_Unregister_FullMethodName = "/registryservice.RegistryService/Unregister"
	RegistryService_Discover_FullMethodName   = "/registryservice.RegistryService/Discover"
)

// RegistryServiceClient is the client API for RegistryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryServiceClient interface {
	Register(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Unregister(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Discover(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*ServiceNodes, error)
}

type registryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryServiceClient(cc grpc.ClientConnInterface) RegistryServiceClient {
	return &registryServiceClient{cc}
}

func (c *registryServiceClient) Register(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, RegistryService_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) Unregister(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, RegistryService_Unregister_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryServiceClient) Discover(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*ServiceNodes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServiceNodes)
	err := c.cc.Invoke(ctx, RegistryService_Discover_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryServiceServer is the server API for RegistryService service.
// All implementations must embed UnimplementedRegistryServiceServer
// for forward compatibility
type RegistryServiceServer interface {
	Register(context.Context, *ServiceRequest) (*empty.Empty, error)
	Unregister(context.Context, *ServiceRequest) (*empty.Empty, error)
	Discover(context.Context, *wrappers.StringValue) (*ServiceNodes, error)
	mustEmbedUnimplementedRegistryServiceServer()
}

// UnimplementedRegistryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryServiceServer struct {
}

func (UnimplementedRegistryServiceServer) Register(context.Context, *ServiceRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistryServiceServer) Unregister(context.Context, *ServiceRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
}
func (UnimplementedRegistryServiceServer) Discover(context.Context, *wrappers.StringValue) (*ServiceNodes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discover not implemented")
}
func (UnimplementedRegistryServiceServer) mustEmbedUnimplementedRegistryServiceServer() {}

// UnsafeRegistryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryServiceServer will
// result in compilation errors.
type UnsafeRegistryServiceServer interface {
	mustEmbedUnimplementedRegistryServiceServer()
}

func RegisterRegistryServiceServer(s grpc.ServiceRegistrar, srv RegistryServiceServer) {
	s.RegisterService(&RegistryService_ServiceDesc, srv)
}

func _RegistryService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Register(ctx, req.(*ServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Unregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryService_Unregister_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Unregister(ctx, req.(*ServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryService_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServiceServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryService_Discover_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServiceServer).Discover(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// RegistryService_ServiceDesc is the grpc.ServiceDesc for RegistryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegistryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "registryservice.RegistryService",
	HandlerType: (*RegistryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _RegistryService_Register_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _RegistryService_Unregister_Handler,
		},
		{
			MethodName: "Discover",
			Handler:    _RegistryService_Discover_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "RegistryService.proto",
}
