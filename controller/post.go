package controller

import (
	"errors"
	"fmt"
	"log"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"nnyd-back/utils"
	"time"

	"gorm.io/gorm"
)

type PostController struct{}

func (pc *PostController) CreatePost(conn *gorm.DB, msg *protosv1.CreatePostRequest, user_id string) (*protosv1.CreatePostResponse, error) {
	if msg.Text == "" {
		err := fmt.Errorf("invalid argument")
		log.Println(err)
		return nil, err
	}

	u := db.Users{}

	if err := conn.First(&u, "id = ? and is_delete = false", user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if msg.ReplyAt != nil {
		reply := db.Posts{}
		if err := conn.First(&reply, "id = ? and is_delete = false", msg.ReplyAt).Error; err != nil {
			log.Println(err)
			return nil, err
		}
		reply.ReplyNumber = reply.ReplyNumber + 1
		if err := conn.Save(&reply).Error; err != nil {
			log.Println(err)
			return nil, err
		}
	}

	nowTime := time.Now()
	durationTime := time.Second * utils.GenerateRandomDelayTime()

	p := db.Posts{
		Text:        msg.Text,
		UserID:      u.ID,
		ReplyAt:     msg.ReplyAt,
		PublishedAt: nowTime.Add(durationTime),
	}

	if err := conn.Create(&p).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	if err := conn.Model(&p).Association("User").Find(&p.User); err != nil {
		log.Println(err)
		return nil, err
	}

	postResponse := &protosv1.CreatePostResponse{
		Post: p.ToProtosModel(false),
	}

	return postResponse, nil
}

func (pc *PostController) GetPostByID(conn *gorm.DB, post_id int32, user_id *string) (*protosv1.GetPostResponse, error) {
	p := db.Posts{}
	if err := conn.Find(&p, "id = ? and is_delete = false", post_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	if time.Now().Before(p.PublishedAt) && *user_id != p.UserID.String() {
		return nil, gorm.ErrRecordNotFound
	}
	if err := conn.Model(&p).Association("User").Find(&p.User); err != nil {
		log.Println(err)
		return nil, err
	}
	isFavorited := false
	if user_id != nil {
		fa := db.Favorites{}
		if err := conn.First(&fa, "user_id = ? and post_id = ?", user_id, post_id).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println(err)
				return nil, err
			}
		} else {
			isFavorited = true
		}
	}

	postResponse := &protosv1.GetPostResponse{
		Post: p.ToProtosModel(isFavorited),
	}

	return postResponse, nil
}

func (pc *PostController) GetPosts(conn *gorm.DB, user_id *string) (*protosv1.GetPostsResponse, error) {
	p := []db.Posts{}
	if err := conn.Order("published_at desc").Find(&p, "is_delete = false and published_at <= now()").Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return nil, err
		}
	}
	for i := 0; i < len(p); i++ {
		if err := conn.Model(&p[i]).Association("User").Find(&p[i].User); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	posts := []*protosv1.Post{}

	for i := 0; i < len(p); i++ {
		isFavorited := false
		if user_id != nil {
			if err := conn.First(&db.Favorites{}, "user_id = ? and post_id = ?", user_id, p[i].ID).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Println(err)
					return nil, err
				}
			} else {
				isFavorited = true
			}
		}
		posts = append(posts, p[i].ToProtosModel(isFavorited))
	}

	postResponse := &protosv1.GetPostsResponse{
		Posts: posts,
	}

	return postResponse, nil
}

func (pc *PostController) DeletePost(conn *gorm.DB, post_id int32) (*protosv1.DeletePostResponse, error) {
	p := db.Posts{}

	if err := conn.Find(&p, "id = ? and is_delete = false", post_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	p.IsDelete = true

	if err := conn.Save(&p).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return &protosv1.DeletePostResponse{
		Status: true,
	}, nil
}

func (pc *PostController) GetReplies(conn *gorm.DB, post_id int32, user_id *string) (*protosv1.GetRepliesResponse, error) {
	p := []db.Posts{}
	if err := conn.Find(&p, "reply_at = ? and is_delete = false and published_at <= now()", post_id).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return nil, err
		}
	}
	for i := 0; i < len(p); i++ {
		if err := conn.Model(&p[i]).Association("User").Find(&p[i].User); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	posts := []*protosv1.Post{}

	for i := 0; i < len(p); i++ {
		isFavorited := false
		if user_id != nil {
			if err := conn.First(&db.Favorites{}, "user_id = ? and post_id = ?", user_id, p[i].ID).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Println(err)
					return nil, err
				}
			} else {
				isFavorited = true
			}
		}
		posts = append(posts, p[i].ToProtosModel(isFavorited))
	}

	postResponse := &protosv1.GetRepliesResponse{
		Replies: posts,
	}

	return postResponse, nil
}
