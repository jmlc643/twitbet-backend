package request

type UpdateMarketStatusRequest struct {
	Status string `json:"status" binding:"required"`
}