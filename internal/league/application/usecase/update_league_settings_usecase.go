package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateLeagueSettingsUseCase struct {
	repo repository.LeagueRepository
}

func NewUpdateLeagueSettingsUseCase(repo repository.LeagueRepository) *UpdateLeagueSettingsUseCase {
	return &UpdateLeagueSettingsUseCase{repo: repo}
}

func (uc *UpdateLeagueSettingsUseCase) Execute(ctx context.Context, in input.UpdateLeagueSettingsInput) error {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("liga no encontrada")
	}

	if league.AdminID != in.AdminID {
		return errors.New("solo el administrador puede modificar los ajustes de la liga")
	}

	league.HideStandings = !in.IsRankingVisible

	return uc.repo.UpdateLeague(ctx, league)
}
