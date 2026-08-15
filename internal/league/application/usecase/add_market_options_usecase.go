package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type AddMarketOptionsUseCase struct {
	matchRepo  repository.MatchRepository
	leagueRepo repository.LeagueRepository
	publisher  port.MarketEventPublisher
}

func NewAddMarketOptionsUseCase(
	matchRepo repository.MatchRepository,
	leagueRepo repository.LeagueRepository,
	publisher port.MarketEventPublisher,
) *AddMarketOptionsUseCase {
	return &AddMarketOptionsUseCase{
		matchRepo:  matchRepo,
		leagueRepo: leagueRepo,
		publisher:  publisher,
	}
}

func (u *AddMarketOptionsUseCase) Execute(ctx context.Context, marketID uuid.UUID, ownerID uuid.UUID, options []input.MarketOptionInput) error {
	if len(options) == 0 {
		return apperror.ErrInvalidMarketOptions
	}

	market, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}
	if market == nil {
		return apperror.ErrMarketNotFound
	}
	if entity.NotResolvableMarketStatus(market) {
		return apperror.ErrMarketNotActive
	}

	league, err := u.leagueRepo.GetLeagueByID(ctx, market.LeagueID)
	if err != nil {
		return err
	}
	if league.OwnerID != ownerID {
		participant, err := u.leagueRepo.GetParticipant(ctx, market.LeagueID, ownerID)
		if err != nil || !participant.IsAdmin {
			return apperror.ErrUnauthorized
		}
	}

	newOptions := make([]entity.MarketOption, 0, len(options))
	for _, opt := range options {
		if opt.Name == "" || opt.Odds <= 1 {
			return apperror.ErrInvalidMarketOptions
		}
		newOptions = append(newOptions, *entity.NewMarketOption(marketID, opt.Name, opt.Odds))
	}

	if err := u.matchRepo.AddMarketOptions(ctx, marketID, newOptions); err != nil {
		return err
	}

	updated, err := u.matchRepo.GetMarketByID(ctx, marketID)
	if err != nil {
		return err
	}

	return u.publisher.PublishMarketOptionsUpdated(ctx, *updated)
}