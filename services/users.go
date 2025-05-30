package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/freekobie/hazel/mail"
	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnverifiedUser     = errors.New("user has an unverified email")
	ErrInvalidToken       = errors.New("token is invalid or expired")
)

const (
	VERIFICATION   = "verification"
	AUTHENTICATION = "authentication"
)

type UserService struct {
	store models.UserStore
	mail  *mail.Mailer
}

func NewUserService(us models.UserStore, m *mail.Mailer) *UserService {
	return &UserService{
		store: us,
		mail:  m,
	}
}

// CreateUser creates a new user with the given details
func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, err
	}
	now := time.Now().UTC()
	user := &models.User{
		Id:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
		Verified:     false,
	}

	if err := s.store.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	otpString := generateOTP()
	otpHash := hashString(otpString)

	userAddr := mail.Address{Name: user.Name, Email: user.Email}
	data := mail.Data{
		Address: userAddr,
		Code:    otpString,
	}

	token := models.UserToken{
		Hash:      otpHash,
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Scope:     VERIFICATION,
	}

	_ = s.store.InsertToken(ctx, &token)

	s.sendEmail([]mail.Address{userAddr}, "verify_email.html", data)

	return user, nil
}

func (us *UserService) VerifyUser(ctx context.Context, code string, email string) (models.User, error) {

	hash := hashString(code)
	user, err := us.store.GetUserForToken(ctx, hash, VERIFICATION, email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.User{}, ErrInvalidToken
		}
		return models.User{}, err
	}

	user.Verified = true

	err = us.store.UpdateUser(ctx, &user)
	if err != nil {
		return models.User{}, err
	}

	// Delete otp after successful verification
	_ = us.store.DeleteToken(ctx, hash, VERIFICATION)

	address := mail.Address{Name: user.Name, Email: user.Email}
	us.sendEmail([]mail.Address{address}, "welcome_email.html", mail.Data{Address: address})

	return user, nil
}

func (us *UserService) ResendOTP(ctx context.Context, email string) error {
	user, err := us.store.GetUserByMail(ctx, email)
	if err != nil {
		return err
	}

	otpString := generateOTP()
	otpHash := hashString(otpString)

	userAddr := mail.Address{Email: email, Name: user.Name}
	data := mail.Data{
		Address: userAddr,
		Code:    otpString,
	}

	token := models.UserToken{
		Hash:      otpHash,
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Scope:     VERIFICATION,
	}

	err = us.store.InsertToken(ctx, &token)
	if err != nil {
		return err
	}

	us.sendEmail([]mail.Address{userAddr}, "verify_email.html", data)

	return nil
}

func (us *UserService) NewSession(context context.Context, email string, password string) (any, error) {
	panic("unimplemented")
}

// UpdateUser updates an existing user's details
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now().UTC()
	return s.store.UpdateUser(ctx, user)
}

// FetchUser retrieves a user by ID or email
func (s *UserService) FetchUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user from the system
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.store.DeleteUser(ctx, id)
}
