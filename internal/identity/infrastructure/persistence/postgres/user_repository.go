package postgres

import (
	"context"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/mapper"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Create(ctx context.Context, user *entity.User) error {
	dbModel := mapper.UserEntityToGORM(user)

	if err := r.db.WithContext(ctx).Create(&dbModel).Error; err != nil {
		return err
	}

	user.ID = dbModel.ID
	user.CreatedAt = dbModel.CreatedAt
	return nil
}

func (r *UserGormRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var dbModel model.UserModel

	result := r.db.WithContext(ctx).Where("email = ?", email).Limit(1).Find(&dbModel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return mapper.UserGORMToEntity(&dbModel), nil
}

func (r *UserGormRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var dbModel model.UserModel

	result := r.db.WithContext(ctx).Where("id = ?", id).Limit(1).Find(&dbModel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return mapper.UserGORMToEntity(&dbModel), nil
}

func (r *UserGormRepository) Update(ctx context.Context, user *entity.User) error {
	dbModel := mapper.UserEntityToGORM(user)

	updateData := map[string]interface{}{
		"username":      dbModel.Username,
		"password_hash": dbModel.PasswordHash,
		"is_verified":   dbModel.IsVerified,
	}

	if dbModel.AvatarURL == nil {
		updateData["avatar_url"] = gorm.Expr("NULL")
	} else {
		updateData["avatar_url"] = dbModel.AvatarURL
	}

	result := r.db.WithContext(ctx).Model(&model.UserModel{}).Where("id = ?", user.ID).Updates(updateData)
	
	if result.Error != nil {
		return result.Error
	}
	
	return nil
}

func (r *UserGormRepository) GetUserStats(ctx context.Context, userID string) (*entity.UserStats, error) {
	var totalLeagues int64
	if err := r.db.WithContext(ctx).Table("league_participants").Where("user_id = ?", userID).Count(&totalLeagues).Error; err != nil {
		return nil, err
	}

	var totalBets int64
	if err := r.db.WithContext(ctx).Table("bets").
		Joins("INNER JOIN league_participants ON bets.participant_id = league_participants.id").
		Where("league_participants.user_id = ?", userID).
		Count(&totalBets).Error; err != nil {
		return nil, err
	}

	var wins int64
	if err := r.db.WithContext(ctx).Table("bets").
		Joins("INNER JOIN league_participants ON bets.participant_id = league_participants.id").
		Where("league_participants.user_id = ? AND bets.status = ?", userID, "WON").
		Count(&wins).Error; err != nil {
		return nil, err
	}

	effectiveness := 0
	if totalBets > 0 {
		effectiveness = int((float64(wins) / float64(totalBets)) * 100)
	}

	return &entity.UserStats{
		Leagues:       int(totalLeagues),
		Wins:          int(wins),
		Effectiveness: effectiveness,
	}, nil
}
