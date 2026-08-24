package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type DeleteMarketUseCase struct {
	matchRepo       repository.MatchRepository
	leagueRepo      repository.LeagueRepository
	marketPublisher port.MarketEventPublisher
}

func NewDeleteMarketUseCase(matchRepo repository.MatchRepository, leagueRepo repository.LeagueRepository, marketPublisher port.MarketEventPublisher) *DeleteMarketUseCase {
	return &DeleteMarketUseCase{matchRepo: matchRepo, leagueRepo: leagueRepo, marketPublisher: marketPublisher}
}

func (uc *DeleteMarketUseCase) Execute(ctx context.Context, marketID uuid.UUID, requesterID uuid.UUID) error {
	market, err := uc.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}
	if market == nil {
		return apperror.ErrMarketNotFound
	}
	if entity.NotResolvableMarketStatus(market) {
		return apperror.ErrMarketNotActive
	}
	league, err := uc.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}
	if league.OwnerID != requesterID {
		p, err := uc.leagueRepo.GetParticipant(ctx, market.LeagueID, requesterID)
		if err != nil || p == nil || !p.IsAdmin {
			return apperror.ErrUnauthorized
		}
	}
	hasBets, err := uc.matchRepo.HasActiveBetsForMarket(ctx, marketID)
	if err != nil {
		return err
	}
	if hasBets {
		return apperror.ErrMarketHasBets
	}
	if err := uc.matchRepo.DeleteMarket(ctx, marketID); err != nil {
		return err
	}
	_ = uc.marketPublisher.PublishMarketDeleted(ctx, marketID)
	_ = uc.marketPublisher.PublishMarketStatusChanged(ctx, marketID, "DELETED")
	return nil
}
