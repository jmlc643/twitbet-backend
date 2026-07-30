package adapter

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	secretKey []byte
}

func NewJWTTokenService(secret string) *JWTTokenService {
	return &JWTTokenService{
		secretKey: []byte(secret),
	}
}

func (j *JWTTokenService) GenerateToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTTokenService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma no válido")
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("token inválido o expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error al leer las claims del token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("id de usuario no encontrado en el token")
	}

	return userID, nil
}
