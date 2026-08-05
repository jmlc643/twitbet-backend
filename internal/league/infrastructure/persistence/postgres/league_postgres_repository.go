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

	OwnerID, _ := uuid.Parse(modelLeague.OwnerID)
	id, _ := uuid.Parse(modelLeague.ID)

	return &entity.League{
		ID:             id,
		OwnerID:        OwnerID,
		Name:           modelLeague.Name,
		Slug:           modelLeague.Slug,
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
		Slug             string
		Name             string
		OwnerID          string
		Balance          float64
		ParticipantCount int
	}

	query := `
		SELECT 
			l.id as league_id, 
			l.slug,
			l.name, 
			l.owner_id, 
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
		if res.OwnerID == userID.String() {
			role = entity.RoleAdmin
		}

		leagueUUID, _ := uuid.Parse(res.LeagueID)

		summaries = append(summaries, entity.LeagueSummary{
			LeagueID:         leagueUUID,
			Slug:             res.Slug,
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

func (r *LeagueRepository) GetLeagueByID(ctx context.Context, id uuid.UUID) (*entity.League, error) {
	var modelLeague model.LeagueModel
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&modelLeague).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	OwnerID, _ := uuid.Parse(modelLeague.OwnerID)
	leagueID, _ := uuid.Parse(modelLeague.ID)

	return &entity.League{
		ID:             leagueID,
		OwnerID:        OwnerID,
		Name:           modelLeague.Name,
		Slug:           modelLeague.Slug,
		InitialBalance: modelLeague.InitialBalance,
		MaxRecharges:   modelLeague.MaxRecharges,
		HideStandings:  modelLeague.HideStandings,
		InviteCode:     modelLeague.InviteCode,
		CreatedAt:      modelLeague.CreatedAt,
		UpdatedAt:      modelLeague.UpdatedAt,
	}, nil
}

func (r *LeagueRepository) GetLeagueBySlug(ctx context.Context, slug string) (*entity.League, error) {
	var modelLeague model.LeagueModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&modelLeague).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	OwnerID, _ := uuid.Parse(modelLeague.OwnerID)
	leagueID, _ := uuid.Parse(modelLeague.ID)

	return &entity.League{
		ID:             leagueID,
		OwnerID:        OwnerID,
		Name:           modelLeague.Name,
		Slug:           modelLeague.Slug,
		InitialBalance: modelLeague.InitialBalance,
		MaxRecharges:   modelLeague.MaxRecharges,
		HideStandings:  modelLeague.HideStandings,
		InviteCode:     modelLeague.InviteCode,
		CreatedAt:      modelLeague.CreatedAt,
		UpdatedAt:      modelLeague.UpdatedAt,
	}, nil
}

func (r *LeagueRepository) UpdateLeague(ctx context.Context, league *entity.League) error {
	leagueModel := mapper.EntityToLeagueModel(league)
	return r.db.WithContext(ctx).Save(leagueModel).Error
}

func (r *LeagueRepository) GetParticipantsByLeagueID(ctx context.Context, id uuid.UUID) ([]entity.ParticipantRanking, error) {
	var results []struct {
		ParticipantID string
		UserID        string
		Username      string
		AvatarURL     *string
		Balance       float64
		Position      int
	}

	query := `
		SELECT 
			p.id as participant_id,
			p.user_id,
			u.username,
			u.avatar_url,
			p.balance,
			RANK() OVER (ORDER BY p.balance DESC) as position
		FROM league_participants p
		JOIN users u ON p.user_id = u.id
		WHERE p.league_id = ?
		ORDER BY p.balance DESC
	`

	err := r.db.WithContext(ctx).Raw(query, id.String()).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var rankings []entity.ParticipantRanking
	for _, res := range results {
		participantID, _ := uuid.Parse(res.ParticipantID)
		userID, _ := uuid.Parse(res.UserID)

		rankings = append(rankings, entity.ParticipantRanking{
			ParticipantID:  participantID,
			UserID:         userID,
			Username:       res.Username,
			ProfilePicture: res.AvatarURL,
			Balance:        res.Balance,
			Position:       res.Position,
		})
	}

	return rankings, nil
}

func (r *LeagueRepository) DeleteLeague(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if err := tx.Where("league_id = ?", id.String()).Delete(&model.TransactionModel{}).Error; err != nil {
		return err
	}
	if err := tx.Where("league_id = ?", id.String()).Delete(&model.ParticipantModel{}).Error; err != nil {
		return err
	}
	if err := tx.Where("id = ?", id.String()).Delete(&model.LeagueModel{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *LeagueRepository) GetParticipant(ctx context.Context, leagueID, userID uuid.UUID) (*entity.Participant, error) {
	var pModel model.ParticipantModel
	err := r.db.WithContext(ctx).Where("league_id = ? AND user_id = ?", leagueID.String(), userID.String()).First(&pModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	pID, _ := uuid.Parse(pModel.ID)
	return &entity.Participant{
		ID:                pID,
		LeagueID:          leagueID,
		UserID:            userID,
		IsAdmin:           pModel.IsAdmin,
		Balance:           pModel.Balance,
		RechargesConsumed: pModel.RechargesConsumed,
		JoinedAt:          pModel.JoinedAt,
	}, nil
}

func (r *LeagueRepository) UpdateParticipant(ctx context.Context, p *entity.Participant) error {
	pModel := mapper.EntityToParticipantModel(p)
	return r.db.WithContext(ctx).Save(pModel).Error
}
