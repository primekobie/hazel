package auth

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("the provided token is no valid")
)

type TokenType string
type CustomClaims struct {
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	TokenTypeAccess  TokenType = "ACCESS"
	TokenTypeRefresh TokenType = "REFRESH"
)

type UserSession struct {
	User         models.User `json:"user"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresAt    time.Time   `json:"expiresAt"`
}

type UserAccess struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func GenerateToken(userID uuid.UUID, email string, duration time.Duration, tokenType TokenType) (string, error) {
	exp := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":        time.Now().UTC().Unix(),
		"exp":        exp.UTC().Unix(),
		"sub":        userID.String(),
		"token_type": tokenType,
		"email":      email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		slog.Error("failed to sign access token", "error", err.Error())
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string, tokenType TokenType) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	if claims.TokenType != string(tokenType) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
