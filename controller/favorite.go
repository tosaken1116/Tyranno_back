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

type FavoriteController struct{}

func (fac *FavoriteController) CreateFavorite(conn *gorm.DB, user_id string, post_id int64) (*protosv1.CreateFavoriteResponse, error) {
	fa := db.Favorites{}
	if err := conn.First(&fa, "user_id = ? and post_id = ?", user_id, post_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user_uuid, err := uuid.Parse(user_id)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			fa = db.Favorites{
				UserID: user_uuid,
				PostID: int32(post_id),
			}
			if err := conn.Create(&fa).Error; err != nil {
				log.Println(err)
				return nil, err
			}

			p := db.Posts{}
			if err := conn.Find(&p, "id = ? and is_delete = false", post_id).Error; err != nil {
				log.Println(err)
				return nil, err
			}
			p.FavoriteNumber = p.FavoriteNumber + 1
			if err := conn.Save(&p).Error; err != nil {
				log.Println(err)
				return nil, err
			}

			return &protosv1.CreateFavoriteResponse{
				Status: true,
			}, nil
		} else {
			log.Println(err)
			return nil, err
		}
	}
	err := fmt.Errorf("this post is already favorite")
	return nil, err
}

func (fac *FavoriteController) DeleteFavorite(conn *gorm.DB, user_id string, post_id int64) (*protosv1.DeleteFavoriteResponse, error) {
	fa := db.Favorites{}
	if err := conn.First(&fa, "user_id = ? and post_id = ?", user_id, post_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := fmt.Errorf("this post is not favorite by you")
			return nil, err
		} else {
			log.Println(err)
			return nil, err
		}
	}
	if err := conn.Delete(&fa).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	p := db.Posts{}
	if err := conn.Find(&p, "id = ? and is_delete = false", post_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	p.FavoriteNumber = p.FavoriteNumber - 1
	if err := conn.Save(&p).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return &protosv1.DeleteFavoriteResponse{
		Status: true,
	}, nil
}
