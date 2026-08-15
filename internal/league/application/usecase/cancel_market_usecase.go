package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CancelMarketUseCase struct {
	betRepo         repository.BetRepository
	matchRepo       repository.MatchRepository
	leagueRepo      repository.LeagueRepository
	marketPublisher port.MarketEventPublisher
}

func NewCancelMarketUseCase(betRepo repository.BetRepository, matchRepo repository.MatchRepository, leagueRepo repository.LeagueRepository, marketPublisher port.MarketEventPublisher) *CancelMarketUseCase {
	return &CancelMarketUseCase{
		betRepo:         betRepo,
		matchRepo:       matchRepo,
		leagueRepo:      leagueRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *CancelMarketUseCase) Execute(ctx context.Context, marketID uuid.UUID, requesterID uuid.UUID, reason string) error {
	market, err := uc.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}
	if market == nil {
		return apperror.ErrMarketNotFound
	}
	if market.Status != string(entity.MarketStatusOpen) && market.Status != string(entity.MarketStatusSuspended) {
		return apperror.ErrMarketNotActive
	}

	league, err := uc.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}
	if league.OwnerID != requesterID {
		participant, err := uc.leagueRepo.GetParticipant(ctx, market.LeagueID, requesterID)
		if err != nil || !participant.IsAdmin {
			return apperror.ErrUnauthorized
		}
	}

	err = uc.betRepo.CancelMarketAtomic(ctx, marketID, reason)
	if err != nil {
		return err
	}

	_ = uc.marketPublisher.PublishMarketStatusChanged(ctx, marketID, "CANCELLED")

	return nil
}
