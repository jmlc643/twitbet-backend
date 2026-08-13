package request

type ResolveMarketRequest struct {
	WinningOptionIDs []string `json:"winning_option_ids" binding:"required,min=1,dive,uuid"`
}