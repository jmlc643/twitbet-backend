package model

import "time"

type MarketOddsHistoryModel struct {
	ID             string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MarketOptionID string    `gorm:"type:uuid;not null;index:idx_moh_option_time,priority:1"`
	OldOdds        *float64  `gorm:"type:decimal(10,4)"`
	NewOdds        float64   `gorm:"type:decimal(10,4);not null"`
	ChangedBy      *string   `gorm:"type:uuid"`
	Reason         string    `gorm:"type:varchar(50);not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_moh_option_time,priority:2,sort:desc"`
}

func (MarketOddsHistoryModel) TableName() string {
	return "market_odds_history"
}
