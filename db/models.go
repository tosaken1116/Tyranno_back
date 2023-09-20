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

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Users struct {
	Base
	DisplayId   string `gorm:"type:string;unique uniqueIndex" json:"uid"`
	FirebaseId  string `gorm:"unique" json:"firebase_id"`
	Name        string `gorm:"not null" json:"name"`
	Icon        string `gorm:"not null" json:"icon"`
	Profile     string `json:"profile"`
	OtpSecret   string `json:"otp_secret"`
	OtpUrl      string `json:"otp_url"`
	OtpEnabled  bool   `gorm:"not null default:false" json:"otp_enabled"`
	OtpVerified bool   `gorm:"not null default:false" json:"otp_verified"`
}

func (u *Users) ToProtosModel() *protosv1.User {
	return &protosv1.User{
		DisplayId: u.DisplayId,
		Name:      u.Name,
		Icon:      u.Icon,
		Profile:   u.Profile,
		CreatedAt: u.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type Posts struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Text        string    `gorm:"not null" json:"text"`
	UserID      uuid.UUID `json:"user_id"`
	ReplyAt     *int64    `gorm:"default:null" json:"reply_at"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        Users     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ReplyAtPost *Posts    `gorm:"foreignKey:ReplyAt;reference:ID"`
}

func (p *Posts) ToProtosModel(FavoriteNumber int32, ReplyNumber int32) *protosv1.Post {
	return &protosv1.Post{
		Id:             p.ID,
		Text:           p.Text,
		User:           p.User.ToProtosModel(),
		FavoriteNumber: FavoriteNumber,
		ReplyAt:        p.ReplyAt,
		ReplyNumber:    ReplyNumber,
		PublishedAt:    p.PublishedAt.Format(time.RFC3339Nano),
		CreatedAt:      p.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339Nano),
	}
}

type Favorites struct {
	Base
	UserID uuid.UUID
	PostID int32
}
