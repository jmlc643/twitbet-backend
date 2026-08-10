package model

import "time"

type LeagueModel struct {
	ID             string    `gorm:"type:uuid;primaryKey;column:id"`
	OwnerID        string    `gorm:"type:uuid;not null;column:owner_id"`
	Name           string    `gorm:"type:varchar(100);not null;column:name"`
	Slug           string    `gorm:"type:varchar(150);uniqueIndex;not null;column:slug"`
	InitialBalance float64   `gorm:"type:numeric(12,2);not null;column:initial_balance"`
	MaxRecharges   int       `gorm:"type:integer;not null;default:2;column:max_recharges"`
	HideStandings  bool      `gorm:"type:boolean;not null;default:false;column:hide_standings"`
	Status         string    `gorm:"type:varchar(20);not null;default:'ACTIVE';column:status"`
	MinBetsToQualify int       `gorm:"type:integer;not null;default:0;column:min_bets_to_qualify"`
	InviteCode     string    `gorm:"type:varchar(20);uniqueIndex;not null;column:invite_code"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (LeagueModel) TableName() string {
	return "leagues"
}
