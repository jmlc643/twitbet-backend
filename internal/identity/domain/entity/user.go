package entity

import "time"

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	AvatarURL    string
	IsVerified   bool
	CreatedAt    time.Time
}
