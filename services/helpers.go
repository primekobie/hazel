package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/primekobie/hazel/mail"
)

// GenerateOTP generates a 6-digit OTP as a string
func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))

	return fmt.Sprintf("%06d", n.Int64())
}

// hashString hashes the token using SHA-256 and returns the hex-encoded hash
func hashString(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (us *UserService) sendEmail(recipients []mail.Address, template string, data any) {
	// send email
	go func() {
		err := us.mail.Send(recipients, template, data)
		if err != nil {
			slog.Error("failed  to send email", "error", err)
			return
		}
	}()
}
