package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type MarketType string

const (
	MarketTypeResult       MarketType = "RESULT"
	MarketTypeTotals       MarketType = "TOTALS"
	MarketTypeHandicap     MarketType = "HANDICAP"
	MarketTypeCorrectScore MarketType = "CORRECT_SCORE"
	MarketTypeOther        MarketType = "OTHER"
)

type MarketStatus string

const (
	MarketStatusOpen      MarketStatus = "OPEN"
	MarketStatusActive    MarketStatus = "ACTIVE"
	MarketStatusSuspended MarketStatus = "SUSPENDED"
	MarketStatusResolved  MarketStatus = "RESOLVED"
	MarketStatusCancelled MarketStatus = "CANCELLED"
	MarketStatusVoided    MarketStatus = "VOIDED"
)

type MarketOptionStatus string

const (
	MarketOptionStatusActive  MarketOptionStatus = "ACTIVE"
	MarketOptionStatusBlocked MarketOptionStatus = "BLOCKED"
)

func IsValidMarketType(marketType string) bool {
	switch MarketType(marketType) {
	case MarketTypeResult, MarketTypeTotals, MarketTypeHandicap, MarketTypeCorrectScore, MarketTypeOther:
		return true
	default:
		return false
	}
}

func NotResolvableMarketStatus(market *Market) bool {
	return market.Status == string(MarketStatusResolved) ||
		market.Status == string(MarketStatusCancelled) ||
		market.Status == string(MarketStatusVoided)
}

type MarketOption struct {
	ID          uuid.UUID
	MarketID    uuid.UUID
	Name        string
	InitialOdds float64
	CurrentOdds float64
	Status      string
}

func (o *MarketOption) IsBlocked() bool {
	return o.Status == string(MarketOptionStatusBlocked)
}

type Market struct {
	ID        uuid.UUID
	LeagueID  uuid.UUID
	MatchID   *uuid.UUID
	Name      string
	Type      string
	Status    string
	CancellationReason *string
	Options   []MarketOption
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MarketOptionCreate struct {
	Name string
	Odds float64
}

func NewMarket(leagueID uuid.UUID, matchID *uuid.UUID, name string, marketType string, options []MarketOptionCreate) (*Market, error) {
	if name == "" {
		return nil, apperror.ErrInvalidMarketName
	}
	if len(options) < 2 {
		return nil, apperror.ErrInvalidMarketOptions
	}
	if !IsValidMarketType(marketType) {
		marketType = string(MarketTypeOther)
	}

	now := time.Now().UTC()
	marketID := uuid.New()

	marketOptions := make([]MarketOption, 0, len(options))
	for _, opt := range options {
		marketOptions = append(marketOptions, MarketOption{
			ID:          uuid.New(),
			MarketID:    marketID,
			Name:        opt.Name,
			InitialOdds: opt.Odds,
			CurrentOdds: opt.Odds,
			Status:      string(MarketOptionStatusActive),
		})
	}

	return &Market{
		ID:        marketID,
		LeagueID:  leagueID,
		MatchID:   matchID,
		Name:      name,
		Type:      marketType,
		Status:    string(MarketStatusOpen),
		Options:   marketOptions,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NewMarketOption(marketID uuid.UUID, name string, odds float64) *MarketOption {
	return &MarketOption{
		ID:          uuid.New(),
		MarketID:    marketID,
		Name:        name,
		InitialOdds: odds,
		CurrentOdds: odds,
		Status:      string(MarketOptionStatusActive),
	}
}