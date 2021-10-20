// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package intstore

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

// IntStoreClient is the client API for IntStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IntStoreClient interface {
	Set(ctx context.Context, in *Item, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Item, error)
	// A server-to-client streaming RPC.
	ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (IntStore_ListItemsClient, error)
	// A client-to-server streaming RPC.
	SetStream(ctx context.Context, opts ...grpc.CallOption) (IntStore_SetStreamClient, error)
	// A Bidirectional streaming RPC.
	StreamChat(ctx context.Context, opts ...grpc.CallOption) (IntStore_StreamChatClient, error)
}

type intStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewIntStoreClient(cc grpc.ClientConnInterface) IntStoreClient {
	return &intStoreClient{cc}
}

func (c *intStoreClient) Set(ctx context.Context, in *Item, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/intstore.IntStore/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *intStoreClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Item, error) {
	out := new(Item)
	err := c.cc.Invoke(ctx, "/intstore.IntStore/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *intStoreClient) ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (IntStore_ListItemsClient, error) {
	stream, err := c.cc.NewStream(ctx, &IntStore_ServiceDesc.Streams[0], "/intstore.IntStore/ListItems", opts...)
	if err != nil {
		return nil, err
	}
	x := &intStoreListItemsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type IntStore_ListItemsClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type intStoreListItemsClient struct {
	grpc.ClientStream
}

func (x *intStoreListItemsClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *intStoreClient) SetStream(ctx context.Context, opts ...grpc.CallOption) (IntStore_SetStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &IntStore_ServiceDesc.Streams[1], "/intstore.IntStore/SetStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &intStoreSetStreamClient{stream}
	return x, nil
}

type IntStore_SetStreamClient interface {
	Send(*Item) error
	CloseAndRecv() (*Summary, error)
	grpc.ClientStream
}

type intStoreSetStreamClient struct {
	grpc.ClientStream
}

func (x *intStoreSetStreamClient) Send(m *Item) error {
	return x.ClientStream.SendMsg(m)
}

func (x *intStoreSetStreamClient) CloseAndRecv() (*Summary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Summary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *intStoreClient) StreamChat(ctx context.Context, opts ...grpc.CallOption) (IntStore_StreamChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &IntStore_ServiceDesc.Streams[2], "/intstore.IntStore/StreamChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &intStoreStreamChatClient{stream}
	return x, nil
}

type IntStore_StreamChatClient interface {
	Send(*Item) error
	Recv() (*Item, error)
	grpc.ClientStream
}

type intStoreStreamChatClient struct {
	grpc.ClientStream
}

func (x *intStoreStreamChatClient) Send(m *Item) error {
	return x.ClientStream.SendMsg(m)
}

func (x *intStoreStreamChatClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// IntStoreServer is the server API for IntStore service.
// All implementations must embed UnimplementedIntStoreServer
// for forward compatibility
type IntStoreServer interface {
	Set(context.Context, *Item) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*Item, error)
	// A server-to-client streaming RPC.
	ListItems(*ListItemsRequest, IntStore_ListItemsServer) error
	// A client-to-server streaming RPC.
	SetStream(IntStore_SetStreamServer) error
	// A Bidirectional streaming RPC.
	StreamChat(IntStore_StreamChatServer) error
	mustEmbedUnimplementedIntStoreServer()
}

// UnimplementedIntStoreServer must be embedded to have forward compatible implementations.
type UnimplementedIntStoreServer struct {
}

func (UnimplementedIntStoreServer) Set(context.Context, *Item) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedIntStoreServer) Get(context.Context, *GetRequest) (*Item, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedIntStoreServer) ListItems(*ListItemsRequest, IntStore_ListItemsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListItems not implemented")
}
func (UnimplementedIntStoreServer) SetStream(IntStore_SetStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SetStream not implemented")
}
func (UnimplementedIntStoreServer) StreamChat(IntStore_StreamChatServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamChat not implemented")
}
func (UnimplementedIntStoreServer) mustEmbedUnimplementedIntStoreServer() {}

// UnsafeIntStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IntStoreServer will
// result in compilation errors.
type UnsafeIntStoreServer interface {
	mustEmbedUnimplementedIntStoreServer()
}

func RegisterIntStoreServer(s grpc.ServiceRegistrar, srv IntStoreServer) {
	s.RegisterService(&IntStore_ServiceDesc, srv)
}

func _IntStore_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Item)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntStoreServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/intstore.IntStore/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntStoreServer).Set(ctx, req.(*Item))
	}
	return interceptor(ctx, in, info, handler)
}

func _IntStore_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IntStoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/intstore.IntStore/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IntStoreServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IntStore_ListItems_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListItemsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(IntStoreServer).ListItems(m, &intStoreListItemsServer{stream})
}

type IntStore_ListItemsServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type intStoreListItemsServer struct {
	grpc.ServerStream
}

func (x *intStoreListItemsServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _IntStore_SetStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IntStoreServer).SetStream(&intStoreSetStreamServer{stream})
}

type IntStore_SetStreamServer interface {
	SendAndClose(*Summary) error
	Recv() (*Item, error)
	grpc.ServerStream
}

type intStoreSetStreamServer struct {
	grpc.ServerStream
}

func (x *intStoreSetStreamServer) SendAndClose(m *Summary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *intStoreSetStreamServer) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _IntStore_StreamChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IntStoreServer).StreamChat(&intStoreStreamChatServer{stream})
}

type IntStore_StreamChatServer interface {
	Send(*Item) error
	Recv() (*Item, error)
	grpc.ServerStream
}

type intStoreStreamChatServer struct {
	grpc.ServerStream
}

func (x *intStoreStreamChatServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func (x *intStoreStreamChatServer) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// IntStore_ServiceDesc is the grpc.ServiceDesc for IntStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IntStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "intstore.IntStore",
	HandlerType: (*IntStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _IntStore_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _IntStore_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListItems",
			Handler:       _IntStore_ListItems_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SetStream",
			Handler:       _IntStore_SetStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamChat",
			Handler:       _IntStore_StreamChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "intstore.proto",
}
