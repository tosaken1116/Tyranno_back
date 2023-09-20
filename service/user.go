package service

import (
	"context"
	"fmt"
	"nnyd-back/config"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct{}

func (us *UserServer) CreateUser(ctx context.Context, req *connect.Request[protosv1.CreateUserRequest]) (*connect.Response[protosv1.CreateUserResponse], error) {
	firebase_id := ctx.Value(config.FIREBASE_ID).(string)
	if req.Msg.DisplayId == "" || req.Msg.Name == "" || req.Msg.Icon == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id, name and icon is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.CreateUser(conn, req.Msg, firebase_id)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) UpdateUser(ctx context.Context, req *connect.Request[protosv1.UpdateUserRequest]) (*connect.Response[protosv1.UpdateUserResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	if req.Msg.DisplayId == "" || req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id and name is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.UpdateUser(conn, user_id, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) DeleteUser(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.DeleteUserResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	conn := db.GetDB()
	uc := &controller.UserController{}
	resultResp, err := uc.DeleteUser(conn, user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetUser(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUserResponse], error) {
	conn := db.GetDB()
	uc := &controller.UserController{}
	resultResp, err := uc.GetUser(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetMe(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUserResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	conn := db.GetDB()
	uc := &controller.UserController{}
	resultResp, err := uc.GetUserById(conn, user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetUsers(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	conn := db.GetDB()
	uc := &controller.UserController{}
	resultResp, err := uc.GetUsers(conn)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) CheckDisplayName(ctx context.Context, req *connect.Request[protosv1.CheckDisplayNameRequest]) (*connect.Response[protosv1.CheckDisplayNameResponse], error) {
	if req.Msg.CheckText == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("check_test is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	resultResp, err := uc.CheckDisplayId(conn, req.Msg.CheckText)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) FollowUser(context.Context, *connect.Request[protosv1.FollowUserRequest]) (*connect.Response[protosv1.FollowUserResponse], error) {
	// mock
	resp := &protosv1.FollowUserResponse{
		User: nil,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) UnfollowUser(context.Context, *connect.Request[protosv1.UnfollowUserRequest]) (*connect.Response[protosv1.UnfollowUserResponse], error) {
	// mock
	resp := &protosv1.UnfollowUserResponse{
		User: nil,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetFollowByID(context.Context, *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetFollowerByID(context.Context, *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetMyFollow(context.Context, *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetMyFollower(context.Context, *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	// mock
	resp := &protosv1.GetUsersResponse{
		Users: []*protosv1.User{},
	}
	return connect.NewResponse(resp), nil
}
