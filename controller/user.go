package controller

import (
	"log"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"time"

	"connectrpc.com/connect"
	"gorm.io/gorm"
)

type UserController struct{}

func (uc *UserController) CreateUser(conn *gorm.DB, msg *protosv1.CreateUserRequest) (*protosv1.CreateUserResponse, error) {
	u := db.Users{
		DisplayId: msg.DisplayId,
		Name:      msg.Name,
	}

	if err := conn.Create(&u).Error; err != nil {
		resp := connect.NewError(connect.CodeInternal, err)
		log.Println(err)
		return nil, resp
	}

	userResp := &protosv1.CreateUserResponse{
		User: &protosv1.User{
			DisplayId: u.DisplayId,
			Name:      msg.GetName(),
			Icon:      "",
			Profile:   "",
			CreatedAt: u.CreatedAt.Format(time.RFC3339Nano),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339Nano),
		},
	}
	return userResp, nil
}
