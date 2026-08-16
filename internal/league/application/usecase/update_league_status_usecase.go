package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateLeagueStatusUseCase struct {
	repo repository.LeagueRepository
}

func NewUpdateLeagueStatusUseCase(repo repository.LeagueRepository) *UpdateLeagueStatusUseCase {
	return &UpdateLeagueStatusUseCase{repo: repo}
}

func (uc *UpdateLeagueStatusUseCase) Execute(ctx context.Context, in input.UpdateLeagueStatusInput) error {
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

	if league.Status == "FINALIZED" {
		return errors.New("La liga ya se encuentra finalizada y es un estado definitivo")
	}

	if in.Status != "ACTIVE" && in.Status != "FINALIZED" {
		return errors.New("Estado inválido")
	}

	league.Status = in.Status

	return uc.repo.UpdateLeague(ctx, league)
}
