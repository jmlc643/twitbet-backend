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
