package postgres

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/persistence/mapper"
	"gorm.io/gorm"
)

type LeagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) *LeagueRepository {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) CreateLeagueWithAdminParticipant(ctx context.Context, l *entity.League, p *entity.Participant) error {
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

	return tx.Commit().Error
}
