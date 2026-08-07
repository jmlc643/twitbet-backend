package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
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

func (uc *ResolveMarketUseCase) Execute(ctx context.Context, marketID, winningOptionID uuid.UUID) error {
	market, err := uc.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}
	if market == nil {
		return apperror.ErrMarketNotFound
	}
	if market.Status == "RESOLVED" {
		return apperror.ErrMarketNotActive
	}

	var validOption bool
	for _, opt := range market.Options {
		if opt.ID == winningOptionID {
			validOption = true
			break
		}
	}
	if !validOption {
		return apperror.ErrMarketOptionNotFound
	}

	err = uc.betRepo.ResolveMarketAtomic(ctx, marketID, winningOptionID)
	if err != nil {
		return err
	}

	_ = uc.marketPublisher.PublishMarketResolved(ctx, marketID, market.LeagueID, winningOptionID)

	_ = uc.marketPublisher.PublishMarketStatusChanged(ctx, marketID, "RESOLVED")

	err = uc.matchRepo.UpdateMatchStatusAtomic(ctx, *market.MatchID, "FINISHED")
	if err == nil {
		_ = uc.marketPublisher.PublishMatchStatusChanged(ctx, *market.MatchID, "FINISHED")
	}

	return nil
}
