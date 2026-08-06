package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeInitialBalance TransactionType = "INITIAL_BALANCE"
	TransactionTypeRecharge       TransactionType = "RECHARGE"
	TransactionTypeBet            TransactionType = "BET"
	TransactionTypeWin            TransactionType = "WIN"
	TransactionTypeRefund         TransactionType = "REFUND"
)

type Transaction struct {
	ID        uuid.UUID
	LeagueID  uuid.UUID
	UserID    uuid.UUID
	Amount    float64
	Type      TransactionType
	CreatedAt time.Time
}

func NewTransaction(leagueID, userID uuid.UUID, amount float64, txType TransactionType) (*Transaction, error) {
	return &Transaction{
		ID:        uuid.New(),
		LeagueID:  leagueID,
		UserID:    userID,
		Amount:    amount,
		Type:      txType,
		CreatedAt: time.Now().UTC(),
	}, nil
}
