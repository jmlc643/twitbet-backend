package model

import "time"

type ParticipantModel struct {
	ID       string    `gorm:"type:uuid;primaryKey;column:id"`
	LeagueID string    `gorm:"type:uuid;not null;column:league_id"`
	UserID   string    `gorm:"type:uuid;not null;column:user_id"`
	Balance  float64   `gorm:"type:numeric(12,2);not null;column:balance"`
	JoinedAt time.Time `gorm:"column:joined_at"`
}

func (ParticipantModel) TableName() string {
	return "league_participants"
}
