package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/afifn11/gopay-x/services/auth-service/internal/domain"
)

type tokenRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTokenRepository(db *gorm.DB, redis *redis.Client) domain.TokenRepository {
	return &tokenRepository{db: db, redis: redis}
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > NOW()", token).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (r *tokenRepository) DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}

func (r *tokenRepository) BlacklistAccessToken(ctx context.Context, token string, expiryMinutes int) error {
	key := fmt.Sprintf("blacklist:%s", token)
	expiry := time.Duration(expiryMinutes) * time.Minute
	return r.redis.Set(ctx, key, "1", expiry).Err()
}

func (r *tokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}