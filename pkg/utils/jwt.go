package utils

import (
	"fmt"
	"log"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
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

// CreateTestingToken creates a token for testing purposes
// with the given username and email and returns it as a string.
// The token is signed with the word "secret".
func CreateTestingToken(username, email string) string {
	claims := CustomClaims{
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user-service",
			Subject:   "8bb48920-bc63-4c84-9bac-6e60cfd06f27",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Panic(err)
	}

	return tokenString
}
