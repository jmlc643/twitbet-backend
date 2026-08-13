package request

type UpdateMarketOptionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}