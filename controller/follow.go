package controller

import (
	"errors"
	"fmt"
	"log"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowController struct{}

func (foc *FollowController) FollowUser(conn *gorm.DB, from_user_id string, to_user_id string) (*protosv1.FollowUserResponse, error) {
	if from_user_id == to_user_id {
		err := fmt.Errorf("can not follow yourself")
		return nil, err
	}
	follow := db.Follows{}
	if err := conn.First(&follow, "from_user_id = ? and to_user_id = ?", from_user_id, to_user_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fui, err := uuid.Parse(from_user_id)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			tui, err := uuid.Parse(to_user_id)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			follow = db.Follows{
				FromUserID: fui,
				ToUserID:   tui,
			}

			if err := conn.Create(&follow).Error; err != nil {
				log.Println(err)
				return nil, err
			}

			fu := db.Users{}
			if err := conn.First(&fu, "id = ? and is_delete = false", from_user_id).Error; err != nil {
				log.Println(err)
				return nil, err
			}
			fu.FollowNumber = fu.FollowNumber + 1
			if err := conn.Save(&fu).Error; err != nil {
				log.Println(err)
				return nil, err
			}

			tu := db.Users{}
			if err := conn.First(&tu, "id = ? and is_delete = false", to_user_id).Error; err != nil {
				log.Println(err)
				return nil, err
			}
			tu.FollowerNumber = tu.FollowerNumber + 1
			if err := conn.Save(&tu).Error; err != nil {
				log.Println(err)
				return nil, err
			}

			return &protosv1.FollowUserResponse{
				User: tu.ToProtosModel(),
			}, nil
		} else {
			log.Println(err)
			return nil, err
		}
	}
	err := fmt.Errorf("this user is already follow")
	return nil, err

}

func (foc *FollowController) UnfollowUser(conn *gorm.DB, from_user_id string, to_user_id string) (*protosv1.UnfollowUserResponse, error) {
	if from_user_id == to_user_id {
		err := fmt.Errorf("can not follow yourself")
		return nil, err
	}
	follow := db.Follows{}
	if err := conn.First(&follow, "from_user_id = ? and to_user_id = ?", from_user_id, to_user_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := fmt.Errorf("this user is already unfollow")
			return nil, err
		} else {
			log.Println(err)
			return nil, err
		}
	}

	if err := conn.Delete(&follow).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	fu := db.Users{}
	if err := conn.First(&fu, "id = ? and is_delete = false", from_user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	fu.FollowNumber = fu.FollowNumber - 1
	if err := conn.Save(&fu).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	tu := db.Users{}
	if err := conn.First(&tu, "id = ? and is_delete = false", to_user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	tu.FollowerNumber = tu.FollowerNumber - 1
	if err := conn.Save(&tu).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return &protosv1.UnfollowUserResponse{
		User: tu.ToProtosModel(),
	}, nil
}

func (foc *FollowController) GetFollowByID(conn *gorm.DB, user_id string) (*protosv1.GetUsersResponse, error) {
	f := []db.Follows{}
	if err := conn.Find(&f, "from_user_id = ?", user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	users := []*protosv1.User{}
	for i := 0; i < len(f); i++ {
		if err := conn.Model(&f[i]).Association("ToUser").Find(&f[i].ToUser); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, f[i].ToUser.ToProtosModel())
	}

	return &protosv1.GetUsersResponse{
		Users: users,
	}, nil
}

func (foc *FollowController) GetFollowerByID(conn *gorm.DB, user_id string) (*protosv1.GetUsersResponse, error) {
	f := []db.Follows{}
	if err := conn.Find(&f, "to_user_id = ?", user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	users := []*protosv1.User{}
	for i := 0; i < len(f); i++ {
		if err := conn.Model(&f[i]).Association("FromUser").Find(&f[i].FromUser); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, f[i].FromUser.ToProtosModel())
	}

	return &protosv1.GetUsersResponse{
		Users: users,
	}, nil
}
