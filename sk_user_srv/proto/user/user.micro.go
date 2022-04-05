// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user/user.proto

/*
Package yuyoung_srv_sk_user_srv is a generated protocol buffer package.

It is generated from these files:
	proto/user/user.proto

It has these top-level messages:
	BasicUserInfo
	UserInfo
	GeneralRequest
	GeneralResponse
	RegisterUserRequest
	GetUserInfoRequest
	GetUserInfoResponse
*/
package yuyoung_srv_sk_user_srv

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for UserService service

type UserService interface {
	Register(ctx context.Context, in *RegisterUserRequest, opts ...client.CallOption) (*GeneralResponse, error)
	Login(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error)
	AdminLogin(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error)
	SendEmail(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error)
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...client.CallOption) (*GetUserInfoResponse, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "yuyoung.srv.sk_user_srv"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) Register(ctx context.Context, in *RegisterUserRequest, opts ...client.CallOption) (*GeneralResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.Register", in)
	out := new(GeneralResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Login(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.Login", in)
	out := new(GeneralResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) AdminLogin(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.AdminLogin", in)
	out := new(GeneralResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) SendEmail(ctx context.Context, in *GeneralRequest, opts ...client.CallOption) (*GeneralResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.SendEmail", in)
	out := new(GeneralResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...client.CallOption) (*GetUserInfoResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.GetUserInfo", in)
	out := new(GetUserInfoResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	Register(context.Context, *RegisterUserRequest, *GeneralResponse) error
	Login(context.Context, *GeneralRequest, *GeneralResponse) error
	AdminLogin(context.Context, *GeneralRequest, *GeneralResponse) error
	SendEmail(context.Context, *GeneralRequest, *GeneralResponse) error
	GetUserInfo(context.Context, *GetUserInfoRequest, *GetUserInfoResponse) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		Register(ctx context.Context, in *RegisterUserRequest, out *GeneralResponse) error
		Login(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error
		AdminLogin(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error
		SendEmail(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error
		GetUserInfo(ctx context.Context, in *GetUserInfoRequest, out *GetUserInfoResponse) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) Register(ctx context.Context, in *RegisterUserRequest, out *GeneralResponse) error {
	return h.UserServiceHandler.Register(ctx, in, out)
}

func (h *userServiceHandler) Login(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error {
	return h.UserServiceHandler.Login(ctx, in, out)
}

func (h *userServiceHandler) AdminLogin(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error {
	return h.UserServiceHandler.AdminLogin(ctx, in, out)
}

func (h *userServiceHandler) SendEmail(ctx context.Context, in *GeneralRequest, out *GeneralResponse) error {
	return h.UserServiceHandler.SendEmail(ctx, in, out)
}

func (h *userServiceHandler) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, out *GetUserInfoResponse) error {
	return h.UserServiceHandler.GetUserInfo(ctx, in, out)
}
