package core_goredis_pool

import (
	"context"
	"fmt"
	"time"

	core_redis_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{
		client: client,
		ttl:    config.TTL,
	}, nil
}

func (p *Pool) Get(
	ctx context.Context,
	key string,
) core_redis_pool.StringCmd {
	return goredisStringCmd{p.client.Get(ctx, key)}
}

func (p *Pool) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) core_redis_pool.StatusCmd {
	return goredisStatusCmd{p.client.Set(ctx, key, value, ttl)}
}

func (p *Pool) Del(
	ctx context.Context,
	keys ...string,
) error {
	return p.client.Del(ctx, keys...).Err()
}

func (p *Pool) Close() error {
	return p.client.Close()
}

func (p *Pool) TTL() time.Duration {
	return p.ttl
}
