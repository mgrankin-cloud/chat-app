// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: message/message.proto

package ssov4

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
	Chat_GetChat_FullMethodName    = "/message.Chat/GetChat"
	Chat_UpdateChat_FullMethodName = "/message.Chat/UpdateChat"
	Chat_DeleteChat_FullMethodName = "/message.Chat/DeleteChat"
	Chat_MuteChat_FullMethodName   = "/message.Chat/MuteChat"
	Chat_PinChat_FullMethodName    = "/message.Chat/PinChat"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	GetChat(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*GetChatResponse, error)
	UpdateChat(ctx context.Context, in *UpdateChatRequest, opts ...grpc.CallOption) (*UpdateChatResponse, error)
	DeleteChat(ctx context.Context, in *DeleteChatRequest, opts ...grpc.CallOption) (*DeleteChatResponse, error)
	MuteChat(ctx context.Context, in *MuteChatRequest, opts ...grpc.CallOption) (*MuteChatResponse, error)
	PinChat(ctx context.Context, in *PinChatRequest, opts ...grpc.CallOption) (*PinChatResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) GetChat(ctx context.Context, in *GetChatRequest, opts ...grpc.CallOption) (*GetChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChatResponse)
	err := c.cc.Invoke(ctx, Chat_GetChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) UpdateChat(ctx context.Context, in *UpdateChatRequest, opts ...grpc.CallOption) (*UpdateChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateChatResponse)
	err := c.cc.Invoke(ctx, Chat_UpdateChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) DeleteChat(ctx context.Context, in *DeleteChatRequest, opts ...grpc.CallOption) (*DeleteChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteChatResponse)
	err := c.cc.Invoke(ctx, Chat_DeleteChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) MuteChat(ctx context.Context, in *MuteChatRequest, opts ...grpc.CallOption) (*MuteChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MuteChatResponse)
	err := c.cc.Invoke(ctx, Chat_MuteChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) PinChat(ctx context.Context, in *PinChatRequest, opts ...grpc.CallOption) (*PinChatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PinChatResponse)
	err := c.cc.Invoke(ctx, Chat_PinChat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility.
type ChatServer interface {
	GetChat(context.Context, *GetChatRequest) (*GetChatResponse, error)
	UpdateChat(context.Context, *UpdateChatRequest) (*UpdateChatResponse, error)
	DeleteChat(context.Context, *DeleteChatRequest) (*DeleteChatResponse, error)
	MuteChat(context.Context, *MuteChatRequest) (*MuteChatResponse, error)
	PinChat(context.Context, *PinChatRequest) (*PinChatResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatServer struct{}

func (UnimplementedChatServer) GetChat(context.Context, *GetChatRequest) (*GetChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChat not implemented")
}
func (UnimplementedChatServer) UpdateChat(context.Context, *UpdateChatRequest) (*UpdateChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChat not implemented")
}
func (UnimplementedChatServer) DeleteChat(context.Context, *DeleteChatRequest) (*DeleteChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChat not implemented")
}
func (UnimplementedChatServer) MuteChat(context.Context, *MuteChatRequest) (*MuteChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MuteChat not implemented")
}
func (UnimplementedChatServer) PinChat(context.Context, *PinChatRequest) (*PinChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PinChat not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}
func (UnimplementedChatServer) testEmbeddedByValue()              {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	// If the following call pancis, it indicates UnimplementedChatServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_GetChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_GetChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetChat(ctx, req.(*GetChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_UpdateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).UpdateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_UpdateChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).UpdateChat(ctx, req.(*UpdateChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_DeleteChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).DeleteChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_DeleteChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).DeleteChat(ctx, req.(*DeleteChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_MuteChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MuteChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).MuteChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_MuteChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).MuteChat(ctx, req.(*MuteChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_PinChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PinChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).PinChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_PinChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).PinChat(ctx, req.(*PinChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChat",
			Handler:    _Chat_GetChat_Handler,
		},
		{
			MethodName: "UpdateChat",
			Handler:    _Chat_UpdateChat_Handler,
		},
		{
			MethodName: "DeleteChat",
			Handler:    _Chat_DeleteChat_Handler,
		},
		{
			MethodName: "MuteChat",
			Handler:    _Chat_MuteChat_Handler,
		},
		{
			MethodName: "PinChat",
			Handler:    _Chat_PinChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message/message.proto",
}
