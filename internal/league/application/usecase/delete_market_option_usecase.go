package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type DeleteMarketOptionUseCase struct {
	matchRepo       repository.MatchRepository
	leagueRepo      repository.LeagueRepository
	marketPublisher port.MarketEventPublisher
}

func NewDeleteMarketOptionUseCase(matchRepo repository.MatchRepository, leagueRepo repository.LeagueRepository, marketPublisher port.MarketEventPublisher) *DeleteMarketOptionUseCase {
	return &DeleteMarketOptionUseCase{matchRepo: matchRepo, leagueRepo: leagueRepo, marketPublisher: marketPublisher}
}

func (uc *DeleteMarketOptionUseCase) Execute(ctx context.Context, marketID uuid.UUID, optionID uuid.UUID, requesterID uuid.UUID) error {
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
	found := false
	for _, opt := range market.Options {
		if opt.ID == optionID {
			found = true
			break
		}
	}
	if !found {
		return apperror.ErrMarketOptionNotFound
	}
	if len(market.Options) <= 2 {
		return apperror.ErrLastMarketOption
	}
	hasBets, err := uc.matchRepo.HasActiveBetsForOption(ctx, optionID)
	if err != nil {
		return err
	}
	if hasBets {
		return apperror.ErrMarketOptionHasBets
	}
	if err := uc.matchRepo.DeleteMarketOption(ctx, marketID, optionID); err != nil {
		return err
	}
	remaining := []entity.MarketOption{}
	for _, opt := range market.Options {
		if opt.ID != optionID {
			remaining = append(remaining, opt)
		}
	}
	_ = uc.marketPublisher.PublishMarketOptionDeleted(ctx, marketID, optionID)
	_ = uc.marketPublisher.PublishOddsUpdated(ctx, marketID, remaining)
	return nil
}
