package request

type ResolveMarketRequest struct {
	WinningOptionID string `json:"winning_option_id" binding:"required,uuid"`
}
