package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetParticipantMeUseCase struct {
	leagueRepo repository.LeagueRepository
}

func NewGetParticipantMeUseCase(leagueRepo repository.LeagueRepository) *GetParticipantMeUseCase {
	return &GetParticipantMeUseCase{leagueRepo: leagueRepo}
}

func (uc *GetParticipantMeUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID) (*entity.Participant, error) {
	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, apperror.ErrUnauthorized
	}

	return participant, nil
}
