package controller

import (
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

	if err := conn.First(&u, "id = ?", user_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if msg.ReplyAt != nil {
		reply := db.Posts{}
		if err := conn.First(&reply, "id = ?", msg.ReplyAt).Error; err != nil {
			reply.ReplyNumber = reply.ReplyNumber + 1
		}
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
		Post: p.ToProtosModel(),
	}

	return postResponse, nil
}

func (pc *PostController) GetPostByID(conn *gorm.DB, post_id int64) (*protosv1.GetPostResponse, error) {
	p := db.Posts{}
	if err := conn.Find(&p, "id = ?", post_id).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	if err := conn.Model(&p).Association("User").Find(&p.User); err != nil {
		log.Println(err)
		return nil, err
	}

	postResponse := &protosv1.GetPostResponse{
		Post: p.ToProtosModel(),
	}

	return postResponse, nil
}

func (pc *PostController) GetPosts(conn *gorm.DB) (*protosv1.GetPostsResponse, error) {
	p := []db.Posts{}
	if err := conn.Find(&p).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	for i := 0; i < len(p); i++ {
		if err := conn.Model(&p[i]).Association("User").Find(&p[i].User); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	posts := []*protosv1.Post{}

	for i := 0; i < len(p); i++ {
		posts = append(posts, p[i].ToProtosModel())
	}

	postResponse := &protosv1.GetPostsResponse{
		Posts: posts,
	}

	return postResponse, nil
}
