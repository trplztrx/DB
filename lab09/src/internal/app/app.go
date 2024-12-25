package app

import (
	"context"
	"fmt"
	"lab09/config"
	task "lab09/internal/tasks"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	query = `
			select a.city, count(o.id) AS total_deliveries
			from orders o
			join addresses a on o.delivery_address_id = a.id
			group by a.city
			order by total_deliveries desc
			limit 5`
	updatePeriod = 5 * time.Second
	cacheTTL     = 5 * time.Second
)

func RunApp(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()

	pool, err := connectDB(cfg, ctx)
	if err != nil {
		log.Fatalf("Ошибка подключения к postgres: %v", err)
	}
	defer pool.Close()
	log.Println("Подключение к postgres успешно выполнено")

	rdb, err := connectRedis(cfg, ctx)
	if err != nil {
		log.Fatalf("Ошибка подключения к redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Ошибка закрытия Redis: %v", err)
		}
	}()
	log.Println("Подключение к redis успешно выполнено")

	scenarioTypes := []task.OperationType{
		task.None,
		task.Insert,
		task.Delete,
		task.Update,
	}

	for _, scenario := range scenarioTypes {
		log.Printf("\nНачало выполнения сценария: %s", scenario)

		aggregator := task.NewStatsAggregator()

		scenarioCtx, scenarioCancel := context.WithTimeout(ctx, 36*time.Second)
		defer scenarioCancel()

		statsConfig := &task.DeliveryStatsConfig{
			Query:        query,
			UpdatePeriod: updatePeriod,
		}
		cacheConfig := &task.CacheConfig{
			CacheKey:     "city_count",
			QuerySQL:     query,
			CacheTTL:     cacheTTL,
			UpdatePeriod: updatePeriod,
		}

		statsController := task.NewDeliveryStatsController(pool, statsConfig)
		cacheController := task.NewCacheController(pool, rdb, cacheConfig)

		go task.RunScenario(scenarioCtx, pool, scenario)
		go statsController.Start(scenarioCtx, aggregator)
		go cacheController.Start(scenarioCtx, aggregator)

		<-scenarioCtx.Done()

		filename := string(scenario) + "_stats.csv"
		if err := aggregator.ExportCSV(filename); err != nil {
			log.Printf("Ошибка экспорта CSV: %v", err)
		} else {
			log.Printf("Результаты сценария %s сохранены в файл %s", scenario, filename)
		}
	}
}


func connectDB(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.DatabaseName)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("ошибка при проверке соединения: %v", err)
	}

	return pool, nil
}

func connectRedis(cfg *config.Config, ctx context.Context) (*redis.Client, error) {
	connString := fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port)
	rdb := redis.NewClient(&redis.Options{Addr: connString})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке соединения: %v", err)
	}

	return rdb, nil
}
