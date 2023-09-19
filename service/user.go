package service

import (
	"context"
	"fmt"
	"log"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
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
