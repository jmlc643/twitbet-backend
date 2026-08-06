package request

type UpdateMatchStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
