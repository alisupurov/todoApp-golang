package core_goredis_pool

import (
	"errors"

	core_redis_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/redis/pool"
	"github.com/redis/go-redis/v9"
)

// goredisStringCmd оборачивает *redis.StringCmd и реализует core_redis_pool.StringCmd.
type goredisStringCmd struct {
	*redis.StringCmd
}

func (c goredisStringCmd) Bytes() ([]byte, error) {
	data, err := c.StringCmd.Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return data, nil
}

// goredisStatusCmd оборачивает *redis.StatusCmd и реализует core_redis_pool.StatusCmd.
type goredisStatusCmd struct {
	*redis.StatusCmd
}

// mapError транслирует redis.Nil (ключ не найден) в core_redis_pool.NotFound,
// чтобы вызывающий код не зависел от конкретной библиотеки клиента.
func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_redis_pool.NotFound
	}

	return err
}
