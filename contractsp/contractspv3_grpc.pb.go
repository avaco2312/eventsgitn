// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package contractsp

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

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	SearchId(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Event, error)
	SearchName(ctx context.Context, in *Name, opts ...grpc.CallOption) (*Event, error)
	SearchAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ArrayEvent, error)
	AddEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Id, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) SearchId(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/contractsp.Events/SearchId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) SearchName(ctx context.Context, in *Name, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/contractsp.Events/SearchName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) SearchAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ArrayEvent, error) {
	out := new(ArrayEvent)
	err := c.cc.Invoke(ctx, "/contractsp.Events/SearchAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) AddEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Id, error) {
	out := new(Id)
	err := c.cc.Invoke(ctx, "/contractsp.Events/AddEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	SearchId(context.Context, *Id) (*Event, error)
	SearchName(context.Context, *Name) (*Event, error)
	SearchAll(context.Context, *Empty) (*ArrayEvent, error)
	AddEvent(context.Context, *Event) (*Id, error)
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) SearchId(context.Context, *Id) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchId not implemented")
}
func (UnimplementedEventsServer) SearchName(context.Context, *Name) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchName not implemented")
}
func (UnimplementedEventsServer) SearchAll(context.Context, *Empty) (*ArrayEvent, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAll not implemented")
}
func (UnimplementedEventsServer) AddEvent(context.Context, *Event) (*Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEvent not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&Events_ServiceDesc, srv)
}

func _Events_SearchId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).SearchId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contractsp.Events/SearchId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).SearchId(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_SearchName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Name)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).SearchName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contractsp.Events/SearchName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).SearchName(ctx, req.(*Name))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_SearchAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).SearchAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contractsp.Events/SearchAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).SearchAll(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_AddEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).AddEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/contractsp.Events/AddEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).AddEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// Events_ServiceDesc is the grpc.ServiceDesc for Events service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Events_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "contractsp.Events",
	HandlerType: (*EventsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchId",
			Handler:    _Events_SearchId_Handler,
		},
		{
			MethodName: "SearchName",
			Handler:    _Events_SearchName_Handler,
		},
		{
			MethodName: "SearchAll",
			Handler:    _Events_SearchAll_Handler,
		},
		{
			MethodName: "AddEvent",
			Handler:    _Events_AddEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contractspv3.proto",
}