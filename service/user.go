package service

import (
	"context"
	"fmt"
	"log"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct{}

func (us *UserServer) CreateUser(ctx context.Context, req *connect.Request[protosv1.CreateUserRequest]) (*connect.Response[protosv1.CreateUserResponse], error) {
	log.Println("Request headers: ", req.Header())

	if req.Msg.Name == "" || req.Msg.Icon == "" {
		// エラーにステータスコードを追加
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name and icon is required."))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.CreateUser(conn, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) Signin(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.SigninResponse], error) {
	// mock
	resp := &protosv1.SigninResponse{
		Token: "",
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) UpdateUser(context.Context, *connect.Request[protosv1.UpdateUserRequest]) (*connect.Response[protosv1.UpdateUserResponse], error) {
	// mock
	resp := &protosv1.UpdateUserResponse{
		User: nil,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) DeleteUser(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.DeleteUserResponse], error) {
	// mock
	resp := &protosv1.DeleteUserResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetUser(context.Context, *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUserResponse], error) {
	// mock
	resp := &protosv1.GetUserResponse{
		User: nil,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetUsers(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) CheckDisplayName(context.Context, *connect.Request[protosv1.CheckDisplayNameRequest]) (*connect.Response[protosv1.CheckDisplayNameResponse], error) {
	// mock
	resp := &protosv1.CheckDisplayNameResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}
