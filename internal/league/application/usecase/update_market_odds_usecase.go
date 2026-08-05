package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type UpdateMarketOddsUseCase struct {
	matchRepo   repository.MatchRepository
	leagueRepo  repository.LeagueRepository
	publisher   port.MarketEventPublisher
}

func NewUpdateMarketOddsUseCase(
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueRepository,
	publisher port.MarketEventPublisher,
) *UpdateMarketOddsUseCase {
	return &UpdateMarketOddsUseCase{
		matchRepo:   matchRepo,
		leagueRepo:  leagueRepo,
		publisher:   publisher,
	}
}

func (u *UpdateMarketOddsUseCase) Execute(ctx context.Context, marketID uuid.UUID, adminID uuid.UUID, optionsOdds map[uuid.UUID]float64) error {
	market, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}

	league, err := u.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}

	if league.AdminID != adminID {
		return errors.New("no autorizado")
	}

	for i, opt := range market.Options {
		if newOdds, ok := optionsOdds[opt.ID]; ok {
			market.Options[i].CurrentOdds = newOdds
		}
	}

	if err := u.matchRepo.UpdateMarket(ctx, market); err != nil {
		return err
	}

	return u.publisher.PublishOddsUpdated(ctx, market.ID, market.Options)
}
