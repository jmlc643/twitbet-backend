package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
)

type CombinedBet struct {
	ID                 uuid.UUID
	ParticipantID      uuid.UUID
	LeagueID           uuid.UUID
	Stake              float64
	UseBonus           bool
	TotalOdds          float64
	PotentialWin       float64
	Status             valueobject.CombinedStatus
	CashoutValue       *float64
	CashoutExpiresAt   *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SettledAt          *time.Time
	Legs               []CombinedBetLeg
	ParticipantBonusID *uuid.UUID
}

func NewCombinedBet(participantID, leagueID uuid.UUID, stake float64, useBonus bool, bonusID *uuid.UUID, legs []CombinedBetLeg) (*CombinedBet, error) {
	if stake <= 0 {
		return nil, apperror.ErrInvalidBetAmount
	}
	if len(legs) < 2 {
		return nil, apperror.ErrInvalidBetAmount
	}

	totalOdds := 1.0
	for _, leg := range legs {
		totalOdds *= leg.OddsAtPlacement
	}

	now := time.Now().UTC()
	return &CombinedBet{
		ID:                 uuid.New(),
		ParticipantID:      participantID,
		LeagueID:           leagueID,
		Stake:              stake,
		UseBonus:           useBonus,
		TotalOdds:          totalOdds,
		PotentialWin:       stake * totalOdds,
		Status:             valueobject.CombinedStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
		Legs:               legs,
		ParticipantBonusID: bonusID,
	}, nil
}

func (cb *CombinedBet) Resolve(legResults map[uuid.UUID]valueobject.LegStatus) {
	hasLoss := false
	allSettled := true
	hasVoid := false

	now := time.Now().UTC()

	for i := range cb.Legs {
		leg := &cb.Legs[i]
		if result, exists := legResults[leg.ID]; exists {
			leg.Status = result
			leg.SettledAt = &now

			if result == valueobject.LegStatusLost {
				hasLoss = true
			} else if result == valueobject.LegStatusVoided {
				hasVoid = true
			}
		}

		if leg.Status == valueobject.LegStatusPending {
			allSettled = false
		}
	}

	if hasLoss {
		cb.Status = valueobject.CombinedStatusLost
		cb.SettledAt = &now
		return
	}

	if allSettled {
		if hasVoid {
			cb.RecalculateOddsAfterVoid()
		}
		cb.Status = valueobject.CombinedStatusWon
		cb.SettledAt = &now
	}
}

func (cb *CombinedBet) RecalculateOddsAfterVoid() {
	newTotalOdds := 1.0

	for _, leg := range cb.Legs {
		if leg.Status != valueobject.LegStatusVoided {
			newTotalOdds *= leg.OddsAtPlacement
		}
	}

	cb.TotalOdds = newTotalOdds
	cb.PotentialWin = cb.Stake * newTotalOdds
}

func (cb *CombinedBet) CalculateCashoutValue(currentMarketOdds map[uuid.UUID]float64) float64 {
	houseMargin := 0.90

	multiplier := 1.0

	for _, leg := range cb.Legs {
		switch leg.Status {
		case valueobject.LegStatusWon:
			multiplier *= leg.OddsAtPlacement

		case valueobject.LegStatusPending:
			currentOdds, ok := currentMarketOdds[leg.MarketID]
			if !ok {
				return 0.0
			}
			adjustedOdds := currentOdds * houseMargin
			multiplier *= adjustedOdds

		case valueobject.LegStatusLost:
			return 0.0

		case valueobject.LegStatusVoided:
			continue
		}
	}

	cashoutValue := cb.Stake * multiplier

	return cashoutValue
}
