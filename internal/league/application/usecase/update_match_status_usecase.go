package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateMatchStatusUseCase struct {
	matchRepo  repository.MatchRepository
	leagueRepo repository.LeagueRepository
	publisher  port.MarketEventPublisher
}

func NewUpdateMatchStatusUseCase(matchRepo repository.MatchRepository, leagueRepo repository.LeagueRepository, publisher port.MarketEventPublisher) *UpdateMatchStatusUseCase {
	return &UpdateMatchStatusUseCase{
		matchRepo:  matchRepo,
		leagueRepo: leagueRepo,
		publisher:  publisher,
	}
}

func (uc *UpdateMatchStatusUseCase) Execute(ctx context.Context, matchID, userID uuid.UUID, newStatus string) error {
	match, err := uc.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return apperror.ErrMatchNotFound
	}

	participant, err := uc.leagueRepo.GetParticipant(ctx, match.LeagueID, userID)
	if err != nil {
		return err
	}
	if participant == nil || !participant.IsAdmin {
		return apperror.ErrUnauthorized
	}

	if err := uc.matchRepo.UpdateMatchStatusAtomic(ctx, matchID, newStatus); err != nil {
		return err
	}

	return uc.publisher.PublishMatchStatusChanged(ctx, matchID, newStatus)
}
