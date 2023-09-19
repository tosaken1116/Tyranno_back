package service

import (
	"context"
	"fmt"
	"nnyd-back/config"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"time"

	"connectrpc.com/connect"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct{}

func (us *UserServer) CreateUser(ctx context.Context, req *connect.Request[protosv1.CreateUserRequest]) (*connect.Response[protosv1.CreateUserResponse], error) {
	firebase_id := ctx.Value(config.FIREBASE_ID).(string)
	if firebase_id == "" {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("verifying failed"))
	}

	if req.Msg.DisplayId == "" || req.Msg.Name == "" || req.Msg.Icon == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("display_id, name and icon is required"))
	}

	conn := db.GetDB()
	uc := &controller.UserController{}
	userResp, err := uc.CreateUser(conn, req.Msg, firebase_id)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(userResp), nil
}

func (us *UserServer) Signin(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.SigninResponse], error) {
	firebase_id := ctx.Value(config.FIREBASE_ID).(string)
	if firebase_id == "" {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("verifying failed"))
	}
	conn := db.GetDB()
	uc := &controller.UserController{}
	user_id, err := uc.CheckVerifyTotp(conn, firebase_id)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(3000 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(config.JST_SECRET_KEY))
	if err != nil {
		return nil, err
	}

	resp := &protosv1.SigninResponse{
		Token: accessToken,
	}
	return connect.NewResponse(resp), nil
}

func (us *UserServer) UpdateUser(ctx context.Context, req *connect.Request[protosv1.UpdateUserRequest]) (*connect.Response[protosv1.UpdateUserResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	if user_id == "" {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("verifying failed"))
	}

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
	if user_id == "" {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("verifying failed"))
	}

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
