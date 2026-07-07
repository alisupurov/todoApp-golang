package statistics_cached_repository

import (
	"context"
	"errors"
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_redis_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/redis/pool"
	"go.uber.org/zap"
)

// GetTasks реализует cache-aside поверх StatisticsRepository.GetTasks:
//  1. Cache hit — список задач отдаётся из Redis, в основное хранилище не идём.
//  2. Cache miss/ошибка — читаем из mainRepository, кладём результат в кеш
//     с TTL пула, отдаём.
//
// Любая ошибка Redis (недоступность, (де)сериализация) только логируется —
// кеш не должен становиться точкой отказа запроса (graceful degradation).
func (r *CachedRepository) GetTasks(
	ctx context.Context,
	userId *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	log := core_logger.FromContext(ctx)

	key := cacheKey(userId, from, to)

	bytes, err := r.pool.Get(ctx, key).Bytes()
	if err != nil {
		if !errors.Is(err, core_redis_pool.NotFound) {
			log.Error("get cached statistics tasks", zap.Error(err))
		}
	} else {
		var model TaskListModel
		if err := model.Deserialize(bytes); err != nil {
			log.Error("deserialize cached statistics tasks", zap.Error(err))
		} else {
			return modelToDomains(model), nil
		}
	}

	tasks, err := r.mainRepository.GetTasks(ctx, userId, from, to)
	if err != nil {
		return nil, err
	}

	model := domainsToModel(tasks)

	data, err := model.Serialize()
	if err != nil {
		log.Error("serialize statistics tasks for cache", zap.Error(err))
	} else if err := r.pool.Set(ctx, key, data, r.pool.TTL()).Err(); err != nil {
		log.Error("set cached statistics tasks", zap.Error(err))
	}

	return tasks, nil
}
