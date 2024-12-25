package task

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"
	"time"
)

type CombinedStats struct {
	Second         int
	DBQueryTime    time.Duration
	CacheQueryTime time.Duration
}

type StatsAggregator struct {
	stats map[int]*CombinedStats
	mu    sync.Mutex
}

func NewStatsAggregator() *StatsAggregator {
	return &StatsAggregator{
		stats: make(map[int]*CombinedStats),
	}
}

func (agg *StatsAggregator) RecordDBStat(second int, duration time.Duration) {
	agg.mu.Lock()
	defer agg.mu.Unlock()

	stat, exists := agg.stats[second]
	if !exists {
		stat = &CombinedStats{Second: second}
		agg.stats[second] = stat
	}
	stat.DBQueryTime = duration
}

func (agg *StatsAggregator) RecordCacheStat(second int, duration time.Duration) {
	agg.mu.Lock()
	defer agg.mu.Unlock()

	stat, exists := agg.stats[second]
	if !exists {
		stat = &CombinedStats{Second: second}
		agg.stats[second] = stat
	}
	stat.CacheQueryTime = duration
}

func (agg *StatsAggregator) ExportCSV(filename string) error {
    agg.mu.Lock()
    defer agg.mu.Unlock()

    filepath := "/app/data/" + filename

    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    writer.Write([]string{"Second", "DBQueryTime", "CacheQueryTime"})

    for _, stat := range agg.stats {
        writer.Write([]string{
            strconv.Itoa(stat.Second),
            stat.DBQueryTime.String(),
            stat.CacheQueryTime.String(),
        })
    }
    return nil
}

