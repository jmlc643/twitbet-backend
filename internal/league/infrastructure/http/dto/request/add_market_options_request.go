package request

type AddMarketOptionsRequest struct {
	Options []MarketOptionRequest `json:"options" binding:"required,min=1,dive"`
}