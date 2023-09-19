package db

import (
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
	DisplayId       string `gorm:"type:string;unique uniqueIndex" json:"uid"`
	FirebaseId      string `gorm:"unique" json:"firebase_id"`
	Name            string `gorm:"not null" json:"name"`
	Icon            string `gorm:"not null" json:"icon"`
	Profile         string `json:"profile"`
	AppVerifyStatus bool   `gorm:"not null default:false" json:"app_verify_status"`
}

type Posts struct {
	ID          int32     `gorm:"primaryKey" json:"id"`
	Text        string    `gorm:"not null" json:"text"`
	UserID      uuid.UUID `json:"user_id"`
	ReplyAt     *int64    `json:"reply_at"`
	PublishedAt string    `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        Users     `gorm:"foreignKey:UserID"`
	ReplyAtPost *Posts    `gorm:"foreignKey:ReplyAt"`
}

type Favorites struct {
	Base
	UserID uuid.UUID
	PostID int32
}
