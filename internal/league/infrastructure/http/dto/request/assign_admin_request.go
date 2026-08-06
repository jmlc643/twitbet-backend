package request

type AssignAdminRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
}
