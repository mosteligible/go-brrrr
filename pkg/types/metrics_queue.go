package types

import (
	"fmt"
)

type DataPoint struct {
	TimeStamp    string // String format timestamp as YYYY-MM-DD HH:MM:SS
	ResponseTime int64
	Success      bool
}

func (dp DataPoint) String() string {
	return fmt.Sprintf(
		"timestamp: %s, resptime: %d, success: %v\n",
		dp.TimeStamp, dp.ResponseTime, dp.Success,
	)
}

type MetricsQueue struct {
	limit           int
	FailureTimeMap  map[string]int64
	SuccessTimeMap  map[string]int64
	ResponseTimeMap map[string][]int64
	MetricData      chan DataPoint
}

func NewMetricsQueue(limit int) *MetricsQueue {
	mq := MetricsQueue{
		limit:           limit,
		FailureTimeMap:  map[string]int64{},
		SuccessTimeMap:  map[string]int64{},
		ResponseTimeMap: map[string][]int64{},
		MetricData:      make(chan DataPoint, limit),
	}

	return &mq
}

func (mq *MetricsQueue) AddMetricData(
	timeStr string, timeTaken int64, success bool,
) {
	dataPoint := DataPoint{
		TimeStamp:    timeStr,
		ResponseTime: timeTaken,
		Success:      success,
	}
	mq.MetricData <- dataPoint
}
