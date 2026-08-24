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

func (v *CombinedBetValidator) ValidateDuplicateTypePerMatch(ctx context.Context, legs []entity.CombinedBetLeg, participant entity.Participant) error {
	seen := make(map[string]bool)
	for _, leg := range legs {
		market, err := v.matchRepo.GetMarketByID(ctx, leg.MarketID)
		if err != nil || market == nil {
			continue
		}
		if market.Type == string(entity.MarketTypeOther) {
			continue
		}
		matchKey := "league"
		if market.MatchID != nil {
			matchKey = market.MatchID.String()
		}
		key := matchKey + ":" + market.Type
		if seen[key] {
			return apperror.ErrDuplicateMarketType
		}
		seen[key] = true
	}
	active, err := v.matchRepo.GetActiveMarketTypesByParticipant(ctx, participant.ID)
	if err != nil {
		return err
	}
	for _, leg := range legs {
		market, err := v.matchRepo.GetMarketByID(ctx, leg.MarketID)
		if err != nil || market == nil {
			continue
		}
		if market.Type == string(entity.MarketTypeOther) {
			continue
		}
		matchKey := "league"
		if market.MatchID != nil {
			matchKey = market.MatchID.String()
		}
		if types, ok := active[matchKey]; ok && types[market.Type] {
			return apperror.ErrDuplicateMarketType
		}
	}
	return nil
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
				if opt.IsBlocked() {
					return nil, apperror.ErrMarketOptionBlocked
				}
				break
			}
		}

		if !found {
			return nil, apperror.ErrMarketOptionNotFound
		}

		leg := entity.NewCombinedBetLeg(uuid.Nil, market.ID, sel.SelectionID, market.MatchID, optionName, sel.AcceptedOdds)
		legs = append(legs, leg)
	}

	if err := v.ValidateDuplicateTypePerMatch(ctx, legs, participant); err != nil {
		return nil, err
	}

	return legs, nil
}
