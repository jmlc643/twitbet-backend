package request

type UpdateLeagueStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
