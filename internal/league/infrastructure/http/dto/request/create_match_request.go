package request

import "time"

type CreateMatchRequest struct {
	Title     string    `json:"title" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
}
