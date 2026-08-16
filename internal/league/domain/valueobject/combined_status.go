package valueobject

type CombinedStatus string

const (
	CombinedStatusPending  CombinedStatus = "PENDING"
	CombinedStatusAccepted CombinedStatus = "ACCEPTED"
	CombinedStatusWon      CombinedStatus = "WON"
	CombinedStatusLost     CombinedStatus = "LOST"
	CombinedStatusCashout  CombinedStatus = "CASHOUT"
)
