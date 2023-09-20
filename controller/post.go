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

func (uc *PostController) CreatePostController(conn *gorm.DB, msg *protosv1.CreatePostRequest, user_id string) (*protosv1.CreatePostResponse, error) {
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
	conn.Model(&p).Association("User").Find(&p.User)

	postResponse := &protosv1.CreatePostResponse{
		Post: p.ToProtosModel(0, 0),
	}

	return postResponse, nil
}
