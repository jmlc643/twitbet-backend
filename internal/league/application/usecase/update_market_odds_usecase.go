package usecase

import (
	"context"
	"errors"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/service"
)

type UpdateMarketOddsUseCase struct {
	matchRepo    repository.MatchRepository
	leagueRepo   repository.LeagueRepository
	publisher    port.MarketEventPublisher
	integritySvc *service.MarketIntegrityService
}

func NewUpdateMarketOddsUseCase(
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueRepository,
	publisher port.MarketEventPublisher,
) *UpdateMarketOddsUseCase {
	return &UpdateMarketOddsUseCase{
		matchRepo:    matchRepo,
		leagueRepo:   leagueRepo,
		publisher:    publisher,
		integritySvc: service.NewMarketIntegrityService(),
	}
}

func (u *UpdateMarketOddsUseCase) Execute(ctx context.Context, marketID uuid.UUID, OwnerID uuid.UUID, optionsOdds map[uuid.UUID]float64) (map[uuid.UUID]float64, error) {
	market, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return nil, err
	}

	league, err := u.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return nil, err
	}

	if league.OwnerID != OwnerID {
		participant, err := u.leagueRepo.GetParticipant(ctx, market.LeagueID, OwnerID)
		if err != nil || !participant.IsAdmin {
			return nil, errors.New("No autorizado")
		}
	}

	mode := os.Getenv("ODDS_UPDATE_MODE")
	if mode == "" {
		mode = "REBALANCE"
	}
	minMargin, _ := strconv.ParseFloat(os.Getenv("MIN_MARGIN"), 64)
	if minMargin == 0 {
		minMargin = 0.03
	}
	cooldownPct, _ := strconv.ParseFloat(os.Getenv("COOLDOWN_DELTA_PCT"), 64)
	if cooldownPct == 0 {
		cooldownPct = 15.0
	}
	cooldownSecs, _ := strconv.Atoi(os.Getenv("COOLDOWN_SECONDS"))
	if cooldownSecs == 0 {
		cooldownSecs = 45
	}

	resultOdds := make(map[uuid.UUID]float64)
	
	if mode == "REBALANCE" {
		rebalanced, err := u.integritySvc.CalculateRebalance(market.Options, optionsOdds, minMargin)
		if err != nil {
			return nil, err
		}
		resultOdds = rebalanced
	} else {
		tempOptions := make([]entity.MarketOption, len(market.Options))
		copy(tempOptions, market.Options)
		
		for i, opt := range tempOptions {
			if newOdd, ok := optionsOdds[opt.ID]; ok {
				tempOptions[i].CurrentOdds = newOdd
			}
		}

		overround := u.integritySvc.CalculateOverround(tempOptions)
		if overround < 1.0 {
			return nil, apperror.ErrArbitrageMarket
		}
		if overround < 1.0+minMargin {
			return nil, &apperror.RebalanceError{
				TouchedMass:    overround,
				MaxTouchedMass: 1.0 + minMargin,
				Hint:           "Las cuotas no alcanzan el overround objetivo y el modo es REJECT",
			}
		}
		
		for _, opt := range tempOptions {
			resultOdds[opt.ID] = opt.CurrentOdds
		}
	}

	tempOptions := make([]entity.MarketOption, len(market.Options))
	copy(tempOptions, market.Options)
	for i, opt := range tempOptions {
		if newOdd, ok := resultOdds[opt.ID]; ok {
			tempOptions[i].CurrentOdds = newOdd
		}
	}
	if u.integritySvc.CalculateOverround(tempOptions) < 1.0 {
		return nil, apperror.ErrArbitrageMarket
	}

	triggeredCooldown := false
	var histories []entity.MarketOddsHistory
	
	for i, opt := range market.Options {
		if newOdd, ok := resultOdds[opt.ID]; ok {
			delta := math.Abs(newOdd - opt.CurrentOdds) / opt.CurrentOdds * 100.0
			if delta > cooldownPct {
				triggeredCooldown = true
			}
			
			oldOdds := opt.CurrentOdds
			if newOdd != opt.CurrentOdds {
				histories = append(histories, *entity.NewMarketOddsHistory(
					opt.ID,
					&oldOdds,
					newOdd,
					&OwnerID,
					string(entity.OddsChangeReasonManual),
				))
			}
			
			market.Options[i].CurrentOdds = newOdd
		}
	}

	market.Seq++

	if triggeredCooldown {
		market.Status = string(entity.MarketStatusSuspended)
		reason := "AUTO_COOLDOWN"
		market.SuspendReason = &reason
		
		go func(mID uuid.UUID) {
			time.Sleep(time.Duration(cooldownSecs) * time.Second)
			ctx := context.Background()
			if m, err := u.matchRepo.GetMarketByID(ctx, mID); err == nil && m.Status == string(entity.MarketStatusSuspended) {
				m.Status = string(entity.MarketStatusOpen)
				m.SuspendReason = nil
				m.Seq++
				_ = u.matchRepo.UpdateMarket(ctx, m)
				_ = u.publisher.PublishOddsUpdated(ctx, m.ID, m.Options) // Or publish market reopened event
			}
		}(market.ID)
	}

	if err := u.matchRepo.UpdateMarketAndHistory(ctx, market, histories); err != nil {
		return nil, err
	}

	_ = u.publisher.PublishOddsUpdated(ctx, market.ID, market.Options)

	return resultOdds, nil
}
