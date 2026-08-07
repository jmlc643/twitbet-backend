package response

type RechargeResponse struct {
	Balance           float64 `json:"balance"`
	RechargesConsumed int     `json:"recharges_consumed"`
}
