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

type MarketOptionEventDTO struct {
	ID          uuid.UUID `json:"id"`
	MarketID    uuid.UUID `json:"market_id"`
	Name        string    `json:"name"`
	InitialOdds float64   `json:"initial_odds"`
	CurrentOdds float64   `json:"current_odds"`
	Status      string    `json:"status"`
}

func newMarketOptionDTO(opt entity.MarketOption) MarketOptionEventDTO {
	status := opt.Status
	if status == "" {
		status = string(entity.MarketOptionStatusActive)
	}
	return MarketOptionEventDTO{
		ID:          opt.ID,
		MarketID:    opt.MarketID,
		Name:        opt.Name,
		InitialOdds: opt.InitialOdds,
		CurrentOdds: opt.CurrentOdds,
		Status:      status,
	}
}

type MarketSnapshotEvent struct {
	Type       string                 `json:"type"`
	MarketID   uuid.UUID              `json:"market_id"`
	LeagueID   uuid.UUID              `json:"league_id"`
	MatchID    *uuid.UUID             `json:"match_id"`
	Name       string                 `json:"name"`
	MarketType string                 `json:"market_type"`
	Status     string                 `json:"status"`
	Options    []MarketOptionEventDTO `json:"options"`
}

func newMarketSnapshot(eventType string, market entity.Market) MarketSnapshotEvent {
	options := make([]MarketOptionEventDTO, 0, len(market.Options))
	for _, opt := range market.Options {
		options = append(options, newMarketOptionDTO(opt))
	}

	return MarketSnapshotEvent{
		Type:       eventType,
		MarketID:   market.ID,
		LeagueID:   market.LeagueID,
		MatchID:    market.MatchID,
		Name:       market.Name,
		MarketType: market.Type,
		Status:     market.Status,
		Options:    options,
	}
}

func (p *marketPublisher) publishSnapshot(ctx context.Context, payloadType string, market entity.Market) error {
	event := newMarketSnapshot(payloadType, market)
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "market_events", payload).Err()
}

func (p *marketPublisher) PublishMarketCreated(ctx context.Context, market entity.Market) error {
	return p.publishSnapshot(ctx, "MARKET_CREATED", market)
}

func (p *marketPublisher) PublishMarketOptionsUpdated(ctx context.Context, market entity.Market) error {
	return p.publishSnapshot(ctx, "MARKET_OPTIONS_UPDATED", market)
}

type MarketStatusEvent struct {
	Type     string    `json:"type"`
	MarketID uuid.UUID `json:"market_id"`
	Status   string    `json:"status"`
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
		dtoOptions = append(dtoOptions, newMarketOptionDTO(opt))
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

type MarketResolvedEvent struct {
	Type             string      `json:"type"`
	MarketID         uuid.UUID   `json:"market_id"`
	LeagueID         uuid.UUID   `json:"league_id"`
	WinningOptionIDs []uuid.UUID `json:"winning_option_ids"`
}

func (p *marketPublisher) PublishMarketResolved(ctx context.Context, marketID uuid.UUID, leagueID uuid.UUID, winningOptionIDs []uuid.UUID) error {
	event := MarketResolvedEvent{
		Type:             "MARKET_RESOLVED",
		MarketID:         marketID,
		LeagueID:         leagueID,
		WinningOptionIDs: winningOptionIDs,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "market_events", payload).Err()
}

type ParticipantBalanceEvent struct {
	Type          string    `json:"type"`
	ParticipantID uuid.UUID `json:"participant_id"`
	LeagueID      uuid.UUID `json:"league_id"`
	UserID        uuid.UUID `json:"user_id"`
}

func (p *marketPublisher) PublishParticipantBalanceUpdated(ctx context.Context, participantID uuid.UUID, leagueID uuid.UUID, userID uuid.UUID) error {
	event := ParticipantBalanceEvent{
		Type:          "PARTICIPANT_BALANCE_UPDATED",
		ParticipantID: participantID,
		LeagueID:      leagueID,
		UserID:        userID,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, "market_events", payload).Err()
}