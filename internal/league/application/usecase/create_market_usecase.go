package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CreateMarketUseCase struct {
	leagueRepo repository.LeagueRepository
	matchRepo  repository.MatchRepository
}

func NewCreateMarketUseCase(leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository) *CreateMarketUseCase {
	return &CreateMarketUseCase{
		leagueRepo: leagueRepo,
		matchRepo:  matchRepo,
	}
}

func (uc *CreateMarketUseCase) Execute(ctx context.Context, OwnerID, leagueID uuid.UUID, matchID *uuid.UUID, name string, options []input.MarketOptionInput) (*entity.Market, error) {
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
		return nil, errors.New("el usuario no es administrador de la liga")
	}

	if matchID != nil {
		match, err := uc.matchRepo.GetMatchByID(ctx, *matchID)
		if err != nil {
			return nil, err
		}
		if match.LeagueID != leagueID {
			return nil, errors.New("el partido no pertenece a esta liga")
		}
	}

	var entityOptions []entity.MarketOptionCreate
	for _, opt := range options {
		entityOptions = append(entityOptions, entity.MarketOptionCreate{Name: opt.Name, Odds: opt.Odds})
	}

	market, err := entity.NewMarket(leagueID, matchID, name, entityOptions)
	if err != nil {
		return nil, err
	}

	if err := uc.matchRepo.CreateMarket(ctx, market); err != nil {
		return nil, err
	}

	return market, nil
}


