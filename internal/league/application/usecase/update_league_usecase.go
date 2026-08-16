package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateLeagueUseCase struct {
	repo repository.LeagueRepository
}

func NewUpdateLeagueUseCase(repo repository.LeagueRepository) *UpdateLeagueUseCase {
	return &UpdateLeagueUseCase{repo: repo}
}

func (uc *UpdateLeagueUseCase) Execute(ctx context.Context, in input.UpdateLeagueInput) error {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Liga no encontrada")
	}

	if league.OwnerID != in.OwnerID {
		participant, err := uc.repo.GetParticipant(ctx, in.LeagueID, in.OwnerID)
		if err != nil || !participant.IsAdmin {
			return errors.New("Solo el administrador puede modificar la liga")
		}
	}

	if in.Name != "" {
		league.Name = in.Name
	}
	if in.InitialBalance >= 0 {
		league.InitialBalance = in.InitialBalance
	} else {
		return apperror.ErrInvalidInitialBalance
	}
	
	if in.MaxRecharges > 0 {
		league.MaxRecharges = in.MaxRecharges
	}
	league.HideStandings = in.HideStandings
	
	if in.MinBetsToQualify >= 0 {
		league.MinBetsToQualify = in.MinBetsToQualify
	}

	return uc.repo.UpdateLeague(ctx, league)
}
