package model

import "time"

type TransactionModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id"`
	LeagueID  string    `gorm:"type:uuid;not null;column:league_id"`
	UserID    string    `gorm:"type:uuid;not null;column:user_id"`
	Amount    float64   `gorm:"type:numeric(12,2);not null;column:amount"`
	Type      string    `gorm:"type:varchar(50);not null;column:type"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}
