package postgres

import (
	"context"
	"errors"

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
