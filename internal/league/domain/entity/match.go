package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type Match struct {
	ID        uuid.UUID
	LeagueID  uuid.UUID
	Title     string
	StartTime time.Time
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMatch(leagueID uuid.UUID, title string, startTime time.Time) (*Match, error) {
	if title == "" {
		return nil, apperror.ErrInvalidMatchTitle
	}

	now := time.Now().UTC()
	return &Match{
		ID:        uuid.New(),
		LeagueID:  leagueID,
		Title:     title,
		StartTime: startTime,
		Status:    "SCHEDULED",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
