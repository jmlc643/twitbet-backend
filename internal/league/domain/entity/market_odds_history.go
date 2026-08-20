package entity

import (
	"time"

	"github.com/google/uuid"
)

type OddsChangeReason string

const (
	OddsChangeReasonManual       OddsChangeReason = "MANUAL"
	OddsChangeReasonRebalance    OddsChangeReason = "REBALANCE"
	OddsChangeReasonAutoCooldown OddsChangeReason = "AUTO_COOLDOWN"
	OddsChangeReasonBetPlacement OddsChangeReason = "BET_PLACEMENT"
)

type MarketOddsHistory struct {
	ID             uuid.UUID
	MarketOptionID uuid.UUID
	OldOdds        *float64
	NewOdds        float64
	ChangedBy      *uuid.UUID
	Reason         string
	CreatedAt      time.Time
}

func NewMarketOddsHistory(optionID uuid.UUID, oldOdds *float64, newOdds float64, changedBy *uuid.UUID, reason string) *MarketOddsHistory {
	return &MarketOddsHistory{
		ID:             uuid.New(),
		MarketOptionID: optionID,
		OldOdds:        oldOdds,
		NewOdds:        newOdds,
		ChangedBy:      changedBy,
		Reason:         reason,
		CreatedAt:      time.Now().UTC(),
	}
}
