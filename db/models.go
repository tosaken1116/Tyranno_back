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
	Uid             string `gorm:"type:string;unique" json:"uid"`
	FirebaseId      string `gorm:"unique" json:"firebase_id"`
	Name            string `gorm:"not null" json:"name"`
	Icon            string `gorm:"not null" json:"icon"`
	AppVerifyStatus bool   `gorm:"not null default:false" json:"app_verify_status"`
}
