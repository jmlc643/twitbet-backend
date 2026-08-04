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
