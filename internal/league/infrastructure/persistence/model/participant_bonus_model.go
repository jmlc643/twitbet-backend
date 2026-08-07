package model

import "time"

type ParticipantBonusModel struct {
	ID            string    `gorm:"primaryKey;type:uuid"`
	LeagueID      string    `gorm:"type:uuid;not null;index"`
	ParticipantID string    `gorm:"type:uuid;not null;index"`
	Amount        float64   `gorm:"not null"`
	Status        string    `gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (ParticipantBonusModel) TableName() string {
	return "participant_bonuses"
}
