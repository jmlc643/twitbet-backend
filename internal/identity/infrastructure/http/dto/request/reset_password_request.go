package request

type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email"`
	OTPCode         string `json:"otp_code" binding:"required,len=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
