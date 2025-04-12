package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jaimesHub/golang-todo-app/internal/config"
)

// Client represents a Redis client
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping Redis to check connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client: client}, nil
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// Set sets a key-value pair in Redis
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return c.client.Set(ctx, key, value, 0).Err()
}

// Get gets a value from Redis by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete deletes a key from Redis
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to a channel
func (c *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.client.Subscribe(ctx, channel)
}
