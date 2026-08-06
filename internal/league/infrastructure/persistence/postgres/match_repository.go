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
		Slug:      match.Slug,
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
		Slug:      dbMatch.Slug,
		Title:     dbMatch.Title,
		StartTime: dbMatch.StartTime,
		Status:    dbMatch.Status,
		CreatedAt: dbMatch.CreatedAt,
		UpdatedAt: dbMatch.UpdatedAt,
	}, nil
}

func (r *matchRepository) GetMatchBySlug(ctx context.Context, slug string) (*entity.Match, error) {
	var dbMatch model.MatchModel
	if err := r.db.WithContext(ctx).First(&dbMatch, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	matchID, _ := uuid.Parse(dbMatch.ID)
	leagueID, _ := uuid.Parse(dbMatch.LeagueID)

	return &entity.Match{
		ID:        matchID,
		LeagueID:  leagueID,
		Slug:      dbMatch.Slug,
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
			Slug:      m.Slug,
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

func (r *matchRepository) GetMarketByID(ctx context.Context, id uuid.UUID) (*entity.Market, error) {
	var dbMarket model.MarketModel
	if err := r.db.WithContext(ctx).Preload("Options").First(&dbMarket, "id = ?", id.String()).Error; err != nil {
		return nil, err
	}
	markets := mapDBMarketsToEntity([]model.MarketModel{dbMarket})
	if len(markets) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &markets[0], nil
}

func (r *matchRepository) UpdateMarket(ctx context.Context, market *entity.Market) error {
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

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(dbMarket).Error; err != nil {
			return err
		}
		for _, opt := range market.Options {
			dbOpt := &model.MarketOptionModel{
				ID:          opt.ID.String(),
				MarketID:    market.ID.String(),
				Name:        opt.Name,
				InitialOdds: opt.InitialOdds,
				CurrentOdds: opt.CurrentOdds,
			}
			if err := tx.Save(dbOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *matchRepository) UpdateMatchStatusAtomic(ctx context.Context, matchID uuid.UUID, newStatus string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Model(&model.MatchModel{}).Where("id = ?", matchID.String()).Update("status", newStatus).Error; err != nil {
		return err
	}

	if newStatus == "VOIDED" {
		if err := tx.Model(&model.MarketModel{}).Where("match_id = ?", matchID.String()).Update("status", "VOIDED").Error; err != nil {
			return err
		}

		var markets []model.MarketModel
		if err := tx.Where("match_id = ?", matchID.String()).Find(&markets).Error; err != nil {
			return err
		}
		
		var marketIDs []string
		for _, m := range markets {
			marketIDs = append(marketIDs, m.ID)
		}

		if len(marketIDs) > 0 {
			var options []model.MarketOptionModel
			if err := tx.Where("market_id IN ?", marketIDs).Find(&options).Error; err != nil {
				return err
			}

			var optionIDs []string
			for _, opt := range options {
				optionIDs = append(optionIDs, opt.ID)
			}

			if len(optionIDs) > 0 {
				var bets []model.BetModel
				if err := tx.Where("market_option_id IN ? AND status = ?", optionIDs, string(entity.BetStatusAccepted)).Find(&bets).Error; err != nil {
					return err
				}

				for _, bet := range bets {
					if err := tx.Model(&bet).Update("status", string(entity.BetStatusVoided)).Error; err != nil {
						return err
					}

					var participant model.ParticipantModel
					if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", bet.ParticipantID).First(&participant).Error; err != nil {
						return err
					}

					participant.Balance += bet.Amount
					if err := tx.Save(&participant).Error; err != nil {
						return err
					}

					var match model.MatchModel
					if err := tx.Where("id = ?", matchID.String()).First(&match).Error; err != nil {
						return err
					}

					txModel := &model.TransactionModel{
						ID:        uuid.New().String(),
						LeagueID:  match.LeagueID,
						UserID:    participant.UserID,
						Amount:    bet.Amount,
						Type:      string(entity.TransactionTypeRefund),
					}
					if err := tx.Create(txModel).Error; err != nil {
						return err
					}
				}
			}
		}
	}

	return tx.Commit().Error
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