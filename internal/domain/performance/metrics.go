package performance

import "time"

type HandlerMetric struct {
	UserID      int64
	HandlerName string
	Command     string
	DurationMs  int64
	Timestamp   time.Time
}

type HandlerStats struct {
	TotalRequests   int64
	AverageDuration float64
	MinDuration     int64
	MaxDuration     int64
	TotalDuration   int64
}

type PercentileStats struct {
	P50 float64
	P75 float64
	P90 float64
	P95 float64
	P99 float64
}

type PerformanceReport struct {
	Date              time.Time
	Period            string // "daily" or "weekly"
	HandlerStats      map[string]HandlerStats
	SlowestRequests   []HandlerMetric
	Percentiles       map[string]PercentileStats
	Trends            map[string]TrendData
	TotalRequests     int64
	AverageSystemTime float64
}

type TrendData struct {
	CurrentAvg  float64
	PreviousAvg float64
	Change      float64 // percentage change
	Direction   string  // "up", "down", "stable"
}
