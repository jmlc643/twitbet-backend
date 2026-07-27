package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"type:varchar(50);not null"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	AvatarURL    *string   `gorm:"type:varchar(255)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}