package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type DeleteLeagueUseCase struct {
	repo repository.LeagueRepository
}

func NewDeleteLeagueUseCase(repo repository.LeagueRepository) *DeleteLeagueUseCase {
	return &DeleteLeagueUseCase{repo: repo}
}

func (uc *DeleteLeagueUseCase) Execute(ctx context.Context, in input.DeleteLeagueInput) error {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Liga no encontrada")
	}

	if league.OwnerID != in.OwnerID {
		return errors.New("Solo el administrador puede eliminar la liga")
	}

	return uc.repo.DeleteLeague(ctx, in.LeagueID)
}
