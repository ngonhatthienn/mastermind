// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: gameservice.proto

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

const (
	Services_CreateGame_FullMethodName     = "/gameservice.Services/CreateGame"
	Services_ListGame_FullMethodName       = "/gameservice.Services/ListGame"
	Services_GetCurrent_FullMethodName     = "/gameservice.Services/GetCurrent"
	Services_PickGame_FullMethodName       = "/gameservice.Services/PickGame"
	Services_PlayGame_FullMethodName       = "/gameservice.Services/PlayGame"
	Services_UpdateGame_FullMethodName     = "/gameservice.Services/UpdateGame"
	Services_HintGame_FullMethodName       = "/gameservice.Services/HintGame"
	Services_CreateUser_FullMethodName     = "/gameservice.Services/CreateUser"
	Services_GetListUser_FullMethodName    = "/gameservice.Services/GetListUser"
	Services_GetLeaderBoard_FullMethodName = "/gameservice.Services/GetLeaderBoard"
	Services_LogIn_FullMethodName          = "/gameservice.Services/LogIn"
)

// ServicesClient is the client API for Services service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesClient interface {
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameReply, error)
	ListGame(ctx context.Context, in *ListGameRequest, opts ...grpc.CallOption) (*ListGameReply, error)
	GetCurrent(ctx context.Context, in *CurrentGameRequest, opts ...grpc.CallOption) (*CurrentGameReply, error)
	PickGame(ctx context.Context, in *PickGameRequest, opts ...grpc.CallOption) (*PickGameReply, error)
	PlayGame(ctx context.Context, in *PlayGameRequest, opts ...grpc.CallOption) (*PlayGameReply, error)
	UpdateGame(ctx context.Context, in *UpdateGameRequest, opts ...grpc.CallOption) (*UpdateGameReply, error)
	HintGame(ctx context.Context, in *HintGameRequest, opts ...grpc.CallOption) (*HintGameReply, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error)
	GetListUser(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserReply, error)
	GetLeaderBoard(ctx context.Context, in *LeaderBoardRequest, opts ...grpc.CallOption) (*LeaderBoardReply, error)
	LogIn(ctx context.Context, in *LogInRequest, opts ...grpc.CallOption) (*LogInReply, error)
}

type servicesClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesClient(cc grpc.ClientConnInterface) ServicesClient {
	return &servicesClient{cc}
}

