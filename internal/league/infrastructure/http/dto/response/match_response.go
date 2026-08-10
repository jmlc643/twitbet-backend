package response

import (
	"time"

	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MatchResponse struct {
	ID        string           `json:"id"`
	LeagueID  string           `json:"league_id"`
	Slug      string           `json:"slug"`
	Title     string           `json:"title"`
	StartTime time.Time        `json:"start_time"`
	Status    string           `json:"status"`
	Markets   []MarketResponse `json:"markets,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func NewMatchResponse(match *entity.Match) MatchResponse {
	var markets []MarketResponse
	if len(match.Markets) > 0 {
		for _, m := range match.Markets {
			markets = append(markets, NewMarketResponse(&m))
		}
	}

	return MatchResponse{
		ID:        match.ID.String(),
		LeagueID:  match.LeagueID.String(),
		Slug:      match.Slug,
		Title:     match.Title,
		StartTime: match.StartTime,
		Status:    match.Status,
		Markets:   markets,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}
}