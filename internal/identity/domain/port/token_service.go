package port

type TokenService interface {
	GenerateToken(userID, email string) (string, error)
	ValidateToken(token string) (string, error)
}
