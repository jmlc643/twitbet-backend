package model

import "time"

type ParticipantModel struct {
	ID                string    `gorm:"type:uuid;primaryKey;column:id"`
	LeagueID          string    `gorm:"type:uuid;not null;column:league_id"`
	UserID            string    `gorm:"type:uuid;not null;column:user_id"`
	IsAdmin           bool      `gorm:"type:boolean;not null;default:false;column:is_admin"`
	Balance           float64   `gorm:"type:numeric(12,2);not null;column:balance"`
	RechargesConsumed int       `gorm:"type:integer;not null;default:0;column:recharges_consumed"`
	JoinedAt          time.Time `gorm:"column:joined_at"`
}

func (ParticipantModel) TableName() string {
	return "league_participants"
}