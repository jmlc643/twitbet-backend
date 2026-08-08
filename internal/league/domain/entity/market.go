package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type MarketOption struct {
	ID          uuid.UUID
	MarketID    uuid.UUID
	Name        string
	InitialOdds float64
	CurrentOdds float64
}

type Market struct {
	ID        uuid.UUID
	LeagueID  uuid.UUID
	MatchID   *uuid.UUID
	Name      string
	Status             string
	CancellationReason *string
	Options            []MarketOption
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MarketOptionCreate struct {
	Name string
	Odds float64
}

func NewMarket(leagueID uuid.UUID, matchID *uuid.UUID, name string, options []MarketOptionCreate) (*Market, error) {
	if name == "" {
		return nil, apperror.ErrInvalidMarketName
	}
	if len(options) < 2 {
		return nil, apperror.ErrInvalidMarketOptions
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
		})
	}

	return &Market{
		ID:        marketID,
		LeagueID:  leagueID,
		MatchID:   matchID,
		Name:      name,
		Status:    "OPEN",
		Options:   marketOptions,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
