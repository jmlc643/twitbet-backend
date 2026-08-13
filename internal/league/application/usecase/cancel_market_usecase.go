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
	marketPublisher port.MarketEventPublisher
}

func NewCancelMarketUseCase(betRepo repository.BetRepository, matchRepo repository.MatchRepository, marketPublisher port.MarketEventPublisher) *CancelMarketUseCase {
	return &CancelMarketUseCase{
		betRepo:         betRepo,
		matchRepo:       matchRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *CancelMarketUseCase) Execute(ctx context.Context, marketID uuid.UUID, reason string) error {
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

	err = uc.betRepo.CancelMarketAtomic(ctx, marketID, reason)
	if err != nil {
		return err
	}

	_ = uc.marketPublisher.PublishMarketStatusChanged(ctx, marketID, "CANCELLED")

	return nil
}
