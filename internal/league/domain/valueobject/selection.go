package valueobject

import "github.com/google/uuid"

type Selection struct {
	MarketID     uuid.UUID
	SelectionID  uuid.UUID
	AcceptedOdds float64
}
