package shortener

import (
	"context"
	"time"

	"shorty/slug"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Exists(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisURLShortener struct {
	client RedisClient
}

func NewRedisURLShortener(client RedisClient) *RedisURLShortener {
	return &RedisURLShortener{client: client}
}

func (r *RedisURLShortener) ShortenURL(longURL string) (string, error) {
	for {
		slug := slug.GenerateRandomSlug()

		exists, err := r.client.Exists(context.Background(), slug)
		if err != nil {
			return "", err
		}

		if exists == 0 {
			err := r.client.Set(context.Background(), slug, longURL, 0)
			if err != nil {
				return "", err
			}

			return slug, nil
		}
	}
}

func (r *RedisURLShortener) Redirect(shortCode string) (string, error) {
	longURL, err := r.client.Get(context.Background(), shortCode)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

type RealRedisClient struct {
	client *redis.Client
}

func NewRealRedisClient(client *redis.Client) *RealRedisClient {
	return &RealRedisClient{client: client}
}

func (r *RealRedisClient) Exists(ctx context.Context, key string) (int64, error) {
	return r.client.Exists(ctx, key).Result()
}

func (r *RealRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RealRedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}
