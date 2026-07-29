package entity

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type League struct {
	ID             uuid.UUID
	AdminID        uuid.UUID
	Name           string
	InitialBalance float64
	MaxRecharges   int
	HideStandings  bool
	InviteCode     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewLeague(adminID uuid.UUID, name string, initialBalance float64, maxRecharges int, hideStandings bool) (*League, error) {
	if name == "" {
		return nil, apperror.ErrInvalidLeagueName
	}
	if initialBalance < 0 {
		return nil, apperror.ErrInvalidInitialBalance
	}

	if maxRecharges <= 0 {
		maxRecharges = 2
	}

	inviteCode, err := generateInviteCode(8)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &League{
		ID:             uuid.New(),
		AdminID:        adminID,
		Name:           name,
		InitialBalance: initialBalance,
		MaxRecharges:   maxRecharges,
		HideStandings:  hideStandings,
		InviteCode:     inviteCode,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func generateInviteCode(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}