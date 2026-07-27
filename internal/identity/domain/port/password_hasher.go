package port

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hashedPassword string) error
}