package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetPendingBonusesUseCase struct {
	leagueRepo repository.LeagueRepository
}

func NewGetPendingBonusesUseCase(leagueRepo repository.LeagueRepository) *GetPendingBonusesUseCase {
	return &GetPendingBonusesUseCase{
		leagueRepo: leagueRepo,
	}
}

func (uc *GetPendingBonusesUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID) ([]*entity.ParticipantBonus, error) {
	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, apperror.ErrLeagueNotFound
	}

	return uc.leagueRepo.GetPendingBonuses(ctx, participant.ID)
}
