package model

import "time"

type MarketModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;column:id"`
	LeagueID  string    `gorm:"type:uuid;not null;column:league_id"`
	MatchID   *string   `gorm:"type:uuid;column:match_id"`
	Name      string    `gorm:"type:varchar(200);not null;column:name"`
	Status    string    `gorm:"type:varchar(50);not null;default:'OPEN';column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	Options []MarketOptionModel `gorm:"foreignKey:MarketID"`
}

func (MarketModel) TableName() string {
	return "markets"
}

type MarketOptionModel struct {
	ID          string  `gorm:"type:uuid;primaryKey;column:id"`
	MarketID    string  `gorm:"type:uuid;not null;column:market_id"`
	Name        string  `gorm:"type:varchar(100);not null;column:name"`
	InitialOdds float64 `gorm:"type:numeric(10,2);not null;column:initial_odds"`
	CurrentOdds float64 `gorm:"type:numeric(10,2);not null;column:current_odds"`
}

func (MarketOptionModel) TableName() string {
	return "market_options"
}
