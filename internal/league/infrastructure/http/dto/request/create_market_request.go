package request

type MarketOptionRequest struct {
	Name string  `json:"name" binding:"required"`
	Odds float64 `json:"odds" binding:"required,gt=1"`
}

type CreateMarketRequest struct {
	Name    string                `json:"name" binding:"required"`
	Options []MarketOptionRequest `json:"options" binding:"required,min=2"`
}