package model

import (
	"time"

	"github.com/google/uuid"
)

type CombinedBetModel struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ParticipantID      uuid.UUID `gorm:"type:uuid;not null;index"`
	LeagueID           uuid.UUID `gorm:"type:uuid;not null;index"`
	Stake              float64   `gorm:"type:numeric(12,2);not null"`
	UseBonus           bool      `gorm:"default:false"`
	TotalOdds          float64   `gorm:"type:numeric(16,4);not null"`
	PotentialWin       float64   `gorm:"type:numeric(12,2);not null"`
	Status             string    `gorm:"type:varchar(20);not null;default:'PENDING';index"`
	CashoutValue       *float64  `gorm:"type:numeric(12,2)"`
	CashoutExpiresAt   *time.Time
	CreatedAt          time.Time          `gorm:"autoCreateTime"`
	UpdatedAt          time.Time          `gorm:"autoUpdateTime"`
	SettledAt          *time.Time
	ParticipantBonusID *uuid.UUID         `gorm:"type:uuid"`
	Legs               []CombinedBetLegModel `gorm:"foreignKey:CombinedBetID;constraint:OnDelete:CASCADE;"`
}

func (CombinedBetModel) TableName() string {
	return "combined_bets"
}

type CombinedBetLegModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CombinedBetID   uuid.UUID `gorm:"type:uuid;not null;index"`
	MarketID        uuid.UUID `gorm:"type:uuid;not null;index"`
	MatchID         *uuid.UUID `gorm:"type:uuid;index"`
	SelectionID     uuid.UUID `gorm:"type:uuid;not null"`
	SelectionName   string    `gorm:"type:varchar(255);not null"`
	OddsAtPlacement float64   `gorm:"type:numeric(10,4);not null"`
	Status          string    `gorm:"type:varchar(20);not null;default:'PENDING'"`
	SettledAt       *time.Time
}

func (CombinedBetLegModel) TableName() string {
	return "combined_bet_legs"
}
