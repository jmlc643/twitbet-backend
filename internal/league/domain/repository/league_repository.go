package repository

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type LeagueRepository interface {
	CreateLeagueWithAdminParticipant(ctx context.Context, l *entity.League, p *entity.Participant) error
}