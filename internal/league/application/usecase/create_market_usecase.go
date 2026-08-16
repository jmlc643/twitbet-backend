package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CreateMarketUseCase struct {
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
	marketPublisher port.MarketEventPublisher
}

func NewCreateMarketUseCase(leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository, marketPublisher port.MarketEventPublisher) *CreateMarketUseCase {
	return &CreateMarketUseCase{
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *CreateMarketUseCase) Execute(ctx context.Context, OwnerID, leagueID uuid.UUID, matchID *uuid.UUID, name string, marketType string, options []input.MarketOptionInput) (*entity.Market, error) {
	if leagueID == uuid.Nil && matchID != nil {
		match, err := uc.matchRepo.GetMatchByID(ctx, *matchID)
		if err != nil {
			return nil, err
		}
		leagueID = match.LeagueID
	}

	league, err := uc.leagueRepo.GetLeagueByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	if league.OwnerID != OwnerID {
		participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, OwnerID)
		if err != nil || !participant.IsAdmin {
			return nil, apperror.ErrUnauthorized
		}
	}

	if matchID != nil {
		match, err := uc.matchRepo.GetMatchByID(ctx, *matchID)
		if err != nil {
			return nil, err
		}
		if match.LeagueID != leagueID {
			return nil, errors.New("El partido no pertenece a esta liga")
		}
	}

	var entityOptions []entity.MarketOptionCreate
	for _, opt := range options {
		entityOptions = append(entityOptions, entity.MarketOptionCreate{Name: opt.Name, Odds: opt.Odds})
	}

	market, err := entity.NewMarket(leagueID, matchID, name, marketType, entityOptions)
	if err != nil {
		return nil, err
	}

	if err := uc.matchRepo.CreateMarket(ctx, market); err != nil {
		return nil, err
	}

	_ = uc.marketPublisher.PublishMarketCreated(ctx, *market)

	return market, nil
}