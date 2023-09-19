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
		Icon:      msg.Icon,
	}

	if err := conn.Create(&u).Error; err != nil {
		resp := connect.NewError(connect.CodeInternal, err)
		log.Println(err)
		return nil, resp
	}

	userResp := &protosv1.CreateUserResponse{
		User: &protosv1.User{
			DisplayId: u.DisplayId,
			Name:      u.Name,
			Icon:      u.Icon,
			Profile:   u.Profile,
			CreatedAt: u.CreatedAt.Format(time.RFC3339Nano),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339Nano),
		},
	}
	return userResp, nil
}

func (uc *UserController) UpdateUser(conn *gorm.DB, id string, msg *protosv1.UpdateUserRequest) (*protosv1.UpdateUserResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "id = ?", id).Error; err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	u.DisplayId = msg.DisplayId
	u.Name = msg.Name
	u.Profile = msg.Profile

	if err := conn.Save(&u).Error; err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	userResp := &protosv1.UpdateUserResponse{
		User: &protosv1.User{
			DisplayId: u.DisplayId,
			Name:      u.Name,
			Icon:      u.Icon,
			Profile:   u.Profile,
			CreatedAt: u.CreatedAt.Format(time.RFC3339Nano),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339Nano),
		},
	}
	return userResp, nil
}
