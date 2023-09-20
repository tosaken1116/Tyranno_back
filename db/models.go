package db

import (
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uuid.New()
	tx.Statement.SetColumn("ID", uuid)
	return
}

type DateTime struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Base struct {
	DateTime
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

type Users struct {
	Base
	DisplayId      string `gorm:"type:string;unique uniqueIndex" json:"uid"`
	FirebaseId     string `gorm:"unique" json:"firebase_id"`
	Name           string `gorm:"not null" json:"name"`
	Icon           string `gorm:"not null" json:"icon"`
	Profile        string `json:"profile"`
	OtpSecret      string `json:"otp_secret"`
	OtpUrl         string `json:"otp_url"`
	OtpEnabled     bool   `gorm:"not null default:false" json:"otp_enabled"`
	OtpVerified    bool   `gorm:"not null default:false" json:"otp_verified"`
	IsDelete       bool   `gorm:"not null default:false" json:"is_delete"`
	FollowNumber   int32  `json:"follow_number"`
	FollowerNumber int32  `json:"follower_number"`
}

func (u *Users) ToProtosModel() *protosv1.User {
	return &protosv1.User{
		DisplayId:      u.DisplayId,
		Name:           u.Name,
		Icon:           u.Icon,
		Profile:        u.Profile,
		FollowNumber:   u.FollowNumber,
		FollowerNumber: u.FollowerNumber,
		CreatedAt:      u.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:      u.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type Posts struct {
	DateTime
	ID             int32     `gorm:"primaryKey" json:"id"`
	Text           string    `gorm:"not null" json:"text"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	ReplyAt        *int32    `gorm:"default:null" json:"reply_at"`
	PublishedAt    time.Time `json:"published_at"`
	FavoriteNumber int32     `json:"favorite_number"`
	ReplyNumber    int32     `json:"reply_number"`
	IsDelete       bool      `gorm:"not null default:false" json:"is_delete"`
	User           Users     `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
	ReplyAtPost    *Posts    `gorm:"foreignKey:ReplyAt;reference:ID"`
}

func (p *Posts) ToProtosModel() *protosv1.Post {
	return &protosv1.Post{
		Id:             p.ID,
		Text:           p.Text,
		User:           p.User.ToProtosModel(),
		FavoriteNumber: p.FavoriteNumber,
		ReplyAt:        p.ReplyAt,
		ReplyNumber:    p.ReplyNumber,
		PublishedAt:    p.PublishedAt.Format(time.RFC3339Nano),
		CreatedAt:      p.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type Favorites struct {
	DateTime
	UserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	PostID int32     `gorm:"primaryKey" json:"post_id"`
	User   Users     `gorm:"foreignKey:UserID;reference:ID"`
	Post   Posts     `gorm:"foreignKey:PostID;reference:ID"`
}

type Follows struct {
	DateTime
	FromUserID uuid.UUID `gorm:"type:uuid;primaryKey" json:"from_user_id"`
	ToUserID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"to_user_id"`
	FromUser   Users     `gorm:"foreignKey:FromUserID;reference:ID"`
	ToUser     Users     `gorm:"foreignKey:ToUserID;reference:ID"`
}
