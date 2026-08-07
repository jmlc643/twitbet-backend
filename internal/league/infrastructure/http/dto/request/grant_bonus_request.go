package request

type GrantBonusRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
