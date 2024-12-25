package task

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CacheConfig struct {
	CacheKey     string
	QuerySQL     string
	CacheTTL     time.Duration
	UpdatePeriod time.Duration
}

type CacheController struct {
	config    *CacheConfig
	pool      *pgxpool.Pool
	rdb       *redis.Client
	startTime time.Time
}

func NewCacheController(pool *pgxpool.Pool, rdb *redis.Client, config *CacheConfig) *CacheController {
	return &CacheController{
		config:    config,
		pool:      pool,
		rdb:       rdb,
		startTime: time.Now(),
	}
}

func (cu *CacheController) Start(ctx context.Context, aggregator *StatsAggregator) {
	ticker := time.NewTicker(cu.config.UpdatePeriod)
	defer ticker.Stop()
	cu.updateCache(ctx, aggregator)
	for {
		select {
		case <-ctx.Done():
			log.Println("Контекст завершен. Остановка выполнения запросов")
			return
		case <-ticker.C:
			cu.updateCache(ctx, aggregator)
		}
	}
}

func (cu *CacheController) updateCache(ctx context.Context, aggregator *StatsAggregator) {
	start := time.Now()

	data, err := cu.getCachedData(ctx)
	if err != nil {
		log.Println("Ошибка получения данных из кэша")
		return
	}

	if data == nil {
		data, err = cu.fetchDatabaseData(ctx)
		if err != nil {
			log.Println("Ошибка получения данных из базы")
			return
		}

		if err := cu.setCacheData(ctx, data); err != nil {
			log.Println("Ошибка обновления кэша")
			return
		}
	}

	elapsed := time.Since(start)
	secondsSinceStart := int(time.Since(cu.startTime).Seconds())

	aggregator.RecordCacheStat(secondsSinceStart, elapsed)

	log.Printf("Запрос к Redis выполнен за %v", elapsed)
}

func (cu *CacheController) getCachedData(ctx context.Context) (map[string]int, error) {
	jsonData, err := cu.rdb.Get(ctx, cu.config.CacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var data map[string]int
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (cu *CacheController) setCacheData(ctx context.Context, data map[string]int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cu.rdb.Set(ctx, cu.config.CacheKey, jsonData, cu.config.CacheTTL).Err()
}

func (cu *CacheController) fetchDatabaseData(ctx context.Context) (map[string]int, error) {
	rows, err := cu.pool.Query(ctx, cu.config.QuerySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]int)
	for rows.Next() {
		var city string
		var deliveries int
		if err := rows.Scan(&city, &deliveries); err != nil {
			return nil, err
		}
		data[city] = deliveries
	}

	return data, rows.Err()
}
