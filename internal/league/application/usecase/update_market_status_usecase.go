package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateMarketStatusUseCase struct {
	matchRepo   repository.MatchRepository
	leagueRepo  repository.LeagueRepository
	publisher   port.MarketEventPublisher
}

func NewUpdateMarketStatusUseCase(
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueRepository,
	publisher port.MarketEventPublisher,
) *UpdateMarketStatusUseCase {
	return &UpdateMarketStatusUseCase{
		matchRepo:   matchRepo,
		leagueRepo:  leagueRepo,
		publisher:   publisher,
	}
}

func (u *UpdateMarketStatusUseCase) Execute(ctx context.Context, marketID uuid.UUID, OwnerID uuid.UUID, newStatus string) error {
	if newStatus != string(entity.MarketStatusActive) && newStatus != string(entity.MarketStatusSuspended) {
		return apperror.ErrInvalidMarketName
	}

	market, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}

	league, err := u.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}

	if league.OwnerID != OwnerID {
		participant, err := u.leagueRepo.GetParticipant(ctx, market.LeagueID, OwnerID)
		if err != nil || !participant.IsAdmin {
			return errors.New("No autorizado")
		}
	}

	market.Status = newStatus

	if err := u.matchRepo.UpdateMarket(ctx, market); err != nil {
		return err
	}

	return u.publisher.PublishMarketStatusChanged(ctx, market.ID, newStatus)
}
