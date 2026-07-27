package output

import "time"

type UserOutput struct {
	ID        string
	Username  string
	Email     string
	AvatarURL string
	CreatedAt time.Time
}