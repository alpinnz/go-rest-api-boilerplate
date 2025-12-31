package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) domain.SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) Set(ctx context.Context, token string, userID int64, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return r.client.Set(ctx, key, userID, expiration).Err()
}

func (r *sessionRepository) Get(ctx context.Context, token string) (int64, error) {
	key := fmt.Sprintf("session:%s", token)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, domain.ErrSessionExpired
	}
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return r.client.Del(ctx, key).Err()
}

func (r *sessionRepository) SetRefreshToken(ctx context.Context, refreshToken string, accessToken string, expiration time.Duration) error {
	key := fmt.Sprintf("refresh:%s", refreshToken)
	return r.client.Set(ctx, key, accessToken, expiration).Err()
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := fmt.Sprintf("refresh:%s", refreshToken)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", domain.ErrSessionExpired
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *sessionRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf("refresh:%s", refreshToken)
	return r.client.Del(ctx, key).Err()
}
