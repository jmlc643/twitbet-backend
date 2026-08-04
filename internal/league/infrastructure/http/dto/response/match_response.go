package response

import (
	"time"

	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MatchResponse struct {
	ID        string    `json:"id"`
	LeagueID  string    `json:"league_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewMatchResponse(match *entity.Match) MatchResponse {
	return MatchResponse{
		ID:        match.ID.String(),
		LeagueID:  match.LeagueID.String(),
		Title:     match.Title,
		StartTime: match.StartTime,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}
}
