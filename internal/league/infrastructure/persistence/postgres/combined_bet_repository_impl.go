package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/valueobject"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/mapper"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type combinedBetRepository struct {
	db *gorm.DB
}

func NewCombinedBetRepository(db *gorm.DB) repository.CombinedBetRepository {
	return &combinedBetRepository{db: db}
}

func (r *combinedBetRepository) Create(ctx context.Context, bet *entity.CombinedBet) error {
	dbModel := mapper.ToCombinedBetModel(bet)
	return r.db.WithContext(ctx).Create(dbModel).Error
}

func (r *combinedBetRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CombinedBet, error) {
	var dbModel model.CombinedBetModel
	if err := r.db.WithContext(ctx).Preload("Legs").First(&dbModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToCombinedBetEntity(&dbModel), nil
}

func (r *combinedBetRepository) GetByParticipantID(ctx context.Context, participantID uuid.UUID, status *string, startDate, endDate *time.Time, limit, offset int) ([]entity.CombinedBet, int64, error) {
	var dbModels []model.CombinedBetModel
	query := r.db.WithContext(ctx).Model(&model.CombinedBetModel{}).Where("participant_id = ?", participantID)

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Order("created_at DESC").Preload("Legs").Find(&dbModels).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]entity.CombinedBet, len(dbModels))
	for i, dbModel := range dbModels {
		entities[i] = *mapper.ToCombinedBetEntity(&dbModel)
	}
	return entities, total, nil
}

func (r *combinedBetRepository) GetByLeagueID(ctx context.Context, leagueID uuid.UUID) ([]entity.CombinedBet, error) {
	var dbModels []model.CombinedBetModel
	if err := r.db.WithContext(ctx).Where("league_id = ?", leagueID).Preload("Legs").Find(&dbModels).Error; err != nil {
		return nil, err
	}

	entities := make([]entity.CombinedBet, len(dbModels))
	for i, dbModel := range dbModels {
		entities[i] = *mapper.ToCombinedBetEntity(&dbModel)
	}
	return entities, nil
}

func (r *combinedBetRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.CombinedBetModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r *combinedBetRepository) UpdateCashout(ctx context.Context, id uuid.UUID, cashoutValue float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.CombinedBetModel{}).Where("id = ?", id).Updates(map[string]interface{}{
			"cashout_value": cashoutValue,
			"status":        string(valueobject.CombinedStatusCashout),
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.CombinedBetLegModel{}).Where("combined_bet_id = ? AND status = ?", id, string(valueobject.LegStatusPending)).Update("status", string(valueobject.LegStatusCashout)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *combinedBetRepository) PlaceCombinedBetAtomic(ctx context.Context, bet *entity.CombinedBet, txEntity *entity.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var participant model.ParticipantModel
		if err := tx.First(&participant, bet.ParticipantID).Error; err != nil {
			return err
		}

		if bet.UseBonus {
			if bet.ParticipantBonusID == nil {
				return errors.New("Falta el id del bono para la apuesta de bono")
			}
			var bonus model.ParticipantBonusModel
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", bet.ParticipantBonusID.String()).First(&bonus).Error; err != nil {
				return err
			}
			if bonus.Status != string(entity.BonusStatusPending) {
				return errors.New("El bono no está pendiente")
			}
			if bonus.Amount != bet.Stake {
				return errors.New("El monto de la apuesta no coincide con el monto del bono")
			}
			bonus.Status = string(entity.BonusStatusUsed)
			if err := tx.Save(&bonus).Error; err != nil {
				return err
			}
		} else {
			if participant.Balance < bet.Stake {
				return apperror.ErrInsufficientBalance
			}
			participant.Balance -= bet.Stake
			if err := tx.Save(&participant).Error; err != nil {
				return err
			}
		}

		for _, leg := range bet.Legs {
			var marketOption model.MarketOptionModel
			if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", leg.SelectionID.String()).First(&marketOption).Error; err != nil {
				return err
			}
			if marketOption.CurrentOdds != leg.OddsAtPlacement {
				return &apperror.OddsChangedError{CurrentOdds: marketOption.CurrentOdds} // We could return an array of current odds, but one is enough to trigger the 409
			}
		}

		betModel := mapper.ToCombinedBetModel(bet)
		if err := tx.Create(betModel).Error; err != nil {
			return err
		}

		txModel := &model.TransactionModel{
			ID:        txEntity.ID.String(),
			LeagueID:  txEntity.LeagueID.String(),
			UserID:    participant.UserID,
			Amount:    -txEntity.Amount,
			Type:      string(txEntity.Type),
			CreatedAt: txEntity.CreatedAt,
		}
		if err := tx.Create(txModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *combinedBetRepository) ResolveCombinedBetsForMarketAtomic(ctx context.Context, marketID uuid.UUID, winningOptionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var affectedLegs []model.CombinedBetLegModel
		if err := tx.Where("market_id = ? AND status = ?", marketID, string(valueobject.LegStatusPending)).Find(&affectedLegs).Error; err != nil {
			return err
		}

		if len(affectedLegs) == 0 {
			return nil
		}

		winningSet := make(map[uuid.UUID]bool)
		for _, w := range winningOptionIDs {
			winningSet[w] = true
		}

		affectedBetIDs := make(map[uuid.UUID]bool)

		for _, leg := range affectedLegs {
			newStatus := string(valueobject.LegStatusLost)
			if winningSet[leg.SelectionID] {
				newStatus = string(valueobject.LegStatusWon)
			}
			
			if err := tx.Model(&model.CombinedBetLegModel{}).Where("id = ?", leg.ID).Update("status", newStatus).Error; err != nil {
				return err
			}
			affectedBetIDs[leg.CombinedBetID] = true
		}

		for betID := range affectedBetIDs {
			var betModel model.CombinedBetModel
			if err := tx.Preload("Legs").First(&betModel, betID).Error; err != nil {
				return err
			}

			cbEntity := mapper.ToCombinedBetEntity(&betModel)
			
			legResults := make(map[uuid.UUID]valueobject.LegStatus)
			for _, l := range cbEntity.Legs {
				legResults[l.ID] = l.Status
			}

			cbEntity.Resolve(legResults)
			
			if err := tx.Model(&model.CombinedBetModel{}).Where("id = ?", betID).Updates(map[string]interface{}{
				"status":        string(cbEntity.Status),
				"total_odds":    cbEntity.TotalOdds,
				"potential_win": cbEntity.PotentialWin,
				"settled_at":    cbEntity.SettledAt,
			}).Error; err != nil {
				return err
			}

			if cbEntity.Status == valueobject.CombinedStatusWon {
				var participant model.ParticipantModel
				if err := tx.First(&participant, cbEntity.ParticipantID).Error; err != nil {
					return err
				}

				participant.Balance += cbEntity.PotentialWin

				if err := tx.Save(&participant).Error; err != nil {
					return err
				}
				
				txWinModel := &model.TransactionModel{
					ID:        uuid.New().String(),
					LeagueID:  cbEntity.LeagueID.String(),
					UserID:    participant.UserID,
					Amount:    cbEntity.PotentialWin,
					Type:      "COMBINED_WIN",
					CreatedAt: cbEntity.SettledAt.UTC(),
				}
				if err := tx.Create(txWinModel).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

