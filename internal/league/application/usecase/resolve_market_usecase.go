package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type ResolveMarketUseCase struct {
	betRepo         repository.BetRepository
	matchRepo       repository.MatchRepository
	leagueRepo      repository.LeagueRepository
	marketPublisher port.MarketEventPublisher
}

func NewResolveMarketUseCase(betRepo repository.BetRepository, matchRepo repository.MatchRepository, leagueRepo repository.LeagueRepository, marketPublisher port.MarketEventPublisher) *ResolveMarketUseCase {
	return &ResolveMarketUseCase{
		betRepo:         betRepo,
		matchRepo:       matchRepo,
		leagueRepo:      leagueRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *ResolveMarketUseCase) Execute(ctx context.Context, marketID uuid.UUID, winningOptionIDs []uuid.UUID) error {
	if len(winningOptionIDs) == 0 {
		return apperror.ErrInvalidMarketOptions
	}

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

	winningSet := make(map[string]bool, len(winningOptionIDs))
	for _, opt := range market.Options {
		for _, wID := range winningOptionIDs {
			if opt.ID == wID {
				winningSet[wID.String()] = true
				break
			}
		}
	}
	for _, wID := range winningOptionIDs {
		if !winningSet[wID.String()] {
			return apperror.ErrMarketOptionNotFound
		}
	}

	err = uc.betRepo.ResolveMarketAtomic(ctx, marketID, winningOptionIDs)
	if err != nil {
		return err
	}

	_ = uc.marketPublisher.PublishMarketResolved(ctx, marketID, market.LeagueID, winningOptionIDs)

	_ = uc.marketPublisher.PublishMarketStatusChanged(ctx, marketID, string(entity.MarketStatusResolved))

	if market.MatchID != nil {
		err = uc.matchRepo.UpdateMatchStatusAtomic(ctx, *market.MatchID, "FINISHED")
		if err == nil {
			_ = uc.marketPublisher.PublishMatchStatusChanged(ctx, *market.MatchID, "FINISHED")
		}
	}

	return nil
}