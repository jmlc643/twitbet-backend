package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/mapper"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type LeagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) *LeagueRepository {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) CreateLeagueWithAdminParticipant(ctx context.Context, l *entity.League, p *entity.Participant, t *entity.Transaction) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	leagueModel := mapper.EntityToLeagueModel(l)
	if err := tx.Create(leagueModel).Error; err != nil {
		return err
	}

	participantModel := mapper.EntityToParticipantModel(p)
	if err := tx.Create(participantModel).Error; err != nil {
		return err
	}

	if t != nil {
		transactionModel := mapper.EntityToTransactionModel(t)
		if err := tx.Create(transactionModel).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *LeagueRepository) GetLeagueByInviteCode(ctx context.Context, code string) (*entity.League, error) {
	var modelLeague model.LeagueModel
	err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&modelLeague).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	adminID, _ := uuid.Parse(modelLeague.AdminID)
	id, _ := uuid.Parse(modelLeague.ID)

	return &entity.League{
		ID:             id,
		AdminID:        adminID,
		Name:           modelLeague.Name,
		InitialBalance: modelLeague.InitialBalance,
		MaxRecharges:   modelLeague.MaxRecharges,
		HideStandings:  modelLeague.HideStandings,
		InviteCode:     modelLeague.InviteCode,
		CreatedAt:      modelLeague.CreatedAt,
		UpdatedAt:      modelLeague.UpdatedAt,
	}, nil
}

func (r *LeagueRepository) GetUserLeagues(ctx context.Context, userID uuid.UUID) ([]entity.LeagueSummary, error) {
	var results []struct {
		LeagueID         string
		Name             string
		AdminID          string
		Balance          float64
		ParticipantCount int
	}

	query := `
		SELECT 
			l.id as league_id, 
			l.name, 
			l.admin_id, 
			p.balance, 
			(SELECT COUNT(*) FROM league_participants p2 WHERE p2.league_id = l.id) as participant_count
		FROM leagues l
		JOIN league_participants p ON l.id = p.league_id
		WHERE p.user_id = ?
	`

	err := r.db.WithContext(ctx).Raw(query, userID.String()).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var summaries []entity.LeagueSummary
	for _, res := range results {
		role := entity.RoleMember
		if res.AdminID == userID.String() {
			role = entity.RoleAdmin
		}

		leagueUUID, _ := uuid.Parse(res.LeagueID)

		summaries = append(summaries, entity.LeagueSummary{
			LeagueID:         leagueUUID,
			Name:             res.Name,
			Role:             role,
			ParticipantCount: res.ParticipantCount,
			Balance:          res.Balance,
		})
	}

	return summaries, nil
}

func (r *LeagueRepository) IsParticipant(ctx context.Context, leagueID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ParticipantModel{}).Where("league_id = ? AND user_id = ?", leagueID.String(), userID.String()).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LeagueRepository) JoinLeague(ctx context.Context, p *entity.Participant, t *entity.Transaction) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	participantModel := mapper.EntityToParticipantModel(p)
	if err := tx.Create(participantModel).Error; err != nil {
		return err
	}

	if t != nil {
		transactionModel := mapper.EntityToTransactionModel(t)
		if err := tx.Create(transactionModel).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
