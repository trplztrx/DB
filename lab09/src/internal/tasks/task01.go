package task

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryStatsConfig struct {
	Query        string
	UpdatePeriod time.Duration
}

type DeliveryStatsController struct {
	config    *DeliveryStatsConfig
	pool      *pgxpool.Pool
	startTime time.Time
}

func NewDeliveryStatsController(pool *pgxpool.Pool, config *DeliveryStatsConfig) *DeliveryStatsController {
	return &DeliveryStatsController{
		config:    config,
		pool:      pool,
		startTime: time.Now(),
	}
}

func (dsc *DeliveryStatsController) Start(ctx context.Context, aggregator *StatsAggregator) {
	ticker := time.NewTicker(dsc.config.UpdatePeriod)
	defer ticker.Stop()
	dsc.executeQuery(ctx, aggregator)
	for {
		select {
		case <-ctx.Done():
			log.Println("Контекст завершен. Остановка выполнения запросов")
			return
		case <-ticker.C:
			dsc.executeQuery(ctx, aggregator)
		}
	}
}

func (dsc *DeliveryStatsController) executeQuery(ctx context.Context, aggregator *StatsAggregator) {
	start := time.Now()

	rows, err := dsc.pool.Query(ctx, dsc.config.Query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return
	}
	defer rows.Close()

	elapsed := time.Since(start)
	secondsSinceStart := int(time.Since(dsc.startTime).Seconds())
	aggregator.RecordDBStat(secondsSinceStart, elapsed)
	log.Printf("Запрос к БД выполнен за %v", elapsed)
}
