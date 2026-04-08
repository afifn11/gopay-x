package domain

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error
	BlacklistAccessToken(ctx context.Context, token string, expiry int) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}