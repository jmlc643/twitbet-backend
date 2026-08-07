package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetUserBetsUseCase struct {
	betRepo    repository.BetRepository
	leagueRepo repository.LeagueRepository
}

func NewGetUserBetsUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository) *GetUserBetsUseCase {
	return &GetUserBetsUseCase{
		betRepo:    betRepo,
		leagueRepo: leagueRepo,
	}
}

func (uc *GetUserBetsUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID, status *string, startDate, endDate *time.Time, page, limit int) ([]entity.BetDetail, int64, error) {
	participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, 0, err
	}
	if participant == nil {
		return nil, 0, apperror.ErrUnauthorized
	}

	var betStatus *entity.BetStatus
	if status != nil && *status != "" {
		st := entity.BetStatus(*status)
		betStatus = &st
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	bets, total, err := uc.betRepo.GetBetsByParticipantID(ctx, participant.ID, betStatus, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if bets == nil {
		bets = []entity.BetDetail{}
	}

	for i := range bets {
		if bets[i].Status == entity.BetStatusAccepted || bets[i].Status == entity.BetStatusPending {
			ca := bets[i].Amount * 0.95
			bets[i].CashoutAmount = &ca
		}
	}

	return bets, total, nil
}
