package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type CombinedBetValidator struct {
	matchRepo repository.MatchRepository
}

func NewCombinedBetValidator(matchRepo repository.MatchRepository) *CombinedBetValidator {
	return &CombinedBetValidator{
		matchRepo: matchRepo,
	}
}

func (v *CombinedBetValidator) Validate(ctx context.Context, selections []valueobject.Selection, participant entity.Participant) ([]entity.CombinedBetLeg, error) {
	if len(selections) < 2 {
		return nil, apperror.ErrInvalidBetAmount
	}

	marketIDs := make(map[uuid.UUID]bool)
	var legs []entity.CombinedBetLeg

	for _, sel := range selections {
		if marketIDs[sel.MarketID] {
			return nil, apperror.ErrInvalidMarketOptions
		}
		marketIDs[sel.MarketID] = true

		market, err := v.matchRepo.GetMarketByID(ctx, sel.MarketID)
		if err != nil {
			return nil, err
		}
		if market == nil {
			return nil, apperror.ErrMarketNotFound
		}

		if market.LeagueID != participant.LeagueID {
			return nil, apperror.ErrInvalidMarketOptions
		}

		if entity.NotResolvableMarketStatus(market) {
			if market.Status != string(entity.MarketStatusActive) {
				return nil, apperror.ErrMarketNotActive
			}
		}

		var optionName string
		found := false
		for _, opt := range market.Options {
			if opt.ID == sel.SelectionID {
				optionName = opt.Name
				found = true
				break
			}
		}

		if !found {
			return nil, apperror.ErrMarketOptionNotFound
		}

		leg := entity.NewCombinedBetLeg(uuid.Nil, market.ID, sel.SelectionID, market.MatchID, optionName, sel.AcceptedOdds)
		legs = append(legs, leg)
	}

	return legs, nil
}
