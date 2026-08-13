package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type PlaceBetUseCase struct {
	betRepo         repository.BetRepository
	leagueRepo      repository.LeagueRepository
	matchRepo       repository.MatchRepository
	marketPublisher port.MarketEventPublisher
}

func NewPlaceBetUseCase(betRepo repository.BetRepository, leagueRepo repository.LeagueRepository, matchRepo repository.MatchRepository, marketPublisher port.MarketEventPublisher) *PlaceBetUseCase {
	return &PlaceBetUseCase{
		betRepo:         betRepo,
		leagueRepo:      leagueRepo,
		matchRepo:       matchRepo,
		marketPublisher: marketPublisher,
	}
}

func (uc *PlaceBetUseCase) Execute(ctx context.Context, userID, leagueID, marketID, marketOptionID uuid.UUID, amount float64, bonusID *uuid.UUID) (*entity.Bet, error) {
	isParticipant, err := uc.leagueRepo.IsParticipant(ctx, leagueID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, apperror.ErrLeagueNotFound
	}

	participants, err := uc.leagueRepo.GetParticipantsByLeagueID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	var participantID uuid.UUID
	var found bool
	for _, p := range participants {
		if p.UserID == userID {
			participantID = p.ParticipantID
			found = true
			break
		}
	}
	if !found {
		return nil, apperror.ErrLeagueNotFound
	}

	market, err := uc.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return nil, err
	}
	if market == nil || market.Status != string(entity.MarketStatusOpen) {
		return nil, apperror.ErrMarketNotActive
	}

	var optionOdds float64
	var optionFound bool
	for _, opt := range market.Options {
		if opt.ID == marketOptionID {
			if opt.IsBlocked() {
				return nil, apperror.ErrMarketOptionBlocked
			}
			optionOdds = opt.CurrentOdds
			optionFound = true
			break
		}
	}
	if !optionFound {
		return nil, apperror.ErrMarketOptionNotFound
	}

	bet, err := entity.NewBet(participantID, marketOptionID, amount, optionOdds, bonusID)
	if err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:        uuid.New(),
		LeagueID:  leagueID,
		UserID:    userID,
		Amount:    amount, 
		Type:      entity.TransactionTypeBet,
		CreatedAt: time.Now().UTC(),
	}

	err = uc.betRepo.PlaceBetAtomic(ctx, bet, transaction)
	if err != nil {
		return nil, err
	}

	const K = 10000.0
	var totalVirtualVolume float64
	virtualVolumes := make([]float64, len(market.Options))
	for i, opt := range market.Options {
		vol := K / opt.CurrentOdds
		if opt.ID == marketOptionID {
			vol += amount
		}
		virtualVolumes[i] = vol
		totalVirtualVolume += vol
	}

	for i := range market.Options {
		market.Options[i].CurrentOdds = totalVirtualVolume / virtualVolumes[i]
	}

	_ = uc.matchRepo.UpdateMarket(ctx, market)

	_ = uc.marketPublisher.PublishOddsUpdated(ctx, market.ID, market.Options)

	_ = uc.marketPublisher.PublishParticipantBalanceUpdated(ctx, participantID, leagueID, userID)

	go func(bID uuid.UUID) {
		time.Sleep(4 * time.Second)
		bgCtx := context.Background()
		_ = uc.betRepo.UpdateBetStatus(bgCtx, bID, entity.BetStatusAccepted)
	}(bet.ID)

	return bet, nil
}
