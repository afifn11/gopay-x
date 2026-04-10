package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/afifn11/gopay-x/services/wallet-service/internal/domain"
)

type lockRepository struct {
	redis *redis.Client
}

func NewLockRepository(rdb *redis.Client) domain.LockRepository {
	return &lockRepository{redis: rdb}
}

func (r *lockRepository) AcquireLock(ctx context.Context, key string, ttlSeconds int) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	result, err := r.redis.SetNX(ctx, lockKey, "1", time.Duration(ttlSeconds)*time.Second).Result()
	return result, err
}

func (r *lockRepository) ReleaseLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return r.redis.Del(ctx, lockKey).Err()
}