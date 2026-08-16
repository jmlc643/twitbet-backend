package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/service"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type PlaceCombinedBetUseCase struct {
	combinedBetRepo repository.CombinedBetRepository
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
	marketPublisher port.MarketEventPublisher
	validator       *service.CombinedBetValidator
}

func NewPlaceCombinedBetUseCase(
	combinedBetRepo repository.CombinedBetRepository,
	leagueRepo repository.LeagueRepository,
	matchRepo repository.MatchRepository,
	marketPublisher port.MarketEventPublisher,
) *PlaceCombinedBetUseCase {
	return &PlaceCombinedBetUseCase{
		combinedBetRepo: combinedBetRepo,
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
		marketPublisher: marketPublisher,
		validator:       service.NewCombinedBetValidator(matchRepo),
	}
}

func (uc *PlaceCombinedBetUseCase) Execute(ctx context.Context, userID uuid.UUID, req input.PlaceCombinedBetInput) (*output.CombinedBetOutput, error) {
	participant, err := uc.leagueRepo.GetParticipant(ctx, req.LeagueID, userID)
	if err != nil {
		return nil, err
	}
	if participant == nil {
		return nil, apperror.ErrUnauthorized
	}

	var participantEntity = entity.Participant{
		ID:                participant.ID,
		UserID:            participant.UserID,
		LeagueID:          participant.LeagueID,
		Balance:           participant.Balance,
		IsAdmin:           participant.IsAdmin,
		RechargesConsumed: participant.RechargesConsumed,
		JoinedAt:          participant.JoinedAt,
	}

	selections := make([]valueobject.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		selections[i] = valueobject.Selection{
			MarketID:    sel.MarketID,
			SelectionID: sel.SelectionID,
		}
	}

	legs, err := uc.validator.Validate(ctx, selections, participantEntity)
	if err != nil {
		return nil, err
	}

	bet, err := entity.NewCombinedBet(participant.ID, req.LeagueID, req.Stake, req.UseBonus, req.BonusID, legs)
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:        uuid.New(),
		LeagueID:  req.LeagueID,
		UserID:    userID,
		Amount:    req.Stake,
		Type:      entity.TransactionTypeBet,
		CreatedAt: time.Now().UTC(),
	}

	err = uc.combinedBetRepo.PlaceCombinedBetAtomic(ctx, bet, transaction)
	if err != nil {
		return nil, err
	}

	_ = uc.marketPublisher.PublishParticipantBalanceUpdated(ctx, participant.ID, req.LeagueID, userID)

	go func(bID uuid.UUID) {
		time.Sleep(4 * time.Second)
		bgCtx := context.Background()
		_ = uc.combinedBetRepo.UpdateStatus(bgCtx, bID, string(valueobject.CombinedStatusAccepted))
	}(bet.ID)

	legOutputs := make([]output.LegOutput, len(bet.Legs))
	for i, leg := range bet.Legs {
		legOutputs[i] = output.LegOutput{
			ID:              leg.ID,
			MarketID:        leg.MarketID,
			MatchID:         leg.MatchID,
			SelectionName:   leg.SelectionName,
			OddsAtPlacement: leg.OddsAtPlacement,
			Status:          string(leg.Status),
			SettledAt:       leg.SettledAt,
		}
	}

	return &output.CombinedBetOutput{
		ID:               bet.ID,
		UserID:           userID,
		LeagueID:         bet.LeagueID,
		Stake:            bet.Stake,
		UseBonus:         bet.UseBonus,
		TotalOdds:        bet.TotalOdds,
		PotentialWin:     bet.PotentialWin,
		Status:           string(bet.Status),
		CashoutValue:     bet.CashoutValue,
		CashoutExpiresAt: bet.CashoutExpiresAt,
		CreatedAt:        bet.CreatedAt,
		SettledAt:        bet.SettledAt,
		Legs:             legOutputs,
	}, nil
}
