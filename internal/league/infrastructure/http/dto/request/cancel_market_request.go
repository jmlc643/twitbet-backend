package request

type CancelMarketRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}
