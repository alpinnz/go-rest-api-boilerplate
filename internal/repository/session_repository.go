package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) repository.SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) SetAccessToken(ctx context.Context, token string, userID uuid.UUID, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return r.client.Set(ctx, key, userID.String(), expiration).Err()
}

func (r *sessionRepository) GetAccessToken(ctx context.Context, token string) (uuid.UUID, error) {
	key := fmt.Sprintf("session:%s", token)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return uuid.Nil, domain.ErrSessionExpired
	}
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *sessionRepository) DeleteAccessToken(ctx context.Context, token string) error {
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