func (c *servicesClient) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameReply, error) {
	out := new(CreateGameReply)
	err := c.cc.Invoke(ctx, Services_CreateGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) ListGame(ctx context.Context, in *ListGameRequest, opts ...grpc.CallOption) (*ListGameReply, error) {
	out := new(ListGameReply)
	err := c.cc.Invoke(ctx, Services_ListGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetCurrent(ctx context.Context, in *CurrentGameRequest, opts ...grpc.CallOption) (*CurrentGameReply, error) {
	out := new(CurrentGameReply)
	err := c.cc.Invoke(ctx, Services_GetCurrent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) PickGame(ctx context.Context, in *PickGameRequest, opts ...grpc.CallOption) (*PickGameReply, error) {
	out := new(PickGameReply)
	err := c.cc.Invoke(ctx, Services_PickGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) PlayGame(ctx context.Context, in *PlayGameRequest, opts ...grpc.CallOption) (*PlayGameReply, error) {
	out := new(PlayGameReply)
	err := c.cc.Invoke(ctx, Services_PlayGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) UpdateGame(ctx context.Context, in *UpdateGameRequest, opts ...grpc.CallOption) (*UpdateGameReply, error) {
	out := new(UpdateGameReply)
	err := c.cc.Invoke(ctx, Services_UpdateGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) HintGame(ctx context.Context, in *HintGameRequest, opts ...grpc.CallOption) (*HintGameReply, error) {
	out := new(HintGameReply)
	err := c.cc.Invoke(ctx, Services_HintGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error) {
	out := new(CreateUserReply)
	err := c.cc.Invoke(ctx, Services_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetListUser(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserReply, error) {
	out := new(ListUserReply)
	err := c.cc.Invoke(ctx, Services_GetListUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetLeaderBoard(ctx context.Context, in *LeaderBoardRequest, opts ...grpc.CallOption) (*LeaderBoardReply, error) {
	out := new(LeaderBoardReply)
	err := c.cc.Invoke(ctx, Services_GetLeaderBoard_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) LogIn(ctx context.Context, in *LogInRequest, opts ...grpc.CallOption) (*LogInReply, error) {
	out := new(LogInReply)
	err := c.cc.Invoke(ctx, Services_LogIn_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesServer is the server API for Services service.
// All implementations must embed UnimplementedServicesServer
// for forward compatibility
type ServicesServer interface {
	CreateGame(context.Context, *CreateGameRequest) (*CreateGameReply, error)
	ListGame(context.Context, *ListGameRequest) (*ListGameReply, error)
	GetCurrent(context.Context, *CurrentGameRequest) (*CurrentGameReply, error)
	PickGame(context.Context, *PickGameRequest) (*PickGameReply, error)
	PlayGame(context.Context, *PlayGameRequest) (*PlayGameReply, error)
	UpdateGame(context.Context, *UpdateGameRequest) (*UpdateGameReply, error)
	HintGame(context.Context, *HintGameRequest) (*HintGameReply, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error)
	GetListUser(context.Context, *ListUserRequest) (*ListUserReply, error)
	GetLeaderBoard(context.Context, *LeaderBoardRequest) (*LeaderBoardReply, error)
	LogIn(context.Context, *LogInRequest) (*LogInReply, error)
	mustEmbedUnimplementedServicesServer()
}

// UnimplementedServicesServer must be embedded to have forward compatible implementations.
type UnimplementedServicesServer struct {
}

func (UnimplementedServicesServer) CreateGame(context.Context, *CreateGameRequest) (*CreateGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGame not implemented")
}
func (UnimplementedServicesServer) ListGame(context.Context, *ListGameRequest) (*ListGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGame not implemented")
}
func (UnimplementedServicesServer) GetCurrent(context.Context, *CurrentGameRequest) (*CurrentGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrent not implemented")
}
func (UnimplementedServicesServer) PickGame(context.Context, *PickGameRequest) (*PickGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PickGame not implemented")
}
func (UnimplementedServicesServer) PlayGame(context.Context, *PlayGameRequest) (*PlayGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlayGame not implemented")
}
func (UnimplementedServicesServer) UpdateGame(context.Context, *UpdateGameRequest) (*UpdateGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGame not implemented")
}
func (UnimplementedServicesServer) HintGame(context.Context, *HintGameRequest) (*HintGameReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HintGame not implemented")
}
func (UnimplementedServicesServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedServicesServer) GetListUser(context.Context, *ListUserRequest) (*ListUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListUser not implemented")
}
func (UnimplementedServicesServer) GetLeaderBoard(context.Context, *LeaderBoardRequest) (*LeaderBoardReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderBoard not implemented")
}
func (UnimplementedServicesServer) LogIn(context.Context, *LogInRequest) (*LogInReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogIn not implemented")
}
func (UnimplementedServicesServer) mustEmbedUnimplementedServicesServer() {}

// UnsafeServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesServer will
// result in compilation errors.
type UnsafeServicesServer interface {
	mustEmbedUnimplementedServicesServer()
}

func RegisterServicesServer(s grpc.ServiceRegistrar, srv ServicesServer) {
	s.RegisterService(&Services_ServiceDesc, srv)
}

func _Services_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_CreateGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).CreateGame(ctx, req.(*CreateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_ListGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).ListGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_ListGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).ListGame(ctx, req.(*ListGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetCurrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrentGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetCurrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_GetCurrent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetCurrent(ctx, req.(*CurrentGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_PickGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PickGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).PickGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_PickGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).PickGame(ctx, req.(*PickGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_PlayGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).PlayGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_PlayGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).PlayGame(ctx, req.(*PlayGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_UpdateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).UpdateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_UpdateGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).UpdateGame(ctx, req.(*UpdateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_HintGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HintGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).HintGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_HintGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).HintGame(ctx, req.(*HintGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetListUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetListUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_GetListUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetListUser(ctx, req.(*ListUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetLeaderBoard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderBoardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetLeaderBoard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_GetLeaderBoard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetLeaderBoard(ctx, req.(*LeaderBoardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_LogIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).LogIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Services_LogIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).LogIn(ctx, req.(*LogInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Services_ServiceDesc is the grpc.ServiceDesc for Services service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Services_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gameservice.Services",
	HandlerType: (*ServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGame",
			Handler:    _Services_CreateGame_Handler,
		},
		{
			MethodName: "ListGame",
			Handler:    _Services_ListGame_Handler,
		},
		{
			MethodName: "GetCurrent",
			Handler:    _Services_GetCurrent_Handler,
		},
		{
			MethodName: "PickGame",
			Handler:    _Services_PickGame_Handler,
		},
		{
			MethodName: "PlayGame",
			Handler:    _Services_PlayGame_Handler,
		},
		{
			MethodName: "UpdateGame",
			Handler:    _Services_UpdateGame_Handler,
		},
		{
			MethodName: "HintGame",
			Handler:    _Services_HintGame_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _Services_CreateUser_Handler,
		},
		{
			MethodName: "GetListUser",
			Handler:    _Services_GetListUser_Handler,
		},
		{
			MethodName: "GetLeaderBoard",
			Handler:    _Services_GetLeaderBoard_Handler,
		},
		{
			MethodName: "LogIn",
			Handler:    _Services_LogIn_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gameservice.proto",
}
