package statistics_cached_repository

import (
	"context"
	"time"

	"github.com/alisupurov/todoApp-golang/internal/core/domain"
	core_redis_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/redis/pool"
)

// StatisticsRepository повторяет контракт statistics_service.StatisticsRepository.
// Дублирование намеренное: пакет реализует интерфейс через структурную
// совместимость Go и не должен зависеть от пакета service.
type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userId *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

// CachedRepository — Decorator над StatisticsRepository: добавляет Redis-кеш
// поверх любой реализации (сейчас — Postgres), не меняя контракт интерфейса.
// Service-слой работает с тем же интерфейсом и не знает о существовании кеша.
//
// Стратегия — cache-aside с TTL: при чтении статистики данные на TTL.Seconds()
// могут быть устаревшими, что приемлемо для агрегации, не требующей real-time
// точности. Любая ошибка Redis (недоступность, (де)сериализация) только
// логируется — кеш не должен становиться точкой отказа запроса.
type CachedRepository struct {
	pool           core_redis_pool.Pool
	mainRepository StatisticsRepository
}

func NewCachedRepository(
	pool core_redis_pool.Pool,
	mainRepository StatisticsRepository,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
	}
}
