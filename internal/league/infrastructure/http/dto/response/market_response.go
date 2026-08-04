package response

import (
	"time"

	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MarketOptionResponse struct {
	ID          string  `json:"id"`
	MarketID    string  `json:"market_id"`
	Name        string  `json:"name"`
	InitialOdds float64 `json:"initial_odds"`
	CurrentOdds float64 `json:"current_odds"`
}

type MarketResponse struct {
	ID        string                 `json:"id"`
	LeagueID  string                 `json:"league_id"`
	MatchID   *string                `json:"match_id"`
	Name      string                 `json:"name"`
	Status    string                 `json:"status"`
	Options   []MarketOptionResponse `json:"options"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewMarketResponse(market *entity.Market) MarketResponse {
	var matchID *string
	if market.MatchID != nil {
		idStr := market.MatchID.String()
		matchID = &idStr
	}

	options := make([]MarketOptionResponse, 0, len(market.Options))
	for _, opt := range market.Options {
		options = append(options, MarketOptionResponse{
			ID:          opt.ID.String(),
			MarketID:    opt.MarketID.String(),
			Name:        opt.Name,
			InitialOdds: opt.InitialOdds,
			CurrentOdds: opt.CurrentOdds,
		})
	}

	return MarketResponse{
		ID:        market.ID.String(),
		LeagueID:  market.LeagueID.String(),
		MatchID:   matchID,
		Name:      market.Name,
		Status:    market.Status,
		Options:   options,
		CreatedAt: market.CreatedAt,
		UpdatedAt: market.UpdatedAt,
	}
}
