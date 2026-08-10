package mapper

import (
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/identity/infrastructure/persistence/model"
)

func UserEntityToGORM(u *entity.User) *model.UserModel {
	var avatar *string
	if u.AvatarURL != "" {
		avatar = &u.AvatarURL
	}

	return &model.UserModel{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		AvatarURL:    avatar,
		IsVerified:   u.IsVerified,
		CreatedAt:    u.CreatedAt,
	}
}

func UserGORMToEntity(m *model.UserModel) *entity.User {
	if m == nil {
		return nil
	}

	avatar := ""
	if m.AvatarURL != nil {
		avatar = *m.AvatarURL
	}

	return &entity.User{
		ID:           m.ID,
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		AvatarURL:    avatar,
		IsVerified:   m.IsVerified,
		CreatedAt:    m.CreatedAt,
	}
}
