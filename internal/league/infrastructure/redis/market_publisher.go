package redis

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/port"
	"github.com/redis/go-redis/v9"
)

type marketPublisher struct {
	client *redis.Client
}

func NewMarketPublisher(client *redis.Client) port.MarketEventPublisher {
	return &marketPublisher{client: client}
}

type MarketStatusEvent struct {
	Type     string    `json:"type"`
	MarketID uuid.UUID `json:"market_id"`
	Status   string    `json:"status"`
}

type MarketOptionEventDTO struct {
	ID          uuid.UUID `json:"id"`
	MarketID    uuid.UUID `json:"market_id"`
	Name        string    `json:"name"`
	InitialOdds float64   `json:"initial_odds"`
	CurrentOdds float64   `json:"current_odds"`
}

type MarketOddsEvent struct {
	Type     string                 `json:"type"`
	MarketID uuid.UUID              `json:"market_id"`
	Options  []MarketOptionEventDTO `json:"options"`
}

func (p *marketPublisher) PublishMarketStatusChanged(ctx context.Context, marketID uuid.UUID, newStatus string) error {
	event := MarketStatusEvent{
		Type:     "MARKET_STATUS_CHANGED",
		MarketID: marketID,
		Status:   newStatus,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "market_events", payload).Err()
}

func (p *marketPublisher) PublishOddsUpdated(ctx context.Context, marketID uuid.UUID, options []entity.MarketOption) error {
	dtoOptions := make([]MarketOptionEventDTO, 0, len(options))
	for _, opt := range options {
		dtoOptions = append(dtoOptions, MarketOptionEventDTO{
			ID:          opt.ID,
			MarketID:    opt.MarketID,
			Name:        opt.Name,
			InitialOdds: opt.InitialOdds,
			CurrentOdds: opt.CurrentOdds,
		})
	}

	event := MarketOddsEvent{
		Type:     "ODDS_UPDATED",
		MarketID: marketID,
		Options:  dtoOptions,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "market_events", payload).Err()
}

type MatchStatusEvent struct {
	Type    string    `json:"type"`
	MatchID uuid.UUID `json:"match_id"`
	Status  string    `json:"status"`
}

func (p *marketPublisher) PublishMatchStatusChanged(ctx context.Context, matchID uuid.UUID, newStatus string) error {
	event := MatchStatusEvent{
		Type:    "MATCH_STATUS_CHANGED",
		MatchID: matchID,
		Status:  newStatus,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "match_events", payload).Err()
}
