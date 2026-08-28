package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis wraps a go-redis client with JSON helpers and pub/sub.
type Redis struct {
	Client *redis.Client
}

// New connects to Redis using a redis:// URL.
func New(url string) (*Redis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &Redis{Client: redis.NewClient(opts)}, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, b, ttl).Err()
}

func (r *Redis) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	b, err := r.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return false, err
	}
	return true, nil
}

// Publish marshals and publishes a message to a channel.
func (r *Redis) Publish(ctx context.Context, channel string, message interface{}) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return r.Client.Publish(ctx, channel, b).Err()
}

// Subscribe returns a PubSub for the given channel.
func (r *Redis) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channel)
}
