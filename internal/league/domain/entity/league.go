package entity

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
)

type League struct {
	ID             uuid.UUID
	OwnerID        uuid.UUID
	Name           string
	Slug           string
	InitialBalance float64
	MaxRecharges   int
	HideStandings  bool
	InviteCode     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewLeague(ownerID uuid.UUID, name string, initialBalance float64, maxRecharges int, hideStandings bool) (*League, error) {
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
		OwnerID:        ownerID,
		Name:           name,
		Slug:           generateSlug(name),
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
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	
	uuidPart := uuid.New().String()[:6]
	return slug + "-" + uuidPart
}
