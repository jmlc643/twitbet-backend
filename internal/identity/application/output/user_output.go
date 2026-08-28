package output

import "time"

type UserStatsOutput struct {
	Leagues       int
	Wins          int
	Effectiveness int
}

type UserOutput struct {
	ID        string
	Username  string
	Email     string
	AvatarURL *string
	CreatedAt time.Time
	Stats     UserStatsOutput
}
