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

	if req.Msg.DisplayId == "" || req.Msg.Name == "" || req.Msg.Icon == "" {
		// エラーにステータスコードを追加
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id, name and icon is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.CreateUser(conn, req.Msg)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) Signin(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.SigninResponse], error) {
	// mock
	resp := &protosv1.SigninResponse{
		Token: "",
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) UpdateUser(ctx context.Context, req *connect.Request[protosv1.UpdateUserRequest]) (*connect.Response[protosv1.UpdateUserResponse], error) {
	user_id := ctx.Value("user_id")

	if user_id == "" {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("verifying failed"))
	}
	if req.Msg.DisplayId == "" || req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id and name is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.UpdateUser(conn, user_id.(string), req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) DeleteUser(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.DeleteUserResponse], error) {
	// mock
	resp := &protosv1.DeleteUserResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetUser(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUserResponse], error) {
	// mock
	resp := &protosv1.GetUserResponse{
		User: nil,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetUsers(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) CheckDisplayName(ctx context.Context, req *connect.Request[protosv1.CheckDisplayNameRequest]) (*connect.Response[protosv1.CheckDisplayNameResponse], error) {
	// mock
	resp := &protosv1.CheckDisplayNameResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}
