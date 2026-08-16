package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type AssignAdminUseCase struct {
	repo repository.LeagueRepository
}

func NewAssignAdminUseCase(repo repository.LeagueRepository) *AssignAdminUseCase {
	return &AssignAdminUseCase{repo: repo}
}

func (uc *AssignAdminUseCase) Execute(ctx context.Context, in input.AssignAdminInput) error {
	league, err := uc.repo.GetLeagueByID(ctx, in.LeagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Liga no encontrada")
	}

	if league.OwnerID != in.OwnerID {
		return errors.New("Solo el propietario puede asignar administradores")
	}

	participant, err := uc.repo.GetParticipant(ctx, in.LeagueID, in.ParticipantID)
	if err != nil {
		return err
	}
	if participant == nil {
		return errors.New("Participante no encontrado en la liga")
	}

	participant.IsAdmin = true

	return uc.repo.UpdateParticipant(ctx, participant)
}
