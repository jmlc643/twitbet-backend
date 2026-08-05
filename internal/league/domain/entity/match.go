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
	Slug      string
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
	slug := generateMatchSlug(title)

	return &Match{
		ID:        uuid.New(),
		LeagueID:  leagueID,
		Title:     title,
		Slug:      slug,
		StartTime: startTime,
		Status:    "SCHEDULED",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func generateMatchSlug(title string) string {
	var slug string
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			slug += string(c)
		} else if c >= 'A' && c <= 'Z' {
			slug += string(c + 32)
		} else if c == ' ' || c == '-' {
			if len(slug) > 0 && slug[len(slug)-1] != '-' {
				slug += "-"
			}
		}
	}
	if len(slug) > 0 && slug[len(slug)-1] == '-' {
		slug = slug[:len(slug)-1]
	}
	
	uuidPart := uuid.New().String()[:6]
	return slug + "-" + uuidPart
}


