package controller

import (
	"log"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
	"gorm.io/gorm"
)

type UserController struct{}

func (uc *UserController) CreateUser(conn *gorm.DB, msg *protosv1.CreateUserRequest) (*protosv1.CreateUserResponse, error) {
	u := db.Users{
		Uid:  msg.Name,
		Name: msg.Name,
		Icon: msg.Icon,
	}

	if err := conn.Create(&u).Error; err != nil {
		resp := connect.NewError(connect.CodeInternal, err)
		log.Fatal(err)
		return nil, resp
	}

	responseUser := &protosv1.User{
		Id:   "dummy",
		Name: msg.GetName(),
		Icon: msg.GetIcon(),
	}
	userResp := &protosv1.CreateUserResponse{
		User: responseUser,
	}
	return userResp, nil
}
