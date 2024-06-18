package handler

import (
	"fmt"
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/mosteligible/go-brrrr/pkg/types"
	"github.com/mosteligible/go-brrrr/pkg/utils"
)

func MetricsBuilder(mq *types.MetricsQueue) {
	for {
		dataPoint := <-mq.MetricData
		if dataPoint.ResponseTime < 0 {
			continue
		}
		if dataPoint.Success {
			if _, ok := mq.SuccessTimeMap[dataPoint.TimeStamp]; ok {
				mq.SuccessTimeMap[dataPoint.TimeStamp]++
			} else {
				mq.SuccessTimeMap[dataPoint.TimeStamp] = 1
			}
			if _, ok := mq.FailureTimeMap[dataPoint.TimeStamp]; !ok {
				mq.FailureTimeMap[dataPoint.TimeStamp] = 0
			}
		} else {
			if _, ok := mq.FailureTimeMap[dataPoint.TimeStamp]; ok {
				mq.FailureTimeMap[dataPoint.TimeStamp]++
			} else {
				mq.FailureTimeMap[dataPoint.TimeStamp] = 1
			}
			if _, ok := mq.SuccessTimeMap[dataPoint.TimeStamp]; !ok {
				mq.SuccessTimeMap[dataPoint.TimeStamp] = 0
			}
		}
		if _, ok := mq.ResponseTimeMap[dataPoint.TimeStamp]; !ok {
			mq.ResponseTimeMap[dataPoint.TimeStamp] = []int64{}
		}
		mq.ResponseTimeMap[dataPoint.TimeStamp] = append(mq.ResponseTimeMap[dataPoint.TimeStamp], dataPoint.ResponseTime)
	}
}

func ShowResults(mq *types.MetricsQueue) {
	timeArr := utils.GetKeysOfMap[int64](mq.FailureTimeMap)
	sort.Slice(timeArr, func(i, j int) bool {
		first, _ := time.Parse(time.DateTime, timeArr[i])
		second, _ := time.Parse(time.DateTime, timeArr[j])
		return first.Before(second)
	})
	respTimes := map[string][]int64{}
	for k, v := range mq.ResponseTimeMap {
		vals := []int64{}
		for _, rTime := range v {
			if rTime < 0 {
				continue
			}
			vals = append(vals, rTime)
		}
		respTimes[k] = vals
	}
	logContent := ""
	for _, k := range timeArr {
		meanResponseTimeRawData := stats.LoadRawData(respTimes[k])
		medianResponseTime, _ := stats.Median(meanResponseTimeRawData)
		meanRespTime, _ := stats.Mean(meanResponseTimeRawData)
		minRespTime, _ := stats.Min(meanResponseTimeRawData)
		maxRespTime, _ := stats.Max(meanResponseTimeRawData)
		row := fmt.Sprintf(
			" [x] %s: success: %d, failures: %d, median: %f, avg resp: %f, min: %f, max: %f\n",
			k, mq.SuccessTimeMap[k],
			mq.FailureTimeMap[k],
			medianResponseTime,
			meanRespTime,
			minRespTime,
			maxRespTime,
		)
		fmt.Println(row)
		logContent = fmt.Sprintf("%s%s", logContent, row)
	}

	utils.WriteToJsonFile[map[string][]int64](
		mq.ResponseTimeMap, "response_time_maps.json", true,
	)
	utils.WriteToJsonFile[map[string]int64](mq.FailureTimeMap, "failures.json", true)
	utils.WriteToJsonFile[map[string]int64](mq.SuccessTimeMap, "successes.json", true)
}
