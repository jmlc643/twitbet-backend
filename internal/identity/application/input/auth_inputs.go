package input

type VerifyAccountInput struct {
	Email   string
	OTPCode string
}

type ForgotPasswordInput struct {
	Email string
}

type ResetPasswordInput struct {
	Email           string
	OTPCode         string
	NewPassword     string
	ConfirmPassword string
}

type ChangePasswordInput struct {
	UserID          string
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

type VerifyResetOtpInput struct {
	Email   string
	OTPCode string
}
