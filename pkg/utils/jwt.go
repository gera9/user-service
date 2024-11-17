package utils

import (
	"fmt"
	"time"

	"github.com/gera9/user-service/config"
	"github.com/gera9/user-service/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = fmt.Errorf("expired token")
)

type CustomClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func CreateToken(cfg config.Config, user *models.User) (string, error) {
	claims := CustomClaims{
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-service",
			Subject:   user.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret"))
}

func ParseAndValidateToken(cfg config.Config, tokenString string) (*models.User, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &models.User{
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}
