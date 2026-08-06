package model

import "time"

type BetModel struct {
	ID             string    `gorm:"primaryKey;type:uuid"`
	ParticipantID  string    `gorm:"type:uuid;not null;index"`
	MarketOptionID string    `gorm:"type:uuid;not null;index"`
	Amount         float64   `gorm:"not null"`
	Odds           float64   `gorm:"not null"`
	PotentialWin   float64   `gorm:"not null"`
	Status         string    `gorm:"type:varchar(20);not null"`
	PlacedAt       time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (BetModel) TableName() string {
	return "bets"
}
