// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.1
// source: server_api_v1.proto

package server_api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserAPIClient is the client API for UserAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserAPIClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (UserAPI_UploadClient, error)
	GetInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FileList, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (UserAPI_DownloadClient, error)
}

type userAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewUserAPIClient(cc grpc.ClientConnInterface) UserAPIClient {
	return &userAPIClient{cc}
}

func (c *userAPIClient) Upload(ctx context.Context, opts ...grpc.CallOption) (UserAPI_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserAPI_ServiceDesc.Streams[0], "/user_api_v1.userAPI/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &userAPIUploadClient{stream}
	return x, nil
}

type UserAPI_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*wrapperspb.BoolValue, error)
	grpc.ClientStream
}

type userAPIUploadClient struct {
	grpc.ClientStream
}

func (x *userAPIUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *userAPIUploadClient) CloseAndRecv() (*wrapperspb.BoolValue, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(wrapperspb.BoolValue)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userAPIClient) GetInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FileList, error) {
	out := new(FileList)
	err := c.cc.Invoke(ctx, "/user_api_v1.userAPI/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAPIClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (UserAPI_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserAPI_ServiceDesc.Streams[1], "/user_api_v1.userAPI/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &userAPIDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserAPI_DownloadClient interface {
	Recv() (*DownloadResponse, error)
	grpc.ClientStream
}

type userAPIDownloadClient struct {
	grpc.ClientStream
}

func (x *userAPIDownloadClient) Recv() (*DownloadResponse, error) {
	m := new(DownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserAPIServer is the server API for UserAPI service.
// All implementations must embed UnimplementedUserAPIServer
// for forward compatibility
type UserAPIServer interface {
	Upload(UserAPI_UploadServer) error
	GetInfo(context.Context, *emptypb.Empty) (*FileList, error)
	Download(*DownloadRequest, UserAPI_DownloadServer) error
	mustEmbedUnimplementedUserAPIServer()
}

// UnimplementedUserAPIServer must be embedded to have forward compatible implementations.
type UnimplementedUserAPIServer struct {
}

func (UnimplementedUserAPIServer) Upload(UserAPI_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedUserAPIServer) GetInfo(context.Context, *emptypb.Empty) (*FileList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedUserAPIServer) Download(*DownloadRequest, UserAPI_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedUserAPIServer) mustEmbedUnimplementedUserAPIServer() {}

// UnsafeUserAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserAPIServer will
// result in compilation errors.
type UnsafeUserAPIServer interface {
	mustEmbedUnimplementedUserAPIServer()
}

func RegisterUserAPIServer(s grpc.ServiceRegistrar, srv UserAPIServer) {
	s.RegisterService(&UserAPI_ServiceDesc, srv)
}

func _UserAPI_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UserAPIServer).Upload(&userAPIUploadServer{stream})
}

type UserAPI_UploadServer interface {
	SendAndClose(*wrapperspb.BoolValue) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type userAPIUploadServer struct {
	grpc.ServerStream
}

func (x *userAPIUploadServer) SendAndClose(m *wrapperspb.BoolValue) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userAPIUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _UserAPI_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAPIServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_api_v1.userAPI/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAPIServer).GetInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAPI_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserAPIServer).Download(m, &userAPIDownloadServer{stream})
}

type UserAPI_DownloadServer interface {
	Send(*DownloadResponse) error
	grpc.ServerStream
}

type userAPIDownloadServer struct {
	grpc.ServerStream
}

func (x *userAPIDownloadServer) Send(m *DownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

// UserAPI_ServiceDesc is the grpc.ServiceDesc for UserAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_api_v1.userAPI",
	HandlerType: (*UserAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _UserAPI_GetInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _UserAPI_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _UserAPI_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server_api_v1.proto",
}
