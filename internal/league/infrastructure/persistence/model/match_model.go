package model

import "time"

type MatchModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id"`
	LeagueID  string    `gorm:"type:uuid;not null;column:league_id"`
	Slug      string    `gorm:"type:varchar(255);unique;not null;column:slug"`
	Title     string    `gorm:"type:varchar(200);not null;column:title"`
	StartTime time.Time `gorm:"column:start_time"`
	Status    string    `gorm:"type:varchar(50);not null;default:'SCHEDULED';column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	League LeagueModel `gorm:"foreignKey:LeagueID"`
}

func (MatchModel) TableName() string {
	return "matches"
}