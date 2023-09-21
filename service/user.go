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

var (
	uc  = &controller.UserController{}
	foc = &controller.FollowController{}
)

type UserServer struct{}

func (us *UserServer) CreateUser(ctx context.Context, req *connect.Request[protosv1.CreateUserRequest]) (*connect.Response[protosv1.CreateUserResponse], error) {
	firebase_id := ctx.Value(config.FIREBASE_ID)
	if firebase_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	firebase_id_str := firebase_id.(string)

	if req.Msg.DisplayId == "" || req.Msg.Name == "" || req.Msg.Icon == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id, name and icon is required"))
	}

	conn := db.GetDB()
	userResp, err := uc.CreateUser(conn, req.Msg, firebase_id_str)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) UpdateUser(ctx context.Context, req *connect.Request[protosv1.UpdateUserRequest]) (*connect.Response[protosv1.UpdateUserResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	if req.Msg.DisplayId == "" || req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id and name is required"))
	}

	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	userResp, err := uc.UpdateUser(conn, user_id_str, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) DeleteUser(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.DeleteUserResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resultResp, err := uc.DeleteUser(conn, user_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetUser(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUserResponse], error) {
	conn := db.GetDB()
	resultResp, err := uc.GetUser(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetMe(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUserResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resultResp, err := uc.GetUserById(conn, user_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) GetUsers(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	conn := db.GetDB()
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
	resultResp, err := uc.CheckDisplayId(conn, req.Msg.CheckText)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("db error"))
	}

	return connect.NewResponse(resultResp), nil
}

func (us *UserServer) FollowUser(ctx context.Context, req *connect.Request[protosv1.FollowUserRequest]) (*connect.Response[protosv1.FollowUserResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	to_user_id, err := uc.GetIdFromDisplayId(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.FollowUser(conn, user_id_str, to_user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) UnfollowUser(ctx context.Context, req *connect.Request[protosv1.UnfollowUserRequest]) (*connect.Response[protosv1.UnfollowUserResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	to_user_id, err := uc.GetIdFromDisplayId(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.UnfollowUser(conn, user_id_str, to_user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetFollowByID(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	conn := db.GetDB()
	user_id, err := uc.GetIdFromDisplayId(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.GetFollowByID(conn, user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetFollowerByID(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	conn := db.GetDB()
	user_id, err := uc.GetIdFromDisplayId(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.GetFollowerByID(conn, user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetMyFollow(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.GetFollowByID(conn, user_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) GetMyFollower(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetUsersResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := foc.GetFollowerByID(conn, user_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}
