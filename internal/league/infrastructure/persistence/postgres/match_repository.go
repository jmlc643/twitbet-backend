package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) repository.MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) CreateMatch(ctx context.Context, match *entity.Match) error {
	dbMatch := &model.MatchModel{
		ID:        match.ID.String(),
		LeagueID:  match.LeagueID.String(),
		Title:     match.Title,
		StartTime: match.StartTime,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}

	return r.db.WithContext(ctx).Create(dbMatch).Error
}

func (r *matchRepository) CreateMarket(ctx context.Context, market *entity.Market) error {
	var matchID *string
	if market.MatchID != nil {
		idStr := market.MatchID.String()
		matchID = &idStr
	}

	dbMarket := &model.MarketModel{
		ID:        market.ID.String(),
		LeagueID:  market.LeagueID.String(),
		MatchID:   matchID,
		Name:      market.Name,
		Status:    market.Status,
		CreatedAt: market.CreatedAt,
		UpdatedAt: market.UpdatedAt,
	}

	for _, opt := range market.Options {
		dbMarket.Options = append(dbMarket.Options, model.MarketOptionModel{
			ID:          opt.ID.String(),
			MarketID:    market.ID.String(),
			Name:        opt.Name,
			InitialOdds: opt.InitialOdds,
			CurrentOdds: opt.CurrentOdds,
		})
	}

	return r.db.WithContext(ctx).Create(dbMarket).Error
}

func (r *matchRepository) GetMatchByID(ctx context.Context, id uuid.UUID) (*entity.Match, error) {
	var dbMatch model.MatchModel
	if err := r.db.WithContext(ctx).First(&dbMatch, "id = ?", id.String()).Error; err != nil {
		return nil, err
	}

	matchID, _ := uuid.Parse(dbMatch.ID)
	leagueID, _ := uuid.Parse(dbMatch.LeagueID)

	return &entity.Match{
		ID:        matchID,
		LeagueID:  leagueID,
		Title:     dbMatch.Title,
		StartTime: dbMatch.StartTime,
		Status:    dbMatch.Status,
		CreatedAt: dbMatch.CreatedAt,
		UpdatedAt: dbMatch.UpdatedAt,
	}, nil
}

func (r *matchRepository) GetMatchesByLeagueID(ctx context.Context, leagueID uuid.UUID, limit, offset int, status string) ([]entity.Match, int64, error) {
	var dbMatches []model.MatchModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.MatchModel{}).Where("league_id = ?", leagueID.String())
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	query = query.Order("start_time ASC")

	if err := query.Find(&dbMatches).Error; err != nil {
		return nil, 0, err
	}

	matches := make([]entity.Match, 0, len(dbMatches))
	for _, m := range dbMatches {
		mID, _ := uuid.Parse(m.ID)
		lID, _ := uuid.Parse(m.LeagueID)
		matches = append(matches, entity.Match{
			ID:        mID,
			LeagueID:  lID,
			Title:     m.Title,
			StartTime: m.StartTime,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}

	return matches, total, nil
}

func (r *matchRepository) GetMarketsByLeagueID(ctx context.Context, leagueID uuid.UUID) ([]entity.Market, error) {
	var dbMarkets []model.MarketModel
	if err := r.db.WithContext(ctx).Preload("Options").Where("league_id = ? AND match_id IS NULL", leagueID.String()).Find(&dbMarkets).Error; err != nil {
		return nil, err
	}
	return mapDBMarketsToEntity(dbMarkets), nil
}

func (r *matchRepository) GetMarketsByMatchID(ctx context.Context, matchID uuid.UUID) ([]entity.Market, error) {
	var dbMarkets []model.MarketModel
	if err := r.db.WithContext(ctx).Preload("Options").Where("match_id = ?", matchID.String()).Find(&dbMarkets).Error; err != nil {
		return nil, err
	}
	return mapDBMarketsToEntity(dbMarkets), nil
}

func mapDBMarketsToEntity(dbMarkets []model.MarketModel) []entity.Market {
	markets := make([]entity.Market, 0, len(dbMarkets))
	for _, m := range dbMarkets {
		mID, _ := uuid.Parse(m.ID)
		lID, _ := uuid.Parse(m.LeagueID)
		var matchID *uuid.UUID
		if m.MatchID != nil {
			id, _ := uuid.Parse(*m.MatchID)
			matchID = &id
		}

		options := make([]entity.MarketOption, 0, len(m.Options))
		for _, o := range m.Options {
			oID, _ := uuid.Parse(o.ID)
			oMarketID, _ := uuid.Parse(o.MarketID)
			options = append(options, entity.MarketOption{
				ID:          oID,
				MarketID:    oMarketID,
				Name:        o.Name,
				InitialOdds: o.InitialOdds,
				CurrentOdds: o.CurrentOdds,
			})
		}

		markets = append(markets, entity.Market{
			ID:        mID,
			LeagueID:  lID,
			MatchID:   matchID,
			Name:      m.Name,
			Status:    m.Status,
			Options:   options,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return markets
}
