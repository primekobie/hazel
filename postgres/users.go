package postgres

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	conn *pgxpool.Pool
}

func NewUserStore(conn *pgxpool.Pool) models.UserStore {
	return &UserStore{
		conn: conn,
	}
}

// InsertUser implements models.UserStore.
func (u *UserStore) InsertUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, profile_photo, created_at, updated_at, verified)
		VALUES ($1, NULLIF($2,''), $3, $4, $5, $6, $7, $8);`

	_, err := u.conn.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.ProfilePhoto,
		user.CreatedAt,
		user.UpdatedAt,
		user.Verified,
	)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return models.ErrDuplicateUser
		}
		slog.Error("failed to insert user", "error", err)
		return err
	}
	return nil
}

// DeleteUser implements models.UserStore.
func (u *UserStore) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := u.conn.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed delete user", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

// GetUser implements models.UserStore.
func (u *UserStore) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	query := `
		SELECT id, name, email, password_hash, profile_photo, created_at, updated_at, verified 
		FROM users 
		WHERE id = $1;`

	var user models.User
	err := u.conn.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePhoto,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, models.ErrNotFound
	}

	return user, nil
}

// GetUserByMail implements models.UserStore.
func (u *UserStore) GetUserByMail(ctx context.Context, email string) (models.User, error) {
	query := `
		SELECT id, name, email, password_hash, profile_photo, created_at, updated_at, verified 
		FROM users 
		WHERE email = $1;`

	var user models.User
	err := u.conn.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePhoto,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, models.ErrNotFound
	}

	return user, nil
}

// UpdateUser implements models.UserStore.
func (u *UserStore) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password_hash = $3, profile_photo = $4, updated_at = $5, verified = $6
		WHERE id = $7;`

	result, err := u.conn.Exec(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.ProfilePhoto,
		user.UpdatedAt,
		user.Verified,
		user.Id,
	)
	if err != nil {
		slog.Error("failed update user", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

// InsertToken implements models.TokenStore.
func (t *UserStore) InsertToken(ctx context.Context, token *models.UserToken) error {
	query := `INSERT INTO user_tokens(token_hash, user_id, scope, expires_at)
	VALUES($1, $2, $3, $4);`

	_, err := t.conn.Exec(ctx, query, token.Hash, token.UserId, token.Scope, token.ExpiresAt)
	if err != nil {
		slog.Error("failed to insert token", "error", err)
		return err
	}

	return nil
}

// GetUserForToken implements models.UserStore
func (t *UserStore) GetUserForToken(ctx context.Context, tokenHash string, scope string, email string) (models.User, error) {
	query := `SELECT
	users.id,
	users.name,
	users.email,
	users.password_hash,
	users.profile_photo,
	users.verified,
	users.created_at,
	users.updated_at
	FROM users
	JOIN user_tokens AS tokens
	ON users.id = tokens.user_id
	WHERE tokens.token_hash = $1
	AND tokens.scope = $2
	AND users.email = $3
	AND tokens.expires_at > now();
	`

	var user models.User
	row := t.conn.QueryRow(ctx, query, tokenHash, scope, email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash, &user.ProfilePhoto, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrNotFound
		}
		slog.Error("failed to fetch token", "error", err)
		return models.User{}, err
	}

	return user, nil
}

// DeleteToken implements models.TokenStore.
func (t *UserStore) DeleteToken(ctx context.Context, tokenHash, scope string) error {
	query := `DELETE FROM user_tokens WHERE token_hash = $1 AND scope = $2;`

	_, err := t.conn.Exec(ctx, query, tokenHash, scope)
	if err != nil {
		slog.Error("failed to delete user token", "error", err)
		return err
	}

	return nil
}
