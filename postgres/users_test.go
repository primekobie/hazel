package postgres_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/primekobie/hazel/models"
	"github.com/primekobie/hazel/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	connStr := os.Getenv("TEST_DB_URL")
	if connStr == "" {
		connStr = "postgres://userone:pa88word@localhost:5432/hazel_test?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	return pool
}

func createTestUser(name, email string) *models.User {
	return &models.User{
		Id:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: []byte("hashedpassword"),
		CreatedAt:    time.Now().UTC(),
		LastModifed:  time.Now().UTC(),
		Verified:     true,
	}
}

func generateTestEmail() string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s@example.com", string(b))
}

func TestUserStore_CreateUser(t *testing.T) {
	pool := setupTestDB(t)
	store := postgres.NewUserStore(pool)
	ctx := context.Background()

	email := generateTestEmail()

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    createTestUser("Test User", email),
			wantErr: false,
		},
		{
			name:    "duplicate email",
			user:    createTestUser("Another User", email),
			wantErr: true,
		},
		{
			name: "empty name",
			user: &models.User{
				Id:           uuid.New(),
				Email:        generateTestEmail(),
				PasswordHash: []byte("hashedpassword"),
				CreatedAt:    time.Now().UTC(),
				LastModifed:  time.Now().UTC(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.InsertUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			got, err := store.GetUser(ctx, tt.user.Id)
			require.NoError(t, err)
			assert.Equal(t, tt.user.Id, got.Id)
			assert.Equal(t, tt.user.Name, got.Name)
			assert.Equal(t, tt.user.Email, got.Email)
		})
	}
}

func TestUserStore_GetUser(t *testing.T) {
	pool := setupTestDB(t)
	store := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Get Test", generateTestEmail())
	require.NoError(t, store.InsertUser(ctx, user))

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      user.Id,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetUser(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, user.Id, got.Id)
		})
	}
}

func TestUserStore_UpdateUser(t *testing.T) {
	pool := setupTestDB(t)
	store := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Update Test", generateTestEmail())
	require.NoError(t, store.InsertUser(ctx, user))

	tests := []struct {
		name    string
		user    *models.User
		updates func(*models.User)
		wantErr bool
	}{
		{
			name: "valid update",
			user: user,
			updates: func(u *models.User) {
				u.Name = "Updated Name"
				u.Email = generateTestEmail()
			},
			wantErr: false,
		},
		{
			name: "non-existent user",
			user: createTestUser("Non-existent", "Kwame@email.com"),
			updates: func(u *models.User) {
				u.Name = "Updated Name"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updates != nil {
				tt.updates(tt.user)
			}
			tt.user.LastModifed = time.Now().UTC()

			err := store.UpdateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			got, err := store.GetUser(ctx, tt.user.Id)
			require.NoError(t, err)
			assert.Equal(t, tt.user.Name, got.Name)
			assert.Equal(t, tt.user.Email, got.Email)
		})
	}
}

func TestUserStore_DeleteUser(t *testing.T) {
	pool := setupTestDB(t)
	store := postgres.NewUserStore(pool)
	ctx := context.Background()

	user := createTestUser("Delete Test", generateTestEmail())
	require.NoError(t, store.InsertUser(ctx, user))

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      user.Id,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      uuid.New(),
			wantErr: true,
		},
		{
			name:    "invalid ID",
			id:      uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.DeleteUser(ctx, tt.id.String())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.NoError(t, err)

			_, err = store.GetUser(ctx, tt.id)
			assert.Error(t, err)
		})
	}
}
