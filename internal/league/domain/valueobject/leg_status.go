package valueobject

type LegStatus string

const (
	LegStatusPending LegStatus = "PENDING"
	LegStatusWon     LegStatus = "WON"
	LegStatusLost    LegStatus = "LOST"
	LegStatusVoided  LegStatus = "VOIDED"
	LegStatusCashout LegStatus = "CASHOUT"
)
