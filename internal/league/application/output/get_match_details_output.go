package output

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type GetMatchDetailsOutput struct {
	ID        uuid.UUID
	LeagueID  uuid.UUID
	Slug      string
	Title     string
	StartTime time.Time
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Markets   []entity.Market
}
