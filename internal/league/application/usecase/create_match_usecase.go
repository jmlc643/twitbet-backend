package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type CreateMatchUseCase struct {
	leagueRepo repository.LeagueRepository
	matchRepo  repository.MatchRepository
}

func NewCreateMatchUseCase(leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository) *CreateMatchUseCase {
	return &CreateMatchUseCase{
		leagueRepo: leagueRepo,
		matchRepo:  matchRepo,
	}
}

func (uc *CreateMatchUseCase) Execute(ctx context.Context, adminID, leagueID uuid.UUID, title string, startTime time.Time) (*entity.Match, error) {
	league, err := uc.leagueRepo.GetLeagueByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	if league.AdminID != adminID {
		return nil, errors.New("el usuario no es administrador de la liga")
	}

	match, err := entity.NewMatch(leagueID, title, startTime)
	if err != nil {
		return nil, err
	}

	if err := uc.matchRepo.CreateMatch(ctx, match); err != nil {
		return nil, err
	}

	return match, nil
}
