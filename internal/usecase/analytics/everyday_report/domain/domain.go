package domain

import "time"

type Metric struct {
	Name   string
	Value  float64
	Suffix *string
}

type EverydayReportResponse struct {
	Date                 time.Time
	Metrics              []Metric
	UserIDsToSendMetrics []int64
}
