package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/mapper"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type BetRepository struct {
	db *gorm.DB
}

func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{db: db}
}

func (r *BetRepository) PlaceBetAtomic(ctx context.Context, bet *entity.Bet, transaction *entity.Transaction) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var participant model.ParticipantModel
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", bet.ParticipantID.String()).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("participant not found")
		}
		return err
	}

	if participant.Balance < bet.Amount {
		return apperror.ErrInsufficientBalance
	}

	participant.Balance -= bet.Amount
	if err := tx.Save(&participant).Error; err != nil {
		return err
	}

	transactionModel := mapper.EntityToTransactionModel(transaction)
	if err := tx.Create(transactionModel).Error; err != nil {
		return err
	}

	betModel := mapper.EntityToBetModel(bet)
	if err := tx.Create(betModel).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *BetRepository) CashoutAtomic(ctx context.Context, bet *entity.Bet, transaction *entity.Transaction) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var participant model.ParticipantModel
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", bet.ParticipantID.String()).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("participant not found")
		}
		return err
	}

	participant.Balance += transaction.Amount
	if err := tx.Save(&participant).Error; err != nil {
		return err
	}

	transactionModel := mapper.EntityToTransactionModel(transaction)
	if err := tx.Create(transactionModel).Error; err != nil {
		return err
	}

	betModel := mapper.EntityToBetModel(bet)
	if err := tx.Save(betModel).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *BetRepository) GetBetByID(ctx context.Context, id uuid.UUID) (*entity.Bet, error) {
	var betModel model.BetModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&betModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	participantID, _ := uuid.Parse(betModel.ParticipantID)
	marketOptionID, _ := uuid.Parse(betModel.MarketOptionID)

	return &entity.Bet{
		ID:             id,
		ParticipantID:  participantID,
		MarketOptionID: marketOptionID,
		Amount:         betModel.Amount,
		Odds:           betModel.Odds,
		PotentialWin:   betModel.PotentialWin,
		Status:         entity.BetStatus(betModel.Status),
		PlacedAt:       betModel.PlacedAt,
		UpdatedAt:      betModel.UpdatedAt,
	}, nil
}

func (r *BetRepository) UpdateBetStatus(ctx context.Context, id uuid.UUID, status entity.BetStatus) error {
	return r.db.WithContext(ctx).Model(&model.BetModel{}).
		Where("id = ?", id.String()).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *BetRepository) ResolveMarketAtomic(ctx context.Context, marketID uuid.UUID, winningOptionID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Model(&model.MarketModel{}).Where("id = ?", marketID.String()).Update("status", "RESOLVED").Error; err != nil {
		return err
	}

	var options []model.MarketOptionModel
	if err := tx.Where("market_id = ?", marketID.String()).Find(&options).Error; err != nil {
		return err
	}

	var optionIDs []string
	for _, opt := range options {
		optionIDs = append(optionIDs, opt.ID)
	}

	if len(optionIDs) == 0 {
		return tx.Commit().Error
	}

	var bets []model.BetModel
	if err := tx.Where("market_option_id IN ? AND status = ?", optionIDs, string(entity.BetStatusAccepted)).Find(&bets).Error; err != nil {
		return err
	}

	var market model.MarketModel
	if err := tx.Where("id = ?", marketID.String()).First(&market).Error; err != nil {
		return err
	}

	for _, bet := range bets {
		if bet.MarketOptionID == winningOptionID.String() {
			if err := tx.Model(&bet).Update("status", string(entity.BetStatusWon)).Error; err != nil {
				return err
			}

			var participant model.ParticipantModel
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", bet.ParticipantID).First(&participant).Error; err != nil {
				return err
			}

			participant.Balance += bet.PotentialWin
			if err := tx.Save(&participant).Error; err != nil {
				return err
			}

			txModel := &model.TransactionModel{
				ID:        uuid.New().String(),
				LeagueID:  market.LeagueID,
				UserID:    participant.UserID,
				Amount:    bet.PotentialWin,
				Type:      string(entity.TransactionTypeWin),
			}
			if err := tx.Create(txModel).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&bet).Update("status", string(entity.BetStatusLost)).Error; err != nil {
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *BetRepository) GetBetsByParticipantID(ctx context.Context, participantID uuid.UUID, status *entity.BetStatus, startDate, endDate *time.Time, limit, offset int) ([]entity.BetDetail, int64, error) {
	var results []struct {
		ID           string
		Amount       float64
		Odds         float64
		PotentialWin float64
		Status       string
		PlacedAt     time.Time
		MatchTitle   string
		MarketID     string
		MarketName   string
		OptionID     string
		OptionName   string
	}

	baseQuery := r.db.WithContext(ctx).Table("bets").
		Joins("JOIN market_options ON bets.market_option_id = market_options.id").
		Joins("JOIN markets ON market_options.market_id = markets.id").
		Joins("JOIN matches ON markets.match_id = matches.id").
		Where("bets.participant_id = ?", participantID.String())

	if status != nil {
		baseQuery = baseQuery.Where("bets.status = ?", string(*status))
	}
	if startDate != nil {
		baseQuery = baseQuery.Where("bets.placed_at >= ?", *startDate)
	}
	if endDate != nil {
		baseQuery = baseQuery.Where("bets.placed_at <= ?", *endDate)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := baseQuery.
		Select(`bets.id, bets.amount, bets.odds, bets.potential_win, bets.status, bets.placed_at,
		        matches.title as match_title, markets.id as market_id, markets.name as market_name, market_options.id as option_id, market_options.name as option_name`).
		Order("bets.placed_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	var betDetails []entity.BetDetail
	for _, row := range results {
		id, _ := uuid.Parse(row.ID)
		marketID, _ := uuid.Parse(row.MarketID)
		optionID, _ := uuid.Parse(row.OptionID)
		
		betDetails = append(betDetails, entity.BetDetail{
			ID:           id,
			Amount:       row.Amount,
			Odds:         row.Odds,
			PotentialWin: row.PotentialWin,
			Status:       entity.BetStatus(row.Status),
			PlacedAt:     row.PlacedAt,
			MatchTitle:   row.MatchTitle,
			MarketID:     marketID,
			MarketName:   row.MarketName,
			OptionID:     optionID,
			OptionName:   row.OptionName,
		})
	}

	return betDetails, total, nil
}
