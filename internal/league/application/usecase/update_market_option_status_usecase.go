package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateMarketOptionStatusUseCase struct {
	matchRepo  repository.MatchRepository
	leagueRepo repository.LeagueRepository
	publisher  port.MarketEventPublisher
}

func NewUpdateMarketOptionStatusUseCase(
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueRepository,
	publisher port.MarketEventPublisher,
) *UpdateMarketOptionStatusUseCase {
	return &UpdateMarketOptionStatusUseCase{
		matchRepo:  matchRepo,
		leagueRepo: leagueRepo,
		publisher:  publisher,
	}
}

func (u *UpdateMarketOptionStatusUseCase) Execute(ctx context.Context, marketID uuid.UUID, optionID uuid.UUID, ownerID uuid.UUID, newStatus string) error {
	if newStatus != string(entity.MarketOptionStatusActive) && newStatus != string(entity.MarketOptionStatusBlocked) {
		return apperror.ErrInvalidMarketOptionStatus
	}

	market, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}
	if market == nil {
		return apperror.ErrMarketNotFound
	}
	if entity.NotResolvableMarketStatus(market) {
		return apperror.ErrMarketNotActive
	}

	var optionFound bool
	for _, opt := range market.Options {
		if opt.ID == optionID {
			optionFound = true
			break
		}
	}
	if !optionFound {
		return apperror.ErrMarketOptionNotFound
	}

	league, err := u.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}
	if league.OwnerID != ownerID {
		participant, err := u.leagueRepo.GetParticipant(ctx, market.LeagueID, ownerID)
		if err != nil || !participant.IsAdmin {
			return apperror.ErrUnauthorized
		}
	}

	if err := u.matchRepo.UpdateMarketOptionStatus(ctx, marketID, optionID, newStatus); err != nil {
		return err
	}

	updated, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}

	return u.publisher.PublishMarketOptionsUpdated(ctx, *updated)
}