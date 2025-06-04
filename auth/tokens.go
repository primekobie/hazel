package auth

import (
	"log/slog"
	"os"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserSession struct {
	User         models.User `json:"user"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresAt    time.Time   `json:"expiresAt"`
}

func GenerateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().UTC().UnixNano(),
		"exp": time.Now().Add(duration).UnixNano(),
		"iss": userID.String(),
	})

	tokenString, err := token.SignedString(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		slog.Error("failed to sign access token", "error", err.Error())
		return "", err
	}

	return tokenString, nil
}
