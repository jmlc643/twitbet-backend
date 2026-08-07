package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type LeagueRepository interface {
	CreateLeagueWithAdminParticipant(ctx context.Context, l *entity.League, p *entity.Participant, t *entity.Transaction) error
	GetLeagueByInviteCode(ctx context.Context, code string) (*entity.League, error)
	GetUserLeagues(ctx context.Context, userID uuid.UUID) ([]entity.LeagueSummary, error)
	IsParticipant(ctx context.Context, leagueID, userID uuid.UUID) (bool, error)
	JoinLeague(ctx context.Context, p *entity.Participant, t *entity.Transaction) error
	GetLeagueByID(ctx context.Context, id uuid.UUID) (*entity.League, error)
	GetLeagueBySlug(ctx context.Context, slug string) (*entity.League, error)
	UpdateLeague(ctx context.Context, league *entity.League) error
	DeleteLeague(ctx context.Context, id uuid.UUID) error
	GetParticipantsByLeagueID(ctx context.Context, id uuid.UUID) ([]entity.ParticipantRanking, error)
	GetParticipant(ctx context.Context, leagueID, userID uuid.UUID) (*entity.Participant, error)
	GetParticipantByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
	UpdateParticipant(ctx context.Context, p *entity.Participant) error
}