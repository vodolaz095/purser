// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: purser.proto

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

// PurserClient is the client API for Purser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PurserClient interface {
	GetSecretByID(ctx context.Context, in *SecretByIDRequest, opts ...grpc.CallOption) (*Secret, error)
	DeleteSecretByID(ctx context.Context, in *SecretByIDRequest, opts ...grpc.CallOption) (*Nothing, error)
	CreateSecret(ctx context.Context, in *NewSecretRequest, opts ...grpc.CallOption) (*Secret, error)
}

type purserClient struct {
	cc grpc.ClientConnInterface
}

func NewPurserClient(cc grpc.ClientConnInterface) PurserClient {
	return &purserClient{cc}
}

func (c *purserClient) GetSecretByID(ctx context.Context, in *SecretByIDRequest, opts ...grpc.CallOption) (*Secret, error) {
	out := new(Secret)
	err := c.cc.Invoke(ctx, "/purser.Purser/GetSecretByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purserClient) DeleteSecretByID(ctx context.Context, in *SecretByIDRequest, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/purser.Purser/DeleteSecretByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purserClient) CreateSecret(ctx context.Context, in *NewSecretRequest, opts ...grpc.CallOption) (*Secret, error) {
	out := new(Secret)
	err := c.cc.Invoke(ctx, "/purser.Purser/CreateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PurserServer is the server API for Purser service.
// All implementations must embed UnimplementedPurserServer
// for forward compatibility
type PurserServer interface {
	GetSecretByID(context.Context, *SecretByIDRequest) (*Secret, error)
	DeleteSecretByID(context.Context, *SecretByIDRequest) (*Nothing, error)
	CreateSecret(context.Context, *NewSecretRequest) (*Secret, error)
	mustEmbedUnimplementedPurserServer()
}

// UnimplementedPurserServer must be embedded to have forward compatible implementations.
type UnimplementedPurserServer struct {
}

func (UnimplementedPurserServer) GetSecretByID(context.Context, *SecretByIDRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecretByID not implemented")
}
func (UnimplementedPurserServer) DeleteSecretByID(context.Context, *SecretByIDRequest) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecretByID not implemented")
}
func (UnimplementedPurserServer) CreateSecret(context.Context, *NewSecretRequest) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSecret not implemented")
}
func (UnimplementedPurserServer) mustEmbedUnimplementedPurserServer() {}

// UnsafePurserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PurserServer will
// result in compilation errors.
type UnsafePurserServer interface {
	mustEmbedUnimplementedPurserServer()
}

func RegisterPurserServer(s grpc.ServiceRegistrar, srv PurserServer) {
	s.RegisterService(&Purser_ServiceDesc, srv)
}

func _Purser_GetSecretByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurserServer).GetSecretByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/purser.Purser/GetSecretByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurserServer).GetSecretByID(ctx, req.(*SecretByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Purser_DeleteSecretByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurserServer).DeleteSecretByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/purser.Purser/DeleteSecretByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurserServer).DeleteSecretByID(ctx, req.(*SecretByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Purser_CreateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurserServer).CreateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/purser.Purser/CreateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurserServer).CreateSecret(ctx, req.(*NewSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Purser_ServiceDesc is the grpc.ServiceDesc for Purser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Purser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "purser.Purser",
	HandlerType: (*PurserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSecretByID",
			Handler:    _Purser_GetSecretByID_Handler,
		},
		{
			MethodName: "DeleteSecretByID",
			Handler:    _Purser_DeleteSecretByID_Handler,
		},
		{
			MethodName: "CreateSecret",
			Handler:    _Purser_CreateSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "purser.proto",
}
