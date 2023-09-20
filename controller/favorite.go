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

func (fac *FavoriteController) CreateFavorite(conn *gorm.DB, user_id string, post_id int32) (*protosv1.CreateFavoriteResponse, error) {
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

func (fac *FavoriteController) DeleteFavorite(conn *gorm.DB, user_id string, post_id int32) (*protosv1.DeleteFavoriteResponse, error) {
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

func (fac *FavoriteController) GetFavoritePosts(conn *gorm.DB, user_id string, my_user_id *string) (*protosv1.GetPostsResponse, error) {
	fa := []db.Favorites{}
	if err := conn.First(&fa, "user_id = ?", user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	posts := []*protosv1.Post{}
	for i := 0; i < len(fa); i++ {
		if err := conn.Model(&fa[i]).Association("Post").Find(&fa[i].Post); err != nil {
			log.Println(err)
			return nil, err
		}
		var isFavorited bool = false
		if my_user_id != nil {
			if user_id == *my_user_id {
				isFavorited = true
			} else {
				if err := conn.First(&db.Favorites{}, "user_id = ? and post_id = ?", my_user_id, fa[i].PostID).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						log.Println(err)
						return nil, err
					}
				} else {
					isFavorited = true
				}
			}
		}
		posts = append(posts, fa[i].Post.ToProtosModel(isFavorited))
	}

	return &protosv1.GetPostsResponse{
		Posts: posts,
	}, nil
}

func (fac *FavoriteController) GetFavoritingUser(conn *gorm.DB, post_id int32) (*protosv1.GetUsersResponse, error) {
	fa := []db.Favorites{}
	if err := conn.First(&fa, "post_id = ?", post_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	users := []*protosv1.User{}
	for i := 0; i < len(fa); i++ {
		if err := conn.Model(&fa[i]).Association("User").Find(&fa[i].User); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, fa[i].User.ToProtosModel())
	}

	return &protosv1.GetUsersResponse{
		Users: users,
	}, nil
}
