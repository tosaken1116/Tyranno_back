package controller

import (
	"errors"
	"fmt"
	"log"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
	"gorm.io/gorm"
)

type UserController struct{}

func (uc *UserController) CreateUser(conn *gorm.DB, msg *protosv1.CreateUserRequest, firebase_id string) (*protosv1.CreateUserResponse, error) {
	u := db.Users{
		DisplayId:  msg.DisplayId,
		Name:       msg.Name,
		FirebaseId: firebase_id,
		Icon:       msg.Icon,
	}

	if err := conn.Create(&u).Error; err != nil {
		resp := connect.NewError(connect.CodeInternal, err)
		log.Println(err)
		return nil, resp
	}

	userResp := &protosv1.CreateUserResponse{
		User: u.ToProtosModel(),
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
		User: u.ToProtosModel(),
	}
	return userResp, nil
}

func (uc *UserController) DeleteUser(conn *gorm.DB, id string) (*protosv1.DeleteUserResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "id = ?", id).Error; err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := conn.Delete(&u).Error; err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &protosv1.DeleteUserResponse{
		Status: true,
	}, nil
}

func (uc *UserController) CheckDisplayId(conn *gorm.DB, display_id string) (*protosv1.CheckDisplayNameResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "display_id = ?", display_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &protosv1.CheckDisplayNameResponse{
				IsNotExist: false,
			}, nil
		} else {
			log.Println(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &protosv1.CheckDisplayNameResponse{
		IsNotExist: true,
	}, nil
}

func (uc *UserController) GetUser(conn *gorm.DB, display_id string) (*protosv1.GetUserResponse, error) {
	u := db.Users{}

	if err := conn.First(&u, "display_id = ?", display_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &protosv1.GetUserResponse{User: nil}, nil
		} else {
			log.Println(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return &protosv1.GetUserResponse{
		User: u.ToProtosModel(),
	}, nil
}

func (uc *UserController) GetUsers(conn *gorm.DB) (*protosv1.GetUsersResponse, error) {
	u := []db.Users{}
	pu := []*protosv1.User{}

	if err := conn.Find(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &protosv1.GetUsersResponse{Users: pu}, nil
		} else {
			log.Println(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	for _, v := range u {
		pu = append(pu, v.ToProtosModel())
	}

	return &protosv1.GetUsersResponse{
		Users: pu,
	}, nil
}

func (uc *UserController) CheckVerifyTotp(conn *gorm.DB, firebase_id string) (string, error) {
	u := db.Users{}

	if err := conn.First(&u, "firebase_id = ?", firebase_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("user having this firebase_id is not found"))
		} else {
			log.Println(err)
			return "", connect.NewError(connect.CodeInternal, err)
		}
	}

	if !u.OtpVerified {
		return "", connect.NewError(connect.CodePermissionDenied, fmt.Errorf("totp check failed"))
	}

	return u.ID.String(), nil
}
