// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: media/media.proto

package ssov6

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
	Media_UploadMedia_FullMethodName   = "/media.Media/UploadMedia"
	Media_DownloadMedia_FullMethodName = "/media.Media/DownloadMedia"
)

// MediaClient is the client API for Media service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MediaClient interface {
	UploadMedia(ctx context.Context, in *UploadMediaRequest, opts ...grpc.CallOption) (*UploadMediaResponse, error)
	DownloadMedia(ctx context.Context, in *DownloadMediaRequest, opts ...grpc.CallOption) (*DownloadMediaResponse, error)
}

type mediaClient struct {
	cc grpc.ClientConnInterface
}

func NewMediaClient(cc grpc.ClientConnInterface) MediaClient {
	return &mediaClient{cc}
}

func (c *mediaClient) UploadMedia(ctx context.Context, in *UploadMediaRequest, opts ...grpc.CallOption) (*UploadMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadMediaResponse)
	err := c.cc.Invoke(ctx, Media_UploadMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mediaClient) DownloadMedia(ctx context.Context, in *DownloadMediaRequest, opts ...grpc.CallOption) (*DownloadMediaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DownloadMediaResponse)
	err := c.cc.Invoke(ctx, Media_DownloadMedia_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MediaServer is the server API for Media service.
// All implementations must embed UnimplementedMediaServer
// for forward compatibility.
type MediaServer interface {
	UploadMedia(context.Context, *UploadMediaRequest) (*UploadMediaResponse, error)
	DownloadMedia(context.Context, *DownloadMediaRequest) (*DownloadMediaResponse, error)
	mustEmbedUnimplementedMediaServer()
}

// UnimplementedMediaServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMediaServer struct{}

func (UnimplementedMediaServer) UploadMedia(context.Context, *UploadMediaRequest) (*UploadMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadMedia not implemented")
}
func (UnimplementedMediaServer) DownloadMedia(context.Context, *DownloadMediaRequest) (*DownloadMediaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadMedia not implemented")
}
func (UnimplementedMediaServer) mustEmbedUnimplementedMediaServer() {}
func (UnimplementedMediaServer) testEmbeddedByValue()               {}

// UnsafeMediaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MediaServer will
// result in compilation errors.
type UnsafeMediaServer interface {
	mustEmbedUnimplementedMediaServer()
}

func RegisterMediaServer(s grpc.ServiceRegistrar, srv MediaServer) {
	// If the following call pancis, it indicates UnimplementedMediaServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Media_ServiceDesc, srv)
}

func _Media_UploadMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServer).UploadMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Media_UploadMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServer).UploadMedia(ctx, req.(*UploadMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Media_DownloadMedia_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadMediaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MediaServer).DownloadMedia(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Media_DownloadMedia_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MediaServer).DownloadMedia(ctx, req.(*DownloadMediaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Media_ServiceDesc is the grpc.ServiceDesc for Media service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Media_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "media.Media",
	HandlerType: (*MediaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadMedia",
			Handler:    _Media_UploadMedia_Handler,
		},
		{
			MethodName: "DownloadMedia",
			Handler:    _Media_DownloadMedia_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "media/media.proto",
}