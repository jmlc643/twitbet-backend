package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type RemoveAdminUseCase struct {
	repo repository.LeagueRepository
}

func NewRemoveAdminUseCase(repo repository.LeagueRepository) *RemoveAdminUseCase {
	return &RemoveAdminUseCase{repo: repo}
}

func (uc *RemoveAdminUseCase) Execute(ctx context.Context, in input.RemoveAdminInput) error {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("liga no encontrada")
	}

	if league.OwnerID != in.OwnerID {
		return errors.New("solo el propietario puede remover administradores")
	}

	if league.OwnerID == in.ParticipantID {
		return errors.New("no se puede remover el rol de admin al propietario")
	}

	participant, err := uc.repo.GetParticipant(ctx, in.LeagueID, in.ParticipantID)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.New("participante no encontrado en la liga")
	}

	participant.IsAdmin = false

	return uc.repo.UpdateParticipant(ctx, participant)
}
